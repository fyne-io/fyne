package fyne

// Settings describes the application configuration available.
type Settings interface {
	Theme() Theme
	SetTheme(Theme)

	AddChangeListener(chan Settings)
}
