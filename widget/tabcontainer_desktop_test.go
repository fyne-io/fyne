// +build !mobile

package widget_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/test"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
)

func TestTabContainer_ApplyTheme(t *testing.T) {
	test.NewApp()
	defer test.NewApp()

	w := test.NewWindow(
		widget.NewTabContainer(&widget.TabItem{Text: "Test", Content: widget.NewLabel("Text")}),
	)
	defer w.Close()
	w.SetPadded(false)
	w.Resize(fyne.NewSize(150, 150))
	c := w.Canvas()

	test.AssertImageMatches(t, "tabcontainer/desktop/single_initial.png", c.Capture())

	test.ApplyTheme(t, test.NewTheme())
	test.AssertImageMatches(t, "tabcontainer/desktop/single_custom_theme.png", c.Capture())
}

func TestTabContainer_ChangeItemContent(t *testing.T) {
	test.NewApp()
	defer test.NewApp()

	item1 := &widget.TabItem{Text: "Test1", Content: widget.NewLabel("Text1")}
	item2 := &widget.TabItem{Text: "Test2", Content: widget.NewLabel("Text2")}
	tabs := widget.NewTabContainer(item1, item2)
	w := test.NewWindow(tabs)
	defer w.Close()
	w.SetPadded(false)
	w.Resize(fyne.NewSize(150, 150))
	c := w.Canvas()

	text1Visible := `
		<canvas size="150x150">
			<content>
				<widget size="150x150" type="*widget.TabContainer">
					<widget pos="0,33" size="150x117" type="*widget.Label">
						<text pos="4,4" size="142x21">Text1</text>
					</widget>
					<container size="150x29">
						<widget size="57x29" type="*widget.tabButton">
							<text alignment="center" bold color="focus" pos="8,4" size="41x21">Test1</text>
						</widget>
						<widget pos="61,0" size="57x29" type="*widget.tabButton">
							<text alignment="center" bold pos="8,4" size="41x21">Test2</text>
						</widget>
					</container>
					<rectangle fillColor="shadow" pos="0,29" size="150x4"/>
					<rectangle fillColor="focus" pos="0,29" size="57x4"/>
				</widget>
			</content>
		</canvas>
	`
	text3Visible := `
		<canvas size="150x150">
			<content>
				<widget size="150x150" type="*widget.TabContainer">
					<widget pos="0,33" size="150x117" type="*widget.Label">
						<text pos="4,4" size="142x21">Text3</text>
					</widget>
					<container size="150x29">
						<widget size="57x29" type="*widget.tabButton">
							<text alignment="center" bold color="focus" pos="8,4" size="41x21">Test1</text>
						</widget>
						<widget pos="61,0" size="57x29" type="*widget.tabButton">
							<text alignment="center" bold pos="8,4" size="41x21">Test2</text>
						</widget>
					</container>
					<rectangle fillColor="shadow" pos="0,29" size="150x4"/>
					<rectangle fillColor="focus" pos="0,29" size="57x4"/>
				</widget>
			</content>
		</canvas>
	`

	test.AssertRendersToMarkup(t, text1Visible, c)

	item1.Content = widget.NewLabel("Text3")
	tabs.Refresh()
	test.AssertRendersToMarkup(t, text3Visible, c)

	item2.Content = widget.NewLabel("Text4")
	tabs.Refresh()
	test.AssertRendersToMarkup(t, text3Visible, c)
}

func TestTabContainer_ChangeItemIcon(t *testing.T) {
	test.NewApp()
	defer test.NewApp()

	item1 := &widget.TabItem{Icon: theme.CancelIcon(), Content: widget.NewLabel("Text1")}
	item2 := &widget.TabItem{Icon: theme.ConfirmIcon(), Content: widget.NewLabel("Text2")}
	tabs := widget.NewTabContainer(item1, item2)
	w := test.NewWindow(tabs)
	defer w.Close()
	w.SetPadded(false)
	w.Resize(fyne.NewSize(150, 150))
	c := w.Canvas()

	test.AssertRendersToMarkup(t, `
		<canvas size="150x150">
			<content>
				<widget size="150x150" type="*widget.TabContainer">
					<widget pos="0,33" size="150x117" type="*widget.Label">
						<text pos="4,4" size="142x21">Text1</text>
					</widget>
					<container size="150x29">
						<widget size="28x29" type="*widget.tabButton">
							<image pos="4,4" rsc="cancelIcon" size="iconInlineSize" themed="primary"/>
						</widget>
						<widget pos="32,0" size="28x29" type="*widget.tabButton">
							<image pos="4,4" rsc="confirmIcon" size="iconInlineSize"/>
						</widget>
					</container>
					<rectangle fillColor="shadow" pos="0,29" size="150x4"/>
					<rectangle fillColor="focus" pos="0,29" size="28x4"/>
				</widget>
			</content>
		</canvas>
	`, c)

	item1.Icon = theme.InfoIcon()
	tabs.Refresh()
	test.AssertRendersToMarkup(t, `
		<canvas size="150x150">
			<content>
				<widget size="150x150" type="*widget.TabContainer">
					<widget pos="0,33" size="150x117" type="*widget.Label">
						<text pos="4,4" size="142x21">Text1</text>
					</widget>
					<container size="150x29">
						<widget size="28x29" type="*widget.tabButton">
							<image pos="4,4" rsc="infoIcon" size="iconInlineSize" themed="primary"/>
						</widget>
						<widget pos="32,0" size="28x29" type="*widget.tabButton">
							<image pos="4,4" rsc="confirmIcon" size="iconInlineSize"/>
						</widget>
					</container>
					<rectangle fillColor="shadow" pos="0,29" size="150x4"/>
					<rectangle fillColor="focus" pos="0,29" size="28x4"/>
				</widget>
			</content>
		</canvas>
	`, c)

	item2.Icon = theme.ContentAddIcon()
	tabs.Refresh()
	test.AssertRendersToMarkup(t, `
		<canvas size="150x150">
			<content>
				<widget size="150x150" type="*widget.TabContainer">
					<widget pos="0,33" size="150x117" type="*widget.Label">
						<text pos="4,4" size="142x21">Text1</text>
					</widget>
					<container size="150x29">
						<widget size="28x29" type="*widget.tabButton">
							<image pos="4,4" rsc="infoIcon" size="iconInlineSize" themed="primary"/>
						</widget>
						<widget pos="32,0" size="28x29" type="*widget.tabButton">
							<image pos="4,4" rsc="contentAddIcon" size="iconInlineSize"/>
						</widget>
					</container>
					<rectangle fillColor="shadow" pos="0,29" size="150x4"/>
					<rectangle fillColor="focus" pos="0,29" size="28x4"/>
				</widget>
			</content>
		</canvas>
	`, c)
}

func TestTabContainer_ChangeItemText(t *testing.T) {
	test.NewApp()
	defer test.NewApp()

	item1 := &widget.TabItem{Text: "Test1", Content: widget.NewLabel("Text1")}
	item2 := &widget.TabItem{Text: "Test2", Content: widget.NewLabel("Text2")}
	tabs := widget.NewTabContainer(item1, item2)
	w := test.NewWindow(tabs)
	defer w.Close()
	w.SetPadded(false)
	w.Resize(fyne.NewSize(150, 150))
	c := w.Canvas()

	test.AssertRendersToMarkup(t, `
		<canvas size="150x150">
			<content>
				<widget size="150x150" type="*widget.TabContainer">
					<widget pos="0,33" size="150x117" type="*widget.Label">
						<text pos="4,4" size="142x21">Text1</text>
					</widget>
					<container size="150x29">
						<widget size="57x29" type="*widget.tabButton">
							<text alignment="center" bold color="focus" pos="8,4" size="41x21">Test1</text>
						</widget>
						<widget pos="61,0" size="57x29" type="*widget.tabButton">
							<text alignment="center" bold pos="8,4" size="41x21">Test2</text>
						</widget>
					</container>
					<rectangle fillColor="shadow" pos="0,29" size="150x4"/>
					<rectangle fillColor="focus" pos="0,29" size="57x4"/>
				</widget>
			</content>
		</canvas>
	`, c)

	item1.Text = "New 1"
	tabs.Refresh()
	test.AssertRendersToMarkup(t, `
		<canvas size="150x150">
			<content>
				<widget size="150x150" type="*widget.TabContainer">
					<widget pos="0,33" size="150x117" type="*widget.Label">
						<text pos="4,4" size="142x21">Text1</text>
					</widget>
					<container size="150x29">
						<widget size="63x29" type="*widget.tabButton">
							<text alignment="center" bold color="focus" pos="8,4" size="47x21">New 1</text>
						</widget>
						<widget pos="67,0" size="57x29" type="*widget.tabButton">
							<text alignment="center" bold pos="8,4" size="41x21">Test2</text>
						</widget>
					</container>
					<rectangle fillColor="shadow" pos="0,29" size="150x4"/>
					<rectangle fillColor="focus" pos="0,29" size="63x4"/>
				</widget>
			</content>
		</canvas>
	`, c)

	item2.Text = "New 2"
	tabs.Refresh()
	test.AssertRendersToMarkup(t, `
		<canvas size="150x150">
			<content>
				<widget size="150x150" type="*widget.TabContainer">
					<widget pos="0,33" size="150x117" type="*widget.Label">
						<text pos="4,4" size="142x21">Text1</text>
					</widget>
					<container size="150x29">
						<widget size="63x29" type="*widget.tabButton">
							<text alignment="center" bold color="focus" pos="8,4" size="47x21">New 1</text>
						</widget>
						<widget pos="67,0" size="63x29" type="*widget.tabButton">
							<text alignment="center" bold pos="8,4" size="47x21">New 2</text>
						</widget>
					</container>
					<rectangle fillColor="shadow" pos="0,29" size="150x4"/>
					<rectangle fillColor="focus" pos="0,29" size="63x4"/>
				</widget>
			</content>
		</canvas>
	`, c)
}

func TestTabContainer_DynamicTabs(t *testing.T) {
	test.NewApp()
	defer test.NewApp()

	item1 := &widget.TabItem{Text: "Test1", Content: widget.NewLabel("Text 1")}
	tabs := widget.NewTabContainer(item1)
	w := test.NewWindow(tabs)
	defer w.Close()
	w.SetPadded(false)
	w.Resize(fyne.NewSize(300, 150))
	c := w.Canvas()

	test.AssertRendersToMarkup(t, `
		<canvas size="300x150">
			<content>
				<widget size="300x150" type="*widget.TabContainer">
					<widget pos="0,33" size="300x117" type="*widget.Label">
						<text pos="4,4" size="292x21">Text 1</text>
					</widget>
					<container size="300x29">
						<widget size="57x29" type="*widget.tabButton">
							<text alignment="center" bold color="focus" pos="8,4" size="41x21">Test1</text>
						</widget>
					</container>
					<rectangle fillColor="shadow" pos="0,29" size="300x4"/>
					<rectangle fillColor="focus" pos="0,29" size="57x4"/>
				</widget>
			</content>
		</canvas>
	`, c)

	appendedItem := widget.NewTabItem("Test2", widget.NewLabel("Text 2"))
	tabs.Append(appendedItem)
	assert.Equal(t, 2, len(tabs.Items))
	assert.Equal(t, "Test2", tabs.Items[1].Text)
	test.AssertRendersToMarkup(t, `
		<canvas size="300x150">
			<content>
				<widget size="300x150" type="*widget.TabContainer">
					<widget pos="0,33" size="300x117" type="*widget.Label">
						<text pos="4,4" size="292x21">Text 1</text>
					</widget>
					<container size="300x29">
						<widget size="57x29" type="*widget.tabButton">
							<text alignment="center" bold color="focus" pos="8,4" size="41x21">Test1</text>
						</widget>
						<widget pos="61,0" size="57x29" type="*widget.tabButton">
							<text alignment="center" bold pos="8,4" size="41x21">Test2</text>
						</widget>
					</container>
					<rectangle fillColor="shadow" pos="0,29" size="300x4"/>
					<rectangle fillColor="focus" pos="0,29" size="57x4"/>
				</widget>
			</content>
		</canvas>
	`, c)

	tabs.RemoveIndex(1)
	assert.Equal(t, len(tabs.Items), 1)
	assert.Equal(t, "Test1", tabs.Items[0].Text)
	test.AssertRendersToMarkup(t, `
		<canvas size="300x150">
			<content>
				<widget size="300x150" type="*widget.TabContainer">
					<widget pos="0,33" size="300x117" type="*widget.Label">
						<text pos="4,4" size="292x21">Text 1</text>
					</widget>
					<container size="300x29">
						<widget size="57x29" type="*widget.tabButton">
							<text alignment="center" bold color="focus" pos="8,4" size="41x21">Test1</text>
						</widget>
					</container>
					<rectangle fillColor="shadow" pos="0,29" size="300x4"/>
					<rectangle fillColor="focus" pos="0,29" size="57x4"/>
				</widget>
			</content>
		</canvas>
	`, c)

	tabs.Append(appendedItem)
	tabs.Remove(tabs.Items[0])
	assert.Equal(t, len(tabs.Items), 1)
	assert.Equal(t, "Test2", tabs.Items[0].Text)
	test.AssertRendersToMarkup(t, `
		<canvas size="300x150">
			<content>
				<widget size="300x150" type="*widget.TabContainer">
					<widget pos="0,33" size="300x117" type="*widget.Label">
						<text pos="4,4" size="292x21">Text 2</text>
					</widget>
					<container size="300x29">
						<widget size="57x29" type="*widget.tabButton">
							<text alignment="center" bold color="focus" pos="8,4" size="41x21">Test2</text>
						</widget>
					</container>
					<rectangle fillColor="shadow" pos="0,29" size="300x4"/>
					<rectangle fillColor="focus" pos="0,29" size="57x4"/>
				</widget>
			</content>
		</canvas>
	`, c)

	tabs.Append(widget.NewTabItem("Test3", canvas.NewCircle(theme.BackgroundColor())))
	tabs.Append(widget.NewTabItem("Test4", canvas.NewCircle(theme.BackgroundColor())))
	tabs.Append(widget.NewTabItem("Test5", canvas.NewCircle(theme.BackgroundColor())))
	assert.Equal(t, 4, len(tabs.Items))
	assert.Equal(t, "Test3", tabs.Items[1].Text)
	assert.Equal(t, "Test4", tabs.Items[2].Text)
	assert.Equal(t, "Test5", tabs.Items[3].Text)
	test.AssertRendersToMarkup(t, `
		<canvas size="300x150">
			<content>
				<widget size="300x150" type="*widget.TabContainer">
					<widget pos="0,33" size="300x117" type="*widget.Label">
						<text pos="4,4" size="292x21">Text 2</text>
					</widget>
					<container size="300x29">
						<widget size="57x29" type="*widget.tabButton">
							<text alignment="center" bold color="focus" pos="8,4" size="41x21">Test2</text>
						</widget>
						<widget pos="61,0" size="57x29" type="*widget.tabButton">
							<text alignment="center" bold pos="8,4" size="41x21">Test3</text>
						</widget>
						<widget pos="122,0" size="57x29" type="*widget.tabButton">
							<text alignment="center" bold pos="8,4" size="41x21">Test4</text>
						</widget>
						<widget pos="183,0" size="57x29" type="*widget.tabButton">
							<text alignment="center" bold pos="8,4" size="41x21">Test5</text>
						</widget>
					</container>
					<rectangle fillColor="shadow" pos="0,29" size="300x4"/>
					<rectangle fillColor="focus" pos="0,29" size="57x4"/>
				</widget>
			</content>
		</canvas>
	`, c)

	tabs.SetItems([]*widget.TabItem{
		widget.NewTabItem("Test6", widget.NewLabel("Text 6")),
		widget.NewTabItem("Test7", widget.NewLabel("Text 7")),
		widget.NewTabItem("Test8", widget.NewLabel("Text 8")),
	})
	assert.Equal(t, 3, len(tabs.Items))
	assert.Equal(t, "Test6", tabs.Items[0].Text)
	assert.Equal(t, "Test7", tabs.Items[1].Text)
	assert.Equal(t, "Test8", tabs.Items[2].Text)
	test.AssertRendersToMarkup(t, `
		<canvas size="300x150">
			<content>
				<widget size="300x150" type="*widget.TabContainer">
					<widget pos="0,33" size="300x117" type="*widget.Label">
						<text pos="4,4" size="292x21">Text 6</text>
					</widget>
					<container size="300x29">
						<widget size="57x29" type="*widget.tabButton">
							<text alignment="center" bold color="focus" pos="8,4" size="41x21">Test6</text>
						</widget>
						<widget pos="61,0" size="57x29" type="*widget.tabButton">
							<text alignment="center" bold pos="8,4" size="41x21">Test7</text>
						</widget>
						<widget pos="122,0" size="57x29" type="*widget.tabButton">
							<text alignment="center" bold pos="8,4" size="41x21">Test8</text>
						</widget>
					</container>
					<rectangle fillColor="shadow" pos="0,29" size="300x4"/>
					<rectangle fillColor="focus" pos="0,29" size="57x4"/>
				</widget>
			</content>
		</canvas>
	`, c)
}

func TestTabContainer_HoverButtons(t *testing.T) {
	test.NewApp()
	defer test.NewApp()

	item1 := &widget.TabItem{Text: "Test1", Content: widget.NewLabel("Text1")}
	item2 := &widget.TabItem{Text: "Test2", Content: widget.NewLabel("Text2")}
	tabs := widget.NewTabContainer(item1, item2)
	w := test.NewWindow(tabs)
	defer w.Close()
	w.SetPadded(false)
	w.Resize(fyne.NewSize(150, 150))
	c := w.Canvas()

	noneHovered := `
		<canvas size="150x150">
			<content>
				<widget size="150x150" type="*widget.TabContainer">
					<widget pos="0,33" size="150x117" type="*widget.Label">
						<text pos="4,4" size="142x21">Text1</text>
					</widget>
					<container size="150x29">
						<widget size="57x29" type="*widget.tabButton">
							<text alignment="center" bold color="focus" pos="8,4" size="41x21">Test1</text>
						</widget>
						<widget pos="61,0" size="57x29" type="*widget.tabButton">
							<text alignment="center" bold pos="8,4" size="41x21">Test2</text>
						</widget>
					</container>
					<rectangle fillColor="shadow" pos="0,29" size="150x4"/>
					<rectangle fillColor="focus" pos="0,29" size="57x4"/>
				</widget>
			</content>
		</canvas>
	`
	firstHovered := `
		<canvas size="150x150">
			<content>
				<widget size="150x150" type="*widget.TabContainer">
					<widget pos="0,33" size="150x117" type="*widget.Label">
						<text pos="4,4" size="142x21">Text1</text>
					</widget>
					<container size="150x29">
						<widget backgroundColor="hover" size="57x29" type="*widget.tabButton">
							<text alignment="center" bold color="focus" pos="8,4" size="41x21">Test1</text>
						</widget>
						<widget pos="61,0" size="57x29" type="*widget.tabButton">
							<text alignment="center" bold pos="8,4" size="41x21">Test2</text>
						</widget>
					</container>
					<rectangle fillColor="shadow" pos="0,29" size="150x4"/>
					<rectangle fillColor="focus" pos="0,29" size="57x4"/>
				</widget>
			</content>
		</canvas>
	`
	secondHovered := `
		<canvas size="150x150">
			<content>
				<widget size="150x150" type="*widget.TabContainer">
					<widget pos="0,33" size="150x117" type="*widget.Label">
						<text pos="4,4" size="142x21">Text1</text>
					</widget>
					<container size="150x29">
						<widget size="57x29" type="*widget.tabButton">
							<text alignment="center" bold color="focus" pos="8,4" size="41x21">Test1</text>
						</widget>
						<widget backgroundColor="hover" pos="61,0" size="57x29" type="*widget.tabButton">
							<text alignment="center" bold pos="8,4" size="41x21">Test2</text>
						</widget>
					</container>
					<rectangle fillColor="shadow" pos="0,29" size="150x4"/>
					<rectangle fillColor="focus" pos="0,29" size="57x4"/>
				</widget>
			</content>
		</canvas>
	`
	test.AssertRendersToMarkup(t, noneHovered, c)

	test.MoveMouse(c, fyne.NewPos(10, 10))
	test.AssertRendersToMarkup(t, firstHovered, c)

	test.MoveMouse(c, fyne.NewPos(75, 10))
	test.AssertRendersToMarkup(t, secondHovered, c)

	test.MoveMouse(c, fyne.NewPos(10, 10))
	test.AssertRendersToMarkup(t, firstHovered, c)
}

func TestTabContainer_Layout(t *testing.T) {
	test.NewApp()
	defer test.NewApp()

	w := test.NewWindow(nil)
	defer w.Close()
	w.SetPadded(false)
	c := w.Canvas()

	tests := []struct {
		name     string
		item     *widget.TabItem
		location widget.TabLocation
		want     string
	}{
		{
			name:     "top: tab with icon and text",
			item:     widget.NewTabItemWithIcon("Text1", theme.CancelIcon(), canvas.NewCircle(theme.BackgroundColor())),
			location: widget.TabLocationTop,
			want: `
				<canvas size="150x150">
					<content>
						<widget size="150x150" type="*widget.TabContainer">
							<circle fillColor="background" pos="0,33" size="150x117"/>
							<container size="150x29">
								<widget size="82x29" type="*widget.tabButton">
									<text alignment="center" bold color="focus" pos="32,4" size="42x21">Text1</text>
									<image pos="8,4" rsc="cancelIcon" size="iconInlineSize" themed="primary"/>
								</widget>
							</container>
							<rectangle fillColor="shadow" pos="0,29" size="150x4"/>
							<rectangle fillColor="focus" pos="0,29" size="82x4"/>
						</widget>
					</content>
				</canvas>
			`,
		},
		{
			name:     "top: tab with text only",
			item:     widget.NewTabItem("Text2", canvas.NewCircle(theme.BackgroundColor())),
			location: widget.TabLocationTop,
			want: `
				<canvas size="150x150">
					<content>
						<widget size="150x150" type="*widget.TabContainer">
							<circle fillColor="background" pos="0,33" size="150x117"/>
							<container size="150x29">
								<widget size="58x29" type="*widget.tabButton">
									<text alignment="center" bold color="focus" pos="8,4" size="42x21">Text2</text>
								</widget>
							</container>
							<rectangle fillColor="shadow" pos="0,29" size="150x4"/>
							<rectangle fillColor="focus" pos="0,29" size="58x4"/>
						</widget>
					</content>
				</canvas>
			`,
		},
		{
			name:     "top: tab with icon only",
			item:     widget.NewTabItemWithIcon("", theme.InfoIcon(), canvas.NewCircle(theme.BackgroundColor())),
			location: widget.TabLocationTop,
			want: `
				<canvas size="150x150">
					<content>
						<widget size="150x150" type="*widget.TabContainer">
							<circle fillColor="background" pos="0,33" size="150x117"/>
							<container size="150x29">
								<widget size="28x29" type="*widget.tabButton">
									<image pos="4,4" rsc="infoIcon" size="iconInlineSize" themed="primary"/>
								</widget>
							</container>
							<rectangle fillColor="shadow" pos="0,29" size="150x4"/>
							<rectangle fillColor="focus" pos="0,29" size="28x4"/>
						</widget>
					</content>
				</canvas>
			`,
		},
		{
			name:     "bottom: tab with icon and text",
			item:     widget.NewTabItemWithIcon("Text1", theme.CancelIcon(), canvas.NewCircle(theme.BackgroundColor())),
			location: widget.TabLocationBottom,
			want: `
				<canvas size="150x150">
					<content>
						<widget size="150x150" type="*widget.TabContainer">
							<circle fillColor="background" size="150x117"/>
							<container pos="0,121" size="150x29">
								<widget size="82x29" type="*widget.tabButton">
									<text alignment="center" bold color="focus" pos="32,4" size="42x21">Text1</text>
									<image pos="8,4" rsc="cancelIcon" size="iconInlineSize" themed="primary"/>
								</widget>
							</container>
							<rectangle fillColor="shadow" pos="0,117" size="150x4"/>
							<rectangle fillColor="focus" pos="0,117" size="82x4"/>
						</widget>
					</content>
				</canvas>
			`,
		},
		{
			name:     "bottom: tab with text only",
			item:     widget.NewTabItem("Text2", canvas.NewCircle(theme.BackgroundColor())),
			location: widget.TabLocationBottom,
			want: `
				<canvas size="150x150">
					<content>
						<widget size="150x150" type="*widget.TabContainer">
							<circle fillColor="background" size="150x117"/>
							<container pos="0,121" size="150x29">
								<widget size="58x29" type="*widget.tabButton">
									<text alignment="center" bold color="focus" pos="8,4" size="42x21">Text2</text>
								</widget>
							</container>
							<rectangle fillColor="shadow" pos="0,117" size="150x4"/>
							<rectangle fillColor="focus" pos="0,117" size="58x4"/>
						</widget>
					</content>
				</canvas>
			`,
		},
		{
			name:     "bottom: tab with icon only",
			item:     widget.NewTabItemWithIcon("", theme.InfoIcon(), canvas.NewCircle(theme.BackgroundColor())),
			location: widget.TabLocationBottom,
			want: `
				<canvas size="150x150">
					<content>
						<widget size="150x150" type="*widget.TabContainer">
							<circle fillColor="background" size="150x117"/>
							<container pos="0,121" size="150x29">
								<widget size="28x29" type="*widget.tabButton">
									<image pos="4,4" rsc="infoIcon" size="iconInlineSize" themed="primary"/>
								</widget>
							</container>
							<rectangle fillColor="shadow" pos="0,117" size="150x4"/>
							<rectangle fillColor="focus" pos="0,117" size="28x4"/>
						</widget>
					</content>
				</canvas>
			`,
		},
		{
			name:     "leading: tab with icon and text",
			item:     widget.NewTabItemWithIcon("Text1", theme.CancelIcon(), canvas.NewCircle(theme.BackgroundColor())),
			location: widget.TabLocationLeading,
			want: `
				<canvas size="150x150">
					<content>
						<widget size="150x150" type="*widget.TabContainer">
							<circle fillColor="background" pos="54,0" size="96x150"/>
							<container size="50x150">
								<widget size="50x73" type="*widget.tabButton">
									<text alignment="center" bold color="focus" pos="4,48" size="42x21">Text1</text>
									<image pos="5,4" rsc="cancelIcon" size="40x40" themed="primary"/>
								</widget>
							</container>
							<rectangle fillColor="shadow" pos="50,0" size="4x150"/>
							<rectangle fillColor="focus" pos="50,0" size="4x73"/>
						</widget>
					</content>
				</canvas>
			`,
		},
		{
			name:     "leading: tab with text only",
			item:     widget.NewTabItem("Text2", canvas.NewCircle(theme.BackgroundColor())),
			location: widget.TabLocationLeading,
			want: `
				<canvas size="150x150">
					<content>
						<widget size="150x150" type="*widget.TabContainer">
							<circle fillColor="background" pos="54,0" size="96x150"/>
							<container size="50x150">
								<widget size="50x29" type="*widget.tabButton">
									<text alignment="center" bold color="focus" pos="4,4" size="42x21">Text2</text>
								</widget>
							</container>
							<rectangle fillColor="shadow" pos="50,0" size="4x150"/>
							<rectangle fillColor="focus" pos="50,0" size="4x29"/>
						</widget>
					</content>
				</canvas>
			`,
		},
		{
			name:     "leading: tab with icon only",
			item:     widget.NewTabItemWithIcon("", theme.InfoIcon(), canvas.NewCircle(theme.BackgroundColor())),
			location: widget.TabLocationLeading,
			want: `
				<canvas size="150x150">
					<content>
						<widget size="150x150" type="*widget.TabContainer">
							<circle fillColor="background" pos="52,0" size="98x150"/>
							<container size="48x150">
								<widget size="48x48" type="*widget.tabButton">
									<image pos="4,4" rsc="infoIcon" size="40x40" themed="primary"/>
								</widget>
							</container>
							<rectangle fillColor="shadow" pos="48,0" size="4x150"/>
							<rectangle fillColor="focus" pos="48,0" size="4x48"/>
						</widget>
					</content>
				</canvas>
			`,
		},
		{
			name:     "trailing: tab with icon and text",
			item:     widget.NewTabItemWithIcon("Text1", theme.CancelIcon(), canvas.NewCircle(theme.BackgroundColor())),
			location: widget.TabLocationTrailing,
			want: `
				<canvas size="150x150">
					<content>
						<widget size="150x150" type="*widget.TabContainer">
							<circle fillColor="background" size="96x150"/>
							<container pos="100,0" size="50x150">
								<widget size="50x73" type="*widget.tabButton">
									<text alignment="center" bold color="focus" pos="4,48" size="42x21">Text1</text>
									<image pos="5,4" rsc="cancelIcon" size="40x40" themed="primary"/>
								</widget>
							</container>
							<rectangle fillColor="shadow" pos="96,0" size="4x150"/>
							<rectangle fillColor="focus" pos="96,0" size="4x73"/>
						</widget>
					</content>
				</canvas>
			`,
		},
		{
			name:     "trailing: tab with text only",
			item:     widget.NewTabItem("Text2", canvas.NewCircle(theme.BackgroundColor())),
			location: widget.TabLocationTrailing,
			want: `
				<canvas size="150x150">
					<content>
						<widget size="150x150" type="*widget.TabContainer">
							<circle fillColor="background" size="96x150"/>
							<container pos="100,0" size="50x150">
								<widget size="50x29" type="*widget.tabButton">
									<text alignment="center" bold color="focus" pos="4,4" size="42x21">Text2</text>
								</widget>
							</container>
							<rectangle fillColor="shadow" pos="96,0" size="4x150"/>
							<rectangle fillColor="focus" pos="96,0" size="4x29"/>
						</widget>
					</content>
				</canvas>
			`,
		},
		{
			name:     "trailing: tab with icon only",
			item:     widget.NewTabItemWithIcon("", theme.InfoIcon(), canvas.NewCircle(theme.BackgroundColor())),
			location: widget.TabLocationTrailing,
			want: `
				<canvas size="150x150">
					<content>
						<widget size="150x150" type="*widget.TabContainer">
							<circle fillColor="background" size="98x150"/>
							<container pos="102,0" size="48x150">
								<widget size="48x48" type="*widget.tabButton">
									<image pos="4,4" rsc="infoIcon" size="40x40" themed="primary"/>
								</widget>
							</container>
							<rectangle fillColor="shadow" pos="98,0" size="4x150"/>
							<rectangle fillColor="focus" pos="98,0" size="4x48"/>
						</widget>
					</content>
				</canvas>
			`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tabs := widget.NewTabContainer(tt.item)
			tabs.SetTabLocation(tt.location)
			w.SetContent(tabs)
			w.Resize(fyne.NewSize(150, 150))

			test.AssertRendersToMarkup(t, tt.want, c)
		})
	}
}

func TestTabContainer_SetTabLocation(t *testing.T) {
	test.NewApp()
	defer test.NewApp()

	item1 := &widget.TabItem{Text: "Test1", Content: widget.NewLabel("Text 1")}
	item2 := &widget.TabItem{Text: "Test2", Content: widget.NewLabel("Text 2")}
	item3 := &widget.TabItem{Text: "Test3", Content: widget.NewLabel("Text 3")}
	tabs := widget.NewTabContainer(item1, item2, item3)
	w := test.NewWindow(tabs)
	defer w.Close()
	w.SetPadded(false)
	c := w.Canvas()

	tabsTop := `
		<canvas size="179x62">
			<content>
				<widget size="179x62" type="*widget.TabContainer">
					<widget pos="0,33" size="179x29" type="*widget.Label">
						<text pos="4,4" size="171x21">Text 1</text>
					</widget>
					<container size="179x29">
						<widget size="57x29" type="*widget.tabButton">
							<text alignment="center" bold color="focus" pos="8,4" size="41x21">Test1</text>
						</widget>
						<widget pos="61,0" size="57x29" type="*widget.tabButton">
							<text alignment="center" bold pos="8,4" size="41x21">Test2</text>
						</widget>
						<widget pos="122,0" size="57x29" type="*widget.tabButton">
							<text alignment="center" bold pos="8,4" size="41x21">Test3</text>
						</widget>
					</container>
					<rectangle fillColor="shadow" pos="0,29" size="179x4"/>
					<rectangle fillColor="focus" pos="0,29" size="57x4"/>
				</widget>
			</content>
		</canvas>
	`
	tabsLeading := `
		<canvas size="105x95">
			<content>
				<widget size="105x95" type="*widget.TabContainer">
					<widget pos="53,0" size="52x95" type="*widget.Label">
						<text pos="4,4" size="44x21">Text 1</text>
					</widget>
					<container size="49x95">
						<widget size="49x29" type="*widget.tabButton">
							<text alignment="center" bold color="focus" pos="4,4" size="41x21">Test1</text>
						</widget>
						<widget pos="0,33" size="49x29" type="*widget.tabButton">
							<text alignment="center" bold pos="4,4" size="41x21">Test2</text>
						</widget>
						<widget pos="0,66" size="49x29" type="*widget.tabButton">
							<text alignment="center" bold pos="4,4" size="41x21">Test3</text>
						</widget>
					</container>
					<rectangle fillColor="shadow" pos="49,0" size="4x95"/>
					<rectangle fillColor="focus" pos="49,0" size="4x29"/>
				</widget>
			</content>
		</canvas>
	`
	tabsBottom := `
		<canvas size="179x62">
			<content>
				<widget size="179x62" type="*widget.TabContainer">
					<widget size="179x29" type="*widget.Label">
						<text pos="4,4" size="171x21">Text 1</text>
					</widget>
					<container pos="0,33" size="179x29">
						<widget size="57x29" type="*widget.tabButton">
							<text alignment="center" bold color="focus" pos="8,4" size="41x21">Test1</text>
						</widget>
						<widget pos="61,0" size="57x29" type="*widget.tabButton">
							<text alignment="center" bold pos="8,4" size="41x21">Test2</text>
						</widget>
						<widget pos="122,0" size="57x29" type="*widget.tabButton">
							<text alignment="center" bold pos="8,4" size="41x21">Test3</text>
						</widget>
					</container>
					<rectangle fillColor="shadow" pos="0,29" size="179x4"/>
					<rectangle fillColor="focus" pos="0,29" size="57x4"/>
				</widget>
			</content>
		</canvas>
	`
	tabsTrailing := `
		<canvas size="105x95">
			<content>
				<widget size="105x95" type="*widget.TabContainer">
					<widget size="52x95" type="*widget.Label">
						<text pos="4,4" size="44x21">Text 1</text>
					</widget>
					<container pos="56,0" size="49x95">
						<widget size="49x29" type="*widget.tabButton">
							<text alignment="center" bold color="focus" pos="4,4" size="41x21">Test1</text>
						</widget>
						<widget pos="0,33" size="49x29" type="*widget.tabButton">
							<text alignment="center" bold pos="4,4" size="41x21">Test2</text>
						</widget>
						<widget pos="0,66" size="49x29" type="*widget.tabButton">
							<text alignment="center" bold pos="4,4" size="41x21">Test3</text>
						</widget>
					</container>
					<rectangle fillColor="shadow" pos="52,0" size="4x95"/>
					<rectangle fillColor="focus" pos="52,0" size="4x29"/>
				</widget>
			</content>
		</canvas>
	`
	w.Resize(tabs.MinSize())
	test.AssertRendersToMarkup(t, tabsTop, c)

	tabs.SetTabLocation(widget.TabLocationLeading)
	w.Resize(tabs.MinSize())
	test.AssertRendersToMarkup(t, tabsLeading, c)

	tabs.SetTabLocation(widget.TabLocationBottom)
	w.Resize(tabs.MinSize())
	test.AssertRendersToMarkup(t, tabsBottom, c)

	tabs.SetTabLocation(widget.TabLocationTrailing)
	w.Resize(tabs.MinSize())
	test.AssertRendersToMarkup(t, tabsTrailing, c)

	tabs.SetTabLocation(widget.TabLocationTop)
	w.Resize(tabs.MinSize())
	test.AssertRendersToMarkup(t, tabsTop, c)
}

func TestTabContainer_Tapped(t *testing.T) {
	test.NewApp()
	defer test.NewApp()

	item1 := &widget.TabItem{Text: "Test1", Content: widget.NewLabel("Text 1")}
	item2 := &widget.TabItem{Text: "Test2", Content: widget.NewLabel("Text 2")}
	item3 := &widget.TabItem{Text: "Test3", Content: widget.NewLabel("Text 3")}
	tabs := widget.NewTabContainer(item1, item2, item3)
	w := test.NewWindow(tabs)
	defer w.Close()
	w.SetPadded(false)
	w.Resize(fyne.NewSize(200, 100))
	c := w.Canvas()

	firstSelected := `
		<canvas size="200x100">
			<content>
				<widget size="200x100" type="*widget.TabContainer">
					<widget pos="0,33" size="200x67" type="*widget.Label">
						<text pos="4,4" size="192x21">Text 1</text>
					</widget>
					<container size="200x29">
						<widget size="57x29" type="*widget.tabButton">
							<text alignment="center" bold color="focus" pos="8,4" size="41x21">Test1</text>
						</widget>
						<widget pos="61,0" size="57x29" type="*widget.tabButton">
							<text alignment="center" bold pos="8,4" size="41x21">Test2</text>
						</widget>
						<widget pos="122,0" size="57x29" type="*widget.tabButton">
							<text alignment="center" bold pos="8,4" size="41x21">Test3</text>
						</widget>
					</container>
					<rectangle fillColor="shadow" pos="0,29" size="200x4"/>
					<rectangle fillColor="focus" pos="0,29" size="57x4"/>
				</widget>
			</content>
		</canvas>
	`
	secondSelected := `
		<canvas size="200x100">
			<content>
				<widget size="200x100" type="*widget.TabContainer">
					<widget pos="0,33" size="200x67" type="*widget.Label">
						<text pos="4,4" size="192x21">Text 2</text>
					</widget>
					<container size="200x29">
						<widget size="57x29" type="*widget.tabButton">
							<text alignment="center" bold pos="8,4" size="41x21">Test1</text>
						</widget>
						<widget pos="61,0" size="57x29" type="*widget.tabButton">
							<text alignment="center" bold color="focus" pos="8,4" size="41x21">Test2</text>
						</widget>
						<widget pos="122,0" size="57x29" type="*widget.tabButton">
							<text alignment="center" bold pos="8,4" size="41x21">Test3</text>
						</widget>
					</container>
					<rectangle fillColor="shadow" pos="0,29" size="200x4"/>
					<rectangle fillColor="focus" pos="61,29" size="57x4"/>
				</widget>
			</content>
		</canvas>
	`
	thirdSelected := `
		<canvas size="200x100">
			<content>
				<widget size="200x100" type="*widget.TabContainer">
					<widget pos="0,33" size="200x67" type="*widget.Label">
						<text pos="4,4" size="192x21">Text 3</text>
					</widget>
					<container size="200x29">
						<widget size="57x29" type="*widget.tabButton">
							<text alignment="center" bold pos="8,4" size="41x21">Test1</text>
						</widget>
						<widget pos="61,0" size="57x29" type="*widget.tabButton">
							<text alignment="center" bold pos="8,4" size="41x21">Test2</text>
						</widget>
						<widget pos="122,0" size="57x29" type="*widget.tabButton">
							<text alignment="center" bold color="focus" pos="8,4" size="41x21">Test3</text>
						</widget>
					</container>
					<rectangle fillColor="shadow" pos="0,29" size="200x4"/>
					<rectangle fillColor="focus" pos="122,29" size="57x4"/>
				</widget>
			</content>
		</canvas>
	`
	require.Equal(t, 0, tabs.CurrentTabIndex())
	test.AssertRendersToMarkup(t, firstSelected, c)

	test.TapCanvas(c, fyne.NewPos(75, 10))
	assert.Equal(t, 1, tabs.CurrentTabIndex())
	test.AssertRendersToMarkup(t, secondSelected, c)

	test.TapCanvas(c, fyne.NewPos(150, 10))
	assert.Equal(t, 2, tabs.CurrentTabIndex())
	test.AssertRendersToMarkup(t, thirdSelected, c)

	test.TapCanvas(c, fyne.NewPos(10, 10))
	require.Equal(t, 0, tabs.CurrentTabIndex())
	test.AssertRendersToMarkup(t, firstSelected, c)
}
