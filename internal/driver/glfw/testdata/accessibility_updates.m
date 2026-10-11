#import <Cocoa/Cocoa.h>
#import <objc/runtime.h>
#import "accessibility_darwin.h"
#include <stdio.h>
#include <string.h>

#define CHECK(condition) do { \
    if (!(condition)) { \
        fprintf(stderr, "%s:%d: %s\n", __FILE__, __LINE__, #condition); \
        return 1; \
    } \
} while (0)

@interface ProbeBaseView : NSView
@property (nonatomic, assign) NSWindow* owner;
@end

@implementation ProbeBaseView
- (NSArray*)accessibilityChildren {
    return @[];
}
- (NSWindow*)window {
    return self.owner;
}
- (NSPoint)convertPoint:(NSPoint)point toView:(NSView*)view {
    return point;
}
@end

@interface ProbeInheritedView : ProbeBaseView
@end

@implementation ProbeInheritedView
@end

@interface ProbeOwnView : ProbeBaseView
@end

@implementation ProbeOwnView
- (NSArray*)accessibilityChildren {
    return [super accessibilityChildren];
}
@end

@interface ProbeWindow : NSObject
@property (nonatomic, retain) ProbeBaseView* contentView;
@property (nonatomic, assign) NSPoint screenOrigin;
@end

@implementation ProbeWindow
- (void)dealloc {
    _contentView.owner = nil;
    [_contentView release];
    [super dealloc];
}
- (NSPoint)convertPointToScreen:(NSPoint)point {
    return NSMakePoint(point.x + self.screenOrigin.x, point.y + self.screenOrigin.y);
}
@end

static char deallocationProbeKey;
static int deallocations;

@interface DeallocationProbe : NSObject
@end
@implementation DeallocationProbe
- (void)dealloc {
    deallocations++;
    [super dealloc];
}
@end

static id newElement(AccessibilityRole role, const char* label, id parent) {
    id element = (id)AccessibilityElementCreate(role, label, label, 5, 7, 100, 30, parent, NULL, NULL);
    DeallocationProbe* probe = [[DeallocationProbe alloc] init];
    objc_setAssociatedObject(element, &deallocationProbeKey, probe, OBJC_ASSOCIATION_RETAIN_NONATOMIC);
    [probe release];
    return element;
}

static ProbeWindow* newWindow(Class viewClass, NSPoint origin) {
    ProbeWindow* window = [[ProbeWindow alloc] init];
    ProbeBaseView* view = [[viewClass alloc] initWithFrame:NSMakeRect(0, 0, 400, 300)];
    window.contentView = view;
    view.owner = (NSWindow*)window;
    window.screenOrigin = origin;
    [view release];
    return window;
}

static int testAttributes(void) {
    id element = newElement(AccessibilityRoleButton, "Button", nil);
    AccessibilityElementSetTitle(element, "Updated button");
    AccessibilityElementSetLabel(element, "Updated label");
    AccessibilityElementSetFrame(element, 20, 40, 120, 50);
    CHECK([[element accessibilityTitle] isEqualToString:@"Updated button"]);
    CHECK([[element accessibilityLabel] isEqualToString:@"Updated label"]);
    CHECK(NSEqualRects([[element valueForKey:@"localFrame"] rectValue], NSMakeRect(20, 40, 120, 50)));
    CHECK(NSEqualRects([element accessibilityFrame], NSMakeRect(20, 40, 120, 50)));
    CHECK([element accessibilityParent] == nil);
    AccessibilityElementDestroy(element);
    CHECK(deallocations == 1);
    return 0;
}

static int testOwnership(void) {
    id parent = newElement(AccessibilityRoleGroup, "Group", nil);
    id child = newElement(AccessibilityRoleButton, "Button", nil);
    AccessibilityElementAddChild(parent, child);
    AccessibilityElementRemoveChild(parent, child);
    CHECK(deallocations == 0);
    CHECK([child accessibilityParent] == nil);
    AccessibilityElementAddChild(parent, child);
    AccessibilityElementDestroy(parent);
    CHECK(deallocations == 1);
    CHECK([child accessibilityParent] == nil);
    CHECK([[child accessibilityRole] isEqualToString:NSAccessibilityButtonRole]);
    AccessibilityElementDestroy(child);
    CHECK(deallocations == 2);

    parent = newElement(AccessibilityRoleGroup, "Group", nil);
    child = newElement(AccessibilityRoleButton, "Button", parent);
    AccessibilityElementDestroy(child);
    AccessibilityElementDestroy(parent);
    CHECK(deallocations == 4);

    parent = newElement(AccessibilityRoleGroup, "First group", nil);
    id secondParent = newElement(AccessibilityRoleGroup, "Second group", nil);
    child = newElement(AccessibilityRoleButton, "Button", nil);
    AccessibilityElementAddChild(parent, child);
    AccessibilityElementAddChild(secondParent, child);
    CHECK([[parent valueForKey:@"children"] count] == 0);
    CHECK([child accessibilityParent] == secondParent);
    AccessibilityElementDestroy(child);
    AccessibilityElementDestroy(parent);
    AccessibilityElementDestroy(secondParent);
    CHECK(deallocations == 7);
    return 0;
}

static int testViewOverride(Class viewClass) {
    ProbeWindow* window = newWindow(viewClass, NSZeroPoint);
    id element = newElement(AccessibilityRoleButton, "Button", nil);
    AccessibilityAttachToWindow(element, window);
    AccessibilityAttachToWindow(element, window);
    CHECK([[window.contentView accessibilityChildren] count] == 1);
    CHECK([[window.contentView accessibilityChildren] firstObject] == element);
    ProbeWindow* unrelated = newWindow(viewClass, NSZeroPoint);
    CHECK([[unrelated.contentView accessibilityChildren] count] == 0);
    AccessibilityElementDestroy(element);
    CHECK([[window.contentView accessibilityChildren] count] == 0);
    [window release];
    [unrelated release];
    return 0;
}

static int testWindows(BOOL mixedClasses) {
    ProbeWindow* firstWindow = newWindow([ProbeOwnView class], NSMakePoint(10, 20));
    ProbeWindow* secondWindow = newWindow(mixedClasses ? [ProbeInheritedView class] : [ProbeOwnView class],
                                        NSMakePoint(500, 600));
    id first = newElement(AccessibilityRoleGroup, "First", nil);
    id second = newElement(AccessibilityRoleButton, "Second", nil);
    id child = newElement(AccessibilityRoleButton, "Child", nil);
    AccessibilityElementAddChild(first, child);
    AccessibilityAttachToWindow(first, firstWindow);
    AccessibilityAttachToWindow(second, secondWindow);
    CHECK([[firstWindow.contentView accessibilityChildren] count] == 1);
    CHECK([[firstWindow.contentView accessibilityChildren] firstObject] == first);
    CHECK([[secondWindow.contentView accessibilityChildren] count] == 1);
    CHECK([[secondWindow.contentView accessibilityChildren] firstObject] == second);
    CHECK(NSEqualRects([first accessibilityFrame], NSMakeRect(15, 283, 100, 30)));
    CHECK(NSEqualRects([child accessibilityFrame], NSMakeRect(15, 283, 100, 30)));
    CHECK(NSEqualRects([second accessibilityFrame], NSMakeRect(505, 863, 100, 30)));
    AccessibilityElementSetFrame(first, 15, 20, 120, 40);
    CHECK(NSEqualRects([first accessibilityFrame], NSMakeRect(25, 260, 120, 40)));
    CHECK(NSEqualRects([second accessibilityFrame], NSMakeRect(505, 863, 100, 30)));
    AccessibilityElementDestroy(child);
    [first retain];
    AccessibilityElementDestroy(first);
    CHECK([first accessibilityParent] == nil);
    [first release];
    CHECK([[firstWindow.contentView accessibilityChildren] count] == 0);
    CHECK([[secondWindow.contentView accessibilityChildren] count] == 1);
    AccessibilityAttachToWindow(second, firstWindow);
    CHECK([[firstWindow.contentView accessibilityChildren] firstObject] == second);
    CHECK([[secondWindow.contentView accessibilityChildren] count] == 0);
    AccessibilityElementDestroy(second);
    [firstWindow release];
    [secondWindow release];
    return 0;
}

static int testSuperMethod(void) {
    CHECK(testViewOverride([ProbeBaseView class]) == 0);
    CHECK(testViewOverride([ProbeOwnView class]) == 0);
    CHECK(testViewOverride([ProbeInheritedView class]) == 0);
    return 0;
}

static int testWindowLifetime(void) {
    ProbeWindow* window = newWindow([ProbeOwnView class], NSMakePoint(10, 20));
    id element = newElement(AccessibilityRoleButton, "Button", nil);
    AccessibilityAttachToWindow(element, window);
    CHECK([element accessibilityParent] == window.contentView);
    [window release];
    CHECK([element accessibilityParent] == nil);
    CHECK(NSEqualRects([element accessibilityFrame], NSMakeRect(5, 7, 100, 30)));
    AccessibilityElementDestroy(element);
    CHECK(deallocations == 1);
    return 0;
}

int main(int argc, const char* argv[]) {
    @autoreleasepool {
        CHECK(argc == 2);
        if (strcmp(argv[1], "attributes") == 0) return testAttributes();
        if (strcmp(argv[1], "ownership") == 0) return testOwnership();
        if (strcmp(argv[1], "window-isolation") == 0) return testWindows(NO);
        if (strcmp(argv[1], "mixed-classes") == 0) return testWindows(YES);
        if (strcmp(argv[1], "super-method") == 0) return testSuperMethod();
        if (strcmp(argv[1], "window-lifetime") == 0) return testWindowLifetime();
        fprintf(stderr, "Unknown native probe: %s\n", argv[1]);
        return 1;
    }
}
