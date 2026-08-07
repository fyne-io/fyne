//go:build windows && directx

package dx

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"testing"
	"unsafe"
)

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
		if x < 0 || y < 0 || x+s.w > atlasSize || y+s.h > atlasSize {
			t.Fatalf("slot %v falls outside the %dx%d atlas", r, atlasSize, atlasSize)
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
			len(placed), atlasSize, atlasSize)
	}
}

func TestAtlasPackerRejectsAndRecovers(t *testing.T) {
	var p atlasPacker
	if _, _, ok := p.add(0, 10); ok {
		t.Error("a zero-width glyph should not be placed")
	}
	if _, _, ok := p.add(atlasSize, atlasSize); ok {
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
	a := newGlyphKey(nil, 42, 12, 1)
	b := newGlyphKey(nil, 42, 24, 1)
	if a == b {
		t.Fatal("keys for two font sizes must differ")
	}
	// Scale folds into the same field: a 12pt glyph at 2x is a 24pt bitmap.
	if newGlyphKey(nil, 42, 12, 2) != b {
		t.Error("fontSize*pixScale should determine the key")
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

// TestGlyphInstMatchesShader pins the Go mirror to the HLSL declaration. A
// mismatch is silent: the shader reads the wrong float4s and every glyph comes
// out at the wrong place, the wrong size or the wrong colour.
func TestGlyphInstMatchesShader(t *testing.T) {
	start := strings.Index(Common, "struct GlyphInst")
	if start < 0 {
		t.Fatal("struct GlyphInst not found in Common")
	}
	end := strings.Index(Common[start:], "};")
	if end < 0 {
		t.Fatal("unterminated struct in Common")
	}
	field := regexp.MustCompile(`(?m)^\s*float4\s+(\w+)\s*;`)
	matches := field.FindAllStringSubmatch(Common[start:start+end], -1)

	typ := reflect.TypeOf(glyphInst{})
	if len(matches) != typ.NumField() {
		t.Fatalf("GlyphInst has %d float4 members, glyphInst has %d fields",
			len(matches), typ.NumField())
	}
	for i, m := range matches {
		want, got := typ.Field(i).Name, m[1]
		if !strings.EqualFold(want, got) {
			t.Errorf("field %d: shader has %q, Go struct has %q", i, got, want)
		}
		if typ.Field(i).Type != reflect.TypeOf([4]float32{}) {
			t.Errorf("field %s must be [4]float32 to match float4 packing", want)
		}
	}

	// D3D11 rejects a constant buffer whose size is not a multiple of 16.
	if size := unsafe.Sizeof(glyphInst{}); size%16 != 0 {
		t.Errorf("glyphInst is %d bytes, must be a multiple of 16", size)
	}
}

// TestGlyphBatchMaxMatchesShader keeps the Go batch cap and the HLSL array in
// step. Too large and DrawInstanced indexes past the array end.
func TestGlyphBatchMaxMatchesShader(t *testing.T) {
	m := regexp.MustCompile(`gGlyphs\[(\d+)\]`).FindStringSubmatch(Common)
	if m == nil {
		t.Fatal("gGlyphs array declaration not found in Common")
	}
	if got := m[1]; got != fmt.Sprint(glyphBatchMax) {
		t.Errorf("shader declares gGlyphs[%s], glyphBatchMax is %d", got, glyphBatchMax)
	}
}
