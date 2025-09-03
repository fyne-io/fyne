package common

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/internal/async"
	"fyne.io/fyne/v2/internal/build"
)

type deduplicatedObjectQueue struct {
	queue *async.CanvasObjectQueue
	dedup async.Map[fyne.CanvasObject, struct{}]
}

// In adds an object to the queue if it is not already present.
func (q *deduplicatedObjectQueue) In(obj fyne.CanvasObject) {
	_, exists := q.dedup.Load(obj)
	if exists {
		return
	}

	q.queue.In(obj)
	q.dedup.Store(obj, struct{}{})
}

// Out removes and returns the next object from the queue.
// It assumes that the whole queue is drained and defers clearing
// the deduplication map until it is empty.
func (q *deduplicatedObjectQueue) Out() fyne.CanvasObject {
	if q.queue.Len() == 0 {
		q.dedup.Clear()
		return nil
	}

	out := q.queue.Out()
	if !build.MigratedToFyneDo() {
		q.dedup.Delete(out)
	}
	return out
}

// Len returns the number of elements in the queue.
func (q *deduplicatedObjectQueue) Len() uint64 {
	return q.queue.Len()
}

type renderCacheTree struct {
	root *RenderCacheNode
}

// RenderCacheNode represents a node in a render cache tree.
type RenderCacheNode struct {
	// structural data
	firstChild  *RenderCacheNode
	nextSibling *RenderCacheNode
	obj         fyne.CanvasObject
	parent      *RenderCacheNode
	// cache data
	minSize fyne.Size
}

// Obj returns the node object.
func (r *RenderCacheNode) Obj() fyne.CanvasObject {
	return r.obj
}

type overlayStack struct {
	OverlayStack

	renderCaches []*renderCacheTree
}

func (o *overlayStack) Add(overlay fyne.CanvasObject) {
	if overlay == nil {
		return
	}
	o.add(overlay)
}

func (o *overlayStack) Remove(overlay fyne.CanvasObject) {
	if overlay == nil || len(o.List()) == 0 {
		return
	}
	o.remove(overlay)
}

func (o *overlayStack) add(overlay fyne.CanvasObject) {
	o.renderCaches = append(o.renderCaches, &renderCacheTree{root: &RenderCacheNode{obj: overlay}})
	o.OverlayStack.Add(overlay)
}

func (o *overlayStack) remove(overlay fyne.CanvasObject) {
	o.OverlayStack.Remove(overlay)
	overlayCount := len(o.List())

	// it is possible that overlays are removed implicitly and render caches already cleared out
	if overlayCount >= len(o.renderCaches) {
		return
	}

	o.renderCaches[overlayCount] = nil // release memory reference to removed element
	o.renderCaches = o.renderCaches[:overlayCount]
}
