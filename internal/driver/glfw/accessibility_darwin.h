//go:build accessibility && darwin

#ifndef ACCESSIBILITY_BRIDGE_H
#define ACCESSIBILITY_BRIDGE_H

#ifdef __OBJC__
#import <Cocoa/Cocoa.h>
#endif

typedef void* AccessibilityElementRef;

typedef enum {
    AccessibilityRoleWindow,
    AccessibilityRoleButton,
    AccessibilityRoleStaticText,
    AccessibilityRoleTextField,
    AccessibilityRoleCheckbox,
    AccessibilityRoleGroup
} AccessibilityRole;

typedef void (*AccessibilityActionCallback)(void* context);

// Create returns a caller-owned reference. Destroy every created element,
// including children; registering a child or root retains it independently.
AccessibilityElementRef AccessibilityElementCreate(
    AccessibilityRole role,
    const char* title,
    const char* label,
    double x, double y, double width, double height,
    AccessibilityElementRef parent,
    AccessibilityActionCallback callback,
    void* callbackContext
);

void AccessibilityElementSetFrame(AccessibilityElementRef elem, double x, double y, double width, double height);
void AccessibilityElementSetTitle(AccessibilityElementRef elem, const char* title);
void AccessibilityElementSetLabel(AccessibilityElementRef elem, const char* label);
void AccessibilityElementSetValue(AccessibilityElementRef elem, const char* value);
void AccessibilityElementSetEnabled(AccessibilityElementRef elem, int enabled);
void AccessibilityElementSetFocused(AccessibilityElementRef elem, int focused);

void AccessibilityElementAddChild(AccessibilityElementRef parent, AccessibilityElementRef child);
void AccessibilityElementRemoveChild(AccessibilityElementRef parent, AccessibilityElementRef child);
void AccessibilityElementSetParent(AccessibilityElementRef child, AccessibilityElementRef parent);

// A NULL window leaves the element unattached, for window-free native testing.
void AccessibilityAttachToWindow(AccessibilityElementRef elem, void* nsWindow);
void AccessibilityPostNotification(AccessibilityElementRef elem, const char* notification);
void AccessibilityElementDestroy(AccessibilityElementRef elem);

#endif // ACCESSIBILITY_BRIDGE_H
