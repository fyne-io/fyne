//go:build wasm

package repository

import (
	"io"
	"syscall/js"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/storage/repository"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// testBrowserFile makes a JavaScript File like those the browser's picker gives.
func testBrowserFile(name, relativePath, content string) js.Value {
	parts := js.Global().Get("Array").New()
	parts.Call("push", content)
	file := js.Global().Get("File").New(parts, name, map[string]any{"type": "text/plain"})
	if relativePath != "" {
		js.Global().Get("Object").Call("defineProperty", file, "webkitRelativePath",
			map[string]any{"value": relativePath})
	}
	return file
}

func names(uris []fyne.URI) []string {
	list := make([]string, len(uris))
	for i, u := range uris {
		list[i] = u.Name()
	}
	return list
}

func TestBrowserFileRepository_File(t *testing.T) {
	repo := BrowserFiles()
	repository.Register("file", repo)

	u := repo.AddFile(testBrowserFile("notes #1.txt", "", "first"))
	assert.Equal(t, "file", u.Scheme())
	assert.Equal(t, "/notes #1.txt", u.Path())

	exists, err := storage.Exists(u)
	require.NoError(t, err)
	assert.True(t, exists)
	listable, err := storage.CanList(u)
	require.NoError(t, err)
	assert.False(t, listable)

	r, err := storage.Reader(u)
	require.NoError(t, err)
	data, err := io.ReadAll(r)
	require.NoError(t, err)
	assert.Equal(t, "first", string(data))
	assert.Equal(t, "text/plain", r.URI().MimeType())

	// choosing a file of the same name again replaces it
	repo.AddFile(testBrowserFile("notes #1.txt", "", "second"))
	r, err = storage.Reader(u)
	require.NoError(t, err)
	data, _ = io.ReadAll(r)
	assert.Equal(t, "second", string(data))

	missing := storage.NewFileURI("/not chosen.txt")
	exists, err = storage.Exists(missing)
	require.NoError(t, err)
	assert.False(t, exists)
	_, err = storage.Reader(missing)
	assert.Error(t, err)
}

func TestBrowserFileRepository_Folder(t *testing.T) {
	repo := BrowserFiles()
	repository.Register("file", repo)

	files := js.Global().Get("Array").New()
	files.Call("push", testBrowserFile("b.txt", "Album/b.txt", "bee"))
	files.Call("push", testBrowserFile("a.txt", "Album/a.txt", "ay"))
	files.Call("push", testBrowserFile("c.txt", "Album/Disc 2/c.txt", "sea"))
	u := repo.AddFolder(files)
	assert.Equal(t, "/Album", u.Path())

	folder, err := storage.ListerForURI(u)
	require.NoError(t, err)
	items, err := folder.List()
	require.NoError(t, err)
	assert.Equal(t, []string{"Disc 2", "a.txt", "b.txt"}, names(items))

	sub, err := storage.ListerForURI(items[0])
	require.NoError(t, err)
	inner, err := sub.List()
	require.NoError(t, err)
	require.Equal(t, []string{"c.txt"}, names(inner))
	r, err := storage.Reader(inner[0])
	require.NoError(t, err)
	data, _ := io.ReadAll(r)
	assert.Equal(t, "sea", string(data))

	parent, err := storage.Parent(inner[0])
	require.NoError(t, err)
	listable, err := storage.CanList(parent)
	require.NoError(t, err)
	assert.True(t, listable)

	// choosing a folder of the same name again replaces all of it
	files = js.Global().Get("Array").New()
	files.Call("push", testBrowserFile("d.txt", "Album/d.txt", "dee"))
	repo.AddFolder(files)
	items, err = folder.List()
	require.NoError(t, err)
	assert.Equal(t, []string{"d.txt"}, names(items))
	exists, _ := storage.Exists(inner[0])
	assert.False(t, exists)

	// files that come without a folder, where a browser cannot choose one, are gathered into one
	files = js.Global().Get("Array").New()
	files.Call("push", testBrowserFile("loose.txt", "", "loose"))
	u = repo.AddFolder(files)
	assert.Equal(t, "/files", u.Path())
}

func TestBrowserFileRepository_Download(t *testing.T) {
	repo := BrowserFiles()
	repository.Register("file", repo)

	// note what is downloaded in place of having the browser do it
	js.Global().Call("eval", `
		window.fyneTestDownloads = [];
		window.fyneTestClick = HTMLAnchorElement.prototype.click;
		HTMLAnchorElement.prototype.click = function() {
			window.fyneTestDownloads.push({name: this.download, href: this.href});
		};`)
	t.Cleanup(func() {
		js.Global().Call("eval", `HTMLAnchorElement.prototype.click = window.fyneTestClick;`)
	})
	downloads := js.Global().Get("fyneTestDownloads")
	content := func(i int) string {
		got := make(chan string, 1)
		then := js.FuncOf(func(_ js.Value, args []js.Value) any {
			got <- args[0].String()
			return nil
		})
		defer then.Release()
		text := js.FuncOf(func(_ js.Value, args []js.Value) any {
			return args[0].Call("text")
		})
		defer text.Release()
		js.Global().Call("fetch", downloads.Index(i).Get("href")).Call("then", text).Call("then", then)
		return <-got
	}

	w := repo.Download("report.txt")
	assert.Equal(t, "/report.txt", w.URI().Path())
	_, err := w.Write([]byte("hello "))
	require.NoError(t, err)
	_, err = w.Write([]byte("world"))
	require.NoError(t, err)
	assert.Equal(t, 0, downloads.Length(), "nothing is downloaded until the file is closed")
	require.NoError(t, w.Close())
	require.NoError(t, w.Close())
	require.Equal(t, 1, downloads.Length())
	assert.Equal(t, "report.txt", downloads.Index(0).Get("name").String())
	assert.Equal(t, "hello world", content(0))
	_, err = w.Write([]byte("more"))
	assert.Error(t, err)

	// saving to the URI again, as an app does to save a document it has open, downloads it again
	again, err := storage.Writer(w.URI())
	require.NoError(t, err)
	_, err = again.Write([]byte("changed"))
	require.NoError(t, err)
	require.NoError(t, again.Close())
	require.Equal(t, 2, downloads.Length())
	assert.Equal(t, "report.txt", downloads.Index(1).Get("name").String())
	assert.Equal(t, "changed", content(1))
}
