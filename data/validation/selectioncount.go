package validation

import (
	"fmt"

	"fyne.io/fyne/v2"
)

// NewSelectionCount returns a fyne.CountValidator that requires
// the given count to be within the set min/max range. Setting
// a negative max value will allow any number larger than min.
func NewSelectionCount(min, max int) fyne.CountValidator {
	var err error
	if min == max {
		err = fmt.Errorf("incorrect number of selections, expected %d", min)
	} else if max < 0 {
		err = fmt.Errorf("too few selections, expected %d or more", min)
	} else {
		err = fmt.Errorf("too few selections, expected between %d and %d", min, max)
	}

	return func(count int) error {
		if count < min || (max >= 0 && count > max) {
			return err
		}

		return nil
	}
}
