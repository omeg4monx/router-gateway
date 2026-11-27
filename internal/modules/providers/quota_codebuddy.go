package providers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/omeg4monx/router-gateway/internal/contracts"
	"github.com/omeg4monx/router-gateway/internal/platform/httpclient"
)

const codebuddyRefillGap = 48 * time.Hour

type codebuddyQuotaFetcher struct {
	providerID string
}

type codebuddyQuotaEnvelope struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		Response struct {
			Data struct {
				Accounts []codebuddyQuotaAccount `json:"Accounts"`
			} `json:"Data"`
		} `json:"Response"`
	} `json:"data"`
}

type codebuddyQuotaAccount struct {
	PackageName              string `json:"PackageName"`
	SubProductName           string `json:"SubProductName"`
	CycleStartTime           any    `json:"CycleStartTime"`
	CycleEndTime             any    `json:"CycleEndTime"`
	DeductionEndTime         any    `json:"DeductionEndTime"`
	CycleCapacityUsedPrecise any    `json:"CycleCapacityUsedPrecise"`
	CycleCapacityUsed        any    `json:"CycleCapacityUsed"`
	CycleCapacitySizePrecise any    `json:"CycleCapacitySizePrecise"`
	CycleCapacitySize        any    `json:"CycleCapacitySize"`
	CapacityUsedPrecise      any    `json:"CapacityUsedPrecise"`
	CapacityUsed             any    `json:"CapacityUsed"`
	CapacitySizePrecise      any    `json:"CapacitySizePrecise"`
	CapacitySize             any    `json:"CapacitySize"`
}

func (f codebuddyQuotaFetcher) ID() string {
	return f.providerID
}

func (f codebuddyQuotaFetcher) Matches(conn ConnectionSummary) bool {
	if strings.EqualFold(strings.TrimSpace(conn.ProviderID), f.providerID) {
		return true
	}
	return isCodeBuddyHostFor(f.providerID, conn.BaseURL)
}

func isCodeBuddyHostFor(providerID, rawURL string) bool {
	trimmed := strings.TrimSpace(rawURL)
	if trimmed == "" {
		return false
	}
	parsed, err := url.Parse(trimmed)
	if err != nil {
		return false
	}
	host := strings.ToLower(strings.TrimSpace(parsed.Hostname()))
	if providerID == providerIDCodeBuddyIntl {
		return host == "codebuddy.ai" || strings.HasSuffix(host, ".codebuddy.ai")
	}
	return host == "copilot.tencent.com" || strings.HasSuffix(host, ".copilot.tencent.com")
}

func resolveCodeBuddyQuotaURL(providerID, baseURL string) string {
	trimmed := strings.TrimSpace(baseURL)
	if trimmed != "" {
		if parsed, err := url.Parse(trimmed); err == nil &&
			(strings.EqualFold(parsed.Scheme, "http") || strings.EqualFold(parsed.Scheme, "https")) &&
			parsed.Host != "" {
			return parsed.Scheme + "://" + parsed.Host + codebuddyQuotaPath
		}
	}
	if providerID == providerIDCodeBuddyIntl {
		return codebuddyIntlOrigin + codebuddyQuotaPath
	}
	return codebuddyCNOrigin + codebuddyQuotaPath
}

func (f codebuddyQuotaFetcher) Fetch(
	ctx context.Context,
	client *http.Client,
	conn ConnectionSummary,
	creds contracts.UpstreamCredentials,
) (quotaFetchResult, error) {
	token := strings.TrimSpace(creds.AccessToken)
	if token == "" {
		token = strings.TrimSpace(creds.APIKey)
	}
	if token == "" {
		return quotaFetchResult{}, errors.New("missing access token or API key")
	}

	baseURL := strings.TrimSpace(creds.BaseURL)
	if baseURL == "" {
		baseURL = conn.BaseURL
	}
	endpoint := resolveCodeBuddyQuotaURL(f.providerID, baseURL)

	reqCtx, cancel := context.WithTimeout(httpclient.WithProviderID(ctx, conn.ProviderID), quotaFetchTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(reqCtx, http.MethodPost, endpoint, bytes.NewReader([]byte("{}")))
	if err != nil {
		return quotaFetchResult{}, fmt.Errorf("build codebuddy quota request: %w", err)
	}
	applyCodeBuddyProbeHeaders(req, f.providerID)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		return quotaFetchResult{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, maxQuotaBodyBytes))
	if err != nil {
		return quotaFetchResult{}, fmt.Errorf("read codebuddy quota response: %w", err)
	}

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return quotaFetchResult{}, fmt.Errorf("upstream rejected credentials (HTTP %d)", resp.StatusCode)
	}
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		snippet := strings.TrimSpace(string(body))
		if snippet != "" {
			return quotaFetchResult{}, fmt.Errorf("upstream returned HTTP %d: %s", resp.StatusCode, snippet)
		}
		return quotaFetchResult{}, fmt.Errorf("upstream returned HTTP %d", resp.StatusCode)
	}

	var payload codebuddyQuotaEnvelope
	if err := json.Unmarshal(body, &payload); err != nil {
		return quotaFetchResult{}, fmt.Errorf("decode CodeBuddy quota response: %w", err)
	}
	if payload.Code != 0 {
		msg := strings.TrimSpace(payload.Msg)
		if msg == "" {
			msg = fmt.Sprintf("code %d", payload.Code)
		}
		return quotaFetchResult{}, fmt.Errorf("CodeBuddy quota error: %s", msg)
	}

	return mapCodeBuddyQuota(payload.Data.Response.Data.Accounts), nil
}

func mapCodeBuddyQuota(accounts []codebuddyQuotaAccount) quotaFetchResult {
	if len(accounts) == 0 {
		plan := "CodeBuddy"
		return quotaFetchResult{
			Status:    "not_activated",
			PlanName:  &plan,
			OnPlan:    false,
			Activated: false,
			Buckets:   []ProviderQuotaBucket{},
			Models:    []string{},
		}
	}

	isRefill := func(acc codebuddyQuotaAccount) bool {
		ce, okCE := parseCodeBuddyTimestamp(acc.CycleEndTime)
		de, okDE := parseCodeBuddyTimestamp(acc.DeductionEndTime)
		return okCE && okDE && de.Sub(ce) > codebuddyRefillGap
	}

	var refills, bonuses []codebuddyQuotaAccount
	for _, acc := range accounts {
		if isRefill(acc) {
			refills = append(refills, acc)
		} else {
			bonuses = append(bonuses, acc)
		}
	}

	sortByCycleEnd := func(list []codebuddyQuotaAccount) {
		sort.SliceStable(list, func(i, j int) bool {
			ti, okI := parseCodeBuddyTimestamp(list[i].CycleEndTime)
			tj, okJ := parseCodeBuddyTimestamp(list[j].CycleEndTime)
			if okI && okJ {
				return ti.Before(tj)
			}
			return okI
		})
	}
	sortByCycleEnd(refills)
	sortByCycleEnd(bonuses)

	buckets := make([]ProviderQuotaBucket, 0, len(refills)+len(bonuses))
	seenCadence := make(map[string]int)

	for _, acc := range refills {
		baseLabel := codebuddyRefillCadence(acc)
		seenCadence[baseLabel]++
		label := baseLabel
		if seenCadence[baseLabel] > 1 {
			label = fmt.Sprintf("%s %d", baseLabel, seenCadence[baseLabel])
		}
		used := parseCodeBuddyNumber(acc.CycleCapacityUsedPrecise, acc.CycleCapacityUsed)
		total := parseCodeBuddyNumber(acc.CycleCapacitySizePrecise, acc.CycleCapacitySize)
		usedPct, remainingPct, exhausted := computeBucketPercentages(used, total)
		startISO := formatCodeBuddyTimestamp(acc.CycleStartTime)
		endISO := formatCodeBuddyTimestamp(acc.CycleEndTime)

		buckets = append(buckets, ProviderQuotaBucket{
			ID:              "refill_" + strings.ToLower(strings.ReplaceAll(label, " ", "_")),
			Label:           label,
			UsedPct:         usedPct,
			RemainingPct:    remainingPct,
			Exhausted:       exhausted,
			PeriodStartedAt: startISO,
			PeriodEndsAt:    cloneOptionalString(endISO),
			ResetAt:         endISO,
		})
	}

	for i, acc := range bonuses {
		used := parseCodeBuddyNumber(acc.CapacityUsedPrecise, acc.CapacityUsed)
		total := parseCodeBuddyNumber(acc.CapacitySizePrecise, acc.CapacitySize)
		usedPct, remainingPct, exhausted := computeBucketPercentages(used, total)
		startISO := formatCodeBuddyTimestamp(acc.CycleStartTime)
		endISO := formatCodeBuddyTimestamp(acc.CycleEndTime)

		buckets = append(buckets, ProviderQuotaBucket{
			ID:              fmt.Sprintf("bonus_%d", i+1),
			Label:           fmt.Sprintf("Bonus Pack %d", i+1),
			UsedPct:         usedPct,
			RemainingPct:    remainingPct,
			Exhausted:       exhausted,
			PeriodStartedAt: startISO,
			PeriodEndsAt:    cloneOptionalString(endISO),
			ResetAt:         endISO,
		})
	}

	basePkg := accounts[0]
	if len(refills) > 0 {
		basePkg = refills[0]
	}
	plan := strings.TrimSpace(basePkg.PackageName)
	if plan == "" {
		plan = strings.TrimSpace(basePkg.SubProductName)
	}
	if plan == "" {
		plan = "CodeBuddy"
	}

	allExhausted := len(buckets) > 0
	anyLow := false
	for _, b := range buckets {
		if !b.Exhausted {
			allExhausted = false
		}
		if b.Exhausted || b.RemainingPct <= quotaLowRemainingPct {
			anyLow = true
		}
	}
	status := "ok"
	switch {
	case allExhausted:
		status = "exhausted"
	case anyLow:
		status = "low"
	}

	return quotaFetchResult{
		Status:    status,
		PlanName:  &plan,
		OnPlan:    len(refills) > 0,
		Activated: true,
		Buckets:   buckets,
		Models:    []string{},
	}
}

func computeBucketPercentages(used, total float64) (usedPct, remainingPct float64, exhausted bool) {
	if total <= 0 {
		if used > 0 {
			return 100, 0, true
		}
		return 0, 100, false
	}
	usedPct = clampPercent((used / total) * 100)
	remainingPct = clampPercent(100 - usedPct)
	exhausted = usedPct >= 100 || used >= total
	return usedPct, remainingPct, exhausted
}

func codebuddyRefillCadence(acc codebuddyQuotaAccount) string {
	start, okStart := parseCodeBuddyTimestamp(acc.CycleStartTime)
	end, okEnd := parseCodeBuddyTimestamp(acc.CycleEndTime)
	if okStart && okEnd && end.After(start) {
		days := end.Sub(start).Hours() / 24
		if days <= 1.5 {
			return "Daily"
		}
		if days <= 10 {
			return "Weekly"
		}
	}
	return "Monthly"
}

func parseCodeBuddyNumber(precise, plain any) float64 {
	for _, candidate := range []any{precise, plain} {
		if candidate == nil {
			continue
		}
		switch v := candidate.(type) {
		case float64:
			if !math.IsNaN(v) && !math.IsInf(v, 0) {
				return v
			}
		case string:
			trimmed := strings.TrimSpace(v)
			if trimmed == "" {
				continue
			}
			if f, err := strconv.ParseFloat(trimmed, 64); err == nil && !math.IsNaN(f) && !math.IsInf(f, 0) {
				return f
			}
		case json.Number:
			if f, err := v.Float64(); err == nil && !math.IsNaN(f) && !math.IsInf(f, 0) {
				return f
			}
		}
	}
	return 0
}

func parseCodeBuddyTimestamp(raw any) (time.Time, bool) {
	if raw == nil {
		return time.Time{}, false
	}
	switch v := raw.(type) {
	case float64:
		return unixNumberToTime(v)
	case int64:
		return unixNumberToTime(float64(v))
	case string:
		trimmed := strings.TrimSpace(v)
		if trimmed == "" {
			return time.Time{}, false
		}
		if f, err := strconv.ParseFloat(trimmed, 64); err == nil {
			return unixNumberToTime(f)
		}
		for _, layout := range []string{time.RFC3339, time.RFC3339Nano, "2006-01-02 15:04:05"} {
			if parsed, err := time.Parse(layout, trimmed); err == nil {
				return parsed.UTC(), true
			}
		}
	}
	return time.Time{}, false
}

func unixNumberToTime(v float64) (time.Time, bool) {
	if math.IsNaN(v) || math.IsInf(v, 0) || v <= 0 {
		return time.Time{}, false
	}
	if v >= 1e12 {
		return time.UnixMilli(int64(v)).UTC(), true
	}
	return time.Unix(int64(v), 0).UTC(), true
}

func formatCodeBuddyTimestamp(raw any) *string {
	t, ok := parseCodeBuddyTimestamp(raw)
	if !ok {
		return nil
	}
	s := t.UTC().Format(time.RFC3339)
	return &s
}
