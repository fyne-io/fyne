package widget

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/internal/cache"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/theme"
)

// Inline code (inline + monospace) renders its text on a background fill using
// the same colour name as the code block, reusing RichTextStyleCodeInline with
// no new style, segment type or field.
func TestTextSegment_InlineCodeVisualHasBackground(t *testing.T) {
	test.NewTempApp(t)

	seg := &TextSegment{Style: RichTextStyleCodeInline, Text: "code"}
	vis := seg.Visual()

	c, ok := vis.(*fyne.Container)
	if !ok {
		t.Fatalf("inline code visual should be a container, got %T", vis)
	}
	assert.Len(t, c.Objects, 2)

	bg, ok := c.Objects[0].(*canvas.Rectangle)
	if !ok {
		t.Fatalf("first object should be the background rectangle, got %T", c.Objects[0])
	}
	assert.Equal(t, theme.Color(theme.ColorNameInputBackground), bg.FillColor)

	txt, ok := c.Objects[1].(*canvas.Text)
	if !ok {
		t.Fatalf("second object should be the text, got %T", c.Objects[1])
	}
	assert.Equal(t, "code", txt.Text)
	assert.True(t, txt.TextStyle.Monospace)
}

// Ordinary inline text is unchanged: a bare canvas.Text with no background.
func TestTextSegment_PlainInlineVisualUnchanged(t *testing.T) {
	test.NewTempApp(t)

	seg := &TextSegment{Style: RichTextStyleInline, Text: "plain"}
	_, ok := seg.Visual().(*canvas.Text)
	assert.True(t, ok, "plain inline text should remain a bare canvas.Text")
}

// Monospace inline text that is not code (e.g. a monospace Entry) must NOT get a
// background: the marker, not the monospace style, drives the backdrop.
func TestTextSegment_MonospaceInlineNotCodeUnchanged(t *testing.T) {
	test.NewTempApp(t)

	style := RichTextStyleInline
	style.TextStyle.Monospace = true
	seg := &TextSegment{Style: style, Text: "mono"}
	_, ok := seg.Visual().(*canvas.Text)
	assert.True(t, ok, "monospace non-code inline text should remain a bare canvas.Text")
}

// Update recolours the backdrop and refreshes the text of an inline code container.
func TestTextSegment_InlineCodeUpdate(t *testing.T) {
	test.NewTempApp(t)

	seg := &TextSegment{Style: RichTextStyleCodeInline, Text: "x"}
	vis := seg.Visual()
	seg.Text = "y"
	seg.Update(vis)

	c := vis.(*fyne.Container)
	assert.Equal(t, "y", c.Objects[1].(*canvas.Text).Text)
	assert.Equal(t, theme.Color(theme.ColorNameInputBackground), c.Objects[0].(*canvas.Rectangle).FillColor)
}

// Emphasised inline code (e.g. **`x`**) keeps its background. Emphasis mutates
// the segment's TextStyle (sets Bold), so its Style no longer equals the
// RichTextStyleCodeInline var — a struct-equality check would miss it, but the
// codeInline marker survives the mutation.
func TestTextSegment_EmphasisedInlineCodeHasBackground(t *testing.T) {
	test.NewTempApp(t)

	r := NewRichTextFromMarkdown("**`code`**")
	var seg *TextSegment
	for _, s := range r.Segments {
		if ts, ok := s.(*TextSegment); ok && ts.Style.TextStyle.Monospace {
			seg = ts
		}
	}
	if seg == nil {
		t.Fatal("no monospace inline-code segment found")
	}
	assert.True(t, seg.Style.TextStyle.Bold, "precondition: emphasis set Bold on the code segment")
	assert.NotEqual(t, RichTextStyleCodeInline, seg.Style, "precondition: style differs from the var once Bold is set")
	_, ok := seg.Visual().(*fyne.Container)
	assert.True(t, ok, "emphasised inline code should still get a background")
}

// The fenced code block shares the inline code background colour.
func TestRichCodeBlock_BackgroundColour(t *testing.T) {
	test.NewTempApp(t)

	cb := newRichCodeBlock("x")
	test.TempWidgetRenderer(t, cb)
	assert.Equal(t, theme.Color(theme.ColorNameInputBackground), cb.bg.FillColor)
}

// A rendered document with inline code lays out without panicking and the code
// fragment appears as a backed container in the renderer's objects.
func TestRichText_InlineCodeRenders(t *testing.T) {
	test.NewTempApp(t)

	r := NewRichTextFromMarkdown("a `code` b")
	w := test.NewTempWindow(t, r)
	w.Resize(fyne.NewSize(200, 100))

	renderer := cache.Renderer(r).(*textRenderer)
	found := false
	for _, o := range renderer.Objects() {
		c, ok := o.(*fyne.Container)
		if !ok || len(c.Objects) != 2 {
			continue
		}
		if _, ok := c.Objects[0].(*canvas.Rectangle); ok {
			found = true
		}
	}
	assert.True(t, found, "expected an inline-code background container among render objects")
}
