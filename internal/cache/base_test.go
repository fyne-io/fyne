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
	destroyedRenderersCnt := 0
	testClearAll()
	tm := &timeMock{}

	for k := 0; k < 2; k++ {
		tm.setTime(10, 10+k*10)
		for i := 0; i < 20; i++ {
			SetSvg(fmt.Sprintf("%d%d", k, i), nil, nil, i, i+1)
			Renderer(&dummyWidget{onDestroy: func() {
				destroyedRenderersCnt++
			}})
			SetCanvasForObject(&dummyWidget{}, &dummyCanvas{}, nil)
		}
	}

	t.Run("no_expired_objects", func(t *testing.T) {
		lastClean = tm.createTime(10, 20)
		Clean()
		assert.Equal(t, svgs.Len(), 40)
		assert.Equal(t, 40, renderers.Len())
		assert.Equal(t, 40, canvases.Len())
		assert.Zero(t, destroyedRenderersCnt)
		assert.Equal(t, tm.now, lastClean)

		tm.setTime(10, 30)
		Clean()
		assert.Equal(t, svgs.Len(), 40)
		assert.Equal(t, 40, renderers.Len())
		assert.Equal(t, 40, canvases.Len())
		assert.Zero(t, destroyedRenderersCnt)
		assert.Equal(t, tm.now, lastClean)
	})

	t.Run("do_not_clean_too_fast", func(t *testing.T) {
		lastClean = tm.createTime(10, 30)
		// when no canvas refresh and has been transcurred less than
		// cleanTaskInterval duration, no clean task should occur.
		tm.setTime(10, 42)
		Clean()
		assert.Less(t, lastClean.UnixNano(), tm.now.UnixNano())

		Clean()
		assert.Equal(t, tm.now, lastClean)

		// when canvas refresh the clean task is only executed if it has been
		// transcurred more than 10 seconds since the lastClean.
		tm.setTime(10, 45)
		Clean()
		assert.Less(t, lastClean.UnixNano(), tm.now.UnixNano())

		tm.setTime(10, 53)
		Clean()
		assert.Equal(t, tm.now, lastClean)

		assert.Equal(t, svgs.Len(), 40)
		assert.Equal(t, 40, renderers.Len())
		assert.Equal(t, 40, canvases.Len())
		assert.Zero(t, destroyedRenderersCnt)
	})

	t.Run("clean_no_canvas_refresh", func(t *testing.T) {
		lastClean = tm.createTime(10, 11)
		tm.setTime(11, 12)
		Clean()
		assert.Equal(t, svgs.Len(), 20)
		assert.Equal(t, renderers.Len(), 40)
		assert.Equal(t, canvases.Len(), 40)
		assert.Zero(t, destroyedRenderersCnt)

		tm.setTime(11, 42)
		Clean()
		assert.Equal(t, svgs.Len(), 0)
		assert.Equal(t, renderers.Len(), 40)
		assert.Equal(t, canvases.Len(), 40)
		assert.Zero(t, destroyedRenderersCnt)
	})

	t.Run("clean_canvas_refresh", func(t *testing.T) {
		lastClean = tm.createTime(10, 11)
		tm.setTime(11, 11)
		Clean()
		assert.Equal(t, svgs.Len(), 0)
		assert.Equal(t, 20, renderers.Len())
		assert.Equal(t, 20, canvases.Len())
		assert.Equal(t, 20, destroyedRenderersCnt)

		tm.setTime(11, 22)
		Clean()
		assert.Equal(t, svgs.Len(), 0)
		assert.Equal(t, 0, renderers.Len())
		assert.Equal(t, 0, canvases.Len())
		assert.Equal(t, 40, destroyedRenderersCnt)
	})

	t.Run("skipped_clean_with_canvas_refresh", func(t *testing.T) {
		testClearAll()
		lastClean = tm.createTime(13, 10)
		tm.setTime(13, 10)
		assert.False(t, skippedCleanWithCanvasRefresh)
		Clean()
		assert.Equal(t, tm.now, lastClean)

		Renderer(&dummyWidget{})

		tm.setTime(13, 15)
		Clean()
		assert.True(t, skippedCleanWithCanvasRefresh)
		assert.Less(t, lastClean.UnixNano(), tm.now.UnixNano())
		assert.Equal(t, 1, renderers.Len())

		tm.setTime(14, 21)
		Clean()
		assert.False(t, skippedCleanWithCanvasRefresh)
		assert.Equal(t, tm.now, lastClean)
		assert.Equal(t, 0, renderers.Len())
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
	skippedCleanWithCanvasRefresh = false
	canvases.Clear()
	svgs.Clear()
	textTextures.Clear()
	objectTextures.Clear()
	renderers.Clear()
	timeNow = time.Now
}
