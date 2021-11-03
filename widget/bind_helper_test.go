package widget

import (
	"sync/atomic"
	"testing"

	"fyne.io/fyne/v2/data/binding"
	"github.com/stretchr/testify/assert"
)

type atomicInt struct{ val int64 }

func (v *atomicInt) Set(val int) { atomic.StoreInt64(&v.val, int64(val)) }
func (v *atomicInt) Get() int    { return int(atomic.LoadInt64(&v.val)) }

func TestBasicBinder(t *testing.T) {
	var binder basicBinder
	value := &atomicInt{-1}

	item1 := binding.NewInt()

	callback := func(data binding.DataItem) {
		val, err := data.(binding.Int).Get()
		assert.NoError(t, err)
		value.Set(val)
	}
	binder.SetCallback(callback)

	assert.Equal(t, -1, value.Get())
	item1.Set(100)
	binder.Bind(item1)
	waitForBinding()
	assert.Equal(t, 100, value.Get()) // callback on the initial value
	item1.Set(150)
	waitForBinding()
	assert.Equal(t, 150, value.Get()) // callback on the changed value

	item2 := binding.NewInt()
	item2.Set(200)
	binder.Bind(item2)
	waitForBinding()
	assert.Equal(t, value.Get(), 200)
	item1.Set(300)
	waitForBinding()
	assert.Equal(t, value.Get(), 200)
	item2.Set(400)
	waitForBinding()
	assert.Equal(t, value.Get(), 400)
}
