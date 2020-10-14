package widget

import (
	"io/ioutil"
	"os"
	"path"
	"testing"

	"fyne.io/fyne/storage"

	"github.com/stretchr/testify/assert"
)

func TestFileTree(t *testing.T) {
	t.Run("Initializer_Empty", func(t *testing.T) {
		tree := &FileTree{}
		var nodes []string
		tree.walkAll(func(uid string, branch bool, depth int) {
			nodes = append(nodes, uid)
		})
		assert.Equal(t, 0, len(nodes))
	})
	t.Run("NewFileTree", func(t *testing.T) {
		tempDir, err := ioutil.TempDir("", "test")
		assert.NoError(t, err)
		defer os.RemoveAll(tempDir)
		err = os.MkdirAll(path.Join(tempDir, "A"), os.ModePerm)
		assert.NoError(t, err)
		err = os.MkdirAll(path.Join(tempDir, "B"), os.ModePerm)
		assert.NoError(t, err)
		err = ioutil.WriteFile(path.Join(tempDir, "B", "C"), []byte("c"), os.ModePerm)
		assert.NoError(t, err)

		root := storage.NewURI("file://" + tempDir)
		tree := NewFileTree(root)
		tree.OpenAllBranches()
		var branches []string
		var leaves []string
		tree.walkAll(func(uid string, branch bool, depth int) {
			if branch {
				branches = append(branches, uid)
			} else {
				leaves = append(leaves, uid)
			}
		})
		assert.Equal(t, 3, len(branches))
		assert.Equal(t, root.String(), branches[0]) // Root
		b1, err := storage.Child(root, "A")
		assert.NoError(t, err)
		assert.Equal(t, b1.String(), branches[1])
		b2, err := storage.Child(root, "B")
		assert.NoError(t, err)
		assert.Equal(t, b2.String(), branches[2])
		assert.Equal(t, 1, len(leaves))
		l1, err := storage.Child(root, "B")
		assert.NoError(t, err)
		l1, err = storage.Child(l1, "C")
		assert.NoError(t, err)
		assert.Equal(t, l1.String(), leaves[0])
	})
}
