// Command textbench drives a deterministic text-heavy scroll so that GL render
// statistics can be compared between builds. It is a measurement harness, not a
// demo app.
//
// Build with -tags glstats to enable the counting GL context:
//
//	go build -tags glstats ./cmd/textbench
//	FYNE_GL_STATS_OUT=before.json ./textbench
//
// The process writes its report and exits on its own once enough frames have
// been recorded (see internal/painter/gl/stats_glstats.go for the environment
// variables that control this).
//
// Every row holds a distinct string, which is the case the per-text-object
// texture cache handles worst: no two rows can share a cached texture, so each
// newly revealed row forces a fresh rasterisation and texture upload.
package main

import (
	"fmt"
	"image"
	"image/png"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// writePNG saves img to path. The path comes from the environment, so it is
// cleaned and required to be a plain .png name rather than something that can
// climb out of the directory the harness was pointed at.
func writePNG(path string, img image.Image) error {
	clean := filepath.Clean(path)
	if strings.Contains(clean, "..") {
		return fmt.Errorf("refusing to write outside the target directory: %s", path)
	}
	if filepath.Ext(clean) != ".png" {
		return fmt.Errorf("capture path must end in .png: %s", path)
	}

	f, err := os.Create(clean)
	if err != nil {
		return err
	}
	if err := png.Encode(f, img); err != nil {
		f.Close()
		return err
	}
	return f.Close()
}

func envInt(name string, fallback int) int {
	if v := os.Getenv(name); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return fallback
}

// words are recombined per row so that rows share glyphs and whole words while
// remaining distinct strings - the pattern real list content tends to follow.
var words = []string{
	"alpha", "bravo", "charlie", "delta", "echo", "foxtrot", "golf", "hotel",
	"india", "juliet", "kilo", "lima", "mike", "november", "oscar", "papa",
}

func rowText(i int) string {
	return fmt.Sprintf("%05d  %s %s %s  (%d items, %0.2f%%)",
		i, words[i%len(words)], words[(i/3)%len(words)], words[(i/7)%len(words)],
		i*3%997, float32(i)/97.0)
}

// Defaults for the scroll workload. The window is roughly phone shaped, and the
// scroll step is deliberately not a whole row height so that rows are drawn at
// many different sub-pixel offsets rather than repeating a short cycle.
const (
	defaultRows   = 5000
	defaultWidth  = 400
	defaultHeight = 800
	defaultStepPx = 7
	defaultTickMs = 8
)

func main() {
	rows := envInt("TEXTBENCH_ROWS", defaultRows)
	width := envInt("TEXTBENCH_WIDTH", defaultWidth)
	height := envInt("TEXTBENCH_HEIGHT", defaultHeight)
	stepPx := envInt("TEXTBENCH_STEP", defaultStepPx)
	tickMs := envInt("TEXTBENCH_TICK_MS", defaultTickMs)

	a := app.New()
	w := a.NewWindow("textbench")

	list := widget.NewList(
		func() int { return rows },
		func() fyne.CanvasObject { return widget.NewLabel("") },
		func(id widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(rowText(id))
		},
	)

	w.SetContent(container.NewStack(list))
	w.Resize(fyne.NewSize(float32(width), float32(height)))

	// TEXTBENCH_CAPTURE writes a PNG of the canvas and exits, so the rendering
	// can be eyeballed for correctness before any timing number is believed.
	if out := os.Getenv("TEXTBENCH_CAPTURE"); out != "" {
		go func() {
			time.Sleep(2 * time.Second)
			fyne.Do(func() {
				list.ScrollToOffset(400)
			})
			time.Sleep(time.Second)
			fyne.Do(func() {
				img := w.Canvas().Capture()
				if err := writePNG(out, img); err != nil {
					fmt.Fprintln(os.Stderr, "capture:", err)
					os.Exit(1)
				}
				fmt.Fprintln(os.Stderr, "captured to", out)
				os.Exit(0)
			})
		}()
		w.ShowAndRun()
		return
	}

	go func() {
		// Let the window map and the first frames settle before scrolling.
		time.Sleep(time.Second)

		offset := float32(0)
		maxOffset := float32(rows) * 40 // approximate; wraps well before the end
		ticker := time.NewTicker(time.Duration(tickMs) * time.Millisecond)
		defer ticker.Stop()

		for range ticker.C {
			offset += float32(stepPx)
			if offset > maxOffset {
				offset = 0
			}
			off := offset
			fyne.Do(func() {
				list.ScrollToOffset(off)
			})
		}
	}()

	w.ShowAndRun()
}
