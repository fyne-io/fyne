package gl

type context interface {
	ActiveTexture(textureUnit uint32)
	CreateTexture() Texture
	GetError() uint32
}
