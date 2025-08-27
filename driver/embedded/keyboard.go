package embedded

import "fyne.io/fyne/v2"

type KeyDirection uint8

const (
	KeyPressed KeyDirection = iota
	KeyReleased
)

type KeyEvent struct {
	Name      fyne.KeyName
	Direction KeyDirection
}

func (d *KeyEvent) isEvent() {}

type CharacterEvent struct {
	Rune rune
}

func (c *CharacterEvent) isEvent() {}
