package combos

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/omeg4monx/router-gateway/internal/contracts"
	"github.com/omeg4monx/router-gateway/internal/platform/httpclient"
	"github.com/omeg4monx/router-gateway/internal/platform/venice"
)

func selectProbeHTTPClient(configured *http.Client, creds contracts.UpstreamCredentials) *http.Client {
	if venice.IsVeniceKeyless(creds) {
		return venice.SelectHTTPClient(configured)
	}
	if isOpenCodeFreeCreds(creds) {
		// Use Chrome/BoringSSL uTLS transport for real upstream requests, while
		// preserving custom mock RoundTrippers injected by unit tests.
		if configured == nil || configured.Transport == nil {
			return httpclient.SharedChromeTLSClient()
		}
		if _, isStdTransport := configured.Transport.(*http.Transport); isStdTransport {
			return httpclient.SharedChromeTLSClient()
		}
		return configured
	}
	if configured != nil {
		return configured
	}
	return &http.Client{Timeout: defaultProbeTimeout}
}

const (
	openCodeFreeUserAgent = "opencode/latest/2.4.1/cli"
	openCodeClientHeader  = "cli"
	openCodeBase62Chars   = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	openCodePrefixMask    = (uint64(1) << 48) - 1
)

func isOpenCodeFreeCreds(creds contracts.UpstreamCredentials) bool {
	pid := strings.ToLower(strings.TrimSpace(creds.ProviderID))
	if pid != "opencode" && pid != "oc" {
		return false
	}
	return strings.EqualFold(strings.TrimSpace(creds.AuthType), "none") ||
		(strings.TrimSpace(creds.APIKey) == "" && strings.TrimSpace(creds.AccessToken) == "")
}

func openCodeRandomHex(numBytes int) string {
	buf := make([]byte, numBytes)
	_, _ = rand.Read(buf)
	return hex.EncodeToString(buf)
}

func generateOpenCodeSessionID(now time.Time) string {
	ms := uint64(max(now.UnixMilli(), 0))
	prefix := (^(ms*0x1000 + 1)) & openCodePrefixMask
	buf := make([]byte, 14)
	_, _ = rand.Read(buf)
	var suffix strings.Builder
	suffix.Grow(14)
	for _, b := range buf {
		suffix.WriteByte(openCodeBase62Chars[int(b)%len(openCodeBase62Chars)])
	}
	return fmt.Sprintf("ses_%012x%s", prefix, suffix.String())
}

func openCodeEnvFingerprintPrompt() string {
	cwd, err := os.Getwd()
	if err != nil || strings.TrimSpace(cwd) == "" {
		cwd = "/tmp"
	}
	tmpDir := filepath.Join(os.TempDir(), "opencode")
	return strings.Join([]string{
		"Here is some useful information about the environment you are running in:",
		"<env>",
		"  Working directory: " + cwd,
		"  Workspace root folder: " + cwd,
		"  Is directory a git repo: no",
		"  Platform: " + runtime.GOOS,
		"  Prefer " + tmpDir + " over generic system temporary directories such as /tmp; it is pre-created and approved for external access.",
		"</env>",
	}, "\n")
}

func openCodeCanonicalTools() []map[string]any {
	return []map[string]any{
		{
			"type": "function",
			"function": map[string]any{
				"name":        "edit",
				"description": "Edit the contents of a file by finding and replacing exact text. When editing text from Read output, preserve the exact indentation (tabs or spaces) and omit the line-number prefix, such as `1: `. Never include the prefix in oldString or newString. Write line breaks as \\n; the file's existing CRLF or LF line endings are matched and preserved. To change line endings themselves, make oldString and newString differ only in \\r; that edit applies to the exact bytes. If there is no exact match, the edit retries while ignoring typographic quote, dash, and space variants, then trailing whitespace on each line. The edit fails if oldString is not found. By default, oldString must identify a UNIQUE location. Multiple matches FAIL unless replaceAll is true. Add more surrounding context to disambiguate, or set replaceAll to true to replace every occurrence. Use replaceAll when the change should apply to every occurrence, such as renaming a variable.",
				"parameters": map[string]any{
					"type": "object",
					"properties": map[string]any{
						"path":       map[string]any{"type": "string", "description": "File to edit"},
						"oldString":  map[string]any{"type": "string", "description": "Exact text to find and replace"},
						"newString":  map[string]any{"type": "string", "description": "Text to replace oldString with (must differ from oldString)"},
						"replaceAll": map[string]any{"type": "boolean", "description": "Whether to replace every occurrence of oldString. When false, oldString must match exactly once. Defaults to false."},
					},
					"required":             []string{"path", "oldString", "newString"},
					"additionalProperties": false,
				},
			},
		},
		{
			"type": "function",
			"function": map[string]any{
				"name":        "glob",
				"description": "Find files matching a glob pattern relative to a directory. Output is a list of file paths. Always use forward slashes (/) as path separators even on Windows. The pattern supports standard glob features: * matches characters within a directory, ** matches across directories, ? matches a single character, and [abc] matches a character set. Prefer specific patterns over broad ones.",
				"parameters": map[string]any{
					"type": "object",
					"properties": map[string]any{
						"pattern": map[string]any{"type": "string", "description": "Glob pattern to match files against"},
						"path":    map[string]any{"type": "string", "description": "Directory to search. Defaults to the current working directory."},
						"hidden":  map[string]any{"type": "boolean", "description": "Include hidden files and directories (default: false)."},
						"limit":   map[string]any{"type": "integer", "exclusiveMinimum": 0, "description": "Maximum number of matching files to return (default: 100)"},
					},
					"required":             []string{"pattern"},
					"additionalProperties": false,
				},
			},
		},
		{
			"type": "function",
			"function": map[string]any{
				"name":        "grep",
				"description": "Search for a regular expression pattern within files. Lines are returned with 1-based line numbers. By default, searches recursively from the working directory. Output is sorted by file path and line number.",
				"parameters": map[string]any{
					"type": "object",
					"properties": map[string]any{
						"pattern":       map[string]any{"type": "string", "minLength": 1, "description": "Regular expression or literal text to match in file contents."},
						"path":          map[string]any{"type": "string", "description": "File or directory to search. Defaults to the current working directory."},
						"include":       map[string]any{"type": "string", "description": "Glob pattern to filter files (for example, \"*.js\" or \"*.{ts,tsx}\")"},
						"literal":       map[string]any{"type": "boolean", "description": "Treat pattern as exact text instead of a regular expression (default: false)."},
						"caseSensitive": map[string]any{"type": "boolean", "description": "Use case-sensitive matching (default: true)."},
						"limit":         map[string]any{"type": "integer", "exclusiveMinimum": 0, "description": "Maximum number of matching lines to return (default: 100)"},
					},
					"required":             []string{"pattern"},
					"additionalProperties": false,
				},
			},
		},
		{
			"type": "function",
			"function": map[string]any{
				"name":        "read",
				"description": "Read the contents of a file or directory. Supports text files, images, and PDFs. Images and PDFs are presented directly to the model. Each text line is prefixed by its 1-based line number as <line>: <content>. The prefix is for reference and is not part of the file content. Directory entries are returned one per line. Use offset and limit to read large files or directories in sections. Prefer one larger read over many small slices, and use grep to find specific content in large files.",
				"parameters": map[string]any{
					"type": "object",
					"properties": map[string]any{
						"path":   map[string]any{"type": "string", "description": "File or directory to read"},
						"offset": map[string]any{"type": "integer", "minimum": 0, "description": "The line or directory entry to start reading from (1-based)"},
						"limit":  map[string]any{"type": "integer", "minimum": 0, "description": "The maximum number of lines or directory entries to read (defaults to and capped at 2000)"},
					},
					"required":             []string{"path"},
					"additionalProperties": false,
				},
			},
		},
		{
			"type": "function",
			"function": map[string]any{
				"name":        "shell",
				"description": "Execute a shell command and return its output. Commands run on Linux using bash. Quote file paths containing spaces or special characters. Prefer dedicated tools over shell commands when possible. When output is large, the full result is saved to a file and a truncated preview is returned. Rely on automatic truncation unless filtering the output is more useful. Commands accept an optional timeout, background commands have no timeout by default. Background commands return immediately, and you will be notified when they complete.",
				"parameters": map[string]any{
					"type": "object",
					"properties": map[string]any{
						"command":    map[string]any{"type": "string", "description": "Shell command string to execute"},
						"workdir":    map[string]any{"type": "string", "description": "Working directory to execute the command in. Defaults to the current working directory. When possible, avoid changing directories in the command and set the working directory here instead."},
						"timeout":    map[string]any{"type": "integer", "minimum": 0, "description": "Timeout in milliseconds. Set to 0 to disable the timeout. Defaults to 120000 for foreground commands. Background commands have no timeout by default."},
						"background": map[string]any{"type": "boolean", "description": "Run the command in the background and return immediately (useful for dev servers and long-running builds). You do not need to use '&' at the end of the command when using this parameter. You will be notified when it completes. DO NOT poll for completion."},
					},
					"required":             []string{"command"},
					"additionalProperties": false,
				},
			},
		},
	}
}

func marshalProbeJSON(payload any) ([]byte, error) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(payload); err != nil {
		return nil, err
	}
	return bytes.TrimSpace(buf.Bytes()), nil
}

func applyOpenCodeFreePayload(payload map[string]any, userPrompt string) {
	payload["stream"] = true
	payload["messages"] = []map[string]string{
		{"role": "system", "content": openCodeEnvFingerprintPrompt()},
		{"role": "user", "content": userPrompt},
	}
	payload["tools"] = openCodeCanonicalTools()
}

func applyOpenCodeFreeHeaders(req *http.Request) {
	sessionID := generateOpenCodeSessionID(time.Now())
	traceID := openCodeRandomHex(16)
	spanID := openCodeRandomHex(8)
	parentSpanID := openCodeRandomHex(8)

	req.Header.Set("Authorization", "Bearer public")
	req.Header.Set("User-Agent", openCodeFreeUserAgent)
	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("b3", traceID+"-"+spanID+"-1-"+parentSpanID)
	req.Header.Set("traceparent", "00-"+traceID+"-"+spanID+"-01")
	req.Header.Set("x-opencode-client", openCodeClientHeader)
	req.Header.Set("x-opencode-session", sessionID)
	req.Header.Set("x-opencode-session-id", sessionID)
	req.Header.Set("x-session-affinity", sessionID)
	req.Header.Set("x-session-id", sessionID)
	req.Header.Set("x-opencode-project", openCodeRandomHex(20))
}

func extractSSETokenUsage(protocol string, body []byte) (inputTokens, outputTokens, cachedTokens int64) {
	for _, line := range bytes.Split(body, []byte("\n")) {
		trimmed := bytes.TrimSpace(line)
		if !bytes.HasPrefix(trimmed, []byte("data:")) {
			continue
		}
		data := bytes.TrimSpace(bytes.TrimPrefix(trimmed, []byte("data:")))
		if len(data) == 0 || bytes.Equal(data, []byte("[DONE]")) || data[0] != '{' {
			continue
		}
		inTok, outTok, cachedTok := extractProbeTokenUsage(protocol, data)
		if inTok > 0 || outTok > 0 || cachedTok > 0 {
			inputTokens, outputTokens, cachedTokens = inTok, outTok, cachedTok
		}
	}
	return inputTokens, outputTokens, cachedTokens
}

func parseSSEIdentityPayload(body []byte) (reportedModel, idPrefix, fingerprint, selfReport string, hasReasoning bool) {
	var textBuf strings.Builder
	for _, line := range bytes.Split(body, []byte("\n")) {
		trimmed := bytes.TrimSpace(line)
		if !bytes.HasPrefix(trimmed, []byte("data:")) {
			continue
		}
		data := bytes.TrimSpace(bytes.TrimPrefix(trimmed, []byte("data:")))
		if len(data) == 0 || bytes.Equal(data, []byte("[DONE]")) || data[0] != '{' {
			continue
		}
		var chunk map[string]any
		if err := json.Unmarshal(data, &chunk); err != nil {
			continue
		}
		if reportedModel == "" {
			if m, ok := chunk["model"].(string); ok {
				reportedModel = strings.TrimSpace(m)
			}
		}
		if fingerprint == "" {
			if fp, ok := chunk["system_fingerprint"].(string); ok {
				fingerprint = strings.TrimSpace(fp)
			}
		}
		if idPrefix == "" {
			if id, ok := chunk["id"].(string); ok {
				idPrefix = extractIDPrefix(strings.TrimSpace(id))
			}
		}
		if choices, ok := chunk["choices"].([]any); ok && len(choices) > 0 {
			if first, ok := choices[0].(map[string]any); ok {
				if delta, ok := first["delta"].(map[string]any); ok {
					if rc, ok := delta["reasoning_content"].(string); ok && strings.TrimSpace(rc) != "" {
						hasReasoning = true
					}
					if c, ok := delta["content"].(string); ok && c != "" {
						textBuf.WriteString(c)
					}
				}
			}
		}
	}
	return reportedModel, idPrefix, fingerprint, strings.TrimSpace(textBuf.String()), hasReasoning
}
