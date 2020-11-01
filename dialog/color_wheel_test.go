package dialog

import (
	"testing"

	"fyne.io/fyne"
	"fyne.io/fyne/test"
	"fyne.io/fyne/theme"
)

func Test_colorWheel_Layout(t *testing.T) {
	test.NewApp()
	defer test.NewApp()
	test.ApplyTheme(t, theme.LightTheme())

	wheel := newColorWheel(nil)
	wheel.SetHSLA(180, 100, 50, 255)
	window := test.NewWindow(wheel)
	window.Resize(wheel.MinSize().Max(fyne.NewSize(100, 100)))

	test.AssertImageMatches(t, "color/wheel_layout.png", window.Canvas().Capture())
	test.AssertRendersToMarkup(t, `
		<canvas padded size="128x128">
			<content>
				<widget pos="4,4" size="120x120" type="*dialog.colorWheel">
					<raster size="120x120"/>
					<line pos="0,60" size="120x0" strokeColor="text"/>
					<line size="0x120" strokeColor="text"/>
				</widget>
			</content>
		</canvas>
	`, window.Canvas())

	window.Close()
}
