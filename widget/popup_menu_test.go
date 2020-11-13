package widget

import (
	"image/color"
	"testing"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/internal/widget"
	"fyne.io/fyne/test"
	"fyne.io/fyne/theme"

	"github.com/stretchr/testify/assert"
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
						<widget clip="79x70@0,0" size="79x70" type="*widget.ScrollContainer">
							<widget clip="79x70@0,0" size="79x70" type="*widget.menuBox">
								<container clip="79x70@0,0" pos="0,4" size="79x78">
									<widget clip="79x70@0,0" size="79x29" type="*widget.menuItem">
										<text clip="79x70@0,0" pos="8,4" size="63x21">Option A</text>
									</widget>
									<widget clip="79x70@0,0" pos="0,33" size="79x29" type="*widget.menuItem">
										<text clip="79x70@0,0" pos="8,4" size="63x21">Option B</text>
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
						<widget clip="79x70@20,20" size="79x70" type="*widget.ScrollContainer">
							<widget clip="79x70@20,20" size="79x70" type="*widget.menuBox">
								<container clip="79x70@20,20" pos="0,4" size="79x78">
									<widget clip="79x70@20,20" size="79x29" type="*widget.menuItem">
										<text clip="79x70@20,20" pos="8,4" size="63x21">Option A</text>
									</widget>
									<widget clip="79x70@20,20" pos="0,33" size="79x29" type="*widget.menuItem">
										<text clip="79x70@20,20" pos="8,4" size="63x21">Option B</text>
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
						<widget clip="79x70@121,10" size="79x70" type="*widget.ScrollContainer">
							<widget clip="79x70@121,10" size="79x70" type="*widget.menuBox">
								<container clip="79x70@121,10" pos="0,4" size="79x78">
									<widget clip="79x70@121,10" size="79x29" type="*widget.menuItem">
										<text clip="79x70@121,10" pos="8,4" size="63x21">Option A</text>
									</widget>
									<widget clip="79x70@121,10" pos="0,33" size="79x29" type="*widget.menuItem">
										<text clip="79x70@121,10" pos="8,4" size="63x21">Option B</text>
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
						<widget clip="79x70@10,130" size="79x70" type="*widget.ScrollContainer">
							<widget clip="79x70@10,130" size="79x70" type="*widget.menuBox">
								<container clip="79x70@10,130" pos="0,4" size="79x78">
									<widget clip="79x70@10,130" size="79x29" type="*widget.menuItem">
										<text clip="79x70@10,130" pos="8,4" size="63x21">Option A</text>
									</widget>
									<widget clip="79x70@10,130" pos="0,33" size="79x29" type="*widget.menuItem">
										<text clip="79x70@10,130" pos="8,4" size="63x21">Option B</text>
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
						<widget clip="79x70@10,10" size="79x70" type="*widget.ScrollContainer">
							<widget clip="79x70@10,10" size="79x70" type="*widget.menuBox">
								<container clip="79x70@10,10" pos="0,4" size="79x78">
									<widget clip="79x70@10,10" size="79x29" type="*widget.menuItem">
										<text clip="79x70@10,10" pos="8,4" size="63x21">Option A</text>
									</widget>
									<widget clip="79x70@10,10" pos="0,33" size="79x29" type="*widget.menuItem">
										<text clip="79x70@10,10" pos="8,4" size="63x21">Option B</text>
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
						<widget clip="89x80@10,10" size="89x80" type="*widget.ScrollContainer">
							<widget clip="89x80@10,10" size="89x80" type="*widget.menuBox">
								<container clip="89x80@10,10" pos="0,4" size="89x88">
									<widget clip="89x80@10,10" size="89x29" type="*widget.menuItem">
										<text clip="89x80@10,10" pos="8,4" size="63x21">Option A</text>
									</widget>
									<widget clip="89x80@10,10" pos="0,33" size="89x29" type="*widget.menuItem">
										<text clip="89x80@10,10" pos="8,4" size="63x21">Option B</text>
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
						<widget clip="210x200@0,0" size="210x200" type="*widget.ScrollContainer">
							<widget clip="210x200@0,0" size="210x200" type="*widget.menuBox">
								<container clip="210x200@0,0" pos="0,4" size="210x208">
									<widget clip="210x200@0,0" size="210x29" type="*widget.menuItem">
										<text clip="210x200@0,0" pos="8,4" size="63x21">Option A</text>
									</widget>
									<widget clip="210x200@0,0" pos="0,33" size="210x29" type="*widget.menuItem">
										<text clip="210x200@0,0" pos="8,4" size="63x21">Option B</text>
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
						<widget clip="79x70@0,0" size="79x70" type="*widget.ScrollContainer">
							<widget clip="79x70@0,0" size="79x70" type="*widget.menuBox">
								<container clip="79x70@0,0" pos="0,4" size="79x78">
									<widget clip="79x70@0,0" size="79x29" type="*widget.menuItem">
										<text clip="79x70@0,0" pos="8,4" size="63x21">Option A</text>
									</widget>
									<widget clip="79x70@0,0" pos="0,33" size="79x29" type="*widget.menuItem">
										<text clip="79x70@0,0" pos="8,4" size="63x21">Option B</text>
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
						<widget clip="79x70@10,10" size="79x70" type="*widget.ScrollContainer">
							<widget clip="79x70@10,10" size="79x70" type="*widget.menuBox">
								<container clip="79x70@10,10" pos="0,4" size="79x78">
									<widget clip="79x70@10,10" size="79x29" type="*widget.menuItem">
										<text clip="79x70@10,10" pos="8,4" size="63x21">Option A</text>
									</widget>
									<widget clip="79x70@10,10" pos="0,33" size="79x29" type="*widget.menuItem">
										<text clip="79x70@10,10" pos="8,4" size="63x21">Option B</text>
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
						<widget clip="79x70@121,10" size="79x70" type="*widget.ScrollContainer">
							<widget clip="79x70@121,10" size="79x70" type="*widget.menuBox">
								<container clip="79x70@121,10" pos="0,4" size="79x78">
									<widget clip="79x70@121,10" size="79x29" type="*widget.menuItem">
										<text clip="79x70@121,10" pos="8,4" size="63x21">Option A</text>
									</widget>
									<widget clip="79x70@121,10" pos="0,33" size="79x29" type="*widget.menuItem">
										<text clip="79x70@121,10" pos="8,4" size="63x21">Option B</text>
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
						<widget clip="79x70@10,130" size="79x70" type="*widget.ScrollContainer">
							<widget clip="79x70@10,130" size="79x70" type="*widget.menuBox">
								<container clip="79x70@10,130" pos="0,4" size="79x78">
									<widget clip="79x70@10,130" size="79x29" type="*widget.menuItem">
										<text clip="79x70@10,130" pos="8,4" size="63x21">Option A</text>
									</widget>
									<widget clip="79x70@10,130" pos="0,33" size="79x29" type="*widget.menuItem">
										<text clip="79x70@10,130" pos="8,4" size="63x21">Option B</text>
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
						<widget clip="210x200@0,0" size="210x200" type="*widget.ScrollContainer">
							<widget clip="210x200@0,0" size="210x200" type="*widget.menuBox">
								<container clip="210x200@0,0" pos="0,4" size="210x208">
									<widget clip="210x200@0,0" size="210x29" type="*widget.menuItem">
										<text clip="210x200@0,0" pos="8,4" size="63x21">Option A</text>
									</widget>
									<widget clip="210x200@0,0" pos="0,33" size="210x29" type="*widget.menuItem">
										<text clip="210x200@0,0" pos="8,4" size="63x21">Option B</text>
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
	m := newPopUpMenu(fyne.NewMenu(
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

//
// Old pop-up menu tests
//

func TestNewPopUpMenu(t *testing.T) {
	c := test.Canvas()
	menu := fyne.NewMenu("Foo", fyne.NewMenuItem("Bar", func() {}))

	pop := NewPopUpMenu(menu, c)
	assert.Equal(t, 1, len(c.Overlays().List()))
	assert.Equal(t, pop, c.Overlays().List()[0])

	pop.Hide()
	assert.Equal(t, 0, len(c.Overlays().List()))
}

func TestPopUpMenu_Size(t *testing.T) {
	win := test.NewWindow(NewLabel("OK"))
	defer win.Close()
	win.Resize(fyne.NewSize(100, 100))
	menu := fyne.NewMenu("Foo",
		fyne.NewMenuItem("A", func() {}),
		fyne.NewMenuItem("A", func() {}),
	)
	menuItemSize := canvas.NewText("A", color.Black).MinSize().Add(fyne.NewSize(theme.Padding()*2, theme.Padding()*2))
	expectedSize := menuItemSize.Add(fyne.NewSize(0, menuItemSize.Height)).Add(fyne.NewSize(0, theme.Padding()))
	c := win.Canvas()

	pop := NewPopUpMenu(menu, c)
	defer pop.Hide()
	assert.Equal(t, expectedSize, pop.Content.Size())

	for _, o := range test.LaidOutObjects(pop) {
		if s, ok := o.(*widget.Shadow); ok {
			assert.Equal(t, expectedSize.Add(fyne.NewSize(theme.Padding()*2, theme.Padding()*2)), s.Size())
		}
	}
}
