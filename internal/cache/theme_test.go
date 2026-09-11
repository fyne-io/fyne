package cache

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"fyne.io/fyne/v2"
)

func TestOverrideTheme(t *testing.T) {
	testClearAll()
	t.Cleanup(testClearAll)

	th := dummyTheme{}
	w := &dummyWidget{}
	child := &dummyCanvasObject{}
	Renderer(w).(*dummyWidgetRenderer).objects = []fyne.CanvasObject{child}

	OverrideTheme(w, th)
	assert.Equal(t, th, WidgetTheme(w))
	assert.NotEmpty(t, WidgetScopeID(w))
	assert.Equal(t, th, WidgetTheme(child))
	assert.Equal(t, WidgetScopeID(w), WidgetScopeID(child))
}

func TestOverrideTheme_ContainerAndPrimitive(t *testing.T) {
	testClearAll()
	t.Cleanup(testClearAll)

	th := dummyTheme{}
	child := &dummyCanvasObject{}
	nested := &dummyWidget{}
	c := &fyne.Container{Objects: []fyne.CanvasObject{child, nested}}

	OverrideTheme(c, th)
	assert.Equal(t, th, WidgetTheme(child))
	assert.Equal(t, th, WidgetTheme(nested))
	assert.Nil(t, WidgetTheme(c))

	primitive := &dummyCanvasObject{}
	OverrideTheme(primitive, th)
	assert.Equal(t, th, WidgetTheme(primitive))
}

func TestOverrideTheme_SkipsMobileScopeAndNilRenderer(t *testing.T) {
	testClearAll()
	t.Cleanup(testClearAll)

	th := dummyTheme{}
	mobile := &dummyMobileScope{}
	OverrideTheme(mobile, th)
	assert.Nil(t, WidgetTheme(mobile))
	assert.Empty(t, WidgetScopeID(mobile))

	w := &nilRendererWidget{}
	OverrideTheme(w, th)
	assert.Equal(t, th, WidgetTheme(w))
}

func TestOverrideThemeMatchingScope(t *testing.T) {
	testClearAll()
	t.Cleanup(testClearAll)

	th := dummyTheme{}
	parent := &dummyCanvasObject{}
	child := &dummyCanvasObject{}

	assert.False(t, OverrideThemeMatchingScope(child, parent))
	assert.Nil(t, WidgetTheme(child))

	OverrideTheme(parent, th)
	assert.True(t, OverrideThemeMatchingScope(child, parent))
	assert.Equal(t, th, WidgetTheme(child))
	assert.Equal(t, WidgetScopeID(parent), WidgetScopeID(child))
}

func TestWidgetThemeAndScopeIDWithoutOverride(t *testing.T) {
	obj := &dummyCanvasObject{}
	assert.Nil(t, WidgetTheme(obj))
	assert.Empty(t, WidgetScopeID(obj))
}
