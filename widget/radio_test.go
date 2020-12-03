package widget_test

import (
	"testing"

	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/test"
	"fyne.io/fyne/widget"
)

func TestRadio_Layout(t *testing.T) {
	test.NewApp()
	defer test.NewApp()

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
							<widget pos="36,81" size="70x29" type="*widget.Radio">
								<circle fillColor="background" size="28x28"/>
								<image pos="4,4" rsc="radioButtonIcon" size="iconInlineSize"/>
								<text pos="32,0" size="70x29">Test</text>
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
							<widget pos="36,81" size="70x29" type="*widget.Radio">
								<circle fillColor="background" size="28x28"/>
								<image pos="4,4" rsc="radioButtonIcon" size="iconInlineSize" themed="disabled"/>
								<text color="disabled" pos="32,0" size="70x29">Test</text>
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
							<widget pos="36,81" size="70x29" type="*widget.Radio">
								<circle fillColor="background" size="28x28"/>
								<image pos="4,4" rsc="radioButtonIcon" size="iconInlineSize"/>
								<text pos="32,0" size="70x29">Test</text>
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
							<widget pos="36,81" size="70x29" type="*widget.Radio">
								<circle fillColor="background" size="28x28"/>
								<image pos="4,4" rsc="radioButtonIcon" size="iconInlineSize" themed="disabled"/>
								<text color="disabled" pos="32,0" size="70x29">Test</text>
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
							<widget pos="36,81" size="70x29" type="*widget.Radio">
								<circle fillColor="background" size="28x28"/>
								<image pos="4,4" rsc="radioButtonCheckedIcon" size="iconInlineSize"/>
								<text pos="32,0" size="70x29">Test</text>
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
							<widget pos="36,81" size="70x29" type="*widget.Radio">
								<circle fillColor="background" size="28x28"/>
								<image pos="4,4" rsc="radioButtonCheckedIcon" size="iconInlineSize" themed="disabled"/>
								<text color="disabled" pos="32,0" size="70x29">Test</text>
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
							<widget pos="36,81" size="70x29" type="*widget.Radio">
								<circle fillColor="background" size="28x28"/>
								<image pos="4,4" rsc="radioButtonCheckedIcon" size="iconInlineSize"/>
								<text pos="32,0" size="70x29">Test</text>
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
							<widget pos="36,81" size="70x29" type="*widget.Radio">
								<circle fillColor="background" size="28x28"/>
								<image pos="4,4" rsc="radioButtonCheckedIcon" size="iconInlineSize" themed="disabled"/>
								<text color="disabled" pos="32,0" size="70x29">Test</text>
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
							<widget pos="37,67" size="67x58" type="*widget.Radio">
								<circle fillColor="background" size="28x28"/>
								<image pos="4,4" rsc="radioButtonIcon" size="iconInlineSize"/>
								<text pos="32,0" size="67x29">Foo</text>
								<circle fillColor="background" pos="0,29" size="28x28"/>
								<image pos="4,33" rsc="radioButtonIcon" size="iconInlineSize"/>
								<text pos="32,29" size="67x29">Bar</text>
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
							<widget pos="37,67" size="67x58" type="*widget.Radio">
								<circle fillColor="background" size="28x28"/>
								<image pos="4,4" rsc="radioButtonIcon" size="iconInlineSize" themed="disabled"/>
								<text color="disabled" pos="32,0" size="67x29">Foo</text>
								<circle fillColor="background" pos="0,29" size="28x28"/>
								<image pos="4,33" rsc="radioButtonIcon" size="iconInlineSize" themed="disabled"/>
								<text color="disabled" pos="32,29" size="67x29">Bar</text>
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
							<widget pos="5,81" size="132x29" type="*widget.Radio">
								<circle fillColor="background" size="28x28"/>
								<image pos="4,4" rsc="radioButtonIcon" size="iconInlineSize"/>
								<text pos="32,0" size="66x29">Foo</text>
								<circle fillColor="background" pos="66,0" size="28x28"/>
								<image pos="70,4" rsc="radioButtonIcon" size="iconInlineSize"/>
								<text pos="98,0" size="66x29">Bar</text>
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
							<widget pos="5,81" size="132x29" type="*widget.Radio">
								<circle fillColor="background" size="28x28"/>
								<image pos="4,4" rsc="radioButtonIcon" size="iconInlineSize" themed="disabled"/>
								<text color="disabled" pos="32,0" size="66x29">Foo</text>
								<circle fillColor="background" pos="66,0" size="28x28"/>
								<image pos="70,4" rsc="radioButtonIcon" size="iconInlineSize" themed="disabled"/>
								<text color="disabled" pos="98,0" size="66x29">Bar</text>
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
							<widget pos="37,67" size="67x58" type="*widget.Radio">
								<circle fillColor="background" size="28x28"/>
								<image pos="4,4" rsc="radioButtonCheckedIcon" size="iconInlineSize"/>
								<text pos="32,0" size="67x29">Foo</text>
								<circle fillColor="background" pos="0,29" size="28x28"/>
								<image pos="4,33" rsc="radioButtonIcon" size="iconInlineSize"/>
								<text pos="32,29" size="67x29">Bar</text>
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
							<widget pos="37,67" size="67x58" type="*widget.Radio">
								<circle fillColor="background" size="28x28"/>
								<image pos="4,4" rsc="radioButtonCheckedIcon" size="iconInlineSize" themed="disabled"/>
								<text color="disabled" pos="32,0" size="67x29">Foo</text>
								<circle fillColor="background" pos="0,29" size="28x28"/>
								<image pos="4,33" rsc="radioButtonIcon" size="iconInlineSize" themed="disabled"/>
								<text color="disabled" pos="32,29" size="67x29">Bar</text>
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
							<widget pos="5,81" size="132x29" type="*widget.Radio">
								<circle fillColor="background" size="28x28"/>
								<image pos="4,4" rsc="radioButtonCheckedIcon" size="iconInlineSize"/>
								<text pos="32,0" size="66x29">Foo</text>
								<circle fillColor="background" pos="66,0" size="28x28"/>
								<image pos="70,4" rsc="radioButtonIcon" size="iconInlineSize"/>
								<text pos="98,0" size="66x29">Bar</text>
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
							<widget pos="5,81" size="132x29" type="*widget.Radio">
								<circle fillColor="background" size="28x28"/>
								<image pos="4,4" rsc="radioButtonCheckedIcon" size="iconInlineSize" themed="disabled"/>
								<text color="disabled" pos="32,0" size="66x29">Foo</text>
								<circle fillColor="background" pos="66,0" size="28x28"/>
								<image pos="70,4" rsc="radioButtonIcon" size="iconInlineSize" themed="disabled"/>
								<text color="disabled" pos="98,0" size="66x29">Bar</text>
							</widget>
						</container>
					</content>
				</canvas>
			`,
		},
	} {
		t.Run(name, func(t *testing.T) {
			radio := &widget.Radio{
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
