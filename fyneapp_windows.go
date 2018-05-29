package fyneapp

import "os/exec"

func (app *eflApp) OpenURL(url string) {
	exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Run()
}
