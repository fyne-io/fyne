package efl

/*
#cgo pkg-config: ecore-evas
#include <Ecore_Evas.h>

// Helper functions

void evas_bridge_image_pixel_set(Evas_Object *img, int x, int y, uint col)
{
	unsigned int *data, *pixel;
	int width;

	data = evas_object_image_data_get(img, EINA_TRUE);
	evas_object_geometry_get(img, NULL, NULL, &width, NULL);

	pixel = data + (y * width + x);
	*pixel = col;
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
