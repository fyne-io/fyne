package async_test

import (
	"image/color"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/internal/async"
)

func TestQueue(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		q := async.NewCanvasObjectQueue()
		if q.Out() != nil {
			t.Fatalf("dequeue empty queue returns non-nil")
		}
	})

	t.Run("length", func(t *testing.T) {
		q := async.NewCanvasObjectQueue()
		if q.Len() != 0 {
			t.Fatalf("empty queue has non-zero length")
		}

		obj := canvas.NewRectangle(color.Black)

		q.In(obj)
		if q.Len() != 1 {
			t.Fatalf("count of enqueue wrong, want %d, got %d.", 1, q.Len())
		}

		q.Out()
		if q.Len() != 0 {
			t.Fatalf("count of dequeue wrong, want %d, got %d", 0, q.Len())
		}
	})

	t.Run("in-out", func(t *testing.T) {
		q := async.NewCanvasObjectQueue()

		want := []fyne.CanvasObject{
			canvas.NewRectangle(color.Black),
			canvas.NewRectangle(color.Black),
			canvas.NewRectangle(color.Black),
		}

		for i := 0; i < len(want); i++ {
			q.In(want[i])
		}

		var x []fyne.CanvasObject
		for {
			e := q.Out()
			if e == nil {
				break
			}
			x = append(x, e)
		}

		for i := 0; i < len(want); i++ {
			if x[i] != want[i] {
				t.Fatalf("input does not match output, want %+v, got %+v", want[i], x[i])
			}
		}
	})
}
