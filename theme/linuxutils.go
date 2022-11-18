//go:build linux
// +build linux

package theme

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"fyne.io/fyne/v2"
)

// convertSVGtoPNG will convert a SVG file to a PNG file. It will try to detect if some common tools to convert SVG exists,
// like inkscape or ImageMagik "convert". If not, the resource is not converted.
func convertSVGtoPNG(filename string) (fyne.Resource, error) {
	tmpfile, err := ioutil.TempFile("", "fyne-theme-gnome-*.png")
	if err != nil {
		return nil, err
	}
	defer os.Remove(tmpfile.Name())

	pngConverterOptions := map[string][]string{
		"inkscape": {"--without-gui", "--export-type=png", "--export-background-opacity=0", filename, "-o", tmpfile.Name()},
		"convert":  {"-background", "transparent", "-flatten", filename, tmpfile.Name()},
	}

	var commandName string
	var opts []string
	for binary, options := range pngConverterOptions {
		if path, err := exec.LookPath(binary); err == nil {
			commandName = path
			opts = options
			break
		}
	}

	if commandName == "" {
		return nil, errors.New("you must install inkscape or imageMagik (convert command) to be able to convert SVG icons to PNG")
	}

	// convert the svg to png, no background
	log.Println("Converting", filename, "to", tmpfile.Name())
	cmd := exec.Command(commandName, opts...)

	err = cmd.Run()
	if err != nil {
		return nil, err
	}

	content, err := ioutil.ReadFile(tmpfile.Name())
	if err != nil {
		return nil, err
	}

	return fyne.NewStaticResource(tmpfile.Name(), content), nil
}

// converToTTF will convert a font to a ttf file. This requires the fontforge package.
func converToTTF(fontpath string) ([]byte, error) {

	// check if fontforge is installed
	fontforge, err := exec.LookPath("fontforge")
	if err != nil {
		return nil, err
	}

	// convert the font to a ttf file
	basename := filepath.Base(fontpath)
	tempTTF := filepath.Join(os.TempDir(), "fyne-"+basename+".ttf")

	// Convert to TTF, this is the FF script to call
	ffScript := `Open("%s");Generate("%s")`
	script := fmt.Sprintf(ffScript, fontpath, tempTTF)

	// call fontforge
	cmd := exec.Command(fontforge, "-c", script)
	cmd.Env = append(cmd.Env, "FONTFORGE_LANGUAGE=ff")

	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Println(err)
		log.Println(string(out))
		return nil, err
	}
	defer os.Remove(tempTTF)

	// read the temporary ttf file
	return ioutil.ReadFile(tempTTF)
}

// getFontPath will detect the font path from the font name taken from gsettings.
// As the font is not exactly the one that fc-match can find, we need to do some
// extra work to rebuild the name with style.
func getFontPath(fontname string) (string, error) {

	// check if fc-list and fc-match are installed
	fcList, err := exec.LookPath("fc-list")
	if err != nil {
		return "", err
	}

	fcMatch, err := exec.LookPath("fc-match")
	if err != nil {
		return "", err
	}

	// This to transoform CamelCase to Camel-Case
	camelRegExp := regexp.MustCompile(`([a-z\-])([A-Z])`)

	// get all possible styles in fc-list
	allstyles := []string{}
	cmd := exec.Command(fcList, "--format", "%{style}\n")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	styles := strings.Split(string(out), "\n")
	for _, style := range styles {
		if style != "" {
			split := strings.Split(style, ",")
			for _, s := range split {
				allstyles = append(allstyles, s)
				// we also need to add a "-" for camel cases
				s = camelRegExp.ReplaceAllString(s, "$1-$2")
				allstyles = append(allstyles, s)
			}
		}
	}

	// Find the styles, remove it from the nmae, this make a correct fc-match query
	fontstyle := []string{}
	for _, style := range allstyles {
		if strings.Contains(fontname, " "+style) {
			fontstyle = append(fontstyle, style)
			fontname = strings.ReplaceAll(fontname, style, "")
		}
	}

	// we can now search
	// fc-match ... "Font Name:Font Style
	var fontpath string
	cmd = exec.Command(fcMatch, "-f", "%{file}", fontname+":"+strings.Join(fontstyle, " "))
	out, err = cmd.CombinedOutput()
	if err != nil {
		log.Println(err)
		log.Println(string(out))
		return "", err
	}

	// get the font path with fc-list command
	fontpath = string(out)
	fontpath = strings.TrimSpace(fontpath)
	return fontpath, nil
}
