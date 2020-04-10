package screens

import (
	"fmt"
	"net/url"
	"strconv"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/binding"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
)

// BindingsScreen loads a data bindings example panel for the demo app
func BindingsScreen() fyne.CanvasObject {
	// TODO list - []binding
	// TODO table - [][]binding
	// TODO tree - map[string][]string

	// Button <-> ProgressBar <-> Button
	goBool := &binding.BoolBinding{}
	goFloat64 := &binding.Float64Binding{}
	goString := &binding.StringBinding{}
	goResource := &binding.ResourceBinding{}

	goLeftButton := (&widget.Button{}).BindTapped(goBool).BindText(goString).BindIcon(goResource)
	goRightButton := (&widget.Button{}).BindTapped(goBool).BindText(goString).BindIcon(goResource)
	goProgressBar := (&widget.ProgressBar{Max: 1}).BindValue(goFloat64)

	goBool.AddBoolListener(func(b bool) {
		if b {
			// Start goroutine to update progress bar
			go func() {
				for num := 0.0; num < 1.0; num += 0.01 {
					time.Sleep(100 * time.Millisecond)
					goFloat64.Set(num)
				}
				goFloat64.Set(1.0)
				goBool.Set(false)
			}()
			goString.Set("")
			goResource.Set(theme.InfoIcon())
		} else {
			goString.Set("Go")
			goResource.Set(theme.MediaPlayIcon())
		}
	})

	goBool.Set(true)

	// Check <-> Label <-> Check
	onOffBool := &binding.BoolBinding{}
	onOffString := &binding.StringBinding{}
	onOffResource := &binding.ResourceBinding{}

	onOffLeftCheck := (&widget.Check{}).BindChecked(onOffBool).BindText(onOffString)
	onOffRightCheck := (&widget.Check{}).BindChecked(onOffBool).BindText(onOffString)
	onOffLabel := (&widget.Label{}).BindText(onOffString)

	onOffBool.AddBoolListener(func(b bool) {
		if b {
			onOffString.Set("On")
			onOffResource.Set(theme.CheckButtonCheckedIcon())
		} else {
			onOffString.Set("Off")
			onOffResource.Set(theme.CheckButtonIcon())
		}
	})
	onOffBool.Set(true)

	// Entry <-> Label <-> Entry
	countString := binding.NewStringBinding("")

	countLeftEntry := widget.NewEntry().BindText(countString)
	countRightEntry := widget.NewEntry().BindText(countString)
	countLabel := widget.NewLabel("0")

	countString.AddStringListener(func(s string) {
		countLabel.SetText(strconv.Itoa(len(s)))
	})

	// Radio <-> Icon <-> Radio
	clipboardOptions := &binding.ListBinding{}
	clipboardString := &binding.StringBinding{}
	clipboardResource := &binding.ResourceBinding{}

	clipboardLeftRadio := (&widget.Radio{}).BindOptions(clipboardOptions).BindSelected(clipboardString)
	clipboardRightRadio := (&widget.Radio{}).BindOptions(clipboardOptions).BindSelected(clipboardString)
	clipboardIcon := (&widget.Icon{}).BindResource(clipboardResource)

	clipboardOptions.Append(
		binding.NewStringBinding("Cut"),
		binding.NewStringBinding("Copy"),
		binding.NewStringBinding("Paste"),
	)
	clipboardString.AddStringListener(func(s string) {
		switch s {
		case "Cut":
			clipboardResource.Set(theme.ContentCutIcon())
		case "Copy":
			clipboardResource.Set(theme.ContentCopyIcon())
		case "Paste":
			clipboardResource.Set(theme.ContentPasteIcon())
		default:
			clipboardResource.Set(theme.QuestionIcon())
		}
	})
	clipboardString.Set("")

	// Select <-> Hyperlink <-> Select
	urlOptions := &binding.ListBinding{}
	urlString := &binding.StringBinding{}
	urlURL := &binding.URLBinding{}

	urlLeftSelect := (&widget.Select{}).BindOptions(urlOptions).BindSelected(urlString)
	urlRightSelect := (&widget.Select{}).BindOptions(urlOptions).BindSelected(urlString)
	urlHyperlink := (&widget.Hyperlink{}).BindText(urlString).BindURL(urlURL)

	urlOptions.Append(
		binding.NewStringBinding("https://fyne.io"),
		binding.NewStringBinding("https://github.com/fyne-io"),
	)
	urlString.AddStringListener(func(s string) {
		u, err := url.Parse(s)
		if err != nil {
			fyne.LogError("Failed to parse URL: "+s, err)
		}
		urlURL.Set(u)
	})
	urlString.Set("")

	// Slider <-> Label <-> Slider
	slideFloat64 := &binding.Float64Binding{}
	slideString := &binding.StringBinding{}

	slideLeftSlider := (&widget.Slider{Max: 1, Step: 0.01}).BindValue(slideFloat64)
	slideRightSlider := (&widget.Slider{Max: 1, Step: 0.01}).BindValue(slideFloat64)
	slideLabel := (&widget.Label{}).BindText(slideString)

	slideFloat64.AddFloat64Listener(func(f float64) {
		slideString.Set(fmt.Sprintf("%f", f))
	})
	slideFloat64.Set(0.25)

	return fyne.NewContainerWithLayout(layout.NewGridLayout(3),
		widget.NewLabel("Left Input"), widget.NewLabel("Output"), widget.NewLabel("Right Input"),
		goLeftButton, goProgressBar, goRightButton,
		onOffLeftCheck, onOffLabel, onOffRightCheck,
		countLeftEntry, countLabel, countRightEntry,
		clipboardLeftRadio, clipboardIcon, clipboardRightRadio,
		urlLeftSelect, urlHyperlink, urlRightSelect,
		slideLeftSlider, slideLabel, slideRightSlider,
	)
}
