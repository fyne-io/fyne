package widget

import (
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
	content string
	// .cursorTextPosition is stored to unite a sequence of several typed runes into a single
	// undoable action. It is not used if actionType != entryActionTypedRune.
	cursorTextPosition int
}

type entryUserAction struct {
	actionType entryActionType
	timestamp  time.Time
	state      entryHistoryState
}

// registerAction creates a new action of the specified type and stores
// the snapshot in the action log. It expects the called to hold .propertyLock().
func (e *Entry) registerAction(actionType entryActionType) {
	if !e.historyEnabled {
		return
	}

	action := entryUserAction{
		actionType: actionType,
		timestamp:  e.timestamper(),
		state:      e.historySnapshot(),
	}

	// TODO: a sequence of typed runes should be a single undoable action.
	if e.redoOffset > 0 {
		actionIndex := len(e.actionLog) - e.redoOffset
		e.redoOffset = 0
		e.actionLog[actionIndex] = action
		e.actionLog = e.actionLog[:actionIndex+1]
	} else {
		e.actionLog = append(e.actionLog, action)
	}
}

func (e *Entry) IsUndoAvailable() bool {
	return e.historyEnabled && (len(e.actionLog)-e.redoOffset > 0)
}

func (e *Entry) IsRedoAvailable() bool {
	return e.historyEnabled && (e.redoOffset > 0)
}

func (e *Entry) Undo() {
	if !e.historyEnabled {
		return
	}

	actionIndex := len(e.actionLog) - e.redoOffset - 1
	if actionIndex == -1 {
		return
	}

	e.restoreHistorySnapshot(e.actionLog[actionIndex].state)
	e.redoOffset++
}

// historySnapshot returns the information sufficient to restore the entry state
// (content, cursor position, scroll position etc) after an undo or redo.
// It expects the caller to hold propertyLock().
func (e *Entry) historySnapshot() entryHistoryState {
	state := entryHistoryState{}
	state.content = e.textProvider().String()
	state.cursorTextPosition = e.cursorTextPos()
	return state
}

func (e *Entry) restoreHistorySnapshot(state entryHistoryState) {
	e.updateText(state.content)
	e.updateCursor()
}

// registerInitialState resets the action history of an entry.
func (e *Entry) registerInitialState() {
	e.historyOrigin = e.historySnapshot()
	e.actionLog = nil
	e.redoOffset = 0
}
