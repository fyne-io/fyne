//go:build !android && !ios && !mobile

package cache_test

import (
	"sync"
	"testing"

	"fyne.io/fyne/v2/internal/cache"
)

// go test -race
func TestTextureDesktop(t *testing.T) {
	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		cache.SetTexture(nil, 0, nil)
	}()
	go func() {
		defer wg.Done()
		cache.SetTexture(nil, 1, nil)
	}()

	wg.Wait()
}
