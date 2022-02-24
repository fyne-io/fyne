package commands

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"

	"fyne.io/fyne/v2/cmd/fyne/internal/templates"
	"fyne.io/fyne/v2/cmd/fyne/internal/util"
)

type webData struct {
	WasmFile     string
	GopherJSFile string
	IsReleased   bool
	HasWasm      bool
	HasGopherJS  bool
}

func packageWebInternal(w webData, appDir string, exeWasmSrc string, exeJSSrc string, icon string, release bool) error {
	index := filepath.Join(appDir, "index.html")
	indexFile, err := os.Create(index)
	if err != nil {
		return err
	}

	err = templates.IndexHTML.Execute(indexFile, w)
	if err != nil {
		return err
	}

	iconDst := filepath.Join(appDir, "icon.png")
	err = util.CopyFile(icon, iconDst)
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
		err := ioutil.WriteFile(webglDebugFile, templates.WebGLDebugJs, 0644)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *Packager) packageWeb() error {
	appDir := util.EnsureSubDir(p.dir, "web")

	tpl := webData{
		GopherJSFile: p.name + ".js",
		WasmFile:     p.name + ".wasm",
		IsReleased:   p.release,
		HasGopherJS:  true,
		HasWasm:      true,
	}

	return packageWebInternal(tpl, appDir, p.exe+".wasm", p.exe+".js", p.icon, p.release)
}

func (p *Packager) packageWasm() error {
	appDir := util.EnsureSubDir(p.dir, "wasm")

	tpl := webData{
		WasmFile:   p.name,
		IsReleased: p.release,
		HasWasm:    true,
	}

	return packageWebInternal(tpl, appDir, p.exe, "", p.icon, p.release)
}

func (p *Packager) packageGopherJS() error {
	appDir := util.EnsureSubDir(p.dir, "gopherjs")

	tpl := webData{
		GopherJSFile: p.name,
		IsReleased:   p.release,
		HasGopherJS:  true,
	}

	return packageWebInternal(tpl, appDir, "", p.exe, p.icon, p.release)
}
