//go:build accessibility && darwin

#import "accessibility_darwin.h"
#import <Cocoa/Cocoa.h>
#import <objc/runtime.h>

static NSMutableArray<NSAccessibilityElement*>* globalAccessibilityElements = nil;
static NSView* targetContentView = nil;
static NSWindow* targetWindow = nil;
static IMP originalAccessibilityChildrenIMP = NULL;

@interface AccessibleElement : NSAccessibilityElement

@property (nonatomic, assign) AccessibilityRole role;
@property (nonatomic, strong) NSString* title;
@property (nonatomic, strong) NSString* label;
@property (nonatomic, strong) NSString* value;
@property (nonatomic, assign) NSRect frame;
@property (nonatomic, assign) NSRect localFrame;
@property (nonatomic, assign) BOOL enabled;
@property (nonatomic, assign) BOOL focused;
@property (nonatomic, strong) NSMutableArray<AccessibleElement*>* children;
@property (nonatomic, assign) id parentElement;
@property (nonatomic, assign) AccessibilityActionCallback actionCallback;
@property (nonatomic, assign) void* callbackContext;

@end

@implementation AccessibleElement

- (instancetype)initWithParent:(id)parent {
    self = [super init];
    if (self) {
        _children = [[NSMutableArray alloc] init];
        _enabled = YES;
        _focused = NO;
        _parentElement = parent;
    }
    return self;
}

- (void)dealloc {
    for (AccessibleElement* child in _children) {
        child.parentElement = nil;
    }
    [_children release];
    [_title release];
    [_label release];
    [_value release];
    [super dealloc];
}

- (NSAccessibilityRole)accessibilityRole {
    switch (self.role) {
        case AccessibilityRoleWindow:
            return NSAccessibilityWindowRole;
        case AccessibilityRoleButton:
            return NSAccessibilityButtonRole;
        case AccessibilityRoleStaticText:
            return NSAccessibilityStaticTextRole;
        case AccessibilityRoleTextField:
            return NSAccessibilityTextFieldRole;
        case AccessibilityRoleCheckbox:
            return NSAccessibilityCheckBoxRole;
        case AccessibilityRoleGroup:
            return NSAccessibilityGroupRole;
        default:
            return NSAccessibilityUnknownRole;
    }
}

- (NSString*)accessibilityLabel {
    return self.label;
}

- (NSString*)accessibilityTitle {
    return self.title;
}

- (id)accessibilityValue {
    return self.value;
}

- (NSRect)accessibilityFrame {
    if (targetWindow && targetContentView) {
        NSRect contentBounds = [targetContentView bounds];
        NSPoint contentBottomLeft = NSMakePoint(0, 0);
        contentBottomLeft = [targetContentView convertPoint:contentBottomLeft toView:nil];
        contentBottomLeft = [targetWindow convertPointToScreen:contentBottomLeft];

        double localX = self.localFrame.origin.x;
        double localY = self.localFrame.origin.y;
        double localWidth = self.localFrame.size.width;
        double localHeight = self.localFrame.size.height;

        // Convert from Fyne coordinates (top-left origin) to screen coordinates (bottom-left origin)
        double screenX = contentBottomLeft.x + localX;
        double screenY = contentBottomLeft.y + (contentBounds.size.height - localY - localHeight);

        return NSMakeRect(screenX, screenY, localWidth, localHeight);
    }

    return self.frame;
}

- (id)accessibilityParent {
    if (self.parentElement) {
        return self.parentElement;
    }
    if (targetContentView) {
        return targetContentView;
    }
    NSWindow* window = [[NSApplication sharedApplication] mainWindow];
    return [window contentView];
}

- (NSArray*)accessibilityChildren {
    return [[self.children copy] autorelease];
}

- (BOOL)isAccessibilityEnabled {
    return self.enabled;
}

- (BOOL)isAccessibilityFocused {
    return self.focused;
}

- (void)setAccessibilityFocused:(BOOL)focused {
    self.focused = focused;
}

- (BOOL)accessibilityPerformPress {
    if (self.actionCallback) {
        self.actionCallback(self.callbackContext);
        return YES;
    }
    return NO;
}

- (BOOL)isAccessibilityElement {
    // Groups/containers with children should return NO so VoiceOver navigates into them
    // Leaf elements (buttons, labels, etc.) should return YES
    if (self.role == AccessibilityRoleGroup && [self.children count] > 0) {
        return NO;
    }
    return YES;
}

@end

static NSArray* customAccessibilityChildren(id self, SEL _cmd) {
    if (globalAccessibilityElements && [globalAccessibilityElements count] > 0) {
        return [[globalAccessibilityElements copy] autorelease];
    }
    if (originalAccessibilityChildrenIMP) {
        return ((NSArray*(*)(id, SEL))originalAccessibilityChildrenIMP)(self, _cmd);
    }
    return @[];
}

AccessibilityElementRef AccessibilityElementCreate(
    AccessibilityRole role,
    const char* title,
    const char* label,
    double x, double y, double width, double height,
    AccessibilityElementRef parent,
    AccessibilityActionCallback callback,
    void* callbackContext
) {
    @autoreleasepool {
        if (!globalAccessibilityElements) {
            globalAccessibilityElements = [[NSMutableArray alloc] init];
        }

        AccessibleElement* elem = [[AccessibleElement alloc] initWithParent:parent];
        elem.role = role;
        elem.title = title ? [NSString stringWithUTF8String:title] : @"";
        elem.label = label ? [NSString stringWithUTF8String:label] : @"";

        // Store local frame (relative to window content view)
        elem.localFrame = CGRectMake(x, y, width, height);
        elem.frame = elem.localFrame; // Will be recalculated dynamically in accessibilityFrame

        elem.actionCallback = callback;
        elem.callbackContext = callbackContext;

        return (void*)elem;
    }
}

void AccessibilityElementSetFrame(AccessibilityElementRef elem, double x, double y, double width, double height) {
    @autoreleasepool {
        AccessibleElement* element = (__bridge AccessibleElement*)elem;
        NSRect frame = NSMakeRect(x, y, width, height);
        BOOL moved = !NSEqualPoints(element.localFrame.origin, frame.origin);
        BOOL resized = !NSEqualSizes(element.localFrame.size, frame.size);
        element.localFrame = frame;
        element.frame = frame;
        if (moved) {
            NSAccessibilityPostNotification(element, NSAccessibilityMovedNotification);
        }
        if (resized) {
            NSAccessibilityPostNotification(element, NSAccessibilityResizedNotification);
        }
    }
}

void AccessibilityElementSetTitle(AccessibilityElementRef elem, const char* title) {
    @autoreleasepool {
        AccessibleElement* element = (__bridge AccessibleElement*)elem;
        element.title = title ? [NSString stringWithUTF8String:title] : @"";
    }
}

void AccessibilityElementSetLabel(AccessibilityElementRef elem, const char* label) {
    @autoreleasepool {
        AccessibleElement* element = (__bridge AccessibleElement*)elem;
        element.label = label ? [NSString stringWithUTF8String:label] : @"";
        NSAccessibilityPostNotification(element, NSAccessibilityTitleChangedNotification);
    }
}

void AccessibilityElementSetValue(AccessibilityElementRef elem, const char* value) {
    @autoreleasepool {
        AccessibleElement* element = (__bridge AccessibleElement*)elem;
        element.value = value ? [NSString stringWithUTF8String:value] : @"";
        NSAccessibilityPostNotification(element, NSAccessibilityValueChangedNotification);
    }
}

void AccessibilityElementSetEnabled(AccessibilityElementRef elem, int enabled) {
    @autoreleasepool {
        AccessibleElement* element = (__bridge AccessibleElement*)elem;
        element.enabled = enabled != 0;
    }
}

void AccessibilityElementSetFocused(AccessibilityElementRef elem, int focused) {
    @autoreleasepool {
        AccessibleElement* element = (__bridge AccessibleElement*)elem;
        element.focused = focused != 0;
        if (focused) {
            NSAccessibilityPostNotification(element, NSAccessibilityFocusedUIElementChangedNotification);
        }
    }
}

void AccessibilityElementAddChild(AccessibilityElementRef parent, AccessibilityElementRef child) {
    @autoreleasepool {
        if (!parent || !child) {
            return;
        }

        AccessibleElement* parentElem = (__bridge AccessibleElement*)parent;
        AccessibleElement* childElem = (__bridge AccessibleElement*)child;

        if (!parentElem || !childElem || !parentElem.children) {
            return;
        }

        if (![parentElem.children containsObject:childElem]) {
            [parentElem.children addObject:childElem];
            childElem.parentElement = parentElem;
            [childElem release];
        }
    }
}

void AccessibilityElementRemoveChild(AccessibilityElementRef parent, AccessibilityElementRef child) {
    @autoreleasepool {
        AccessibleElement* parentElem = (__bridge AccessibleElement*)parent;
        AccessibleElement* childElem = (__bridge AccessibleElement*)child;
        childElem.parentElement = nil;
        [parentElem.children removeObject:childElem];
    }
}

void AccessibilityElementSetParent(AccessibilityElementRef child, AccessibilityElementRef parent) {
    @autoreleasepool {
        AccessibleElement* childElem = (__bridge AccessibleElement*)child;
        AccessibleElement* parentElem = (__bridge AccessibleElement*)parent;
        childElem.parentElement = parentElem;
    }
}

void AccessibilitySetTargetWindow(void* nsWindow) {
    @autoreleasepool {
        targetWindow = (NSWindow*)nsWindow;
        if (targetWindow) {
            targetContentView = [targetWindow contentView];
        }
    }
}

void AccessibilityAttachToWindow(AccessibilityElementRef elem) {
    @autoreleasepool {
        AccessibleElement* element = (AccessibleElement*)elem;

        if (!globalAccessibilityElements) {
            globalAccessibilityElements = [[NSMutableArray alloc] init];
        }

        NSWindow* window = targetWindow;
        if (!window) {
            return;
        }

        NSView* contentView = targetContentView;
        if (!contentView) {
            return;
        }

        targetContentView = contentView;

        if (!originalAccessibilityChildrenIMP) {
            Class viewClass = [contentView class];
            SEL selector = @selector(accessibilityChildren);
            Method originalMethod = class_getInstanceMethod(viewClass, selector);
            if (originalMethod) {
                originalAccessibilityChildrenIMP = method_getImplementation(originalMethod);
                class_replaceMethod(viewClass, selector, (IMP)customAccessibilityChildren,
                                    method_getTypeEncoding(originalMethod));
            }
        }

        element.parentElement = contentView;

        if (![globalAccessibilityElements containsObject:element]) {
            [globalAccessibilityElements addObject:element];
        }

        NSAccessibilityPostNotification(contentView, NSAccessibilityCreatedNotification);
        NSAccessibilityPostNotification(element, NSAccessibilityCreatedNotification);
    }
}

void AccessibilityPostNotification(AccessibilityElementRef elem, const char* notification) {
    @autoreleasepool {
        AccessibleElement* element = (__bridge AccessibleElement*)elem;
        NSString* notificationName = notification ? [NSString stringWithUTF8String:notification] : nil;
        if (notificationName) {
            NSAccessibilityPostNotification(element, notificationName);
        }
    }
}

void AccessibilityElementDestroy(AccessibilityElementRef elem) {
    @autoreleasepool {
        if (!elem) return;

        AccessibleElement* element = (AccessibleElement*)elem;
        [globalAccessibilityElements removeObject:element];
        if ([element.parentElement isKindOfClass:[AccessibleElement class]]) {
            AccessibleElement* parent = element.parentElement;
            element.parentElement = nil;
            [parent.children removeObject:element];
            return;
        }
        element.parentElement = nil;
        [element release];
    }
}
