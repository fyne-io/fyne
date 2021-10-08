package widget

import (
	"time"

	"fyne.io/fyne/v2"
)

type entryUserAction struct {
	actionType entryActionType
	timestamp  time.Time
	state      entryHistoryState
}

type entryActionType int

const (
	entryActionTypedRune entryActionType = iota
	entryActionCut
	entryActionPaste
	entryActionErasing
	entryActionSetText
	entryActionOriginal
)

type entryHistoryState struct {
	content                 string
	cursorRow, cursorColumn int
	scrollOffset            fyne.Position
	contentScrollOffset     fyne.Position
}

// registerAction creates a new action of the specified type and stores
// the snapshot in the action log. It expects the caller to hold .propertyLock().
func (e *Entry) registerAction(actionType entryActionType) {
	if !e.HistoryEnabled {
		return
	}
	e.ensureHistoryDefined()

	action := entryUserAction{
		actionType: actionType,
		timestamp:  e.timestamper(),
		state:      e.historySnapshot(),
	}

	if e.redoOffset > 0 {
		actionIndex := len(e.actionLog) - e.redoOffset
		e.actionLog = e.actionLog[:actionIndex]
		e.redoOffset = 0
	}

	if e.shouldMergeAction(action) {
		e.actionLog = e.actionLog[:len(e.actionLog)-1] // shouldMergeAction guarantees len>0
	}
	e.actionLog = append(e.actionLog, action)
}

// shouldMergeAction is an internal test that checks if the suggested action
// should replace the last recorded action or just be added as a next action.
func (e *Entry) shouldMergeAction(action entryUserAction) bool {
	const historyMergeInterval time.Duration = 1000 * time.Millisecond

	if len(e.actionLog) == 0 {
		return false
	}

	lastAction := e.actionLog[len(e.actionLog)-1]
	if action.timestamp.After(lastAction.timestamp.Add(historyMergeInterval)) {
		return false
	}

	areBothOfType := func(actionType entryActionType) bool {
		return (action.actionType == actionType) && (lastAction.actionType == actionType)
	}
	shouldMergeTyped := areBothOfType(entryActionTypedRune)
	shouldMergeErased := areBothOfType(entryActionErasing)

	return (shouldMergeTyped || shouldMergeErased)
}

// IsUndoAvailable returns true if Undo() may be called.
func (e *Entry) IsUndoAvailable() bool {
	return e.HistoryEnabled && (len(e.actionLog)-e.redoOffset > 0)
}

// IsRedoAvailable returns true if Redo() may be called, i.e.,
// if some action has just been undone.
func (e *Entry) IsRedoAvailable() bool {
	return e.HistoryEnabled && (e.redoOffset > 0)
}

// Undo rolls back one user action at a time if history tracking is enabled.
func (e *Entry) Undo() {
	if !e.HistoryEnabled {
		return
	}
	e.ensureHistoryDefined()

	actionIndex := len(e.actionLog) - 2 - e.redoOffset
	newState := *e.historyOrigin
	if actionIndex >= 0 {
		newState = e.actionLog[actionIndex].state
	}
	if actionIndex >= -1 {
		e.redoOffset++
	}

	e.restoreHistorySnapshot(newState)
}

// Redo replicates the recently undone action.
func (e *Entry) Redo() {
	if !e.HistoryEnabled || (e.redoOffset == 0) {
		return
	}
	e.ensureHistoryDefined()

	actionIndex := len(e.actionLog) - e.redoOffset
	e.restoreHistorySnapshot(e.actionLog[actionIndex].state)
	e.redoOffset--
}

// historySnapshot returns the information sufficient to restore the entry state
// (content, cursor position, scroll position etc) after an undo or redo.
// It expects the caller to hold propertyLock().
func (e *Entry) historySnapshot() entryHistoryState {
	state := entryHistoryState{}
	state.content = e.textProvider().String()
	state.cursorRow = e.CursorRow
	state.cursorColumn = e.CursorColumn
	state.scrollOffset = e.scroll.Offset
	state.contentScrollOffset = e.content.scroll.Offset
	return state
}

func (e *Entry) restoreHistorySnapshot(state entryHistoryState) {
	e.updateText(state.content)
	e.selecting = false
	e.selectKeyDown = false

	e.CursorRow = state.cursorRow
	e.CursorColumn = state.cursorColumn
	e.updateCursor()

	e.scroll.Offset = state.scrollOffset
	e.content.scroll.Offset = state.scrollOffset
	e.content.scroll.Refresh()
	e.scroll.Refresh()

	e.Refresh()
}

// registerInitialState resets the action history of an entry.
func (e *Entry) registerInitialState() {
	snapshot := e.historySnapshot()
	e.historyOrigin = &snapshot
	e.actionLog = nil
	e.redoOffset = 0
}

// ensureHistoryDefined handles the situation when user just switched HistoryEnabled flag
// and various history-associated pieces of data may need to be initialized.
func (e *Entry) ensureHistoryDefined() {
	if e.historyOrigin == nil {
		e.registerInitialState()
	}
	if e.timestamper == nil {
		e.timestamper = time.Now
	}
}
