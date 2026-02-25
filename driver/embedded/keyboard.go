package embedded

import "fyne.io/fyne/v2"

// KeyDirection specifies the press/release of a key event
//
// Since: 2.7
type KeyDirection uint8

const (
	// KeyPressed specifies that a key was pushed down.
	//
	// Since: 2.7
	KeyPressed KeyDirection = iota

	// KeyReleased indicates a key was let back up.
	//
	// Since: 2.7
	KeyReleased
)

// KeyEvent is an event from keyboard actions occurring in an embedded device keyboard.
//
// Since: 2.7
type KeyEvent struct {
	Name      fyne.KeyName
	Direction KeyDirection
}

func (d *KeyEvent) isEvent() {}

// CharacterEvent is an event specifying that a character was created by a hardware or virtual keyboard.
//
// Since: 2.7
type CharacterEvent struct {
	Rune rune
}

func (c *CharacterEvent) isEvent() {}
