package combos

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"strings"

	"github.com/omeg4monx/router-gateway/internal/contracts"
)

var _ contracts.ComboSnapshotter = (*Module)(nil)

// ExportCombos returns portable snapshots of all configured combos and model aliases.
func (m *Module) ExportCombos(ctx context.Context) ([]contracts.ComboSnapshot, []contracts.AliasSnapshot, error) {
	comboRows, err := m.repo.ListCombos(ctx)
	if err != nil {
		return nil, nil, err
	}
	aliasRows, err := m.repo.ListAliases(ctx)
	if err != nil {
		return nil, nil, err
	}

	combosOut := make([]contracts.ComboSnapshot, 0, len(comboRows))
	for _, c := range comboRows {
		targets := make([]string, len(c.Targets))
		copy(targets, c.Targets)
		combosOut = append(combosOut, contracts.ComboSnapshot{
			Name:        c.Name,
			Description: c.Description,
			Strategy:    c.Strategy,
			Targets:     targets,
			Enabled:     c.Enabled,
		})
	}

	aliasesOut := make([]contracts.AliasSnapshot, 0, len(aliasRows))
	for _, a := range aliasRows {
		aliasesOut = append(aliasesOut, contracts.AliasSnapshot{
			Alias:       a.Alias,
			TargetModel: a.TargetModel,
		})
	}
	return combosOut, aliasesOut, nil
}

// ImportCombos imports fallback combos and model aliases into the combos store.
// When overwrite is false, existing combos and aliases with the same name are preserved.
func (m *Module) ImportCombos(
	ctx context.Context,
	inCombos []contracts.ComboSnapshot,
	inAliases []contracts.AliasSnapshot,
	overwrite bool,
	actor string,
) (combosImported, aliasesImported int, err error) {
	actor = strings.TrimSpace(actor)
	if actor == "" {
		actor = defaultAdminSubject
	}

	existingCombos, err := m.repo.ListCombos(ctx)
	if err != nil {
		return 0, 0, err
	}
	comboByName := make(map[string]Combo, len(existingCombos))
	for _, c := range existingCombos {
		comboByName[strings.ToLower(strings.TrimSpace(c.Name))] = c
	}

	existingAliases, err := m.repo.ListAliases(ctx)
	if err != nil {
		return 0, 0, err
	}
	aliasSet := make(map[string]struct{}, len(existingAliases))
	for _, a := range existingAliases {
		aliasSet[strings.ToLower(strings.TrimSpace(a.Alias))] = struct{}{}
	}

	for _, item := range inCombos {
		name := strings.TrimSpace(item.Name)
		if name == "" || len(item.Targets) == 0 {
			continue
		}
		strategy := strings.TrimSpace(item.Strategy)
		if strategy == "" {
			strategy = strategyOrdered
		}
		targets := make([]string, 0, len(item.Targets))
		for _, t := range item.Targets {
			if trimmed := strings.TrimSpace(t); trimmed != "" {
				targets = append(targets, trimmed)
			}
		}
		if len(targets) == 0 {
			continue
		}

		norm := strings.ToLower(name)
		if cur, exists := comboByName[norm]; exists {
			if !overwrite {
				continue
			}
			updated, err := m.repo.UpdateCombo(ctx, cur.ID, func(_ Combo) (updateComboValues, error) {
				return updateComboValues{
					Name:        name,
					Description: strings.TrimSpace(item.Description),
					Targets:     targets,
					Strategy:    strategy,
					Enabled:     item.Enabled,
				}, nil
			})
			if err != nil {
				return combosImported, 0, err
			}
			comboByName[norm] = updated
			combosImported++
			continue
		}

		raw := make([]byte, idRandomBytes)
		if _, err := rand.Read(raw); err != nil {
			return combosImported, 0, err
		}
		created, err := m.repo.CreateCombo(ctx, createComboParams{
			ID:          defaultComboPrefix + hex.EncodeToString(raw),
			Name:        name,
			Description: strings.TrimSpace(item.Description),
			Targets:     targets,
			Strategy:    strategy,
			Enabled:     item.Enabled,
		})
		if err != nil {
			return combosImported, 0, err
		}
		comboByName[norm] = created
		combosImported++
	}

	for _, item := range inAliases {
		alias := strings.TrimSpace(item.Alias)
		target := strings.TrimSpace(item.TargetModel)
		if alias == "" || target == "" {
			continue
		}
		norm := strings.ToLower(alias)
		if _, exists := aliasSet[norm]; exists && !overwrite {
			continue
		}
		if _, err := m.repo.UpsertAlias(ctx, alias, target); err != nil {
			return combosImported, aliasesImported, err
		}
		aliasSet[norm] = struct{}{}
		aliasesImported++
	}

	if (combosImported > 0 || aliasesImported > 0) && m.bus != nil {
		m.bus.Publish(TopicCombosChanged, ChangedEvent{
			Action:       "import",
			ResourceType: "snapshot",
			ResourceID:   "batch",
			ChangedBy:    actor,
			UpdatedAt:    m.now().UTC().Format(rfc3339NanoFixed),
		})
	}
	return combosImported, aliasesImported, nil
}
