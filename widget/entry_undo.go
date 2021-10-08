package widget

import (
	"log"
	"time"
)

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
}

type entryUserAction struct {
	actionType entryActionType
	timestamp  time.Time
	state      entryHistoryState
}

// registerAction creates a new action of the specified type and stores
// the snapshot in the action log. It expects the caller to hold .propertyLock().
func (e *Entry) registerAction(actionType entryActionType) {
	if !e.historyEnabled {
		return
	}

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
	return e.historyEnabled && (len(e.actionLog)-e.redoOffset > 0)
}

// IsRedoAvailable returns true if Redo() may be called, i.e.,
// if some action has just been undone.
func (e *Entry) IsRedoAvailable() bool {
	return e.historyEnabled && (e.redoOffset > 0)
}

func (e *Entry) Undo() {
	if !e.historyEnabled {
		return
	}

	actionIndex := len(e.actionLog) - 2 - e.redoOffset
	newState := e.historyOrigin
	if actionIndex >= 0 {
		newState = e.actionLog[actionIndex].state
	}
	if actionIndex >= -1 {
		e.redoOffset++
	}

	e.restoreHistorySnapshot(newState)
}

func (e *Entry) Redo() {
	if !e.historyEnabled || (e.redoOffset == 0) {
		return
	}

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
	return state
}

func (e *Entry) restoreHistorySnapshot(state entryHistoryState) {
	e.updateText(state.content)
	e.CursorRow = state.cursorRow
	e.CursorColumn = state.cursorColumn
	e.updateCursor()
}

// registerInitialState resets the action history of an entry.
func (e *Entry) registerInitialState() {
	e.historyOrigin = e.historySnapshot()
	e.actionLog = nil
	e.redoOffset = 0
}

func (e *Entry) DumpHistoryState() {
	if !e.historyEnabled {
		log.Printf("History is disabled, nothing to dump\n")
		return
	}
	log.Printf("Entry history action log with %d items\n", len(e.actionLog))
	for i := range e.actionLog {
		log.Printf("Action %d of type %d\n", i, int(e.actionLog[i].actionType))
		log.Printf("Performed at %s\n", e.actionLog[i].timestamp.String())
		log.Printf("Contents: %s\n", e.actionLog[i].state.content)
	}
}
