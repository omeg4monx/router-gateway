// Package telegram implements the vaulted Telegram Bot controller, fail-closed
// webhook receiver, dashboard command simulator, and eventbus alert forwarder.
package telegram

import (
	"context"
	"crypto/rand"
	"crypto/subtle"
	"embed"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/omeg4monx/router-gateway/internal/contracts"
	"github.com/omeg4monx/router-gateway/internal/platform/crypto"
	"github.com/omeg4monx/router-gateway/internal/platform/db"
	"github.com/omeg4monx/router-gateway/internal/platform/eventbus"
	"github.com/omeg4monx/router-gateway/internal/platform/httpserver"
	"github.com/omeg4monx/router-gateway/internal/platform/lifecycle"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

const (
	moduleName       = "telegram"
	maxJSONBodyBytes = 1 << 20
	busSubBufSize    = 32

	topicNotificationsEmit = "notifications.emit"
	topicJobsFailed        = "jobs.failed"
	topicAuthPolicyChanged = "auth.policy_changed"
)

var watchedTopics = []string{
	topicJobsFailed,
	topicAuthPolicyChanged,
	topicNotificationsEmit,
}

// Options configures a telegram Module instance.
type Options struct {
	DB            *db.DB
	Vault         *crypto.Vault
	HTTPClient    *http.Client
	Bus           *eventbus.Bus
	Logger        *slog.Logger
	Providers     contracts.ProviderStore
	Models        contracts.ModelLister
	Jobs          contracts.JobInspector
	ActiveModules func() []string
	// APIBaseURL optionally overrides https://api.telegram.org for tests.
	APIBaseURL string
}

// ConfigView is the sanitized JSON representation of the Telegram Bot config.
type ConfigView struct {
	Enabled                 bool   `json:"enabled"`
	BotTokenConfigured      bool   `json:"bot_token_configured"`
	BotTokenMask            string `json:"bot_token_mask"`
	ChatID                  string `json:"chat_id"`
	UserID                  string `json:"user_id"`
	WebhookURL              string `json:"webhook_url"`
	WebhookPath             string `json:"webhook_path"`
	WebhookSecretConfigured bool   `json:"webhook_secret_configured"`
	NotifyOnJobFailure      bool   `json:"notify_on_job_failure"`
	NotifyOnAuthChange      bool   `json:"notify_on_auth_change"`
	NotifyOnEmit            bool   `json:"notify_on_emit"`
	UpdatedAt               string `json:"updated_at"`
}

// Module implements lifecycle.Module for Telegram Bot notifications and commands.
type Module struct {
	repo          *repo
	vault         *crypto.Vault
	client        *botAPIClient
	bus           *eventbus.Bus
	logger        *slog.Logger
	providers     contracts.ProviderStore
	models        contracts.ModelLister
	jobs          contracts.JobInspector
	activeModules func() []string
	customAPIBase bool
	startedAt     time.Time

	cfgMu sync.Mutex

	lifeMu  sync.Mutex
	started bool
	cancel  context.CancelFunc
	wg      sync.WaitGroup
}

var _ lifecycle.Module = (*Module)(nil)

// New constructs a telegram Module.
func New(opts Options) *Module {
	logger := opts.Logger
	if logger == nil {
		logger = slog.Default()
	}
	return &Module{
		repo:          newRepo(opts.DB),
		vault:         opts.Vault,
		client:        newBotAPIClient(opts.APIBaseURL, opts.HTTPClient),
		bus:           opts.Bus,
		logger:        logger,
		providers:     opts.Providers,
		models:        opts.Models,
		jobs:          opts.Jobs,
		activeModules: opts.ActiveModules,
		customAPIBase: strings.TrimSpace(opts.APIBaseURL) != "",
		startedAt:     time.Now(),
	}
}

// Name returns the canonical module name.
func (m *Module) Name() string { return moduleName }

// Dependencies declares modules required before telegram initializes.
func (m *Module) Dependencies() []string { return []string{"system"} }

// Migrations returns the embedded SQLite migrations for telegram.
func (m *Module) Migrations() fs.FS {
	sub, err := fs.Sub(migrationsFS, "migrations")
	if err != nil {
		panic(fmt.Sprintf("telegram: load embedded migrations: %v", err))
	}
	return sub
}

// RegisterRoutes mounts the fail-closed public webhook and admin control routes.
func (m *Module) RegisterRoutes(mounts lifecycle.RouteMounts) {
	if mounts.PublicAPI != nil {
		mounts.PublicAPI.Post("/telegram/webhook", m.handleWebhook)
	}
	if mounts.AdminAPI != nil {
		mounts.AdminAPI.Get("/telegram/config", m.handleGetConfig)
		mounts.AdminAPI.Patch("/telegram/config", m.handleUpdateConfig)
		mounts.AdminAPI.Post("/telegram/secret/rotate", m.handleRotateSecret)
		mounts.AdminAPI.Post("/telegram/webhook/register", m.handleRegisterWebhook)
		mounts.AdminAPI.Post("/telegram/webhook/delete", m.handleDeleteWebhook)
		mounts.AdminAPI.Post("/telegram/test", m.handleSendTest)
		mounts.AdminAPI.Post("/telegram/command", m.handleSimulateCommand)
		mounts.AdminAPI.Get("/telegram/events", m.handleListEvents)
		mounts.AdminAPI.Delete("/telegram/events", m.handleClearEvents)
	}
}

// Start launches the eventbus alert subscriber loop.
func (m *Module) Start(context.Context) error {
	m.lifeMu.Lock()
	defer m.lifeMu.Unlock()
	if m.started {
		return nil
	}
	m.startedAt = time.Now()
	m.started = true
	if m.bus == nil {
		return nil
	}

	subCtx, cancel := context.WithCancel(context.Background())
	m.cancel = cancel

	for _, topic := range watchedTopics {
		ch, unsub := m.bus.Subscribe(topic, busSubBufSize)
		m.wg.Add(1)
		go func(t string, events <-chan eventbus.Event, unsubscribe func()) {
			defer m.wg.Done()
			defer unsubscribe()
			for {
				select {
				case <-subCtx.Done():
					return
				case ev, ok := <-events:
					if !ok {
						return
					}
					m.handleBusEvent(subCtx, t, ev)
				}
			}
		}(topic, ch, unsub)
	}
	return nil
}

// Stop cancels eventbus subscriptions and waits for background workers to exit.
func (m *Module) Stop(ctx context.Context) error {
	m.lifeMu.Lock()
	if !m.started {
		m.lifeMu.Unlock()
		return nil
	}
	m.started = false
	cancel := m.cancel
	m.cancel = nil
	m.lifeMu.Unlock()

	if cancel != nil {
		cancel()
	}

	done := make(chan struct{})
	go func() {
		m.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (m *Module) handleBusEvent(ctx context.Context, topic string, ev eventbus.Event) {
	m.cfgMu.Lock()
	cfg, err := m.repo.getConfig(ctx)
	m.cfgMu.Unlock()
	if err != nil || !cfg.Enabled || cfg.BotTokenEnc == "" || strings.TrimSpace(cfg.ChatID) == "" {
		return
	}

	switch topic {
	case topicJobsFailed:
		if !cfg.NotifyOnJobFailure {
			return
		}
	case topicAuthPolicyChanged:
		if !cfg.NotifyOnAuthChange {
			return
		}
	case topicNotificationsEmit:
		if !cfg.NotifyOnEmit {
			return
		}
	default:
		return
	}

	botToken, err := m.decryptSecret(cfg.BotTokenEnc)
	if err != nil {
		return
	}

	summary := formatBusAlert(topic, ev.Payload)
	sendErr := m.client.sendMessage(ctx, botToken, cfg.ChatID, summary)
	status := "delivered"
	errText := ""
	if sendErr != nil {
		status = "failed"
		errText = sendErr.Error()
	}
	_ = m.repo.recordEvent(ctx, Event{
		Direction: "outbound",
		Kind:      "alert",
		Command:   topic,
		Status:    status,
		Summary:   truncateText(summary, 240),
		Error:     errText,
	})
}

func formatBusAlert(topic string, payload any) string {
	if p, ok := payload.(map[string]any); ok {
		title, _ := p["title"].(string)
		msg, _ := p["message"].(string)
		if title != "" && msg != "" {
			return fmt.Sprintf("[%s] %s: %s", topic, title, msg)
		}
		if title != "" {
			return fmt.Sprintf("[%s] %s", topic, title)
		}
	}
	raw, err := json.Marshal(payload)
	if err != nil || len(raw) == 0 || string(raw) == "null" {
		return fmt.Sprintf("[%s] event triggered", topic)
	}
	return fmt.Sprintf("[%s] %s", topic, truncateText(string(raw), 320))
}

func (m *Module) toConfigView(cfg storedConfig) ConfigView {
	return ConfigView{
		Enabled:                 cfg.Enabled,
		BotTokenConfigured:      cfg.BotTokenEnc != "",
		BotTokenMask:            cfg.BotTokenMask,
		ChatID:                  cfg.ChatID,
		UserID:                  cfg.UserID,
		WebhookURL:              cfg.WebhookURL,
		WebhookPath:             "/api/telegram/webhook",
		WebhookSecretConfigured: cfg.WebhookSecretEnc != "",
		NotifyOnJobFailure:      cfg.NotifyOnJobFailure,
		NotifyOnAuthChange:      cfg.NotifyOnAuthChange,
		NotifyOnEmit:            cfg.NotifyOnEmit,
		UpdatedAt:               cfg.UpdatedAt,
	}
}

func (m *Module) encryptSecret(plaintext string) (string, error) {
	if m.vault == nil {
		return "", errors.New("secret vault is not configured")
	}
	return m.vault.Encrypt([]byte(plaintext))
}

func (m *Module) decryptSecret(ciphertext string) (string, error) {
	if m.vault == nil {
		return "", errors.New("secret vault is not configured")
	}
	raw, err := m.vault.Decrypt(ciphertext)
	if err != nil {
		return "", err
	}
	return string(raw), nil
}

func generateWebhookSecret() (string, error) {
	var buf [32]byte
	if _, err := rand.Read(buf[:]); err != nil {
		return "", fmt.Errorf("generate webhook secret: %w", err)
	}
	return hex.EncodeToString(buf[:]), nil
}

func maskToken(token string) string {
	trimmed := strings.TrimSpace(token)
	if len(trimmed) <= 8 {
		return "****"
	}
	return trimmed[:4] + "…" + trimmed[len(trimmed)-4:]
}

func (m *Module) validateWebhookURL(raw string) error {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return nil
	}
	u, err := url.Parse(trimmed)
	if err != nil || u.Host == "" {
		return errors.New("webhook_url must be a valid URL")
	}
	if u.Scheme == "https" {
		return nil
	}
	if m.customAPIBase && u.Scheme == "http" {
		return nil
	}
	return errors.New("webhook_url must use https://")
}

func (m *Module) handleGetConfig(w http.ResponseWriter, r *http.Request) {
	m.cfgMu.Lock()
	cfg, err := m.repo.getConfig(r.Context())
	m.cfgMu.Unlock()
	if err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "telegram_config_read_failed", "failed to load telegram config")
		return
	}
	httpserver.WriteJSON(w, http.StatusOK, m.toConfigView(cfg))
}

func (m *Module) handleUpdateConfig(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Enabled            *bool   `json:"enabled"`
		BotToken           *string `json:"bot_token"`
		ClearBotToken      *bool   `json:"clear_bot_token"`
		ChatID             *string `json:"chat_id"`
		UserID             *string `json:"user_id"`
		WebhookURL         *string `json:"webhook_url"`
		NotifyOnJobFailure *bool   `json:"notify_on_job_failure"`
		NotifyOnAuthChange *bool   `json:"notify_on_auth_change"`
		NotifyOnEmit       *bool   `json:"notify_on_emit"`
	}
	if err := decodeJSON(r, &body); err != nil {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_json", err.Error())
		return
	}

	if body.WebhookURL != nil {
		if err := m.validateWebhookURL(*body.WebhookURL); err != nil {
			httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_webhook_url", err.Error())
			return
		}
	}

	m.cfgMu.Lock()
	defer m.cfgMu.Unlock()

	cfg, err := m.repo.getConfig(r.Context())
	if err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "telegram_config_read_failed", "failed to load telegram config")
		return
	}

	if body.Enabled != nil {
		cfg.Enabled = *body.Enabled
	}
	if body.ClearBotToken != nil && *body.ClearBotToken {
		cfg.BotTokenEnc = ""
		cfg.BotTokenMask = ""
	}
	if body.BotToken != nil {
		tok := strings.TrimSpace(*body.BotToken)
		if tok != "" {
			enc, encErr := m.encryptSecret(tok)
			if encErr != nil {
				httpserver.WriteError(w, r, http.StatusInternalServerError, "vault_encrypt_failed", "failed to encrypt bot token")
				return
			}
			cfg.BotTokenEnc = enc
			cfg.BotTokenMask = maskToken(tok)
		}
	}
	if body.ChatID != nil {
		cfg.ChatID = strings.TrimSpace(*body.ChatID)
	}
	if body.UserID != nil {
		cfg.UserID = strings.TrimSpace(*body.UserID)
	}
	if body.WebhookURL != nil {
		cfg.WebhookURL = strings.TrimSpace(*body.WebhookURL)
	}
	if body.NotifyOnJobFailure != nil {
		cfg.NotifyOnJobFailure = *body.NotifyOnJobFailure
	}
	if body.NotifyOnAuthChange != nil {
		cfg.NotifyOnAuthChange = *body.NotifyOnAuthChange
	}
	if body.NotifyOnEmit != nil {
		cfg.NotifyOnEmit = *body.NotifyOnEmit
	}

	// Ensure a strong webhook secret exists once a bot token or enablement is set.
	if cfg.WebhookSecretEnc == "" && (cfg.BotTokenEnc != "" || cfg.Enabled) {
		sec, secErr := generateWebhookSecret()
		if secErr != nil {
			httpserver.WriteError(w, r, http.StatusInternalServerError, "secret_generate_failed", "failed to generate webhook secret")
			return
		}
		encSec, encErr := m.encryptSecret(sec)
		if encErr != nil {
			httpserver.WriteError(w, r, http.StatusInternalServerError, "vault_encrypt_failed", "failed to encrypt webhook secret")
			return
		}
		cfg.WebhookSecretEnc = encSec
	}

	if err := m.repo.saveConfig(r.Context(), cfg); err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "telegram_config_save_failed", "failed to save telegram config")
		return
	}
	updated, err := m.repo.getConfig(r.Context())
	if err != nil {
		updated = cfg
	}
	httpserver.WriteJSON(w, http.StatusOK, m.toConfigView(updated))
}

func (m *Module) handleRotateSecret(w http.ResponseWriter, r *http.Request) {
	m.cfgMu.Lock()
	defer m.cfgMu.Unlock()

	cfg, err := m.repo.getConfig(r.Context())
	if err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "telegram_config_read_failed", "failed to load telegram config")
		return
	}
	sec, err := generateWebhookSecret()
	if err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "secret_generate_failed", "failed to generate webhook secret")
		return
	}
	enc, err := m.encryptSecret(sec)
	if err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "vault_encrypt_failed", "failed to encrypt webhook secret")
		return
	}
	cfg.WebhookSecretEnc = enc
	if err := m.repo.saveConfig(r.Context(), cfg); err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "telegram_config_save_failed", "failed to save rotated webhook secret")
		return
	}
	_ = m.repo.recordEvent(r.Context(), Event{
		Direction: "system",
		Kind:      "secret_rotated",
		Status:    "ok",
		Summary:   "Webhook secret token rotated; re-register webhook to sync Telegram.",
	})
	updated, _ := m.repo.getConfig(r.Context())
	httpserver.WriteJSON(w, http.StatusOK, m.toConfigView(updated))
}

func (m *Module) handleRegisterWebhook(w http.ResponseWriter, r *http.Request) {
	m.cfgMu.Lock()
	cfg, err := m.repo.getConfig(r.Context())
	m.cfgMu.Unlock()
	if err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "telegram_config_read_failed", "failed to load telegram config")
		return
	}
	if cfg.BotTokenEnc == "" {
		httpserver.WriteError(w, r, http.StatusBadRequest, "bot_token_missing", "configure bot_token before registering webhook")
		return
	}
	if strings.TrimSpace(cfg.WebhookURL) == "" {
		httpserver.WriteError(w, r, http.StatusBadRequest, "webhook_url_missing", "configure webhook_url before registering webhook")
		return
	}
	if err := m.validateWebhookURL(cfg.WebhookURL); err != nil {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_webhook_url", err.Error())
		return
	}

	botToken, err := m.decryptSecret(cfg.BotTokenEnc)
	if err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "vault_decrypt_failed", "failed to decrypt bot token")
		return
	}
	webhookSecret, err := m.decryptSecret(cfg.WebhookSecretEnc)
	if err != nil || webhookSecret == "" {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "vault_decrypt_failed", "webhook secret is not initialized")
		return
	}

	if err := m.client.setWebhook(r.Context(), botToken, cfg.WebhookURL, webhookSecret); err != nil {
		_ = m.repo.recordEvent(r.Context(), Event{
			Direction: "outbound",
			Kind:      "webhook_register",
			Status:    "failed",
			Summary:   cfg.WebhookURL,
			Error:     err.Error(),
		})
		httpserver.WriteError(w, r, http.StatusBadGateway, "telegram_api_failed", err.Error())
		return
	}

	_ = m.repo.recordEvent(r.Context(), Event{
		Direction: "outbound",
		Kind:      "webhook_register",
		Status:    "ok",
		Summary:   "Webhook registered at " + cfg.WebhookURL,
	})
	httpserver.WriteJSON(w, http.StatusOK, map[string]any{
		"ok":          true,
		"webhook_url": cfg.WebhookURL,
	})
}

func (m *Module) handleDeleteWebhook(w http.ResponseWriter, r *http.Request) {
	m.cfgMu.Lock()
	cfg, err := m.repo.getConfig(r.Context())
	m.cfgMu.Unlock()
	if err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "telegram_config_read_failed", "failed to load telegram config")
		return
	}
	if cfg.BotTokenEnc == "" {
		httpserver.WriteError(w, r, http.StatusBadRequest, "bot_token_missing", "configure bot_token before deleting webhook")
		return
	}
	botToken, err := m.decryptSecret(cfg.BotTokenEnc)
	if err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "vault_decrypt_failed", "failed to decrypt bot token")
		return
	}
	if err := m.client.deleteWebhook(r.Context(), botToken); err != nil {
		_ = m.repo.recordEvent(r.Context(), Event{
			Direction: "outbound",
			Kind:      "webhook_delete",
			Status:    "failed",
			Error:     err.Error(),
		})
		httpserver.WriteError(w, r, http.StatusBadGateway, "telegram_api_failed", err.Error())
		return
	}
	_ = m.repo.recordEvent(r.Context(), Event{
		Direction: "outbound",
		Kind:      "webhook_delete",
		Status:    "ok",
		Summary:   "Webhook removed from Telegram Bot API.",
	})
	httpserver.WriteJSON(w, http.StatusOK, map[string]any{
		"ok": true,
	})
}

func (m *Module) handleSendTest(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Message string `json:"message"`
	}
	if r.Body != nil {
		_ = decodeJSON(r, &body)
	}
	msg := strings.TrimSpace(body.Message)
	if msg == "" {
		msg = "OMEG4MON ROUTER Telegram connectivity probe succeeded."
	}

	m.cfgMu.Lock()
	cfg, err := m.repo.getConfig(r.Context())
	m.cfgMu.Unlock()
	if err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "telegram_config_read_failed", "failed to load telegram config")
		return
	}
	if cfg.BotTokenEnc == "" || strings.TrimSpace(cfg.ChatID) == "" {
		httpserver.WriteError(w, r, http.StatusBadRequest, "telegram_unconfigured", "configure bot_token and chat_id before sending a test message")
		return
	}
	botToken, err := m.decryptSecret(cfg.BotTokenEnc)
	if err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "vault_decrypt_failed", "failed to decrypt bot token")
		return
	}
	if err := m.client.sendMessage(r.Context(), botToken, cfg.ChatID, msg); err != nil {
		_ = m.repo.recordEvent(r.Context(), Event{
			Direction: "outbound",
			Kind:      "test",
			Status:    "failed",
			Summary:   truncateText(msg, 200),
			Error:     err.Error(),
		})
		httpserver.WriteError(w, r, http.StatusBadGateway, "telegram_api_failed", err.Error())
		return
	}
	_ = m.repo.recordEvent(r.Context(), Event{
		Direction: "outbound",
		Kind:      "test",
		Status:    "delivered",
		Summary:   truncateText(msg, 200),
	})
	httpserver.WriteJSON(w, http.StatusOK, map[string]any{
		"ok":      true,
		"chat_id": cfg.ChatID,
		"message": msg,
	})
}

func (m *Module) handleSimulateCommand(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Command string `json:"command"`
	}
	if err := decodeJSON(r, &body); err != nil {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_json", err.Error())
		return
	}
	cmd, reply := m.executeCommand(r.Context(), body.Command)
	executedAt := nowTimestamp()
	_ = m.repo.recordEvent(r.Context(), Event{
		Direction: "internal",
		Kind:      "simulator",
		Command:   "/" + cmd,
		Status:    "ok",
		Summary:   truncateText(reply, 200),
	})
	httpserver.WriteJSON(w, http.StatusOK, map[string]any{
		"command":     "/" + cmd,
		"reply":       reply,
		"executed_at": executedAt,
	})
}

func (m *Module) handleWebhook(w http.ResponseWriter, r *http.Request) {
	m.cfgMu.Lock()
	cfg, err := m.repo.getConfig(r.Context())
	m.cfgMu.Unlock()
	if err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "telegram_config_read_failed", "failed to load telegram config")
		return
	}
	if !cfg.Enabled || cfg.BotTokenEnc == "" || cfg.WebhookSecretEnc == "" {
		httpserver.WriteError(w, r, http.StatusForbidden, "telegram_disabled", "telegram bot webhook is disabled or not configured")
		return
	}

	expectedSecret, err := m.decryptSecret(cfg.WebhookSecretEnc)
	if err != nil || expectedSecret == "" {
		httpserver.WriteError(w, r, http.StatusForbidden, "telegram_unconfigured", "webhook secret is unavailable")
		return
	}
	providedSecret := r.Header.Get("X-Telegram-Bot-Api-Secret-Token")
	if !safeEqual(providedSecret, expectedSecret) {
		_ = m.repo.recordEvent(r.Context(), Event{
			Direction: "inbound",
			Kind:      "webhook",
			Status:    "rejected",
			Error:     "invalid webhook secret token",
		})
		httpserver.WriteError(w, r, http.StatusUnauthorized, "invalid_webhook_secret", "unauthorized telegram webhook request")
		return
	}

	// Fail closed if no admin user_id is configured.
	configuredUserID := strings.TrimSpace(cfg.UserID)
	if configuredUserID == "" {
		_ = m.repo.recordEvent(r.Context(), Event{
			Direction: "inbound",
			Kind:      "webhook",
			Status:    "rejected",
			Error:     "admin user_id is not configured",
		})
		httpserver.WriteError(w, r, http.StatusForbidden, "user_id_unconfigured", "telegram admin user_id must be configured")
		return
	}

	var update struct {
		UpdateID int64 `json:"update_id"`
		Message  *struct {
			MessageID int64 `json:"message_id"`
			From      *struct {
				ID int64 `json:"id"`
			} `json:"from"`
			Chat *struct {
				ID int64 `json:"id"`
			} `json:"chat"`
			Text string `json:"text"`
		} `json:"message"`
	}
	if err := decodeJSON(r, &update); err != nil {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_json", err.Error())
		return
	}
	if update.Message == nil || update.Message.From == nil || update.Message.Chat == nil {
		httpserver.WriteJSON(w, http.StatusOK, map[string]any{"ok": true, "ignored": true})
		return
	}

	senderID := strconv.FormatInt(update.Message.From.ID, 10)
	if !safeEqual(senderID, configuredUserID) {
		_ = m.repo.recordEvent(r.Context(), Event{
			Direction: "inbound",
			Kind:      "webhook",
			Status:    "rejected",
			Summary:   "unauthorized sender " + senderID,
			Error:     "sender user_id does not match configured admin",
		})
		httpserver.WriteError(w, r, http.StatusForbidden, "unauthorized_user", "sender is not authorized")
		return
	}

	cmd, reply := m.executeCommand(r.Context(), update.Message.Text)
	botToken, err := m.decryptSecret(cfg.BotTokenEnc)
	if err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "vault_decrypt_failed", "failed to decrypt bot token")
		return
	}

	chatID := strconv.FormatInt(update.Message.Chat.ID, 10)
	sendErr := m.client.sendMessage(r.Context(), botToken, chatID, reply)
	status := "ok"
	errText := ""
	if sendErr != nil {
		status = "failed"
		errText = sendErr.Error()
		m.logger.Warn("telegram webhook reply delivery failed", "error", sendErr)
	}
	_ = m.repo.recordEvent(r.Context(), Event{
		Direction: "inbound",
		Kind:      "webhook",
		Command:   "/" + cmd,
		Status:    status,
		Summary:   truncateText(reply, 200),
		Error:     errText,
	})
	httpserver.WriteJSON(w, http.StatusOK, map[string]any{
		"ok":      sendErr == nil,
		"command": "/" + cmd,
	})
}

func (m *Module) handleListEvents(w http.ResponseWriter, r *http.Request) {
	limit := 50
	if raw := strings.TrimSpace(r.URL.Query().Get("limit")); raw != "" {
		if parsed, err := strconv.Atoi(raw); err == nil && parsed > 0 {
			limit = parsed
		}
	}
	events, err := m.repo.listEvents(r.Context(), limit)
	if err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "telegram_events_read_failed", "failed to load telegram events")
		return
	}
	httpserver.WriteJSON(w, http.StatusOK, map[string]any{
		"events": events,
	})
}

func (m *Module) handleClearEvents(w http.ResponseWriter, r *http.Request) {
	cleared, err := m.repo.clearEvents(r.Context())
	if err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "telegram_events_clear_failed", "failed to clear telegram events")
		return
	}
	httpserver.WriteJSON(w, http.StatusOK, map[string]any{
		"cleared": cleared,
	})
}

func safeEqual(a, b string) bool {
	if a == "" || len(a) != len(b) {
		return false
	}
	return subtle.ConstantTimeCompare([]byte(a), []byte(b)) == 1
}

func decodeJSON(r *http.Request, dst any) error {
	raw, err := io.ReadAll(io.LimitReader(r.Body, maxJSONBodyBytes+1))
	if err != nil {
		return errors.New("failed to read request body")
	}
	if len(raw) > maxJSONBodyBytes {
		return errors.New("request body exceeds limit")
	}
	if strings.TrimSpace(string(raw)) == "" {
		return errors.New("request body is empty")
	}
	if err := json.Unmarshal(raw, dst); err != nil {
		return fmt.Errorf("invalid JSON body: %w", err)
	}
	return nil
}
