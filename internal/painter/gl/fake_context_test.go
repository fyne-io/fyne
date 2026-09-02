//go:build !windows || !ci

package gl

// fakeContext records the GL calls the painter makes, so the drawing paths can
// be exercised without a real context.
type fakeContext struct {
	seq                       uint32
	draws, vertices           int
	bufferUploads, floats     int
	buffersFreed              int
	imageBytes, subImageBytes int
}

func (c *fakeContext) next() uint32 {
	c.seq++
	return c.seq
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

func (c *fakeContext) CreateBuffer() Buffer {
	return Buffer(c.next())
}

func (c *fakeContext) CreateProgram() Program {
	return Program(c.next())
}

func (c *fakeContext) CreateShader(uint32) Shader {
	return Shader(c.next())
}

func (c *fakeContext) CreateTexture() Texture {
	return Texture(c.next())
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

func (c *fakeContext) GetAttribLocation(Program, string) Attribute {
	return Attribute(c.next())
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

func (c *fakeContext) GetUniformLocation(Program, string) Uniform {
	return Uniform(c.next()) //gosec:disable G115 -- ids in a test never approach the limit
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
