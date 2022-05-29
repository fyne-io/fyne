package app

import (
	"encoding/base64"
	"io/ioutil"
	"strings"

	"fyne.io/fyne/v2"
)

var meta = fyne.AppMetadata{
	ID:      "com.example",
	Name:    "Fyne App",
	Version: "1.0.0",
	Build:   1,
}

// SetMetadata overrides the packaged application metadata.
// This data can be used in many places like notifications and about screens.
func SetMetadata(m fyne.AppMetadata) {
	meta = m
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
