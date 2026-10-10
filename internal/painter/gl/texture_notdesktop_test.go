//go:build (android || ios || mobile || test_web_driver) && !wasm

package gl

import "fyne.io/fyne/v2/internal/cache"

// testTexture returns a non-zero texture handle for use in tests.
func testTexture() Texture {
	return Texture(cache.TextureType{Value: 1})
}
