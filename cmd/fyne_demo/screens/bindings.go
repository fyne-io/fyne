package screens

import (
	"fmt"
	"net/url"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/binding"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
)

// BindingsScreen loads a data bindings example panel for the demo app
func BindingsScreen() fyne.CanvasObject {
	// TODO entry
	// TODO form
	// TODO group
	// TODO list
	// TODO menu
	// TODO popup
	// TODO scroller
	// TODO splitter
	// TODO tabber
	// TODO table
	// TODO toolbar
	// TODO tree

	// Button <-> ProgressBar <-> Button
	goLeftButton := &widget.Button{}
	goRightButton := &widget.Button{}
	goProgressBar := &widget.ProgressBar{Max: 1}

	goBool := &binding.BoolBinding{}
	goFloat64 := &binding.Float64Binding{}
	goString := &binding.StringBinding{}
	goResource := &binding.ResourceBinding{}

	binding.BindButtonTapped(goLeftButton, goBool)
	binding.BindButtonTapped(goRightButton, goBool)
	binding.BindButtonText(goLeftButton, goString)
	binding.BindButtonText(goRightButton, goString)
	binding.BindButtonIcon(goLeftButton, goResource)
	binding.BindButtonIcon(goRightButton, goResource)
	binding.BindProgressBarValue(goProgressBar, goFloat64)

	goBool.AddBoolListener(func(b bool) {
		if b {
			// Start goroutine to update progress bar
			go func() {
				for num := 0.0; num < 1.0; num += 0.01 {
					time.Sleep(100 * time.Millisecond)
					goFloat64.SetFloat64(num)
				}
				goFloat64.SetFloat64(1.0)
				goBool.SetBool(false)
			}()
			goString.SetString("")
			goResource.SetResource(theme.InfoIcon())
		} else {
			goString.SetString("Go")
			goResource.SetResource(theme.MediaPlayIcon())
		}
	})

	goBool.SetBool(true)

	// Check <-> Label <-> Check
	// FIXME checkRenderer.Layout can be called before checkRenderer.Refresh
	//  meaning check box is not initially rendered checked
	onOffLeftCheck := &widget.Check{}
	onOffRightCheck := &widget.Check{}
	onOffLabel := &widget.Label{}

	onOffBool := &binding.BoolBinding{}
	onOffString := &binding.StringBinding{}
	onOffResource := &binding.ResourceBinding{}

	binding.BindCheckChanged(onOffLeftCheck, onOffBool)
	binding.BindCheckChanged(onOffRightCheck, onOffBool)
	binding.BindCheckText(onOffLeftCheck, onOffString)
	binding.BindCheckText(onOffRightCheck, onOffString)
	binding.BindLabelText(onOffLabel, onOffString)

	onOffBool.AddBoolListener(func(b bool) {
		if b {
			onOffString.SetString("On")
			onOffResource.SetResource(theme.CheckButtonCheckedIcon())
		} else {
			onOffString.SetString("Off")
			onOffResource.SetResource(theme.CheckButtonIcon())
		}
	})
	onOffBool.SetBool(true)

	// Radio <-> Icon <-> Radio
	// FIXME radioRenderer.Layout can be called before radioRenderer.Refresh
	//  meaning radio is not initially rendered selected
	// FIXME iconRenderer.Layout can be called before iconRenderer.Refresh
	//  meaning icon is not initially rendered
	options := []string{
		"Cut",
		"Copy",
		"Paste",
	}
	clipboardLeftRadio := &widget.Radio{Options: options}
	clipboardRightRadio := &widget.Radio{Options: options}
	clipboardIcon := &widget.Icon{}

	clipboardString := &binding.StringBinding{}
	clipboardResource := &binding.ResourceBinding{}

	binding.BindRadioChanged(clipboardLeftRadio, clipboardString)
	binding.BindRadioChanged(clipboardRightRadio, clipboardString)
	binding.BindIconResource(clipboardIcon, clipboardResource)

	clipboardString.AddStringListener(func(s string) {
		switch s {
		case "Cut":
			clipboardResource.SetResource(theme.ContentCutIcon())
		case "Copy":
			clipboardResource.SetResource(theme.ContentCopyIcon())
		case "Paste":
			clipboardResource.SetResource(theme.ContentPasteIcon())
		default:
			clipboardResource.SetResource(theme.QuestionIcon())
		}
	})
	clipboardString.SetString(options[0])

	// Select <-> Hyperlink <-> Select
	urls := []string{
		"https://fyne.io",
		"https://github.com/fyne-io",
	}
	urlLeftSelect := &widget.Select{Options: urls}
	urlRightSelect := &widget.Select{Options: urls}
	urlHyperlink := &widget.Hyperlink{}

	urlString := &binding.StringBinding{}
	urlURL := &binding.URLBinding{}

	binding.BindSelectChanged(urlLeftSelect, urlString)
	binding.BindSelectChanged(urlRightSelect, urlString)
	binding.BindHyperlinkText(urlHyperlink, urlString)
	binding.BindHyperlinkURL(urlHyperlink, urlURL)

	urlString.AddStringListener(func(s string) {
		u, err := url.Parse(s)
		if err != nil {
			fyne.LogError("Failed to parse URL: "+s, err)
		}
		urlURL.SetURL(u)
	})
	urlString.SetString(urls[0])

	// Slider <-> Label <-> Slider
	slideLeftSlider := &widget.Slider{Max: 1, Step: 0.01}
	slideRightSlider := &widget.Slider{Max: 1, Step: 0.01}
	slideLabel := &widget.Label{}

	slideFloat64 := &binding.Float64Binding{}
	slideString := &binding.StringBinding{}

	binding.BindSliderChanged(slideLeftSlider, slideFloat64)
	binding.BindSliderChanged(slideRightSlider, slideFloat64)
	binding.BindLabelText(slideLabel, slideString)

	slideFloat64.AddFloat64Listener(func(f float64) {
		slideString.SetString(fmt.Sprintf("%f", f))
	})
	slideFloat64.SetFloat64(0.25)

	return fyne.NewContainerWithLayout(layout.NewGridLayout(3),
		widget.NewLabel("Left Input"), widget.NewLabel("Output"), widget.NewLabel("Right Input"),
		goLeftButton, goProgressBar, goRightButton,
		onOffLeftCheck, onOffLabel, onOffRightCheck,
		clipboardLeftRadio, clipboardIcon, clipboardRightRadio,
		urlLeftSelect, urlHyperlink, urlRightSelect,
		slideLeftSlider, slideLabel, slideRightSlider,
	)
}
