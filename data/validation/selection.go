package validation

import (
	"fmt"

	"fyne.io/fyne/v2"
)

// NewSelection returns a fyne.SelectionValidator that requires
// that "require" amount or more items are selected.
func NewSelection(require int) fyne.SelectionValidator {
	err := fmt.Errorf("too few selections, expected %d or more", require)
	return func(got int) error {
		if got < require {
			return err
		}

		return nil
	}
}
