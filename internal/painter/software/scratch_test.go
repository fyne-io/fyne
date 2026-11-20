package software

import (
	"image"
	"testing"
)

func TestReleaseScratch_DropsOversizedFrames(t *testing.T) {
	p := NewPainter()
	big := image.NewNRGBA(image.Rect(0, 0, 1, maxPooledFrameBytes/4+1))
	p.ReleaseScratch(big)
	if got, _ := framePool.Get().(*image.NRGBA); got != nil && cap(got.Pix) > maxPooledFrameBytes {
		t.Fatalf("oversized frame was kept in the pool (cap=%d)", cap(got.Pix))
	}
	p.ReleaseScratch(nil)
}
