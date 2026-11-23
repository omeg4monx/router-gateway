package telegram

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	defaultAPIBaseURL  = "https://api.telegram.org"
	apiCallTimeout     = 10 * time.Second
	maxAPIResponseSize = 64 << 10 // 64 KiB
)

type botAPIClient struct {
	baseURL    string
	httpClient *http.Client
}

func newBotAPIClient(baseURL string, client *http.Client) *botAPIClient {
	base := strings.TrimRight(strings.TrimSpace(baseURL), "/")
	if base == "" {
		base = defaultAPIBaseURL
	}
	if client == nil {
		client = &http.Client{Timeout: apiCallTimeout}
	}
	return &botAPIClient{
		baseURL:    base,
		httpClient: client,
	}
}

func (c *botAPIClient) sendMessage(ctx context.Context, botToken, chatID, text string) error {
	payload := map[string]any{
		"chat_id":                  chatID,
		"text":                     truncateText(text, 4000),
		"disable_web_page_preview": true,
	}
	return c.post(ctx, botToken, "sendMessage", payload)
}

func (c *botAPIClient) setWebhook(ctx context.Context, botToken, webhookURL, secretToken string) error {
	payload := map[string]any{
		"url":          webhookURL,
		"secret_token": secretToken,
	}
	return c.post(ctx, botToken, "setWebhook", payload)
}

func (c *botAPIClient) deleteWebhook(ctx context.Context, botToken string) error {
	return c.post(ctx, botToken, "deleteWebhook", map[string]any{
		"drop_pending_updates": false,
	})
}

func (c *botAPIClient) post(ctx context.Context, botToken, method string, payload map[string]any) error {
	if strings.TrimSpace(botToken) == "" {
		return errors.New("telegram bot token is not configured")
	}
	callCtx, cancel := context.WithTimeout(ctx, apiCallTimeout)
	defer cancel()

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal telegram %s request: %w", method, err)
	}
	endpoint := fmt.Sprintf("%s/bot%s/%s", c.baseURL, botToken, method)
	req, err := http.NewRequestWithContext(callCtx, http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("build telegram %s request: %w", method, err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("telegram %s request failed: %w", method, err)
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(io.LimitReader(resp.Body, maxAPIResponseSize))
	if err != nil {
		return fmt.Errorf("read telegram %s response: %w", method, err)
	}

	var apiResp struct {
		OK          bool   `json:"ok"`
		Description string `json:"description"`
	}
	if err := json.Unmarshal(raw, &apiResp); err != nil {
		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			return fmt.Errorf("telegram %s returned HTTP %d", method, resp.StatusCode)
		}
		return fmt.Errorf("decode telegram %s response: %w", method, err)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 || !apiResp.OK {
		desc := strings.TrimSpace(apiResp.Description)
		if desc == "" {
			desc = fmt.Sprintf("HTTP %d", resp.StatusCode)
		}
		return fmt.Errorf("telegram %s error: %s", method, desc)
	}
	return nil
}

func truncateText(s string, maxLen int) string {
	if maxLen <= 0 || len(s) <= maxLen {
		return s
	}
	return s[:maxLen]
}
