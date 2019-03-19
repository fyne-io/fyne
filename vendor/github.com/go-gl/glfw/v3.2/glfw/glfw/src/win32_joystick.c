//========================================================================
// GLFW 3.1 Win32 - www.glfw.org
//------------------------------------------------------------------------
// Copyright (c) 2002-2006 Marcus Geelnard
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

#include "internal.h"

#include <math.h>

#include <initguid.h>

#define _GLFW_PRESENCE_ONLY 1
#define _GLFW_UPDATE_STATE  2

#define _GLFW_TYPE_AXIS     0
#define _GLFW_TYPE_SLIDER   1
#define _GLFW_TYPE_BUTTON   2
#define _GLFW_TYPE_POV      3

// Data produced with DirectInput device object enumeration
//
typedef struct _GLFWobjenumWin32
{
    IDirectInputDevice8W*   device;
    _GLFWjoyobjectWin32*    objects;
    int                     objectCount;
    int                     axisCount;
    int                     sliderCount;
    int                     buttonCount;
    int                     povCount;
} _GLFWobjenumWin32;

// Define only the necessary GUIDs (it's bad enough that we're exporting these)
//
DEFINE_GUID(IID_IDirectInput8W,0xbf798031,0x483a,0x4da2,0xaa,0x99,0x5d,0x64,0xed,0x36,0x97,0x00);
DEFINE_GUID(GUID_XAxis,0xa36d02e0,0xc9f3,0x11cf,0xbf,0xc7,0x44,0x45,0x53,0x54,0x00,0x00);
DEFINE_GUID(GUID_YAxis,0xa36d02e1,0xc9f3,0x11cf,0xbf,0xc7,0x44,0x45,0x53,0x54,0x00,0x00);
DEFINE_GUID(GUID_ZAxis,0xa36d02e2,0xc9f3,0x11cf,0xbf,0xc7,0x44,0x45,0x53,0x54,0x00,0x00);
DEFINE_GUID(GUID_RxAxis,0xa36d02f4,0xc9f3,0x11cf,0xbf,0xc7,0x44,0x45,0x53,0x54,0x00,0x00);
DEFINE_GUID(GUID_RyAxis,0xa36d02f5,0xc9f3,0x11cf,0xbf,0xc7,0x44,0x45,0x53,0x54,0x00,0x00);
DEFINE_GUID(GUID_RzAxis,0xa36d02e3,0xc9f3,0x11cf,0xbf,0xc7,0x44,0x45,0x53,0x54,0x00,0x00);
DEFINE_GUID(GUID_Slider,0xa36d02e4,0xc9f3,0x11cf,0xbf,0xc7,0x44,0x45,0x53,0x54,0x00,0x00);
DEFINE_GUID(GUID_Button,0xa36d02f0,0xc9f3,0x11cf,0xbf,0xc7,0x44,0x45,0x53,0x54,0x00,0x00);
DEFINE_GUID(GUID_POV,0xa36d02f2,0xc9f3,0x11cf,0xbf,0xc7,0x44,0x45,0x53,0x54,0x00,0x00);

// Object data array for our clone of c_dfDIJoystick
// Generated with https://github.com/elmindreda/c_dfDIJoystick2
//
static DIOBJECTDATAFORMAT _glfwObjectDataFormats[] =
{
    { &GUID_XAxis,DIJOFS_X,DIDFT_AXIS|DIDFT_OPTIONAL|DIDFT_ANYINSTANCE,DIDOI_ASPECTPOSITION },
    { &GUID_YAxis,DIJOFS_Y,DIDFT_AXIS|DIDFT_OPTIONAL|DIDFT_ANYINSTANCE,DIDOI_ASPECTPOSITION },
    { &GUID_ZAxis,DIJOFS_Z,DIDFT_AXIS|DIDFT_OPTIONAL|DIDFT_ANYINSTANCE,DIDOI_ASPECTPOSITION },
    { &GUID_RxAxis,DIJOFS_RX,DIDFT_AXIS|DIDFT_OPTIONAL|DIDFT_ANYINSTANCE,DIDOI_ASPECTPOSITION },
    { &GUID_RyAxis,DIJOFS_RY,DIDFT_AXIS|DIDFT_OPTIONAL|DIDFT_ANYINSTANCE,DIDOI_ASPECTPOSITION },
    { &GUID_RzAxis,DIJOFS_RZ,DIDFT_AXIS|DIDFT_OPTIONAL|DIDFT_ANYINSTANCE,DIDOI_ASPECTPOSITION },
    { &GUID_Slider,DIJOFS_SLIDER(0),DIDFT_AXIS|DIDFT_OPTIONAL|DIDFT_ANYINSTANCE,DIDOI_ASPECTPOSITION },
    { &GUID_Slider,DIJOFS_SLIDER(1),DIDFT_AXIS|DIDFT_OPTIONAL|DIDFT_ANYINSTANCE,DIDOI_ASPECTPOSITION },
    { &GUID_POV,DIJOFS_POV(0),DIDFT_POV|DIDFT_OPTIONAL|DIDFT_ANYINSTANCE,0 },
    { &GUID_POV,DIJOFS_POV(1),DIDFT_POV|DIDFT_OPTIONAL|DIDFT_ANYINSTANCE,0 },
    { &GUID_POV,DIJOFS_POV(2),DIDFT_POV|DIDFT_OPTIONAL|DIDFT_ANYINSTANCE,0 },
    { &GUID_POV,DIJOFS_POV(3),DIDFT_POV|DIDFT_OPTIONAL|DIDFT_ANYINSTANCE,0 },
    { NULL,DIJOFS_BUTTON(0),DIDFT_BUTTON|DIDFT_OPTIONAL|DIDFT_ANYINSTANCE,0 },
    { NULL,DIJOFS_BUTTON(1),DIDFT_BUTTON|DIDFT_OPTIONAL|DIDFT_ANYINSTANCE,0 },
    { NULL,DIJOFS_BUTTON(2),DIDFT_BUTTON|DIDFT_OPTIONAL|DIDFT_ANYINSTANCE,0 },
    { NULL,DIJOFS_BUTTON(3),DIDFT_BUTTON|DIDFT_OPTIONAL|DIDFT_ANYINSTANCE,0 },
    { NULL,DIJOFS_BUTTON(4),DIDFT_BUTTON|DIDFT_OPTIONAL|DIDFT_ANYINSTANCE,0 },
    { NULL,DIJOFS_BUTTON(5),DIDFT_BUTTON|DIDFT_OPTIONAL|DIDFT_ANYINSTANCE,0 },
    { NULL,DIJOFS_BUTTON(6),DIDFT_BUTTON|DIDFT_OPTIONAL|DIDFT_ANYINSTANCE,0 },
    { NULL,DIJOFS_BUTTON(7),DIDFT_BUTTON|DIDFT_OPTIONAL|DIDFT_ANYINSTANCE,0 },
    { NULL,DIJOFS_BUTTON(8),DIDFT_BUTTON|DIDFT_OPTIONAL|DIDFT_ANYINSTANCE,0 },
    { NULL,DIJOFS_BUTTON(9),DIDFT_BUTTON|DIDFT_OPTIONAL|DIDFT_ANYINSTANCE,0 },
    { NULL,DIJOFS_BUTTON(10),DIDFT_BUTTON|DIDFT_OPTIONAL|DIDFT_ANYINSTANCE,0 },
    { NULL,DIJOFS_BUTTON(11),DIDFT_BUTTON|DIDFT_OPTIONAL|DIDFT_ANYINSTANCE,0 },
    { NULL,DIJOFS_BUTTON(12),DIDFT_BUTTON|DIDFT_OPTIONAL|DIDFT_ANYINSTANCE,0 },
    { NULL,DIJOFS_BUTTON(13),DIDFT_BUTTON|DIDFT_OPTIONAL|DIDFT_ANYINSTANCE,0 },
    { NULL,DIJOFS_BUTTON(14),DIDFT_BUTTON|DIDFT_OPTIONAL|DIDFT_ANYINSTANCE,0 },
    { NULL,DIJOFS_BUTTON(15),DIDFT_BUTTON|DIDFT_OPTIONAL|DIDFT_ANYINSTANCE,0 },
    { NULL,DIJOFS_BUTTON(16),DIDFT_BUTTON|DIDFT_OPTIONAL|DIDFT_ANYINSTANCE,0 },
    { NULL,DIJOFS_BUTTON(17),DIDFT_BUTTON|DIDFT_OPTIONAL|DIDFT_ANYINSTANCE,0 },
    { NULL,DIJOFS_BUTTON(18),DIDFT_BUTTON|DIDFT_OPTIONAL|DIDFT_ANYINSTANCE,0 },
    { NULL,DIJOFS_BUTTON(19),DIDFT_BUTTON|DIDFT_OPTIONAL|DIDFT_ANYINSTANCE,0 },
    { NULL,DIJOFS_BUTTON(20),DIDFT_BUTTON|DIDFT_OPTIONAL|DIDFT_ANYINSTANCE,0 },
    { NULL,DIJOFS_BUTTON(21),DIDFT_BUTTON|DIDFT_OPTIONAL|DIDFT_ANYINSTANCE,0 },
    { NULL,DIJOFS_BUTTON(22),DIDFT_BUTTON|DIDFT_OPTIONAL|DIDFT_ANYINSTANCE,0 },
    { NULL,DIJOFS_BUTTON(23),DIDFT_BUTTON|DIDFT_OPTIONAL|DIDFT_ANYINSTANCE,0 },
    { NULL,DIJOFS_BUTTON(24),DIDFT_BUTTON|DIDFT_OPTIONAL|DIDFT_ANYINSTANCE,0 },
    { NULL,DIJOFS_BUTTON(25),DIDFT_BUTTON|DIDFT_OPTIONAL|DIDFT_ANYINSTANCE,0 },
    { NULL,DIJOFS_BUTTON(26),DIDFT_BUTTON|DIDFT_OPTIONAL|DIDFT_ANYINSTANCE,0 },
    { NULL,DIJOFS_BUTTON(27),DIDFT_BUTTON|DIDFT_OPTIONAL|DIDFT_ANYINSTANCE,0 },
    { NULL,DIJOFS_BUTTON(28),DIDFT_BUTTON|DIDFT_OPTIONAL|DIDFT_ANYINSTANCE,0 },
    { NULL,DIJOFS_BUTTON(29),DIDFT_BUTTON|DIDFT_OPTIONAL|DIDFT_ANYINSTANCE,0 },
    { NULL,DIJOFS_BUTTON(30),DIDFT_BUTTON|DIDFT_OPTIONAL|DIDFT_ANYINSTANCE,0 },
    { NULL,DIJOFS_BUTTON(31),DIDFT_BUTTON|DIDFT_OPTIONAL|DIDFT_ANYINSTANCE,0 },
};

// Our clone of c_dfDIJoystick
//
static const DIDATAFORMAT _glfwDataFormat =
{
    sizeof(DIDATAFORMAT),
    sizeof(DIOBJECTDATAFORMAT),
    DIDFT_ABSAXIS,
    sizeof(DIJOYSTATE),
    sizeof(_glfwObjectDataFormats) / sizeof(DIOBJECTDATAFORMAT),
    _glfwObjectDataFormats
};

// Returns a description fitting the specified XInput capabilities
//
static const char* getDeviceDescription(const XINPUT_CAPABILITIES* xic)
{
    switch (xic->SubType)
    {
        case XINPUT_DEVSUBTYPE_WHEEL:
            return "XInput Wheel";
        case XINPUT_DEVSUBTYPE_ARCADE_STICK:
            return "XInput Arcade Stick";
        case XINPUT_DEVSUBTYPE_FLIGHT_STICK:
            return "XInput Flight Stick";
        case XINPUT_DEVSUBTYPE_DANCE_PAD:
            return "XInput Dance Pad";
        case XINPUT_DEVSUBTYPE_GUITAR:
            return "XInput Guitar";
        case XINPUT_DEVSUBTYPE_DRUM_KIT:
            return "XInput Drum Kit";
        case XINPUT_DEVSUBTYPE_GAMEPAD:
        {
            if (xic->Flags & XINPUT_CAPS_WIRELESS)
                return "Wireless Xbox 360 Controller";
            else
                return "Xbox 360 Controller";
        }
    }

    return "Unknown XInput Device";
}

// Lexically compare device objects
//
static int compareJoystickObjects(const void* first, const void* second)
{
    const _GLFWjoyobjectWin32* fo = first;
    const _GLFWjoyobjectWin32* so = second;

    if (fo->type != so->type)
        return fo->type - so->type;

    return fo->offset - so->offset;
}

// Checks whether the specified device supports XInput
// Technique from FDInputJoystickManager::IsXInputDeviceFast in ZDoom
//
static GLFWbool supportsXInput(const GUID* guid)
{
    UINT i, count = 0;
    RAWINPUTDEVICELIST* ridl;
    GLFWbool result = GLFW_FALSE;

    if (GetRawInputDeviceList(NULL, &count, sizeof(RAWINPUTDEVICELIST)) != 0)
        return GLFW_FALSE;

    ridl = calloc(count, sizeof(RAWINPUTDEVICELIST));

    if (GetRawInputDeviceList(ridl, &count, sizeof(RAWINPUTDEVICELIST)) == (UINT) -1)
    {
        free(ridl);
        return GLFW_FALSE;
    }

    for (i = 0;  i < count;  i++)
    {
        RID_DEVICE_INFO rdi;
        char name[256];
        UINT size;

        if (ridl[i].dwType != RIM_TYPEHID)
            continue;

        ZeroMemory(&rdi, sizeof(rdi));
        rdi.cbSize = sizeof(rdi);
        size = sizeof(rdi);

        if ((INT) GetRawInputDeviceInfoA(ridl[i].hDevice,
                                         RIDI_DEVICEINFO,
                                         &rdi, &size) == -1)
        {
            continue;
        }

        if (MAKELONG(rdi.hid.dwVendorId, rdi.hid.dwProductId) != (LONG) guid->Data1)
            continue;

        memset(name, 0, sizeof(name));
        size = sizeof(name);

        if ((INT) GetRawInputDeviceInfoA(ridl[i].hDevice,
                                         RIDI_DEVICENAME,
                                         name, &size) == -1)
        {
            break;
        }

        name[sizeof(name) - 1] = '\0';
        if (strstr(name, "IG_"))
        {
            result = GLFW_TRUE;
            break;
        }
    }

    free(ridl);
    return result;
}

// Frees all resources associated with the specified joystick
//
static void closeJoystick(_GLFWjoystickWin32* js)
{
    if (js->device)
    {
        IDirectInputDevice8_Unacquire(js->device);
        IDirectInputDevice8_Release(js->device);
    }

    free(js->name);
    free(js->axes);
    free(js->buttons);
    free(js->objects);
    memset(js, 0, sizeof(_GLFWjoystickWin32));

    _glfwInputJoystickChange((int) (js - _glfw.win32_js), GLFW_DISCONNECTED);
}

// DirectInput device object enumeration callback
// Insights gleaned from SDL2
//
static BOOL CALLBACK deviceObjectCallback(const DIDEVICEOBJECTINSTANCEW* doi,
                                          void* user)
{
    _GLFWobjenumWin32* data = user;
    _GLFWjoyobjectWin32* object = data->objects + data->objectCount;

    if (DIDFT_GETTYPE(doi->dwType) & DIDFT_AXIS)
    {
        DIPROPRANGE dipr;

        if (memcmp(&doi->guidType, &GUID_Slider, sizeof(GUID)) == 0)
            object->offset = DIJOFS_SLIDER(data->sliderCount);
        else if (memcmp(&doi->guidType, &GUID_XAxis, sizeof(GUID)) == 0)
            object->offset = DIJOFS_X;
        else if (memcmp(&doi->guidType, &GUID_YAxis, sizeof(GUID)) == 0)
            object->offset = DIJOFS_Y;
        else if (memcmp(&doi->guidType, &GUID_ZAxis, sizeof(GUID)) == 0)
            object->offset = DIJOFS_Z;
        else if (memcmp(&doi->guidType, &GUID_RxAxis, sizeof(GUID)) == 0)
            object->offset = DIJOFS_RX;
        else if (memcmp(&doi->guidType, &GUID_RyAxis, sizeof(GUID)) == 0)
            object->offset = DIJOFS_RY;
        else if (memcmp(&doi->guidType, &GUID_RzAxis, sizeof(GUID)) == 0)
            object->offset = DIJOFS_RZ;
        else
            return DIENUM_CONTINUE;

        ZeroMemory(&dipr, sizeof(dipr));
        dipr.diph.dwSize = sizeof(dipr);
        dipr.diph.dwHeaderSize = sizeof(dipr.diph);
        dipr.diph.dwObj = doi->dwType;
        dipr.diph.dwHow = DIPH_BYID;
        dipr.lMin = -32768;
        dipr.lMax =  32767;

        if (FAILED(IDirectInputDevice8_SetProperty(data->device,
                                                   DIPROP_RANGE,
                                                   &dipr.diph)))
        {
            return DIENUM_CONTINUE;
        }

        if (memcmp(&doi->guidType, &GUID_Slider, sizeof(GUID)) == 0)
        {
            object->type = _GLFW_TYPE_SLIDER;
            data->sliderCount++;
        }
        else
        {
            object->type = _GLFW_TYPE_AXIS;
            data->axisCount++;
        }
    }
    else if (DIDFT_GETTYPE(doi->dwType) & DIDFT_BUTTON)
    {
        object->offset = DIJOFS_BUTTON(data->buttonCount);
        object->type = _GLFW_TYPE_BUTTON;
        data->buttonCount++;
    }
    else if (DIDFT_GETTYPE(doi->dwType) & DIDFT_POV)
    {
        object->offset = DIJOFS_POV(data->povCount);
        object->type = _GLFW_TYPE_POV;
        data->povCount++;
    }

    data->objectCount++;
    return DIENUM_CONTINUE;
}

// DirectInput device enumeration callback
//
static BOOL CALLBACK deviceCallback(const DIDEVICEINSTANCE* di, void* user)
{
    int joy = 0;
    DIDEVCAPS dc;
    DIPROPDWORD dipd;
    IDirectInputDevice8* device;
    _GLFWobjenumWin32 data;
    _GLFWjoystickWin32* js;

    for (joy = GLFW_JOYSTICK_1;  joy <= GLFW_JOYSTICK_LAST;  joy++)
    {
        if (memcmp(&_glfw.win32_js[joy].guid, &di->guidInstance, sizeof(GUID)) == 0)
            return DIENUM_CONTINUE;
    }

    for (joy = GLFW_JOYSTICK_1;  joy <= GLFW_JOYSTICK_LAST;  joy++)
    {
        if (!_glfw.win32_js[joy].present)
            break;
    }

    if (joy > GLFW_JOYSTICK_LAST)
        return DIENUM_STOP;

    if (supportsXInput(&di->guidProduct))
        return DIENUM_CONTINUE;

    if (FAILED(IDirectInput8_CreateDevice(_glfw.win32.dinput8.api,
                                          &di->guidInstance,
                                          &device,
                                          NULL)))
    {
        _glfwInputError(GLFW_PLATFORM_ERROR, "DI: Failed to create device");
        return DIENUM_CONTINUE;
    }

    if (FAILED(IDirectInputDevice8_SetDataFormat(device, &_glfwDataFormat)))
    {
        _glfwInputError(GLFW_PLATFORM_ERROR,
                        "DI: Failed to set device data format");

        IDirectInputDevice8_Release(device);
        return DIENUM_CONTINUE;
    }

    ZeroMemory(&dc, sizeof(dc));
    dc.dwSize = sizeof(dc);

    if (FAILED(IDirectInputDevice8_GetCapabilities(device, &dc)))
    {
        _glfwInputError(GLFW_PLATFORM_ERROR,
                        "DI: Failed to query device capabilities");

        IDirectInputDevice8_Release(device);
        return DIENUM_CONTINUE;
    }

    ZeroMemory(&dipd, sizeof(dipd));
    dipd.diph.dwSize = sizeof(dipd);
    dipd.diph.dwHeaderSize = sizeof(dipd.diph);
    dipd.diph.dwHow = DIPH_DEVICE;
    dipd.dwData = DIPROPAXISMODE_ABS;

    if (FAILED(IDirectInputDevice8_SetProperty(device,
                                               DIPROP_AXISMODE,
                                               &dipd.diph)))
    {
        _glfwInputError(GLFW_PLATFORM_ERROR,
                        "DI: Failed to set device axis mode");

        IDirectInputDevice8_Release(device);
        return DIENUM_CONTINUE;
    }

    memset(&data, 0, sizeof(data));
    data.device = device;
    data.objects = calloc(dc.dwAxes + dc.dwButtons + dc.dwPOVs,
                          sizeof(_GLFWjoyobjectWin32));

    if (FAILED(IDirectInputDevice8_EnumObjects(device,
                                               deviceObjectCallback,
                                               &data,
                                               DIDFT_AXIS | DIDFT_BUTTON | DIDFT_POV)))
    {
        _glfwInputError(GLFW_PLATFORM_ERROR,
                        "DI: Failed to enumerate device objects");

        IDirectInputDevice8_Release(device);
        free(data.objects);
        return DIENUM_CONTINUE;
    }

    qsort(data.objects, data.objectCount,
          sizeof(_GLFWjoyobjectWin32),
          compareJoystickObjects);

    js = _glfw.win32_js + joy;
    js->device = device;
    js->guid = di->guidInstance;
    js->axisCount = data.axisCount + data.sliderCount;
    js->axes = calloc(js->axisCount, sizeof(float));
    js->buttonCount += data.buttonCount + data.povCount * 4;
    js->buttons = calloc(js->buttonCount, 1);
    js->objects = data.objects;
    js->objectCount = data.objectCount;
    js->name = _glfwCreateUTF8FromWideStringWin32(di->tszInstanceName);
    js->present = GLFW_TRUE;

    _glfwInputJoystickChange(joy, GLFW_CONNECTED);
    return DIENUM_CONTINUE;
}

// Attempt to open the specified joystick device
// TODO: Pack state arrays for non-gamepad devices
//
static GLFWbool openXinputDevice(DWORD index)
{
    int joy;
    XINPUT_CAPABILITIES xic;
    _GLFWjoystickWin32* js;

    for (joy = GLFW_JOYSTICK_1;  joy <= GLFW_JOYSTICK_LAST;  joy++)
    {
        if (_glfw.win32_js[joy].present &&
            _glfw.win32_js[joy].device == NULL &&
            _glfw.win32_js[joy].index == index)
        {
            return GLFW_FALSE;
        }
    }

    for (joy = GLFW_JOYSTICK_1;  joy <= GLFW_JOYSTICK_LAST;  joy++)
    {
        if (!_glfw.win32_js[joy].present)
            break;
    }

    if (joy > GLFW_JOYSTICK_LAST)
        return GLFW_FALSE;

    if (_glfw_XInputGetCapabilities(index, 0, &xic) != ERROR_SUCCESS)
        return GLFW_FALSE;

    js = _glfw.win32_js + joy;
    js->axisCount = 6;
    js->axes = calloc(js->axisCount, sizeof(float));
    js->buttonCount = 14;
    js->buttons = calloc(js->buttonCount, 1);
    js->present = GLFW_TRUE;
    js->name = strdup(getDeviceDescription(&xic));
    js->index = index;

    _glfwInputJoystickChange(joy, GLFW_CONNECTED);

    return GLFW_TRUE;
}

// Polls for and processes events the specified joystick
//
static GLFWbool pollJoystickState(_GLFWjoystickWin32* js, int mode)
{
    if (!js->present)
        return GLFW_FALSE;

    if (js->device)
    {
        int i, j, ai = 0, bi = 0;
        HRESULT result;
        DIJOYSTATE state;

        IDirectInputDevice8_Poll(js->device);
        result = IDirectInputDevice8_GetDeviceState(js->device,
                                                    sizeof(state),
                                                    &state);
        if (result == DIERR_NOTACQUIRED || result == DIERR_INPUTLOST)
        {
            IDirectInputDevice8_Acquire(js->device);
            IDirectInputDevice8_Poll(js->device);
            result = IDirectInputDevice8_GetDeviceState(js->device,
                                                        sizeof(state),
                                                        &state);
        }

        if (FAILED(result))
        {
            closeJoystick(js);
            return GLFW_FALSE;
        }

        if (mode == _GLFW_PRESENCE_ONLY)
            return GLFW_TRUE;

        for (i = 0;  i < js->objectCount;  i++)
        {
            const void* data = (char*) &state + js->objects[i].offset;

            switch (js->objects[i].type)
            {
                case _GLFW_TYPE_AXIS:
                case _GLFW_TYPE_SLIDER:
                {
                    js->axes[ai++] = (*((LONG*) data) + 0.5f) / 32767.5f;
                    break;
                }

                case _GLFW_TYPE_BUTTON:
                {
                    if (*((BYTE*) data) & 0x80)
                        js->buttons[bi++] = GLFW_PRESS;
                    else
                        js->buttons[bi++] = GLFW_RELEASE;

                    break;
                }

                case _GLFW_TYPE_POV:
                {
                    const int directions[9] = { 1, 3, 2, 6, 4, 12, 8, 9, 0 };
                    // Screams of horror are appropriate at this point
                    int value = LOWORD(*(DWORD*) data) / (45 * DI_DEGREES);
                    if (value < 0 || value > 8)
                        value = 8;

                    for (j = 0;  j < 4;  j++)
                    {
                        if (directions[value] & (1 << j))
                            js->buttons[bi++] = GLFW_PRESS;
                        else
                            js->buttons[bi++] = GLFW_RELEASE;
                    }

                    break;
                }
            }
        }

        return GLFW_TRUE;
    }
    else
    {
        int i;
        DWORD result;
        XINPUT_STATE xis;
        const WORD buttons[14] =
        {
            XINPUT_GAMEPAD_A,
            XINPUT_GAMEPAD_B,
            XINPUT_GAMEPAD_X,
            XINPUT_GAMEPAD_Y,
            XINPUT_GAMEPAD_LEFT_SHOULDER,
            XINPUT_GAMEPAD_RIGHT_SHOULDER,
            XINPUT_GAMEPAD_BACK,
            XINPUT_GAMEPAD_START,
            XINPUT_GAMEPAD_LEFT_THUMB,
            XINPUT_GAMEPAD_RIGHT_THUMB,
            XINPUT_GAMEPAD_DPAD_UP,
            XINPUT_GAMEPAD_DPAD_RIGHT,
            XINPUT_GAMEPAD_DPAD_DOWN,
            XINPUT_GAMEPAD_DPAD_LEFT
        };

        result = _glfw_XInputGetState(js->index, &xis);
        if (result != ERROR_SUCCESS)
        {
            if (result == ERROR_DEVICE_NOT_CONNECTED)
                closeJoystick(js);

            return GLFW_FALSE;
        }

        if (mode == _GLFW_PRESENCE_ONLY)
            return GLFW_TRUE;

        if (sqrt((double) (xis.Gamepad.sThumbLX * xis.Gamepad.sThumbLX +
                           xis.Gamepad.sThumbLY * xis.Gamepad.sThumbLY)) >
            (double) XINPUT_GAMEPAD_LEFT_THUMB_DEADZONE)
        {
            js->axes[0] = (xis.Gamepad.sThumbLX + 0.5f) / 32767.f;
            js->axes[1] = (xis.Gamepad.sThumbLY + 0.5f) / 32767.f;
        }
        else
        {
            js->axes[0] = 0.f;
            js->axes[1] = 0.f;
        }

        if (sqrt((double) (xis.Gamepad.sThumbRX * xis.Gamepad.sThumbRX +
                           xis.Gamepad.sThumbRY * xis.Gamepad.sThumbRY)) >
            (double) XINPUT_GAMEPAD_RIGHT_THUMB_DEADZONE)
        {
            js->axes[2] = (xis.Gamepad.sThumbRX + 0.5f) / 32767.f;
            js->axes[3] = (xis.Gamepad.sThumbRY + 0.5f) / 32767.f;
        }
        else
        {
            js->axes[2] = 0.f;
            js->axes[3] = 0.f;
        }

        if (xis.Gamepad.bLeftTrigger > XINPUT_GAMEPAD_TRIGGER_THRESHOLD)
            js->axes[4] = xis.Gamepad.bLeftTrigger / 127.5f - 1.f;
        else
            js->axes[4] = -1.f;

        if (xis.Gamepad.bRightTrigger > XINPUT_GAMEPAD_TRIGGER_THRESHOLD)
            js->axes[5] = xis.Gamepad.bRightTrigger / 127.5f - 1.f;
        else
            js->axes[5] = -1.f;

        for (i = 0;  i < 14;  i++)
            js->buttons[i] = (xis.Gamepad.wButtons & buttons[i]) ? 1 : 0;

        return GLFW_TRUE;
    }
}


//////////////////////////////////////////////////////////////////////////
//////                       GLFW internal API                      //////
//////////////////////////////////////////////////////////////////////////

// Initialize joystick interface
//
void _glfwInitJoysticksWin32(void)
{
    if (_glfw.win32.dinput8.instance)
    {
        if (FAILED(_glfw_DirectInput8Create(GetModuleHandle(NULL),
                                            DIRECTINPUT_VERSION,
                                            &IID_IDirectInput8W,
                                            (void**) &_glfw.win32.dinput8.api,
                                            NULL)))
        {
            _glfwInputError(GLFW_PLATFORM_ERROR,
                            "DI: Failed to create interface");
        }
    }

    _glfwDetectJoystickConnectionWin32();
}

// Close all opened joystick handles
//
void _glfwTerminateJoysticksWin32(void)
{
    int joy;

    for (joy = GLFW_JOYSTICK_1;  joy <= GLFW_JOYSTICK_LAST;  joy++)
        closeJoystick(_glfw.win32_js + joy);

    if (_glfw.win32.dinput8.api)
        IDirectInput8_Release(_glfw.win32.dinput8.api);
}

// Checks for new joysticks after DBT_DEVICEARRIVAL
//
void _glfwDetectJoystickConnectionWin32(void)
{
    if (_glfw.win32.xinput.instance)
    {
        DWORD i;

        for (i = 0;  i < XUSER_MAX_COUNT;  i++)
            openXinputDevice(i);
    }

    if (_glfw.win32.dinput8.api)
    {
        if (FAILED(IDirectInput8_EnumDevices(_glfw.win32.dinput8.api,
                                             DI8DEVCLASS_GAMECTRL,
                                             deviceCallback,
                                             NULL,
                                             DIEDFL_ALLDEVICES)))
        {
            _glfwInputError(GLFW_PLATFORM_ERROR,
                            "Failed to enumerate DirectInput8 devices");
            return;
        }
    }
}

// Checks for joystick disconnection after DBT_DEVICEREMOVECOMPLETE
//
void _glfwDetectJoystickDisconnectionWin32(void)
{
    int joy;

    for (joy = GLFW_JOYSTICK_1;  joy <= GLFW_JOYSTICK_LAST;  joy++)
        pollJoystickState(_glfw.win32_js + joy, _GLFW_PRESENCE_ONLY);
}


//////////////////////////////////////////////////////////////////////////
//////                       GLFW platform API                      //////
//////////////////////////////////////////////////////////////////////////

int _glfwPlatformJoystickPresent(int joy)
{
    _GLFWjoystickWin32* js = _glfw.win32_js + joy;
    return pollJoystickState(js, _GLFW_PRESENCE_ONLY);
}

const float* _glfwPlatformGetJoystickAxes(int joy, int* count)
{
    _GLFWjoystickWin32* js = _glfw.win32_js + joy;
    if (!pollJoystickState(js, _GLFW_UPDATE_STATE))
        return NULL;

    *count = js->axisCount;
    return js->axes;
}

const unsigned char* _glfwPlatformGetJoystickButtons(int joy, int* count)
{
    _GLFWjoystickWin32* js = _glfw.win32_js + joy;
    if (!pollJoystickState(js, _GLFW_UPDATE_STATE))
        return NULL;

    *count = js->buttonCount;
    return js->buttons;
}

const char* _glfwPlatformGetJoystickName(int joy)
{
    _GLFWjoystickWin32* js = _glfw.win32_js + joy;
    if (!pollJoystickState(js, _GLFW_PRESENCE_ONLY))
        return NULL;

    return js->name;
}

