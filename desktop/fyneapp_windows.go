package fyneapp

import "os/exec"

import _ "github.com/fyne-io/fyne/desktop/efl" // import our Fyne driver for Windows

func (app *fyneApp) OpenURL(url string) {
	exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Run()
}
