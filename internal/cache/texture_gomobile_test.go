//go:build android || ios || mobile
// +build android ios mobile

package cache_test

import (
	"sync"
	"testing"

	"fyne.io/fyne/v2/internal/cache"
	"github.com/fyne-io/mobile/gl"
)

// go test -race
func TestTextureGomobile(t *testing.T) {
	wg := sync.WaitGroup{}

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
