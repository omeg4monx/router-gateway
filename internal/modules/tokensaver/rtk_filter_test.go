package tokensaver

import (
	"strings"
	"testing"
)

func TestRTKFilterStripsANSIAndProgressOverwrites(t *testing.T) {
	raw := "\x1b[32mPASS\x1b[0m internal/auth\nDownloading 10%\rDownloading 50%\rDownloading 100%\n\x1b[1;31mFAIL\x1b[0m internal/db"
	got := filterToolOutput(raw, 8000)

	if strings.Contains(got, "\x1b[") {
		t.Fatalf("expected ANSI escape sequences stripped, got: %q", got)
	}
	if strings.Contains(got, "Downloading 10%") || !strings.Contains(got, "Downloading 100%") {
		t.Fatalf("expected CR progress overwrite collapsed to final frame, got: %q", got)
	}
	if !strings.Contains(got, "PASS internal/auth") || !strings.Contains(got, "FAIL internal/db") {
		t.Fatalf("expected test statuses preserved, got: %q", got)
	}
}

func TestRTKFilterDeduplicatesRepeatedLogLines(t *testing.T) {
	lines := make([]string, 0, 22)
	lines = append(lines, "Starting server...")
	for range 20 {
		lines = append(lines, "WARN retrying connection to upstream")
	}
	lines = append(lines, "Connected.")

	got := filterToolOutput(strings.Join(lines, "\n"), 8000)
	if !strings.Contains(got, "WARN retrying connection to upstream (x20)") {
		t.Fatalf("expected repeated log lines collapsed to (x20), got: %q", got)
	}
	if len(got) >= len(strings.Join(lines, "\n"))/2 {
		t.Fatalf("expected >50%% reduction on repeated logs, got len=%d", len(got))
	}
}

func TestRTKFilterCompactsGitStatusHints(t *testing.T) {
	raw := `On branch main
Your branch is up to date with 'origin/main'.

Changes not staged for commit:
  (use "git add <file>..." to update what will be committed)
  (use "git restore <file>..." to discard changes in working directory)
	modified:   internal/modules/tokensaver/compressor.go

Untracked files:
  (use "git add <file>..." to include in what will be committed)
	internal/modules/tokensaver/rtk_filter.go

no changes added to commit (use "git add" and/or "git commit -a")`

	got := filterToolOutput(raw, 8000)
	if strings.Contains(got, `use "git add`) || strings.Contains(got, `use "git restore`) {
		t.Fatalf("expected git status hint lines stripped, got: %q", got)
	}
	if !strings.Contains(got, "compressor.go") || !strings.Contains(got, "rtk_filter.go") {
		t.Fatalf("expected modified/untracked file paths preserved, got: %q", got)
	}
}

func TestRTKFilterCompactsGitDiffHeaders(t *testing.T) {
	raw := `diff --git a/internal/app/server.go b/internal/app/server.go
index 1111111..2222222 100644
--- a/internal/app/server.go
+++ b/internal/app/server.go
@@ -10,3 +10,4 @@ func Start() {
 	listen()
-	oldCall()
+	newCall()
 }`

	got := filterToolOutput(raw, 8000)
	if strings.Contains(got, "index 1111111..2222222") || strings.Contains(got, "diff --git ") {
		t.Fatalf("expected redundant git diff header lines stripped, got: %q", got)
	}
	if !strings.Contains(got, "+++ b/internal/app/server.go") || !strings.Contains(got, "-oldCall()") && !strings.Contains(got, "-\toldCall()") {
		t.Fatalf("expected diff file header and hunks preserved, got: %q", got)
	}
	if !strings.Contains(got, "+") || !strings.Contains(got, "newCall()") {
		t.Fatalf("expected added line preserved, got: %q", got)
	}
}

func TestRTKFilterGroupsGrepOutputByFile(t *testing.T) {
	raw := strings.Join([]string{
		"internal/modules/tokensaver/compressor.go:10:const StrategyOff = \"off\"",
		"internal/modules/tokensaver/compressor.go:11:const StrategyWhitespace = \"whitespace\"",
		"internal/modules/tokensaver/compressor.go:12:const StrategyDedupe = \"dedupe\"",
		"internal/modules/tokensaver/compressor.go:13:const StrategyAggressive = \"aggressive\"",
		"internal/modules/tokensaver/headroom.go:35:type headroomCache struct {",
		"internal/modules/tokensaver/headroom.go:43:func newHeadroomCache(",
	}, "\n")

	got := filterToolOutput(raw, 8000)
	if strings.Count(got, "internal/modules/tokensaver/compressor.go") != 1 {
		t.Fatalf("expected repeated grep file path grouped once, got:\n%s", got)
	}
	if !strings.Contains(got, "10:const StrategyOff") || !strings.Contains(got, "43:func newHeadroomCache") {
		t.Fatalf("expected line numbers and matches preserved, got:\n%s", got)
	}
	if len(got) >= len(raw) {
		t.Fatalf("expected grouped grep output shorter than raw (%d >= %d)", len(got), len(raw))
	}
}
