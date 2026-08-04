//go:build windows

package directx

import (
	"image"
	"testing"
	"unsafe"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

// meanOpaqueInk averages the colour of the fully opaque pixels of an icon, which
// for a monochrome glyph is the ink it was drawn in.
func meanOpaqueInk(t *testing.T, img image.Image) (r, g, b uint32) {
	t.Helper()

	bounds := img.Bounds()
	var sr, sg, sb, n uint32
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			pr, pg, pb, pa := img.At(x, y).RGBA()
			if pa>>8 < 250 {
				continue
			}
			sr, sg, sb, n = sr+pr>>8, sg+pg>>8, sb+pb>>8, n+1
		}
	}
	if n == 0 {
		t.Fatal("icon rasterised with no opaque pixels")
	}
	return sr / n, sg / n, sb / n
}

// TestMenuIconMatchesMenuText covers icons coming out invisible in a native menu.
// Fyne's default theme is dark, so a themed icon rasterises white; the Win32 menu
// it is drawn into is painted by Windows, which on a light system means white ink
// on a white background. Nothing errors - the icon is simply not there.
func TestMenuIconMatchesMenuText(t *testing.T) {
	icon := theme.ContentCopyIcon()

	wr, wg, wb, _ := sysColour(colorMenuText).RGBA()
	wantR, wantG, wantB := wr>>8, wg>>8, wb>>8

	gotR, gotG, gotB := meanOpaqueInk(t, rasterise(menuIconResource(icon), 16))
	if gotR != wantR || gotG != wantG || gotB != wantB {
		t.Errorf("menu icon ink is RGB(%d, %d, %d), want COLOR_MENUTEXT RGB(%d, %d, %d)",
			gotR, gotG, gotB, wantR, wantG, wantB)
	}

	// The theme foreground is the colour this used to come out in. If the two ever
	// agree the test proves nothing, so only compare when they genuinely differ.
	fr, fg, fb, _ := theme.Color(theme.ColorNameForeground).RGBA()
	if fr>>8 != wantR || fg>>8 != wantG || fb>>8 != wantB {
		if gotR == fr>>8 && gotG == fg>>8 && gotB == fb>>8 {
			t.Error("menu icon still uses the app theme foreground rather than the menu text colour")
		}
	}
}

// TestMenuIconLeavesArtworkAlone checks that only themed resources are recoloured.
// A caller's own icon may be multi coloured, and flattening it to one ink would
// destroy it.
func TestMenuIconLeavesArtworkAlone(t *testing.T) {
	raw := fyne.NewStaticResource("artwork.svg", theme.ContentCopyIcon().Content())
	if got := menuIconResource(raw); got != fyne.Resource(raw) {
		t.Errorf("non-themed resource was replaced by %v, want it passed through", got.Name())
	}
}

// TestMenuItemInfoLayout pins the Go struct to the C MENUITEMINFOW it is passed as.
// cbSize is validated by Windows, so a layout drift makes every SetMenuItemInfoW
// call fail - and the icon quietly not appear.
func TestMenuItemInfoLayout(t *testing.T) {
	var m menuItemInfoW
	if got := unsafe.Sizeof(m); got != 80 {
		t.Errorf("sizeof(menuItemInfoW) = %d, want 80 (C MENUITEMINFOW on x64)", got)
	}

	var mi measureItemStruct
	if got := unsafe.Sizeof(mi); got != 32 {
		t.Errorf("sizeof(measureItemStruct) = %d, want 32 (C MEASUREITEMSTRUCT on x64)", got)
	}
	var di drawItemStruct
	if got := unsafe.Sizeof(di); got != 64 {
		t.Errorf("sizeof(drawItemStruct) = %d, want 64 (C DRAWITEMSTRUCT on x64)", got)
	}

	for _, c := range []struct {
		name string
		got  uintptr
		want uintptr
	}{
		{"measureItemStruct.itemWidth", unsafe.Offsetof(mi.itemWidth), 12},
		{"measureItemStruct.itemHeight", unsafe.Offsetof(mi.itemHeight), 16},
		{"measureItemStruct.itemData", unsafe.Offsetof(mi.itemData), 24},
		{"drawItemStruct.itemState", unsafe.Offsetof(di.itemState), 16},
		{"drawItemStruct.hDC", unsafe.Offsetof(di.hDC), 32},
		{"drawItemStruct.rcItem", unsafe.Offsetof(di.rcItem), 40},
		{"drawItemStruct.itemData", unsafe.Offsetof(di.itemData), 56},
		{"cbSize", unsafe.Offsetof(m.cbSize), 0},
		{"fMask", unsafe.Offsetof(m.fMask), 4},
		{"fType", unsafe.Offsetof(m.fType), 8},
		{"fState", unsafe.Offsetof(m.fState), 12},
		{"wID", unsafe.Offsetof(m.wID), 16},
		{"hSubMenu", unsafe.Offsetof(m.hSubMenu), 24},
		{"hbmpChecked", unsafe.Offsetof(m.hbmpChecked), 32},
		{"hbmpUnchecked", unsafe.Offsetof(m.hbmpUnchecked), 40},
		{"dwItemData", unsafe.Offsetof(m.dwItemData), 48},
		{"dwTypeData", unsafe.Offsetof(m.dwTypeData), 56},
		{"cch", unsafe.Offsetof(m.cch), 64},
		{"hbmpItem", unsafe.Offsetof(m.hbmpItem), 72},
	} {
		if c.got != c.want {
			t.Errorf("menuItemInfoW.%s is at offset %d, want %d", c.name, c.got, c.want)
		}
	}
}
