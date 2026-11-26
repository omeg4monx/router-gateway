package usage

import (
	"maps"
	"sort"
	"strings"
	"sync"
)

const (
	pricingSourceBuiltin = "builtin"
	pricingSourceCustom  = "custom"
	tokensPerMillion     = 1_000_000.0
)

// ModelPricing describes per-million-token USD rates for a model.
type ModelPricing struct {
	Model            string  `json:"model"`
	InputPerMillion  float64 `json:"input_per_million"`
	OutputPerMillion float64 `json:"output_per_million"`
	CachedPerMillion float64 `json:"cached_per_million"`
	Source           string  `json:"source"`
	UpdatedAt        string  `json:"updated_at,omitempty"`
}

var defaultBuiltinPricing = []ModelPricing{
	{
		Model:            "gpt-4o",
		InputPerMillion:  2.50,
		OutputPerMillion: 10.00,
		CachedPerMillion: 1.25,
		Source:           pricingSourceBuiltin,
	},
	{
		Model:            "gpt-4o-mini",
		InputPerMillion:  0.15,
		OutputPerMillion: 0.60,
		CachedPerMillion: 0.075,
		Source:           pricingSourceBuiltin,
	},
	{
		Model:            "gpt-4.1",
		InputPerMillion:  2.00,
		OutputPerMillion: 8.00,
		CachedPerMillion: 0.50,
		Source:           pricingSourceBuiltin,
	},
	{
		Model:            "o3-mini",
		InputPerMillion:  1.10,
		OutputPerMillion: 4.40,
		CachedPerMillion: 0.55,
		Source:           pricingSourceBuiltin,
	},
	{
		Model:            "claude-sonnet-4",
		InputPerMillion:  3.00,
		OutputPerMillion: 15.00,
		CachedPerMillion: 0.30,
		Source:           pricingSourceBuiltin,
	},
	{
		Model:            "claude-3-7-sonnet",
		InputPerMillion:  3.00,
		OutputPerMillion: 15.00,
		CachedPerMillion: 0.30,
		Source:           pricingSourceBuiltin,
	},
	{
		Model:            "claude-3-5-haiku",
		InputPerMillion:  0.80,
		OutputPerMillion: 4.00,
		CachedPerMillion: 0.08,
		Source:           pricingSourceBuiltin,
	},
	{
		Model:            "gemini-2.5-pro",
		InputPerMillion:  1.25,
		OutputPerMillion: 10.00,
		CachedPerMillion: 0.3125,
		Source:           pricingSourceBuiltin,
	},
	{
		Model:            "gemini-2.5-flash",
		InputPerMillion:  0.30,
		OutputPerMillion: 2.50,
		CachedPerMillion: 0.075,
		Source:           pricingSourceBuiltin,
	},
	{
		Model:            "gemini-2.0-flash",
		InputPerMillion:  0.10,
		OutputPerMillion: 0.40,
		CachedPerMillion: 0.025,
		Source:           pricingSourceBuiltin,
	},
}

type pricingStore struct {
	mu        sync.RWMutex
	builtins  map[string]ModelPricing
	overrides map[string]ModelPricing
}

func newPricingStore() *pricingStore {
	builtins := make(map[string]ModelPricing, len(defaultBuiltinPricing))
	for _, entry := range defaultBuiltinPricing {
		builtins[entry.Model] = entry
	}
	return &pricingStore{
		builtins:  builtins,
		overrides: make(map[string]ModelPricing),
	}
}

func (s *pricingStore) SetOverrides(items []ModelPricing) {
	next := make(map[string]ModelPricing, len(items))
	for _, item := range items {
		item.Source = pricingSourceCustom
		next[item.Model] = item
	}
	s.mu.Lock()
	s.overrides = next
	s.mu.Unlock()
}

func (s *pricingStore) UpsertOverride(item ModelPricing) {
	item.Source = pricingSourceCustom
	s.mu.Lock()
	s.overrides[item.Model] = item
	s.mu.Unlock()
}

func (s *pricingStore) DeleteOverride(model string) {
	s.mu.Lock()
	delete(s.overrides, model)
	s.mu.Unlock()
}

func (s *pricingStore) Clear() {
	s.mu.Lock()
	s.overrides = make(map[string]ModelPricing)
	s.mu.Unlock()
}

func (s *pricingStore) lookup(model string) (ModelPricing, bool) {
	clean := strings.TrimSpace(model)
	if clean == "" {
		return ModelPricing{}, false
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	if entry, ok := s.lookupExactLocked(clean); ok {
		return entry, true
	}
	if _, bare, ok := strings.Cut(clean, "/"); ok && bare != "" {
		if entry, found := s.lookupExactLocked(bare); found {
			return entry, true
		}
	}
	return ModelPricing{}, false
}

func (s *pricingStore) lookupExactLocked(model string) (ModelPricing, bool) {
	if entry, ok := s.overrides[model]; ok {
		return entry, true
	}
	if entry, ok := s.builtins[model]; ok {
		return entry, true
	}
	return ModelPricing{}, false
}

// CalculateCost computes USD cost for token counts using custom overrides first
// and built-in model rates second, falling back to the bare model name after
// stripping any "<provider>/" prefix.
func (s *pricingStore) CalculateCost(model string, inputTokens, outputTokens, cachedTokens int64) float64 {
	rate, ok := s.lookup(model)
	if !ok {
		return 0
	}
	in := max(inputTokens, 0)
	out := max(outputTokens, 0)
	cached := max(cachedTokens, 0)

	total := float64(in)*rate.InputPerMillion +
		float64(out)*rate.OutputPerMillion +
		float64(cached)*rate.CachedPerMillion
	return total / tokensPerMillion
}

// List returns all built-in and custom model pricing entries merged by model
// name (custom overrides take precedence) and sorted alphabetically.
func (s *pricingStore) List() []ModelPricing {
	s.mu.RLock()
	merged := make(map[string]ModelPricing, len(s.builtins)+len(s.overrides))
	maps.Copy(merged, s.builtins)
	maps.Copy(merged, s.overrides)
	s.mu.RUnlock()

	out := make([]ModelPricing, 0, len(merged))
	for _, item := range merged {
		out = append(out, item)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].Model < out[j].Model
	})
	return out
}
