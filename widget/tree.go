package widget

import (
	"image/color"
	"log"
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

const treeDividerHeight = 1

var _ fyne.Widget = (*Tree)(nil)

// Tree widget displays hierarchical data.
// Each node of the tree must be identified by a Unique ID.
type Tree struct {
	BaseWidget
	Root     string
	Selected string
	Offset   fyne.Position

	open          map[string]bool
	branchMinSize fyne.Size
	leafMinSize   fyne.Size

	Children       func(uid string) (c []string)                         // Return a sorted slice of Children Unique IDs for the given Node Unique ID
	IsBranch       func(uid string) (ok bool)                            // Return true if the given Unique ID represents a Branch
	CreateNode     func(branch bool) (o fyne.CanvasObject)               // Return a CanvasObject that can represent a Branch (if branch is true), or a Leaf (if branch is false)
	UpdateNode     func(uid string, branch bool, node fyne.CanvasObject) // Update the given CanvasObject to represent the data at the given Unique ID
	OnBranchOpened func(uid string)                                      // Called when a Branch is opened
	OnBranchClosed func(uid string)                                      // Called when a Branch is closed
	OnNodeSelected func(uid string, node fyne.CanvasObject)              // Called when the Node with the given Unique ID is selected.
}

// NewTreeWithStrings creates a new tree with the given string map.
// Data must contain a mapping for the root, which defaults to empty string ("").
func NewTreeWithStrings(data map[string][]string) (t *Tree) {
	t = &Tree{
		Children: func(uid string) (c []string) {
			c = data[uid]
			return
		},
		IsBranch: func(uid string) (b bool) {
			_, b = data[uid]
			return
		},
		CreateNode: func(branch bool) fyne.CanvasObject {
			log.Println("Creating Node:", branch)
			return NewLabel("Template Object")
		},
		UpdateNode: func(uid string, branch bool, node fyne.CanvasObject) {
			node.(*Label).SetText(uid)
		},
	}
	t.ExtendBaseWidget(t)
	return
}

// NewTreeWithFiles creates a new tree with the given file system URI.
// TODO contents of a directory are sorted with the given sorter
func NewTreeWithFiles(root fyne.URI) (t *Tree) {
	t = &Tree{
		Root: root.String(),
		Children: func(uid string) (c []string) {
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
			return
		},
		IsBranch: func(uid string) bool {
			path := strings.TrimPrefix(uid, "file://")
			fi, err := os.Lstat(path)
			if err != nil {
				fyne.LogError("Unable to stat path "+path, err)
				return false
			}
			return fi.IsDir()
		},
		CreateNode: func(branch bool) fyne.CanvasObject {
			return NewHBox(NewIcon(theme.FileIcon()), NewLabel("Template Object"))
		},
	}
	t.UpdateNode = func(uid string, branch bool, node fyne.CanvasObject) {
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
	}
	t.ExtendBaseWidget(t)
	return
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
	t.ensureOpenMap()
	t.Walk(func(uid string, branch bool, depth int) {
		if branch {
			t.propertyLock.Lock()
			t.open[uid] = true
			t.propertyLock.Unlock()
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

// MinSize returns the size that this widget should not shrink below.
func (t *Tree) MinSize() fyne.Size {
	t.ExtendBaseWidget(t)
	return t.BaseWidget.MinSize()
}

// Walk visits every open node of the tree and calls the given callback with node Unique ID, whether node is branch, and the depth of node.
func (t *Tree) Walk(onNode func(string, bool, int)) {
	t.walk(t.Root, 0, onNode)
}

func (t *Tree) walk(uid string, depth int, onNode func(string, bool, int)) {
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

func (t *Tree) ensureOpenMap() {
	t.propertyLock.Lock()
	defer t.propertyLock.Unlock()
	if t.open == nil {
		t.open = make(map[string]bool)
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
	r.updateMinSizes()
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
	min = min.Max(r.tree.branchMinSize)
	min = min.Max(r.tree.leafMinSize)
	return
}

func (r *treeRenderer) Layout(size fyne.Size) {
	r.scroller.Resize(size)
}

func (r *treeRenderer) Refresh() {
	r.updateMinSizes()
	s := r.tree.Size()
	if s.IsZero() {
		r.tree.Resize(r.tree.MinSize())
	}
	r.content.Refresh()
	canvas.Refresh(r.tree.super())
}

func (r *treeRenderer) updateMinSizes() {
	r.tree.propertyLock.RLock()
	defer r.tree.propertyLock.RUnlock()
	if f := r.tree.CreateNode; f != nil {
		log.Println("UpdateMinSizes")
		r.tree.branchMinSize = newBranch(r.tree, f(true)).MinSize()
		r.tree.leafMinSize = newLeaf(r.tree, f(false)).MinSize()
	}
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
	dividers    []*canvas.Rectangle
	branches    map[string]*branch
	leaves      map[string]*leaf
	branchPool  pool
	leafPool    pool
}

func (r *treeContentRenderer) MinSize() (min fyne.Size) {
	r.treeContent.propertyLock.Lock()
	defer r.treeContent.propertyLock.Unlock()

	maxDepth := 0
	r.treeContent.tree.Walk(func(uid string, isBranch bool, depth int) {
		// Root node is not rendered unless it has been customized
		if r.treeContent.tree.Root == "" {
			depth = depth - 1
			if uid == "" {
				// This is root node, skip
				return
			}
		}
		maxDepth = fyne.Max(maxDepth, depth)

		// If this is not the first item, add a divider
		if min.Height > 0 {
			min.Height += treeDividerHeight
		}

		m := r.treeContent.tree.leafMinSize
		if isBranch {
			m = r.treeContent.tree.branchMinSize
		}
		min.Width = fyne.Max(min.Width, m.Width)
		min.Height += m.Height
	})
	min.Width += maxDepth * (theme.IconInlineSize() + theme.Padding())
	return
}

func (r *treeContentRenderer) Layout(size fyne.Size) {

	r.treeContent.propertyLock.Lock()
	defer r.treeContent.propertyLock.Unlock()

	branches := make(map[string]*branch)
	leaves := make(map[string]*leaf)

	offsetY := r.treeContent.tree.Offset.Y
	y := 0
	numDividers := 0
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
			s := fyne.NewSize(size.Width-2*theme.Padding(), treeDividerHeight)
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
		} else if y > offsetY+size.Height {
			// Node is below viewport and not visible
		} else {
			// Node is in viewport
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
				n.Move(fyne.NewPos(0, y))
				n.Resize(fyne.NewSize(size.Width, m.Height))
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
			log.Println("Branch Released")
			r.branchPool.Release(b)
		}
	}
	for uid, l := range r.leaves {
		if _, ok := leaves[uid]; !ok {
			l.Hide()
			log.Println("Leaf Released")
			r.leafPool.Release(l)
		}
	}

	r.branches = branches
	r.leaves = leaves
}

func (r *treeContentRenderer) Refresh() {
	s := r.treeContent.Size()
	if s.IsZero() {
		m := r.treeContent.MinSize().Max(r.treeContent.tree.Size())
		r.treeContent.Resize(m)
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

func (r *treeContentRenderer) getBranch() (b *branch) {
	o := r.branchPool.Obtain()
	if o != nil {
		log.Println("Branch Obtained")
		b = o.(*branch)
	} else {
		log.Println("Branch Created")
		b = newBranch(r.treeContent.tree, r.treeContent.tree.CreateNode(true))
	}
	return
}

func (r *treeContentRenderer) getLeaf() (l *leaf) {
	o := r.leafPool.Obtain()
	if o != nil {
		log.Println("Leaf Obtained")
		l = o.(*leaf)
	} else {
		log.Println("Leaf Created")
		l = newLeaf(r.treeContent.tree, r.treeContent.tree.CreateNode(false))
	}
	return
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
	temp := n.content
	n.content = nil
	if n.Visible() {
		n.Refresh()
	} else {
		n.Show()
	}
	n.content = temp
}

func (n *treeNode) Content() fyne.CanvasObject {
	return n.content
}

func (n *treeNode) Indent() int {
	return n.depth * (theme.IconInlineSize() + theme.Padding())
}

func (n *treeNode) Tapped(*fyne.PointEvent) {
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
		indicator:    canvas.NewRectangle(theme.BackgroundColor()),
	}
}

var _ fyne.WidgetRenderer = (*treeNodeRenderer)(nil)

type treeNodeRenderer struct {
	widget.BaseRenderer
	treeNode  *treeNode
	indicator *canvas.Rectangle
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
	min.Width += theme.Padding() + r.treeNode.Indent() + theme.IconInlineSize()
	min.Width += 2 * theme.Padding()
	min.Height = fyne.Max(min.Height, theme.IconInlineSize())
	min.Height += 2 * theme.Padding()
	return
}

func (r *treeNodeRenderer) Layout(size fyne.Size) {
	x := 0
	y := 0
	if i := r.indicator; i != nil {
		i.Move(fyne.NewPos(x, y))
		s := fyne.NewSize(theme.Padding(), size.Height)
		i.SetMinSize(s)
		i.Resize(s)
	}
	h := size.Height - 2*theme.Padding()
	x += theme.Padding() + r.treeNode.Indent()
	y += theme.Padding()
	if i := r.treeNode.icon; i != nil {
		i.Move(fyne.NewPos(x, y))
		i.Resize(fyne.NewSize(theme.IconInlineSize(), h))
	}
	x += theme.IconInlineSize()
	x += theme.Padding()
	if c := r.treeNode.content; c != nil {
		c.Move(fyne.NewPos(x, y))
		c.Resize(fyne.NewSize(size.Width-x-theme.Padding(), h))
	}
}

func (r *treeNodeRenderer) Refresh() {
	if i := r.treeNode.icon; i != nil {
		i.Refresh()
	}
	if i := r.indicator; i != nil {
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
	if i := r.indicator; i != nil {
		objects = append(objects, i)
	}
	if c := r.treeNode.content; c != nil {
		objects = append(objects, c)
	}
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
	i.tree.ToggleBranch(i.uid)
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
