//go:build darwin

// Command main_loop is a manual fixture for measuring macOS accessibility latency.
//
// Build the fixture and the AX probe on each revision being compared:
//
//	go build -tags accessibility,migrated_fynedo -ldflags=-linkmode=external -o /tmp/fyne-main-loop ./internal/driver/glfw/testdata/main_loop
//	clang -framework Cocoa -framework ApplicationServices internal/driver/glfw/testdata/accessibility_query.m -o /tmp/fyne-ax-query
//	/tmp/fyne-main-loop -duration=30s -elements=16 -repaint -repaint-delay=1s
//
// After the fixture prints "READY <pid>", run "/tmp/fyne-ax-query <pid> 200" from
// another terminal. The probe needs macOS Accessibility permission and reports
// median, p95 and maximum window-attribute and hit-test latency. The repaint delay
// lets the probe retain an element before repainting starts. Repeat without
// -repaint and with -elements=256. On exit the fixture prints animation frames,
// accessibility updates and CPU time as JSON.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"sync/atomic"
	"syscall"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type accessibleText struct {
	widget.Label
	updates atomic.Int64
}

func (t *accessibleText) AccessibilityLabel() string {
	t.updates.Add(1)
	return t.Text
}

func (*accessibleText) AccessibilityRole() fyne.AccessibleRole {
	return fyne.AccessibleRoleText
}

type result struct {
	Frames  int64         `json:"frames"`
	Updates int64         `json:"updates"`
	Elapsed time.Duration `json:"elapsed_ns"`
	CPU     time.Duration `json:"cpu_ns"`
}

func cpuTime() time.Duration {
	var usage syscall.Rusage
	if err := syscall.Getrusage(syscall.RUSAGE_SELF, &usage); err != nil {
		panic(err)
	}
	return time.Duration(usage.Utime.Nano() + usage.Stime.Nano())
}

func main() {
	const defaultElements = 16
	const windowSize = 400
	elements := flag.Int("elements", defaultElements, "Number of accessible text elements")
	duration := flag.Duration("duration", 5*time.Second, "Measurement duration after startup")
	repaint := flag.Bool("repaint", false, "Continuously repaint at the animation cadence")
	repaintDelay := flag.Duration("repaint-delay", 0, "Delay repainting after readiness so an AX client can retain an element")
	flag.Parse()
	if *elements < 1 || *duration <= 0 {
		panic("elements and duration must be positive")
	}

	a := app.NewWithID("io.fyne.test.main-loop")
	w := a.NewWindow("Fyne native event-loop probe")
	objects := make([]fyne.CanvasObject, *elements)
	var first *accessibleText
	for i := range objects {
		text := &accessibleText{Label: widget.Label{Text: fmt.Sprintf("Element %d", i)}}
		text.ExtendBaseWidget(text)
		objects[i] = text
		if i == 0 {
			first = text
		}
	}
	w.SetContent(container.NewVBox(objects...))
	w.Resize(fyne.NewSize(windowSize, windowSize))
	var frames atomic.Int64
	animation := fyne.NewAnimation(time.Second, func(float32) {
		frames.Add(1)
		first.Refresh()
	})
	animation.RepeatCount = fyne.AnimationRepeatForever

	var measured result
	a.Lifecycle().SetOnStarted(func() {
		go func() {
			time.Sleep(time.Second)
			if *repaint && *repaintDelay == 0 {
				fyne.DoAndWait(animation.Start)
			}
			time.Sleep(time.Second)
			frames.Store(0)
			first.updates.Store(0)
			start, cpu := time.Now(), cpuTime()
			fmt.Printf("READY %d\n", os.Getpid())
			if *repaint && *repaintDelay > 0 {
				time.AfterFunc(*repaintDelay, func() { fyne.Do(animation.Start) })
			}
			time.Sleep(time.Until(start.Add(*duration)))
			fyne.DoAndWait(func() {
				animation.Stop()
				measured.Frames = frames.Load()
				measured.Updates = first.updates.Load()
				measured.Elapsed = time.Since(start)
				measured.CPU = cpuTime() - cpu
			})
			a.Quit()
		}()
	})
	w.ShowAndRun()
	if err := json.NewEncoder(os.Stdout).Encode(measured); err != nil {
		panic(err)
	}
}
