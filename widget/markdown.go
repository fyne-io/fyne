package widget

import (
	"html"
	"io"
	"net/url"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	ast2 "github.com/yuin/goldmark/extension/ast"
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
	segs, err := renderNode(source, n, 0, 0)
	*m = segs
	return err
}

func renderNode(source []byte, n ast.Node, quotingDepth int, listDepth int) ([]RichTextSegment, error) {
	switch t := n.(type) {
	case *ast.Document:
		return renderChildren(source, n, quotingDepth, listDepth)
	case *ast.Paragraph:
		children, err := renderChildren(source, n, quotingDepth, listDepth)
		linebreak := &TextSegment{Style: RichTextStyleParagraph}
		children = append(children, linebreak)
		return children, err
	case *ast.List:
		items, err := renderChildren(source, n, quotingDepth, listDepth+1)
		return []RichTextSegment{
			&ListSegment{startIndex: t.Start - 1, Items: items, Ordered: t.Marker != '*' && t.Marker != '-' && t.Marker != '+', indentationLevel: listDepth},
		}, err
	case *ast.ListItem:
		children, err := renderChildren(source, n, quotingDepth, listDepth)
		var texts []RichTextSegment
		var sublist RichTextSegment
		for _, child := range children {
			// check if child is a sub-list
			if _, ok := child.(*ListSegment); ok {
				sublist = child
			} else {
				texts = append(texts, child)
			}
		}
		result := []RichTextSegment{&ParagraphSegment{Texts: texts}}
		if sublist != nil {
			result = append(result, sublist)
		}
		return result, err
	case *ast.TextBlock:
		if c, ok := t.FirstChild().(*ast2.TaskCheckBox); ok {
			child := c.NextSibling()
			text := ""
			if child != nil {
				text = string(child.(*ast.Text).Value(source))
			}
			return []RichTextSegment{&CheckBoxSegment{Text: decodeText(text), Checked: c.IsChecked}}, nil
		}
		return renderChildren(source, n, quotingDepth, listDepth)
	case *ast.Heading:
		return renderHeading(source, n, quotingDepth, listDepth)
	case *ast.ThematicBreak:
		return []RichTextSegment{&SeparatorSegment{}}, nil
	case *ast.Link:
		link, _ := url.Parse(string(t.Destination))
		text := forceIntoText(source, n)
		return []RichTextSegment{&HyperlinkSegment{Alignment: fyne.TextAlignLeading, Text: decodeText(text), URL: link}}, nil
	case *ast.AutoLink:
		link, _ := url.Parse(string(t.URL(source)))
		text := string(t.Label(source))
		return []RichTextSegment{&HyperlinkSegment{Alignment: fyne.TextAlignLeading, Text: decodeText(text), URL: link}}, nil
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
		return renderEmphasis(source, n, quotingDepth, n.(*ast.Emphasis).Level, listDepth)
	case *ast2.Strikethrough:
		return renderEmphasis(source, n, quotingDepth, 3, listDepth)
	case *ast.Text:
		text := string(t.Value(source))
		if text == "" {
			// These empty text elements indicate single line breaks after non-text elements in goldmark.
			return []RichTextSegment{&TextSegment{Style: RichTextStyleInline, Text: " "}}, nil
		}
		if n.(*ast.Text).SoftLineBreak() {
			text = text + " "
		}
		if quotingDepth > 0 {
			style := RichTextStyleBlockquote
			style.QuotingDepth = quotingDepth
			return []RichTextSegment{&TextSegment{Style: style, Text: text}}, nil
		}
		return []RichTextSegment{&TextSegment{Style: RichTextStyleInline, Text: decodeText(text)}}, nil
	case *ast.Blockquote:
		quotingDepth++
		return renderChildren(source, n, quotingDepth, listDepth)
	case *ast.Image:
		return parseMarkdownImage(t), nil
	}
	return nil, nil
}

func renderChildren(source []byte, n ast.Node, quotingDepth int, listDepth int) ([]RichTextSegment, error) {
	children := make([]RichTextSegment, 0, n.ChildCount())
	for childCount, child := n.ChildCount(), n.FirstChild(); childCount > 0; childCount-- {
		segs, err := renderNode(source, child, quotingDepth, listDepth)
		if err != nil {
			return children, err
		}
		children = append(children, segs...)
		child = child.NextSibling()
	}
	return children, nil
}

func renderEmphasis(source []byte, n ast.Node, quotingDepth int, strength, listDepth int) ([]RichTextSegment, error) {
	style := RichTextStyleInline
	switch strength {
	case 1:
		style = RichTextStyleEmphasis
		if _, ok := n.Parent().(*ast2.Strikethrough); ok {
			style.TextStyle.Strikethrough = true
		}
	case 2:
		style = RichTextStyleStrong
		if _, ok := n.Parent().(*ast2.Strikethrough); ok {
			style.TextStyle.Strikethrough = true
		}
	case 3:
		style = RichTextStyleInline
		style.TextStyle.Strikethrough = true
		if emp, ok := n.Parent().(*ast.Emphasis); ok {
			if emp.Level == 1 {
				style.TextStyle.Italic = true
			} else if emp.Level == 2 {
				style.TextStyle.Bold = true
			}
		}
	}

	children, err := renderChildren(source, n, quotingDepth, listDepth)
	for _, child := range children {
		switch t := child.(type) {
		case *TextSegment:
			t.Style.TextStyle.Bold = t.Style.TextStyle.Bold || style.TextStyle.Bold
			t.Style.TextStyle.Italic = t.Style.TextStyle.Italic || style.TextStyle.Italic
			t.Style.TextStyle.Strikethrough = t.Style.TextStyle.Strikethrough || style.TextStyle.Strikethrough
		case *HyperlinkSegment:
			t.TextStyle.Bold = t.TextStyle.Bold || style.TextStyle.Bold
			t.TextStyle.Italic = t.TextStyle.Italic || style.TextStyle.Italic
			t.TextStyle.Strikethrough = t.TextStyle.Strikethrough || style.TextStyle.Strikethrough
		}
	}
	return children, err
}

func renderHeading(source []byte, n ast.Node, quotingDepth int, listDepth int) ([]RichTextSegment, error) {
	var style RichTextStyle
	switch n.(*ast.Heading).Level {
	case 1:
		style = RichTextStyleHeading
	case 2:
		style = RichTextStyleSubHeading
	default:
		style = RichTextStyleStrong
	}

	children := make([]RichTextSegment, 0, n.ChildCount())
	for childCount, child := n.ChildCount(), n.FirstChild(); childCount > 0; childCount-- {
		switch t := child.(type) {
		case *ast.Text:
			text := string(t.Value(source))
			children = append(children, &TextSegment{Style: style, Text: decodeText(text)})
		default:
			segs, err := renderNode(source, child, quotingDepth, listDepth)
			if err != nil {
				return children, err
			}
			for _, seg := range segs {
				if t, ok := seg.(*TextSegment); ok { // apply heading to other text
					t.Style.SizeName = style.SizeName
					t.Style.TextStyle.Bold = true
				}
			}
			children = append(children, segs...)
		}
		child = child.NextSibling()
	}
	if len(children) == 0 {
		children = append(children, &TextSegment{Style: style, Text: ""})
	}

	for _, child := range children {
		switch t := child.(type) {
		case *HyperlinkSegment:
			t.TextStyle = style.TextStyle
			t.SizeName = style.SizeName
		}
	}
	linebreak := &TextSegment{Style: RichTextStyleParagraph}
	children = append(children, linebreak)
	return children, nil
}

func forceIntoText(source []byte, n ast.Node) string {
	text := strings.Builder{}
	ast.Walk(n, func(n2 ast.Node, entering bool) (ast.WalkStatus, error) {
		if entering {
			switch t := n2.(type) {
			case *ast.Text:
				text.Write(t.Value(source))
				text.WriteByte(' ')
			}
		}
		return ast.WalkContinue, nil
	})
	return strings.TrimSuffix(text.String(), " ")
}

func parseMarkdown(content string) []RichTextSegment {
	r := markdownRenderer{}
	md := goldmark.New(goldmark.WithRenderer(&r), goldmark.WithExtensions(extension.Strikethrough, extension.TaskList))
	err := md.Convert([]byte(content), nil)
	if err != nil {
		fyne.LogError("Failed to parse markdown", err)
	}
	return r
}

func decodeText(text string) string {
	return html.UnescapeString(text)
}
