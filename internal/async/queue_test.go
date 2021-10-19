package async_test

import (
	"strings"
	"testing"

	"fyne.io/fyne/v2/internal/async"
)

func TestQueue(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		q := async.NewQueue()
		if q.Out() != nil {
			t.Fatalf("dequeue empty queue returns non-nil")
		}
	})

	t.Run("length", func(t *testing.T) {
		q := async.NewQueue()
		if q.Len() != 0 {
			t.Fatalf("empty queue has non-zero length")
		}

		q.In(1)
		if q.Len() != 1 {
			t.Fatalf("count of enqueue wrong, want %d, got %d.", 1, q.Len())
		}

		q.Out()
		if q.Len() != 0 {
			t.Fatalf("count of dequeue wrong, want %d, got %d", 0, q.Len())
		}
	})

	t.Run("in-out", func(t *testing.T) {
		q := async.NewQueue()

		want := []string{
			"1st item",
			"2nd item",
			"3rd item",
		}

		for i := 0; i < len(want); i++ {
			q.In(want[i])
		}

		var x []string
		for {
			e := q.Out()
			if e == nil {
				break
			}
			x = append(x, e.(string))
		}

		for i := 0; i < len(want); i++ {
			if strings.Compare(x[i], want[i]) != 0 {
				t.Fatalf("input does not match output, want %+v, got %+v", want[i], x[i])
			}
		}
	})
}
