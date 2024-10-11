package widget

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/internal/cache"
	"fyne.io/fyne/v2/internal/widget"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/theme"
)

func richTextRenderTexts(rich fyne.Widget) []*canvas.Text {
	renderer := cache.Renderer(rich).(*textRenderer)
	texts := make([]*canvas.Text, len(renderer.Objects()))
	for i, obj := range renderer.Objects() {
		texts[i] = obj.(*canvas.Text)
	}
	return texts
}

func trailingBoldErrorSegment() *TextSegment {
	return &TextSegment{Style: RichTextStyle{
		Alignment: fyne.TextAlignTrailing,
		ColorName: theme.ColorNameError,
		TextStyle: fyne.TextStyle{Bold: true},
	}}
}

func TestRichText_Hyperlink_Endline(t *testing.T) {
	u, _ := url.Parse("https://github.com/fyne-io/fyne")
	r := NewRichText(
		&TextSegment{Text: "Text", Style: RichTextStyleInline},
		&HyperlinkSegment{Text: "Link", URL: u},
	)
	r.Resize(r.MinSize())
	view := cache.Renderer(r)

	assert.Equal(t, 2, len(view.Objects()))
	assert.Equal(t, view.Objects()[0].Position().Y, view.Objects()[1].Position().Y)   // same baseline
	assert.Greater(t, view.Objects()[1].Position().X, view.Objects()[0].Position().X) // to the right
}

func TestText_Alignment(t *testing.T) {
	seg := trailingBoldErrorSegment()
	seg.Text = "Test"
	text := NewRichText(seg)
	assert.Equal(t, fyne.TextAlignTrailing, test.TempWidgetRenderer(t, text).Objects()[0].(*canvas.Text).Alignment)
}

func TestText_Row(t *testing.T) {
	text := NewRichTextWithText("")
	text.Segments[0].(*TextSegment).Text = "test"
	text.Refresh()

	assert.Nil(t, text.row(-1))
	assert.Nil(t, text.row(1))

	assert.Equal(t, []rune("test"), text.row(0))
}

func TestText_Rows(t *testing.T) {
	text := NewRichTextWithText("test")
	assert.Equal(t, 1, text.rows())
	textSeg := text.Segments[0].(*TextSegment)

	textSeg.Text = "test\ntest"
	text.Refresh()
	assert.Equal(t, 2, text.rows())

	textSeg.Text = "test\ntest\ntest"
	text.Refresh()
	assert.Equal(t, 3, text.rows())

	textSeg.Text = "test\ntest\ntest\n"
	text.Refresh()
	assert.Equal(t, 4, text.rows())

	textSeg.Text = "\n"
	text.Refresh()
	assert.Equal(t, 2, text.rows())
}

func TestText_RowLength(t *testing.T) {
	text := NewRichTextWithText("test")

	rl := text.rowLength(0)
	assert.Equal(t, 4, rl)
	textSeg := text.Segments[0].(*TextSegment)

	textSeg.Text = "test\ntèsts"
	text.Refresh()
	rl = text.rowLength(0)
	assert.Equal(t, 4, rl)

	rl = text.rowLength(1)
	assert.Equal(t, 5, rl)

	textSeg.Text = ""
	text.Refresh()
	rl = text.rowLength(0)
	assert.Equal(t, 0, rl)

	textSeg.Text = "\nhello"
	text.Refresh()
	rl = text.rowLength(0)
	assert.Equal(t, 0, rl)

	rl = text.rowLength(1)
	assert.Equal(t, 5, rl)
}

func TestText_Scroll(t *testing.T) {
	text1 := NewRichTextWithText("test1\ntest2")
	text2 := NewRichTextWithText("test1\ntest2")
	text2.Scroll = widget.ScrollBoth

	assert.Less(t, text2.MinSize().Width, text1.MinSize().Width)
	assert.Less(t, text2.MinSize().Height, text1.MinSize().Height)

	text3 := NewRichTextWithText("test1\ntest2")
	text3.Scroll = widget.ScrollVerticalOnly
	assert.Equal(t, text3.MinSize().Width, text1.MinSize().Width)
	assert.Less(t, text3.MinSize().Height, text1.MinSize().Height)

	text4 := NewRichTextWithText("test1\ntest2")
	text4.Scroll = widget.ScrollVerticalOnly
	text4.Wrapping = fyne.TextWrapBreak

	text3.Resize(fyne.NewSize(32, 32))
	text4.Resize(fyne.NewSize(32, 32))
	assert.Less(t, text4.MinSize().Width, text3.MinSize().Width)
	assert.Equal(t, text4.MinSize().Height, text3.MinSize().Height)

	content3 := test.TempWidgetRenderer(t, text3).Objects()[0].(*widget.Scroll).Content
	content4 := test.TempWidgetRenderer(t, text4).Objects()[0].(*widget.Scroll).Content
	assert.Less(t, content4.MinSize().Width, content3.MinSize().Width)
	assert.Greater(t, content4.MinSize().Height, content3.MinSize().Height)
}

func TestText_InsertAt(t *testing.T) {
	type fields struct {
		buffer string
	}
	type args struct {
		pos   int
		runes string
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantBuffer string
	}{
		{
			name:   "case_1",
			fields: fields{buffer: "A\n1"},
			args: args{
				pos:   0,
				runes: "\n",
			},
			wantBuffer: "\nA\n1",
		},
		{
			name:   "case_2",
			fields: fields{buffer: "hello\nèé+^#"},
			args: args{
				pos:   5,
				runes: "\naddme",
			},
			wantBuffer: "hello\naddme\nèé+^#",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			text := NewRichTextWithText(tt.fields.buffer)
			text.insertAt(tt.args.pos, []rune(tt.args.runes))
			assert.Equal(t, tt.wantBuffer, text.String())
		})
	}
}

func TestText_Insert(t *testing.T) {
	text := NewRichTextWithText("")
	text.insertAt(0, []rune("a"))
	assert.Equal(t, "a", text.String())
	text.insertAt(1, []rune("\n"))
	assert.Equal(t, "a\n", text.String())
	text.insertAt(2, []rune("b"))
	assert.Equal(t, "a\nb", text.String())
}

func TestText_DeleteFromTo(t *testing.T) {
	type fields struct {
		buffer string
	}
	type args struct {
		lowBound  int
		highBound int
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		want       string
		wantBuffer string
	}{
		{
			name:   "case_1",
			fields: fields{buffer: "A\n1"},
			args: args{
				lowBound:  0,
				highBound: 1,
			},
			want:       "A",
			wantBuffer: "\n1",
		},
		{
			name:   "case_2",
			fields: fields{buffer: "A\n1"},
			args: args{
				lowBound:  1,
				highBound: 2,
			},
			want:       "\n",
			wantBuffer: "A1",
		},
		{
			name:   "case_3",
			fields: fields{buffer: "A\nè1"},
			args: args{
				lowBound:  1,
				highBound: 3,
			},
			want:       "\nè",
			wantBuffer: "A1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			text := NewRichTextWithText(tt.fields.buffer)
			got := text.deleteFromTo(tt.args.lowBound, tt.args.highBound)
			assert.Equal(t, tt.want, string(got))
			assert.Equal(t, tt.wantBuffer, text.String())
		})
	}
}

func TestText_DeleteFromTo_Segments(t *testing.T) {
	type args struct {
		lowBound  int
		highBound int
	}
	tests := []struct {
		name         string
		segments     []RichTextSegment
		args         args
		want         string
		wantSegments []RichTextSegment
	}{
		{
			name: "remove begin",
			segments: []RichTextSegment{
				&TextSegment{Text: "A\n"},
				&TextSegment{Text: "1"},
			},
			args: args{
				lowBound:  0,
				highBound: 1,
			},
			want: "A",
			wantSegments: []RichTextSegment{
				&TextSegment{Text: "\n"},
				&TextSegment{Text: "1"},
			},
		},
		{
			name: "remove end",
			segments: []RichTextSegment{
				&TextSegment{Text: "A"},
				&TextSegment{Text: "\n1"},
			},
			args: args{
				lowBound:  1,
				highBound: 2,
			},
			want: "\n",
			wantSegments: []RichTextSegment{
				&TextSegment{Text: "A"},
				&TextSegment{Text: "1"},
			},
		},
		{
			name: "remove both",
			segments: []RichTextSegment{
				&TextSegment{Text: "A\n"},
				&TextSegment{Text: "è1"},
			},
			args: args{
				lowBound:  1,
				highBound: 3,
			},
			want: "\nè",
			wantSegments: []RichTextSegment{
				&TextSegment{Text: "A"},
				&TextSegment{Text: "1"},
			},
		},
		{
			name: "remove nontext",
			segments: []RichTextSegment{
				&TextSegment{Text: "A\n"},
				&SeparatorSegment{},
				&TextSegment{Text: "B1"},
			},
			args: args{
				lowBound:  1,
				highBound: 3,
			},
			want: "\nB",
			wantSegments: []RichTextSegment{
				&TextSegment{Text: "A"},
				&TextSegment{Text: "1"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			text := NewRichText(tt.segments...)
			got := text.deleteFromTo(tt.args.lowBound, tt.args.highBound)
			assert.Equal(t, tt.want, string(got))
			for _, s := range tt.wantSegments {
				if txt, ok := s.(*TextSegment); ok {
					txt.parent = text
				}
			}
			assert.Equal(t, tt.wantSegments, text.Segments)
		})
	}
}

func TestText_Multiline(t *testing.T) {
	text := NewRichText(
		&TextSegment{Text: "line1\nli", Style: RichTextStyleStrong},
		&TextSegment{Text: "ne2\nline3", Style: RichTextStyleInline})

	w := test.NewTempWindow(t, text)
	w.Resize(fyne.NewSize(64, 90))
	test.AssertImageMatches(t, "richtext/richtext_multiline.png", w.Canvas().Capture())
}

func TestText_Color(t *testing.T) {
	text := NewRichText(trailingBoldErrorSegment())

	assert.Equal(t, theme.Color(theme.ColorNameError), richTextRenderTexts(text)[0].Color)
}

func TestTextRenderer_ApplyTheme(t *testing.T) {
	label := NewLabel("Test\nLine2")
	texts := labelTextRenderTexts(label)

	text1 := texts[0]
	text2 := texts[1]
	textSize1 := text1.TextSize
	textSize2 := text2.TextSize
	customTextSize1 := textSize1
	customTextSize2 := textSize2
	test.WithTestTheme(t, func() {
		label.Refresh()
		text1 := texts[0]
		text2 := texts[1]
		customTextSize1 = text1.TextSize
		customTextSize2 = text2.TextSize
	})

	assert.NotEqual(t, textSize1, customTextSize1)
	assert.NotEqual(t, textSize2, customTextSize2)
}

func TestTextProvider_LineSizeToColumn(t *testing.T) {
	label := NewLabel("Test")
	label.CreateRenderer() // TODO make this a simple refresh call once it's in
	provider := label.provider

	inPad := theme.InnerPadding()
	textSize := theme.TextSize()
	fullSize := provider.lineSizeToColumn(4, 0, textSize, inPad)
	assert.Equal(t, fullSize, provider.lineSizeToColumn(10, 0, textSize, inPad))
	assert.Greater(t, fullSize.Width, provider.lineSizeToColumn(2, 0, textSize, inPad).Width)

	out := provider.lineSizeToColumn(-1, -1, textSize, inPad)
	assert.Equal(t, out, provider.lineSizeToColumn(0, 0, textSize, inPad))
}

func TestText_splitLines(t *testing.T) {
	tests := []struct {
		name string
		text string
		want [][2]int
	}{
		{
			name: "Empty",
			text: "",
			want: [][2]int{
				{0, 0},
			},
		},
		{
			name: "Single",
			text: "foo",
			want: [][2]int{
				{0, 3},
			},
		},
		{
			name: "Multiple",
			text: "foo\nbar",
			want: [][2]int{
				{0, 3},
				{4, 7},
			},
		},
		{
			name: "Multiple_Trailing",
			text: "foo\nbar\n",
			want: [][2]int{
				{0, 3},
				{4, 7},
				{8, 8},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := splitLines(&TextSegment{Text: tt.text})
			for i, wantRow := range tt.want {
				assert.Equal(t, wantRow[0], got[i].begin)
				assert.Equal(t, wantRow[1], got[i].end)
			}
		})
	}
}

func TestText_lineBounds(t *testing.T) {
	measurer := func(text []rune) fyne.Size {
		return fyne.MeasureText(string(text), 14, fyne.TextStyle{})
	}
	tests := []struct {
		name     string
		text     string
		wrap     fyne.TextWrap
		trunc    fyne.TextTruncation
		want     [][2]int
		ellipses int
	}{
		{
			name: "Empty_WrapOff",
			text: "",
			wrap: fyne.TextWrapOff,
			want: [][2]int{
				{0, 0},
			},
		},
		{
			name: "Empty_Truncate",
			text: "",
			wrap: fyne.TextWrap(fyne.TextTruncateClip),
			want: [][2]int{
				{0, 0},
			},
		},
		{
			name:  "Empty_TruncateClip",
			text:  "",
			trunc: fyne.TextTruncateClip,
			want: [][2]int{
				{0, 0},
			},
		},
		{
			name: "Empty_WrapBreak",
			text: "",
			wrap: fyne.TextWrapBreak,
			want: [][2]int{
				{0, 0},
			},
		},
		{
			name: "Empty_WrapWord",
			text: "",
			wrap: fyne.TextWrapWord,
			want: [][2]int{
				{0, 0},
			},
		},
		{
			name: "Single_Short_WrapOff",
			text: "foobar",
			wrap: fyne.TextWrapOff,
			want: [][2]int{
				{0, 6},
			},
		},
		{
			name: "Single_Short_Truncate",
			text: "foobar",
			wrap: fyne.TextWrap(fyne.TextTruncateClip),
			want: [][2]int{
				{0, 6},
			},
		},
		{
			name:  "Single_Short_TruncateClip",
			text:  "foobar",
			trunc: fyne.TextTruncateClip,
			want: [][2]int{
				{0, 6},
			},
		},
		{
			name:  "Single_Short_TruncateEllipsis",
			text:  "foobar",
			trunc: fyne.TextTruncateEllipsis,
			want: [][2]int{
				{0, 6},
			},
			ellipses: 0, // too short to cut
		},
		{
			name: "Single_Short_WrapBreak",
			text: "foobar",
			wrap: fyne.TextWrapBreak,
			want: [][2]int{
				{0, 6},
			},
		},
		{
			name: "Single_Short_WrapWord",
			text: "foobar",
			wrap: fyne.TextWrapWord,
			want: [][2]int{
				{0, 6},
			},
		},
		{
			name: "Single_Long_WrapOff",
			text: "foobar foobar",
			wrap: fyne.TextWrapOff,
			want: [][2]int{
				{0, 13},
			},
		},
		{
			name: "Single_Long_Truncate",
			text: "foobar foobar",
			wrap: fyne.TextWrap(fyne.TextTruncateClip),
			want: [][2]int{
				{0, 10},
			},
		},
		{
			name:  "Single_Long_TruncateClip",
			text:  "foobar foobar",
			trunc: fyne.TextTruncateClip,
			want: [][2]int{
				{0, 10},
			},
		},
		{
			name:  "Single_Long_TruncateEllipsis",
			text:  "foobar foobar",
			trunc: fyne.TextTruncateEllipsis,
			want: [][2]int{
				{0, 9},
			},
			ellipses: 1,
		},
		{
			name: "Single_Long_WrapBreak",
			text: "foobar foobar",
			wrap: fyne.TextWrapBreak,
			want: [][2]int{
				{0, 10},
				{10, 13},
			},
		},
		{
			name: "Single_Long_WrapWord",
			text: "foobar foobar",
			wrap: fyne.TextWrapWord,
			want: [][2]int{
				{0, 6},
				{7, 13},
			},
		},
		{
			name: "Multiple_Short_WrapOff",
			text: "foo\nbar",
			wrap: fyne.TextWrapOff,
			want: [][2]int{
				{0, 3},
				{4, 7},
			},
		},
		{
			name: "Multiple_Short_Truncate",
			text: "foo\nbar",
			wrap: fyne.TextWrap(fyne.TextTruncateClip),
			want: [][2]int{
				{0, 3},
				{4, 7},
			},
		},
		{
			name:  "Multiple_Short_TruncateClip",
			text:  "foo\nbar",
			trunc: fyne.TextTruncateClip,
			want: [][2]int{
				{0, 3},
				{4, 7},
			},
		},
		{
			name:  "Multiple_Short_TruncateEllipsis",
			text:  "foo\nbar",
			trunc: fyne.TextTruncateEllipsis,
			want: [][2]int{
				{0, 3},
				{4, 7},
			},
			ellipses: 0, // too wide
		},
		{
			name: "Multiple_Short_WrapBreak",
			text: "foo\nbar",
			wrap: fyne.TextWrapBreak,
			want: [][2]int{
				{0, 3},
				{4, 7},
			},
		},
		{
			name: "Multiple_Short_WrapWord",
			text: "foo\nbar",
			wrap: fyne.TextWrapWord,
			want: [][2]int{
				{0, 3},
				{4, 7},
			},
		},
		{
			name: "Multiple_Long_WrapOff",
			text: "foobar\nfoobar foobar foobar\nfoobar foobar",
			wrap: fyne.TextWrapOff,
			want: [][2]int{
				{0, 6},
				{7, 27},
				{28, 41},
			},
		},
		{
			name: "Multiple_Long_Truncate",
			text: "foobar\nfoobar foobar foobar\nfoobar foobar",
			wrap: fyne.TextWrap(fyne.TextTruncateClip),
			want: [][2]int{
				{0, 6},
				{7, 17},
				{28, 38},
			},
		},
		{
			name:  "Multiple_Long_TruncateClip",
			text:  "foobar\nfoobar foobar foobar\nfoobar foobar",
			trunc: fyne.TextTruncateClip,
			want: [][2]int{
				{0, 6},
				{7, 17},
				{28, 38},
			},
		},
		{
			name:  "Multiple_Long_TruncateEllipsis",
			text:  "foobar\nfoobar foobar foobar\nfoobar foobar",
			trunc: fyne.TextTruncateEllipsis,
			want: [][2]int{
				{0, 6},
				{7, 16},
				{28, 37},
			},
			ellipses: 2,
		},
		{
			name: "Multiple_Long_WrapBreak",
			text: "foobar\nfoobar foobar foobar\nfoobar foobar",
			wrap: fyne.TextWrapBreak,
			want: [][2]int{
				{0, 6},
				{7, 17},
				{17, 27},
				{28, 38},
				{38, 41},
			},
		},
		{
			name: "Multiple_Long_WrapWord",
			text: "foobar\nfoobar foobar foobar\nfoobar foobar",
			wrap: fyne.TextWrapWord,
			want: [][2]int{
				{0, 6},
				{7, 13},
				{14, 20},
				{21, 27},
				{28, 34},
				{35, 41},
			},
		},
		{
			name: "Multiple_Contiguous_Long_WrapOff",
			text: "foobar\nfoobarfoobarfoobar\nfoobarfoobar\n",
			wrap: fyne.TextWrapOff,
			want: [][2]int{
				{0, 6},
				{7, 25},
				{26, 38},
				{39, 39},
			},
		},
		{
			name: "Multiple_Contiguous_Long_Truncate",
			text: "foobar\nfoobarfoobarfoobar\nfoobarfoobar\n",
			wrap: fyne.TextWrap(fyne.TextTruncateClip),
			want: [][2]int{
				{0, 6},
				{7, 17},
				{26, 36},
				{39, 39},
			},
		},
		{
			name:  "Multiple_Contiguous_Long_TruncateClip",
			text:  "foobar\nfoobarfoobarfoobar\nfoobarfoobar\n",
			trunc: fyne.TextTruncateClip,
			want: [][2]int{
				{0, 6},
				{7, 17},
				{26, 36},
				{39, 39},
			},
		},
		{
			name:  "Multiple_Contiguous_Long_TruncateEllipsis",
			text:  "foobar\nfoobarfoobarfoobar\nfoobarfoobar\n",
			trunc: fyne.TextTruncateEllipsis,
			want: [][2]int{
				{0, 6},
				{7, 15},
				{26, 34},
				{39, 39},
			},
			ellipses: 2,
		},
		{
			name: "Multiple_Contiguous_Long_WrapBreak",
			text: "foobar\nfoobarfoobarfoobar\nfoobarfoobar\n",
			wrap: fyne.TextWrapBreak,
			want: [][2]int{
				{0, 6},
				{7, 17},
				{17, 25},
				{26, 36},
				{36, 38},
				{39, 39},
			},
		},
		{
			name: "Multiple_Contiguous_Long_WrapWord",
			text: "foobar\nfoobarfoobarfoobar\nfoobarfoobar\n",
			wrap: fyne.TextWrapWord,
			want: [][2]int{
				{0, 6},
				{7, 17},
				{17, 25},
				{26, 36},
				{36, 38},
				{39, 39},
			},
		},
		{
			name: "Multiple_Trailing_Short_WrapOff",
			text: "foo\nbar\n",
			wrap: fyne.TextWrapOff,
			want: [][2]int{
				{0, 3},
				{4, 7},
				{8, 8},
			},
		},
		{
			name: "Multiple_Trailing_Short_Truncate",
			text: "foo\nbar\n",
			wrap: fyne.TextWrap(fyne.TextTruncateClip),
			want: [][2]int{
				{0, 3},
				{4, 7},
				{8, 8},
			},
		},
		{
			name:  "Multiple_Trailing_Short_TruncateClip",
			text:  "foo\nbar\n",
			trunc: fyne.TextTruncateClip,
			want: [][2]int{
				{0, 3},
				{4, 7},
				{8, 8},
			},
		},
		{
			name:  "Multiple_Trailing_Short_TruncateEllipsis",
			text:  "foo\nbar\n",
			trunc: fyne.TextTruncateEllipsis,
			want: [][2]int{
				{0, 3},
				{4, 7},
				{8, 8},
			},
		},
		{
			name: "Multiple_Trailing_Short_WrapBreak",
			text: "foo\nbar\n",
			wrap: fyne.TextWrapBreak,
			want: [][2]int{
				{0, 3},
				{4, 7},
				{8, 8},
			},
		},
		{
			name: "Multiple_Trailing_Short_WrapWord",
			text: "foo\nbar\n",
			wrap: fyne.TextWrapWord,
			want: [][2]int{
				{0, 3},
				{4, 7},
				{8, 8},
			},
		},
		{
			name:  "Multiple_Trailing_Short_WrapTruncateClip",
			text:  "foo\nbar\nbaz",
			wrap:  fyne.TextWrapBreak,
			trunc: fyne.TextTruncateClip,
			want: [][2]int{
				{0, 3},
				{4, 7},
			},
		},
		{
			name:  "Multiple_Trailing_Short_WrapTruncateEllipsis",
			text:  "foo\nbar\nbaz",
			wrap:  fyne.TextWrapBreak,
			trunc: fyne.TextTruncateEllipsis,
			want: [][2]int{
				{0, 3},
				{4, 7},
			},
			ellipses: 0,
		},
		{
			name: "Multiple_Trailing_Long_WrapOff",
			text: "foobar\nfoobar foobar foobar\nfoobar foobar\n",
			wrap: fyne.TextWrapOff,
			want: [][2]int{
				{0, 6},
				{7, 27},
				{28, 41},
				{42, 42},
			},
		},
		{
			name: "Multiple_Trailing_Long_Truncate",
			text: "foobar\nfoobar foobar foobar\nfoobar foobar\n",
			wrap: fyne.TextWrap(fyne.TextTruncateClip),
			want: [][2]int{
				{0, 6},
				{7, 17},
				{28, 38},
				{42, 42},
			},
		},
		{
			name:  "Multiple_Trailing_Long_TruncateClip",
			text:  "foobar\nfoobar foobar foobar\nfoobar foobar\n",
			trunc: fyne.TextTruncateClip,
			want: [][2]int{
				{0, 6},
				{7, 17},
				{28, 38},
				{42, 42},
			},
		},
		{
			name:  "Multiple_Trailing_Long_TruncateEllipsis",
			text:  "foobar\nfoobar foobar foobar\nfoobar foobar\n",
			trunc: fyne.TextTruncateEllipsis,
			want: [][2]int{
				{0, 6},
				{7, 16},
				{28, 37},
				{42, 42},
			},
			ellipses: 2,
		},
		{
			name: "Multiple_Trailing_Long_WrapBreak",
			text: "foobar\nfoobar foobar foobar\nfoobar foobar\n",
			wrap: fyne.TextWrapBreak,
			want: [][2]int{
				{0, 6},
				{7, 17},
				{17, 27},
				{28, 38},
				{38, 41},
				{42, 42},
			},
		},
		{
			name: "Multiple_Trailing_Long_WrapWord",
			text: "foobar\nfoobar foobar foobar\nfoobar foobar\n",
			wrap: fyne.TextWrapWord,
			want: [][2]int{
				{0, 6},
				{7, 13},
				{14, 20},
				{21, 27},
				{28, 34},
				{35, 41},
				{42, 42},
			},
		},
		{
			name:  "Multiple_Trailing_Long_WrapTruncateClip",
			text:  "foo\nfoobar foobar foobar\nbaz",
			wrap:  fyne.TextWrapBreak,
			trunc: fyne.TextTruncateClip,
			want: [][2]int{
				{0, 3},
				{4, 14},
			},
		},
		{
			name:  "Multiple_Trailing_Long_WrapTruncateEllipsis",
			text:  "foo\nfoobar foobar foobar\nbazzzzzzzzzzzzzzzzzzzzzzzzzzzzzzz",
			wrap:  fyne.TextWrapBreak,
			trunc: fyne.TextTruncateEllipsis,
			want: [][2]int{
				{0, 3},
				{4, 14},
				{14, 23},
			},
			ellipses: 1,
		},
		{
			name:  "Multiple_Trailing_Long_WrapWordTruncateEllipsis",
			text:  "foo\nfoobar foobar foobar\nbaz",
			wrap:  fyne.TextWrapWord,
			trunc: fyne.TextTruncateEllipsis,
			want: [][2]int{
				{0, 3},
				{4, 10},
				{11, 17},
			},
			ellipses: 1,
		},
		{
			name:  "Multi_byte_ellipsis_not_truncated",
			text:  "🪃 234",
			trunc: fyne.TextTruncateEllipsis,
			wrap:  fyne.TextWrapOff,
			want: [][2]int{
				{0, 5},
			},
			ellipses: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ellipses := 0
			got, _ := lineBounds(&TextSegment{Text: tt.text}, tt.wrap, tt.trunc, 76, fyne.NewSize(76, 64), measurer)
			for i, wantRow := range tt.want {
				assert.Equal(t, wantRow[0], got[i].begin)
				assert.Equal(t, wantRow[1], got[i].end)

				if got[i].ellipsis {
					ellipses++
				}
			}
			assert.Equal(t, tt.ellipses, ellipses)
		})
	}
}

func TestText_lineBounds_variable_char_width(t *testing.T) {
	tests := []struct {
		name  string
		text  string
		wrap  fyne.TextWrap
		trunc fyne.TextTruncation
		want  [][2]int
	}{
		{
			name: "IM_WrapOff",
			text: "iiiiiiiiiimmmmmmmmmm",
			wrap: fyne.TextWrapOff,
			want: [][2]int{
				{0, 20},
			},
		},
		{
			name: "IM_Truncate",
			text: "iiiiiiiiiimmmmmmmmmm",
			wrap: fyne.TextWrap(fyne.TextTruncateClip),
			want: [][2]int{
				{0, 12},
			},
		},
		{
			name:  "IM_TruncateClip",
			text:  "iiiiiiiiiimmmmmmmmmm",
			trunc: fyne.TextTruncateClip,
			want: [][2]int{
				{0, 12},
			},
		},
		{
			name:  "IM_TruncateEllipsis",
			text:  "iiiiiiiiiimmmmmmmmmm",
			trunc: fyne.TextTruncateEllipsis,
			want: [][2]int{
				{0, 9},
			},
		},
		{
			name: "IM_WrapBreak",
			text: "iiiiiiiiiimmmmmmmmmm",
			wrap: fyne.TextWrapBreak,
			want: [][2]int{
				{0, 12},
				{12, 16},
				{16, 20},
			},
		},
		{
			name: "IM_WrapWord",
			text: "iiiiiiiiiimmmmmmmmmm",
			wrap: fyne.TextWrapWord,
			want: [][2]int{
				{0, 12},
				{12, 16},
				{16, 20},
			},
		},
	}
	textSize := float32(10)
	textStyle := fyne.TextStyle{}
	measurer := func(text []rune) fyne.Size {
		return fyne.MeasureText(string(text), textSize, textStyle)
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := lineBounds(&TextSegment{Text: tt.text}, tt.wrap, tt.trunc, 46, fyne.NewSize(46, 184), measurer)
			for i, wantRow := range tt.want {
				assert.Equal(t, wantRow[0], got[i].begin)
				assert.Equal(t, wantRow[1], got[i].end)
			}
		})
	}
}

func TestText_binarySearch(t *testing.T) {
	maxWidth := float32(46)
	textSize := float32(10)
	textStyle := fyne.TextStyle{}
	measurer := func(text []rune) float32 {
		return fyne.MeasureText(string(text), textSize, textStyle).Width
	}
	for name, tt := range map[string]struct {
		text string
		want int
	}{
		"IM": {
			text: "iiiiiiiiiimmmmmmmmmm",
			want: 12,
		},
		"Single_Line": {
			text: "foobar foobar",
			want: 9,
		},
		"WH": {
			text: "wwwww hhhhhh",
			want: 6,
		},
		"DS": {
			text: "dddddd sssssss",
			want: 8,
		},
		"DI": {
			text: "dididi dididd",
			want: 10,
		},
		"XW": {
			text: "xwxwxwxw xwxw",
			want: 7,
		},
		"W": {
			text: "WWWWW",
			want: 4,
		},
		"Empty": {
			text: "",
			want: 0,
		},
	} {
		checker := func(low int, high int) bool {
			return measurer([]rune(tt.text[low:high])) <= maxWidth
		}
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tt.want, binarySearch(checker, 0, len(tt.text)))
		})
	}
}

func TestText_findSpaceIndex(t *testing.T) {
	for name, tt := range map[string]struct {
		text string
		want int
	}{
		"no_space_fallback": {
			text: "iiiiiiiiiimmmmmmmmmm",
			want: 19,
		},
		"single_space": {
			text: "foobar foobar",
			want: 6,
		},
		"double_space": {
			text: "ww wwww www",
			want: 7,
		},
		"many_spaces": {
			text: "ww wwww www wwwww",
			want: 11,
		},
		"space beginning": {
			text: " ww",
			want: 0,
		},
	} {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tt.want, findSpaceIndex([]rune(tt.text), len(tt.text)-1))
		})
	}
}
