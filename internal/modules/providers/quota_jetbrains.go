package providers

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/omeg4monx/router-gateway/internal/contracts"
	"github.com/omeg4monx/router-gateway/internal/platform/httpclient"
)

const jetbrainsQuotaProbeTimeout = 6 * time.Second

type jetbrainsQuotaFetcher struct{}

func (jetbrainsQuotaFetcher) ID() string {
	return providerIDJetBrains
}

func (jetbrainsQuotaFetcher) Matches(conn ConnectionSummary) bool {
	if strings.EqualFold(strings.TrimSpace(conn.ProviderID), providerIDJetBrains) {
		return true
	}
	parsed, err := url.Parse(strings.TrimSpace(conn.BaseURL))
	if err != nil {
		return false
	}
	host := strings.ToLower(strings.TrimSpace(parsed.Hostname()))
	return host == "api.jetbrains.ai"
}

func resolveJetBrainsModelsURL(baseURL string) string {
	trimmed := strings.TrimRight(strings.TrimSpace(baseURL), "/")
	if trimmed == "" {
		return jetbrainsDefaultBaseURL + "/models"
	}
	trimmed = strings.TrimSuffix(trimmed, "/chat/completions")
	if strings.HasSuffix(trimmed, "/models") {
		return trimmed
	}
	return trimmed + "/models"
}

func extractJWTExpISO(token string) *string {
	payload := decodeJWTPayload(token)
	if payload == nil {
		return nil
	}
	rawExp, ok := payload["exp"].(float64)
	if !ok || rawExp <= 0 {
		return nil
	}
	iso := time.Unix(int64(rawExp), 0).UTC().Format(time.RFC3339)
	return &iso
}

func (jetbrainsQuotaFetcher) Fetch(
	ctx context.Context,
	client *http.Client,
	conn ConnectionSummary,
	creds contracts.UpstreamCredentials,
) (quotaFetchResult, error) {
	jwt := strings.TrimSpace(creds.ExtraHeaders["Grazie-Authenticate-JWT"])
	if jwt == "" {
		jwt = strings.TrimSpace(creds.AccessToken)
	}
	if jwt == "" {
		jwt = strings.TrimSpace(creds.APIKey)
	}
	if jwt == "" {
		return quotaFetchResult{}, errors.New("missing JetBrains Grazie JWT token")
	}

	baseURL := strings.TrimSpace(creds.BaseURL)
	if baseURL == "" {
		baseURL = conn.BaseURL
	}
	modelsURL := resolveJetBrainsModelsURL(baseURL)

	reqCtx, cancel := context.WithTimeout(httpclient.WithProviderID(ctx, conn.ProviderID), jetbrainsQuotaProbeTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(reqCtx, http.MethodGet, modelsURL, http.NoBody)
	if err != nil {
		return quotaFetchResult{}, fmt.Errorf("build jetbrains quota request: %w", err)
	}
	req.Header.Set("Grazie-Authenticate-JWT", jwt)
	req.Header.Set("Grazie-Agent", jetbrainsGrazieAgentHeader)
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err == nil {
		defer resp.Body.Close()
		_, _ = io.Copy(io.Discard, io.LimitReader(resp.Body, maxQuotaBodyBytes))
		if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
			return quotaFetchResult{}, fmt.Errorf("%w (HTTP %d)", errQuotaUnauthorized, resp.StatusCode)
		}
	}

	resetAt := extractJWTExpISO(jwt)
	planName := "JetBrains AI (Grazie Lite)"

	return quotaFetchResult{
		Status:    "ok",
		PlanName:  &planName,
		OnPlan:    true,
		Activated: true,
		Buckets: []ProviderQuotaBucket{
			{
				ID:           jetbrainsDefaultModel,
				Label:        "GPT-6 Luna (Grazie Lite)",
				UsedPct:      0,
				RemainingPct: 100,
				Exhausted:    false,
				ResetAt:      resetAt,
			},
		},
		Models: []string{jetbrainsDefaultModel},
	}, nil
}
