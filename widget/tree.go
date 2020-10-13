package widget

import (
	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/driver/desktop"
	"fyne.io/fyne/internal/cache"
	"fyne.io/fyne/internal/widget"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/storage"
	"fyne.io/fyne/theme"
)

const treeDividerHeight = 1

var _ fyne.Widget = (*Tree)(nil)

// Tree widget displays hierarchical data.
// Each node of the tree must be identified by a Unique ID.
type Tree struct {
	BaseWidget
	Root     string
	Selected string
	Offset   fyne.Position

	ChildUIDs      func(uid string) (c []string)                         // Return a sorted slice of Children Unique IDs for the given Node Unique ID
	IsBranch       func(uid string) (ok bool)                            // Return true if the given Unique ID represents a Branch
	CreateNode     func(branch bool) (o fyne.CanvasObject)               // Return a CanvasObject that can represent a Branch (if branch is true), or a Leaf (if branch is false)
	UpdateNode     func(uid string, branch bool, node fyne.CanvasObject) // Called to update the given CanvasObject to represent the data at the given Unique ID
	OnBranchOpened func(uid string)                                      // Called when a Branch is opened
	OnBranchClosed func(uid string)                                      // Called when a Branch is closed
	OnNodeSelected func(uid string)                                      // Called when the Node with the given Unique ID is selected.

	open          map[string]bool
	branchMinSize fyne.Size
	leafMinSize   fyne.Size
}

// NewTreeWithFiles creates a new tree with the given file system URI.
func NewTreeWithFiles(root fyne.URI) (t *Tree) {
	t = &Tree{
		Root: root.String(),
		ChildUIDs: func(uid string) (c []string) {
			luri, err := storage.ListerForURI(storage.NewURI(uid))
			if err != nil {
				fyne.LogError("Unable to get lister for "+uid, err)
			} else {
				uris, err := luri.List()
				if err != nil {
					fyne.LogError("Unable to list "+luri.String(), err)
				} else {
					for _, u := range uris {
						c = append(c, u.String())
					}
				}
			}
			return
		},
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
	return
}

// NewTreeWithStrings creates a new tree with the given string map.
// Data must contain a mapping for the root, which defaults to empty string ("").
func NewTreeWithStrings(data map[string][]string) (t *Tree) {
	t = &Tree{
		ChildUIDs: func(uid string) (c []string) {
			c = data[uid]
			return
		},
		IsBranch: func(uid string) (b bool) {
			_, b = data[uid]
			return
		},
		CreateNode: func(branch bool) fyne.CanvasObject {
			return NewLabel("Template Object")
		},
		UpdateNode: func(uid string, branch bool, node fyne.CanvasObject) {
			node.(*Label).SetText(uid)
		},
	}
	t.ExtendBaseWidget(t)
	return
}

// CloseAllBranches closes all branches in the tree.
func (t *Tree) CloseAllBranches() {
	t.propertyLock.Lock()
	t.open = make(map[string]bool)
	t.propertyLock.Unlock()
	t.Refresh()
}

// CloseBranch closes the branch with the given Unique ID.
func (t *Tree) CloseBranch(uid string) {
	t.ensureOpenMap()
	t.propertyLock.Lock()
	t.open[uid] = false
	t.propertyLock.Unlock()
	if f := t.OnBranchClosed; f != nil {
		f(uid)
	}
	t.Refresh()
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
		c.Refresh()
	}
	r.updateMinSizes()
	r.content.viewport = r.MinSize()
	return r
}

// IsBranchOpen returns true if the branch with the given Unique ID is expanded.
func (t *Tree) IsBranchOpen(uid string) bool {
	if uid == t.Root {
		return true // Root is always open
	}
	t.ensureOpenMap()
	t.propertyLock.RLock()
	defer t.propertyLock.RUnlock()
	return t.open[uid]
}

// MinSize returns the size that this widget should not shrink below.
func (t *Tree) MinSize() fyne.Size {
	t.ExtendBaseWidget(t)
	return t.BaseWidget.MinSize()
}

// OpenAllBranches opens all branches in the tree.
func (t *Tree) OpenAllBranches() {
	t.ensureOpenMap()
	t.walkAll(func(uid string, branch bool, depth int) {
		if branch {
			t.propertyLock.Lock()
			t.open[uid] = true
			t.propertyLock.Unlock()
		}
	})
	t.Refresh()
}

// OpenBranch opens the branch with the given Unique ID.
func (t *Tree) OpenBranch(uid string) {
	t.ensureOpenMap()
	t.propertyLock.Lock()
	t.open[uid] = true
	t.propertyLock.Unlock()
	if f := t.OnBranchOpened; f != nil {
		f(uid)
	}
	t.Refresh()
}

// Resize sets a new size for a widget.
func (t *Tree) Resize(size fyne.Size) {
	t.propertyLock.RLock()
	s := t.size
	t.propertyLock.RUnlock()

	if s == size {
		return
	}

	t.propertyLock.Lock()
	t.size = size
	t.propertyLock.Unlock()

	t.Refresh() // trigger a redraw
}

// SetSelectedNode updates the current selection to the node with the given Unique ID.
func (t *Tree) SetSelectedNode(uid string) {
	t.Selected = uid
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

func (t *Tree) ensureOpenMap() {
	t.propertyLock.Lock()
	defer t.propertyLock.Unlock()
	if t.open == nil {
		t.open = make(map[string]bool)
	}
}

func (t *Tree) walk(uid string, depth int, onNode func(string, bool, int)) {
	if isBranch := t.IsBranch; isBranch != nil {
		if isBranch(uid) {
			onNode(uid, true, depth)
			if t.IsBranchOpen(uid) {
				if childUIDs := t.ChildUIDs; childUIDs != nil {
					for _, c := range childUIDs(uid) {
						t.walk(c, depth+1, onNode)
					}
				}
			}
		} else {
			onNode(uid, false, depth)
		}
	}
}

// walkAll visits every open node of the tree and calls the given callback with node Unique ID, whether node is branch, and the depth of node.
func (t *Tree) walkAll(onNode func(string, bool, int)) {
	t.walk(t.Root, 0, onNode)
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
	min = min.Max(r.tree.branchMinSize)
	min = min.Max(r.tree.leafMinSize)
	return
}

func (r *treeRenderer) Layout(size fyne.Size) {
	r.content.viewport = size
	r.scroller.Resize(size)
}

func (r *treeRenderer) Refresh() {
	r.updateMinSizes()
	s := r.tree.Size()
	if s.IsZero() {
		r.tree.Resize(r.tree.MinSize())
	} else {
		r.Layout(s)
	}
	r.scroller.Refresh()
	r.content.Refresh()
	canvas.Refresh(r.tree.super())
}

func (r *treeRenderer) updateMinSizes() {
	if f := r.tree.CreateNode; f != nil {
		r.tree.branchMinSize = newBranch(r.tree, f(true)).MinSize()
		r.tree.leafMinSize = newLeaf(r.tree, f(false)).MinSize()
	}
}

var _ fyne.Widget = (*treeContent)(nil)

type treeContent struct {
	BaseWidget
	tree     *Tree
	viewport fyne.Size
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
		branches:     make(map[string]*branch),
		leaves:       make(map[string]*leaf),
		branchPool:   &syncPool{},
		leafPool:     &syncPool{},
	}
}

func (c *treeContent) Resize(size fyne.Size) {
	c.propertyLock.RLock()
	s := c.size
	c.propertyLock.RUnlock()

	if s == size {
		return
	}

	c.propertyLock.Lock()
	c.size = size
	c.propertyLock.Unlock()

	c.Refresh() // trigger a redraw
}

var _ fyne.WidgetRenderer = (*treeContentRenderer)(nil)

type treeContentRenderer struct {
	widget.BaseRenderer
	treeContent *treeContent
	dividers    []*canvas.Rectangle
	branches    map[string]*branch
	leaves      map[string]*leaf
	branchPool  pool
	leafPool    pool
}

func (r *treeContentRenderer) Layout(size fyne.Size) {
	r.treeContent.propertyLock.Lock()
	defer r.treeContent.propertyLock.Unlock()

	branches := make(map[string]*branch)
	leaves := make(map[string]*leaf)

	offsetY := r.treeContent.tree.Offset.Y
	viewport := r.treeContent.viewport
	width := fyne.Max(size.Width, viewport.Width)
	y := 0
	numDividers := 0
	// walkAll open branches and obtain nodes to render in scroller's viewport
	r.treeContent.tree.walkAll(func(uid string, isBranch bool, depth int) {
		// Root node is not rendered unless it has been customized
		if r.treeContent.tree.Root == "" {
			depth = depth - 1
			if uid == "" {
				// This is root node, skip
				return
			}
		}

		// If this is not the first item, add a divider
		if y > 0 {
			var divider *canvas.Rectangle
			if numDividers < len(r.dividers) {
				divider = r.dividers[numDividers]
			} else {
				divider = canvas.NewRectangle(theme.ShadowColor())
				r.dividers = append(r.dividers, divider)
			}
			divider.Move(fyne.NewPos(theme.Padding(), y))
			s := fyne.NewSize(width-2*theme.Padding(), treeDividerHeight)
			divider.SetMinSize(s)
			divider.Resize(s)
			divider.Show()
			y += treeDividerHeight
			numDividers++
		}

		m := r.treeContent.tree.leafMinSize
		if isBranch {
			m = r.treeContent.tree.branchMinSize
		}
		if y+m.Height < offsetY {
			// Node is above viewport and not visible
		} else if y > offsetY+viewport.Height {
			// Node is below viewport and not visible
		} else {
			// Node is in viewport
			var n *treeNode
			if isBranch {
				b, ok := r.branches[uid]
				if !ok {
					b = r.getBranch()
					b.update(uid, depth)
					if f := r.treeContent.tree.UpdateNode; f != nil {
						f(uid, true, b.Content())
					}
				}
				branches[uid] = b
				n = b.treeNode
			} else {
				l, ok := r.leaves[uid]
				if !ok {
					l = r.getLeaf()
					l.update(uid, depth)
					if f := r.treeContent.tree.UpdateNode; f != nil {
						f(uid, false, l.Content())
					}
				}
				leaves[uid] = l
				n = l.treeNode
			}
			if n != nil {
				n.Move(fyne.NewPos(0, y))
				n.Resize(fyne.NewSize(width, m.Height))
			}
		}
		y += m.Height
	})

	// Hide any dividers that haven't been reused
	for ; numDividers < len(r.dividers); numDividers++ {
		r.dividers[numDividers].Hide()
	}

	// Release any nodes that haven't been reused
	for uid, b := range r.branches {
		if _, ok := branches[uid]; !ok {
			b.Hide()
			r.branchPool.Release(b)
		}
	}
	for uid, l := range r.leaves {
		if _, ok := leaves[uid]; !ok {
			l.Hide()
			r.leafPool.Release(l)
		}
	}

	r.branches = branches
	r.leaves = leaves
}

func (r *treeContentRenderer) MinSize() (min fyne.Size) {
	r.treeContent.propertyLock.Lock()
	defer r.treeContent.propertyLock.Unlock()

	r.treeContent.tree.walkAll(func(uid string, isBranch bool, depth int) {
		// Root node is not rendered unless it has been customized
		if r.treeContent.tree.Root == "" {
			depth = depth - 1
			if uid == "" {
				// This is root node, skip
				return
			}
		}

		// If this is not the first item, add a divider
		if min.Height > 0 {
			min.Height += treeDividerHeight
		}

		m := r.treeContent.tree.leafMinSize
		if isBranch {
			m = r.treeContent.tree.branchMinSize
		}
		m.Width += depth * (theme.IconInlineSize() + theme.Padding())
		min.Width = fyne.Max(min.Width, m.Width)
		min.Height += m.Height
	})
	return
}

func (r *treeContentRenderer) Objects() (objects []fyne.CanvasObject) {
	r.treeContent.propertyLock.RLock()
	for _, d := range r.dividers {
		objects = append(objects, d)
	}
	for _, b := range r.branches {
		objects = append(objects, b)
	}
	for _, l := range r.leaves {
		objects = append(objects, l)
	}
	r.treeContent.propertyLock.RUnlock()
	return
}

func (r *treeContentRenderer) Refresh() {
	s := r.treeContent.Size()
	if s.IsZero() {
		r.treeContent.Resize(r.treeContent.MinSize().Max(r.treeContent.tree.Size()))
	} else {
		r.Layout(s)
	}
	r.treeContent.propertyLock.RLock()
	for _, d := range r.dividers {
		d.FillColor = theme.ShadowColor()
	}
	for _, b := range r.branches {
		b.Refresh()
	}
	for _, l := range r.leaves {
		l.Refresh()
	}
	r.treeContent.propertyLock.RUnlock()
	canvas.Refresh(r.treeContent.super())
}

func (r *treeContentRenderer) getBranch() (b *branch) {
	o := r.branchPool.Obtain()
	if o != nil {
		b = o.(*branch)
	} else {
		var content fyne.CanvasObject
		if f := r.treeContent.tree.CreateNode; f != nil {
			content = f(true)
		}
		b = newBranch(r.treeContent.tree, content)
	}
	return
}

func (r *treeContentRenderer) getLeaf() (l *leaf) {
	o := r.leafPool.Obtain()
	if o != nil {
		l = o.(*leaf)
	} else {
		var content fyne.CanvasObject
		if f := r.treeContent.tree.CreateNode; f != nil {
			content = f(false)
		}
		l = newLeaf(r.treeContent.tree, content)
	}
	return
}

var _ desktop.Hoverable = (*treeNode)(nil)
var _ fyne.CanvasObject = (*treeNode)(nil)
var _ fyne.Tappable = (*treeNode)(nil)

type treeNode struct {
	BaseWidget
	tree    *Tree
	uid     string
	depth   int
	hovered bool
	icon    fyne.CanvasObject
	content fyne.CanvasObject
}

func (n *treeNode) Content() fyne.CanvasObject {
	return n.content
}

func (n *treeNode) CreateRenderer() fyne.WidgetRenderer {
	return &treeNodeRenderer{
		BaseRenderer: widget.BaseRenderer{},
		treeNode:     n,
		indicator:    canvas.NewRectangle(theme.BackgroundColor()),
	}
}

func (n *treeNode) Indent() int {
	return n.depth * (theme.IconInlineSize() + theme.Padding())
}

// MouseIn is called when a desktop pointer enters the widget
func (n *treeNode) MouseIn(*desktop.MouseEvent) {
	n.hovered = true
	n.partialRefresh()
}

// MouseMoved is called when a desktop pointer hovers over the widget
func (n *treeNode) MouseMoved(*desktop.MouseEvent) {
}

// MouseOut is called when a desktop pointer exits the widget
func (n *treeNode) MouseOut() {
	n.hovered = false
	n.partialRefresh()
}

func (n *treeNode) Tapped(*fyne.PointEvent) {
	n.tree.SetSelectedNode(n.uid)
	if f := n.tree.OnNodeSelected; f != nil {
		f(n.uid)
	}
}

func (n *treeNode) partialRefresh() {
	if r := cache.Renderer(n.super()); r != nil {
		r.(*treeNodeRenderer).partialRefresh()
	}
}

func (n *treeNode) update(uid string, depth int) {
	n.uid = uid
	n.depth = depth
	n.propertyLock.Lock()
	n.Hidden = false
	n.propertyLock.Unlock()
	n.partialRefresh()
}

var _ fyne.WidgetRenderer = (*treeNodeRenderer)(nil)

type treeNodeRenderer struct {
	widget.BaseRenderer
	treeNode  *treeNode
	indicator *canvas.Rectangle
}

func (r *treeNodeRenderer) Layout(size fyne.Size) {
	x := 0
	y := 0
	r.indicator.Move(fyne.NewPos(x, y))
	s := fyne.NewSize(theme.Padding(), size.Height)
	r.indicator.SetMinSize(s)
	r.indicator.Resize(s)
	h := size.Height - 2*theme.Padding()
	x += theme.Padding() + r.treeNode.Indent()
	y += theme.Padding()
	if r.treeNode.icon != nil {
		r.treeNode.icon.Move(fyne.NewPos(x, y))
		r.treeNode.icon.Resize(fyne.NewSize(theme.IconInlineSize(), h))
	}
	x += theme.IconInlineSize()
	x += theme.Padding()
	if r.treeNode.content != nil {
		r.treeNode.content.Move(fyne.NewPos(x, y))
		r.treeNode.content.Resize(fyne.NewSize(size.Width-x-theme.Padding(), h))
	}
}

func (r *treeNodeRenderer) MinSize() (min fyne.Size) {
	if r.treeNode.content != nil {
		min = r.treeNode.content.MinSize()
	}
	min.Width += theme.Padding() + r.treeNode.Indent() + theme.IconInlineSize()
	min.Width += 2 * theme.Padding()
	min.Height = fyne.Max(min.Height, theme.IconInlineSize())
	min.Height += 2 * theme.Padding()
	return
}

func (r *treeNodeRenderer) Objects() (objects []fyne.CanvasObject) {
	if r.treeNode.content != nil {
		objects = append(objects, r.treeNode.content)
	}
	if r.treeNode.icon != nil {
		objects = append(objects, r.treeNode.icon)
	}
	objects = append(objects, r.indicator)
	return
}

func (r *treeNodeRenderer) Refresh() {
	if c := r.treeNode.content; c != nil {
		c.Refresh()
	}
	r.partialRefresh()
}

func (r *treeNodeRenderer) partialRefresh() {
	if r.treeNode.icon != nil {
		r.treeNode.icon.Refresh()
	}
	if r.treeNode.uid == r.treeNode.tree.Selected {
		r.indicator.FillColor = theme.PrimaryColor()
	} else if r.treeNode.hovered {
		r.indicator.FillColor = theme.HoverColor()
	} else {
		r.indicator.FillColor = theme.BackgroundColor()
	}
	r.indicator.Refresh()
	canvas.Refresh(r.treeNode.super())
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

func (b *branch) DoubleTapped(*fyne.PointEvent) {
	b.tree.ToggleBranch(b.uid)
}

func (b *branch) update(uid string, depth int) {
	b.treeNode.update(uid, depth)
	b.icon.(*branchIcon).update(uid, depth)
}

var _ fyne.DoubleTappable = (*branchIcon)(nil)
var _ fyne.Tappable = (*branchIcon)(nil)

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

func (i *branchIcon) DoubleTapped(*fyne.PointEvent) {
	// Do nothing - this stops the event propagating to branch
}

func (i *branchIcon) Refresh() {
	if i.tree.IsBranchOpen(i.uid) {
		i.Resource = theme.MoveDownIcon()
	} else {
		i.Resource = theme.NavigateNextIcon()
	}
	i.Icon.Refresh()
}

func (i *branchIcon) Tapped(*fyne.PointEvent) {
	i.tree.ToggleBranch(i.uid)
}

func (i *branchIcon) update(uid string, depth int) {
	i.uid = uid
	i.Refresh()
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
