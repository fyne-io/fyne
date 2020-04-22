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
	goFloat64 := binding.EmptyFloat64()
	goString := binding.EmptyString()
	goResource := binding.EmptyResource()

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

	goLeftButton := &widget.Button{
		//Text:     goString,
		//Icon:     goResource,
		OnTapped: trigger,
	}
	goRightButton := &widget.Button{
		//Text:     goString,
		//Icon:     goResource,
		OnTapped: trigger,
	}
	goProgressBar := &widget.ProgressBar{
		//Max:   binding.NewFloat64(1),
		//Value: goFloat64,
	}

	reset()

	// Check <-> Label <-> Check
	onOffBool := binding.NewBool(true)
	onOffString := binding.EmptyString()
	onOffResource := binding.EmptyResource()

	onOffLeftCheck := &widget.Check{
		//Checked: onOffBool,
		//Text:    onOffString,
	}
	onOffRightCheck := &widget.Check{
		//Checked: onOffBool,
		//Text:    onOffString,
	}
	onOffLabel := &widget.Label{
		//Text: onOffString,
	}

	// Create transformation pipeline
	onOffBool.OnUpdate(func(b bool) {
		if b {
			onOffString.Set("On")
			onOffResource.Set(theme.CheckButtonCheckedIcon())
		} else {
			onOffString.Set("Off")
			onOffResource.Set(theme.CheckButtonIcon())
		}
	})

	// Entry <-> Label <-> Entry
	countString := binding.EmptyString()
	countStringLength := binding.EmptyString()

	countLeftEntry := &widget.Entry{
		//Text: countString,
	}
	countRightEntry := &widget.Entry{
		//Text: countString,
	}
	countLabel := &widget.Label{
		//Text: countStringLength,
	}

	// Create transformation pipeline
	countString.OnUpdate(func(s string) {
		countStringLength.Set(strconv.Itoa(len(s)))
	})

	// Radio <-> Icon <-> Radio
	cut := "Cut"
	copy := "Copy"
	paste := "Paste"
	//clipboardOptions := binding.NewStringList(cut, copy, paste)
	clipboardString := binding.EmptyString()
	clipboardResource := binding.EmptyResource()

	clipboardLeftRadio := &widget.Radio{
		//Options:  clipboardOptions,
		//Selected: clipboardString,
	}
	clipboardRightRadio := &widget.Radio{
		//Options:  clipboardOptions,
		//Selected: clipboardString,
	}
	clipboardIcon := &widget.Icon{
		//Resource: clipboardResource,
	}

	// Create transformation pipeline
	clipboardString.OnUpdate(func(s string) {
		switch s {
		case cut:
			clipboardResource.Set(theme.ContentCutIcon())
		case copy:
			clipboardResource.Set(theme.ContentCopyIcon())
		case paste:
			clipboardResource.Set(theme.ContentPasteIcon())
		default:
			clipboardResource.Set(theme.QuestionIcon())
		}
	})

	// Select <-> Hyperlink <-> Select
	//url1 := "https://fyne.io"
	//url2 := "https://github.com/fyne-io"
	//urlOptions := binding.NewStringList(url1, url2)
	urlString := binding.EmptyString()
	urlURL := binding.EmptyURL()

	urlLeftSelect := &widget.Select{
		//Options:  urlOptions,
		//Selected: urlString,
	}
	urlRightSelect := &widget.Select{
		//Options:  urlOptions,
		//Selected: urlString,
	}
	urlHyperlink := &widget.Hyperlink{
		//Text: urlString,
		//URL:  urlURL,
	}

	// Create transformation pipeline
	urlString.OnUpdate(func(s string) {
		u, err := url.Parse(s)
		if err != nil {
			fyne.LogError("Failed to parse URL: "+s, err)
		} else {
			urlURL.Set(u)
		}
	})

	// Slider <-> Label <-> Slider
	//slideMax := binding.NewFloat64(1)
	//slideStep := binding.NewFloat64(0.01)
	slideValue := binding.NewFloat64(0.25)
	slideString := binding.EmptyString()

	slideLeftSlider := &widget.Slider{
		//Max:   slideMax,
		//Step:  slideStep,
		//Value: slideValue,
	}
	slideRightSlider := &widget.Slider{
		//Max:   slideMax,
		//Step:  slideStep,
		//Value: slideValue,
	}
	slideLabel := &widget.Label{
		//Text: slideString,
	}

	// Create transformation pipeline
	slideValue.OnUpdate(func(f float64) {
		slideString.Set(fmt.Sprintf("%f", f))
	})

	return fyne.NewContainerWithLayout(layout.NewGridLayout(3),
		widget.NewLabel("Left Input"), widget.NewLabel("Output"), widget.NewLabel("Right Input"),
		goLeftButton, goProgressBar, goRightButton,
		onOffLeftCheck, onOffLabel, onOffRightCheck,
		countLeftEntry, countLabel, countRightEntry,
		clipboardLeftRadio, clipboardIcon, clipboardRightRadio,
		//trackNameSelect, trackInfo, trackUnitSelect,
		urlLeftSelect, urlHyperlink, urlRightSelect,
		slideLeftSlider, slideLabel, slideRightSlider,
	)
}
