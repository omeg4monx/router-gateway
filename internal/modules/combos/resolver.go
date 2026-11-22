package combos

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"sort"
	"strings"

	"github.com/omeg4monx/router-gateway/internal/contracts"
)

const maxResolveDepth = 16

type resolveSnapshot struct {
	activeCombos map[string]Combo
	aliases      map[string]string
	disabled     map[string]struct{}
	connections  []contracts.ProviderConnection
}

type targetKey struct {
	connectionID  string
	upstreamModel string
}

// Resolve implements contracts.ProviderResolver with deterministic 4-tier
// precedence:
//  1. Active Combo (combos.name == model && enabled)
//  2. Model Alias (model_aliases.alias == model)
//  3. Explicit <provider>/<model>
//  4. Bare <model> ordered by (Priority ASC, ID ASC)
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

	visited := make(map[string]bool)
	rawCandidates := resolveIdentifier(requested, requested, snap, visited, 0)
	if len(rawCandidates) == 0 {
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

	now := m.now().UTC()
	available := make([]contracts.ResolvedTarget, 0, len(deduped))
	for _, cand := range deduped {
		if m.cooldown.IsInCooldown(cand.ConnectionID, now) {
			continue
		}
		available = append(available, cand)
	}
	if len(available) == 0 {
		return nil, contracts.ErrAllAccountsCoolingDown
	}
	return available, nil
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
			bareModel := strings.TrimPrefix(modelID, conn.ProviderID+"/")
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
		activeCombos: activeCombos,
		aliases:      aliases,
		disabled:     disabled,
		connections:  enabledConns,
	}, nil
}

func resolveIdentifier(
	ident string,
	logicalModel string,
	snap *resolveSnapshot,
	visited map[string]bool,
	depth int,
) []contracts.ResolvedTarget {
	clean := strings.TrimSpace(ident)
	if clean == "" || depth >= maxResolveDepth || visited[clean] {
		return nil
	}
	visited[clean] = true
	defer delete(visited, clean)

	// Tier 1: Active Combo.
	if combo, ok := snap.activeCombos[clean]; ok {
		var out []contracts.ResolvedTarget
		for _, target := range combo.Targets {
			resolved := resolveIdentifier(target, logicalModel, snap, visited, depth+1)
			out = append(out, resolved...)
		}
		return out
	}

	// Tier 2: Model Alias.
	if aliasTarget, ok := snap.aliases[clean]; ok {
		return resolveIdentifier(aliasTarget, logicalModel, snap, visited, depth+1)
	}

	// Tier 3: Explicit <provider>/<model>.
	if providerID, upstreamModel, hasSlash := strings.Cut(clean, "/"); hasSlash {
		if providerID == "" || upstreamModel == "" {
			return nil
		}
		var out []contracts.ResolvedTarget
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
			out = append(out, buildResolvedTarget(conn, logicalModel, upstreamModel))
		}
		return out
	}

	// Tier 4: Bare <model> ordered by (Priority ASC, ID ASC).
	var out []contracts.ResolvedTarget
	for _, conn := range snap.connections {
		if !connectionSupportsBare(conn, clean) {
			continue
		}
		if isModelDisabled(snap.disabled, conn.ProviderID, clean) {
			continue
		}
		out = append(out, buildResolvedTarget(conn, logicalModel, clean))
	}
	return out
}

func connectionSupportsExplicit(conn contracts.ProviderConnection, upstreamModel, fullSpec string) bool {
	if len(conn.Models) == 0 {
		return true
	}
	return slices.Contains(conn.Models, "*") ||
		slices.Contains(conn.Models, upstreamModel) ||
		slices.Contains(conn.Models, fullSpec)
}

func connectionSupportsBare(conn contracts.ProviderConnection, bareModel string) bool {
	if len(conn.Models) == 0 {
		return false
	}
	return slices.Contains(conn.Models, "*") || slices.Contains(conn.Models, bareModel)
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
