package metadata

import (
	"bytes"
	"image/png"
	"strconv"

	"github.com/nfnt/resize"

	"fyne.io/fyne/v2"
)

func ScaleIcon(data fyne.Resource, size int) fyne.Resource {
	img, err := png.Decode(bytes.NewReader(data.Content()))
	if err != nil {
		fyne.LogError("Failed to decode app icon", err)
		return data
	}

	if img.Bounds().Dx() <= size {
		return data
	}

	sized := resize.Resize(uint(size), uint(size), img, resize.Lanczos3)
	smallData := &bytes.Buffer{}
	err = png.Encode(smallData, sized)
	if err != nil {
		fyne.LogError("Failed to encode smaller app icon", err)
		return data
	}

	name := data.Name()
	nameLen := len(name)
	suffix := "-" + strconv.Itoa(size) + ".png"
	if nameLen <= 4 || name[nameLen-4] != '.' {
		name = "appicon" + suffix
	} else {
		name = name[:nameLen-4] + suffix
	}
	return fyne.NewStaticResource(name, smallData.Bytes())
}
