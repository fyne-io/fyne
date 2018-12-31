package widget

import (
	"github.com/fyne-io/fyne"
)

// Label widget is a label component with appropriate padding and layout.
type Label = Text

// NewLabel creates a new layout widget with the set text content
func NewLabel(text string) *Label {
	return NewLabelStyled(text, fyne.TextAlignLeading, fyne.TextStyle{})
}

// NewLabelStyled creates a new layout widget with the set text content
func NewLabelStyled(text string, alignment fyne.TextAlign, style fyne.TextStyle) *Label {
	l := &Label{
		Alignment: alignment,
		TextStyle: style,
	}
	l.SetText(text)
	return l
}
