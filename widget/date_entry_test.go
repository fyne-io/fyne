package widget

import (
	"testing"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"
	"github.com/stretchr/testify/assert"
)

func TestDateEntry_SetDate(t *testing.T) {
	test.NewApp()

	e := NewDateEntry()
	date := time.Date(2026, 8, 28, 0, 0, 0, 0, time.UTC)

	e.SetDate(&date)
	assert.NotNil(t, e.Date)
	assert.Equal(t, date.Format(getLocaleDateFormat()), e.Text)

	e.SetDate(nil)
	assert.Nil(t, e.Date)
	assert.Equal(t, "", e.Text)
}

func TestDateEntry_EnableDisable(t *testing.T) {
	test.NewApp()

	e := NewDateEntry()
	_ = e.CreateRenderer()

	e.Disable()
	assert.True(t, e.Disabled())

	e.Enable()
	assert.False(t, e.Disabled())
}

func TestDateEntry_MinSize(t *testing.T) {
	e := NewDateEntry()
	size := e.MinSize()

	assert.True(t, size.Width > 0)
	assert.True(t, size.Height > 0)
}

func TestDateEntry_MoveAndResize(t *testing.T) {
	w := test.NewWindow(nil)
	defer w.Close()

	e := NewDateEntry()
	w.SetContent(e)

	_ = e.CreateRenderer()
	btn := e.ActionItem.(*Button)
	btn.OnTapped()

	e.Move(fyne.NewPos(10, 20))
	assert.Equal(t, float32(10), e.Position().X)
	assert.Equal(t, float32(20), e.Position().Y)

	e.Resize(fyne.NewSize(120, 40))
	assert.Equal(t, float32(120), e.Size().Width)
	assert.Equal(t, float32(40), e.Size().Height)
}

func TestDateEntry_OnChangedAndValidation(t *testing.T) {
	test.NewApp()

	var changedDate *time.Time
	e := NewDateEntry()
	e.OnChanged = func(d *time.Time) {
		changedDate = d
	}
	_ = e.CreateRenderer()

	dateFormat := getLocaleDateFormat()
	targetDate := time.Date(2026, 8, 28, 0, 0, 0, 0, time.UTC)
	formatted := targetDate.Format(dateFormat)

	assert.NoError(t, e.Validator(formatted))
	assert.Error(t, e.Validator("invalid-date"))

	e.SetText(formatted)
	assert.NotNil(t, changedDate)

	e.SetText("invalid-date")

	e.SetText("")
	assert.Nil(t, e.Date)
	assert.Nil(t, changedDate)
}

func TestDateEntry_CoverageEdges(_ *testing.T) {
	e := NewDateEntry()
	_ = e.CreateRenderer()

	e.Disable()
	e.Enable()

	w := test.NewWindow(e)
	w.Resize(fyne.NewSize(200, 200))

	e.popUp = NewPopUp(NewLabel("dummy"), w.Canvas())

	e.Move(fyne.NewPos(10, 20))
	e.Resize(fyne.NewSize(120, 40))

	d := time.Now()
	e.SetDate(&d)

	e.OnChanged = func(d *time.Time) {}
	e.Entry.OnChanged("invalid")
	e.Entry.OnChanged("")
	e.Entry.OnChanged(d.Format("02/01/2006"))
}
