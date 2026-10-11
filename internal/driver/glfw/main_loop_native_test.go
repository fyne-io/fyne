//go:build accessibility && darwin && no_glfw

package glfw

import (
	"context"
	"encoding/json"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestMainLoopNativeLifecycle(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	for _, tags := range []string{"migrated_fynedo", "accessibility,migrated_fynedo"} {
		t.Run(tags, func(t *testing.T) {
			checkNativeLifecycle(ctx, t, tags)
		})
	}
}

func checkNativeLifecycle(ctx context.Context, t *testing.T, tags string) {
	t.Helper()
	binary := filepath.Join(t.TempDir(), "main-loop")
	//nolint:gosec // Fixed repository fixture, compiled into a test-owned temporary directory.
	build := exec.CommandContext(ctx, "go", "build", "-tags", tags,
		"-ldflags=-linkmode=external", "-o", binary, "./testdata/main_loop")
	output, err := build.CombinedOutput()
	require.NoError(t, err, "%s", output)
	for _, mode := range []string{"idle", "repaint"} {
		t.Run(mode, func(t *testing.T) {
			args := []string{"-duration=1s", "-exercise"}
			if mode == "repaint" {
				args = append(args, "-repaint")
			}
			//nolint:gosec // Executes the just-compiled fixture with fixed arguments.
			cmd := exec.CommandContext(ctx, binary, args...)
			output, err := cmd.CombinedOutput()
			require.NoError(t, err, "%s", output)
			var result struct {
				Frames        int64 `json:"frames"`
				Updates       int64 `json:"updates"`
				Elapsed       int64 `json:"elapsed_ns"`
				CPU           int64 `json:"cpu_ns"`
				Queued        int   `json:"queued"`
				SecondaryDone bool  `json:"secondary_closed"`
			}
			lines := strings.Split(strings.TrimSpace(string(output)), "\n")
			require.NoError(t, json.Unmarshal([]byte(lines[len(lines)-1]), &result), "%s", output)
			require.Equal(t, 100, result.Queued)
			require.True(t, result.SecondaryDone)
			require.GreaterOrEqual(t, result.Elapsed, int64(time.Second))
			if mode == "repaint" {
				require.InDelta(t, 60, float64(result.Frames)*float64(time.Second)/float64(result.Elapsed), 20)
				if tags == "accessibility,migrated_fynedo" {
					require.Greater(t, result.Updates, int64(20))
				}
			} else {
				require.Zero(t, result.Frames)
				require.Less(t, float64(result.CPU)/float64(result.Elapsed), 0.25, "idle native loop kept the CPU busy")
			}
			t.Logf("%s: frames=%d updates=%d cpu=%.1f%%", mode, result.Frames, result.Updates,
				100*float64(result.CPU)/float64(result.Elapsed))
		})
	}
}
