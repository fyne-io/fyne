package internal

import "fyne.io/fyne/v2"

// ClipStack keeps track of the areas that should be clipped when drawing a canvas.
// If no clips are present then adding one will be added as-is.
// Subsequent items pushed will be completely within the previous clip.
type ClipStack struct {
	clips []*ClipItem
}

// Pop removes the current top clip and returns it.
func (c *ClipStack) Pop() *ClipItem {
	if len(c.clips) == 0 {
		return nil
	}

	ret := c.clips[len(c.clips)-1]
	c.clips = c.clips[:len(c.clips)-1]
	return ret
}

// Length  returns the number of items in this clip stack. 0 means no clip.
func (c *ClipStack) Length() int {
	return len(c.clips)
}

// Push a new clip onto this stack at position and size specified.
// The returned clip item is the result of calculating the intersection of the requested clip and it's parent.
func (c *ClipStack) Push(p fyne.Position, s fyne.Size) *ClipItem {
	outer := c.Top()
	inner := outer.Intersect(p, s)

	c.clips = append(c.clips, inner)
	return inner
}

// Top returns the current clip item - it will always be within the bounds of any parent clips.
func (c *ClipStack) Top() *ClipItem {
	if len(c.clips) == 0 {
		return nil
	}

	return c.clips[len(c.clips)-1]
}

// ClipItem represents a single clip in a clip stack, denoted by a size and position.
type ClipItem struct {
	pos  fyne.Position
	size fyne.Size
}

// Rect returns the position and size parameters of the clip.
func (i *ClipItem) Rect() (fyne.Position, fyne.Size) {
	return i.pos, i.size
}

// Intersect returns a new clip item that is the intersection of the requested parameters and this clip.
func (i *ClipItem) Intersect(p fyne.Position, s fyne.Size) *ClipItem {
	ret := &ClipItem{p, s}
	if i == nil {
		return ret
	}

	if ret.pos.X < i.pos.X {
		ret.pos.X = i.pos.X
		ret.size.Width -= i.pos.X - p.X
	}
	if ret.pos.Y < i.pos.Y {
		ret.pos.Y = i.pos.Y
		ret.size.Height -= i.pos.Y - p.Y
	}

	if p.X+s.Width > i.pos.X+i.size.Width {
		ret.size.Width = (i.pos.X + i.size.Width) - ret.pos.X
	}
	if p.Y+s.Height > i.pos.Y+i.size.Height {
		ret.size.Height = (i.pos.Y + i.size.Height) - ret.pos.Y
	}

	if ret.size.Width < 0 || ret.size.Height < 0 {
		ret.size = fyne.NewSize(0, 0)
		return ret
	}

	return ret
}
