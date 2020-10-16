package fyne

// SettingsScaleAuto is a specific scale value that indicates a canvas should
// scale according to the DPI of the window that contains it.
//
// Deprecated: Automatic scaling is now handled in the drivers and is not a user setting.
const SettingsScaleAuto = float32(-1.0)

// BuildType defines different modes that an application can be built using.
type BuildType int

const (
	// StandardBuild is the normal build mode - it is not debug, test or release mode.
	StandardBuild BuildType = iota
	// DebugBuild is used when a developer would like more information and visual output for app debugging.
	DebugBuild
	// ReleaseBuild is a final production build, it is like StandardBuild but will use distribution certificates.
	// A release build is typically going to connect to live services and is not usually used during development.
	ReleaseBuild
)

// Settings describes the application configuration available.
type Settings interface {
	Theme() Theme
	SetTheme(Theme)
	Scale() float32
	PrimaryColor() string

	AddChangeListener(chan Settings)
	BuildType() BuildType
}
