package efl

// #cgo pkg-config: ecore ecore-evas ecore-wl2 evas
// #cgo CFLAGS: -DEFL_BETA_API_SUPPORT=1
// #include <Ecore.h>
// #include <Ecore_Evas.h>
// #include <Ecore_Wl2.h>
// #include <Evas.h>
import "C"
import "fmt"

import "github.com/fyne-io/fyne/ui"

type window struct {
	ee *C.Ecore_Evas
}

func (w *window) Init(ee *C.Ecore_Evas) {
	w.ee = ee
}

func (w *window) Title() string {
	return C.GoString(C.ecore_evas_title_get(w.ee))
}

func (w *window) SetTitle(title string) {
	C.ecore_evas_title_set(w.ee, C.CString(title))
}

type EFLDriver struct {
}

func (d EFLDriver) CreateWindow(title string) ui.Window {
	engine := "wayland_shm"
	env := C.getenv(C.CString("WAYLAND_DISPLAY"))

	if env == nil {
		fmt.Println("Unable to connect to Wayland - attempting X")
		engine = "software_x11"
	}

	C.evas_init()
	C.ecore_init()
	C.ecore_evas_init()

	ee := C.ecore_evas_new(C.CString(engine), 10, 10, 300, 200, nil)
	w := &window{
		ee: ee,
	}

	w.SetTitle(title)
	if engine == "wayland_shm" {
		win := C.ecore_evas_wayland2_window_get(w.ee)
		C.ecore_wl2_window_type_set(win, C.ECORE_WL2_WINDOW_TYPE_TOPLEVEL)
	}

	C.ecore_evas_show(ee)
	return w
}

func (d EFLDriver) Run() {
	C.ecore_main_loop_begin()
}
