package app

import (
	"encoding/base64"
	"io/ioutil"
	"strconv"
	"strings"

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
		Icon:    convertIcon(intapp.MetaIcon),
		ID:      intapp.MetaID,
		Name:    intapp.MetaName,
		Version: intapp.MetaVersion,
		Build:   build,
	}
}

func (a *fyneApp) Metadata() fyne.AppMetadata {
	return meta
}

func convertIcon(bytes string) fyne.Resource {
	if bytes == "" {
		return nil
	}
	data := base64.NewDecoder(base64.StdEncoding, strings.NewReader(bytes))
	img, err := ioutil.ReadAll(data)
	if err != nil {
		fyne.LogError("Failed to decode icon from embedded data", err)
		return nil
	}

	return &fyne.StaticResource{
		StaticName:    "FyneAppIcon",
		StaticContent: img,
	}
}
