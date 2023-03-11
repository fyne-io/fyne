//go:build !wasm
// +build !wasm

package gl

func (v Attrib) IsValid() bool       { return v != NoAttrib }
func (v Program) IsValid() bool      { return v != NoProgram }
func (v Shader) IsValid() bool       { return v != NoShader }
func (v Buffer) IsValid() bool       { return v != NoBuffer }
func (v Framebuffer) IsValid() bool  { return v != NoFramebuffer }
func (v Renderbuffer) IsValid() bool { return v != NoRenderbuffer }
func (v Texture) IsValid() bool      { return v != NoTexture }
func (v Uniform) IsValid() bool      { return v != NoUniform }
