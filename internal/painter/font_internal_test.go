//go:build test

package painter

import "github.com/go-text/typesetting/fontscan"

//
//func Test_compositeFace_Close(t *testing.T) {
//	chosenFont := &truetype.Font{}
//	fallbackFont := &truetype.Font{}
//
//	t.Run("happy case", func(t *testing.T) {
//		chosen := &mockFace{CloseFunc: func() error { return nil }}
//		fallback := &mockFace{CloseFunc: func() error { return nil }}
//		c := newFontWithFallback(chosen, fallback, chosenFont, fallbackFont)
//		if assert.NoError(t, c.Close()) {
//			assert.True(t, chosen.CloseInvoked)
//			assert.True(t, fallback.CloseInvoked)
//		}
//	})
//
//	t.Run("with Fallback failing", func(t *testing.T) {
//		chosen := &mockFace{CloseFunc: func() error { return nil }}
//		e := errors.New("oops")
//		fallback := &mockFace{CloseFunc: func() error { return e }}
//		c := newFontWithFallback(chosen, fallback, chosenFont, fallbackFont)
//		assert.ErrorIs(t, c.Close(), e)
//		assert.True(t, chosen.CloseInvoked)
//		assert.True(t, fallback.CloseInvoked)
//	})
//
//	t.Run("with primary failing", func(t *testing.T) {
//		e := errors.New("oops")
//		chosen := &mockFace{CloseFunc: func() error { return e }}
//		fallback := &mockFace{CloseFunc: func() error { return nil }}
//		c := newFontWithFallback(chosen, fallback, chosenFont, fallbackFont)
//		assert.ErrorIs(t, c.Close(), e)
//		assert.True(t, chosen.CloseInvoked)
//		assert.True(t, fallback.CloseInvoked)
//	})
//}
//
//func Test_compositeFace_GlyphFunctions(t *testing.T) {
//	chosenFont := &mockFont{IndexFunc: func(r rune) truetype.Index {
//		if r == 'e' || r == 'p' {
//			return 1
//		}
//		return 0
//	}}
//	fallbackFont := &mockFont{IndexFunc: func(r rune) truetype.Index {
//		if r == 'e' || r == 's' {
//			return 1
//		}
//		return 0
//	}}
//
//	t.Run("when primary has glyph", func(t *testing.T) {
//		chosen := &mockFace{
//			GlyphFunc: func(p fixed.Point26_6, r rune) (image.Rectangle, image.Image, image.Point, fixed.Int26_6, bool) {
//				assert.Equal(t, fixed.Point26_6{X: 13}, p)
//				assert.Equal(t, 'e', r)
//				return image.Rect(1, 2, 3, 4), image.NewUniform(color.NRGBA{R: 5, G: 6, B: 7, A: 8}), image.Pt(9, 10), 11, true
//			},
//			GlyphAdvanceFunc: func(r rune) (advance fixed.Int26_6, ok bool) {
//				assert.Equal(t, 'e', r)
//				return 11, true
//			},
//			GlyphBoundsFunc: func(r rune) (fixed.Rectangle26_6, fixed.Int26_6, bool) {
//				assert.Equal(t, 'e', r)
//				return fixed.Rectangle26_6{Min: fixed.Point26_6{X: 12, Y: 13}, Max: fixed.Point26_6{X: 14, Y: 15}}, 11, true
//			},
//		}
//		fallback := &mockFace{
//			GlyphFunc: func(p fixed.Point26_6, r rune) (image.Rectangle, image.Image, image.Point, fixed.Int26_6, bool) {
//				return image.Rectangle{}, nil, image.Point{}, 0, false
//			},
//			GlyphAdvanceFunc: func(r rune) (advance fixed.Int26_6, ok bool) {
//				return 0, false
//			},
//			GlyphBoundsFunc: func(r rune) (fixed.Rectangle26_6, fixed.Int26_6, bool) {
//				return fixed.Rectangle26_6{}, 0, false
//			},
//		}
//		c := newFontWithFallback(chosen, fallback, chosenFont, fallbackFont)
//		t.Run("Glyph", func(t *testing.T) {
//			r, m, mp, a, ok := c.Glyph(fixed.Point26_6{X: 13}, 'e')
//			assert.True(t, chosen.GlyphInvoked)
//			assert.False(t, fallback.GlyphInvoked)
//			assert.Equal(t, image.Rect(1, 2, 3, 4), r)
//			assert.Equal(t, image.NewUniform(color.NRGBA{R: 5, G: 6, B: 7, A: 8}), m)
//			assert.Equal(t, image.Pt(9, 10), mp)
//			assert.Equal(t, fixed.Int26_6(11), a)
//			assert.True(t, ok)
//		})
//		t.Run("GlyphAdvance", func(t *testing.T) {
//			a, ok := c.GlyphAdvance('e')
//			assert.True(t, chosen.GlyphAdvanceInvoked)
//			assert.False(t, fallback.GlyphAdvanceInvoked)
//			assert.Equal(t, fixed.Int26_6(11), a)
//			assert.True(t, ok)
//		})
//		t.Run("GlyphBounds", func(t *testing.T) {
//			b, a, ok := c.GlyphBounds('e')
//			assert.True(t, chosen.GlyphBoundsInvoked)
//			assert.False(t, fallback.GlyphBoundsInvoked)
//			assert.Equal(t, fixed.Rectangle26_6{Min: fixed.Point26_6{X: 12, Y: 13}, Max: fixed.Point26_6{X: 14, Y: 15}}, b)
//			assert.Equal(t, fixed.Int26_6(11), a)
//			assert.True(t, ok)
//		})
//	})
//
//	t.Run("when only Fallback has glyph", func(t *testing.T) {
//		chosen := &mockFace{
//			GlyphFunc: func(p fixed.Point26_6, r rune) (image.Rectangle, image.Image, image.Point, fixed.Int26_6, bool) {
//				return image.Rectangle{}, nil, image.Point{}, 0, false
//			},
//			GlyphAdvanceFunc: func(r rune) (advance fixed.Int26_6, ok bool) {
//				return 0, false
//			},
//			GlyphBoundsFunc: func(r rune) (fixed.Rectangle26_6, fixed.Int26_6, bool) {
//				return fixed.Rectangle26_6{}, 0, false
//			},
//		}
//		fallback := &mockFace{
//			GlyphFunc: func(p fixed.Point26_6, r rune) (image.Rectangle, image.Image, image.Point, fixed.Int26_6, bool) {
//				assert.Equal(t, fixed.Point26_6{X: 42}, p)
//				assert.Equal(t, 's', r)
//				return image.Rect(21, 22, 23, 24), image.NewUniform(color.NRGBA{R: 25, G: 26, B: 27, A: 28}), image.Pt(29, 30), 31, true
//			},
//			GlyphAdvanceFunc: func(r rune) (advance fixed.Int26_6, ok bool) {
//				assert.Equal(t, 's', r)
//				return 31, true
//			},
//			GlyphBoundsFunc: func(r rune) (fixed.Rectangle26_6, fixed.Int26_6, bool) {
//				assert.Equal(t, 's', r)
//				return fixed.Rectangle26_6{Min: fixed.Point26_6{X: 32, Y: 33}, Max: fixed.Point26_6{X: 34, Y: 35}}, 31, true
//			},
//		}
//		c := newFontWithFallback(chosen, fallback, chosenFont, fallbackFont)
//		t.Run("Glyph", func(t *testing.T) {
//			r, m, mp, a, ok := c.Glyph(fixed.Point26_6{X: 42}, 's')
//			assert.False(t, chosen.GlyphInvoked)
//			assert.True(t, fallback.GlyphInvoked)
//			assert.Equal(t, image.Rect(21, 22, 23, 24), r)
//			assert.Equal(t, image.NewUniform(color.NRGBA{R: 25, G: 26, B: 27, A: 28}), m)
//			assert.Equal(t, image.Pt(29, 30), mp)
//			assert.Equal(t, fixed.Int26_6(31), a)
//			assert.True(t, ok)
//		})
//		t.Run("GlyphAdvance", func(t *testing.T) {
//			a, ok := c.GlyphAdvance('s')
//			assert.False(t, chosen.GlyphAdvanceInvoked)
//			assert.True(t, fallback.GlyphAdvanceInvoked)
//			assert.Equal(t, fixed.Int26_6(31), a)
//			assert.True(t, ok)
//		})
//		t.Run("GlyphBounds", func(t *testing.T) {
//			b, a, ok := c.GlyphBounds('s')
//			assert.False(t, chosen.GlyphBoundsInvoked)
//			assert.True(t, fallback.GlyphBoundsInvoked)
//			assert.Equal(t, fixed.Rectangle26_6{Min: fixed.Point26_6{X: 32, Y: 33}, Max: fixed.Point26_6{X: 34, Y: 35}}, b)
//			assert.Equal(t, fixed.Int26_6(31), a)
//			assert.True(t, ok)
//		})
//	})
//
//	t.Run("when no Font has glyph", func(t *testing.T) {
//		chosen := &mockFace{
//			GlyphFunc: func(p fixed.Point26_6, r rune) (image.Rectangle, image.Image, image.Point, fixed.Int26_6, bool) {
//				return image.Rectangle{}, nil, image.Point{}, 0, false
//			},
//			GlyphAdvanceFunc: func(r rune) (advance fixed.Int26_6, ok bool) {
//				return 0, false
//			},
//			GlyphBoundsFunc: func(r rune) (fixed.Rectangle26_6, fixed.Int26_6, bool) {
//				return fixed.Rectangle26_6{}, 0, false
//			},
//		}
//		fallback := &mockFace{
//			GlyphFunc: func(p fixed.Point26_6, r rune) (image.Rectangle, image.Image, image.Point, fixed.Int26_6, bool) {
//				return image.Rectangle{}, nil, image.Point{}, 0, false
//			},
//			GlyphAdvanceFunc: func(r rune) (advance fixed.Int26_6, ok bool) {
//				return 0, false
//			},
//			GlyphBoundsFunc: func(r rune) (fixed.Rectangle26_6, fixed.Int26_6, bool) {
//				return fixed.Rectangle26_6{}, 0, false
//			},
//		}
//		c := newFontWithFallback(chosen, fallback, chosenFont, fallbackFont)
//		t.Run("Glyph", func(t *testing.T) {
//			_, _, _, _, ok := c.Glyph(fixed.Point26_6{X: 42}, 'x')
//			assert.False(t, chosen.GlyphInvoked)
//			assert.False(t, fallback.GlyphInvoked)
//			assert.False(t, ok)
//		})
//		t.Run("GlyphAdvance", func(t *testing.T) {
//			_, ok := c.GlyphAdvance('x')
//			assert.False(t, chosen.GlyphAdvanceInvoked)
//			assert.False(t, fallback.GlyphAdvanceInvoked)
//			assert.False(t, ok)
//		})
//		t.Run("GlyphBounds", func(t *testing.T) {
//			_, _, ok := c.GlyphBounds('x')
//			assert.False(t, chosen.GlyphBoundsInvoked)
//			assert.False(t, fallback.GlyphBoundsInvoked)
//			assert.False(t, ok)
//		})
//	})
//}
//
//func Test_compositeFace_Kern(t *testing.T) {
//	chosenFont := &mockFont{IndexFunc: func(r rune) truetype.Index {
//		if r == 'e' || r == 'p' {
//			return 1
//		}
//		return 0
//	}}
//	fallbackFont := &mockFont{IndexFunc: func(r rune) truetype.Index {
//		if r == 'e' || r == 's' {
//			return 1
//		}
//		return 0
//	}}
//
//	t.Run("when primary has both glyphs", func(t *testing.T) {
//		chosen := &mockFace{KernFunc: func(r1, r2 rune) fixed.Int26_6 {
//			assert.Equal(t, 'e', r1)
//			assert.Equal(t, 'p', r2)
//			return 1
//		}}
//		fallback := &mockFace{KernFunc: func(r1, r2 rune) fixed.Int26_6 { return 0 }}
//		c := newFontWithFallback(chosen, fallback, chosenFont, fallbackFont)
//		k := c.Kern('e', 'p')
//		assert.True(t, chosen.KernInvoked)
//		assert.False(t, fallback.KernInvoked)
//		assert.Equal(t, fixed.Int26_6(1), k)
//	})
//
//	t.Run("when primary misses first glyph", func(t *testing.T) {
//		chosen := &mockFace{KernFunc: func(r1, r2 rune) fixed.Int26_6 { return 0 }}
//		fallback := &mockFace{KernFunc: func(r1, r2 rune) fixed.Int26_6 {
//			assert.Equal(t, 's', r1)
//			assert.Equal(t, 'e', r2)
//			return 2
//		}}
//		c := newFontWithFallback(chosen, fallback, chosenFont, fallbackFont)
//		k := c.Kern('s', 'e')
//		assert.False(t, chosen.KernInvoked)
//		assert.True(t, fallback.KernInvoked)
//		assert.Equal(t, fixed.Int26_6(2), k)
//	})
//
//	t.Run("when primary misses second glyph", func(t *testing.T) {
//		chosen := &mockFace{KernFunc: func(r1, r2 rune) fixed.Int26_6 { return 0 }}
//		fallback := &mockFace{KernFunc: func(r1, r2 rune) fixed.Int26_6 {
//			assert.Equal(t, 'e', r1)
//			assert.Equal(t, 's', r2)
//			return 2
//		}}
//		c := newFontWithFallback(chosen, fallback, chosenFont, fallbackFont)
//		k := c.Kern('e', 's')
//		assert.False(t, chosen.KernInvoked)
//		assert.True(t, fallback.KernInvoked)
//		assert.Equal(t, fixed.Int26_6(2), k)
//	})
//}
//
//func Test_compositeFace_Metrics(t *testing.T) {
//	chosenFont := &truetype.Font{}
//	fallbackFont := &truetype.Font{}
//	chosen := &mockFace{MetricsFunc: func() font.Metrics {
//		return font.Metrics{Height: 1, Ascent: 2, Descent: 3, XHeight: 4, CapHeight: 5, CaretSlope: image.Pt(6, 7)}
//	}}
//	fallback := &mockFace{MetricsFunc: func() font.Metrics { return font.Metrics{} }}
//	c := newFontWithFallback(chosen, fallback, chosenFont, fallbackFont)
//	m := c.Metrics()
//	assert.True(t, chosen.MetricsInvoked)
//	assert.False(t, fallback.MetricsInvoked)
//	assert.Equal(t, font.Metrics{Height: 1, Ascent: 2, Descent: 3, XHeight: 4, CapHeight: 5, CaretSlope: image.Pt(6, 7)}, m)
//}
//
//type mockFace struct {
//	CloseFunc           func() error
//	CloseInvoked        bool
//	GlyphFunc           func(fixed.Point26_6, rune) (image.Rectangle, image.Image, image.Point, fixed.Int26_6, bool)
//	GlyphInvoked        bool
//	GlyphAdvanceFunc    func(rune) (fixed.Int26_6, bool)
//	GlyphAdvanceInvoked bool
//	GlyphBoundsFunc     func(rune) (fixed.Rectangle26_6, fixed.Int26_6, bool)
//	GlyphBoundsInvoked  bool
//	KernFunc            func(rune, rune) fixed.Int26_6
//	KernInvoked         bool
//	MetricsFunc         func() font.Metrics
//	MetricsInvoked      bool
//}
//
//var _ font.Face = (*mockFace)(nil)
//
//func (f *mockFace) Close() error {
//	f.CloseInvoked = true
//	return f.CloseFunc()
//}
//
//func (f *mockFace) Glyph(dot fixed.Point26_6, r rune) (dr image.Rectangle, mask image.Image, maskp image.Point, advance fixed.Int26_6, ok bool) {
//	f.GlyphInvoked = true
//	return f.GlyphFunc(dot, r)
//}
//
//func (f *mockFace) GlyphAdvance(r rune) (advance fixed.Int26_6, ok bool) {
//	f.GlyphAdvanceInvoked = true
//	return f.GlyphAdvanceFunc(r)
//}
//
//func (f *mockFace) GlyphBounds(r rune) (bounds fixed.Rectangle26_6, advance fixed.Int26_6, ok bool) {
//	f.GlyphBoundsInvoked = true
//	return f.GlyphBoundsFunc(r)
//}
//
//func (f *mockFace) Kern(r0, r1 rune) fixed.Int26_6 {
//	f.KernInvoked = true
//	return f.KernFunc(r0, r1)
//}
//
//func (f *mockFace) Metrics() font.Metrics {
//	f.MetricsInvoked = true
//	return f.MetricsFunc()
//}
//
//type mockFont struct {
//	IndexFunc    func(rune) truetype.Index
//	IndexInvoked bool
//}
//
//var _ ttfFont = (*mockFont)(nil)
//
//func (f *mockFont) Index(r rune) truetype.Index {
//	f.IndexInvoked = true
//	return f.IndexFunc(r)
//}

func loadSystemFonts(fm *fontscan.FontMap) error {
	return nil
}
