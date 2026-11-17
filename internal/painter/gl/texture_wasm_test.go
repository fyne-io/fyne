//go:build wasm

package gl

import (
	"syscall/js"

	"fyne.io/fyne/v2/internal/cache"
)

// testTexture returns a non-zero texture handle for use in tests.
func testTexture() Texture {
	return Texture(cache.TextureType{Value: js.ValueOf(1)})
}
