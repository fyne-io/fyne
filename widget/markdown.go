package widget

import (
	"io"
	"net/url"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/renderer"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/storage"
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

// AppendMarkdown parses the given markdown string and appends the
// content to the widget, with the appropriate formatting.
// This API is intended for appending complete markdown documents or
// standalone fragments, and should not be used to parse a single
// markdown document piecewise.
//
// Since: 2.5
func (t *RichText) AppendMarkdown(content string) {
	t.Segments = append(t.Segments, parseMarkdown(content)...)
	t.Refresh()
}

type markdownRenderer []RichTextSegment

func (m *markdownRenderer) AddOptions(...renderer.Option) {}

func (m *markdownRenderer) Render(_ io.Writer, source []byte, n ast.Node) error {
	segs, err := renderNode(source, n, false)
	*m = segs
	return err
}

func renderNode(source []byte, n ast.Node, blockquote bool) ([]RichTextSegment, error) {
	switch t := n.(type) {
	case *ast.Document:
		return renderChildren(source, n, blockquote)
	case *ast.Paragraph:
		children, err := renderChildren(source, n, blockquote)
		if !blockquote {
			linebreak := &TextSegment{Style: RichTextStyleParagraph}
			children = append(children, linebreak)
		}
		return children, err
	case *ast.List:
		items, err := renderChildren(source, n, blockquote)
		return []RichTextSegment{
			&ListSegment{Items: items, Ordered: t.Marker != '*' && t.Marker != '-' && t.Marker != '+'},
		}, err
	case *ast.ListItem:
		texts, err := renderChildren(source, n, blockquote)
		return []RichTextSegment{&ParagraphSegment{Texts: texts}}, err
	case *ast.TextBlock:
		return renderChildren(source, n, blockquote)
	case *ast.Heading:
		text := forceIntoHeadingText(source, n)
		switch t.Level {
		case 1:
			return []RichTextSegment{&TextSegment{Style: RichTextStyleHeading, Text: text}}, nil
		case 2:
			return []RichTextSegment{&TextSegment{Style: RichTextStyleSubHeading, Text: text}}, nil
		default:
			textSegment := TextSegment{Style: RichTextStyleParagraph, Text: text}
			textSegment.Style.TextStyle.Bold = true
			return []RichTextSegment{&textSegment}, nil
		}
	case *ast.ThematicBreak:
		return []RichTextSegment{&SeparatorSegment{}}, nil
	case *ast.Link:
		link, _ := url.Parse(string(t.Destination))
		text := forceIntoText(source, n)
		return []RichTextSegment{&HyperlinkSegment{Alignment: fyne.TextAlignLeading, Text: text, URL: link}}, nil
	case *ast.CodeSpan:
		text := forceIntoText(source, n)
		return []RichTextSegment{&TextSegment{Style: RichTextStyleCodeInline, Text: text}}, nil
	case *ast.CodeBlock, *ast.FencedCodeBlock:
		var data []byte
		lines := n.Lines()
		for i := 0; i < lines.Len(); i++ {
			line := lines.At(i)
			data = append(data, line.Value(source)...)
		}
		if len(data) == 0 {
			return nil, nil
		}
		if data[len(data)-1] == '\n' {
			data = data[:len(data)-1]
		}
		return []RichTextSegment{&TextSegment{Style: RichTextStyleCodeBlock, Text: string(data)}}, nil
	case *ast.Emphasis:
		text := string(forceIntoText(source, n))
		switch t.Level {
		case 2:
			return []RichTextSegment{&TextSegment{Style: RichTextStyleStrong, Text: text}}, nil
		default:
			return []RichTextSegment{&TextSegment{Style: RichTextStyleEmphasis, Text: text}}, nil
		}
	case *ast.Text:
		text := string(t.Text(source))
		if text == "" {
			// These empty text elements indicate single line breaks after non-text elements in goldmark.
			return []RichTextSegment{&TextSegment{Style: RichTextStyleInline, Text: " "}}, nil
		}
		text = suffixSpaceIfAppropriate(text, n)
		if blockquote {
			return []RichTextSegment{&TextSegment{Style: RichTextStyleBlockquote, Text: text}}, nil
		}
		return []RichTextSegment{&TextSegment{Style: RichTextStyleInline, Text: text}}, nil
	case *ast.Blockquote:
		return renderChildren(source, n, true)
	case *ast.Image:
		dest := string(t.Destination)
		u, err := storage.ParseURI(dest)
		if err != nil {
			u = storage.NewFileURI(dest)
		}
		return []RichTextSegment{&ImageSegment{Source: u, Title: string(t.Title), Alignment: fyne.TextAlignCenter}}, nil
	}
	return nil, nil
}

func suffixSpaceIfAppropriate(text string, n ast.Node) string {
	next := n.NextSibling()
	if next != nil && next.Type() == ast.TypeInline && !strings.HasSuffix(text, " ") {
		return text + " "
	}
	return text
}

func renderChildren(source []byte, n ast.Node, blockquote bool) ([]RichTextSegment, error) {
	children := make([]RichTextSegment, 0, n.ChildCount())
	for childCount, child := n.ChildCount(), n.FirstChild(); childCount > 0; childCount-- {
		segs, err := renderNode(source, child, blockquote)
		if err != nil {
			return children, err
		}
		children = append(children, segs...)
		child = child.NextSibling()
	}
	return children, nil
}

func forceIntoText(source []byte, n ast.Node) string {
	texts := make([]string, 0)
	ast.Walk(n, func(n2 ast.Node, entering bool) (ast.WalkStatus, error) {
		if entering {
			switch t := n2.(type) {
			case *ast.Text:
				texts = append(texts, string(t.Text(source)))
			}
		}
		return ast.WalkContinue, nil
	})
	return strings.Join(texts, " ")
}

func forceIntoHeadingText(source []byte, n ast.Node) string {
	var text strings.Builder
	ast.Walk(n, func(n2 ast.Node, entering bool) (ast.WalkStatus, error) {
		if entering {
			switch t := n2.(type) {
			case *ast.Text:
				text.Write(t.Text(source))
			}
		}
		return ast.WalkContinue, nil
	})
	return text.String()
}

func parseMarkdown(content string) []RichTextSegment {
	r := markdownRenderer{}
	md := goldmark.New(goldmark.WithRenderer(&r))
	err := md.Convert([]byte(content), nil)
	if err != nil {
		fyne.LogError("Failed to parse markdown", err)
	}
	return r
}
