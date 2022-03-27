package gl

type context interface {
	CreateTexture() Texture
	GetError() uint32
}
