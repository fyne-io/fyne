package theme

import (
	"bytes"
	"encoding/xml"
	"image/color"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSVG_ReplaceFillColor(t *testing.T) {
	src, err := ioutil.ReadFile("testdata/cancel_Paths.svg")
	if err != nil {
		t.Fatal(err)
	}
	red := color.NRGBA{0xff, 0x00, 0x00, 0xff}
	rdr := bytes.NewReader(src)
	s, err := svgFromXML(rdr)
	if err != nil {
		t.Fatal(err)
	}
	if err := s.replaceFillColor(red); err != nil {
		t.Fatal(err)
	}
	res, err := xml.Marshal(s)
	if err != nil {
		t.Fatal(err)
	}
	assert.NotEqual(t, string(src), string(res))
	assert.True(t, strings.Contains(string(res), "#ff0000"))
}

func TestSVG_ReplaceFillColor_Ellipse(t *testing.T) {
	src, err := ioutil.ReadFile("testdata/ellipse.svg")
	if err != nil {
		t.Fatal(err)
	}
	red := color.NRGBA{0xff, 0x00, 0x00, 0xff}
	rdr := bytes.NewReader(src)
	s, err := svgFromXML(rdr)
	if err != nil {
		t.Fatal(err)
	}
	if err := s.replaceFillColor(red); err != nil {
		t.Fatal(err)
	}
	res, err := xml.Marshal(s)
	if err != nil {
		t.Fatal(err)
	}
	assert.NotEqual(t, string(src), string(res))
	assert.True(t, strings.Contains(string(res), "#ff0000"))
}
