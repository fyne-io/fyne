//go:build !android && !ios && !mobile && !wasm && !test_web_driver

package gl

// testTexture returns a non-zero texture handle for use in tests.
func testTexture() Texture {
	return Texture(1)
}
