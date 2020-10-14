package widget

import (
	"sort"

	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/storage"
	"fyne.io/fyne/theme"
)

// FileTree widget displays a hierarchical file system.
// Each node of the tree must be identified by a Unique ID.
type FileTree struct {
	Tree
	Filter func(string) bool
	Sorter func(string, string) bool
}

// NewFileTree creates a new tree with the given file system URI.
func NewFileTree(root fyne.URI) *FileTree {
	t := &FileTree{
		Tree: Tree{
			Root: root.String(),
			IsBranch: func(uid string) bool {
				_, err := storage.ListerForURI(storage.NewURI(uid))
				return err == nil
			},
			CreateNode: func(branch bool) fyne.CanvasObject {
				var icon fyne.CanvasObject
				if branch {
					icon = NewIcon(nil)
				} else {
					icon = NewFileIcon(nil)
				}
				return fyne.NewContainerWithLayout(layout.NewHBoxLayout(), icon, NewLabel("Template Object"))
			},
		},
	}
	t.ChildUIDs = func(uid string) (c []string) {
		luri, err := storage.ListerForURI(storage.NewURI(uid))
		if err != nil {
			fyne.LogError("Unable to get lister for "+uid, err)
		} else {
			uris, err := luri.List()
			if err != nil {
				fyne.LogError("Unable to list "+luri.String(), err)
			} else {
				for _, u := range uris {
					s := u.String()
					if f := t.Filter; f == nil || f(s) {
						c = append(c, s)
					}
				}
			}
		}
		if s := t.Sorter; s != nil {
			sort.Slice(c, func(i, j int) bool {
				return s(c[i], c[j])
			})
		}
		return
	}
	t.UpdateNode = func(uid string, branch bool, node fyne.CanvasObject) {
		uri := storage.NewURI(uid)
		c := node.(*fyne.Container)
		if branch {
			var r fyne.Resource
			if t.IsBranchOpen(uid) {
				// Set open folder icon
				r = theme.FolderOpenIcon()
			} else {
				// Set folder icon
				r = theme.FolderIcon()
			}
			c.Objects[0].(*Icon).SetResource(r)
		} else {
			// Set file uri to update icon
			c.Objects[0].(*FileIcon).SetURI(uri)
		}
		l := c.Objects[1].(*Label)
		if t.Root == uid {
			l.SetText(uid)
		} else {
			l.SetText(uri.Name())
		}
	}
	t.ExtendBaseWidget(t)
	return t
}
