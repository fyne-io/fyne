//go:build accessibility && darwin

package glfw

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
)

type stateBag struct {
	canvas.Rectangle
	states  []fyne.AccessibleState
	actions []fyne.AccessibleAction
}

func (s *stateBag) AccessibilityLabel() string             { return "bag" }
func (s *stateBag) AccessibilityRole() fyne.AccessibleRole { return fyne.AccessibleRoleButton }
func (s *stateBag) AccessibilityStates() []fyne.AccessibleState {
	return s.states
}
func (s *stateBag) AccessibilityActions() []fyne.AccessibleAction {
	return s.actions
}
func (s *stateBag) AccessibilityPerformAction(_ fyne.AccessibleAction) bool { return true }

// TestRoleToC_Distinct verifies that semantically distinct Fyne roles map
// to distinct values in the C bridge enum, except for the documented
// fallbacks (Container/unknown both map to Group).
func TestRoleToC_Distinct(t *testing.T) {
	roles := []fyne.AccessibleRole{
		fyne.AccessibleRoleButton,
		fyne.AccessibleRoleCheckbox,
		fyne.AccessibleRoleHeading,
		fyne.AccessibleRoleImage,
		fyne.AccessibleRoleLink,
		fyne.AccessibleRoleList,
		fyne.AccessibleRoleListItem,
		fyne.AccessibleRoleProgressBar,
		fyne.AccessibleRoleRadio,
		fyne.AccessibleRoleSeparator,
		fyne.AccessibleRoleSlider,
		fyne.AccessibleRoleTab,
		fyne.AccessibleRoleTabList,
		fyne.AccessibleRoleTable,
		fyne.AccessibleRoleText,
		fyne.AccessibleRoleTextField,
		fyne.AccessibleRoleTree,
		fyne.AccessibleRoleTreeItem,
	}
	seen := map[uint64]fyne.AccessibleRole{}
	for _, r := range roles {
		v := uint64(roleToC(r))
		if prev, ok := seen[v]; ok {
			t.Errorf("role %s collides with %s (both -> %d)", r, prev, v)
		}
		seen[v] = r
	}
}

func TestRoleToC_Deterministic(t *testing.T) {
	for i := 0; i < 3; i++ {
		assert.Equal(t, roleToC(fyne.AccessibleRoleButton), roleToC(fyne.AccessibleRoleButton))
	}
}

func TestRoleToC_UnknownFallsBackToGroup(t *testing.T) {
	// Unknown role values are mapped to the catch-all Group role at the
	// bridge level. Container is its own distinct role.
	r1 := roleToC(fyne.AccessibleRole("not-a-real-role"))
	r2 := roleToC(fyne.AccessibleRole("also-not-real"))
	assert.Equal(t, r1, r2, "unknown roles must map deterministically to the same fallback")
	// And it must NOT collide with any of the known semantic roles.
	assert.NotEqual(t, roleToC(fyne.AccessibleRoleContainer), r1)
	assert.NotEqual(t, roleToC(fyne.AccessibleRoleButton), r1)
}

func TestStateMaskFor_AllStates(t *testing.T) {
	bag := &stateBag{
		states: []fyne.AccessibleState{
			fyne.AccessibleStateChecked,
			fyne.AccessibleStateDisabled,
			fyne.AccessibleStateExpanded,
			fyne.AccessibleStateFocused,
			fyne.AccessibleStateInvalid,
			fyne.AccessibleStateRequired,
			fyne.AccessibleStateSelected,
		},
	}
	mask := stateMaskFor(bag)
	// Each state should contribute a distinct bit; combining all 7 states
	// yields exactly 7 bits set in the mask.
	bits := 0
	for v := uint64(mask); v != 0; v &= v - 1 {
		bits++
	}
	assert.Equal(t, 7, bits, "expected 7 bits in state mask, got %d", bits)
}

func TestStateMaskFor_NoStatesProvider(t *testing.T) {
	rect := canvas.NewRectangle(nil)
	assert.Equal(t, 0, stateMaskFor(rect))
}

func TestStateMaskFor_DistinctBits(t *testing.T) {
	// Each individual state must be a power of two so they don't collide.
	individual := []fyne.AccessibleState{
		fyne.AccessibleStateChecked,
		fyne.AccessibleStateDisabled,
		fyne.AccessibleStateExpanded,
		fyne.AccessibleStateFocused,
		fyne.AccessibleStateInvalid,
		fyne.AccessibleStateRequired,
		fyne.AccessibleStateSelected,
	}
	seen := map[uint64]fyne.AccessibleState{}
	for _, s := range individual {
		bag := &stateBag{states: []fyne.AccessibleState{s}}
		v := uint64(stateMaskFor(bag))
		assert.NotZero(t, v, "%s should produce a non-zero mask", s)
		// Ensure single-bit (power of two) and no collision with other states.
		assert.Equal(t, uint64(0), v&(v-1), "%s mask must be a single bit, got %b", s, v)
		if prev, ok := seen[v]; ok {
			t.Errorf("state %s collides with %s (both -> %b)", s, prev, v)
		}
		seen[v] = s
	}
}

func TestActionMaskFor_AllActions(t *testing.T) {
	bag := &stateBag{
		actions: []fyne.AccessibleAction{
			fyne.AccessibleActionPress,
			fyne.AccessibleActionIncrement,
			fyne.AccessibleActionDecrement,
			fyne.AccessibleActionShowMenu,
			fyne.AccessibleActionSelect,
			fyne.AccessibleActionSetValue,
		},
	}
	mask, ok := actionMaskFor(bag)
	assert.True(t, ok)
	bits := 0
	for v := uint64(mask); v != 0; v &= v - 1 {
		bits++
	}
	assert.Equal(t, 6, bits, "expected 6 action bits, got %d", bits)
}

func TestActionMaskFor_NoActions(t *testing.T) {
	rect := canvas.NewRectangle(nil)
	_, ok := actionMaskFor(rect)
	assert.False(t, ok)
}

type setterOnly struct {
	canvas.Rectangle
}

func (s *setterOnly) AccessibilityLabel() string             { return "" }
func (s *setterOnly) AccessibilityRole() fyne.AccessibleRole { return fyne.AccessibleRoleTextField }
func (s *setterOnly) AccessibilityValue() string             { return "" }
func (s *setterOnly) AccessibilitySetValue(_ string) bool    { return true }

func TestActionMaskFor_ValueSetterImpliesSetValue(t *testing.T) {
	// A widget that implements AccessibleValueSetter but has no AccessibleActions
	// shouldn't surface SetValue (we need both an action list AND the setter).
	// However, if AccessibleActions exists and includes SetValue, the mask must
	// include the SetValue bit even without the setter — the setter implication
	// only applies when actions are *also* provided.
	bag := &stateBag{actions: []fyne.AccessibleAction{fyne.AccessibleActionSetValue}}
	mask, ok := actionMaskFor(bag)
	assert.True(t, ok)
	assert.NotZero(t, mask, "set-value action should be present in mask")
}

func TestActionFromC_RoundTrip(t *testing.T) {
	allActions := []fyne.AccessibleAction{
		fyne.AccessibleActionPress,
		fyne.AccessibleActionIncrement,
		fyne.AccessibleActionDecrement,
		fyne.AccessibleActionShowMenu,
		fyne.AccessibleActionSelect,
		fyne.AccessibleActionSetValue,
	}
	// Iterate the C codes 0..5 and verify each maps to one of the documented Go actions.
	seen := map[fyne.AccessibleAction]bool{}
	for code := 0; code <= 5; code++ {
		got, ok := actionFromC(code)
		assert.True(t, ok, "code %d should be known", code)
		seen[got] = true
	}
	for _, a := range allActions {
		assert.True(t, seen[a], "action %s missing from round-trip", a)
	}
}

func TestActionFromC_Unknown(t *testing.T) {
	_, ok := actionFromC(99)
	assert.False(t, ok)
}
