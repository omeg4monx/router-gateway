package combos

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand/v2"
	"sort"
	"strings"
	"sync"

	"github.com/omeg4monx/router-gateway/internal/contracts"
)

const (
	defaultComboStrategy = "priority"

	// DefaultProviderStrategy is the active default account-selection strategy
	// for every provider before any per-provider override is saved.
	DefaultProviderStrategy = "fill-first"

	maxStrategyStateEntries = 4096
)

// ProviderAccountStrategies lists the account-level routing strategies
// supported on the Provider Detail page for distributing requests across
// multiple connections belonging to a single provider.
var ProviderAccountStrategies = []string{
	"fill-first",
	"priority",
	"round-robin",
	"p2c",
	"least-used",
	"random",
	"strict-random",
	"lkgp",
	"cache-optimized",
	"headroom",
	"reset-aware",
	"reset-window",
	"expiry-first",
}

var validProviderAccountStrategies = func() map[string]struct{} {
	m := make(map[string]struct{}, len(ProviderAccountStrategies))
	for _, s := range ProviderAccountStrategies {
		m[s] = struct{}{}
	}
	return m
}()

// ProviderStrategyState describes the effective routing strategy for a single
// provider, including whether it is currently using the active default.
type ProviderStrategyState struct {
	ProviderID      string `json:"provider_id"`
	Strategy        string `json:"strategy"`
	DefaultStrategy string `json:"default_strategy"`
	IsDefault       bool   `json:"is_default"`
	UpdatedAt       string `json:"updated_at,omitempty"`
}

func validateProviderStrategy(raw string) (string, error) {
	trimmed := strings.ToLower(strings.TrimSpace(raw))
	if trimmed == "" {
		return DefaultProviderStrategy, nil
	}
	if _, ok := validProviderAccountStrategies[trimmed]; !ok {
		return "", fmt.Errorf("unsupported provider strategy %q", raw)
	}
	return trimmed, nil
}

// OmniRouteRoutingStrategies lists all 22 user-selectable OmniRoute v3.8.54 combo strategies.
var OmniRouteRoutingStrategies = []string{
	"priority",
	"weighted",
	"round-robin",
	"context-relay",
	"fill-first",
	"p2c",
	"random",
	"least-used",
	"cost-optimized",
	"reset-aware",
	"reset-window",
	"expiry-first",
	"headroom",
	"quota-weighted",
	"strict-random",
	"auto",
	"lkgp",
	"context-optimized",
	"cache-optimized",
	"fusion",
	"pipeline",
	"jev",
}

var validRoutingStrategies = func() map[string]struct{} {
	m := make(map[string]struct{}, len(OmniRouteRoutingStrategies))
	for _, s := range OmniRouteRoutingStrategies {
		m[s] = struct{}{}
	}
	return m
}()

// ComboStep represents a structured OmniRoute combo step (either a model target or a nested combo-ref).
type ComboStep struct {
	Kind                          string  `json:"kind"` // "model" | "combo-ref"
	Model                         string  `json:"model,omitempty"`
	ComboName                     string  `json:"comboName,omitempty"`
	Weight                        float64 `json:"weight,omitempty"`
	Label                         string  `json:"label,omitempty"`
	ConnectionID                  string  `json:"connectionId,omitempty"`
	Prompt                        string  `json:"prompt,omitempty"`
	FallbackOnlyOnQuotaExhaustion bool    `json:"fallbackOnlyOnQuotaExhaustion,omitempty"`
}

func (s ComboStep) targetIdentifier() string {
	if s.Kind == "combo-ref" {
		return strings.TrimSpace(s.ComboName)
	}
	return strings.TrimSpace(s.Model)
}

func validateStrategy(raw string) (string, error) {
	trimmed := strings.ToLower(strings.TrimSpace(raw))
	if trimmed == "" || trimmed == "ordered" {
		return defaultComboStrategy, nil
	}
	switch trimmed {
	case "usage":
		return "least-used", nil
	case "context":
		return "context-optimized", nil
	case "weekly-reset", "reset-window-order":
		return "reset-window", nil
	}
	if _, ok := validRoutingStrategies[trimmed]; !ok {
		return "", fmt.Errorf("unsupported strategy %q", raw)
	}
	return trimmed, nil
}

func normalizeAndValidateSteps(comboName string, rawTargets []string, rawSteps []ComboStep) ([]string, []ComboStep, error) {
	steps := make([]ComboStep, 0)
	if len(rawSteps) > 0 {
		if len(rawSteps) > maxTargetsCount {
			return nil, nil, fmt.Errorf("steps must not exceed %d items", maxTargetsCount)
		}
		seen := make(map[string]struct{}, len(rawSteps))
		for _, st := range rawSteps {
			kind := strings.ToLower(strings.TrimSpace(st.Kind))
			if kind == "" {
				if strings.TrimSpace(st.ComboName) != "" && strings.TrimSpace(st.Model) == "" {
					kind = "combo-ref"
				} else {
					kind = "model"
				}
			}
			if kind != "model" && kind != "combo-ref" {
				return nil, nil, fmt.Errorf("step kind must be 'model' or 'combo-ref', got %q", st.Kind)
			}
			if st.Weight < 0 || st.Weight > 100 {
				return nil, nil, errors.New("step weight must be between 0 and 100")
			}

			clean := ComboStep{
				Kind:                          kind,
				Weight:                        st.Weight,
				Label:                         strings.TrimSpace(st.Label),
				ConnectionID:                  strings.TrimSpace(st.ConnectionID),
				Prompt:                        strings.TrimSpace(st.Prompt),
				FallbackOnlyOnQuotaExhaustion: st.FallbackOnlyOnQuotaExhaustion,
			}

			var ident string
			var err error
			if kind == "combo-ref" {
				ident, err = validateModelIdentifier("steps.comboName", st.ComboName)
				if err != nil {
					return nil, nil, err
				}
				clean.ComboName = ident
			} else {
				ident, err = validateModelIdentifier("steps.model", st.Model)
				if err != nil {
					return nil, nil, err
				}
				clean.Model = ident
			}
			if ident == comboName {
				return nil, nil, errors.New("targets must not reference the combo itself")
			}
			dedupKey := kind + ":" + ident + ":" + clean.ConnectionID
			if _, dup := seen[dedupKey]; dup {
				return nil, nil, fmt.Errorf("steps contains duplicate entry %q", ident)
			}
			seen[dedupKey] = struct{}{}
			steps = append(steps, clean)
		}

		targets := make([]string, 0, len(steps))
		for _, st := range steps {
			targets = append(targets, st.targetIdentifier())
		}
		return targets, steps, nil
	}

	targets, err := validateComboTargets(comboName, rawTargets)
	if err != nil {
		return nil, nil, err
	}
	steps = stepsFromTargets(targets)
	return targets, steps, nil
}

func stepsFromTargets(targets []string) []ComboStep {
	out := make([]ComboStep, 0, len(targets))
	for _, t := range targets {
		clean := strings.TrimSpace(t)
		if clean == "" {
			continue
		}
		out = append(out, ComboStep{
			Kind:  "model",
			Model: clean,
		})
	}
	return out
}

func encodeStepsJSON(steps []ComboStep) (string, error) {
	if steps == nil {
		steps = []ComboStep{}
	}
	raw, err := json.Marshal(steps)
	if err != nil {
		return "", err
	}
	return string(raw), nil
}

func decodeStepsJSON(raw string, fallbackTargets []string) []ComboStep {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" || trimmed == "[]" || trimmed == "null" {
		return stepsFromTargets(fallbackTargets)
	}
	var out []ComboStep
	if err := json.Unmarshal([]byte(trimmed), &out); err != nil || len(out) == 0 {
		return stepsFromTargets(fallbackTargets)
	}
	return out
}

func encodeConfigJSON(cfg map[string]any) (string, error) {
	if len(cfg) == 0 {
		return "{}", nil
	}
	raw, err := json.Marshal(cfg)
	if err != nil {
		return "", err
	}
	return string(raw), nil
}

func decodeConfigJSON(raw string) map[string]any {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" || trimmed == "{}" || trimmed == "null" {
		return map[string]any{}
	}
	var out map[string]any
	if err := json.Unmarshal([]byte(trimmed), &out); err != nil || out == nil {
		return map[string]any{}
	}
	return out
}

func cloneSteps(in []ComboStep) []ComboStep {
	if len(in) == 0 {
		return []ComboStep{}
	}
	out := make([]ComboStep, len(in))
	copy(out, in)
	return out
}

func cloneConfig(in map[string]any) map[string]any {
	if len(in) == 0 {
		return map[string]any{}
	}
	raw, err := json.Marshal(in)
	if err != nil {
		return map[string]any{}
	}
	var out map[string]any
	if err := json.Unmarshal(raw, &out); err != nil || out == nil {
		return map[string]any{}
	}
	return out
}

// strategyState tracks lightweight, bounded in-memory routing counters for
// round-robin, least-used, p2c, and last-known-good-provider (lkgp) strategies.
type strategyState struct {
	mu          sync.Mutex
	rrCounters  map[string]uint64
	usageCounts map[string]uint64
	lastGoodSeq map[string]uint64
	goodSeq     uint64
}

func newStrategyState() *strategyState {
	return &strategyState{
		rrCounters:  make(map[string]uint64),
		usageCounts: make(map[string]uint64),
		lastGoodSeq: make(map[string]uint64),
	}
}

func (s *strategyState) nextRR(scopeKey string) uint64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.rrCounters[scopeKey]; !exists && len(s.rrCounters) >= maxStrategyStateEntries {
		clear(s.rrCounters)
	}
	val := s.rrCounters[scopeKey]
	s.rrCounters[scopeKey] = val + 1
	return val
}

func (s *strategyState) recordSelection(connectionID string) {
	if connectionID == "" {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.usageCounts[connectionID]; !exists && len(s.usageCounts) >= maxStrategyStateEntries {
		clear(s.usageCounts)
	}
	s.usageCounts[connectionID]++
}

func (s *strategyState) recordSuccess(connectionID string) {
	if connectionID == "" {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.lastGoodSeq[connectionID]; !exists && len(s.lastGoodSeq) >= maxStrategyStateEntries {
		clear(s.lastGoodSeq)
	}
	s.goodSeq++
	s.lastGoodSeq[connectionID] = s.goodSeq
}

func (s *strategyState) recordFailure(connectionID string) {
	if connectionID == "" {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.lastGoodSeq, connectionID)
}

func (s *strategyState) snapshotUsage(targets []contracts.ResolvedTarget) []uint64 {
	counts := make([]uint64, len(targets))
	s.mu.Lock()
	defer s.mu.Unlock()
	for i := range targets {
		counts[i] = s.usageCounts[targets[i].ConnectionID]
	}
	return counts
}

func (s *strategyState) usagePair(connA, connB string) (usageA, usageB uint64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.usageCounts[connA], s.usageCounts[connB]
}

func (s *strategyState) bestLKGPIndex(targets []contracts.ResolvedTarget) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	bestIdx := -1
	var bestSeq uint64
	for i := range targets {
		if seq := s.lastGoodSeq[targets[i].ConnectionID]; seq > bestSeq {
			bestSeq = seq
			bestIdx = i
		}
	}
	return bestIdx
}

func fastRandIntN(n int) int {
	if n <= 1 {
		return 0
	}
	return rand.IntN(n) //nolint:gosec // G404: non-cryptographic PRNG is intentional for high-throughput load-balancing selection
}

func fastRandFloat64() float64 {
	return rand.Float64() //nolint:gosec // G404: non-cryptographic PRNG is intentional for high-throughput load-balancing selection
}

type weightedCandidate struct {
	target contracts.ResolvedTarget
	weight float64
}

func applyComboStrategy(
	combo Combo,
	items []weightedCandidate,
	state *strategyState,
) []contracts.ResolvedTarget {
	return applyStrategy("combo:"+combo.Name, combo.Strategy, combo.Config, items, state)
}

func applyProviderStrategy(
	providerID string,
	strategy string,
	targets []contracts.ResolvedTarget,
	state *strategyState,
) []contracts.ResolvedTarget {
	if len(targets) == 0 {
		return nil
	}
	normalized := strings.ToLower(strings.TrimSpace(strategy))
	if normalized == "" {
		normalized = DefaultProviderStrategy
	}
	// Fast path: single candidate or order-preserving strategies require zero reordering allocations.
	if len(targets) == 1 || normalized == "fill-first" || normalized == "priority" {
		if state != nil {
			state.recordSelection(targets[0].ConnectionID)
		}
		return targets
	}
	items := make([]weightedCandidate, len(targets))
	for i, target := range targets {
		items[i] = weightedCandidate{target: target, weight: 1}
	}
	return applyStrategy("provider:"+providerID, normalized, nil, items, state)
}

func applyStrategy(
	scopeKey string,
	rawStrategy string,
	config map[string]any,
	items []weightedCandidate,
	state *strategyState,
) []contracts.ResolvedTarget {
	if len(items) == 0 {
		return nil
	}
	if len(items) == 1 {
		if state != nil {
			state.recordSelection(items[0].target.ConnectionID)
		}
		return []contracts.ResolvedTarget{items[0].target}
	}

	strategy := strings.ToLower(strings.TrimSpace(rawStrategy))
	if strategy == "" || strategy == "ordered" {
		strategy = defaultComboStrategy
	}
	if strategy == "auto" {
		if rs, ok := config["routerStrategy"].(string); ok {
			switch strings.ToLower(strings.TrimSpace(rs)) {
			case "cost", "eco":
				strategy = "cost-optimized"
			case "latency", "fast":
				strategy = "p2c"
			case "lkgp":
				strategy = "lkgp"
			case "score", "sla", "sla-aware":
				strategy = "least-used"
			default:
				strategy = "priority"
			}
		} else {
			strategy = "priority"
		}
	}

	out := make([]contracts.ResolvedTarget, len(items))
	for i, it := range items {
		out[i] = it.target
	}

	switch strategy {
	case "round-robin":
		if state != nil {
			offset := int(state.nextRR(scopeKey) % uint64(len(out))) //nolint:gosec // G115: modulo by len(out) guarantees result fits in int
			if offset > 0 {
				rotated := make([]contracts.ResolvedTarget, 0, len(out))
				rotated = append(rotated, out[offset:]...)
				rotated = append(rotated, out[:offset]...)
				out = rotated
			}
		}

	case "weighted", "quota-weighted":
		totalWeight := 0.0
		weights := make([]float64, len(items))
		for i, it := range items {
			w := it.weight
			if w <= 0 {
				w = 1
			}
			weights[i] = w
			totalWeight += w
		}
		pick := 0
		r := fastRandFloat64() * totalWeight
		acc := 0.0
		for i, w := range weights {
			acc += w
			if r <= acc {
				pick = i
				break
			}
		}
		ordered := make([]contracts.ResolvedTarget, 0, len(out))
		ordered = append(ordered, items[pick].target)
		rest := make([]weightedCandidate, 0, len(items)-1)
		for i, it := range items {
			if i != pick {
				rest = append(rest, it)
			}
		}
		sort.SliceStable(rest, func(i, j int) bool {
			return rest[i].weight > rest[j].weight
		})
		for _, it := range rest {
			ordered = append(ordered, it.target)
		}
		out = ordered

	case "random", "strict-random":
		for i := len(out) - 1; i > 0; i-- {
			j := fastRandIntN(i + 1)
			out[i], out[j] = out[j], out[i]
		}

	case "p2c":
		a := fastRandIntN(len(out))
		b := fastRandIntN(len(out) - 1)
		if b >= a {
			b++
		}
		winner := a
		if state != nil {
			ua, ub := state.usagePair(out[a].ConnectionID, out[b].ConnectionID)
			if ub < ua {
				winner = b
			}
		}
		if winner != 0 {
			out[0], out[winner] = out[winner], out[0]
		}

	case "least-used", "headroom", "reset-aware", "reset-window", "expiry-first":
		if state != nil {
			usage := state.snapshotUsage(out)
			type rankedTarget struct {
				target contracts.ResolvedTarget
				usage  uint64
			}
			ranked := make([]rankedTarget, len(out))
			for i := range out {
				ranked[i] = rankedTarget{target: out[i], usage: usage[i]}
			}
			sort.SliceStable(ranked, func(i, j int) bool {
				if ranked[i].usage != ranked[j].usage {
					return ranked[i].usage < ranked[j].usage
				}
				return ranked[i].target.Priority < ranked[j].target.Priority
			})
			for i := range ranked {
				out[i] = ranked[i].target
			}
		}

	case "lkgp":
		if state != nil {
			if bestIdx := state.bestLKGPIndex(out); bestIdx > 0 {
				chosen := out[bestIdx]
				copy(out[1:bestIdx+1], out[0:bestIdx])
				out[0] = chosen
			}
		}

	case "cost-optimized":
		type costRanked struct {
			target contracts.ResolvedTarget
			free   bool
		}
		ranked := make([]costRanked, len(out))
		for i := range out {
			ranked[i] = costRanked{
				target: out[i],
				free:   inferCatalogModel(out[i].ProviderID, out[i].UpstreamModel, out[i].ProviderType, false).Free,
			}
		}
		sort.SliceStable(ranked, func(i, j int) bool {
			if ranked[i].free != ranked[j].free {
				return ranked[i].free
			}
			return ranked[i].target.Priority < ranked[j].target.Priority
		})
		for i := range ranked {
			out[i] = ranked[i].target
		}

	case "context-optimized":
		type ctxRanked struct {
			target        contracts.ResolvedTarget
			contextLength int
		}
		ranked := make([]ctxRanked, len(out))
		for i := range out {
			ranked[i] = ctxRanked{
				target:        out[i],
				contextLength: inferCatalogModel(out[i].ProviderID, out[i].UpstreamModel, out[i].ProviderType, false).ContextLength,
			}
		}
		sort.SliceStable(ranked, func(i, j int) bool {
			if ranked[i].contextLength != ranked[j].contextLength {
				return ranked[i].contextLength > ranked[j].contextLength
			}
			return ranked[i].target.Priority < ranked[j].target.Priority
		})
		for i := range ranked {
			out[i] = ranked[i].target
		}

	case "cache-optimized":
		sort.SliceStable(out, func(i, j int) bool {
			if out[i].ProviderID != out[j].ProviderID {
				return out[i].ProviderID < out[j].ProviderID
			}
			if out[i].UpstreamModel != out[j].UpstreamModel {
				return out[i].UpstreamModel < out[j].UpstreamModel
			}
			return out[i].ConnectionID < out[j].ConnectionID
		})

	default:
		// priority, fill-first, context-relay, fusion, pipeline, jev preserve explicit step order.
	}

	if state != nil && len(out) > 0 {
		state.recordSelection(out[0].ConnectionID)
	}
	return out
}
