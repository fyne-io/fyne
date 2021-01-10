package widget_test

import (
	"errors"
	"testing"

	"fyne.io/fyne"
	"fyne.io/fyne/data/validation"
	"fyne.io/fyne/test"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"

	"github.com/stretchr/testify/assert"
)

var validator = validation.NewRegexp(`^\d{4}-\d{2}-\d{2}$`, "Input is not a valid date")

func TestEntry_ValidatedEntry(t *testing.T) {
	entry, window := setupImageTest(t, false)
	defer teardownImageTest(window)
	c := window.Canvas()

	r := validation.NewRegexp(`^\d{4}-\d{2}-\d{2}`, "Input is not a valid date")
	entry.Validator = r
	test.AssertRendersToMarkup(t, `
		<canvas padded size="150x200">
			<content>
				<widget pos="10,10" size="120x37" type="*widget.Entry">
					<rectangle fillColor="shadow" pos="0,33" size="120x4"/>
					<widget pos="4,4" size="112x29" type="*widget.ScrollContainer">
						<widget size="112x29" type="*widget.entryContent">
							<widget size="112x29" type="*widget.textProvider">
								<text color="placeholder" pos="4,4" size="104x21"></text>
							</widget>
							<widget size="112x29" type="*widget.textProvider">
								<text pos="4,4" size="104x21"></text>
							</widget>
						</widget>
					</widget>
				</widget>
			</content>
		</canvas>
	`, c)

	test.Type(entry, "2020-02")
	assert.Error(t, r(entry.Text))
	entry.FocusLost()
	test.AssertRendersToMarkup(t, `
		<canvas padded size="150x200">
			<content>
				<widget pos="10,10" size="120x37" type="*widget.Entry">
					<rectangle fillColor="error" pos="0,33" size="120x4"/>
					<widget pos="4,4" size="112x29" type="*widget.ScrollContainer">
						<widget size="112x29" type="*widget.entryContent">
							<widget size="112x29" type="*widget.textProvider">
								<text pos="4,4" size="104x21">2020-02</text>
							</widget>
						</widget>
					</widget>
					<widget pos="92,8" size="20x20" type="*widget.validationStatus">
						<image rsc="errorIcon" size="iconInlineSize" themed="error"/>
					</widget>
				</widget>
			</content>
		</canvas>
	`, c)

	test.Type(entry, "-12")
	assert.NoError(t, r(entry.Text))
	test.AssertRendersToMarkup(t, `
		<canvas padded size="150x200">
			<content>
				<widget pos="10,10" size="120x37" type="*widget.Entry">
					<rectangle fillColor="focus" pos="0,33" size="120x4"/>
					<widget pos="4,4" size="112x29" type="*widget.ScrollContainer">
						<widget size="112x29" type="*widget.entryContent">
							<widget size="112x29" type="*widget.textProvider">
								<text pos="4,4" size="104x21">2020-02-12</text>
							</widget>
							<rectangle fillColor="focus" pos="83,4" size="2x21"/>
						</widget>
					</widget>
					<widget pos="92,8" size="20x20" type="*widget.validationStatus">
						<image rsc="confirmIcon" size="iconInlineSize"/>
					</widget>
				</widget>
			</content>
		</canvas>
	`, c)
}

func TestEntry_Validate(t *testing.T) {
	entry := widget.NewEntry()
	entry.Validator = validator

	test.Type(entry, "2020-02")
	assert.Error(t, entry.Validate())
	assert.Equal(t, entry.Validate(), entry.Validator(entry.Text))

	test.Type(entry, "-12")
	assert.NoError(t, entry.Validate())
	assert.Equal(t, entry.Validate(), entry.Validator(entry.Text))

	entry.SetText("incorrect")
	assert.Error(t, entry.Validate())
	assert.Equal(t, entry.Validate(), entry.Validator(entry.Text))
}

func TestEntry_SetValidationError(t *testing.T) {
	entry, window := setupImageTest(t, false)
	fyne.CurrentApp().Settings().SetTheme(theme.LightTheme())
	defer teardownImageTest(window)
	c := window.Canvas()

	entry.Validator = validator

	entry.SetText("2020-30-30")
	entry.SetValidationError(errors.New("set invalid"))
	test.AssertImageMatches(t, "entry/validation_set_invalid.png", c.Capture())

	entry.SetText("set valid")
	entry.SetValidationError(nil)
	test.AssertImageMatches(t, "entry/validation_set_valid.png", c.Capture())
}

func TestEntry_SetOnValidationChanged(t *testing.T) {
	entry := widget.NewEntry()
	entry.Validator = validator

	modified := false
	entry.SetOnValidationChanged(func(err error) {
		assert.Equal(t, err, entry.Validator(entry.Text))
		modified = true
	})

	test.Type(entry, "2020")
	assert.True(t, modified)

	modified = false
	test.Type(entry, "-01-01")
	assert.True(t, modified)
}
