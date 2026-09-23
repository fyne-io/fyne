package fyne

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAccessibleRole_Values(t *testing.T) {
	// These values are part of Fyne's public API and assistive technologies
	// may map directly from them. Keep this table in sync with
	// accessibility.go and the platform-specific role mappers under
	// internal/driver/.
	cases := map[AccessibleRole]string{
		AccessibleRoleButton:      "button",
		AccessibleRoleCheckbox:    "checkbox",
		AccessibleRoleContainer:   "container",
		AccessibleRoleHeading:     "heading",
		AccessibleRoleImage:       "image",
		AccessibleRoleLink:        "link",
		AccessibleRoleList:        "list",
		AccessibleRoleListItem:    "listItem",
		AccessibleRoleProgressBar: "progressBar",
		AccessibleRoleRadio:       "radio",
		AccessibleRoleSeparator:   "separator",
		AccessibleRoleSlider:      "slider",
		AccessibleRoleTab:         "tab",
		AccessibleRoleTabList:     "tabList",
		AccessibleRoleTable:       "table",
		AccessibleRoleText:        "text",
		AccessibleRoleTextField:   "textField",
		AccessibleRoleTree:        "tree",
		AccessibleRoleTreeItem:    "treeItem",
	}
	for role, want := range cases {
		assert.Equal(t, want, string(role))
	}
}

func TestAccessibleAction_Values(t *testing.T) {
	cases := map[AccessibleAction]string{
		AccessibleActionDecrement: "decrement",
		AccessibleActionIncrement: "increment",
		AccessibleActionPress:     "press",
		AccessibleActionSelect:    "select",
		AccessibleActionSetValue:  "setValue",
		AccessibleActionShowMenu:  "showMenu",
	}
	for action, want := range cases {
		assert.Equal(t, want, string(action))
	}
}

func TestAccessibleState_Values(t *testing.T) {
	cases := map[AccessibleState]string{
		AccessibleStateChecked:  "checked",
		AccessibleStateDisabled: "disabled",
		AccessibleStateExpanded: "expanded",
		AccessibleStateFocused:  "focused",
		AccessibleStateInvalid:  "invalid",
		AccessibleStateRequired: "required",
		AccessibleStateSelected: "selected",
	}
	for state, want := range cases {
		assert.Equal(t, want, string(state))
	}
}

// stubAccessible exercises every optional accessibility interface so the
// compiler will report any drift in their signatures.
type stubAccessible struct {
	label     string
	role      AccessibleRole
	children  []CanvasObject
	value     string
	actions   []AccessibleAction
	states    []AccessibleState
	setOK     bool
	performed AccessibleAction
}

func (s *stubAccessible) AccessibilityLabel() string               { return s.label }
func (s *stubAccessible) AccessibilityRole() AccessibleRole        { return s.role }
func (s *stubAccessible) AccessibilityChildren() []CanvasObject    { return s.children }
func (s *stubAccessible) AccessibilityValue() string               { return s.value }
func (s *stubAccessible) AccessibilitySetValue(v string) bool      { s.value = v; return s.setOK }
func (s *stubAccessible) AccessibilityActions() []AccessibleAction { return s.actions }
func (s *stubAccessible) AccessibilityPerformAction(a AccessibleAction) bool {
	s.performed = a
	return true
}
func (s *stubAccessible) AccessibilityStates() []AccessibleState { return s.states }

func TestAccessible_OptionalInterfaces(t *testing.T) {
	var s stubAccessible
	assert.Implements(t, (*Accessible)(nil), &s)
	assert.Implements(t, (*AccessibleChildren)(nil), &s)
	assert.Implements(t, (*AccessibleValue)(nil), &s)
	assert.Implements(t, (*AccessibleValueSetter)(nil), &s)
	assert.Implements(t, (*AccessibleActions)(nil), &s)
	assert.Implements(t, (*AccessibleStates)(nil), &s)

	s.setOK = true
	assert.True(t, s.AccessibilitySetValue("hello"))
	assert.Equal(t, "hello", s.value)

	assert.True(t, s.AccessibilityPerformAction(AccessibleActionPress))
	assert.Equal(t, AccessibleActionPress, s.performed)
}
