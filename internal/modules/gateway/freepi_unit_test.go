package gateway

import (
	"bytes"
	"encoding/json"
	"io"
	"strings"
	"testing"
)

func TestFreePIClampMaxTokensFloor(t *testing.T) {
	t.Run("raises small max_tokens to default floor", func(t *testing.T) {
		body := []byte(`{"model":"particle","max_tokens":30,"messages":[{"role":"user","content":"hi"}]}`)
		out, raisedFrom := clampFreePIMaxTokens(body, 256)
		if raisedFrom != 30 {
			t.Fatalf("raisedFrom = %d, want 30", raisedFrom)
		}
		var parsed map[string]any
		if err := json.Unmarshal(out, &parsed); err != nil {
			t.Fatalf("unmarshal: %v", err)
		}
		if got, ok := parsed["max_tokens"].(float64); !ok || int(got) != 256 {
			t.Fatalf("max_tokens = %v, want 256", parsed["max_tokens"])
		}
	})

	t.Run("raises small max_completion_tokens to floor", func(t *testing.T) {
		body := []byte(`{"model":"particle","max_completion_tokens":60,"messages":[{"role":"user","content":"hi"}]}`)
		out, raisedFrom := clampFreePIMaxTokens(body, 256)
		if raisedFrom != 60 {
			t.Fatalf("raisedFrom = %d, want 60", raisedFrom)
		}
		var parsed map[string]any
		if err := json.Unmarshal(out, &parsed); err != nil {
			t.Fatalf("unmarshal: %v", err)
		}
		if got, ok := parsed["max_completion_tokens"].(float64); !ok || int(got) != 256 {
			t.Fatalf("max_completion_tokens = %v, want 256", parsed["max_completion_tokens"])
		}
	})

	t.Run("leaves adequate budget untouched", func(t *testing.T) {
		body := []byte(`{"model":"particle","max_tokens":4096}`)
		out, raisedFrom := clampFreePIMaxTokens(body, 256)
		if raisedFrom != 0 {
			t.Fatalf("raisedFrom = %d, want 0", raisedFrom)
		}
		if !bytes.Equal(out, body) {
			t.Fatalf("expected body untouched, got %s", string(out))
		}
	})

	t.Run("disabled when minOut is zero", func(t *testing.T) {
		body := []byte(`{"model":"particle","max_tokens":30}`)
		out, raisedFrom := clampFreePIMaxTokens(body, 0)
		if raisedFrom != 0 {
			t.Fatalf("raisedFrom = %d, want 0", raisedFrom)
		}
		if !bytes.Equal(out, body) {
			t.Fatalf("expected body untouched, got %s", string(out))
		}
	})
}

func TestFreePIParseRequiredVersionAndCompare(t *testing.T) {
	msg := `{"error":"upgrade_required","message":"your free-pi-cli 0.6.25 is below the required minimum (0.6.28); update with: npx free-pi-cli@latest"}`
	if got := parseFreePIRequiredVersion([]byte(msg)); got != "0.6.28" {
		t.Fatalf("parseFreePIRequiredVersion = %q, want 0.6.28", got)
	}
	if compareSemver("0.6.28", "0.6.25") <= 0 {
		t.Fatalf("expected 0.6.28 > 0.6.25")
	}
	if compareSemver("0.6.28", "0.6.28") != 0 {
		t.Fatalf("expected 0.6.28 == 0.6.28")
	}
	if compareSemver("0.6.25", "0.6.28") >= 0 {
		t.Fatalf("expected 0.6.25 < 0.6.28")
	}
}

func TestPeekFreePIPinMismatch(t *testing.T) {
	t.Run("detects policy sentence in SSE chunk", func(t *testing.T) {
		payload := "data: {\"choices\":[{\"delta\":{\"content\":\"This conversation belongs to another free-pi account. Sharing or rotating accounts isn't allowed.\"}}]}\n\n"
		body := io.NopCloser(strings.NewReader(payload))
		replayed, matched := peekFreePIPinMismatch(body, freePIPinPeekLimit)
		if !matched {
			t.Fatalf("matched = false, want true")
		}
		if replayed != nil {
			t.Fatalf("expected nil replay reader on match")
		}
	})

	t.Run("replays normal SSE stream byte-for-byte", func(t *testing.T) {
		parts := []string{
			`data: {"choices":[{"delta":{"content":"Here is a long enough answer to stop the peek immediately, as designed."}}]}` + "\n\n",
			`data: {"choices":[{"delta":{"content":" second part"}}]}` + "\n\n",
			"data: [DONE]\n\n",
		}
		want := strings.Join(parts, "")
		pr, pw := io.Pipe()
		go func() {
			for _, p := range parts {
				_, _ = pw.Write([]byte(p))
			}
			_ = pw.Close()
		}()

		replayed, matched := peekFreePIPinMismatch(pr, freePIPinPeekLimit)
		if matched {
			t.Fatalf("matched = true, want false")
		}
		defer replayed.Close()
		got, err := io.ReadAll(replayed)
		if err != nil {
			t.Fatalf("read replayed stream: %v", err)
		}
		if string(got) != want {
			t.Fatalf("replayed = %q, want %q", string(got), want)
		}
	})

	t.Run("detects policy sentence in non-streaming 200 JSON body", func(t *testing.T) {
		payload := `{"id":"chatcmpl-1","choices":[{"message":{"role":"assistant","content":"This conversation belongs to another free-pi account."}}]}`
		replayed, matched := peekFreePIPinMismatch(io.NopCloser(strings.NewReader(payload)), freePIPinPeekLimit)
		if !matched {
			t.Fatalf("matched = false, want true")
		}
		if replayed != nil {
			t.Fatalf("expected nil replay reader on match")
		}
	})
}
