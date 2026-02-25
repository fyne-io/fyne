package theme

import (
	_ "embed"

	"fyne.io/fyne/v2"
)

//go:embed font/NotoSans-Regular.ttf
var notoSansRegular []byte

var regular = &fyne.StaticResource{
	StaticName:    "NotoSans-Regular.ttf",
	StaticContent: notoSansRegular,
}

//go:embed font/NotoSans-Bold.ttf
var notoSansBold []byte

var bold = &fyne.StaticResource{
	StaticName:    "NotoSans-Bold.ttf",
	StaticContent: notoSansBold,
}

//go:embed font/NotoSans-Italic.ttf
var notoSansItalic []byte

var italic = &fyne.StaticResource{
	StaticName:    "NotoSans-Italic.ttf",
	StaticContent: notoSansItalic,
}

//go:embed font/NotoSans-BoldItalic.ttf
var notoSansBoldItalic []byte

var bolditalic = &fyne.StaticResource{
	StaticName:    "NotoSans-BoldItalic.ttf",
	StaticContent: notoSansBoldItalic,
}

//go:embed font/DejaVuSansMono-Powerline.ttf
var dejaVuSansMono []byte

var monospace = &fyne.StaticResource{
	StaticName:    "DejaVuSansMono-Powerline.ttf",
	StaticContent: dejaVuSansMono,
}

//go:embed font/InterSymbols-Regular.ttf
var interSymbolsRegular []byte

var symbol = &fyne.StaticResource{
	StaticName:    "InterSymbols-Regular.ttf",
	StaticContent: interSymbolsRegular,
}
