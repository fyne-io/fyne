package theme

import "image/color"
import "path"
import "runtime"

func BackgroundColor() color.RGBA {
	return color.RGBA{0xcc, 0xcc, 0xcc, 0xff}
}

func ButtonColor() color.RGBA {
	return color.RGBA{0x80, 0x80, 0x80, 0xff}
}

func TextColor() color.RGBA {
	return color.RGBA{0x10, 0x10, 0x10, 0xff}
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
