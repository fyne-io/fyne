package noos_test

import (
	"image"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/noos"
	intNoos "fyne.io/fyne/v2/internal/driver/noos"
	"github.com/stretchr/testify/assert"
)

func TestNoOSDriver(t *testing.T) {
	count := 0
	render := func(img image.Image) {
		count++
	}
	queue := make(chan noos.Event, 1)
	d := intNoos.NewNoOSDriver(render, queue)

	go d.Run()
	d.Quit()
	assert.Equal(t, 0, count)

	d = intNoos.NewNoOSDriver(render, queue)
	w := d.CreateWindow("Test")
	keyed := make(chan bool)
	w.Canvas().SetOnTypedKey(func(k *fyne.KeyEvent) {
		keyed <- true
	})
	go d.Run()
	queue <- &noos.CharacterEvent{
		Rune: 'a',
	}
	queue <- &noos.KeyEvent{
		Name: fyne.KeyEscape,
	}
	<-keyed
	d.Quit()
	assert.Equal(t, 2, count)
}
