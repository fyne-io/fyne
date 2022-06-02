package commands

import (
	"bytes"
	"path/filepath"
	"runtime"

	"fyne.io/fyne/v2/cmd/fyne/internal/templates"
)

func (p *Packager) packageWeb() error {
	appDir := util.EnsureSubDir(p.dir, "web")

	tpl := webData{
		AppName:      p.Name,
		AppVersion:   p.AppVersion,
		GopherJSFile: p.Name + ".js",
		WasmFile:     p.Name + ".wasm",
		IsReleased:   p.release,
		HasGopherJS:  true,
		HasWasm:      true,
	}

	return tpl.packageWebInternal(appDir, p.exe+".wasm", p.exe+".js", p.icon, p.release)
}

func (p *Packager) packageWasm() error {
	appDir := util.EnsureSubDir(p.dir, "wasm")

	tpl := webData{
		AppName:    p.Name,
		AppVersion: p.AppVersion,
		WasmFile:   p.Name + ".wasm",
		IsReleased: p.release,
		HasWasm:    true,
	}

	return tpl.packageWebInternal(appDir, p.exe, "", p.icon, p.release)
}

func (p *Packager) packageGopherJS() error {
	appDir := util.EnsureSubDir(p.dir, "gopherjs")

	tpl := webData{
		AppName:      p.Name,
		AppVersion:   p.AppVersion,
		GopherJSFile: p.Name + ".js",
		IsReleased:   p.release,
		HasGopherJS:  true,
	}

	return tpl.packageWebInternal(appDir, "", p.exe, p.icon, p.release)
}

type webData struct {
	AppName      string
	AppVersion   string
	WasmFile     string
	GopherJSFile string
	IsReleased   bool
	HasWasm      bool
	HasGopherJS  bool
}

func (w webData) packageWebInternal(appDir string, exeWasmSrc string, exeJSSrc string, icon string, release bool) error {
	var tpl bytes.Buffer
	err := templates.IndexHTML.Execute(&tpl, w)
	if err != nil {
		return err
	}

	index := filepath.Join(appDir, "index.html")
	err = util.WriteFile(index, tpl.Bytes())
	if err != nil {
		return err
	}

	iconDst := filepath.Join(appDir, "icon.png")
	err = util.CopyFile(icon, iconDst)
	if err != nil {
		return err
	}

	spinnerLightFile := filepath.Join(appDir, "spinner_light.gif")
	err = util.WriteFile(spinnerLightFile, templates.SpinnerLight)
	if err != nil {
		return err
	}

	spinnerDarkFile := filepath.Join(appDir, "spinner_dark.gif")
	err = util.WriteFile(spinnerDarkFile, templates.SpinnerDark)
	if err != nil {
		return err
	}

	lightCSSFile := filepath.Join(appDir, "light.css")
	err = util.WriteFile(lightCSSFile, templates.CSSLight)
	if err != nil {
		return err
	}

	darkCSSFile := filepath.Join(appDir, "dark.css")
	err = util.WriteFile(darkCSSFile, templates.CSSDark)
	if err != nil {
		return err
	}

	if w.HasGopherJS {
		exeJSDst := filepath.Join(appDir, w.GopherJSFile)
		err = util.CopyFile(exeJSSrc, exeJSDst)
		if err != nil {
			return err
		}
	}

	if w.HasWasm {
		wasmExecSrc := filepath.Join(runtime.GOROOT(), "misc", "wasm", "wasm_exec.js")
		wasmExecDst := filepath.Join(appDir, "wasm_exec.js")
		err = util.CopyFile(wasmExecSrc, wasmExecDst)
		if err != nil {
			return err
		}

		exeWasmDst := filepath.Join(appDir, w.WasmFile)
		err = util.CopyFile(exeWasmSrc, exeWasmDst)
		if err != nil {
			return err
		}
	}

	// Download webgl-debug.js directly from the KhronosGroup repository when needed
	if !release {
		webglDebugFile := filepath.Join(appDir, "webgl-debug.js")
		err := util.WriteFile(webglDebugFile, templates.WebGLDebugJs)
		if err != nil {
			return err
		}
	}

	return nil
}
