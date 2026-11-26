package combos_test

import (
	"context"
	"testing"

	"github.com/omeg4monx/router-gateway/internal/platform/venice"
)

func TestResolverKeylessVeniceFallbackWithoutDBConnections(t *testing.T) {
	t.Parallel()

	// Empty provider store: zero DB connections configured.
	h := newTestHarness(t, nil)

	for _, tc := range []struct {
		input        string
		wantUpstream string
	}{
		{"venice/venice-uncensored-1-2", "venice-uncensored-1-2"},
		{"venice-ai/deepseek-v4-1-flash-uncensored", "deepseek-v4-1-flash-uncensored"},
		{"veniceai/qwen3-coder-480b-a35b-instruct-turbo", "qwen3-coder-480b-a35b-instruct-turbo"},
	} {
		targets, err := h.mod.Resolve(context.Background(), tc.input, nil)
		if err != nil {
			t.Fatalf("Resolve(%q) err = %v", tc.input, err)
		}
		if len(targets) != 1 {
			t.Fatalf("Resolve(%q) len = %d, want 1", tc.input, len(targets))
		}
		got := targets[0]
		if got.ConnectionID != venice.VirtualConnectionID {
			t.Errorf("Resolve(%q) ConnectionID = %q, want %q", tc.input, got.ConnectionID, venice.VirtualConnectionID)
		}
		if got.ProviderID != "venice" {
			t.Errorf("Resolve(%q) ProviderID = %q, want venice", tc.input, got.ProviderID)
		}
		if got.UpstreamModel != tc.wantUpstream {
			t.Errorf("Resolve(%q) UpstreamModel = %q, want %q", tc.input, got.UpstreamModel, tc.wantUpstream)
		}
		if got.Endpoint != venice.DefaultInferenceURL {
			t.Errorf("Resolve(%q) Endpoint = %q, want %q", tc.input, got.Endpoint, venice.DefaultInferenceURL)
		}
	}
}
