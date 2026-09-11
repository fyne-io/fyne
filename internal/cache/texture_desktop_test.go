//go:build !android && !ios && !mobile && !wasm && !test_web_driver

package cache

func testTexture(value uint32) TextureType {
	return TextureType(value)
}
