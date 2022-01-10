//go:build !wasm
// +build !wasm

package gl

func (v Attrib) Valid() bool       { return v != NoAttrib }
func (v Program) Valid() bool      { return v != NoProgram }
func (v Shader) Valid() bool       { return v != NoShader }
func (v Buffer) Valid() bool       { return v != NoBuffer }
func (v Framebuffer) Valid() bool  { return v != NoFramebuffer }
func (v Renderbuffer) Valid() bool { return v != NoRenderbuffer }
func (v Texture) Valid() bool      { return v != NoTexture }
func (v Uniform) Valid() bool      { return v != NoUniform }
