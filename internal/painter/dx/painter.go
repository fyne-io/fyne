//go:build windows

package dx

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"log"
	"math"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/internal"
	"fyne.io/fyne/v2/internal/cache"
	paint "fyne.io/fyne/v2/internal/painter"
	"fyne.io/fyne/v2/internal/painter/geom"
	"fyne.io/fyne/v2/theme"
)

func init() {
	log.SetFlags(log.Lshortfile | log.Lmicroseconds)
}

const (
	edgeSoftness = 0.5
	max16bit     = float32(255 * 256)

	usageStaging  = 3
	cpuAccessRead = 0x20000
	mapRead       = 1
)

// constants mirrors the `Constants` cbuffer in shaders.go. Every member is a
// float4 so HLSL's packing rules cannot silently shift a field; the two
// declarations must be kept in the same order.
type constants struct {
	Frame        [4]float32
	Bounds       [4]float32
	FillColor    [4]float32
	StrokeColor  [4]float32
	ShadowColor  [4]float32
	Radius       [4]float32
	RectHalf     [4]float32
	Misc         [4]float32
	ShadowOffset [4]float32
	TexParams    [4]float32
	Inset        [4]float32
	// NdcRect carries the quad corners in clip space (x1, y1, x2, y2); the
	// vertex shader expands them via SV_VertexID, so draws need no vertex
	// buffer. The line shader reuses it as the two endpoints.
	NdcRect [4]float32
}

// GPU owns the Direct3D 11 device and swap chain backing one window.
type GPU struct {
	dev  *device
	ctx  *deviceContext
	swap *swapChain
	rtv  *renderTargetView
	back *texture2D

	width, height uint32
	featureLevel  uint32

	// presentDur and presentN accumulate Present wall time between
	// FYNE_DX_DEBUG reports; blocking in Present means the compositor or GPU
	// is gating the frame, which a CPU profile cannot separate from spin.
	presentDur time.Duration
	presentN   int

	// frameStart and frameDur time Clear-to-Present: the CPU cost of walking
	// the tree and issuing every draw. A frame time far above avg present means
	// the app is CPU-bound in the painter, not waiting on the GPU.
	frameStart time.Time
	frameDur   time.Duration
}

func NewGPU(hwnd windows.Handle, width, height uint32) (*GPU, error) {
	if width == 0 {
		width = 1
	}
	if height == 0 {
		height = 1
	}
	dev, ctx, swap, level, err := createDeviceAndSwapChain(hwnd, width, height)
	if err != nil {
		return nil, err
	}
	// Frame latency bounds how many presents may be in flight before Present
	// blocks - and Intel's driver spins that block on the CPU. On battery the
	// GPU is down-clocked to near the frame budget, so tight budgets (1 or 2,
	// both profiled) turn into milliseconds of spinning per frame; 3 gives the
	// queue depth the GL driver effectively enjoys. Latency only binds when
	// presents outpace retirement, so this costs no drag lag while the GPU
	// keeps up.
	if err := setMaximumFrameLatency(dev, 3); err != nil {
		fyne.LogError("directx: could not set the frame latency", err)
	}

	g := &GPU{dev: dev, ctx: ctx, swap: swap, width: width, height: height, featureLevel: level}
	if err := g.createTarget(); err != nil {
		g.Release()
		return nil, err
	}
	return g, nil
}

// MaxTextureSize reports the largest texture dimension the device supports,
// which Direct3D 11 fixes per feature level.
func (g *GPU) MaxTextureSize() int {
	switch {
	case g.featureLevel >= featureLevel11_0:
		return 16384
	case g.featureLevel >= featureLevel10_0:
		return 8192
	default:
		return 4096
	}
}

// createTarget (re)binds a render target view to the current back buffer.
func (g *GPU) createTarget() error {
	back, err := g.swap.GetBuffer(0)
	if err != nil {
		return err
	}
	rtv, err := g.dev.CreateRenderTargetView(back)
	if err != nil {
		back.Release()
		return err
	}
	g.back, g.rtv = back, rtv
	g.ctx.OMSetRenderTargets(rtv)
	return nil
}

// resize drops the old back buffer view, resizes the swap chain and rebuilds the
// view. The view must be released first or ResizeBuffers fails with
// DXGI_ERROR_INVALID_CALL because the buffer is still referenced.
func (g *GPU) Resize(width, height uint32) error {
	if width == 0 || height == 0 || (width == g.width && height == g.height) {
		return nil
	}
	g.ctx.OMSetRenderTargets(nil)
	releaseCOM(&g.rtv)
	releaseCOM(&g.back)

	if err := g.swap.ResizeBuffers(width, height); err != nil {
		// Leave the painter with a usable target even on failure, rather than no
		// render target at all.
		g.createTarget()
		return err
	}
	g.width, g.height = width, height
	return g.createTarget()
}

// Present publishes the back buffer with sync interval 0, matching the GL
// driver's SwapInterval(0): frame pacing comes from the run loop's ticker, and
// a vsynced Present would instead block (often spin-waiting) inside the driver
// for the rest of the refresh interval, burning a core during any continuous
// redraw. Windowed flip-model presents are still composed tear-free by DWM.
//
// It reports whether the device was lost (removed or reset - a GPU driver
// upgrade or timeout kills every resource); the caller must then rebuild the
// GPU and its painter.
func (g *GPU) Present() (deviceLost bool) {
	var start time.Time
	if dxDebug {
		start = time.Now()
		if !g.frameStart.IsZero() {
			g.frameDur += start.Sub(g.frameStart)
		}
	}
	hr := g.swap.Present(0)
	if dxDebug {
		g.presentDur += time.Since(start)
		g.presentN++
	}
	if hr.deviceLost() {
		fyne.LogError("directx: device lost", g.dev.RemovedReason().error("GetDeviceRemovedReason"))
		return true
	}
	return false
}

// Size reports the back buffer dimensions in pixels.
func (g *GPU) Size() (width, height uint32) {
	return g.width, g.height
}

func (g *GPU) Release() {
	releaseCOM(&g.rtv)
	releaseCOM(&g.back)
	releaseCOM(&g.swap)
	releaseCOM(&g.ctx)
	releaseCOM(&g.dev)
}

// The Painter needs two blend modes, mirroring the two glBlendFunc calls in the
// GL Painter. Getting this wrong is silent - nothing errors, output just looks
// subtly wrong - so both are built by named functions the tests can check.

// straightAlphaBlend is for shapes, whose fragment shaders emit un-premultiplied
// colour (see fragmentColor).
func straightAlphaBlend() blendDesc {
	d := blendDesc{}
	d.RenderTarget[0] = renderTargetBlendDesc{
		BlendEnable:           1,
		SrcBlend:              blendSrcAlpha,
		DestBlend:             blendInvSrcAlpha,
		BlendOp:               blendOpAdd,
		SrcBlendAlpha:         blendSrcAlpha,
		DestBlendAlpha:        blendInvSrcAlpha,
		BlendOpAlpha:          blendOpAdd,
		RenderTargetWriteMask: colorWriteEnableAll,
	}
	return d
}

// premultipliedAlphaBlend is for textures. Every texture here comes from a Go
// image.RGBA, and Go's image/color is alpha-premultiplied by definition, so glyph
// coverage and image alpha are already folded into the colour channels. Blending
// those with SRC_ALPHA would apply alpha twice and eat the anti-aliased edges,
// which shows up as thin, ragged text.
func premultipliedAlphaBlend() blendDesc {
	d := straightAlphaBlend()
	d.RenderTarget[0].SrcBlend = blendOne
	d.RenderTarget[0].SrcBlendAlpha = blendOne
	return d
}

// gpuTexture is one uploaded image plus the view the shader samples. width and
// height are set only for textures created by imgToTexture, marking them safe
// for the same-size reuse pool.
type gpuTexture struct {
	tex           *texture2D
	srv           *shaderResourceView
	width, height uint32
}

func (t *gpuTexture) release() {
	if t == nil { // negative cache entries store nil
		return
	}
	releaseCOM(&t.srv)
	releaseCOM(&t.tex)
}

// Painter renders a Fyne canvas with Direct3D 11.
type Painter struct {
	canvas fyne.Canvas
	g      *GPU

	vsQuad  *vertexShader
	vsLine  *vertexShader
	vsGlyph *vertexShader

	psTextured *pixelShader
	psGlyph    *pixelShader
	psRect     *pixelShader
	psRound    *pixelShader
	psEllipse  *pixelShader
	psArc      *pixelShader
	psPolygon  *pixelShader
	psArbPoly  *pixelShader
	psBezier   *pixelShader
	psBlur     *pixelShader
	psLine     *pixelShader

	blend        *blendState // straight alpha, for shapes
	blendPremul  *blendState // premultiplied alpha, for textures
	blendReplace *blendState // overwrite the destination, for the blur passes
	rsPlain      *rasterizerState
	rsScissor    *rasterizerState
	sampLinear   *samplerState
	sampNearest  *samplerState

	cbuf *buffer
	// auxbuf is the register(b1) buffer carrying the variable-length data the
	// polygon and curve shaders need, which does not fit the shared cbuffer.
	auxbuf *buffer
	// glyphbuf is the register(b2) buffer holding one batch of glyph quads. It
	// is dynamic rather than default because a batch writes fewer entries than
	// the buffer holds and UpdateSubresource cannot partially update a constant
	// buffer without D3D11.1; Map(WRITE_DISCARD) writes only what is used.
	glyphbuf *buffer
	// glyphPending accumulates glyph quads until something forces a flush.
	glyphPending []glyphInst
	// glyphScratch and entryScratch are reused by drawText to collect one
	// string's shaped glyphs and their atlas slots.
	glyphScratch []paint.PlacedGlyph
	entryScratch []glyphEntry
	// shaped caches the result of shaping one string, keyed by everything the
	// shaper reads. Shaping (segmenting plus HarfBuzz) is the dominant cost of
	// text once rasterisation and draws are cached and batched, and a UI's
	// strings mostly repeat frame over frame.
	shaped map[shapedKey]shapedEntry

	atlas glyphAtlas

	blurSnap   blurSnapshot
	blurKernel blurKernelTexture

	// textures maps the uint32 handle stored in Fyne's texture cache to the COM
	// objects behind it. The shared cache is typed to uint32 for the GL driver, so
	// D3D resources are indirected through an id rather than stored directly.
	textures  map[uint32]*gpuTexture
	nextTexID uint32

	// texPool parks released textures for same-size reuse. Playback-style
	// workloads (gauges, plots, changing labels) free and recreate equal-sized
	// textures every frame, and a D3D11 create/destroy pair is a kernel round
	// trip - far dearer than refilling existing storage with UpdateSubresource.
	texPool     map[[2]uint32][]pooledTexture
	texPoolSize int
	frameTick   uint32

	// drawCount tallies draw calls between FYNE_DX_DEBUG reports. glyphCount and
	// batchCount tally the batched quads and the DrawInstanced calls that
	// carried them; their ratio is the only direct evidence that batching is
	// collapsing anything, since drawCount already counts a batch as one call.
	// rectCount is the subset of those quads that were solid rectangles riding
	// the batch rather than glyphs.
	drawCount  int
	glyphCount int
	batchCount int
	rectCount  int

	// clippedTextTextures holds the windowed textures for text runs wider than
	// the device texture limit, keyed by the text object as the GL painter does.
	clippedTextTextures map[*canvas.Text]clippedTextEntry

	// userShaders caches the compiled program, uniform buffer and textures of
	// each canvas.Shader, keyed by Shader.Name as the GL Painter does.
	userShaders map[string]*userShader

	pixScale, texScale float32
	clipping           bool

	// last* mirror the pipeline state bound on the immediate context, which
	// persists across draws and frames. The hot path only issues a COM call -
	// a syscall each - for state that actually changed.
	lastBlend    *blendState
	lastTopology uint32
	lastVS       *vertexShader
	lastPS       *pixelShader
	lastSRV      *shaderResourceView
	lastSampler  *samplerState
}

// Repainter is a canvas that can redraw itself into the back buffer on demand.
//
// Capture needs this because the swap chain discards the back buffer contents on
// Present, so by the time a capture is requested there is nothing valid to read.
type Repainter interface {
	Repaint(fyne.Size)
}

// Declare conformity with the Painter interface.
var _ paint.Painter = (*Painter)(nil)

// pooledTexture is one parked texture plus the frame it was parked on, so the
// sweep in Clear can drop sizes that stopped recurring.
type pooledTexture struct {
	tex  *gpuTexture
	tick uint32
}

const (
	// texPoolMax bounds how many textures sit parked at once; beyond it a
	// released texture is destroyed rather than pooled.
	texPoolMax = 32
	// texPoolTTL is how many frames a parked texture survives unused (~2s at
	// 60fps): long enough to ride out a pause, short enough to return VRAM.
	texPoolTTL = 120
)

func NewPainter(c fyne.Canvas, g *GPU) *Painter {
	p := &Painter{canvas: c, g: g,
		textures: map[uint32]*gpuTexture{},
		texPool:  map[[2]uint32][]pooledTexture{},
		shaped:   map[shapedKey]shapedEntry{},
	}
	p.SetFrameBufferScale(1.0)
	return p
}

// Init compiles the shader set and builds the fixed pipeline state. It is called
// once, after the device exists. A failure here is a bug (the shaders ship with
// the binary), so it panics the way the GL painter does on a shader compile
// error, rather than leaving a painter that silently draws nothing.
func (p *Painter) Init() {
	if p.vsQuad != nil {
		return
	}
	if err := p.initPipeline(); err != nil {
		panic("directx: painter init failed: " + err.Error())
	}
	// Bind everything that never changes once; upload only touches the rest.
	// There is no input layout or vertex buffer at all: the vertex shaders
	// generate geometry from SV_VertexID and the constant buffer.
	p.g.ctx.VSSetConstantBuffer(p.cbuf)
	p.g.ctx.PSSetConstantBuffer(p.cbuf)
	// Only the vertex stage reads the glyph batch; the pixel stage gets each
	// glyph's colour down the interpolator instead of indexing the array again.
	p.g.ctx.VSSetConstantBufferAt(2, p.glyphbuf)
	p.g.ctx.OMSetBlendState(p.blend)
	p.g.ctx.RSSetState(p.rsPlain)
	p.lastBlend = p.blend
}

func (p *Painter) initPipeline() error {
	vsCode, err := compileShader(VertexPassthrough2D, "main", "vs_4_0")
	if err != nil {
		return fmt.Errorf("vertex shader: %w", err)
	}
	if p.vsQuad, err = p.g.dev.CreateVertexShader(vsCode); err != nil {
		return fmt.Errorf("vertex shader: %w", err)
	}
	lineCode, err := compileShader(VertexLine, "main", "vs_4_0")
	if err != nil {
		return fmt.Errorf("line vertex shader: %w", err)
	}
	if p.vsLine, err = p.g.dev.CreateVertexShader(lineCode); err != nil {
		return fmt.Errorf("line vertex shader: %w", err)
	}
	glyphCode, err := compileShader(VertexGlyph, "main", "vs_4_0")
	if err != nil {
		return fmt.Errorf("glyph vertex shader: %w", err)
	}
	if p.vsGlyph, err = p.g.dev.CreateVertexShader(glyphCode); err != nil {
		return fmt.Errorf("glyph vertex shader: %w", err)
	}

	for _, s := range []struct {
		src string
		out **pixelShader
	}{
		{PixelSimple, &p.psTextured},
		{PixelRectangle, &p.psRect},
		{PixelRoundRectangle, &p.psRound},
		{PixelEllipse, &p.psEllipse},
		{PixelArc, &p.psArc},
		{PixelRegularPolygon, &p.psPolygon},
		{PixelArbitraryPolygon, &p.psArbPoly},
		{PixelBezierCurve, &p.psBezier},
		{PixelBlur, &p.psBlur},
		{PixelLine, &p.psLine},
		{PixelGlyph, &p.psGlyph},
	} {
		code, err := compileShader(s.src, "main", "ps_4_0")
		if err != nil {
			return fmt.Errorf("pixel shader: %w", err)
		}
		if *s.out, err = p.g.dev.CreatePixelShader(code); err != nil {
			return fmt.Errorf("pixel shader: %w", err)
		}
	}

	straight := straightAlphaBlend()
	if p.blend, err = p.g.dev.CreateBlendState(&straight); err != nil {
		return fmt.Errorf("blend state: %w", err)
	}
	premul := premultipliedAlphaBlend()
	if p.blendPremul, err = p.g.dev.CreateBlendState(&premul); err != nil {
		return fmt.Errorf("premultiplied blend state: %w", err)
	}
	replace := replaceBlend()
	if p.blendReplace, err = p.g.dev.CreateBlendState(&replace); err != nil {
		return fmt.Errorf("replace blend state: %w", err)
	}

	// Two rasterizer states rather than a dynamic one: scissoring is a pipeline
	// state in D3D11, so clipping toggles between these.
	plain := rasterizerDesc{FillMode: fillSolid, CullMode: cullNone, DepthClipEnable: 1}
	if p.rsPlain, err = p.g.dev.CreateRasterizerState(&plain); err != nil {
		return fmt.Errorf("rasterizer state: %w", err)
	}
	scissor := plain
	scissor.ScissorEnable = 1
	if p.rsScissor, err = p.g.dev.CreateRasterizerState(&scissor); err != nil {
		return fmt.Errorf("scissor rasterizer state: %w", err)
	}

	linear := samplerDesc{
		Filter: filterMinMagMipLinear, AddressU: addressClamp, AddressV: addressClamp,
		AddressW: addressClamp, MaxLOD: math.MaxFloat32,
	}
	if p.sampLinear, err = p.g.dev.CreateSamplerState(&linear); err != nil {
		return fmt.Errorf("sampler: %w", err)
	}
	nearest := linear
	nearest.Filter = filterMinMagMipPoint
	if p.sampNearest, err = p.g.dev.CreateSamplerState(&nearest); err != nil {
		return fmt.Errorf("sampler: %w", err)
	}

	// Default usage, written with UpdateSubresource: one syscall per write
	// instead of a Map/Unmap pair, and the driver handles versioning. Only the
	// polygon and curve shaders touch this one, so the extra syscall the shared
	// cbuf avoids below is not worth paying for here.
	if p.auxbuf, err = p.g.dev.CreateBuffer(&bufferDesc{
		ByteWidth: uint32(unsafe.Sizeof(auxConstants{})),
		Usage:     usageDefault,
		BindFlags: bindConstantBuffer,
	}, nil); err != nil {
		return fmt.Errorf("aux constant buffer: %w", err)
	}
	if p.cbuf, err = p.g.dev.CreateBuffer(&bufferDesc{
		ByteWidth:      uint32(unsafe.Sizeof(constants{})),
		Usage:          usageDynamic,
		BindFlags:      bindConstantBuffer,
		CPUAccessFlags: cpuAccessWrite,
	}, nil); err != nil {
		return fmt.Errorf("constant buffer: %w", err)
	}
	if p.glyphbuf, err = p.g.dev.CreateBuffer(&bufferDesc{
		ByteWidth:      uint32(unsafe.Sizeof(glyphInst{})) * glyphBatchMax,
		Usage:          usageDynamic,
		BindFlags:      bindConstantBuffer,
		CPUAccessFlags: cpuAccessWrite,
	}, nil); err != nil {
		return fmt.Errorf("glyph instance buffer: %w", err)
	}
	if err = p.initAtlas(); err != nil {
		return fmt.Errorf("glyph atlas: %w", err)
	}

	return nil
}

func (p *Painter) ready() bool { return p.vsQuad != nil }

// textureStats totals the live textures and the video memory they occupy, so
// FYNE_DX_DEBUG can show whether a workload is accumulating textures. This is
// the one thing a Go heap profile cannot see: a texture costs a pointer on the
// Go side and megabytes on the GPU.
//
// ponytail: walked on demand rather than kept as a running counter. It only
// runs once per debug report, and a counter maintained across every create,
// pool, evict and device-loss path is exactly the kind of bookkeeping that
// drifts and then lies to you.
func (p *Painter) textureStats() (count int, bytes uint64) {
	add := func(t *gpuTexture) {
		if t == nil { // negative cache entries store nil
			return
		}
		count++
		bytes += uint64(t.width) * uint64(t.height) * 4 // every format here is 32bpp
	}
	for _, t := range p.textures {
		add(t)
	}
	for _, entries := range p.texPool {
		for _, e := range entries {
			add(e.tex)
		}
	}
	for _, s := range p.userShaders {
		for _, t := range s.textures {
			add(t.GPU)
		}
	}
	return count, bytes
}

const debugReportFrames = 120

func (p *Painter) Clear() {
	p.frameTick++
	p.sweepTexPool()
	// Anything still queued belongs to a frame that is about to be wiped, so
	// drop it rather than drawing it into the cleared target.
	p.glyphPending = p.glyphPending[:0]

	if dxDebug {
		if p.frameTick%debugReportFrames == 0 {
			avgPresent := time.Duration(0)
			avgFrame := time.Duration(0)
			if p.g.presentN > 0 {
				avgPresent = p.g.presentDur / time.Duration(p.g.presentN)
				avgFrame = p.g.frameDur / time.Duration(p.g.presentN)
			}
			perBatch := float64(0)
			if p.batchCount > 0 {
				perBatch = float64(p.glyphCount) / float64(p.batchCount)
			}
			texCount, texBytes := p.textureStats()
			log.Printf("directx: %d draw calls/frame (%d glyphs + %d rects in %d batches, %.1f per batch), "+
				"avg frame cpu %v, avg present %v, %d textures live using %dMB (over %d frames)",
				p.drawCount/debugReportFrames,
				(p.glyphCount-p.rectCount)/debugReportFrames,
				p.rectCount/debugReportFrames,
				p.batchCount/debugReportFrames, perBatch,
				avgFrame.Round(10*time.Microsecond),
				avgPresent.Round(10*time.Microsecond),
				texCount, texBytes/(1024*1024), debugReportFrames)
			p.drawCount, p.glyphCount, p.batchCount, p.rectCount = 0, 0, 0, 0
			p.g.presentDur, p.g.presentN, p.g.frameDur = 0, 0, 0
		}
		p.g.frameStart = time.Now()
	}

	r, g, b, a := theme.Color(theme.ColorNameBackground).RGBA()
	rgba := [4]float32{
		float32(r) / max16bit, float32(g) / max16bit,
		float32(b) / max16bit, float32(a) / max16bit,
	}
	p.g.ctx.OMSetRenderTargets(p.g.rtv)
	p.g.ctx.ClearRenderTargetView(p.g.rtv, &rgba)
}

func (p *Painter) SetFrameBufferScale(scale float32) {
	p.texScale = scale
	p.pixScale = p.canvas.Scale() * p.texScale
}

// SetOutputSize sets the viewport. Callers should pass the back buffer size, not
// a size derived from the canvas: the two can disagree by a pixel at fractional
// scales, and momentarily disagree by much more mid-resize.
func (p *Painter) SetOutputSize(width, height int) {
	v := viewport{Width: float32(width), Height: float32(height), MaxDepth: 1}
	p.g.ctx.RSSetViewports(&v)
}

// StartClipping restricts drawing to the given canvas rectangle. Fyne passes
// top-origin canvas coordinates and D3D scissor rects are also top-origin, so
// unlike the GL Painter no vertical flip is needed.
func (p *Painter) StartClipping(pos fyne.Position, size fyne.Size) {
	p.flushGlyphs() // queued glyphs belong to the outgoing scissor rectangle
	x := p.textureScale(pos.X)
	y := p.textureScale(pos.Y)
	w := p.textureScale(size.Width)
	h := p.textureScale(size.Height)
	if w < 0 {
		w = 0
	}
	if h < 0 {
		h = 0
	}
	r := d3dRect{Left: int32(x), Top: int32(y), Right: int32(x + w), Bottom: int32(y + h)}
	p.g.ctx.RSSetScissorRects(&r)
	if !p.clipping {
		p.g.ctx.RSSetState(p.rsScissor)
		p.clipping = true
	}
}

func (p *Painter) StopClipping() {
	if !p.clipping {
		return
	}
	p.flushGlyphs() // as in StartClipping: emit under the rectangle still in force
	p.g.ctx.RSSetState(p.rsPlain)
	p.clipping = false
}

func (p *Painter) Paint(obj fyne.CanvasObject, pos fyne.Position, frame fyne.Size, clip *internal.ClipItem) {
	if !p.ready() || !obj.Visible() {
		return
	}

	size := obj.Size()
	var clipPos fyne.Position
	var clipSize fyne.Size
	if clip != nil {
		clipPos, clipSize = clip.Rect()
	} else {
		clipSize = frame
	}
	if pos.Y > clipPos.Y+clipSize.Height || pos.Y+size.Height < clipPos.Y ||
		pos.X > clipPos.X+clipSize.Width || pos.X+size.Width < clipPos.X {
		return
	}

	p.drawObject(obj, pos, frame, clip)
}

func (p *Painter) drawObject(o fyne.CanvasObject, pos fyne.Position, frame fyne.Size, clip *internal.ClipItem) {
	switch obj := o.(type) {
	case *canvas.Circle:
		p.drawCircle(obj, pos, frame)
	case *canvas.Ellipse:
		p.drawEllipse(obj, pos, frame)
	case *canvas.Arc:
		p.drawArc(obj, pos, frame)
	case *canvas.RegularPolygon:
		p.drawPolygon(obj, pos, frame)
	case *canvas.ArbitraryPolygon:
		p.drawArbitraryPolygon(obj, pos, frame)
	case *canvas.BezierCurve:
		p.drawBezierCurve(obj, pos, frame)
	case *canvas.Blur:
		p.drawBlur(obj, pos, frame)
	case *canvas.Line:
		p.drawLine(obj, pos, frame)
	case *canvas.Image:
		aspect := obj.Aspect()
		if aspect == 0 {
			aspect = 1 // fallback, should not occur - normally an image load error
		}
		p.drawTexture(obj, p.imageTexture, pos, obj.Size(), frame, obj.FillMode, float32(obj.Alpha()), aspect)
	case *canvas.Raster:
		p.drawTexture(obj, p.rasterTexture, pos, obj.Size(), frame, canvas.ImageFillStretch, float32(obj.Alpha()), 0)
	case *canvas.Rectangle:
		p.drawRectangle(obj, pos, frame)
	case *canvas.Text:
		p.drawText(obj, pos, frame, clip)
	case *canvas.LinearGradient:
		p.drawTexture(obj, p.linearGradientTexture, pos, obj.Size(), frame, canvas.ImageFillStretch, 1, 0)
	case *canvas.RadialGradient:
		p.drawTexture(obj, p.radialGradientTexture, pos, obj.Size(), frame, canvas.ImageFillStretch, 1, 0)
	case *canvas.Shader:
		p.drawShader(obj, pos, frame)
	}
}

// ---------------------------------------------------------------------------
// Pipeline helpers
// ---------------------------------------------------------------------------

// upload writes the constants and issues the draw; the vertex shader generates
// the geometry from c.NdcRect, so no vertex data is uploaded at all. blend
// selects straight or premultiplied alpha, which differs between shapes and
// textures the same way it does in the GL Painter.
func (p *Painter) upload(c *constants, ndcRect [4]float32, vs *vertexShader, ps *pixelShader,
	topology uint32, blend *blendState, vertexCount uint32,
) {
	// Glyphs queued before this object have to reach the back buffer before it
	// does, or alpha blending sees the wrong order.
	p.flushGlyphs()

	c.NdcRect = ndcRect
	// Map(WRITE_DISCARD) rather than UpdateSubresource: this buffer is rewritten
	// once per draw, and discard is the rename path the driver is tuned for -
	// it hands back fresh storage instead of versioning a default-usage copy.
	m, err := p.g.ctx.Map(unsafe.Pointer(p.cbuf), mapWriteDiscard)
	if err != nil {
		// Realistically only device removal, which Present catches next; drawing
		// with whatever the buffer held last would draw the wrong thing.
		fyne.LogError("directx: mapping the constant buffer", err)
		return
	}
	*(*constants)(m.Data) = *c
	p.g.ctx.Unmap(unsafe.Pointer(p.cbuf))

	p.setState(vs, ps, topology, blend)
	p.g.ctx.Draw(vertexCount, 0)
	p.drawCount++
}

// setState binds the pipeline objects a draw needs. The constant buffer
// bindings never change after Init; everything else is set only when it differs
// from what is already bound.
func (p *Painter) setState(vs *vertexShader, ps *pixelShader, topology uint32, blend *blendState) {
	if blend != p.lastBlend {
		p.g.ctx.OMSetBlendState(blend)
		p.lastBlend = blend
	}
	if topology != p.lastTopology {
		p.g.ctx.IASetPrimitiveTopology(topology)
		p.lastTopology = topology
	}
	if vs != p.lastVS {
		p.g.ctx.VSSetShader(vs)
		p.lastVS = vs
	}
	if ps != p.lastPS {
		p.g.ctx.PSSetShader(ps)
		p.lastPS = ps
	}
}

// pushGlyph queues one glyph quad for the current batch, emitting the batch
// first if it is full.
func (p *Painter) pushGlyph(inst glyphInst) {
	if len(p.glyphPending) == glyphBatchMax {
		p.flushGlyphs()
	}
	p.glyphPending = append(p.glyphPending, inst)
}

// flushGlyphs draws every queued glyph quad as one instanced draw. Callers must
// invoke it before any other draw, before the scissor rectangle moves, before
// the atlas contents are overwritten, and at the end of the frame; every one of
// those would otherwise reorder, mis-clip or mis-sample the queue.
func (p *Painter) flushGlyphs() {
	n := len(p.glyphPending)
	if n == 0 {
		return
	}

	m, err := p.g.ctx.Map(unsafe.Pointer(p.glyphbuf), mapWriteDiscard)
	if err != nil {
		fyne.LogError("directx: mapping the glyph batch", err)
		p.glyphPending = p.glyphPending[:0]
		return
	}
	copy(unsafe.Slice((*glyphInst)(m.Data), n), p.glyphPending)
	p.g.ctx.Unmap(unsafe.Pointer(p.glyphbuf))
	p.glyphPending = p.glyphPending[:0]

	if p.atlas.tex.srv != p.lastSRV {
		p.g.ctx.PSSetShaderResource(p.atlas.tex.srv)
		p.lastSRV = p.atlas.tex.srv
	}
	if p.sampLinear != p.lastSampler {
		p.g.ctx.PSSetSampler(p.sampLinear)
		p.lastSampler = p.sampLinear
	}

	p.setState(p.vsGlyph, p.psGlyph, topologyTriangleList, p.blend)
	p.g.ctx.DrawInstanced(6, uint32(n))
	p.drawCount++
	p.glyphCount += n
	p.batchCount++
}

// Flush emits anything still queued. The canvas calls it once the render tree
// has been walked, because the last object drawn is often text.
func (p *Painter) Flush() {
	if !p.ready() {
		return
	}
	p.flushGlyphs()
}

// baseConstants fills the fields every shape shader reads.
func (p *Painter) baseConstants(frame fyne.Size, bounds [4]float32, fill, stroke color.Color,
	shadow canvas.Shadow,
) constants {
	fw, fh := p.scaleFrameSize(frame)
	x1, x2, y1, y2 := p.scaleRectCoords(bounds[0], bounds[2], bounds[1], bounds[3])

	c := constants{}
	c.Frame = [4]float32{fw, fh, 0, 0}
	c.Bounds = [4]float32{x1, y1, x2, y2}

	r, g, b, a := fragmentColor(fill)
	c.FillColor = [4]float32{r, g, b, a}

	strokeCol := stroke
	if strokeCol == nil {
		strokeCol = color.Transparent
	}
	r, g, b, a = fragmentColor(strokeCol)
	c.StrokeColor = [4]float32{r, g, b, a}

	c.RectHalf[3] = roundToPixel(edgeSoftness*p.pixScale, 1.0)

	if paint.IsShadowVisible(shadow) {
		r, g, b, a = fragmentColor(shadow.Color)
		c.ShadowColor = [4]float32{r, g, b, a}
		c.ShadowOffset[0] = roundToPixel(shadow.Offset.X*p.pixScale, 1.0)
		c.ShadowOffset[1] = roundToPixel(shadow.Offset.Y*p.pixScale, 1.0)
		c.ShadowOffset[2] = float32(shadow.Variant)
		c.Misc[1] = 1 // addShadow
		c.Misc[2] = roundToPixel(shadow.BlurRadius*p.pixScale, 1.0)
		c.Misc[3] = roundToPixel(shadow.Spread*p.pixScale, 1.0)
	}
	return c
}

// ---------------------------------------------------------------------------
// Shape drawing
// ---------------------------------------------------------------------------

func (p *Painter) drawRectangle(d3dRect *canvas.Rectangle, pos fyne.Position, frame fyne.Size) {
	topRight := paint.GetCornerRadius(d3dRect.TopRightCornerRadius, d3dRect.CornerRadius)
	topLeft := paint.GetCornerRadius(d3dRect.TopLeftCornerRadius, d3dRect.CornerRadius)
	bottomRight := paint.GetCornerRadius(d3dRect.BottomRightCornerRadius, d3dRect.CornerRadius)
	bottomLeft := paint.GetCornerRadius(d3dRect.BottomLeftCornerRadius, d3dRect.CornerRadius)
	if topRight == 0 && topLeft == 0 && bottomRight == 0 && bottomLeft == 0 &&
		p.pushSolidRect(d3dRect, pos, frame) {
		return
	}
	p.drawOblong(d3dRect, d3dRect.FillColor, d3dRect.StrokeColor, d3dRect.StrokeWidth,
		topRight, topLeft, bottomRight, bottomLeft, d3dRect.Aspect, d3dRect.Shadow, pos, frame)
}

// pushSolidRect queues a plain filled rectangle on the glyph batch instead of
// issuing a draw of its own: a run of table cells and their labels then costs
// one DrawInstanced rather than a constant-buffer Map and Draw per cell. Only
// the square-cornered, unstroked, unshadowed case qualifies - for it, the
// generic path's pixel shader draws a hard-edged fill (no antialiasing), so a
// quad snapped to the same device pixels is identical output. Reports false
// when the rectangle needs the generic path.
func (p *Painter) pushSolidRect(r *canvas.Rectangle, pos fyne.Position, frame fyne.Size) bool {
	// A nonzero stroke width always goes to the generic path: its shader draws
	// the rim in the stroke colour even when that colour is transparent, which
	// visually insets the fill - not something a plain quad can reproduce.
	if r.Aspect != 0 || r.StrokeWidth != 0 || paint.IsShadowVisible(r.Shadow) {
		return false
	}
	if r.FillColor == nil || r.FillColor == color.Transparent {
		return true // nothing to draw, and no stroke or shadow either
	}
	if !p.ensureAtlasWhite() {
		return false
	}

	// The same snapping vecRectCoords applies, without the antialiasing skirt
	// it adds for the generic quad (whose shader clips it back to bounds).
	size := r.Size()
	x := roundToPixel(pos.X, p.pixScale)
	y := roundToPixel(pos.Y, p.pixScale)
	w := roundToPixel(size.Width, p.pixScale)
	h := roundToPixel(size.Height, p.pixScale)
	x1, x2, y1, y2 := p.scaleRectCoords(x, x+w, y, y+h)
	fw, fh := p.scaleFrameSize(frame)

	cr, cg, cb, ca := fragmentColor(r.FillColor)
	p.pushGlyph(glyphInst{
		NDC: [4]float32{
			x1/fw*2 - 1, 1 - y1/fh*2,
			x2/fw*2 - 1, 1 - y2/fh*2,
		},
		UV:    [4]float32{p.atlas.whiteU, p.atlas.whiteV, p.atlas.whiteU, p.atlas.whiteV},
		Color: [4]float32{cr, cg, cb, ca},
	})
	p.rectCount++
	return true
}

func (p *Painter) drawOblong(obj fyne.CanvasObject, fill, stroke color.Color, strokeWidth,
	topRightRadius, topLeftRadius, bottomRightRadius, bottomLeftRadius, aspect float32,
	shadow canvas.Shadow, pos fyne.Position, frame fyne.Size,
) {
	if !paint.IsShadowVisible(shadow) && (fill == color.Transparent || fill == nil) &&
		(stroke == color.Transparent || stroke == nil || strokeWidth == 0) {
		return
	}

	rounded := topRightRadius != 0 || topLeftRadius != 0 || bottomRightRadius != 0 || bottomLeftRadius != 0
	points, bounds := p.vecRectCoords(pos, obj, frame, aspect, shadow)
	c := p.baseConstants(frame, bounds, fill, stroke, shadow)

	strokeScaled := roundToPixel(strokeWidth*p.pixScale, 1.0)
	ps := p.psRect
	if rounded {
		ps = p.psRound
		c.RectHalf[2] = strokeScaled * 0.5

		width := c.Bounds[2] - c.Bounds[0] - strokeScaled
		height := c.Bounds[3] - c.Bounds[1] - strokeScaled
		c.RectHalf[0] = width * 0.5
		c.RectHalf[1] = height * 0.5

		size := fyne.NewSize(bounds[2]-bounds[0], bounds[3]-bounds[1])
		c.Radius = [4]float32{
			roundToPixel(paint.GetMaximumCornerRadius(topRightRadius, topLeftRadius, bottomRightRadius, size)*p.pixScale, 1.0),
			roundToPixel(paint.GetMaximumCornerRadius(bottomRightRadius, bottomLeftRadius, topRightRadius, size)*p.pixScale, 1.0),
			roundToPixel(paint.GetMaximumCornerRadius(topLeftRadius, topRightRadius, bottomLeftRadius, size)*p.pixScale, 1.0),
			roundToPixel(paint.GetMaximumCornerRadius(bottomLeftRadius, bottomRightRadius, topLeftRadius, size)*p.pixScale, 1.0),
		}
	} else {
		c.Misc[0] = strokeScaled
	}

	p.upload(&c, points, p.vsQuad, ps, topologyTriangleStrip, p.blend, 4)
}

func (p *Painter) drawCircle(circle *canvas.Circle, pos fyne.Position, frame fyne.Size) {
	if (circle.FillColor == color.Transparent || circle.FillColor == nil) &&
		(circle.StrokeColor == color.Transparent || circle.StrokeColor == nil || circle.StrokeWidth == 0) &&
		!paint.IsShadowVisible(circle.Shadow) {
		return
	}
	// Aspect 1 squares the drawn bounds inside the object's rectangle, matching
	// the GL painter's vecSquareCoords: a circle stays circular even when its
	// bounding box is not square.
	points, bounds := p.vecRectCoords(pos, circle, frame, 1, circle.Shadow)
	c := p.baseConstants(frame, bounds, circle.FillColor, circle.StrokeColor, circle.Shadow)

	strokeScaled := roundToPixel(circle.StrokeWidth*p.pixScale, 1.0)
	c.Misc[0] = strokeScaled
	// A circle is an ellipse with equal semi-axes taken from the drawn bounds.
	halfW := (c.Bounds[2] - c.Bounds[0]) * 0.5
	halfH := (c.Bounds[3] - c.Bounds[1]) * 0.5
	softness := c.RectHalf[3]
	c.Radius = [4]float32{halfW - softness, halfH - softness, 0, 0}

	p.upload(&c, points, p.vsQuad, p.psEllipse, topologyTriangleStrip, p.blend, 4)
}

func (p *Painter) drawEllipse(ellipse *canvas.Ellipse, pos fyne.Position, frame fyne.Size) {
	if (ellipse.FillColor == color.Transparent || ellipse.FillColor == nil) &&
		(ellipse.StrokeColor == color.Transparent || ellipse.StrokeColor == nil || ellipse.StrokeWidth == 0) &&
		!paint.IsShadowVisible(ellipse.Shadow) {
		return
	}
	points, bounds := p.vecRectCoords(pos, ellipse, frame, 0, ellipse.Shadow)
	c := p.baseConstants(frame, bounds, ellipse.FillColor, ellipse.StrokeColor, ellipse.Shadow)

	c.Misc[0] = roundToPixel(ellipse.StrokeWidth*p.pixScale, 1.0)
	softness := c.RectHalf[3]
	c.Radius = [4]float32{
		(c.Bounds[2]-c.Bounds[0])*0.5 - softness,
		(c.Bounds[3]-c.Bounds[1])*0.5 - softness,
		0, 0,
	}

	p.upload(&c, points, p.vsQuad, p.psEllipse, topologyTriangleStrip, p.blend, 4)
}

func (p *Painter) drawArc(arc *canvas.Arc, pos fyne.Position, frame fyne.Size) {
	if ((arc.FillColor == color.Transparent || arc.FillColor == nil) &&
		(arc.StrokeColor == color.Transparent || arc.StrokeColor == nil || arc.StrokeWidth == 0)) ||
		arc.StartAngle == arc.EndAngle {
		return
	}

	points, bounds := p.vecRectCoords(pos, arc, frame, 0, canvas.Shadow{})
	c := p.baseConstants(frame, bounds, arc.FillColor, arc.StrokeColor, canvas.Shadow{})

	size := arc.Size()
	outerRadius := fyne.Min(size.Width, size.Height) / 2
	innerRadius := outerRadius * float32(math.Min(1.0, math.Max(0.0, float64(arc.CutoutRatio))))
	cornerRadius := fyne.Min(
		paint.GetMaximumRadiusArc(outerRadius, innerRadius, arc.EndAngle-arc.StartAngle),
		arc.CornerRadius,
	)

	c.Radius = [4]float32{
		roundToPixel(innerRadius*p.pixScale, 1.0),
		roundToPixel(outerRadius*p.pixScale, 1.0),
		roundToPixel(cornerRadius*p.pixScale, 1.0),
		0,
	}

	// PixelArc samples no texture, so texParams.xy carries the angles in degrees.
	startAngle, endAngle := paint.NormalizeArcAngles(arc.StartAngle, arc.EndAngle)
	c.TexParams = [4]float32{startAngle, endAngle, 0, 0}

	c.Misc[0] = roundToPixel(arc.StrokeWidth*p.pixScale, 1.0)

	p.upload(&c, points, p.vsQuad, p.psArc, topologyTriangleStrip, p.blend, 4)
}

func (p *Painter) drawLine(line *canvas.Line, pos fyne.Position, frame fyne.Size) {
	if line.StrokeColor == color.Transparent || line.StrokeColor == nil || line.StrokeWidth == 0 {
		return
	}
	ndc, normal, halfWidth, feather := p.lineGeometry(pos, line.Position1, line.Position2, line.StrokeWidth, 0.5, frame)
	if halfWidth == 0 {
		return // degenerate: both endpoints on the same pixel
	}

	c := constants{}
	r, g, b, a := fragmentColor(line.StrokeColor)
	c.FillColor = [4]float32{r, g, b, a}
	c.Misc[0] = halfWidth
	c.ShadowOffset[3] = feather
	// The line vertex shader reads the edge normal from inset, which lines
	// never use for texturing.
	c.Inset[0], c.Inset[1] = normal[0], normal[1]

	p.upload(&c, ndc, p.vsLine, p.psLine, topologyTriangleList, p.blend, 6)
}

func (p *Painter) drawText(text *canvas.Text, pos fyne.Position, frame fyne.Size, clip *internal.ClipItem) {
	if text.Text == "" {
		return
	}
	decorated := text.TextStyle.Underline || text.TextStyle.Strikethrough
	if text.Text == " " && !decorated {
		return
	}

	size := text.MinSize()
	containerSize := text.Size()
	switch text.Alignment {
	case fyne.TextAlignTrailing:
		pos = fyne.NewPos(pos.X+containerSize.Width-size.Width, pos.Y)
	case fyne.TextAlignCenter:
		pos = fyne.NewPos(pos.X+(containerSize.Width-size.Width)/2, pos.Y)
	}
	if containerSize.Height > size.Height {
		pos = fyne.NewPos(pos.X, pos.Y+(containerSize.Height-size.Height)/2)
	}

	// Text is sensitive to its position on screen: the quad has to land on whole
	// pixels or the glyph texture gets resampled and turns blurry. The padding is
	// what the texture was generated with - italic overspill on the right, room for
	// descenders and underlines below - and is snapped the same way.
	size.Width = roundToPixel(size.Width, p.pixScale)
	size.Height = roundToPixel(size.Height, p.pixScale)
	size.Width += roundToPixel(paint.VectorPad(text), p.pixScale)
	size.Height += roundToPixel(paint.TextVectorPad, p.pixScale)

	// The atlas draws a quad per glyph out of one shared texture, so it needs no
	// texture of its own and no width limit - both of the paths below exist only
	// for strings it cannot represent.
	if p.drawTextFromAtlas(text, pos, frame) {
		// An object that used to be too wide for one texture may have a windowed
		// texture parked from before; the atlas has no width limit, so it is dead.
		p.freeClippedTextTexture(text)
		if decorated {
			p.drawTextDecoration(text, pos, size, frame)
		}
		return
	}

	// A run wider than the device's texture limit cannot upload whole; render
	// only a window around the visible part, the way the GL painter does.
	fullWidth := int(math.Ceil(float64(size.Width * p.pixScale)))
	if maxSize := p.g.MaxTextureSize(); fullWidth <= maxSize || maxSize <= 0 {
		p.freeClippedTextTexture(text)
		p.drawTexture(text, p.textTexture, pos, size, frame, canvas.ImageFillStretch, 1, 0)
	} else {
		visibleOffset, visibleWidth := visibleTextPixels(pos, size, frame, clip, p.pixScale)
		height := int(math.Ceil(float64(size.Height * p.pixScale)))
		cached := p.clippedTextTexture(text, visibleOffset, visibleWidth, fullWidth, height)
		if cached.tex != nil {
			clipPos := fyne.NewPos(pos.X+float32(cached.offset)/p.pixScale, pos.Y)
			clipSize := fyne.NewSize(float32(cached.width)/p.pixScale, size.Height)
			p.drawGPUTexture(text, cached.tex, clipPos, clipSize, frame, canvas.ImageFillStretch, 1, 1)
		}
	}

	if !decorated {
		return
	}
	p.drawTextDecoration(text, pos, size, frame)
}

func (p *Painter) drawTextDecoration(text *canvas.Text, pos fyne.Position, size, frame fyne.Size) {
	_, baseline := cache.GetFontMetrics(text.Text, text.TextSize, text.TextStyle, text.FontSource)
	line := canvas.NewLine(text.Color)
	line.Resize(fyne.NewSize(size.Width, 0))
	if text.TextStyle.Underline {
		p.drawLine(line, fyne.NewPos(pos.X, pos.Y+baseline+paint.UnderlineOffsetFromBaseline), frame)
	}
	if text.TextStyle.Strikethrough {
		p.drawLine(line, fyne.NewPos(pos.X, pos.Y+baseline*paint.StrikethroughToBaselineFactor), frame)
	}
}

// shapedKey identifies one shaped string: the text and everything else the
// shaper reads. The face is the cached *FontCacheItem pointer, whose identity
// is stable until the font caches are cleared - a theme or font change hands
// out new pointers, orphaning (not corrupting) old entries.
type shapedKey struct {
	text  string
	face  *paint.FontCacheItem
	size  float32
	scale float32
	style fyne.TextStyle
}

// shapedEntry is one cached shaping result. ok is false when the string has to
// take the whole-run texture path (an emoji, in practice), cached so the
// losing shape is not re-run every frame just to fail again.
type shapedEntry struct {
	glyphs []paint.PlacedGlyph
	ok     bool
}

// shapedCacheMax caps the shaped-run cache. A UI whose labels churn (counters,
// live values) grows entries without bound; past the cap the whole cache is
// dropped and rebuilt from live strings, the same reset-and-repack answer the
// atlas uses. 4096 entries of a dozen glyphs is roughly 2MB.
const shapedCacheMax = 4096

// drawTextFromAtlas queues one quad per glyph of the run, reporting false when
// the string has to go down the whole-run texture path instead.
func (p *Painter) drawTextFromAtlas(text *canvas.Text, pos fyne.Position, frame fyne.Size) bool {
	col := text.Color
	if col == nil {
		col = theme.Color(theme.ColorNameForeground)
	}

	face := paint.CachedFontFace(text.TextStyle, text.FontSource, text)
	key := shapedKey{text: text.Text, face: face, size: text.TextSize,
		scale: p.pixScale, style: text.TextStyle}
	entry, ok := p.shaped[key]
	if !ok {
		p.glyphScratch = p.glyphScratch[:0]
		paint.WalkGlyphs(face.Fonts, text.Text, text.TextSize, p.pixScale, text.TextStyle,
			func(g paint.PlacedGlyph) { p.glyphScratch = append(p.glyphScratch, g) })
		entry.ok = len(p.glyphScratch) > 0 && shapeable(p.glyphScratch)
		if entry.ok {
			entry.glyphs = append([]paint.PlacedGlyph(nil), p.glyphScratch...)
		}
		if len(p.shaped) >= shapedCacheMax {
			clear(p.shaped)
		}
		p.shaped[key] = entry
	}
	if !entry.ok {
		return false
	}

	// The run's own pixel origin. WalkGlyphs positions are relative to it, and
	// the same rounding the texture path applies to the quad is applied here so
	// glyphs land on whole pixels either way.
	originX := roundToPixel(pos.X*p.pixScale, 1.0)
	originY := roundToPixel(pos.Y*p.pixScale, 1.0)

	// Every glyph is resolved before any quad is queued. Resolving can reset a
	// full atlas, and failing part way through after queueing would leave half a
	// string on screen with no way to take it back.
	p.entryScratch = p.entryScratch[:0]
	for _, pg := range entry.glyphs {
		e, ok := p.glyphEntry(pg, text.TextSize)
		if !ok {
			return false
		}
		p.entryScratch = append(p.entryScratch, e)
	}

	r, g, b, a := fragmentColor(col)
	fw, fh := p.scaleFrameSize(frame)

	for i, pg := range entry.glyphs {
		e := p.entryScratch[i]
		if e.empty() {
			continue
		}

		// Snapped to whole device pixels. The bitmap was rasterised at a whole
		// pixel origin, so a quad at a fractional position would resample it and
		// the glyph would come out soft and fringed - the pen position drifts
		// fractional as the shaped advances accumulate, so this is every glyph
		// but the first, not an edge case.
		//
		// ponytail: whole-pixel placement, so within-run subpixel positioning is
		// lost and spacing can differ from the whole-run path by under a pixel.
		// The upgrade, if that is ever visible, is to cache each glyph at a few
		// horizontal subpixel phases and pick by the fractional pen position.
		x1 := roundToPixel(originX+pg.X, 1.0) + float32(e.bearX)
		y1 := roundToPixel(originY+pg.Y, 1.0) + float32(e.bearY)
		x2, y2 := x1+float32(e.width), y1+float32(e.height)

		p.pushGlyph(glyphInst{
			// Pixels to clip space: x doubles and shifts, y additionally flips,
			// which is why y1 becomes the larger clip-space value.
			NDC: [4]float32{
				x1/fw*2 - 1, 1 - y1/fh*2,
				x2/fw*2 - 1, 1 - y2/fh*2,
			},
			UV: [4]float32{
				float32(e.x) / atlasSize, float32(e.y) / atlasSize,
				float32(e.x+e.width) / atlasSize, float32(e.y+e.height) / atlasSize,
			},
			Color: [4]float32{r, g, b, a},
		})
	}
	return true
}

// clippedTextEntry is one windowed texture for a text run wider than the device
// texture limit, remembering which pixel range of the full run it covers.
type clippedTextEntry struct {
	tex           *gpuTexture
	offset        int
	width, height int
	scale         float32
}

func (t clippedTextEntry) covers(offset, width, height int, scale float32) bool {
	return t.height == height && t.scale == scale &&
		t.offset <= offset && t.offset+t.width >= offset+width
}

// visibleTextPixels reports the horizontal pixel range of the text run that is
// inside the clip (or frame), in texture pixels from the run's left edge.
func visibleTextPixels(pos fyne.Position, size, frame fyne.Size, clip *internal.ClipItem, scale float32) (offset, width int) {
	clipPos := fyne.Position{}
	clipSize := frame
	if clip != nil {
		clipPos, clipSize = clip.Rect()
	}

	left := fyne.Max(pos.X, clipPos.X)
	right := fyne.Min(pos.X+size.Width, clipPos.X+clipSize.Width)
	if right <= left {
		return 0, 0
	}

	offset = int(math.Floor(float64((left - pos.X) * scale)))
	width = int(math.Ceil(float64((right-pos.X)*scale))) - offset
	return offset, width
}

// textTextureWindow centres a maxWidth-wide window on the visible range so
// moderate scrolling reuses the texture instead of re-rendering every frame.
func textTextureWindow(visibleOffset, visibleWidth, fullWidth, maxWidth int) (offset, width int) {
	width = maxWidth
	if fullWidth < width {
		width = fullWidth
	}
	if visibleWidth > width {
		visibleWidth = width
	}
	offset = visibleOffset - (width-visibleWidth)/2
	if offset < 0 {
		offset = 0
	}
	if maxOffset := fullWidth - width; offset > maxOffset {
		offset = maxOffset
	}
	return offset, width
}

// clippedTextTexture returns a texture covering the visible window of an
// over-wide text run, rendering a new one when the cached window no longer
// covers what is visible.
func (p *Painter) clippedTextTexture(text *canvas.Text, visibleOffset, visibleWidth, fullWidth, height int) clippedTextEntry {
	if cached, ok := p.clippedTextTextures[text]; ok {
		if cached.covers(visibleOffset, visibleWidth, height, p.pixScale) {
			cache.GetTexture(text) // keep the expiry marker alive while still in use
			return cached
		}
		p.releaseTexture(cached.tex)
		delete(p.clippedTextTextures, text)
	}

	offset, width := textTextureWindow(visibleOffset, visibleWidth, fullWidth, p.g.MaxTextureSize())

	col := text.Color
	if col == nil {
		col = theme.Color(theme.ColorNameForeground)
	}
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	face := paint.CachedFontFace(text.TextStyle, text.FontSource, text)
	paint.DrawStringOffset(img, text.Text, col, face.Fonts, text.TextSize, p.pixScale, text.TextStyle, offset)
	tex := p.imgToTexture(img)
	if tex == nil {
		return clippedTextEntry{}
	}

	if p.clippedTextTextures == nil {
		p.clippedTextTextures = make(map[*canvas.Text]clippedTextEntry)
	}
	cached := clippedTextEntry{tex: tex, offset: offset, width: width, height: height, scale: p.pixScale}
	p.clippedTextTextures[text] = cached
	// A no-texture marker in the shared cache makes FreeDirtyTextures call Free
	// for this text when it refreshes or expires, which drops the entry above.
	cache.SetTexture(text, cache.NoTexture, p.canvas)
	return cached
}

func (p *Painter) freeClippedTextTexture(text *canvas.Text) {
	cached, ok := p.clippedTextTextures[text]
	if !ok {
		return
	}
	p.releaseTexture(cached.tex)
	delete(p.clippedTextTextures, text)
}

// drawTexture is the shared path for every texture-backed object.
func (p *Painter) drawTexture(obj fyne.CanvasObject, creator func(fyne.CanvasObject) *gpuTexture,
	pos fyne.Position, size, frame fyne.Size, fill canvas.ImageFill, alpha, aspect float32,
) {
	tex := p.getTexture(obj, creator)
	if tex == nil {
		return
	}
	p.drawGPUTexture(obj, tex, pos, size, frame, fill, alpha, aspect)
}

// drawGPUTexture draws an already-uploaded texture as a quad.
func (p *Painter) drawGPUTexture(obj fyne.CanvasObject, tex *gpuTexture,
	pos fyne.Position, size, frame fyne.Size, fill canvas.ImageFill, alpha, aspect float32,
) {
	points, insets, inner := p.rectCoords(size, pos, frame, fill, aspect, 0)

	c := constants{}
	fw, fh := p.scaleFrameSize(frame)
	c.Frame = [4]float32{fw, fh, 0, 0}
	// The shader rounds corners against the drawn quad, so it needs the inner size
	// rectCoords actually used - which differs from `size` for contain/original fills.
	c.TexParams = [4]float32{alpha, 0, inner.Width * p.pixScale, inner.Height * p.pixScale}
	c.Inset = insets

	if img, ok := obj.(*canvas.Image); ok && img.CornerRadius > 0 {
		c.TexParams[1] = fyne.Min(paint.GetMaximumRadius(size), img.CornerRadius) * p.pixScale
	}

	// Before binding, not after: flushing draws the glyph batch, which leaves the
	// atlas in shader slot 0. Binding first would have this draw sample the atlas.
	p.flushGlyphs()

	sampler := p.sampLinear
	switch o := obj.(type) {
	case *canvas.Image:
		if o.ScaleMode == canvas.ImageScalePixels {
			sampler = p.sampNearest
		}
	case *canvas.Raster:
		if o.ScaleMode == canvas.ImageScalePixels {
			sampler = p.sampNearest
		}
	}
	if tex.srv != p.lastSRV {
		p.g.ctx.PSSetShaderResource(tex.srv)
		p.lastSRV = tex.srv
	}
	if sampler != p.lastSampler {
		p.g.ctx.PSSetSampler(sampler)
		p.lastSampler = sampler
	}

	p.upload(&c, points, p.vsQuad, p.psTextured, topologyTriangleStrip, p.blendPremul, 4)
}

// ---------------------------------------------------------------------------
// Coordinate maths - ported from internal/Painter/gl/draw.go so both painters
// place geometry identically.
// ---------------------------------------------------------------------------

func (p *Painter) vecRectCoords(pos fyne.Position, obj fyne.CanvasObject, frame fyne.Size,
	aspect float32, shadow canvas.Shadow,
) ([4]float32, [4]float32) {
	xPad, yPad := float32(0), float32(0)
	if aspect != 0 {
		inner := obj.Size()
		frameAspect := inner.Width / inner.Height
		if frameAspect > aspect {
			xPad = (inner.Width - inner.Height*aspect) / 2
		} else if frameAspect < aspect {
			yPad = (inner.Height - inner.Width/aspect) / 2
		}
	}

	size := obj.Size()
	pos1 := obj.Position()

	xPosDiff := pos.X - pos1.X + xPad
	yPosDiff := pos.Y - pos1.Y + yPad
	pos1.X = roundToPixel(pos1.X+xPosDiff, p.pixScale)
	pos1.Y = roundToPixel(pos1.Y+yPosDiff, p.pixScale)
	size.Width = roundToPixel(size.Width-2*xPad, p.pixScale)
	size.Height = roundToPixel(size.Height-2*yPad, p.pixScale)

	pads := paint.GetShadowPaddings(shadow)
	padLeft := roundToPixel(pads[0], p.pixScale)
	padTop := roundToPixel(pads[1], p.pixScale)
	padRight := roundToPixel(pads[2], p.pixScale)
	padBottom := roundToPixel(pads[3], p.pixScale)

	softness := roundToPixel(edgeSoftness*p.pixScale, 1.0)
	x1Pos := pos1.X
	x1Norm := -1 + (x1Pos-softness-padLeft)*2/frame.Width
	x2Pos := pos1.X + size.Width
	x2Norm := -1 + (x2Pos+softness+padRight)*2/frame.Width
	y1Pos := pos1.Y
	y1Norm := 1 - (y1Pos-softness-padTop)*2/frame.Height
	y2Pos := pos1.Y + size.Height
	y2Norm := 1 - (y2Pos+softness+padBottom)*2/frame.Height

	return [4]float32{x1Norm, y1Norm, x2Norm, y2Norm}, [4]float32{x1Pos, y1Pos, x2Pos, y2Pos}
}

// rectCoords returns the clip-space rectangle of a textured quad, plus the
// texture insets used by the sampling and rounded-corner paths.
func (p *Painter) rectCoords(size fyne.Size, pos fyne.Position, frame fyne.Size,
	fill canvas.ImageFill, aspect, pad float32,
) ([4]float32, [4]float32, fyne.Size) {
	innerSize, innerPos := rectInnerCoords(size, pos, fill, aspect)
	pixelSize, pixelPos := roundToPixelCoords(innerSize, innerPos, p.pixScale)

	x1 := -1 + (pixelPos.X-pad)*2/frame.Width
	x2 := -1 + (pixelPos.X+pixelSize.Width+pad)*2/frame.Width
	y1 := 1 - (pixelPos.Y-pad)*2/frame.Height
	y2 := 1 - (pixelPos.Y+pixelSize.Height+pad)*2/frame.Height

	xInset, yInset := float32(0), float32(0)
	if fill == canvas.ImageFillCover {
		viewAspect := pixelSize.Width / pixelSize.Height
		if viewAspect > aspect {
			newHeight := pixelSize.Width / aspect
			yInset = ((newHeight - pixelSize.Height) / 2) / newHeight
		} else if viewAspect < aspect {
			newWidth := pixelSize.Height * aspect
			xInset = ((newWidth - pixelSize.Width) / 2) / newWidth
		}
	}
	insets := [4]float32{xInset, yInset, 1 - xInset, 1 - yInset}

	return [4]float32{x1, y1, x2, y2}, insets, innerSize
}

func rectInnerCoords(size fyne.Size, pos fyne.Position, fill canvas.ImageFill, aspect float32) (fyne.Size, fyne.Position) {
	if fill != canvas.ImageFillContain && fill != canvas.ImageFillOriginal {
		return size, pos
	}
	viewAspect := size.Width / size.Height
	if viewAspect == aspect {
		return size, pos
	}
	newWidth, newHeight := size.Width, size.Height
	newX, newY := pos.X, pos.Y
	if viewAspect > aspect {
		newWidth = size.Height * aspect
		newX += (size.Width - newWidth) / 2
	} else if viewAspect < aspect {
		newHeight = size.Width / aspect
		newY += (size.Height - newHeight) / 2
	}
	return fyne.NewSize(newWidth, newHeight), fyne.NewPos(newX, newY)
}

// lineGeometry converts a line to the clip-space endpoints and outward edge
// normal that VertexLine expands into a quad, scaled by the half width.
func (p *Painter) lineGeometry(pos, pos1, pos2 fyne.Position, lineWidth, feather float32,
	frame fyne.Size,
) (ndc [4]float32, normal [2]float32, halfWidth, featherWidth float32) {
	xPosDiff := pos.X - fyne.Min(pos1.X, pos2.X)
	yPosDiff := pos.Y - fyne.Min(pos1.Y, pos2.Y)
	pos1.X = roundToPixel(pos1.X+xPosDiff, p.pixScale)
	pos1.Y = roundToPixel(pos1.Y+yPosDiff, p.pixScale)
	pos2.X = roundToPixel(pos2.X+xPosDiff, p.pixScale)
	pos2.Y = roundToPixel(pos2.Y+yPosDiff, p.pixScale)

	if lineWidth <= 1 {
		offset := float32(0.5)
		if lineWidth <= 0.5 && p.pixScale > 1 {
			offset = 0.25
		}
		if pos1.X == pos2.X {
			pos1.X -= offset
			pos2.X -= offset
		}
		if pos1.Y == pos2.Y {
			pos1.Y -= offset
			pos2.Y -= offset
		}
	}

	x1 := -1 + pos1.X*2/frame.Width
	y1 := 1 - pos1.Y*2/frame.Height
	x2 := -1 + pos2.X*2/frame.Width
	y2 := 1 - pos2.Y*2/frame.Height

	normalX := (pos2.Y - pos1.Y) / frame.Width
	normalY := (pos2.X - pos1.X) / frame.Height
	dirLength := float32(math.Sqrt(float64(normalX*normalX + normalY*normalY)))
	if dirLength == 0 {
		return ndc, normal, 0, 0
	}
	normalX /= dirLength
	normalY /= dirLength

	normalObjX := normalX * 0.5 * frame.Width
	normalObjY := normalY * 0.5 * frame.Height
	widthMultiplier := float32(math.Sqrt(float64(normalObjX*normalObjX + normalObjY*normalObjY)))
	halfWidth = (roundToPixel(lineWidth+feather, p.pixScale) * 0.5) / widthMultiplier
	featherWidth = feather / widthMultiplier

	return [4]float32{x1, y1, x2, y2}, [2]float32{normalX, normalY}, halfWidth, featherWidth
}

func roundToPixel(v, pixScale float32) float32 {
	if pixScale == 1.0 {
		return float32(math.Round(float64(v)))
	}
	return float32(math.Round(float64(v*pixScale))) / pixScale
}

func roundToPixelCoords(size fyne.Size, pos fyne.Position, pixScale float32) (fyne.Size, fyne.Position) {
	end := pos.Add(size)
	end.X = roundToPixel(end.X, pixScale)
	end.Y = roundToPixel(end.Y, pixScale)
	pos.X = roundToPixel(pos.X, pixScale)
	pos.Y = roundToPixel(pos.Y, pixScale)
	return fyne.NewSize(end.X-pos.X, end.Y-pos.Y), pos
}

func (p *Painter) textureScale(v float32) float32 {
	if p.pixScale == 1.0 {
		return float32(math.Round(float64(v)))
	}
	return float32(math.Round(float64(v * p.pixScale)))
}

func (p *Painter) scaleFrameSize(frame fyne.Size) (width, height float32) {
	return roundToPixel(frame.Width*p.pixScale, 1.0), roundToPixel(frame.Height*p.pixScale, 1.0)
}

func (p *Painter) scaleRectCoords(x1, x2, y1, y2 float32) (float32, float32, float32, float32) {
	return roundToPixel(x1*p.pixScale, 1.0), roundToPixel(x2*p.pixScale, 1.0),
		roundToPixel(y1*p.pixScale, 1.0), roundToPixel(y2*p.pixScale, 1.0)
}

// fragmentColor splits a colour into un-premultiplied RGB plus alpha, matching
// getFragmentColor in the GL Painter.
func fragmentColor(col color.Color) (r, g, b, a float32) {
	if col == nil {
		return 0, 0, 0, 0
	}
	cr, cg, cb, ca := col.RGBA()
	if ca == 0 {
		return 0, 0, 0, 0
	}
	alpha := float32(ca)
	return float32(cr) / alpha, float32(cg) / alpha, float32(cb) / alpha, alpha / 0xffff
}

// ---------------------------------------------------------------------------
// Capture
// ---------------------------------------------------------------------------

// Capture reads the back buffer back to the CPU through a staging texture, the
// only resource type D3D11 allows the CPU to map for reading.
func (p *Painter) Capture(c fyne.Canvas) image.Image {
	width, height := c.PixelCoordinateForPosition(fyne.NewPos(c.Size().Width, c.Size().Height))
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	if !p.ready() || p.g.back == nil {
		return img
	}

	// The swap chain uses DXGI_SWAP_EFFECT_DISCARD, which leaves the back buffer
	// contents undefined once Present has run - so whatever is sitting there now is
	// not the last frame, and reading it straight out yields garbage or black.
	// Redraw into it first. (The GL painter avoids this by reading the front
	// buffer, which D3D11 gives no access to.)
	if r, ok := c.(Repainter); ok {
		r.Repaint(c.Size())
	}

	staging, err := p.g.dev.CreateTexture2D(&texture2DDesc{
		Width: p.g.width, Height: p.g.height, MipLevels: 1, ArraySize: 1,
		Format: formatB8G8R8A8Unorm, SampleDesc: dxgiSampleDesc{Count: 1},
		Usage: usageStaging, CPUAccessFlags: cpuAccessRead,
	}, nil)
	if err != nil {
		fyne.LogError("directx: capture staging texture", err)
		return img
	}
	defer staging.Release()

	p.g.ctx.CopyResource(staging, p.g.back)
	m, err := p.g.ctx.Map(unsafe.Pointer(staging), mapRead)
	if err != nil {
		fyne.LogError("directx: capture map", err)
		return img
	}
	defer p.g.ctx.Unmap(unsafe.Pointer(staging))

	src := unsafe.Slice((*byte)(m.Data), int(m.RowPitch)*int(p.g.height))
	for y := 0; y < height && y < int(p.g.height); y++ {
		row := src[y*int(m.RowPitch):]
		out := img.Pix[y*img.Stride:]
		for x := 0; x < width && x < int(p.g.width); x++ {
			// The back buffer is BGRA; Go's RGBA wants the first two swapped.
			out[x*4+0] = row[x*4+2]
			out[x*4+1] = row[x*4+1]
			out[x*4+2] = row[x*4+0]
			out[x*4+3] = row[x*4+3]
		}
	}
	return img
}

// ---------------------------------------------------------------------------
// Textures
// ---------------------------------------------------------------------------

// releaseTexture retires a texture: poolable ones (created by imgToTexture)
// are parked for same-size reuse, the rest are destroyed outright.
func (p *Painter) releaseTexture(t *gpuTexture) {
	if t == nil {
		return
	}
	if t.width > 0 && p.texPoolSize < texPoolMax {
		key := [2]uint32{t.width, t.height}
		p.texPool[key] = append(p.texPool[key], pooledTexture{tex: t, tick: p.frameTick})
		p.texPoolSize++
		return
	}
	p.destroyTexture(t)
}

// destroyTexture frees a texture's COM objects and drops it from the bound-SRV
// cache: a later allocation could reuse the same address, and a stale cache hit
// would skip a bind the context actually needs.
func (p *Painter) destroyTexture(t *gpuTexture) {
	if t.srv == p.lastSRV {
		p.lastSRV = nil
	}
	t.release()
}

// pooledTextureFor pops a parked texture of exactly the given size, most
// recently parked first.
func (p *Painter) pooledTextureFor(width, height uint32) *gpuTexture {
	key := [2]uint32{width, height}
	entries := p.texPool[key]
	if len(entries) == 0 {
		return nil
	}
	t := entries[len(entries)-1].tex
	if len(entries) == 1 {
		delete(p.texPool, key)
	} else {
		p.texPool[key] = entries[:len(entries)-1]
	}
	p.texPoolSize--
	return t
}

// sweepTexPool destroys parked textures whose size has not recurred within
// texPoolTTL frames, so a burst of odd sizes (say, mid-resize plots) does not
// pin VRAM forever.
func (p *Painter) sweepTexPool() {
	for key, entries := range p.texPool {
		kept := entries[:0]
		for _, e := range entries {
			if p.frameTick-e.tick > texPoolTTL {
				p.destroyTexture(e.tex)
				p.texPoolSize--
			} else {
				kept = append(kept, e)
			}
		}
		if len(kept) == 0 {
			delete(p.texPool, key)
		} else {
			p.texPool[key] = kept
		}
	}
}

func (p *Painter) Free(obj fyne.CanvasObject) {
	if text, ok := obj.(*canvas.Text); ok {
		p.freeClippedTextTexture(text)
	}
	id, ok := cache.GetTexture(obj)
	if !ok {
		return
	}
	if tex, found := p.textures[id]; found {
		p.releaseTexture(tex)
		delete(p.textures, id)
	}
	cache.DeleteTexture(obj)
}

// freeAll releases every uploaded texture, used when the device goes away.
func (p *Painter) freeAll() {
	for id, tex := range p.textures {
		if tex != nil {
			p.destroyTexture(tex)
		}
		delete(p.textures, id)
	}
	for text, cached := range p.clippedTextTextures {
		p.destroyTexture(cached.tex)
		delete(p.clippedTextTextures, text)
	}
	for key, entries := range p.texPool {
		for _, e := range entries {
			p.destroyTexture(e.tex)
		}
		delete(p.texPool, key)
	}
	p.texPoolSize = 0
}

// getTexture returns the cached upload for obj, creating it on first use. Text
// runs share a separate cache keyed by content and style rather than by object.
func (p *Painter) getTexture(obj fyne.CanvasObject, creator func(fyne.CanvasObject) *gpuTexture) *gpuTexture {
	if t, ok := obj.(*canvas.Text); ok {
		custom := ""
		if t.FontSource != nil {
			custom = t.FontSource.Name()
		}
		ent := cache.FontCacheEntry{Color: t.Color, Canvas: p.canvas}
		ent.Text = t.Text
		ent.Size = t.TextSize
		ent.Style = t.TextStyle
		ent.Source = custom

		if id, ok := cache.GetTextTexture(ent); ok {
			if tex := p.textures[id]; tex != nil {
				return tex
			}
		}
		tex := creator(obj)
		if tex == nil {
			return nil
		}
		id := p.storeTexture(tex)
		cache.SetTextTexture(ent, id, p.canvas, func() {
			if t := p.textures[id]; t != nil {
				p.releaseTexture(t)
				delete(p.textures, id)
			}
		})
		return tex
	}

	if id, ok := cache.GetTexture(obj); ok {
		if tex, found := p.textures[id]; found {
			return tex // may be nil: a cached failure, until the object refreshes
		}
	}
	// A nil result is stored too, so a failing creator (e.g. a Raster generator
	// returning an empty image) does not re-run every frame.
	tex := creator(obj)
	id := p.storeTexture(tex)
	cache.SetTexture(obj, id, p.canvas)
	return tex
}

// storeTexture registers a texture and returns its cache handle. Ids start at 1
// because cache.NoTexture is 0.
func (p *Painter) storeTexture(tex *gpuTexture) uint32 {
	p.nextTexID++
	id := p.nextTexID
	p.textures[id] = tex
	return id
}

// imgToTexture uploads a Go image as an RGBA texture, refilling a pooled
// same-size texture when one is available instead of creating a new resource.
func (p *Painter) imgToTexture(img image.Image) *gpuTexture {
	if uni, ok := img.(*image.Uniform); ok {
		// A Uniform's bounds can be effectively infinite (canvas.Raster returns
		// one for non-overlapping rects), so upload a single stretched pixel.
		one := image.NewRGBA(image.Rect(0, 0, 1, 1))
		one.Set(0, 0, uni.C)
		img = one
	}
	rgba, ok := img.(*image.RGBA)
	if !ok {
		if img == nil {
			return nil
		}
		b := img.Bounds()
		conv := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
		draw.Draw(conv, conv.Rect, img, b.Min, draw.Over)
		rgba = conv
	}
	if len(rgba.Pix) == 0 {
		return nil
	}
	w := uint32(rgba.Rect.Dx())
	h := uint32(rgba.Rect.Dy())
	if w == 0 || h == 0 {
		return nil
	}

	if t := p.pooledTextureFor(w, h); t != nil {
		p.g.ctx.UpdateTexture2D(unsafe.Pointer(t.tex), unsafe.Pointer(&rgba.Pix[0]), uint32(rgba.Stride))
		return t
	}

	data := subresourceData{
		SysMem:      unsafe.Pointer(&rgba.Pix[0]),
		SysMemPitch: uint32(rgba.Stride),
	}
	// Default usage rather than immutable, so a pooled texture can be refilled.
	tex, err := p.g.dev.CreateTexture2D(&texture2DDesc{
		Width: w, Height: h, MipLevels: 1, ArraySize: 1,
		Format: formatR8G8B8A8Unorm, SampleDesc: dxgiSampleDesc{Count: 1},
		Usage: usageDefault, BindFlags: bindShaderResource,
	}, &data)
	if err != nil {
		fyne.LogError("directx: texture upload", err)
		return nil
	}
	srv, err := p.g.dev.CreateShaderResourceView(tex)
	if err != nil {
		tex.Release()
		fyne.LogError("directx: shader resource view", err)
		return nil
	}
	return &gpuTexture{tex: tex, srv: srv, width: w, height: h}
}

func (p *Painter) imageTexture(obj fyne.CanvasObject) *gpuTexture {
	img := obj.(*canvas.Image)
	width := p.textureScale(img.Size().Width)
	height := p.textureScale(img.Size().Height)
	tex := paint.PaintImage(img, p.canvas, int(width), int(height))
	if tex == nil {
		return nil
	}
	return p.imgToTexture(tex)
}

func (p *Painter) rasterTexture(obj fyne.CanvasObject) *gpuTexture {
	rast := obj.(*canvas.Raster)
	width := p.textureScale(rast.Size().Width)
	height := p.textureScale(rast.Size().Height)
	return p.imgToTexture(rast.Generator(int(width), int(height)))
}

func (p *Painter) linearGradientTexture(obj fyne.CanvasObject) *gpuTexture {
	gradient := obj.(*canvas.LinearGradient)
	w := gradient.Size().Width
	h := gradient.Size().Height
	// An axis-aligned gradient only varies along one axis; a 1px strip stretched
	// by the sampler renders identically for a fraction of the memory.
	switch a := gradient.Angle; {
	case almostEqual(a, geom.AngleQuarter), almostEqual(a, geom.AngleThreeQuarter):
		h = 1
	case almostEqual(a, geom.AngleNone), almostEqual(a, geom.AngleHalf):
		w = 1
	}
	width := p.textureScale(w)
	height := p.textureScale(h)
	return p.imgToTexture(gradient.Generate(int(width), int(height)))
}

func almostEqual(a, b float64) bool {
	return math.Abs(a-b) < 1e-9
}

func (p *Painter) radialGradientTexture(obj fyne.CanvasObject) *gpuTexture {
	gradient := obj.(*canvas.RadialGradient)
	width := p.textureScale(gradient.Size().Width)
	height := p.textureScale(gradient.Size().Height)
	return p.imgToTexture(gradient.Generate(int(width), int(height)))
}

func (p *Painter) textTexture(obj fyne.CanvasObject) *gpuTexture {
	text := obj.(*canvas.Text)
	col := text.Color
	if col == nil {
		col = theme.Color(theme.ColorNameForeground)
	}

	bounds := text.MinSize()
	width := int(math.Ceil(float64(p.textureScale(bounds.Width) + paint.VectorPad(text))))
	height := int(math.Ceil(float64(p.textureScale(bounds.Height) + paint.TextVectorPad)))
	if width <= 0 || height <= 0 {
		return nil
	}
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	face := paint.CachedFontFace(text.TextStyle, text.FontSource, text)
	paint.DrawString(img, text.Text, col, face.Fonts, text.TextSize, p.pixScale, text.TextStyle)
	return p.imgToTexture(img)
}

func (p *Painter) Release() {
	p.freeAll()
	p.releaseUserShaders()
	p.blurSnap.release()
	p.blurKernel.release()
	releaseCOM(&p.auxbuf)
	releaseCOM(&p.cbuf)
	releaseCOM(&p.glyphbuf)
	p.atlas.release()
	releaseCOM(&p.sampNearest)
	releaseCOM(&p.sampLinear)
	releaseCOM(&p.rsScissor)
	releaseCOM(&p.rsPlain)
	releaseCOM(&p.blendPremul)
	releaseCOM(&p.blendReplace)
	releaseCOM(&p.blend)
	releaseCOM(&p.psGlyph)
	releaseCOM(&p.psLine)
	releaseCOM(&p.psBlur)
	releaseCOM(&p.psBezier)
	releaseCOM(&p.psArbPoly)
	releaseCOM(&p.psPolygon)
	releaseCOM(&p.psArc)
	releaseCOM(&p.psEllipse)
	releaseCOM(&p.psRound)
	releaseCOM(&p.psRect)
	releaseCOM(&p.psTextured)
	releaseCOM(&p.vsGlyph)
	releaseCOM(&p.vsLine)
	releaseCOM(&p.vsQuad)
}

// ---------------------------------------------------------------------------
// Polygons, curves and blur
// ---------------------------------------------------------------------------

// maxPolygonVertices matches MAX_VERTICES in arbitrary_polygon.ps.hlsl.
const maxPolygonVertices = 16

// auxConstants mirrors the register(b1) buffers declared by the polygon and
// curve shaders. One buffer serves both: the curve shader declares only the
// first two float4s and simply ignores the rest.
//
// The polygon shader reads Verts as float4[16] with xy the vertex and z its
// corner radius, so the curve's start/end/control points overlay the first two
// entries exactly.
type auxConstants struct {
	Verts [maxPolygonVertices][4]float32
	Info  [4]float32 // x: vertex count
}

// replaceBlend overwrites the destination. The blur passes use it because each
// one replaces the exact rectangle it just sampled.
func replaceBlend() blendDesc {
	d := straightAlphaBlend()
	d.RenderTarget[0].SrcBlend = blendOne
	d.RenderTarget[0].DestBlend = blendZero
	d.RenderTarget[0].SrcBlendAlpha = blendOne
	d.RenderTarget[0].DestBlendAlpha = blendZero
	return d
}

// uploadAux writes the register(b1) buffer and binds it.
func (p *Painter) uploadAux(a *auxConstants) {
	p.g.ctx.UpdateSubresource(unsafe.Pointer(p.auxbuf), unsafe.Pointer(a), nil)
	p.g.ctx.PSSetConstantBufferAt(1, p.auxbuf)
}

func (p *Painter) drawPolygon(polygon *canvas.RegularPolygon, pos fyne.Position, frame fyne.Size) {
	if polygon.Sides < 3 ||
		((polygon.FillColor == color.Transparent || polygon.FillColor == nil) &&
			(polygon.StrokeColor == color.Transparent || polygon.StrokeColor == nil || polygon.StrokeWidth == 0)) {
		return
	}

	points, bounds := p.vecRectCoords(pos, polygon, frame, 0, canvas.Shadow{})
	c := p.baseConstants(frame, bounds, polygon.FillColor, polygon.StrokeColor, canvas.Shadow{})

	size := polygon.Size()
	outerRadius := fyne.Min(size.Width, size.Height) / 2
	cornerRadius := fyne.Min(paint.GetMaximumRadius(size), polygon.CornerRadius)

	c.Radius = [4]float32{
		roundToPixel(outerRadius*p.pixScale, 1.0),
		0,
		roundToPixel(cornerRadius*p.pixScale, 1.0),
		float32(polygon.Sides),
	}
	c.ShadowOffset[3] = polygon.Angle // rotation in degrees
	c.Misc[0] = roundToPixel(polygon.StrokeWidth*p.pixScale, 1.0)

	p.upload(&c, points, p.vsQuad, p.psPolygon, topologyTriangleStrip, p.blend, 4)
}

func (p *Painter) drawArbitraryPolygon(polygon *canvas.ArbitraryPolygon, pos fyne.Position, frame fyne.Size) {
	if len(polygon.Points) < 3 ||
		((polygon.FillColor == color.Transparent || polygon.FillColor == nil) &&
			(polygon.StrokeColor == color.Transparent || polygon.StrokeColor == nil || polygon.StrokeWidth == 0)) {
		return
	}

	points, bounds := p.vecRectCoords(pos, polygon, frame, 0, canvas.Shadow{})
	c := p.baseConstants(frame, bounds, polygon.FillColor, polygon.StrokeColor, canvas.Shadow{})
	c.Misc[0] = roundToPixel(polygon.StrokeWidth*p.pixScale, 1.0)

	num := int(fyne.Min(paint.ArbitraryPolygonVerticesMaximum, float32(len(polygon.Points))))
	size := polygon.Size()
	clamp := func(v fyne.Position) fyne.Position {
		return fyne.NewPos(
			fyne.Min(fyne.Max(v.X, 0), fyne.Max(size.Width, 0)),
			fyne.Min(fyne.Max(v.Y, 0), fyne.Max(size.Height, 0)),
		)
	}

	fixed := make([]fyne.Position, num)
	radii := make([]float32, num)
	for i := 0; i < num; i++ {
		px, py := polygon.Points[i].X, polygon.Points[i].Y
		if polygon.NormalizedPoints {
			px, py = px*size.Width, py*size.Height
		}
		fixed[i] = clamp(fyne.NewPos(px, py))
		if i < len(polygon.CornerRadii) {
			radii[i] = polygon.CornerRadii[i]
		}
	}
	radii = paint.GetMaximumCornerRadii(fixed, radii)

	var aux auxConstants
	aux.Info[0] = float32(num)
	for i := 0; i < num; i++ {
		aux.Verts[i] = [4]float32{
			roundToPixel(fixed[i].X*p.pixScale, 1.0),
			roundToPixel(fixed[i].Y*p.pixScale, 1.0),
			roundToPixel(radii[i]*p.pixScale, 1.0),
			0,
		}
	}
	p.uploadAux(&aux)

	p.upload(&c, points, p.vsQuad, p.psArbPoly, topologyTriangleStrip, p.blend, 4)
}

func (p *Painter) drawBezierCurve(curve *canvas.BezierCurve, pos fyne.Position, frame fyne.Size) {
	if curve.StrokeColor == color.Transparent || curve.StrokeColor == nil || curve.StrokeWidth == 0 {
		return
	}

	points, bounds := p.vecRectCoords(pos, curve, frame, 0, canvas.Shadow{})
	c := p.baseConstants(frame, bounds, color.Transparent, curve.StrokeColor, canvas.Shadow{})

	// Keep the stroke inside the object, matching the GL Painter.
	size := curve.Size()
	strokeWidth := fyne.Min(curve.StrokeWidth, fyne.Min(size.Width, size.Height))
	if strokeWidth < 1 {
		strokeWidth = 1
	}
	p1, p2, cp := paint.NormalizeBezierCurvePoints(curve.StartPoint, curve.EndPoint,
		curve.ControlPoints, size, strokeWidth/2.0)

	scaled := func(v fyne.Position) (float32, float32) {
		return roundToPixel(v.X*p.pixScale, 1.0), roundToPixel(v.Y*p.pixScale, 1.0)
	}

	var aux auxConstants
	sx, sy := scaled(p1)
	ex, ey := scaled(p2)
	aux.Verts[0] = [4]float32{sx, sy, ex, ey}
	if len(cp) >= 1 {
		aux.Verts[1][0], aux.Verts[1][1] = scaled(cp[0])
	}
	if len(cp) >= 2 {
		aux.Verts[1][2], aux.Verts[1][3] = scaled(cp[1])
	}
	p.uploadAux(&aux)

	c.TexParams[0] = fyne.Min(float32(len(cp)), 2)
	c.RectHalf[2] = roundToPixel(strokeWidth*p.pixScale, 1.0) * 0.5

	p.upload(&c, points, p.vsQuad, p.psBezier, topologyTriangleStrip, p.blend, 4)
}

// blurSnapshot is the copy of the back buffer that the blur passes sample. D3D11
// forbids sampling a texture that is currently bound as a render target, so the
// region has to be copied out before each pass - the same shape as the GL
// Painter's glCopyTexSubImage2D dance.
type blurSnapshot struct {
	gpuTexture
	width, height uint32
}

func (b *blurSnapshot) release() {
	b.gpuTexture.release()
	b.width, b.height = 0, 0
}

// ensure (re)allocates the snapshot when the blurred area changes size.
func (b *blurSnapshot) ensure(dev *device, width, height uint32) bool {
	if b.tex != nil && b.width == width && b.height == height {
		return true
	}
	b.release()

	tex, err := dev.CreateTexture2D(&texture2DDesc{
		Width: width, Height: height, MipLevels: 1, ArraySize: 1,
		// Must match the back buffer format for CopySubresourceRegion to be legal.
		Format: formatB8G8R8A8Unorm, SampleDesc: dxgiSampleDesc{Count: 1},
		Usage: usageDefault, BindFlags: bindShaderResource,
	}, nil)
	if err != nil {
		fyne.LogError("directx: blur snapshot texture", err)
		return false
	}
	srv, err := dev.CreateShaderResourceView(tex)
	if err != nil {
		tex.Release()
		fyne.LogError("directx: blur snapshot view", err)
		return false
	}
	b.tex, b.srv, b.width, b.height = tex, srv, width, height
	return true
}

// blurKernelTexture holds the Gaussian weights as a one pixel tall texture,
// rebuilt only when the radius changes.
type blurKernelTexture struct {
	gpuTexture
	radius float32
}

func (k *blurKernelTexture) release() {
	k.gpuTexture.release()
	k.radius = 0
}

func (k *blurKernelTexture) ensure(dev *device, radius float32) bool {
	if k.tex != nil && k.radius == radius {
		return true
	}

	values, ok := cache.GetBlurKernel(radius)
	if !ok {
		values = createBlurKernel(radius)
		cache.SetBlurKernel(radius, values)
	}
	if len(values) == 0 {
		return false
	}
	k.release()

	data := kernelToRGBA(values)
	tex, err := dev.CreateTexture2D(&texture2DDesc{
		Width: uint32(len(values)), Height: 1, MipLevels: 1, ArraySize: 1,
		Format: formatR8G8B8A8Unorm, SampleDesc: dxgiSampleDesc{Count: 1},
		Usage: usageImmutable, BindFlags: bindShaderResource,
	}, &subresourceData{SysMem: unsafe.Pointer(&data[0]), SysMemPitch: uint32(len(data))})
	if err != nil {
		fyne.LogError("directx: blur kernel texture", err)
		return false
	}
	srv, err := dev.CreateShaderResourceView(tex)
	if err != nil {
		tex.Release()
		fyne.LogError("directx: blur kernel view", err)
		return false
	}
	k.tex, k.srv, k.radius = tex, srv, radius
	return true
}

// createBlurKernel builds normalised Gaussian weights, matching createKernel in
// the GL Painter so both backends blur identically.
func createBlurKernel(radius float32) []float32 {
	sum := float32(0)
	length := int(radius)*2 + 1
	values := make([]float32, length)
	for i, x := 0, float64(-radius); i < length; i, x = i+1, x+1 {
		value := float32(math.Exp(-(x * x / 4 / float64(radius))))
		values[i] = value
		sum += value
	}
	for i := range values {
		values[i] /= sum
	}
	return values
}

// kernelToRGBA quantises the weights into RGBA bytes; the shader reads .r.
func kernelToRGBA(values []float32) []uint8 {
	data := make([]uint8, len(values)*4)
	for i, v := range values {
		b := uint8(v*math.MaxUint8 + 0.5)
		off := i * 4
		data[off+0], data[off+1], data[off+2] = b, b, b
		data[off+3] = math.MaxUint8
	}
	return data
}

// drawBlur blurs whatever is already on screen behind the object, in two
// separable passes. Each pass snapshots the region, then overwrites it.
func (p *Painter) drawBlur(b *canvas.Blur, pos fyne.Position, frame fyne.Size) {
	if b.Radius == 0 || p.g.back == nil {
		return
	}
	radius := b.Radius * p.pixScale

	x := int32(roundToPixel(pos.X*p.pixScale, 1.0))
	y := int32(roundToPixel(pos.Y*p.pixScale, 1.0))
	bw := int32(roundToPixel(b.Size().Width*p.pixScale, 1.0))
	bh := int32(roundToPixel(b.Size().Height*p.pixScale, 1.0))

	// Clamp to the back buffer: CopySubresourceRegion rejects out-of-bounds
	// boxes, and a partially offscreen blur should still blur its visible strip.
	if x < 0 {
		bw += x
		x = 0
	}
	if y < 0 {
		bh += y
		y = 0
	}
	if m := int32(p.g.width) - x; bw > m {
		bw = m
	}
	if m := int32(p.g.height) - y; bh > m {
		bh = m
	}
	if bw <= 0 || bh <= 0 {
		return
	}
	pos = fyne.NewPos(float32(x)/p.pixScale, float32(y)/p.pixScale)
	size := fyne.NewSize(float32(bw)/p.pixScale, float32(bh)/p.pixScale)

	// Cap the kernel at 101 samples per pass. Beyond that the samples are spread
	// wider instead, downsampling but staying smooth thanks to bilinear filtering.
	kernelRadius := radius
	sampleScale := float32(1)
	const maxKernelRadius = 50.0
	if kernelRadius > maxKernelRadius {
		sampleScale = kernelRadius / maxKernelRadius
		kernelRadius = maxKernelRadius
	}

	if !p.blurSnap.ensure(p.g.dev, uint32(bw), uint32(bh)) {
		return
	}
	if !p.blurKernel.ensure(p.g.dev, kernelRadius) {
		return
	}

	points, _, _ := p.rectCoords(size, pos, frame, canvas.ImageFillStretch, 1.0, 0)

	cornerRadius := fyne.Min(paint.GetMaximumRadius(b.Size()), b.CornerRadius)

	c := constants{}
	fw, fh := p.scaleFrameSize(frame)
	c.Frame = [4]float32{fw, fh, 0, 0}
	c.TexParams = [4]float32{
		kernelRadius,
		roundToPixel(cornerRadius*p.pixScale, 1.0),
		float32(bw), float32(bh),
	}

	region := box{
		Left: uint32(x), Top: uint32(y), Front: 0,
		Right: uint32(x + bw), Bottom: uint32(y + bh), Back: 1,
	}

	// Ahead of both the binding and the snapshot below: queued glyphs have to be
	// in the back buffer before it is copied, or the blur samples a frame that
	// is missing its text.
	p.flushGlyphs()

	p.g.ctx.PSSetSamplersAt(0, []*samplerState{p.sampLinear, p.sampNearest})
	views := []*shaderResourceView{p.blurSnap.srv, p.blurKernel.srv}

	// The vertex shader owns inset as the texture-coordinate source, so the
	// blur direction rides in rectHalf instead, which blur never uses.
	c.Inset = [4]float32{0, 0, 1, 1}

	// Horizontal then vertical. Each pass reads a fresh snapshot of the region and
	// replaces it, so the second pass sees the first pass's output.
	for _, direction := range [2][2]float32{{1 / float32(bw), 0}, {0, 1 / float32(bh)}} {
		p.g.ctx.CopySubresourceRegion(p.blurSnap.tex, p.g.back, &region)
		p.g.ctx.PSSetShaderResourcesAt(0, views)

		c.RectHalf = [4]float32{direction[0], direction[1], sampleScale, 0}
		p.upload(&c, points, p.vsQuad, p.psBlur, topologyTriangleStrip, p.blendReplace, 4)
	}

	// Leave only the default texture bound so a later draw cannot sample the
	// snapshot by accident. Slots 0-1 changed behind the hot path's back; keep
	// its cache truthful: SRV slot 0 is now nil, sampler slot 0 is sampLinear.
	p.g.ctx.PSSetShaderResourcesAt(0, []*shaderResourceView{nil, nil})
	p.lastSRV, p.lastSampler = nil, p.sampLinear
}
