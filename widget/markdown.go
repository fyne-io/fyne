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
	blockquote  bool
	heading     bool
	nextSeg     RichTextSegment
	parentStack [][]RichTextSegment
	segs        []RichTextSegment
}

func (m *markdownRenderer) AddOptions(...renderer.Option) {}

func (m *markdownRenderer) Render(_ io.Writer, source []byte, n ast.Node) error {
	m.nextSeg = &TextSegment{}
	err := ast.Walk(n, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			if n.Kind().String() == "Heading" {
				m.segs = append(m.segs, m.nextSeg)
				m.heading = false
			}
			return ast.WalkContinue, m.handleExitNode(n)
		}

		switch n.Kind().String() {
		case "List":
			// prepare a new child level
			m.parentStack = append(m.parentStack, m.segs)
			m.segs = nil
		case "ListItem":
			// prepare a new item level
			m.parentStack = append(m.parentStack, m.segs)
			m.segs = nil
		case "Heading":
			m.heading = true
			switch n.(*ast.Heading).Level {
			case 1:
				m.nextSeg = &TextSegment{
					Style: RichTextStyleHeading,
				}
			case 2:
				m.nextSeg = &TextSegment{
					Style: RichTextStyleSubHeading,
				}
			default:
				m.nextSeg = &TextSegment{
					Style: RichTextStyleParagraph,
				}
				m.nextSeg.(*TextSegment).Style.TextStyle.Bold = true
			}
		case "HorizontalRule", "ThematicBreak":
			m.segs = append(m.segs, &SeparatorSegment{})
		case "Link":
			link, _ := url.Parse(string(n.(*ast.Link).Destination))
			m.nextSeg = &HyperlinkSegment{fyne.TextAlignLeading, "", link}
		case "Paragraph":
			m.nextSeg = &TextSegment{
				Style: RichTextStyleInline, // we make it a paragraph at the end if there are no more elements
			}
			if m.blockquote {
				m.nextSeg.(*TextSegment).Style = RichTextStyleBlockquote
			}
		case "CodeSpan":
			m.nextSeg = &TextSegment{
				Style: RichTextStyleCodeInline,
			}
		case "CodeBlock", "FencedCodeBlock":
			var data []byte
			lines := n.Lines()
			for i := 0; i < lines.Len(); i++ {
				line := lines.At(i)
				data = append(data, line.Value(source)...)
			}
			if len(data) == 0 {
				return ast.WalkContinue, nil
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
				m.nextSeg = &TextSegment{
					Style: RichTextStyleStrong,
				}
			default:
				m.nextSeg = &TextSegment{
					Style: RichTextStyleEmphasis,
				}
			}
		case "Strong":
			m.nextSeg = &TextSegment{
				Style: RichTextStyleStrong,
			}
		case "Text":
			trimmed := string(n.Text(source))
			trimmed = strings.ReplaceAll(trimmed, "\n", " ") // newline inside paragraph is not newline
			if trimmed == "" {
				return ast.WalkContinue, nil
			}
			if t, ok := m.nextSeg.(*TextSegment); ok {
				next := n.(*ast.Text).NextSibling()
				if next != nil {
					if nextText, ok := next.(*ast.Text); ok {
						if nextText.Segment.Start > n.(*ast.Text).Segment.Stop { // detect presence of a trailing newline
							trimmed = trimmed + " "
						}
					}
				}

				t.Text = t.Text + trimmed
			}
			if link, ok := m.nextSeg.(*HyperlinkSegment); ok {
				link.Text = link.Text + trimmed
			}

			if !m.heading {
				m.segs = append(m.segs, m.nextSeg)
			}
		case "Blockquote":
			m.blockquote = true
		}

		return ast.WalkContinue, nil
	})
	return err
}

func (m *markdownRenderer) handleExitNode(n ast.Node) error {
	if n.Kind().String() == "Blockquote" {
		m.blockquote = false
	} else if n.Kind().String() == "List" {
		listSegs := m.segs
		m.segs = m.parentStack[len(m.parentStack)-1]
		m.parentStack = m.parentStack[:len(m.parentStack)-1]
		marker := n.(*ast.List).Marker
		m.segs = append(m.segs, &ListSegment{Items: listSegs, Ordered: marker != '*' && marker != '-' && marker != '+'})
	} else if n.Kind().String() == "ListItem" {
		itemSegs := m.segs
		m.segs = m.parentStack[len(m.parentStack)-1]
		m.parentStack = m.parentStack[:len(m.parentStack)-1]
		m.segs = append(m.segs, &ParagraphSegment{Texts: itemSegs})
	} else if !m.blockquote && !m.heading {
		if len(m.segs) > 0 {
			if text, ok := m.segs[len(m.segs)-1].(*TextSegment); ok && n.Kind().String() == "Paragraph" {
				text.Style.Inline = false
			}
		}
		m.nextSeg = &TextSegment{
			Style: RichTextStyleInline,
		}
	}
	return nil
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
