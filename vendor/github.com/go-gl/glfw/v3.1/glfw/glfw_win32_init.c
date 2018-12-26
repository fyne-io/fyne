#ifdef _GLFW_WIN32

	// Bug fix for 6l/8l on Windows, see:
	//  https://github.com/go-gl/glfw3/issues/82#issuecomment-53859967
	#include <stdlib.h>
	#include <string.h>

	#ifndef strdup
	char *strdup (const char *str) {
		char *new = malloc(strlen(str) + 1);
		strcpy(new, str);
		return new;
	}
	#endif
	
	// _get_output_format is described at MSDN:
	// http://msdn.microsoft.com/en-us/library/571yb472.aspx
	//
	// MinGW32 does not define it. But MinGW-W64 does. Unlike strdup above
	// compilers don't #define _get_output_format if it exists.
	#ifndef __MINGW64__
	unsigned int _get_output_format(void) {
		return 0;
	};
	#endif

	#include "glfw/src/win32_init.c"
#endif

