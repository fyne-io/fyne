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
	// TODO entry - stringBinding
	// TODO list - []binding
	// TODO table - [][]binding
	// TODO tree - map[string][]string

	// Button <-> ProgressBar <-> Button
	goLeftButton := &widget.Button{}
	goRightButton := &widget.Button{}
	goProgressBar := &widget.ProgressBar{Max: 1}

	goBool := &binding.BoolBinding{}
	goFloat64 := &binding.Float64Binding{}
	goString := &binding.StringBinding{}
	goResource := &binding.ResourceBinding{}

	goLeftButton.BindTapped(goBool)
	goRightButton.BindTapped(goBool)
	goLeftButton.BindText(goString)
	goRightButton.BindText(goString)
	goLeftButton.BindIcon(goResource)
	goRightButton.BindIcon(goResource)
	goProgressBar.BindValue(goFloat64)

	goBool.AddListener(func(b bool) {
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
	// FIXME checkRenderer.Layout can be called before checkRenderer.Refresh
	//  meaning check box is not initially rendered checked
	onOffLeftCheck := &widget.Check{}
	onOffRightCheck := &widget.Check{}
	onOffLabel := &widget.Label{}

	onOffBool := &binding.BoolBinding{}
	onOffString := &binding.StringBinding{}
	onOffResource := &binding.ResourceBinding{}

	onOffLeftCheck.BindChecked(onOffBool)
	onOffRightCheck.BindChecked(onOffBool)
	onOffLeftCheck.BindText(onOffString)
	onOffRightCheck.BindText(onOffString)
	onOffLabel.BindText(onOffString)

	onOffBool.AddListener(func(b bool) {
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
	countLeftEntry := widget.NewEntry()
	countRightEntry := widget.NewEntry()
	countLabel := widget.NewLabel("0")

	countString := binding.NewStringBinding("")

	countLeftEntry.BindText(countString)
	countRightEntry.BindText(countString)

	countString.AddListener(func(s string) {
		countLabel.SetText(strconv.Itoa(len(s)))
	})

	// Radio <-> Icon <-> Radio
	// FIXME radioRenderer.Layout can be called before radioRenderer.Refresh
	//  meaning radio is not initially rendered selected
	// FIXME iconRenderer.Layout can be called before iconRenderer.Refresh
	//  meaning icon is not initially rendered

	clipboardLeftRadio := &widget.Radio{}
	clipboardRightRadio := &widget.Radio{}
	clipboardIcon := &widget.Icon{}

	clipboardOptions := &binding.SliceBinding{}
	clipboardString := &binding.StringBinding{}
	clipboardResource := &binding.ResourceBinding{}

	clipboardLeftRadio.BindOptions(clipboardOptions)
	clipboardRightRadio.BindOptions(clipboardOptions)
	clipboardLeftRadio.BindSelected(clipboardString)
	clipboardRightRadio.BindSelected(clipboardString)
	clipboardIcon.BindResource(clipboardResource)

	clipboardOptions.Append(
		binding.NewStringBinding("Cut"),
		binding.NewStringBinding("Copy"),
		binding.NewStringBinding("Paste"),
	)
	clipboardString.AddListener(func(s string) {
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
	urlLeftSelect := &widget.Select{}
	urlRightSelect := &widget.Select{}
	urlHyperlink := &widget.Hyperlink{}

	urlOptions := &binding.SliceBinding{}
	urlString := &binding.StringBinding{}
	urlURL := &binding.URLBinding{}

	urlLeftSelect.BindOptions(urlOptions)
	urlRightSelect.BindOptions(urlOptions)
	urlLeftSelect.BindSelected(urlString)
	urlRightSelect.BindSelected(urlString)
	urlHyperlink.BindText(urlString)
	urlHyperlink.BindURL(urlURL)

	urlOptions.Append(
		binding.NewStringBinding("https://fyne.io"),
		binding.NewStringBinding("https://github.com/fyne-io"),
	)
	urlString.AddListener(func(s string) {
		u, err := url.Parse(s)
		if err != nil {
			fyne.LogError("Failed to parse URL: "+s, err)
		}
		urlURL.Set(u)
	})
	urlString.Set("")

	// Slider <-> Label <-> Slider
	slideLeftSlider := &widget.Slider{Max: 1, Step: 0.01}
	slideRightSlider := &widget.Slider{Max: 1, Step: 0.01}
	slideLabel := &widget.Label{}

	slideFloat64 := &binding.Float64Binding{}
	slideString := &binding.StringBinding{}

	slideLeftSlider.BindValue(slideFloat64)
	slideRightSlider.BindValue(slideFloat64)
	slideLabel.BindText(slideString)

	slideFloat64.AddListener(func(f float64) {
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
