package widget

import (
	"image/color"
	"sort"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/driver/desktop"
	"fyne.io/fyne/theme"
)

var _ fyne.Widget = (*Tree)(nil)

// Tree widget displays hierarchical data.
type Tree struct {
	BaseWidget
	root *branch
	// GetBranchIcon is called to get the icon to display for the branch at the given path.
	GetBranchIcon func(path []string, open bool) fyne.Resource
	// GetLeafIcon is called to get the icon to display for the leaf at the given path.
	GetLeafIcon func(path []string) fyne.Resource
	// OnBranchChanged is called after a branch is opened or closed.
	OnBranchChanged func(path []string, open bool)
	// OnLeafSelected is called after a leaf is selected.
	OnLeafSelected func(path []string)
}

// Add the given path to the Tree.
func (t *Tree) Add(path ...string) {
	if t.root == nil {
		t.root = newBranch(t, nil, "")
	}
	t.root.add(path...)
	t.root.setOpen(true)
}

// CreateRenderer is a private method to Fyne which links this widget to its renderer.
func (t *Tree) CreateRenderer() fyne.WidgetRenderer {
	return &treeRenderer{
		tree: t,
	}
}

// UseArrowIcons configures the tree to display a right arrow for closed branches,
// a down arrow for open branches, and a circle for leaves.
func (t *Tree) UseArrowIcons() {
	t.GetBranchIcon = func(path []string, open bool) fyne.Resource {
		if open {
			return theme.MoveDownIcon()
		}
		return theme.NavigateNextIcon()
	}
	t.GetLeafIcon = func(path []string) fyne.Resource {
		return theme.RadioButtonCheckedIcon()
	}
}

// UseFileSystemIcons configures the tree to display a folder for closed branches,
// an open folder for open branches, and a file for leaves.
func (t *Tree) UseFileSystemIcons() {
	t.GetBranchIcon = func(path []string, open bool) fyne.Resource {
		if open {
			return theme.FolderOpenIcon()
		}
		return theme.FolderIcon()
	}
	t.GetLeafIcon = func(path []string) fyne.Resource {
		return theme.FileIcon()
	}
}

// NewTree creates a new empty tree.
func NewTree() *Tree {
	t := &Tree{}
	t.ExtendBaseWidget(t)
	return t
}

type treeRenderer struct {
	tree *Tree
}

func (r *treeRenderer) MinSize() fyne.Size {
	width := 0
	height := 0
	measure := func(depth int, node treeNode) {
		size := node.MinSize()
		sw := size.Width + theme.Padding()*depth
		if sw > width {
			width = sw
		}
		height += size.Height
	}
	if r.tree.root != nil {
		r.tree.root.walk(0, func(d int, b *branch) {
			measure(d, b)
		}, func(d int, l *leaf) {
			measure(d, l)
		})
	}
	return fyne.NewSize(width+theme.Padding()*2, height+theme.Padding()*2)
}

func (r *treeRenderer) Layout(size fyne.Size) {
	x := theme.Padding()
	y := theme.Padding()
	layout := func(depth int, node treeNode) {
		nx := x
		ny := y
		nw := size.Width - nx - theme.Padding()
		nh := node.MinSize().Height
		for n := node.getParent(); n != nil; n = n.getParent() {
			p := n.Position()
			ny -= p.Y
			nw -= p.X
		}
		node.Move(fyne.NewPos(nx, ny))
		node.Resize(fyne.NewSize(nw, nh))
		y += nh
	}
	if r.tree.root != nil {
		r.tree.root.walk(0, func(d int, b *branch) {
			layout(d, b)
		}, func(d int, l *leaf) {
			layout(d, l)
		})
	}
}

func (r *treeRenderer) BackgroundColor() color.Color {
	return theme.BackgroundColor()
}

func (r *treeRenderer) Refresh() {
	r.Layout(r.tree.Size())
	canvas.Refresh(r.tree)
}

func (r *treeRenderer) Objects() (objects []fyne.CanvasObject) {
	if r.tree.root != nil {
		objects = append(objects, r.tree.root)
	}
	return
}

func (r *treeRenderer) Destroy() {
}

type treeNode interface {
	fyne.Widget
	getTree() *Tree
	getParent() *branch
	getPath() []string
	getText() string
	getIcon() fyne.Resource
	isHovered() bool
	walk(int, func(int, *branch), func(int, *leaf))
}

type treeNodeRenderer struct {
	node  treeNode
	image *canvas.Image
	text  *canvas.Text
}

func (r *treeNodeRenderer) Layout(size fyne.Size) {
	if r.text.Text == "" {
		return
	}
	x := theme.Padding()
	y := theme.Padding()
	width := 0
	height := 0
	if r.image != nil {
		width += theme.IconInlineSize()
		height += theme.IconInlineSize()
		r.image.Move(fyne.NewPos(x, y))
		r.image.Resize(fyne.NewSize(width, height))
		x += width + theme.Padding()
	}
	r.text.Move(fyne.NewPos(x, y))
	r.text.Resize(fyne.NewSize(size.Width-x-theme.Padding(), r.text.MinSize().Height))
}

func (r *treeNodeRenderer) MinSize() fyne.Size {
	width := 0
	height := 0
	if r.text.Text != "" {
		if r.image != nil {
			width += theme.IconInlineSize() + theme.Padding()
			height += theme.IconInlineSize()
		}
		s := r.text.MinSize()
		width += s.Width
		if s.Height > height {
			height = s.Height
		}
		width += theme.Padding() * 2
		height += theme.Padding() * 2
	}
	return fyne.NewSize(width, height)
}

func (r *treeNodeRenderer) Refresh() {
	r.createCanvasObjects()
	r.Layout(r.node.Size())
	canvas.Refresh(r.node)
	r.node.getTree().Refresh()
}

func (r *treeNodeRenderer) createCanvasObjects() {
	if icon := r.node.getIcon(); icon != nil {
		r.image = canvas.NewImageFromResource(icon)
	}
	r.text = &canvas.Text{
		Text:      r.node.getText(),
		Color:     theme.TextColor(),
		TextSize:  theme.TextSize(),
		TextStyle: fyne.TextStyle{},
	}
}

func (r *treeNodeRenderer) BackgroundColor() color.Color {
	if r.node.isHovered() {
		return theme.HoverColor()
	}
	return theme.BackgroundColor()
}

func (r *treeNodeRenderer) Objects() []fyne.CanvasObject {
	objects := []fyne.CanvasObject{r.image, r.text}
	b, ok := r.node.(*branch)
	if ok && b.open {
		for _, c := range b.branches {
			objects = append(objects, c)
		}
		for _, c := range b.leaves {
			objects = append(objects, c)
		}
	}
	return objects
}

func (r *treeNodeRenderer) Destroy() {
}

var _ fyne.Tappable = (*branch)(nil)

type branch struct {
	BaseWidget
	tree        *Tree
	parent      *branch
	text        string
	branchNames []string
	leafNames   []string
	branches    map[string]*branch
	leaves      map[string]*leaf
	open        bool
	hovered     bool
}

func (b *branch) getTree() *Tree {
	return b.tree
}

func (b *branch) getParent() *branch {
	return b.parent
}

func (b *branch) getPath() (result []string) {
	if b.parent != nil {
		result = append(result, b.parent.getPath()...)
	}
	if b.text != "" {
		result = append(result, b.text)
	}
	return
}

func (b *branch) getText() string {
	return b.text
}

func (b *branch) getIcon() fyne.Resource {
	if g := b.tree.GetBranchIcon; g != nil {
		return g(b.getPath(), b.open)
	}
	return nil
}

func (b *branch) isHovered() bool {
	return b.hovered
}

func (b *branch) walk(depth int, onBranch func(int, *branch), onLeaf func(int, *leaf)) {
	onBranch(depth, b)
	if b.open {
		for _, n := range b.branchNames {
			b.branches[n].walk(depth+1, onBranch, onLeaf)
		}
		for _, n := range b.leafNames {
			b.leaves[n].walk(depth+1, onBranch, onLeaf)
		}
	}
}

func (b *branch) add(path ...string) {
	if len(path) == 0 {
		return
	}
	p := path[0]
	if len(path) > 1 {
		c, ok := b.branches[p]
		if !ok {
			c = newBranch(b.tree, b, p)
			b.branches[p] = c
			b.branchNames = append(b.branchNames, p)
			sort.Strings(b.branchNames)
		}
		c.add(path[1:]...)
	} else {
		_, ok := b.leaves[p]
		if !ok {
			b.leaves[p] = newLeaf(b.tree, b, p)
			b.leafNames = append(b.leafNames, p)
			sort.Strings(b.leafNames)
		}
	}
	b.update()
}

func (b *branch) update() {
	if b.open {
		for _, c := range b.branches {
			c.Show()
		}
		for _, c := range b.leaves {
			c.Show()
		}
	} else {
		for _, c := range b.branches {
			c.Hide()
		}
		for _, c := range b.leaves {
			c.Hide()
		}
	}
}

func (b *branch) setOpen(open bool) {
	b.open = open
	b.update()
	b.Refresh()
}

func (b *branch) Hide() {
	for _, c := range b.branches {
		c.Hide()
	}
	for _, c := range b.leaves {
		c.Hide()
	}
	b.BaseWidget.Hide()
}

func (b *branch) Show() {
	if b.open {
		for _, c := range b.branches {
			c.Show()
		}
		for _, c := range b.leaves {
			c.Show()
		}
	}
	b.BaseWidget.Show()
}

func (b *branch) Tapped(event *fyne.PointEvent) {
	path := b.getPath()
	b.setOpen(!b.open)
	if c := b.tree.OnBranchChanged; c != nil {
		c(path, b.open)
	}
}

func (b *branch) MouseIn(event *desktop.MouseEvent) {
	b.hovered = true
	b.Refresh()
}

func (b *branch) MouseOut() {
	b.hovered = false
	b.Refresh()
}

func (b *branch) MouseMoved(event *desktop.MouseEvent) {
}

func (b *branch) CreateRenderer() fyne.WidgetRenderer {
	r := &treeNodeRenderer{node: b}
	r.createCanvasObjects()
	return r
}

func newBranch(tree *Tree, parent *branch, text string) *branch {
	b := &branch{
		tree:     tree,
		parent:   parent,
		text:     text,
		branches: make(map[string]*branch),
		leaves:   make(map[string]*leaf),
	}
	b.ExtendBaseWidget(b)
	return b
}

var _ fyne.Tappable = (*leaf)(nil)

type leaf struct {
	BaseWidget
	tree    *Tree
	parent  *branch
	text    string
	hovered bool
}

func (l *leaf) getTree() *Tree {
	return l.tree
}

func (l *leaf) getParent() *branch {
	return l.parent
}

func (l *leaf) getPath() (result []string) {
	if l.parent != nil {
		result = append(result, l.parent.getPath()...)
	}
	result = append(result, l.text)
	return
}

func (l *leaf) getText() string {
	return l.text
}

func (l *leaf) getIcon() fyne.Resource {
	if g := l.tree.GetLeafIcon; g != nil {
		return g(l.getPath())
	}
	return nil
}

func (l *leaf) isHovered() bool {
	return l.hovered
}

func (l *leaf) walk(depth int, onBranch func(int, *branch), onLeaf func(int, *leaf)) {
	onLeaf(depth, l)
}

func (l *leaf) Tapped(event *fyne.PointEvent) {
	path := l.getPath()
	if s := l.tree.OnLeafSelected; s != nil {
		s(path)
	}
}

func (l *leaf) MouseIn(event *desktop.MouseEvent) {
	l.hovered = true
	l.Refresh()
}

func (l *leaf) MouseOut() {
	l.hovered = false
	l.Refresh()
}

func (l *leaf) MouseMoved(event *desktop.MouseEvent) {
}

func (l *leaf) CreateRenderer() fyne.WidgetRenderer {
	r := &treeNodeRenderer{node: l}
	r.createCanvasObjects()
	return r
}

func newLeaf(tree *Tree, parent *branch, text string) *leaf {
	l := &leaf{
		tree:   tree,
		parent: parent,
		text:   text,
	}
	l.ExtendBaseWidget(l)
	return l
}
