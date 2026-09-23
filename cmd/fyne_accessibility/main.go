// Fyne Accessibility is a deterministic widget gallery for assistive technology testing.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/theme"
)

func main() {
	scenario := flag.String("scenario", "all", "Screen to show: all, "+strings.Join(scenarioNames(), ", "))
	quitAfter := flag.Duration("quit-after", 0, "Automatically quit after this duration (0 disables the test watchdog)")
	background := flag.Bool("background", false, "Show a non-focusing, mouse-transparent window for macOS automation")
	flag.Parse()
	if flag.NArg() != 0 || !validScenario(*scenario) || *quitAfter < 0 {
		flag.Usage()
		os.Exit(2)
	}

	a := app.NewWithID("io.fyne.accessibility")
	a.Settings().SetTheme(theme.DefaultTheme())
	w := a.NewWindow("Fyne Accessibility - " + *scenario)
	if *background {
		if err := configureBackground(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
	content, err := makeGallery(*scenario, w)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	w.SetContent(content)
	w.Resize(fyne.NewSize(900, 650))
	if *quitAfter > 0 {
		timer := time.AfterFunc(*quitAfter, func() {
			log.Print("Accessibility gallery watchdog expired")
			fyne.Do(a.Quit)
		})
		defer timer.Stop()
	}
	w.ShowAndRun()
}
