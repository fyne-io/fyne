package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"fyne.io/fyne/v2"
)

func TestClipStack_Intersect(t *testing.T) {
	p1 := fyne.NewPos(5, 25)
	s1 := fyne.NewSize(100, 100)
	c := &ClipStack{
		clips: []*ClipItem{
			{p1, s1},
		},
	}

	p2 := fyne.NewPos(25, 0)
	s2 := fyne.NewSize(50, 50)
	i := c.Push(p2, s2)

	assert.Equal(t, fyne.NewPos(25, 25), i.pos)
	assert.Equal(t, fyne.NewSize(50, 25), i.size)
	assert.Equal(t, 2, len(c.clips))

	_ = c.Pop()
	p2 = fyne.NewPos(50, 50)
	s2 = fyne.NewSize(150, 50)
	i = c.Push(p2, s2)

	assert.Equal(t, fyne.NewPos(50, 50), i.pos)
	assert.Equal(t, fyne.NewSize(55, 50), i.size)
	assert.Equal(t, 2, len(c.clips))
}

func TestClipStack_Pop(t *testing.T) {
	p := fyne.NewPos(5, 5)
	s := fyne.NewSize(100, 100)
	c := &ClipStack{
		clips: []*ClipItem{
			{p, s},
		},
	}

	i := c.Pop()
	assert.Equal(t, p, i.pos)
	assert.Equal(t, s, i.size)
	assert.Equal(t, 0, len(c.clips))
}

func TestClipStack_Push(t *testing.T) {
	c := &ClipStack{}
	p := fyne.NewPos(5, 5)
	s := fyne.NewSize(100, 100)

	i := c.Push(p, s)
	assert.Equal(t, p, i.pos)
	assert.Equal(t, s, i.size)
	assert.Equal(t, 1, len(c.clips))
}
