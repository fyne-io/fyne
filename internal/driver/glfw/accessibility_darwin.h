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
    AccessibilityRoleCheckbox,
    AccessibilityRoleContainer,
    AccessibilityRoleHeading,
    AccessibilityRoleImage,
    AccessibilityRoleLink,
    AccessibilityRoleList,
    AccessibilityRoleListItem,
    AccessibilityRoleProgressBar,
    AccessibilityRoleRadio,
    AccessibilityRoleSeparator,
    AccessibilityRoleSlider,
    AccessibilityRoleStaticText,
    AccessibilityRoleTab,
    AccessibilityRoleTabList,
    AccessibilityRoleTable,
    AccessibilityRoleTextField,
    AccessibilityRoleTree,
    AccessibilityRoleTreeItem,
    AccessibilityRoleGroup
} AccessibilityRole;

// AccessibilityActionCode mirrors fyne.AccessibleAction values.
// Keep in sync with roleToC / actionFromC in accessibility_darwin.go.
typedef enum {
    AccessibilityActionPress     = 0,
    AccessibilityActionIncrement = 1,
    AccessibilityActionDecrement = 2,
    AccessibilityActionShowMenu  = 3,
    AccessibilityActionSelect    = 4,
    AccessibilityActionSetValue  = 5
} AccessibilityActionCode;

// AccessibilityActionMask is a bitmask of supported actions.
enum {
    AccessibilityActionMaskPress     = 1 << 0,
    AccessibilityActionMaskIncrement = 1 << 1,
    AccessibilityActionMaskDecrement = 1 << 2,
    AccessibilityActionMaskShowMenu  = 1 << 3,
    AccessibilityActionMaskSelect    = 1 << 4,
    AccessibilityActionMaskSetValue  = 1 << 5
};

// AccessibilityStateMask is a bitmask of accessibility state flags.
enum {
    AccessibilityStateMaskChecked  = 1 << 0,
    AccessibilityStateMaskDisabled = 1 << 1,
    AccessibilityStateMaskExpanded = 1 << 2,
    AccessibilityStateMaskFocused  = 1 << 3,
    AccessibilityStateMaskInvalid  = 1 << 4,
    AccessibilityStateMaskRequired = 1 << 5,
    AccessibilityStateMaskSelected = 1 << 6
};

// AccessibilityActionCallback returns 1 if the action was handled, 0 otherwise.
// setValueArg is non-NULL only for AccessibilityActionSetValue.
typedef int (*AccessibilityActionCallback)(void* context, int actionCode, const char* setValueArg);

// AccessibilityContextDestroy is invoked when the Obj-C element is finally
// deallocated so the Go side can release any resources tied to the context
// (typically a runtime/cgo Handle).
typedef void (*AccessibilityContextDestroy)(void* context);

AccessibilityElementRef AccessibilityElementCreate(
    AccessibilityRole role,
    const char* title,
    const char* label,
    double x, double y, double width, double height,
    AccessibilityElementRef parent,
    AccessibilityActionCallback callback,
    void* callbackContext,
    AccessibilityContextDestroy contextDestroy
);

void AccessibilityElementSetFrame(AccessibilityElementRef elem, double x, double y, double width, double height);
void AccessibilityElementSetTitle(AccessibilityElementRef elem, const char* title);
void AccessibilityElementSetLabel(AccessibilityElementRef elem, const char* label);
void AccessibilityElementSetValue(AccessibilityElementRef elem, const char* value);
void AccessibilityElementSetEnabled(AccessibilityElementRef elem, int enabled);
void AccessibilityElementSetFocused(AccessibilityElementRef elem, int focused);
void AccessibilityElementSetStates(AccessibilityElementRef elem, int stateMask);
void AccessibilityElementSetSupportedActions(AccessibilityElementRef elem, int actionMask);

void AccessibilityElementAddChild(AccessibilityElementRef parent, AccessibilityElementRef child);
void AccessibilityElementRemoveChild(AccessibilityElementRef parent, AccessibilityElementRef child);
void AccessibilityElementSetParent(AccessibilityElementRef child, AccessibilityElementRef parent);

void AccessibilitySetTargetWindow(void* nsWindow);
void AccessibilityAttachToWindow(AccessibilityElementRef elem);
void AccessibilityPostNotification(AccessibilityElementRef elem, const char* notification);
void AccessibilityElementDestroy(AccessibilityElementRef elem);

#endif // ACCESSIBILITY_BRIDGE_H
