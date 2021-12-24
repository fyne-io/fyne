package icns

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"image"
	"io"
	"io/ioutil"
	"sort"
)

var jpeg2000header = []byte{0x00, 0x00, 0x00, 0x0c, 0x6a, 0x50, 0x20, 0x20}

// Decode finds the largest icon listed in the icns file and returns it,
// ignoring all other sizes. The format returned will be whatever the icon data
// is, typically jpeg or png.
func Decode(r io.Reader) (image.Image, error) {
	icons, err := decode(r)
	if err != nil {
		return nil, err
	}
	sort.Slice(icons, func(ii, jj int) bool {
		return icons[ii].OsType.Size > icons[jj].OsType.Size
	})
	img, _, err := image.Decode(icons[0].r)
	if err != nil {
		return nil, fmt.Errorf("decoding largest image: %w", err)
	}
	return img, nil
}

// DecodeAll extracts all icon resolutions present in the icns data.
func DecodeAll(r io.Reader) (images []image.Image, err error) {
	icons, err := decode(r)
	if err != nil {
		return nil, err
	}
	for _, icon := range icons {
		img, _, err := image.Decode(icon.r)
		if err != nil {
			return nil, fmt.Errorf("decoding %q icon: %w", icon.OsType, err)
		}
		images = append(images, img)
	}
	sort.Slice(images, func(ii, jj int) bool {
		var (
			left  = images[ii].Bounds().Size()
			right = images[jj].Bounds().Size()
		)
		return (left.X * left.Y) > (right.X * right.Y)
	})
	return images, nil
}

func decode(r io.Reader) (icons []iconReader, err error) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	var (
		header   = data[0:4]
		fileSize = binary.BigEndian.Uint32(data[4:8])
		read     = uint32(8)
	)
	if string(header) != "icns" {
		return nil, fmt.Errorf("invalid header for icns file")
	}
	for read < fileSize {
		next := data[read : read+4]
		read += 4
		switch string(next) {
		case "TOC ":
			tocSize := binary.BigEndian.Uint32(data[read : read+4])
			read += tocSize - 4 // size includes header and size fields
			continue
		case "icnV":
			read += 4
			continue
		}
		dataSize := binary.BigEndian.Uint32(data[read : read+4])
		read += 4
		if dataSize == 0 {
			continue // no content, we're not interested
		}
		iconData := data[read : read+dataSize-8]
		read += dataSize - 8 // size includes header and size fields
		if isOsType(string(next)) {
			if bytes.Equal(iconData[:8], jpeg2000header) {
				continue // skipping JPEG2000
			}
			icons = append(icons, iconReader{
				OsType: osTypeFromID(string(next)),
				r:      bytes.NewBuffer(iconData),
			})
		}
	}
	if len(icons) == 0 {
		return nil, fmt.Errorf("no icons found")
	}
	return icons, nil
}

type iconReader struct {
	OsType
	r io.Reader
}

func isOsType(ID string) bool {
	_, ok := getTypeFromID(ID)
	return ok
}

func init() {
	image.RegisterFormat("icns", "icns", Decode, nil)
}
