package combos

import (
	"crypto/rand"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"
	"sync"

	"github.com/omeg4monx/router-gateway/internal/contracts"
)

const defaultComboStrategy = "priority"

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

// strategyState tracks lightweight in-memory routing counters for round-robin,
// least-used, p2c, and last-known-good-provider (lkgp) combo strategies.
type strategyState struct {
	mu          sync.Mutex
	rrCounters  map[string]uint64
	usageCounts map[string]uint64
	lkgpConn    map[string]string
}

func newStrategyState() *strategyState {
	return &strategyState{
		rrCounters:  make(map[string]uint64),
		usageCounts: make(map[string]uint64),
		lkgpConn:    make(map[string]string),
	}
}

func (s *strategyState) nextRR(comboName string) uint64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	val := s.rrCounters[comboName]
	s.rrCounters[comboName] = val + 1
	return val
}

func (s *strategyState) usageOf(connectionID string) uint64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.usageCounts[connectionID]
}

func (s *strategyState) recordSelection(connectionID string) {
	if connectionID == "" {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.usageCounts[connectionID]++
}

func (s *strategyState) recordLKGP(comboName, connectionID string) {
	if comboName == "" || connectionID == "" {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.lkgpConn[comboName] = connectionID
}

func (s *strategyState) getLKGP(comboName string) string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.lkgpConn[comboName]
}

func cryptoRandIntn(n int) int {
	if n <= 1 {
		return 0
	}
	var b [8]byte
	if _, err := rand.Read(b[:]); err != nil {
		return 0
	}
	return int(binary.LittleEndian.Uint64(b[:]) % uint64(n))
}

func cryptoRandFloat64() float64 {
	var b [8]byte
	if _, err := rand.Read(b[:]); err != nil {
		return 0
	}
	return float64(binary.LittleEndian.Uint64(b[:])>>11) / (1 << 53)
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
	if len(items) == 0 {
		return nil
	}
	if len(items) == 1 {
		if state != nil {
			state.recordSelection(items[0].target.ConnectionID)
		}
		return []contracts.ResolvedTarget{items[0].target}
	}

	strat := strings.ToLower(strings.TrimSpace(combo.Strategy))
	if strat == "" || strat == "ordered" {
		strat = defaultComboStrategy
	}
	if strat == "auto" {
		if rs, ok := combo.Config["routerStrategy"].(string); ok {
			switch strings.ToLower(strings.TrimSpace(rs)) {
			case "cost", "eco":
				strat = "cost-optimized"
			case "latency", "fast":
				strat = "p2c"
			case "lkgp":
				strat = "lkgp"
			case "score", "sla", "sla-aware":
				strat = "least-used"
			default:
				strat = "priority"
			}
		} else {
			strat = "priority"
		}
	}

	out := make([]contracts.ResolvedTarget, len(items))
	for i, it := range items {
		out[i] = it.target
	}

	switch strat {
	case "round-robin":
		if state != nil {
			offset := int(state.nextRR(combo.Name) % uint64(len(out)))
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
		r := cryptoRandFloat64() * totalWeight
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
			j := cryptoRandIntn(i + 1)
			out[i], out[j] = out[j], out[i]
		}

	case "p2c":
		a := cryptoRandIntn(len(out))
		b := cryptoRandIntn(len(out) - 1)
		if b >= a {
			b++
		}
		winner := a
		if state != nil && state.usageOf(out[b].ConnectionID) < state.usageOf(out[a].ConnectionID) {
			winner = b
		}
		if winner != 0 {
			out[0], out[winner] = out[winner], out[0]
		}

	case "least-used", "headroom", "reset-aware", "reset-window", "expiry-first":
		if state != nil {
			sort.SliceStable(out, func(i, j int) bool {
				ui := state.usageOf(out[i].ConnectionID)
				uj := state.usageOf(out[j].ConnectionID)
				if ui != uj {
					return ui < uj
				}
				return out[i].Priority < out[j].Priority
			})
		}

	case "lkgp":
		if state != nil {
			if goodConn := state.getLKGP(combo.Name); goodConn != "" {
				for i := range out {
					if out[i].ConnectionID == goodConn {
						if i > 0 {
							chosen := out[i]
							copy(out[1:i+1], out[0:i])
							out[0] = chosen
						}
						break
					}
				}
			}
		}

	case "cost-optimized":
		sort.SliceStable(out, func(i, j int) bool {
			fi := inferCatalogModel(out[i].ProviderID, out[i].UpstreamModel, out[i].ProviderType, false).Free
			fj := inferCatalogModel(out[j].ProviderID, out[j].UpstreamModel, out[j].ProviderType, false).Free
			if fi != fj {
				return fi
			}
			return out[i].Priority < out[j].Priority
		})

	case "context-optimized":
		sort.SliceStable(out, func(i, j int) bool {
			ci := inferCatalogModel(out[i].ProviderID, out[i].UpstreamModel, out[i].ProviderType, false).ContextLength
			cj := inferCatalogModel(out[j].ProviderID, out[j].UpstreamModel, out[j].ProviderType, false).ContextLength
			if ci != cj {
				return ci > cj
			}
			return out[i].Priority < out[j].Priority
		})

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
