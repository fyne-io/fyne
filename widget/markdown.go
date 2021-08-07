package widget

import (
	"io"
	"net/url"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/renderer"

	"fyne.io/fyne/v2"
)

// NewRichTextFromMarkdown configures a RichText widget by parsing the provided markdown content.
//
// Since: 2.1
func NewRichTextFromMarkdown(content string) *RichText {
	return NewRichText(parseMarkdown(content)...)
}

// ParseMarkdown allows setting the content of this RichText widget from a markdown string.
// It will replace the content of this widget similarly to SetText, but with the appropriate formatting.
func (t *RichText) ParseMarkdown(content string) {
	t.Segments = parseMarkdown(content)
	t.Refresh()
}

type markdownRenderer struct {
	segs []RichTextSegment
}

func (m *markdownRenderer) AddOptions(...renderer.Option) {}

func (m *markdownRenderer) Render(_ io.Writer, source []byte, n ast.Node) error {
	m.segs = nil
	var nextSeg RichTextSegment
	blockquote := false
	return ast.Walk(n, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			if n.Kind().String() == "Blockquote" {
				blockquote = false
			} else if !blockquote {
				if text, ok := m.segs[len(m.segs)-1].(*TextSegment); ok && n.Kind().String() == "Paragraph" {
					text.Style = RichTextStyleParagraph
				}
				nextSeg = &TextSegment{
					Style: RichTextStyleInline,
				}
			}
			return ast.WalkContinue, nil
		}

		switch n.Kind().String() {
		case "Heading":
			switch n.(*ast.Heading).Level {
			case 1:
				nextSeg = &TextSegment{
					Style: RichTextStyleHeading,
					Text:  string(n.Text(source)),
				}
			case 2:
				nextSeg = &TextSegment{
					Style: RichTextStyleSubHeading,
					Text:  string(n.Text(source)),
				}
			}
		case "HorizontalRule", "ThematicBreak":
			m.segs = append(m.segs, &SeparatorSegment{})
		case "Link":
			link, _ := url.Parse(string(n.(*ast.Link).Destination))
			nextSeg = &HyperlinkSegment{fyne.TextAlignLeading, strings.TrimSpace(string(n.Text(source))), link}
		case "Paragraph":
			nextSeg = &TextSegment{
				Style: RichTextStyleInline, // we make it a paragraph at the end if there are no more elements
				Text:  string(n.Text(source)),
			}
			if blockquote {
				nextSeg.(*TextSegment).Style = RichTextStyleBlockquote
			}
		case "CodeSpan":
			nextSeg = &TextSegment{
				Style: RichTextStyleCodeInline,
				Text:  string(n.Text(source)),
			}
		case "CodeBlock", "FencedCodeBlock":
			var data []byte
			lines := n.Lines()
			for i := 0; i < lines.Len(); i++ {
				line := lines.At(i)
				data = append(data, line.Value(source)...)
			}
			if data[len(data)-1] == '\n' {
				data = data[:len(data)-1]
			}
			m.segs = append(m.segs, &TextSegment{
				Style: RichTextStyleCodeBlock,
				Text:  string(data),
			})
		case "Emph", "Emphasis":
			switch n.(*ast.Emphasis).Level {
			case 2:
				nextSeg = &TextSegment{
					Style: RichTextStyleStrong,
					Text:  string(n.Text(source)),
				}
			default:
				nextSeg = &TextSegment{
					Style: RichTextStyleEmphasis,
					Text:  string(n.Text(source)),
				}
			}
		case "Strong":
			nextSeg = &TextSegment{
				Style: RichTextStyleStrong,
				Text:  string(n.Text(source)),
			}
		case "Text":
			trimmed := string(n.Text(source))
			trimmed = strings.ReplaceAll(trimmed, "\n", " ") // newline inside paragraph is not newline
			if trimmed == "" {
				return ast.WalkContinue, nil
			}
			if text, ok := nextSeg.(*TextSegment); ok {
				text.Text = trimmed
			}
			if link, ok := nextSeg.(*HyperlinkSegment); ok {
				link.Text = trimmed
			}
			m.segs = append(m.segs, nextSeg)
		case "Blockquote":
			blockquote = true
		case "Document":
		default:
			fyne.LogError("Unrecognised MarkDown type "+n.Kind().String(), nil)
		}

		return ast.WalkContinue, nil
	})
}

func parseMarkdown(content string) []RichTextSegment {
	r := &markdownRenderer{}
	if content == "" {
		return r.segs
	}

	md := goldmark.New(goldmark.WithRenderer(r))
	err := md.Convert([]byte(content), nil)
	if err != nil {
		fyne.LogError("Failed to parse markdown", err)
	}
	return r.segs
}
