//go:build wasm || test_web_driver

package cache

import "syscall/js"

func testTexture(value uint32) TextureType {
	return TextureType{Value: js.ValueOf(value)}
}
