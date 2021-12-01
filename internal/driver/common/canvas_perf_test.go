package common_test

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
)

const sizeTree = 1_000

func Example() {
	a := app.New()
	w := a.NewWindow("Perf tree")

	w.Resize(fyne.NewSize(600, 600))
	content := tree()
	w.SetContent(content)
	time.AfterFunc(200*time.Millisecond, w.Close)
	w.ShowAndRun()
	fmt.Println("Number of isBranch calls:", count/100*100) // rounding

	// with master version of Canvas.EnsureMinSize : 1028000
	// with patched version of Canvas.EnsureMinSize : 29000

	// Output: Number of isBranch calls: 1028000
}

var count int

func tree() fyne.CanvasObject {
	childUIDs := func(uid widget.TreeNodeID) (c []widget.TreeNodeID) {
		if uid == "" {
			out := make([]string, sizeTree)
			for i := range out {
				out[i] = fmt.Sprintf("node_%d", i)
			}
			return out
		} else if uid[0] == 'n' {
			return []string{"kid" + uid}
		} else {
			return nil
		}
	}
	// Return a CanvasObject that can represent a Branch (if branch is true), or a Leaf (if branch is false)
	createNode := func(branch bool) (o fyne.CanvasObject) {
		return widget.NewLabel("")
	}
	// Return true if the given widget.TreeNodeID represents a Branch
	isBranch := func(uid widget.TreeNodeID) (ok bool) {
		count++
		return uid == "" || uid[0] != 'k'
	}
	// Called to update the given CanvasObject to represent the data at the given widget.TreeNodeID
	updateNode := func(uid widget.TreeNodeID, _ bool, node fyne.CanvasObject) {
		node.(*widget.Label).SetText(uid)
	}

	return widget.NewTree(childUIDs, isBranch, createNode, updateNode)
}
