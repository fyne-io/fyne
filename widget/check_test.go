package widget

import (
	"testing"

	"fyne.io/fyne/test"
	"fyne.io/fyne/theme"
	"github.com/stretchr/testify/assert"
)

func TestCheckSize(t *testing.T) {
	check := NewCheck("Hi", nil)
	min := check.MinSize()

	assert.True(t, min.Width > theme.Padding()*2)
	assert.True(t, min.Height > theme.Padding()*2)
}

func TestCheckChecked(t *testing.T) {
	checked := false
	check := NewCheck("Hi", func(on bool) {
		checked = on
	})
	test.Tap(check)

	assert.True(t, checked)
}

func TestCheckUnChecked(t *testing.T) {
	checked := true
	check := NewCheck("Hi", func(on bool) {
		checked = on
	})
	check.Checked = checked
	test.Tap(check)

	assert.False(t, checked)
}

func TestCheckIsDisabledByDefault(t *testing.T) {
	check := NewCheck("", func(on bool) {})
	assert.False(t, check.Checked)
}

func TestCheckIsEnabledAfterUpdating(t *testing.T) {
	check := NewCheck("", func(on bool) {})

	check.SetChecked(true)

	assert.True(t, check.Checked)
}

func TestCheckStateIsCorrectAfterMultipleUpdates(t *testing.T) {
	check := NewCheck("", func(on bool) {})

	checked := false
	for i := 0; i < 5; i++ {
		check.SetChecked(checked)
		assert.True(t, check.Checked == checked)

		checked = !checked
	}
}
