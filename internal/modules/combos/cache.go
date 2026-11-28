package combos

import "sync"

type cachedRoutingState struct {
	activeCombos       map[string]Combo
	aliases            map[string]string
	disabled           map[string]struct{}
	providerStrategies map[string]string
}

// routingStateCache caches the immutable DB-backed routing tables (combos,
// aliases, disabled models, provider strategies) using a monotonically
// increasing version counter so concurrent reads never repopulate stale state
// after a mutation.
type routingStateCache struct {
	mu      sync.RWMutex
	version uint64
	state   *cachedRoutingState
}

func newRoutingStateCache() *routingStateCache {
	return &routingStateCache{}
}

func (c *routingStateCache) Version() uint64 {
	c.mu.RLock()
	v := c.version
	c.mu.RUnlock()
	return v
}

func (c *routingStateCache) Get() (*cachedRoutingState, bool) {
	c.mu.RLock()
	st := c.state
	c.mu.RUnlock()
	return st, st != nil
}

func (c *routingStateCache) HasEntry() bool {
	c.mu.RLock()
	ok := c.state != nil
	c.mu.RUnlock()
	return ok
}

func (c *routingStateCache) SetIfVersion(state *cachedRoutingState, expectedVersion uint64) {
	c.mu.Lock()
	if c.version == expectedVersion {
		c.state = state
	}
	c.mu.Unlock()
}

func (c *routingStateCache) Invalidate() {
	c.mu.Lock()
	c.version++
	c.state = nil
	c.mu.Unlock()
}
