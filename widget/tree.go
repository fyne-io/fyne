package widget

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/internal/cache"
	"fyne.io/fyne/v2/internal/widget"
	"fyne.io/fyne/v2/theme"
)

// allTreeNodesID represents all tree nodes when refreshing requested nodes
const allTreeNodesID = "_ALLNODES"

// TreeNodeID represents the unique id of a tree node.
type TreeNodeID = string

// Declare conformity with interfaces
var _ fyne.Focusable = (*Tree)(nil)
var _ fyne.Widget = (*Tree)(nil)

// Tree widget displays hierarchical data.
// Each node of the tree must be identified by a Unique TreeNodeID.
//
// Since: 1.4
type Tree struct {
	BaseWidget
	Root TreeNodeID

	ChildUIDs      func(uid TreeNodeID) (c []TreeNodeID)                     `json:"-"` // Return a sorted slice of Children TreeNodeIDs for the given Node TreeNodeID
	CreateNode     func(branch bool) (o fyne.CanvasObject)                   `json:"-"` // Return a CanvasObject that can represent a Branch (if branch is true), or a Leaf (if branch is false)
	IsBranch       func(uid TreeNodeID) (ok bool)                            `json:"-"` // Return true if the given TreeNodeID represents a Branch
	OnBranchClosed func(uid TreeNodeID)                                      `json:"-"` // Called when a Branch is closed
	OnBranchOpened func(uid TreeNodeID)                                      `json:"-"` // Called when a Branch is opened
	OnSelected     func(uid TreeNodeID)                                      `json:"-"` // Called when the Node with the given TreeNodeID is selected.
	OnUnselected   func(uid TreeNodeID)                                      `json:"-"` // Called when the Node with the given TreeNodeID is unselected.
	UpdateNode     func(uid TreeNodeID, branch bool, node fyne.CanvasObject) `json:"-"` // Called to update the given CanvasObject to represent the data at the given TreeNodeID

	branchMinSize fyne.Size
	currentFocus  TreeNodeID
	focused       bool
	leafMinSize   fyne.Size
	offset        fyne.Position
	open          map[TreeNodeID]bool
	scroller      *widget.Scroll
	selected      []TreeNodeID
}

// NewTree returns a new performant tree widget defined by the passed functions.
// childUIDs returns the child TreeNodeIDs of the given node.
// isBranch returns true if the given node is a branch, false if it is a leaf.
// create returns a new template object that can be cached.
// update is used to apply data at specified data location to the passed template CanvasObject.
//
// Since: 1.4
func NewTree(childUIDs func(TreeNodeID) []TreeNodeID, isBranch func(TreeNodeID) bool, create func(bool) fyne.CanvasObject, update func(TreeNodeID, bool, fyne.CanvasObject)) *Tree {
	t := &Tree{ChildUIDs: childUIDs, IsBranch: isBranch, CreateNode: create, UpdateNode: update}
	t.ExtendBaseWidget(t)
	return t
}

// NewTreeWithData creates a new tree widget that will display the contents of the provided data.
//
// Since: 2.4
func NewTreeWithData(data binding.DataTree, createItem func(bool) fyne.CanvasObject, updateItem func(binding.DataItem, bool, fyne.CanvasObject)) *Tree {
	t := NewTree(
		data.ChildIDs,
		func(id TreeNodeID) bool {
			children := data.ChildIDs(id)
			return len(children) > 0
		},
		createItem,
		func(i TreeNodeID, branch bool, o fyne.CanvasObject) {
			item, err := data.GetItem(i)
			if err != nil {
				fyne.LogError(fmt.Sprintf("Error getting data item %s", i), err)
				return
			}
			updateItem(item, branch, o)
		})

	data.AddListener(binding.NewDataListener(t.Refresh))
	return t
}

// NewTreeWithStrings creates a new tree with the given string map.
// Data must contain a mapping for the root, which defaults to empty string ("").
//
// Since: 1.4
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
	t.open = make(map[TreeNodeID]bool)
	t.propertyLock.Unlock()
	t.Refresh()
}

// CloseBranch closes the branch with the given TreeNodeID.
func (t *Tree) CloseBranch(uid TreeNodeID) {
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
	s := widget.NewScroll(c)
	t.scroller = s
	r := &treeRenderer{
		BaseRenderer: widget.NewBaseRenderer([]fyne.CanvasObject{s}),
		tree:         t,
		content:      c,
		scroller:     s,
	}
	s.OnScrolled = t.offsetUpdated
	r.updateMinSizes()
	r.content.viewport = r.MinSize()
	return r
}

// IsBranchOpen returns true if the branch with the given TreeNodeID is expanded.
func (t *Tree) IsBranchOpen(uid TreeNodeID) bool {
	if uid == t.Root {
		return true // Root is always open
	}
	t.ensureOpenMap()
	t.propertyLock.RLock()
	defer t.propertyLock.RUnlock()
	return t.open[uid]
}

// FocusGained is called after this Tree has gained focus.
//
// Implements: fyne.Focusable
func (t *Tree) FocusGained() {
	if t.currentFocus == "" {
		if childUIDs := t.ChildUIDs; childUIDs != nil {
			if ids := childUIDs(""); len(ids) > 0 {
				t.currentFocus = ids[0]
			}
		}
	}

	t.focused = true
	t.ScrollTo(t.currentFocus)
	t.RefreshItem(t.currentFocus)
}

// FocusLost is called after this Tree has lost focus.
//
// Implements: fyne.Focusable
func (t *Tree) FocusLost() {
	t.focused = false
	t.Refresh() //Item(t.currentFocus)
}

// MinSize returns the size that this widget should not shrink below.
func (t *Tree) MinSize() fyne.Size {
	t.ExtendBaseWidget(t)
	return t.BaseWidget.MinSize()
}

// RefreshItem refreshes a single item, specified by the item ID passed in.
//
// Since: 2.4
func (t *Tree) RefreshItem(id TreeNodeID) {
	if t.scroller == nil {
		return
	}
	r := cache.Renderer(t.scroller.Content.(*treeContent))
	if r == nil {
		return
	}

	r.(*treeContentRenderer).refreshForID(id)
}

// OpenAllBranches opens all branches in the tree.
func (t *Tree) OpenAllBranches() {
	t.ensureOpenMap()
	t.walkAll(func(uid, parent TreeNodeID, branch bool, depth int) {
		if branch {
			t.propertyLock.Lock()
			t.open[uid] = true
			t.propertyLock.Unlock()
		}
	})
	t.Refresh()
}

// OpenBranch opens the branch with the given TreeNodeID.
func (t *Tree) OpenBranch(uid TreeNodeID) {
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

// ScrollToBottom scrolls to the bottom of the tree.
//
// Since 2.1
func (t *Tree) ScrollToBottom() {
	if t.scroller == nil {
		return
	}

	y, size := t.findBottom()
	t.scroller.Offset.Y = y + size.Height - t.scroller.Size().Height

	t.offsetUpdated(t.scroller.Offset)
	t.Refresh()
}

// ScrollTo scrolls to the node with the given id.
//
// Since 2.1
func (t *Tree) ScrollTo(uid TreeNodeID) {
	if t.scroller == nil {
		return
	}

	y, size, ok := t.offsetAndSize(uid)
	if !ok {
		return
	}

	// TODO scrolling to a node should open all parents if they aren't already
	if y < t.scroller.Offset.Y {
		t.scroller.Offset.Y = y
	} else if y+size.Height > t.scroller.Offset.Y+t.scroller.Size().Height {
		t.scroller.Offset.Y = y + size.Height - t.scroller.Size().Height
	}

	t.offsetUpdated(t.scroller.Offset)
	t.Refresh()
}

// ScrollToTop scrolls to the top of the tree.
//
// Since 2.1
func (t *Tree) ScrollToTop() {
	if t.scroller == nil {
		return
	}

	t.scroller.Offset.Y = 0
	t.offsetUpdated(t.scroller.Offset)
	t.Refresh()
}

// Select marks the specified node to be selected.
func (t *Tree) Select(uid TreeNodeID) {
	if len(t.selected) > 0 {
		if uid == t.selected[0] {
			return // no change
		}
		if f := t.OnUnselected; f != nil {
			f(t.selected[0])
		}
	}
	t.selected = []TreeNodeID{uid}
	t.ScrollTo(uid)
	if f := t.OnSelected; f != nil {
		f(uid)
	}
}

// ToggleBranch flips the state of the branch with the given TreeNodeID.
func (t *Tree) ToggleBranch(uid string) {
	if t.IsBranchOpen(uid) {
		t.CloseBranch(uid)
	} else {
		t.OpenBranch(uid)
	}
}

// TypedKey is called if a key event happens while this Tree is focused.
//
// Implements: fyne.Focusable
func (t *Tree) TypedKey(event *fyne.KeyEvent) {
	switch event.Name {
	case fyne.KeySpace:
		t.Select(t.currentFocus)
	case fyne.KeyDown:
		t.RefreshItem(t.currentFocus)
		next := false
		t.walk(t.Root, "", 0, func(id, p TreeNodeID, _ bool, _ int) {
			if next {
				t.currentFocus = id
				next = false
			} else if id == t.currentFocus {
				next = true
			}
		})

		t.ScrollTo(t.currentFocus)
		t.RefreshItem(t.currentFocus)
	case fyne.KeyLeft:
		// If the current focus is on a branch which is open, just close it
		if t.IsBranch(t.currentFocus) && t.IsBranchOpen(t.currentFocus) {
			t.CloseBranch(t.currentFocus)
		} else {
			// Every other case should move the focus to the current parent node
			t.walk(t.Root, "", 0, func(id, p TreeNodeID, _ bool, _ int) {
				if id == t.currentFocus && p != "" {
					t.currentFocus = p
				}
			})
		}

		t.RefreshItem(t.currentFocus)
		t.ScrollTo(t.currentFocus)
		t.RefreshItem(t.currentFocus)
	case fyne.KeyRight:
		if t.IsBranch(t.currentFocus) {
			t.OpenBranch(t.currentFocus)
		}
		children := []TreeNodeID{}
		if childUIDs := t.ChildUIDs; childUIDs != nil {
			children = childUIDs(t.currentFocus)
		}

		if len(children) > 0 {
			t.currentFocus = children[0]
		}

		t.RefreshItem(t.currentFocus)
		t.ScrollTo(t.currentFocus)
		t.RefreshItem(t.currentFocus)
	case fyne.KeyUp:
		t.RefreshItem(t.currentFocus)
		previous := ""
		t.walk(t.Root, "", 0, func(id, p TreeNodeID, _ bool, _ int) {
			if id == t.currentFocus && previous != "" {
				t.currentFocus = previous
			}
			previous = id
		})

		t.ScrollTo(t.currentFocus)
		t.RefreshItem(t.currentFocus)
	}
}

// TypedRune is called if a text event happens while this Tree is focused.
//
// Implements: fyne.Focusable
func (t *Tree) TypedRune(_ rune) {
	// intentionally left blank
}

// Unselect marks the specified node to be not selected.
func (t *Tree) Unselect(uid TreeNodeID) {
	if len(t.selected) == 0 || t.selected[0] != uid {
		return
	}

	t.selected = nil
	t.Refresh()
	if f := t.OnUnselected; f != nil {
		f(uid)
	}
}

// UnselectAll sets all nodes to be not selected.
//
// Since: 2.1
func (t *Tree) UnselectAll() {
	if len(t.selected) == 0 {
		return
	}

	selected := t.selected
	t.selected = nil
	t.Refresh()
	if f := t.OnUnselected; f != nil {
		for _, uid := range selected {
			f(uid)
		}
	}
}

func (t *Tree) ensureOpenMap() {
	t.propertyLock.Lock()
	defer t.propertyLock.Unlock()
	if t.open == nil {
		t.open = make(map[string]bool)
	}
}

func (t *Tree) findBottom() (y float32, size fyne.Size) {
	sep := theme.Padding()
	t.walkAll(func(id, _ TreeNodeID, branch bool, _ int) {
		size = t.leafMinSize
		if branch {
			size = t.branchMinSize
		}

		// Root node is not rendered unless it has been customized
		if t.Root == "" && id == "" {
			// This is root node, skip
			return
		}

		// If this is not the first item, add a separator
		if y > 0 {
			y += sep
		}

		y += size.Height
	})
	if y > 0 {
		y -= sep
	}
	return
}

func (t *Tree) offsetAndSize(uid TreeNodeID) (y float32, size fyne.Size, found bool) {
	t.walkAll(func(id, _ TreeNodeID, branch bool, _ int) {
		m := t.leafMinSize
		if branch {
			m = t.branchMinSize
		}
		if id == uid {
			found = true
			size = m
		} else if !found {
			// Root node is not rendered unless it has been customized
			if t.Root == "" && id == "" {
				// This is root node, skip
				return
			}
			// If this is not the first item, add a separator
			if y > 0 {
				y += theme.Padding()
			}

			y += m.Height
		}
	})
	return
}

func (t *Tree) offsetUpdated(pos fyne.Position) {
	if t.offset == pos {
		return
	}
	t.offset = pos
	t.scroller.Content.Refresh()
}

func (t *Tree) walk(uid, parent TreeNodeID, depth int, onNode func(TreeNodeID, TreeNodeID, bool, int)) {
	if isBranch := t.IsBranch; isBranch != nil {
		if isBranch(uid) {
			onNode(uid, parent, true, depth)
			if t.IsBranchOpen(uid) {
				if childUIDs := t.ChildUIDs; childUIDs != nil {
					for _, c := range childUIDs(uid) {
						t.walk(c, uid, depth+1, onNode)
					}
				}
			}
		} else {
			onNode(uid, parent, false, depth)
		}
	}
}

// walkAll visits every open node of the tree and calls the given callback with TreeNodeID, whether node is branch, and the depth of node.
func (t *Tree) walkAll(onNode func(TreeNodeID, TreeNodeID, bool, int)) {
	t.walk(t.Root, "", 0, onNode)
}

var _ fyne.WidgetRenderer = (*treeRenderer)(nil)

type treeRenderer struct {
	widget.BaseRenderer
	tree     *Tree
	content  *treeContent
	scroller *widget.Scroll
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
	separators  []fyne.CanvasObject
	objects     []fyne.CanvasObject
	branches    map[string]*branch
	leaves      map[string]*leaf
	branchPool  pool
	leafPool    pool
}

func (r *treeContentRenderer) Layout(size fyne.Size) {
	r.treeContent.propertyLock.Lock()
	defer r.treeContent.propertyLock.Unlock()

	r.objects = nil
	branches := make(map[string]*branch)
	leaves := make(map[string]*leaf)

	pad := theme.Padding()
	offsetY := r.treeContent.tree.offset.Y
	viewport := r.treeContent.viewport
	width := fyne.Max(size.Width, viewport.Width)
	separatorCount := 0
	separatorThickness := theme.SeparatorThicknessSize()
	separatorSize := fyne.NewSize(width, separatorThickness)
	separatorOff := (pad + separatorThickness) / 2
	y := float32(0)
	// walkAll open branches and obtain nodes to render in scroller's viewport
	r.treeContent.tree.walkAll(func(uid, _ string, isBranch bool, depth int) {
		// Root node is not rendered unless it has been customized
		if r.treeContent.tree.Root == "" {
			depth = depth - 1
			if uid == "" {
				// This is root node, skip
				return
			}
		}

		// If this is not the first item, add a separator
		addSeparator := y > 0
		if addSeparator {
			y += pad
			separatorCount++
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

			if addSeparator {
				var separator fyne.CanvasObject
				if separatorCount < len(r.separators) {
					separator = r.separators[separatorCount]
				} else {
					separator = NewSeparator()
					r.separators = append(r.separators, separator)
				}
				separator.Move(fyne.NewPos(0, y-separatorOff))
				separator.Resize(separatorSize)
				r.objects = append(r.objects, separator)
				separatorCount++
			}

			var n fyne.CanvasObject
			if isBranch {
				b, ok := r.branches[uid]
				if !ok {
					b = r.getBranch()
					if f := r.treeContent.tree.UpdateNode; f != nil {
						f(uid, true, b.Content())
					}
					b.update(uid, depth)
				}
				branches[uid] = b
				n = b
				r.objects = append(r.objects, b)
			} else {
				l, ok := r.leaves[uid]
				if !ok {
					l = r.getLeaf()
					if f := r.treeContent.tree.UpdateNode; f != nil {
						f(uid, false, l.Content())
					}
					l.update(uid, depth)
				}
				leaves[uid] = l
				n = l
				r.objects = append(r.objects, l)
			}
			if n != nil {
				n.Move(fyne.NewPos(0, y))
				n.Resize(fyne.NewSize(width, m.Height))
			}
		}
		y += m.Height
	})

	// Hide any separators that haven't been reused
	for ; separatorCount < len(r.separators); separatorCount++ {
		r.separators[separatorCount].Hide()
	}

	// Release any nodes that haven't been reused
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

func (r *treeContentRenderer) MinSize() (min fyne.Size) {
	r.treeContent.propertyLock.Lock()
	defer r.treeContent.propertyLock.Unlock()

	r.treeContent.tree.walkAll(func(uid, _ string, isBranch bool, depth int) {
		// Root node is not rendered unless it has been customized
		if r.treeContent.tree.Root == "" {
			depth = depth - 1
			if uid == "" {
				// This is root node, skip
				return
			}
		}

		// If this is not the first item, add a separator
		if min.Height > 0 {
			min.Height += theme.Padding()
		}

		m := r.treeContent.tree.leafMinSize
		if isBranch {
			m = r.treeContent.tree.branchMinSize
		}
		m.Width += float32(depth) * (theme.IconInlineSize() + theme.Padding())
		min.Width = fyne.Max(min.Width, m.Width)
		min.Height += m.Height
	})
	return
}

func (r *treeContentRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

func (r *treeContentRenderer) Refresh() {
	r.refreshForID(allTreeNodesID)
}

func (r *treeContentRenderer) refreshForID(toDraw TreeNodeID) {
	s := r.treeContent.Size()
	if s.IsZero() {
		r.treeContent.Resize(r.treeContent.MinSize().Max(r.treeContent.tree.Size()))
	} else {
		r.Layout(s)
	}
	r.treeContent.propertyLock.RLock()
	for id, b := range r.branches {
		if toDraw != allTreeNodesID && id != toDraw {
			continue
		}

		b.Refresh()
	}
	for id, l := range r.leaves {
		if toDraw != allTreeNodesID && id != toDraw {
			continue
		}

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
	tree     *Tree
	uid      string
	depth    int
	hovered  bool
	icon     fyne.CanvasObject
	isBranch bool
	content  fyne.CanvasObject
}

func (n *treeNode) Content() fyne.CanvasObject {
	return n.content
}

func (n *treeNode) CreateRenderer() fyne.WidgetRenderer {
	background := canvas.NewRectangle(theme.HoverColor())
	background.CornerRadius = theme.SelectionRadiusSize()
	background.Hide()
	return &treeNodeRenderer{
		BaseRenderer: widget.BaseRenderer{},
		treeNode:     n,
		background:   background,
	}
}

func (n *treeNode) Indent() float32 {
	return float32(n.depth) * (theme.IconInlineSize() + theme.Padding())
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
	if n.tree.currentFocus != "" {
		n.tree.RefreshItem(n.tree.currentFocus)
	}

	n.tree.Select(n.uid)
	if !fyne.CurrentDevice().IsMobile() {
		canvas := fyne.CurrentApp().Driver().CanvasForObject(n.tree)
		if canvas != nil {
			canvas.Focus(n.tree)
		}
		n.tree.currentFocus = n.uid
		n.Refresh()
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
	treeNode   *treeNode
	background *canvas.Rectangle
}

func (r *treeNodeRenderer) Layout(size fyne.Size) {
	x := theme.Padding() + r.treeNode.Indent()
	y := float32(0)
	r.background.Resize(size)
	if r.treeNode.icon != nil {
		r.treeNode.icon.Move(fyne.NewPos(x, y))
		r.treeNode.icon.Resize(fyne.NewSize(theme.IconInlineSize(), size.Height))
	}
	x += theme.IconInlineSize()
	x += theme.Padding()
	if r.treeNode.content != nil {
		r.treeNode.content.Move(fyne.NewPos(x, y))
		r.treeNode.content.Resize(fyne.NewSize(size.Width-x, size.Height))
	}
}

func (r *treeNodeRenderer) MinSize() (min fyne.Size) {
	if r.treeNode.content != nil {
		min = r.treeNode.content.MinSize()
	}
	min.Width += theme.InnerPadding() + r.treeNode.Indent() + theme.IconInlineSize()
	min.Height = fyne.Max(min.Height, theme.IconInlineSize())
	return
}

func (r *treeNodeRenderer) Objects() (objects []fyne.CanvasObject) {
	objects = append(objects, r.background)
	if r.treeNode.content != nil {
		objects = append(objects, r.treeNode.content)
	}
	if r.treeNode.icon != nil {
		objects = append(objects, r.treeNode.icon)
	}
	return
}

func (r *treeNodeRenderer) Refresh() {
	if c := r.treeNode.content; c != nil {
		if f := r.treeNode.tree.UpdateNode; f != nil {
			f(r.treeNode.uid, r.treeNode.isBranch, c)
		}
	}
	r.partialRefresh()
}

func (r *treeNodeRenderer) partialRefresh() {
	if r.treeNode.icon != nil {
		r.treeNode.icon.Refresh()
	}
	r.background.CornerRadius = theme.SelectionRadiusSize()
	if len(r.treeNode.tree.selected) > 0 && r.treeNode.uid == r.treeNode.tree.selected[0] {
		r.background.FillColor = theme.SelectionColor()
		r.background.Show()
	} else if r.treeNode.hovered || (r.treeNode.tree.focused && r.treeNode.tree.currentFocus == r.treeNode.uid) {
		r.background.FillColor = theme.HoverColor()
		r.background.Show()
	} else {
		r.background.Hide()
	}
	r.background.Refresh()
	r.Layout(r.treeNode.size)
	canvas.Refresh(r.treeNode.super())
}

var _ fyne.Widget = (*branch)(nil)

type branch struct {
	*treeNode
}

func newBranch(tree *Tree, content fyne.CanvasObject) (b *branch) {
	b = &branch{
		treeNode: &treeNode{
			tree:     tree,
			icon:     newBranchIcon(tree),
			isBranch: true,
			content:  content,
		},
	}
	b.ExtendBaseWidget(b)
	return
}

func (b *branch) update(uid string, depth int) {
	b.treeNode.update(uid, depth)
	b.icon.(*branchIcon).update(uid, depth)
}

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
			tree:     tree,
			content:  content,
			isBranch: false,
		},
	}
	l.ExtendBaseWidget(l)
	return
}
