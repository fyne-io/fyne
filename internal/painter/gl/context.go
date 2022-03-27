package gl

type context interface {
	ActiveTexture(textureUnit uint32)
	BindTexture(target uint32, texture Texture)
	CreateTexture() Texture
	GetError() uint32
	TexParameteri(target, param uint32, value int32)
}
