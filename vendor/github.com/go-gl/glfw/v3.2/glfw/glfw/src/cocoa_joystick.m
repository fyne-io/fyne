//========================================================================
// GLFW 3.2 Cocoa - www.glfw.org
//------------------------------------------------------------------------
// Copyright (c) 2009-2016 Camilla Berglund <elmindreda@glfw.org>
// Copyright (c) 2012 Torsten Walluhn <tw@mad-cad.net>
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

#include <unistd.h>
#include <ctype.h>
#include <string.h>

#include <mach/mach.h>
#include <mach/mach_error.h>

#include <CoreFoundation/CoreFoundation.h>
#include <Kernel/IOKit/hidsystem/IOHIDUsageTables.h>


// Joystick element information
//
typedef struct _GLFWjoyelementNS
{
    IOHIDElementRef elementRef;

    long min;
    long max;

    long minReport;
    long maxReport;

} _GLFWjoyelementNS;


static void getElementsCFArrayHandler(const void* value, void* parameter);

// Adds an element to the specified joystick
//
static void addJoystickElement(_GLFWjoystickNS* js,
                               IOHIDElementRef elementRef)
{
    IOHIDElementType elementType;
    long usagePage, usage;
    CFMutableArrayRef elementsArray = NULL;

    elementType = IOHIDElementGetType(elementRef);
    usagePage = IOHIDElementGetUsagePage(elementRef);
    usage = IOHIDElementGetUsage(elementRef);

    if ((elementType != kIOHIDElementTypeInput_Axis) &&
        (elementType != kIOHIDElementTypeInput_Button) &&
        (elementType != kIOHIDElementTypeInput_Misc))
    {
        return;
    }

    switch (usagePage)
    {
        case kHIDPage_GenericDesktop:
        {
            switch (usage)
            {
                case kHIDUsage_GD_X:
                case kHIDUsage_GD_Y:
                case kHIDUsage_GD_Z:
                case kHIDUsage_GD_Rx:
                case kHIDUsage_GD_Ry:
                case kHIDUsage_GD_Rz:
                case kHIDUsage_GD_Slider:
                case kHIDUsage_GD_Dial:
                case kHIDUsage_GD_Wheel:
                    elementsArray = js->axisElements;
                    break;
                case kHIDUsage_GD_Hatswitch:
                    elementsArray = js->hatElements;
                    break;
            }

            break;
        }

        case kHIDPage_Button:
            elementsArray = js->buttonElements;
            break;
        default:
            break;
    }

    if (elementsArray)
    {
        _GLFWjoyelementNS* element = calloc(1, sizeof(_GLFWjoyelementNS));

        CFArrayAppendValue(elementsArray, element);

        element->elementRef = elementRef;

        element->minReport = IOHIDElementGetLogicalMin(elementRef);
        element->maxReport = IOHIDElementGetLogicalMax(elementRef);
    }
}

// Adds an element to the specified joystick
//
static void getElementsCFArrayHandler(const void* value, void* parameter)
{
    if (CFGetTypeID(value) == IOHIDElementGetTypeID())
    {
        addJoystickElement((_GLFWjoystickNS*) parameter,
                           (IOHIDElementRef) value);
    }
}

// Returns the value of the specified element of the specified joystick
//
static long getElementValue(_GLFWjoystickNS* js, _GLFWjoyelementNS* element)
{
    IOReturn result = kIOReturnSuccess;
    IOHIDValueRef valueRef;
    long value = 0;

    if (js && element && js->deviceRef)
    {
        result = IOHIDDeviceGetValue(js->deviceRef,
                                     element->elementRef,
                                     &valueRef);

        if (kIOReturnSuccess == result)
        {
            value = IOHIDValueGetIntegerValue(valueRef);

            // Record min and max for auto calibration
            if (value < element->minReport)
                element->minReport = value;
            if (value > element->maxReport)
                element->maxReport = value;
        }
    }

    // Auto user scale
    return value;
}

// Removes the specified joystick
//
static void removeJoystick(_GLFWjoystickNS* js)
{
    int i;

    if (!js->present)
        return;

    for (i = 0;  i < CFArrayGetCount(js->axisElements);  i++)
        free((void*) CFArrayGetValueAtIndex(js->axisElements, i));
    CFArrayRemoveAllValues(js->axisElements);
    CFRelease(js->axisElements);

    for (i = 0;  i < CFArrayGetCount(js->buttonElements);  i++)
        free((void*) CFArrayGetValueAtIndex(js->buttonElements, i));
    CFArrayRemoveAllValues(js->buttonElements);
    CFRelease(js->buttonElements);

    for (i = 0;  i < CFArrayGetCount(js->hatElements);  i++)
        free((void*) CFArrayGetValueAtIndex(js->hatElements, i));
    CFArrayRemoveAllValues(js->hatElements);
    CFRelease(js->hatElements);

    free(js->axes);
    free(js->buttons);

    memset(js, 0, sizeof(_GLFWjoystickNS));

    _glfwInputJoystickChange(js - _glfw.ns_js, GLFW_DISCONNECTED);
}

// Polls for joystick axis events and updates GLFW state
//
static GLFWbool pollJoystickAxisEvents(_GLFWjoystickNS* js)
{
    CFIndex i;

    if (!js->present)
        return GLFW_FALSE;

    for (i = 0;  i < CFArrayGetCount(js->axisElements);  i++)
    {
        _GLFWjoyelementNS* axis = (_GLFWjoyelementNS*)
            CFArrayGetValueAtIndex(js->axisElements, i);

        long value = getElementValue(js, axis);
        long readScale = axis->maxReport - axis->minReport;

        if (readScale == 0)
            js->axes[i] = value;
        else
            js->axes[i] = (2.f * (value - axis->minReport) / readScale) - 1.f;
    }

    return GLFW_TRUE;
}

// Polls for joystick button events and updates GLFW state
//
static GLFWbool pollJoystickButtonEvents(_GLFWjoystickNS* js)
{
    CFIndex i;
    int buttonIndex = 0;

    if (!js->present)
        return GLFW_FALSE;

    for (i = 0;  i < CFArrayGetCount(js->buttonElements);  i++)
    {
        _GLFWjoyelementNS* button = (_GLFWjoyelementNS*)
            CFArrayGetValueAtIndex(js->buttonElements, i);

        if (getElementValue(js, button))
            js->buttons[buttonIndex++] = GLFW_PRESS;
        else
            js->buttons[buttonIndex++] = GLFW_RELEASE;
    }

    for (i = 0;  i < CFArrayGetCount(js->hatElements);  i++)
    {
        _GLFWjoyelementNS* hat = (_GLFWjoyelementNS*)
            CFArrayGetValueAtIndex(js->hatElements, i);

        // Bit fields of button presses for each direction, including nil
        const int directions[9] = { 1, 3, 2, 6, 4, 12, 8, 9, 0 };

        long j, value = getElementValue(js, hat);
        if (value < 0 || value > 8)
            value = 8;

        for (j = 0;  j < 4;  j++)
        {
            if (directions[value] & (1 << j))
                js->buttons[buttonIndex++] = GLFW_PRESS;
            else
                js->buttons[buttonIndex++] = GLFW_RELEASE;
        }
    }

    return GLFW_TRUE;
}

// Callback for user-initiated joystick addition
//
static void matchCallback(void* context,
                          IOReturn result,
                          void* sender,
                          IOHIDDeviceRef deviceRef)
{
    _GLFWjoystickNS* js;
    int joy;

    for (joy = GLFW_JOYSTICK_1;  joy <= GLFW_JOYSTICK_LAST;  joy++)
    {
        if (_glfw.ns_js[joy].present && _glfw.ns_js[joy].deviceRef == deviceRef)
            return;
    }

    for (joy = GLFW_JOYSTICK_1;  joy <= GLFW_JOYSTICK_LAST;  joy++)
    {
        if (!_glfw.ns_js[joy].present)
            break;
    }

    if (joy > GLFW_JOYSTICK_LAST)
        return;

    js = _glfw.ns_js + joy;
    js->present = GLFW_TRUE;
    js->deviceRef = deviceRef;

    CFStringRef name = IOHIDDeviceGetProperty(deviceRef,
                                              CFSTR(kIOHIDProductKey));
    if (name)
    {
        CFStringGetCString(name,
                           js->name,
                           sizeof(js->name),
                           kCFStringEncodingUTF8);
    }
    else
        strncpy(js->name, "Unknown", sizeof(js->name));

    js->axisElements = CFArrayCreateMutable(NULL, 0, NULL);
    js->buttonElements = CFArrayCreateMutable(NULL, 0, NULL);
    js->hatElements = CFArrayCreateMutable(NULL, 0, NULL);

    CFArrayRef arrayRef = IOHIDDeviceCopyMatchingElements(deviceRef,
                                                          NULL,
                                                          kIOHIDOptionsTypeNone);
    CFRange range = { 0, CFArrayGetCount(arrayRef) };
    CFArrayApplyFunction(arrayRef,
                         range,
                         getElementsCFArrayHandler,
                         (void*) js);

    CFRelease(arrayRef);

    js->axes = calloc(CFArrayGetCount(js->axisElements), sizeof(float));
    js->buttons = calloc(CFArrayGetCount(js->buttonElements) +
                         CFArrayGetCount(js->hatElements) * 4, 1);

    _glfwInputJoystickChange(joy, GLFW_CONNECTED);
}

// Callback for user-initiated joystick removal
//
static void removeCallback(void* context,
                           IOReturn result,
                           void* sender,
                           IOHIDDeviceRef deviceRef)
{
    int joy;

    for (joy = GLFW_JOYSTICK_1;  joy <= GLFW_JOYSTICK_LAST;  joy++)
    {
        if (_glfw.ns_js[joy].deviceRef == deviceRef)
        {
            removeJoystick(_glfw.ns_js + joy);
            break;
        }
    }
}

// Creates a dictionary to match against devices with the specified usage page
// and usage
//
static CFMutableDictionaryRef createMatchingDictionary(long usagePage,
                                                       long usage)
{
    CFMutableDictionaryRef result =
        CFDictionaryCreateMutable(kCFAllocatorDefault,
                                  0,
                                  &kCFTypeDictionaryKeyCallBacks,
                                  &kCFTypeDictionaryValueCallBacks);

    if (result)
    {
        CFNumberRef pageRef = CFNumberCreate(kCFAllocatorDefault,
                                             kCFNumberLongType,
                                             &usagePage);
        if (pageRef)
        {
            CFDictionarySetValue(result,
                                 CFSTR(kIOHIDDeviceUsagePageKey),
                                 pageRef);
            CFRelease(pageRef);

            CFNumberRef usageRef = CFNumberCreate(kCFAllocatorDefault,
                                                  kCFNumberLongType,
                                                  &usage);
            if (usageRef)
            {
                CFDictionarySetValue(result,
                                     CFSTR(kIOHIDDeviceUsageKey),
                                     usageRef);
                CFRelease(usageRef);
            }
        }
    }

    return result;
}


//////////////////////////////////////////////////////////////////////////
//////                       GLFW internal API                      //////
//////////////////////////////////////////////////////////////////////////

// Initialize joystick interface
//
void _glfwInitJoysticksNS(void)
{
    CFMutableArrayRef matchingCFArrayRef;

    _glfw.ns.hidManager = IOHIDManagerCreate(kCFAllocatorDefault,
                                             kIOHIDOptionsTypeNone);

    matchingCFArrayRef = CFArrayCreateMutable(kCFAllocatorDefault,
                                              0,
                                              &kCFTypeArrayCallBacks);
    if (matchingCFArrayRef)
    {
        CFDictionaryRef matchingCFDictRef =
            createMatchingDictionary(kHIDPage_GenericDesktop,
                                     kHIDUsage_GD_Joystick);
        if (matchingCFDictRef)
        {
            CFArrayAppendValue(matchingCFArrayRef, matchingCFDictRef);
            CFRelease(matchingCFDictRef);
        }

        matchingCFDictRef = createMatchingDictionary(kHIDPage_GenericDesktop,
                                                     kHIDUsage_GD_GamePad);
        if (matchingCFDictRef)
        {
            CFArrayAppendValue(matchingCFArrayRef, matchingCFDictRef);
            CFRelease(matchingCFDictRef);
        }

        matchingCFDictRef =
            createMatchingDictionary(kHIDPage_GenericDesktop,
                                     kHIDUsage_GD_MultiAxisController);
        if (matchingCFDictRef)
        {
            CFArrayAppendValue(matchingCFArrayRef, matchingCFDictRef);
            CFRelease(matchingCFDictRef);
        }

        IOHIDManagerSetDeviceMatchingMultiple(_glfw.ns.hidManager,
                                              matchingCFArrayRef);
        CFRelease(matchingCFArrayRef);
    }

    IOHIDManagerRegisterDeviceMatchingCallback(_glfw.ns.hidManager,
                                               &matchCallback, NULL);
    IOHIDManagerRegisterDeviceRemovalCallback(_glfw.ns.hidManager,
                                              &removeCallback, NULL);

    IOHIDManagerScheduleWithRunLoop(_glfw.ns.hidManager,
                                    CFRunLoopGetMain(),
                                    kCFRunLoopDefaultMode);

    IOHIDManagerOpen(_glfw.ns.hidManager, kIOHIDOptionsTypeNone);

    // Execute the run loop once in order to register any initially-attached
    // joysticks
    CFRunLoopRunInMode(kCFRunLoopDefaultMode, 0, false);
}

// Close all opened joystick handles
//
void _glfwTerminateJoysticksNS(void)
{
    int joy;

    for (joy = GLFW_JOYSTICK_1;  joy <= GLFW_JOYSTICK_LAST;  joy++)
    {
        _GLFWjoystickNS* js = _glfw.ns_js + joy;
        removeJoystick(js);
    }

    CFRelease(_glfw.ns.hidManager);
    _glfw.ns.hidManager = NULL;
}


//////////////////////////////////////////////////////////////////////////
//////                       GLFW platform API                      //////
//////////////////////////////////////////////////////////////////////////

int _glfwPlatformJoystickPresent(int joy)
{
    _GLFWjoystickNS* js = _glfw.ns_js + joy;
    return js->present;
}

const float* _glfwPlatformGetJoystickAxes(int joy, int* count)
{
    _GLFWjoystickNS* js = _glfw.ns_js + joy;
    if (!pollJoystickAxisEvents(js))
        return NULL;

    *count = (int) CFArrayGetCount(js->axisElements);
    return js->axes;
}

const unsigned char* _glfwPlatformGetJoystickButtons(int joy, int* count)
{
    _GLFWjoystickNS* js = _glfw.ns_js + joy;
    if (!pollJoystickButtonEvents(js))
        return NULL;

    *count = (int) CFArrayGetCount(js->buttonElements) +
             (int) CFArrayGetCount(js->hatElements) * 4;
    return js->buttons;
}

const char* _glfwPlatformGetJoystickName(int joy)
{
    _GLFWjoystickNS* js = _glfw.ns_js + joy;
    if (!js->present)
        return NULL;

    return js->name;
}

