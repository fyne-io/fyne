package app

import (
	"strconv"

	"fyne.io/fyne/v2"
	intapp "fyne.io/fyne/v2/internal/app"
)

var (
	meta fyne.AppMetadata
)

func init() {
	build, err := strconv.Atoi(intapp.MetaBuild)
	if err != nil {
		build = 1
	}

	meta = fyne.AppMetadata{
		ID:      intapp.MetaID,
		Name:    intapp.MetaName,
		Version: intapp.MetaVersion,
		Build:   build,
	}
}

func (a *fyneApp) Metadata() fyne.AppMetadata {
	return meta
}
