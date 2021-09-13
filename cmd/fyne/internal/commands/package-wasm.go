package commands

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"

	"fyne.io/fyne/v2/cmd/fyne/internal/util"
)

func (p *Packager) packageWasm() error {
	appDir := util.EnsureSubDir(p.dir, "wasm")

	index := filepath.Join(appDir, "index.html")
	indexFile, err := os.Create(index)
	if err != nil {
		return err
	}

	io.WriteString(indexFile, "<!DOCTYPE html>\n"+
		"<html>\n"+
		"\t<head>\n"+
		"\t\t<meta charset=\"utf-8\">\n"+
		"\t\t<link rel=\"icon\" type=\"image/png\" href=\"icon.png\">\n"+
		"\t\t<script src=\"wasm_exec.js\"></script>\n")
	if !p.release {
		io.WriteString(indexFile, "\t\t<script src=\"webgl-debug.js\"></script>\n")
	}
	io.WriteString(indexFile, "\t\t<script>\n"+
		"\t\t\tfunction webgl_support () {\n"+
		"\t\t\t\ttry {\n"+
		"\t\t\t\t\tvar canvas = document.createElement('canvas');\n"+
		"\t\t\t\t\treturn !!window.WebGLRenderingContext &&\n"+
		"\t\t\t\t\t\t(canvas.getContext('webgl') || canvas.getContext('experimental-webgl'));\n"+
		"\t\t\t\t} catch(e) {\n"+
		"\t\t\t\t\treturn false;\n"+
		"\t\t\t\t}\n"+
		"\t\t\t}\n"+
		"\n"+
		"\t\t\tif (webgl_support() && WebAssembly) {\n"+
		"\t\t\t\t// WebAssembly.instantiateStreaming is not currently available in Safari\n"+
		"\t\t\t\tif (WebAssembly && !WebAssembly.instantiateStreaming) { // polyfill\n"+
		"\t\t\t\t\tWebAssembly.instantiateStreaming = async (resp, importObject) => {\n"+
		"\t\t\t\t\t\tconst source = await (await resp).arrayBuffer();\n"+
		"\t\t\t\t\t\treturn await WebAssembly.instantiate(source, importObject);\n"+
		"\t\t\t\t\t};\n"+
		"\t\t\t\t}\n"+
		"\n"+
		"\t\t\t\tconst go = new Go();\n"+
		"\t\t\t\tWebAssembly.instantiateStreaming(fetch(\""+p.name+"\"), go.importObject).then((result) => {\n"+
		"\t\t\t\t\tgo.run(result.instance);\n"+
		"\t\t\t\t});\n"+
		"\t\t\t} else {\n"+
		"\t\t\t\tconsole.log(\"WebAssembly and WebGL are not supported in your browser\")\n"+
		"\t\t\t\tvar main = document.querySelector('#main');\n"+
		"\t\t\t\tmain.innerHTML = \"WebAssembly and WebGL are not supported in your browser\";"+
		"\t\t\t}\n"+
		"\n"+
		"\t\t</script>\n"+
		"\t</head>\n"+
		"\t<body>\n"+
		"\t\t<main id=\"wasm\"></main>\n"+
		"\t</body>\n"+
		"</html>\n")

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
