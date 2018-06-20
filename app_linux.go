package fyne

import "os/exec"

func (app *fyneApp) OpenURL(url string) {
	exec.Command("xdg-open", url).Run()
}
