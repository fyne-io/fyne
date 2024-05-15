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

func genMask(c fyne.Canvas, obj fyne.CanvasObject, bounds image.Rectangle) *image.Alpha {
	t := time.Now()
	mask := image.NewAlpha(bounds)
	driver.ReverseWalkVisibleObjectTree(obj, func(obj fyne.CanvasObject, pos, clipPos fyne.Position, clipSize fyne.Size) bool {
		var objStart image.Point
		switch o := obj.(type) {
		case *fyne.Container, fyne.Widget:
			return false
		case *canvas.Image:
			_, _, b := imageCords(c, o, pos)
			objStart = b.Min
		case *canvas.Text:
			_, _, b, _ := textCords(c, o, pos, bounds)
			objStart = b.Min
		case gradient:
			// p.drawGradient(c, o, pos, base, clip)
		case *canvas.Circle:
			_, b := circleCords(c, o, pos, bounds)
			objStart = b.Min
		case *canvas.Line:
			_, b := lineCords(c, o, pos, bounds)
			objStart = b.Min
		case *canvas.Raster:
			b := rasterCords(c, o, pos)
			objStart = b.Min
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
				_, b := rectangleStrokeCords(c, o, pos, clip)
				objStart = b.Min
			} else {
				b := rectangleCords(c, o, pos, clip)
				objStart = b.Min
			}
		}
		tex, ok := cache.GetTexture(obj)
		if !ok {
			return true
		}
		switch o := tex.(type) {
		case *image.NRGBA:
			for y := 0; y < o.Bounds().Dy(); y++ {
				for x := 0; x < o.Bounds().Dx(); x++ {
					da := &mask.Pix[mask.PixOffset(objStart.X+x, objStart.Y+y)]
					sa := uint32(o.Pix[o.PixOffset(x, y)+3]) | uint32(o.Pix[o.PixOffset(x, y)+3])<<8
					a := ((1<<16 - 1) - uint32(sa)) * 0x101
					if sa > 0 {
						*da = uint8((uint32(*da)*a/(1<<16-1) + sa) >> 8)
					}
				}
			}
		case *image.RGBA:
			for y := 0; y < o.Bounds().Dy(); y++ {
				for x := 0; x < o.Bounds().Dx(); x++ {
					da := &mask.Pix[mask.PixOffset(objStart.X+x, objStart.Y+y)]
					sa := uint32(o.Pix[o.PixOffset(x, y)+3]) | uint32(o.Pix[o.PixOffset(x, y)+3])<<8
					a := ((1<<16 - 1) - uint32(sa)) * 0x101
					if sa > 0 {
						*da = uint8((uint32(*da)*a/(1<<16-1) + sa) >> 8)
					}
				}
			}
		}

		return false
	}, nil)

	for i := 0; i < len(mask.Pix); i++ {
		mask.Pix[i] = 0xff - mask.Pix[i]
	}

	// tex := image.NewRGBA(bounds)
	// draw.Draw(tex, tex.Bounds(), image.NewUniform(color.Black), image.Point{}, draw.Src)
	// draw.DrawMask(tex, bounds, image.NewUniform(color.White), image.Point{}, mask, image.Point{}, draw.Over)

	fmt.Println("Generated mask in", time.Since(t))
	return mask
}
