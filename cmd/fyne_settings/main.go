package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
)

type settings struct {
	fyneSettings app.SettingsSchema

	preview *canvas.Image
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

func (s *settings) chooseTheme(name string) {
	s.fyneSettings.ThemeName = name

	switch name {
	case "light":
		s.preview.Resource = themeLightPreview
	default:
		s.preview.Resource = themeDarkPreview
	}
	canvas.Refresh(s.preview)
}

func (s *settings) chooseScale(value string) {
	if value == "" {
		s.fyneSettings.Scale = fyne.SettingsScaleAuto
		return
	}

	scale, err := strconv.ParseFloat(value, 32)
	if err != nil {
		log.Println("Cannot set scale to:", value)
	}
	s.fyneSettings.Scale = float32(scale)
}

func main() {
	s := &settings{}
	s.load()

	a := app.New()
	w := a.NewWindow("Fyne Settings")

	s.preview = canvas.NewImageFromResource(themeDarkPreview)
	s.preview.FillMode = canvas.ImageFillContain

	def := s.fyneSettings.ThemeName
	themes := widget.NewSelect([]string{"dark", "light"}, s.chooseTheme)
	themes.SetSelected(def)
	scale := widget.NewEntry()
	scale.OnChanged = s.chooseScale
	scale.SetPlaceHolder("Auto")
	if s.fyneSettings.Scale != fyne.SettingsScaleAuto {
		scale.SetText(fmt.Sprintf("%.2f", s.fyneSettings.Scale))
	}

	top := widget.NewForm(
		&widget.FormItem{Text: "Scale", Widget: scale},
		&widget.FormItem{Text: "Theme", Widget: themes})
	bottom := widget.NewHBox(layout.NewSpacer(),
		&widget.Button{Text: "Apply", Style: widget.PrimaryButton, OnTapped: func() {
			s.save()
		}})

	appearance := fyne.NewContainerWithLayout(layout.NewBorderLayout(top, bottom, nil, nil),
		top, bottom, s.preview)

	tabs := widget.NewTabContainer(
		&widget.TabItem{Text: "Appearance", Icon: theme.NewThemedResource(appearanceIcon, nil), Content: appearance})
	tabs.SetTabLocation(widget.TabLocationLeading)
	w.SetContent(tabs)

	w.Resize(fyne.NewSize(480, 320))
	w.ShowAndRun()
}
