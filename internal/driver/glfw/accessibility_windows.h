//go:build accessibility && windows

#ifndef ACCESSIBILITY_WINDOWS_H
#define ACCESSIBILITY_WINDOWS_H

typedef enum {
    WinAccessibilityRoleButton = 0,
    WinAccessibilityRoleText,
    WinAccessibilityRoleLink,
    WinAccessibilityRoleGroup,
    WinAccessibilityRoleCheckbox,
    WinAccessibilityRoleRadio,
    WinAccessibilityRoleSlider,
    WinAccessibilityRoleProgressBar,
    WinAccessibilityRoleTextField,
    WinAccessibilityRoleList,
    WinAccessibilityRoleListItem,
    WinAccessibilityRoleTree,
    WinAccessibilityRoleTreeItem,
    WinAccessibilityRoleTable,
    WinAccessibilityRoleTab,
    WinAccessibilityRoleTabList,
    WinAccessibilityRoleImage,
    WinAccessibilityRoleHeading,
    WinAccessibilityRoleSeparator
} WinAccessibilityRole;

void WinAccessibilitySetWindow(void* hwnd);
void WinAccessibilityAddElement(const char* name, WinAccessibilityRole role,
    double x, double y, double width, double height);
void WinAccessibilityClearElements(void);
void WinAccessibilityUpdate(void);
void WinAccessibilityCleanup(void);

#endif
