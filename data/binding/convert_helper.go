package binding

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/storage"
)

func uriFromString(in string) (fyne.URI, error) {
	return storage.ParseURI(in)
}

func uriToString(in fyne.URI) (string, error) {
	if in == nil {
		return "", nil
	}

	return in.String(), nil
}
