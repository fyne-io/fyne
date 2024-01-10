//go:build android || ios || mobile

package cache_test

import (
	"sync"
	"testing"

	"fyne.io/fyne/v2/internal/cache"
	"fyne.io/fyne/v2/internal/driver/mobile/gl"
)

// go test -race
func TestTextureGomobile(t *testing.T) {
	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		cache.SetTexture(nil, gl.Texture{0}, nil)
	}()
	go func() {
		defer wg.Done()
		cache.SetTexture(nil, gl.Texture{1}, nil)
	}()

	wg.Wait()
}
