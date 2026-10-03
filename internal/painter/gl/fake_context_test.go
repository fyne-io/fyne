//go:build !windows || !ci

package gl

// fakeContext records the GL calls the painter makes, so the drawing paths can
// be exercised without a real context. The handle types are numeric on desktop
// and structs on mobile, so the fake hands back zero values: nothing here reads
// them back, and the zero value is the one form that compiles everywhere.
type fakeContext struct {
	draws, vertices           int
	bufferUploads, floats     int
	buffersFreed              int
	imageBytes, subImageBytes int
}

var _ context = (*fakeContext)(nil)

func (*fakeContext) ActiveTexture(uint32) {}

func (*fakeContext) AttachShader(Program, Shader) {}

func (*fakeContext) BindBuffer(uint32, Buffer) {}

func (*fakeContext) BindTexture(uint32, Texture) {}

func (*fakeContext) BlendColor(float32, float32, float32, float32) {}

func (*fakeContext) BlendFunc(uint32, uint32) {}

func (c *fakeContext) BufferData(_ uint32, points []float32, _ uint32) {
	c.bufferUploads++
	c.floats += len(points)
}

func (*fakeContext) BufferSubData(uint32, []float32) {}

func (*fakeContext) Clear(uint32) {}

func (*fakeContext) ClearColor(float32, float32, float32, float32) {}

func (*fakeContext) CompileShader(Shader) {}

func (*fakeContext) CreateBuffer() Buffer {
	var v Buffer
	return v
}

func (*fakeContext) CreateProgram() Program {
	var v Program
	return v
}

func (*fakeContext) CreateShader(uint32) Shader {
	var v Shader
	return v
}

func (*fakeContext) CreateTexture() Texture {
	var v Texture
	return v
}

func (c *fakeContext) DeleteBuffer(Buffer) {
	c.buffersFreed++
}

func (*fakeContext) DeleteProgram(Program) {}

func (*fakeContext) DeleteTexture(Texture) {}

func (*fakeContext) Disable(uint32) {}

func (c *fakeContext) DrawArrays(_ uint32, _, count int) {
	c.draws++
	c.vertices += count
}

func (*fakeContext) Enable(uint32) {}

func (*fakeContext) EnableVertexAttribArray(Attribute) {}

func (*fakeContext) GetAttribLocation(Program, string) Attribute {
	var v Attribute
	return v
}

func (*fakeContext) GetError() uint32 {
	return 0
}

func (*fakeContext) GetInteger(uint32) int {
	return 2048
}

func (*fakeContext) GetProgrami(Program, uint32) int {
	return 1
}

func (*fakeContext) GetProgramInfoLog(Program) string {
	return ""
}

func (*fakeContext) GetShaderi(Shader, uint32) int {
	return 1
}

func (*fakeContext) GetShaderInfoLog(Shader) string {
	return ""
}

func (*fakeContext) GetUniformLocation(Program, string) Uniform {
	var v Uniform
	return v
}

func (*fakeContext) LinkProgram(Program) {}

func (*fakeContext) CopyTexSubImage2D(uint32, int, int, int, int, int, int, int) {}

func (*fakeContext) ReadBuffer(uint32) {}

func (*fakeContext) ReadPixels(int, int, int, int, uint32, uint32, []uint8) {}

func (*fakeContext) Scissor(int32, int32, int32, int32) {}

func (*fakeContext) ShaderSource(Shader, string) {}

func (c *fakeContext) TexImage2D(_ uint32, _, _, _ int, _, _ uint32, data []uint8) {
	c.imageBytes += len(data)
}

func (c *fakeContext) TexSubImage2D(_ uint32, _, _, _, _, _ int, _, _ uint32, data []uint8) {
	c.subImageBytes += len(data)
}

func (*fakeContext) TexParameteri(uint32, uint32, int32) {}

func (*fakeContext) Uniform1f(Uniform, float32) {}

func (*fakeContext) Uniform1fv(Uniform, []float32) {}

func (*fakeContext) Uniform1i(Uniform, int32) {}

func (*fakeContext) Uniform2f(Uniform, float32, float32) {}

func (*fakeContext) Uniform2fv(Uniform, []float32) {}

func (*fakeContext) Uniform4f(Uniform, float32, float32, float32, float32) {}

func (*fakeContext) UseProgram(Program) {}

func (*fakeContext) VertexAttribPointerWithOffset(Attribute, int, uint32, bool, int, int) {}

func (*fakeContext) Viewport(int, int, int, int) {}
