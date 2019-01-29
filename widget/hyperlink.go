package widget

import (
	"net/url"
	"os/exec"
	"runtime"

	"fyne.io/fyne"
)

// Hyperlink widget is a text component with appropriate padding and layout.
// When clicked, the default web browser should open with a URL
type Hyperlink struct {
	textWidget
	Text      string
	Url       string
	Alignment fyne.TextAlign // The alignment of the Text
	TextStyle fyne.TextStyle // The style of the hyperlink text
}

// NewHyperlink creates a new layout widget with the set text content
func NewHyperlink(text string, sUrl string) *Hyperlink {
	// ignore the error here, we'll continue anyway
	return NewHyperlinkWithStyle(text, sUrl, fyne.TextAlignLeading, fyne.TextStyle{})
}

// NewHyperlinkWithStyle creates a new layout widget with the set text content
func NewHyperlinkWithStyle(text string, sUrl string, alignment fyne.TextAlign, style fyne.TextStyle) *Hyperlink {
	hl := &Hyperlink{
		Text:      text,
		Url:       sUrl,
		Alignment: alignment,
		TextStyle: style,
	}

	hl.TextType = TextWidgetType_Hyperlink

	Renderer(hl).Refresh()
	return hl
}

// SetText sets the text of the hyperlink
func (hl *Hyperlink) SetText(text string) {
	hl.Text = text
	hl.textWidget.SetText(text)
	Renderer(hl).Refresh()
}

// SetUrl sets the URL of the hyperlink, taking in a string type
func (hl *Hyperlink) SetUrl(sUrl string) {
	hl.Url = sUrl
}

// SetUrl sets the URL of the hyperlink, taking in a Url type
func (hl *Hyperlink) SetUrlFromUrl(uUrl *url.URL) {
	hl.Url = uUrl.String()
}

// OnMouseDown is called when a mouse down event is captured and triggers any change handler
func (hl *Hyperlink) OnMouseDown(*fyne.MouseEvent) {
	if hl.Url != "" {
		open(hl.Url)
	}
}

// CreateRenderer is a private method to Fyne which links this widget to it's renderer
func (hl *Hyperlink) CreateRenderer() fyne.WidgetRenderer {
	hl.textWidget = textWidget{
		TextType:  hl.TextType,
		Alignment: hl.Alignment,
		TextStyle: hl.TextStyle,
	}
	hl.textWidget.SetText(hl.Text)
	r := hl.textWidget.CreateRenderer()
	r.Refresh()
	return r
}

// Resize sets a new size for a widget.
// Note this should not be used if the widget is being managed by a Layout within a Container.
func (hl *Hyperlink) Resize(size fyne.Size) {
	hl.resize(size, hl)
}

// Move the widget to a new position, relative to it's parent.
// Note this should not be used if the widget is being managed by a Layout within a Container.
func (hl *Hyperlink) Move(pos fyne.Position) {
	hl.move(pos, hl)
}

// MinSize returns the smallest size this widget can shrink to
func (hl *Hyperlink) MinSize() fyne.Size {
	return hl.minSize(hl)
}

// Show this widget, if it was previously hidden
func (hl *Hyperlink) Show() {
	hl.show(hl)
}

// Hide this widget, if it was previously visible
func (hl *Hyperlink) Hide() {
	hl.hide(hl)
}

// taken from https://github.com/icza/gowut/blob/085680a418c9a92dcf2c72d48df2df3cf2a1e88f/gwu/server_start.go#L29
// open opens the specified URL in the default browser of the user.
func open(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}
