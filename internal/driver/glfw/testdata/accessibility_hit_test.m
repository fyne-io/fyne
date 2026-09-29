#import <Cocoa/Cocoa.h>
#import "accessibility_darwin.h"
#include <stdio.h>

static int expectHit(id root, NSPoint point, id expected, const char *description) {
    id hit = [root accessibilityHitTest:point];
    if (hit == expected) {
        return 0;
    }
    fprintf(stderr, "%s: expected %s, got %s\n", description,
        [[expected accessibilityTitle] UTF8String], [[hit accessibilityTitle] UTF8String]);
    return 1;
}

int main(void) {
    @autoreleasepool {
        id tabs = (id)AccessibilityElementCreate(
            AccessibilityRoleTabList, "Tabs", "", 10, 20, 300, 200, NULL, NULL, NULL, NULL);
        id group = (id)AccessibilityElementCreate(
            AccessibilityRoleContainer, "Layout", "", 20, 40, 200, 100, NULL, NULL, NULL, NULL);
        id button = (id)AccessibilityElementCreate(
            AccessibilityRoleButton, "Button", "", 30, 50, 80, 30, NULL, NULL, NULL, NULL);
        id entry = (id)AccessibilityElementCreate(
            AccessibilityRoleTextField, "Entry", "", 120, 50, 80, 30, NULL, NULL, NULL, NULL);
        AccessibilityElementAddChild(tabs, group);
        AccessibilityElementAddChild(group, button);
        AccessibilityElementAddChild(group, entry);
        int failed = 0;
        failed |= expectHit(tabs, NSMakePoint(50, 60), button, "Nested button");
        failed |= expectHit(tabs, NSMakePoint(150, 60), entry, "Nested entry");
        failed |= expectHit(tabs, NSMakePoint(25, 100), tabs, "Ignored layout background");
        failed |= expectHit(tabs, NSMakePoint(250, 180), tabs, "Tab group background");
        failed |= expectHit(tabs, NSMakePoint(0, 0), nil, "Outside root");
        AccessibilityElementSetEnabled(button, 0);
        failed |= expectHit(tabs, NSMakePoint(50, 60), button, "Disabled button remains inspectable");
        id overlay = (id)AccessibilityElementCreate(
            AccessibilityRoleButton, "Overlay", "", 30, 50, 40, 30, NULL, NULL, NULL, NULL);
        AccessibilityElementAddChild(tabs, overlay);
        failed |= expectHit(tabs, NSMakePoint(50, 60), overlay, "Last child is frontmost");
        failed |= expectHit(tabs, NSMakePoint(90, 60), button, "Uncovered child");
        AccessibilityElementDestroy(tabs);
        return failed;
    }
}
