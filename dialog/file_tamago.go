//go:build tamago || noos

package dialog

import "fyne.io/fyne/v2"

func getFavoriteLocations() (map[string]fyne.ListableURI, error) {
	return map[string]fyne.ListableURI{}, nil
}
