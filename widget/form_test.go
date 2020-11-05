package widget

import (
	"testing"

	"fyne.io/fyne"
	"fyne.io/fyne/test"
	"fyne.io/fyne/theme"

	"github.com/stretchr/testify/assert"
)

func TestFormSize(t *testing.T) {
	form := &Form{Items: []*FormItem{
		{Text: "test1", Widget: NewEntry()},
		{Text: "test2", Widget: NewEntry()},
	}}

	assert.Equal(t, 2, len(form.Items))
}

func TestForm_CreateRenderer(t *testing.T) {
	form := &Form{Items: []*FormItem{{Text: "test1", Widget: NewEntry()}}}
	assert.NotNil(t, test.WidgetRenderer(form))
	assert.Equal(t, 2, len(form.itemGrid.Objects))

	form.Append("test2", NewEntry())
	assert.Equal(t, 4, len(form.itemGrid.Objects))
}

func TestForm_Append(t *testing.T) {
	form := &Form{Items: []*FormItem{{Text: "test1", Widget: NewEntry()}}}
	assert.Equal(t, 1, len(form.Items))

	form.Append("test2", NewEntry())
	assert.True(t, len(form.Items) == 2)

	item := &FormItem{Text: "test3", Widget: NewEntry()}
	form.AppendItem(item)
	assert.True(t, len(form.Items) == 3)
	assert.Equal(t, item, form.Items[2])
}

func TestForm_CustomButtonsText(t *testing.T) {
	form := &Form{OnSubmit: func() {}, OnCancel: func() {}}
	form.Append("test", NewEntry())
	assert.Equal(t, "Submit", form.SubmitText)
	assert.Equal(t, "Cancel", form.CancelText)

	form = &Form{OnSubmit: func() {}, SubmitText: "Apply",
		OnCancel: func() {}, CancelText: "Close"}
	assert.Equal(t, "Apply", form.SubmitText)
	assert.Equal(t, "Close", form.CancelText)
}

func TestForm_AddRemoveButton(t *testing.T) {
	scount := 0
	ccount := 0
	sscount := 10
	form := &Form{OnSubmit: func() {}, OnCancel: func() {}}
	form.Append("test", NewEntry())
	form.OnSubmit = func() { scount++ }
	form.OnCancel = func() { ccount++ }
	form.Refresh()

	test.Tap(form.submitButton)
	assert.Equal(t, 1, scount, "tapping submit should incr scount")

	test.Tap(form.cancelButton)
	assert.Equal(t, 1, ccount, "tapping cancel should incr ccount")

	form.OnSubmit = func() { sscount++ }
	form.Refresh()
	test.Tap(form.submitButton)
	assert.Equal(t, 11, sscount, "tapping new submit should incr sscount from 10 to 11")

	form.OnCancel = func() { sscount = sscount - 6 }
	form.Refresh()
	test.Tap(form.cancelButton)
	assert.Equal(t, 5, sscount, "tapping new cancel should decr ssount from 11 down to 5")
}

func TestForm_Renderer(t *testing.T) {
	app := test.NewApp()
	defer test.NewApp()
	app.Settings().SetTheme(theme.LightTheme())

	form := &Form{
		Items: []*FormItem{
			{Text: "test1", Widget: NewEntry()},
			{Text: "test2", Widget: NewEntry()},
		},
		OnSubmit: func() {}, OnCancel: func() {}}
	w := test.NewWindow(form)
	defer w.Close()

	test.AssertRendersToMarkup(t, `
		<canvas padded size="214x127">
			<content>
				<widget pos="4,4" size="206x119" type="*widget.Form">
					<container size="206x78">
						<widget size="47x37" type="*widget.Label">
							<text alignment="trailing" bold pos="4,4" size="39x21">test1</text>
						</widget>
						<widget pos="51,0" size="155x37" type="*widget.Entry">
							<rectangle fillColor="shadow" pos="0,33" size="155x4"/>
							<widget pos="4,4" size="147x29" type="*widget.textProvider">
								<text color="placeholder" pos="4,4" size="139x21"></text>
							</widget>
							<widget pos="4,4" size="147x29" type="*widget.textProvider">
								<text pos="4,4" size="139x21"></text>
							</widget>
						</widget>
						<widget pos="0,41" size="47x37" type="*widget.Label">
							<text alignment="trailing" bold pos="4,4" size="39x21">test2</text>
						</widget>
						<widget pos="51,41" size="155x37" type="*widget.Entry">
							<rectangle fillColor="shadow" pos="0,33" size="155x4"/>
							<widget pos="4,4" size="147x29" type="*widget.textProvider">
								<text color="placeholder" pos="4,4" size="139x21"></text>
							</widget>
							<widget pos="4,4" size="147x29" type="*widget.textProvider">
								<text pos="4,4" size="139x21"></text>
							</widget>
						</widget>
					</container>
					<widget pos="0,82" size="206x37" type="*widget.Box">
						<spacer size="0x0"/>
						<widget size="99x37" type="*widget.Button">
							<widget pos="2,2" size="95x33" type="*widget.Shadow">
								<radialGradient centerOffset="0.5,0.5" pos="-2,-2" size="2x2" startColor="shadow"/>
								<linearGradient endColor="shadow" pos="0,-2" size="95x2"/>
								<radialGradient centerOffset="-0.5,0.5" pos="95,-2" size="2x2" startColor="shadow"/>
								<linearGradient angle="270" pos="95,0" size="2x33" startColor="shadow"/>
								<radialGradient centerOffset="-0.5,-0.5" pos="95,33" size="2x2" startColor="shadow"/>
								<linearGradient pos="0,33" size="95x2" startColor="shadow"/>
								<radialGradient centerOffset="0.5,-0.5" pos="-2,33" size="2x2" startColor="shadow"/>
								<linearGradient angle="270" endColor="shadow" pos="-2,0" size="2x33"/>
							</widget>
							<rectangle pos="2,2" size="95x33"/>
							<text bold pos="36,8" size="51x21">Cancel</text>
							<image fillMode="contain" pos="12,8" rsc="cancelIcon" size="20x21"/>
						</widget>
						<widget pos="103,0" size="103x37" type="*widget.Button">
							<widget pos="2,2" size="99x33" type="*widget.Shadow">
								<radialGradient centerOffset="0.5,0.5" pos="-2,-2" size="2x2" startColor="shadow"/>
								<linearGradient endColor="shadow" pos="0,-2" size="99x2"/>
								<radialGradient centerOffset="-0.5,0.5" pos="99,-2" size="2x2" startColor="shadow"/>
								<linearGradient angle="270" pos="99,0" size="2x33" startColor="shadow"/>
								<radialGradient centerOffset="-0.5,-0.5" pos="99,33" size="2x2" startColor="shadow"/>
								<linearGradient pos="0,33" size="99x2" startColor="shadow"/>
								<radialGradient centerOffset="0.5,-0.5" pos="-2,33" size="2x2" startColor="shadow"/>
								<linearGradient angle="270" endColor="shadow" pos="-2,0" size="2x33"/>
							</widget>
							<rectangle fillColor="focus" pos="2,2" size="99x33"/>
							<text bold color="background" pos="36,8" size="55x21">Submit</text>
							<image fillMode="contain" pos="12,8" rsc="confirmIcon" size="20x21" themed="inverted"/>
						</widget>
					</widget>
				</widget>
			</content>
		</canvas>
	`, w.Canvas())
}

func TestForm_ChangeText(t *testing.T) {
	item := &FormItem{Text: "Test", Widget: NewEntry()}
	form := &Form{Items: []*FormItem{item}}

	renderer := test.WidgetRenderer(form)
	c := renderer.Objects()[0].(*fyne.Container)
	assert.Equal(t, "Test", c.Objects[0].(*Label).Text)

	item.Text = "Changed"
	form.Refresh()
	assert.Equal(t, "Changed", c.Objects[0].(*Label).Text)
}

func TestForm_ChangeTheme(t *testing.T) {
	app := test.NewApp()
	defer test.NewApp()
	app.Settings().SetTheme(theme.LightTheme())

	form := &Form{
		Items: []*FormItem{
			{Text: "test1", Widget: NewEntry()},
			{Text: "test2", Widget: NewEntry()},
		},
		OnSubmit: func() {}, OnCancel: func() {}}
	w := test.NewWindow(form)
	w.Resize(fyne.NewSize(340, 240))
	defer w.Close()

	test.AssertImageMatches(t, "form/theme_initial.png", w.Canvas().Capture())

	test.WithTestTheme(t, func() {
		form.Refresh()
		test.AssertImageMatches(t, "form/theme_changed.png", w.Canvas().Capture())
	})
}
