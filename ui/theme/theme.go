package theme

import "image/color"
import "os"
import "path"
import "runtime"

var loadedColors *themeColors

type themeColors struct {
	Background color.RGBA

	Button, Text, Primary color.RGBA
}

func loadLightColors() *themeColors {
	return &themeColors{
		Background: color.RGBA{0xff, 0xff, 0xff, 0xff},
		Button:     color.RGBA{0xee, 0xee, 0xee, 0xff},
		Text:       color.RGBA{0x0, 0x0, 0x0, 0xdd},
		Primary:    color.RGBA{0x9f, 0xa8, 0xda, 0xff},
	}
}

func loadDarkColors() *themeColors {
	return &themeColors{
		Background: color.RGBA{0x42, 0x42, 0x42, 0xff},
		Button:     color.RGBA{0x21, 0x21, 0x21, 0xff},
		Text:       color.RGBA{0xff, 0xff, 0xff, 0xff},
		Primary:    color.RGBA{0x1a, 0x23, 0x7e, 0xff},
	}
}

func colors() *themeColors {
	if loadedColors != nil {
		return loadedColors
	}

	env := os.Getenv("FYNE_THEME")
	if env == "light" {
		loadedColors = loadLightColors()
	} else {
		loadedColors = loadDarkColors()
	}

	return loadedColors
}

func BackgroundColor() color.RGBA {
	return colors().Background
}

func ButtonColor() color.RGBA {
	return colors().Button
}

func TextColor() color.RGBA {
	return colors().Text
}

func PrimaryColor() color.RGBA {
	return colors().Primary
}

func TextSize() int {
	return 14
}

func fontPath(style string) string {
	_, dirname, _, _ := runtime.Caller(0)
	return path.Join(path.Dir(dirname), "font/NotoSans-"+style+".ttf")
}

func TextFont() string {
	return fontPath("Regular")
}

func TextBoldFont() string {
	return fontPath("Bold")
}

func TextItalicFont() string {
	return fontPath("Italic")
}

func TextBoldItalicFont() string {
	return fontPath("BoldItalic")
}

func Padding() int {
	return 5
}
