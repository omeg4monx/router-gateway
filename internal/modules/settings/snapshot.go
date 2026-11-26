package settings

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/omeg4monx/router-gateway/internal/contracts"
)

var _ contracts.SettingsSnapshotter = (*Module)(nil)

// ExportSettings returns a map of all known typed settings keys to their current values.
func (m *Module) ExportSettings(ctx context.Context) (map[string]any, error) {
	items, err := m.repo.List(ctx)
	if err != nil {
		return nil, err
	}
	out := make(map[string]any, len(items))
	for _, item := range items {
		out[item.Key] = item.Value
	}
	return out, nil
}

// ImportSettings applies known settings keys from values through the validated
// settings repository, ignoring unknown keys so snapshots from older or newer
// versions remain forward-compatible.
func (m *Module) ImportSettings(ctx context.Context, values map[string]any, actor, reason string) (int, error) {
	if len(values) == 0 {
		return 0, nil
	}
	existing, err := m.repo.List(ctx)
	if err != nil {
		return 0, err
	}
	known := make(map[string]struct{}, len(existing))
	for _, item := range existing {
		known[item.Key] = struct{}{}
	}

	rawUpdates := make(map[string]json.RawMessage)
	for k, v := range values {
		if _, ok := known[k]; !ok {
			continue
		}
		b, err := json.Marshal(v)
		if err != nil {
			continue
		}
		rawUpdates[k] = json.RawMessage(b)
	}
	if len(rawUpdates) == 0 {
		return 0, nil
	}

	actor = strings.TrimSpace(actor)
	if actor == "" {
		actor = defaultAdminSubject
	}
	reason = strings.TrimSpace(reason)
	if reason == "" {
		reason = "snapshot import"
	}

	res, err := m.repo.ApplyUpdates(ctx, rawUpdates, actor, reason)
	if err != nil {
		return 0, err
	}
	if m.bus != nil {
		m.bus.Publish(TopicSettingsUpdated, res.Event)
	}
	return len(rawUpdates), nil
}
