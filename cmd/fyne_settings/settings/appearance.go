package settings

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
)

// Settings gives access to user interfaces to control Fyne settings
type Settings struct {
	fyneSettings app.SettingsSchema

	preview *canvas.Image
}

func (s *Settings) save() error {
	return s.saveToFile(s.fyneSettings.StoragePath())
}

func (s *Settings) saveToFile(path string) error {
	err := os.MkdirAll(filepath.Dir(path), 0700)
	if err != nil { // this is not an exists error according to docs
		return err
	}

	file, err := os.Create(path)
	if err != nil {
		if !os.IsExist(err) {
			return err
		}
		file, err = os.Open(path) // #nosec
		if err != nil {
			return err
		}
	}
	encode := json.NewEncoder(file)

	return encode.Encode(&s.fyneSettings)
}

func (s *Settings) load() {
	err := s.loadFromFile(s.fyneSettings.StoragePath())
	if err != nil {
		fyne.LogError("Settings load error:", err)
	}
}

func (s *Settings) loadFromFile(path string) error {
	file, err := os.Open(path) // #nosec
	if err != nil {
		if os.IsNotExist(err) {
			err := os.MkdirAll(filepath.Dir(path), 0700)
			if err != nil {
				return err
			}
			return nil
		}
		return err
	}
	decode := json.NewDecoder(file)

	return decode.Decode(&s.fyneSettings)
}

func (s *Settings) chooseTheme(name string) {
	s.fyneSettings.ThemeName = name

	switch name {
	case "light":
		s.preview.Resource = themeLightPreview
	default:
		s.preview.Resource = themeDarkPreview
	}
	canvas.Refresh(s.preview)
}

func (s *Settings) chooseScale(value string) {
	if value == "" || strings.EqualFold(value, "auto") {
		s.fyneSettings.Scale = fyne.SettingsScaleAuto
		return
	}

	scale, err := strconv.ParseFloat(value, 32)
	if err != nil {
		log.Println("Cannot set scale to:", value)
	}
	s.fyneSettings.Scale = float32(scale)
}

// LoadAppearanceScreen creates a new settings screen to handle appearance configuration
func (s *Settings) LoadAppearanceScreen() fyne.CanvasObject {
	s.preview = canvas.NewImageFromResource(themeDarkPreview)
	s.preview.FillMode = canvas.ImageFillContain

	def := s.fyneSettings.ThemeName
	themes := widget.NewSelect([]string{"dark", "light"}, s.chooseTheme)
	themes.SetSelected(def)
	scale := widget.NewEntry()
	scale.SetText("Auto")
	if s.fyneSettings.Scale != fyne.SettingsScaleAuto {
		scale.SetText(fmt.Sprintf("%.2f", s.fyneSettings.Scale))
	}
	scale.OnChanged = s.chooseScale

	top := widget.NewForm(
		&widget.FormItem{Text: "Scale", Widget: scale},
		&widget.FormItem{Text: "Theme", Widget: themes})
	bottom := widget.NewHBox(layout.NewSpacer(),
		&widget.Button{Text: "Apply", Style: widget.PrimaryButton, OnTapped: func() {
			err := s.save()
			if err != nil {
				fyne.LogError("Failed on saving", err)
			}
		}})

	return fyne.NewContainerWithLayout(layout.NewBorderLayout(top, bottom, nil, nil),
		top, bottom, s.preview)
}

// AppearanceIcon returns the icon for appearance settings
func (s *Settings) AppearanceIcon() fyne.Resource {
	return theme.NewThemedResource(appearanceIcon, nil)
}

// NewSettings returns a new settings instance with the current configuration loaded
func NewSettings() *Settings {
	s := &Settings{}
	s.load()

	return s
}
