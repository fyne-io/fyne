package widget

import (
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"
)

// layoutCountingSegment is a rich text segment that records how often the row
// bounds walk measures it while an entry lays its content out. One walk asks a
// segment for its text once, so the counter is a probe for the number of layout
// passes a single refresh performs.
type layoutCountingSegment struct {
	*TextSegment
	walks int
}

func (s *layoutCountingSegment) Textual() string {
	if calledFromRowBoundsWalk() {
		s.walks++
	}
	return s.TextSegment.Textual()
}

// calledFromRowBoundsWalk reports whether the current call stack passes through
// the row bounds builder, so that only layout passes are counted.
func calledFromRowBoundsWalk() bool {
	callers := make([]uintptr, 16)
	n := runtime.Callers(2, callers)
	frames := runtime.CallersFrames(callers[:n])
	for {
		frame, more := frames.Next()
		if strings.Contains(frame.Function, "rowBoundsBuilder).appendObject") {
			return true
		}
		if !more {
			return false
		}
	}
}

// TestEntryRefreshLaysOutOnce reproduces the text-entry lag: entryRenderer.Refresh
// laid the text out twice on every keystroke by calling updateRowBounds before
// the rich text refreshed itself, which lays out again. A steady-state refresh
// must walk the rows once, not twice.
func TestEntryRefreshLaysOutOnce(t *testing.T) {
	test.NewApp()

	e := NewMultiLineEntry()
	e.Wrapping = fyne.TextWrapWord
	e.Resize(fyne.NewSize(200, 200))
	e.SetText("editable text that has to be laid out for every keystroke")

	probe := &layoutCountingSegment{TextSegment: &TextSegment{Style: RichTextStyleInline, Text: e.Text}}
	e.textWidget.Segments = []RichTextSegment{e.textWidget.Segments[0], probe}

	e.Refresh() // warm up the one-off renderer and segment builds

	probe.walks = 0
	e.textWidget.Refresh()
	widgetWalks := probe.walks

	probe.walks = 0
	e.Refresh()
	entryWalks := probe.walks

	assert.Equal(t, 1, widgetWalks, "a rich text refresh must walk its rows once")
	assert.Equal(t, widgetWalks, entryWalks,
		"entryRenderer.Refresh walked the rows %d times, not the %d of a single layout", entryWalks, widgetWalks)
}
