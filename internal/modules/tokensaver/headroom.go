package tokensaver

import (
	"container/list"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"strings"
	"sync"
	"time"
)

// WarmPrefixSummary exposes a single warm Headroom prefix entry for observability.
type WarmPrefixSummary struct {
	Key          string `json:"key"`
	Protocol     string `json:"protocol"`
	Model        string `json:"model"`
	ProviderID   string `json:"provider_id"`
	ConnectionID string `json:"connection_id"`
	Hits         int64  `json:"hits"`
	LastSeenAt   string `json:"last_seen_at"`
}

type cacheEntry struct {
	key          string
	protocol     string
	model        string
	providerID   string
	connectionID string
	hits         int64
	lastSeen     time.Time
}

// headroomCache is a bounded, thread-safe LRU prefix affinity cache.
type headroomCache struct {
	mu      sync.Mutex
	maxSize int
	ttl     time.Duration
	ll      *list.List
	entries map[string]*list.Element
}

func newHeadroomCache(maxSize int, ttl time.Duration) *headroomCache {
	if maxSize <= 0 {
		maxSize = 1024
	}
	if ttl <= 0 {
		ttl = 30 * time.Minute
	}
	return &headroomCache{
		maxSize: maxSize,
		ttl:     ttl,
		ll:      list.New(),
		entries: make(map[string]*list.Element),
	}
}

func (c *headroomCache) configure(maxSize int, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if maxSize > 0 {
		c.maxSize = maxSize
	}
	if ttl > 0 {
		c.ttl = ttl
	}
	c.evictLocked(time.Now().UTC())
}

func (c *headroomCache) lookup(key, protocol, model string, now time.Time) (hit bool, providerID, connectionID string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.evictExpiredLocked(now)
	if elem, ok := c.entries[key]; ok {
		ent, ok := elem.Value.(*cacheEntry)
		if !ok {
			return false, "", ""
		}
		ent.hits++
		ent.lastSeen = now
		c.ll.MoveToFront(elem)
		return true, ent.providerID, ent.connectionID
	}

	ent := &cacheEntry{
		key:      key,
		protocol: protocol,
		model:    model,
		lastSeen: now,
	}
	elem := c.ll.PushFront(ent)
	c.entries[key] = elem
	c.evictLocked(now)
	return false, "", ""
}

func (c *headroomCache) recordAffinity(key, providerID, connectionID string, now time.Time) {
	if key == "" {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if elem, ok := c.entries[key]; ok {
		ent, ok := elem.Value.(*cacheEntry)
		if !ok {
			return
		}
		ent.providerID = providerID
		ent.connectionID = connectionID
		ent.lastSeen = now
		c.ll.MoveToFront(elem)
	}
}

func (c *headroomCache) clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.ll.Init()
	c.entries = make(map[string]*list.Element)
}

func (c *headroomCache) snapshot(limit int) (active, maxEntries int, top []WarmPrefixSummary) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.evictExpiredLocked(time.Now().UTC())
	if limit <= 0 {
		limit = 10
	}
	top = make([]WarmPrefixSummary, 0, min(limit, c.ll.Len()))
	for e := c.ll.Front(); e != nil && len(top) < limit; e = e.Next() {
		ent, ok := e.Value.(*cacheEntry)
		if !ok {
			continue
		}
		top = append(top, WarmPrefixSummary{
			Key:          ent.key,
			Protocol:     ent.protocol,
			Model:        ent.model,
			ProviderID:   ent.providerID,
			ConnectionID: ent.connectionID,
			Hits:         ent.hits,
			LastSeenAt:   ent.lastSeen.UTC().Format(rfc3339NanoUTC),
		})
	}
	return c.ll.Len(), c.maxSize, top
}

func (c *headroomCache) evictExpiredLocked(now time.Time) {
	cutoff := now.Add(-c.ttl)
	for {
		back := c.ll.Back()
		if back == nil {
			break
		}
		ent, ok := back.Value.(*cacheEntry)
		if !ok {
			c.ll.Remove(back)
			continue
		}
		if ent.lastSeen.After(cutoff) {
			break
		}
		delete(c.entries, ent.key)
		c.ll.Remove(back)
	}
}

func (c *headroomCache) evictLocked(now time.Time) {
	c.evictExpiredLocked(now)
	for c.ll.Len() > c.maxSize {
		back := c.ll.Back()
		if back == nil {
			break
		}
		if ent, ok := back.Value.(*cacheEntry); ok {
			delete(c.entries, ent.key)
		}
		c.ll.Remove(back)
	}
}

// computePrefixFingerprint derives a stable short SHA-256 fingerprint of the
// normalized system prompt and initial conversation prefix.
func computePrefixFingerprint(protocol, model string, body []byte, prefixDepth int) string {
	if prefixDepth <= 0 {
		prefixDepth = 2
	}
	var root map[string]any
	if err := json.Unmarshal(body, &root); err != nil || root == nil {
		return ""
	}

	var b strings.Builder
	b.WriteString(protocol)
	b.WriteByte('|')
	b.WriteString(model)
	b.WriteByte('|')

	if sys, ok := root["system"]; ok {
		raw, _ := json.Marshal(sys)
		b.Write(raw)
		b.WriteByte('|')
	}
	if inst, ok := root["instructions"]; ok {
		raw, _ := json.Marshal(inst)
		b.Write(raw)
		b.WriteByte('|')
	}
	if si, ok := root["systemInstruction"]; ok {
		raw, _ := json.Marshal(si)
		b.Write(raw)
		b.WriteByte('|')
	}
	if msgs, ok := root["messages"].([]any); ok && len(msgs) > 0 {
		n := min(prefixDepth, len(msgs))
		raw, _ := json.Marshal(msgs[:n])
		b.Write(raw)
	} else if contents, ok := root["contents"].([]any); ok && len(contents) > 0 {
		n := min(prefixDepth, len(contents))
		raw, _ := json.Marshal(contents[:n])
		b.Write(raw)
	}

	sum := sha256.Sum256([]byte(b.String()))
	return hex.EncodeToString(sum[:12])
}
