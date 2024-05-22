//go:build software

package software

import (
	"fmt"
	"image"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/internal/cache"
	"fyne.io/fyne/v2/internal/driver"
	"fyne.io/fyne/v2/internal/scale"
)

func genMask(c fyne.Canvas, sObj fyne.CanvasObject, objBounds, bounds image.Rectangle) *image.Alpha {
	t := time.Now()
	mask := image.NewAlpha(bounds)

	driver.ReverseWalkVisibleObjectTree(c.Content(), func(obj fyne.CanvasObject, pos, clipPos fyne.Position, clipSize fyne.Size) bool {
		if obj == sObj {
			return true
		}

		var b image.Rectangle
		switch o := obj.(type) {
		case *fyne.Container, fyne.Widget:
			return false
		case *canvas.Image:
			_, _, b = imageCords(c, o, pos)
		case *canvas.Text:
			b, _ = textCords(c, o, pos, bounds)
		case gradient:
			b = gradientCords(c, o, pos)
		case *canvas.Circle:
			_, b, _ = circleCords(c, o, pos, bounds)
		case *canvas.Line:
			_, b, _ = lineCords(c, o, pos, bounds)
		case *canvas.Raster:
			b = rasterCords(c, o, pos)
		case *canvas.Rectangle:
			w := fyne.Min(clipPos.X+clipSize.Width, c.Size().Width)
			h := fyne.Min(clipPos.Y+clipSize.Height, c.Size().Height)
			clip := image.Rect(
				scale.ToScreenCoordinate(c, clipPos.X),
				scale.ToScreenCoordinate(c, clipPos.Y),
				scale.ToScreenCoordinate(c, w),
				scale.ToScreenCoordinate(c, h),
			)
			if o.StrokeWidth > 0 {
				_, _, b = rectangleStrokeCords(c, o, pos, clip)
			} else {
				_, b = rectangleCords(c, o, pos, clip)
			}
		}

		if !b.Overlaps(objBounds) {
			return false
		}

		tex, ok := cache.GetTexture(obj)
		if !ok {
			return false
		}
		m := uint32(1<<16 - 1)

		// TODO: The color mixing seems wrong, but I'm not quite sure why.
		switch o := tex.(type) {
		case *image.NRGBA:
			for y := 0; y < o.Bounds().Dy(); y++ {
				for x := 0; x < o.Bounds().Dx(); x++ {
					da := &mask.Pix[mask.PixOffset(b.Min.X+x, b.Min.Y+y)]
					sa := uint32(o.Pix[o.PixOffset(x, y)+3]) * 0x101
					a := (m - sa) * 0x101
					if sa > 0 {
						*da = uint8((uint32(*da)*a/m + sa) >> 8)
					}
				}
			}
		case *image.RGBA:
			for y := 0; y < o.Bounds().Dy(); y++ {
				for x := 0; x < o.Bounds().Dx(); x++ {
					da := &mask.Pix[mask.PixOffset(b.Min.X+x, b.Min.Y+y)]
					sa := uint32(o.Pix[o.PixOffset(x, y)+3]) * 0x101
					a := (m - sa) * 0x101
					if sa > 0 {
						*da = uint8((uint32(*da)*a/m + sa) >> 8)
					}
				}
			}
		}

		return false
	}, nil)

	for i := 0; i < len(mask.Pix); i++ {
		mask.Pix[i] = 0xff - mask.Pix[i]
	}

	fmt.Println("Generated mask in", time.Since(t))
	return mask
}
