package settings

import (
	"encoding/json"
	"image"
	"image/color"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/internal/cache"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/tools/playground"
	"fyne.io/fyne/v2/widget"
)

const (
	systemThemeName = "system default"
)

// Settings gives access to user interfaces to control Fyne settings
type Settings struct {
	fyneSettings app.SettingsSchema

	preview *canvas.Image
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
	s.userTheme = fyne.CurrentApp().Settings().Theme()
	if s.userTheme == nil {
		s.userTheme = theme.DefaultTheme()
	}

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
	box := container.NewVBox(scale)

	for _, c := range theme.PrimaryColorNames() {
		b := newColorButton(c, theme.PrimaryColorNamed(c), s)
		s.colors = append(s.colors, b)
	}
	swatch := container.NewGridWithColumns(len(s.colors), s.colors...)
	appearance := widget.NewForm(widget.NewFormItem("Main Color", swatch),
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

	return container.NewBorder(box, bottom, nil, nil, s.preview)
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

	variant := theme.VariantDark
	if s.fyneSettings.ThemeName == "light" {
		variant = theme.VariantLight
	}

	cache.ResetThemeCaches() // reset icon cache
	fyne.CurrentApp().Settings().(overrideTheme).OverrideTheme(&previewTheme{s.userTheme, variant}, s.fyneSettings.PrimaryColor)

	empty := widget.NewLabel("")
	tabs := container.NewAppTabs(
		container.NewTabItemWithIcon("Home", theme.HomeIcon(), widget.NewLabel("Home")),
		container.NewTabItemWithIcon("Browse", theme.ComputerIcon(), empty),
		container.NewTabItemWithIcon("Settings", theme.SettingsIcon(), empty),
		container.NewTabItemWithIcon("Help", theme.HelpIcon(), empty))
	tabs.SetTabLocation(container.TabLocationLeading)
	showOverlay(c)

	c.SetContent(tabs)
	c.Resize(fyne.NewSize(380, 380))
	// wait for indicator animation
	time.Sleep(canvas.DurationShort)
	img := c.Capture()

	cache.ResetThemeCaches() // ensure we re-create the correct cached assets
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
	r := canvas.NewRectangle(c.color)
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
	c.rect.FillColor = c.c.color

	c.rect.Refresh()
}

func (c *colorRenderer) Objects() []fyne.CanvasObject {
	return c.objs
}

func (c *colorRenderer) Destroy() {
}

type previewTheme struct {
	t fyne.Theme
	v fyne.ThemeVariant
}

func (p *previewTheme) Color(n fyne.ThemeColorName, v fyne.ThemeVariant) color.Color {
	return p.t.Color(n, p.v)
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

func showOverlay(c fyne.Canvas) {
	username := widget.NewEntry()
	password := widget.NewPasswordEntry()
	form := widget.NewForm(widget.NewFormItem("Username", username),
		widget.NewFormItem("Password", password))
	form.OnCancel = func() {}
	form.OnSubmit = func() {}
	content := container.NewVBox(
		widget.NewLabelWithStyle("Login demo", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}), form)
	wrap := container.NewWithoutLayout(content)
	wrap.Resize(content.MinSize().Add(fyne.NewSize(theme.Padding()*2, theme.Padding()*2)))
	content.Resize(content.MinSize())
	content.Move(fyne.NewPos(theme.Padding(), theme.Padding()))

	over := container.NewMax(
		canvas.NewRectangle(theme.ShadowColor()), container.NewCenter(wrap),
	)

	c.Overlays().Add(over)
	c.Focus(username)
}
