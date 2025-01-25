//go:build ignore

package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"fyne.io/fyne/v2"
)

func main() {
	symbolFont := "InterSymbols-Regular.ttf"
	err := createFontByStripping(symbolFont, "Inter-Regular.ttf", []rune{
		'←',
		'↑',
		'→',
		'↓',
		'↖',
		'↘',
		'↩',
		'↪',
		'↳',
		'↵',
		'⇞',
		'⇟',
		'⇤',
		'⇥',
		'⇧',
		'⌃',
		'⌘',
		'⌥',
		'⌦',
		'⌫',
		'⎋',
		'␣',
		'❖',
	})
	if err != nil {
		fyne.LogError("symbol font creation failed", err)
		os.Exit(1)
	}
}

func createFontByStripping(newFontFile, fontFile string, runes []rune) error {
	unicodes := make([]string, 0, len(runes))
	for _, r := range runes {
		unicodes = append(unicodes, fmt.Sprintf(`%04X`, r))
	}
	cmd := exec.Command(
		"pyftsubset",
		fontPath(fontFile),
		"--output-file="+fontPath(newFontFile),
		"--unicodes="+strings.Join(unicodes, ","),
	)
	fmt.Println("creating font by executing:", cmd.String())
	output, err := cmd.CombinedOutput()
	if len(output) > 0 {
		fmt.Println("output:")
		fmt.Println(string(output))
	}
	return err
}

func fontPath(filename string) string {
	dirname, _ := os.Getwd()
	return filepath.Join(dirname, "font", filename)
}
