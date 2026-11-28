package providers

import (
	"slices"
	"sync"

	"github.com/omeg4monx/router-gateway/internal/contracts"
)

const maxCachedCredentialEntries = 2048

// providerRuntimeCache provides a versioned, read-through in-memory cache for
// hot-path ResolveCredentials and ListEnabledConnections calls. Every mutation
// bumps version and clears stale entries so concurrent readers cannot repopulate
// outdated state after a write.
type providerRuntimeCache struct {
	mu            sync.RWMutex
	version       uint64
	enabledLoaded bool
	enabledConns  []contracts.ProviderConnection
	creds         map[string]contracts.UpstreamCredentials
}

func newProviderRuntimeCache() *providerRuntimeCache {
	return &providerRuntimeCache{
		creds: make(map[string]contracts.UpstreamCredentials),
	}
}

func (c *providerRuntimeCache) Version() uint64 {
	c.mu.RLock()
	v := c.version
	c.mu.RUnlock()
	return v
}

func (c *providerRuntimeCache) GetEnabled() ([]contracts.ProviderConnection, bool) {
	c.mu.RLock()
	if !c.enabledLoaded {
		c.mu.RUnlock()
		return nil, false
	}
	out := cloneProviderConnections(c.enabledConns)
	c.mu.RUnlock()
	return out, true
}

func (c *providerRuntimeCache) SetEnabledIfVersion(conns []contracts.ProviderConnection, expectedVersion uint64) {
	cloned := cloneProviderConnections(conns)
	c.mu.Lock()
	if c.version == expectedVersion {
		c.enabledConns = cloned
		c.enabledLoaded = true
	}
	c.mu.Unlock()
}

func (c *providerRuntimeCache) GetCredentials(connectionID string) (contracts.UpstreamCredentials, bool) {
	c.mu.RLock()
	creds, ok := c.creds[connectionID]
	c.mu.RUnlock()
	if !ok {
		return contracts.UpstreamCredentials{}, false
	}
	return cloneUpstreamCredentials(creds), true
}

func (c *providerRuntimeCache) SetCredentialsIfVersion(connectionID string, creds contracts.UpstreamCredentials, expectedVersion uint64) {
	cloned := cloneUpstreamCredentials(creds)
	c.mu.Lock()
	if c.version == expectedVersion {
		if _, exists := c.creds[connectionID]; !exists && len(c.creds) >= maxCachedCredentialEntries {
			clear(c.creds)
		}
		c.creds[connectionID] = cloned
	}
	c.mu.Unlock()
}

func (c *providerRuntimeCache) Invalidate() {
	c.mu.Lock()
	c.version++
	c.enabledLoaded = false
	c.enabledConns = nil
	clear(c.creds)
	c.mu.Unlock()
}

func cloneProviderConnections(in []contracts.ProviderConnection) []contracts.ProviderConnection {
	if len(in) == 0 {
		return []contracts.ProviderConnection{}
	}
	out := make([]contracts.ProviderConnection, len(in))
	for i, conn := range in {
		conn.Models = slices.Clone(conn.Models)
		out[i] = conn
	}
	return out
}
