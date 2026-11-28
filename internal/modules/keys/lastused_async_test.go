package keys

import (
	"context"
	"path/filepath"
	"sync"
	"testing"

	"github.com/omeg4monx/router-gateway/internal/contracts"
	"github.com/omeg4monx/router-gateway/internal/platform/db"
)

func TestValidateKeyQueuesLastUsedAsynchronouslyAndFlushesOnStop(t *testing.T) {
	t.Parallel()

	database, err := db.Open(filepath.Join(t.TempDir(), "keys_async_touch.db"))
	if err != nil {
		t.Fatalf("db.Open: %v", err)
	}
	t.Cleanup(func() { _ = database.Close() })

	mod := New(Options{DB: database})
	ctx := context.Background()
	if err := database.MigrateModule(ctx, mod.Name(), mod.Migrations()); err != nil {
		t.Fatalf("MigrateModule: %v", err)
	}
	if err := mod.Start(ctx); err != nil {
		t.Fatalf("Start: %v", err)
	}

	prov, err := mod.ProvisionKey(ctx, contracts.ProvisionKeyRequest{Name: "async-touch"})
	if err != nil {
		t.Fatalf("ProvisionKey: %v", err)
	}

	// Verify touchCh is active when started.
	mod.lifeMu.Lock()
	hasTouchWorker := mod.touchCh != nil
	mod.lifeMu.Unlock()
	if !hasTouchWorker {
		t.Fatal("expected async touchCh worker to be active after Start")
	}

	var wg sync.WaitGroup
	for range 16 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if _, err := mod.ValidateKey(ctx, prov.Secret); err != nil {
				t.Errorf("ValidateKey: %v", err)
			}
		}()
	}
	wg.Wait()

	if err := mod.Stop(ctx); err != nil {
		t.Fatalf("Stop: %v", err)
	}

	summary, err := mod.repo.GetByID(ctx, prov.ID)
	if err != nil {
		t.Fatalf("GetByID: %v", err)
	}
	if summary.LastUsedAt == nil || *summary.LastUsedAt == "" {
		t.Fatal("expected last_used_at to be flushed to SQLite on Stop")
	}
}
