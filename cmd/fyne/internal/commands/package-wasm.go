package commands

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"

	"fyne.io/fyne/v2/cmd/fyne/internal/templates"
	"fyne.io/fyne/v2/cmd/fyne/internal/util"
)

type indexData struct {
	Name       string
	IsReleased bool
}

func (p *Packager) packageWasm() error {
	appDir := util.EnsureSubDir(p.dir, "wasm")

	index := filepath.Join(appDir, "index.html")
	indexFile, err := os.Create(index)
	if err != nil {
		return err
	}

	tplData := indexData{Name: p.name, IsReleased: p.release}
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
		r, err := http.Get("https://raw.githubusercontent.com/KhronosGroup/WebGLDeveloperTools/b42e702487d02d5278814e0fe2e2888d234893e6/src/debug/webgl-debug.js")
		if err != nil {
			return err
		}
		defer r.Body.Close()

		webglDebugFile := filepath.Join(appDir, "webgl-debug.js")
		out, err := os.Create(webglDebugFile)
		if err != nil {
			return err
		}
		defer out.Close()

		_, err = io.Copy(out, r.Body)
		return err
	}

	return nil
}
