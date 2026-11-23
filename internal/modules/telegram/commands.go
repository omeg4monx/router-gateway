package telegram

import (
	"context"
	"fmt"
	"runtime"
	"strings"
	"time"
)

func (m *Module) executeCommand(ctx context.Context, raw string) (cmdName, reply string) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		trimmed = "/status"
	}
	fields := strings.Fields(trimmed)
	cmd := strings.ToLower(strings.TrimLeft(fields[0], "/!"))
	if atIdx := strings.IndexByte(cmd, '@'); atIdx > 0 {
		cmd = cmd[:atIdx]
	}

	switch cmd {
	case "start", "help":
		return cmd, strings.Join([]string{
			"OMEG4MON ROUTER Telegram Controller",
			"Available commands:",
			"- /status    : Live runtime health, uptime, memory, and module count",
			"- /providers : List active upstream provider connections",
			"- /models    : List available gateway models, aliases, and combos",
			"- /jobs      : Inspect background job status and run counts",
			"- /ping      : Heartbeat probe",
			"- /help      : Show this command list",
		}, "\n")

	case "ping":
		return "ping", fmt.Sprintf("Pong. OMEG4MON ROUTER is active (uptime %s).", formatDuration(time.Since(m.startedAt)))

	case "status":
		var mem runtime.MemStats
		runtime.ReadMemStats(&mem)
		memMB := float64(mem.Alloc) / (1024 * 1024)

		connCount := 0
		if m.providers != nil {
			if conns, err := m.providers.ListEnabledConnections(ctx); err == nil {
				connCount = len(conns)
			}
		}
		modelCount := 0
		if m.models != nil {
			if list, err := m.models.ListModels(ctx, nil); err == nil {
				modelCount = len(list)
			}
		}
		modCount := 0
		if m.activeModules != nil {
			modCount = len(m.activeModules())
		}
		return "status", fmt.Sprintf(
			"OMEG4MON ROUTER Status\n- Uptime: %s\n- Heap Alloc: %.1f MB\n- Active Modules: %d\n- Enabled Connections: %d\n- Routable Models: %d",
			formatDuration(time.Since(m.startedAt)),
			memMB,
			modCount,
			connCount,
			modelCount,
		)

	case "providers":
		if m.providers == nil {
			return "providers", "Provider inventory is unavailable."
		}
		conns, err := m.providers.ListEnabledConnections(ctx)
		if err != nil {
			return "providers", fmt.Sprintf("Failed to query providers: %v", err)
		}
		if len(conns) == 0 {
			return "providers", "No enabled provider connections configured."
		}
		var lines []string
		lines = append(lines, fmt.Sprintf("Enabled Provider Connections (%d):", len(conns)))
		for i, c := range conns {
			if i >= 15 {
				lines = append(lines, fmt.Sprintf("... and %d more", len(conns)-15))
				break
			}
			lines = append(lines, fmt.Sprintf("- %s [%s] (%d models)", c.Name, c.ProviderID, len(c.Models)))
		}
		return "providers", strings.Join(lines, "\n")

	case "models":
		if m.models == nil {
			return "models", "Model catalog is unavailable."
		}
		list, err := m.models.ListModels(ctx, nil)
		if err != nil {
			return "models", fmt.Sprintf("Failed to list models: %v", err)
		}
		if len(list) == 0 {
			return "models", "No models currently exposed."
		}
		var lines []string
		lines = append(lines, fmt.Sprintf("Routable Models (%d):", len(list)))
		for i, d := range list {
			if i >= 20 {
				lines = append(lines, fmt.Sprintf("... and %d more", len(list)-20))
				break
			}
			lines = append(lines, fmt.Sprintf("- %s (%s / %s)", d.ID, d.Kind, d.OwnedBy))
		}
		return "models", strings.Join(lines, "\n")

	case "jobs":
		if m.jobs == nil {
			return "jobs", "Background job inspector is unavailable."
		}
		snap := m.jobs.Snapshot()
		if len(snap) == 0 {
			return "jobs", "No background jobs registered."
		}
		var lines []string
		lines = append(lines, fmt.Sprintf("Background Jobs (%d):", len(snap)))
		for _, j := range snap {
			status := j.LastStatus
			if status == "" {
				status = "idle"
			}
			lines = append(lines, fmt.Sprintf("- %s: %s (runs=%d, errors=%d)", j.Name, status, j.RunCount, j.ErrorCount))
		}
		return "jobs", strings.Join(lines, "\n")

	default:
		return cmd, fmt.Sprintf("Unknown command /%s. Send /help to see supported commands.", cmd)
	}
}

func formatDuration(d time.Duration) string {
	if d < 0 {
		d = 0
	}
	total := int64(d.Seconds())
	days := total / 86400
	hours := (total % 86400) / 3600
	mins := (total % 3600) / 60
	secs := total % 60
	if days > 0 {
		return fmt.Sprintf("%dd %dh %dm %ds", days, hours, mins, secs)
	}
	if hours > 0 {
		return fmt.Sprintf("%dh %dm %ds", hours, mins, secs)
	}
	if mins > 0 {
		return fmt.Sprintf("%dm %ds", mins, secs)
	}
	return fmt.Sprintf("%ds", secs)
}
