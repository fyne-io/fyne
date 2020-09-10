package widget

import (
	"image/color"
	//"log"
	"os"
	"path/filepath"
	"strings"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/driver/desktop"
	"fyne.io/fyne/internal/widget"
	"fyne.io/fyne/storage"
	"fyne.io/fyne/theme"
)

const (
	indentation = 4 // Multiplier of theme.Padding() to indent children by. // TODO consider moving to theme.Indentation()
)

var _ fyne.Widget = (*Tree)(nil)

// Tree widget displays hierarchical data.
// Each node of the tree must be identified by a Unique ID.
type Tree struct {
	BaseWidget
	Root     string
	Selected string
	Offset   fyne.Position

	open map[string]bool

	Children       func(uid string) (c []string)                         // Return a sorted slice of Children Unique IDs for the given Node Unique ID
	IsBranch       func(uid string) (ok bool)                            // Return true if the given Unique ID represents a Branch
	NewNode        func(branch bool) (o fyne.CanvasObject)               // Return a CanvasObject that can represent a Branch (if branch is true), or a Leaf (if branch is false)
	UpdateNode     func(uid string, branch bool, node fyne.CanvasObject) // Update the given CanvasObject to represent the data at the given Unique ID
	OnBranchOpened func(uid string)                                      // Called when a Branch is opened
	OnBranchClosed func(uid string)                                      // Called when a Branch is closed
	OnNodeSelected func(uid string, node fyne.CanvasObject)              // Called when the Node with the given Unique ID is selected.
}

// NewTreeWithStrings creates a new tree with the given string map.
// Data must contain a mapping for the root, which defaults to empty string ("").
func NewTreeWithStrings(data map[string][]string) (t *Tree) {
	t = &Tree{
		open: make(map[string]bool),
		Children: func(uid string) (c []string) {
			c, _ = data[uid]
			//log.Println("StringTree.Children:", uid, c)
			return
		},
		IsBranch: func(uid string) (b bool) {
			_, b = data[uid]
			//log.Println("StringTree.IsBranch:", uid, b)
			return
		},
		NewNode: func(branch bool) (o fyne.CanvasObject) {
			o = NewLabel("")
			//log.Println("StringTree.NewNode:", branch, o)
			return
		},
	}
	t.UpdateNode = func(uid string, branch bool, node fyne.CanvasObject) {
		//log.Println("StringTree.UpdateNode:", uid, branch, node)
		l := node.(*Label)
		l.SetText(uid)
		//log.Println("StringTree.Label:", l.Text)
	}
	t.ExtendBaseWidget(t)
	return
}

// NewTreeWithFiles creates a new tree with the given file system URI.
// TODO contents of a directory are sorted with the given sorter
func NewTreeWithFiles(root fyne.URI) (t *Tree) {
	t = &Tree{
		Root: root.String(),
		open: make(map[string]bool),
	}
	t.Children = func(uid string) (c []string) {
		luri, err := storage.ListerForURI(storage.NewURI(uid))
		if err != nil {
			fyne.LogError("Unable to get lister for "+uid, err)
		} else {
			uris, err := luri.List()
			if err != nil {
				fyne.LogError("Unable to list "+luri.String(), err)
			} else {
				// TODO sort.Slice(uris, sorter)
				for _, u := range uris {
					c = append(c, u.String())
				}
			}
		}
		//log.Println("FileTree.Children:", uid, c)
		return
	}
	t.IsBranch = func(uid string) (b bool) {
		if strings.HasPrefix(uid, "file://") {
			uid = uid[7:]
		}
		fi, err := os.Lstat(uid)
		if err != nil {
			fyne.LogError("Unable to stat path "+uid, err)
		} else {
			b = fi.IsDir()
		}
		//log.Println("FileTree.IsBranch:", uid, b)
		return
	}
	t.NewNode = func(branch bool) (o fyne.CanvasObject) {
		//log.Println("FileTree.NewNode:", branch, o)
		i := NewIcon(theme.FileIcon())
		l := NewLabel("Name")
		o = NewHBox(i, l)
		return
	}
	t.UpdateNode = func(uid string, branch bool, node fyne.CanvasObject) {
		//log.Println("FileTree.UpdateNode:", uid, branch, node)
		b := node.(*Box)
		i := b.Children[0].(*Icon)
		var r fyne.Resource
		if branch {
			if t.IsBranchOpen(uid) {
				// Set open folder icon
				r = theme.FolderOpenIcon()
			} else {
				// Set folder icon
				r = theme.FolderIcon()
			}
		} else {
			// Set file icon
			r = theme.FileIcon()

			/* TODO Copied from dialog/fileicon.go - perhaps this should be a utility?
			var res fyne.Resource
			switch strings.Split(mime.TypeByExtension(filepath.Extension(uid)), "/")[0] {
			case "application":
				res = theme.FileApplicationIcon()
			case "audio":
				res = theme.FileAudioIcon()
			case "image":
				res = theme.FileImageIcon()
			case "text":
				res = theme.FileTextIcon()
			case "video":
				res = theme.FileVideoIcon()
			default:
				res = theme.FileIcon()
			}
			*/
		}
		i.SetResource(r)
		l := b.Children[1].(*Label)
		if t.Root == uid {
			l.SetText(uid)
		} else {
			l.SetText(filepath.Base(uid))
		}
		//log.Println("FileTree.Label:", l.Text)
	}
	t.ExtendBaseWidget(t)
	return
}

// IsBranchOpen returns true if the branch with the given Unique ID is expanded.
func (t *Tree) IsBranchOpen(uid string) bool {
	if uid == t.Root {
		return true // Root is always open
	}
	t.propertyLock.RLock()
	defer t.propertyLock.RUnlock()
	return t.open[uid]
}

// OpenBranch opens the branch with the given Unique ID.
func (t *Tree) OpenBranch(uid string) {
	t.propertyLock.Lock()
	t.open[uid] = true
	t.propertyLock.Unlock()
	if f := t.OnBranchOpened; f != nil {
		f(uid)
	}
	t.Refresh()
}

// CloseBranch closes the branch with the given Unique ID.
func (t *Tree) CloseBranch(uid string) {
	t.propertyLock.Lock()
	t.open[uid] = false
	t.propertyLock.Unlock()
	if f := t.OnBranchClosed; f != nil {
		f(uid)
	}
	t.Refresh()
}

// ToggleBranch flips the state of the branch with the given Unique ID.
func (t *Tree) ToggleBranch(uid string) {
	if t.IsBranchOpen(uid) {
		t.CloseBranch(uid)
	} else {
		t.OpenBranch(uid)
	}
}

// OpenAllBranches opens all branches in the tree.
func (t *Tree) OpenAllBranches() {
	t.Walk(func(uid string, branch bool, depth int) {
		if branch {
			// TODO this triggers a refresh for each branch
			t.OpenBranch(uid)
		}
	})
	t.Refresh()
}

// CloseAllBranches closes all branches in the tree.
func (t *Tree) CloseAllBranches() {
	t.propertyLock.Lock()
	t.open = make(map[string]bool)
	t.propertyLock.Unlock()
	t.Refresh()
}

// SetSelectedNode updates the current selection to the node with the given Unique ID.
func (t *Tree) SetSelectedNode(uid string) {
	t.Selected = uid
	t.Refresh()
}

// Walk visits every open node of the tree and calls the given callback with node Unique ID, whether node is branch, and the depth of node.
func (t *Tree) Walk(onNode func(string, bool, int)) {
	t.walk(t.Root, 0, onNode)
}

func (t *Tree) walk(uid string, depth int, onNode func(string, bool, int)) {
	//log.Println("Tree.walk:", uid, depth)
	if isBranch := t.IsBranch; isBranch != nil {
		if isBranch(uid) {
			onNode(uid, true, depth)
			if t.IsBranchOpen(uid) {
				if children := t.Children; children != nil {
					for _, c := range children(uid) {
						t.walk(c, depth+1, onNode)
					}
				}
			}
		} else {
			onNode(uid, false, depth)
		}
	}
}

// CreateRenderer is a private method to Fyne which links this widget to its renderer.
func (t *Tree) CreateRenderer() fyne.WidgetRenderer {
	t.ExtendBaseWidget(t)
	c := newTreeContent(t)
	s := NewScrollContainer(c)
	r := &treeRenderer{
		BaseRenderer: widget.NewBaseRenderer([]fyne.CanvasObject{s}),
		tree:         t,
		content:      c,
		scroller:     s,
	}
	s.onOffsetChanged = func() {
		if t.Offset == s.Offset {
			return
		}
		t.Offset = s.Offset
		t.Refresh()
	}
	return r
}

var _ fyne.WidgetRenderer = (*treeRenderer)(nil)

type treeRenderer struct {
	widget.BaseRenderer
	tree     *Tree
	content  *treeContent
	scroller *ScrollContainer
}

func (r *treeRenderer) MinSize() (min fyne.Size) {
	min = r.scroller.MinSize()
	//log.Println("treeRenderer.MinSize:", min)
	return
}

func (r *treeRenderer) Layout(size fyne.Size) {
	//log.Println("treeRenderer.Layout:", size)
	r.scroller.Resize(size)
}

func (r *treeRenderer) Refresh() {
	//log.Println("treeRenderer.Refresh")
	s := r.tree.Size()
	if s.IsZero() {
		r.tree.Resize(r.tree.MinSize())
	}
	r.content.Refresh()
	canvas.Refresh(r.tree.super())
}

var _ fyne.Widget = (*treeContent)(nil)

type treeContent struct {
	BaseWidget
	tree *Tree
}

func newTreeContent(tree *Tree) (c *treeContent) {
	c = &treeContent{
		tree: tree,
	}
	c.ExtendBaseWidget(c)
	return
}

func (c *treeContent) CreateRenderer() fyne.WidgetRenderer {
	return &treeContentRenderer{
		BaseRenderer: widget.BaseRenderer{},
		treeContent:  c,
		minSizes:     make(map[string]fyne.Size),
		branches:     make(map[string]*branch),
		leaves:       make(map[string]*leaf),
		branchPool:   &syncPool{},
		leafPool:     &syncPool{},
	}
}

var _ fyne.WidgetRenderer = (*treeContentRenderer)(nil)

type treeContentRenderer struct {
	widget.BaseRenderer
	treeContent *treeContent
	minSizes    map[string]fyne.Size // Holds cache of children min sizes updated by MinSize and used by Layout
	branches    map[string]*branch
	leaves      map[string]*leaf
	branchPool  pool
	leafPool    pool
}

func (r *treeContentRenderer) MinSize() (min fyne.Size) {
	r.treeContent.propertyLock.Lock()
	defer r.treeContent.propertyLock.Unlock()

	r.treeContent.tree.Walk(func(uid string, isBranch bool, depth int) {
		// Root node is not rendered unless it has been customized
		if r.treeContent.tree.Root == "" {
			depth = depth - 1
			if uid == "" {
				// This is root node, skip
				return
			}
		}
		// TODO FIXME this assumes a node with the given Unique ID will always have the same MinSize.
		m := r.minSizeOf(uid, isBranch, depth)
		min.Width = fyne.Max(min.Width, m.Width)
		min.Height = min.Height + m.Height
	})
	min.Width += 2 * theme.Padding()
	min.Height += 2 * theme.Padding()
	//log.Println("treeContentRenderer.MinSize:", min)
	return
}

func (r *treeContentRenderer) Layout(size fyne.Size) {
	//log.Println("treeContentRenderer.Layout:", size)

	r.treeContent.propertyLock.Lock()
	defer r.treeContent.propertyLock.Unlock()

	branches := make(map[string]*branch)
	leaves := make(map[string]*leaf)

	offsetY := r.treeContent.tree.Offset.Y
	y := theme.Padding()
	// Walk open branches and obtain nodes to render in scroller's viewport
	r.treeContent.tree.Walk(func(uid string, isBranch bool, depth int) {
		// Root node is not rendered unless it has been customized
		if r.treeContent.tree.Root == "" {
			depth = depth - 1
			if uid == "" {
				// This is root node, skip
				return
			}
		}
		m := r.minSizeOf(uid, isBranch, depth)
		if y+m.Height < offsetY {
			// node is above viewport and not visible
		} else if y > offsetY+size.Height {
			// node is below viewport and not visible
		} else {
			// node is in viewport
			var n *treeNode
			if isBranch {
				b, ok := r.branches[uid]
				if !ok {
					b = r.getBranch()
					b.Update(uid, depth)
					r.treeContent.tree.UpdateNode(uid, true, b.Content())
				}
				branches[uid] = b
				n = b.treeNode
			} else {
				l, ok := r.leaves[uid]
				if !ok {
					l = r.getLeaf()
					l.Update(uid, depth)
					r.treeContent.tree.UpdateNode(uid, false, l.Content())
				}
				leaves[uid] = l
				n = l.treeNode
			}
			if n != nil {
				n.Move(fyne.NewPos(theme.Padding(), y))
				//log.Println("Node", uid, "Pos", n.Position())
				n.Resize(fyne.NewSize(size.Width-n.Position().X-theme.Padding(), m.Height))
				//log.Println("Node", uid, "Size", n.Size())
			}
		}
		y += m.Height
	})

	// Release any treeNodes that haven't been reused
	for uid, b := range r.branches {
		if _, ok := branches[uid]; !ok {
			r.branchPool.Release(b)
		}
	}
	for uid, l := range r.leaves {
		if _, ok := leaves[uid]; !ok {
			r.leafPool.Release(l)
		}
	}

	r.branches = branches
	r.leaves = leaves
}

func (r *treeContentRenderer) Refresh() {
	//log.Println("treeContentRenderer.Refresh")
	s := r.treeContent.Size()
	if s.IsZero() {
		r.treeContent.Resize(r.treeContent.MinSize())
	} else {
		r.Layout(s)
	}
	r.treeContent.propertyLock.RLock()
	for _, b := range r.branches {
		b.Refresh()
	}
	for _, l := range r.leaves {
		l.Refresh()
	}
	r.treeContent.propertyLock.RUnlock()
	canvas.Refresh(r.treeContent.super())
}

func (r *treeContentRenderer) Objects() (objects []fyne.CanvasObject) {
	r.treeContent.propertyLock.RLock()
	for _, b := range r.branches {
		objects = append(objects, b)
	}
	for _, l := range r.leaves {
		objects = append(objects, l)
	}
	r.treeContent.propertyLock.RUnlock()
	//log.Println("treeContentRenderer.Objects:", objects)
	return
}

func (r *treeContentRenderer) getBranch() (b *branch) {
	o := r.branchPool.Obtain()
	if o != nil {
		b = o.(*branch)
	} else {
		b = newBranch(r.treeContent.tree, r.treeContent.tree.NewNode(true))
	}
	return
}

func (r *treeContentRenderer) getLeaf() (l *leaf) {
	o := r.leafPool.Obtain()
	if o != nil {
		l = o.(*leaf)
	} else {
		l = newLeaf(r.treeContent.tree, r.treeContent.tree.NewNode(false))
	}
	return
}

func (r *treeContentRenderer) minSizeOf(uid string, isBranch bool, depth int) fyne.Size {
	m, ok := r.minSizes[uid]
	if !ok {
		var n fyne.CanvasObject
		if isBranch {
			b := r.getBranch()
			b.Update(uid, depth)
			r.treeContent.tree.UpdateNode(uid, true, b.Content())
			n = b
		} else {
			l := r.getLeaf()
			l.Update(uid, depth)
			r.treeContent.tree.UpdateNode(uid, false, l.Content())
			n = l
		}
		m = n.MinSize()
		if isBranch {
			r.branchPool.Release(n)
		} else {
			r.leafPool.Release(n)
		}
		r.minSizes[uid] = m
	}
	return m
}

var _ fyne.CanvasObject = (*treeNode)(nil)
var _ fyne.Tappable = (*treeNode)(nil)
var _ desktop.Hoverable = (*treeNode)(nil)

type treeNode struct {
	BaseWidget
	tree    *Tree
	uid     string
	depth   int
	hovered bool
	icon    fyne.CanvasObject
	content fyne.CanvasObject
}

func (n *treeNode) Update(uid string, depth int) {
	n.uid = uid
	n.depth = depth
}

func (n *treeNode) Content() fyne.CanvasObject {
	return n.content
}

func (n *treeNode) Indent() int {
	return n.depth * indentation * theme.Padding()
}

func (n *treeNode) Tapped(*fyne.PointEvent) {
	//log.Println("treeNode.Tapped")
	n.tree.SetSelectedNode(n.uid)
	if f := n.tree.OnNodeSelected; f != nil {
		f(n.uid, n.content)
	}
}

// MouseIn is called when a desktop pointer enters the widget
func (n *treeNode) MouseIn(*desktop.MouseEvent) {
	n.hovered = true
	n.Refresh()
}

// MouseOut is called when a desktop pointer exits the widget
func (n *treeNode) MouseOut() {
	n.hovered = false
	n.Refresh()
}

// MouseMoved is called when a desktop pointer hovers over the widget
func (n *treeNode) MouseMoved(*desktop.MouseEvent) {
}

func (n *treeNode) CreateRenderer() fyne.WidgetRenderer {
	return &treeNodeRenderer{
		BaseRenderer: widget.BaseRenderer{},
		treeNode:     n,
		divider:      canvas.NewRectangle(theme.ShadowColor()),
		indicator:    canvas.NewRectangle(theme.BackgroundColor()),
	}
}

var _ fyne.WidgetRenderer = (*treeNodeRenderer)(nil)

type treeNodeRenderer struct {
	widget.BaseRenderer
	treeNode           *treeNode
	divider, indicator *canvas.Rectangle
}

func (r *treeNodeRenderer) BackgroundColor() color.Color {
	/* TODO FIXME removed until transparent BG becomes the default
	if r.treeNode.hovered {
		return theme.HoverColor()
	}
	*/
	return theme.BackgroundColor()
}

func (r *treeNodeRenderer) MinSize() (min fyne.Size) {
	if c := r.treeNode.content; c != nil {
		min = c.MinSize()
	}
	min.Width += r.treeNode.Indent() + theme.IconInlineSize() + theme.Padding()
	min.Height = fyne.Max(min.Height, theme.IconInlineSize())
	min.Width += 2 * theme.Padding()
	min.Height += 2 * theme.Padding()
	//log.Println("treeNodeRenderer.MinSize:", min)
	return
}

func (r *treeNodeRenderer) Layout(size fyne.Size) {
	//log.Println("treeNodeRenderer.Layout:", size)
	if d := r.divider; d != nil {
		d.Move(fyne.NewPos(0, size.Height-1))
		s := fyne.NewSize(size.Width, 1)
		d.SetMinSize(s)
		d.Resize(s)
	}
	x := theme.Padding() + r.treeNode.Indent()
	y := theme.Padding()
	height := size.Height - 2*theme.Padding()
	if i := r.treeNode.icon; i != nil {
		i.Move(fyne.NewPos(x, y))
		i.Resize(fyne.NewSize(theme.IconInlineSize(), height))
	}
	x += theme.IconInlineSize()
	if i := r.indicator; i != nil {
		i.Move(fyne.NewPos(x, 0))
		s := fyne.NewSize(theme.Padding(), size.Height-1)
		i.SetMinSize(s)
		i.Resize(s)
	}
	x += theme.Padding()
	if c := r.treeNode.content; c != nil {
		c.Move(fyne.NewPos(x, y))
		c.Resize(fyne.NewSize(size.Width-x-theme.Padding(), height))
	}
}

func (r *treeNodeRenderer) Refresh() {
	//log.Println(r.treeNode.uid, "treeNodeRenderer.Refresh")
	if i := r.treeNode.icon; i != nil {
		i.Refresh()
	}
	if i := r.divider; i != nil {
		i.FillColor = theme.ShadowColor()
		i.Refresh()
	}
	if i := r.indicator; i != nil {
		//log.Println(r.treeNode.uid, "treeNodeRenderer.Refresh:", r.treeNode.tree.Selected)
		if r.treeNode.uid == r.treeNode.tree.Selected {
			i.FillColor = theme.PrimaryColor()
		} else if r.treeNode.hovered {
			i.FillColor = theme.HoverColor()
		} else {
			i.FillColor = theme.BackgroundColor()
		}
		i.Refresh()
	}
	if c := r.treeNode.content; c != nil {
		c.Refresh()
	}
	canvas.Refresh(r.treeNode.super())
}

func (r *treeNodeRenderer) Objects() (objects []fyne.CanvasObject) {
	if i := r.treeNode.icon; i != nil {
		objects = append(objects, i)
	}
	if d := r.divider; d != nil {
		objects = append(objects, d)
	}
	if i := r.indicator; i != nil {
		objects = append(objects, i)
	}
	if c := r.treeNode.content; c != nil {
		objects = append(objects, c)
	}
	//log.Println("branchRenderer.Objects:", objects)
	return
}

var _ fyne.DoubleTappable = (*branch)(nil)
var _ fyne.Widget = (*branch)(nil)

type branch struct {
	*treeNode
}

func newBranch(tree *Tree, content fyne.CanvasObject) (b *branch) {
	b = &branch{
		treeNode: &treeNode{
			tree:    tree,
			icon:    newBranchIcon(tree),
			content: content,
		},
	}
	b.ExtendBaseWidget(b)
	return
}

func (b *branch) Update(uid string, depth int) {
	b.treeNode.Update(uid, depth)
	b.icon.(*branchIcon).Update(uid, depth)
}

func (b *branch) DoubleTapped(*fyne.PointEvent) {
	//log.Println("branch.DoubleTapped")
	b.tree.ToggleBranch(b.uid)
}

var _ fyne.Tappable = (*branchIcon)(nil)
var _ fyne.DoubleTappable = (*branchIcon)(nil)

type branchIcon struct {
	Icon
	tree *Tree
	uid  string
}

func newBranchIcon(tree *Tree) (i *branchIcon) {
	i = &branchIcon{
		tree: tree,
	}
	i.ExtendBaseWidget(i)
	return
}

func (i *branchIcon) Update(uid string, depth int) {
	i.uid = uid
	i.Refresh()
}

func (i *branchIcon) Tapped(*fyne.PointEvent) {
	//log.Println("branchIcon.Tapped")
	i.tree.ToggleBranch(i.uid)
}

func (i *branchIcon) DoubleTapped(*fyne.PointEvent) {
	//log.Println("branchIcon.DoubleTapped")
	// Do nothing - this stops the event propagating to branch
}

func (i *branchIcon) Refresh() {
	//log.Println(i.uid, "branchIcon.Refresh")
	if i.tree.IsBranchOpen(i.uid) {
		i.Resource = theme.MoveDownIcon()
	} else {
		i.Resource = theme.NavigateNextIcon()
	}
	i.Icon.Refresh()
}

var _ fyne.Widget = (*leaf)(nil)

type leaf struct {
	*treeNode
}

func newLeaf(tree *Tree, content fyne.CanvasObject) (l *leaf) {
	l = &leaf{
		&treeNode{
			tree:    tree,
			content: content,
		},
	}
	l.ExtendBaseWidget(l)
	return
}
