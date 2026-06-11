package cache

import (
	"fmt"
	"os"
	"testing"
	"time"

	"fyne.io/fyne/v2"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	ret := m.Run()
	testClearAll()
	os.Exit(ret)
}

func TestCacheClean(t *testing.T) {
	tm := &timeMock{}

	t.Run("frame_counter_keeps_objects_alive", func(t *testing.T) {
		testClearAll()
		destroyedRenderersCnt := 0
		tm.setTime(10, 10)
		for i := 0; i < 20; i++ {
			Renderer(&dummyWidget{onDestroy: func() { destroyedRenderersCnt++ }})
			SetCanvasForObject(&dummyWidget{}, &dummyCanvas{}, nil)
		}

		// Force full clean without advancing the frame counter
		lastClean = tm.createTime(10, 10)
		tm.setTime(12, 11) // ValidDuration (2 min) elapsed
		Clean()

		// Frame-based caches kept alive because framecounter hasn't advanced
		assert.Equal(t, 20, renderers.Len())
		assert.Equal(t, 20, canvases.Len())
		assert.Zero(t, destroyedRenderersCnt)
	})

	t.Run("advancing_frame_counter_expires_objects", func(t *testing.T) {
		testClearAll()
		destroyedRenderersCnt := 0
		tm.setTime(10, 10)
		for i := 0; i < 20; i++ {
			Renderer(&dummyWidget{onDestroy: func() { destroyedRenderersCnt++ }})
			SetCanvasForObject(&dummyWidget{}, &dummyCanvas{}, nil)
		}

		// Advance frame counter — previous objects are now expired
		IncrementFrameCounter()

		// Force full clean
		lastClean = tm.createTime(10, 10)
		tm.setTime(12, 11) // ValidDuration (2 min) elapsed
		Clean()

		assert.Equal(t, 0, renderers.Len())
		assert.Equal(t, 0, canvases.Len())
		assert.Equal(t, 20, destroyedRenderersCnt)
	})

	t.Run("full_clean_expires_svgs_after_valid_duration", func(t *testing.T) {
		testClearAll()
		tm.setTime(10, 10)
		for i := 0; i < 20; i++ {
			SetSvg(fmt.Sprintf("%d", i), nil, nil, i, i+1)
		}

		// No full clean yet — ValidDuration hasn't elapsed
		lastClean = tm.createTime(10, 10)
		tm.setTime(10, 20)
		Clean()
		assert.Equal(t, 20, svgs.Len())

		// After ValidDuration, full clean removes expired SVGs
		tm.setTime(12, 11)
		Clean()
		assert.Equal(t, 0, svgs.Len())
	})

	t.Run("dynamic_flags_trigger_targeted_clean", func(t *testing.T) {
		testClearAll()
		destroyedRenderersCnt := 0
		tm.setTime(10, 10)

		// Set a non-zero baseline clean size
		rendererCacheLastCleanSize = 5
		for i := 0; i < 11; i++ {
			Renderer(&dummyWidget{onDestroy: func() { destroyedRenderersCnt++ }})
		}
		assert.True(t, shouldCleanRenderers) // 11 > 2*5

		// Targeted clean without full clean (ValidDuration not elapsed)
		lastClean = tm.createTime(10, 10)
		tm.setTime(10, 20)
		Clean()

		// Flag cleared and last-clean size updated
		assert.False(t, shouldCleanRenderers)
		assert.Equal(t, 11, rendererCacheLastCleanSize)
		// Objects not expired (frame counter not advanced)
		assert.Equal(t, 11, renderers.Len())
		assert.Zero(t, destroyedRenderersCnt)
	})
}

func TestCleanCanvas(t *testing.T) {
	destroyedRenderersCnt := 0
	testClearAll()

	dcanvas1 := &dummyCanvas{}
	dcanvas2 := &dummyCanvas{}

	for i := 0; i < 20; i++ {
		dwidget := &dummyWidget{onDestroy: func() {
			destroyedRenderersCnt++
		}}
		Renderer(dwidget)
		SetCanvasForObject(dwidget, dcanvas1, nil)
	}

	for i := 0; i < 22; i++ {
		dwidget := &dummyWidget{onDestroy: func() {
			destroyedRenderersCnt++
		}}
		Renderer(dwidget)
		SetCanvasForObject(dwidget, dcanvas2, nil)
	}

	assert.Equal(t, 42, renderers.Len())
	assert.Equal(t, 42, canvases.Len())

	CleanCanvas(dcanvas1)
	assert.Equal(t, 22, renderers.Len())
	assert.Equal(t, 22, canvases.Len())
	assert.Equal(t, 20, destroyedRenderersCnt)
	canvases.Range(func(_ fyne.CanvasObject, cinfo *canvasInfo) bool {
		assert.Equal(t, dcanvas2, cinfo.canvas)
		return true
	})

	CleanCanvas(dcanvas2)
	assert.Equal(t, 0, renderers.Len())
	assert.Equal(t, 0, canvases.Len())
	assert.Equal(t, 42, destroyedRenderersCnt)
}

func Test_expiringCache(t *testing.T) {
	tm := &timeMock{}
	tm.setTime(10, 10)

	c := &expiringCache{}
	assert.True(t, c.isExpired(tm.now))

	c.setAlive()

	tm.setTime(10, 20)
	assert.False(t, c.isExpired(tm.now))

	tm.setTime(10, 11)
	tm.now = tm.now.Add(ValidDuration)
	assert.True(t, c.isExpired(tm.now))
}

type dummyCanvas struct {
	fyne.Canvas
}

type dummyWidget struct {
	fyne.Widget
	onDestroy func()
}

func (w *dummyWidget) CreateRenderer() fyne.WidgetRenderer {
	return &dummyWidgetRenderer{widget: w}
}

type dummyWidgetRenderer struct {
	widget  *dummyWidget
	objects []fyne.CanvasObject
}

func (r *dummyWidgetRenderer) Destroy() {
	if r.widget.onDestroy != nil {
		r.widget.onDestroy()
	}
}

func (r *dummyWidgetRenderer) Layout(size fyne.Size) {
}

func (r *dummyWidgetRenderer) MinSize() fyne.Size {
	return fyne.NewSize(0, 0)
}

func (r *dummyWidgetRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

func (r *dummyWidgetRenderer) Refresh() {
}

type timeMock struct {
	now time.Time
}

func (t *timeMock) createTime(min, sec int) time.Time {
	return time.Date(2021, time.June, 15, 2, min, sec, 0, time.UTC)
}

func (t *timeMock) setTime(min, sec int) {
	t.now = time.Date(2021, time.June, 15, 2, min, sec, 0, time.UTC)
	timeNow = func() time.Time {
		return t.now
	}
}

func testClearAll() {
	canvases.Clear()
	svgs.Clear()
	textTextures.Clear()
	objectTextures.Clear()
	renderers.Clear()
	blurKernels.Clear()
	timeNow = time.Now
	lastClean = time.Time{}
	framecounter = 1
	rendererCacheLastCleanSize = 0
	canvasCacheLastCleanSize = 0
	fontSizeCacheLastCleanSize = 0
	textTextureLastCleanSize = 0
	objectTexturesLastCleanSize = 0
	shouldCleanRenderers = false
	shouldCleanCanvases = false
	shouldCleanFontSizeCache = false
	shouldCleanTextTextures = false
	shouldCleanObjectTextures = false
}
