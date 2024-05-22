//go:build software

package cache_test

import (
	"image"
	"sync"
	"testing"

	"fyne.io/fyne/v2/internal/cache"
)

// go test -race
func TestTextureSoftware(t *testing.T) {
	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		cache.SetTexture(nil, image.White, nil)
	}()
	go func() {
		defer wg.Done()
		cache.SetTexture(nil, image.Black, nil)
	}()

	wg.Wait()
}
