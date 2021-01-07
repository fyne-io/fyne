package widget_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"fyne.io/fyne"
	"fyne.io/fyne/data/binding"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/test"
	"fyne.io/fyne/widget"
)

func TestCheck_Binding(t *testing.T) {
	c := widget.NewCheck("", nil)
	c.SetChecked(true)
	assert.Equal(t, true, c.Checked)

	val := binding.NewBool()
	c.Bind(val)
	waitForBinding()
	assert.Equal(t, false, c.Checked)

	err := val.Set(true)
	assert.Nil(t, err)
	waitForBinding()
	assert.Equal(t, true, c.Checked)

	c.SetChecked(false)
	v, err := val.Get()
	assert.Nil(t, err)
	assert.Equal(t, false, v)

	c.Unbind()
	waitForBinding()
	assert.Equal(t, false, c.Checked)
}

func TestCheck_Layout(t *testing.T) {
	test.NewApp()
	defer test.NewApp()

	for name, tt := range map[string]struct {
		text     string
		checked  bool
		disabled bool
		want     string
	}{
		"checked": {
			text:    "Test",
			checked: true,
			want: `
				<canvas padded size="150x200">
					<content>
						<container pos="4,4" size="142x192">
							<widget pos="36,81" size="70x29" type="*widget.Check">
								<circle fillColor="background" size="28x28"/>
								<image pos="4,4" rsc="checkButtonCheckedIcon" size="iconInlineSize"/>
								<text pos="32,0" size="42x29">Test</text>
							</widget>
						</container>
					</content>
				</canvas>
			`,
		},
		"unchecked": {
			text: "Test",
			want: `
				<canvas padded size="150x200">
					<content>
						<container pos="4,4" size="142x192">
							<widget pos="36,81" size="70x29" type="*widget.Check">
								<circle fillColor="background" size="28x28"/>
								<image pos="4,4" rsc="checkButtonIcon" size="iconInlineSize"/>
								<text pos="32,0" size="42x29">Test</text>
							</widget>
						</container>
					</content>
				</canvas>
			`,
		},
		"checked_disabled": {
			text:     "Test",
			checked:  true,
			disabled: true,
			want: `
				<canvas padded size="150x200">
					<content>
						<container pos="4,4" size="142x192">
							<widget pos="36,81" size="70x29" type="*widget.Check">
								<circle fillColor="background" size="28x28"/>
								<image pos="4,4" rsc="checkButtonCheckedIcon" size="iconInlineSize" themed="disabled"/>
								<text color="disabled" pos="32,0" size="42x29">Test</text>
							</widget>
						</container>
					</content>
				</canvas>
			`,
		},
		"unchecked_disabled": {
			text:     "Test",
			disabled: true,
			want: `
				<canvas padded size="150x200">
					<content>
						<container pos="4,4" size="142x192">
							<widget pos="36,81" size="70x29" type="*widget.Check">
								<circle fillColor="background" size="28x28"/>
								<image pos="4,4" rsc="checkButtonIcon" size="iconInlineSize" themed="disabled"/>
								<text color="disabled" pos="32,0" size="42x29">Test</text>
							</widget>
						</container>
					</content>
				</canvas>
			`,
		},
	} {
		t.Run(name, func(t *testing.T) {
			check := &widget.Check{
				Text:    tt.text,
				Checked: tt.checked,
			}
			if tt.disabled {
				check.Disable()
			}

			window := test.NewWindow(fyne.NewContainerWithLayout(layout.NewCenterLayout(), check))
			window.Resize(check.MinSize().Max(fyne.NewSize(150, 200)))

			test.AssertRendersToMarkup(t, tt.want, window.Canvas())

			window.Close()
		})
	}
}

func TestNewCheckWithData(t *testing.T) {
	val := binding.NewBool()
	err := val.Set(true)
	assert.Nil(t, err)

	c := widget.NewCheckWithData("", val)
	waitForBinding()
	assert.Equal(t, true, c.Checked)

	c.SetChecked(false)
	v, err := val.Get()
	assert.Nil(t, err)
	assert.Equal(t, false, v)
}
