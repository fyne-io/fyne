package ui

type Window interface {
	Title() string
	SetTitle(string)
	Show()
	Hide()
	Close()

	Canvas() Canvas
}
