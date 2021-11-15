package painter

import (
	"errors"
	"image"
	"testing"

	"github.com/goki/freetype/truetype"
	"github.com/stretchr/testify/assert"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

func Test_compositeFace_Close(t *testing.T) {
	chosenFont := &truetype.Font{}
	fallbackFont := &truetype.Font{}

	t.Run("happy case", func(t *testing.T) {
		chosen := &mockFace{CloseFunc: func() error { return nil }}
		fallback := &mockFace{CloseFunc: func() error { return nil }}
		c := newFontWithFallback(chosen, fallback, chosenFont, fallbackFont)
		if assert.NoError(t, c.Close()) {
			assert.True(t, chosen.CloseInvoked)
			assert.True(t, fallback.CloseInvoked)
		}
	})

	t.Run("with fallback failing", func(t *testing.T) {
		chosen := &mockFace{CloseFunc: func() error { return nil }}
		e := errors.New("oops")
		fallback := &mockFace{CloseFunc: func() error { return e }}
		c := newFontWithFallback(chosen, fallback, chosenFont, fallbackFont)
		assert.ErrorIs(t, c.Close(), e)
		assert.True(t, chosen.CloseInvoked)
		assert.True(t, fallback.CloseInvoked)
	})

	t.Run("with primary failing", func(t *testing.T) {
		e := errors.New("oops")
		chosen := &mockFace{CloseFunc: func() error { return e }}
		fallback := &mockFace{CloseFunc: func() error { return nil }}
		c := newFontWithFallback(chosen, fallback, chosenFont, fallbackFont)
		assert.ErrorIs(t, c.Close(), e)
		assert.True(t, chosen.CloseInvoked)
		assert.True(t, fallback.CloseInvoked)
	})
}

type mockFace struct {
	CloseFunc    func() error
	CloseInvoked bool
	GlyphFunc    func(fixed.Point26_6, rune) (image.Rectangle, image.Image, image.Point, fixed.Int26_6, bool)
	GlyphInvoked bool
}

var _ font.Face = (*mockFace)(nil)

func (f *mockFace) Close() error {
	f.CloseInvoked = true
	return f.CloseFunc()
}

func (f *mockFace) Glyph(dot fixed.Point26_6, r rune) (dr image.Rectangle, mask image.Image, maskp image.Point, advance fixed.Int26_6, ok bool) {
	f.GlyphInvoked = true
	return f.GlyphFunc(dot, r)
}

func (f *mockFace) GlyphAdvance(r rune) (advance fixed.Int26_6, ok bool) {
	panic("implement me")
}

func (f *mockFace) GlyphBounds(r rune) (bounds fixed.Rectangle26_6, advance fixed.Int26_6, ok bool) {
	panic("implement me")
}

func (f *mockFace) Kern(r0, r1 rune) fixed.Int26_6 {
	panic("implement me")
}

func (f *mockFace) Metrics() font.Metrics {
	panic("implement me")
}
