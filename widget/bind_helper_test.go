package widget

import (
	"testing"

	"fyne.io/fyne/v2/data/binding"
	"github.com/stretchr/testify/assert"
)

func TestBasicBinder(t *testing.T) {
	var binder basicBinder
	var value int = -1

	item1 := binding.NewInt()

	callback := func(data binding.DataItem) {
		val, err := data.(binding.Int).Get()
		assert.NoError(t, err)
		value = val
	}
	binder.SetCallback(callback)

	assert.Equal(t, -1, value)
	item1.Set(100)
	binder.Bind(item1)
	waitForBinding()
	assert.Equal(t, 100, value) // callback on the initial value
	item1.Set(150)
	waitForBinding()
	assert.Equal(t, 150, value) // callback on the changed value

	item2 := binding.NewInt()
	item2.Set(200)
	binder.Bind(item2)
	waitForBinding()
	assert.Equal(t, 200, value)
	item1.Set(300)
	waitForBinding()
	assert.Equal(t, 200, value)
	item2.Set(400)
	waitForBinding()
	assert.Equal(t, 400, value)
}
