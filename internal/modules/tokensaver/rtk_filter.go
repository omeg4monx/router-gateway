package tokensaver

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	ansiCSIRe = regexp.MustCompile(`\x1b\[[0-9;?]*[\x20-\x2f]*[\x40-\x7e]`)
	ansiOSCRe = regexp.MustCompile(`\x1b\][^\x07\x1b]*(?:\x07|\x1b\\)`)
	grepRe    = regexp.MustCompile(`^([^:\s][^:]*?):(\d+):(.*)$`)
)

// filterToolOutput applies lean Go-native RTK compression to tool / tool_result
// outputs: ANSI & CR progress stripping, git status/diff compaction, grep file
// grouping, repeated log-line collapsing, and bounded UTF-8 safe truncation.
func filterToolOutput(input string, maxChars int) string {
	if input == "" {
		return ""
	}
	s := stripANSIAndProgress(input)

	switch {
	case isGitDiffOutput(s):
		s = compactGitDiff(s)
	case isGitStatusOutput(s):
		s = compactGitStatus(s)
	default:
		if grouped, ok := groupGrepOutput(s); ok {
			s = grouped
		}
	}

	s = collapseRepeatedLines(s)
	s = strings.TrimSpace(s)
	return truncateToolOutput(s, maxChars)
}

func stripANSIAndProgress(s string) string {
	if strings.ContainsRune(s, '\x1b') {
		s = ansiOSCRe.ReplaceAllString(s, "")
		s = ansiCSIRe.ReplaceAllString(s, "")
	}
	if !strings.ContainsRune(s, '\r') {
		return s
	}
	lines := strings.Split(s, "\n")
	for i, line := range lines {
		// Normalize CRLF first, then keep only the last carriage-return frame
		// so progress-bar overwrites ("10%\r50%\r100%") collapse to the final state.
		line = strings.TrimSuffix(line, "\r")
		if idx := strings.LastIndexByte(line, '\r'); idx >= 0 {
			line = line[idx+1:]
		}
		lines[i] = line
	}
	return strings.Join(lines, "\n")
}

func isGitStatusOutput(s string) bool {
	return strings.Contains(s, "On branch ") &&
		(strings.Contains(s, "Changes not staged for commit:") ||
			strings.Contains(s, "Changes to be committed:") ||
			strings.Contains(s, "Untracked files:") ||
			strings.Contains(s, "nothing to commit"))
}

func compactGitStatus(s string) string {
	lines := strings.Split(s, "\n")
	out := make([]string, 0, len(lines))
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, `(use "git `) ||
			strings.HasPrefix(trimmed, "no changes added to commit ") {
			continue
		}
		out = append(out, line)
	}
	return strings.Join(out, "\n")
}

func isGitDiffOutput(s string) bool {
	return strings.Contains(s, "diff --git ") && strings.Contains(s, "@@ ")
}

func compactGitDiff(s string) string {
	lines := strings.Split(s, "\n")
	out := make([]string, 0, len(lines))
	for _, line := range lines {
		if strings.HasPrefix(line, "diff --git ") ||
			strings.HasPrefix(line, "index ") ||
			strings.HasPrefix(line, "--- a/") {
			continue
		}
		out = append(out, line)
	}
	return strings.Join(out, "\n")
}

func groupGrepOutput(s string) (string, bool) {
	lines := strings.Split(s, "\n")
	if len(lines) < 4 {
		return s, false
	}

	type grepMatch struct {
		file string
		line string
		text string
	}
	matches := make([]grepMatch, 0, len(lines))
	nonEmpty := 0
	for _, raw := range lines {
		if strings.TrimSpace(raw) == "" {
			continue
		}
		nonEmpty++
		m := grepRe.FindStringSubmatch(raw)
		if m == nil {
			continue
		}
		matches = append(matches, grepMatch{file: m[1], line: m[2], text: m[3]})
	}
	if nonEmpty < 4 || len(matches)*10 < nonEmpty*8 {
		return s, false
	}

	var b strings.Builder
	lastFile := ""
	for _, m := range matches {
		if m.file != lastFile {
			if lastFile != "" {
				b.WriteByte('\n')
			}
			b.WriteString(m.file)
			b.WriteString(":\n")
			lastFile = m.file
		}
		b.WriteString("  ")
		b.WriteString(m.line)
		b.WriteByte(':')
		b.WriteString(strings.TrimRight(m.text, " \t"))
		b.WriteByte('\n')
	}
	grouped := strings.TrimRight(b.String(), "\n")
	if len(grouped) >= len(s) {
		return s, false
	}
	return grouped, true
}

func collapseRepeatedLines(s string) string {
	lines := strings.Split(s, "\n")
	if len(lines) < 3 {
		return s
	}
	out := make([]string, 0, len(lines))
	flush := func(prev string, count int) {
		if count <= 0 {
			return
		}
		if count >= 3 && strings.TrimSpace(prev) != "" {
			out = append(out, fmt.Sprintf("%s (x%d)", prev, count))
			return
		}
		for range count {
			out = append(out, prev)
		}
	}

	prev := lines[0]
	count := 1
	for _, line := range lines[1:] {
		if line == prev {
			count++
			continue
		}
		flush(prev, count)
		prev = line
		count = 1
	}
	flush(prev, count)
	return strings.Join(out, "\n")
}
