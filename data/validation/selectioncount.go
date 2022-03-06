package validation

import (
	"fmt"

	"fyne.io/fyne/v2"
)

// NewSelectionCount returns a fyne.SelectionValidator that requires
// that "required" amount or more items are selected.
func NewSelectionCount(required int) fyne.SelectionValidator {
	err := fmt.Errorf("too few selections, expected %d or more", required)
	return func(got int) error {
		if got < required {
			return err
		}

		return nil
	}
}
