package promptinjection

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

const (
	ModeBlock    = "block"
	ModeSanitize = "sanitize"
	ModeAudit    = "audit"

	SensitivityLow    = "low"
	SensitivityMedium = "medium"
	SensitivityHigh   = "high"

	MatchSubstring = "substring"
	MatchRegex     = "regex"

	SeverityLow    = "low"
	SeverityMedium = "medium"
	SeverityHigh   = "high"

	sanitizedMarker  = "[REDACTED_INJECTION]"
	maxPatternBytes  = 512
	maxSnippetLength = 200
)

var secretRedactionPattern = regexp.MustCompile(`(?i)\b(sk-[a-z0-9_-]{8,}|bearer\s+[a-z0-9._-]{8,}|ghp_[a-z0-9]{8,}|AIza[0-9A-Za-z_-]{8,})`)

// MatchHit describes a single rule match during prompt inspection.
type MatchHit struct {
	RuleID   string `json:"rule_id"`
	RuleName string `json:"rule_name"`
	Action   string `json:"action"`
	Severity string `json:"severity"`
	Matched  string `json:"matched"`
}

// TestResult is returned by POST /api/prompt-guard/test.
type TestResult struct {
	Matched       bool       `json:"matched"`
	Action        string     `json:"action"`
	Hits          []MatchHit `json:"hits"`
	SanitizedText string     `json:"sanitized_text"`
}

type compiledRule struct {
	rule Rule
	re   *regexp.Regexp
}

func builtinRules(now string) []Rule {
	return []Rule{
		{
			ID:          "builtin-ignore-instructions",
			Name:        "Instruction Override",
			Pattern:     `(?i)\b(ignore|disregard|forget)\s+(all\s+)?(previous|prior|above|system)\s+(instructions|prompts|rules|directives)\b`,
			MatchType:   MatchRegex,
			Action:      ModeBlock,
			Severity:    SeverityHigh,
			Enabled:     true,
			Builtin:     true,
			Description: "Detects attempts to override system or prior instructions.",
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		{
			ID:          "builtin-system-prompt-exfil",
			Name:        "System Prompt Exfiltration",
			Pattern:     `(?i)\b(reveal|print|dump|repeat|output|show)\s+(me\s+)?(your|the)\s+(full\s+)?(system\s+prompt|hidden\s+instructions|developer\s+message)\b`,
			MatchType:   MatchRegex,
			Action:      ModeBlock,
			Severity:    SeverityHigh,
			Enabled:     true,
			Builtin:     true,
			Description: "Detects attempts to extract system or developer prompts.",
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		{
			ID:          "builtin-jailbreak-roleplay",
			Name:        "Jailbreak Persona Override",
			Pattern:     `(?i)\b(you\s+are\s+now\s+(DAN|unrestricted|jailbroken)|enter\s+developer\s+mode|bypass\s+safety\s+filters)\b`,
			MatchType:   MatchRegex,
			Action:      ModeBlock,
			Severity:    SeverityMedium,
			Enabled:     true,
			Builtin:     true,
			Description: "Detects classic persona jailbreak and safety bypass prompts.",
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		{
			ID:          "builtin-delimiter-spoofing",
			Name:        "Control Delimiter Spoofing",
			Pattern:     `(?i)(<\|im_start\|>system|<\|system\|>|###\s*SYSTEM\s*OVERRIDE\s*###)`,
			MatchType:   MatchRegex,
			Action:      ModeSanitize,
			Severity:    SeverityLow,
			Enabled:     true,
			Builtin:     true,
			Description: "Detects injected chat-template control delimiters.",
			CreatedAt:   now,
			UpdatedAt:   now,
		},
	}
}

func compileRule(r Rule) (compiledRule, error) {
	if r.Pattern == "" || len(r.Pattern) > maxPatternBytes {
		return compiledRule{}, fmt.Errorf("pattern must be 1..%d bytes", maxPatternBytes)
	}
	if r.MatchType == MatchRegex {
		re, err := regexp.Compile(r.Pattern)
		if err != nil {
			return compiledRule{}, fmt.Errorf("invalid regex pattern: %w", err)
		}
		return compiledRule{rule: r, re: re}, nil
	}
	return compiledRule{rule: r}, nil
}

func severityRank(s string) int {
	switch s {
	case SeverityHigh:
		return 3
	case SeverityMedium:
		return 2
	default:
		return 1
	}
}

func minSeverityForSensitivity(sensitivity string) int {
	switch sensitivity {
	case SensitivityLow:
		return severityRank(SeverityHigh)
	case SensitivityHigh:
		return severityRank(SeverityLow)
	default:
		return severityRank(SeverityMedium)
	}
}

func evaluateText(text string, cfg Config, rules []compiledRule) (string, []MatchHit) {
	minRank := minSeverityForSensitivity(cfg.Sensitivity)
	sanitized := text
	hits := make([]MatchHit, 0)

	for _, cr := range rules {
		if !cr.rule.Enabled || severityRank(cr.rule.Severity) < minRank {
			continue
		}
		if cr.rule.MatchType == MatchRegex && cr.re != nil {
			if loc := cr.re.FindString(sanitized); loc != "" {
				hits = append(hits, MatchHit{
					RuleID:   cr.rule.ID,
					RuleName: cr.rule.Name,
					Action:   cr.rule.Action,
					Severity: cr.rule.Severity,
					Matched:  redactAndTruncate(loc),
				})
				sanitized = cr.re.ReplaceAllString(sanitized, sanitizedMarker)
			}
		} else {
			lowerText := strings.ToLower(sanitized)
			lowerPat := strings.ToLower(cr.rule.Pattern)
			if idx := strings.Index(lowerText, lowerPat); idx >= 0 {
				matched := sanitized[idx : idx+len(cr.rule.Pattern)]
				hits = append(hits, MatchHit{
					RuleID:   cr.rule.ID,
					RuleName: cr.rule.Name,
					Action:   cr.rule.Action,
					Severity: cr.rule.Severity,
					Matched:  redactAndTruncate(matched),
				})
				sanitized = replaceCaseInsensitive(sanitized, cr.rule.Pattern, sanitizedMarker)
			}
		}
	}
	return sanitized, hits
}

func replaceCaseInsensitive(src, pat, repl string) string {
	re, err := regexp.Compile("(?i)" + regexp.QuoteMeta(pat))
	if err != nil {
		return src
	}
	return re.ReplaceAllString(src, repl)
}

func redactAndTruncate(s string) string {
	redacted := secretRedactionPattern.ReplaceAllString(s, "[REDACTED_SECRET]")
	redacted = strings.TrimSpace(redacted)
	if len(redacted) <= maxSnippetLength {
		return redacted
	}
	return redacted[:maxSnippetLength-3] + "..."
}

// inspectAndMaybeSanitizeBody inspects textual prompt fields inside an inbound
// JSON request body (OpenAI Chat, Anthropic, OpenAI Responses, Gemini).
// When sanitize is true and matches occur, it returns the mutated JSON body.
func inspectAndMaybeSanitizeBody(body []byte, cfg Config, rules []compiledRule, sanitize bool) (outBody []byte, hits []MatchHit, snippet string) {
	var root map[string]any
	if err := json.Unmarshal(body, &root); err != nil {
		return body, nil, ""
	}

	var mutated bool

	var walk func(v any, key string) any
	walk = func(v any, key string) any {
		switch x := v.(type) {
		case string:
			if isPromptTextKey(key) {
				san, textHits := evaluateText(x, cfg, rules)
				if len(textHits) > 0 {
					hits = append(hits, textHits...)
					if snippet == "" {
						snippet = redactAndTruncate(x)
					}
					if sanitize && san != x {
						mutated = true
						return san
					}
				}
			}
			return x
		case []any:
			for i := range x {
				x[i] = walk(x[i], key)
			}
			return x
		case map[string]any:
			for k, val := range x {
				x[k] = walk(val, k)
			}
			return x
		default:
			return v
		}
	}

	for _, topKey := range []string{"messages", "input", "contents", "system", "systemInstruction", "prompt"} {
		if val, ok := root[topKey]; ok {
			root[topKey] = walk(val, topKey)
		}
	}

	if !mutated {
		return body, hits, snippet
	}
	out, err := json.Marshal(root)
	if err != nil {
		return body, hits, snippet
	}
	return out, hits, snippet
}

func isPromptTextKey(key string) bool {
	switch key {
	case "content", "text", "input", "system", "prompt":
		return true
	default:
		return false
	}
}
