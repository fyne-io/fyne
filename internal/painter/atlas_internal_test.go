package painter

import "testing"

// rect is a placed slot, used to check that nothing the packer hands out
// overlaps anything else it has handed out.
type rect struct{ x, y, w, h int }

func (r rect) overlaps(o rect) bool {
	return r.x < o.x+o.w && o.x < r.x+r.w && r.y < o.y+o.h && o.y < r.y+r.h
}

// TestAtlasPackerNoOverlap is the property that matters: two glyphs must never
// share a pixel, or one draws part of the other. Overlap is invisible in a
// screenshot until it is not, so it is checked exhaustively here.
func TestAtlasPackerNoOverlap(t *testing.T) {
	var p atlasPacker
	var placed []rect

	// Sizes chosen to force several shelves and a mid-shelf height change,
	// which is where a shelf packer goes wrong if it goes wrong at all.
	sizes := []struct{ w, h int }{{8, 12}, {30, 9}, {5, 20}, {17, 20}, {40, 6}, {9, 13}}
	for i := 0; i < 500; i++ {
		s := sizes[i%len(sizes)]
		x, y, ok := p.add(s.w, s.h)
		if !ok {
			break
		}
		r := rect{x, y, s.w, s.h}
		if x < 0 || y < 0 || x+s.w > AtlasSize || y+s.h > AtlasSize {
			t.Fatalf("slot %v falls outside the %dx%d atlas", r, AtlasSize, AtlasSize)
		}
		for _, o := range placed {
			if r.overlaps(o) {
				t.Fatalf("slot %v overlaps %v", r, o)
			}
		}
		placed = append(placed, r)
	}

	if len(placed) < 100 {
		t.Errorf("only packed %d glyphs into %dx%d, expected the atlas to hold far more",
			len(placed), AtlasSize, AtlasSize)
	}
}

func TestAtlasPackerRejectsAndRecovers(t *testing.T) {
	var p atlasPacker
	if _, _, ok := p.add(0, 10); ok {
		t.Error("a zero-width glyph should not be placed")
	}
	if _, _, ok := p.add(AtlasSize, AtlasSize); ok {
		t.Error("a glyph needing the whole atlas leaves no room for padding")
	}

	// Fill it, then check reset makes it usable again - that is the response to
	// a full atlas, so it has to actually work.
	for {
		if _, _, ok := p.add(64, 64); !ok {
			break
		}
	}
	p.reset()
	if x, y, ok := p.add(64, 64); !ok || x != 0 || y != 0 {
		t.Errorf("after reset got (%d,%d,%v), want (0,0,true)", x, y, ok)
	}
}

func TestGlyphKeyDistinguishesSize(t *testing.T) {
	// Body text and a heading share glyph IDs, so a key that ignored size would
	// draw the heading at body size. Same face and GID, different size only.
	a := newGlyphKey(nil, 42, 12, 1, 0)
	b := newGlyphKey(nil, 42, 24, 1, 0)
	if a == b {
		t.Fatal("keys for two font sizes must differ")
	}
	// Scale folds into the same field: a 12pt glyph at 2x is a 24pt bitmap.
	if newGlyphKey(nil, 42, 12, 2, 0) != b {
		t.Error("fontSize*pixScale should determine the key")
	}
	// The same glyph rasterised at another sub-pixel offset is another bitmap.
	if newGlyphKey(nil, 42, 12, 1, 16) == a {
		t.Error("keys for two sub-pixel offsets must differ")
	}
}

func TestSnapX(t *testing.T) {
	for _, tc := range []struct {
		x          float32
		steps      int
		pixel, sub int
	}{
		{10, 4, 10, 0},
		{10.2, 4, 10, 16}, // nearest quarter is 0.25
		{10.4, 4, 10, 32}, // 0.5
		{10.9, 4, 11, 0},  // rounds up past the last quarter into the next pixel
		{10.4, 2, 10, 32}, // halves at 2x
		{10.2, 2, 10, 0},  // ...where 0.2 is nearer the whole pixel
		{10.7, 1, 11, 0},  // one step is plain rounding
		{-0.2, 4, -1, 48}, // floors below zero: -0.25 is pixel -1 plus 0.75
		{-1.1, 4, -1, 0},  // ...and -1.0 is the nearest quarter to -1.1
	} {
		pixel, sub := snapX(tc.x, tc.steps)
		if pixel != tc.pixel || sub != tc.sub {
			t.Errorf("snapX(%v, %d) = %d, %d; want %d, %d", tc.x, tc.steps, pixel, sub, tc.pixel, tc.sub)
		}
	}
}

func TestSubpixelSteps(t *testing.T) {
	for _, tc := range []struct {
		scale float32
		steps int
	}{{1, 4}, {1.25, 4}, {1.5, 3}, {2, 2}, {3, 2}, {4, 1}, {6, 1}} {
		if got := subpixelSteps(tc.scale); got != tc.steps {
			t.Errorf("subpixelSteps(%v) = %d, want %d", tc.scale, got, tc.steps)
		}
		// However many steps, the worst snap stays within an eighth of a logical
		// pixel of where the shaper put the glyph.
		if worst := 1 / (2 * float32(tc.steps) * tc.scale); worst > 0.125 {
			t.Errorf("at %vx, %d steps can misplace a glyph by %v logical pixels", tc.scale, tc.steps, worst)
		}
	}
}

func TestGlyphEntryEmpty(t *testing.T) {
	// A space has no ink and must not produce a quad; a real glyph must.
	if !(glyphEntry{width: 0, height: 8}).empty() || !(glyphEntry{width: 8, height: 0}).empty() {
		t.Error("a zero dimension means no ink")
	}
	if (glyphEntry{width: 8, height: 8}).empty() {
		t.Error("a glyph with both dimensions is not empty")
	}
}
