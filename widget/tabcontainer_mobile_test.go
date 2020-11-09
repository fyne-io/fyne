// +build mobile

package widget_test

import (
	"testing"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/test"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTabContainer_ApplyTheme(t *testing.T) {
	test.NewApp()
	defer test.NewApp()
	test.ApplyTheme(t, theme.LightTheme())

	w := test.NewWindow(
		widget.NewTabContainer(&widget.TabItem{Text: "Test", Content: widget.NewLabel("Text")}),
	)
	defer w.Close()
	w.SetPadded(false)
	w.Resize(fyne.NewSize(150, 150))
	c := w.Canvas()

	test.AssertImageMatches(t, "tabcontainer/mobile/single_initial.png", c.Capture())

	test.ApplyTheme(t, theme.DarkTheme())
	test.AssertImageMatches(t, "tabcontainer/mobile/single_dark.png", c.Capture())

	test.ApplyTheme(t, test.NewTheme())
	test.AssertImageMatches(t, "tabcontainer/mobile/single_custom_theme.png", c.Capture())
}

func TestTabContainer_ChangeItemContent(t *testing.T) {
	test.NewApp()
	defer test.NewApp()
	test.ApplyTheme(t, theme.LightTheme())

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
					<widget size="73x29" type="*widget.tabButton">
						<text alignment="1" bold color="focus" pos="4,4" size="65x21">Test1</text>
					</widget>
					<widget pos="77,0" size="73x29" type="*widget.tabButton">
						<text alignment="1" bold pos="4,4" size="65x21">Test2</text>
					</widget>
					<rectangle fillColor="shadow" pos="0,29" size="150x4"/>
					<rectangle fillColor="focus" pos="0,29" size="73x4"/>
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
					<widget size="73x29" type="*widget.tabButton">
						<text alignment="1" bold color="focus" pos="4,4" size="65x21">Test1</text>
					</widget>
					<widget pos="77,0" size="73x29" type="*widget.tabButton">
						<text alignment="1" bold pos="4,4" size="65x21">Test2</text>
					</widget>
					<rectangle fillColor="shadow" pos="0,29" size="150x4"/>
					<rectangle fillColor="focus" pos="0,29" size="73x4"/>
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
	test.ApplyTheme(t, theme.LightTheme())

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
					<widget pos="0,52" size="150x98" type="*widget.Label">
						<text pos="4,4" size="142x21">Text1</text>
					</widget>
					<widget size="73x48" type="*widget.tabButton">
						<image pos="16,4" rsc="cancelIcon" size="40x40" themed="primary"/>
					</widget>
					<widget pos="77,0" size="73x48" type="*widget.tabButton">
						<image pos="16,4" rsc="confirmIcon" size="40x40"/>
					</widget>
					<rectangle fillColor="shadow" pos="0,48" size="150x4"/>
					<rectangle fillColor="focus" pos="0,48" size="73x4"/>
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
					<widget pos="0,52" size="150x98" type="*widget.Label">
						<text pos="4,4" size="142x21">Text1</text>
					</widget>
					<widget size="73x48" type="*widget.tabButton">
						<image pos="16,4" rsc="infoIcon" size="40x40" themed="primary"/>
					</widget>
					<widget pos="77,0" size="73x48" type="*widget.tabButton">
						<image pos="16,4" rsc="confirmIcon" size="40x40"/>
					</widget>
					<rectangle fillColor="shadow" pos="0,48" size="150x4"/>
					<rectangle fillColor="focus" pos="0,48" size="73x4"/>
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
					<widget pos="0,52" size="150x98" type="*widget.Label">
						<text pos="4,4" size="142x21">Text1</text>
					</widget>
					<widget size="73x48" type="*widget.tabButton">
						<image pos="16,4" rsc="infoIcon" size="40x40" themed="primary"/>
					</widget>
					<widget pos="77,0" size="73x48" type="*widget.tabButton">
						<image pos="16,4" rsc="contentAddIcon" size="40x40"/>
					</widget>
					<rectangle fillColor="shadow" pos="0,48" size="150x4"/>
					<rectangle fillColor="focus" pos="0,48" size="73x4"/>
				</widget>
			</content>
		</canvas>
	`, c)
}

func TestTabContainer_ChangeItemText(t *testing.T) {
	test.NewApp()
	defer test.NewApp()
	test.ApplyTheme(t, theme.LightTheme())

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
					<widget size="73x29" type="*widget.tabButton">
						<text alignment="1" bold color="focus" pos="4,4" size="65x21">Test1</text>
					</widget>
					<widget pos="77,0" size="73x29" type="*widget.tabButton">
						<text alignment="1" bold pos="4,4" size="65x21">Test2</text>
					</widget>
					<rectangle fillColor="shadow" pos="0,29" size="150x4"/>
					<rectangle fillColor="focus" pos="0,29" size="73x4"/>
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
					<widget size="73x29" type="*widget.tabButton">
						<text alignment="1" bold color="focus" pos="4,4" size="65x21">New 1</text>
					</widget>
					<widget pos="77,0" size="73x29" type="*widget.tabButton">
						<text alignment="1" bold pos="4,4" size="65x21">Test2</text>
					</widget>
					<rectangle fillColor="shadow" pos="0,29" size="150x4"/>
					<rectangle fillColor="focus" pos="0,29" size="73x4"/>
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
					<widget size="73x29" type="*widget.tabButton">
						<text alignment="1" bold color="focus" pos="4,4" size="65x21">New 1</text>
					</widget>
					<widget pos="77,0" size="73x29" type="*widget.tabButton">
						<text alignment="1" bold pos="4,4" size="65x21">New 2</text>
					</widget>
					<rectangle fillColor="shadow" pos="0,29" size="150x4"/>
					<rectangle fillColor="focus" pos="0,29" size="73x4"/>
				</widget>
			</content>
		</canvas>
	`, c)
}

func TestTabContainer_DynamicTabs(t *testing.T) {
	test.NewApp()
	defer test.NewApp()
	test.ApplyTheme(t, theme.LightTheme())

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
					<widget size="300x29" type="*widget.tabButton">
						<text alignment="1" bold color="focus" pos="4,4" size="292x21">Test1</text>
					</widget>
					<rectangle fillColor="shadow" pos="0,29" size="300x4"/>
					<rectangle fillColor="focus" pos="0,29" size="300x4"/>
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
					<widget size="148x29" type="*widget.tabButton">
						<text alignment="1" bold color="focus" pos="4,4" size="140x21">Test1</text>
					</widget>
					<widget pos="152,0" size="148x29" type="*widget.tabButton">
						<text alignment="1" bold pos="4,4" size="140x21">Test2</text>
					</widget>
					<rectangle fillColor="shadow" pos="0,29" size="300x4"/>
					<rectangle fillColor="focus" pos="0,29" size="148x4"/>
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
					<widget size="300x29" type="*widget.tabButton">
						<text alignment="1" bold color="focus" pos="4,4" size="292x21">Test1</text>
					</widget>
					<rectangle fillColor="shadow" pos="0,29" size="300x4"/>
					<rectangle fillColor="focus" pos="0,29" size="300x4"/>
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
					<widget size="300x29" type="*widget.tabButton">
						<text alignment="1" bold color="focus" pos="4,4" size="292x21">Test2</text>
					</widget>
					<rectangle fillColor="shadow" pos="0,29" size="300x4"/>
					<rectangle fillColor="focus" pos="0,29" size="300x4"/>
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
					<widget size="72x29" type="*widget.tabButton">
						<text alignment="1" bold color="focus" pos="4,4" size="64x21">Test2</text>
					</widget>
					<widget pos="76,0" size="72x29" type="*widget.tabButton">
						<text alignment="1" bold pos="4,4" size="64x21">Test3</text>
					</widget>
					<widget pos="152,0" size="72x29" type="*widget.tabButton">
						<text alignment="1" bold pos="4,4" size="64x21">Test4</text>
					</widget>
					<widget pos="228,0" size="72x29" type="*widget.tabButton">
						<text alignment="1" bold pos="4,4" size="64x21">Test5</text>
					</widget>
					<rectangle fillColor="shadow" pos="0,29" size="300x4"/>
					<rectangle fillColor="focus" pos="0,29" size="72x4"/>
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
					<widget size="97x29" type="*widget.tabButton">
						<text alignment="1" bold color="focus" pos="4,4" size="89x21">Test6</text>
					</widget>
					<widget pos="101,0" size="98x29" type="*widget.tabButton">
						<text alignment="1" bold pos="4,4" size="90x21">Test7</text>
					</widget>
					<widget pos="203,0" size="97x29" type="*widget.tabButton">
						<text alignment="1" bold pos="4,4" size="89x21">Test8</text>
					</widget>
					<rectangle fillColor="shadow" pos="0,29" size="300x4"/>
					<rectangle fillColor="focus" pos="0,29" size="97x4"/>
				</widget>
			</content>
		</canvas>
	`, c)
}

func TestTabContainer_HoverButtons(t *testing.T) {
	test.NewApp()
	defer test.NewApp()
	test.ApplyTheme(t, theme.LightTheme())

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
					<widget size="73x29" type="*widget.tabButton">
						<text alignment="1" bold color="focus" pos="4,4" size="65x21">Test1</text>
					</widget>
					<widget pos="77,0" size="73x29" type="*widget.tabButton">
						<text alignment="1" bold pos="4,4" size="65x21">Test2</text>
					</widget>
					<rectangle fillColor="shadow" pos="0,29" size="150x4"/>
					<rectangle fillColor="focus" pos="0,29" size="73x4"/>
				</widget>
			</content>
		</canvas>
	`
	test.AssertRendersToMarkup(t, noneHovered, c)

	test.MoveMouse(c, fyne.NewPos(10, 10))
	test.AssertRendersToMarkup(t, noneHovered, c, "no hovering on mobile")

	test.MoveMouse(c, fyne.NewPos(75, 10))
	test.AssertRendersToMarkup(t, noneHovered, c, "no hovering on mobile")

	test.MoveMouse(c, fyne.NewPos(10, 10))
	test.AssertRendersToMarkup(t, noneHovered, c, "no hovering on mobile")
}

func TestTabContainer_Layout(t *testing.T) {
	test.NewApp()
	defer test.NewApp()
	test.ApplyTheme(t, theme.LightTheme())

	w := test.NewWindow(nil)
	defer w.Close()
	w.SetPadded(false)
	c := w.Canvas()

	bottomIcon := `
		<canvas size="150x150">
			<content>
				<widget size="150x150" type="*widget.TabContainer">
					<circle fillColor="background" size="150x98"/>
					<widget size="150x48" type="*widget.tabButton">
						<image pos="55,4" rsc="infoIcon" size="40x40" themed="primary"/>
					</widget>
					<rectangle fillColor="shadow" pos="0,98" size="150x4"/>
					<rectangle fillColor="focus" pos="0,98" size="150x4"/>
				</widget>
			</content>
		</canvas>
	`
	bottomIconAndText := `
		<canvas size="150x150">
			<content>
				<widget size="150x150" type="*widget.TabContainer">
					<circle fillColor="background" size="150x73"/>
					<widget size="150x73" type="*widget.tabButton">
						<text alignment="1" bold color="focus" pos="4,48" size="142x21">Text1</text>
						<image pos="55,4" rsc="cancelIcon" size="40x40" themed="primary"/>
					</widget>
					<rectangle fillColor="shadow" pos="0,73" size="150x4"/>
					<rectangle fillColor="focus" pos="0,73" size="150x4"/>
				</widget>
			</content>
		</canvas>
	`
	bottomText := `
		<canvas size="150x150">
			<content>
				<widget size="150x150" type="*widget.TabContainer">
					<circle fillColor="background" size="150x117"/>
					<widget size="150x29" type="*widget.tabButton">
						<text alignment="1" bold color="focus" pos="4,4" size="142x21">Text2</text>
					</widget>
					<rectangle fillColor="shadow" pos="0,117" size="150x4"/>
					<rectangle fillColor="focus" pos="0,117" size="150x4"/>
				</widget>
			</content>
		</canvas>
	`
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
							<circle fillColor="background" pos="0,77" size="150x73"/>
							<widget size="150x73" type="*widget.tabButton">
								<text alignment="1" bold color="focus" pos="4,48" size="142x21">Text1</text>
								<image pos="55,4" rsc="cancelIcon" size="40x40" themed="primary"/>
							</widget>
							<rectangle fillColor="shadow" pos="0,73" size="150x4"/>
							<rectangle fillColor="focus" pos="0,73" size="150x4"/>
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
							<widget size="150x29" type="*widget.tabButton">
								<text alignment="1" bold color="focus" pos="4,4" size="142x21">Text2</text>
							</widget>
							<rectangle fillColor="shadow" pos="0,29" size="150x4"/>
							<rectangle fillColor="focus" pos="0,29" size="150x4"/>
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
							<circle fillColor="background" pos="0,52" size="150x98"/>
							<widget size="150x48" type="*widget.tabButton">
								<image pos="55,4" rsc="infoIcon" size="40x40" themed="primary"/>
							</widget>
							<rectangle fillColor="shadow" pos="0,48" size="150x4"/>
							<rectangle fillColor="focus" pos="0,48" size="150x4"/>
						</widget>
					</content>
				</canvas>
			`,
		},
		{
			name:     "bottom: tab with icon and text",
			item:     widget.NewTabItemWithIcon("Text1", theme.CancelIcon(), canvas.NewCircle(theme.BackgroundColor())),
			location: widget.TabLocationBottom,
			want:     bottomIconAndText,
		},
		{
			name:     "bottom: tab with text only",
			item:     widget.NewTabItem("Text2", canvas.NewCircle(theme.BackgroundColor())),
			location: widget.TabLocationBottom,
			want:     bottomText,
		},
		{
			name:     "bottom: tab with icon only",
			item:     widget.NewTabItemWithIcon("", theme.InfoIcon(), canvas.NewCircle(theme.BackgroundColor())),
			location: widget.TabLocationBottom,
			want:     bottomIcon,
		},
		{
			name:     "leading: tab with icon and text",
			item:     widget.NewTabItemWithIcon("Text1", theme.CancelIcon(), canvas.NewCircle(theme.BackgroundColor())),
			location: widget.TabLocationLeading,
			want:     bottomIconAndText,
		},
		{
			name:     "leading: tab with text only",
			item:     widget.NewTabItem("Text2", canvas.NewCircle(theme.BackgroundColor())),
			location: widget.TabLocationLeading,
			want:     bottomText,
		},
		{
			name:     "leading: tab with icon only",
			item:     widget.NewTabItemWithIcon("", theme.InfoIcon(), canvas.NewCircle(theme.BackgroundColor())),
			location: widget.TabLocationLeading,
			want:     bottomIcon,
		},
		{
			name:     "trailing: tab with icon and text",
			item:     widget.NewTabItemWithIcon("Text1", theme.CancelIcon(), canvas.NewCircle(theme.BackgroundColor())),
			location: widget.TabLocationTrailing,
			want:     bottomIconAndText,
		},
		{
			name:     "trailing: tab with text only",
			item:     widget.NewTabItem("Text2", canvas.NewCircle(theme.BackgroundColor())),
			location: widget.TabLocationTrailing,
			want:     bottomText,
		},
		{
			name:     "trailing: tab with icon only",
			item:     widget.NewTabItemWithIcon("", theme.InfoIcon(), canvas.NewCircle(theme.BackgroundColor())),
			location: widget.TabLocationTrailing,
			want:     bottomIcon,
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
	test.ApplyTheme(t, theme.LightTheme())

	item1 := &widget.TabItem{Text: "Test1", Content: widget.NewLabel("Text 1")}
	item2 := &widget.TabItem{Text: "Test2", Content: widget.NewLabel("Text 2")}
	item3 := &widget.TabItem{Text: "Test3", Content: widget.NewLabel("Text 3")}
	tabs := widget.NewTabContainer(item1, item2, item3)
	w := test.NewWindow(tabs)
	defer w.Close()
	w.SetPadded(false)
	c := w.Canvas()

	tabsTop := `
		<canvas size="155x62">
			<content>
				<widget size="155x62" type="*widget.TabContainer">
					<widget pos="0,33" size="155x29" type="*widget.Label">
						<text pos="4,4" size="147x21">Text 1</text>
					</widget>
					<widget size="49x29" type="*widget.tabButton">
						<text alignment="1" bold color="focus" pos="4,4" size="41x21">Test1</text>
					</widget>
					<widget pos="53,0" size="49x29" type="*widget.tabButton">
						<text alignment="1" bold pos="4,4" size="41x21">Test2</text>
					</widget>
					<widget pos="106,0" size="49x29" type="*widget.tabButton">
						<text alignment="1" bold pos="4,4" size="41x21">Test3</text>
					</widget>
					<rectangle fillColor="shadow" pos="0,29" size="155x4"/>
					<rectangle fillColor="focus" pos="0,29" size="49x4"/>
				</widget>
			</content>
		</canvas>
	`
	tabsBottom := `
		<canvas size="155x62">
			<content>
				<widget size="155x62" type="*widget.TabContainer">
					<widget size="155x29" type="*widget.Label">
						<text pos="4,4" size="147x21">Text 1</text>
					</widget>
					<widget size="49x29" type="*widget.tabButton">
						<text alignment="1" bold color="focus" pos="4,4" size="41x21">Test1</text>
					</widget>
					<widget pos="53,0" size="49x29" type="*widget.tabButton">
						<text alignment="1" bold pos="4,4" size="41x21">Test2</text>
					</widget>
					<widget pos="106,0" size="49x29" type="*widget.tabButton">
						<text alignment="1" bold pos="4,4" size="41x21">Test3</text>
					</widget>
					<rectangle fillColor="shadow" pos="0,29" size="155x4"/>
					<rectangle fillColor="focus" pos="0,29" size="49x4"/>
				</widget>
			</content>
		</canvas>
	`

	w.Resize(tabs.MinSize())
	test.AssertRendersToMarkup(t, tabsTop, c)

	tabs.SetTabLocation(widget.TabLocationLeading)
	w.Resize(tabs.MinSize())
	test.AssertRendersToMarkup(t, tabsBottom, c, "leading is the same as bottom on mobile")

	tabs.SetTabLocation(widget.TabLocationBottom)
	w.Resize(tabs.MinSize())
	test.AssertRendersToMarkup(t, tabsBottom, c)

	tabs.SetTabLocation(widget.TabLocationTrailing)
	w.Resize(tabs.MinSize())
	test.AssertRendersToMarkup(t, tabsBottom, c, "trailing is the same as bottom on mobile")

	tabs.SetTabLocation(widget.TabLocationTop)
	w.Resize(tabs.MinSize())
	test.AssertRendersToMarkup(t, tabsTop, c)
}

func TestTabContainer_Tapped(t *testing.T) {
	test.NewApp()
	defer test.NewApp()
	test.ApplyTheme(t, theme.LightTheme())

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
					<widget size="64x29" type="*widget.tabButton">
						<text alignment="1" bold color="focus" pos="4,4" size="56x21">Test1</text>
					</widget>
					<widget pos="68,0" size="64x29" type="*widget.tabButton">
						<text alignment="1" bold pos="4,4" size="56x21">Test2</text>
					</widget>
					<widget pos="136,0" size="64x29" type="*widget.tabButton">
						<text alignment="1" bold pos="4,4" size="56x21">Test3</text>
					</widget>
					<rectangle fillColor="shadow" pos="0,29" size="200x4"/>
					<rectangle fillColor="focus" pos="0,29" size="64x4"/>
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
					<widget size="64x29" type="*widget.tabButton">
						<text alignment="1" bold pos="4,4" size="56x21">Test1</text>
					</widget>
					<widget pos="68,0" size="64x29" type="*widget.tabButton">
						<text alignment="1" bold color="focus" pos="4,4" size="56x21">Test2</text>
					</widget>
					<widget pos="136,0" size="64x29" type="*widget.tabButton">
						<text alignment="1" bold pos="4,4" size="56x21">Test3</text>
					</widget>
					<rectangle fillColor="shadow" pos="0,29" size="200x4"/>
					<rectangle fillColor="focus" pos="68,29" size="64x4"/>
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
					<widget size="64x29" type="*widget.tabButton">
						<text alignment="1" bold pos="4,4" size="56x21">Test1</text>
					</widget>
					<widget pos="68,0" size="64x29" type="*widget.tabButton">
						<text alignment="1" bold pos="4,4" size="56x21">Test2</text>
					</widget>
					<widget pos="136,0" size="64x29" type="*widget.tabButton">
						<text alignment="1" bold color="focus" pos="4,4" size="56x21">Test3</text>
					</widget>
					<rectangle fillColor="shadow" pos="0,29" size="200x4"/>
					<rectangle fillColor="focus" pos="136,29" size="64x4"/>
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
