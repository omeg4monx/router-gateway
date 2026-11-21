package connection

import (
	"bytes"
	"crypto/subtle"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strings"
	"time"

	"github.com/omeg4monx/router-gateway/internal/database"
)

var startTime = time.Now()

func formatUptime(d time.Duration) string {
	days := int(d.Hours()) / 24
	hours := int(d.Hours()) % 24
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60

	parts := []string{}
	if days > 0 {
		parts = append(parts, fmt.Sprintf("%dd", days))
	}
	if hours > 0 {
		parts = append(parts, fmt.Sprintf("%dh", hours))
	}
	if minutes > 0 {
		parts = append(parts, fmt.Sprintf("%dm", minutes))
	}
	parts = append(parts, fmt.Sprintf("%ds", seconds))
	return strings.Join(parts, " ")
}

type TelegramUpdate struct {
	UpdateID int64 `json:"update_id"`
	Message  *struct {
		MessageID int64  `json:"message_id"`
		Text      string `json:"text"`
		Chat      struct {
			ID int64 `json:"id"`
		} `json:"chat"`
		From struct {
			ID int64 `json:"id"`
		} `json:"from"`
	} `json:"message"`
}

func HandleCommand(rawCmd, requestUserId, requestToken string, db *database.DB) (string, error) {
	cfg, err := db.GetTelegramConfig()
	if err != nil {
		return "", err
	}

	// The bot token is the only secret here; a Telegram user id is public, so it
	// can never authenticate on its own.
	if cfg.BotToken == "" || subtle.ConstantTimeCompare([]byte(requestToken), []byte(cfg.BotToken)) != 1 {
		return "", fmt.Errorf("unauthorized: invalid token")
	}

	parts := strings.Fields(strings.TrimSpace(rawCmd))
	cmd := "status"
	if len(parts) > 0 {
		cmd = strings.TrimPrefix(strings.ToLower(parts[0]), "/")
		cmd = strings.TrimPrefix(cmd, "!")
	}

	switch cmd {
	case "start", "help":
		return "⚡ *OMEG4MON ROUTER Go Core Controller*\n" +
			"High-Performance Go Engine · Control Panel\n\n" +
			"*Available Commands:*\n" +
			"• `/status` - Live system health, uptime & metrics\n" +
			"• `/providers` - List active providers & connections\n" +
			"• `/pools` - Check proxy pool health\n" +
			"• `/prune` - Delete dead proxies from pools\n" +
			"• `/ping` - Heartbeat probe\n" +
			"• `/help` - Show this command list", nil

	case "ping":
		return fmt.Sprintf("🏓 Pong! OMEG4MON ROUTER Go Core is active (Uptime: %s)", formatUptime(time.Since(startTime))), nil

	case "status":
		activeConns, totalConns, _ := db.GetActiveConnectionCount()
		livePools, totalPools, _ := db.GetProxyPoolStats()

		var mem runtime.MemStats
		runtime.ReadMemStats(&mem)
		allocMb := float64(mem.Alloc) / (1024 * 1024)

		return fmt.Sprintf(`📊 *OMEG4MON ROUTER Go Telemetry*
• *Engine:* `+"`Golang %s`"+`
• *Uptime:* %s
• *Memory Allocated:* %.2f MB
• *Active Connections:* %d / %d
• *Proxy Pools:* %d live (%d total)`,
			runtime.Version(), formatUptime(time.Since(startTime)), allocMb, activeConns, totalConns, livePools, totalPools), nil

	case "providers":
		activeConns, totalConns, _ := db.GetActiveConnectionCount()
		return fmt.Sprintf(`🔌 *Active Providers*
• *Total Registered:* %d
• *Active Connections:* %d`, totalConns, activeConns), nil

	case "pools":
		live, total, _ := db.GetProxyPoolStats()
		dead := total - live
		return fmt.Sprintf(`🌐 *Proxy Pool Health*
• *Total Proxies:* %d
• *Live Proxies:* %d
• *Dead / Failed:* %d

_Gunakan `+"`/prune`"+` untuk membersihkan proxy mati._`, total, live, dead), nil

	case "prune":
		deleted, err := db.PruneDeadProxyPools()
		if err != nil {
			return fmt.Sprintf("❌ Gagal prune: %s", err.Error()), nil
		}
		if deleted == 0 {
			return "✅ Pool bersih! Tidak ada proxy mati yang perlu di-prune.", nil
		}
		return fmt.Sprintf("🧹 *Pruning Selesai!*\nBerhasil menghapus *%d* proxy mati dari database.", deleted), nil

	default:
		return fmt.Sprintf("❓ Perintah tidak dikenal: `/%s`\nKetik `/help` untuk melihat daftar perintah.", cmd), nil
	}
}

func HandleWebhook(w http.ResponseWriter, r *http.Request, db *database.DB) {
	cfg, err := db.GetTelegramConfig()
	if err != nil || !cfg.Enabled || cfg.BotToken == "" {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"ok":false,"error":"bot disabled"}`))
		return
	}

	// Only updates carrying the secret registered with setWebhook come from Telegram.
	secret := r.Header.Get("X-Telegram-Bot-Api-Secret-Token")
	if cfg.WebhookSecret == "" || subtle.ConstantTimeCompare([]byte(secret), []byte(cfg.WebhookSecret)) != 1 {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"ok":false,"error":"unauthorized"}`))
		return
	}
	allowed := strings.TrimSpace(cfg.UserID)
	if allowed == "" {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(`{"ok":false,"error":"userId not configured"}`))
		return
	}

	var update TelegramUpdate
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil || update.Message == nil {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"ok":true,"ignored":true}`))
		return
	}

	chatId := update.Message.Chat.ID
	fromId := update.Message.From.ID
	text := strings.TrimSpace(update.Message.Text)

	if chatId == 0 || text == "" {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"ok":true}`))
		return
	}

	if fmt.Sprintf("%d", fromId) != allowed && fmt.Sprintf("%d", chatId) != allowed {
		sendTelegramMessage(cfg.BotToken, chatId, "⛔ *Akses Ditolak*\nAnda tidak memiliki otorisasi untuk mengontrol OMEG4MON ROUTER ini.")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"ok":true,"authorized":false}`))
		return
	}

	replyText, err := HandleCommand(text, fmt.Sprintf("%d", fromId), cfg.BotToken, db)
	if err != nil {
		replyText = fmt.Sprintf("⚠️ Error: %s", err.Error())
	}

	sendTelegramMessage(cfg.BotToken, chatId, replyText)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"ok":true}`))
}

var telegramClient = &http.Client{Timeout: 10 * time.Second}

func sendTelegramMessage(botToken string, chatId int64, text string) {
	payload := map[string]interface{}{
		"chat_id":    chatId,
		"text":       text,
		"parse_mode": "Markdown",
	}
	bodyBytes, _ := json.Marshal(payload)
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", botToken)
	resp, err := telegramClient.Post(url, "application/json", bytes.NewBuffer(bodyBytes))
	if err != nil {
		// err embeds the URL, which carries the bot token: log only that it failed.
		log.Printf("[Telegram] sendMessage failed (network error)")
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Printf("[Telegram] sendMessage failed: HTTP %d", resp.StatusCode)
	}
}

func TestBotToken(botToken string) (map[string]interface{}, error) {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/getMe", botToken)
	client := &http.Client{Timeout: 8 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		OK          bool                   `json:"ok"`
		Description string                 `json:"description"`
		Result      map[string]interface{} `json:"result"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	if !result.OK {
		return nil, fmt.Errorf("telegram api error: %s", result.Description)
	}
	return result.Result, nil
}
