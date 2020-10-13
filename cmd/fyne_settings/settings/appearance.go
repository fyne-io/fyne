package settings

import (
	"encoding/json"
	"image"
	"image/color"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/internal/painter"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/tools/playground"
	"fyne.io/fyne/widget"
)

const (
	systemThemeName = "system default"
)

// Settings gives access to user interfaces to control Fyne settings
type Settings struct {
	fyneSettings app.SettingsSchema

	preview *canvas.Image
	colors  []fyne.CanvasObject
}

// NewSettings returns a new settings instance with the current configuration loaded
func NewSettings() *Settings {
	s := &Settings{}
	s.load()

	return s
}

// AppearanceIcon returns the icon for appearance settings
func (s *Settings) AppearanceIcon() fyne.Resource {
	return theme.NewThemedResource(appearanceIcon, nil)
}

// LoadAppearanceScreen creates a new settings screen to handle appearance configuration
func (s *Settings) LoadAppearanceScreen(w fyne.Window) fyne.CanvasObject {
	s.preview = canvas.NewImageFromImage(s.createPreview())
	s.preview.FillMode = canvas.ImageFillContain

	def := s.fyneSettings.ThemeName
	themeNames := []string{"dark", "light"}
	if runtime.GOOS == "darwin" || runtime.GOOS == "windows" {
		themeNames = append(themeNames, systemThemeName)
		if s.fyneSettings.ThemeName == "" {
			def = systemThemeName
		}
	}
	themes := widget.NewSelect(themeNames, s.chooseTheme)
	themes.SetSelected(def)

	scale := s.makeScaleGroup(w.Canvas().Scale())

	current := s.fyneSettings.PrimaryColor
	if current == "" {
		current = theme.ColorBlue
	}
	for _, c := range theme.PrimaryColorNames() {
		b := newColorButton(c, theme.PrimaryColorNamed(c), s)
		s.colors = append(s.colors, b)
	}
	swatch := fyne.NewContainerWithLayout(layout.NewGridLayout(len(s.colors)), s.colors...)

	scale.Append(widget.NewGroup("Main Color", swatch))
	scale.Append(widget.NewGroup("Theme", themes))

	bottom := widget.NewHBox(layout.NewSpacer(),
		&widget.Button{Text: "Apply", Style: widget.PrimaryButton, OnTapped: func() {
			err := s.save()
			if err != nil {
				fyne.LogError("Failed on saving", err)
			}

			s.appliedScale(s.fyneSettings.Scale)
		}})

	return fyne.NewContainerWithLayout(layout.NewBorderLayout(scale, bottom, nil, nil),
		scale, bottom, s.preview)
}

func (s *Settings) chooseTheme(name string) {
	if name == systemThemeName {
		name = ""
	}
	s.fyneSettings.ThemeName = name

	s.preview.Image = s.createPreview()
	canvas.Refresh(s.preview)
}

type overrideTheme interface {
	OverrideTheme(fyne.Theme, string)
}

func (s *Settings) createPreview() image.Image {
	c := playground.NewSoftwareCanvas()
	oldTheme := fyne.CurrentApp().Settings().Theme()
	oldColor := fyne.CurrentApp().Settings().PrimaryColor()

	th := oldTheme
	if s.fyneSettings.ThemeName == "light" {
		th = theme.LightTheme()
	} else if s.fyneSettings.ThemeName == "dark" {
		th = theme.DarkTheme()
	}

	painter.SvgCacheReset() // reset icon cache
	fyne.CurrentApp().Settings().(overrideTheme).OverrideTheme(th, s.fyneSettings.PrimaryColor)

	empty := widget.NewLabel("")
	tabs := widget.NewTabContainer(
		widget.NewTabItemWithIcon("Home", theme.HomeIcon(), widget.NewLabel("Home")),
		widget.NewTabItemWithIcon("Browse", theme.ComputerIcon(), empty),
		widget.NewTabItemWithIcon("Settings", theme.SettingsIcon(), empty),
		widget.NewTabItemWithIcon("Help", theme.HelpIcon(), empty))
	tabs.SetTabLocation(widget.TabLocationLeading)
	showOverlay(c)

	c.SetContent(tabs)
	c.Resize(fyne.NewSize(380, 380))
	img := c.Capture()

	painter.SvgCacheReset() // ensure we re-create the correct cached assets
	fyne.CurrentApp().Settings().(overrideTheme).OverrideTheme(oldTheme, oldColor)
	return img
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

	return ioutil.WriteFile(path, data, 0644)
}

type colorButton struct {
	widget.BaseWidget
	name  string
	color color.Color

	s *Settings
}

func newColorButton(n string, c color.Color, s *Settings) *colorButton {
	b := &colorButton{name: n, color: c, s: s}
	b.ExtendBaseWidget(b)
	return b
}

func (c *colorButton) CreateRenderer() fyne.WidgetRenderer {
	r := canvas.NewRectangle(color.Transparent)
	r.StrokeWidth = 5

	if c.name == c.s.fyneSettings.PrimaryColor {
		r.StrokeColor = theme.PrimaryColor()
	}

	return &colorRenderer{c: c, rect: r, objs: []fyne.CanvasObject{r}}
}

func (c *colorButton) Tapped(_ *fyne.PointEvent) {
	c.s.fyneSettings.PrimaryColor = c.name
	for _, child := range c.s.colors {
		child.Refresh()
	}

	c.s.preview.Image = c.s.createPreview()
	canvas.Refresh(c.s.preview)
}

type colorRenderer struct {
	c    *colorButton
	rect *canvas.Rectangle
	objs []fyne.CanvasObject
}

func (c *colorRenderer) Layout(s fyne.Size) {
	c.rect.Resize(s)
}

func (c *colorRenderer) MinSize() fyne.Size {
	return fyne.NewSize(20, 20)
}

func (c *colorRenderer) Refresh() {
	if c.c.name == c.c.s.fyneSettings.PrimaryColor {
		c.rect.StrokeColor = theme.PrimaryColor()
	} else {
		c.rect.StrokeColor = color.Transparent
	}

	c.rect.Refresh()
}

func (c *colorRenderer) BackgroundColor() color.Color {
	return c.c.color
}

func (c *colorRenderer) Objects() []fyne.CanvasObject {
	return c.objs
}

func (c *colorRenderer) Destroy() {
}

func showOverlay(c fyne.Canvas) {
	username := widget.NewEntry()
	password := widget.NewPasswordEntry()
	form := widget.NewForm(widget.NewFormItem("Username", username),
		widget.NewFormItem("Password", password))
	form.OnCancel = func() {}
	form.OnSubmit = func() {}
	content := widget.NewVBox(
		widget.NewLabelWithStyle("Login demo", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}), form)
	wrap := fyne.NewContainerWithoutLayout(content)
	wrap.Resize(content.MinSize().Add(fyne.NewSize(theme.Padding()*2, theme.Padding()*2)))
	content.Resize(content.MinSize())
	content.Move(fyne.NewPos(theme.Padding(), theme.Padding()))

	over := fyne.NewContainerWithLayout(layout.NewMaxLayout(),
		canvas.NewRectangle(theme.ShadowColor()), fyne.NewContainerWithLayout(layout.NewCenterLayout(),
			wrap))

	c.Overlays().Add(over)
	c.Focus(username)
}
