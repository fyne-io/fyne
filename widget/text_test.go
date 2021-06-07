package widget

import (
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/internal/cache"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/theme"
	"github.com/stretchr/testify/assert"
)

func textRenderTexts(p fyne.Widget) []*canvas.Text {
	return cache.Renderer(p).(*textRenderer).texts
}

func trailingBoldErrorSegment() RichTextSegment {
	return RichTextSegment{
		Alignment: fyne.TextAlignTrailing,
		ColorName: theme.ColorNameError,
		TextStyle: fyne.TextStyle{Bold: true},
	}
}

func TestText_Alignment(t *testing.T) {
	seg := trailingBoldErrorSegment()
	seg.Text = "Test"
	text := NewRichText(seg)
	assert.Equal(t, fyne.TextAlignTrailing, test.WidgetRenderer(text).(*textRenderer).texts[0].Alignment)
}

func TestText_Row(t *testing.T) {
	text := NewRichTextWithText("")
	text.Segments[0].Text = "test"
	text.Refresh()

	assert.Nil(t, text.row(-1))
	assert.Nil(t, text.row(1))

	assert.Equal(t, []rune("test"), text.row(0))
}

func TestText_Rows(t *testing.T) {
	text := NewRichTextWithText("test")
	assert.Equal(t, 1, text.rows())

	text.Segments[0].Text = "test\ntest"
	text.Refresh()
	assert.Equal(t, 2, text.rows())

	text.Segments[0].Text = "test\ntest\ntest"
	text.Refresh()
	assert.Equal(t, 3, text.rows())

	text.Segments[0].Text = "test\ntest\ntest\n"
	text.Refresh()
	assert.Equal(t, 4, text.rows())

	text.Segments[0].Text = "\n"
	text.Refresh()
	assert.Equal(t, 2, text.rows())
}

func TestText_RowLength(t *testing.T) {
	text := NewRichTextWithText("test")

	rl := text.rowLength(0)
	assert.Equal(t, 4, rl)

	text.Segments[0].Text = "test\ntèsts"
	text.Refresh()
	rl = text.rowLength(0)
	assert.Equal(t, 4, rl)

	rl = text.rowLength(1)
	assert.Equal(t, 5, rl)

	text.Segments[0].Text = ""
	text.Refresh()
	rl = text.rowLength(0)
	assert.Equal(t, 0, rl)

	text.Segments[0].Text = "\nhello"
	text.Refresh()
	rl = text.rowLength(0)
	assert.Equal(t, 0, rl)

	rl = text.rowLength(1)
	assert.Equal(t, 5, rl)
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
			text.insertAt(tt.args.pos, tt.args.runes)
			assert.Equal(t, tt.wantBuffer, text.String())
		})
	}
}

func TestText_Insert(t *testing.T) {
	text := NewRichTextWithText("")
	text.insertAt(0, "a")
	assert.Equal(t, "a", text.String())
	text.insertAt(1, "\n")
	assert.Equal(t, "a\n", text.String())
	text.insertAt(2, "b")
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
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantBuffer, text.String())
		})
	}
}

func TestText_Color(t *testing.T) {
	text := NewRichText(trailingBoldErrorSegment())

	assert.Equal(t, theme.ErrorColor(), textRenderTexts(text)[0].Color)
}

func TestTextRenderer_ApplyTheme(t *testing.T) {
	label := NewLabel("Test\nLine2")
	render := test.WidgetRenderer(label).(*textRenderer)

	text1 := render.Objects()[0].(*canvas.Text)
	text2 := render.Objects()[0].(*canvas.Text)
	textSize1 := text1.TextSize
	textSize2 := text2.TextSize
	customTextSize1 := textSize1
	customTextSize2 := textSize2
	test.WithTestTheme(t, func() {
		render.Refresh()
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

	fullSize := provider.lineSizeToColumn(4, 0)
	assert.Equal(t, fullSize, provider.lineSizeToColumn(10, 0))
	assert.Greater(t, fullSize.Width, provider.lineSizeToColumn(2, 0).Width)
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
			got := splitLines(RichTextSegment{Text: tt.text})
			for i, wantRow := range tt.want {
				assert.Equal(t, wantRow[0], got[i].begin)
				assert.Equal(t, wantRow[1], got[i].end)
			}
		})
	}
}

func TestText_lineBounds(t *testing.T) {
	mockMeasurer := func(text []rune) float32 {
		return float32(len(text))
	}
	tests := []struct {
		name string
		text string
		wrap fyne.TextWrap
		want [][2]int
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
			wrap: fyne.TextTruncate,
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
			wrap: fyne.TextTruncate,
			want: [][2]int{
				{0, 6},
			},
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
			wrap: fyne.TextTruncate,
			want: [][2]int{
				{0, 10},
			},
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
			wrap: fyne.TextTruncate,
			want: [][2]int{
				{0, 3},
				{4, 7},
			},
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
			wrap: fyne.TextTruncate,
			want: [][2]int{
				{0, 6},
				{7, 17},
				{28, 38},
			},
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
			wrap: fyne.TextTruncate,
			want: [][2]int{
				{0, 6},
				{7, 17},
				{26, 36},
				{39, 39},
			},
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
			wrap: fyne.TextTruncate,
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
			wrap: fyne.TextTruncate,
			want: [][2]int{
				{0, 6},
				{7, 17},
				{28, 38},
				{42, 42},
			},
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := lineBounds(RichTextSegment{Text: tt.text}, tt.wrap, 10, mockMeasurer)
			for i, wantRow := range tt.want {
				assert.Equal(t, wantRow[0], got[i].begin)
				assert.Equal(t, wantRow[1], got[i].end)
			}
		})
	}
}

func TestText_lineBounds_variable_char_width(t *testing.T) {
	tests := []struct {
		name string
		text string
		wrap fyne.TextWrap
		want [][2]int
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
			wrap: fyne.TextTruncate,
			want: [][2]int{
				{0, 12},
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
	measurer := func(text []rune) float32 {
		return fyne.MeasureText(string(text), textSize, textStyle).Width
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := lineBounds(RichTextSegment{Text: tt.text}, tt.wrap, 50, measurer)
			for i, wantRow := range tt.want {
				assert.Equal(t, wantRow[0], got[i].begin)
				assert.Equal(t, wantRow[1], got[i].end)
			}
		})
	}
}

func TestText_binarySearch(t *testing.T) {
	maxWidth := float32(50)
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
	} {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tt.want, findSpaceIndex([]rune(tt.text), len(tt.text)-1))
		})
	}
}
