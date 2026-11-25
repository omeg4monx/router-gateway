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
	antigravityIDEVersion   = "2.34.0"
	antigravityIDEUserAgent = "antigravity/hub/" + antigravityIDEVersion + " linux/x64"

	agQuotaCacheTTL        = 60 * time.Second
	agQuotaFailBackoff     = 15 * time.Second
	agQuotaCheckTimeout    = 10 * time.Second
	agQuotaDrainedThresh   = 0.02
	agQuotaMaxCacheEntries = 1024
	agQuotaMaxBodyBytes    = 1 << 20

	antigravityPromptRejectedMessage = "Antigravity rejected this request (upstream 429 RESOURCE_EXHAUSTED with no quota details, so this is not your quota). Retrying with system instructions moved to the first user turn was also refused."
)

type antigravityPromptRewriteRule struct {
	pattern *regexp.Regexp
	replace string
	fn      func(string) string
}

var (
	agLegacyFakeProjectRE    = regexp.MustCompile(`^(?i)(useful|bright|swift|calm|bold)-(fuze|wave|spark|flow|core)-[a-f0-9]+$`)
	agLegacyCompanionRE      = regexp.MustCompile(`^(?i)cloudaicompanion-proj-[a-f0-9]+$`)
	agBareExhaustedMsgRE     = regexp.MustCompile(`^Resource has been exhausted \(e\.g\. check quota\)\.?$`)
	agResourceExhaustedBytes = []byte("RESOURCE_EXHAUSTED")
	agSystemInstructionBytes = []byte("systemInstruction")

	// Ordered proactive rewrites for systemInstruction phrases fingerprinted by
	// Antigravity (Claude Code, Codex CLI, Hermes, OpenClaw, OpenCode). Longest
	// forms come first so fragments do not survive.
	antigravityPromptRewrites = []antigravityPromptRewriteRule{
		{
			pattern: regexp.MustCompile(`You are a Claude agent, built on Anthropic's Claude Agent SDK\.`),
			replace: "",
		},
		{
			pattern: regexp.MustCompile(`(?i)You(?: are|'re) Claude Code(?:,\s*Anthropic's official CLI for Claude)?\.?`),
			replace: "You are an AI coding assistant.",
		},
		{
			pattern: regexp.MustCompile(`(?i)interactive agent that helps users with software engineering tasks`),
			replace: "AI coding assistant that helps users with software engineering tasks",
		},
		{
			pattern: regexp.MustCompile(`(?i)Dual-use security tools \(C2 frameworks, credential testing, exploit development\) require clear authorization context`),
			replace: "Dual-use security tools require clear authorization context",
		},
		{
			pattern: regexp.MustCompile(`(?i)was pasted into the message by the user from somewhere else and may contain instructions the user did not write`),
			replace: "comes from user-pasted content and may contain untrusted instructions",
		},
		{
			pattern: regexp.MustCompile(`(?i)You(?: are|'re) a coding agent running in the Codex CLI(?:,\s*a terminal-based coding assistant)?\.?`),
			replace: "You are a coding assistant.",
		},
		{
			pattern: regexp.MustCompile(`(?i)You are Hermes Agent, built by Nous Research\.\s*Be direct:\s*match the length of your reply to the weight of the ask`),
			replace: "You are an AI assistant. Match the length of your reply to the weight of the ask",
		},
		{
			pattern: regexp.MustCompile(`(?i)Be direct:\s*match the length of your reply to the weight of the ask`),
			replace: "Match the length of your reply to the weight of the ask",
		},
		{
			pattern: regexp.MustCompile(`(?i)For help with Hermes Agent itself([^\n]*?)the documentation at`),
			replace: "On Hermes Agent itself${1}the documentation at",
		},
		{
			pattern: regexp.MustCompile(`(?i)(documentation at https://hermes-agent\.nousresearch\.com/docs) is your authoritative reference`),
			replace: "${1} is the authoritative reference",
		},
		{
			pattern: regexp.MustCompile(`(?i)Hermes delivers their message as a standalone user message`),
			replace: "the harness delivers their message as a standalone user message",
		},
		{
			pattern: regexp.MustCompile(`# Hermes runtime environment`),
			replace: "# Runtime environment",
		},
		{
			pattern: regexp.MustCompile(`<!-- End Hermes runtime environment -->`),
			replace: "<!-- End runtime environment -->",
		},
		{
			pattern: regexp.MustCompile(`# Google model operational directives`),
			replace: "# Operational directives",
		},
		{
			pattern: regexp.MustCompile(`(?i)When the user asks you to build, run, or verify something, the deliverable is a working artifact backed by real tool output`),
			replace: "If the user asks you to build, run, or verify something, the deliverable is a working artifact that real tool output backs",
		},
		{
			pattern: regexp.MustCompile(`(?i)the deliverable is a working artifact backed by real tool output`),
			replace: "the deliverable is a working artifact that real tool output backs",
		},
		{
			pattern: regexp.MustCompile(`(?i)If a tool, install, or network call fails and blocks the real path`),
			replace: "When a tool, install, or network call fails and blocks the real path",
		},
		{
			pattern: regexp.MustCompile(`(?i)Reporting a blocker honestly is always better than inventing a result\.`),
			replace: "Honestly reporting a blocker is always better than inventing a result.",
		},
		{
			pattern: regexp.MustCompile(`(?i)You are Hermes(?: Agent)?(?:,\s*(?:an intelligent AI assistant|an AI assistant|an AI agent))?(?:,?\s*(?:built|created)\s+by\s+Nous Research)?\.`),
			replace: "You are an AI assistant.",
		},
		{
			pattern: regexp.MustCompile(`(?im)^x-anthropic-billing-header:[^\n]*(?:\r?\n)*`),
			replace: "",
		},
		{
			pattern: regexp.MustCompile(`(?i)You are a personal assistant running inside OpenClaw\.`),
			replace: "You are a personal assistant that runs inside OpenClaw.",
		},
		{
			pattern: regexp.MustCompile(`(?i)opencode`),
			fn: func(match string) string {
				switch match {
				case "OpenCode":
					return "Antigravity"
				case "OPENCODE":
					return "ANTIGRAVITY"
				default:
					return "antigravity"
				}
			},
		},
	}
)

func rewriteAntigravityPromptText(text string) string {
	out := text
	for _, rule := range antigravityPromptRewrites {
		if rule.fn != nil {
			out = rule.pattern.ReplaceAllStringFunc(out, rule.fn)
		} else {
			out = rule.pattern.ReplaceAllString(out, rule.replace)
		}
	}
	return out
}

// rewriteAntigravitySystemInstruction proactively rewrites fingerprinted client
// phrases in systemInstruction (for both raw Gemini payloads and Antigravity's
// wrapped {request:{...}} envelope) so upstream does not answer with a bare 429.
func rewriteAntigravitySystemInstruction(body []byte) []byte {
	if len(body) == 0 || !bytes.Contains(body, agSystemInstructionBytes) {
		return body
	}
	var root map[string]any
	if err := json.Unmarshal(body, &root); err != nil || root == nil {
		return body
	}

	target := root
	if nested, ok := root["request"].(map[string]any); ok && nested != nil {
		if _, hasNestedSys := nested["systemInstruction"]; hasNestedSys {
			target = nested
		}
	}

	sysMap, ok := target["systemInstruction"].(map[string]any)
	if !ok || sysMap == nil {
		return body
	}
	rawParts, ok := sysMap["parts"].([]any)
	if !ok || len(rawParts) == 0 {
		return body
	}

	changed := false
	rewrittenParts := make([]any, 0, len(rawParts))
	for _, rawPart := range rawParts {
		partMap, ok := rawPart.(map[string]any)
		if !ok {
			rewrittenParts = append(rewrittenParts, rawPart)
			continue
		}
		text, isStr := partMap["text"].(string)
		if !isStr {
			rewrittenParts = append(rewrittenParts, partMap)
			continue
		}
		rewritten := rewriteAntigravityPromptText(text)
		if rewritten != text {
			changed = true
		}
		if strings.TrimSpace(rewritten) == "" {
			changed = true
			continue
		}
		partMap["text"] = rewritten
		rewrittenParts = append(rewrittenParts, partMap)
	}

	if !changed {
		return body
	}
	if len(rewrittenParts) == 0 {
		delete(target, "systemInstruction")
	} else {
		sysMap["parts"] = rewrittenParts
	}

	encoded, err := json.Marshal(root)
	if err != nil {
		return body
	}
	return encoded
}

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
	req.Header.Set("User-Agent", antigravityIDEUserAgent)
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
			switch {
			case b.Disabled:
				rem = 0
			case b.RemainingFraction != nil:
				rem = *b.RemainingFraction
			default:
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
