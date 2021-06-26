package widget

import (
	"net/url"
	"strings"

	"github.com/russross/blackfriday/v2"

	"fyne.io/fyne/v2"
)

// NewRichTextFromMarkdown configures a RichText widget by parsing the provided markdown content.
//
// Since: 2.1
func NewRichTextFromMarkdown(content string) *RichText {
	return NewRichText(parseMarkdown(content)...)
}

func parseMarkdown(content string) []RichTextSegment {
	nodes := blackfriday.New().Parse([]byte(content))
	var segs []RichTextSegment

	var nextSeg RichTextSegment
	nodes.Walk(func(node *blackfriday.Node, entering bool) blackfriday.WalkStatus {
		if !entering {
			if text, ok := nextSeg.(*TextSegment); ok && text.Style == RichTextStyleInline {
				text.Style = RichTextStyleParagraph
			}
			nextSeg = &TextSegment{
				Style: RichTextStyleInline,
			}

			return blackfriday.GoToNext
		}

		switch node.Type {
		case blackfriday.Heading:
			switch node.HeadingData.Level {
			case 1:
				nextSeg = &TextSegment{
					Style: RichTextStyleHeading,
					Text:  string(node.Literal),
				}
			case 2:
				nextSeg = &TextSegment{
					Style: RichTextStyleSubHeading,
					Text:  string(node.Literal),
				}
			}
		case blackfriday.HorizontalRule:
			segs = append(segs, &SeparatorSegment{})
		case blackfriday.Link:
			link, _ := url.Parse(string(node.LinkData.Destination))
			nextSeg = &HyperlinkSegment{fyne.TextAlignLeading, strings.TrimSpace(string(node.LinkData.Title)), link}
		case blackfriday.Paragraph:
			nextSeg = &TextSegment{
				Style: RichTextStyleInline, // we make it a paragraph at the end if there are no more elements
				Text:  string(node.Literal),
			}
		case blackfriday.Code:
			segs = append(segs, &TextSegment{
				Style: RichTextStyleCodeInline,
				Text:  string(node.Literal),
			})
			nextSeg = &TextSegment{
				Style: RichTextStyleInline, // we make it a paragraph at the end if there are no more elements
				Text:  string(node.Literal),
			}
		case blackfriday.Emph:
			nextSeg = &TextSegment{
				Style: RichTextStyleEmphasis,
				Text:  string(node.Literal),
			}
		case blackfriday.Strong:
			nextSeg = &TextSegment{
				Style: RichTextStyleStrong,
				Text:  string(node.Literal),
			}
		case blackfriday.Text:
			trimmed := string(node.Literal)
			trimmed = strings.ReplaceAll(trimmed, "\n", " ") // newline inside paragraph is not newline
			if trimmed == "" {
				return blackfriday.GoToNext
			}
			if text, ok := nextSeg.(*TextSegment); ok {
				text.Text = trimmed
			}
			if link, ok := nextSeg.(*HyperlinkSegment); ok {
				link.Text = trimmed
			}
			segs = append(segs, nextSeg)
		}

		return blackfriday.GoToNext
	})
	return segs
}
