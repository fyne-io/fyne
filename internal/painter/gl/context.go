package gl

type context interface {
	ActiveTexture(textureUnit uint32)
	BindBuffer(target uint32, buf Buffer)
	BindTexture(target uint32, texture Texture)
	BlendColor(r, g, b, a float32)
	BlendFunc(srcFactor, destFactor uint32)
	BufferData(target uint32, points []float32, usage uint32)
	Clear(mask uint32)
	ClearColor(r, g, b, a float32)
	CreateBuffer() Buffer
	CreateTexture() Texture
	DeleteBuffer(buffer Buffer)
	DeleteTexture(texture Texture)
	Disable(capability uint32)
	Enable(capability uint32)
	EnableVertexAttribArray(attribute Attribute)
	GetAttribLocation(program Program, name string) Attribute
	GetError() uint32
	Scissor(x, y, w, h int32)
	TexImage2D(target uint32, level, width, height int, colorFormat, typ uint32, data []uint8)
	TexParameteri(target, param uint32, value int32)
	UseProgram(program Program)
	VertexAttribPointerWithOffset(attribute Attribute, size int, typ uint32, normalized bool, stride, offset int)
	Viewport(x, y, width, height int)
}
