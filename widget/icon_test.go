package widget_test

import (
	"testing"

	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/test"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
)

func TestIcon_Layout(t *testing.T) {
	test.NewApp()
	defer test.NewApp()
	test.ApplyTheme(t, theme.LightTheme())

	for name, tt := range map[string]struct {
		resource fyne.Resource
		want     string
	}{
		"empty": {
			want: `
				<canvas padded size="150x200">
					<content>
						<container pos="4,4" size="142x192">
							<widget pos="61,86" size="20x20" type="*widget.Icon">
							</widget>
						</container>
					</content>
				</canvas>
			`,
		},
		"resource": {
			resource: theme.CancelIcon(),
			want: `
				<canvas padded size="150x200">
					<content>
						<container pos="4,4" size="142x192">
							<widget pos="61,86" size="20x20" type="*widget.Icon">
								<image fillMode="contain" rsc="cancelIcon" size="iconInlineSize"/>
							</widget>
						</container>
					</content>
				</canvas>
			`,
		},
	} {
		t.Run(name, func(t *testing.T) {
			icon := &widget.Icon{
				Resource: tt.resource,
			}

			window := test.NewWindow(fyne.NewContainerWithLayout(layout.NewCenterLayout(), icon))
			window.Resize(icon.MinSize().Max(fyne.NewSize(150, 200)))

			test.AssertRendersToMarkup(t, tt.want, window.Canvas())

			window.Close()
		})
	}
}
