package noos

import "fyne.io/fyne/v2"

type KeyDirection uint8

const (
	KeyPressed KeyDirection = iota
	KeyReleased
)

type HardwareKeyEvent struct {
	Name      fyne.KeyName
	Direction KeyDirection
}

func (d *HardwareKeyEvent) isEvent() {}
