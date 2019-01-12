package widget

import (
	"image/color"
	"testing"

	"fyne.io/fyne"
	"github.com/stretchr/testify/assert"
)

func TestText_Alignment(t *testing.T) {
	text := &textWidget{
		Alignment: fyne.TextAlignTrailing,
	}
	text.SetText("Test")
	assert.Equal(t, fyne.TextAlignTrailing, Renderer(text).(*textRenderer).texts[0].Alignment)
}

func TestText_Rows(t *testing.T) {
	text := textWidget{}
	text.SetText("test")
	assert.Equal(t, 1, text.rows())

	text.SetText("test\ntest")
	assert.Equal(t, text.rows(), 2)

	text.SetText("test\ntest\ntest")
	assert.Equal(t, text.rows(), 3)

	text.SetText("\n")
	assert.Equal(t, text.rows(), 2)
}

func TestText_RowLength(t *testing.T) {
	text := textWidget{}
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

func TestText_InsertAt(t *testing.T) {
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
			text := &textWidget{
				buffer: tt.fields.buffer,
			}
			text.insertAt(tt.args.pos, tt.args.runes)
			assert.Equal(t, tt.wantBuffer, text.buffer)
		})
	}
}

func TestText_Insert(t *testing.T) {
	text := &textWidget{}
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
			text := &textWidget{
				buffer: tt.fields.buffer,
			}
			got := text.deleteFromTo(tt.args.lowBound, tt.args.highBound)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantBuffer, text.buffer)
		})
	}
}

func TestText_Color(t *testing.T) {
	text := newTextWidget("test")
	text.color = color.White
	renderer := Renderer(text).(*textRenderer)
	renderer.Refresh()
	assert.Equal(t, color.White, renderer.texts[0].Color)
}
