package container

import (
	"testing"

	internalTest "fyne.io/fyne/v2/internal/test"
	"github.com/stretchr/testify/assert"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/internal/cache"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func TestTabButton_Icon_Change(t *testing.T) {
	b := &tabButton{icon: theme.CancelIcon()}
	r := cache.Renderer(b)
	icon := r.Objects()[3].(*canvas.Image)
	oldResource := icon.Resource

	b.icon = theme.ConfirmIcon()
	b.Refresh()
	assert.NotEqual(t, oldResource, icon.Resource)
}

func TestTabButton_ContentFitsMinSize(t *testing.T) {
	for name, th := range map[string]fyne.Theme{
		"default": theme.DefaultTheme(),
		"test":    test.Theme(),
		"custom":  test.NewTheme(),
	} {
		t.Run(name, func(t *testing.T) {
			a := test.NewTempApp(t)
			a.Settings().SetTheme(th)
			for _, tt := range []struct {
				name     string
				position buttonIconPosition
				text     string
				closable bool
			}{
				{"top icon", buttonIconTop, "", false},
				{"top icon and text", buttonIconTop, "Tab", false},
				{"top icon with close", buttonIconTop, "", true},
				{"top icon and text with close", buttonIconTop, "Tab", true},
				{"inline icon", buttonIconInline, "", false},
				{"inline icon and text", buttonIconInline, "Tab", false},
			} {
				t.Run(tt.name, func(t *testing.T) {
					b := &tabButton{icon: theme.InfoIcon(), iconPosition: tt.position, text: tt.text}
					if tt.closable {
						b.onClosed = func() {}
					}
					r := test.TempWidgetRenderer(t, b).(*tabButtonRenderer)
					r.Refresh()
					b.Resize(b.MinSize())

					assert.GreaterOrEqual(t, r.icon.Position().X, float32(0))
					assert.GreaterOrEqual(t, r.icon.Position().Y, float32(0))
					assert.LessOrEqual(t, r.icon.Position().X+r.icon.Size().Width, b.Size().Width)
					assert.LessOrEqual(t, r.icon.Position().Y+r.icon.Size().Height, b.Size().Height)
					if tt.text != "" {
						assert.LessOrEqual(t, r.label.Position().Y+r.label.Size().Height,
							b.Size().Height-r.padding().Height/2)
					}
				})
			}
		})
	}
}

func TestAppTabs_IconOnlyFitsButton(t *testing.T) {
	test.NewTempApp(t)
	for _, mode := range []struct {
		name   string
		mobile bool
	}{
		{"desktop", false},
		{"mobile", true},
	} {
		for name, location := range map[string]TabLocation{
			"top":      TabLocationTop,
			"leading":  TabLocationLeading,
			"bottom":   TabLocationBottom,
			"trailing": TabLocationTrailing,
		} {
			t.Run(name+" "+mode.name, func(t *testing.T) {
				item := NewTabItemWithIcon("", theme.InfoIcon(), widget.NewLabel("Content"))
				tabs := NewAppTabs(item)
				override := NewThemeOverride(tabs, test.Theme())
				override.SetDeviceIsMobile(mode.mobile)
				tabs.SetTabLocation(location)
				w := test.NewTempWindow(t, override)
				w.SetPadded(false)
				w.Resize(fyne.NewSize(150, 150))

				b := item.button
				r := cache.Renderer(b).(*tabButtonRenderer)
				assert.LessOrEqual(t, r.icon.Position().Y+r.icon.Size().Height, b.Size().Height)
			})
		}
	}
}

func TestTab_ThemeChange(t *testing.T) {
	a := test.NewTempApp(t)
	a.Settings().SetTheme(internalTest.LightTheme(theme.DefaultTheme()))

	tabs := NewAppTabs(
		NewTabItem("a", widget.NewLabel("a")),
		NewTabItem("b", widget.NewLabel("b")),
	)
	w := test.NewTempWindow(t, tabs)
	w.Resize(fyne.NewSize(180, 120))

	initial := w.Canvas().Capture()

	a.Settings().SetTheme(internalTest.DarkTheme(theme.DefaultTheme()))
	tabs.SelectIndex(1)
	second := w.Canvas().Capture()
	assert.NotEqual(t, initial, second)

	a.Settings().SetTheme(internalTest.LightTheme(theme.DefaultTheme()))
	tabs.SelectIndex(0)
	assert.Equal(t, initial, w.Canvas().Capture())
}
