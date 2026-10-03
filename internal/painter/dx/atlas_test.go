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
