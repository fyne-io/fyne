//go:build !wasm && !test_web_driver

package widget

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/storage"
	"github.com/yuin/goldmark/ast"
)

func parseMarkdownImage(t *ast.Image) []RichTextSegment {
	dest := string(t.Destination)
	u, err := storage.ParseURI(dest)
	if err != nil {
		u = storage.NewFileURI(dest)
	}
	return []RichTextSegment{&ImageSegment{Source: u, Title: string(t.Title), Alignment: fyne.TextAlignCenter}}
}
