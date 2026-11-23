package security

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

const (
	redactedPlaceholder = "[REDACTED]"
	maxSummaryRunes     = 500
)

var secretValuePatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?i)(bearer\s+)[A-Za-z0-9._~+/=-]{6,}`),
	regexp.MustCompile(`(?i)(basic\s+)[A-Za-z0-9+/=]{6,}`),
	regexp.MustCompile(`\bsk-[A-Za-z0-9_-]{8,}\b`),
	regexp.MustCompile(`\b(?:ghp|gho|ghu|ghs|ghr)_[A-Za-z0-9]{12,}\b`),
	regexp.MustCompile(`\bAIza[0-9A-Za-z_-]{16,}\b`),
	regexp.MustCompile(`\b(?:xoxb|xoxp|xoxa)-[A-Za-z0-9-]{10,}\b`),
	regexp.MustCompile(`(?i)((?:api[_-]?key|password|passwd|secret|access[_-]?token|refresh[_-]?token|token)["']?\s*[:=]\s*["']?)([^\s,;'"}\]]+)`),
}

var sensitiveKeyTokens = []string{
	"secret",
	"password",
	"passwd",
	"token",
	"api_key",
	"apikey",
	"authorization",
	"cookie",
	"credential",
	"private_key",
	"master_key",
}

// RedactText replaces recognizable credentials and secrets in text with
// [REDACTED] and truncates the output to maxSummaryRunes UTF-8 runes.
func RedactText(input string) string {
	out := strings.TrimSpace(input)
	if out == "" {
		return ""
	}
	for i, re := range secretValuePatterns {
		switch i {
		case 0, 1, 6:
			out = re.ReplaceAllString(out, "${1}"+redactedPlaceholder)
		default:
			out = re.ReplaceAllString(out, redactedPlaceholder)
		}
	}
	return truncateRunes(out, maxSummaryRunes)
}

// RedactMetadata recursively clones a metadata map and redacts any sensitive
// keys or secret-bearing string values.
func RedactMetadata(in map[string]any) map[string]any {
	if len(in) == 0 {
		return map[string]any{}
	}
	out := make(map[string]any, len(in))
	for k, v := range in {
		if isSensitiveKey(k) {
			out[k] = redactedPlaceholder
			continue
		}
		out[k] = redactValue(v)
	}
	return out
}

func redactValue(v any) any {
	switch tv := v.(type) {
	case string:
		return RedactText(tv)
	case map[string]any:
		return RedactMetadata(tv)
	case map[string]string:
		out := make(map[string]any, len(tv))
		for k, val := range tv {
			if isSensitiveKey(k) {
				out[k] = redactedPlaceholder
				continue
			}
			out[k] = RedactText(val)
		}
		return out
	case []string:
		out := make([]string, len(tv))
		for i, item := range tv {
			out[i] = RedactText(item)
		}
		return out
	case []any:
		out := make([]any, len(tv))
		for i, item := range tv {
			out[i] = redactValue(item)
		}
		return out
	default:
		return v
	}
}

func isSensitiveKey(key string) bool {
	lower := strings.ToLower(strings.TrimSpace(key))
	for _, token := range sensitiveKeyTokens {
		if strings.Contains(lower, token) {
			return true
		}
	}
	return false
}

func truncateRunes(s string, maxRunes int) string {
	if maxRunes <= 0 || utf8.RuneCountInString(s) <= maxRunes {
		return s
	}
	var b strings.Builder
	b.Grow(len(s))
	count := 0
	for _, r := range s {
		if count >= maxRunes {
			break
		}
		b.WriteRune(r)
		count++
	}
	return b.String()
}
