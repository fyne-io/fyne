package widget

import (
	"image/color"
	//"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

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

var _ fyne.Widget = (*TreeContainer)(nil)

// TreeContainer widget displays hierarchical data.
type TreeContainer struct {
	BaseWidget
	Root     string
	Selected string
	Offset   fyne.Position

	open map[string]bool

	Children       func(id string) (c []string)                         // Return a sorted slice of Children IDs for the given Node ID
	IsBranch       func(id string) (ok bool)                            // Return true if the given ID represents a Branch
	NewNode        func(branch bool) (o fyne.CanvasObject)              // Return a CanvasObject that can represent a Branch (if branch is true), or a Leaf (if branch is false)
	UpdateNode     func(id string, branch bool, node fyne.CanvasObject) // Update the given CanvasObject to represent the data at the given ID
	OnBranchOpened func(id string)                                      // Called when a Branch is opened
	OnBranchClosed func(id string)                                      // Called when a Branch is closed
	OnNodeSelected func(id string, node fyne.CanvasObject)              // Called when the Node with the given ID is selected.
}

// NewTreeWithStrings creates a new tree with the given string map.
func NewTreeWithStrings(ss map[string][]string) (t *TreeContainer) {
	t = &TreeContainer{
		open: make(map[string]bool),
		Children: func(id string) (c []string) {
			c, _ = ss[id]
			//log.Println("StringTree.Children:", id, c)
			return
		},
		IsBranch: func(id string) (b bool) {
			_, b = ss[id]
			//log.Println("StringTree.IsBranch:", id, b)
			return
		},
		NewNode: func(branch bool) (o fyne.CanvasObject) {
			o = NewLabel("")
			//log.Println("StringTree.NewNode:", branch, o)
			return
		},
	}
	t.UpdateNode = func(id string, branch bool, node fyne.CanvasObject) {
		//log.Println("StringTree.UpdateNode:", id, branch, node)
		l := node.(*Label)
		l.SetText(id)
		//log.Println("StringTree.Label:", l.Text)
	}
	t.ExtendBaseWidget(t)
	return
}

// NewTreeWithFiles creates a new tree with the given file system URI.
// TODO contents of a directory are sorted with the given sorter
func NewTreeWithFiles(root fyne.URI) (t *TreeContainer) {
	t = &TreeContainer{
		Root: root.String(),
		open: make(map[string]bool),
	}
	t.Children = func(id string) (c []string) {
		luri, err := storage.ListerForURI(storage.NewURI(id))
		if err != nil {
			fyne.LogError("Unable to get lister for "+id, err)
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
		//log.Println("FileTree.Children:", id, c)
		return
	}
	t.IsBranch = func(id string) (b bool) {
		if strings.HasPrefix(id, "file://") {
			id = id[7:]
		}
		fi, err := os.Lstat(id)
		if err != nil {
			fyne.LogError("Unable to stat path "+id, err)
		} else {
			b = fi.IsDir()
		}
		//log.Println("FileTree.IsBranch:", id, b)
		return
	}
	t.NewNode = func(branch bool) (o fyne.CanvasObject) {
		//log.Println("FileTree.NewNode:", branch, o)
		i := NewIcon(theme.FileIcon())
		l := NewLabel("Name")
		o = NewHBox(i, l)
		return
	}
	t.UpdateNode = func(id string, branch bool, node fyne.CanvasObject) {
		//log.Println("FileTree.UpdateNode:", id, branch, node)
		b := node.(*Box)
		i := b.Children[0].(*Icon)
		var r fyne.Resource
		if branch {
			if t.IsBranchOpen(id) {
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
			switch strings.Split(mime.TypeByExtension(filepath.Extension(id)), "/")[0] {
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
		if t.Root == id {
			l.SetText(id)
		} else {
			l.SetText(filepath.Base(id))
		}
		//log.Println("FileTree.Label:", l.Text)
	}
	t.ExtendBaseWidget(t)
	return
}

// AddTreePath adds the given path to the given parent->children map
func AddTreePath(data map[string][]string, path ...string) {
	parent := ""
	for _, p := range path {
		children := data[parent]
		add := true
		for _, c := range children {
			if c == p {
				add = false
				break
			}
		}
		if add {
			data[parent] = append(children, p)
		}
		parent = p
	}
}

// IsBranchOpen returns true if the branch with the given ID is expanded.
func (t *TreeContainer) IsBranchOpen(id string) bool {
	if id == t.Root {
		return true // Root is always open
	}
	t.propertyLock.RLock()
	defer t.propertyLock.RUnlock()
	return t.open[id]
}

// OpenBranch opens the branch with the given ID.
func (t *TreeContainer) OpenBranch(id string) {
	t.propertyLock.Lock()
	t.open[id] = true
	t.propertyLock.Unlock()
	if f := t.OnBranchOpened; f != nil {
		f(id)
	}
	t.Refresh()
}

// CloseBranch closes the branch with the given ID.
func (t *TreeContainer) CloseBranch(id string) {
	t.propertyLock.Lock()
	t.open[id] = false
	t.propertyLock.Unlock()
	if f := t.OnBranchClosed; f != nil {
		f(id)
	}
	t.Refresh()
}

// ToggleBranch flips the state of the branch with the given ID.
func (t *TreeContainer) ToggleBranch(id string) {
	if t.IsBranchOpen(id) {
		t.CloseBranch(id)
	} else {
		t.OpenBranch(id)
	}
}

// OpenAllBranches opens all branches in the tree.
func (t *TreeContainer) OpenAllBranches() {
	t.Walk(func(id string, branch bool, depth int) {
		if branch {
			// TODO this triggers a refresh for each branch
			t.OpenBranch(id)
		}
	})
	t.Refresh()
}

// CloseAllBranches closes all branches in the tree.
func (t *TreeContainer) CloseAllBranches() {
	t.propertyLock.Lock()
	t.open = make(map[string]bool)
	t.propertyLock.Unlock()
	t.Refresh()
}

// SetSelectedNode updates the current selection to the node with the given ID.
func (t *TreeContainer) SetSelectedNode(id string) {
	t.Selected = id
	t.Refresh()
}

// Walk visits every open node of the tree and calls the given callback with node ID, whether node is branch, and the depth of node.
func (t *TreeContainer) Walk(onNode func(string, bool, int)) {
	t.walk(t.Root, 0, onNode)
}

func (t *TreeContainer) walk(id string, depth int, onNode func(string, bool, int)) {
	//log.Println("TreeContainer.walk:", id, depth)
	if isBranch := t.IsBranch; isBranch != nil {
		if isBranch(id) {
			onNode(id, true, depth)
			if t.IsBranchOpen(id) {
				if children := t.Children; children != nil {
					for _, c := range children(id) {
						t.walk(c, depth+1, onNode)
					}
				}
			}
		} else {
			onNode(id, false, depth)
		}
	}
}

// CreateRenderer is a private method to Fyne which links this widget to its renderer.
func (t *TreeContainer) CreateRenderer() fyne.WidgetRenderer {
	t.ExtendBaseWidget(t)
	c := newTreeContentContainer(t)
	s := NewScrollContainer(c)
	r := &treeContainerRenderer{
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

var _ fyne.WidgetRenderer = (*treeContainerRenderer)(nil)

type treeContainerRenderer struct {
	widget.BaseRenderer
	tree     *TreeContainer
	content  *treeContentContainer
	scroller *ScrollContainer
}

func (r *treeContainerRenderer) MinSize() (min fyne.Size) {
	min = r.scroller.MinSize()
	//log.Println("treeContainerRenderer.MinSize:", min)
	return
}

func (r *treeContainerRenderer) Layout(size fyne.Size) {
	//log.Println("treeContainerRenderer.Layout:", size)
	r.scroller.Resize(size)
}

func (r *treeContainerRenderer) Refresh() {
	//log.Println("treeContainerRenderer.Refresh")
	r.content.Refresh()
	s := r.tree.Size()
	if s.IsZero() {
		r.tree.Resize(r.tree.MinSize())
	}
	r.content.Resize(r.content.MinSize().Max(s))
	canvas.Refresh(r.tree.super())
}

var _ fyne.Widget = (*treeContentContainer)(nil)

type treeContentContainer struct {
	BaseWidget
	tree *TreeContainer
}

func newTreeContentContainer(tree *TreeContainer) (c *treeContentContainer) {
	c = &treeContentContainer{
		tree: tree,
	}
	c.ExtendBaseWidget(c)
	return
}

func (c *treeContentContainer) CreateRenderer() fyne.WidgetRenderer {
	return &treeContentRenderer{
		BaseRenderer: widget.BaseRenderer{},
		treeContent:  c,
		minSizes:     make(map[string]fyne.Size),
		branches:     make(map[string]*branch),
		leaves:       make(map[string]*leaf),
		branchPool:   &libPool{},
		leafPool:     &libPool{},
	}
}

var _ fyne.WidgetRenderer = (*treeContentRenderer)(nil)

type treeContentRenderer struct {
	widget.BaseRenderer
	treeContent *treeContentContainer
	minSizes    map[string]fyne.Size // Holds cache of children min sizes updated by MinSize and used by Layout
	branches    map[string]*branch
	leaves      map[string]*leaf
	branchPool  pool
	leafPool    pool
}

func (r *treeContentRenderer) MinSize() (min fyne.Size) {
	r.treeContent.propertyLock.Lock()
	defer r.treeContent.propertyLock.Unlock()

	r.treeContent.tree.Walk(func(id string, isBranch bool, depth int) {
		// Root node is not rendered unless it has been customized
		if r.treeContent.tree.Root == "" {
			depth = depth - 1
			if id == "" {
				// This is root node, skip
				return
			}
		}
		// TODO FIXME this assumes a node with the given ID will always have the same MinSize.
		m, ok := r.minSizes[id]
		if !ok {
			var c, n fyne.CanvasObject
			if isBranch {
				b := r.getBranch()
				b.Update(id, depth)
				c = b.Content()
				n = b
			} else {
				l := r.getLeaf()
				l.Update(id, depth)
				c = l.Content()
				n = l
			}
			if c != nil {
				r.treeContent.tree.UpdateNode(id, isBranch, c)
			}
			if n != nil {
				m = n.MinSize()
				if isBranch {
					r.branchPool.Release(n)
				} else {
					r.leafPool.Release(n)
				}
			} else {
				m = fyne.Size{}
			}
			r.minSizes[id] = m
		}
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

	offsetY := r.treeContent.tree.Offset.Y
	y := theme.Padding()
	// Walk open branches and obtain nodes to render in scroller's viewport
	r.treeContent.tree.Walk(func(id string, isBranch bool, depth int) {
		// Root node is not rendered unless it has been customized
		if r.treeContent.tree.Root == "" {
			depth = depth - 1
			if id == "" {
				// This is root node, skip
				return
			}
		}
		m, ok := r.minSizes[id]
		if !ok {
			m = fyne.Size{}
		}
		if y+m.Height < offsetY {
			// node is above viewport and not visible
			r.release(id, isBranch)
		} else if y > offsetY+size.Height {
			// node is below viewport and not visible
			r.release(id, isBranch)
		} else {
			// node is in viewport
			var n *treeNode
			if isBranch {
				b, ok := r.branches[id]
				if !ok {
					b = r.getBranch()
					b.Update(id, depth)
					r.treeContent.tree.UpdateNode(id, true, b.Content())
					r.branches[id] = b
				}
				n = b.treeNode
			} else {
				l, ok := r.leaves[id]
				if !ok {
					l = r.getLeaf()
					l.Update(id, depth)
					r.treeContent.tree.UpdateNode(id, false, l.Content())
					r.leaves[id] = l
				}
				n = l.treeNode
			}
			if n != nil {
				n.Move(fyne.NewPos(theme.Padding(), y))
				//log.Println("Node", id, "Pos", n.Position())
				n.Resize(fyne.NewSize(size.Width-n.Position().X-theme.Padding(), m.Height))
				//log.Println("Node", id, "Size", n.Size())
				n.Refresh()
			}
		}
		y += m.Height
	})
}

func (r *treeContentRenderer) Refresh() {
	//log.Println("treeContentRenderer.Refresh")
	if s := r.treeContent.Size(); s.IsZero() {
		r.treeContent.Resize(r.treeContent.MinSize())
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

func (r *treeContentRenderer) release(id string, isBranch bool) {
	if isBranch {
		if b, ok := r.branches[id]; ok {
			r.branchPool.Release(b)
			delete(r.branches, id)
		}
	} else {
		if l, ok := r.leaves[id]; ok {
			r.leafPool.Release(l)
			delete(r.leaves, id)
		}
	}
}

var _ fyne.CanvasObject = (*treeNode)(nil)
var _ fyne.Tappable = (*treeNode)(nil)
var _ desktop.Hoverable = (*treeNode)(nil)

type treeNode struct {
	BaseWidget
	tree    *TreeContainer
	id      string
	depth   int
	hovered bool
	icon    fyne.CanvasObject
	content fyne.CanvasObject
}

func (n *treeNode) Update(id string, depth int) {
	n.id = id
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
	n.tree.SetSelectedNode(n.id)
	if f := n.tree.OnNodeSelected; f != nil {
		f(n.id, n.content)
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
	//log.Println(r.treeNode.id, "treeNodeRenderer.Refresh")
	if i := r.treeNode.icon; i != nil {
		i.Refresh()
	}
	if i := r.divider; i != nil {
		i.FillColor = theme.ShadowColor()
		i.Refresh()
	}
	if i := r.indicator; i != nil {
		//log.Println(r.treeNode.id, "treeNodeRenderer.Refresh:", r.treeNode.tree.Selected)
		if r.treeNode.id == r.treeNode.tree.Selected {
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

func newBranch(tree *TreeContainer, content fyne.CanvasObject) (b *branch) {
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

func (b *branch) Update(id string, depth int) {
	b.treeNode.Update(id, depth)
	b.icon.(*branchIcon).id = id
}

func (b *branch) DoubleTapped(*fyne.PointEvent) {
	//log.Println("branch.DoubleTapped")
	b.tree.ToggleBranch(b.id)
}

var _ fyne.Tappable = (*branchIcon)(nil)
var _ fyne.DoubleTappable = (*branchIcon)(nil)

type branchIcon struct {
	Icon
	tree *TreeContainer
	id   string
}

func newBranchIcon(tree *TreeContainer) (i *branchIcon) {
	i = &branchIcon{
		tree: tree,
	}
	i.ExtendBaseWidget(i)
	return
}

func (i *branchIcon) Tapped(*fyne.PointEvent) {
	//log.Println("branchIcon.Tapped")
	i.tree.ToggleBranch(i.id)
}

func (i *branchIcon) DoubleTapped(*fyne.PointEvent) {
	//log.Println("branchIcon.DoubleTapped")
	// Do nothing - this stops the event propagating to branch
}

func (i *branchIcon) Refresh() {
	//log.Println(i.id, "branchIcon.Refresh")
	if i.tree.IsBranchOpen(i.id) {
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

func newLeaf(tree *TreeContainer, content fyne.CanvasObject) (l *leaf) {
	l = &leaf{
		&treeNode{
			tree:    tree,
			content: content,
		},
	}
	l.ExtendBaseWidget(l)
	return
}

type pool interface {
	Obtain() fyne.CanvasObject
	Release(fyne.CanvasObject)
}

var _ pool = (*libPool)(nil)

type libPool struct {
	sync.Pool
}

// Obtain returns an item from the pool for use
func (p *libPool) Obtain() (item fyne.CanvasObject) {
	o := p.Get()
	if o != nil {
		item = o.(fyne.CanvasObject)
		//item.Show()
	}
	return
}

// Release adds an item into the pool to be used later
func (p *libPool) Release(item fyne.CanvasObject) {
	//item.Hide()
	p.Put(item)
}

var _ pool = (*slicePool)(nil)

type slicePool struct {
	contents []fyne.CanvasObject
}

// Obtain returns an item from the pool for use
func (p *slicePool) Obtain() fyne.CanvasObject {
	//log.Println("slicePool:", len(p.contents))
	if len(p.contents) == 0 {
		return nil
	}
	item := p.contents[0]
	p.contents = p.contents[1:]
	//item.Show()
	return item
}

// Release adds an item into the pool to be used later
func (p *slicePool) Release(item fyne.CanvasObject) {
	//log.Println("Pool.Release:", item)
	p.contents = append(p.contents, item)
	//item.Hide()
}
