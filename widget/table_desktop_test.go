// +build !mobile

package widget

import (
	"fmt"
	"testing"

	"fyne.io/fyne"
	"fyne.io/fyne/test"

	"github.com/stretchr/testify/assert"
)

func TestTable_Hovered(t *testing.T) {
	test.NewApp()
	defer test.NewApp()

	table := NewTable(
		func() (int, int) { return 2, 2 },
		func() fyne.CanvasObject {
			return NewLabel("placeholder")
		},
		func(id TableCellID, c fyne.CanvasObject) {
			c.(*Label).SetText(fmt.Sprintf("Cell %d, %d", id.Row, id.Col))
		})

	w := test.NewWindow(table)
	defer w.Close()
	w.Resize(fyne.NewSize(180, 180))

	test.MoveMouse(w.Canvas(), fyne.NewPos(35, 50))
	test.MoveMouse(w.Canvas(), fyne.NewPos(35, 100))

	assert.Nil(t, table.hoveredCell)

	test.AssertRendersToMarkup(t, `
		<canvas padded size="180x180">
			<content>
				<widget pos="4,4" size="172x172" type="*widget.Table">
					<widget pos="4,4" size="168x168" type="*widget.Scroll">
						<widget size="203x168" type="*widget.tableCells">
							<widget pos="4,4" size="93x29" type="*widget.Label">
								<text pos="4,4" size="85x21">Cell 0, 0</text>
							</widget>
							<widget pos="106,4" size="93x29" type="*widget.Label">
								<text pos="4,4" size="85x21">Cell 0, 1</text>
							</widget>
							<widget pos="4,42" size="93x29" type="*widget.Label">
								<text pos="4,4" size="85x21">Cell 1, 0</text>
							</widget>
							<widget pos="106,42" size="93x29" type="*widget.Label">
								<text pos="4,4" size="85x21">Cell 1, 1</text>
							</widget>
						</widget>
						<widget pos="0,162" size="168x6" type="*widget.scrollBarArea">
							<widget backgroundColor="scrollbar" pos="0,3" size="139x3" type="*widget.scrollBar">
							</widget>
						</widget>
						<widget pos="168,0" size="0x168" type="*widget.Shadow">
							<linearGradient angle="270" endColor="shadow" pos="-8,0" size="8x168"/>
						</widget>
					</widget>
					<widget pos="105,4" size="1x168" type="*widget.Separator">
						<rectangle fillColor="disabled" size="1x168"/>
					</widget>
					<widget pos="4,41" size="168x1" type="*widget.Separator">
						<rectangle fillColor="disabled" size="168x1"/>
					</widget>
					<widget size="0x0" type="*widget.Separator">
						<rectangle fillColor="disabled" size="0x0"/>
					</widget>
					<widget size="0x0" type="*widget.Separator">
						<rectangle fillColor="disabled" size="0x0"/>
					</widget>
					<widget size="0x0" type="*widget.Separator">
						<rectangle fillColor="disabled" size="0x0"/>
					</widget>
				</widget>
			</content>
		</canvas>
	`, w.Canvas())

	table.Length = func() (int, int) { return 3, 5 }
	table.Refresh()

	w.SetContent(table)
	w.Resize(fyne.NewSize(180, 180))
	test.MoveMouse(w.Canvas(), fyne.NewPos(35, 50))

	assert.Equal(t, 0, table.hoveredCell.Col)
	assert.Equal(t, 1, table.hoveredCell.Row)

	test.AssertRendersToMarkup(t, `
		<canvas padded size="180x180">
			<content>
				<widget pos="4,4" size="172x172" type="*widget.Table">
					<widget pos="4,4" size="168x168" type="*widget.Scroll">
						<widget size="509x168" type="*widget.tableCells">
							<widget pos="4,4" size="93x29" type="*widget.Label">
								<text pos="4,4" size="85x21">Cell 0, 0</text>
							</widget>
							<widget pos="106,4" size="93x29" type="*widget.Label">
								<text pos="4,4" size="85x21">Cell 0, 1</text>
							</widget>
							<widget pos="4,42" size="93x29" type="*widget.Label">
								<text pos="4,4" size="85x21">Cell 1, 0</text>
							</widget>
							<widget pos="106,42" size="93x29" type="*widget.Label">
								<text pos="4,4" size="85x21">Cell 1, 1</text>
							</widget>
							<widget pos="4,80" size="93x29" type="*widget.Label">
								<text pos="4,4" size="85x21">Cell 2, 0</text>
							</widget>
							<widget pos="106,80" size="93x29" type="*widget.Label">
								<text pos="4,4" size="85x21">Cell 2, 1</text>
							</widget>
						</widget>
						<widget pos="0,162" size="168x6" type="*widget.scrollBarArea">
							<widget backgroundColor="scrollbar" pos="0,3" size="55x3" type="*widget.scrollBar">
							</widget>
						</widget>
						<widget pos="168,0" size="0x168" type="*widget.Shadow">
							<linearGradient angle="270" endColor="shadow" pos="-8,0" size="8x168"/>
						</widget>
					</widget>
					<rectangle fillColor="hover" pos="4,0" size="101x4"/>
					<rectangle fillColor="hover" pos="0,42" size="4x37"/>
					<widget pos="105,4" size="1x168" type="*widget.Separator">
						<rectangle fillColor="disabled" size="1x168"/>
					</widget>
					<widget pos="4,41" size="168x1" type="*widget.Separator">
						<rectangle fillColor="disabled" size="168x1"/>
					</widget>
					<widget pos="4,79" size="168x1" type="*widget.Separator">
						<rectangle fillColor="disabled" size="168x1"/>
					</widget>
					<widget size="0x0" type="*widget.Separator">
						<rectangle fillColor="disabled" size="0x0"/>
					</widget>
					<widget size="0x0" type="*widget.Separator">
						<rectangle fillColor="disabled" size="0x0"/>
					</widget>
				</widget>
			</content>
		</canvas>
	`, w.Canvas())
}
