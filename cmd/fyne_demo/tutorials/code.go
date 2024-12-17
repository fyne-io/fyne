package tutorials

import (
	"image/color"
	"log"

	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"

	"fyne.io/fyne/v2/widget"
)

func highlightTextGrid(grid *widget.TextGrid) {
	content := grid.Text()
	lex := lexers.Get("go")

	style := styles.Get("solarized-dark")
	if style == nil {
		style = styles.Fallback
	}

	iterator, err := lex.Tokenise(nil, string(content))
	if err != nil {
		log.Println("Token error", err)
		return
	}

	row, col := 0, 0
	textColor := styleColor(chroma.Background, style)
	grid.SetRowStyle(0, &widget.CustomTextGridStyle{
		FGColor: textColor})
	for _, tok := range iterator.Tokens() {
		length := len(tok.Value)

		if tok.Value == "\n" {
			row++
			col = 0

			grid.SetRowStyle(row, &widget.CustomTextGridStyle{
				FGColor: textColor})

			continue
		}

		c := resolveColor(style.Get(tok.Type).Colour)
		grid.SetStyleRange(row, col, row, col+length,
			&widget.CustomTextGridStyle{FGColor: c})
		col += length
	}
}

func resolveColor(colour chroma.Colour) color.Color {
	r, g, b := colour.Red(), colour.Green(), colour.Blue()

	return &color.NRGBA{R: r, G: g, B: b, A: 0xff}
}

func styleBackgroundColor(name chroma.TokenType, style *chroma.Style) color.Color {
	entry := style.Get(name)
	return resolveColor(entry.Background)
}

func styleColor(name chroma.TokenType, style *chroma.Style) color.Color {
	entry := style.Get(name)
	return resolveColor(entry.Colour)
}
