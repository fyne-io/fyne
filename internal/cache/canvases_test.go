package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"fyne.io/fyne/v2"
)

func TestSetCanvasForObject(t *testing.T) {
	c := &dummyCanvas{}
	obj := &dummyWidget{}

	called := 0
	setup := func() {
		called++
	}

	SetCanvasForObject(obj, c, setup)
	assert.Equal(t, 1, called)
	SetCanvasForObject(obj, c, setup)
	assert.Equal(t, 1, called)

	// a different canvas (object moved window) - the stored canvas must follow,
	// an object can only be on one canvas at a time
	c = &dummyCanvas{}
	SetCanvasForObject(obj, c, setup)
	assert.Equal(t, 2, called)
	assert.Equal(t, c, GetCanvasForObject(obj))

	SetCanvasForObject(obj, c, setup)
	assert.Equal(t, 2, called)
}

func TestGetCanvasForObject_Missing(t *testing.T) {
	testClearAll()
	t.Cleanup(testClearAll)

	obj := &dummyWidget{}
	assert.Nil(t, GetCanvasForObject(obj))

	canvases.Store(obj, nil)
	assert.Nil(t, GetCanvasForObject(obj))
}

func TestSetCanvasForObject_keepsAlive(t *testing.T) {
	testClearAll()
	defer testClearAll()

	tm := &timeMock{}
	tm.setTime(10, 10)

	obj := &dummyWidget{}
	SetCanvasForObject(obj, &dummyCanvas{}, nil)

	// a window that keeps repainting must not expire its own cache entry
	tm.advance(2 * ValidDuration)
	SetCanvasForObject(obj, canvasForObject(obj), nil)

	destroyExpiredCanvases(tm.now)
	assert.Equal(t, 1, canvases.Len())
}

func canvasForObject(obj fyne.CanvasObject) fyne.Canvas {
	cinfo, _ := canvases.Load(obj)
	return cinfo.canvas
}

// setCanvasForObjectLoadOrStore is the implementation SetCanvasForObject replaced.
// It is kept here only so the benchmark below can show the difference.
func setCanvasForObjectLoadOrStore(obj fyne.CanvasObject, c fyne.Canvas, setup func()) {
	cinfo := &canvasInfo{canvas: c}
	cinfo.setAlive()

	old, found := canvases.LoadOrStore(obj, cinfo)
	if (!found || old.canvas != c) && setup != nil {
		setup()
	}
}

// BenchmarkSetCanvasForObject models a single repaint: Canvas.EnsureMinSize passes
// every visible object of a window through SetCanvasForObject on every dirty frame,
// so one iteration here is one frame of a window holding objectCount objects.
func BenchmarkSetCanvasForObject(b *testing.B) {
	const objectCount = 200

	run := func(b *testing.B, set func(fyne.CanvasObject, fyne.Canvas, func())) {
		testClearAll()
		defer testClearAll()

		c := &dummyCanvas{}
		objects := make([]fyne.CanvasObject, objectCount)
		for i := range objects {
			objects[i] = &dummyWidget{}
			set(objects[i], c, nil) // first frame attaches, the rest are re-visits
		}

		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			for _, obj := range objects {
				set(obj, c, nil)
			}
		}
	}

	b.Run("LoadOrStore", func(b *testing.B) { run(b, setCanvasForObjectLoadOrStore) })
	b.Run("Load", func(b *testing.B) { run(b, SetCanvasForObject) })
}

// BenchmarkSetCanvasForObjectRepaints models a window repainting for a second
// without interruption, long enough for the cache to expire underneath it. The
// setups/op metric counts how often an object had to be re-attached to its
// canvas - for a canvas.Image that is a full texture re-upload.
func BenchmarkSetCanvasForObjectRepaints(b *testing.B) {
	const (
		objectCount = 200
		frames      = 60
		frameTime   = time.Second / 60
	)

	run := func(b *testing.B, set func(fyne.CanvasObject, fyne.Canvas, func())) {
		testClearAll()
		validDuration := ValidDuration
		ValidDuration = time.Second / 2 // expire twice per benchmarked second
		defer func() {
			ValidDuration = validDuration
			testClearAll()
		}()

		tm := &timeMock{}
		tm.setTime(0, 0)

		c := &dummyCanvas{}
		objects := make([]fyne.CanvasObject, objectCount)
		for i := range objects {
			objects[i] = &dummyWidget{}
		}

		setups := 0
		countSetup := func() { setups++ }

		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			for f := 0; f < frames; f++ {
				for _, obj := range objects {
					set(obj, c, countSetup)
				}

				tm.now = tm.now.Add(frameTime)
				destroyExpiredCanvases(tm.now) // stands in for the driver's periodic Clean
			}
		}

		b.StopTimer()
		b.ReportMetric(float64(setups)/float64(b.N), "setups/op")
	}

	b.Run("LoadOrStore", func(b *testing.B) { run(b, setCanvasForObjectLoadOrStore) })
	b.Run("Load", func(b *testing.B) { run(b, SetCanvasForObject) })
}
