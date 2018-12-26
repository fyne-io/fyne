//========================================================================
// GLFW 3.2 Win32 - www.glfw.org
//------------------------------------------------------------------------
// Copyright (c) 2006-2016 Camilla Berglund <elmindreda@glfw.org>
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

#ifndef _glfw3_win32_joystick_h_
#define _glfw3_win32_joystick_h_

#define _GLFW_PLATFORM_LIBRARY_JOYSTICK_STATE \
    _GLFWjoystickWin32 win32_js[GLFW_JOYSTICK_LAST + 1]

// Joystick element (axis, button or slider)
//
typedef struct _GLFWjoyobjectWin32
{
    int                     offset;
    int                     type;
} _GLFWjoyobjectWin32;

// Win32-specific per-joystick data
//
typedef struct _GLFWjoystickWin32
{
    GLFWbool                present;
    float*                  axes;
    int                     axisCount;
    unsigned char*          buttons;
    int                     buttonCount;
    _GLFWjoyobjectWin32*    objects;
    int                     objectCount;
    char*                   name;
    IDirectInputDevice8W*   device;
    DWORD                   index;
    GUID                    guid;
} _GLFWjoystickWin32;


void _glfwInitJoysticksWin32(void);
void _glfwTerminateJoysticksWin32(void);
void _glfwDetectJoystickConnectionWin32(void);
void _glfwDetectJoystickDisconnectionWin32(void);

#endif // _glfw3_win32_joystick_h_
