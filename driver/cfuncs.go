package driver

/*
#cgo pkg-config: ecore-evas
#include <Ecore_Evas.h>

// Gateway callback functions

void onWindowResize_cgo(Ecore_Evas *ee)
{
	void onWindowResize(Ecore_Evas*);
	onWindowResize(ee);
}

void onWindowClose_cgo(Ecore_Evas *ee)
{
	void onWindowClose(Ecore_Evas*);
	onWindowClose(ee);
}

void onObjectMouseDown_cgo(void *data, Evas *e, Evas_Object *obj, void *event_info)
{
	void onObjectMouseDown(Evas_Object*, Evas_Event_Mouse_Down*);
	onObjectMouseDown(obj, event_info);
}
*/
import "C"
