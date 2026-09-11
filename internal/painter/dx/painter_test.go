//go:build windows && directx

package dx

import (
	"fmt"
	"image/color"
	"reflect"
	"regexp"
	"strings"
	"testing"
	"unsafe"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
)

// The constant buffer is the one place where a silent mistake corrupts every
// single draw without any error: HLSL reads whatever bytes sit at each offset.
// These checks pin the Go mirror to the shader declaration.

var cbufferField = regexp.MustCompile(`(?m)^\s*float4\s+(\w+)\s*;`)

func TestConstantsMatchShaderCBuffer(t *testing.T) {
	start := strings.Index(Common, "cbuffer Constants")
	if start < 0 {
		t.Fatal("cbuffer Constants not found in Common")
	}
	end := strings.Index(Common[start:], "};")
	if end < 0 {
		t.Fatal("unterminated cbuffer in Common")
	}
	matches := cbufferField.FindAllStringSubmatch(Common[start:start+end], -1)

	typ := reflect.TypeOf(constants{})
	if len(matches) != typ.NumField() {
		t.Fatalf("cbuffer has %d float4 members, constants has %d fields",
			len(matches), typ.NumField())
	}

	for i, m := range matches {
		want := typ.Field(i).Name
		got := m[1]
		if !strings.EqualFold(want, got) {
			t.Errorf("field %d: shader has %q, Go struct has %q", i, got, want)
		}
		if typ.Field(i).Type != reflect.TypeOf([4]float32{}) {
			t.Errorf("field %s must be [4]float32 to match float4 packing", want)
		}
	}
}

func TestConstantBufferSizeIsAligned(t *testing.T) {
	// D3D11 rejects CreateBuffer for a constant buffer whose size is not a
	// multiple of 16 bytes.
	size := unsafe.Sizeof(constants{})
	if size%16 != 0 {
		t.Errorf("constants is %d bytes, must be a multiple of 16", size)
	}
}

func TestRoundToPixel(t *testing.T) {
	cases := []struct {
		v, scale, want float32
	}{
		{1.4, 1.0, 1},
		{1.5, 1.0, 2},
		{1.25, 2.0, 1.5}, // 2.5 rounds to 3, /2 = 1.5
		{-1.5, 1.0, -2},
	}
	for _, c := range cases {
		if got := roundToPixel(c.v, c.scale); got != c.want {
			t.Errorf("roundToPixel(%v, %v) = %v, want %v", c.v, c.scale, got, c.want)
		}
	}
}

// TestVecRectCoords checks the NDC conversion against hand-computed values: a
// 100x50 rect at (10,20) inside a 200x100 frame, with edge softness folded in.
func TestVecRectCoords(t *testing.T) {
	p := &Painter{pixScale: 1, texScale: 1}
	obj := &testObject{size: fyne.NewSize(100, 50), pos: fyne.NewPos(10, 20)}
	frame := fyne.NewSize(200, 100)

	coords, bounds := p.vecRectCoords(fyne.NewPos(10, 20), obj, frame, 0, canvas.Shadow{})

	if bounds != [4]float32{10, 20, 110, 70} {
		t.Fatalf("bounds = %v, want [10 20 110 70]", bounds)
	}

	// Edge softness is 0.5pt but roundToPixel rounds half away from zero, so the
	// quad is inflated by a whole pixel on each side:
	//   x1 = -1 + (10-1)*2/200 = -0.91,  y1 = 1 - (20-1)*2/100 = 0.62
	const wantX1, wantY1 = -0.91, 0.62
	if !closeTo(coords[0], wantX1) || !closeTo(coords[1], wantY1) {
		t.Errorf("first corner = (%v, %v), want (%v, %v)", coords[0], coords[1], wantX1, wantY1)
	}
	// coords is the clip-space rect {x1, y1, x2, y2} the vertex shader expands.
	if coords[2] <= coords[0] {
		t.Error("x2 should be right of x1 in NDC")
	}
	if coords[1] <= coords[3] {
		t.Error("y1 should be above y2 in NDC (larger value)")
	}
}

func TestFragmentColorUnpremultiplies(t *testing.T) {
	// The shaders expect straight alpha, but color.RGBA() returns premultiplied
	// components, so fragmentColor has to divide the alpha back out.
	r, g, b, a := fragmentColor(color.NRGBA{R: 0xff, G: 0x00, B: 0x00, A: 0x80})
	if !closeTo(r, 1) || !closeTo(g, 0) || !closeTo(b, 0) {
		t.Errorf("rgb = (%v, %v, %v), want full red", r, g, b)
	}
	if !closeTo(a, float32(0x8080)/0xffff) {
		t.Errorf("alpha = %v, want ~0.502", a)
	}

	if _, _, _, a := fragmentColor(color.Transparent); a != 0 {
		t.Errorf("transparent alpha = %v, want 0", a)
	}
	if _, _, _, a := fragmentColor(nil); a != 0 {
		t.Errorf("nil colour alpha = %v, want 0", a)
	}
}

func closeTo(a, b float32) bool {
	d := a - b
	return d < 1e-4 && d > -1e-4
}

type testObject struct {
	fyne.CanvasObject
	size fyne.Size
	pos  fyne.Position
}

func (t *testObject) Size() fyne.Size { return t.size }

func (t *testObject) Position() fyne.Position { return t.pos }

// TestShadersCompile runs every shader in internal/painter/dx/shaders through the
// real compiler. This is the only check that the ported GLSL is valid HLSL - a
// typo there would otherwise only show up as a blank window at runtime. It also
// covers the shared prelude, which every source here has prepended.
func TestShadersCompile(t *testing.T) {
	if err := procD3DCompile.Find(); err != nil {
		t.Skip("d3dcompiler_47.dll unavailable:", err)
	}

	shaders := []struct {
		name, src, target string
	}{
		{"passthrough_2d.vs.hlsl", VertexPassthrough2D, "vs_4_0"},
		{"line.vs.hlsl", VertexLine, "vs_4_0"},
		{"simple.ps.hlsl", PixelSimple, "ps_4_0"},
		{"rectangle.ps.hlsl", PixelRectangle, "ps_4_0"},
		{"round_rectangle.ps.hlsl", PixelRoundRectangle, "ps_4_0"},
		{"ellipse.ps.hlsl", PixelEllipse, "ps_4_0"},
		{"arc.ps.hlsl", PixelArc, "ps_4_0"},
		{"regular_polygon.ps.hlsl", PixelRegularPolygon, "ps_4_0"},
		{"arbitrary_polygon.ps.hlsl", PixelArbitraryPolygon, "ps_4_0"},
		{"bezier_curve.ps.hlsl", PixelBezierCurve, "ps_4_0"},
		{"blur.ps.hlsl", PixelBlur, "ps_4_0"},
		{"line.ps.hlsl", PixelLine, "ps_4_0"},
	}
	for _, s := range shaders {
		t.Run(s.name, func(t *testing.T) {
			code, err := compileShader(s.src, "main", s.target)
			if err != nil {
				t.Fatalf("compile failed: %v", err)
			}
			if len(code) == 0 {
				t.Error("compiler returned empty bytecode")
			}
		})
	}
}

// TestBlendModes pins the two blend equations. Textures come from Go's
// alpha-premultiplied image.RGBA, so they must blend with ONE; using SRC_ALPHA
// applies alpha a second time and quietly destroys anti-aliased glyph edges.
// Nothing errors when this is wrong, so only a test catches it.
func TestBlendModes(t *testing.T) {
	shape := straightAlphaBlend().RenderTarget[0]
	if shape.SrcBlend != blendSrcAlpha || shape.SrcBlendAlpha != blendSrcAlpha {
		t.Error("shapes emit straight alpha and must blend with SRC_ALPHA")
	}

	tex := premultipliedAlphaBlend().RenderTarget[0]
	if tex.SrcBlend != blendOne || tex.SrcBlendAlpha != blendOne {
		t.Error("textures are premultiplied and must blend with ONE, not SRC_ALPHA")
	}

	for name, rt := range map[string]renderTargetBlendDesc{"shape": shape, "texture": tex} {
		if rt.BlendEnable != 1 {
			t.Errorf("%s blending is disabled", name)
		}
		if rt.DestBlend != blendInvSrcAlpha || rt.DestBlendAlpha != blendInvSrcAlpha {
			t.Errorf("%s must composite over the destination with INV_SRC_ALPHA", name)
		}
		if rt.BlendOp != blendOpAdd || rt.BlendOpAlpha != blendOpAdd {
			t.Errorf("%s blend op should be ADD", name)
		}
		if rt.RenderTargetWriteMask != colorWriteEnableAll {
			t.Errorf("%s does not write all channels", name)
		}
	}
}

// TestAuxConstantsMatchShaders pins the register(b1) buffer the polygon and
// curve shaders read. Same hazard as the shared cbuffer: a size or ordering
// mismatch is read as garbage offsets rather than reported as an error.
func TestAuxConstantsMatchShaders(t *testing.T) {
	if size := unsafe.Sizeof(auxConstants{}); size%16 != 0 {
		t.Errorf("auxConstants is %d bytes, must be a multiple of 16", size)
	}

	// MAX_VERTICES in arbitrary_polygon.ps.hlsl must match the Go array.
	want := regexp.MustCompile(`#define\s+MAX_VERTICES\s+(\d+)`).
		FindStringSubmatch(PixelArbitraryPolygon)
	if want == nil {
		t.Fatal("MAX_VERTICES not found in arbitrary_polygon.ps.hlsl")
	}
	if got := want[1]; got != "16" || maxPolygonVertices != 16 {
		t.Errorf("shader MAX_VERTICES is %s, Go maxPolygonVertices is %d",
			got, maxPolygonVertices)
	}

	// The curve shader overlays its two float4s on the first polygon vertices.
	if unsafe.Offsetof(auxConstants{}.Info) != maxPolygonVertices*16 {
		t.Error("Info must follow the full vertex array, as polyInfo does in HLSL")
	}
}

// TestBlurKernelIsNormalised checks the Gaussian weights sum to one; if they did
// not, blurring would visibly brighten or darken what it blurs.
func TestBlurKernelIsNormalised(t *testing.T) {
	for _, radius := range []float32{1, 4, 12.5, 50} {
		values := createBlurKernel(radius)
		if len(values) != int(radius)*2+1 {
			t.Errorf("radius %v: got %d weights, want %d", radius, len(values), int(radius)*2+1)
		}
		var sum float32
		for _, v := range values {
			sum += v
		}
		if !closeTo(sum, 1) {
			t.Errorf("radius %v: weights sum to %v, want 1", radius, sum)
		}
	}
}

// TestTextTextureWindow pins the windowing maths for over-wide text runs:
// the window must cover the visible span and stay inside the full run.
func TestTextTextureWindow(t *testing.T) {
	cases := []struct {
		visibleOffset, visibleWidth, fullWidth, maxWidth int
		wantOffset, wantWidth                            int
	}{
		{0, 400, 20000, 4096, 0, 4096},
		{5000, 400, 20000, 4096, 3152, 4096},
		{19800, 200, 20000, 4096, 15904, 4096},
	}
	for _, c := range cases {
		offset, width := textTextureWindow(c.visibleOffset, c.visibleWidth, c.fullWidth, c.maxWidth)
		if offset != c.wantOffset || width != c.wantWidth {
			t.Errorf("textTextureWindow(%d, %d, %d, %d) = %d, %d, want %d, %d",
				c.visibleOffset, c.visibleWidth, c.fullWidth, c.maxWidth,
				offset, width, c.wantOffset, c.wantWidth)
		}
	}

	cached := clippedTextEntry{offset: 3152, width: 4096, height: 40, scale: 1}
	if !cached.covers(5000, 400, 40, 1) {
		t.Error("cached window should cover a range inside it")
	}
	if cached.covers(7200, 400, 40, 1) {
		t.Error("cached window should not cover a range beyond it")
	}
}

// TestFrameLatencySlotDerivation pins the SetMaximumFrameLatency vtable index to
// the interface chain it is counted from. An off-by-one here does not fail - it
// calls a neighbouring method with the wrong argument, on a COM object, in a
// process that then keeps running.
func TestFrameLatencySlotDerivation(t *testing.T) {
	const (
		iUnknown    = 3 // QueryInterface, AddRef, Release
		iDXGIObject = 4 // SetPrivateData, SetPrivateDataInterface, GetPrivateData, GetParent
		iDXGIDevice = 5 // GetAdapter, CreateSurface, QueryResourceResidency,
		// SetGPUThreadPriority, GetGPUThreadPriority
	)

	// SetMaximumFrameLatency is the first method IDXGIDevice1 adds of its own.
	if want := iUnknown + iDXGIObject + iDXGIDevice; slotSetMaximumFrameLatency != want {
		t.Errorf("slotSetMaximumFrameLatency = %d, want %d", slotSetMaximumFrameLatency, want)
	}
}

// TestDXGIDevice1IID checks the GUID byte layout against the canonical form. A
// typo here is silent in a different way: QueryInterface just reports
// E_NOINTERFACE, setMaximumFrameLatency logs and carries on, and the frame queue
// quietly stays at its laggy default.
func TestDXGIDevice1IID(t *testing.T) {
	// {77db970f-6276-48ba-ba28-070143b4392c}
	if got := fmt.Sprintf("%08x-%04x-%04x-%02x-%02x", iidDXGIDevice1.Data1, iidDXGIDevice1.Data2,
		iidDXGIDevice1.Data3, iidDXGIDevice1.Data4[:2], iidDXGIDevice1.Data4[2:]); got != "77db970f-6276-48ba-ba28-070143b4392c" {
		t.Errorf("iidDXGIDevice1 = %s, want 77db970f-6276-48ba-ba28-070143b4392c", got)
	}
}

// poolPainter is a Painter with just enough state for the texture pool; the
// parked textures carry nil COM pointers, which release() ignores.
func poolPainter() *Painter {
	return &Painter{texPool: map[[2]uint32][]pooledTexture{}}
}

func TestTexturePool_ParkAndReuse(t *testing.T) {
	p := poolPainter()
	tex := &gpuTexture{width: 8, height: 4}

	p.releaseTexture(tex)
	if p.texPoolSize != 1 {
		t.Fatalf("texPoolSize = %d, want 1", p.texPoolSize)
	}
	if got := p.pooledTextureFor(8, 8); got != nil {
		t.Errorf("pooledTextureFor(8, 8) = %v, want nil for a size miss", got)
	}
	if got := p.pooledTextureFor(8, 4); got != tex {
		t.Errorf("pooledTextureFor(8, 4) = %v, want the parked texture", got)
	}
	if p.texPoolSize != 0 || len(p.texPool) != 0 {
		t.Errorf("after reuse: texPoolSize = %d, entries = %d, want an empty pool", p.texPoolSize, len(p.texPool))
	}
}

func TestTexturePool_LIFOAndUnpoolable(t *testing.T) {
	p := poolPainter()
	older := &gpuTexture{width: 8, height: 4}
	newer := &gpuTexture{width: 8, height: 4}
	p.releaseTexture(older)
	p.releaseTexture(newer)
	p.releaseTexture(&gpuTexture{}) // width 0: not from imgToTexture, never pooled

	if p.texPoolSize != 2 {
		t.Fatalf("texPoolSize = %d, want 2 (unpoolable texture parked?)", p.texPoolSize)
	}
	if got := p.pooledTextureFor(8, 4); got != newer {
		t.Errorf("first pop = %v, want the most recently parked", got)
	}
	if got := p.pooledTextureFor(8, 4); got != older {
		t.Errorf("second pop = %v, want the older texture", got)
	}
}

func TestTexturePool_CapAndSweep(t *testing.T) {
	p := poolPainter()
	for i := 0; i < texPoolMax+5; i++ {
		p.releaseTexture(&gpuTexture{width: uint32(i + 1), height: 1})
	}
	if p.texPoolSize != texPoolMax {
		t.Fatalf("texPoolSize = %d, want the cap %d", p.texPoolSize, texPoolMax)
	}

	p.frameTick += texPoolTTL + 1
	p.sweepTexPool()
	if p.texPoolSize != 0 || len(p.texPool) != 0 {
		t.Errorf("after sweep: texPoolSize = %d, entries = %d, want an empty pool", p.texPoolSize, len(p.texPool))
	}

	// A freshly parked texture must survive a sweep.
	p.releaseTexture(&gpuTexture{width: 8, height: 4})
	p.sweepTexPool()
	if p.texPoolSize != 1 {
		t.Errorf("fresh texture swept: texPoolSize = %d, want 1", p.texPoolSize)
	}
}
