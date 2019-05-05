package theme

import (
	"bytes"
	"encoding/xml"
	"image/color"
	"io/ioutil"
	"testing"

	"fyne.io/fyne"
	"github.com/stretchr/testify/assert"
)

func TestSVG_ReplaceFillColor(t *testing.T) {
	src, err := ioutil.ReadFile("testdata/cancel_Paths.svg")
	if err != nil {
		t.Fatal(err)
	}
	sRes := fyne.NewStaticResource("cancel", src)
	red := color.RGBA{0xff, 0x00, 0x00, 0xff}
	rdr := bytes.NewReader(src)
	var s SVG
	if err := s.ReplaceFillColor(rdr, red); err != nil {
		t.Fatal(err)
	}
	res, err := xml.Marshal(s)
	if err != nil {
		t.Fatal(err)
	}
	assert.NotEqual(t, string(sRes.Content()), string(res))
}
