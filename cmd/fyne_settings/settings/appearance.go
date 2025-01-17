package settings

import (
	"encoding/json"
	"image/color"
	"os"
	"path/filepath"
	"runtime"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	internalapp "fyne.io/fyne/v2/internal/app"
	internaltheme "fyne.io/fyne/v2/internal/theme"
	intWidget "fyne.io/fyne/v2/internal/widget"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

const (
	themeNameDark        = "dark"
	themeNameLight       = "light"
	themeNameSystem      = ""
	themeNameSystemLabel = "system default"
)

// Settings gives access to user interfaces to control Fyne settings
type Settings struct {
	fyneSettings app.SettingsSchema

	preview *fyne.Container
	colors  []fyne.CanvasObject

	userTheme fyne.Theme
}

// NewSettings returns a new settings instance with the current configuration loaded
func NewSettings() *Settings {
	s := &Settings{}
	s.load()
	if s.fyneSettings.Scale == 0 {
		s.fyneSettings.Scale = 1
	}
	return s
}

// AppearanceIcon returns the icon for appearance settings
func (s *Settings) AppearanceIcon() fyne.Resource {
	return theme.NewThemedResource(resourceAppearanceSvg)
}

// LoadAppearanceScreen creates a new settings screen to handle appearance configuration
func (s *Settings) LoadAppearanceScreen(w fyne.Window) fyne.CanvasObject {
	fallback := fyne.CurrentApp().Settings().Theme()
	if fallback == nil {
		fallback = theme.DefaultTheme()
	}
	s.userTheme = &previewTheme{s: s, t: fallback}
	s.preview = s.createPreview()

	def := s.fyneSettings.ThemeName
	themeNames := []string{themeNameDark, themeNameLight}
	if runtime.GOOS == "darwin" || runtime.GOOS == "windows" {
		themeNames = append(themeNames, themeNameSystemLabel)
		if s.fyneSettings.ThemeName == themeNameSystem {
			def = themeNameSystemLabel
		}
	}
	themes := widget.NewSelect(themeNames, s.chooseTheme)
	themes.SetSelected(def)

	scale := s.makeScaleGroup(w.Canvas().Scale())
	box := container.NewVBox(scale)

	animations := widget.NewCheck("Animate widgets", func(on bool) {
		s.fyneSettings.DisableAnimations = !on
	})
	animations.Checked = !s.fyneSettings.DisableAnimations
	for _, c := range theme.PrimaryColorNames() {
		b := newPrimaryColorButton(c, s)
		s.colors = append(s.colors, b)
	}
	swatch := container.NewGridWithColumns(len(s.colors), s.colors...)
	appearance := widget.NewForm(
		widget.NewFormItem("Animations", animations),
		widget.NewFormItem("Main Color", swatch),
		widget.NewFormItem("Theme", themes))

	box.Add(widget.NewCard("Appearance", "", appearance))
	bottom := container.NewHBox(layout.NewSpacer(),
		&widget.Button{Text: "Apply", Importance: widget.HighImportance, OnTapped: func() {
			if s.fyneSettings.Scale == 0.0 {
				s.chooseScale(1.0)
			}
			err := s.save()
			if err != nil {
				fyne.LogError("Failed on saving", err)
			}

			s.appliedScale(s.fyneSettings.Scale)
		}})

	return container.NewBorder(box, bottom, nil, nil,
		container.NewCenter(s.preview))
}

func (s *Settings) chooseTheme(name string) {
	if name == themeNameSystemLabel {
		name = themeNameSystem
	}
	s.fyneSettings.ThemeName = name

	s.refreshPreview()
}

func (s *Settings) createPreview() *fyne.Container {
	preview := createPreviewWidget()

	content := container.NewThemeOverride(preview, s.userTheme)
	bg := canvas.NewRectangle(s.userTheme.Color(theme.ColorNameBackground, theme.VariantDark))
	content.Refresh()

	return container.NewStack(bg, content)
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

func (s *Settings) refreshPreview() {
	s.preview.Objects[0].(*canvas.Rectangle).FillColor = s.userTheme.Color(theme.ColorNameBackground, theme.VariantDark)
	s.preview.Refresh()
}

func (s *Settings) save() error {
	return s.saveToFile(s.fyneSettings.StoragePath())
}

func (s *Settings) saveToFile(path string) error {
	err := os.MkdirAll(filepath.Dir(path), 0700)
	if err != nil { // this is not an exists error according to docs
		return err
	}

	data, err := json.Marshal(&s.fyneSettings)
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

type primaryColorButton struct {
	widget.BaseWidget
	name string

	s *Settings
}

func newPrimaryColorButton(name string, s *Settings) *primaryColorButton {
	b := &primaryColorButton{name: name, s: s}
	b.ExtendBaseWidget(b)
	return b
}

func (c *primaryColorButton) CreateRenderer() fyne.WidgetRenderer {
	r := canvas.NewRectangle(internaltheme.PrimaryColorNamed(c.name))
	r.CornerRadius = theme.SelectionRadiusSize()
	r.StrokeWidth = 5

	if c.name == c.s.fyneSettings.PrimaryColor {
		r.StrokeColor = theme.Color(theme.ColorNamePrimary)
	}

	return &primaryColorButtonRenderer{c: c, rect: r, objs: []fyne.CanvasObject{r}}
}

func (c *primaryColorButton) Tapped(_ *fyne.PointEvent) {
	c.s.fyneSettings.PrimaryColor = c.name
	for _, child := range c.s.colors {
		child.Refresh()
	}

	c.s.refreshPreview()
}

type primaryColorButtonRenderer struct {
	c    *primaryColorButton
	rect *canvas.Rectangle
	objs []fyne.CanvasObject
}

func (c *primaryColorButtonRenderer) Layout(s fyne.Size) {
	c.rect.Resize(s)
}

func (c *primaryColorButtonRenderer) MinSize() fyne.Size {
	return fyne.NewSize(20, 32)
}

func (c *primaryColorButtonRenderer) Refresh() {
	if c.c.name == c.c.s.fyneSettings.PrimaryColor {
		c.rect.StrokeColor = theme.Color(theme.ColorNamePrimary)
	} else {
		c.rect.StrokeColor = color.Transparent
	}
	c.rect.FillColor = internaltheme.PrimaryColorNamed(c.c.name)
	c.rect.CornerRadius = theme.SelectionRadiusSize()

	c.rect.Refresh()
}

func (c *primaryColorButtonRenderer) Objects() []fyne.CanvasObject {
	return c.objs
}

func (c *primaryColorButtonRenderer) Destroy() {
}

type previewTheme struct {
	s *Settings
	t fyne.Theme
}

func (p *previewTheme) Color(n fyne.ThemeColorName, _ fyne.ThemeVariant) color.Color {
	variant := theme.VariantDark
	switch p.s.fyneSettings.ThemeName {
	case themeNameLight:
		variant = theme.VariantLight
	case themeNameSystem:
		variant = internalapp.DefaultVariant()
	}

	switch n {
	case theme.ColorNamePrimary:
		return internaltheme.PrimaryColorNamed(p.s.fyneSettings.PrimaryColor)
	case theme.ColorNameForegroundOnPrimary:
		return internaltheme.ForegroundOnPrimaryColorNamed(p.s.fyneSettings.PrimaryColor)
	}

	return p.t.Color(n, variant)
}

func (p *previewTheme) Font(s fyne.TextStyle) fyne.Resource {
	return p.t.Font(s)
}

func (p *previewTheme) Icon(n fyne.ThemeIconName) fyne.Resource {
	return p.t.Icon(n)
}

func (p *previewTheme) Size(n fyne.ThemeSizeName) float32 {
	return p.t.Size(n)
}

func createPreviewWidget() fyne.CanvasObject {
	username := widget.NewEntry()
	password := widget.NewPasswordEntry()
	form := widget.NewForm(widget.NewFormItem("Username", username),
		widget.NewFormItem("Password", password))
	form.OnCancel = func() {}
	form.OnSubmit = func() {}
	content := container.NewVBox(
		widget.NewLabelWithStyle("Theme Preview", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}), form)

	over := container.NewStack(intWidget.NewShadow(intWidget.ShadowAround, intWidget.DialogLevel),
		container.NewPadded(content))

	return over
}
