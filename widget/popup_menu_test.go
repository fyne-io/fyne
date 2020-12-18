package widget

import (
	"image/color"
	"testing"

	"github.com/stretchr/testify/assert"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/test"
)

func TestPopUpMenu_Move(t *testing.T) {
	m, w := setupPopUpMenuTest()
	defer tearDownPopUpMenuTest(w)
	c := w.Canvas()

	m.Show()
	test.AssertRendersToMarkup(t, `
		<canvas padded size="200x200">
			<content>
				<rectangle fillColor="rgba(0,150,150,255)" pos="4,4" size="192x192"/>
			</content>
			<overlay>
				<widget size="200x200" type="*widget.OverlayContainer">
					<widget size="79x70" type="*widget.Menu">
						<widget size="79x70" type="*widget.Shadow">
							<radialGradient centerOffset="0.5,0.5" pos="-4,-4" size="4x4" startColor="shadow"/>
							<linearGradient endColor="shadow" pos="0,-4" size="79x4"/>
							<radialGradient centerOffset="-0.5,0.5" pos="79,-4" size="4x4" startColor="shadow"/>
							<linearGradient angle="270" pos="79,0" size="4x70" startColor="shadow"/>
							<radialGradient centerOffset="-0.5,-0.5" pos="79,70" size="4x4" startColor="shadow"/>
							<linearGradient pos="0,70" size="79x4" startColor="shadow"/>
							<radialGradient centerOffset="0.5,-0.5" pos="-4,70" size="4x4" startColor="shadow"/>
							<linearGradient angle="270" endColor="shadow" pos="-4,0" size="4x70"/>
						</widget>
						<widget size="79x70" type="*widget.ScrollContainer">
							<widget size="79x70" type="*widget.menuBox">
								<container pos="0,4" size="79x78">
									<widget size="79x29" type="*widget.menuItem">
										<text pos="8,4" size="63x21">Option A</text>
									</widget>
									<widget pos="0,33" size="79x29" type="*widget.menuItem">
										<text pos="8,4" size="63x21">Option B</text>
									</widget>
								</container>
							</widget>
						</widget>
					</widget>
				</widget>
			</overlay>
		</canvas>
	`, c)

	m.Move(fyne.NewPos(20, 20))
	test.AssertRendersToMarkup(t, `
		<canvas padded size="200x200">
			<content>
				<rectangle fillColor="rgba(0,150,150,255)" pos="4,4" size="192x192"/>
			</content>
			<overlay>
				<widget size="200x200" type="*widget.OverlayContainer">
					<widget pos="20,20" size="79x70" type="*widget.Menu">
						<widget size="79x70" type="*widget.Shadow">
							<radialGradient centerOffset="0.5,0.5" pos="-4,-4" size="4x4" startColor="shadow"/>
							<linearGradient endColor="shadow" pos="0,-4" size="79x4"/>
							<radialGradient centerOffset="-0.5,0.5" pos="79,-4" size="4x4" startColor="shadow"/>
							<linearGradient angle="270" pos="79,0" size="4x70" startColor="shadow"/>
							<radialGradient centerOffset="-0.5,-0.5" pos="79,70" size="4x4" startColor="shadow"/>
							<linearGradient pos="0,70" size="79x4" startColor="shadow"/>
							<radialGradient centerOffset="0.5,-0.5" pos="-4,70" size="4x4" startColor="shadow"/>
							<linearGradient angle="270" endColor="shadow" pos="-4,0" size="4x70"/>
						</widget>
						<widget size="79x70" type="*widget.ScrollContainer">
							<widget size="79x70" type="*widget.menuBox">
								<container pos="0,4" size="79x78">
									<widget size="79x29" type="*widget.menuItem">
										<text pos="8,4" size="63x21">Option A</text>
									</widget>
									<widget pos="0,33" size="79x29" type="*widget.menuItem">
										<text pos="8,4" size="63x21">Option B</text>
									</widget>
								</container>
							</widget>
						</widget>
					</widget>
				</widget>
			</overlay>
		</canvas>
	`, c)

	m.Move(fyne.NewPos(190, 10))
	test.AssertRendersToMarkup(t, `
		<canvas padded size="200x200">
			<content>
				<rectangle fillColor="rgba(0,150,150,255)" pos="4,4" size="192x192"/>
			</content>
			<overlay>
				<widget size="200x200" type="*widget.OverlayContainer">
					<widget pos="121,10" size="79x70" type="*widget.Menu">
						<widget size="79x70" type="*widget.Shadow">
							<radialGradient centerOffset="0.5,0.5" pos="-4,-4" size="4x4" startColor="shadow"/>
							<linearGradient endColor="shadow" pos="0,-4" size="79x4"/>
							<radialGradient centerOffset="-0.5,0.5" pos="79,-4" size="4x4" startColor="shadow"/>
							<linearGradient angle="270" pos="79,0" size="4x70" startColor="shadow"/>
							<radialGradient centerOffset="-0.5,-0.5" pos="79,70" size="4x4" startColor="shadow"/>
							<linearGradient pos="0,70" size="79x4" startColor="shadow"/>
							<radialGradient centerOffset="0.5,-0.5" pos="-4,70" size="4x4" startColor="shadow"/>
							<linearGradient angle="270" endColor="shadow" pos="-4,0" size="4x70"/>
						</widget>
						<widget size="79x70" type="*widget.ScrollContainer">
							<widget size="79x70" type="*widget.menuBox">
								<container pos="0,4" size="79x78">
									<widget size="79x29" type="*widget.menuItem">
										<text pos="8,4" size="63x21">Option A</text>
									</widget>
									<widget pos="0,33" size="79x29" type="*widget.menuItem">
										<text pos="8,4" size="63x21">Option B</text>
									</widget>
								</container>
							</widget>
						</widget>
					</widget>
				</widget>
			</overlay>
		</canvas>
	`, c)

	m.Move(fyne.NewPos(10, 190))
	test.AssertRendersToMarkup(t, `
		<canvas padded size="200x200">
			<content>
				<rectangle fillColor="rgba(0,150,150,255)" pos="4,4" size="192x192"/>
			</content>
			<overlay>
				<widget size="200x200" type="*widget.OverlayContainer">
					<widget pos="10,130" size="79x70" type="*widget.Menu">
						<widget size="79x70" type="*widget.Shadow">
							<radialGradient centerOffset="0.5,0.5" pos="-4,-4" size="4x4" startColor="shadow"/>
							<linearGradient endColor="shadow" pos="0,-4" size="79x4"/>
							<radialGradient centerOffset="-0.5,0.5" pos="79,-4" size="4x4" startColor="shadow"/>
							<linearGradient angle="270" pos="79,0" size="4x70" startColor="shadow"/>
							<radialGradient centerOffset="-0.5,-0.5" pos="79,70" size="4x4" startColor="shadow"/>
							<linearGradient pos="0,70" size="79x4" startColor="shadow"/>
							<radialGradient centerOffset="0.5,-0.5" pos="-4,70" size="4x4" startColor="shadow"/>
							<linearGradient angle="270" endColor="shadow" pos="-4,0" size="4x70"/>
						</widget>
						<widget size="79x70" type="*widget.ScrollContainer">
							<widget size="79x70" type="*widget.menuBox">
								<container pos="0,4" size="79x78">
									<widget size="79x29" type="*widget.menuItem">
										<text pos="8,4" size="63x21">Option A</text>
									</widget>
									<widget pos="0,33" size="79x29" type="*widget.menuItem">
										<text pos="8,4" size="63x21">Option B</text>
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

func TestPopUpMenu_Resize(t *testing.T) {
	m, w := setupPopUpMenuTest()
	defer tearDownPopUpMenuTest(w)
	c := w.Canvas()

	m.ShowAtPosition(fyne.NewPos(10, 10))
	test.AssertRendersToMarkup(t, `
		<canvas padded size="200x200">
			<content>
				<rectangle fillColor="rgba(0,150,150,255)" pos="4,4" size="192x192"/>
			</content>
			<overlay>
				<widget size="200x200" type="*widget.OverlayContainer">
					<widget pos="10,10" size="79x70" type="*widget.Menu">
						<widget size="79x70" type="*widget.Shadow">
							<radialGradient centerOffset="0.5,0.5" pos="-4,-4" size="4x4" startColor="shadow"/>
							<linearGradient endColor="shadow" pos="0,-4" size="79x4"/>
							<radialGradient centerOffset="-0.5,0.5" pos="79,-4" size="4x4" startColor="shadow"/>
							<linearGradient angle="270" pos="79,0" size="4x70" startColor="shadow"/>
							<radialGradient centerOffset="-0.5,-0.5" pos="79,70" size="4x4" startColor="shadow"/>
							<linearGradient pos="0,70" size="79x4" startColor="shadow"/>
							<radialGradient centerOffset="0.5,-0.5" pos="-4,70" size="4x4" startColor="shadow"/>
							<linearGradient angle="270" endColor="shadow" pos="-4,0" size="4x70"/>
						</widget>
						<widget size="79x70" type="*widget.ScrollContainer">
							<widget size="79x70" type="*widget.menuBox">
								<container pos="0,4" size="79x78">
									<widget size="79x29" type="*widget.menuItem">
										<text pos="8,4" size="63x21">Option A</text>
									</widget>
									<widget pos="0,33" size="79x29" type="*widget.menuItem">
										<text pos="8,4" size="63x21">Option B</text>
									</widget>
								</container>
							</widget>
						</widget>
					</widget>
				</widget>
			</overlay>
		</canvas>
	`, c)

	m.Resize(m.Size().Add(fyne.NewSize(10, 10)))
	test.AssertRendersToMarkup(t, `
		<canvas padded size="200x200">
			<content>
				<rectangle fillColor="rgba(0,150,150,255)" pos="4,4" size="192x192"/>
			</content>
			<overlay>
				<widget size="200x200" type="*widget.OverlayContainer">
					<widget pos="10,10" size="89x80" type="*widget.Menu">
						<widget size="89x80" type="*widget.Shadow">
							<radialGradient centerOffset="0.5,0.5" pos="-4,-4" size="4x4" startColor="shadow"/>
							<linearGradient endColor="shadow" pos="0,-4" size="89x4"/>
							<radialGradient centerOffset="-0.5,0.5" pos="89,-4" size="4x4" startColor="shadow"/>
							<linearGradient angle="270" pos="89,0" size="4x80" startColor="shadow"/>
							<radialGradient centerOffset="-0.5,-0.5" pos="89,80" size="4x4" startColor="shadow"/>
							<linearGradient pos="0,80" size="89x4" startColor="shadow"/>
							<radialGradient centerOffset="0.5,-0.5" pos="-4,80" size="4x4" startColor="shadow"/>
							<linearGradient angle="270" endColor="shadow" pos="-4,0" size="4x80"/>
						</widget>
						<widget size="89x80" type="*widget.ScrollContainer">
							<widget size="89x80" type="*widget.menuBox">
								<container pos="0,4" size="89x88">
									<widget size="89x29" type="*widget.menuItem">
										<text pos="8,4" size="63x21">Option A</text>
									</widget>
									<widget pos="0,33" size="89x29" type="*widget.menuItem">
										<text pos="8,4" size="63x21">Option B</text>
									</widget>
								</container>
							</widget>
						</widget>
					</widget>
				</widget>
			</overlay>
		</canvas>
	`, c)

	largeSize := c.Size().Add(fyne.NewSize(10, 10))
	m.Resize(largeSize)
	test.AssertRendersToMarkup(t, `
		<canvas padded size="200x200">
			<content>
				<rectangle fillColor="rgba(0,150,150,255)" pos="4,4" size="192x192"/>
			</content>
			<overlay>
				<widget size="200x200" type="*widget.OverlayContainer">
					<widget size="210x200" type="*widget.Menu">
						<widget size="210x200" type="*widget.Shadow">
							<radialGradient centerOffset="0.5,0.5" pos="-4,-4" size="4x4" startColor="shadow"/>
							<linearGradient endColor="shadow" pos="0,-4" size="210x4"/>
							<radialGradient centerOffset="-0.5,0.5" pos="210,-4" size="4x4" startColor="shadow"/>
							<linearGradient angle="270" pos="210,0" size="4x200" startColor="shadow"/>
							<radialGradient centerOffset="-0.5,-0.5" pos="210,200" size="4x4" startColor="shadow"/>
							<linearGradient pos="0,200" size="210x4" startColor="shadow"/>
							<radialGradient centerOffset="0.5,-0.5" pos="-4,200" size="4x4" startColor="shadow"/>
							<linearGradient angle="270" endColor="shadow" pos="-4,0" size="4x200"/>
						</widget>
						<widget size="210x200" type="*widget.ScrollContainer">
							<widget size="210x200" type="*widget.menuBox">
								<container pos="0,4" size="210x208">
									<widget size="210x29" type="*widget.menuItem">
										<text pos="8,4" size="63x21">Option A</text>
									</widget>
									<widget pos="0,33" size="210x29" type="*widget.menuItem">
										<text pos="8,4" size="63x21">Option B</text>
									</widget>
								</container>
							</widget>
						</widget>
					</widget>
				</widget>
			</overlay>
		</canvas>
	`, c)
	assert.Equal(t, fyne.NewSize(largeSize.Width, c.Size().Height), m.Size(), "width is larger than canvas; height is limited by canvas (menu scrolls)")
}

func TestPopUpMenu_Show(t *testing.T) {
	m, w := setupPopUpMenuTest()
	defer tearDownPopUpMenuTest(w)
	c := w.Canvas()

	test.AssertRendersToMarkup(t, `
		<canvas padded size="200x200">
			<content>
				<rectangle fillColor="rgba(0,150,150,255)" pos="4,4" size="192x192"/>
			</content>
		</canvas>
	`, c)

	m.Show()
	test.AssertRendersToMarkup(t, `
		<canvas padded size="200x200">
			<content>
				<rectangle fillColor="rgba(0,150,150,255)" pos="4,4" size="192x192"/>
			</content>
			<overlay>
				<widget size="200x200" type="*widget.OverlayContainer">
					<widget size="79x70" type="*widget.Menu">
						<widget size="79x70" type="*widget.Shadow">
							<radialGradient centerOffset="0.5,0.5" pos="-4,-4" size="4x4" startColor="shadow"/>
							<linearGradient endColor="shadow" pos="0,-4" size="79x4"/>
							<radialGradient centerOffset="-0.5,0.5" pos="79,-4" size="4x4" startColor="shadow"/>
							<linearGradient angle="270" pos="79,0" size="4x70" startColor="shadow"/>
							<radialGradient centerOffset="-0.5,-0.5" pos="79,70" size="4x4" startColor="shadow"/>
							<linearGradient pos="0,70" size="79x4" startColor="shadow"/>
							<radialGradient centerOffset="0.5,-0.5" pos="-4,70" size="4x4" startColor="shadow"/>
							<linearGradient angle="270" endColor="shadow" pos="-4,0" size="4x70"/>
						</widget>
						<widget size="79x70" type="*widget.ScrollContainer">
							<widget size="79x70" type="*widget.menuBox">
								<container pos="0,4" size="79x78">
									<widget size="79x29" type="*widget.menuItem">
										<text pos="8,4" size="63x21">Option A</text>
									</widget>
									<widget pos="0,33" size="79x29" type="*widget.menuItem">
										<text pos="8,4" size="63x21">Option B</text>
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

func TestPopUpMenu_ShowAtPosition(t *testing.T) {
	m, w := setupPopUpMenuTest()
	defer tearDownPopUpMenuTest(w)
	c := w.Canvas()

	hidden := `
		<canvas padded size="200x200">
			<content>
				<rectangle fillColor="rgba(0,150,150,255)" pos="4,4" size="192x192"/>
			</content>
		</canvas>
	`
	test.AssertRendersToMarkup(t, hidden, c)

	m.ShowAtPosition(fyne.NewPos(10, 10))
	test.AssertRendersToMarkup(t, `
		<canvas padded size="200x200">
			<content>
				<rectangle fillColor="rgba(0,150,150,255)" pos="4,4" size="192x192"/>
			</content>
			<overlay>
				<widget size="200x200" type="*widget.OverlayContainer">
					<widget pos="10,10" size="79x70" type="*widget.Menu">
						<widget size="79x70" type="*widget.Shadow">
							<radialGradient centerOffset="0.5,0.5" pos="-4,-4" size="4x4" startColor="shadow"/>
							<linearGradient endColor="shadow" pos="0,-4" size="79x4"/>
							<radialGradient centerOffset="-0.5,0.5" pos="79,-4" size="4x4" startColor="shadow"/>
							<linearGradient angle="270" pos="79,0" size="4x70" startColor="shadow"/>
							<radialGradient centerOffset="-0.5,-0.5" pos="79,70" size="4x4" startColor="shadow"/>
							<linearGradient pos="0,70" size="79x4" startColor="shadow"/>
							<radialGradient centerOffset="0.5,-0.5" pos="-4,70" size="4x4" startColor="shadow"/>
							<linearGradient angle="270" endColor="shadow" pos="-4,0" size="4x70"/>
						</widget>
						<widget size="79x70" type="*widget.ScrollContainer">
							<widget size="79x70" type="*widget.menuBox">
								<container pos="0,4" size="79x78">
									<widget size="79x29" type="*widget.menuItem">
										<text pos="8,4" size="63x21">Option A</text>
									</widget>
									<widget pos="0,33" size="79x29" type="*widget.menuItem">
										<text pos="8,4" size="63x21">Option B</text>
									</widget>
								</container>
							</widget>
						</widget>
					</widget>
				</widget>
			</overlay>
		</canvas>
	`, c)

	m.Hide()
	test.AssertRendersToMarkup(t, hidden, c)

	m.ShowAtPosition(fyne.NewPos(190, 10))
	test.AssertRendersToMarkup(t, `
		<canvas padded size="200x200">
			<content>
				<rectangle fillColor="rgba(0,150,150,255)" pos="4,4" size="192x192"/>
			</content>
			<overlay>
				<widget size="200x200" type="*widget.OverlayContainer">
					<widget pos="121,10" size="79x70" type="*widget.Menu">
						<widget size="79x70" type="*widget.Shadow">
							<radialGradient centerOffset="0.5,0.5" pos="-4,-4" size="4x4" startColor="shadow"/>
							<linearGradient endColor="shadow" pos="0,-4" size="79x4"/>
							<radialGradient centerOffset="-0.5,0.5" pos="79,-4" size="4x4" startColor="shadow"/>
							<linearGradient angle="270" pos="79,0" size="4x70" startColor="shadow"/>
							<radialGradient centerOffset="-0.5,-0.5" pos="79,70" size="4x4" startColor="shadow"/>
							<linearGradient pos="0,70" size="79x4" startColor="shadow"/>
							<radialGradient centerOffset="0.5,-0.5" pos="-4,70" size="4x4" startColor="shadow"/>
							<linearGradient angle="270" endColor="shadow" pos="-4,0" size="4x70"/>
						</widget>
						<widget size="79x70" type="*widget.ScrollContainer">
							<widget size="79x70" type="*widget.menuBox">
								<container pos="0,4" size="79x78">
									<widget size="79x29" type="*widget.menuItem">
										<text pos="8,4" size="63x21">Option A</text>
									</widget>
									<widget pos="0,33" size="79x29" type="*widget.menuItem">
										<text pos="8,4" size="63x21">Option B</text>
									</widget>
								</container>
							</widget>
						</widget>
					</widget>
				</widget>
			</overlay>
		</canvas>
	`, c)

	m.Hide()
	test.AssertRendersToMarkup(t, hidden, c)

	m.ShowAtPosition(fyne.NewPos(10, 190))
	test.AssertRendersToMarkup(t, `
		<canvas padded size="200x200">
			<content>
				<rectangle fillColor="rgba(0,150,150,255)" pos="4,4" size="192x192"/>
			</content>
			<overlay>
				<widget size="200x200" type="*widget.OverlayContainer">
					<widget pos="10,130" size="79x70" type="*widget.Menu">
						<widget size="79x70" type="*widget.Shadow">
							<radialGradient centerOffset="0.5,0.5" pos="-4,-4" size="4x4" startColor="shadow"/>
							<linearGradient endColor="shadow" pos="0,-4" size="79x4"/>
							<radialGradient centerOffset="-0.5,0.5" pos="79,-4" size="4x4" startColor="shadow"/>
							<linearGradient angle="270" pos="79,0" size="4x70" startColor="shadow"/>
							<radialGradient centerOffset="-0.5,-0.5" pos="79,70" size="4x4" startColor="shadow"/>
							<linearGradient pos="0,70" size="79x4" startColor="shadow"/>
							<radialGradient centerOffset="0.5,-0.5" pos="-4,70" size="4x4" startColor="shadow"/>
							<linearGradient angle="270" endColor="shadow" pos="-4,0" size="4x70"/>
						</widget>
						<widget size="79x70" type="*widget.ScrollContainer">
							<widget size="79x70" type="*widget.menuBox">
								<container pos="0,4" size="79x78">
									<widget size="79x29" type="*widget.menuItem">
										<text pos="8,4" size="63x21">Option A</text>
									</widget>
									<widget pos="0,33" size="79x29" type="*widget.menuItem">
										<text pos="8,4" size="63x21">Option B</text>
									</widget>
								</container>
							</widget>
						</widget>
					</widget>
				</widget>
			</overlay>
		</canvas>
	`, c)

	m.Hide()
	test.AssertRendersToMarkup(t, hidden, c)
	menuSize := c.Size().Add(fyne.NewSize(10, 10))
	m.Resize(menuSize)

	m.ShowAtPosition(fyne.NewPos(10, 10))
	test.AssertRendersToMarkup(t, `
		<canvas padded size="200x200">
			<content>
				<rectangle fillColor="rgba(0,150,150,255)" pos="4,4" size="192x192"/>
			</content>
			<overlay>
				<widget size="200x200" type="*widget.OverlayContainer">
					<widget size="210x200" type="*widget.Menu">
						<widget size="210x200" type="*widget.Shadow">
							<radialGradient centerOffset="0.5,0.5" pos="-4,-4" size="4x4" startColor="shadow"/>
							<linearGradient endColor="shadow" pos="0,-4" size="210x4"/>
							<radialGradient centerOffset="-0.5,0.5" pos="210,-4" size="4x4" startColor="shadow"/>
							<linearGradient angle="270" pos="210,0" size="4x200" startColor="shadow"/>
							<radialGradient centerOffset="-0.5,-0.5" pos="210,200" size="4x4" startColor="shadow"/>
							<linearGradient pos="0,200" size="210x4" startColor="shadow"/>
							<radialGradient centerOffset="0.5,-0.5" pos="-4,200" size="4x4" startColor="shadow"/>
							<linearGradient angle="270" endColor="shadow" pos="-4,0" size="4x200"/>
						</widget>
						<widget size="210x200" type="*widget.ScrollContainer">
							<widget size="210x200" type="*widget.menuBox">
								<container pos="0,4" size="210x208">
									<widget size="210x29" type="*widget.menuItem">
										<text pos="8,4" size="63x21">Option A</text>
									</widget>
									<widget pos="0,33" size="210x29" type="*widget.menuItem">
										<text pos="8,4" size="63x21">Option B</text>
									</widget>
								</container>
							</widget>
						</widget>
					</widget>
				</widget>
			</overlay>
		</canvas>
	`, c)
	assert.Equal(t, fyne.NewSize(menuSize.Width, c.Size().Height), m.Size(), "width is larger than canvas; height is limited by canvas (menu scrolls)")
}

func setupPopUpMenuTest() (*PopUpMenu, fyne.Window) {
	test.NewApp()

	w := test.NewWindow(canvas.NewRectangle(color.NRGBA{G: 150, B: 150, A: 255}))
	w.Resize(fyne.NewSize(200, 200))
	m := NewPopUpMenu(fyne.NewMenu(
		"",
		fyne.NewMenuItem("Option A", nil),
		fyne.NewMenuItem("Option B", nil),
	), w.Canvas())
	return m, w
}

func tearDownPopUpMenuTest(w fyne.Window) {
	w.Close()
	test.NewApp()
}
