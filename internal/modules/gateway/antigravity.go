package gateway

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/omeg4monx/router-gateway/internal/platform/httpclient"
)

const (
	agQuotaCacheTTL        = 60 * time.Second
	agQuotaFailBackoff     = 15 * time.Second
	agQuotaCheckTimeout    = 10 * time.Second
	agQuotaDrainedThresh   = 0.02
	agQuotaMaxCacheEntries = 1024
	agQuotaMaxBodyBytes    = 1 << 20

	antigravityPromptRejectedMessage = "Antigravity rejected this request (upstream 429 RESOURCE_EXHAUSTED with no quota details, so this is not your quota). Retrying with system instructions moved to the first user turn was also refused."
)

var (
	agLegacyFakeProjectRE    = regexp.MustCompile(`^(?i)(useful|bright|swift|calm|bold)-(fuze|wave|spark|flow|core)-[a-f0-9]+$`)
	agLegacyCompanionRE      = regexp.MustCompile(`^(?i)cloudaicompanion-proj-[a-f0-9]+$`)
	agBareExhaustedMsgRE     = regexp.MustCompile(`^Resource has been exhausted \(e\.g\. check quota\)\.?$`)
	agResourceExhaustedBytes = []byte("RESOURCE_EXHAUSTED")
)

func isBlockedAntigravityProjectID(projectID string) bool {
	trimmed := strings.TrimSpace(projectID)
	if trimmed == "" || strings.EqualFold(trimmed, "aicode-consumers") {
		return true
	}
	return agLegacyFakeProjectRE.MatchString(trimmed) || agLegacyCompanionRE.MatchString(trimmed)
}

// isAntigravityPromptRejection returns true only for a 429 whose JSON body is
// a bare RESOURCE_EXHAUSTED with no structured error details. Real quota 429s
// carry RetryInfo / ErrorInfo / quotaResetTimeStamp details and return false.
func isAntigravityPromptRejection(statusCode int, body []byte) bool {
	if statusCode != http.StatusTooManyRequests || len(body) == 0 || !bytes.Contains(body, agResourceExhaustedBytes) {
		return false
	}
	var env struct {
		Error struct {
			Status  string            `json:"status"`
			Message string            `json:"message"`
			Details []json.RawMessage `json:"details"`
		} `json:"error"`
	}
	if err := json.Unmarshal(body, &env); err != nil {
		return false
	}
	if env.Error.Status != "RESOURCE_EXHAUSTED" {
		return false
	}
	if !agBareExhaustedMsgRE.MatchString(strings.TrimSpace(env.Error.Message)) {
		return false
	}
	return len(env.Error.Details) == 0
}

// moveGeminiSystemInstructionToContents moves non-empty text parts from
// systemInstruction into a leading user turn in contents (handling both raw
// Gemini payloads and Antigravity's wrapped {request:{...}} envelope).
func moveGeminiSystemInstructionToContents(body []byte) ([]byte, bool) {
	var root map[string]any
	if err := json.Unmarshal(body, &root); err != nil || root == nil {
		return nil, false
	}

	target := root
	if nested, ok := root["request"].(map[string]any); ok && nested != nil {
		if _, hasNestedSys := nested["systemInstruction"]; hasNestedSys {
			target = nested
		}
	}

	sysRaw, ok := target["systemInstruction"]
	if !ok || sysRaw == nil {
		return nil, false
	}
	sysMap, ok := sysRaw.(map[string]any)
	if !ok {
		return nil, false
	}
	rawParts, ok := sysMap["parts"].([]any)
	if !ok || len(rawParts) == 0 {
		return nil, false
	}

	var keptParts []any
	for _, p := range rawParts {
		pm, ok := p.(map[string]any)
		if !ok {
			continue
		}
		if text, ok := pm["text"].(string); ok && strings.TrimSpace(text) != "" {
			keptParts = append(keptParts, pm)
		}
	}
	if len(keptParts) == 0 {
		return nil, false
	}

	delete(target, "systemInstruction")
	var existingContents []any
	if c, ok := target["contents"].([]any); ok {
		existingContents = c
	}
	userTurn := map[string]any{
		"role":  "user",
		"parts": keptParts,
	}
	newContents := make([]any, 0, 1+len(existingContents))
	newContents = append(newContents, userTurn)
	newContents = append(newContents, existingContents...)
	target["contents"] = newContents

	encoded, err := json.Marshal(root)
	if err != nil {
		return nil, false
	}
	return encoded, true
}

type agQuotaBucket struct {
	Window            string   `json:"window"`
	BucketID          string   `json:"bucketId"`
	DisplayName       string   `json:"displayName"`
	RemainingFraction *float64 `json:"remainingFraction"`
	ResetTime         string   `json:"resetTime"`
	Disabled          bool     `json:"disabled"`
}

type agQuotaGroup struct {
	DisplayName string          `json:"displayName"`
	Buckets     []agQuotaBucket `json:"buckets"`
}

type agQuotaSummary struct {
	Groups           []agQuotaGroup `json:"groups"`
	QuotaGroups      []agQuotaGroup `json:"quotaGroups"`
	ModelQuotaGroups []agQuotaGroup `json:"modelQuotaGroups"`
}

func (s *agQuotaSummary) allGroups() []agQuotaGroup {
	if s == nil {
		return nil
	}
	switch {
	case len(s.Groups) > 0:
		return s.Groups
	case len(s.QuotaGroups) > 0:
		return s.QuotaGroups
	case len(s.ModelQuotaGroups) > 0:
		return s.ModelQuotaGroups
	default:
		return nil
	}
}

type agQuotaCacheEntry struct {
	summary   *agQuotaSummary
	expiresAt time.Time
	done      chan struct{}
}

type antigravityQuotaMonitor struct {
	mu    sync.Mutex
	cache map[string]*agQuotaCacheEntry
}

func newAntigravityQuotaMonitor() *antigravityQuotaMonitor {
	return &antigravityQuotaMonitor{
		cache: make(map[string]*agQuotaCacheEntry),
	}
}

func agQuotaCacheKey(connectionID, accessToken string) string {
	sum := sha256.Sum256([]byte(accessToken))
	return connectionID + ":" + hex.EncodeToString(sum[:8])
}

func (q *antigravityQuotaMonitor) checkDrained(
	ctx context.Context,
	client *http.Client,
	baseURL string,
	accessToken string,
	model string,
	connectionID string,
) (bool, time.Duration) {
	if q == nil || client == nil || ctx.Err() != nil || strings.TrimSpace(baseURL) == "" || strings.TrimSpace(accessToken) == "" || strings.TrimSpace(model) == "" {
		return false, 0
	}

	key := agQuotaCacheKey(connectionID, accessToken)
	now := time.Now().UTC()

	q.mu.Lock()
	if entry, ok := q.cache[key]; ok {
		if entry.done != nil {
			done := entry.done
			q.mu.Unlock()
			select {
			case <-done:
				return evaluateAntigravityQuota(entry.summary, model, time.Now().UTC())
			case <-ctx.Done():
				return false, 0
			}
		}
		if now.Before(entry.expiresAt) {
			summary := entry.summary
			q.mu.Unlock()
			return evaluateAntigravityQuota(summary, model, now)
		}
		delete(q.cache, key)
	}

	if len(q.cache) >= agQuotaMaxCacheEntries {
		for k, v := range q.cache {
			if v.done == nil && !now.Before(v.expiresAt) {
				delete(q.cache, k)
			}
		}
		if len(q.cache) >= agQuotaMaxCacheEntries {
			var oldestKey string
			var oldestExp time.Time
			for k, v := range q.cache {
				if v.done != nil {
					continue
				}
				if oldestKey == "" || v.expiresAt.Before(oldestExp) {
					oldestKey = k
					oldestExp = v.expiresAt
				}
			}
			if oldestKey != "" {
				delete(q.cache, oldestKey)
			}
		}
	}

	inflight := &agQuotaCacheEntry{done: make(chan struct{})}
	q.cache[key] = inflight
	q.mu.Unlock()

	summary := fetchAntigravityQuotaSummary(ctx, client, baseURL, accessToken)
	ttl := agQuotaCacheTTL
	if summary == nil {
		ttl = agQuotaFailBackoff
	}

	q.mu.Lock()
	inflight.summary = summary
	inflight.expiresAt = time.Now().UTC().Add(ttl)
	done := inflight.done
	inflight.done = nil
	close(done)
	q.mu.Unlock()

	return evaluateAntigravityQuota(summary, model, time.Now().UTC())
}

func fetchAntigravityQuotaSummary(ctx context.Context, client *http.Client, baseURL, accessToken string) *agQuotaSummary {
	endpoint := strings.TrimRight(strings.TrimSpace(baseURL), "/") + "/v1internal:retrieveUserQuotaSummary"
	reqCtx, cancel := context.WithTimeout(httpclient.WithProviderID(ctx, "antigravity"), agQuotaCheckTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(reqCtx, http.MethodPost, endpoint, bytes.NewReader([]byte(`{}`)))
	if err != nil {
		return nil
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-Client-Name", "antigravity")

	resp, err := client.Do(req)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil
	}
	raw, err := io.ReadAll(io.LimitReader(resp.Body, agQuotaMaxBodyBytes))
	if err != nil || len(raw) == 0 {
		return nil
	}

	var summary agQuotaSummary
	if err := json.Unmarshal(raw, &summary); err != nil {
		return nil
	}
	if len(summary.allGroups()) == 0 {
		return nil
	}
	return &summary
}

func evaluateAntigravityQuota(summary *agQuotaSummary, model string, now time.Time) (bool, time.Duration) {
	groups := summary.allGroups()
	if len(groups) == 0 {
		return false, 0
	}

	family := agModelFamily(model)
	found := false
	minRemaining := 1.0
	var resetStr string

	for _, g := range groups {
		if !agGroupMatchesFamily(g.DisplayName, family) {
			continue
		}
		for _, b := range g.Buckets {
			isWeekly := agIsWeeklyBucket(b)
			isSession := agIsSessionBucket(b)
			if !isWeekly && !isSession {
				continue
			}
			if b.Disabled && isWeekly {
				continue
			}
			var rem float64
			if b.Disabled {
				rem = 0
			} else if b.RemainingFraction != nil {
				rem = *b.RemainingFraction
			} else {
				continue
			}
			if !found || rem < minRemaining {
				found = true
				minRemaining = rem
				resetStr = b.ResetTime
			}
		}
	}

	if !found || minRemaining > agQuotaDrainedThresh {
		return false, 0
	}

	var retryAfter time.Duration
	if parsed, err := parseAntigravityResetTime(resetStr); err == nil {
		if d := parsed.Sub(now); d > 0 {
			retryAfter = d
		}
	}
	return true, retryAfter
}

func agModelFamily(model string) string {
	lower := strings.ToLower(strings.TrimSpace(model))
	switch {
	case strings.Contains(lower, "claude"), strings.Contains(lower, "gpt"), strings.HasPrefix(lower, "o1"), strings.HasPrefix(lower, "o3"), strings.HasPrefix(lower, "o4"):
		return "claude-gpt"
	case strings.Contains(lower, "argon"):
		return "argon"
	default:
		return "gemini"
	}
}

func agGroupMatchesFamily(displayName, family string) bool {
	lower := strings.ToLower(displayName)
	switch family {
	case "claude-gpt":
		return strings.Contains(lower, "claude") || strings.Contains(lower, "gpt")
	case "argon":
		return strings.Contains(lower, "argon")
	default:
		return strings.Contains(lower, "gemini")
	}
}

func agIsSessionBucket(b agQuotaBucket) bool {
	w := strings.ToLower(strings.TrimSpace(b.Window))
	text := strings.ToLower(b.BucketID + " " + b.DisplayName)
	return w == "5h" || w == "daily" || strings.Contains(text, "5h") || strings.Contains(text, "five hour") || strings.Contains(text, "daily")
}

func agIsWeeklyBucket(b agQuotaBucket) bool {
	w := strings.ToLower(strings.TrimSpace(b.Window))
	text := strings.ToLower(b.BucketID + " " + b.DisplayName)
	return w == "weekly" || strings.Contains(text, "weekly")
}

func parseAntigravityResetTime(raw string) (time.Time, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return time.Time{}, io.EOF
	}
	if ts, err := time.Parse(time.RFC3339Nano, trimmed); err == nil {
		return ts, nil
	}
	return time.Parse(time.RFC3339, trimmed)
}
