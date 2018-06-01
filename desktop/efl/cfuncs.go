package efl

/*
#cgo pkg-config: ecore-evas ecore-input
#include <Ecore.h>
#include <Ecore_Evas.h>
#include <Ecore_Input.h>

// Gateway callback functions

void onWindowResize_cgo(Ecore_Evas *ee)
{
	void onWindowResize(Ecore_Evas*);
	onWindowResize(ee);
}

void onWindowMove_cgo(Ecore_Evas *ee)
{
	void onWindowMove(Ecore_Evas*);
	onWindowMove(ee);
}

void onWindowClose_cgo(Ecore_Evas *ee)
{
	void onWindowClose(Ecore_Evas*);
	onWindowClose(ee);
}

void onWindowKeyDown_cgo(void *data, int type, void *event_info)
{
	void onWindowKeyDown(Ecore_Window, Ecore_Event_Key*);
	Ecore_Event_Key *key_ev = (Ecore_Event_Key *) event_info;
	onWindowKeyDown(key_ev->window, key_ev);
}

void onObjectMouseDown_cgo(void *data, Evas *e, Evas_Object *obj, void *event_info)
{
	void onObjectMouseDown(Evas_Object*, Evas_Event_Mouse_Down*);
	onObjectMouseDown(obj, event_info);
}

void onExit_cgo(Ecore_Event_Signal_Exit *sig)
{
	void DoQuit();
	DoQuit();
}
*/
import "C"
