package binding

import (
	"fmt"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/storage"
	"github.com/stretchr/testify/assert"
)

func TestSprintfConversionRead(t *testing.T) {
	var b bool = true
	var f float64 = 42
	var i int = 7
	var r rune = 'B'
	var s string = "this is a string"
	var u fyne.URI = storage.NewFileURI("/")

	bb := BindBool(&b)
	bf := BindFloat(&f)
	bi := BindInt(&i)
	br := BindRune(&r)
	bs := BindString(&s)
	bu := BindURI(&u)

	format := "Bool %v, Float %f, Int %i, Rune %v, String '%s', URI '%s'"
	sp := NewSprintf(format, bb, bf, bi, br, bs, bu)
	expected := fmt.Sprintf(format, b, f, i, r, s, u)
	assert.NotNil(t, sp)

	sGenerated, err := sp.Get()
	assert.Nil(t, err)
	assert.NotNil(t, sGenerated)
	assert.Equal(t, expected, sGenerated)

	bb.Set(false)
	bf.Set(7)
	bi.Set(42)
	br.Set('C')
	bs.Set("a different string")

	expectedChange := fmt.Sprintf(format, b, f, i, r, s, u)

	sChange, err := sp.Get()
	assert.Nil(t, err)
	assert.NotNil(t, sChange)
	assert.Equal(t, expectedChange, sChange)
	assert.NotEqual(t, sGenerated, sChange)
}

func TestSprintfConversionReadWrite(t *testing.T) {
	var b bool = true
	var f float64 = 42
	var i int = 7
	var r rune = 'B'
	var s string = "this is a string"
	var u fyne.URI = storage.NewFileURI("/")

	bb := BindBool(&b)
	bf := BindFloat(&f)
	bi := BindInt(&i)
	br := BindRune(&r)
	bs := BindString(&s)
	bu := BindURI(&u)

	format := "Bool %v , Float %f , Int %v , Rune %v , String %s , URI %s"
	sp := NewSprintf(format, bb, bf, bi, br, bs, bu)
	expected := fmt.Sprintf(format, b, f, i, r, s, u)
	assert.NotNil(t, sp)

	sGenerated, err := sp.Get()
	assert.Nil(t, err)
	assert.NotNil(t, sGenerated)
	assert.Equal(t, expected, sGenerated)

	err = sp.Set("Bool false , Float 7.000000 , Int 42 , Rune 67 , String nospacestring , URI file:///var/")
	assert.Nil(t, err)

	assert.Equal(t, b, false)
	assert.Equal(t, f, float64(7))
	assert.Equal(t, i, 42)
	assert.Equal(t, r, 'C')
	assert.Equal(t, s, "nospacestring")
	assert.Equal(t, u.String(), "file:///var/")

	expectedChange := fmt.Sprintf(format, b, f, i, r, s, u)

	sChange, err := sp.Get()

	assert.Nil(t, err)
	assert.NotNil(t, sChange)
	assert.Equal(t, expectedChange, sChange)
	assert.NotEqual(t, sGenerated, sChange)
}

func TestNewStringWithFormat(t *testing.T) {
	s := "this is a string"
	bs := BindString(&s)

	format := "String %s"
	sp := StringToStringWithFormat(bs, format)
	expected := fmt.Sprintf(format, s)
	assert.NotNil(t, sp)

	sGenerated, err := sp.Get()
	assert.Nil(t, err)
	assert.NotNil(t, sGenerated)
	assert.Equal(t, expected, sGenerated)

	err = sp.Set("String nospacestring")
	assert.Nil(t, err)

	assert.Equal(t, s, "nospacestring")
	expectedChange := fmt.Sprintf(format, s)
	sChange, err := sp.Get()
	assert.Nil(t, err)
	assert.NotNil(t, sChange)
	assert.Equal(t, expectedChange, sChange)
	assert.NotEqual(t, sGenerated, sChange)

	emptyFormat := "%s"
	wrapped := StringToStringWithFormat(bs, emptyFormat)
	assert.Equal(t, bs, wrapped)
}
