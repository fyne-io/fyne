//go:build !ci && (!android || !ios || !mobile) && (wasm || test_web_driver)

package widget

import (
	"strings"
	"syscall/js"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/storage"
	"github.com/yuin/goldmark/ast"
)

func parseMarkdownImage(t *ast.Image) []RichTextSegment {
	dest := string(t.Destination)
	u, err := storage.ParseURI(dest)
	if err != nil {
		if !strings.HasPrefix(dest, "/") {
			dest = "/" + dest
		}
		origin := js.Global().Get("location").Get("origin").String()
		u, err = storage.ParseURI(origin + dest)
		if err != nil {
			fyne.LogError("Can't load image in markdown", err)
			return []RichTextSegment{}
		}
	}
	return []RichTextSegment{&ImageSegment{Source: u, Title: string(t.Title), Alignment: fyne.TextAlignCenter}}
}
