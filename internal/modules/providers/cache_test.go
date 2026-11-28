package providers

import (
	"context"
	"path/filepath"
	"sync"
	"testing"

	"github.com/omeg4monx/router-gateway/internal/contracts"
	"github.com/omeg4monx/router-gateway/internal/platform/crypto"
	"github.com/omeg4monx/router-gateway/internal/platform/db"
)

func setupProvidersCacheTest(t *testing.T) *Module {
	t.Helper()

	database, err := db.Open(filepath.Join(t.TempDir(), "providers_cache.db"))
	if err != nil {
		t.Fatalf("db.Open: %v", err)
	}
	t.Cleanup(func() { _ = database.Close() })

	masterKey := make([]byte, 32)
	for i := range masterKey {
		masterKey[i] = byte(i + 1)
	}
	vault, err := crypto.NewVault(masterKey)
	if err != nil {
		t.Fatalf("crypto.NewVault: %v", err)
	}

	mod := New(Options{DB: database, Vault: vault})
	if err := database.MigrateModule(context.Background(), mod.Name(), mod.Migrations()); err != nil {
		t.Fatalf("MigrateModule: %v", err)
	}
	return mod
}

func TestProviderRuntimeCacheVersionedAndCloned(t *testing.T) {
	t.Parallel()

	cache := newProviderRuntimeCache()

	ver := cache.Version()
	creds := contracts.UpstreamCredentials{
		ConnectionID: "conn-1",
		ProviderID:   "openai",
		APIKey:       "sk-initial",
		ExtraHeaders: map[string]string{"X-Test": "1"},
	}
	cache.SetCredentialsIfVersion("conn-1", creds, ver)

	got, ok := cache.GetCredentials("conn-1")
	if !ok || got.APIKey != "sk-initial" {
		t.Fatalf("GetCredentials = (%+v, %v), want sk-initial", got, ok)
	}
	// Mutating caller's map must not mutate cached copy.
	got.ExtraHeaders["X-Test"] = "mutated"
	got2, _ := cache.GetCredentials("conn-1")
	if got2.ExtraHeaders["X-Test"] != "1" {
		t.Fatalf("cached ExtraHeaders mutated to %q, want 1", got2.ExtraHeaders["X-Test"])
	}

	// Stale SetIfVersion after Invalidate must be rejected.
	staleVer := cache.Version()
	cache.Invalidate()
	cache.SetCredentialsIfVersion("conn-1", creds, staleVer)
	if _, ok := cache.GetCredentials("conn-1"); ok {
		t.Fatal("stale SetCredentialsIfVersion repopulated invalidated cache")
	}
}

func TestProviderModuleCacheServesAndInvalidatesConcurrently(t *testing.T) {
	t.Parallel()

	mod := setupProvidersCacheTest(t)
	ctx := context.Background()

	enc1, err := mod.sealCredentials(storedCredentials{APIKey: "sk-v1"})
	if err != nil {
		t.Fatalf("sealCredentials: %v", err)
	}
	_, err = mod.repo.CreateConnection(ctx, createConnectionParams{
		ID:             "conn-cache-1",
		ProviderID:     "openai",
		Name:           "Primary",
		Protocol:       "openai-chat",
		BaseURL:        "https://api.openai.com/v1",
		AuthType:       "api_key",
		CredentialsEnc: enc1,
		Priority:       1,
		Enabled:        true,
		Models:         []string{"gpt-4o"},
		Headers:        map[string]string{"X-Env": "v1"},
	})
	if err != nil {
		t.Fatalf("CreateConnection: %v", err)
	}

	// Warm caches.
	conns, err := mod.ListEnabledConnections(ctx)
	if err != nil || len(conns) != 1 {
		t.Fatalf("ListEnabledConnections = (%d, %v), want 1", len(conns), err)
	}
	creds, err := mod.ResolveCredentials(ctx, "conn-cache-1")
	if err != nil || creds.APIKey != "sk-v1" {
		t.Fatalf("ResolveCredentials = (%q, %v), want sk-v1", creds.APIKey, err)
	}

	// Verify cache hit without DB: even under heavy concurrent reads + update,
	// there must be zero races and the updated key must be observed immediately.
	enc2, err := mod.sealCredentials(storedCredentials{APIKey: "sk-v2"})
	if err != nil {
		t.Fatalf("sealCredentials v2: %v", err)
	}

	var wg sync.WaitGroup
	for range 16 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range 25 {
				c, err := mod.ResolveCredentials(ctx, "conn-cache-1")
				if err != nil {
					t.Errorf("concurrent ResolveCredentials: %v", err)
					return
				}
				if c.APIKey != "sk-v1" && c.APIKey != "sk-v2" {
					t.Errorf("unexpected APIKey %q", c.APIKey)
					return
				}
				list, err := mod.ListEnabledConnections(ctx)
				if err != nil || len(list) != 1 {
					t.Errorf("concurrent ListEnabledConnections = (%d, %v)", len(list), err)
					return
				}
			}
		}()
	}

	_, err = mod.repo.UpdateConnection(ctx, "conn-cache-1", func(cur connectionRecord) (updateConnectionValues, error) {
		return updateConnectionValues{
			Name:           cur.Summary.Name,
			BaseURL:        cur.Summary.BaseURL,
			CredentialsEnc: enc2,
			Priority:       cur.Summary.Priority,
			Enabled:        true,
			Models:         cur.Summary.Models,
			Headers:        map[string]string{"X-Env": "v2"},
		}, nil
	})
	if err != nil {
		t.Fatalf("UpdateConnection: %v", err)
	}
	wg.Wait()

	after, err := mod.ResolveCredentials(ctx, "conn-cache-1")
	if err != nil || after.APIKey != "sk-v2" || after.ExtraHeaders["X-Env"] != "v2" {
		t.Fatalf("ResolveCredentials after update = (%+v, %v), want sk-v2 / X-Env=v2", after, err)
	}
}
