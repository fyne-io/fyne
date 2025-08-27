package embedded_test

import (
	"image"
	"testing"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/embedded"
	intNoos "fyne.io/fyne/v2/internal/driver/embedded"
	"github.com/stretchr/testify/assert"
)

func TestNoOSDriver(t *testing.T) {
	count := 0
	render := func(img image.Image) {
		count++
	}
	queue := make(chan embedded.Event, 1)
	d := intNoos.NewNoOSDriver(render, queue)

	go func() {
		time.Sleep(100 * time.Millisecond)
		close(queue)
	}()
	d.Run()
	assert.Equal(t, 0, count)

	queue = make(chan embedded.Event, 1)
	d = intNoos.NewNoOSDriver(render, queue)
	_ = d.CreateWindow("Test")
	go func() {
		time.Sleep(100 * time.Millisecond)
		close(queue)
	}()
	d.Run()
	assert.Equal(t, 1, count)

	count = 0
	queue = make(chan embedded.Event, 1)
	d = intNoos.NewNoOSDriver(render, queue)
	w := d.CreateWindow("Test")
	keyed := make(chan bool)
	w.Canvas().SetOnTypedKey(func(k *fyne.KeyEvent) {
		keyed <- true
	})
	go func() {
		queue <- &embedded.CharacterEvent{
			Rune: 'a',
		}
		queue <- &embedded.KeyEvent{
			Name: fyne.KeyEscape,
		}
		<-keyed
		close(queue)
	}()
	d.Run()
	assert.Equal(t, 3, count)
}
