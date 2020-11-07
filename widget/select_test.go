package widget_test

import (
	"testing"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/test"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewSelect(t *testing.T) {
	combo := widget.NewSelect([]string{"1", "2"}, func(string) {})

	assert.Equal(t, 2, len(combo.Options))
	assert.Equal(t, "", combo.Selected)
}

func TestSelect_ChangeTheme(t *testing.T) {
	app := test.NewApp()
	defer test.NewApp()
	app.Settings().SetTheme(theme.LightTheme())

	combo := widget.NewSelect([]string{"1", "2"}, func(s string) {})
	w := test.NewWindow(combo)
	defer w.Close()
	w.Resize(fyne.NewSize(220, 220))
	combo.Resize(combo.MinSize())
	combo.Move(fyne.NewPos(10, 10))
	test.Tap(combo)

	test.AssertImageMatches(t, "select/theme_initial.png", w.Canvas().Capture())

	test.WithTestTheme(t, func() {
		combo.Resize(combo.MinSize())
		combo.Refresh()
		time.Sleep(100 * time.Millisecond)
		test.AssertImageMatches(t, "select/theme_changed.png", w.Canvas().Capture())
	})
}

func TestSelect_ClearSelected(t *testing.T) {
	const (
		opt1     = "1"
		opt2     = "2"
		optClear = ""
	)
	combo := widget.NewSelect([]string{opt1, opt2}, func(string) {})
	assert.NotEmpty(t, combo.PlaceHolder)

	combo.SetSelected(opt1)
	assert.Equal(t, opt1, combo.Selected)

	var triggered bool
	var triggeredValue string
	combo.OnChanged = func(s string) {
		triggered = true
		triggeredValue = s
	}
	combo.ClearSelected()
	assert.Equal(t, optClear, combo.Selected)
	assert.True(t, triggered)
	assert.Equal(t, optClear, triggeredValue)
}

func TestSelect_Disable(t *testing.T) {
	app := test.NewApp()
	defer test.NewApp()
	app.Settings().SetTheme(theme.LightTheme())

	sel := widget.NewSelect([]string{"Hi"}, func(string) {})
	w := test.NewWindow(fyne.NewContainerWithLayout(layout.NewCenterLayout(), sel))
	defer w.Close()
	w.Resize(fyne.NewSize(200, 150))
	c := fyne.CurrentApp().Driver().CanvasForObject(sel)

	sel.Disable()
	disabled := `
		<canvas padded size="200x150">
			<content>
				<container pos="4,4" size="192x142">
					<widget pos="28,52" size="136x37" type="*widget.Select">
						<widget pos="4,4" size="128x29" type="*widget.Shadow">
							<radialGradient centerOffset="0.5,0.5" pos="-2,-2" size="2x2" startColor="shadow"/>
							<linearGradient endColor="shadow" pos="0,-2" size="128x2"/>
							<radialGradient centerOffset="-0.5,0.5" pos="128,-2" size="2x2" startColor="shadow"/>
							<linearGradient angle="270" pos="128,0" size="2x29" startColor="shadow"/>
							<radialGradient centerOffset="-0.5,-0.5" pos="128,29" size="2x2" startColor="shadow"/>
							<linearGradient pos="0,29" size="128x2" startColor="shadow"/>
							<radialGradient centerOffset="0.5,-0.5" pos="-2,29" size="2x2" startColor="shadow"/>
							<linearGradient angle="270" endColor="shadow" pos="-2,0" size="2x29"/>
						</widget>
						<rectangle pos="4,4" size="128x29"/>
						<widget pos="8,4" size="100x29" type="*widget.textProvider">
							<text bold color="disabled text" pos="4,4" size="92x21">(Select one)</text>
						</widget>
						<widget pos="108,8" size="20x20" type="*widget.Icon">
							<image fillMode="contain" rsc="menuDropDownIcon" size="iconInlineSize" themed="disabled"/>
						</widget>
					</widget>
				</container>
			</content>
		</canvas>
	`
	test.AssertRendersToMarkup(t, disabled, c)
	test.Tap(sel)
	assert.Nil(t, c.Overlays().Top(), "no pop-up for disabled Select")
	test.AssertRendersToMarkup(t, disabled, c)
}

func TestSelect_Disabled(t *testing.T) {
	sel := widget.NewSelect([]string{"Hi"}, func(string) {})
	assert.False(t, sel.Disabled())
	sel.Disable()
	assert.True(t, sel.Disabled())
	sel.Enable()
	assert.False(t, sel.Disabled())
}

func TestSelect_Enable(t *testing.T) {
	selected := ""
	sel := widget.NewSelect([]string{"Hi"}, func(sel string) {
		selected = sel
	})
	sel.Disable()
	require.True(t, sel.Disabled())

	sel.Enable()
	test.Tap(sel)
	c := fyne.CurrentApp().Driver().CanvasForObject(sel)
	ovl := c.Overlays().Top()
	if assert.NotNil(t, ovl, "pop-up for enabled Select") {
		test.TapCanvas(c, ovl.Position().Add(fyne.NewPos(theme.Padding()*2, theme.Padding()*2)))
		assert.Equal(t, "Hi", selected, "Radio should have been re-enabled.")
	}
}

func TestSelect_FocusRendering(t *testing.T) {
	test.NewApp()
	defer test.NewApp()
	test.ApplyTheme(t, theme.LightTheme())

	focusedNoneSelected := `
			<canvas padded size="200x150">
				<content>
					<container pos="4,4" size="192x142">
						<widget pos="28,52" size="136x37" type="*widget.Select">
							<widget pos="4,4" size="128x29" type="*widget.Shadow">
								<radialGradient centerOffset="0.5,0.5" pos="-2,-2" size="2x2" startColor="shadow"/>
								<linearGradient endColor="shadow" pos="0,-2" size="128x2"/>
								<radialGradient centerOffset="-0.5,0.5" pos="128,-2" size="2x2" startColor="shadow"/>
								<linearGradient angle="270" pos="128,0" size="2x29" startColor="shadow"/>
								<radialGradient centerOffset="-0.5,-0.5" pos="128,29" size="2x2" startColor="shadow"/>
								<linearGradient pos="0,29" size="128x2" startColor="shadow"/>
								<radialGradient centerOffset="0.5,-0.5" pos="-2,29" size="2x2" startColor="shadow"/>
								<linearGradient angle="270" endColor="shadow" pos="-2,0" size="2x29"/>
							</widget>
							<rectangle fillColor="focus" pos="4,4" size="128x29"/>
							<widget pos="8,4" size="100x29" type="*widget.textProvider">
								<text bold pos="4,4" size="92x21">(Select one)</text>
							</widget>
							<widget pos="108,8" size="20x20" type="*widget.Icon">
								<image fillMode="contain" rsc="menuDropDownIcon" size="iconInlineSize"/>
							</widget>
						</widget>
					</container>
				</content>
			</canvas>
		`
	t.Run("gain/lose focus", func(t *testing.T) {
		sel := widget.NewSelect([]string{"Option A", "Option B", "Option C"}, nil)
		w := test.NewWindow(fyne.NewContainerWithLayout(layout.NewCenterLayout(), sel))
		defer w.Close()
		w.Resize(fyne.NewSize(200, 150))

		c := w.Canvas()
		unfocusedNoneSelected := `
				<canvas padded size="200x150">
					<content>
						<container pos="4,4" size="192x142">
							<widget pos="28,52" size="136x37" type="*widget.Select">
								<widget pos="4,4" size="128x29" type="*widget.Shadow">
									<radialGradient centerOffset="0.5,0.5" pos="-2,-2" size="2x2" startColor="shadow"/>
									<linearGradient endColor="shadow" pos="0,-2" size="128x2"/>
									<radialGradient centerOffset="-0.5,0.5" pos="128,-2" size="2x2" startColor="shadow"/>
									<linearGradient angle="270" pos="128,0" size="2x29" startColor="shadow"/>
									<radialGradient centerOffset="-0.5,-0.5" pos="128,29" size="2x2" startColor="shadow"/>
									<linearGradient pos="0,29" size="128x2" startColor="shadow"/>
									<radialGradient centerOffset="0.5,-0.5" pos="-2,29" size="2x2" startColor="shadow"/>
									<linearGradient angle="270" endColor="shadow" pos="-2,0" size="2x29"/>
								</widget>
								<rectangle pos="4,4" size="128x29"/>
								<widget pos="8,4" size="100x29" type="*widget.textProvider">
									<text bold pos="4,4" size="92x21">(Select one)</text>
								</widget>
								<widget pos="108,8" size="20x20" type="*widget.Icon">
									<image fillMode="contain" rsc="menuDropDownIcon" size="iconInlineSize"/>
								</widget>
							</widget>
						</container>
					</content>
				</canvas>
			`
		test.AssertRendersToMarkup(t, unfocusedNoneSelected, c)
		c.Focus(sel)
		test.AssertRendersToMarkup(t, focusedNoneSelected, c)
		c.Unfocus()
		test.AssertRendersToMarkup(t, unfocusedNoneSelected, c)

		sel.SetSelected("Option B")
		assert.Equal(t, "Option B", sel.Selected)
		unfocusedBSelected := `
			<canvas padded size="200x150">
				<content>
					<container pos="4,4" size="192x142">
						<widget pos="28,52" size="136x37" type="*widget.Select">
							<widget pos="4,4" size="128x29" type="*widget.Shadow">
								<radialGradient centerOffset="0.5,0.5" pos="-2,-2" size="2x2" startColor="shadow"/>
								<linearGradient endColor="shadow" pos="0,-2" size="128x2"/>
								<radialGradient centerOffset="-0.5,0.5" pos="128,-2" size="2x2" startColor="shadow"/>
								<linearGradient angle="270" pos="128,0" size="2x29" startColor="shadow"/>
								<radialGradient centerOffset="-0.5,-0.5" pos="128,29" size="2x2" startColor="shadow"/>
								<linearGradient pos="0,29" size="128x2" startColor="shadow"/>
								<radialGradient centerOffset="0.5,-0.5" pos="-2,29" size="2x2" startColor="shadow"/>
								<linearGradient angle="270" endColor="shadow" pos="-2,0" size="2x29"/>
							</widget>
							<rectangle pos="4,4" size="128x29"/>
							<widget pos="8,4" size="100x29" type="*widget.textProvider">
								<text bold pos="4,4" size="92x21">Option B</text>
							</widget>
							<widget pos="108,8" size="20x20" type="*widget.Icon">
								<image fillMode="contain" rsc="menuDropDownIcon" size="iconInlineSize"/>
							</widget>
						</widget>
					</container>
				</content>
			</canvas>
		`
		test.AssertRendersToMarkup(t, unfocusedBSelected, c)
		c.Focus(sel)
		test.AssertRendersToMarkup(t, `
			<canvas padded size="200x150">
				<content>
					<container pos="4,4" size="192x142">
						<widget pos="28,52" size="136x37" type="*widget.Select">
							<widget pos="4,4" size="128x29" type="*widget.Shadow">
								<radialGradient centerOffset="0.5,0.5" pos="-2,-2" size="2x2" startColor="shadow"/>
								<linearGradient endColor="shadow" pos="0,-2" size="128x2"/>
								<radialGradient centerOffset="-0.5,0.5" pos="128,-2" size="2x2" startColor="shadow"/>
								<linearGradient angle="270" pos="128,0" size="2x29" startColor="shadow"/>
								<radialGradient centerOffset="-0.5,-0.5" pos="128,29" size="2x2" startColor="shadow"/>
								<linearGradient pos="0,29" size="128x2" startColor="shadow"/>
								<radialGradient centerOffset="0.5,-0.5" pos="-2,29" size="2x2" startColor="shadow"/>
								<linearGradient angle="270" endColor="shadow" pos="-2,0" size="2x29"/>
							</widget>
							<rectangle fillColor="focus" pos="4,4" size="128x29"/>
							<widget pos="8,4" size="100x29" type="*widget.textProvider">
								<text bold pos="4,4" size="92x21">Option B</text>
							</widget>
							<widget pos="108,8" size="20x20" type="*widget.Icon">
								<image fillMode="contain" rsc="menuDropDownIcon" size="iconInlineSize"/>
							</widget>
						</widget>
					</container>
				</content>
			</canvas>
		`, c)
		c.Unfocus()
		test.AssertRendersToMarkup(t, unfocusedBSelected, c)
	})
	t.Run("disable/enable focused", func(t *testing.T) {
		sel := widget.NewSelect([]string{"Option A", "Option B", "Option C"}, nil)
		w := test.NewWindow(fyne.NewContainerWithLayout(layout.NewCenterLayout(), sel))
		defer w.Close()
		w.Resize(fyne.NewSize(200, 150))

		c := w.Canvas().(test.WindowlessCanvas)
		c.FocusNext()
		test.AssertRendersToMarkup(t, focusedNoneSelected, c)
		sel.Disable()
		test.AssertRendersToMarkup(t, `
			<canvas padded size="200x150">
				<content>
					<container pos="4,4" size="192x142">
						<widget pos="28,52" size="136x37" type="*widget.Select">
							<widget pos="4,4" size="128x29" type="*widget.Shadow">
								<radialGradient centerOffset="0.5,0.5" pos="-2,-2" size="2x2" startColor="shadow"/>
								<linearGradient endColor="shadow" pos="0,-2" size="128x2"/>
								<radialGradient centerOffset="-0.5,0.5" pos="128,-2" size="2x2" startColor="shadow"/>
								<linearGradient angle="270" pos="128,0" size="2x29" startColor="shadow"/>
								<radialGradient centerOffset="-0.5,-0.5" pos="128,29" size="2x2" startColor="shadow"/>
								<linearGradient pos="0,29" size="128x2" startColor="shadow"/>
								<radialGradient centerOffset="0.5,-0.5" pos="-2,29" size="2x2" startColor="shadow"/>
								<linearGradient angle="270" endColor="shadow" pos="-2,0" size="2x29"/>
							</widget>
							<rectangle pos="4,4" size="128x29"/>
							<widget pos="8,4" size="100x29" type="*widget.textProvider">
								<text bold color="disabled text" pos="4,4" size="92x21">(Select one)</text>
							</widget>
							<widget pos="108,8" size="20x20" type="*widget.Icon">
								<image fillMode="contain" rsc="menuDropDownIcon" size="iconInlineSize" themed="disabled"/>
							</widget>
						</widget>
					</container>
				</content>
			</canvas>
		`, c)
		sel.Enable()
		test.AssertRendersToMarkup(t, focusedNoneSelected, c)
	})
}

func TestSelect_KeyboardControl(t *testing.T) {
	test.NewApp()
	defer test.NewApp()
	test.ApplyTheme(t, theme.LightTheme())

	noneSelected := `
		<canvas padded size="150x200">
			<content>
				<container pos="4,4" size="142x192">
					<widget pos="3,77" size="136x37" type="*widget.Select">
						<widget pos="4,4" size="128x29" type="*widget.Shadow">
							<radialGradient centerOffset="0.5,0.5" pos="-2,-2" size="2x2" startColor="shadow"/>
							<linearGradient endColor="shadow" pos="0,-2" size="128x2"/>
							<radialGradient centerOffset="-0.5,0.5" pos="128,-2" size="2x2" startColor="shadow"/>
							<linearGradient angle="270" pos="128,0" size="2x29" startColor="shadow"/>
							<radialGradient centerOffset="-0.5,-0.5" pos="128,29" size="2x2" startColor="shadow"/>
							<linearGradient pos="0,29" size="128x2" startColor="shadow"/>
							<radialGradient centerOffset="0.5,-0.5" pos="-2,29" size="2x2" startColor="shadow"/>
							<linearGradient angle="270" endColor="shadow" pos="-2,0" size="2x29"/>
						</widget>
						<rectangle fillColor="focus" pos="4,4" size="128x29"/>
						<widget pos="8,4" size="100x29" type="*widget.textProvider">
							<text bold pos="4,4" size="92x21">(Select one)</text>
						</widget>
						<widget pos="108,8" size="20x20" type="*widget.Icon">
							<image fillMode="contain" rsc="menuDropDownIcon" size="iconInlineSize"/>
						</widget>
					</widget>
				</container>
			</content>
		</canvas>
	`
	t.Run("activate pop-up", func(t *testing.T) {
		sel := widget.NewSelect([]string{"Option A", "Option B"}, nil)
		w := test.NewWindow(fyne.NewContainerWithLayout(layout.NewCenterLayout(), sel))
		defer w.Close()
		w.Resize(fyne.NewSize(150, 200))
		c := w.Canvas()
		c.Focus(sel)

		noneSelectedPopup := `
			<canvas padded size="150x200">
				<content>
					<container pos="4,4" size="142x192">
						<widget pos="3,77" size="136x37" type="*widget.Select">
							<widget pos="4,4" size="128x29" type="*widget.Shadow">
								<radialGradient centerOffset="0.5,0.5" pos="-2,-2" size="2x2" startColor="shadow"/>
								<linearGradient endColor="shadow" pos="0,-2" size="128x2"/>
								<radialGradient centerOffset="-0.5,0.5" pos="128,-2" size="2x2" startColor="shadow"/>
								<linearGradient angle="270" pos="128,0" size="2x29" startColor="shadow"/>
								<radialGradient centerOffset="-0.5,-0.5" pos="128,29" size="2x2" startColor="shadow"/>
								<linearGradient pos="0,29" size="128x2" startColor="shadow"/>
								<radialGradient centerOffset="0.5,-0.5" pos="-2,29" size="2x2" startColor="shadow"/>
								<linearGradient angle="270" endColor="shadow" pos="-2,0" size="2x29"/>
							</widget>
							<rectangle fillColor="focus" pos="4,4" size="128x29"/>
							<widget pos="8,4" size="100x29" type="*widget.textProvider">
								<text bold pos="4,4" size="92x21">(Select one)</text>
							</widget>
							<widget pos="108,8" size="20x20" type="*widget.Icon">
								<image fillMode="contain" rsc="menuDropDownIcon" size="iconInlineSize"/>
							</widget>
						</widget>
					</container>
				</content>
				<overlay>
					<widget size="150x200" type="*widget.OverlayContainer">
						<widget pos="11,114" size="128x70" type="*widget.Menu">
							<widget size="128x70" type="*widget.Shadow">
								<radialGradient centerOffset="0.5,0.5" pos="-4,-4" size="4x4" startColor="shadow"/>
								<linearGradient endColor="shadow" pos="0,-4" size="128x4"/>
								<radialGradient centerOffset="-0.5,0.5" pos="128,-4" size="4x4" startColor="shadow"/>
								<linearGradient angle="270" pos="128,0" size="4x70" startColor="shadow"/>
								<radialGradient centerOffset="-0.5,-0.5" pos="128,70" size="4x4" startColor="shadow"/>
								<linearGradient pos="0,70" size="128x4" startColor="shadow"/>
								<radialGradient centerOffset="0.5,-0.5" pos="-4,70" size="4x4" startColor="shadow"/>
								<linearGradient angle="270" endColor="shadow" pos="-4,0" size="4x70"/>
							</widget>
							<widget clip="128x70@11,114" size="128x70" type="*widget.ScrollContainer">
								<widget clip="128x70@11,114" size="128x70" type="*widget.menuBox">
									<container clip="128x70@11,114" pos="0,4" size="128x78">
										<widget clip="128x70@11,114" size="128x29" type="*widget.menuItem">
											<text clip="128x70@11,114" pos="8,4" size="63x21">Option A</text>
										</widget>
										<widget clip="128x70@11,114" pos="0,33" size="128x29" type="*widget.menuItem">
											<text clip="128x70@11,114" pos="8,4" size="63x21">Option B</text>
										</widget>
									</container>
								</widget>
							</widget>
						</widget>
					</widget>
				</overlay>
			</canvas>
		`
		test.AssertRendersToMarkup(t, noneSelected, c)
		sel.TypedKey(&fyne.KeyEvent{Name: fyne.KeySpace})
		test.AssertRendersToMarkup(t, noneSelectedPopup, c)
		test.TapCanvas(c, fyne.NewPos(0, 0))

		test.AssertRendersToMarkup(t, noneSelected, c)
		sel.TypedKey(&fyne.KeyEvent{Name: fyne.KeyUp})
		test.AssertRendersToMarkup(t, noneSelectedPopup, c)
		test.TapCanvas(c, fyne.NewPos(0, 0))

		test.AssertRendersToMarkup(t, noneSelected, c)
		sel.TypedKey(&fyne.KeyEvent{Name: fyne.KeyDown})
		test.AssertRendersToMarkup(t, noneSelectedPopup, c)
		test.TapCanvas(c, fyne.NewPos(0, 0))

		test.AssertRendersToMarkup(t, noneSelected, c)
		sel.TypedKey(&fyne.KeyEvent{Name: fyne.KeyEnter})
		test.AssertRendersToMarkup(t, noneSelected, c)
	})

	t.Run("traverse options without pop-up", func(t *testing.T) {
		sel := widget.NewSelect([]string{"Option A", "Option B", "Option C"}, nil)
		w := test.NewWindow(fyne.NewContainerWithLayout(layout.NewCenterLayout(), sel))
		defer w.Close()
		w.Resize(fyne.NewSize(150, 200))
		c := w.Canvas()
		c.Focus(sel)
		test.AssertRendersToMarkup(t, noneSelected, c)

		aSelected := `
			<canvas padded size="150x200">
				<content>
					<container pos="4,4" size="142x192">
						<widget pos="3,77" size="136x37" type="*widget.Select">
							<widget pos="4,4" size="128x29" type="*widget.Shadow">
								<radialGradient centerOffset="0.5,0.5" pos="-2,-2" size="2x2" startColor="shadow"/>
								<linearGradient endColor="shadow" pos="0,-2" size="128x2"/>
								<radialGradient centerOffset="-0.5,0.5" pos="128,-2" size="2x2" startColor="shadow"/>
								<linearGradient angle="270" pos="128,0" size="2x29" startColor="shadow"/>
								<radialGradient centerOffset="-0.5,-0.5" pos="128,29" size="2x2" startColor="shadow"/>
								<linearGradient pos="0,29" size="128x2" startColor="shadow"/>
								<radialGradient centerOffset="0.5,-0.5" pos="-2,29" size="2x2" startColor="shadow"/>
								<linearGradient angle="270" endColor="shadow" pos="-2,0" size="2x29"/>
							</widget>
							<rectangle fillColor="focus" pos="4,4" size="128x29"/>
							<widget pos="8,4" size="100x29" type="*widget.textProvider">
								<text bold pos="4,4" size="92x21">Option A</text>
							</widget>
							<widget pos="108,8" size="20x20" type="*widget.Icon">
								<image fillMode="contain" rsc="menuDropDownIcon" size="iconInlineSize"/>
							</widget>
						</widget>
					</container>
				</content>
			</canvas>
		`
		bSelected := `
			<canvas padded size="150x200">
				<content>
					<container pos="4,4" size="142x192">
						<widget pos="3,77" size="136x37" type="*widget.Select">
							<widget pos="4,4" size="128x29" type="*widget.Shadow">
								<radialGradient centerOffset="0.5,0.5" pos="-2,-2" size="2x2" startColor="shadow"/>
								<linearGradient endColor="shadow" pos="0,-2" size="128x2"/>
								<radialGradient centerOffset="-0.5,0.5" pos="128,-2" size="2x2" startColor="shadow"/>
								<linearGradient angle="270" pos="128,0" size="2x29" startColor="shadow"/>
								<radialGradient centerOffset="-0.5,-0.5" pos="128,29" size="2x2" startColor="shadow"/>
								<linearGradient pos="0,29" size="128x2" startColor="shadow"/>
								<radialGradient centerOffset="0.5,-0.5" pos="-2,29" size="2x2" startColor="shadow"/>
								<linearGradient angle="270" endColor="shadow" pos="-2,0" size="2x29"/>
							</widget>
							<rectangle fillColor="focus" pos="4,4" size="128x29"/>
							<widget pos="8,4" size="100x29" type="*widget.textProvider">
								<text bold pos="4,4" size="92x21">Option B</text>
							</widget>
							<widget pos="108,8" size="20x20" type="*widget.Icon">
								<image fillMode="contain" rsc="menuDropDownIcon" size="iconInlineSize"/>
							</widget>
						</widget>
					</container>
				</content>
			</canvas>
		`
		cSelected := `
			<canvas padded size="150x200">
				<content>
					<container pos="4,4" size="142x192">
						<widget pos="3,77" size="136x37" type="*widget.Select">
							<widget pos="4,4" size="128x29" type="*widget.Shadow">
								<radialGradient centerOffset="0.5,0.5" pos="-2,-2" size="2x2" startColor="shadow"/>
								<linearGradient endColor="shadow" pos="0,-2" size="128x2"/>
								<radialGradient centerOffset="-0.5,0.5" pos="128,-2" size="2x2" startColor="shadow"/>
								<linearGradient angle="270" pos="128,0" size="2x29" startColor="shadow"/>
								<radialGradient centerOffset="-0.5,-0.5" pos="128,29" size="2x2" startColor="shadow"/>
								<linearGradient pos="0,29" size="128x2" startColor="shadow"/>
								<radialGradient centerOffset="0.5,-0.5" pos="-2,29" size="2x2" startColor="shadow"/>
								<linearGradient angle="270" endColor="shadow" pos="-2,0" size="2x29"/>
							</widget>
							<rectangle fillColor="focus" pos="4,4" size="128x29"/>
							<widget pos="8,4" size="100x29" type="*widget.textProvider">
								<text bold pos="4,4" size="92x21">Option C</text>
							</widget>
							<widget pos="108,8" size="20x20" type="*widget.Icon">
								<image fillMode="contain" rsc="menuDropDownIcon" size="iconInlineSize"/>
							</widget>
						</widget>
					</container>
				</content>
			</canvas>
		`
		sel.TypedKey(&fyne.KeyEvent{Name: fyne.KeyLeft})
		assert.Equal(t, "Option C", sel.Selected)
		test.AssertRendersToMarkup(t, cSelected, c)
		sel.TypedKey(&fyne.KeyEvent{Name: fyne.KeyLeft})
		assert.Equal(t, "Option B", sel.Selected)
		test.AssertRendersToMarkup(t, bSelected, c)
		sel.TypedKey(&fyne.KeyEvent{Name: fyne.KeyLeft})
		assert.Equal(t, "Option A", sel.Selected)
		test.AssertRendersToMarkup(t, aSelected, c)
		sel.TypedKey(&fyne.KeyEvent{Name: fyne.KeyLeft})
		assert.Equal(t, "Option C", sel.Selected)
		test.AssertRendersToMarkup(t, cSelected, c)

		sel.TypedKey(&fyne.KeyEvent{Name: fyne.KeyRight})
		assert.Equal(t, "Option A", sel.Selected)
		test.AssertRendersToMarkup(t, aSelected, c)
		sel.TypedKey(&fyne.KeyEvent{Name: fyne.KeyRight})
		assert.Equal(t, "Option B", sel.Selected)
		test.AssertRendersToMarkup(t, bSelected, c)
		sel.TypedKey(&fyne.KeyEvent{Name: fyne.KeyRight})
		assert.Equal(t, "Option C", sel.Selected)
		test.AssertRendersToMarkup(t, cSelected, c)
		sel.TypedKey(&fyne.KeyEvent{Name: fyne.KeyRight})
		assert.Equal(t, "Option A", sel.Selected)
		test.AssertRendersToMarkup(t, aSelected, c)
	})

	t.Run("trying to traverse empty options without pop-up", func(t *testing.T) {
		sel := widget.NewSelect([]string{}, nil)
		w := test.NewWindow(fyne.NewContainerWithLayout(layout.NewCenterLayout(), sel))
		defer w.Close()
		w.Resize(fyne.NewSize(150, 200))
		c := w.Canvas()
		c.Focus(sel)
		test.AssertRendersToMarkup(t, noneSelected, c)

		sel.TypedKey(&fyne.KeyEvent{Name: fyne.KeyLeft})
		assert.Equal(t, "", sel.Selected)
		test.AssertRendersToMarkup(t, noneSelected, c)
		sel.TypedKey(&fyne.KeyEvent{Name: fyne.KeyLeft})
		assert.Equal(t, "", sel.Selected)
		test.AssertRendersToMarkup(t, noneSelected, c)

		sel.TypedKey(&fyne.KeyEvent{Name: fyne.KeyRight})
		assert.Equal(t, "", sel.Selected)
		test.AssertRendersToMarkup(t, noneSelected, c)
		sel.TypedKey(&fyne.KeyEvent{Name: fyne.KeyRight})
		assert.Equal(t, "", sel.Selected)
		test.AssertRendersToMarkup(t, noneSelected, c)
	})
}

func TestSelect_Move(t *testing.T) {
	app := test.NewApp()
	defer test.NewApp()
	app.Settings().SetTheme(theme.LightTheme())

	combo := widget.NewSelect([]string{"1", "2"}, nil)
	w := test.NewWindow(combo)
	defer w.Close()
	w.Resize(fyne.NewSize(200, 150))

	combo.Resize(combo.MinSize())
	combo.Move(fyne.NewPos(10, 10))
	test.AssertRendersToMarkup(t, `
		<canvas padded size="200x150">
			<content>
				<widget pos="10,10" size="136x37" type="*widget.Select">
					<widget pos="4,4" size="128x29" type="*widget.Shadow">
						<radialGradient centerOffset="0.5,0.5" pos="-2,-2" size="2x2" startColor="shadow"/>
						<linearGradient endColor="shadow" pos="0,-2" size="128x2"/>
						<radialGradient centerOffset="-0.5,0.5" pos="128,-2" size="2x2" startColor="shadow"/>
						<linearGradient angle="270" pos="128,0" size="2x29" startColor="shadow"/>
						<radialGradient centerOffset="-0.5,-0.5" pos="128,29" size="2x2" startColor="shadow"/>
						<linearGradient pos="0,29" size="128x2" startColor="shadow"/>
						<radialGradient centerOffset="0.5,-0.5" pos="-2,29" size="2x2" startColor="shadow"/>
						<linearGradient angle="270" endColor="shadow" pos="-2,0" size="2x29"/>
					</widget>
					<rectangle pos="4,4" size="128x29"/>
					<widget pos="8,4" size="100x29" type="*widget.textProvider">
						<text bold pos="4,4" size="92x21">(Select one)</text>
					</widget>
					<widget pos="108,8" size="20x20" type="*widget.Icon">
						<image fillMode="contain" rsc="menuDropDownIcon" size="iconInlineSize"/>
					</widget>
				</widget>
			</content>
		</canvas>
	`, w.Canvas())

	combo.Tapped(&fyne.PointEvent{})
	test.AssertRendersToMarkup(t, `
		<canvas padded size="200x150">
			<content>
				<widget pos="10,10" size="136x37" type="*widget.Select">
					<widget pos="4,4" size="128x29" type="*widget.Shadow">
						<radialGradient centerOffset="0.5,0.5" pos="-2,-2" size="2x2" startColor="shadow"/>
						<linearGradient endColor="shadow" pos="0,-2" size="128x2"/>
						<radialGradient centerOffset="-0.5,0.5" pos="128,-2" size="2x2" startColor="shadow"/>
						<linearGradient angle="270" pos="128,0" size="2x29" startColor="shadow"/>
						<radialGradient centerOffset="-0.5,-0.5" pos="128,29" size="2x2" startColor="shadow"/>
						<linearGradient pos="0,29" size="128x2" startColor="shadow"/>
						<radialGradient centerOffset="0.5,-0.5" pos="-2,29" size="2x2" startColor="shadow"/>
						<linearGradient angle="270" endColor="shadow" pos="-2,0" size="2x29"/>
					</widget>
					<rectangle pos="4,4" size="128x29"/>
					<widget pos="8,4" size="100x29" type="*widget.textProvider">
						<text bold pos="4,4" size="92x21">(Select one)</text>
					</widget>
					<widget pos="108,8" size="20x20" type="*widget.Icon">
						<image fillMode="contain" rsc="menuDropDownIcon" size="iconInlineSize"/>
					</widget>
				</widget>
			</content>
			<overlay>
				<widget size="200x150" type="*widget.OverlayContainer">
					<widget pos="14,43" size="128x70" type="*widget.Menu">
						<widget size="128x70" type="*widget.Shadow">
							<radialGradient centerOffset="0.5,0.5" pos="-4,-4" size="4x4" startColor="shadow"/>
							<linearGradient endColor="shadow" pos="0,-4" size="128x4"/>
							<radialGradient centerOffset="-0.5,0.5" pos="128,-4" size="4x4" startColor="shadow"/>
							<linearGradient angle="270" pos="128,0" size="4x70" startColor="shadow"/>
							<radialGradient centerOffset="-0.5,-0.5" pos="128,70" size="4x4" startColor="shadow"/>
							<linearGradient pos="0,70" size="128x4" startColor="shadow"/>
							<radialGradient centerOffset="0.5,-0.5" pos="-4,70" size="4x4" startColor="shadow"/>
							<linearGradient angle="270" endColor="shadow" pos="-4,0" size="4x70"/>
						</widget>
						<widget clip="128x70@14,43" size="128x70" type="*widget.ScrollContainer">
							<widget clip="128x70@14,43" size="128x70" type="*widget.menuBox">
								<container clip="128x70@14,43" pos="0,4" size="128x78">
									<widget clip="128x70@14,43" size="128x29" type="*widget.menuItem">
										<text clip="128x70@14,43" pos="8,4" size="9x21">1</text>
									</widget>
									<widget clip="128x70@14,43" pos="0,33" size="128x29" type="*widget.menuItem">
										<text clip="128x70@14,43" pos="8,4" size="9x21">2</text>
									</widget>
								</container>
							</widget>
						</widget>
					</widget>
				</widget>
			</overlay>
		</canvas>
	`, w.Canvas())

	combo.Move(fyne.NewPos(20, 20))
	test.AssertRendersToMarkup(t, `
		<canvas padded size="200x150">
			<content>
				<widget pos="20,20" size="136x37" type="*widget.Select">
					<widget pos="4,4" size="128x29" type="*widget.Shadow">
						<radialGradient centerOffset="0.5,0.5" pos="-2,-2" size="2x2" startColor="shadow"/>
						<linearGradient endColor="shadow" pos="0,-2" size="128x2"/>
						<radialGradient centerOffset="-0.5,0.5" pos="128,-2" size="2x2" startColor="shadow"/>
						<linearGradient angle="270" pos="128,0" size="2x29" startColor="shadow"/>
						<radialGradient centerOffset="-0.5,-0.5" pos="128,29" size="2x2" startColor="shadow"/>
						<linearGradient pos="0,29" size="128x2" startColor="shadow"/>
						<radialGradient centerOffset="0.5,-0.5" pos="-2,29" size="2x2" startColor="shadow"/>
						<linearGradient angle="270" endColor="shadow" pos="-2,0" size="2x29"/>
					</widget>
					<rectangle pos="4,4" size="128x29"/>
					<widget pos="8,4" size="100x29" type="*widget.textProvider">
						<text bold pos="4,4" size="92x21">(Select one)</text>
					</widget>
					<widget pos="108,8" size="20x20" type="*widget.Icon">
						<image fillMode="contain" rsc="menuDropDownIcon" size="iconInlineSize"/>
					</widget>
				</widget>
			</content>
			<overlay>
				<widget size="200x150" type="*widget.OverlayContainer">
					<widget pos="24,53" size="128x70" type="*widget.Menu">
						<widget size="128x70" type="*widget.Shadow">
							<radialGradient centerOffset="0.5,0.5" pos="-4,-4" size="4x4" startColor="shadow"/>
							<linearGradient endColor="shadow" pos="0,-4" size="128x4"/>
							<radialGradient centerOffset="-0.5,0.5" pos="128,-4" size="4x4" startColor="shadow"/>
							<linearGradient angle="270" pos="128,0" size="4x70" startColor="shadow"/>
							<radialGradient centerOffset="-0.5,-0.5" pos="128,70" size="4x4" startColor="shadow"/>
							<linearGradient pos="0,70" size="128x4" startColor="shadow"/>
							<radialGradient centerOffset="0.5,-0.5" pos="-4,70" size="4x4" startColor="shadow"/>
							<linearGradient angle="270" endColor="shadow" pos="-4,0" size="4x70"/>
						</widget>
						<widget clip="128x70@24,53" size="128x70" type="*widget.ScrollContainer">
							<widget clip="128x70@24,53" size="128x70" type="*widget.menuBox">
								<container clip="128x70@24,53" pos="0,4" size="128x78">
									<widget clip="128x70@24,53" size="128x29" type="*widget.menuItem">
										<text clip="128x70@24,53" pos="8,4" size="9x21">1</text>
									</widget>
									<widget clip="128x70@24,53" pos="0,33" size="128x29" type="*widget.menuItem">
										<text clip="128x70@24,53" pos="8,4" size="9x21">2</text>
									</widget>
								</container>
							</widget>
						</widget>
					</widget>
				</widget>
			</overlay>
		</canvas>
	`, w.Canvas())
}

func TestSelect_PlaceHolder(t *testing.T) {
	combo := widget.NewSelect([]string{"1", "2"}, func(string) {})
	assert.NotEmpty(t, combo.PlaceHolder)

	combo.PlaceHolder = "changed!"
	assert.Equal(t, "changed!", combo.PlaceHolder)
}

func TestSelect_SelectedIndex(t *testing.T) {
	combo := widget.NewSelect([]string{"1", "2"}, func(string) {})

	assert.Equal(t, -1, combo.SelectedIndex())

	combo.SetSelected("2")
	assert.Equal(t, 1, combo.SelectedIndex())

	combo.Selected = "invalid"
	assert.Equal(t, -1, combo.SelectedIndex())
}

func TestSelect_SetSelected(t *testing.T) {
	test.NewApp()
	defer test.NewApp()
	test.ApplyTheme(t, theme.LightTheme())

	var triggered bool
	var triggeredValue string
	combo := widget.NewSelect([]string{"1", "2"}, func(s string) {
		triggered = true
		triggeredValue = s
	})
	w := test.NewWindow(fyne.NewContainerWithLayout(layout.NewCenterLayout(), combo))
	defer w.Close()
	w.Resize(fyne.NewSize(200, 150))

	c := w.Canvas()
	test.AssertRendersToMarkup(t, `
		<canvas padded size="200x150">
			<content>
				<container pos="4,4" size="192x142">
					<widget pos="28,52" size="136x37" type="*widget.Select">
						<widget pos="4,4" size="128x29" type="*widget.Shadow">
							<radialGradient centerOffset="0.5,0.5" pos="-2,-2" size="2x2" startColor="shadow"/>
							<linearGradient endColor="shadow" pos="0,-2" size="128x2"/>
							<radialGradient centerOffset="-0.5,0.5" pos="128,-2" size="2x2" startColor="shadow"/>
							<linearGradient angle="270" pos="128,0" size="2x29" startColor="shadow"/>
							<radialGradient centerOffset="-0.5,-0.5" pos="128,29" size="2x2" startColor="shadow"/>
							<linearGradient pos="0,29" size="128x2" startColor="shadow"/>
							<radialGradient centerOffset="0.5,-0.5" pos="-2,29" size="2x2" startColor="shadow"/>
							<linearGradient angle="270" endColor="shadow" pos="-2,0" size="2x29"/>
						</widget>
						<rectangle pos="4,4" size="128x29"/>
						<widget pos="8,4" size="100x29" type="*widget.textProvider">
							<text bold pos="4,4" size="92x21">(Select one)</text>
						</widget>
						<widget pos="108,8" size="20x20" type="*widget.Icon">
							<image fillMode="contain" rsc="menuDropDownIcon" size="iconInlineSize"/>
						</widget>
					</widget>
				</container>
			</content>
		</canvas>
	`, c)
	combo.SetSelected("2")
	test.AssertRendersToMarkup(t, `
		<canvas padded size="200x150">
			<content>
				<container pos="4,4" size="192x142">
					<widget pos="28,52" size="136x37" type="*widget.Select">
						<widget pos="4,4" size="128x29" type="*widget.Shadow">
							<radialGradient centerOffset="0.5,0.5" pos="-2,-2" size="2x2" startColor="shadow"/>
							<linearGradient endColor="shadow" pos="0,-2" size="128x2"/>
							<radialGradient centerOffset="-0.5,0.5" pos="128,-2" size="2x2" startColor="shadow"/>
							<linearGradient angle="270" pos="128,0" size="2x29" startColor="shadow"/>
							<radialGradient centerOffset="-0.5,-0.5" pos="128,29" size="2x2" startColor="shadow"/>
							<linearGradient pos="0,29" size="128x2" startColor="shadow"/>
							<radialGradient centerOffset="0.5,-0.5" pos="-2,29" size="2x2" startColor="shadow"/>
							<linearGradient angle="270" endColor="shadow" pos="-2,0" size="2x29"/>
						</widget>
						<rectangle pos="4,4" size="128x29"/>
						<widget pos="8,4" size="100x29" type="*widget.textProvider">
							<text bold pos="4,4" size="92x21">2</text>
						</widget>
						<widget pos="108,8" size="20x20" type="*widget.Icon">
							<image fillMode="contain" rsc="menuDropDownIcon" size="iconInlineSize"/>
						</widget>
					</widget>
				</container>
			</content>
		</canvas>
	`, c)
	assert.Equal(t, "2", combo.Selected)
	assert.True(t, triggered)
	assert.Equal(t, "2", triggeredValue)
}

func TestSelect_SetSelected_Callback(t *testing.T) {
	selected := ""
	combo := widget.NewSelect([]string{"1", "2"}, func(s string) {
		selected = s
	})
	combo.SetSelected("2")

	assert.Equal(t, "2", selected)
}

func TestSelect_SetSelected_Invalid(t *testing.T) {
	combo := widget.NewSelect([]string{"1", "2"}, func(string) {})
	combo.SetSelected("3")

	assert.Equal(t, "", combo.Selected)
}

func TestSelect_SetSelected_InvalidNoCallback(t *testing.T) {
	var triggered bool
	combo := widget.NewSelect([]string{"1", "2"}, func(string) {
		triggered = true
	})
	combo.SetSelected("")

	assert.False(t, triggered)
}

func TestSelect_SetSelected_InvalidReplace(t *testing.T) {
	combo := widget.NewSelect([]string{"1", "2"}, func(string) {})
	combo.SetSelected("2")
	combo.SetSelected("3")

	assert.Equal(t, "2", combo.Selected)
}

func TestSelect_SetSelected_NoChangeOnEmpty(t *testing.T) {
	var triggered bool
	combo := widget.NewSelect([]string{"1", "2"}, func(string) { triggered = true })
	combo.SetSelected("")

	assert.False(t, triggered)
}

func TestSelect_SetSelectedIndex(t *testing.T) {
	var triggered bool
	var triggeredValue string
	combo := widget.NewSelect([]string{"1", "2"}, func(s string) {
		triggered = true
		triggeredValue = s
	})
	combo.SetSelectedIndex(1)

	assert.Equal(t, "2", combo.Selected)
	assert.True(t, triggered)
	assert.Equal(t, "2", triggeredValue)
}

func TestSelect_SetSelectedIndex_Invalid(t *testing.T) {
	var triggered bool
	combo := widget.NewSelect([]string{"1", "2"}, func(s string) {
		triggered = true
	})
	combo.SetSelectedIndex(-1)
	assert.Equal(t, -1, combo.SelectedIndex())
	assert.False(t, triggered)
	combo.SetSelectedIndex(2)
	assert.Equal(t, -1, combo.SelectedIndex())
	assert.False(t, triggered)
}

func TestSelect_Tapped(t *testing.T) {
	app := test.NewApp()
	defer test.NewApp()
	app.Settings().SetTheme(theme.LightTheme())

	combo := widget.NewSelect([]string{"1", "2"}, func(s string) {})
	w := test.NewWindow(combo)
	defer w.Close()
	w.Resize(fyne.NewSize(200, 150))
	combo.Resize(combo.MinSize())

	test.Tap(combo)
	canvas := fyne.CurrentApp().Driver().CanvasForObject(combo)
	assert.Equal(t, 1, len(canvas.Overlays().List()))
	test.AssertRendersToMarkup(t, `
		<canvas padded size="200x150">
			<content>
				<widget pos="4,4" size="136x37" type="*widget.Select">
					<widget pos="4,4" size="128x29" type="*widget.Shadow">
						<radialGradient centerOffset="0.5,0.5" pos="-2,-2" size="2x2" startColor="shadow"/>
						<linearGradient endColor="shadow" pos="0,-2" size="128x2"/>
						<radialGradient centerOffset="-0.5,0.5" pos="128,-2" size="2x2" startColor="shadow"/>
						<linearGradient angle="270" pos="128,0" size="2x29" startColor="shadow"/>
						<radialGradient centerOffset="-0.5,-0.5" pos="128,29" size="2x2" startColor="shadow"/>
						<linearGradient pos="0,29" size="128x2" startColor="shadow"/>
						<radialGradient centerOffset="0.5,-0.5" pos="-2,29" size="2x2" startColor="shadow"/>
						<linearGradient angle="270" endColor="shadow" pos="-2,0" size="2x29"/>
					</widget>
					<rectangle fillColor="focus" pos="4,4" size="128x29"/>
					<widget pos="8,4" size="100x29" type="*widget.textProvider">
						<text bold pos="4,4" size="92x21">(Select one)</text>
					</widget>
					<widget pos="108,8" size="20x20" type="*widget.Icon">
						<image fillMode="contain" rsc="menuDropDownIcon" size="iconInlineSize"/>
					</widget>
				</widget>
			</content>
			<overlay>
				<widget size="200x150" type="*widget.OverlayContainer">
					<widget pos="8,37" size="128x70" type="*widget.Menu">
						<widget size="128x70" type="*widget.Shadow">
							<radialGradient centerOffset="0.5,0.5" pos="-4,-4" size="4x4" startColor="shadow"/>
							<linearGradient endColor="shadow" pos="0,-4" size="128x4"/>
							<radialGradient centerOffset="-0.5,0.5" pos="128,-4" size="4x4" startColor="shadow"/>
							<linearGradient angle="270" pos="128,0" size="4x70" startColor="shadow"/>
							<radialGradient centerOffset="-0.5,-0.5" pos="128,70" size="4x4" startColor="shadow"/>
							<linearGradient pos="0,70" size="128x4" startColor="shadow"/>
							<radialGradient centerOffset="0.5,-0.5" pos="-4,70" size="4x4" startColor="shadow"/>
							<linearGradient angle="270" endColor="shadow" pos="-4,0" size="4x70"/>
						</widget>
						<widget clip="128x70@8,37" size="128x70" type="*widget.ScrollContainer">
							<widget clip="128x70@8,37" size="128x70" type="*widget.menuBox">
								<container clip="128x70@8,37" pos="0,4" size="128x78">
									<widget clip="128x70@8,37" size="128x29" type="*widget.menuItem">
										<text clip="128x70@8,37" pos="8,4" size="9x21">1</text>
									</widget>
									<widget clip="128x70@8,37" pos="0,33" size="128x29" type="*widget.menuItem">
										<text clip="128x70@8,37" pos="8,4" size="9x21">2</text>
									</widget>
								</container>
							</widget>
						</widget>
					</widget>
				</widget>
			</overlay>
		</canvas>
	`, w.Canvas())
}

func TestSelect_Tapped_Constrained(t *testing.T) {
	app := test.NewApp()
	defer test.NewApp()
	app.Settings().SetTheme(theme.LightTheme())

	combo := widget.NewSelect([]string{"1", "2"}, func(s string) {})
	w := test.NewWindow(combo)
	defer w.Close()
	w.Resize(fyne.NewSize(200, 150))
	combo.Resize(combo.MinSize())

	canvas := w.Canvas()
	combo.Move(fyne.NewPos(canvas.Size().Width-10, canvas.Size().Height-10))
	test.Tap(combo)
	assert.Equal(t, 1, len(canvas.Overlays().List()))
	test.AssertRendersToMarkup(t, `
		<canvas padded size="200x150">
			<content>
				<widget pos="190,140" size="136x37" type="*widget.Select">
					<widget pos="4,4" size="128x29" type="*widget.Shadow">
						<radialGradient centerOffset="0.5,0.5" pos="-2,-2" size="2x2" startColor="shadow"/>
						<linearGradient endColor="shadow" pos="0,-2" size="128x2"/>
						<radialGradient centerOffset="-0.5,0.5" pos="128,-2" size="2x2" startColor="shadow"/>
						<linearGradient angle="270" pos="128,0" size="2x29" startColor="shadow"/>
						<radialGradient centerOffset="-0.5,-0.5" pos="128,29" size="2x2" startColor="shadow"/>
						<linearGradient pos="0,29" size="128x2" startColor="shadow"/>
						<radialGradient centerOffset="0.5,-0.5" pos="-2,29" size="2x2" startColor="shadow"/>
						<linearGradient angle="270" endColor="shadow" pos="-2,0" size="2x29"/>
					</widget>
					<rectangle fillColor="focus" pos="4,4" size="128x29"/>
					<widget pos="8,4" size="100x29" type="*widget.textProvider">
						<text bold pos="4,4" size="92x21">(Select one)</text>
					</widget>
					<widget pos="108,8" size="20x20" type="*widget.Icon">
						<image fillMode="contain" rsc="menuDropDownIcon" size="iconInlineSize"/>
					</widget>
				</widget>
			</content>
			<overlay>
				<widget size="200x150" type="*widget.OverlayContainer">
					<widget pos="72,80" size="128x70" type="*widget.Menu">
						<widget size="128x70" type="*widget.Shadow">
							<radialGradient centerOffset="0.5,0.5" pos="-4,-4" size="4x4" startColor="shadow"/>
							<linearGradient endColor="shadow" pos="0,-4" size="128x4"/>
							<radialGradient centerOffset="-0.5,0.5" pos="128,-4" size="4x4" startColor="shadow"/>
							<linearGradient angle="270" pos="128,0" size="4x70" startColor="shadow"/>
							<radialGradient centerOffset="-0.5,-0.5" pos="128,70" size="4x4" startColor="shadow"/>
							<linearGradient pos="0,70" size="128x4" startColor="shadow"/>
							<radialGradient centerOffset="0.5,-0.5" pos="-4,70" size="4x4" startColor="shadow"/>
							<linearGradient angle="270" endColor="shadow" pos="-4,0" size="4x70"/>
						</widget>
						<widget clip="128x70@72,80" size="128x70" type="*widget.ScrollContainer">
							<widget clip="128x70@72,80" size="128x70" type="*widget.menuBox">
								<container clip="128x70@72,80" pos="0,4" size="128x78">
									<widget clip="128x70@72,80" size="128x29" type="*widget.menuItem">
										<text clip="128x70@72,80" pos="8,4" size="9x21">1</text>
									</widget>
									<widget clip="128x70@72,80" pos="0,33" size="128x29" type="*widget.menuItem">
										<text clip="128x70@72,80" pos="8,4" size="9x21">2</text>
									</widget>
								</container>
							</widget>
						</widget>
					</widget>
				</widget>
			</overlay>
		</canvas>
	`, w.Canvas())
}

func TestSelect_Layout(t *testing.T) {
	test.NewApp()
	defer test.NewApp()
	test.ApplyTheme(t, theme.LightTheme())

	for name, tt := range map[string]struct {
		placeholder string
		options     []string
		selected    string
		expanded    bool
		want        string
	}{
		"empty": {
			want: `
				<canvas padded size="150x200">
					<content>
						<container pos="4,4" size="142x192">
							<widget pos="3,77" size="136x37" type="*widget.Select">
								<widget pos="4,4" size="128x29" type="*widget.Shadow">
									<radialGradient centerOffset="0.5,0.5" pos="-2,-2" size="2x2" startColor="shadow"/>
									<linearGradient endColor="shadow" pos="0,-2" size="128x2"/>
									<radialGradient centerOffset="-0.5,0.5" pos="128,-2" size="2x2" startColor="shadow"/>
									<linearGradient angle="270" pos="128,0" size="2x29" startColor="shadow"/>
									<radialGradient centerOffset="-0.5,-0.5" pos="128,29" size="2x2" startColor="shadow"/>
									<linearGradient pos="0,29" size="128x2" startColor="shadow"/>
									<radialGradient centerOffset="0.5,-0.5" pos="-2,29" size="2x2" startColor="shadow"/>
									<linearGradient angle="270" endColor="shadow" pos="-2,0" size="2x29"/>
								</widget>
								<rectangle pos="4,4" size="128x29"/>
								<widget pos="8,4" size="100x29" type="*widget.textProvider">
									<text bold pos="4,4" size="92x21">(Select one)</text>
								</widget>
								<widget pos="108,8" size="20x20" type="*widget.Icon">
									<image fillMode="contain" rsc="menuDropDownIcon" size="iconInlineSize"/>
								</widget>
							</widget>
						</container>
					</content>
				</canvas>
			`,
		},
		"empty_placeholder": {
			placeholder: "(Pick 1)",
			want: `
				<canvas padded size="150x200">
					<content>
						<container pos="4,4" size="142x192">
							<widget pos="19,77" size="103x37" type="*widget.Select">
								<widget pos="4,4" size="95x29" type="*widget.Shadow">
									<radialGradient centerOffset="0.5,0.5" pos="-2,-2" size="2x2" startColor="shadow"/>
									<linearGradient endColor="shadow" pos="0,-2" size="95x2"/>
									<radialGradient centerOffset="-0.5,0.5" pos="95,-2" size="2x2" startColor="shadow"/>
									<linearGradient angle="270" pos="95,0" size="2x29" startColor="shadow"/>
									<radialGradient centerOffset="-0.5,-0.5" pos="95,29" size="2x2" startColor="shadow"/>
									<linearGradient pos="0,29" size="95x2" startColor="shadow"/>
									<radialGradient centerOffset="0.5,-0.5" pos="-2,29" size="2x2" startColor="shadow"/>
									<linearGradient angle="270" endColor="shadow" pos="-2,0" size="2x29"/>
								</widget>
								<rectangle pos="4,4" size="95x29"/>
								<widget pos="8,4" size="67x29" type="*widget.textProvider">
									<text bold pos="4,4" size="59x21">(Pick 1)</text>
								</widget>
								<widget pos="75,8" size="20x20" type="*widget.Icon">
									<image fillMode="contain" rsc="menuDropDownIcon" size="iconInlineSize"/>
								</widget>
							</widget>
						</container>
					</content>
				</canvas>
			`,
		},
		"empty_expanded": {
			expanded: true,
			want: `
				<canvas padded size="150x200">
					<content>
						<container pos="4,4" size="142x192">
							<widget pos="3,77" size="136x37" type="*widget.Select">
								<widget pos="4,4" size="128x29" type="*widget.Shadow">
									<radialGradient centerOffset="0.5,0.5" pos="-2,-2" size="2x2" startColor="shadow"/>
									<linearGradient endColor="shadow" pos="0,-2" size="128x2"/>
									<radialGradient centerOffset="-0.5,0.5" pos="128,-2" size="2x2" startColor="shadow"/>
									<linearGradient angle="270" pos="128,0" size="2x29" startColor="shadow"/>
									<radialGradient centerOffset="-0.5,-0.5" pos="128,29" size="2x2" startColor="shadow"/>
									<linearGradient pos="0,29" size="128x2" startColor="shadow"/>
									<radialGradient centerOffset="0.5,-0.5" pos="-2,29" size="2x2" startColor="shadow"/>
									<linearGradient angle="270" endColor="shadow" pos="-2,0" size="2x29"/>
								</widget>
								<rectangle fillColor="focus" pos="4,4" size="128x29"/>
								<widget pos="8,4" size="100x29" type="*widget.textProvider">
									<text bold pos="4,4" size="92x21">(Select one)</text>
								</widget>
								<widget pos="108,8" size="20x20" type="*widget.Icon">
									<image fillMode="contain" rsc="menuDropDownIcon" size="iconInlineSize"/>
								</widget>
							</widget>
						</container>
					</content>
					<overlay>
						<widget size="150x200" type="*widget.OverlayContainer">
							<widget pos="11,114" size="128x8" type="*widget.Menu">
								<widget size="128x8" type="*widget.Shadow">
									<radialGradient centerOffset="0.5,0.5" pos="-4,-4" size="4x4" startColor="shadow"/>
									<linearGradient endColor="shadow" pos="0,-4" size="128x4"/>
									<radialGradient centerOffset="-0.5,0.5" pos="128,-4" size="4x4" startColor="shadow"/>
									<linearGradient angle="270" pos="128,0" size="4x8" startColor="shadow"/>
									<radialGradient centerOffset="-0.5,-0.5" pos="128,8" size="4x4" startColor="shadow"/>
									<linearGradient pos="0,8" size="128x4" startColor="shadow"/>
									<radialGradient centerOffset="0.5,-0.5" pos="-4,8" size="4x4" startColor="shadow"/>
									<linearGradient angle="270" endColor="shadow" pos="-4,0" size="4x8"/>
								</widget>
								<widget clip="128x8@11,114" size="128x8" type="*widget.ScrollContainer">
									<widget clip="128x8@11,114" size="128x8" type="*widget.menuBox">
										<container clip="128x8@11,114" pos="0,4" size="128x16">
										</container>
									</widget>
								</widget>
							</widget>
						</widget>
					</overlay>
				</canvas>
			`,
		},
		"empty_expanded_placeholder": {
			placeholder: "(Pick 1)",
			expanded:    true,
			want: `
				<canvas padded size="150x200">
					<content>
						<container pos="4,4" size="142x192">
							<widget pos="19,77" size="103x37" type="*widget.Select">
								<widget pos="4,4" size="95x29" type="*widget.Shadow">
									<radialGradient centerOffset="0.5,0.5" pos="-2,-2" size="2x2" startColor="shadow"/>
									<linearGradient endColor="shadow" pos="0,-2" size="95x2"/>
									<radialGradient centerOffset="-0.5,0.5" pos="95,-2" size="2x2" startColor="shadow"/>
									<linearGradient angle="270" pos="95,0" size="2x29" startColor="shadow"/>
									<radialGradient centerOffset="-0.5,-0.5" pos="95,29" size="2x2" startColor="shadow"/>
									<linearGradient pos="0,29" size="95x2" startColor="shadow"/>
									<radialGradient centerOffset="0.5,-0.5" pos="-2,29" size="2x2" startColor="shadow"/>
									<linearGradient angle="270" endColor="shadow" pos="-2,0" size="2x29"/>
								</widget>
								<rectangle fillColor="focus" pos="4,4" size="95x29"/>
								<widget pos="8,4" size="67x29" type="*widget.textProvider">
									<text bold pos="4,4" size="59x21">(Pick 1)</text>
								</widget>
								<widget pos="75,8" size="20x20" type="*widget.Icon">
									<image fillMode="contain" rsc="menuDropDownIcon" size="iconInlineSize"/>
								</widget>
							</widget>
						</container>
					</content>
					<overlay>
						<widget size="150x200" type="*widget.OverlayContainer">
							<widget pos="27,114" size="95x8" type="*widget.Menu">
								<widget size="95x8" type="*widget.Shadow">
									<radialGradient centerOffset="0.5,0.5" pos="-4,-4" size="4x4" startColor="shadow"/>
									<linearGradient endColor="shadow" pos="0,-4" size="95x4"/>
									<radialGradient centerOffset="-0.5,0.5" pos="95,-4" size="4x4" startColor="shadow"/>
									<linearGradient angle="270" pos="95,0" size="4x8" startColor="shadow"/>
									<radialGradient centerOffset="-0.5,-0.5" pos="95,8" size="4x4" startColor="shadow"/>
									<linearGradient pos="0,8" size="95x4" startColor="shadow"/>
									<radialGradient centerOffset="0.5,-0.5" pos="-4,8" size="4x4" startColor="shadow"/>
									<linearGradient angle="270" endColor="shadow" pos="-4,0" size="4x8"/>
								</widget>
								<widget clip="95x8@27,114" size="95x8" type="*widget.ScrollContainer">
									<widget clip="95x8@27,114" size="95x8" type="*widget.menuBox">
										<container clip="95x8@27,114" pos="0,4" size="95x16">
										</container>
									</widget>
								</widget>
							</widget>
						</widget>
					</overlay>
				</canvas>
			`,
		},
		"single": {
			options: []string{"Test"},
			want: `
				<canvas padded size="150x200">
					<content>
						<container pos="4,4" size="142x192">
							<widget pos="3,77" size="136x37" type="*widget.Select">
								<widget pos="4,4" size="128x29" type="*widget.Shadow">
									<radialGradient centerOffset="0.5,0.5" pos="-2,-2" size="2x2" startColor="shadow"/>
									<linearGradient endColor="shadow" pos="0,-2" size="128x2"/>
									<radialGradient centerOffset="-0.5,0.5" pos="128,-2" size="2x2" startColor="shadow"/>
									<linearGradient angle="270" pos="128,0" size="2x29" startColor="shadow"/>
									<radialGradient centerOffset="-0.5,-0.5" pos="128,29" size="2x2" startColor="shadow"/>
									<linearGradient pos="0,29" size="128x2" startColor="shadow"/>
									<radialGradient centerOffset="0.5,-0.5" pos="-2,29" size="2x2" startColor="shadow"/>
									<linearGradient angle="270" endColor="shadow" pos="-2,0" size="2x29"/>
								</widget>
								<rectangle pos="4,4" size="128x29"/>
								<widget pos="8,4" size="100x29" type="*widget.textProvider">
									<text bold pos="4,4" size="92x21">(Select one)</text>
								</widget>
								<widget pos="108,8" size="20x20" type="*widget.Icon">
									<image fillMode="contain" rsc="menuDropDownIcon" size="iconInlineSize"/>
								</widget>
							</widget>
						</container>
					</content>
				</canvas>
			`,
		},
		"single_placeholder": {
			placeholder: "(Pick 1)",
			options:     []string{"Test"},
			want: `
				<canvas padded size="150x200">
					<content>
						<container pos="4,4" size="142x192">
							<widget pos="19,77" size="103x37" type="*widget.Select">
								<widget pos="4,4" size="95x29" type="*widget.Shadow">
									<radialGradient centerOffset="0.5,0.5" pos="-2,-2" size="2x2" startColor="shadow"/>
									<linearGradient endColor="shadow" pos="0,-2" size="95x2"/>
									<radialGradient centerOffset="-0.5,0.5" pos="95,-2" size="2x2" startColor="shadow"/>
									<linearGradient angle="270" pos="95,0" size="2x29" startColor="shadow"/>
									<radialGradient centerOffset="-0.5,-0.5" pos="95,29" size="2x2" startColor="shadow"/>
									<linearGradient pos="0,29" size="95x2" startColor="shadow"/>
									<radialGradient centerOffset="0.5,-0.5" pos="-2,29" size="2x2" startColor="shadow"/>
									<linearGradient angle="270" endColor="shadow" pos="-2,0" size="2x29"/>
								</widget>
								<rectangle pos="4,4" size="95x29"/>
								<widget pos="8,4" size="67x29" type="*widget.textProvider">
									<text bold pos="4,4" size="59x21">(Pick 1)</text>
								</widget>
								<widget pos="75,8" size="20x20" type="*widget.Icon">
									<image fillMode="contain" rsc="menuDropDownIcon" size="iconInlineSize"/>
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
							<widget pos="3,77" size="136x37" type="*widget.Select">
								<widget pos="4,4" size="128x29" type="*widget.Shadow">
									<radialGradient centerOffset="0.5,0.5" pos="-2,-2" size="2x2" startColor="shadow"/>
									<linearGradient endColor="shadow" pos="0,-2" size="128x2"/>
									<radialGradient centerOffset="-0.5,0.5" pos="128,-2" size="2x2" startColor="shadow"/>
									<linearGradient angle="270" pos="128,0" size="2x29" startColor="shadow"/>
									<radialGradient centerOffset="-0.5,-0.5" pos="128,29" size="2x2" startColor="shadow"/>
									<linearGradient pos="0,29" size="128x2" startColor="shadow"/>
									<radialGradient centerOffset="0.5,-0.5" pos="-2,29" size="2x2" startColor="shadow"/>
									<linearGradient angle="270" endColor="shadow" pos="-2,0" size="2x29"/>
								</widget>
								<rectangle pos="4,4" size="128x29"/>
								<widget pos="8,4" size="100x29" type="*widget.textProvider">
									<text bold pos="4,4" size="92x21">Test</text>
								</widget>
								<widget pos="108,8" size="20x20" type="*widget.Icon">
									<image fillMode="contain" rsc="menuDropDownIcon" size="iconInlineSize"/>
								</widget>
							</widget>
						</container>
					</content>
				</canvas>
			`,
		},
		"single_selected_placeholder": {
			placeholder: "(Pick 1)",
			options:     []string{"Test"},
			selected:    "Test",
			want: `
				<canvas padded size="150x200">
					<content>
						<container pos="4,4" size="142x192">
							<widget pos="19,77" size="103x37" type="*widget.Select">
								<widget pos="4,4" size="95x29" type="*widget.Shadow">
									<radialGradient centerOffset="0.5,0.5" pos="-2,-2" size="2x2" startColor="shadow"/>
									<linearGradient endColor="shadow" pos="0,-2" size="95x2"/>
									<radialGradient centerOffset="-0.5,0.5" pos="95,-2" size="2x2" startColor="shadow"/>
									<linearGradient angle="270" pos="95,0" size="2x29" startColor="shadow"/>
									<radialGradient centerOffset="-0.5,-0.5" pos="95,29" size="2x2" startColor="shadow"/>
									<linearGradient pos="0,29" size="95x2" startColor="shadow"/>
									<radialGradient centerOffset="0.5,-0.5" pos="-2,29" size="2x2" startColor="shadow"/>
									<linearGradient angle="270" endColor="shadow" pos="-2,0" size="2x29"/>
								</widget>
								<rectangle pos="4,4" size="95x29"/>
								<widget pos="8,4" size="67x29" type="*widget.textProvider">
									<text bold pos="4,4" size="59x21">Test</text>
								</widget>
								<widget pos="75,8" size="20x20" type="*widget.Icon">
									<image fillMode="contain" rsc="menuDropDownIcon" size="iconInlineSize"/>
								</widget>
							</widget>
						</container>
					</content>
				</canvas>
			`,
		},
		"single_expanded": {
			options:  []string{"Test"},
			expanded: true,
			want: `
				<canvas padded size="150x200">
					<content>
						<container pos="4,4" size="142x192">
							<widget pos="3,77" size="136x37" type="*widget.Select">
								<widget pos="4,4" size="128x29" type="*widget.Shadow">
									<radialGradient centerOffset="0.5,0.5" pos="-2,-2" size="2x2" startColor="shadow"/>
									<linearGradient endColor="shadow" pos="0,-2" size="128x2"/>
									<radialGradient centerOffset="-0.5,0.5" pos="128,-2" size="2x2" startColor="shadow"/>
									<linearGradient angle="270" pos="128,0" size="2x29" startColor="shadow"/>
									<radialGradient centerOffset="-0.5,-0.5" pos="128,29" size="2x2" startColor="shadow"/>
									<linearGradient pos="0,29" size="128x2" startColor="shadow"/>
									<radialGradient centerOffset="0.5,-0.5" pos="-2,29" size="2x2" startColor="shadow"/>
									<linearGradient angle="270" endColor="shadow" pos="-2,0" size="2x29"/>
								</widget>
								<rectangle fillColor="focus" pos="4,4" size="128x29"/>
								<widget pos="8,4" size="100x29" type="*widget.textProvider">
									<text bold pos="4,4" size="92x21">(Select one)</text>
								</widget>
								<widget pos="108,8" size="20x20" type="*widget.Icon">
									<image fillMode="contain" rsc="menuDropDownIcon" size="iconInlineSize"/>
								</widget>
							</widget>
						</container>
					</content>
					<overlay>
						<widget size="150x200" type="*widget.OverlayContainer">
							<widget pos="11,114" size="128x37" type="*widget.Menu">
								<widget size="128x37" type="*widget.Shadow">
									<radialGradient centerOffset="0.5,0.5" pos="-4,-4" size="4x4" startColor="shadow"/>
									<linearGradient endColor="shadow" pos="0,-4" size="128x4"/>
									<radialGradient centerOffset="-0.5,0.5" pos="128,-4" size="4x4" startColor="shadow"/>
									<linearGradient angle="270" pos="128,0" size="4x37" startColor="shadow"/>
									<radialGradient centerOffset="-0.5,-0.5" pos="128,37" size="4x4" startColor="shadow"/>
									<linearGradient pos="0,37" size="128x4" startColor="shadow"/>
									<radialGradient centerOffset="0.5,-0.5" pos="-4,37" size="4x4" startColor="shadow"/>
									<linearGradient angle="270" endColor="shadow" pos="-4,0" size="4x37"/>
								</widget>
								<widget clip="128x37@11,114" size="128x37" type="*widget.ScrollContainer">
									<widget clip="128x37@11,114" size="128x37" type="*widget.menuBox">
										<container clip="128x37@11,114" pos="0,4" size="128x45">
											<widget clip="128x37@11,114" size="128x29" type="*widget.menuItem">
												<text clip="128x37@11,114" pos="8,4" size="30x21">Test</text>
											</widget>
										</container>
									</widget>
								</widget>
							</widget>
						</widget>
					</overlay>
				</canvas>
			`,
		},
		"single_expanded_placeholder": {
			placeholder: "(Pick 1)",
			options:     []string{"Test"},
			expanded:    true,
			want: `
				<canvas padded size="150x200">
					<content>
						<container pos="4,4" size="142x192">
							<widget pos="19,77" size="103x37" type="*widget.Select">
								<widget pos="4,4" size="95x29" type="*widget.Shadow">
									<radialGradient centerOffset="0.5,0.5" pos="-2,-2" size="2x2" startColor="shadow"/>
									<linearGradient endColor="shadow" pos="0,-2" size="95x2"/>
									<radialGradient centerOffset="-0.5,0.5" pos="95,-2" size="2x2" startColor="shadow"/>
									<linearGradient angle="270" pos="95,0" size="2x29" startColor="shadow"/>
									<radialGradient centerOffset="-0.5,-0.5" pos="95,29" size="2x2" startColor="shadow"/>
									<linearGradient pos="0,29" size="95x2" startColor="shadow"/>
									<radialGradient centerOffset="0.5,-0.5" pos="-2,29" size="2x2" startColor="shadow"/>
									<linearGradient angle="270" endColor="shadow" pos="-2,0" size="2x29"/>
								</widget>
								<rectangle fillColor="focus" pos="4,4" size="95x29"/>
								<widget pos="8,4" size="67x29" type="*widget.textProvider">
									<text bold pos="4,4" size="59x21">(Pick 1)</text>
								</widget>
								<widget pos="75,8" size="20x20" type="*widget.Icon">
									<image fillMode="contain" rsc="menuDropDownIcon" size="iconInlineSize"/>
								</widget>
							</widget>
						</container>
					</content>
					<overlay>
						<widget size="150x200" type="*widget.OverlayContainer">
							<widget pos="27,114" size="95x37" type="*widget.Menu">
								<widget size="95x37" type="*widget.Shadow">
									<radialGradient centerOffset="0.5,0.5" pos="-4,-4" size="4x4" startColor="shadow"/>
									<linearGradient endColor="shadow" pos="0,-4" size="95x4"/>
									<radialGradient centerOffset="-0.5,0.5" pos="95,-4" size="4x4" startColor="shadow"/>
									<linearGradient angle="270" pos="95,0" size="4x37" startColor="shadow"/>
									<radialGradient centerOffset="-0.5,-0.5" pos="95,37" size="4x4" startColor="shadow"/>
									<linearGradient pos="0,37" size="95x4" startColor="shadow"/>
									<radialGradient centerOffset="0.5,-0.5" pos="-4,37" size="4x4" startColor="shadow"/>
									<linearGradient angle="270" endColor="shadow" pos="-4,0" size="4x37"/>
								</widget>
								<widget clip="95x37@27,114" size="95x37" type="*widget.ScrollContainer">
									<widget clip="95x37@27,114" size="95x37" type="*widget.menuBox">
										<container clip="95x37@27,114" pos="0,4" size="95x45">
											<widget clip="95x37@27,114" size="95x29" type="*widget.menuItem">
												<text clip="95x37@27,114" pos="8,4" size="30x21">Test</text>
											</widget>
										</container>
									</widget>
								</widget>
							</widget>
						</widget>
					</overlay>
				</canvas>
			`,
		},
		"single_expanded_selected": {
			options:  []string{"Test"},
			selected: "Test",
			expanded: true,
			want: `
				<canvas padded size="150x200">
					<content>
						<container pos="4,4" size="142x192">
							<widget pos="3,77" size="136x37" type="*widget.Select">
								<widget pos="4,4" size="128x29" type="*widget.Shadow">
									<radialGradient centerOffset="0.5,0.5" pos="-2,-2" size="2x2" startColor="shadow"/>
									<linearGradient endColor="shadow" pos="0,-2" size="128x2"/>
									<radialGradient centerOffset="-0.5,0.5" pos="128,-2" size="2x2" startColor="shadow"/>
									<linearGradient angle="270" pos="128,0" size="2x29" startColor="shadow"/>
									<radialGradient centerOffset="-0.5,-0.5" pos="128,29" size="2x2" startColor="shadow"/>
									<linearGradient pos="0,29" size="128x2" startColor="shadow"/>
									<radialGradient centerOffset="0.5,-0.5" pos="-2,29" size="2x2" startColor="shadow"/>
									<linearGradient angle="270" endColor="shadow" pos="-2,0" size="2x29"/>
								</widget>
								<rectangle fillColor="focus" pos="4,4" size="128x29"/>
								<widget pos="8,4" size="100x29" type="*widget.textProvider">
									<text bold pos="4,4" size="92x21">Test</text>
								</widget>
								<widget pos="108,8" size="20x20" type="*widget.Icon">
									<image fillMode="contain" rsc="menuDropDownIcon" size="iconInlineSize"/>
								</widget>
							</widget>
						</container>
					</content>
					<overlay>
						<widget size="150x200" type="*widget.OverlayContainer">
							<widget pos="11,114" size="128x37" type="*widget.Menu">
								<widget size="128x37" type="*widget.Shadow">
									<radialGradient centerOffset="0.5,0.5" pos="-4,-4" size="4x4" startColor="shadow"/>
									<linearGradient endColor="shadow" pos="0,-4" size="128x4"/>
									<radialGradient centerOffset="-0.5,0.5" pos="128,-4" size="4x4" startColor="shadow"/>
									<linearGradient angle="270" pos="128,0" size="4x37" startColor="shadow"/>
									<radialGradient centerOffset="-0.5,-0.5" pos="128,37" size="4x4" startColor="shadow"/>
									<linearGradient pos="0,37" size="128x4" startColor="shadow"/>
									<radialGradient centerOffset="0.5,-0.5" pos="-4,37" size="4x4" startColor="shadow"/>
									<linearGradient angle="270" endColor="shadow" pos="-4,0" size="4x37"/>
								</widget>
								<widget clip="128x37@11,114" size="128x37" type="*widget.ScrollContainer">
									<widget clip="128x37@11,114" size="128x37" type="*widget.menuBox">
										<container clip="128x37@11,114" pos="0,4" size="128x45">
											<widget clip="128x37@11,114" size="128x29" type="*widget.menuItem">
												<text clip="128x37@11,114" pos="8,4" size="30x21">Test</text>
											</widget>
										</container>
									</widget>
								</widget>
							</widget>
						</widget>
					</overlay>
				</canvas>
			`,
		},
		"single_expanded_selected_placeholder": {
			placeholder: "(Pick 1)",
			options:     []string{"Test"},
			selected:    "Test",
			expanded:    true,
			want: `
				<canvas padded size="150x200">
					<content>
						<container pos="4,4" size="142x192">
							<widget pos="19,77" size="103x37" type="*widget.Select">
								<widget pos="4,4" size="95x29" type="*widget.Shadow">
									<radialGradient centerOffset="0.5,0.5" pos="-2,-2" size="2x2" startColor="shadow"/>
									<linearGradient endColor="shadow" pos="0,-2" size="95x2"/>
									<radialGradient centerOffset="-0.5,0.5" pos="95,-2" size="2x2" startColor="shadow"/>
									<linearGradient angle="270" pos="95,0" size="2x29" startColor="shadow"/>
									<radialGradient centerOffset="-0.5,-0.5" pos="95,29" size="2x2" startColor="shadow"/>
									<linearGradient pos="0,29" size="95x2" startColor="shadow"/>
									<radialGradient centerOffset="0.5,-0.5" pos="-2,29" size="2x2" startColor="shadow"/>
									<linearGradient angle="270" endColor="shadow" pos="-2,0" size="2x29"/>
								</widget>
								<rectangle fillColor="focus" pos="4,4" size="95x29"/>
								<widget pos="8,4" size="67x29" type="*widget.textProvider">
									<text bold pos="4,4" size="59x21">Test</text>
								</widget>
								<widget pos="75,8" size="20x20" type="*widget.Icon">
									<image fillMode="contain" rsc="menuDropDownIcon" size="iconInlineSize"/>
								</widget>
							</widget>
						</container>
					</content>
					<overlay>
						<widget size="150x200" type="*widget.OverlayContainer">
							<widget pos="27,114" size="95x37" type="*widget.Menu">
								<widget size="95x37" type="*widget.Shadow">
									<radialGradient centerOffset="0.5,0.5" pos="-4,-4" size="4x4" startColor="shadow"/>
									<linearGradient endColor="shadow" pos="0,-4" size="95x4"/>
									<radialGradient centerOffset="-0.5,0.5" pos="95,-4" size="4x4" startColor="shadow"/>
									<linearGradient angle="270" pos="95,0" size="4x37" startColor="shadow"/>
									<radialGradient centerOffset="-0.5,-0.5" pos="95,37" size="4x4" startColor="shadow"/>
									<linearGradient pos="0,37" size="95x4" startColor="shadow"/>
									<radialGradient centerOffset="0.5,-0.5" pos="-4,37" size="4x4" startColor="shadow"/>
									<linearGradient angle="270" endColor="shadow" pos="-4,0" size="4x37"/>
								</widget>
								<widget clip="95x37@27,114" size="95x37" type="*widget.ScrollContainer">
									<widget clip="95x37@27,114" size="95x37" type="*widget.menuBox">
										<container clip="95x37@27,114" pos="0,4" size="95x45">
											<widget clip="95x37@27,114" size="95x29" type="*widget.menuItem">
												<text clip="95x37@27,114" pos="8,4" size="30x21">Test</text>
											</widget>
										</container>
									</widget>
								</widget>
							</widget>
						</widget>
					</overlay>
				</canvas>
			`,
		},
		"multiple": {
			options: []string{"Foo", "Bar"},
			want: `
				<canvas padded size="150x200">
					<content>
						<container pos="4,4" size="142x192">
							<widget pos="3,77" size="136x37" type="*widget.Select">
								<widget pos="4,4" size="128x29" type="*widget.Shadow">
									<radialGradient centerOffset="0.5,0.5" pos="-2,-2" size="2x2" startColor="shadow"/>
									<linearGradient endColor="shadow" pos="0,-2" size="128x2"/>
									<radialGradient centerOffset="-0.5,0.5" pos="128,-2" size="2x2" startColor="shadow"/>
									<linearGradient angle="270" pos="128,0" size="2x29" startColor="shadow"/>
									<radialGradient centerOffset="-0.5,-0.5" pos="128,29" size="2x2" startColor="shadow"/>
									<linearGradient pos="0,29" size="128x2" startColor="shadow"/>
									<radialGradient centerOffset="0.5,-0.5" pos="-2,29" size="2x2" startColor="shadow"/>
									<linearGradient angle="270" endColor="shadow" pos="-2,0" size="2x29"/>
								</widget>
								<rectangle pos="4,4" size="128x29"/>
								<widget pos="8,4" size="100x29" type="*widget.textProvider">
									<text bold pos="4,4" size="92x21">(Select one)</text>
								</widget>
								<widget pos="108,8" size="20x20" type="*widget.Icon">
									<image fillMode="contain" rsc="menuDropDownIcon" size="iconInlineSize"/>
								</widget>
							</widget>
						</container>
					</content>
				</canvas>
			`,
		},
		"multiple_placeholder": {
			placeholder: "(Pick 1)",
			options:     []string{"Foo", "Bar"},
			want: `
				<canvas padded size="150x200">
					<content>
						<container pos="4,4" size="142x192">
							<widget pos="19,77" size="103x37" type="*widget.Select">
								<widget pos="4,4" size="95x29" type="*widget.Shadow">
									<radialGradient centerOffset="0.5,0.5" pos="-2,-2" size="2x2" startColor="shadow"/>
									<linearGradient endColor="shadow" pos="0,-2" size="95x2"/>
									<radialGradient centerOffset="-0.5,0.5" pos="95,-2" size="2x2" startColor="shadow"/>
									<linearGradient angle="270" pos="95,0" size="2x29" startColor="shadow"/>
									<radialGradient centerOffset="-0.5,-0.5" pos="95,29" size="2x2" startColor="shadow"/>
									<linearGradient pos="0,29" size="95x2" startColor="shadow"/>
									<radialGradient centerOffset="0.5,-0.5" pos="-2,29" size="2x2" startColor="shadow"/>
									<linearGradient angle="270" endColor="shadow" pos="-2,0" size="2x29"/>
								</widget>
								<rectangle pos="4,4" size="95x29"/>
								<widget pos="8,4" size="67x29" type="*widget.textProvider">
									<text bold pos="4,4" size="59x21">(Pick 1)</text>
								</widget>
								<widget pos="75,8" size="20x20" type="*widget.Icon">
									<image fillMode="contain" rsc="menuDropDownIcon" size="iconInlineSize"/>
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
							<widget pos="3,77" size="136x37" type="*widget.Select">
								<widget pos="4,4" size="128x29" type="*widget.Shadow">
									<radialGradient centerOffset="0.5,0.5" pos="-2,-2" size="2x2" startColor="shadow"/>
									<linearGradient endColor="shadow" pos="0,-2" size="128x2"/>
									<radialGradient centerOffset="-0.5,0.5" pos="128,-2" size="2x2" startColor="shadow"/>
									<linearGradient angle="270" pos="128,0" size="2x29" startColor="shadow"/>
									<radialGradient centerOffset="-0.5,-0.5" pos="128,29" size="2x2" startColor="shadow"/>
									<linearGradient pos="0,29" size="128x2" startColor="shadow"/>
									<radialGradient centerOffset="0.5,-0.5" pos="-2,29" size="2x2" startColor="shadow"/>
									<linearGradient angle="270" endColor="shadow" pos="-2,0" size="2x29"/>
								</widget>
								<rectangle pos="4,4" size="128x29"/>
								<widget pos="8,4" size="100x29" type="*widget.textProvider">
									<text bold pos="4,4" size="92x21">Foo</text>
								</widget>
								<widget pos="108,8" size="20x20" type="*widget.Icon">
									<image fillMode="contain" rsc="menuDropDownIcon" size="iconInlineSize"/>
								</widget>
							</widget>
						</container>
					</content>
				</canvas>
			`,
		},
		"multiple_selected_placeholder": {
			placeholder: "(Pick 1)",
			options:     []string{"Foo", "Bar"},
			selected:    "Foo",
			want: `
				<canvas padded size="150x200">
					<content>
						<container pos="4,4" size="142x192">
							<widget pos="19,77" size="103x37" type="*widget.Select">
								<widget pos="4,4" size="95x29" type="*widget.Shadow">
									<radialGradient centerOffset="0.5,0.5" pos="-2,-2" size="2x2" startColor="shadow"/>
									<linearGradient endColor="shadow" pos="0,-2" size="95x2"/>
									<radialGradient centerOffset="-0.5,0.5" pos="95,-2" size="2x2" startColor="shadow"/>
									<linearGradient angle="270" pos="95,0" size="2x29" startColor="shadow"/>
									<radialGradient centerOffset="-0.5,-0.5" pos="95,29" size="2x2" startColor="shadow"/>
									<linearGradient pos="0,29" size="95x2" startColor="shadow"/>
									<radialGradient centerOffset="0.5,-0.5" pos="-2,29" size="2x2" startColor="shadow"/>
									<linearGradient angle="270" endColor="shadow" pos="-2,0" size="2x29"/>
								</widget>
								<rectangle pos="4,4" size="95x29"/>
								<widget pos="8,4" size="67x29" type="*widget.textProvider">
									<text bold pos="4,4" size="59x21">Foo</text>
								</widget>
								<widget pos="75,8" size="20x20" type="*widget.Icon">
									<image fillMode="contain" rsc="menuDropDownIcon" size="iconInlineSize"/>
								</widget>
							</widget>
						</container>
					</content>
				</canvas>
			`,
		},
		"multiple_expanded": {
			options:  []string{"Foo", "Bar"},
			expanded: true,
			want: `
				<canvas padded size="150x200">
					<content>
						<container pos="4,4" size="142x192">
							<widget pos="3,77" size="136x37" type="*widget.Select">
								<widget pos="4,4" size="128x29" type="*widget.Shadow">
									<radialGradient centerOffset="0.5,0.5" pos="-2,-2" size="2x2" startColor="shadow"/>
									<linearGradient endColor="shadow" pos="0,-2" size="128x2"/>
									<radialGradient centerOffset="-0.5,0.5" pos="128,-2" size="2x2" startColor="shadow"/>
									<linearGradient angle="270" pos="128,0" size="2x29" startColor="shadow"/>
									<radialGradient centerOffset="-0.5,-0.5" pos="128,29" size="2x2" startColor="shadow"/>
									<linearGradient pos="0,29" size="128x2" startColor="shadow"/>
									<radialGradient centerOffset="0.5,-0.5" pos="-2,29" size="2x2" startColor="shadow"/>
									<linearGradient angle="270" endColor="shadow" pos="-2,0" size="2x29"/>
								</widget>
								<rectangle fillColor="focus" pos="4,4" size="128x29"/>
								<widget pos="8,4" size="100x29" type="*widget.textProvider">
									<text bold pos="4,4" size="92x21">(Select one)</text>
								</widget>
								<widget pos="108,8" size="20x20" type="*widget.Icon">
									<image fillMode="contain" rsc="menuDropDownIcon" size="iconInlineSize"/>
								</widget>
							</widget>
						</container>
					</content>
					<overlay>
						<widget size="150x200" type="*widget.OverlayContainer">
							<widget pos="11,114" size="128x70" type="*widget.Menu">
								<widget size="128x70" type="*widget.Shadow">
									<radialGradient centerOffset="0.5,0.5" pos="-4,-4" size="4x4" startColor="shadow"/>
									<linearGradient endColor="shadow" pos="0,-4" size="128x4"/>
									<radialGradient centerOffset="-0.5,0.5" pos="128,-4" size="4x4" startColor="shadow"/>
									<linearGradient angle="270" pos="128,0" size="4x70" startColor="shadow"/>
									<radialGradient centerOffset="-0.5,-0.5" pos="128,70" size="4x4" startColor="shadow"/>
									<linearGradient pos="0,70" size="128x4" startColor="shadow"/>
									<radialGradient centerOffset="0.5,-0.5" pos="-4,70" size="4x4" startColor="shadow"/>
									<linearGradient angle="270" endColor="shadow" pos="-4,0" size="4x70"/>
								</widget>
								<widget clip="128x70@11,114" size="128x70" type="*widget.ScrollContainer">
									<widget clip="128x70@11,114" size="128x70" type="*widget.menuBox">
										<container clip="128x70@11,114" pos="0,4" size="128x78">
											<widget clip="128x70@11,114" size="128x29" type="*widget.menuItem">
												<text clip="128x70@11,114" pos="8,4" size="27x21">Foo</text>
											</widget>
											<widget clip="128x70@11,114" pos="0,33" size="128x29" type="*widget.menuItem">
												<text clip="128x70@11,114" pos="8,4" size="25x21">Bar</text>
											</widget>
										</container>
									</widget>
								</widget>
							</widget>
						</widget>
					</overlay>
				</canvas>
			`,
		},
		"multiple_expanded_placeholder": {
			placeholder: "(Pick 1)",
			options:     []string{"Foo", "Bar"},
			expanded:    true,
			want: `
				<canvas padded size="150x200">
					<content>
						<container pos="4,4" size="142x192">
							<widget pos="19,77" size="103x37" type="*widget.Select">
								<widget pos="4,4" size="95x29" type="*widget.Shadow">
									<radialGradient centerOffset="0.5,0.5" pos="-2,-2" size="2x2" startColor="shadow"/>
									<linearGradient endColor="shadow" pos="0,-2" size="95x2"/>
									<radialGradient centerOffset="-0.5,0.5" pos="95,-2" size="2x2" startColor="shadow"/>
									<linearGradient angle="270" pos="95,0" size="2x29" startColor="shadow"/>
									<radialGradient centerOffset="-0.5,-0.5" pos="95,29" size="2x2" startColor="shadow"/>
									<linearGradient pos="0,29" size="95x2" startColor="shadow"/>
									<radialGradient centerOffset="0.5,-0.5" pos="-2,29" size="2x2" startColor="shadow"/>
									<linearGradient angle="270" endColor="shadow" pos="-2,0" size="2x29"/>
								</widget>
								<rectangle fillColor="focus" pos="4,4" size="95x29"/>
								<widget pos="8,4" size="67x29" type="*widget.textProvider">
									<text bold pos="4,4" size="59x21">(Pick 1)</text>
								</widget>
								<widget pos="75,8" size="20x20" type="*widget.Icon">
									<image fillMode="contain" rsc="menuDropDownIcon" size="iconInlineSize"/>
								</widget>
							</widget>
						</container>
					</content>
					<overlay>
						<widget size="150x200" type="*widget.OverlayContainer">
							<widget pos="27,114" size="95x70" type="*widget.Menu">
								<widget size="95x70" type="*widget.Shadow">
									<radialGradient centerOffset="0.5,0.5" pos="-4,-4" size="4x4" startColor="shadow"/>
									<linearGradient endColor="shadow" pos="0,-4" size="95x4"/>
									<radialGradient centerOffset="-0.5,0.5" pos="95,-4" size="4x4" startColor="shadow"/>
									<linearGradient angle="270" pos="95,0" size="4x70" startColor="shadow"/>
									<radialGradient centerOffset="-0.5,-0.5" pos="95,70" size="4x4" startColor="shadow"/>
									<linearGradient pos="0,70" size="95x4" startColor="shadow"/>
									<radialGradient centerOffset="0.5,-0.5" pos="-4,70" size="4x4" startColor="shadow"/>
									<linearGradient angle="270" endColor="shadow" pos="-4,0" size="4x70"/>
								</widget>
								<widget clip="95x70@27,114" size="95x70" type="*widget.ScrollContainer">
									<widget clip="95x70@27,114" size="95x70" type="*widget.menuBox">
										<container clip="95x70@27,114" pos="0,4" size="95x78">
											<widget clip="95x70@27,114" size="95x29" type="*widget.menuItem">
												<text clip="95x70@27,114" pos="8,4" size="27x21">Foo</text>
											</widget>
											<widget clip="95x70@27,114" pos="0,33" size="95x29" type="*widget.menuItem">
												<text clip="95x70@27,114" pos="8,4" size="25x21">Bar</text>
											</widget>
										</container>
									</widget>
								</widget>
							</widget>
						</widget>
					</overlay>
				</canvas>
			`,
		},
		"multiple_expanded_selected": {
			options:  []string{"Foo", "Bar"},
			selected: "Foo",
			expanded: true,
			want: `
				<canvas padded size="150x200">
					<content>
						<container pos="4,4" size="142x192">
							<widget pos="3,77" size="136x37" type="*widget.Select">
								<widget pos="4,4" size="128x29" type="*widget.Shadow">
									<radialGradient centerOffset="0.5,0.5" pos="-2,-2" size="2x2" startColor="shadow"/>
									<linearGradient endColor="shadow" pos="0,-2" size="128x2"/>
									<radialGradient centerOffset="-0.5,0.5" pos="128,-2" size="2x2" startColor="shadow"/>
									<linearGradient angle="270" pos="128,0" size="2x29" startColor="shadow"/>
									<radialGradient centerOffset="-0.5,-0.5" pos="128,29" size="2x2" startColor="shadow"/>
									<linearGradient pos="0,29" size="128x2" startColor="shadow"/>
									<radialGradient centerOffset="0.5,-0.5" pos="-2,29" size="2x2" startColor="shadow"/>
									<linearGradient angle="270" endColor="shadow" pos="-2,0" size="2x29"/>
								</widget>
								<rectangle fillColor="focus" pos="4,4" size="128x29"/>
								<widget pos="8,4" size="100x29" type="*widget.textProvider">
									<text bold pos="4,4" size="92x21">Foo</text>
								</widget>
								<widget pos="108,8" size="20x20" type="*widget.Icon">
									<image fillMode="contain" rsc="menuDropDownIcon" size="iconInlineSize"/>
								</widget>
							</widget>
						</container>
					</content>
					<overlay>
						<widget size="150x200" type="*widget.OverlayContainer">
							<widget pos="11,114" size="128x70" type="*widget.Menu">
								<widget size="128x70" type="*widget.Shadow">
									<radialGradient centerOffset="0.5,0.5" pos="-4,-4" size="4x4" startColor="shadow"/>
									<linearGradient endColor="shadow" pos="0,-4" size="128x4"/>
									<radialGradient centerOffset="-0.5,0.5" pos="128,-4" size="4x4" startColor="shadow"/>
									<linearGradient angle="270" pos="128,0" size="4x70" startColor="shadow"/>
									<radialGradient centerOffset="-0.5,-0.5" pos="128,70" size="4x4" startColor="shadow"/>
									<linearGradient pos="0,70" size="128x4" startColor="shadow"/>
									<radialGradient centerOffset="0.5,-0.5" pos="-4,70" size="4x4" startColor="shadow"/>
									<linearGradient angle="270" endColor="shadow" pos="-4,0" size="4x70"/>
								</widget>
								<widget clip="128x70@11,114" size="128x70" type="*widget.ScrollContainer">
									<widget clip="128x70@11,114" size="128x70" type="*widget.menuBox">
										<container clip="128x70@11,114" pos="0,4" size="128x78">
											<widget clip="128x70@11,114" size="128x29" type="*widget.menuItem">
												<text clip="128x70@11,114" pos="8,4" size="27x21">Foo</text>
											</widget>
											<widget clip="128x70@11,114" pos="0,33" size="128x29" type="*widget.menuItem">
												<text clip="128x70@11,114" pos="8,4" size="25x21">Bar</text>
											</widget>
										</container>
									</widget>
								</widget>
							</widget>
						</widget>
					</overlay>
				</canvas>
			`,
		},
		"multiple_expanded_selected_placeholder": {
			placeholder: "(Pick 1)",
			options:     []string{"Foo", "Bar"},
			selected:    "Foo",
			expanded:    true,
			want: `
				<canvas padded size="150x200">
					<content>
						<container pos="4,4" size="142x192">
							<widget pos="19,77" size="103x37" type="*widget.Select">
								<widget pos="4,4" size="95x29" type="*widget.Shadow">
									<radialGradient centerOffset="0.5,0.5" pos="-2,-2" size="2x2" startColor="shadow"/>
									<linearGradient endColor="shadow" pos="0,-2" size="95x2"/>
									<radialGradient centerOffset="-0.5,0.5" pos="95,-2" size="2x2" startColor="shadow"/>
									<linearGradient angle="270" pos="95,0" size="2x29" startColor="shadow"/>
									<radialGradient centerOffset="-0.5,-0.5" pos="95,29" size="2x2" startColor="shadow"/>
									<linearGradient pos="0,29" size="95x2" startColor="shadow"/>
									<radialGradient centerOffset="0.5,-0.5" pos="-2,29" size="2x2" startColor="shadow"/>
									<linearGradient angle="270" endColor="shadow" pos="-2,0" size="2x29"/>
								</widget>
								<rectangle fillColor="focus" pos="4,4" size="95x29"/>
								<widget pos="8,4" size="67x29" type="*widget.textProvider">
									<text bold pos="4,4" size="59x21">Foo</text>
								</widget>
								<widget pos="75,8" size="20x20" type="*widget.Icon">
									<image fillMode="contain" rsc="menuDropDownIcon" size="iconInlineSize"/>
								</widget>
							</widget>
						</container>
					</content>
					<overlay>
						<widget size="150x200" type="*widget.OverlayContainer">
							<widget pos="27,114" size="95x70" type="*widget.Menu">
								<widget size="95x70" type="*widget.Shadow">
									<radialGradient centerOffset="0.5,0.5" pos="-4,-4" size="4x4" startColor="shadow"/>
									<linearGradient endColor="shadow" pos="0,-4" size="95x4"/>
									<radialGradient centerOffset="-0.5,0.5" pos="95,-4" size="4x4" startColor="shadow"/>
									<linearGradient angle="270" pos="95,0" size="4x70" startColor="shadow"/>
									<radialGradient centerOffset="-0.5,-0.5" pos="95,70" size="4x4" startColor="shadow"/>
									<linearGradient pos="0,70" size="95x4" startColor="shadow"/>
									<radialGradient centerOffset="0.5,-0.5" pos="-4,70" size="4x4" startColor="shadow"/>
									<linearGradient angle="270" endColor="shadow" pos="-4,0" size="4x70"/>
								</widget>
								<widget clip="95x70@27,114" size="95x70" type="*widget.ScrollContainer">
									<widget clip="95x70@27,114" size="95x70" type="*widget.menuBox">
										<container clip="95x70@27,114" pos="0,4" size="95x78">
											<widget clip="95x70@27,114" size="95x29" type="*widget.menuItem">
												<text clip="95x70@27,114" pos="8,4" size="27x21">Foo</text>
											</widget>
											<widget clip="95x70@27,114" pos="0,33" size="95x29" type="*widget.menuItem">
												<text clip="95x70@27,114" pos="8,4" size="25x21">Bar</text>
											</widget>
										</container>
									</widget>
								</widget>
							</widget>
						</widget>
					</overlay>
				</canvas>
			`,
		},
	} {
		t.Run(name, func(t *testing.T) {
			combo := &widget.Select{
				PlaceHolder: tt.placeholder,
				Options:     tt.options,
				Selected:    tt.selected,
			}

			window := test.NewWindow(fyne.NewContainerWithLayout(layout.NewCenterLayout(), combo))
			if tt.expanded {
				test.Tap(combo)
			}
			window.Resize(combo.MinSize().Max(fyne.NewSize(150, 200)))

			test.AssertRendersToMarkup(t, tt.want, window.Canvas())

			window.Close()
		})
	}
}
