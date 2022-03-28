package gl

type context interface {
	ActiveTexture(textureUnit uint32)
	BindTexture(target uint32, texture Texture)
	CreateTexture() Texture
	GetError() uint32
	TexImage2D(target uint32, level, width, height int, colorFormat, typ uint32, data []uint8)
	TexParameteri(target, param uint32, value int32)
	Viewport(x, y, width, height int)
}
