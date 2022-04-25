package gl

type context interface {
	ActiveTexture(textureUnit uint32)
	BindTexture(target uint32, texture Texture)
	Clear(mask uint32)
	ClearColor(r, g, b, a float32)
	CreateTexture() Texture
	DeleteTexture(texture Texture)
	GetError() uint32
	Scissor(x, y, w, h int32)
	TexImage2D(target uint32, level, width, height int, colorFormat, typ uint32, data []uint8)
	TexParameteri(target, param uint32, value int32)
	UseProgram(program Program)
	Viewport(x, y, width, height int)
}
