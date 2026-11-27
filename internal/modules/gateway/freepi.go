package gateway

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

const (
	providerIDFreePI        = "freepi"
	defaultFreePICliVersion = "0.6.28"
	defaultFreePIMinTokens  = 256
	maxFreePIMinTokens      = 4096
)

var (
	freepiVersionPattern = regexp.MustCompile(`(?i)below the required minimum \((\d+\.\d+\.\d+)\)`)

	freepiVersionMu  sync.RWMutex
	freepiCliVersion = resolveInitialFreePICliVersion()
)

func resolveInitialFreePICliVersion() string {
	if v := strings.TrimSpace(os.Getenv("FREEPI_CLI_VERSION")); v != "" {
		return v
	}
	return defaultFreePICliVersion
}

// GetFreePIClientVersion returns the active FreePI x-client-version value.
func GetFreePIClientVersion() string {
	freepiVersionMu.RLock()
	v := freepiCliVersion
	freepiVersionMu.RUnlock()
	if v == "" {
		return defaultFreePICliVersion
	}
	return v
}

// maybeLiftFreePIClientVersion updates the active FreePI client version when
// upstream reports a higher required minimum in a 426 response body.
func maybeLiftFreePIClientVersion(body []byte) (prev, next string, lifted bool) {
	required := parseFreePIRequiredVersion(body)
	if required == "" {
		return "", "", false
	}
	freepiVersionMu.Lock()
	defer freepiVersionMu.Unlock()
	cur := freepiCliVersion
	if cur == "" {
		cur = defaultFreePICliVersion
	}
	if compareSemver(required, cur) <= 0 {
		return cur, cur, false
	}
	freepiCliVersion = required
	return cur, required, true
}

// SetFreePIClientVersionForTest overrides the active FreePI client version in tests.
// Passing an empty string restores the environment/default version.
func SetFreePIClientVersionForTest(v string) {
	freepiVersionMu.Lock()
	if v == "" {
		freepiCliVersion = resolveInitialFreePICliVersion()
	} else {
		freepiCliVersion = v
	}
	freepiVersionMu.Unlock()
}

func parseFreePIRequiredVersion(body []byte) string {
	matches := freepiVersionPattern.FindSubmatch(body)
	if len(matches) < 2 {
		return ""
	}
	return strings.TrimSpace(string(matches[1]))
}

func compareSemver(a, b string) int {
	pa := strings.Split(strings.TrimSpace(a), ".")
	pb := strings.Split(strings.TrimSpace(b), ".")
	for i := range 3 {
		var va, vb int
		if i < len(pa) {
			va, _ = strconv.Atoi(pa[i])
		}
		if i < len(pb) {
			vb, _ = strconv.Atoi(pb[i])
		}
		if va != vb {
			if va > vb {
				return 1
			}
			return -1
		}
	}
	return 0
}

func resolveFreePIMinMaxTokens() int {
	raw, ok := os.LookupEnv("FREEPI_MIN_MAX_TOKENS")
	if !ok || strings.TrimSpace(raw) == "" {
		return defaultFreePIMinTokens
	}
	n, err := strconv.Atoi(strings.TrimSpace(raw))
	if err != nil || n < 0 {
		return 0
	}
	if n > maxFreePIMinTokens {
		return maxFreePIMinTokens
	}
	return n
}

// clampFreePIMaxTokens raises a too-small max_tokens / max_completion_tokens
// budget to minOut so reasoning-first FreePI models (such as particle) still
// produce visible content instead of exhausting the budget on reasoning_content.
func clampFreePIMaxTokens(body []byte, minOut int) (adjusted []byte, raisedFrom int) {
	if minOut <= 0 || len(body) == 0 {
		return body, 0
	}
	var payload map[string]any
	if err := json.Unmarshal(body, &payload); err != nil || payload == nil {
		return body, 0
	}

	key := ""
	if _, ok := payload["max_tokens"]; ok {
		key = "max_tokens"
	} else if _, ok := payload["max_completion_tokens"]; ok {
		key = "max_completion_tokens"
	}
	if key == "" {
		return body, 0
	}

	rawVal, ok := payload[key].(float64)
	if !ok {
		return body, 0
	}
	requested := int(rawVal)
	if requested <= 0 || requested >= minOut {
		return body, 0
	}

	payload[key] = minOut
	encoded, err := json.Marshal(payload)
	if err != nil {
		return body, 0
	}
	return encoded, requested
}

func applyFreePIHeaders(req *http.Request) {
	req.Header.Set("x-client-version", GetFreePIClientVersion())
	if strings.TrimSpace(req.Header.Get("x-session-id")) == "" {
		req.Header.Set("x-session-id", newRandomHexID(16))
	}
	req.Header.Set("x-turn-id", newRandomHexID(16))
	req.Header.Set("x-ad-context", "none")
}

func newRandomHexID(nBytes int) string {
	buf := make([]byte, nBytes)
	if _, err := rand.Read(buf); err != nil {
		return "00000000000000000000000000000000"
	}
	return hex.EncodeToString(buf)
}
