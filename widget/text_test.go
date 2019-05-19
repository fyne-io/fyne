package widget

import (
	"image/color"
	"testing"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/test"
	"github.com/stretchr/testify/assert"
)

func textRenderTexts(p fyne.Widget) []*canvas.Text {
	return Renderer(p).(*textRenderer).texts
}

type testTextPresenter struct {
	obj   fyne.Widget
	fg    color.Color
	style fyne.TextStyle
	align fyne.TextAlign
}

func (t *testTextPresenter) textAlign() fyne.TextAlign {
	return t.align
}

func (t *testTextPresenter) textStyle() fyne.TextStyle {
	return t.style
}

func (t *testTextPresenter) textColor() color.Color {
	return t.fg
}

func (t *testTextPresenter) password() bool {
	return false
}

func (t *testTextPresenter) object() fyne.Widget {
	return t.obj
}

func newTestTextPresenter() textPresenter {
	t := &testTextPresenter{}
	t.obj = NewLabel("")

	return t
}

func newTrailingBoldWhiteTextPresenter() textPresenter {
	t := &testTextPresenter{}
	t.style = fyne.TextStyle{Bold: true}
	t.align = fyne.TextAlignTrailing
	t.fg = color.White

	t.obj = NewLabel("")
	return t
}

func TestText_Alignment(t *testing.T) {
	text := &textProvider{presenter: newTrailingBoldWhiteTextPresenter()}
	text.SetText("Test")
	assert.Equal(t, fyne.TextAlignTrailing, test.WidgetRenderer(text).(*textRenderer).texts[0].Alignment)
}

func TestText_Row(t *testing.T) {
	text := &textProvider{presenter: newTestTextPresenter()}
	text.SetText("test")

	assert.Nil(t, text.row(-1))
	assert.Nil(t, text.row(1))

	assert.Equal(t, []rune("test"), text.row(0))
}

func TestText_Color(t *testing.T) {
	text := &textProvider{presenter: newTrailingBoldWhiteTextPresenter()}
	Refresh(text.presenter.object())

	assert.Equal(t, color.White, textRenderTexts(text)[0].Color)
}

func TestTextHandler_Rows(t *testing.T) {
	text := &textHandler{}
	text.SetText("test")
	assert.Equal(t, 1, text.rows())

	text.SetText("test\ntest")
	assert.Equal(t, text.rows(), 2)

	text.SetText("test\ntest\ntest")
	assert.Equal(t, text.rows(), 3)

	text.SetText("\n")
	assert.Equal(t, text.rows(), 2)
}

func TestTextHandler_RowLength(t *testing.T) {
	text := &textHandler{}
	text.SetText("test")

	rl := text.rowLength(0)
	assert.Equal(t, 4, rl)

	text.SetText("test\ntèsts")
	rl = text.rowLength(0)
	assert.Equal(t, 4, rl)

	rl = text.rowLength(1)
	assert.Equal(t, 5, rl)

	text.SetText("")
	rl = text.rowLength(0)
	assert.Equal(t, 0, rl)

	text.SetText("\nhello")
	rl = text.rowLength(0)
	assert.Equal(t, 0, rl)

	rl = text.rowLength(1)
	assert.Equal(t, 5, rl)
}

func TestTextHandler_InsertAt(t *testing.T) {
	type fields struct {
		buffer []rune
	}
	type args struct {
		pos   int
		runes []rune
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantBuffer []rune
	}{
		{
			name:   "case_1",
			fields: fields{buffer: []rune("A\n1")},
			args: args{
				pos:   0,
				runes: []rune("\n"),
			},
			wantBuffer: []rune("\nA\n1"),
		},
		{
			name:   "case_2",
			fields: fields{buffer: []rune("hello\nèé+^#")},
			args: args{
				pos:   5,
				runes: []rune("\naddme"),
			},
			wantBuffer: []rune("hello\naddme\nèé+^#"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			text := &textHandler{}
			text.buffer = tt.fields.buffer

			text.insertAt(tt.args.pos, tt.args.runes)
			assert.Equal(t, tt.wantBuffer, text.buffer)
		})
	}
}

func TestTextHandler_Insert(t *testing.T) {
	text := &textHandler{}
	text.insertAt(0, []rune("a"))
	assert.Equal(t, []rune("a"), text.buffer)
	text.insertAt(1, []rune("\n"))
	assert.Equal(t, []rune("a\n"), text.buffer)
	text.insertAt(2, []rune("b"))
	assert.Equal(t, []rune("a\nb"), text.buffer)
}

func TestText_DeleteFromTo(t *testing.T) {
	type fields struct {
		buffer []rune
	}
	type args struct {
		lowBound  int
		highBound int
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		want       []rune
		wantBuffer []rune
	}{
		{
			name:   "case_1",
			fields: fields{buffer: []rune("A\n1")},
			args: args{
				lowBound:  0,
				highBound: 1,
			},
			want:       []rune("A"),
			wantBuffer: []rune("\n1"),
		},
		{
			name:   "case_2",
			fields: fields{buffer: []rune("A\n1")},
			args: args{
				lowBound:  1,
				highBound: 2,
			},
			want:       []rune("\n"),
			wantBuffer: []rune("A1"),
		},
		{
			name:   "case_3",
			fields: fields{buffer: []rune("A\nè1")},
			args: args{
				lowBound:  1,
				highBound: 3,
			},
			want:       []rune("\nè"),
			wantBuffer: []rune("A1"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			text := &textHandler{}
			text.buffer = tt.fields.buffer

			got := text.deleteFromTo(tt.args.lowBound, tt.args.highBound)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantBuffer, text.buffer)
		})
	}
}

func TestTextHandler_SetText(t *testing.T) {
	text := &textHandler{}
	text.SetText("test")

	assert.Equal(t, "test", string(text.buffer))
	assert.Equal(t, 1, len(text.rowBounds))
	assert.Equal(t, 0, text.rowBounds[0][0])
	assert.Equal(t, 4, text.rowBounds[0][1])
}

func TestTextHandler_String(t *testing.T) {
	text := &textHandler{}
	text.SetText("test")

	assert.Equal(t, "test", text.String())
}

func TestTextHandler_MaxCols(t *testing.T) {
	text := &textHandler{}
	text.SetText("test")

	assert.Equal(t, 4, text.maxCols)
}

func TestTextRenderer_ApplyTheme(t *testing.T) {
	label := NewLabel("Test\nLine2")
	render := test.WidgetRenderer(label).(*textRenderer)

	text1 := render.objects[0].(*canvas.Text)
	text2 := render.objects[0].(*canvas.Text)
	customTextSize1 := text1.TextSize
	customTextSize2 := text2.TextSize
	withTestTheme(func() {
		render.applyTheme()
		customTextSize1 = text1.TextSize
		customTextSize2 = text2.TextSize
	})

	assert.Equal(t, testTextSize, customTextSize1)
	assert.Equal(t, testTextSize, customTextSize2)
}

func TestTextProvider_LineSizeToColumn(t *testing.T) {
	label := NewLabel("Test")
	label.CreateRenderer() // TODO make this a simple refresh call once it's in
	provider := label.textProvider

	fullSize := provider.lineSizeToColumn(4, 0)
	assert.Equal(t, fullSize, provider.lineSizeToColumn(10, 0))
	assert.Greater(t, fullSize.Width, provider.lineSizeToColumn(2, 0).Width)
}
