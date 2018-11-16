// +build !ci

package efl

/*
#cgo pkg-config: ecore-evas ecore-input
#include <Ecore.h>
#include <Ecore_Evas.h>
#include <Ecore_Input.h>

// a callback if we need to force an immediate redraw on the main thread
static Eina_Bool
_immediate_iterate(void *data EINA_UNUSED)
{
	ecore_main_loop_iterate();

	// only tick once
	return ECORE_CALLBACK_CANCEL;
}

// included so that darwin specific code can push an immediate refresh
void
force_render()
{
	ecore_timer_add(0.001, _immediate_iterate, NULL);
}

// hook into go logging from EFL
void
log_callback(const Eina_Log_Domain *d, Eina_Log_Level level, const char *file, const char *fnc, int line, const char *fmt, void *data, va_list args)
{
	void onLogCallback(const char*, int, const char*, int, char*);

	size_t len = snprintf(NULL, 0, fmt, args);
	char  *buffer = malloc(len + 1);
	sprintf(buffer, fmt, args);

	onLogCallback(d->name, level, file, line, buffer);
	free(buffer);
}

void
setup_log() {
	eina_log_print_cb_set(log_callback, NULL);
}

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

void onWindowFocusIn_cgo(Ecore_Evas *ee)
{
	void onWindowFocusGained(Ecore_Evas*);
	onWindowFocusGained(ee);
}

void onWindowFocusOut_cgo(Ecore_Evas *ee)
{
	void onWindowFocusLost(Ecore_Evas*);
	onWindowFocusLost(ee);
}

void onWindowClose_cgo(Ecore_Evas *ee)
{
	void onWindowClose(Ecore_Evas*);
	onWindowClose(ee);
}

void onKeyDown_cgo(void *data, int type, void *event_info)
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
