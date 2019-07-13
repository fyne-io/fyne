package main

import (
	"encoding/json"
	"image/color"
	"os"
	"path/filepath"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
)

type settings struct {
	fyneSettings app.SettingsSchema
}

func (s *settings) save() error {
	return s.saveToFile(s.fyneSettings.StoragePath())
}

func (s *settings) saveToFile(path string) error {
	err := os.MkdirAll(filepath.Dir(path), 0700)
	if err != nil { // this is not an exists error according to docs
		return err
	}

	file, err := os.Create(path)
	if err != nil {
		if !os.IsExist(err) {
			return err
		}
		file, err = os.Open(path)
		if err != nil {
			return err
		}
	}
	encode := json.NewEncoder(file)

	return encode.Encode(&s.fyneSettings)
}

func (s *settings) load() {
	err := s.loadFromFile(s.fyneSettings.StoragePath())
	if err != nil {
		fyne.LogError("Settings load error:", err)
	}
}

func (s *settings) loadFromFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			os.MkdirAll(filepath.Dir(path), 0700)
			return nil
		}
		return err
	}
	decode := json.NewDecoder(file)

	return decode.Decode(&s.fyneSettings)
}

var preview *canvas.Image

func main() {
	s := &settings{}
	s.load()

	a := app.New()
	w := a.NewWindow("Fyne Settings")

	preview = canvas.NewImageFromResource(themeDarkPreview)
	preview.FillMode = canvas.ImageFillContain

	themes := widget.NewSelect([]string{"light", "dark"}, func(name string) {
		s.fyneSettings.ThemeName = name

		switch name {
		case "light":
			preview.Resource = themeLightPreview
		default:
			preview.Resource = themeDarkPreview
		}
		canvas.Refresh(preview)
	})
	themes.SetSelected("dark")
	top := widget.NewForm(
		&widget.FormItem{Text: "Scale", Widget: widget.NewSelect([]string{"Auto"}, func(string) {})},
		&widget.FormItem{Text: "Theme", Widget: themes})
	bottom := widget.NewHBox(layout.NewSpacer(),
		&widget.Button{Text: "Apply", Style: widget.PrimaryButton, OnTapped: func() {
			s.save()
		}})

	appearance := fyne.NewContainerWithLayout(layout.NewBorderLayout(top, bottom, nil, nil),
		top, bottom, preview)

	tabs := widget.NewTabContainer(
		&widget.TabItem{Text: "Appearance", Icon: theme.NewThemedResource(appearanceIcon, nil), Content: appearance},
		&widget.TabItem{Text: "Language", Icon: theme.NewThemedResource(languageIcon, nil), Content: &canvas.Rectangle{FillColor: color.White}})
	tabs.SetTabLocation(widget.TabLocationLeading)
	w.SetContent(tabs)

	w.Resize(fyne.NewSize(480, 320))
	w.ShowAndRun()
}
