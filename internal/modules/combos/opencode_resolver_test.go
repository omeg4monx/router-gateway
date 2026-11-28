package combos_test

import (
	"context"
	"testing"

	"github.com/omeg4monx/router-gateway/internal/platform/opencode"
)

func TestResolverKeylessOpenCodeFallbackWithoutDBConnections(t *testing.T) {
	t.Parallel()

	h := newTestHarness(t, nil)

	for _, tc := range []struct {
		input        string
		wantUpstream string
	}{
		{"opencode/big-pickle", "big-pickle"},
		{"oc/deepseek-v4.1-flash-free", "deepseek-v4.1-flash-free"},
	} {
		targets, err := h.mod.Resolve(context.Background(), tc.input, nil)
		if err != nil {
			t.Fatalf("Resolve(%q) err = %v", tc.input, err)
		}
		if len(targets) != 1 {
			t.Fatalf("Resolve(%q) len = %d, want 1", tc.input, len(targets))
		}
		got := targets[0]
		if got.ConnectionID != opencode.VirtualConnectionID {
			t.Errorf("Resolve(%q) ConnectionID = %q, want %q", tc.input, got.ConnectionID, opencode.VirtualConnectionID)
		}
		if got.ProviderID != "opencode" {
			t.Errorf("Resolve(%q) ProviderID = %q, want opencode", tc.input, got.ProviderID)
		}
		if got.UpstreamModel != tc.wantUpstream {
			t.Errorf("Resolve(%q) UpstreamModel = %q, want %q", tc.input, got.UpstreamModel, tc.wantUpstream)
		}
		if got.Endpoint != opencode.DefaultFreeBaseURL {
			t.Errorf("Resolve(%q) Endpoint = %q, want %q", tc.input, got.Endpoint, opencode.DefaultFreeBaseURL)
		}
	}
}
