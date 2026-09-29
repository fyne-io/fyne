package widget

import (
	"reflect"
	"slices"
)

// Declare conformity with entryMergeableUndoAction interface
var _ entryMergeableUndoAction = (*richTextUndoAction)(nil)

// richTextUndoAction is a change to the content of a [RichTextEntry]. It holds the segments that
// were replaced, augmenting the text changes that the undo history of an [Entry].
type richTextUndoAction struct {
	change richTextChange

	// the rune offsets of the cursor on either side of this change
	cursorBefore, cursorAfter int

	// typed is the text change that this action was recorded for
	typed *entryModifyAction
	// base is the content that it was applied to - used for merging.
	base []RichTextSegment

	// merged is set when more typing has been added since the change was calculated.
	merged bool
	// transient marks a change that cannot be seen, such as to the style that will be typed in.
	transient bool
}

// Undo does not change the text, the segments are restored by the [RichTextEntry].
func (*richTextUndoAction) Undo(s string) string {
	return s
}

// Redo does not change the text, the segments are restored by the [RichTextEntry].
func (*richTextUndoAction) Redo(s string) string {
	return s
}

// TryMerge adds the next action to this one if it carries on the same run of typing.
func (a *richTextUndoAction) TryMerge(next entryMergeableUndoAction) bool {
	if a.typed == nil || !a.typed.TryMerge(next) {
		return false
	}

	a.merged = true
	return true
}

// close stops any more typing from being merged into this action.
func (a *richTextUndoAction) close() {
	a.typed, a.base = nil, nil
}

// richTextEdit describes what an edit changed in the content of a [RichTextEntry]
// that was not recorded by the [Entry] as a change to its text.
type richTextEdit int

const (
	// richTextEditText is an edit that changed no more than the text.
	richTextEditText richTextEdit = iota
	// richTextEditHidden is an edit that cannot be seen, such as choosing the style
	// that will be typed in or removing segments that were left with no text.
	richTextEditHidden
	// richTextEditStyled is an edit to the style or the structure of the content.
	richTextEditStyled
)

// richTextChange describes how a list of segments was changed, keeping only the
// segments that differ and not the content on either side of them.
type richTextChange struct {
	at            int // the index of the first segment that was changed
	before, after []RichTextSegment

	// head and tail count the runes that did not change at the start of the first
	// segment and the end of the last, which are left out of before and after.
	head, tail int
}

// diffSegments works out the change that turns one list of segments into another.
func diffSegments(before, after []RichTextSegment) richTextChange {
	start := 0
	for start < len(before) && start < len(after) && sameSegment(before[start], after[start]) {
		start++
	}

	endBefore, endAfter := len(before), len(after)
	for endBefore > start && endAfter > start && sameSegment(before[endBefore-1], after[endAfter-1]) {
		endBefore--
		endAfter--
	}

	c := richTextChange{
		at:     start,
		before: cloneSegments(before[start:endBefore]),
		after:  cloneSegments(after[start:endAfter]),
	}
	c.trim()
	return c
}

// trim takes the text that was not changed off the ends of the segments in this
// change, so that typing in a long run of text does not keep a copy of all of it.
func (c *richTextChange) trim() {
	if len(c.before) == 0 || len(c.after) == 0 {
		return
	}

	firstBefore, firstAfter := c.before[0], c.after[0]
	lastBefore, lastAfter := c.before[len(c.before)-1], c.after[len(c.after)-1]

	var startBefore, startAfter, endBefore, endAfter []rune
	if sameHolder(firstBefore, firstAfter) {
		startBefore, startAfter = []rune(firstBefore.Textual()), []rune(firstAfter.Textual())
		for c.head < len(startBefore) && c.head < len(startAfter) && startBefore[c.head] == startAfter[c.head] {
			c.head++
		}
	}
	if sameHolder(lastBefore, lastAfter) {
		endBefore, endAfter = []rune(lastBefore.Textual()), []rune(lastAfter.Textual())
		limit := min(len(endBefore), len(endAfter))
		if len(c.before) == 1 { // the head was counted in this segment as well
			limit = min(limit, len(endBefore)-c.head)
		}
		if len(c.after) == 1 {
			limit = min(limit, len(endAfter)-c.head)
		}
		for c.tail < limit && endBefore[len(endBefore)-1-c.tail] == endAfter[len(endAfter)-1-c.tail] {
			c.tail++
		}
	}

	trimHolder(firstBefore, c.head, 0)
	trimHolder(firstAfter, c.head, 0)
	trimHolder(lastBefore, 0, c.tail)
	trimHolder(lastAfter, 0, c.tail)
}

// trimHolder removes the given number of runes from the start and end of the text in a segment.
func trimHolder(seg RichTextSegment, head, tail int) {
	holder, ok := seg.(textHolder)
	if !ok || head+tail == 0 {
		return
	}

	runes := []rune(holder.content())
	holder.setContent(string(runes[head : len(runes)-tail]))
}

// apply returns the given content with this change made to it.
func (c *richTextChange) apply(state []RichTextSegment) []RichTextSegment {
	return c.replace(state, c.before, c.after)
}

// revert returns the given content with this change taken out of it.
func (c *richTextChange) revert(state []RichTextSegment) []RichTextSegment {
	return c.replace(state, c.after, c.before)
}

func (c *richTextChange) replace(state, from, to []RichTextSegment) []RichTextSegment {
	end := c.at + len(from)
	if end > len(state) {
		return state // this is not the content that the change was recorded against
	}

	insert := cloneSegments(to)
	if c.head > 0 {
		runes := []rune(state[c.at].Textual())
		if first, ok := insert[0].(textHolder); ok && c.head <= len(runes) {
			first.setContent(string(runes[:c.head]) + first.content())
		}
	}
	if c.tail > 0 {
		runes := []rune(state[end-1].Textual())
		if last, ok := insert[len(insert)-1].(textHolder); ok && c.tail <= len(runes) {
			last.setContent(last.content() + string(runes[len(runes)-c.tail:]))
		}
	}

	return slices.Concat(state[:c.at], insert, state[end:])
}

// cloneSegments copies segments so that later edits to the content do not change
// the copy. Segments that are not edited in place, such as images, are shared.
func cloneSegments(in []RichTextSegment) []RichTextSegment {
	out := make([]RichTextSegment, len(in))
	for i, seg := range in {
		switch t := seg.(type) {
		case *TextSegment:
			text := *t
			out[i] = &text
		case *CodeBlockSegment:
			out[i] = &CodeBlockSegment{Text: t.Text, quotingLevel: t.quotingLevel}
		case *ListSegment:
			list := *t
			list.Items = cloneSegments(t.Items)
			list.markers = nil
			out[i] = &list
		case *ParagraphSegment:
			out[i] = &ParagraphSegment{Texts: cloneSegments(t.Texts)}
		default:
			out[i] = seg
		}
	}
	return out
}

// sameSegments reports whether two lists of segments hold the same content and style.
func sameSegments(a, b []RichTextSegment) bool {
	return slices.EqualFunc(a, b, sameSegment)
}

func sameSegment(a, b RichTextSegment) bool {
	switch first := a.(type) {
	case *TextSegment:
		second, ok := b.(*TextSegment)
		return ok && first.Text == second.Text && first.Style == second.Style
	case *CodeBlockSegment:
		second, ok := b.(*CodeBlockSegment)
		return ok && first.Text == second.Text && first.quotingLevel == second.quotingLevel
	case *ListSegment:
		second, ok := b.(*ListSegment)
		return ok && first.Ordered == second.Ordered && first.startIndex == second.startIndex &&
			first.indentationLevel == second.indentationLevel && first.quotingLevel == second.quotingLevel &&
			sameSegments(first.Items, second.Items)
	case *ParagraphSegment:
		second, ok := b.(*ParagraphSegment)
		return ok && sameSegments(first.Texts, second.Texts)
	}

	kind := reflect.TypeOf(a)
	return kind == reflect.TypeOf(b) && kind.Comparable() && a == b
}

// sameHolder reports whether two segments hold text in the same way, so that they
// differ by no more than the text inside them.
func sameHolder(a, b RichTextSegment) bool {
	switch first := a.(type) {
	case *TextSegment:
		second, ok := b.(*TextSegment)
		return ok && first.Style == second.Style
	case *CodeBlockSegment:
		second, ok := b.(*CodeBlockSegment)
		return ok && first.quotingLevel == second.quotingLevel
	}
	return false
}

// beginUndo prepares the undo history for an edit, returning the rune offset of
// the cursor before it takes place to be passed to [RichTextEntry.commitUndo].
func (e *RichTextEntry) beginUndo() int {
	provider := e.richProvider()
	if !sameSegments(e.undoState, provider.Segments) {
		// the content was replaced from outside, so the history no longer leads to it
		e.resetUndo()
	}

	return e.CursorTextOffset()
}

// resetUndo forgets the undo history, which starts again from the current content.
func (e *RichTextEntry) resetUndo() {
	e.undoStack.Clear()
	e.undoState = cloneSegments(e.richProvider().Segments)
	e.undoEdit = richTextEditText
}

// noteEdit records that the edit taking place changed more than the text.
func (e *RichTextEntry) noteEdit(edit richTextEdit) {
	e.undoEdit = max(e.undoEdit, edit)
}

// commitUndo adds what changed during an edit to the undo history.
// The text changes that were recorded by the [Entry] are replaced by an action
// that holds the segments, including any styling that happened along with them.
func (e *RichTextEntry) commitUndo(cursor int) {
	stack := &e.undoStack
	var typed *entryModifyAction
	for stack.index > 0 {
		text, ok := stack.items[stack.index-1].(*entryModifyAction)
		if !ok {
			break
		}

		if typed == nil {
			typed = text // further typing carries on from the last of them
		}
		stack.index--
		stack.items = stack.items[:stack.index]
	}

	top := e.lastUndoAction()
	merged := top != nil && top.merged
	edit := e.undoEdit
	e.undoEdit = richTextEditText
	if !merged && typed == nil && edit == richTextEditText {
		return
	}

	state := cloneSegments(e.richProvider().Segments)
	action := top
	if !merged {
		if sameSegments(e.undoState, state) {
			return
		}
		if top != nil {
			top.close()
		}

		hidden := typed == nil && edit == richTextEditHidden
		action = &richTextUndoAction{base: e.undoState, cursorBefore: cursor, transient: hidden}
		stack.Add(action)
	}

	action.change = diffSegments(action.base, state)
	action.cursorAfter = e.CursorTextOffset()
	action.merged = false
	if typed != nil {
		action.typed = typed
	}
	if edit == richTextEditStyled { // the offsets in the text change no longer match the content
		action.close()
	}
	e.undoState = state
}

// lastUndoAction returns the action that undo would take out, if there is one.
func (e *RichTextEntry) lastUndoAction() *richTextUndoAction {
	if !e.undoStack.CanUndo() {
		return nil
	}

	action, _ := e.undoStack.items[e.undoStack.index-1].(*richTextUndoAction)
	return action
}

// nextUndoAction returns the action that redo would apply, if there is one.
func (e *RichTextEntry) nextUndoAction() *richTextUndoAction {
	if !e.undoStack.CanRedo() {
		return nil
	}

	action, _ := e.undoStack.items[e.undoStack.index].(*richTextUndoAction)
	return action
}

// Undo un-does the last modifying user-action, bringing back the style of the
// content that it changed along with the text.
//
// Since: 2.9
func (e *RichTextEntry) Undo() {
	e.commitUndo(e.beginUndo())

	var last *richTextUndoAction
	for action := e.lastUndoAction(); action != nil; action = e.lastUndoAction() {
		e.undoStack.Undo(e.Text)
		action.close()
		e.undoState = action.change.revert(e.undoState)

		last = action
		if !action.transient {
			break // a change that could be seen was taken out
		}
	}

	if last != nil {
		e.restoreUndo(last.cursorBefore)
	}
}

// Redo re-applies the last undone user-action.
//
// Since: 2.9
func (e *RichTextEntry) Redo() {
	e.commitUndo(e.beginUndo())

	var last *richTextUndoAction
	visible := false
	for action := e.nextUndoAction(); action != nil; action = e.nextUndoAction() {
		if visible && !action.transient {
			break // the next change that can be seen waits for another redo
		}

		e.undoStack.Redo(e.Text)
		e.undoState = action.change.apply(e.undoState)

		last = action
		visible = visible || !action.transient
	}

	if last != nil {
		e.restoreUndo(last.cursorAfter)
	}
}

// restoreUndo shows the content that the undo history has moved to.
func (e *RichTextEntry) restoreUndo(cursor int) {
	provider := e.richProvider()
	provider.Segments = cloneSegments(e.undoState)

	content := provider.String()
	e.updateText(content, false)
	e.setCursorOffset(cursor)
	e.sel.selectRow, e.sel.selectColumn = e.CursorRow, e.CursorColumn
	e.sel.selecting = false

	if e.OnChanged != nil {
		e.OnChanged(content)
	}
	e.validateWithoutRefresh()
	e.Refresh()
}
