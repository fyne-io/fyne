package storage

import (
	"io/ioutil"

	"fyne.io/fyne/v2"
)

// LoadResourceFromURI creates a new StaticResource in memory using the contents of the specified URI.
// The URI will be opened using the current driver, so valid schemas will vary from platform to platform.
// The file:// schema will always work.
func LoadResourceFromURI(u fyne.URI) (fyne.Resource, error) {
	read, err := Reader(u)
	if err != nil {
		return nil, err
	}

	defer read.Close()
	bytes, err := ioutil.ReadAll(read)
	if err != nil {
		return nil, err
	}

	return fyne.NewStaticResource(u.Name(), bytes), nil
}
