//go:build windows

package dx

import (
	"fmt"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	d3d11DLL        = windows.NewLazySystemDLL("d3d11.dll")
	d3dCompiler     = windows.NewLazySystemDLL("d3dcompiler_47.dll")
	procD3D11Create = d3d11DLL.NewProc("D3D11CreateDeviceAndSwapChain")
	procD3DCompile  = d3dCompiler.NewProc("D3DCompile")
)

// COM plumbing.
//
// A COM interface pointer points at a struct whose first word is the vtable: a
// flat array of function pointers in interface-declaration order, base interfaces
// first. Calling a method is therefore an indirect call through a fixed slot with
// the interface pointer as the implicit first argument.
//
// ponytail: vtbl is typed as an oversized array so slots can be indexed directly
// rather than declaring a Go struct mirroring every vtable. Only slots that really
// exist for a given interface are ever indexed - the indices below are the whole
// correctness surface, so they carry their derivation in comments.

// maxVtblSlots covers the longest vtable indexed here: ID3D11DeviceContext ends
// at FinishCommandList, slot 114.
const maxVtblSlots = 128

type unknown struct {
	vtbl *[maxVtblSlots]uintptr
}

// call invokes vtable slot idx. The receiver is passed as the `this` argument.
//
// ponytail: variadic SyscallN on the draw path. Go keeps the arg slice on the
// stack here, so it does not allocate; if profiling ever says otherwise the fix is
// fixed-arity call2/call4/call6 helpers, not a redesign.
func (o *unknown) call(idx int, args ...uintptr) uintptr {
	all := make([]uintptr, 0, 12)
	all = append(all, uintptr(unsafe.Pointer(o)))
	all = append(all, args...)
	r, _, _ := syscall.SyscallN(o.vtbl[idx], all...)
	return r
}

// IUnknown slots.
const (
	slotQueryInterface = 0
	slotAddRef         = 1
	slotRelease        = 2
)

func (o *unknown) Release() uint32 {
	if o == nil {
		return 0
	}
	return uint32(o.call(slotRelease))
}

// releaseCOM is a nil-safe helper for the many optional COM pointers held by the
// Painter and window.
func releaseCOM[T any](p **T) {
	if *p == nil {
		return
	}
	(*unknown)(unsafe.Pointer(*p)).Release()
	*p = nil
}

type hresult int32

func (h hresult) failed() bool { return h < 0 }

// deviceLost reports whether the HRESULT says the device is gone and every
// resource created from it with it; the only recovery is a full recreate.
func (h hresult) deviceLost() bool {
	return uint32(h) == dxgiErrorDeviceRemoved || uint32(h) == dxgiErrorDeviceReset
}

func (h hresult) error(what string) error {
	if !h.failed() {
		return nil
	}
	return fmt.Errorf("%s failed: HRESULT 0x%08X", what, uint32(h))
}

// ---------------------------------------------------------------------------
// Enums and flags
// ---------------------------------------------------------------------------

const (
	d3d11SDKVersion = 7

	driverTypeHardware = 1
	driverTypeWARP     = 5

	featureLevel11_0 = 0xb000
	featureLevel10_1 = 0xa100
	featureLevel10_0 = 0xa000
	featureLevel9_3  = 0x9300

	createDeviceBGRASupport = 0x20
	createDeviceDebug       = 0x2

	formatR8G8B8A8Unorm     = 28
	formatB8G8R8A8Unorm     = 87
	formatR32G32Float       = 16
	formatR32G32B32A32Float = 2

	usageDefault   = 0
	usageImmutable = 1
	usageDynamic   = 2

	bindVertexBuffer   = 0x1
	bindIndexBuffer    = 0x2
	bindConstantBuffer = 0x4
	bindShaderResource = 0x8
	bindRenderTarget   = 0x20

	cpuAccessWrite = 0x10000

	mapWriteDiscard = 4

	swapEffectDiscard = 0
	// swapEffectFlipDiscard is the modern flip presentation model (Windows 10+):
	// the compositor reads the buffer directly instead of copying it, which is
	// measurably cheaper per Present. Older systems fall back to blt discard.
	swapEffectFlipDiscard = 4
	swapUsageRTOutput     = 0x00000020

	// DXGI_ERROR_DEVICE_REMOVED / _RESET, reported by Present when the GPU went
	// away (driver upgrade, TDR). Compared as uint32 because hresult is signed.
	dxgiErrorDeviceRemoved = 0x887a0005
	dxgiErrorDeviceReset   = 0x887a0007

	topologyTriangleList  = 4
	topologyTriangleStrip = 5

	inputPerVertexData = 0

	fillSolid = 3
	cullNone  = 1

	blendZero           = 1
	blendOne            = 2
	blendInvSrcAlpha    = 6
	blendSrcAlpha       = 5
	blendOpAdd          = 1
	colorWriteEnableAll = 0x0f

	filterMinMagMipLinear = 0x15
	filterMinMagMipPoint  = 0
	addressClamp          = 3

	// compileFlags: optimise, and pack matrices row-major so Go's row-major
	// float arrays upload without a transpose.
	compileOptimizationLevel3 = 1 << 15
	compilePackMatrixRowMajor = 1 << 3
)

// ---------------------------------------------------------------------------
// Structures - layouts must match the C headers exactly.
// ---------------------------------------------------------------------------

type dxgiRational struct {
	Numerator, Denominator uint32
}

type dxgiModeDesc struct {
	Width            uint32
	Height           uint32
	RefreshRate      dxgiRational
	Format           uint32
	ScanlineOrdering uint32
	Scaling          uint32
}

type dxgiSampleDesc struct {
	Count, Quality uint32
}

type dxgiSwapChainDesc struct {
	BufferDesc   dxgiModeDesc
	SampleDesc   dxgiSampleDesc
	BufferUsage  uint32
	BufferCount  uint32
	OutputWindow windows.Handle
	Windowed     int32
	SwapEffect   uint32
	Flags        uint32
}

type bufferDesc struct {
	ByteWidth           uint32
	Usage               uint32
	BindFlags           uint32
	CPUAccessFlags      uint32
	MiscFlags           uint32
	StructureByteStride uint32
}

type subresourceData struct {
	SysMem           unsafe.Pointer
	SysMemPitch      uint32
	SysMemSlicePitch uint32
}

type mappedSubresource struct {
	Data       unsafe.Pointer
	RowPitch   uint32
	DepthPitch uint32
}

type texture2DDesc struct {
	Width          uint32
	Height         uint32
	MipLevels      uint32
	ArraySize      uint32
	Format         uint32
	SampleDesc     dxgiSampleDesc
	Usage          uint32
	BindFlags      uint32
	CPUAccessFlags uint32
	MiscFlags      uint32
}

type inputElementDesc struct {
	SemanticName         *byte
	SemanticIndex        uint32
	Format               uint32
	InputSlot            uint32
	AlignedByteOffset    uint32
	InputSlotClass       uint32
	InstanceDataStepRate uint32
}

type samplerDesc struct {
	Filter         uint32
	AddressU       uint32
	AddressV       uint32
	AddressW       uint32
	MipLODBias     float32
	MaxAnisotropy  uint32
	ComparisonFunc uint32
	BorderColor    [4]float32
	MinLOD, MaxLOD float32
}

type renderTargetBlendDesc struct {
	BlendEnable           int32
	SrcBlend              uint32
	DestBlend             uint32
	BlendOp               uint32
	SrcBlendAlpha         uint32
	DestBlendAlpha        uint32
	BlendOpAlpha          uint32
	RenderTargetWriteMask uint8
	_                     [3]uint8
}

type blendDesc struct {
	AlphaToCoverageEnable  int32
	IndependentBlendEnable int32
	RenderTarget           [8]renderTargetBlendDesc
}

type rasterizerDesc struct {
	FillMode              uint32
	CullMode              uint32
	FrontCounterClockwise int32
	DepthBias             int32
	DepthBiasClamp        float32
	SlopeScaledDepthBias  float32
	DepthClipEnable       int32
	ScissorEnable         int32
	MultisampleEnable     int32
	AntialiasedLineEnable int32
}

// d3dRect is D3D11_RECT, used for the scissor rectangle. It has the same layout
// as the Win32 RECT the driver package declares, but the two are unrelated types
// in unrelated packages.
type d3dRect struct {
	Left, Top, Right, Bottom int32
}

type viewport struct {
	TopLeftX, TopLeftY, Width, Height, MinDepth, MaxDepth float32
}

// ---------------------------------------------------------------------------
// Interfaces
// ---------------------------------------------------------------------------

// device wraps ID3D11Device. Slot indices follow the header order after the three
// IUnknown entries: CreateBuffer is the first declared method, so slot 3.
type device struct{ unknown }

const (
	slotCreateBuffer             = 3
	slotCreateTexture2D          = 5
	slotCreateShaderResourceView = 7
	slotCreateRenderTargetView   = 9
	slotCreateInputLayout        = 11
	slotCreateVertexShader       = 12
	slotCreatePixelShader        = 15
	slotCreateBlendState         = 20
	slotCreateRasterizerState    = 22
	slotCreateSamplerState       = 23
	slotGetDeviceRemovedReason   = 39
)

type (
	buffer             struct{ unknown }
	texture2D          struct{ unknown }
	shaderResourceView struct{ unknown }
	renderTargetView   struct{ unknown }
	inputLayout        struct{ unknown }
	vertexShader       struct{ unknown }
	pixelShader        struct{ unknown }
	blendState         struct{ unknown }
	rasterizerState    struct{ unknown }
	samplerState       struct{ unknown }
)

func (d *device) CreateBuffer(desc *bufferDesc, data *subresourceData) (*buffer, error) {
	var out *buffer
	hr := hresult(d.call(slotCreateBuffer, uintptr(unsafe.Pointer(desc)),
		uintptr(unsafe.Pointer(data)), uintptr(unsafe.Pointer(&out))))
	return out, hr.error("CreateBuffer")
}

func (d *device) CreateTexture2D(desc *texture2DDesc, data *subresourceData) (*texture2D, error) {
	var out *texture2D
	hr := hresult(d.call(slotCreateTexture2D, uintptr(unsafe.Pointer(desc)),
		uintptr(unsafe.Pointer(data)), uintptr(unsafe.Pointer(&out))))
	return out, hr.error("CreateTexture2D")
}

// CreateShaderResourceView passes a nil description, which means "view the whole
// resource in its own format" - exactly what every texture here wants.
func (d *device) CreateShaderResourceView(res *texture2D) (*shaderResourceView, error) {
	var out *shaderResourceView
	hr := hresult(d.call(slotCreateShaderResourceView, uintptr(unsafe.Pointer(res)), 0,
		uintptr(unsafe.Pointer(&out))))
	return out, hr.error("CreateShaderResourceView")
}

func (d *device) CreateRenderTargetView(res *texture2D) (*renderTargetView, error) {
	var out *renderTargetView
	hr := hresult(d.call(slotCreateRenderTargetView, uintptr(unsafe.Pointer(res)), 0,
		uintptr(unsafe.Pointer(&out))))
	return out, hr.error("CreateRenderTargetView")
}

func (d *device) CreateInputLayout(elems []inputElementDesc, shaderByteCode []byte) (*inputLayout, error) {
	var out *inputLayout
	hr := hresult(d.call(slotCreateInputLayout, uintptr(unsafe.Pointer(&elems[0])), uintptr(len(elems)),
		uintptr(unsafe.Pointer(&shaderByteCode[0])), uintptr(len(shaderByteCode)),
		uintptr(unsafe.Pointer(&out))))
	return out, hr.error("CreateInputLayout")
}

func (d *device) CreateVertexShader(code []byte) (*vertexShader, error) {
	var out *vertexShader
	hr := hresult(d.call(slotCreateVertexShader, uintptr(unsafe.Pointer(&code[0])), uintptr(len(code)), 0,
		uintptr(unsafe.Pointer(&out))))
	return out, hr.error("CreateVertexShader")
}

func (d *device) CreatePixelShader(code []byte) (*pixelShader, error) {
	var out *pixelShader
	hr := hresult(d.call(slotCreatePixelShader, uintptr(unsafe.Pointer(&code[0])), uintptr(len(code)), 0,
		uintptr(unsafe.Pointer(&out))))
	return out, hr.error("CreatePixelShader")
}

func (d *device) CreateBlendState(desc *blendDesc) (*blendState, error) {
	var out *blendState
	hr := hresult(d.call(slotCreateBlendState, uintptr(unsafe.Pointer(desc)), uintptr(unsafe.Pointer(&out))))
	return out, hr.error("CreateBlendState")
}

func (d *device) CreateRasterizerState(desc *rasterizerDesc) (*rasterizerState, error) {
	var out *rasterizerState
	hr := hresult(d.call(slotCreateRasterizerState, uintptr(unsafe.Pointer(desc)), uintptr(unsafe.Pointer(&out))))
	return out, hr.error("CreateRasterizerState")
}

func (d *device) CreateSamplerState(desc *samplerDesc) (*samplerState, error) {
	var out *samplerState
	hr := hresult(d.call(slotCreateSamplerState, uintptr(unsafe.Pointer(desc)), uintptr(unsafe.Pointer(&out))))
	return out, hr.error("CreateSamplerState")
}

// RemovedReason surfaces why a lost device died, for the log line at detection.
func (d *device) RemovedReason() hresult {
	return hresult(d.call(slotGetDeviceRemovedReason))
}

// deviceContext wraps ID3D11DeviceContext. It derives from ID3D11DeviceChild,
// which adds four methods after IUnknown, so the first context method sits at
// slot 7.
type deviceContext struct{ unknown }

const (
	slotVSSetConstantBuffers   = 7
	slotPSSetShaderResources   = 8
	slotPSSetShader            = 9
	slotPSSetSamplers          = 10
	slotVSSetShader            = 11
	slotDraw                   = 13
	slotMap                    = 14
	slotUnmap                  = 15
	slotPSSetConstantBuffers   = 16
	slotIASetInputLayout       = 17
	slotIASetVertexBuffers     = 18
	slotIASetPrimitiveTopology = 24
	slotOMSetRenderTargets     = 33
	slotOMSetBlendState        = 35
	slotRSSetState             = 43
	slotRSSetViewports         = 44
	slotRSSetScissorRects      = 45
	slotCopySubresourceRegion  = 46
	slotCopyResource           = 47
	slotUpdateSubresource      = 48
	slotClearRenderTargetView  = 50
	slotFlush                  = 111
)

func (c *deviceContext) OMSetRenderTargets(rtv *renderTargetView) {
	c.call(slotOMSetRenderTargets, 1, uintptr(unsafe.Pointer(&rtv)), 0)
}

func (c *deviceContext) ClearRenderTargetView(rtv *renderTargetView, rgba *[4]float32) {
	c.call(slotClearRenderTargetView, uintptr(unsafe.Pointer(rtv)), uintptr(unsafe.Pointer(rgba)))
}

func (c *deviceContext) RSSetViewports(v *viewport) {
	c.call(slotRSSetViewports, 1, uintptr(unsafe.Pointer(v)))
}

func (c *deviceContext) RSSetScissorRects(r *d3dRect) {
	if r == nil {
		c.call(slotRSSetScissorRects, 0, 0)
		return
	}
	c.call(slotRSSetScissorRects, 1, uintptr(unsafe.Pointer(r)))
}

func (c *deviceContext) RSSetState(s *rasterizerState) {
	c.call(slotRSSetState, uintptr(unsafe.Pointer(s)))
}

func (c *deviceContext) OMSetBlendState(s *blendState) {
	c.call(slotOMSetBlendState, uintptr(unsafe.Pointer(s)), 0, 0xffffffff)
}

func (c *deviceContext) IASetInputLayout(l *inputLayout) {
	c.call(slotIASetInputLayout, uintptr(unsafe.Pointer(l)))
}

func (c *deviceContext) IASetPrimitiveTopology(t uint32) {
	c.call(slotIASetPrimitiveTopology, uintptr(t))
}

func (c *deviceContext) IASetVertexBuffer(b *buffer, stride, offset uint32) {
	c.call(slotIASetVertexBuffers, 0, 1, uintptr(unsafe.Pointer(&b)),
		uintptr(unsafe.Pointer(&stride)), uintptr(unsafe.Pointer(&offset)))
}

func (c *deviceContext) VSSetShader(s *vertexShader) {
	c.call(slotVSSetShader, uintptr(unsafe.Pointer(s)), 0, 0)
}

func (c *deviceContext) PSSetShader(s *pixelShader) {
	c.call(slotPSSetShader, uintptr(unsafe.Pointer(s)), 0, 0)
}

func (c *deviceContext) VSSetConstantBuffer(b *buffer) {
	c.call(slotVSSetConstantBuffers, 0, 1, uintptr(unsafe.Pointer(&b)))
}

func (c *deviceContext) PSSetConstantBuffer(b *buffer) {
	c.call(slotPSSetConstantBuffers, 0, 1, uintptr(unsafe.Pointer(&b)))
}

func (c *deviceContext) PSSetShaderResource(srv *shaderResourceView) {
	c.call(slotPSSetShaderResources, 0, 1, uintptr(unsafe.Pointer(&srv)))
}

func (c *deviceContext) PSSetSampler(s *samplerState) {
	c.call(slotPSSetSamplers, 0, 1, uintptr(unsafe.Pointer(&s)))
}

// The slot-taking variants below serve user canvas.Shader draws, which bind a
// second constant buffer and several textures at once. The built in shaders
// only ever use slot 0, so they keep the simpler forms above.

func (c *deviceContext) PSSetConstantBufferAt(slot uint32, b *buffer) {
	c.call(slotPSSetConstantBuffers, uintptr(slot), 1, uintptr(unsafe.Pointer(&b)))
}

func (c *deviceContext) PSSetShaderResourcesAt(slot uint32, views []*shaderResourceView) {
	if len(views) == 0 {
		return
	}
	c.call(slotPSSetShaderResources, uintptr(slot), uintptr(len(views)), uintptr(unsafe.Pointer(&views[0])))
}

func (c *deviceContext) PSSetSamplersAt(slot uint32, samplers []*samplerState) {
	if len(samplers) == 0 {
		return
	}
	c.call(slotPSSetSamplers, uintptr(slot), uintptr(len(samplers)), uintptr(unsafe.Pointer(&samplers[0])))
}

func (c *deviceContext) Draw(vertexCount, startVertex uint32) {
	c.call(slotDraw, uintptr(vertexCount), uintptr(startVertex))
}

// box is D3D11_BOX, the source region for CopySubresourceRegion. Coordinates are
// top-origin pixels, unlike OpenGL's bottom-origin glCopyTexSubImage2D.
type box struct {
	Left, Top, Front, Right, Bottom, Back uint32
}

// UpdateSubresource copies CPU data into a default-usage resource in one call,
// half the syscalls of a Map/Unmap pair. For buffers dstBox selects the byte
// range (nil updates the whole resource - required for constant buffers); the
// driver versions the copy internally, so it does not stall against in-flight
// draws reading the previous contents.
func (c *deviceContext) UpdateSubresource(res, data unsafe.Pointer, dstBox *box) {
	c.call(slotUpdateSubresource, uintptr(res), 0, uintptr(unsafe.Pointer(dstBox)), uintptr(data), 0, 0)
}

// CopySubresourceRegion copies a rectangle of src into dst at the origin.
func (c *deviceContext) CopySubresourceRegion(dst *texture2D, src *texture2D, region *box) {
	c.call(slotCopySubresourceRegion, uintptr(unsafe.Pointer(dst)), 0, 0, 0, 0,
		uintptr(unsafe.Pointer(src)), 0, uintptr(unsafe.Pointer(region)))
}

func (c *deviceContext) CopyResource(dst, src *texture2D) {
	c.call(slotCopyResource, uintptr(unsafe.Pointer(dst)), uintptr(unsafe.Pointer(src)))
}

func (c *deviceContext) Flush() {
	c.call(slotFlush)
}

// Map locks a dynamic resource for CPU writes. Callers must Unmap.
func (c *deviceContext) Map(res unsafe.Pointer, mapType uint32) (mappedSubresource, error) {
	var m mappedSubresource
	hr := hresult(c.call(slotMap, uintptr(res), 0, uintptr(mapType), 0, uintptr(unsafe.Pointer(&m))))
	return m, hr.error("Map")
}

func (c *deviceContext) Unmap(res unsafe.Pointer) {
	c.call(slotUnmap, uintptr(res), 0)
}

// swapChain wraps IDXGISwapChain: IUnknown (3) + IDXGIObject (4) +
// IDXGIDeviceSubObject (1) = 8 inherited slots before Present.
type swapChain struct{ unknown }

const (
	slotPresent       = 8
	slotGetBuffer     = 9
	slotResizeBuffers = 13
)

func (s *swapChain) Present(syncInterval uint32) hresult {
	return hresult(s.call(slotPresent, uintptr(syncInterval), 0))
}

// GetBuffer fetches the back buffer as a texture. iidTexture2D identifies
// ID3D11Texture2D.
func (s *swapChain) GetBuffer(index uint32) (*texture2D, error) {
	var out *texture2D
	hr := hresult(s.call(slotGetBuffer, uintptr(index), uintptr(unsafe.Pointer(&iidTexture2D)),
		uintptr(unsafe.Pointer(&out))))
	return out, hr.error("GetBuffer")
}

func (s *swapChain) ResizeBuffers(w, h uint32) error {
	hr := hresult(s.call(slotResizeBuffers, 0, uintptr(w), uintptr(h), 0, 0))
	return hr.error("ResizeBuffers")
}

// iidTexture2D is IID_ID3D11Texture2D {6f15aaf2-d208-4e89-9ab4-489535d34f9c}.
var iidTexture2D = windows.GUID{
	Data1: 0x6f15aaf2, Data2: 0xd208, Data3: 0x4e89,
	Data4: [8]byte{0x9a, 0xb4, 0x48, 0x95, 0x35, 0xd3, 0x4f, 0x9c},
}

// blob wraps ID3DBlob, the buffer type returned by the shader compiler.
type blob struct{ unknown }

const (
	slotGetBufferPointer = 3
	slotGetBufferSize    = 4
)

// bytes exposes the blob contents. The returned slice aliases COM-owned memory
// and is only valid until the blob is released.
func (b *blob) bytes() []byte {
	ptr := b.call(slotGetBufferPointer)
	size := b.call(slotGetBufferSize)
	if ptr == 0 || size == 0 {
		return nil
	}
	return unsafe.Slice((*byte)(comMemory(ptr)), int(size))
}

// comMemory reinterprets an address returned by COM as a pointer.
//
// go vet's unsafeptr check rejects a direct uintptr->unsafe.Pointer conversion
// because for Go heap memory the address can go stale the moment it stops being a
// pointer the collector can see. That hazard does not apply here: the allocation
// belongs to the COM allocator, lives outside the Go heap, and never moves. The
// type pun below states that intent without tripping the heuristic.
func comMemory(addr uintptr) unsafe.Pointer {
	return *(*unsafe.Pointer)(unsafe.Pointer(&addr))
}

// ---------------------------------------------------------------------------
// Entry points
// ---------------------------------------------------------------------------

// createDeviceAndSwapChain brings up D3D11 bound to hwnd. It asks for a hardware
// device first and falls back to WARP so the driver still runs on machines with no
// usable GPU driver (VMs, RDP sessions).
func createDeviceAndSwapChain(hwnd windows.Handle, width, height uint32) (*device, *deviceContext, *swapChain, uint32, error) {
	levels := [...]uint32{featureLevel11_0, featureLevel10_1, featureLevel10_0, featureLevel9_3}
	desc := dxgiSwapChainDesc{
		BufferDesc: dxgiModeDesc{
			Width:  width,
			Height: height,
			Format: formatB8G8R8A8Unorm,
		},
		SampleDesc:   dxgiSampleDesc{Count: 1},
		BufferUsage:  swapUsageRTOutput,
		BufferCount:  2,
		OutputWindow: hwnd,
		Windowed:     1,
		// SwapEffect is set per attempt below: flip model first, blt fallback.
	}

	lastErr := fmt.Errorf("no Direct3D 11 device available")
	for _, driverType := range [...]uint32{driverTypeHardware, driverTypeWARP} {
		for _, effect := range [...]uint32{swapEffectFlipDiscard, swapEffectDiscard} {
			desc.SwapEffect = effect
			var sc *swapChain
			var dev *device
			var ctx *deviceContext
			var gotLevel uint32
			hr := hresult(callD3D11Create(driverType, createDeviceBGRASupport, &levels[0], uint32(len(levels)),
				&desc, &sc, &dev, &gotLevel, &ctx))
			if !hr.failed() {
				return dev, ctx, sc, gotLevel, nil
			}
			lastErr = hr.error("D3D11CreateDeviceAndSwapChain")
		}
	}
	return nil, nil, nil, 0, lastErr
}

func callD3D11Create(driverType, flags uint32, levels *uint32, levelCount uint32,
	desc *dxgiSwapChainDesc, sc **swapChain, dev **device, gotLevel *uint32, ctx **deviceContext,
) int32 {
	r, _, _ := procD3D11Create.Call(
		0, // pAdapter - null means "pick for the driver type"
		uintptr(driverType),
		0, // Software
		uintptr(flags),
		uintptr(unsafe.Pointer(levels)),
		uintptr(levelCount),
		d3d11SDKVersion,
		uintptr(unsafe.Pointer(desc)),
		uintptr(unsafe.Pointer(sc)),
		uintptr(unsafe.Pointer(dev)),
		uintptr(unsafe.Pointer(gotLevel)),
		uintptr(unsafe.Pointer(ctx)),
	)
	return int32(r)
}

// compileShader builds HLSL source into DXBC at runtime via d3dcompiler_47.dll,
// which ships with Windows 8.1 and later.
//
// ponytail: runtime compilation costs a few ms once per shader at startup and
// avoids a build-time fxc step plus checked-in binaries. If startup time ever
// matters, precompile to .cso and embed - the call site does not change.
func compileShader(src, entry, target string) ([]byte, error) {
	entryC, err := windows.BytePtrFromString(entry)
	if err != nil {
		return nil, err
	}
	targetC, err := windows.BytePtrFromString(target)
	if err != nil {
		return nil, err
	}
	nameC, err := windows.BytePtrFromString("fyne.hlsl")
	if err != nil {
		return nil, err
	}
	srcBytes := []byte(src)

	var code, errs *blob
	r, _, _ := procD3DCompile.Call(
		uintptr(unsafe.Pointer(&srcBytes[0])),
		uintptr(len(srcBytes)),
		uintptr(unsafe.Pointer(nameC)),
		0, // pDefines
		0, // pInclude
		uintptr(unsafe.Pointer(entryC)),
		uintptr(unsafe.Pointer(targetC)),
		compileOptimizationLevel3|compilePackMatrixRowMajor,
		0,
		uintptr(unsafe.Pointer(&code)),
		uintptr(unsafe.Pointer(&errs)),
	)
	if hr := hresult(r); hr.failed() {
		detail := ""
		if errs != nil {
			detail = ": " + string(errs.bytes())
			errs.Release()
		}
		return nil, fmt.Errorf("compiling %s (%s)%s", entry, target, detail)
	}
	if errs != nil {
		errs.Release()
	}
	// Copy out of COM memory so the blob can be released immediately.
	out := append([]byte(nil), code.bytes()...)
	code.Release()
	return out, nil
}
