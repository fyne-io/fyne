package commands

import (
	"os"
	"path/filepath"

	"fyne.io/fyne/v2/cmd/fyne/internal/templates"
	"fyne.io/fyne/v2/cmd/fyne/internal/util"
)

func (p *Packager) packageGopherJS() error {
	appDir := util.EnsureSubDir(p.dir, "gopherjs")

	index := filepath.Join(appDir, "index.html")
	indexFile, err := os.Create(index)
	if err != nil {
		return err
	}

	tplData := indexData{GopherJSFile: p.name, IsReleased: p.release, HasGopherJS: true}
	err = templates.IndexHTML.Execute(indexFile, tplData)
	if err != nil {
		return err
	}

	iconDst := filepath.Join(appDir, "icon.png")
	err = util.CopyFile(p.icon, iconDst)
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
		err := util.WriteFile(webglDebugFile, templates.WebGLDebugJs)
		if err != nil {
			return err
		}
	}

	return nil
}
