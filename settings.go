package fyne

// SettingsScaleAuto is a specific scale value that indicates a canvas should
// scale according to the DPI of the window that contains it.
// Deprecated: Automatic scaling is now handled in the drivers and is not a user setting.
const SettingsScaleAuto = float32(-1.0)

// Settings describes the application configuration available.
type Settings interface {
	Theme() Theme
	SetTheme(Theme)
	Scale() float32

	AddChangeListener(chan Settings)
}
