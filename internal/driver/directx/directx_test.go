//go:build windows && directx

package directx

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"testing"
	"unsafe"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/theme"
)

func TestKeyToName(t *testing.T) {
	cases := []struct {
		vk   uintptr
		want fyne.KeyName
	}{
		{'A', fyne.KeyA},
		{'Z', fyne.KeyZ},
		{'0', fyne.Key0},
		{'9', fyne.Key9},
		{vkF1, fyne.KeyF1},
		{vkF1 + 11, fyne.KeyF12},
		{vkReturn, fyne.KeyReturn},
		{vkEscape, fyne.KeyEscape},
		{vkBack, fyne.KeyBackspace},
		{vkLeft, fyne.KeyLeft},
		{vkPrior, fyne.KeyPageUp},
		{vkNumpad0 + 5, fyne.Key5},
		{vkLControl, desktop.KeyControlLeft},
		{vkOEM3, fyne.KeyBackTick},
		{0x07, fyne.KeyUnknown}, // unmapped
	}
	for _, c := range cases {
		if got := keyToName(c.vk); got != c.want {
			t.Errorf("keyToName(%#x) = %q, want %q", c.vk, got, c.want)
		}
	}
}

// TestCursorForCoversEveryStandardCursor pins the cursor table. The failure this
// guards against is silent: an unmapped cursor falls through to the arrow, which
// looks like "resize doesn't work" rather than like a bug. Iterating to
// HiddenCursor also means a cursor added to Fyne later fails here instead of
// quietly becoming an arrow.
func TestCursorForCoversEveryStandardCursor(t *testing.T) {
	want := map[desktop.StandardCursor]uint16{
		desktop.DefaultCursor:    idcArrow,
		desktop.TextCursor:       idcIBeam,
		desktop.CrosshairCursor:  idcCross,
		desktop.PointerCursor:    idcHand,
		desktop.HResizeCursor:    idcSizeWE,
		desktop.VResizeCursor:    idcSizeNS,
		desktop.NESWResizeCursor: idcSizeNESW,
		desktop.NWSEResizeCursor: idcSizeNWSE,
	}

	for c := desktop.DefaultCursor; c <= desktop.HiddenCursor; c++ {
		id, visible := cursorFor(c)

		if c == desktop.HiddenCursor {
			if visible {
				t.Error("HiddenCursor must report not visible so SetCursor gets a null handle")
			}
			continue
		}
		if !visible {
			t.Errorf("cursor %d reported hidden", c)
			continue
		}
		expect, ok := want[c]
		if !ok {
			t.Errorf("cursor %d is not covered by the test table - was a new cursor added?", c)
			continue
		}
		if id != expect {
			t.Errorf("cursor %d mapped to %d, want %d", c, id, expect)
		}
		if c != desktop.DefaultCursor && id == idcArrow {
			t.Errorf("cursor %d fell through to the default arrow", c)
		}
	}
}

// TestMenuLabel checks the accelerator text. A menu item whose shortcut is not a
// KeyboardShortcut has nothing to show, and the tab separator is what makes Win32
// right-align the rest, so a plain space here would silently render as one run.
func TestMenuLabel(t *testing.T) {
	plain := fyne.NewMenuItem("Open", nil)
	if got := menuLabel(plain); got != "Open" {
		t.Errorf("no shortcut = %q, want %q", got, "Open")
	}

	saved := fyne.NewMenuItem("Save", nil)
	saved.Shortcut = &desktop.CustomShortcut{
		KeyName:  fyne.KeyS,
		Modifier: fyne.KeyModifierControl | fyne.KeyModifierShift,
	}
	if got, want := menuLabel(saved), "Save\tCtrl+Shift+S"; got != want {
		t.Errorf("menuLabel = %q, want %q", got, want)
	}
}

// TestTriggerMenuShortcut covers the recursion into submenus and the two ways an
// item declines the key: no shortcut at all, or a shortcut with no action.
func TestTriggerMenuShortcut(t *testing.T) {
	sh := &desktop.CustomShortcut{KeyName: fyne.KeyN, Modifier: fyne.KeyModifierControl}

	ran := false
	deep := fyne.NewMenuItem("New", func() { ran = true })
	deep.Shortcut = sh

	noAction := fyne.NewMenuItem("Dead", nil)
	noAction.Shortcut = sh

	parent := fyne.NewMenuItem("File", nil)
	parent.ChildMenu = fyne.NewMenu("File", deep)

	if triggerMenuShortcut(sh, fyne.NewMenu("", fyne.NewMenuItem("Plain", nil))) {
		t.Error("an item with no shortcut claimed the key")
	}
	if triggerMenuShortcut(sh, fyne.NewMenu("", noAction)) {
		t.Error("an item with no action claimed the key")
	}
	if !triggerMenuShortcut(sh, fyne.NewMenu("", parent)) || !ran {
		t.Error("shortcut in a submenu was not triggered")
	}
}

// TestMenuIcon drives a real menu through the icon path. It fails silently in
// production - a wrong MENUITEMINFOW layout makes SetMenuItemInfoW reject cbSize,
// and nothing else reports it - so calling the API is the only way to catch it.
func TestMenuIcon(t *testing.T) {
	test.NewTempApp(t)

	menu, _, _ := procCreatePopupMenu.Call()
	if menu == 0 {
		t.Fatal("CreatePopupMenu failed")
	}
	defer procDestroyMenu.Call(menu)
	appendMenu(menu, mfString, firstMenuID, "Home")

	s := &menuState{iconSize: 16}
	defer s.release()
	s.setIcon(menu, 0, &fyne.MenuItem{Label: "Home", Icon: theme.HomeIcon()})

	if len(s.icons) != 1 || s.icons[0] == 0 {
		t.Fatal("setIcon did not build an icon")
	}

	// Read the item back: hbmpItem must be the callback marker and the icon must
	// have survived in the item data, which is where WM_DRAWITEM looks for it.
	info := menuItemInfoW{
		cbSize: uint32(unsafe.Sizeof(menuItemInfoW{})),
		fMask:  miimBitmap | miimData,
	}
	r, _, err := procGetMenuItemInfoW.Call(menu, 0, 1, uintptr(unsafe.Pointer(&info)))
	if r == 0 {
		t.Fatalf("GetMenuItemInfoW failed: %v", err)
	}
	if uintptr(info.hbmpItem) != hbmMenuCallback {
		t.Errorf("hbmpItem is %#x, want HBMMENU_CALLBACK %#x", info.hbmpItem, hbmMenuCallback)
	}
	if info.dwItemData != uintptr(s.icons[0]) {
		t.Errorf("item data is %#x, want the icon handle %#x", info.dwItemData, s.icons[0])
	}
}

// TestMenuIconOwnerDraw runs the two owner-draw replies against real Win32
// structures. A layout drift here reads garbage out of the message and the icon
// silently never appears.
func TestMenuIconOwnerDraw(t *testing.T) {
	test.NewTempApp(t)

	s := &menuState{iconSize: 16}
	defer s.release()

	measure := measureItemStruct{ctlType: odtMenu, itemData: 1}
	if !s.measureMenuIcon(uintptr(unsafe.Pointer(&measure))) {
		t.Fatal("measureMenuIcon did not claim the menu item")
	}
	if measure.itemWidth != 16 || measure.itemHeight != 16 {
		t.Errorf("measured %dx%d, want 16x16", measure.itemWidth, measure.itemHeight)
	}

	// A non-menu owner-draw message belongs to someone else and must be declined.
	other := measureItemStruct{ctlType: odtMenu + 1, itemData: 1}
	if s.measureMenuIcon(uintptr(unsafe.Pointer(&other))) {
		t.Error("measureMenuIcon claimed a message that was not for a menu")
	}
	noIcon := measureItemStruct{ctlType: odtMenu}
	if s.measureMenuIcon(uintptr(unsafe.Pointer(&noIcon))) {
		t.Error("measureMenuIcon claimed an item that carries no icon")
	}

	draw := drawItemStruct{ctlType: odtMenu}
	if s.drawMenuIcon(uintptr(unsafe.Pointer(&draw))) {
		t.Error("drawMenuIcon claimed an item that carries no icon")
	}
}

// TestBlendModes pins the two blend equations. Textures come from Go's
// alpha-premultiplied image.RGBA, so they must blend with ONE; using SRC_ALPHA
// applies alpha a second time and quietly destroys anti-aliased glyph edges.
// Nothing errors when this is wrong, so only a test catches it.

// TestCharInputSurrogatePairs pins the WM_CHAR handling: a non-BMP character
// arrives as two messages carrying a UTF-16 surrogate pair, which must reach
// TypedRune as one rune. Getting this wrong turns every emoji into garbage.
func TestCharInputSurrogatePairs(t *testing.T) {
	w := &window{canvas: newCanvas()}
	var got []rune
	w.canvas.onTypedRune = func(r rune) { got = append(got, r) }

	// "a", then U+1F600 as its surrogate halves, then a lone low surrogate
	// (must be dropped), then "b".
	for _, wp := range []uintptr{'a', 0xd83d, 0xde00, 0xde00, 'b'} {
		w.handleMessage(wmChar, wp, 0)
	}
	if want := "a\U0001F600b"; string(got) != want {
		t.Errorf("typed runes = %q, want %q", string(got), want)
	}
}

// TestKeyNameFor pins the lParam-sensitive key resolution: the numpad Enter
// shares VK_RETURN with the main Return key and only the extended bit tells
// them apart.
func TestKeyNameFor(t *testing.T) {
	if got := keyNameFor(vkReturn, 0); got != fyne.KeyReturn {
		t.Errorf("plain Return = %q, want %q", got, fyne.KeyReturn)
	}
	if got := keyNameFor(vkReturn, 1<<24); got != fyne.KeyEnter {
		t.Errorf("extended Return = %q, want %q", got, fyne.KeyEnter)
	}
	if got := keyNameFor(vkApps, 0); got != desktop.KeyMenu {
		t.Errorf("VK_APPS = %q, want %q", got, desktop.KeyMenu)
	}
	if got := scanCodeFor(0x1C<<16 | 1<<24); got != 0x11C {
		t.Errorf("scanCodeFor = %#x, want 0x11C", got)
	}
}

// TestToOSIconWritesClassicBMPEntry pins the hand-rolled ICO writer: always a
// 32px uncompressed 32bpp BMP entry (Wine cannot decode PNG-in-ICO, and
// LoadImage materialises SM_CXICON anyway), bottom-up BGRA with a doubled-height
// header and an AND mask derived from the alpha channel.
func TestToOSIconWritesClassicBMPEntry(t *testing.T) {
	src := image.NewNRGBA(image.Rect(0, 0, systrayIconSize, systrayIconSize))
	src.SetNRGBA(1, 1, color.NRGBA{R: 10, G: 20, B: 30, A: 255})
	buf := &bytes.Buffer{}
	if err := png.Encode(buf, src); err != nil {
		t.Fatal(err)
	}

	data, err := toOSIcon(buf.Bytes())
	if err != nil {
		t.Fatal(err)
	}

	// 22 byte directory + 40 byte header + 32*32*4 XOR + 32 rows * 4 byte mask.
	if len(data) != 4286 {
		t.Fatalf("encoded size = %d, want 4286", len(data))
	}
	head := []byte{0, 0, 1, 0, 1, 0, 32, 32, 0, 0, 1, 0, 32, 0}
	if !bytes.Equal(data[:14], head) {
		t.Errorf("directory = %v, want %v", data[:14], head)
	}
	if data[22] != 40 || data[26] != 32 || data[30] != 64 {
		t.Errorf("info header size/width/doubled height = %d/%d/%d, want 40/32/64",
			data[22], data[26], data[30])
	}
	// Rows are bottom-up: y=1 lands 30 rows in; x=1 is 4 bytes along.
	px := data[62+30*32*4+4:]
	if px[0] != 30 || px[1] != 20 || px[2] != 10 || px[3] != 255 {
		t.Errorf("pixel = BGRA %v, want [30 20 10 255]", px[:4])
	}
	// Mask row for y=1: x=0 transparent (bit set), x=1 opaque (bit clear),
	// x=2..7 transparent again.
	if got := data[62+32*32*4+30*4]; got != 0xbf {
		t.Errorf("mask byte = %#x, want 0xbf", got)
	}
}
