package driver

/*
#cgo pkg-config: ecore-evas
#include <Ecore_Evas.h>

// The gateway function
void onWindowResize_cgo(Ecore_Evas *ee)
{
	void onWindowResize(Ecore_Evas*);
	onWindowResize(ee);
}
*/
import "C"
