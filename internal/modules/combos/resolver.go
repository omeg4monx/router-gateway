package combos

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/omeg4monx/router-gateway/internal/contracts"
)

const maxResolveDepth = 16

type resolveSnapshot struct {
	activeCombos       map[string]Combo
	aliases            map[string]string
	disabled           map[string]struct{}
	providerStrategies map[string]string
	connections        []contracts.ProviderConnection
}

type targetKey struct {
	connectionID  string
	upstreamModel string
}

// Resolve implements contracts.ProviderResolver with deterministic 4-tier
// precedence:
//  1. Active Combo (combos.name == model && enabled)
//  2. Model Alias (model_aliases.alias == model)
//  3. Explicit <provider>/<model> (ordered by provider's active routing strategy)
//  4. Bare <model> (ordered by provider priority & each provider's active strategy)
func (m *Module) Resolve(ctx context.Context, model string, allowed []string) ([]contracts.ResolvedTarget, error) {
	requested := strings.TrimSpace(model)
	if requested == "" {
		return nil, contracts.ErrModelNotFound
	}
	if !isModelAllowed(requested, allowed) {
		return nil, contracts.ErrModelForbidden
	}

	snap, err := m.loadSnapshot(ctx)
	if err != nil {
		return nil, err
	}

	now := m.now().UTC()
	visited := make(map[string]bool)
	rawCandidates, matchedAny := m.resolveIdentifier(requested, requested, snap, now, visited, 0)
	if !matchedAny {
		return nil, contracts.ErrModelNotFound
	}

	// Deduplicate candidates while preserving first-seen fallback order.
	deduped := make([]contracts.ResolvedTarget, 0, len(rawCandidates))
	seen := make(map[targetKey]struct{}, len(rawCandidates))
	for _, cand := range rawCandidates {
		key := targetKey{
			connectionID:  cand.ConnectionID,
			upstreamModel: cand.UpstreamModel,
		}
		if _, dup := seen[key]; dup {
			continue
		}
		seen[key] = struct{}{}
		deduped = append(deduped, cand)
	}

	if len(deduped) == 0 {
		return nil, contracts.ErrAllAccountsCoolingDown
	}
	return deduped, nil
}

// ListModels implements contracts.ModelLister returning a deduplicated,
// allowed_models-filtered list of active combos, aliases, and enabled
// non-disabled provider models.
func (m *Module) ListModels(ctx context.Context, allowed []string) ([]contracts.ModelDescriptor, error) {
	snap, err := m.loadSnapshot(ctx)
	if err != nil {
		return nil, err
	}

	seen := make(map[string]struct{})
	out := make([]contracts.ModelDescriptor, 0)

	// 1. Active combos (Tier 1 precedence).
	comboNames := make([]string, 0, len(snap.activeCombos))
	for name := range snap.activeCombos {
		comboNames = append(comboNames, name)
	}
	sort.Strings(comboNames)
	for _, name := range comboNames {
		if !isModelAllowed(name, allowed) {
			continue
		}
		seen[name] = struct{}{}
		out = append(out, contracts.ModelDescriptor{
			ID:      name,
			OwnedBy: "combo",
			Kind:    "combo",
		})
	}

	// 2. Model aliases (Tier 2 precedence).
	aliasNames := make([]string, 0, len(snap.aliases))
	for alias := range snap.aliases {
		aliasNames = append(aliasNames, alias)
	}
	sort.Strings(aliasNames)
	for _, alias := range aliasNames {
		if _, dup := seen[alias]; dup {
			continue
		}
		if !isModelAllowed(alias, allowed) {
			continue
		}
		seen[alias] = struct{}{}
		out = append(out, contracts.ModelDescriptor{
			ID:      alias,
			OwnedBy: "alias",
			Kind:    "alias",
		})
	}

	// 3. Enabled provider models (qualified <provider>/<model> and bare <model>),
	// deduplicated by ID with highest-priority connection winning ownership.
	modelDescriptors := make([]contracts.ModelDescriptor, 0)
	for _, conn := range snap.connections {
		for _, raw := range conn.Models {
			modelID := strings.TrimSpace(raw)
			if modelID == "" || modelID == "*" {
				continue
			}
			bareModel := strings.TrimSpace(strings.TrimPrefix(modelID, conn.ProviderID+"/"))
			if bareModel == "" || bareModel == "*" {
				continue
			}
			if isModelDisabled(snap.disabled, conn.ProviderID, bareModel) {
				continue
			}

			qualifiedID := conn.ProviderID + "/" + bareModel
			if _, dup := seen[qualifiedID]; !dup && isModelAllowed(qualifiedID, allowed) {
				seen[qualifiedID] = struct{}{}
				modelDescriptors = append(modelDescriptors, contracts.ModelDescriptor{
					ID:      qualifiedID,
					OwnedBy: conn.ProviderID,
					Kind:    "model",
				})
			}

			if !strings.Contains(bareModel, "/") {
				if _, dup := seen[bareModel]; !dup && isModelAllowed(bareModel, allowed) {
					seen[bareModel] = struct{}{}
					modelDescriptors = append(modelDescriptors, contracts.ModelDescriptor{
						ID:      bareModel,
						OwnedBy: conn.ProviderID,
						Kind:    "model",
					})
				}
			}
		}
	}

	sort.Slice(modelDescriptors, func(i, j int) bool {
		return modelDescriptors[i].ID < modelDescriptors[j].ID
	})
	out = append(out, modelDescriptors...)
	return out, nil
}

func (m *Module) loadSnapshot(ctx context.Context) (*resolveSnapshot, error) {
	if m.providers == nil {
		return nil, errors.New("combos: provider store is not initialized")
	}

	comboList, err := m.repo.ListCombos(ctx)
	if err != nil {
		return nil, fmt.Errorf("combos: load combos: %w", err)
	}
	aliasList, err := m.repo.ListAliases(ctx)
	if err != nil {
		return nil, fmt.Errorf("combos: load aliases: %w", err)
	}
	disabledList, err := m.repo.ListDisabledModels(ctx)
	if err != nil {
		return nil, fmt.Errorf("combos: load disabled models: %w", err)
	}
	providerStratList, err := m.repo.ListProviderStrategies(ctx)
	if err != nil {
		return nil, fmt.Errorf("combos: load provider strategies: %w", err)
	}
	conns, err := m.providers.ListEnabledConnections(ctx)
	if err != nil {
		return nil, fmt.Errorf("combos: list enabled provider connections: %w", err)
	}

	activeCombos := make(map[string]Combo, len(comboList))
	for _, c := range comboList {
		if c.Enabled {
			activeCombos[c.Name] = c
		}
	}

	aliases := make(map[string]string, len(aliasList))
	for _, a := range aliasList {
		aliases[a.Alias] = a.TargetModel
	}

	disabled := make(map[string]struct{}, len(disabledList))
	for _, d := range disabledList {
		disabled[d.Model] = struct{}{}
	}

	providerStrategies := make(map[string]string, len(providerStratList))
	for _, ps := range providerStratList {
		providerStrategies[ps.ProviderID] = ps.Strategy
	}

	enabledConns := make([]contracts.ProviderConnection, 0, len(conns))
	for _, c := range conns {
		if c.Enabled {
			enabledConns = append(enabledConns, c)
		}
	}
	sort.Slice(enabledConns, func(i, j int) bool {
		if enabledConns[i].Priority != enabledConns[j].Priority {
			return enabledConns[i].Priority < enabledConns[j].Priority
		}
		return enabledConns[i].ID < enabledConns[j].ID
	})

	return &resolveSnapshot{
		activeCombos:       activeCombos,
		aliases:            aliases,
		disabled:           disabled,
		providerStrategies: providerStrategies,
		connections:        enabledConns,
	}, nil
}

func (m *Module) resolveIdentifier(
	ident string,
	logicalModel string,
	snap *resolveSnapshot,
	now time.Time,
	visited map[string]bool,
	depth int,
) ([]contracts.ResolvedTarget, bool) {
	clean := strings.TrimSpace(ident)
	if clean == "" || depth >= maxResolveDepth || visited[clean] {
		return nil, false
	}
	visited[clean] = true
	defer delete(visited, clean)

	// Tier 1: Active Combo (with full OmniRoute strategy ordering and structured steps).
	if combo, ok := snap.activeCombos[clean]; ok {
		steps := combo.Steps
		if len(steps) == 0 {
			steps = stepsFromTargets(combo.Targets)
		}
		weighted := make([]weightedCandidate, 0, len(steps))
		matchedAny := false
		for _, step := range steps {
			targetID := step.targetIdentifier()
			if targetID == "" {
				continue
			}
			resolved, stepMatched := m.resolveIdentifier(targetID, logicalModel, snap, now, visited, depth+1)
			if stepMatched {
				matchedAny = true
			}
			for _, cand := range resolved {
				if step.ConnectionID != "" && cand.ConnectionID != step.ConnectionID {
					continue
				}
				weighted = append(weighted, weightedCandidate{
					target: cand,
					weight: step.Weight,
				})
			}
		}
		if len(weighted) == 0 {
			return nil, matchedAny
		}
		return applyComboStrategy(combo, weighted, m.stratState), true
	}

	// Tier 2: Model Alias.
	if aliasTarget, ok := snap.aliases[clean]; ok {
		return m.resolveIdentifier(aliasTarget, logicalModel, snap, now, visited, depth+1)
	}

	// Tier 3: Explicit <provider>/<model> using the provider's active strategy.
	if providerID, upstreamModel, hasSlash := strings.Cut(clean, "/"); hasSlash {
		if providerID == "" || upstreamModel == "" {
			return nil, false
		}
		var available []contracts.ResolvedTarget
		matchedAny := false
		for _, conn := range snap.connections {
			if conn.ProviderID != providerID {
				continue
			}
			if !connectionSupportsExplicit(conn, upstreamModel, clean) {
				continue
			}
			if isModelDisabled(snap.disabled, conn.ProviderID, upstreamModel) {
				continue
			}
			matchedAny = true
			if m.cooldown.IsInCooldown(conn.ID, now) {
				continue
			}
			available = append(available, buildResolvedTarget(conn, logicalModel, upstreamModel))
		}
		if len(available) == 0 {
			return nil, matchedAny
		}
		strategy := snap.providerStrategies[providerID]
		return applyProviderStrategy(providerID, strategy, available, m.stratState), true
	}

	// Tier 4: Bare <model> ordered by provider priority, applying each
	// provider's active strategy across its own matching connections.
	var available []contracts.ResolvedTarget
	matchedAny := false
	for _, conn := range snap.connections {
		if !connectionSupportsBare(conn, clean) {
			continue
		}
		if isModelDisabled(snap.disabled, conn.ProviderID, clean) {
			continue
		}
		matchedAny = true
		if m.cooldown.IsInCooldown(conn.ID, now) {
			continue
		}
		available = append(available, buildResolvedTarget(conn, logicalModel, clean))
	}
	if len(available) <= 1 {
		if len(available) == 1 && m.stratState != nil {
			m.stratState.recordSelection(available[0].ConnectionID)
		}
		return available, matchedAny
	}

	// Group available candidates by provider while preserving first-seen provider
	// priority order, then apply each provider's active strategy to its pool.
	providerOrder := make([]string, 0, 4)
	byProvider := make(map[string][]contracts.ResolvedTarget, 4)
	for _, cand := range available {
		if _, exists := byProvider[cand.ProviderID]; !exists {
			providerOrder = append(providerOrder, cand.ProviderID)
		}
		byProvider[cand.ProviderID] = append(byProvider[cand.ProviderID], cand)
	}
	ordered := make([]contracts.ResolvedTarget, 0, len(available))
	for _, providerID := range providerOrder {
		group := byProvider[providerID]
		strategy := snap.providerStrategies[providerID]
		ordered = append(ordered, applyProviderStrategy(providerID, strategy, group, m.stratState)...)
	}
	return ordered, true
}

func connectionSupportsExplicit(conn contracts.ProviderConnection, upstreamModel, fullSpec string) bool {
	if len(conn.Models) == 0 {
		return true
	}
	providerWildcard := conn.ProviderID + "/*"
	for _, raw := range conn.Models {
		entry := strings.TrimSpace(raw)
		if entry == "*" || entry == providerWildcard || entry == upstreamModel || entry == fullSpec {
			return true
		}
	}
	return false
}

func connectionSupportsBare(conn contracts.ProviderConnection, bareModel string) bool {
	if len(conn.Models) == 0 {
		return false
	}
	qualified := conn.ProviderID + "/" + bareModel
	for _, raw := range conn.Models {
		entry := strings.TrimSpace(raw)
		if entry == "*" || entry == bareModel || entry == qualified {
			return true
		}
	}
	return false
}

func isModelDisabled(disabled map[string]struct{}, providerID, upstreamModel string) bool {
	if len(disabled) == 0 {
		return false
	}
	if _, ok := disabled[upstreamModel]; ok {
		return true
	}
	if _, ok := disabled[providerID+"/"+upstreamModel]; ok {
		return true
	}
	return false
}

func isModelAllowed(model string, allowed []string) bool {
	if len(allowed) == 0 {
		return true
	}
	for _, rawRule := range allowed {
		rule := strings.TrimSpace(rawRule)
		if rule == "" {
			continue
		}
		if rule == "*" || rule == model {
			return true
		}
		if prefix, ok := strings.CutSuffix(rule, "/*"); ok && prefix != "" {
			qualifiedPrefix := prefix + "/"
			if strings.HasPrefix(model, qualifiedPrefix) && len(model) > len(qualifiedPrefix) {
				return true
			}
		}
	}
	return false
}

func buildResolvedTarget(conn contracts.ProviderConnection, logicalModel, upstreamModel string) contracts.ResolvedTarget {
	return contracts.ResolvedTarget{
		ProviderID:    conn.ProviderID,
		ProviderType:  conn.Protocol,
		ConnectionID:  conn.ID,
		Model:         logicalModel,
		UpstreamModel: upstreamModel,
		Endpoint:      conn.BaseURL,
		Priority:      conn.Priority,
		Metadata: map[string]string{
			"auth_type":       conn.AuthType,
			"connection_name": conn.Name,
		},
	}
}
