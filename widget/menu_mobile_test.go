// +build mobile

package widget_test

import (
	"image/color"
	"testing"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	internalWidget "fyne.io/fyne/internal/widget"
	"fyne.io/fyne/test"
	"fyne.io/fyne/widget"
)

func TestMenu_Layout(t *testing.T) {
	test.NewApp()
	defer test.NewApp()

	w := test.NewWindow(canvas.NewRectangle(color.Transparent))
	defer w.Close()
	w.SetPadded(false)
	c := w.Canvas()

	item1 := fyne.NewMenuItem("A", nil)
	item2 := fyne.NewMenuItem("B (long)", nil)
	sep := fyne.NewMenuItemSeparator()
	item3 := fyne.NewMenuItem("C", nil)
	subItem1 := fyne.NewMenuItem("subitem A", nil)
	subItem2 := fyne.NewMenuItem("subitem B", nil)
	subItem3 := fyne.NewMenuItem("subitem C (long)", nil)
	subsubItem1 := fyne.NewMenuItem("subsubitem A (long)", nil)
	subsubItem2 := fyne.NewMenuItem("subsubitem B", nil)
	subItem3.ChildMenu = fyne.NewMenu("", subsubItem1, subsubItem2)
	item3.ChildMenu = fyne.NewMenu("", subItem1, subItem2, subItem3)
	menu := fyne.NewMenu("", item1, sep, item2, item3)

	for name, tt := range map[string]struct {
		windowSize   fyne.Size
		menuPos      fyne.Position
		tapPositions []fyne.Position
		useTestTheme bool
		want         string
	}{
		"normal": {
			windowSize: fyne.NewSize(500, 300),
			menuPos:    fyne.NewPos(10, 10),
			want: `
				<canvas size="500x300">
					<content>
						<rectangle size="500x300"/>
					</content>
					<overlay>
						<widget size="500x300" type="*widget.OverlayContainer">
							<widget pos="10,10" size="71x108" type="*widget.Menu">
								<widget size="71x108" type="*widget.Shadow">
									<radialGradient centerOffset="0.5,0.5" pos="-4,-4" size="4x4" startColor="shadow"/>
									<linearGradient endColor="shadow" pos="0,-4" size="71x4"/>
									<radialGradient centerOffset="-0.5,0.5" pos="71,-4" size="4x4" startColor="shadow"/>
									<linearGradient angle="270" pos="71,0" size="4x108" startColor="shadow"/>
									<radialGradient centerOffset="-0.5,-0.5" pos="71,108" size="4x4" startColor="shadow"/>
									<linearGradient pos="0,108" size="71x4" startColor="shadow"/>
									<radialGradient centerOffset="0.5,-0.5" pos="-4,108" size="4x4" startColor="shadow"/>
									<linearGradient angle="270" endColor="shadow" pos="-4,0" size="4x108"/>
								</widget>
								<widget size="71x108" type="*widget.ScrollContainer">
									<widget size="71x108" type="*widget.menuBox">
										<container pos="0,4" size="71x116">
											<widget size="71x29" type="*widget.menuItem">
												<text pos="8,4" size="10x21">A</text>
											</widget>
											<widget pos="0,33" size="71x1" type="*widget.Separator">
												<rectangle fillColor="disabled text" size="71x1"/>
											</widget>
											<widget pos="0,38" size="71x29" type="*widget.menuItem">
												<text pos="8,4" size="55x21">B (long)</text>
											</widget>
											<widget pos="0,71" size="71x29" type="*widget.menuItem">
												<text pos="8,4" size="10x21">C</text>
												<image pos="51,4" rsc="menuExpandIcon" size="iconInlineSize"/>
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
		"normal with submenus": {
			windowSize: fyne.NewSize(500, 300),
			menuPos:    fyne.NewPos(10, 10),
			tapPositions: []fyne.Position{
				fyne.NewPos(30, 100),
				fyne.NewPos(100, 170),
			},
			want: `
				<canvas size="500x300">
					<content>
						<rectangle size="500x300"/>
					</content>
					<overlay>
						<widget size="500x300" type="*widget.OverlayContainer">
							<widget pos="10,10" size="71x108" type="*widget.Menu">
								<widget size="71x108" type="*widget.Shadow">
									<radialGradient centerOffset="0.5,0.5" pos="-4,-4" size="4x4" startColor="shadow"/>
									<linearGradient endColor="shadow" pos="0,-4" size="71x4"/>
									<radialGradient centerOffset="-0.5,0.5" pos="71,-4" size="4x4" startColor="shadow"/>
									<linearGradient angle="270" pos="71,0" size="4x108" startColor="shadow"/>
									<radialGradient centerOffset="-0.5,-0.5" pos="71,108" size="4x4" startColor="shadow"/>
									<linearGradient pos="0,108" size="71x4" startColor="shadow"/>
									<radialGradient centerOffset="0.5,-0.5" pos="-4,108" size="4x4" startColor="shadow"/>
									<linearGradient angle="270" endColor="shadow" pos="-4,0" size="4x108"/>
								</widget>
								<widget size="71x108" type="*widget.ScrollContainer">
									<widget size="71x108" type="*widget.menuBox">
										<container pos="0,4" size="71x116">
											<widget size="71x29" type="*widget.menuItem">
												<text pos="8,4" size="10x21">A</text>
											</widget>
											<widget pos="0,33" size="71x1" type="*widget.Separator">
												<rectangle fillColor="disabled text" size="71x1"/>
											</widget>
											<widget pos="0,38" size="71x29" type="*widget.menuItem">
												<text pos="8,4" size="55x21">B (long)</text>
											</widget>
											<widget pos="0,71" size="71x29" type="*widget.menuItem">
												<text pos="8,4" size="10x21">C</text>
												<image pos="51,4" rsc="menuExpandIcon" size="iconInlineSize"/>
											</widget>
										</container>
									</widget>
								</widget>
								<widget pos="71,71" size="153x103" type="*widget.Menu">
									<widget size="153x103" type="*widget.Shadow">
										<radialGradient centerOffset="0.5,0.5" pos="-4,-4" size="4x4" startColor="shadow"/>
										<linearGradient endColor="shadow" pos="0,-4" size="153x4"/>
										<radialGradient centerOffset="-0.5,0.5" pos="153,-4" size="4x4" startColor="shadow"/>
										<linearGradient angle="270" pos="153,0" size="4x103" startColor="shadow"/>
										<radialGradient centerOffset="-0.5,-0.5" pos="153,103" size="4x4" startColor="shadow"/>
										<linearGradient pos="0,103" size="153x4" startColor="shadow"/>
										<radialGradient centerOffset="0.5,-0.5" pos="-4,103" size="4x4" startColor="shadow"/>
										<linearGradient angle="270" endColor="shadow" pos="-4,0" size="4x103"/>
									</widget>
									<widget size="153x103" type="*widget.ScrollContainer">
										<widget size="153x103" type="*widget.menuBox">
											<container pos="0,4" size="153x111">
												<widget size="153x29" type="*widget.menuItem">
													<text pos="8,4" size="72x21">subitem A</text>
												</widget>
												<widget pos="0,33" size="153x29" type="*widget.menuItem">
													<text pos="8,4" size="72x21">subitem B</text>
												</widget>
												<widget pos="0,66" size="153x29" type="*widget.menuItem">
													<text pos="8,4" size="117x21">subitem C (long)</text>
													<image pos="133,4" rsc="menuExpandIcon" size="iconInlineSize"/>
												</widget>
											</container>
										</widget>
									</widget>
									<widget pos="153,66" size="159x70" type="*widget.Menu">
										<widget size="159x70" type="*widget.Shadow">
											<radialGradient centerOffset="0.5,0.5" pos="-4,-4" size="4x4" startColor="shadow"/>
											<linearGradient endColor="shadow" pos="0,-4" size="159x4"/>
											<radialGradient centerOffset="-0.5,0.5" pos="159,-4" size="4x4" startColor="shadow"/>
											<linearGradient angle="270" pos="159,0" size="4x70" startColor="shadow"/>
											<radialGradient centerOffset="-0.5,-0.5" pos="159,70" size="4x4" startColor="shadow"/>
											<linearGradient pos="0,70" size="159x4" startColor="shadow"/>
											<radialGradient centerOffset="0.5,-0.5" pos="-4,70" size="4x4" startColor="shadow"/>
											<linearGradient angle="270" endColor="shadow" pos="-4,0" size="4x70"/>
										</widget>
										<widget size="159x70" type="*widget.ScrollContainer">
											<widget size="159x70" type="*widget.menuBox">
												<container pos="0,4" size="159x78">
													<widget size="159x29" type="*widget.menuItem">
														<text pos="8,4" size="143x21">subsubitem A (long)</text>
													</widget>
													<widget pos="0,33" size="159x29" type="*widget.menuItem">
														<text pos="8,4" size="98x21">subsubitem B</text>
													</widget>
												</container>
											</widget>
										</widget>
									</widget>
								</widget>
							</widget>
						</widget>
					</overlay>
				</canvas>
			`,
		},
		"background of active submenu parents resets if sibling is hovered": {
			windowSize: fyne.NewSize(500, 300),
			menuPos:    fyne.NewPos(10, 10),
			tapPositions: []fyne.Position{
				fyne.NewPos(30, 100),  // open submenu
				fyne.NewPos(100, 170), // open subsubmenu
				fyne.NewPos(300, 170), // hover subsubmenu item
				fyne.NewPos(30, 60),   // hover sibling of submenu parent
			},
			want: `
				<canvas size="500x300">
					<content>
						<rectangle size="500x300"/>
					</content>
					<overlay>
						<widget size="500x300" type="*widget.OverlayContainer">
							<widget pos="10,10" size="71x108" type="*widget.Menu">
								<widget size="71x108" type="*widget.Shadow">
									<radialGradient centerOffset="0.5,0.5" pos="-4,-4" size="4x4" startColor="shadow"/>
									<linearGradient endColor="shadow" pos="0,-4" size="71x4"/>
									<radialGradient centerOffset="-0.5,0.5" pos="71,-4" size="4x4" startColor="shadow"/>
									<linearGradient angle="270" pos="71,0" size="4x108" startColor="shadow"/>
									<radialGradient centerOffset="-0.5,-0.5" pos="71,108" size="4x4" startColor="shadow"/>
									<linearGradient pos="0,108" size="71x4" startColor="shadow"/>
									<radialGradient centerOffset="0.5,-0.5" pos="-4,108" size="4x4" startColor="shadow"/>
									<linearGradient angle="270" endColor="shadow" pos="-4,0" size="4x108"/>
								</widget>
								<widget size="71x108" type="*widget.ScrollContainer">
									<widget size="71x108" type="*widget.menuBox">
										<container pos="0,4" size="71x116">
											<widget size="71x29" type="*widget.menuItem">
												<text pos="8,4" size="10x21">A</text>
											</widget>
											<widget pos="0,33" size="71x1" type="*widget.Separator">
												<rectangle fillColor="disabled text" size="71x1"/>
											</widget>
											<widget pos="0,38" size="71x29" type="*widget.menuItem">
												<text pos="8,4" size="55x21">B (long)</text>
											</widget>
											<widget pos="0,71" size="71x29" type="*widget.menuItem">
												<text pos="8,4" size="10x21">C</text>
												<image pos="51,4" rsc="menuExpandIcon" size="iconInlineSize"/>
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
		"no space on right side for submenu": {
			windowSize: fyne.NewSize(500, 300),
			menuPos:    fyne.NewPos(410, 10),
			tapPositions: []fyne.Position{
				fyne.NewPos(430, 100), // open submenu
				fyne.NewPos(300, 170), // open subsubmenu
			},
			want: `
				<canvas size="500x300">
					<content>
						<rectangle size="500x300"/>
					</content>
					<overlay>
						<widget size="500x300" type="*widget.OverlayContainer">
							<widget pos="410,10" size="71x108" type="*widget.Menu">
								<widget size="71x108" type="*widget.Shadow">
									<radialGradient centerOffset="0.5,0.5" pos="-4,-4" size="4x4" startColor="shadow"/>
									<linearGradient endColor="shadow" pos="0,-4" size="71x4"/>
									<radialGradient centerOffset="-0.5,0.5" pos="71,-4" size="4x4" startColor="shadow"/>
									<linearGradient angle="270" pos="71,0" size="4x108" startColor="shadow"/>
									<radialGradient centerOffset="-0.5,-0.5" pos="71,108" size="4x4" startColor="shadow"/>
									<linearGradient pos="0,108" size="71x4" startColor="shadow"/>
									<radialGradient centerOffset="0.5,-0.5" pos="-4,108" size="4x4" startColor="shadow"/>
									<linearGradient angle="270" endColor="shadow" pos="-4,0" size="4x108"/>
								</widget>
								<widget size="71x108" type="*widget.ScrollContainer">
									<widget size="71x108" type="*widget.menuBox">
										<container pos="0,4" size="71x116">
											<widget size="71x29" type="*widget.menuItem">
												<text pos="8,4" size="10x21">A</text>
											</widget>
											<widget pos="0,33" size="71x1" type="*widget.Separator">
												<rectangle fillColor="disabled text" size="71x1"/>
											</widget>
											<widget pos="0,38" size="71x29" type="*widget.menuItem">
												<text pos="8,4" size="55x21">B (long)</text>
											</widget>
											<widget pos="0,71" size="71x29" type="*widget.menuItem">
												<text pos="8,4" size="10x21">C</text>
												<image pos="51,4" rsc="menuExpandIcon" size="iconInlineSize"/>
											</widget>
										</container>
									</widget>
								</widget>
								<widget pos="-153,71" size="153x103" type="*widget.Menu">
									<widget size="153x103" type="*widget.Shadow">
										<radialGradient centerOffset="0.5,0.5" pos="-4,-4" size="4x4" startColor="shadow"/>
										<linearGradient endColor="shadow" pos="0,-4" size="153x4"/>
										<radialGradient centerOffset="-0.5,0.5" pos="153,-4" size="4x4" startColor="shadow"/>
										<linearGradient angle="270" pos="153,0" size="4x103" startColor="shadow"/>
										<radialGradient centerOffset="-0.5,-0.5" pos="153,103" size="4x4" startColor="shadow"/>
										<linearGradient pos="0,103" size="153x4" startColor="shadow"/>
										<radialGradient centerOffset="0.5,-0.5" pos="-4,103" size="4x4" startColor="shadow"/>
										<linearGradient angle="270" endColor="shadow" pos="-4,0" size="4x103"/>
									</widget>
									<widget size="153x103" type="*widget.ScrollContainer">
										<widget size="153x103" type="*widget.menuBox">
											<container pos="0,4" size="153x111">
												<widget size="153x29" type="*widget.menuItem">
													<text pos="8,4" size="72x21">subitem A</text>
												</widget>
												<widget pos="0,33" size="153x29" type="*widget.menuItem">
													<text pos="8,4" size="72x21">subitem B</text>
												</widget>
												<widget pos="0,66" size="153x29" type="*widget.menuItem">
													<text pos="8,4" size="117x21">subitem C (long)</text>
													<image pos="133,4" rsc="menuExpandIcon" size="iconInlineSize"/>
												</widget>
											</container>
										</widget>
									</widget>
									<widget pos="-159,66" size="159x70" type="*widget.Menu">
										<widget size="159x70" type="*widget.Shadow">
											<radialGradient centerOffset="0.5,0.5" pos="-4,-4" size="4x4" startColor="shadow"/>
											<linearGradient endColor="shadow" pos="0,-4" size="159x4"/>
											<radialGradient centerOffset="-0.5,0.5" pos="159,-4" size="4x4" startColor="shadow"/>
											<linearGradient angle="270" pos="159,0" size="4x70" startColor="shadow"/>
											<radialGradient centerOffset="-0.5,-0.5" pos="159,70" size="4x4" startColor="shadow"/>
											<linearGradient pos="0,70" size="159x4" startColor="shadow"/>
											<radialGradient centerOffset="0.5,-0.5" pos="-4,70" size="4x4" startColor="shadow"/>
											<linearGradient angle="270" endColor="shadow" pos="-4,0" size="4x70"/>
										</widget>
										<widget size="159x70" type="*widget.ScrollContainer">
											<widget size="159x70" type="*widget.menuBox">
												<container pos="0,4" size="159x78">
													<widget size="159x29" type="*widget.menuItem">
														<text pos="8,4" size="143x21">subsubitem A (long)</text>
													</widget>
													<widget pos="0,33" size="159x29" type="*widget.menuItem">
														<text pos="8,4" size="98x21">subsubitem B</text>
													</widget>
												</container>
											</widget>
										</widget>
									</widget>
								</widget>
							</widget>
						</widget>
					</overlay>
				</canvas>
			`,
		},
		"no space on left & right side for submenu": {
			windowSize: fyne.NewSize(200, 300),
			menuPos:    fyne.NewPos(10, 10),
			tapPositions: []fyne.Position{
				fyne.NewPos(30, 100),  // open submenu
				fyne.NewPos(100, 170), // open subsubmenu
			},
			want: `
				<canvas size="200x300">
					<content>
						<rectangle size="200x300"/>
					</content>
					<overlay>
						<widget size="200x300" type="*widget.OverlayContainer">
							<widget pos="10,10" size="71x108" type="*widget.Menu">
								<widget size="71x108" type="*widget.Shadow">
									<radialGradient centerOffset="0.5,0.5" pos="-4,-4" size="4x4" startColor="shadow"/>
									<linearGradient endColor="shadow" pos="0,-4" size="71x4"/>
									<radialGradient centerOffset="-0.5,0.5" pos="71,-4" size="4x4" startColor="shadow"/>
									<linearGradient angle="270" pos="71,0" size="4x108" startColor="shadow"/>
									<radialGradient centerOffset="-0.5,-0.5" pos="71,108" size="4x4" startColor="shadow"/>
									<linearGradient pos="0,108" size="71x4" startColor="shadow"/>
									<radialGradient centerOffset="0.5,-0.5" pos="-4,108" size="4x4" startColor="shadow"/>
									<linearGradient angle="270" endColor="shadow" pos="-4,0" size="4x108"/>
								</widget>
								<widget size="71x108" type="*widget.ScrollContainer">
									<widget size="71x108" type="*widget.menuBox">
										<container pos="0,4" size="71x116">
											<widget size="71x29" type="*widget.menuItem">
												<text pos="8,4" size="10x21">A</text>
											</widget>
											<widget pos="0,33" size="71x1" type="*widget.Separator">
												<rectangle fillColor="disabled text" size="71x1"/>
											</widget>
											<widget pos="0,38" size="71x29" type="*widget.menuItem">
												<text pos="8,4" size="55x21">B (long)</text>
											</widget>
											<widget pos="0,71" size="71x29" type="*widget.menuItem">
												<text pos="8,4" size="10x21">C</text>
												<image pos="51,4" rsc="menuExpandIcon" size="iconInlineSize"/>
											</widget>
										</container>
									</widget>
								</widget>
								<widget pos="37,71" size="153x103" type="*widget.Menu">
									<widget size="153x103" type="*widget.Shadow">
										<radialGradient centerOffset="0.5,0.5" pos="-4,-4" size="4x4" startColor="shadow"/>
										<linearGradient endColor="shadow" pos="0,-4" size="153x4"/>
										<radialGradient centerOffset="-0.5,0.5" pos="153,-4" size="4x4" startColor="shadow"/>
										<linearGradient angle="270" pos="153,0" size="4x103" startColor="shadow"/>
										<radialGradient centerOffset="-0.5,-0.5" pos="153,103" size="4x4" startColor="shadow"/>
										<linearGradient pos="0,103" size="153x4" startColor="shadow"/>
										<radialGradient centerOffset="0.5,-0.5" pos="-4,103" size="4x4" startColor="shadow"/>
										<linearGradient angle="270" endColor="shadow" pos="-4,0" size="4x103"/>
									</widget>
									<widget size="153x103" type="*widget.ScrollContainer">
										<widget size="153x103" type="*widget.menuBox">
											<container pos="0,4" size="153x111">
												<widget size="153x29" type="*widget.menuItem">
													<text pos="8,4" size="72x21">subitem A</text>
												</widget>
												<widget pos="0,33" size="153x29" type="*widget.menuItem">
													<text pos="8,4" size="72x21">subitem B</text>
												</widget>
												<widget pos="0,66" size="153x29" type="*widget.menuItem">
													<text pos="8,4" size="117x21">subitem C (long)</text>
													<image pos="133,4" rsc="menuExpandIcon" size="iconInlineSize"/>
												</widget>
											</container>
										</widget>
									</widget>
									<widget pos="-6,66" size="159x70" type="*widget.Menu">
										<widget size="159x70" type="*widget.Shadow">
											<radialGradient centerOffset="0.5,0.5" pos="-4,-4" size="4x4" startColor="shadow"/>
											<linearGradient endColor="shadow" pos="0,-4" size="159x4"/>
											<radialGradient centerOffset="-0.5,0.5" pos="159,-4" size="4x4" startColor="shadow"/>
											<linearGradient angle="270" pos="159,0" size="4x70" startColor="shadow"/>
											<radialGradient centerOffset="-0.5,-0.5" pos="159,70" size="4x4" startColor="shadow"/>
											<linearGradient pos="0,70" size="159x4" startColor="shadow"/>
											<radialGradient centerOffset="0.5,-0.5" pos="-4,70" size="4x4" startColor="shadow"/>
											<linearGradient angle="270" endColor="shadow" pos="-4,0" size="4x70"/>
										</widget>
										<widget size="159x70" type="*widget.ScrollContainer">
											<widget size="159x70" type="*widget.menuBox">
												<container pos="0,4" size="159x78">
													<widget size="159x29" type="*widget.menuItem">
														<text pos="8,4" size="143x21">subsubitem A (long)</text>
													</widget>
													<widget pos="0,33" size="159x29" type="*widget.menuItem">
														<text pos="8,4" size="98x21">subsubitem B</text>
													</widget>
												</container>
											</widget>
										</widget>
									</widget>
								</widget>
							</widget>
						</widget>
					</overlay>
				</canvas>
			`,
		},
		"window too short for submenu": {
			windowSize: fyne.NewSize(500, 150),
			menuPos:    fyne.NewPos(10, 10),
			tapPositions: []fyne.Position{
				fyne.NewPos(30, 100),  // open submenu
				fyne.NewPos(100, 130), // open subsubmenu
			},
			want: `
				<canvas size="500x150">
					<content>
						<rectangle size="500x150"/>
					</content>
					<overlay>
						<widget size="500x150" type="*widget.OverlayContainer">
							<widget pos="10,10" size="71x108" type="*widget.Menu">
								<widget size="71x108" type="*widget.Shadow">
									<radialGradient centerOffset="0.5,0.5" pos="-4,-4" size="4x4" startColor="shadow"/>
									<linearGradient endColor="shadow" pos="0,-4" size="71x4"/>
									<radialGradient centerOffset="-0.5,0.5" pos="71,-4" size="4x4" startColor="shadow"/>
									<linearGradient angle="270" pos="71,0" size="4x108" startColor="shadow"/>
									<radialGradient centerOffset="-0.5,-0.5" pos="71,108" size="4x4" startColor="shadow"/>
									<linearGradient pos="0,108" size="71x4" startColor="shadow"/>
									<radialGradient centerOffset="0.5,-0.5" pos="-4,108" size="4x4" startColor="shadow"/>
									<linearGradient angle="270" endColor="shadow" pos="-4,0" size="4x108"/>
								</widget>
								<widget size="71x108" type="*widget.ScrollContainer">
									<widget size="71x108" type="*widget.menuBox">
										<container pos="0,4" size="71x116">
											<widget size="71x29" type="*widget.menuItem">
												<text pos="8,4" size="10x21">A</text>
											</widget>
											<widget pos="0,33" size="71x1" type="*widget.Separator">
												<rectangle fillColor="disabled text" size="71x1"/>
											</widget>
											<widget pos="0,38" size="71x29" type="*widget.menuItem">
												<text pos="8,4" size="55x21">B (long)</text>
											</widget>
											<widget pos="0,71" size="71x29" type="*widget.menuItem">
												<text pos="8,4" size="10x21">C</text>
												<image pos="51,4" rsc="menuExpandIcon" size="iconInlineSize"/>
											</widget>
										</container>
									</widget>
								</widget>
								<widget pos="71,37" size="153x103" type="*widget.Menu">
									<widget size="153x103" type="*widget.Shadow">
										<radialGradient centerOffset="0.5,0.5" pos="-4,-4" size="4x4" startColor="shadow"/>
										<linearGradient endColor="shadow" pos="0,-4" size="153x4"/>
										<radialGradient centerOffset="-0.5,0.5" pos="153,-4" size="4x4" startColor="shadow"/>
										<linearGradient angle="270" pos="153,0" size="4x103" startColor="shadow"/>
										<radialGradient centerOffset="-0.5,-0.5" pos="153,103" size="4x4" startColor="shadow"/>
										<linearGradient pos="0,103" size="153x4" startColor="shadow"/>
										<radialGradient centerOffset="0.5,-0.5" pos="-4,103" size="4x4" startColor="shadow"/>
										<linearGradient angle="270" endColor="shadow" pos="-4,0" size="4x103"/>
									</widget>
									<widget size="153x103" type="*widget.ScrollContainer">
										<widget size="153x103" type="*widget.menuBox">
											<container pos="0,4" size="153x111">
												<widget size="153x29" type="*widget.menuItem">
													<text pos="8,4" size="72x21">subitem A</text>
												</widget>
												<widget pos="0,33" size="153x29" type="*widget.menuItem">
													<text pos="8,4" size="72x21">subitem B</text>
												</widget>
												<widget pos="0,66" size="153x29" type="*widget.menuItem">
													<text pos="8,4" size="117x21">subitem C (long)</text>
													<image pos="133,4" rsc="menuExpandIcon" size="iconInlineSize"/>
												</widget>
											</container>
										</widget>
									</widget>
									<widget pos="153,33" size="159x70" type="*widget.Menu">
										<widget size="159x70" type="*widget.Shadow">
											<radialGradient centerOffset="0.5,0.5" pos="-4,-4" size="4x4" startColor="shadow"/>
											<linearGradient endColor="shadow" pos="0,-4" size="159x4"/>
											<radialGradient centerOffset="-0.5,0.5" pos="159,-4" size="4x4" startColor="shadow"/>
											<linearGradient angle="270" pos="159,0" size="4x70" startColor="shadow"/>
											<radialGradient centerOffset="-0.5,-0.5" pos="159,70" size="4x4" startColor="shadow"/>
											<linearGradient pos="0,70" size="159x4" startColor="shadow"/>
											<radialGradient centerOffset="0.5,-0.5" pos="-4,70" size="4x4" startColor="shadow"/>
											<linearGradient angle="270" endColor="shadow" pos="-4,0" size="4x70"/>
										</widget>
										<widget size="159x70" type="*widget.ScrollContainer">
											<widget size="159x70" type="*widget.menuBox">
												<container pos="0,4" size="159x78">
													<widget size="159x29" type="*widget.menuItem">
														<text pos="8,4" size="143x21">subsubitem A (long)</text>
													</widget>
													<widget pos="0,33" size="159x29" type="*widget.menuItem">
														<text pos="8,4" size="98x21">subsubitem B</text>
													</widget>
												</container>
											</widget>
										</widget>
									</widget>
								</widget>
							</widget>
						</widget>
					</overlay>
				</canvas>
			`,
		},
		"theme change": {
			windowSize:   fyne.NewSize(500, 300),
			menuPos:      fyne.NewPos(10, 10),
			useTestTheme: true,
			want: `
				<canvas size="500x300">
					<content>
						<rectangle size="500x300"/>
					</content>
					<overlay>
						<widget size="500x300" type="*widget.OverlayContainer">
							<widget pos="10,10" size="71x108" type="*widget.Menu">
								<widget size="71x108" type="*widget.Shadow">
									<radialGradient centerOffset="0.5,0.5" pos="-4,-4" size="4x4" startColor="shadow"/>
									<linearGradient endColor="shadow" pos="0,-4" size="71x4"/>
									<radialGradient centerOffset="-0.5,0.5" pos="71,-4" size="4x4" startColor="shadow"/>
									<linearGradient angle="270" pos="71,0" size="4x108" startColor="shadow"/>
									<radialGradient centerOffset="-0.5,-0.5" pos="71,108" size="4x4" startColor="shadow"/>
									<linearGradient pos="0,108" size="71x4" startColor="shadow"/>
									<radialGradient centerOffset="0.5,-0.5" pos="-4,108" size="4x4" startColor="shadow"/>
									<linearGradient angle="270" endColor="shadow" pos="-4,0" size="4x108"/>
								</widget>
								<widget size="71x108" type="*widget.ScrollContainer">
									<widget size="71x108" type="*widget.menuBox">
										<container pos="0,4" size="71x116">
											<widget size="71x29" type="*widget.menuItem">
												<text pos="8,4" size="10x21">A</text>
											</widget>
											<widget pos="0,33" size="71x1" type="*widget.Separator">
												<rectangle fillColor="disabled text" size="71x1"/>
											</widget>
											<widget pos="0,38" size="71x29" type="*widget.menuItem">
												<text pos="8,4" size="55x21">B (long)</text>
											</widget>
											<widget pos="0,71" size="71x29" type="*widget.menuItem">
												<text pos="8,4" size="10x21">C</text>
												<image pos="51,4" rsc="menuExpandIcon" size="iconInlineSize"/>
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
		"window too short for menu": {
			windowSize: fyne.NewSize(100, 50),
			menuPos:    fyne.NewPos(10, 10),
			want: `
				<canvas size="100x50">
					<content>
						<rectangle size="100x50"/>
					</content>
					<overlay>
						<widget size="100x50" type="*widget.OverlayContainer">
							<widget pos="10,10" size="71x40" type="*widget.Menu">
								<widget size="71x40" type="*widget.Shadow">
									<radialGradient centerOffset="0.5,0.5" pos="-4,-4" size="4x4" startColor="shadow"/>
									<linearGradient endColor="shadow" pos="0,-4" size="71x4"/>
									<radialGradient centerOffset="-0.5,0.5" pos="71,-4" size="4x4" startColor="shadow"/>
									<linearGradient angle="270" pos="71,0" size="4x40" startColor="shadow"/>
									<radialGradient centerOffset="-0.5,-0.5" pos="71,40" size="4x4" startColor="shadow"/>
									<linearGradient pos="0,40" size="71x4" startColor="shadow"/>
									<radialGradient centerOffset="0.5,-0.5" pos="-4,40" size="4x4" startColor="shadow"/>
									<linearGradient angle="270" endColor="shadow" pos="-4,0" size="4x40"/>
								</widget>
								<widget size="71x40" type="*widget.ScrollContainer">
									<widget size="71x108" type="*widget.menuBox">
										<container pos="0,4" size="71x116">
											<widget size="71x29" type="*widget.menuItem">
												<text pos="8,4" size="10x21">A</text>
											</widget>
											<widget pos="0,33" size="71x1" type="*widget.Separator">
												<rectangle fillColor="disabled text" size="71x1"/>
											</widget>
											<widget pos="0,38" size="71x29" type="*widget.menuItem">
												<text pos="8,4" size="55x21">B (long)</text>
											</widget>
											<widget pos="0,71" size="71x29" type="*widget.menuItem">
												<text pos="8,4" size="10x21">C</text>
												<image pos="51,4" rsc="menuExpandIcon" size="iconInlineSize"/>
											</widget>
										</container>
									</widget>
									<widget pos="65,0" size="6x40" type="*widget.scrollBarArea">
										<widget backgroundColor="scrollbar" pos="3,0" size="3x16" type="*widget.scrollBar">
										</widget>
									</widget>
									<widget pos="0,40" size="71x0" type="*widget.Shadow">
										<linearGradient endColor="shadow" pos="0,-8" size="71x8"/>
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
			w.Resize(tt.windowSize)
			m := widget.NewMenu(menu)
			o := internalWidget.NewOverlayContainer(m, c, nil)
			c.Overlays().Add(o)
			defer c.Overlays().Remove(o)
			m.Move(tt.menuPos)
			m.Resize(m.MinSize())
			for _, pos := range tt.tapPositions {
				test.TapCanvas(c, pos)
			}
			test.AssertRendersToMarkup(t, tt.want, w.Canvas())
			if tt.useTestTheme {
				test.AssertImageMatches(t, "menu/layout_normal.png", c.Capture())
				test.WithTestTheme(t, func() {
					test.AssertImageMatches(t, "menu/layout_theme_changed.png", c.Capture())
				})
			}
		})
	}
}

func TestMenu_Dragging(t *testing.T) {
	test.NewApp()
	defer test.NewApp()

	w := test.NewWindow(canvas.NewRectangle(color.Transparent))
	defer w.Close()
	w.SetPadded(false)
	c := w.Canvas()

	menu := fyne.NewMenu("",
		fyne.NewMenuItem("A", nil),
		fyne.NewMenuItem("B", nil),
		fyne.NewMenuItem("C", nil),
		fyne.NewMenuItem("D", nil),
		fyne.NewMenuItem("E", nil),
		fyne.NewMenuItem("F", nil),
	)

	w.Resize(fyne.NewSize(100, 100))
	m := widget.NewMenu(menu)
	o := internalWidget.NewOverlayContainer(m, c, nil)
	c.Overlays().Add(o)
	defer c.Overlays().Remove(o)
	m.Move(fyne.NewPos(10, 10))
	m.Resize(m.MinSize())
	maxDragDistance := m.MinSize().Height - 90
	test.AssertRendersToMarkup(t, `
		<canvas size="100x100">
			<content>
				<rectangle size="100x100"/>
			</content>
			<overlay>
				<widget size="100x100" type="*widget.OverlayContainer">
					<widget pos="10,10" size="28x90" type="*widget.Menu">
						<widget size="28x90" type="*widget.Shadow">
							<radialGradient centerOffset="0.5,0.5" pos="-4,-4" size="4x4" startColor="shadow"/>
							<linearGradient endColor="shadow" pos="0,-4" size="28x4"/>
							<radialGradient centerOffset="-0.5,0.5" pos="28,-4" size="4x4" startColor="shadow"/>
							<linearGradient angle="270" pos="28,0" size="4x90" startColor="shadow"/>
							<radialGradient centerOffset="-0.5,-0.5" pos="28,90" size="4x4" startColor="shadow"/>
							<linearGradient pos="0,90" size="28x4" startColor="shadow"/>
							<radialGradient centerOffset="0.5,-0.5" pos="-4,90" size="4x4" startColor="shadow"/>
							<linearGradient angle="270" endColor="shadow" pos="-4,0" size="4x90"/>
						</widget>
						<widget size="28x90" type="*widget.ScrollContainer">
							<widget size="28x202" type="*widget.menuBox">
								<container pos="0,4" size="28x210">
									<widget size="28x29" type="*widget.menuItem">
										<text pos="8,4" size="10x21">A</text>
									</widget>
									<widget pos="0,33" size="28x29" type="*widget.menuItem">
										<text pos="8,4" size="10x21">B</text>
									</widget>
									<widget pos="0,66" size="28x29" type="*widget.menuItem">
										<text pos="8,4" size="10x21">C</text>
									</widget>
									<widget pos="0,99" size="28x29" type="*widget.menuItem">
										<text pos="8,4" size="12x21">D</text>
									</widget>
									<widget pos="0,132" size="28x29" type="*widget.menuItem">
										<text pos="8,4" size="9x21">E</text>
									</widget>
									<widget pos="0,165" size="28x29" type="*widget.menuItem">
										<text pos="8,4" size="8x21">F</text>
									</widget>
								</container>
							</widget>
							<widget pos="22,0" size="6x90" type="*widget.scrollBarArea">
								<widget backgroundColor="scrollbar" pos="3,0" size="3x40" type="*widget.scrollBar">
								</widget>
							</widget>
							<widget pos="0,90" size="28x0" type="*widget.Shadow">
								<linearGradient endColor="shadow" pos="0,-8" size="28x8"/>
							</widget>
						</widget>
					</widget>
				</widget>
			</overlay>
		</canvas>
	`, w.Canvas())

	test.Drag(c, fyne.NewPos(20, 20), 0, -50)
	test.AssertRendersToMarkup(t, `
		<canvas size="100x100">
			<content>
				<rectangle size="100x100"/>
			</content>
			<overlay>
				<widget size="100x100" type="*widget.OverlayContainer">
					<widget pos="10,10" size="28x90" type="*widget.Menu">
						<widget size="28x90" type="*widget.Shadow">
							<radialGradient centerOffset="0.5,0.5" pos="-4,-4" size="4x4" startColor="shadow"/>
							<linearGradient endColor="shadow" pos="0,-4" size="28x4"/>
							<radialGradient centerOffset="-0.5,0.5" pos="28,-4" size="4x4" startColor="shadow"/>
							<linearGradient angle="270" pos="28,0" size="4x90" startColor="shadow"/>
							<radialGradient centerOffset="-0.5,-0.5" pos="28,90" size="4x4" startColor="shadow"/>
							<linearGradient pos="0,90" size="28x4" startColor="shadow"/>
							<radialGradient centerOffset="0.5,-0.5" pos="-4,90" size="4x4" startColor="shadow"/>
							<linearGradient angle="270" endColor="shadow" pos="-4,0" size="4x90"/>
						</widget>
						<widget size="28x90" type="*widget.ScrollContainer">
							<widget pos="0,-50" size="28x202" type="*widget.menuBox">
								<container pos="0,4" size="28x210">
									<widget size="28x29" type="*widget.menuItem">
										<text pos="8,4" size="10x21">A</text>
									</widget>
									<widget pos="0,33" size="28x29" type="*widget.menuItem">
										<text pos="8,4" size="10x21">B</text>
									</widget>
									<widget pos="0,66" size="28x29" type="*widget.menuItem">
										<text pos="8,4" size="10x21">C</text>
									</widget>
									<widget pos="0,99" size="28x29" type="*widget.menuItem">
										<text pos="8,4" size="12x21">D</text>
									</widget>
									<widget pos="0,132" size="28x29" type="*widget.menuItem">
										<text pos="8,4" size="9x21">E</text>
									</widget>
									<widget pos="0,165" size="28x29" type="*widget.menuItem">
										<text pos="8,4" size="8x21">F</text>
									</widget>
								</container>
							</widget>
							<widget pos="22,0" size="6x90" type="*widget.scrollBarArea">
								<widget backgroundColor="scrollbar" pos="3,22" size="3x40" type="*widget.scrollBar">
								</widget>
							</widget>
							<widget size="28x0" type="*widget.Shadow">
								<linearGradient size="28x8" startColor="shadow"/>
							</widget>
							<widget pos="0,90" size="28x0" type="*widget.Shadow">
								<linearGradient endColor="shadow" pos="0,-8" size="28x8"/>
							</widget>
						</widget>
					</widget>
				</widget>
			</overlay>
		</canvas>
	`, w.Canvas())

	test.Drag(c, fyne.NewPos(20, 20), 0, -maxDragDistance)
	test.AssertRendersToMarkup(t, `
		<canvas size="100x100">
			<content>
				<rectangle size="100x100"/>
			</content>
			<overlay>
				<widget size="100x100" type="*widget.OverlayContainer">
					<widget pos="10,10" size="28x90" type="*widget.Menu">
						<widget size="28x90" type="*widget.Shadow">
							<radialGradient centerOffset="0.5,0.5" pos="-4,-4" size="4x4" startColor="shadow"/>
							<linearGradient endColor="shadow" pos="0,-4" size="28x4"/>
							<radialGradient centerOffset="-0.5,0.5" pos="28,-4" size="4x4" startColor="shadow"/>
							<linearGradient angle="270" pos="28,0" size="4x90" startColor="shadow"/>
							<radialGradient centerOffset="-0.5,-0.5" pos="28,90" size="4x4" startColor="shadow"/>
							<linearGradient pos="0,90" size="28x4" startColor="shadow"/>
							<radialGradient centerOffset="0.5,-0.5" pos="-4,90" size="4x4" startColor="shadow"/>
							<linearGradient angle="270" endColor="shadow" pos="-4,0" size="4x90"/>
						</widget>
						<widget size="28x90" type="*widget.ScrollContainer">
							<widget pos="0,-112" size="28x202" type="*widget.menuBox">
								<container pos="0,4" size="28x210">
									<widget size="28x29" type="*widget.menuItem">
										<text pos="8,4" size="10x21">A</text>
									</widget>
									<widget pos="0,33" size="28x29" type="*widget.menuItem">
										<text pos="8,4" size="10x21">B</text>
									</widget>
									<widget pos="0,66" size="28x29" type="*widget.menuItem">
										<text pos="8,4" size="10x21">C</text>
									</widget>
									<widget pos="0,99" size="28x29" type="*widget.menuItem">
										<text pos="8,4" size="12x21">D</text>
									</widget>
									<widget pos="0,132" size="28x29" type="*widget.menuItem">
										<text pos="8,4" size="9x21">E</text>
									</widget>
									<widget pos="0,165" size="28x29" type="*widget.menuItem">
										<text pos="8,4" size="8x21">F</text>
									</widget>
								</container>
							</widget>
							<widget pos="22,0" size="6x90" type="*widget.scrollBarArea">
								<widget backgroundColor="scrollbar" pos="3,50" size="3x40" type="*widget.scrollBar">
								</widget>
							</widget>
							<widget size="28x0" type="*widget.Shadow">
								<linearGradient size="28x8" startColor="shadow"/>
							</widget>
						</widget>
					</widget>
				</widget>
			</overlay>
		</canvas>
	`, w.Canvas())

	test.Drag(c, fyne.NewPos(20, 20), 0, maxDragDistance-50)
	test.AssertRendersToMarkup(t, `
		<canvas size="100x100">
			<content>
				<rectangle size="100x100"/>
			</content>
			<overlay>
				<widget size="100x100" type="*widget.OverlayContainer">
					<widget pos="10,10" size="28x90" type="*widget.Menu">
						<widget size="28x90" type="*widget.Shadow">
							<radialGradient centerOffset="0.5,0.5" pos="-4,-4" size="4x4" startColor="shadow"/>
							<linearGradient endColor="shadow" pos="0,-4" size="28x4"/>
							<radialGradient centerOffset="-0.5,0.5" pos="28,-4" size="4x4" startColor="shadow"/>
							<linearGradient angle="270" pos="28,0" size="4x90" startColor="shadow"/>
							<radialGradient centerOffset="-0.5,-0.5" pos="28,90" size="4x4" startColor="shadow"/>
							<linearGradient pos="0,90" size="28x4" startColor="shadow"/>
							<radialGradient centerOffset="0.5,-0.5" pos="-4,90" size="4x4" startColor="shadow"/>
							<linearGradient angle="270" endColor="shadow" pos="-4,0" size="4x90"/>
						</widget>
						<widget size="28x90" type="*widget.ScrollContainer">
							<widget pos="0,-50" size="28x202" type="*widget.menuBox">
								<container pos="0,4" size="28x210">
									<widget size="28x29" type="*widget.menuItem">
										<text pos="8,4" size="10x21">A</text>
									</widget>
									<widget pos="0,33" size="28x29" type="*widget.menuItem">
										<text pos="8,4" size="10x21">B</text>
									</widget>
									<widget pos="0,66" size="28x29" type="*widget.menuItem">
										<text pos="8,4" size="10x21">C</text>
									</widget>
									<widget pos="0,99" size="28x29" type="*widget.menuItem">
										<text pos="8,4" size="12x21">D</text>
									</widget>
									<widget pos="0,132" size="28x29" type="*widget.menuItem">
										<text pos="8,4" size="9x21">E</text>
									</widget>
									<widget pos="0,165" size="28x29" type="*widget.menuItem">
										<text pos="8,4" size="8x21">F</text>
									</widget>
								</container>
							</widget>
							<widget pos="22,0" size="6x90" type="*widget.scrollBarArea">
								<widget backgroundColor="scrollbar" pos="3,22" size="3x40" type="*widget.scrollBar">
								</widget>
							</widget>
							<widget size="28x0" type="*widget.Shadow">
								<linearGradient size="28x8" startColor="shadow"/>
							</widget>
							<widget pos="0,90" size="28x0" type="*widget.Shadow">
								<linearGradient endColor="shadow" pos="0,-8" size="28x8"/>
							</widget>
						</widget>
					</widget>
				</widget>
			</overlay>
		</canvas>
	`, w.Canvas())

	test.Drag(c, fyne.NewPos(20, 20), 0, 50)
	test.AssertRendersToMarkup(t, `
		<canvas size="100x100">
			<content>
				<rectangle size="100x100"/>
			</content>
			<overlay>
				<widget size="100x100" type="*widget.OverlayContainer">
					<widget pos="10,10" size="28x90" type="*widget.Menu">
						<widget size="28x90" type="*widget.Shadow">
							<radialGradient centerOffset="0.5,0.5" pos="-4,-4" size="4x4" startColor="shadow"/>
							<linearGradient endColor="shadow" pos="0,-4" size="28x4"/>
							<radialGradient centerOffset="-0.5,0.5" pos="28,-4" size="4x4" startColor="shadow"/>
							<linearGradient angle="270" pos="28,0" size="4x90" startColor="shadow"/>
							<radialGradient centerOffset="-0.5,-0.5" pos="28,90" size="4x4" startColor="shadow"/>
							<linearGradient pos="0,90" size="28x4" startColor="shadow"/>
							<radialGradient centerOffset="0.5,-0.5" pos="-4,90" size="4x4" startColor="shadow"/>
							<linearGradient angle="270" endColor="shadow" pos="-4,0" size="4x90"/>
						</widget>
						<widget size="28x90" type="*widget.ScrollContainer">
							<widget size="28x202" type="*widget.menuBox">
								<container pos="0,4" size="28x210">
									<widget size="28x29" type="*widget.menuItem">
										<text pos="8,4" size="10x21">A</text>
									</widget>
									<widget pos="0,33" size="28x29" type="*widget.menuItem">
										<text pos="8,4" size="10x21">B</text>
									</widget>
									<widget pos="0,66" size="28x29" type="*widget.menuItem">
										<text pos="8,4" size="10x21">C</text>
									</widget>
									<widget pos="0,99" size="28x29" type="*widget.menuItem">
										<text pos="8,4" size="12x21">D</text>
									</widget>
									<widget pos="0,132" size="28x29" type="*widget.menuItem">
										<text pos="8,4" size="9x21">E</text>
									</widget>
									<widget pos="0,165" size="28x29" type="*widget.menuItem">
										<text pos="8,4" size="8x21">F</text>
									</widget>
								</container>
							</widget>
							<widget pos="22,0" size="6x90" type="*widget.scrollBarArea">
								<widget backgroundColor="scrollbar" pos="3,0" size="3x40" type="*widget.scrollBar">
								</widget>
							</widget>
							<widget pos="0,90" size="28x0" type="*widget.Shadow">
								<linearGradient endColor="shadow" pos="0,-8" size="28x8"/>
							</widget>
						</widget>
					</widget>
				</widget>
			</overlay>
		</canvas>
	`, w.Canvas())
}
