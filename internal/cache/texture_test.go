package cache_test

import (
	"sync"
	"testing"

	"fyne.io/fyne/v2/internal/cache"
)

// go test -race
func TestTexture(t *testing.T) {
	wg := sync.WaitGroup{}

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
