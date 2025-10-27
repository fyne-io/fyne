package common

import (
	"fmt"
	"image/color"
	"testing"

	"github.com/stretchr/testify/assert"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/internal/driver"
	"fyne.io/fyne/v2/internal/painter/gl"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/theme"
)

func TestCanvas_walkTree(t *testing.T) {
	test.NewTempApp(t)

	leftObj1 := canvas.NewRectangle(color.Gray16{Y: 1})
	leftObj2 := canvas.NewRectangle(color.Gray16{Y: 2})
	leftCol := container.NewWithoutLayout(leftObj1, leftObj2)
	rightObj1 := canvas.NewRectangle(color.Gray16{Y: 10})
	rightObj2 := canvas.NewRectangle(color.Gray16{Y: 20})
	rightCol := container.NewWithoutLayout(rightObj1, rightObj2)
	content := container.NewWithoutLayout(leftCol, rightCol)
	content.Move(fyne.NewPos(17, 42))
	leftCol.Move(fyne.NewPos(300, 400))
	leftObj1.Move(fyne.NewPos(1, 2))
	leftObj2.Move(fyne.NewPos(20, 30))
	rightObj1.Move(fyne.NewPos(500, 600))
	rightObj2.Move(fyne.NewPos(60, 70))
	rightCol.Move(fyne.NewPos(7, 8))

	tree := &renderCacheTree{root: &RenderCacheNode{obj: content}}
	c := &Canvas{}
	c.Initialize(nil, func() {})
	c.SetContentTreeAndFocusMgr(&canvas.Rectangle{FillColor: theme.Color(theme.ColorNameBackground)})

	type nodeInfo struct {
		obj                                     fyne.CanvasObject
		lastBeforeCallIndex, lastAfterCallIndex int
	}

	painterData := make(map[*RenderCacheNode]nodeInfo)
	updateInfoBefore := func(node *RenderCacheNode, index int) {
		pd, ok := painterData[node]
		if ok && pd.obj != node.obj {
			panic("node cache does not match node obj - nodes should not be reused for different objects")
		}
		pd.obj = node.obj
		pd.lastBeforeCallIndex = index
		painterData[node] = pd
	}
	updateInfoAfter := func(node *RenderCacheNode, index int) {
		pd := painterData[node]
		if pd.obj != node.obj {
			panic("node cache does not match node obj - nodes should not be reused for different objects")
		}
		pd.lastAfterCallIndex = index
		painterData[node] = pd
	}

	//
	// test that first walk calls the hooks correctly
	//
	type beforeCall struct {
		obj    fyne.CanvasObject
		parent fyne.CanvasObject
		pos    fyne.Position
	}
	beforeCalls := []beforeCall{}
	type afterCall struct {
		obj    fyne.CanvasObject
		parent fyne.CanvasObject
	}
	afterCalls := []afterCall{}

	var i int
	c.walkTree(tree, func(node *RenderCacheNode, pos fyne.Position) {
		var parent fyne.CanvasObject
		if node.parent != nil {
			parent = node.parent.obj
		}
		i++
		updateInfoBefore(node, i)
		beforeCalls = append(beforeCalls, beforeCall{obj: node.obj, parent: parent, pos: pos})
	}, func(node *RenderCacheNode, _ fyne.Position) {
		var parent fyne.CanvasObject
		if node.parent != nil {
			parent = node.parent.obj
		}
		i++
		updateInfoAfter(node, i)
		node.minSize.Height = node.obj.Position().Y
		afterCalls = append(afterCalls, afterCall{obj: node.obj, parent: parent})
	})

	assert.Equal(t, []beforeCall{
		{obj: content, pos: fyne.NewPos(17, 42)},
		{obj: leftCol, parent: content, pos: fyne.NewPos(317, 442)},
		{obj: leftObj1, parent: leftCol, pos: fyne.NewPos(318, 444)},
		{obj: leftObj2, parent: leftCol, pos: fyne.NewPos(337, 472)},
		{obj: rightCol, parent: content, pos: fyne.NewPos(24, 50)},
		{obj: rightObj1, parent: rightCol, pos: fyne.NewPos(524, 650)},
		{obj: rightObj2, parent: rightCol, pos: fyne.NewPos(84, 120)},
	}, beforeCalls, "calls before children hook with the correct node and position")
	assert.Equal(t, []afterCall{
		{obj: leftObj1, parent: leftCol},
		{obj: leftObj2, parent: leftCol},
		{obj: leftCol, parent: content},
		{obj: rightObj1, parent: rightCol},
		{obj: rightObj2, parent: rightCol},
		{obj: rightCol, parent: content},
		{obj: content},
	}, afterCalls, "calls after children hook with the correct node")

	//
	// test that second walk gives access to the cache
	//
	secondRunBeforePainterData := []nodeInfo{}
	secondRunAfterPainterData := []nodeInfo{}
	nodes := []*RenderCacheNode{}

	c.walkTree(tree, func(node *RenderCacheNode, pos fyne.Position) {
		secondRunBeforePainterData = append(secondRunBeforePainterData, painterData[node])
		nodes = append(nodes, node)
	}, func(node *RenderCacheNode, _ fyne.Position) {
		secondRunAfterPainterData = append(secondRunAfterPainterData, painterData[node])
	})

	assert.Equal(t, []nodeInfo{
		{obj: content, lastBeforeCallIndex: 1, lastAfterCallIndex: 14},
		{obj: leftCol, lastBeforeCallIndex: 2, lastAfterCallIndex: 7},
		{obj: leftObj1, lastBeforeCallIndex: 3, lastAfterCallIndex: 4},
		{obj: leftObj2, lastBeforeCallIndex: 5, lastAfterCallIndex: 6},
		{obj: rightCol, lastBeforeCallIndex: 8, lastAfterCallIndex: 13},
		{obj: rightObj1, lastBeforeCallIndex: 9, lastAfterCallIndex: 10},
		{obj: rightObj2, lastBeforeCallIndex: 11, lastAfterCallIndex: 12},
	}, secondRunBeforePainterData, "second run uses cached nodes")
	assert.Equal(t, []nodeInfo{
		{obj: leftObj1, lastBeforeCallIndex: 3, lastAfterCallIndex: 4},
		{obj: leftObj2, lastBeforeCallIndex: 5, lastAfterCallIndex: 6},
		{obj: leftCol, lastBeforeCallIndex: 2, lastAfterCallIndex: 7},
		{obj: rightObj1, lastBeforeCallIndex: 9, lastAfterCallIndex: 10},
		{obj: rightObj2, lastBeforeCallIndex: 11, lastAfterCallIndex: 12},
		{obj: rightCol, lastBeforeCallIndex: 8, lastAfterCallIndex: 13},
		{obj: content, lastBeforeCallIndex: 1, lastAfterCallIndex: 14},
	}, secondRunAfterPainterData, "second run uses cached nodes")
	leftObj1Node := nodes[2]
	leftObj2Node := nodes[3]
	assert.Equal(t, leftObj2Node, leftObj1Node.nextSibling, "correct sibling relation")
	assert.Nil(t, leftObj2Node.nextSibling, "no surplus nodes")
	rightObj1Node := nodes[5]
	rightObj2Node := nodes[6]
	assert.Equal(t, rightObj2Node, rightObj1Node.nextSibling, "correct sibling relation")
	rightColNode := nodes[4]
	assert.Nil(t, rightColNode.nextSibling, "no surplus nodes")

	//
	// test that removal, replacement and adding at the end of a children list works
	//
	deleteAt(leftCol, 1)
	leftNewObj2 := canvas.NewRectangle(color.Gray16{Y: 3})
	leftCol.Add(leftNewObj2)
	deleteAt(rightCol, 1)
	thirdCol := container.NewVBox()
	content.Add(thirdCol)
	thirdRunBeforePainterData := []nodeInfo{}
	thirdRunAfterPainterData := []nodeInfo{}

	i = 0
	c.walkTree(tree, func(node *RenderCacheNode, pos fyne.Position) {
		i++
		updateInfoBefore(node, i)
		thirdRunBeforePainterData = append(thirdRunBeforePainterData, painterData[node])
	}, func(node *RenderCacheNode, _ fyne.Position) {
		i++
		updateInfoAfter(node, i)
		thirdRunAfterPainterData = append(thirdRunAfterPainterData, painterData[node])
	})

	assert.Equal(t, []nodeInfo{
		{obj: content, lastBeforeCallIndex: 1, lastAfterCallIndex: 14},
		{obj: leftCol, lastBeforeCallIndex: 2, lastAfterCallIndex: 7},
		{obj: leftObj1, lastBeforeCallIndex: 3, lastAfterCallIndex: 4},
		{obj: leftNewObj2, lastBeforeCallIndex: 5, lastAfterCallIndex: 0}, // new node for replaced obj
		{obj: rightCol, lastBeforeCallIndex: 8, lastAfterCallIndex: 13},
		{obj: rightObj1, lastBeforeCallIndex: 9, lastAfterCallIndex: 10},
		{obj: thirdCol, lastBeforeCallIndex: 12, lastAfterCallIndex: 0}, // new node for third column
	}, thirdRunBeforePainterData, "third run uses cached nodes if possible")
	assert.Equal(t, []nodeInfo{
		{obj: leftObj1, lastBeforeCallIndex: 3, lastAfterCallIndex: 4},
		{obj: leftNewObj2, lastBeforeCallIndex: 5, lastAfterCallIndex: 6}, // new node for replaced obj
		{obj: leftCol, lastBeforeCallIndex: 2, lastAfterCallIndex: 7},
		{obj: rightObj1, lastBeforeCallIndex: 9, lastAfterCallIndex: 10},
		{obj: rightCol, lastBeforeCallIndex: 8, lastAfterCallIndex: 11},
		{obj: thirdCol, lastBeforeCallIndex: 12, lastAfterCallIndex: 13}, // new node for third column
		{obj: content, lastBeforeCallIndex: 1, lastAfterCallIndex: 14},
	}, thirdRunAfterPainterData, "third run uses cached nodes if possible")
	assert.NotEqual(t, leftObj2Node, leftObj1Node.nextSibling, "new node for replaced object")
	assert.Nil(t, rightObj1Node.nextSibling, "node for removed object has been removed, too")
	assert.NotNil(t, rightColNode.nextSibling, "new node for new object")

	//
	// test that insertion at the beginning or in the middle of a children list
	// removes all following siblings and their subtrees
	//
	leftNewObj2a := canvas.NewRectangle(color.Gray16{Y: 4})
	insert(leftCol, leftNewObj2a, 1)
	rightNewObj0 := canvas.NewRectangle(color.Gray16{Y: 30})
	Prepend(rightCol, rightNewObj0)
	fourthRunBeforePainterData := []nodeInfo{}
	fourthRunAfterPainterData := []nodeInfo{}
	nodes = []*RenderCacheNode{}

	i = 0
	c.walkTree(tree, func(node *RenderCacheNode, pos fyne.Position) {
		i++
		updateInfoBefore(node, i)
		fourthRunBeforePainterData = append(fourthRunBeforePainterData, painterData[node])
		nodes = append(nodes, node)
	}, func(node *RenderCacheNode, _ fyne.Position) {
		i++
		updateInfoAfter(node, i)
		fourthRunAfterPainterData = append(fourthRunAfterPainterData, painterData[node])
	})

	assert.Equal(t, []nodeInfo{
		{obj: content, lastBeforeCallIndex: 1, lastAfterCallIndex: 14},
		{obj: leftCol, lastBeforeCallIndex: 2, lastAfterCallIndex: 7},
		{obj: leftObj1, lastBeforeCallIndex: 3, lastAfterCallIndex: 4},
		{obj: leftNewObj2a, lastBeforeCallIndex: 5, lastAfterCallIndex: 0}, // new node for inserted obj
		{obj: leftNewObj2, lastBeforeCallIndex: 7, lastAfterCallIndex: 0},  // new node because of tail cut
		{obj: rightCol, lastBeforeCallIndex: 10, lastAfterCallIndex: 11},
		{obj: rightNewObj0, lastBeforeCallIndex: 11, lastAfterCallIndex: 0}, // new node for inserted obj
		{obj: rightObj1, lastBeforeCallIndex: 13, lastAfterCallIndex: 0},    // new node because of tail cut
		{obj: thirdCol, lastBeforeCallIndex: 16, lastAfterCallIndex: 13},
	}, fourthRunBeforePainterData, "fourth run uses cached nodes if possible")
	assert.Equal(t, []nodeInfo{
		{obj: leftObj1, lastBeforeCallIndex: 3, lastAfterCallIndex: 4},
		{obj: leftNewObj2a, lastBeforeCallIndex: 5, lastAfterCallIndex: 6},
		{obj: leftNewObj2, lastBeforeCallIndex: 7, lastAfterCallIndex: 8},
		{obj: leftCol, lastBeforeCallIndex: 2, lastAfterCallIndex: 9},
		{obj: rightNewObj0, lastBeforeCallIndex: 11, lastAfterCallIndex: 12},
		{obj: rightObj1, lastBeforeCallIndex: 13, lastAfterCallIndex: 14},
		{obj: rightCol, lastBeforeCallIndex: 10, lastAfterCallIndex: 15},
		{obj: thirdCol, lastBeforeCallIndex: 16, lastAfterCallIndex: 17},
		{obj: content, lastBeforeCallIndex: 1, lastAfterCallIndex: 18},
	}, fourthRunAfterPainterData, "fourth run uses cached nodes if possible")
	// check cache tree integrity
	// content node
	assert.Equal(t, content, nodes[0].obj)
	assert.Equal(t, leftCol, nodes[0].firstChild.obj)
	assert.Nil(t, nodes[0].nextSibling)
	// leftCol node
	assert.Equal(t, leftCol, nodes[1].obj)
	assert.Equal(t, leftObj1, nodes[1].firstChild.obj)
	assert.Equal(t, rightCol, nodes[1].nextSibling.obj)
	// leftObj1 node
	assert.Equal(t, leftObj1, nodes[2].obj)
	assert.Nil(t, nodes[2].firstChild)
	assert.Equal(t, leftNewObj2a, nodes[2].nextSibling.obj)
	// leftNewObj2a node
	assert.Equal(t, leftNewObj2a, nodes[3].obj)
	assert.Nil(t, nodes[3].firstChild)
	assert.Equal(t, leftNewObj2, nodes[3].nextSibling.obj)
	// leftNewObj2 node
	assert.Equal(t, leftNewObj2, nodes[4].obj)
	assert.Nil(t, nodes[4].firstChild)
	assert.Nil(t, nodes[4].nextSibling)
	// rightCol node
	assert.Equal(t, rightCol, nodes[5].obj)
	assert.Equal(t, rightNewObj0, nodes[5].firstChild.obj)
	assert.Equal(t, thirdCol, nodes[5].nextSibling.obj)
	// rightNewObj0 node
	assert.Equal(t, rightNewObj0, nodes[6].obj)
	assert.Nil(t, nodes[6].firstChild)
	assert.Equal(t, rightObj1, nodes[6].nextSibling.obj)
	// rightObj1 node
	assert.Equal(t, rightObj1, nodes[7].obj)
	assert.Nil(t, nodes[7].firstChild)
	assert.Nil(t, nodes[7].nextSibling)
	// thirdCol node
	assert.Equal(t, thirdCol, nodes[8].obj)
	assert.Nil(t, nodes[8].firstChild)
	assert.Nil(t, nodes[8].nextSibling)

	//
	// test that removal at the beginning or in the middle of a children list
	// removes all following siblings and their subtrees
	//
	deleteAt(leftCol, 1)
	deleteAt(rightCol, 0)
	fifthRunBeforePainterData := []nodeInfo{}
	fifthRunAfterPainterData := []nodeInfo{}
	nodes = []*RenderCacheNode{}

	i = 0
	c.walkTree(tree, func(node *RenderCacheNode, pos fyne.Position) {
		i++
		updateInfoBefore(node, i)
		fifthRunBeforePainterData = append(fifthRunBeforePainterData, painterData[node])
		nodes = append(nodes, node)
	}, func(node *RenderCacheNode, _ fyne.Position) {
		i++
		updateInfoAfter(node, i)
		fifthRunAfterPainterData = append(fifthRunAfterPainterData, painterData[node])
	})

	assert.Equal(t, []nodeInfo{
		{obj: content, lastBeforeCallIndex: 1, lastAfterCallIndex: 18},
		{obj: leftCol, lastBeforeCallIndex: 2, lastAfterCallIndex: 9},
		{obj: leftObj1, lastBeforeCallIndex: 3, lastAfterCallIndex: 4},
		{obj: leftNewObj2, lastBeforeCallIndex: 5, lastAfterCallIndex: 0}, // new node because of tail cut
		{obj: rightCol, lastBeforeCallIndex: 8, lastAfterCallIndex: 15},
		{obj: rightObj1, lastBeforeCallIndex: 9, lastAfterCallIndex: 0}, // new node because of tail cut
		{obj: thirdCol, lastBeforeCallIndex: 12, lastAfterCallIndex: 17},
	}, fifthRunBeforePainterData, "fifth run uses cached nodes if possible")
	assert.Equal(t, []nodeInfo{
		{obj: leftObj1, lastBeforeCallIndex: 3, lastAfterCallIndex: 4},
		{obj: leftNewObj2, lastBeforeCallIndex: 5, lastAfterCallIndex: 6},
		{obj: leftCol, lastBeforeCallIndex: 2, lastAfterCallIndex: 7},
		{obj: rightObj1, lastBeforeCallIndex: 9, lastAfterCallIndex: 10},
		{obj: rightCol, lastBeforeCallIndex: 8, lastAfterCallIndex: 11},
		{obj: thirdCol, lastBeforeCallIndex: 12, lastAfterCallIndex: 13},
		{obj: content, lastBeforeCallIndex: 1, lastAfterCallIndex: 14},
	}, fifthRunAfterPainterData, "fifth run uses cached nodes if possible")
	// check cache tree integrity
	// content node
	assert.Equal(t, content, nodes[0].obj)
	assert.Equal(t, leftCol, nodes[0].firstChild.obj)
	assert.Nil(t, nodes[0].nextSibling)
	// leftCol node
	assert.Equal(t, leftCol, nodes[1].obj)
	assert.Equal(t, leftObj1, nodes[1].firstChild.obj)
	assert.Equal(t, rightCol, nodes[1].nextSibling.obj)
	// leftObj1 node
	assert.Equal(t, leftObj1, nodes[2].obj)
	assert.Nil(t, nodes[2].firstChild)
	assert.Equal(t, leftNewObj2, nodes[2].nextSibling.obj)
	// leftNewObj2 node
	assert.Equal(t, leftNewObj2, nodes[3].obj)
	assert.Nil(t, nodes[3].firstChild)
	assert.Nil(t, nodes[3].nextSibling)
	// rightCol node
	assert.Equal(t, rightCol, nodes[4].obj)
	assert.Equal(t, rightObj1, nodes[4].firstChild.obj)
	assert.Equal(t, thirdCol, nodes[4].nextSibling.obj)
	// rightObj1 node
	assert.Equal(t, rightObj1, nodes[5].obj)
	assert.Nil(t, nodes[5].firstChild)
	assert.Nil(t, nodes[5].nextSibling)
	// thirdCol node
	assert.Equal(t, thirdCol, nodes[6].obj)
	assert.Nil(t, nodes[6].firstChild)
	assert.Nil(t, nodes[6].nextSibling)
}

func TestCanvas_OverlayStack(t *testing.T) {
	o := &overlayStack{}
	a := canvas.NewRectangle(color.Black)
	b := canvas.NewCircle(color.Black)
	c := canvas.NewRectangle(color.White)
	o.Add(a)
	o.Add(b)
	o.Add(c)
	assert.Equal(t, 3, len(o.List()))
	o.Remove(c)
	assert.Equal(t, 2, len(o.List()))
	o.Remove(a)
	assert.Equal(t, 0, len(o.List()))
}

func deleteAt(c *fyne.Container, index int) {
	if index < len(c.Objects)-1 {
		c.Objects = append(c.Objects[:index], c.Objects[index+1:]...)
	} else {
		c.Objects = c.Objects[:index]
	}
	c.Refresh()
}

func insert(c *fyne.Container, object fyne.CanvasObject, index int) {
	tail := append([]fyne.CanvasObject{object}, c.Objects[index:]...)
	c.Objects = append(c.Objects[:index], tail...)
	c.Refresh()
}

func Prepend(c *fyne.Container, object fyne.CanvasObject) {
	c.Objects = append([]fyne.CanvasObject{object}, c.Objects...)
	c.Refresh()
}

type dummyCanvas struct {
	fyne.Canvas
}

func (dummyCanvas) Scale() float32 {
	return 1.0
}

func TestRefreshCount(t *testing.T) { // Issue 2548.
	c := &Canvas{}
	c.Initialize(nil, func() {})

	freed := c.FreeDirtyTextures()
	assert.Zero(t, freed)

	c.SetPainter(gl.NewPainter(&dummyCanvas{}, struct{ driver.WithContext }{}))

	const refresh = uint64(1000)
	for i := uint64(0); i < refresh; i++ {
		c.Refresh(canvas.NewRectangle(color.Gray16{Y: 1}))
	}

	freed = c.FreeDirtyTextures()
	assert.Equal(t, refresh, freed)
}

func BenchmarkRefresh(b *testing.B) {
	c := &Canvas{}
	c.Initialize(nil, func() {})
	c.SetPainter(gl.NewPainter(&dummyCanvas{}, struct{ driver.WithContext }{}))

	for i := uint64(1); i < 1<<15; i *= 2 {
		b.Run(fmt.Sprintf("#%d", i), func(b *testing.B) {
			b.ReportAllocs()

			for j := 0; j < b.N; j++ {
				for n := uint64(0); n < i; n++ {
					c.Refresh(canvas.NewRectangle(color.Black))
				}
				c.FreeDirtyTextures()
			}
		})
	}
}
