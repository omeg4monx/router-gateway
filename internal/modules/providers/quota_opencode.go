package providers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/omeg4monx/router-gateway/internal/contracts"
	"github.com/omeg4monx/router-gateway/internal/platform/httpclient"
)

const opencodeQuotaTimeout = 8 * time.Second

type opencodeQuotaFetcher struct {
	providerID string
	planName   string
	usageURL   string
}

func (f opencodeQuotaFetcher) ID() string {
	return f.providerID
}

func (f opencodeQuotaFetcher) Matches(conn ConnectionSummary) bool {
	return strings.EqualFold(strings.TrimSpace(conn.ProviderID), f.providerID)
}

type opencodeUsagePeriod struct {
	Percent  any `json:"percent"`
	ResetsAt any `json:"resetsAt"`
}

type opencodeUsagePayload struct {
	Usage map[string]opencodeUsagePeriod `json:"usage"`
}

func parseFlexibleFloat(v any) (float64, bool) {
	switch n := v.(type) {
	case float64:
		if !math.IsNaN(n) && !math.IsInf(n, 0) {
			return n, true
		}
	case string:
		if parsed, err := strconv.ParseFloat(strings.TrimSpace(n), 64); err == nil && !math.IsNaN(parsed) && !math.IsInf(parsed, 0) {
			return parsed, true
		}
	}
	return 0, false
}

func parseFlexibleResetISO(v any) *string {
	switch raw := v.(type) {
	case string:
		trimmed := strings.TrimSpace(raw)
		if trimmed == "" {
			return nil
		}
		if t, err := time.Parse(time.RFC3339, trimmed); err == nil {
			iso := t.UTC().Format(time.RFC3339)
			return &iso
		}
		return &trimmed
	case float64:
		if raw <= 0 || math.IsNaN(raw) || math.IsInf(raw, 0) {
			return nil
		}
		sec := int64(raw)
		if sec > 1e12 {
			sec /= 1000
		}
		iso := time.Unix(sec, 0).UTC().Format(time.RFC3339)
		return &iso
	}
	return nil
}

func resolveOpenCodeUsageURL(baseURL, defaultUsageURL string) string {
	trimmed := strings.TrimRight(strings.TrimSpace(baseURL), "/")
	if trimmed == "" {
		return defaultUsageURL
	}
	for _, suffix := range []string{"/chat/completions", "/responses", "/messages", "/usage"} {
		trimmed = strings.TrimSuffix(trimmed, suffix)
	}
	return trimmed + "/usage"
}

func (f opencodeQuotaFetcher) Fetch(
	ctx context.Context,
	client *http.Client,
	conn ConnectionSummary,
	creds contracts.UpstreamCredentials,
) (quotaFetchResult, error) {
	apiKey := strings.TrimSpace(creds.APIKey)
	if apiKey == "" {
		apiKey = strings.TrimSpace(creds.AccessToken)
	}
	if apiKey == "" {
		return quotaFetchResult{}, fmt.Errorf("%s API key not available", f.planName)
	}

	baseURL := strings.TrimSpace(creds.BaseURL)
	if baseURL == "" {
		baseURL = conn.BaseURL
	}
	usageURL := resolveOpenCodeUsageURL(baseURL, f.usageURL)

	reqCtx, cancel := context.WithTimeout(httpclient.WithProviderID(ctx, conn.ProviderID), opencodeQuotaTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(reqCtx, http.MethodGet, usageURL, http.NoBody)
	if err != nil {
		return quotaFetchResult{}, fmt.Errorf("build %s usage request: %w", f.providerID, err)
	}
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return quotaFetchResult{}, fmt.Errorf("%s usage request failed: %w", f.providerID, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, maxQuotaBodyBytes))
	if err != nil {
		return quotaFetchResult{}, fmt.Errorf("read %s usage response: %w", f.providerID, err)
	}

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return quotaFetchResult{}, fmt.Errorf("%w (HTTP %d)", errQuotaUnauthorized, resp.StatusCode)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return quotaFetchResult{}, fmt.Errorf("%s usage API error (HTTP %d)", f.planName, resp.StatusCode)
	}

	var parsed opencodeUsagePayload
	if err := json.Unmarshal(body, &parsed); err != nil || len(parsed.Usage) == 0 {
		return quotaFetchResult{}, errors.New("opencode usage response did not contain quota data")
	}

	periodOrder := []struct {
		key   string
		label string
	}{
		{"rolling", "Rolling"},
		{"weekly", "Weekly"},
		{"monthly", "Monthly"},
	}

	buckets := make([]ProviderQuotaBucket, 0, len(periodOrder))
	for _, p := range periodOrder {
		entry, ok := parsed.Usage[p.key]
		if !ok {
			continue
		}
		pct, ok := parseFlexibleFloat(entry.Percent)
		if !ok {
			continue
		}
		used := math.Max(0, math.Min(100, pct))
		remaining := 100 - used
		buckets = append(buckets, ProviderQuotaBucket{
			ID:           p.key,
			Label:        p.label,
			UsedPct:      used,
			RemainingPct: remaining,
			Exhausted:    remaining <= 0,
			ResetAt:      parseFlexibleResetISO(entry.ResetsAt),
		})
	}
	if len(buckets) == 0 {
		return quotaFetchResult{}, errors.New("opencode usage response did not contain valid quota buckets")
	}

	planName := f.planName
	return quotaFetchResult{
		Status:    "ok",
		PlanName:  &planName,
		OnPlan:    true,
		Activated: true,
		Buckets:   buckets,
	}, nil
}
