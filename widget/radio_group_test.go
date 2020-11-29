package widget_test

import (
	"testing"

	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/test"
	"fyne.io/fyne/widget"

	"github.com/stretchr/testify/assert"
)

func TestRadioGroup_FocusRendering(t *testing.T) {
	t.Run("gain/lose focus", func(t *testing.T) {
		radio := widget.NewRadioGroup([]string{"Option A", "Option B", "Option C"}, nil)
		window := test.NewWindow(fyne.NewContainerWithLayout(layout.NewCenterLayout(), radio))
		defer window.Close()
		window.Resize(radio.MinSize().Max(fyne.NewSize(150, 200)))

		canvas := window.Canvas().(test.WindowlessCanvas)
		test.AssertRendersToMarkup(t, `
			<canvas padded size="150x200">
				<content>
					<container pos="4,4" size="142x192">
						<widget pos="19,52" size="103x87" type="*widget.RadioGroup">
							<widget size="103x29" type="*widget.radioItem">
								<circle fillColor="background" size="28x28"/>
								<image pos="4,4" rsc="radioButtonIcon" size="iconInlineSize"/>
								<text pos="32,0" size="103x29">Option A</text>
							</widget>
							<widget pos="0,29" size="103x29" type="*widget.radioItem">
								<circle fillColor="background" size="28x28"/>
								<image pos="4,4" rsc="radioButtonIcon" size="iconInlineSize"/>
								<text pos="32,0" size="103x29">Option B</text>
							</widget>
							<widget pos="0,58" size="103x29" type="*widget.radioItem">
								<circle fillColor="background" size="28x28"/>
								<image pos="4,4" rsc="radioButtonIcon" size="iconInlineSize"/>
								<text pos="32,0" size="103x29">Option C</text>
							</widget>
						</widget>
					</container>
				</content>
			</canvas>
		`, canvas)
		canvas.FocusNext()
		test.AssertRendersToMarkup(t, `
			<canvas padded size="150x200">
				<content>
					<container pos="4,4" size="142x192">
						<widget pos="19,52" size="103x87" type="*widget.RadioGroup">
							<widget size="103x29" type="*widget.radioItem">
								<circle fillColor="focus" size="28x28"/>
								<image pos="4,4" rsc="radioButtonIcon" size="iconInlineSize"/>
								<text pos="32,0" size="103x29">Option A</text>
							</widget>
							<widget pos="0,29" size="103x29" type="*widget.radioItem">
								<circle fillColor="background" size="28x28"/>
								<image pos="4,4" rsc="radioButtonIcon" size="iconInlineSize"/>
								<text pos="32,0" size="103x29">Option B</text>
							</widget>
							<widget pos="0,58" size="103x29" type="*widget.radioItem">
								<circle fillColor="background" size="28x28"/>
								<image pos="4,4" rsc="radioButtonIcon" size="iconInlineSize"/>
								<text pos="32,0" size="103x29">Option C</text>
							</widget>
						</widget>
					</container>
				</content>
			</canvas>
		`, canvas)
		canvas.FocusNext()
		test.AssertRendersToMarkup(t, `
			<canvas padded size="150x200">
				<content>
					<container pos="4,4" size="142x192">
						<widget pos="19,52" size="103x87" type="*widget.RadioGroup">
							<widget size="103x29" type="*widget.radioItem">
								<circle fillColor="background" size="28x28"/>
								<image pos="4,4" rsc="radioButtonIcon" size="iconInlineSize"/>
								<text pos="32,0" size="103x29">Option A</text>
							</widget>
							<widget pos="0,29" size="103x29" type="*widget.radioItem">
								<circle fillColor="focus" size="28x28"/>
								<image pos="4,4" rsc="radioButtonIcon" size="iconInlineSize"/>
								<text pos="32,0" size="103x29">Option B</text>
							</widget>
							<widget pos="0,58" size="103x29" type="*widget.radioItem">
								<circle fillColor="background" size="28x28"/>
								<image pos="4,4" rsc="radioButtonIcon" size="iconInlineSize"/>
								<text pos="32,0" size="103x29">Option C</text>
							</widget>
						</widget>
					</container>
				</content>
			</canvas>
		`, canvas)
		canvas.Unfocus()
		test.AssertRendersToMarkup(t, `
			<canvas padded size="150x200">
				<content>
					<container pos="4,4" size="142x192">
						<widget pos="19,52" size="103x87" type="*widget.RadioGroup">
							<widget size="103x29" type="*widget.radioItem">
								<circle fillColor="background" size="28x28"/>
								<image pos="4,4" rsc="radioButtonIcon" size="iconInlineSize"/>
								<text pos="32,0" size="103x29">Option A</text>
							</widget>
							<widget pos="0,29" size="103x29" type="*widget.radioItem">
								<circle fillColor="background" size="28x28"/>
								<image pos="4,4" rsc="radioButtonIcon" size="iconInlineSize"/>
								<text pos="32,0" size="103x29">Option B</text>
							</widget>
							<widget pos="0,58" size="103x29" type="*widget.radioItem">
								<circle fillColor="background" size="28x28"/>
								<image pos="4,4" rsc="radioButtonIcon" size="iconInlineSize"/>
								<text pos="32,0" size="103x29">Option C</text>
							</widget>
						</widget>
					</container>
				</content>
			</canvas>
		`, canvas)

		radio.SetSelected("Option B")
		test.AssertRendersToMarkup(t, `
			<canvas padded size="150x200">
				<content>
					<container pos="4,4" size="142x192">
						<widget pos="19,52" size="103x87" type="*widget.RadioGroup">
							<widget size="103x29" type="*widget.radioItem">
								<circle fillColor="background" size="28x28"/>
								<image pos="4,4" rsc="radioButtonIcon" size="iconInlineSize"/>
								<text pos="32,0" size="103x29">Option A</text>
							</widget>
							<widget pos="0,29" size="103x29" type="*widget.radioItem">
								<circle fillColor="background" size="28x28"/>
								<image pos="4,4" rsc="radioButtonCheckedIcon" size="iconInlineSize"/>
								<text pos="32,0" size="103x29">Option B</text>
							</widget>
							<widget pos="0,58" size="103x29" type="*widget.radioItem">
								<circle fillColor="background" size="28x28"/>
								<image pos="4,4" rsc="radioButtonIcon" size="iconInlineSize"/>
								<text pos="32,0" size="103x29">Option C</text>
							</widget>
						</widget>
					</container>
				</content>
			</canvas>
		`, canvas)
		canvas.FocusNext()
		test.AssertRendersToMarkup(t, `
			<canvas padded size="150x200">
				<content>
					<container pos="4,4" size="142x192">
						<widget pos="19,52" size="103x87" type="*widget.RadioGroup">
							<widget size="103x29" type="*widget.radioItem">
								<circle fillColor="focus" size="28x28"/>
								<image pos="4,4" rsc="radioButtonIcon" size="iconInlineSize"/>
								<text pos="32,0" size="103x29">Option A</text>
							</widget>
							<widget pos="0,29" size="103x29" type="*widget.radioItem">
								<circle fillColor="background" size="28x28"/>
								<image pos="4,4" rsc="radioButtonCheckedIcon" size="iconInlineSize"/>
								<text pos="32,0" size="103x29">Option B</text>
							</widget>
							<widget pos="0,58" size="103x29" type="*widget.radioItem">
								<circle fillColor="background" size="28x28"/>
								<image pos="4,4" rsc="radioButtonIcon" size="iconInlineSize"/>
								<text pos="32,0" size="103x29">Option C</text>
							</widget>
						</widget>
					</container>
				</content>
			</canvas>
		`, canvas)
		canvas.FocusNext()
		test.AssertRendersToMarkup(t, `
			<canvas padded size="150x200">
				<content>
					<container pos="4,4" size="142x192">
						<widget pos="19,52" size="103x87" type="*widget.RadioGroup">
							<widget size="103x29" type="*widget.radioItem">
								<circle fillColor="background" size="28x28"/>
								<image pos="4,4" rsc="radioButtonIcon" size="iconInlineSize"/>
								<text pos="32,0" size="103x29">Option A</text>
							</widget>
							<widget pos="0,29" size="103x29" type="*widget.radioItem">
								<circle fillColor="focus" size="28x28"/>
								<image pos="4,4" rsc="radioButtonCheckedIcon" size="iconInlineSize"/>
								<text pos="32,0" size="103x29">Option B</text>
							</widget>
							<widget pos="0,58" size="103x29" type="*widget.radioItem">
								<circle fillColor="background" size="28x28"/>
								<image pos="4,4" rsc="radioButtonIcon" size="iconInlineSize"/>
								<text pos="32,0" size="103x29">Option C</text>
							</widget>
						</widget>
					</container>
				</content>
			</canvas>
		`, canvas)
		canvas.Unfocus()
		test.AssertRendersToMarkup(t, `
			<canvas padded size="150x200">
				<content>
					<container pos="4,4" size="142x192">
						<widget pos="19,52" size="103x87" type="*widget.RadioGroup">
							<widget size="103x29" type="*widget.radioItem">
								<circle fillColor="background" size="28x28"/>
								<image pos="4,4" rsc="radioButtonIcon" size="iconInlineSize"/>
								<text pos="32,0" size="103x29">Option A</text>
							</widget>
							<widget pos="0,29" size="103x29" type="*widget.radioItem">
								<circle fillColor="background" size="28x28"/>
								<image pos="4,4" rsc="radioButtonCheckedIcon" size="iconInlineSize"/>
								<text pos="32,0" size="103x29">Option B</text>
							</widget>
							<widget pos="0,58" size="103x29" type="*widget.radioItem">
								<circle fillColor="background" size="28x28"/>
								<image pos="4,4" rsc="radioButtonIcon" size="iconInlineSize"/>
								<text pos="32,0" size="103x29">Option C</text>
							</widget>
						</widget>
					</container>
				</content>
			</canvas>
		`, canvas)
	})

	t.Run("disable/enable focused", func(t *testing.T) {
		radio := &widget.RadioGroup{Options: []string{"Option A", "Option B", "Option C"}}
		window := test.NewWindow(fyne.NewContainerWithLayout(layout.NewCenterLayout(), radio))
		defer window.Close()
		window.Resize(radio.MinSize().Max(fyne.NewSize(150, 200)))

		canvas := window.Canvas().(test.WindowlessCanvas)
		canvas.FocusNext()
		radio.Disable()
		test.AssertRendersToMarkup(t, `
			<canvas padded size="150x200">
				<content>
					<container pos="4,4" size="142x192">
						<widget pos="19,52" size="103x87" type="*widget.RadioGroup">
							<widget size="103x29" type="*widget.radioItem">
								<circle fillColor="background" size="28x28"/>
								<image pos="4,4" rsc="radioButtonIcon" size="iconInlineSize" themed="disabled"/>
								<text color="disabled" pos="32,0" size="103x29">Option A</text>
							</widget>
							<widget pos="0,29" size="103x29" type="*widget.radioItem">
								<circle fillColor="background" size="28x28"/>
								<image pos="4,4" rsc="radioButtonIcon" size="iconInlineSize" themed="disabled"/>
								<text color="disabled" pos="32,0" size="103x29">Option B</text>
							</widget>
							<widget pos="0,58" size="103x29" type="*widget.radioItem">
								<circle fillColor="background" size="28x28"/>
								<image pos="4,4" rsc="radioButtonIcon" size="iconInlineSize" themed="disabled"/>
								<text color="disabled" pos="32,0" size="103x29">Option C</text>
							</widget>
						</widget>
					</container>
				</content>
			</canvas>
		`, canvas)
		radio.Enable()
		test.AssertRendersToMarkup(t, `
			<canvas padded size="150x200">
				<content>
					<container pos="4,4" size="142x192">
						<widget pos="19,52" size="103x87" type="*widget.RadioGroup">
							<widget size="103x29" type="*widget.radioItem">
								<circle fillColor="focus" size="28x28"/>
								<image pos="4,4" rsc="radioButtonIcon" size="iconInlineSize"/>
								<text pos="32,0" size="103x29">Option A</text>
							</widget>
							<widget pos="0,29" size="103x29" type="*widget.radioItem">
								<circle fillColor="background" size="28x28"/>
								<image pos="4,4" rsc="radioButtonIcon" size="iconInlineSize"/>
								<text pos="32,0" size="103x29">Option B</text>
							</widget>
							<widget pos="0,58" size="103x29" type="*widget.radioItem">
								<circle fillColor="background" size="28x28"/>
								<image pos="4,4" rsc="radioButtonIcon" size="iconInlineSize"/>
								<text pos="32,0" size="103x29">Option C</text>
							</widget>
						</widget>
					</container>
				</content>
			</canvas>
		`, canvas)
	})
}

func TestRadioGroup_Layout(t *testing.T) {
	for name, tt := range map[string]struct {
		disabled   bool
		horizontal bool
		options    []string
		selected   string
		want       string
	}{
		"single": {
			options: []string{"Test"},
			want: `
				<canvas padded size="150x200">
					<content>
						<container pos="4,4" size="142x192">
							<widget pos="36,81" size="70x29" type="*widget.RadioGroup">
								<widget size="70x29" type="*widget.radioItem">
									<circle fillColor="background" size="28x28"/>
									<image pos="4,4" rsc="radioButtonIcon" size="iconInlineSize"/>
									<text pos="32,0" size="70x29">Test</text>
								</widget>
							</widget>
						</container>
					</content>
				</canvas>
			`,
		},
		"single_disabled": {
			disabled: true,
			options:  []string{"Test"},
			want: `
				<canvas padded size="150x200">
					<content>
						<container pos="4,4" size="142x192">
							<widget pos="36,81" size="70x29" type="*widget.RadioGroup">
								<widget size="70x29" type="*widget.radioItem">
									<circle fillColor="background" size="28x28"/>
									<image pos="4,4" rsc="radioButtonIcon" size="iconInlineSize" themed="disabled"/>
									<text color="disabled" pos="32,0" size="70x29">Test</text>
								</widget>
							</widget>
						</container>
					</content>
				</canvas>
			`,
		},
		"single_horizontal": {
			horizontal: true,
			options:    []string{"Test"},
			want: `
				<canvas padded size="150x200">
					<content>
						<container pos="4,4" size="142x192">
							<widget pos="36,81" size="70x29" type="*widget.RadioGroup">
								<widget size="70x29" type="*widget.radioItem">
									<circle fillColor="background" size="28x28"/>
									<image pos="4,4" rsc="radioButtonIcon" size="iconInlineSize"/>
									<text pos="32,0" size="70x29">Test</text>
								</widget>
							</widget>
						</container>
					</content>
				</canvas>
			`,
		},
		"single_horizontal_disabled": {
			disabled:   true,
			horizontal: true,
			options:    []string{"Test"},
			want: `
				<canvas padded size="150x200">
					<content>
						<container pos="4,4" size="142x192">
							<widget pos="36,81" size="70x29" type="*widget.RadioGroup">
								<widget size="70x29" type="*widget.radioItem">
									<circle fillColor="background" size="28x28"/>
									<image pos="4,4" rsc="radioButtonIcon" size="iconInlineSize" themed="disabled"/>
									<text color="disabled" pos="32,0" size="70x29">Test</text>
								</widget>
							</widget>
						</container>
					</content>
				</canvas>
			`,
		},
		"single_selected": {
			options:  []string{"Test"},
			selected: "Test",
			want: `
				<canvas padded size="150x200">
					<content>
						<container pos="4,4" size="142x192">
							<widget pos="36,81" size="70x29" type="*widget.RadioGroup">
								<widget size="70x29" type="*widget.radioItem">
									<circle fillColor="background" size="28x28"/>
									<image pos="4,4" rsc="radioButtonCheckedIcon" size="iconInlineSize"/>
									<text pos="32,0" size="70x29">Test</text>
								</widget>
							</widget>
						</container>
					</content>
				</canvas>
			`,
		},
		"single_selected_disabled": {
			disabled: true,
			options:  []string{"Test"},
			selected: "Test",
			want: `
				<canvas padded size="150x200">
					<content>
						<container pos="4,4" size="142x192">
							<widget pos="36,81" size="70x29" type="*widget.RadioGroup">
								<widget size="70x29" type="*widget.radioItem">
									<circle fillColor="background" size="28x28"/>
									<image pos="4,4" rsc="radioButtonCheckedIcon" size="iconInlineSize" themed="disabled"/>
									<text color="disabled" pos="32,0" size="70x29">Test</text>
								</widget>
							</widget>
						</container>
					</content>
				</canvas>
			`,
		},
		"single_selected_horizontal": {
			horizontal: true,
			options:    []string{"Test"},
			selected:   "Test",
			want: `
				<canvas padded size="150x200">
					<content>
						<container pos="4,4" size="142x192">
							<widget pos="36,81" size="70x29" type="*widget.RadioGroup">
								<widget size="70x29" type="*widget.radioItem">
									<circle fillColor="background" size="28x28"/>
									<image pos="4,4" rsc="radioButtonCheckedIcon" size="iconInlineSize"/>
									<text pos="32,0" size="70x29">Test</text>
								</widget>
							</widget>
						</container>
					</content>
				</canvas>
			`,
		},
		"single_selected_horizontal_disabled": {
			disabled:   true,
			horizontal: true,
			options:    []string{"Test"},
			selected:   "Test",
			want: `
				<canvas padded size="150x200">
					<content>
						<container pos="4,4" size="142x192">
							<widget pos="36,81" size="70x29" type="*widget.RadioGroup">
								<widget size="70x29" type="*widget.radioItem">
									<circle fillColor="background" size="28x28"/>
									<image pos="4,4" rsc="radioButtonCheckedIcon" size="iconInlineSize" themed="disabled"/>
									<text color="disabled" pos="32,0" size="70x29">Test</text>
								</widget>
							</widget>
						</container>
					</content>
				</canvas>
			`,
		},
		"multiple": {
			options: []string{"Foo", "Bar"},
			want: `
				<canvas padded size="150x200">
					<content>
						<container pos="4,4" size="142x192">
							<widget pos="37,67" size="67x58" type="*widget.RadioGroup">
								<widget size="67x29" type="*widget.radioItem">
									<circle fillColor="background" size="28x28"/>
									<image pos="4,4" rsc="radioButtonIcon" size="iconInlineSize"/>
									<text pos="32,0" size="67x29">Foo</text>
								</widget>
								<widget pos="0,29" size="67x29" type="*widget.radioItem">
									<circle fillColor="background" size="28x28"/>
									<image pos="4,4" rsc="radioButtonIcon" size="iconInlineSize"/>
									<text pos="32,0" size="67x29">Bar</text>
								</widget>
							</widget>
						</container>
					</content>
				</canvas>
			`,
		},
		"multiple_disabled": {
			disabled: true,
			options:  []string{"Foo", "Bar"},
			want: `
				<canvas padded size="150x200">
					<content>
						<container pos="4,4" size="142x192">
							<widget pos="37,67" size="67x58" type="*widget.RadioGroup">
								<widget size="67x29" type="*widget.radioItem">
									<circle fillColor="background" size="28x28"/>
									<image pos="4,4" rsc="radioButtonIcon" size="iconInlineSize" themed="disabled"/>
									<text color="disabled" pos="32,0" size="67x29">Foo</text>
								</widget>
								<widget pos="0,29" size="67x29" type="*widget.radioItem">
									<circle fillColor="background" size="28x28"/>
									<image pos="4,4" rsc="radioButtonIcon" size="iconInlineSize" themed="disabled"/>
									<text color="disabled" pos="32,0" size="67x29">Bar</text>
								</widget>
							</widget>
						</container>
					</content>
				</canvas>
			`,
		},
		"multiple_horizontal": {
			horizontal: true,
			options:    []string{"Foo", "Bar"},
			want: `
				<canvas padded size="150x200">
					<content>
						<container pos="4,4" size="142x192">
							<widget pos="5,81" size="132x29" type="*widget.RadioGroup">
								<widget size="66x29" type="*widget.radioItem">
									<circle fillColor="background" size="28x28"/>
									<image pos="4,4" rsc="radioButtonIcon" size="iconInlineSize"/>
									<text pos="32,0" size="66x29">Foo</text>
								</widget>
								<widget pos="66,0" size="66x29" type="*widget.radioItem">
									<circle fillColor="background" size="28x28"/>
									<image pos="4,4" rsc="radioButtonIcon" size="iconInlineSize"/>
									<text pos="32,0" size="66x29">Bar</text>
								</widget>
							</widget>
						</container>
					</content>
				</canvas>
			`,
		},
		"multiple_horizontal_disabled": {
			disabled:   true,
			horizontal: true,
			options:    []string{"Foo", "Bar"},
			want: `
				<canvas padded size="150x200">
					<content>
						<container pos="4,4" size="142x192">
							<widget pos="5,81" size="132x29" type="*widget.RadioGroup">
								<widget size="66x29" type="*widget.radioItem">
									<circle fillColor="background" size="28x28"/>
									<image pos="4,4" rsc="radioButtonIcon" size="iconInlineSize" themed="disabled"/>
									<text color="disabled" pos="32,0" size="66x29">Foo</text>
								</widget>
								<widget pos="66,0" size="66x29" type="*widget.radioItem">
									<circle fillColor="background" size="28x28"/>
									<image pos="4,4" rsc="radioButtonIcon" size="iconInlineSize" themed="disabled"/>
									<text color="disabled" pos="32,0" size="66x29">Bar</text>
								</widget>
							</widget>
						</container>
					</content>
				</canvas>
			`,
		},
		"multiple_selected": {
			options:  []string{"Foo", "Bar"},
			selected: "Foo",
			want: `
				<canvas padded size="150x200">
					<content>
						<container pos="4,4" size="142x192">
							<widget pos="37,67" size="67x58" type="*widget.RadioGroup">
								<widget size="67x29" type="*widget.radioItem">
									<circle fillColor="background" size="28x28"/>
									<image pos="4,4" rsc="radioButtonCheckedIcon" size="iconInlineSize"/>
									<text pos="32,0" size="67x29">Foo</text>
								</widget>
								<widget pos="0,29" size="67x29" type="*widget.radioItem">
									<circle fillColor="background" size="28x28"/>
									<image pos="4,4" rsc="radioButtonIcon" size="iconInlineSize"/>
									<text pos="32,0" size="67x29">Bar</text>
								</widget>
							</widget>
						</container>
					</content>
				</canvas>
			`,
		},
		"multiple_selected_disabled": {
			disabled: true,
			options:  []string{"Foo", "Bar"},
			selected: "Foo",
			want: `
				<canvas padded size="150x200">
					<content>
						<container pos="4,4" size="142x192">
							<widget pos="37,67" size="67x58" type="*widget.RadioGroup">
								<widget size="67x29" type="*widget.radioItem">
									<circle fillColor="background" size="28x28"/>
									<image pos="4,4" rsc="radioButtonCheckedIcon" size="iconInlineSize" themed="disabled"/>
									<text color="disabled" pos="32,0" size="67x29">Foo</text>
								</widget>
								<widget pos="0,29" size="67x29" type="*widget.radioItem">
									<circle fillColor="background" size="28x28"/>
									<image pos="4,4" rsc="radioButtonIcon" size="iconInlineSize" themed="disabled"/>
									<text color="disabled" pos="32,0" size="67x29">Bar</text>
								</widget>
							</widget>
						</container>
					</content>
				</canvas>
			`,
		},
		"multiple_selected_horizontal": {
			horizontal: true,
			options:    []string{"Foo", "Bar"},
			selected:   "Foo",
			want: `
				<canvas padded size="150x200">
					<content>
						<container pos="4,4" size="142x192">
							<widget pos="5,81" size="132x29" type="*widget.RadioGroup">
								<widget size="66x29" type="*widget.radioItem">
									<circle fillColor="background" size="28x28"/>
									<image pos="4,4" rsc="radioButtonCheckedIcon" size="iconInlineSize"/>
									<text pos="32,0" size="66x29">Foo</text>
								</widget>
								<widget pos="66,0" size="66x29" type="*widget.radioItem">
									<circle fillColor="background" size="28x28"/>
									<image pos="4,4" rsc="radioButtonIcon" size="iconInlineSize"/>
									<text pos="32,0" size="66x29">Bar</text>
								</widget>
							</widget>
						</container>
					</content>
				</canvas>
			`,
		},
		"multiple_selected_horizontal_disabled": {
			disabled:   true,
			horizontal: true,
			options:    []string{"Foo", "Bar"},
			selected:   "Foo",
			want: `
				<canvas padded size="150x200">
					<content>
						<container pos="4,4" size="142x192">
							<widget pos="5,81" size="132x29" type="*widget.RadioGroup">
								<widget size="66x29" type="*widget.radioItem">
									<circle fillColor="background" size="28x28"/>
									<image pos="4,4" rsc="radioButtonCheckedIcon" size="iconInlineSize" themed="disabled"/>
									<text color="disabled" pos="32,0" size="66x29">Foo</text>
								</widget>
								<widget pos="66,0" size="66x29" type="*widget.radioItem">
									<circle fillColor="background" size="28x28"/>
									<image pos="4,4" rsc="radioButtonIcon" size="iconInlineSize" themed="disabled"/>
									<text color="disabled" pos="32,0" size="66x29">Bar</text>
								</widget>
							</widget>
						</container>
					</content>
				</canvas>
			`,
		},
	} {
		t.Run(name, func(t *testing.T) {
			radio := &widget.RadioGroup{
				Horizontal: tt.horizontal,
				Options:    tt.options,
				Selected:   tt.selected,
			}
			if tt.disabled {
				radio.Disable()
			}

			window := test.NewWindow(fyne.NewContainerWithLayout(layout.NewCenterLayout(), radio))
			window.Resize(radio.MinSize().Max(fyne.NewSize(150, 200)))

			test.AssertRendersToMarkup(t, tt.want, window.Canvas())

			window.Close()
		})
	}
}

func TestRadioGroup_ToggleSelectionWithSpaceKey(t *testing.T) {
	radio := &widget.RadioGroup{Options: []string{"Option A", "Option B", "Option C"}}
	window := test.NewWindow(radio)
	defer window.Close()

	assert.Equal(t, "", radio.Selected)

	canvas := window.Canvas().(test.WindowlessCanvas)
	canvas.FocusNext()
	canvas.FocusNext()
	test.Type(canvas.Focused(), " ")
	assert.Equal(t, "Option B", radio.Selected)

	canvas.FocusNext()
	canvas.FocusNext()
	test.Type(canvas.Focused(), " ")
	assert.Equal(t, "Option A", radio.Selected)

	test.Type(canvas.Focused(), " ")
	assert.Equal(t, "", radio.Selected)

	canvas.FocusNext()
	canvas.FocusNext()
	test.Type(canvas.Focused(), " ")
	assert.Equal(t, "Option C", radio.Selected)

	radio.Required = true
	test.Type(canvas.Focused(), " ")
	assert.Equal(t, "Option C", radio.Selected, "cannot unselect required radio")
}
