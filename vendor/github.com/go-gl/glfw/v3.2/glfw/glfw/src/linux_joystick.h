//========================================================================
// GLFW 3.2 Linux - www.glfw.org
//------------------------------------------------------------------------
// Copyright (c) 2014 Jonas Ådahl <jadahl@gmail.com>
//
// This software is provided 'as-is', without any express or implied
// warranty. In no event will the authors be held liable for any damages
// arising from the use of this software.
//
// Permission is granted to anyone to use this software for any purpose,
// including commercial applications, and to alter it and redistribute it
// freely, subject to the following restrictions:
//
// 1. The origin of this software must not be misrepresented; you must not
//    claim that you wrote the original software. If you use this software
//    in a product, an acknowledgment in the product documentation would
//    be appreciated but is not required.
//
// 2. Altered source versions must be plainly marked as such, and must not
//    be misrepresented as being the original software.
//
// 3. This notice may not be removed or altered from any source
//    distribution.
//
//========================================================================

#ifndef _glfw3_linux_joystick_h_
#define _glfw3_linux_joystick_h_

#include <regex.h>

#define _GLFW_PLATFORM_LIBRARY_JOYSTICK_STATE _GLFWjoylistLinux linux_js


// Linux-specific joystick data
//
typedef struct _GLFWjoystickLinux
{
    GLFWbool        present;
    int             fd;
    float*          axes;
    int             axisCount;
    unsigned char*  buttons;
    int             buttonCount;
    char*           name;
    char*           path;
} _GLFWjoystickLinux;

// Linux-specific joystick API data
//
typedef struct _GLFWjoylistLinux
{
    _GLFWjoystickLinux js[GLFW_JOYSTICK_LAST + 1];

#if defined(__linux__)
    int             inotify;
    int             watch;
    regex_t         regex;
#endif /*__linux__*/
} _GLFWjoylistLinux;


GLFWbool _glfwInitJoysticksLinux(void);
void _glfwTerminateJoysticksLinux(void);

void _glfwPollJoystickEvents(void);

#endif // _glfw3_linux_joystick_h_
