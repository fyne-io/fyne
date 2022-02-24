package commands

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"

	"fyne.io/fyne/v2/cmd/fyne/internal/templates"
	"fyne.io/fyne/v2/cmd/fyne/internal/util"
)

type indexData struct {
	WasmFile     string
	GopherJSFile string
	IsReleased   bool
	HasWasm      bool
	HasGopherJS  bool
}

func (p *Packager) packageWasm() error {
	appDir := util.EnsureSubDir(p.dir, "wasm")

	index := filepath.Join(appDir, "index.html")
	indexFile, err := os.Create(index)
	if err != nil {
		return err
	}

	tplData := indexData{WasmFile: p.name, IsReleased: p.release, HasWasm: true}
	err = templates.IndexHTML.Execute(indexFile, tplData)
	if err != nil {
		return err
	}

	iconDst := filepath.Join(appDir, "icon.png")
	err = util.CopyFile(p.icon, iconDst)
	if err != nil {
		return err
	}

	wasmExecSrc := filepath.Join(runtime.GOROOT(), "misc", "wasm", "wasm_exec.js")
	wasmExecDst := filepath.Join(appDir, "wasm_exec.js")
	err = util.CopyFile(wasmExecSrc, wasmExecDst)
	if err != nil {
		return err
	}

	exeDst := filepath.Join(appDir, p.name)
	err = util.CopyFile(p.exe, exeDst)
	if err != nil {
		return err
	}

	// Download webgl-debug.js directly from the KhronosGroup repository when needed
	if !p.release {
		webglDebugFile := filepath.Join(appDir, "webgl-debug.js")
		err := ioutil.WriteFile(webglDebugFile, templates.WebGLDebugJs, 0644)
		if err != nil {
			return err
		}
	}

	return nil
}
