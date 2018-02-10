package ui

type Window interface {
	Title() string
	SetTitle(string)
}
