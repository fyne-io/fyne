//go:build android || ios || mobile

package cache

func testTexture(value uint32) TextureType {
	return TextureType{Value: value}
}
