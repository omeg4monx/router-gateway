package combos

import (
	"context"
	"path/filepath"
	"sync"
	"testing"

	"github.com/omeg4monx/router-gateway/internal/contracts"
	"github.com/omeg4monx/router-gateway/internal/platform/db"
)

type staticProviderStore struct {
	mu    sync.RWMutex
	conns []contracts.ProviderConnection
}

func (s *staticProviderStore) ListEnabledConnections(context.Context) ([]contracts.ProviderConnection, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]contracts.ProviderConnection, len(s.conns))
	copy(out, s.conns)
	return out, nil
}

func (s *staticProviderStore) ResolveCredentials(context.Context, string) (contracts.UpstreamCredentials, error) {
	return contracts.UpstreamCredentials{}, nil
}

func (s *staticProviderStore) RefreshOAuthToken(context.Context, string) (contracts.UpstreamCredentials, error) {
	return contracts.UpstreamCredentials{}, nil
}

func TestCombosRoutingStateCacheHitsAndInvalidates(t *testing.T) {
	t.Parallel()

	database, err := db.Open(filepath.Join(t.TempDir(), "combos_cache.db"))
	if err != nil {
		t.Fatalf("db.Open: %v", err)
	}
	t.Cleanup(func() { _ = database.Close() })

	store := &staticProviderStore{
		conns: []contracts.ProviderConnection{
			{
				ID:         "conn-1",
				ProviderID: "openai",
				Protocol:   "openai-chat",
				BaseURL:    "https://api.openai.com/v1",
				Priority:   1,
				Enabled:    true,
				Models:     []string{"gpt-4o", "gpt-4o-mini"},
			},
		},
	}

	mod := New(Options{DB: database, Providers: store})
	ctx := context.Background()
	if err := database.MigrateModule(ctx, mod.Name(), mod.Migrations()); err != nil {
		t.Fatalf("MigrateModule: %v", err)
	}

	// First Resolve populates cache.
	if _, err := mod.Resolve(ctx, "openai/gpt-4o", nil); err != nil {
		t.Fatalf("initial Resolve: %v", err)
	}
	if !mod.routingCache.HasEntry() {
		t.Fatal("expected routingCache to be populated after Resolve")
	}

	// UpsertAlias must invalidate cache so new alias resolves immediately.
	if _, err := mod.repo.UpsertAlias(ctx, "fast", "openai/gpt-4o-mini"); err != nil {
		t.Fatalf("UpsertAlias: %v", err)
	}
	if mod.routingCache.HasEntry() {
		t.Fatal("expected routingCache to be invalidated after UpsertAlias")
	}

	targets, err := mod.Resolve(ctx, "fast", nil)
	if err != nil || len(targets) != 1 || targets[0].UpstreamModel != "gpt-4o-mini" {
		t.Fatalf("Resolve(fast) = (%+v, %v), want gpt-4o-mini", targets, err)
	}

	// Concurrent Resolve + mutation under -race.
	var wg sync.WaitGroup
	for range 16 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range 25 {
				_, _ = mod.Resolve(ctx, "openai/gpt-4o", nil)
				_, _ = mod.ListModels(ctx, nil)
			}
		}()
	}
	if _, err := mod.repo.UpsertProviderStrategy(ctx, "openai", "round-robin"); err != nil {
		t.Fatalf("UpsertProviderStrategy: %v", err)
	}
	wg.Wait()
}
