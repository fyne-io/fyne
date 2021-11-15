package painter

import (
	"errors"
	"image"
	"image/color"
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

func Test_compositeFace_Glyph(t *testing.T) {
	chosenFont := &mockFont{IndexFunc: func(r rune) truetype.Index {
		if r == 'e' || r == 'p' {
			return 1
		}
		return 0
	}}
	fallbackFont := &mockFont{IndexFunc: func(r rune) truetype.Index {
		if r == 'e' || r == 's' {
			return 1
		}
		return 0
	}}

	t.Run("when primary has glyph", func(t *testing.T) {
		chosen := &mockFace{GlyphFunc: func(p fixed.Point26_6, r rune) (image.Rectangle, image.Image, image.Point, fixed.Int26_6, bool) {
			assert.Equal(t, fixed.Point26_6{X: 13}, p)
			assert.Equal(t, 'e', r)
			return image.Rect(1, 2, 3, 4), image.NewUniform(color.NRGBA{R: 5, G: 6, B: 7, A: 8}), image.Pt(9, 10), 11, true
		}}
		fallback := &mockFace{GlyphFunc: func(p fixed.Point26_6, r rune) (image.Rectangle, image.Image, image.Point, fixed.Int26_6, bool) {
			return image.Rectangle{}, nil, image.Point{}, 0, false
		}}
		c := newFontWithFallback(chosen, fallback, chosenFont, fallbackFont)
		r, m, mp, a, ok := c.Glyph(fixed.Point26_6{X: 13}, 'e')
		assert.True(t, chosen.GlyphInvoked)
		assert.False(t, fallback.GlyphInvoked)
		assert.Equal(t, image.Rect(1, 2, 3, 4), r)
		assert.Equal(t, image.NewUniform(color.NRGBA{R: 5, G: 6, B: 7, A: 8}), m)
		assert.Equal(t, image.Pt(9, 10), mp)
		assert.Equal(t, fixed.Int26_6(11), a)
		assert.True(t, ok)
	})

	t.Run("when only fallback has glyph", func(t *testing.T) {
		chosen := &mockFace{GlyphFunc: func(p fixed.Point26_6, r rune) (image.Rectangle, image.Image, image.Point, fixed.Int26_6, bool) {
			return image.Rectangle{}, nil, image.Point{}, 0, false
		}}
		fallback := &mockFace{GlyphFunc: func(p fixed.Point26_6, r rune) (image.Rectangle, image.Image, image.Point, fixed.Int26_6, bool) {
			assert.Equal(t, fixed.Point26_6{X: 42}, p)
			assert.Equal(t, 's', r)
			return image.Rect(21, 22, 23, 24), image.NewUniform(color.NRGBA{R: 25, G: 26, B: 27, A: 28}), image.Pt(29, 30), 31, true
		}}
		c := newFontWithFallback(chosen, fallback, chosenFont, fallbackFont)
		r, m, mp, a, ok := c.Glyph(fixed.Point26_6{X: 42}, 's')
		assert.False(t, chosen.GlyphInvoked)
		assert.True(t, fallback.GlyphInvoked)
		assert.Equal(t, image.Rect(21, 22, 23, 24), r)
		assert.Equal(t, image.NewUniform(color.NRGBA{R: 25, G: 26, B: 27, A: 28}), m)
		assert.Equal(t, image.Pt(29, 30), mp)
		assert.Equal(t, fixed.Int26_6(31), a)
		assert.True(t, ok)
	})

	t.Run("when no font has glyph", func(t *testing.T) {
		chosen := &mockFace{GlyphFunc: func(p fixed.Point26_6, r rune) (image.Rectangle, image.Image, image.Point, fixed.Int26_6, bool) {
			return image.Rectangle{}, nil, image.Point{}, 0, false
		}}
		fallback := &mockFace{GlyphFunc: func(p fixed.Point26_6, r rune) (image.Rectangle, image.Image, image.Point, fixed.Int26_6, bool) {
			return image.Rectangle{}, nil, image.Point{}, 0, false
		}}
		c := newFontWithFallback(chosen, fallback, chosenFont, fallbackFont)
		_, _, _, _, ok := c.Glyph(fixed.Point26_6{X: 42}, 'x')
		assert.False(t, chosen.GlyphInvoked)
		assert.False(t, fallback.GlyphInvoked)
		assert.False(t, ok)
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

type mockFont struct {
	IndexFunc    func(rune) truetype.Index
	IndexInvoked bool
}

var _ ttfFont = (*mockFont)(nil)

func (f *mockFont) Index(r rune) truetype.Index {
	f.IndexInvoked = true
	return f.IndexFunc(r)
}
