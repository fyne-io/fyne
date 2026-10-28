package widget

import (
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/internal/cache"
	"fyne.io/fyne/v2/test"
)

// longWrappedRichText builds one wrapping TextSegment long enough to fill many
// rows, so a refresh has to slice the same backing segment onto every row. That
// is the shape the text handling hot path is measured on.
func longWrappedRichText(tb testing.TB) *RichText {
	tb.Helper()

	text := NewRichTextWithText(strings.Repeat("the quick brown fox jumps over the lazy dog ", 70))
	text.Wrapping = fyne.TextWrapWord
	text.Resize(fyne.NewSize(50, 400))
	return text
}

// TestTextRendererRefreshReusesSegmentRunes reproduces the Android text-entry
// lag: textRenderer.Refresh converted the whole backing segment to runes once
// for every visible row, so a long wrapped segment cost O(rows * runes) work per
// keystroke. Refreshing a wrapped segment must reuse its runes across rows.
func TestTextRendererRefreshReusesSegmentRunes(t *testing.T) {
	test.NewApp()

	long := longWrappedRichText(t)
	short := NewRichTextWithText("short")
	short.Wrapping = fyne.TextWrapWord
	short.Resize(fyne.NewSize(200, 400))

	longRenderer := cache.Renderer(long).(*textRenderer)
	shortRenderer := cache.Renderer(short).(*textRenderer)
	longRenderer.Refresh()
	shortRenderer.Refresh()
	longRenderer.Refresh()
	shortRenderer.Refresh()

	// The bytes one refresh allocates, taking the cheapest of a few samples so a
	// single noisy run cannot decide the result.
	measure := func(r *textRenderer) uint64 {
		least := ^uint64(0)
		for i := 0; i < 3; i++ {
			var before, after runtime.MemStats
			runtime.ReadMemStats(&before)
			r.Refresh()
			runtime.ReadMemStats(&after)
			if delta := after.TotalAlloc - before.TotalAlloc; delta < least {
				least = delta
			}
		}
		return least
	}

	shortBytes := measure(shortRenderer)
	longBytes := measure(longRenderer)
	extra := uint64(0)
	if longBytes > shortBytes {
		extra = longBytes - shortBytes
	}

	// The long and short texts share everything but the wrapped segment, so the
	// extra bytes are the price of laying that segment out. One conversion of the
	// segment's runes costs a few tens of kilobytes; converting them again for
	// every row costs megabytes.
	assert.Less(t, extra, uint64(1<<20),
		"textRenderer.Refresh converted the wrapped segment to runes once per visible row (extra %d bytes over %d rows)", extra, long.rows())
}

// TestRichTextRefreshKeepsRowText checks the happy path: sharing the runes across
// the rows of a wrapped segment must still slice the correct text onto each row,
// before and after a refresh.
func TestRichTextRefreshKeepsRowText(t *testing.T) {
	test.NewApp()

	text := longWrappedRichText(t)
	rows := richTextRowTexts(text)
	if !assert.Greater(t, len(rows), 5, "expected the text to wrap onto many rows") {
		return
	}

	drawn := renderTexts(text)
	assert.Equal(t, rows, drawn, "each row must render exactly the runes the bounds assigned it")

	text.Refresh()
	assert.Equal(t, drawn, renderTexts(text), "a refresh must reproduce the same rows, not lose or repeat runes")
}

// TestRichTextRefreshEmpty checks the empty boundary: refreshing an empty rich
// text stays at most one empty row, draws no text and does not panic.
func TestRichTextRefreshEmpty(t *testing.T) {
	test.NewApp()

	text := NewRichTextWithText("")
	text.Wrapping = fyne.TextWrapWord
	text.Resize(fyne.NewSize(200, 400))

	r := cache.Renderer(text).(*textRenderer)
	r.Refresh()
	r.Refresh()

	drawn := renderTexts(text)
	assert.LessOrEqual(t, len(drawn), 1, "an empty rich text must draw at most one row")
	for _, row := range drawn {
		assert.Empty(t, row, "an empty rich text must draw no text")
	}
	assert.Zero(t, len(text.Segments[0].Textual()), "an empty rich text must keep no content")
}

// renderTexts returns the text each of the rich text's rendered rows draws.
func renderTexts(rich *RichText) []string {
	var drawn []string
	for _, obj := range richTextRenderTexts(rich) {
		drawn = append(drawn, obj.Text)
	}
	return drawn
}
