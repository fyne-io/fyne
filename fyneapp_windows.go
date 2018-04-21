package fyneapp

import "os/exec"

func (app *eflApp) OpenURL(url string) {
	exec.Command("start", "", url).Run()
}
