package fyne

// AccessibleRole describes the different roles an accessible element can take.
//
// Since: 2.8
type AccessibleRole string

const (
	AccessibleRoleButton      AccessibleRole = "button"
	AccessibleRoleCheckbox    AccessibleRole = "checkbox"
	AccessibleRoleContainer   AccessibleRole = "container"
	AccessibleRoleHeading     AccessibleRole = "heading"
	AccessibleRoleImage       AccessibleRole = "image"
	AccessibleRoleLink        AccessibleRole = "link"
	AccessibleRoleList        AccessibleRole = "list"
	AccessibleRoleListItem    AccessibleRole = "listItem"
	AccessibleRoleProgressBar AccessibleRole = "progressBar"
	AccessibleRoleRadio       AccessibleRole = "radio"
	AccessibleRoleSeparator   AccessibleRole = "separator"
	AccessibleRoleSlider      AccessibleRole = "slider"
	AccessibleRoleTab         AccessibleRole = "tab"
	AccessibleRoleTabList     AccessibleRole = "tabList"
	AccessibleRoleTable       AccessibleRole = "table"
	AccessibleRoleText        AccessibleRole = "text"
	AccessibleRoleTextField   AccessibleRole = "textField"
	AccessibleRoleTree        AccessibleRole = "tree"
	AccessibleRoleTreeItem    AccessibleRole = "treeItem"
)

// AccessibleAction names a user-triggered action that an [Accessible] widget
// can support. Action names are passed to [AccessibleActions.AccessibilityPerformAction].
//
// Since: 2.8
type AccessibleAction string

const (
	AccessibleActionDecrement AccessibleAction = "decrement"
	AccessibleActionIncrement AccessibleAction = "increment"
	AccessibleActionPress     AccessibleAction = "press"
	AccessibleActionSelect    AccessibleAction = "select"
	AccessibleActionSetValue  AccessibleAction = "setValue"
	AccessibleActionShowMenu  AccessibleAction = "showMenu"
)

// AccessibleState reports an additional flag about an [Accessible] widget
// that an assistive technology may convey to the user (for example,
// "checked" or "disabled").
//
// Since: 2.8
type AccessibleState string

const (
	AccessibleStateChecked  AccessibleState = "checked"
	AccessibleStateDisabled AccessibleState = "disabled"
	AccessibleStateExpanded AccessibleState = "expanded"
	AccessibleStateFocused  AccessibleState = "focused"
	AccessibleStateInvalid  AccessibleState = "invalid"
	AccessibleStateRequired AccessibleState = "required"
	AccessibleStateSelected AccessibleState = "selected"
)

// Accessible interface should be implemented for a widget that should be accessible
//
// Since: 2.8
type Accessible interface {
	AccessibilityLabel() string
	AccessibilityRole() AccessibleRole
}

// AccessibleChildren may be implemented by an [Accessible] widget that
// contains accessible descendants outside the standard [Container] tree
// (for example, the active tab content of a tabbed container or the
// visible rows of a tree).
//
// Implementations should return only the descendants that should appear in
// the accessibility tree at this point. Returning nil is equivalent to
// having no accessible children. Objects that do not implement
// AccessibleChildren may still be traversed through any [Container.Objects].
//
// Since: 2.8
type AccessibleChildren interface {
	AccessibilityChildren() []CanvasObject
}

// AccessibleValue may be implemented by an [Accessible] widget that exposes
// a textual value distinct from its label (for example the typed text of an
// entry, the current value of a slider, or the percentage of a progress bar).
//
// Since: 2.8
type AccessibleValue interface {
	AccessibilityValue() string
}

// AccessibleValueSetter may be implemented by an [Accessible] widget whose
// value can be replaced by an assistive technology (for example, voice
// dictation typing into a text entry).
//
// Implementations should return true if the value was applied.
//
// Since: 2.8
type AccessibleValueSetter interface {
	AccessibilitySetValue(value string) bool
}

// AccessibleActions may be implemented by an [Accessible] widget that
// supports one or more user-triggered actions (clicks, increments, menu
// invocations).
//
// Implementations report the supported actions via AccessibilityActions
// and execute them via AccessibilityPerformAction. The latter should
// return true if the action was handled.
//
// Since: 2.8
type AccessibleActions interface {
	AccessibilityActions() []AccessibleAction
	AccessibilityPerformAction(action AccessibleAction) bool
}

// AccessibleStates may be implemented by an [Accessible] widget to expose
// additional flags about its current state (for example, "checked" for a
// checkbox or "expanded" for a disclosure widget).
//
// Since: 2.8
type AccessibleStates interface {
	AccessibilityStates() []AccessibleState
}
