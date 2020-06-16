package screens

import (
	"fmt"
	"image/color"
	"net/url"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
)

const (
	loremIpsum = `Lorem ipsum dolor sit amet, consectetur adipiscing elit. Pellentesque quis consectetur nisi. Suspendisse id interdum felis.
Sed egestas eget tellus eu pharetra. Praesent pulvinar sed massa id placerat. Etiam sem libero, semper vitae consequat ut, volutpat id mi.
Mauris volutpat pellentesque convallis. Curabitur rutrum venenatis orci nec ornare. Maecenas quis pellentesque neque.
Aliquam consectetur dapibus nulla, id maximus odio ultrices ac. Sed luctus at felis sed faucibus. Cras leo augue, congue in velit ut, mattis rhoncus lectus.

Praesent viverra, mauris ut ullamcorper semper, leo urna auctor lectus, vitae vehicula mi leo quis lorem.
Nullam condimentum, massa at tempor feugiat, metus enim lobortis velit, eget suscipit eros ipsum quis tellus. Aenean fermentum diam vel felis dictum semper.
Duis nisl orci, tincidunt ut leo quis, luctus vehicula diam. Sed velit justo, congue id augue eu, euismod dapibus lacus. Proin sit amet imperdiet sapien.
Mauris erat urna, fermentum et quam rhoncus, fringilla consequat ante. Vivamus consectetur molestie odio, ac rutrum erat finibus a.
Suspendisse id maximus felis. Sed mauris odio, mattis eget mi eu, consequat tempus purus.`
)

var (
	progress    *widget.ProgressBar
	infProgress *widget.ProgressBarInfinite
	endProgress chan interface{}
)

func makeButtonTab() fyne.Widget {
	disabled := widget.NewButton("Disabled", func() {})
	disabled.Disable()

	shareItem := fyne.NewMenuItem("Share via", nil)
	shareItem.ChildMenu = fyne.NewMenu("",
		fyne.NewMenuItem("Twitter", func() { fmt.Println("context menu Share->Twitter") }),
		fyne.NewMenuItem("Reddit", func() { fmt.Println("context menu Share->Reddit") }),
	)
	menuLabel := newContextMenuButton("tap me for pop-up menu with submenus", fyne.NewMenu("",
		fyne.NewMenuItem("Copy", func() { fmt.Println("context menu copy") }),
		shareItem,
	))

	return widget.NewVBox(
		widget.NewButton("Button (text only)", func() { fmt.Println("tapped text button") }),
		widget.NewButtonWithIcon("Button (text & leading icon)", theme.ConfirmIcon(), func() { fmt.Println("tapped text & leading icon button") }),
		&widget.Button{
			Alignment: widget.ButtonAlignLeading,
			Text:      "Button (leading-aligned, text only)",
			OnTapped:  func() { fmt.Println("tapped leading-aligned, text only button") },
		},
		&widget.Button{
			Alignment:     widget.ButtonAlignTrailing,
			IconPlacement: widget.ButtonIconTrailingText,
			Text:          "Button (trailing-aligned, text & trailing icon)",
			Icon:          theme.ConfirmIcon(),
			OnTapped:      func() { fmt.Println("tapped trailing-aligned, text & trailing icon button") },
		},
		disabled,
		layout.NewSpacer(),
		layout.NewSpacer(),
		menuLabel,
		layout.NewSpacer(),
	)
}

func makeTextGrid() *widget.TextGrid {
	grid := widget.NewTextGridFromString("TextGrid\n  Content\nZebra")
	grid.SetStyleRange(0, 4, 0, 7,
		&widget.CustomTextGridStyle{BGColor: &color.NRGBA{R: 64, G: 64, B: 192, A: 128}})
	grid.Rows[1].Style = &widget.CustomTextGridStyle{BGColor: &color.NRGBA{R: 64, G: 192, B: 64, A: 128}}

	white := &widget.CustomTextGridStyle{FGColor: color.White, BGColor: color.Black}
	black := &widget.CustomTextGridStyle{FGColor: color.Black, BGColor: color.White}
	grid.Rows[2].Cells[0].Style = white
	grid.Rows[2].Cells[1].Style = black
	grid.Rows[2].Cells[2].Style = white
	grid.Rows[2].Cells[3].Style = black
	grid.Rows[2].Cells[4].Style = white

	grid.ShowLineNumbers = true
	grid.ShowWhitespace = true

	return grid
}

func makeTextTab() fyne.CanvasObject {
	label := widget.NewLabel("Label")

	link, err := url.Parse("https://fyne.io/")
	if err != nil {
		fyne.LogError("Could not parse URL", err)
	}
	hyperlink := widget.NewHyperlink("Hyperlink", link)

	entryLoremIpsum := widget.NewMultiLineEntry()
	entryLoremIpsum.SetText(loremIpsum)
	entryLoremIpsumScroller := widget.NewVScrollContainer(entryLoremIpsum)

	label.Alignment = fyne.TextAlignLeading
	hyperlink.Alignment = fyne.TextAlignLeading

	label.Wrapping = fyne.TextWrapWord
	hyperlink.Wrapping = fyne.TextWrapWord
	entryLoremIpsum.Wrapping = fyne.TextWrapWord

	radioAlign := widget.NewRadio([]string{"Text Alignment Leading", "Text Alignment Center", "Text Alignment Trailing"}, func(s string) {
		var align fyne.TextAlign
		switch s {
		case "Text Alignment Leading":
			align = fyne.TextAlignLeading
		case "Text Alignment Center":
			align = fyne.TextAlignCenter
		case "Text Alignment Trailing":
			align = fyne.TextAlignTrailing
		}

		label.Alignment = align
		hyperlink.Alignment = align

		label.Refresh()
		hyperlink.Refresh()
	})
	radioAlign.SetSelected("Text Alignment Leading")

	radioWrap := widget.NewRadio([]string{"Text Wrapping Off", "Text Wrapping Truncate", "Text Wrapping Break", "Text Wrapping Word"}, func(s string) {
		var wrap fyne.TextWrap
		switch s {
		case "Text Wrapping Off":
			wrap = fyne.TextWrapOff
		case "Text Wrapping Truncate":
			wrap = fyne.TextTruncate
		case "Text Wrapping Break":
			wrap = fyne.TextWrapBreak
		case "Text Wrapping Word":
			wrap = fyne.TextWrapWord
		}

		label.Wrapping = wrap
		hyperlink.Wrapping = wrap
		if wrap != fyne.TextTruncate {
			entryLoremIpsum.Wrapping = wrap
		}

		label.Refresh()
		hyperlink.Refresh()
		entryLoremIpsum.Refresh()
		entryLoremIpsumScroller.Refresh()
	})
	radioWrap.SetSelected("Text Wrapping Word")

	fixed := widget.NewVBox(
		widget.NewHBox(
			radioAlign,
			layout.NewSpacer(),
			radioWrap,
		),
		label,
		hyperlink,
	)

	grid := makeTextGrid()
	return fyne.NewContainerWithLayout(layout.NewBorderLayout(fixed, grid, nil, nil),
		fixed, entryLoremIpsumScroller, grid)
}

func makeInputTab() fyne.Widget {
	entry := widget.NewEntry()
	entry.SetPlaceHolder("Entry")
	entryDisabled := widget.NewEntry()
	entryDisabled.SetText("Entry (disabled)")
	entryDisabled.Disable()
	entryMultiLine := widget.NewMultiLineEntry()
	entryMultiLine.SetPlaceHolder("MultiLine Entry")
	selectEntry := widget.NewSelectEntry([]string{"Option A", "Option B", "Option C"})
	selectEntry.PlaceHolder = "Type or select"
	disabledCheck := widget.NewCheck("Disabled check", func(bool) {})
	disabledCheck.Disable()
	radio := widget.NewRadio([]string{"Radio Item 1", "Radio Item 2"}, func(s string) { fmt.Println("selected", s) })
	radio.Horizontal = true
	disabledRadio := widget.NewRadio([]string{"Disabled radio"}, func(string) {})
	disabledRadio.Disable()

	return widget.NewVBox(
		entry,
		entryDisabled,
		entryMultiLine,
		widget.NewSelect([]string{"Option 1", "Option 2", "Option 3"}, func(s string) { fmt.Println("selected", s) }),
		selectEntry,
		widget.NewCheck("Check", func(on bool) { fmt.Println("checked", on) }),
		disabledCheck,
		radio,
		disabledRadio,
		widget.NewSlider(0, 100),
	)
}

func makeProgressTab() fyne.Widget {
	progress = widget.NewProgressBar()
	infProgress = widget.NewProgressBarInfinite()
	endProgress = make(chan interface{}, 1)

	return widget.NewVBox(
		widget.NewLabel("Percent"), progress,
		widget.NewLabel("Infinite"), infProgress)
}

func makeFormTab() fyne.Widget {
	name := widget.NewEntry()
	name.SetPlaceHolder("John Smith")
	email := widget.NewEntry()
	email.SetPlaceHolder("test@example.com")
	password := widget.NewPasswordEntry()
	password.SetPlaceHolder("Password")
	largeText := widget.NewMultiLineEntry()

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Name", Widget: name},
			{Text: "Email", Widget: email},
		},
		OnCancel: func() {
			fmt.Println("Cancelled")
		},
		OnSubmit: func() {
			fmt.Println("Form submitted")
			fyne.CurrentApp().SendNotification(&fyne.Notification{
				Title:   "Form for: " + name.Text,
				Content: largeText.Text,
			})
		},
	}
	form.Append("Password", password)
	form.Append("Message", largeText)
	return form
}

func startProgress() {
	progress.SetValue(0)
	select { // ignore stale end message
	case <-endProgress:
	default:
	}

	go func() {
		num := 0.0
		for num < 1.0 {
			time.Sleep(100 * time.Millisecond)
			select {
			case <-endProgress:
				return
			default:
			}

			progress.SetValue(num)
			num += 0.01
		}

		progress.SetValue(1)
	}()
	infProgress.Start()
}

func stopProgress() {
	if !infProgress.Running() {
		return
	}

	infProgress.Stop()
	endProgress <- struct{}{}
}

// WidgetScreen shows a panel containing widget demos
func WidgetScreen() fyne.CanvasObject {
	toolbar := widget.NewToolbar(widget.NewToolbarAction(theme.MailComposeIcon(), func() { fmt.Println("New") }),
		widget.NewToolbarSeparator(),
		widget.NewToolbarSpacer(),
		widget.NewToolbarAction(theme.ContentCutIcon(), func() { fmt.Println("Cut") }),
		widget.NewToolbarAction(theme.ContentCopyIcon(), func() { fmt.Println("Copy") }),
		widget.NewToolbarAction(theme.ContentPasteIcon(), func() { fmt.Println("Paste") }),
	)

	progress := makeProgressTab()
	tabs := widget.NewTabContainer(
		widget.NewTabItem("Buttons", makeButtonTab()),
		widget.NewTabItem("Text", makeTextTab()),
		widget.NewTabItem("Input", makeInputTab()),
		widget.NewTabItem("Progress", progress),
		widget.NewTabItem("Form", makeFormTab()),
	)
	tabs.OnChanged = func(t *widget.TabItem) {
		if t.Content == progress {
			startProgress()
		} else {
			stopProgress()
		}
	}

	return fyne.NewContainerWithLayout(layout.NewBorderLayout(toolbar, nil, nil, nil),
		toolbar, tabs,
	)
}

type contextMenuButton struct {
	widget.Button
	menu *fyne.Menu
}

func (b *contextMenuButton) Tapped(e *fyne.PointEvent) {
	widget.ShowPopUpMenuAtPosition(b.menu, fyne.CurrentApp().Driver().CanvasForObject(b), e.AbsolutePosition)
}

func newContextMenuButton(label string, menu *fyne.Menu) *contextMenuButton {
	b := &contextMenuButton{menu: menu}
	b.Text = label

	b.ExtendBaseWidget(b)
	return b
}
