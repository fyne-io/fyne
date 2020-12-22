package widget_test

import (
	"testing"

	"fyne.io/fyne"
	"fyne.io/fyne/test"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"

	"github.com/stretchr/testify/assert"
)

func TestSelectEntry_Disableable(t *testing.T) {
	test.NewApp()
	defer test.NewApp()

	options := []string{"A", "B", "C"}
	e := widget.NewSelectEntry(options)
	w := test.NewWindow(e)
	defer w.Close()
	w.Resize(fyne.NewSize(150, 200))
	e.Resize(e.MinSize().Max(fyne.NewSize(130, 0)))
	e.Move(fyne.NewPos(10, 10))
	c := w.Canvas()

	enabled := `
		<canvas padded size="150x200">
			<content>
				<widget pos="10,10" size="130x37" type="*widget.SelectEntry">
					<rectangle fillColor="shadow" pos="0,33" size="130x4"/>
					<widget pos="4,4" size="142x29" type="*widget.textProvider">
						<text color="placeholder" pos="4,4" size="134x21"></text>
					</widget>
					<widget pos="4,4" size="142x29" type="*widget.textProvider">
						<text pos="4,4" size="134x21"></text>
					</widget>
					<widget pos="102,8" size="20x20" type="*widget.Button">
						<rectangle fillColor="button" pos="2,2" size="16x16"/>
						<image fillMode="contain" rsc="menuDropDownIcon" size="iconInlineSize"/>
					</widget>
				</widget>
			</content>
		</canvas>
	`
	enabledOpened := `
		<canvas padded size="150x200">
			<content>
				<widget pos="10,10" size="130x37" type="*widget.SelectEntry">
					<rectangle fillColor="shadow" pos="0,33" size="130x4"/>
					<widget pos="4,4" size="142x29" type="*widget.textProvider">
						<text color="placeholder" pos="4,4" size="134x21"></text>
					</widget>
					<widget pos="4,4" size="142x29" type="*widget.textProvider">
						<text pos="4,4" size="134x21"></text>
					</widget>
					<widget pos="102,8" size="20x20" type="*widget.Button">
						<rectangle fillColor="button" pos="2,2" size="16x16"/>
						<image fillMode="contain" rsc="menuDropDownIcon" size="iconInlineSize"/>
					</widget>
				</widget>
			</content>
			<overlay>
				<widget size="150x200" type="*widget.OverlayContainer">
					<widget pos="14,47" size="122x103" type="*widget.Menu">
						<widget size="122x103" type="*widget.Shadow">
							<radialGradient centerOffset="0.5,0.5" pos="-4,-4" size="4x4" startColor="shadow"/>
							<linearGradient endColor="shadow" pos="0,-4" size="122x4"/>
							<radialGradient centerOffset="-0.5,0.5" pos="122,-4" size="4x4" startColor="shadow"/>
							<linearGradient angle="270" pos="122,0" size="4x103" startColor="shadow"/>
							<radialGradient centerOffset="-0.5,-0.5" pos="122,103" size="4x4" startColor="shadow"/>
							<linearGradient pos="0,103" size="122x4" startColor="shadow"/>
							<radialGradient centerOffset="0.5,-0.5" pos="-4,103" size="4x4" startColor="shadow"/>
							<linearGradient angle="270" endColor="shadow" pos="-4,0" size="4x103"/>
						</widget>
						<widget size="122x103" type="*widget.ScrollContainer">
							<widget size="122x103" type="*widget.menuBox">
								<container pos="0,4" size="122x111">
									<widget size="122x29" type="*widget.menuItem">
										<text pos="8,4" size="10x21">A</text>
									</widget>
									<widget pos="0,33" size="122x29" type="*widget.menuItem">
										<text pos="8,4" size="10x21">B</text>
									</widget>
									<widget pos="0,66" size="122x29" type="*widget.menuItem">
										<text pos="8,4" size="10x21">C</text>
									</widget>
								</container>
							</widget>
						</widget>
					</widget>
				</widget>
			</overlay>
		</canvas>
	`
	disabled := `
		<canvas padded size="150x200">
			<content>
				<widget pos="10,10" size="130x37" type="*widget.SelectEntry">
					<rectangle fillColor="disabled" pos="0,33" size="130x4"/>
					<widget pos="4,4" size="142x29" type="*widget.textProvider">
						<text color="disabled" pos="4,4" size="134x21"></text>
					</widget>
					<widget pos="4,4" size="142x29" type="*widget.textProvider">
						<text color="disabled" pos="4,4" size="134x21"></text>
					</widget>
					<widget pos="102,8" size="20x20" type="*widget.Button">
						<rectangle fillColor="disabled button" pos="2,2" size="16x16"/>
						<image fillMode="contain" rsc="menuDropDownIcon" size="iconInlineSize" themed="disabled"/>
					</widget>
				</widget>
			</content>
		</canvas>
	`

	assert.False(t, e.Disabled())
	test.AssertRendersToMarkup(t, enabled, c)

	switchPos := fyne.NewPos(140-theme.Padding()-theme.IconInlineSize()/2, 10+theme.Padding()+theme.IconInlineSize()/2)
	test.TapCanvas(c, switchPos)
	test.AssertRendersToMarkup(t, enabledOpened, c)

	test.TapCanvas(c, fyne.NewPos(0, 0))
	test.AssertRendersToMarkup(t, enabled, c)

	e.Disable()
	assert.True(t, e.Disabled())
	test.AssertRendersToMarkup(t, disabled, c)

	test.TapCanvas(c, switchPos)
	test.AssertRendersToMarkup(t, disabled, c, "no drop-down when disabled")

	e.Enable()
	assert.False(t, e.Disabled())
	test.AssertRendersToMarkup(t, enabled, c)

	test.TapCanvas(c, switchPos)
	test.AssertRendersToMarkup(t, enabledOpened, c)
}

func TestSelectEntry_DropDown(t *testing.T) {
	test.NewApp()
	defer test.NewApp()

	options := []string{"A", "B", "C"}
	e := widget.NewSelectEntry(options)
	w := test.NewWindow(e)
	defer w.Close()
	w.Resize(fyne.NewSize(150, 200))
	e.Resize(e.MinSize().Max(fyne.NewSize(130, 0)))
	e.Move(fyne.NewPos(10, 10))
	c := w.Canvas()

	test.AssertRendersToMarkup(t, `
		<canvas padded size="150x200">
			<content>
				<widget pos="10,10" size="130x37" type="*widget.SelectEntry">
					<rectangle fillColor="shadow" pos="0,33" size="130x4"/>
					<widget pos="4,4" size="142x29" type="*widget.textProvider">
						<text color="placeholder" pos="4,4" size="134x21"></text>
					</widget>
					<widget pos="4,4" size="142x29" type="*widget.textProvider">
						<text pos="4,4" size="134x21"></text>
					</widget>
					<widget pos="102,8" size="20x20" type="*widget.Button">
						<rectangle fillColor="button" pos="2,2" size="16x16"/>
						<image fillMode="contain" rsc="menuDropDownIcon" size="iconInlineSize"/>
					</widget>
				</widget>
			</content>
		</canvas>
	`, c)
	assert.Nil(t, c.Overlays().Top())

	switchPos := fyne.NewPos(140-theme.Padding()-theme.IconInlineSize()/2, 10+theme.Padding()+theme.IconInlineSize()/2)
	test.TapCanvas(c, switchPos)
	test.AssertRendersToMarkup(t, `
		<canvas padded size="150x200">
			<content>
				<widget pos="10,10" size="130x37" type="*widget.SelectEntry">
					<rectangle fillColor="shadow" pos="0,33" size="130x4"/>
					<widget pos="4,4" size="142x29" type="*widget.textProvider">
						<text color="placeholder" pos="4,4" size="134x21"></text>
					</widget>
					<widget pos="4,4" size="142x29" type="*widget.textProvider">
						<text pos="4,4" size="134x21"></text>
					</widget>
					<widget pos="102,8" size="20x20" type="*widget.Button">
						<rectangle fillColor="button" pos="2,2" size="16x16"/>
						<image fillMode="contain" rsc="menuDropDownIcon" size="iconInlineSize"/>
					</widget>
				</widget>
			</content>
			<overlay>
				<widget size="150x200" type="*widget.OverlayContainer">
					<widget pos="14,47" size="122x103" type="*widget.Menu">
						<widget size="122x103" type="*widget.Shadow">
							<radialGradient centerOffset="0.5,0.5" pos="-4,-4" size="4x4" startColor="shadow"/>
							<linearGradient endColor="shadow" pos="0,-4" size="122x4"/>
							<radialGradient centerOffset="-0.5,0.5" pos="122,-4" size="4x4" startColor="shadow"/>
							<linearGradient angle="270" pos="122,0" size="4x103" startColor="shadow"/>
							<radialGradient centerOffset="-0.5,-0.5" pos="122,103" size="4x4" startColor="shadow"/>
							<linearGradient pos="0,103" size="122x4" startColor="shadow"/>
							<radialGradient centerOffset="0.5,-0.5" pos="-4,103" size="4x4" startColor="shadow"/>
							<linearGradient angle="270" endColor="shadow" pos="-4,0" size="4x103"/>
						</widget>
						<widget size="122x103" type="*widget.ScrollContainer">
							<widget size="122x103" type="*widget.menuBox">
								<container pos="0,4" size="122x111">
									<widget size="122x29" type="*widget.menuItem">
										<text pos="8,4" size="10x21">A</text>
									</widget>
									<widget pos="0,33" size="122x29" type="*widget.menuItem">
										<text pos="8,4" size="10x21">B</text>
									</widget>
									<widget pos="0,66" size="122x29" type="*widget.menuItem">
										<text pos="8,4" size="10x21">C</text>
									</widget>
								</container>
							</widget>
						</widget>
					</widget>
				</widget>
			</overlay>
		</canvas>
	`, c)

	test.TapCanvas(c, fyne.NewPos(50, 15+2*(theme.Padding()+e.Size().Height)))
	test.AssertRendersToMarkup(t, `
		<canvas padded size="150x200">
			<content>
				<widget pos="10,10" size="130x37" type="*widget.SelectEntry">
					<rectangle fillColor="shadow" pos="0,33" size="130x4"/>
					<widget pos="4,4" size="142x29" type="*widget.textProvider">
						<text pos="4,4" size="134x21">B</text>
					</widget>
					<widget pos="102,8" size="20x20" type="*widget.Button">
						<rectangle fillColor="button" pos="2,2" size="16x16"/>
						<image fillMode="contain" rsc="menuDropDownIcon" size="iconInlineSize"/>
					</widget>
				</widget>
			</content>
		</canvas>
	`, c)
	assert.Equal(t, "B", e.Text)

	test.TapCanvas(c, switchPos)
	test.AssertRendersToMarkup(t, `
		<canvas padded size="150x200">
			<content>
				<widget pos="10,10" size="130x37" type="*widget.SelectEntry">
					<rectangle fillColor="shadow" pos="0,33" size="130x4"/>
					<widget pos="4,4" size="142x29" type="*widget.textProvider">
						<text pos="4,4" size="134x21">B</text>
					</widget>
					<widget pos="102,8" size="20x20" type="*widget.Button">
						<rectangle fillColor="button" pos="2,2" size="16x16"/>
						<image fillMode="contain" rsc="menuDropDownIcon" size="iconInlineSize"/>
					</widget>
				</widget>
			</content>
			<overlay>
				<widget size="150x200" type="*widget.OverlayContainer">
					<widget pos="14,47" size="122x103" type="*widget.Menu">
						<widget size="122x103" type="*widget.Shadow">
							<radialGradient centerOffset="0.5,0.5" pos="-4,-4" size="4x4" startColor="shadow"/>
							<linearGradient endColor="shadow" pos="0,-4" size="122x4"/>
							<radialGradient centerOffset="-0.5,0.5" pos="122,-4" size="4x4" startColor="shadow"/>
							<linearGradient angle="270" pos="122,0" size="4x103" startColor="shadow"/>
							<radialGradient centerOffset="-0.5,-0.5" pos="122,103" size="4x4" startColor="shadow"/>
							<linearGradient pos="0,103" size="122x4" startColor="shadow"/>
							<radialGradient centerOffset="0.5,-0.5" pos="-4,103" size="4x4" startColor="shadow"/>
							<linearGradient angle="270" endColor="shadow" pos="-4,0" size="4x103"/>
						</widget>
						<widget size="122x103" type="*widget.ScrollContainer">
							<widget size="122x103" type="*widget.menuBox">
								<container pos="0,4" size="122x111">
									<widget size="122x29" type="*widget.menuItem">
										<text pos="8,4" size="10x21">A</text>
									</widget>
									<widget pos="0,33" size="122x29" type="*widget.menuItem">
										<text pos="8,4" size="10x21">B</text>
									</widget>
									<widget pos="0,66" size="122x29" type="*widget.menuItem">
										<text pos="8,4" size="10x21">C</text>
									</widget>
								</container>
							</widget>
						</widget>
					</widget>
				</widget>
			</overlay>
		</canvas>
	`, c)

	test.TapCanvas(c, fyne.NewPos(50, 15+3*(theme.Padding()+e.Size().Height)))
	test.AssertRendersToMarkup(t, `
		<canvas padded size="150x200">
			<content>
				<widget pos="10,10" size="130x37" type="*widget.SelectEntry">
					<rectangle fillColor="shadow" pos="0,33" size="130x4"/>
					<widget pos="4,4" size="142x29" type="*widget.textProvider">
						<text pos="4,4" size="134x21">C</text>
					</widget>
					<widget pos="102,8" size="20x20" type="*widget.Button">
						<rectangle fillColor="button" pos="2,2" size="16x16"/>
						<image fillMode="contain" rsc="menuDropDownIcon" size="iconInlineSize"/>
					</widget>
				</widget>
			</content>
		</canvas>
	`, c)
	assert.Equal(t, "C", e.Text)
}

func TestSelectEntry_DropDownResize(t *testing.T) {
	test.NewApp()
	defer test.NewApp()

	options := []string{"A", "B", "C"}
	e := widget.NewSelectEntry(options)
	w := test.NewWindow(e)
	defer w.Close()
	w.Resize(fyne.NewSize(150, 200))
	e.Resize(e.MinSize().Max(fyne.NewSize(130, 0)))
	e.Move(fyne.NewPos(10, 10))
	c := w.Canvas()

	test.AssertRendersToMarkup(t, `
		<canvas padded size="150x200">
			<content>
				<widget pos="10,10" size="130x37" type="*widget.SelectEntry">
					<rectangle fillColor="shadow" pos="0,33" size="130x4"/>
					<widget pos="4,4" size="142x29" type="*widget.textProvider">
						<text color="placeholder" pos="4,4" size="134x21"></text>
					</widget>
					<widget pos="4,4" size="142x29" type="*widget.textProvider">
						<text pos="4,4" size="134x21"></text>
					</widget>
					<widget pos="102,8" size="20x20" type="*widget.Button">
						<rectangle fillColor="button" pos="2,2" size="16x16"/>
						<image fillMode="contain" rsc="menuDropDownIcon" size="iconInlineSize"/>
					</widget>
				</widget>
			</content>
		</canvas>
	`, c)
	assert.Nil(t, c.Overlays().Top())

	switchPos := fyne.NewPos(140-theme.Padding()-theme.IconInlineSize()/2, 10+theme.Padding()+theme.IconInlineSize()/2)
	test.TapCanvas(c, switchPos)
	opened := `
		<canvas padded size="150x200">
			<content>
				<widget pos="10,10" size="130x37" type="*widget.SelectEntry">
					<rectangle fillColor="shadow" pos="0,33" size="130x4"/>
					<widget pos="4,4" size="142x29" type="*widget.textProvider">
						<text color="placeholder" pos="4,4" size="134x21"></text>
					</widget>
					<widget pos="4,4" size="142x29" type="*widget.textProvider">
						<text pos="4,4" size="134x21"></text>
					</widget>
					<widget pos="102,8" size="20x20" type="*widget.Button">
						<rectangle fillColor="button" pos="2,2" size="16x16"/>
						<image fillMode="contain" rsc="menuDropDownIcon" size="iconInlineSize"/>
					</widget>
				</widget>
			</content>
			<overlay>
				<widget size="150x200" type="*widget.OverlayContainer">
					<widget pos="14,47" size="122x103" type="*widget.Menu">
						<widget size="122x103" type="*widget.Shadow">
							<radialGradient centerOffset="0.5,0.5" pos="-4,-4" size="4x4" startColor="shadow"/>
							<linearGradient endColor="shadow" pos="0,-4" size="122x4"/>
							<radialGradient centerOffset="-0.5,0.5" pos="122,-4" size="4x4" startColor="shadow"/>
							<linearGradient angle="270" pos="122,0" size="4x103" startColor="shadow"/>
							<radialGradient centerOffset="-0.5,-0.5" pos="122,103" size="4x4" startColor="shadow"/>
							<linearGradient pos="0,103" size="122x4" startColor="shadow"/>
							<radialGradient centerOffset="0.5,-0.5" pos="-4,103" size="4x4" startColor="shadow"/>
							<linearGradient angle="270" endColor="shadow" pos="-4,0" size="4x103"/>
						</widget>
						<widget size="122x103" type="*widget.ScrollContainer">
							<widget size="122x103" type="*widget.menuBox">
								<container pos="0,4" size="122x111">
									<widget size="122x29" type="*widget.menuItem">
										<text pos="8,4" size="10x21">A</text>
									</widget>
									<widget pos="0,33" size="122x29" type="*widget.menuItem">
										<text pos="8,4" size="10x21">B</text>
									</widget>
									<widget pos="0,66" size="122x29" type="*widget.menuItem">
										<text pos="8,4" size="10x21">C</text>
									</widget>
								</container>
							</widget>
						</widget>
					</widget>
				</widget>
			</overlay>
		</canvas>
	`
	test.AssertRendersToMarkup(t, opened, c)

	e.Resize(e.Size().Subtract(fyne.NewSize(20, 0)))
	test.AssertRendersToMarkup(t, `
		<canvas padded size="150x200">
			<content>
				<widget pos="10,10" size="110x37" type="*widget.SelectEntry">
					<rectangle fillColor="shadow" pos="0,33" size="110x4"/>
					<widget pos="4,4" size="122x29" type="*widget.textProvider">
						<text color="placeholder" pos="4,4" size="114x21"></text>
					</widget>
					<widget pos="4,4" size="122x29" type="*widget.textProvider">
						<text pos="4,4" size="114x21"></text>
					</widget>
					<widget pos="82,8" size="20x20" type="*widget.Button">
						<rectangle fillColor="button" pos="2,2" size="16x16"/>
						<image fillMode="contain" rsc="menuDropDownIcon" size="iconInlineSize"/>
					</widget>
				</widget>
			</content>
			<overlay>
				<widget size="150x200" type="*widget.OverlayContainer">
					<widget pos="14,47" size="102x103" type="*widget.Menu">
						<widget size="102x103" type="*widget.Shadow">
							<radialGradient centerOffset="0.5,0.5" pos="-4,-4" size="4x4" startColor="shadow"/>
							<linearGradient endColor="shadow" pos="0,-4" size="102x4"/>
							<radialGradient centerOffset="-0.5,0.5" pos="102,-4" size="4x4" startColor="shadow"/>
							<linearGradient angle="270" pos="102,0" size="4x103" startColor="shadow"/>
							<radialGradient centerOffset="-0.5,-0.5" pos="102,103" size="4x4" startColor="shadow"/>
							<linearGradient pos="0,103" size="102x4" startColor="shadow"/>
							<radialGradient centerOffset="0.5,-0.5" pos="-4,103" size="4x4" startColor="shadow"/>
							<linearGradient angle="270" endColor="shadow" pos="-4,0" size="4x103"/>
						</widget>
						<widget size="102x103" type="*widget.ScrollContainer">
							<widget size="102x103" type="*widget.menuBox">
								<container pos="0,4" size="102x111">
									<widget size="102x29" type="*widget.menuItem">
										<text pos="8,4" size="10x21">A</text>
									</widget>
									<widget pos="0,33" size="102x29" type="*widget.menuItem">
										<text pos="8,4" size="10x21">B</text>
									</widget>
									<widget pos="0,66" size="102x29" type="*widget.menuItem">
										<text pos="8,4" size="10x21">C</text>
									</widget>
								</container>
							</widget>
						</widget>
					</widget>
				</widget>
			</overlay>
		</canvas>
	`, c)

	e.Resize(e.Size().Add(fyne.NewSize(20, 0)))
	test.AssertRendersToMarkup(t, opened, c)
}

func TestSelectEntry_MinSize(t *testing.T) {
	smallOptions := []string{"A", "B", "C"}

	largeOptions := []string{"Large Option A", "Larger Option B", "Very Large Option C"}
	largeOptionsMinWidth := optionsMinSize(largeOptions).Width

	minTextHeight := widget.NewLabel("W").MinSize().Height

	tests := map[string]struct {
		placeholder string
		value       string
		options     []string
		want        fyne.Size
	}{
		"empty": {
			want: fyne.NewSize(emptyTextWidth()+dropDownIconWidth()+4*theme.Padding(), minTextHeight+2*theme.Padding()),
		},
		"empty + small options": {
			options: smallOptions,
			want:    fyne.NewSize(emptyTextWidth()+dropDownIconWidth()+4*theme.Padding(), minTextHeight+2*theme.Padding()),
		},
		"empty + large options": {
			options: largeOptions,
			want:    fyne.NewSize(largeOptionsMinWidth+2*theme.Padding(), minTextHeight+2*theme.Padding()),
		},
		"value": {
			value: "foo",
			want:  widget.NewLabel("foo").MinSize().Add(fyne.NewSize(dropDownIconWidth()+4*theme.Padding(), 2*theme.Padding())),
		},
		"large value + small options": {
			value:   "large",
			options: smallOptions,
			want:    widget.NewLabel("large").MinSize().Add(fyne.NewSize(dropDownIconWidth()+4*theme.Padding(), 2*theme.Padding())),
		},
		"small value + large options": {
			value:   "small",
			options: largeOptions,
			want:    fyne.NewSize(largeOptionsMinWidth+2*theme.Padding(), minTextHeight+2*theme.Padding()),
		},
		"placeholder": {
			placeholder: "example",
			want:        widget.NewLabel("example").MinSize().Add(fyne.NewSize(dropDownIconWidth()+4*theme.Padding(), 2*theme.Padding())),
		},
		"large placeholder + small options": {
			placeholder: "large",
			options:     smallOptions,
			want:        widget.NewLabel("large").MinSize().Add(fyne.NewSize(dropDownIconWidth()+4*theme.Padding(), 2*theme.Padding())),
		},
		"small placeholder + large options": {
			placeholder: "small",
			options:     largeOptions,
			want:        fyne.NewSize(largeOptionsMinWidth+2*theme.Padding(), minTextHeight+2*theme.Padding()),
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			e := widget.NewSelectEntry(tt.options)
			e.PlaceHolder = tt.placeholder
			e.Text = tt.value
			assert.Equal(t, tt.want, e.MinSize())
		})
	}
}

func TestSelectEntry_SetOptions(t *testing.T) {
	test.NewApp()
	defer test.NewApp()

	e := widget.NewSelectEntry([]string{"A", "B", "C"})
	w := test.NewWindow(e)
	defer w.Close()
	w.Resize(fyne.NewSize(150, 200))
	e.Resize(e.MinSize().Max(fyne.NewSize(130, 0)))
	e.Move(fyne.NewPos(10, 10))
	c := w.Canvas()

	switchPos := fyne.NewPos(140-theme.Padding()-theme.IconInlineSize()/2, 10+theme.Padding()+theme.IconInlineSize()/2)
	test.TapCanvas(c, switchPos)
	test.AssertRendersToMarkup(t, `
		<canvas padded size="150x200">
			<content>
				<widget pos="10,10" size="130x37" type="*widget.SelectEntry">
					<rectangle fillColor="shadow" pos="0,33" size="130x4"/>
					<widget pos="4,4" size="142x29" type="*widget.textProvider">
						<text color="placeholder" pos="4,4" size="134x21"></text>
					</widget>
					<widget pos="4,4" size="142x29" type="*widget.textProvider">
						<text pos="4,4" size="134x21"></text>
					</widget>
					<widget pos="102,8" size="20x20" type="*widget.Button">
						<rectangle fillColor="button" pos="2,2" size="16x16"/>
						<image fillMode="contain" rsc="menuDropDownIcon" size="iconInlineSize"/>
					</widget>
				</widget>
			</content>
			<overlay>
				<widget size="150x200" type="*widget.OverlayContainer">
					<widget pos="14,47" size="122x103" type="*widget.Menu">
						<widget size="122x103" type="*widget.Shadow">
							<radialGradient centerOffset="0.5,0.5" pos="-4,-4" size="4x4" startColor="shadow"/>
							<linearGradient endColor="shadow" pos="0,-4" size="122x4"/>
							<radialGradient centerOffset="-0.5,0.5" pos="122,-4" size="4x4" startColor="shadow"/>
							<linearGradient angle="270" pos="122,0" size="4x103" startColor="shadow"/>
							<radialGradient centerOffset="-0.5,-0.5" pos="122,103" size="4x4" startColor="shadow"/>
							<linearGradient pos="0,103" size="122x4" startColor="shadow"/>
							<radialGradient centerOffset="0.5,-0.5" pos="-4,103" size="4x4" startColor="shadow"/>
							<linearGradient angle="270" endColor="shadow" pos="-4,0" size="4x103"/>
						</widget>
						<widget size="122x103" type="*widget.ScrollContainer">
							<widget size="122x103" type="*widget.menuBox">
								<container pos="0,4" size="122x111">
									<widget size="122x29" type="*widget.menuItem">
										<text pos="8,4" size="10x21">A</text>
									</widget>
									<widget pos="0,33" size="122x29" type="*widget.menuItem">
										<text pos="8,4" size="10x21">B</text>
									</widget>
									<widget pos="0,66" size="122x29" type="*widget.menuItem">
										<text pos="8,4" size="10x21">C</text>
									</widget>
								</container>
							</widget>
						</widget>
					</widget>
				</widget>
			</overlay>
		</canvas>
	`, c)
	test.TapCanvas(c, switchPos)

	e.SetOptions([]string{"1", "2", "3"})
	test.TapCanvas(c, switchPos)
	test.AssertRendersToMarkup(t, `
		<canvas padded size="150x200">
			<content>
				<widget pos="10,10" size="130x37" type="*widget.SelectEntry">
					<rectangle fillColor="shadow" pos="0,33" size="130x4"/>
					<widget pos="4,4" size="142x29" type="*widget.textProvider">
						<text color="placeholder" pos="4,4" size="134x21"></text>
					</widget>
					<widget pos="4,4" size="142x29" type="*widget.textProvider">
						<text pos="4,4" size="134x21"></text>
					</widget>
					<widget pos="102,8" size="20x20" type="*widget.Button">
						<rectangle fillColor="button" pos="2,2" size="16x16"/>
						<image fillMode="contain" rsc="menuDropDownIcon" size="iconInlineSize"/>
					</widget>
				</widget>
			</content>
			<overlay>
				<widget size="150x200" type="*widget.OverlayContainer">
					<widget pos="14,47" size="122x103" type="*widget.Menu">
						<widget size="122x103" type="*widget.Shadow">
							<radialGradient centerOffset="0.5,0.5" pos="-4,-4" size="4x4" startColor="shadow"/>
							<linearGradient endColor="shadow" pos="0,-4" size="122x4"/>
							<radialGradient centerOffset="-0.5,0.5" pos="122,-4" size="4x4" startColor="shadow"/>
							<linearGradient angle="270" pos="122,0" size="4x103" startColor="shadow"/>
							<radialGradient centerOffset="-0.5,-0.5" pos="122,103" size="4x4" startColor="shadow"/>
							<linearGradient pos="0,103" size="122x4" startColor="shadow"/>
							<radialGradient centerOffset="0.5,-0.5" pos="-4,103" size="4x4" startColor="shadow"/>
							<linearGradient angle="270" endColor="shadow" pos="-4,0" size="4x103"/>
						</widget>
						<widget size="122x103" type="*widget.ScrollContainer">
							<widget size="122x103" type="*widget.menuBox">
								<container pos="0,4" size="122x111">
									<widget size="122x29" type="*widget.menuItem">
										<text pos="8,4" size="9x21">1</text>
									</widget>
									<widget pos="0,33" size="122x29" type="*widget.menuItem">
										<text pos="8,4" size="9x21">2</text>
									</widget>
									<widget pos="0,66" size="122x29" type="*widget.menuItem">
										<text pos="8,4" size="9x21">3</text>
									</widget>
								</container>
							</widget>
						</widget>
					</widget>
				</widget>
			</overlay>
		</canvas>
	`, c)
}

func TestSelectEntry_SetOptions_Empty(t *testing.T) {
	test.NewApp()
	defer test.NewApp()

	e := widget.NewSelectEntry([]string{})
	w := test.NewWindow(e)
	defer w.Close()
	w.Resize(fyne.NewSize(150, 200))
	e.Resize(e.MinSize().Max(fyne.NewSize(130, 0)))
	e.Move(fyne.NewPos(10, 10))
	c := w.Canvas()

	switchPos := fyne.NewPos(140-theme.Padding()-theme.IconInlineSize()/2, 10+theme.Padding()+theme.IconInlineSize()/2)
	e.SetOptions([]string{"1", "2", "3"})
	test.TapCanvas(c, switchPos)
	test.AssertRendersToMarkup(t, `
		<canvas padded size="150x200">
			<content>
				<widget pos="10,10" size="130x37" type="*widget.SelectEntry">
					<rectangle fillColor="shadow" pos="0,33" size="130x4"/>
					<widget pos="4,4" size="142x29" type="*widget.textProvider">
						<text color="placeholder" pos="4,4" size="134x21"></text>
					</widget>
					<widget pos="4,4" size="142x29" type="*widget.textProvider">
						<text pos="4,4" size="134x21"></text>
					</widget>
					<widget pos="102,8" size="20x20" type="*widget.Button">
						<rectangle fillColor="button" pos="2,2" size="16x16"/>
						<image fillMode="contain" rsc="menuDropDownIcon" size="iconInlineSize"/>
					</widget>
				</widget>
			</content>
			<overlay>
				<widget size="150x200" type="*widget.OverlayContainer">
					<widget pos="14,47" size="122x103" type="*widget.Menu">
						<widget size="122x103" type="*widget.Shadow">
							<radialGradient centerOffset="0.5,0.5" pos="-4,-4" size="4x4" startColor="shadow"/>
							<linearGradient endColor="shadow" pos="0,-4" size="122x4"/>
							<radialGradient centerOffset="-0.5,0.5" pos="122,-4" size="4x4" startColor="shadow"/>
							<linearGradient angle="270" pos="122,0" size="4x103" startColor="shadow"/>
							<radialGradient centerOffset="-0.5,-0.5" pos="122,103" size="4x4" startColor="shadow"/>
							<linearGradient pos="0,103" size="122x4" startColor="shadow"/>
							<radialGradient centerOffset="0.5,-0.5" pos="-4,103" size="4x4" startColor="shadow"/>
							<linearGradient angle="270" endColor="shadow" pos="-4,0" size="4x103"/>
						</widget>
						<widget size="122x103" type="*widget.ScrollContainer">
							<widget size="122x103" type="*widget.menuBox">
								<container pos="0,4" size="122x111">
									<widget size="122x29" type="*widget.menuItem">
										<text pos="8,4" size="9x21">1</text>
									</widget>
									<widget pos="0,33" size="122x29" type="*widget.menuItem">
										<text pos="8,4" size="9x21">2</text>
									</widget>
									<widget pos="0,66" size="122x29" type="*widget.menuItem">
										<text pos="8,4" size="9x21">3</text>
									</widget>
								</container>
							</widget>
						</widget>
					</widget>
				</widget>
			</overlay>
		</canvas>
	`, c)
}

func dropDownIconWidth() float32 {
	return theme.IconInlineSize() + theme.Padding()
}

func emptyTextWidth() float32 {
	return widget.NewLabel("M").MinSize().Width
}

func optionsMinSize(options []string) fyne.Size {
	var labels []*widget.Label
	for _, option := range options {
		labels = append(labels, widget.NewLabel(option))
	}
	minWidth := float32(0)
	minHeight := float32(0)
	for _, label := range labels {
		if minWidth < label.MinSize().Width {
			minWidth = label.MinSize().Width
		}
		minHeight += label.MinSize().Height
	}
	// padding between all options
	minHeight += float32(len(labels)-1) * theme.Padding()
	return fyne.NewSize(minWidth, minHeight)
}
