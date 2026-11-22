package keys

import (
	"slices"
	"sync"
	"time"

	"github.com/omeg4monx/router-gateway/internal/contracts"
)

type cachedEntry struct {
	info      contracts.KeyInfo
	expiresAt *time.Time
}

type keyCache struct {
	mu      sync.RWMutex
	version uint64
	entries map[string]cachedEntry
}

func newKeyCache() *keyCache {
	return &keyCache{
		entries: make(map[string]cachedEntry),
	}
}

func (c *keyCache) Get(hash string) (cachedEntry, bool) {
	c.mu.RLock()
	entry, ok := c.entries[hash]
	c.mu.RUnlock()
	if !ok {
		return cachedEntry{}, false
	}
	return cloneCachedEntry(entry), true
}

func (c *keyCache) Version() uint64 {
	c.mu.RLock()
	v := c.version
	c.mu.RUnlock()
	return v
}

// SetIfVersion caches entry only if no invalidation occurred since expectedVersion,
// preventing a concurrent DB read from repopulating a freshly invalidated hash.
func (c *keyCache) SetIfVersion(hash string, entry cachedEntry, expectedVersion uint64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.version != expectedVersion {
		return
	}
	c.entries[hash] = cloneCachedEntry(entry)
}

func (c *keyCache) Delete(hashes ...string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.version++
	for _, h := range hashes {
		if h != "" {
			delete(c.entries, h)
		}
	}
}

func (c *keyCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.version++
	clear(c.entries)
}

func cloneCachedEntry(e cachedEntry) cachedEntry {
	out := cachedEntry{
		info: cloneKeyInfo(e.info),
	}
	if e.expiresAt != nil {
		t := *e.expiresAt
		out.expiresAt = &t
	}
	return out
}

func cloneKeyInfo(info contracts.KeyInfo) contracts.KeyInfo {
	out := info
	out.AllowedModels = cloneStringSlice(info.AllowedModels)
	out.Scopes = cloneStringSlice(info.Scopes)
	out.Tags = cloneStringSlice(info.Tags)
	return out
}

func cloneStringSlice(in []string) []string {
	if len(in) == 0 {
		return []string{}
	}
	return slices.Clone(in)
}
