//go:build linux
// +build linux

package theme

import (
	"image/color"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
)

// KDETheme theme is based on the KDETheme or Plasma theme.
type KDETheme struct {
	variant         fyne.ThemeVariant
	bgColor         color.Color
	fgColor         color.Color
	viewColor       color.Color
	buttonColor     color.Color
	buttonAlternate color.Color
	fontConfig      string
	fontSize        float32
	font            fyne.Resource
}

// Color returns the color for the specified name.
//
// Implements: fyne.Theme
func (k *KDETheme) Color(name fyne.ThemeColorName, _ fyne.ThemeVariant) color.Color {
	switch name {
	case ColorNameBackground:
		return k.bgColor
	case ColorNameForeground:
		return k.fgColor
	case ColorNameButton:
		return k.buttonColor
	case ColorNameDisabledButton:
		return k.buttonAlternate
	case ColorNameInputBackground:
		return k.viewColor

	}
	return DefaultTheme().Color(name, k.variant)
}

// Font returns the font for the specified name.
//
// Implements: fyne.Theme
func (k *KDETheme) Font(s fyne.TextStyle) fyne.Resource {
	if k.font != nil {
		return k.font
	}
	return DefaultTheme().Font(s)
}

// Icon returns the icon for the specified name.
//
// Implements: fyne.Theme
func (k *KDETheme) Icon(i fyne.ThemeIconName) fyne.Resource {
	return DefaultTheme().Icon(i)
}

// Size returns the size of the font for the specified text style.
//
// Implements: fyne.Theme
func (k *KDETheme) Size(s fyne.ThemeSizeName) float32 {
	if s == SizeNameText {
		return k.fontSize
	}
	return DefaultTheme().Size(s)
}

// decodeTheme initialize the theme.
func (k *KDETheme) decodeTheme() error {
	if err := k.loadScheme(); err != nil {
		return err
	}
	k.setFont()
	return nil
}

// loadScheme loads the KDE theme from kdeglobals if it is found.
func (k *KDETheme) loadScheme() error {
	// the theme name is declared in ~/.config/kdedefaults/kdeglobals
	// in the ini section [General] as "ColorScheme" entry
	homedir, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	content, err := ioutil.ReadFile(filepath.Join(homedir, ".config/kdeglobals"))
	if err != nil {
		return err
	}

	section := ""
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "[") {
			section = strings.ReplaceAll(line, "[", "")
			section = strings.ReplaceAll(section, "]", "")
		}
		if section == "Colors:Window" {
			if strings.HasPrefix(line, "BackgroundNormal=") {
				k.bgColor = k.parseColor(strings.ReplaceAll(line, "BackgroundNormal=", ""))
			}
			if strings.HasPrefix(line, "ForegroundNormal=") {
				k.fgColor = k.parseColor(strings.ReplaceAll(line, "ForegroundNormal=", ""))
			}
		}
		if section == "Colors:Button" {
			if strings.HasPrefix(line, "BackgroundNormal=") {
				k.buttonColor = k.parseColor(strings.ReplaceAll(line, "BackgroundNormal=", ""))
			}
			if strings.HasPrefix(line, "BackgroundAlternate=") {
				k.buttonAlternate = k.parseColor(strings.ReplaceAll(line, "BackgroundAlternate=", ""))
			}
		}
		if section == "Colors:View" {
			if strings.HasPrefix(line, "BackgroundNormal=") {
				k.viewColor = k.parseColor(strings.ReplaceAll(line, "BackgroundNormal=", ""))
			}
		}
		if section == "General" {
			if strings.HasPrefix(line, "font=") {
				k.fontConfig = strings.ReplaceAll(line, "font=", "")
			}
		}
	}

	return nil
}

// parseColor parses a color from a string in form r,g,b or r,g,b,a.
func (k *KDETheme) parseColor(col string) color.Color {
	// the color is in the form r,g,b,
	// we need to convert it to a color.Color

	// split the string
	cols := strings.Split(col, ",")
	// convert the string to int
	r, _ := strconv.Atoi(cols[0])
	g, _ := strconv.Atoi(cols[1])
	b, _ := strconv.Atoi(cols[2])
	a := 0xff
	if len(cols) > 3 {
		a, _ = strconv.Atoi(cols[3])
	}

	// convert the int to a color.Color
	return color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)}
}

// setFont sets the font for the theme.
func (k *KDETheme) setFont() {

	if k.fontConfig == "" {
		return
	}
	// the fontline is in the form "fontline,size,...", so we can split it
	fontline := strings.Split(k.fontConfig, ",")
	name := fontline[0]
	size, _ := strconv.ParseFloat(fontline[1], 32)
	k.fontSize = float32(size)

	// we need to load the font, Gnome struct has got some nice methods
	fontpath, err := getFontPath(name)
	if err != nil {
		log.Println(err)
		return
	}

	var font []byte
	if filepath.Ext(fontpath) == ".ttf" {
		font, err = ioutil.ReadFile(fontpath)
		if err != nil {
			log.Println(err)
			return
		}
	} else {
		font, err = converToTTF(fontpath)
		if err != nil {
			log.Println(err)
			return
		}
	}
	k.font = fyne.NewStaticResource(fontpath, font)
}

// NewKDETheme returns a new KDE theme.
func NewKDETheme() fyne.Theme {
	kde := &KDETheme{
		variant: VariantDark,
	}
	if err := kde.decodeTheme(); err != nil {
		log.Println(err)
		return DefaultTheme()
	}

	return kde
}
