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
	// TODO form - map[string]binding
	// TODO list - []binding
	// TODO table - [][]binding
	// TODO tree - []binding

	// Button <-> ProgressBar <-> Button
	var goProgress float64
	var goText string
	var goRes fyne.Resource
	goFloat64 := binding.NewFloat64(&goProgress)
	goString := binding.NewString(&goText)
	goResource := binding.NewResource(&goRes)

	reset := func() {
		goString.Set("Go")
		goResource.Set(theme.MediaPlayIcon())
	}
	trigger := func() {
		// Start goroutine to update progress bar
		go func() {
			for num := 0.0; num < 1.0; num += 0.01 {
				time.Sleep(100 * time.Millisecond)
				goFloat64.Set(num)
			}
			goFloat64.Set(1.0)
			reset()
		}()
		goString.Set("")
		goResource.Set(theme.InfoIcon())
	}

	goLeftButton := widget.NewButton("", trigger).BindText(goString).BindIcon(goResource)
	goRightButton := widget.NewButton("", trigger).BindText(goString).BindIcon(goResource)
	goProgressBar := (&widget.ProgressBar{Max: 1}).BindValue(goFloat64)

	reset()

	// Check <-> Label <-> Check
	var onOffChecked bool
	var onOffText string
	var onOffRes fyne.Resource
	onOffBool := binding.NewBool(&onOffChecked)
	onOffString := binding.NewString(&onOffText)
	onOffResource := binding.NewResource(&onOffRes)

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
	var countText string
	countString := binding.NewString(&countText)

	countLeftEntry := widget.NewEntry().BindText(countString)
	countRightEntry := widget.NewEntry().BindText(countString)
	countLabel := widget.NewLabel("0")

	countString.AddStringListener(func(s string) {
		countLabel.SetText(strconv.Itoa(len(s)))
	})

	// Radio <-> Icon <-> Radio
	clipboardSlice := []string{"Cut", "Copy", "Paste"}
	var clipboardText string
	var clipboardRes fyne.Resource
	clipboardOptions := binding.NewStringList(clipboardSlice...)
	clipboardString := binding.NewString(&clipboardText)
	clipboardResource := binding.NewResource(&clipboardRes)

	clipboardLeftRadio := (&widget.Radio{}).BindOptions(clipboardOptions).BindSelected(clipboardString)
	clipboardRightRadio := (&widget.Radio{}).BindOptions(clipboardOptions).BindSelected(clipboardString)
	clipboardIcon := (&widget.Icon{}).BindResource(clipboardResource)

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

	// Select <-> Hyperlink <-> Select
	urlSlice := []string{"https://fyne.io", "https://github.com/fyne-io"}
	var urlText string
	var urlU *url.URL
	urlOptions := binding.NewStringList(urlSlice...)
	urlString := binding.NewString(&urlText)
	urlURL := binding.NewURL(&urlU)

	urlLeftSelect := (&widget.Select{}).BindOptions(urlOptions).BindSelected(urlString)
	urlRightSelect := (&widget.Select{}).BindOptions(urlOptions).BindSelected(urlString)
	urlHyperlink := (&widget.Hyperlink{}).BindText(urlString).BindURL(urlURL)

	urlString.AddStringListener(func(s string) {
		u, err := url.Parse(s)
		if err != nil {
			fyne.LogError("Failed to parse URL: "+s, err)
		} else {
			urlURL.Set(u)
		}
	})

	// Slider <-> Label <-> Slider
	slideValue := 0.0
	var slideText string
	slideFloat64 := binding.NewFloat64(&slideValue)
	slideString := binding.NewString(&slideText)

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
