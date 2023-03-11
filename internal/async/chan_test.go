package async_test

import (
	"runtime"
	"sync"
	"testing"

	"fyne.io/fyne/v2/internal/async"
)

func TestUnboundedChann(t *testing.T) {
	N := 200
	if testing.Short() {
		N = 20
	}

	wg := sync.WaitGroup{}
	for i := 0; i < N; i++ {
		t.Run("interface{}", func(t *testing.T) {
			t.Run("send", func(t *testing.T) {
				// Ensure send to an unbounded channel does not block.
				c := async.NewUnboundedInterfaceChan()
				blocked := false
				wg.Add(1)
				go func() {
					defer wg.Done()
					select {
					case c.In() <- true:
					default:
						blocked = true
					}
				}()
				wg.Wait()
				if blocked {
					t.Fatalf("send op to an unbounded channel blocked")
				}
				c.Close()
			})

			t.Run("recv", func(t *testing.T) {
				// Ensure that receive op from unbounded chan can happen on
				// the same goroutine of send op.
				c := async.NewUnboundedInterfaceChan()
				wg.Add(1)
				go func() {
					defer wg.Done()
					c.In() <- true
					<-c.Out()
				}()
				wg.Wait()
				c.Close()
			})

			t.Run("order", func(t *testing.T) {
				// Ensure that the unbounded channel processes everything FIFO.
				c := async.NewUnboundedInterfaceChan()
				for i := 0; i < 1<<11; i++ {
					c.In() <- i
				}
				for i := 0; i < 1<<11; i++ {
					if val := <-c.Out(); val != i {
						t.Fatalf("unbounded channel passes messages in a non-FIFO order, got %v want %v", val, i)
					}
				}
				c.Close()
			})

			t.Run("graceful-close", func(t *testing.T) {
				grs := runtime.NumGoroutine()
				N := 10
				n := 0
				done := make(chan struct{})
				ch := async.NewUnboundedInterfaceChan()
				for i := 0; i < N; i++ {
					ch.In() <- true
				}
				go func() {
					for range ch.Out() {
						n++
					}
					done <- struct{}{}
				}()
				ch.Close()
				<-done
				runtime.GC()
				if runtime.NumGoroutine() > grs+2 {
					t.Fatalf("leaking goroutines: %v", n)
				}
				if n != N {
					t.Fatalf("After close, not all elements are received, got %v, want %v", n, N)
				}
			})
		})
		t.Run("struct{}", func(t *testing.T) {
			t.Run("send", func(t *testing.T) {
				// Ensure send to an unbounded channel does not block.
				c := async.NewUnboundedStructChan()
				blocked := false
				wg.Add(1)
				go func() {
					defer wg.Done()
					select {
					case c.In() <- struct{}{}:
					default:
						blocked = true
					}
				}()
				<-c.Out()
				wg.Wait()
				if blocked {
					t.Fatalf("send op to an unbounded channel blocked")
				}
				c.Close()
			})

			t.Run("recv", func(t *testing.T) {
				// Ensure that receive op from unbounded chan can happen on
				// the same goroutine of send op.
				c := async.NewUnboundedStructChan()
				wg.Add(1)
				go func() {
					defer wg.Done()
					c.In() <- struct{}{}
					<-c.Out()
				}()
				wg.Wait()
				c.Close()
			})

			t.Run("order", func(t *testing.T) {
				// Ensure that the unbounded channel processes everything FIFO.
				c := async.NewUnboundedStructChan()
				for i := 0; i < 1<<11; i++ {
					c.In() <- struct{}{}
				}
				n := 0
				for i := 0; i < 1<<11; i++ {
					if _, ok := <-c.Out(); ok {
						n++
					}
				}
				if n != 1<<11 {
					t.Fatalf("unbounded channel missed a message, got %v want %v", n, 1<<11)
				}
				c.Close()
			})

			t.Run("graceful-close", func(t *testing.T) {
				grs := runtime.NumGoroutine()
				N := 10
				n := 0
				done := make(chan struct{})
				ch := async.NewUnboundedStructChan()
				for i := 0; i < N; i++ {
					ch.In() <- struct{}{}
				}
				go func() {
					for range ch.Out() {
						n++
					}
					done <- struct{}{}
				}()
				ch.Close()
				<-done
				runtime.GC()
				if runtime.NumGoroutine() > grs+2 {
					t.Fatalf("leaking goroutines: %v", n)
				}
				if n != N {
					t.Fatalf("After close, not all elements are received, got %v, want %v", n, N)
				}
			})
		})
	}
}

func BenchmarkUnboundedChann(b *testing.B) {
	b.Run("interface{}", func(b *testing.B) {
		b.Run("sync", func(b *testing.B) {
			c := async.NewUnboundedInterfaceChan()
			b.ResetTimer()
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				c.In() <- struct{}{}
				<-c.Out()
			}
		})
		b.Run("async", func(b *testing.B) {
			c := async.NewUnboundedInterfaceChan()
			b.ResetTimer()
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				go func() { c.In() <- struct{}{} }()
				<-c.Out()
			}
		})
	})
	b.Run("struct{}", func(b *testing.B) {
		b.Run("sync", func(b *testing.B) {
			c := async.NewUnboundedStructChan()
			b.ResetTimer()
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				c.In() <- struct{}{}
				<-c.Out()
			}
		})
		b.Run("async", func(b *testing.B) {
			c := async.NewUnboundedStructChan()
			b.ResetTimer()
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				go func() { c.In() <- struct{}{} }()
				<-c.Out()
			}
		})
	})
}
