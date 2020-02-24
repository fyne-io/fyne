package fyne

import "unicode"

// TextAlign represents the horizontal alignment of text within a widget or
// canvas object.
type TextAlign int

const (
	// TextAlignLeading specifies a left alignment for left-to-right languages.
	TextAlignLeading TextAlign = iota
	// TextAlignCenter places the text centrally within the available space.
	TextAlignCenter
	// TextAlignTrailing will align the text right for a left-to-right language.
	TextAlignTrailing
)

// TextWrap represents how text longer than the widget's width will be wrapped.
type TextWrap int

const (
	// TextWrapOff extends the widget's width to fit the text, no wrapping is applied.
	TextWrapOff TextWrap = iota
	// TextTruncate trims the text to the widget's width, no wrapping is applied.
	TextTruncate
	// TextWrapBreak trims the line of characters to the widget's width adding the excess as new line.
	TextWrapBreak
	// TextWrapWord trims the line of words to the widget's width adding the excess as new line.
	TextWrapWord
)

// TextStyle represents the styles that can be applied to a text canvas object
// or text based widget.
type TextStyle struct {
	Bold      bool // Should text be bold
	Italic    bool // Should text be italic
	Monospace bool // Use the system monospace font instead of regular
}

// MeasureText uses the current driver to calculate the size of text when rendered.
func MeasureText(text string, size int, style TextStyle) Size {
	return CurrentApp().Driver().RenderedTextSize(text, size, style)
}

// SplitLines accepts a slice of runes and returns a slice containing the
// start and end indicies of each line delimited by the newline character.
func SplitLines(text []rune) [][2]int {
	var low, high int
	var lines [][2]int
	length := len(text)
	for i := 0; i < length; i++ {
		c := text[i]
		if c == '\r' {
			// Ignore Carriage Return unless followed by New Line
			i++
			if i < length {
				c = text[i]
			} else {
				break
			}
		}
		if c == '\n' {
			high = i
			lines = append(lines, [2]int{low, high})
			low = i + 1
		}
	}
	return append(lines, [2]int{low, length})
}

// LastWordBoundary accepts a slice of runes and returns the last index of a word delimiter.
func LastWordBoundary(text []rune) int {
	i := len(text) - 1
	for ; i >= 0; i-- {
		if unicode.IsSpace(text[i]) {
			break
		}
	}
	return i
}

// LineBounds accepts a slice of runes, a wrapping mode, a maximum line width and a function to measure line width.
// LineBounds returns a slice containing the start and end indicies of each line with the given wrapping applied.
func LineBounds(text []rune, wrap TextWrap, maxWidth int, measurer func([]rune) int) [][2]int {
	var bounds [][2]int
	lines := SplitLines(text)
	switch wrap {
	case TextTruncate:
		if maxWidth > 0 {
			for _, l := range lines {
				low := l[0]
				high := l[1]
				if low == high {
					bounds = append(bounds, l)
				} else {
					for {
						size := measurer(text[low:high])
						if size <= maxWidth {
							bounds = append(bounds, [2]int{low, high})
							break
						} else {
							high--
						}
					}
				}
			}
			return bounds
		}
		// No truncation
		fallthrough
	case TextWrapBreak:
		if maxWidth > 0 {
			for _, l := range lines {
				low := l[0]
				high := l[1]
				if low == high {
					bounds = append(bounds, l)
				} else {
					for low < high {
						size := measurer(text[low:high])
						if size <= maxWidth {
							bounds = append(bounds, [2]int{low, high})
							low = high
							high = l[1]
						} else {
							high--
						}
					}
				}
			}
			return bounds
		}
		// No wrapping
		fallthrough
	case TextWrapWord:
		if maxWidth > 0 {
			for _, l := range lines {
				low := l[0]
				high := l[1]
				if low == high {
					bounds = append(bounds, l)
				} else {
					for low < high {
						sub := text[low:high]
						size := measurer(sub)
						if size <= maxWidth {
							bounds = append(bounds, [2]int{low, high})
							low = high
							high = l[1]
							if low < high && unicode.IsSpace(text[low]) {
								low++
							}
						} else {
							last := LastWordBoundary(sub)
							if last < 0 {
								high--
							} else {
								high = low + last
							}
						}
					}
				}
			}
			return bounds
		}
		// No wrapping
		fallthrough
	case TextWrapOff:
		// No wrapping
		fallthrough
	default:
		return lines
	}
}
