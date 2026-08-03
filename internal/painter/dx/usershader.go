//go:build windows

package dx

// User canvas.Shader support, mirroring drawShader in internal/Painter/gl/draw.go:
// the program is compiled once per Shader.Name and cached for the lifetime of
// the device, textures are uploaded once and reused until the image is replaced,
// and the uniforms are pushed every frame.
//
// The contract a user shader has to meet is documented on canvas.Shader.SourceHLSL:
//
//   - the shared prelude (internal/Painter/dx, shaders/common.hlsl) is prepended,
//     so `frame`, `bounds` and PSIn arrive exactly as
//     the built in shaders see them. `bounds` is already top-origin and so is
//     SV_Position, so a GLSL original's `frame.y - gl_FragCoord.y` flip is simply
//     dropped when porting - the same deviation shaders.go documents.
//   - the shader declares its own Uniforms as a flat `cbuffer X : register(b1)`.
//     Offsets are read out of the source text rather than by reflection, so the
//     names in Shader.Uniforms can be matched to buffer offsets without pulling
//     in ID3D11ShaderReflection.
//   - Textures bind to t0 upwards with a sampler each at s0 upwards, ordered by
//     sorted name. That is the same assignment the GL Painter makes to texture
//     units, so one shader can serve both backends.

import (
	"encoding/binary"
	"image"
	"math"
	"regexp"
	"sort"
	"strings"
	"unsafe"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
)

// userShader is the cached device state for one canvas.Shader.
type userShader struct {
	ps   *pixelShader
	cbuf *buffer // the b1 uniform buffer, nil when the shader declares none

	offsets map[string]uint32 // uniform name -> byte offset into cbuf
	data    []byte            // staging copy of cbuf, refilled each draw

	textures map[string]*shaderTexture
	srvs     []*shaderResourceView // rebuilt per draw, kept to avoid the garbage
	samplers []*samplerState

	valid bool
}

// shaderTexture pairs an upload with the image it came from, so a Textures entry
// replaced with a different image re-uploads and an unchanged one does not.
type shaderTexture struct {
	GPU *gpuTexture
	src image.Image
}

func (s *userShader) release() {
	for name, t := range s.textures {
		t.GPU.release()
		delete(s.textures, name)
	}
	releaseCOM(&s.cbuf)
	releaseCOM(&s.ps)
}

func (p *Painter) drawShader(shader *canvas.Shader, pos fyne.Position, frame fyne.Size) {
	state, ok := p.userShader(shader)
	if !ok {
		return
	}
	if !state.bindTextures(p, shader) {
		return
	}
	state.uploadUniforms(p, shader)

	points, bounds := p.vecRectCoords(pos, shader, frame, 0, canvas.Shadow{})

	c := constants{}
	fw, fh := p.scaleFrameSize(frame)
	c.Frame = [4]float32{fw, fh, 0, 0}
	x1, x2, y1, y2 := p.scaleRectCoords(bounds[0], bounds[2], bounds[1], bounds[3])
	c.Bounds = [4]float32{x1, y1, x2, y2}

	p.upload(p.quadVertices(points), &c, p.vsQuad, state.ps, topologyTriangleStrip, p.blend)
}

// userShader returns the cached state for a shader, compiling it on first use.
// A compile failure is cached as an invalid entry so it is logged once rather
// than every frame, like the GL Painter's shaderPrograms cache.
func (p *Painter) userShader(shader *canvas.Shader) (*userShader, bool) {
	if p.userShaders == nil {
		p.userShaders = make(map[string]*userShader)
	}
	if state, ok := p.userShaders[shader.Name]; ok {
		return state, state.valid
	}

	state := &userShader{textures: make(map[string]*shaderTexture, len(shader.Textures))}
	p.userShaders[shader.Name] = state

	if len(shader.SourceHLSL) == 0 {
		fyne.LogError("directx: shader "+shader.Name+" has no HLSL source, not drawn", nil)
		return state, false
	}

	code, err := compileShader(Common+string(shader.SourceHLSL), "main", "ps_4_0")
	if err != nil {
		fyne.LogError("directx: compile shader "+shader.Name, err)
		return state, false
	}
	if state.ps, err = p.g.dev.CreatePixelShader(code); err != nil {
		fyne.LogError("directx: create shader "+shader.Name, err)
		return state, false
	}

	if offsets, size := parseUserUniforms(string(shader.SourceHLSL)); size > 0 {
		state.cbuf, err = p.g.dev.CreateBuffer(&bufferDesc{
			ByteWidth: size,
			Usage:     usageDefault,
			BindFlags: bindConstantBuffer,
		}, nil)
		if err != nil {
			fyne.LogError("directx: uniform buffer for shader "+shader.Name, err)
			releaseCOM(&state.ps)
			return state, false
		}
		state.offsets = offsets
		state.data = make([]byte, size)
	} else if len(shader.Uniforms) > 0 {
		fyne.LogError("directx: shader "+shader.Name+" sets uniforms but its register(b1) cbuffer "+
			"is missing or not parseable (flat float/floatN declarations only), they are ignored", nil)
	}

	state.valid = true
	return state, true
}

// bindTextures uploads any new or replaced image and binds every texture to its
// slot. It reports false if an upload failed, in which case the draw is skipped
// rather than sampling a stale or missing texture.
func (s *userShader) bindTextures(p *Painter, shader *canvas.Shader) bool {
	if len(shader.Textures) == 0 {
		return true
	}

	names := make([]string, 0, len(shader.Textures))
	for name, img := range shader.Textures {
		names = append(names, name)
		if cached := s.textures[name]; cached == nil || cached.src != img {
			if cached != nil {
				p.releaseTexture(cached.GPU)
			}
			GPU := p.imgToTexture(img)
			if GPU == nil {
				delete(s.textures, name)
				return false
			}
			s.textures[name] = &shaderTexture{GPU: GPU, src: img}
		}
	}
	sort.Strings(names)

	s.srvs, s.samplers = s.srvs[:0], s.samplers[:0]
	for _, name := range names {
		s.srvs = append(s.srvs, s.textures[name].GPU.srv)
		// Smooth sampling for every texture, matching the GL Painter, which
		// uploads user textures with canvas.ImageScaleSmooth.
		s.samplers = append(s.samplers, p.sampLinear)
	}
	p.g.ctx.PSSetShaderResourcesAt(0, s.srvs)
	p.g.ctx.PSSetSamplersAt(0, s.samplers)
	// Slot 0 changed behind the hot path's back; keep its cache truthful.
	if len(s.srvs) > 0 {
		p.lastSRV, p.lastSampler = s.srvs[0], s.samplers[0]
	}
	return true
}

// uploadUniforms writes the shader's float uniforms into the b1 buffer and binds
// it. Names with no matching declaration are skipped; unwritten slots keep the
// previous frame's value, which matches GL leaving an unset uniform alone.
func (s *userShader) uploadUniforms(p *Painter, shader *canvas.Shader) {
	if s.cbuf == nil {
		return
	}
	for name, v := range shader.Uniforms {
		if off, ok := s.offsets[name]; ok {
			binary.LittleEndian.PutUint32(s.data[off:], math.Float32bits(v))
		}
	}

	p.g.ctx.UpdateSubresource(unsafe.Pointer(s.cbuf), unsafe.Pointer(&s.data[0]), nil)
	p.g.ctx.PSSetConstantBufferAt(1, s.cbuf)
}

func (p *Painter) releaseUserShaders() {
	for name, state := range p.userShaders {
		state.release()
		delete(p.userShaders, name)
	}
}

var (
	userCBufferRE  = regexp.MustCompile(`(?s)cbuffer\s+\w+\s*:\s*register\s*\(\s*b1\s*\)\s*\{(.*?)\}`)
	lineCommentRE  = regexp.MustCompile(`//[^\n]*`)
	blockCommentRE = regexp.MustCompile(`(?s)/\*.*?\*/`)
)

// parseUserUniforms maps each name in the shader's register(b1) cbuffer to its
// byte offset, and returns the buffer size rounded up to a 16 byte register.
//
// ponytail: a flat cbuffer of scalars and vectors is all any shader here needs,
// so the packing rules implemented are the two that matter - elements pack
// tightly, and a vector that would straddle a 16 byte boundary is bumped past
// it. Arrays, matrices and structs each occupy whole registers per element and
// are not handled: an unknown type abandons the parse rather than silently
// returning offsets that are wrong from that point on. Swap in
// ID3D11ShaderReflection if a shader ever needs them.
func parseUserUniforms(src string) (map[string]uint32, uint32) {
	m := userCBufferRE.FindStringSubmatch(stripHLSLComments(src))
	if m == nil {
		return nil, 0
	}

	offsets := make(map[string]uint32)
	var off uint32
	for _, decl := range strings.Split(m[1], ";") {
		decl = strings.TrimSpace(decl)
		if decl == "" {
			continue
		}
		typ, rest, ok := strings.Cut(decl, " ")
		if !ok {
			return nil, 0
		}
		size := hlslScalarSize(typ)
		if size == 0 {
			return nil, 0
		}
		for _, name := range strings.Split(rest, ",") {
			name = strings.TrimSpace(name)
			if name == "" || !isHLSLIdent(name) {
				return nil, 0
			}
			// a vector may not straddle a 16 byte register boundary
			if size > 4 && off/16 != (off+size-1)/16 {
				off = (off/16 + 1) * 16
			}
			offsets[name] = off
			off += size
		}
	}
	if off == 0 {
		return nil, 0
	}
	return offsets, (off + 15) / 16 * 16
}

func hlslScalarSize(typ string) uint32 {
	switch typ {
	case "float", "int", "uint", "bool":
		return 4
	case "float2", "int2", "uint2":
		return 8
	case "float3", "int3", "uint3":
		return 12
	case "float4", "int4", "uint4":
		return 16
	}
	return 0
}

func isHLSLIdent(s string) bool {
	for i, r := range s {
		switch {
		case r >= 'a' && r <= 'z', r >= 'A' && r <= 'Z', r == '_':
		case i > 0 && r >= '0' && r <= '9':
		default:
			return false
		}
	}
	return true
}

func stripHLSLComments(src string) string {
	return lineCommentRE.ReplaceAllString(blockCommentRE.ReplaceAllString(src, " "), "")
}
