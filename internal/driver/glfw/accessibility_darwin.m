//go:build accessibility && darwin

#import "accessibility_darwin.h"
#import <Cocoa/Cocoa.h>
#import <objc/runtime.h>

static NSMutableArray<NSAccessibilityElement*>* globalAccessibilityElements = nil;
static NSView* targetContentView = nil;
static NSWindow* targetWindow = nil;
static IMP originalAccessibilityChildrenIMP = NULL;
static BOOL contentViewSwizzled = NO;
static BOOL appAccessibilitySwizzled = NO;

@interface AccessibleElement : NSAccessibilityElement

@property (nonatomic, assign) AccessibilityRole role;
@property (nonatomic, strong) NSString* title;
@property (nonatomic, strong) NSString* label;
@property (nonatomic, strong) NSString* value;
@property (nonatomic, assign) NSRect frame;
@property (nonatomic, assign) NSRect localFrame;
@property (nonatomic, assign) BOOL enabled;
@property (nonatomic, assign) BOOL focused;
@property (nonatomic, assign) BOOL selected;
@property (nonatomic, assign) BOOL expanded;
@property (nonatomic, assign) BOOL checked;
@property (nonatomic, assign) int supportedActions;
@property (nonatomic, strong) NSMutableArray<AccessibleElement*>* children;
@property (nonatomic, assign) id parentElement;
@property (nonatomic, assign) AccessibilityActionCallback actionCallback;
@property (nonatomic, assign) void* callbackContext;
@property (nonatomic, assign) AccessibilityContextDestroy contextDestroy;

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
    if (_contextDestroy && _callbackContext) {
        _contextDestroy(_callbackContext);
        _callbackContext = NULL;
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
        case AccessibilityRoleCheckbox:
            return NSAccessibilityCheckBoxRole;
        case AccessibilityRoleHeading:
        case AccessibilityRoleStaticText:
            return NSAccessibilityStaticTextRole;
        case AccessibilityRoleImage:
            return NSAccessibilityImageRole;
        case AccessibilityRoleLink:
            return NSAccessibilityLinkRole;
        case AccessibilityRoleList:
            return NSAccessibilityListRole;
        case AccessibilityRoleListItem:
        case AccessibilityRoleTreeItem:
            return NSAccessibilityRowRole;
        case AccessibilityRoleProgressBar:
            return NSAccessibilityProgressIndicatorRole;
        case AccessibilityRoleRadio:
            return NSAccessibilityRadioButtonRole;
        case AccessibilityRoleSeparator:
            return @"AXSeparator";
        case AccessibilityRoleSlider:
            return NSAccessibilitySliderRole;
        case AccessibilityRoleTab:
            return NSAccessibilityRadioButtonRole;
        case AccessibilityRoleTabList:
            return NSAccessibilityTabGroupRole;
        case AccessibilityRoleTable:
            return NSAccessibilityTableRole;
        case AccessibilityRoleTextField:
            return NSAccessibilityTextFieldRole;
        case AccessibilityRoleTree:
            return NSAccessibilityOutlineRole;
        case AccessibilityRoleContainer:
        case AccessibilityRoleGroup:
        default:
            return NSAccessibilityGroupRole;
    }
}

- (NSAccessibilitySubrole)accessibilitySubrole {
    if (self.role == AccessibilityRoleTab) {
        return NSAccessibilityTabButtonSubrole;
    }
    return [super accessibilitySubrole];
}

- (NSString*)accessibilityRoleDescription {
    if (self.role == AccessibilityRoleHeading) {
        return @"heading";
    }
    return [super accessibilityRoleDescription];
}

- (NSString*)accessibilityLabel {
    return self.label;
}

- (NSString*)accessibilityTitle {
    return self.title;
}

- (id)accessibilityValue {
    if (self.role == AccessibilityRoleCheckbox || self.role == AccessibilityRoleRadio) {
        return @(self.checked ? 1 : 0);
    }
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

- (id)accessibilityWindow {
    return targetWindow;
}

- (id)accessibilityTopLevelUIElement {
    return targetWindow;
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

- (BOOL)isAccessibilitySelected {
    return self.selected;
}

- (BOOL)isAccessibilityExpanded {
    return self.expanded;
}

- (BOOL)invokeAction:(AccessibilityActionCode)code mask:(int)mask withValue:(const char*)value {
    if (!self.actionCallback) {
        return NO;
    }
    if ((self.supportedActions & mask) == 0) {
        return NO;
    }
    return self.actionCallback(self.callbackContext, (int)code, value) != 0;
}

- (BOOL)accessibilityPerformPress {
    return [self invokeAction:AccessibilityActionPress mask:AccessibilityActionMaskPress withValue:NULL];
}

- (BOOL)accessibilityPerformIncrement {
    return [self invokeAction:AccessibilityActionIncrement mask:AccessibilityActionMaskIncrement withValue:NULL];
}

- (BOOL)accessibilityPerformDecrement {
    return [self invokeAction:AccessibilityActionDecrement mask:AccessibilityActionMaskDecrement withValue:NULL];
}

- (BOOL)accessibilityPerformShowMenu {
    return [self invokeAction:AccessibilityActionShowMenu mask:AccessibilityActionMaskShowMenu withValue:NULL];
}

- (void)setAccessibilityValue:(id)value {
    NSString* strValue = nil;
    if ([value isKindOfClass:[NSString class]]) {
        strValue = (NSString*)value;
    } else if ([value respondsToSelector:@selector(stringValue)]) {
        strValue = [value stringValue];
    }
    if (!strValue) {
        return;
    }
    if ([self invokeAction:AccessibilityActionSetValue mask:AccessibilityActionMaskSetValue withValue:[strValue UTF8String]]) {
        self.value = [[strValue copy] autorelease];
        NSAccessibilityPostNotification(self, NSAccessibilityValueChangedNotification);
    }
}

- (BOOL)isAccessibilityElement {
    // Containers without children should not appear as a leaf.
    if ((self.role == AccessibilityRoleContainer || self.role == AccessibilityRoleGroup) &&
        [self.children count] > 0) {
        return NO;
    }
    return YES;
}

// Legacy API: ensure elements are not ignored.
#pragma clang diagnostic push
#pragma clang diagnostic ignored "-Wdeprecated-declarations"
- (BOOL)accessibilityIsIgnored {
    return ![self isAccessibilityElement];
}

- (NSArray*)accessibilityAttributeNames {
    NSMutableArray* attrs = [NSMutableArray arrayWithArray:@[
        NSAccessibilityRoleAttribute,
        NSAccessibilityRoleDescriptionAttribute,
        NSAccessibilityTitleAttribute,
        NSAccessibilityDescriptionAttribute,
        NSAccessibilityParentAttribute,
        NSAccessibilityChildrenAttribute,
        NSAccessibilityWindowAttribute,
        NSAccessibilityTopLevelUIElementAttribute,
        NSAccessibilityPositionAttribute,
        NSAccessibilitySizeAttribute,
        NSAccessibilityFocusedAttribute,
        NSAccessibilityEnabledAttribute,
    ]];
    if (self.value) {
        [attrs addObject:NSAccessibilityValueAttribute];
    }
    return attrs;
}

- (id)accessibilityAttributeValue:(NSString*)attr {
    if ([attr isEqualToString:NSAccessibilityRoleAttribute]) {
        return [self accessibilityRole];
    }
    if ([attr isEqualToString:NSAccessibilityRoleDescriptionAttribute]) {
        return [self accessibilityRoleDescription];
    }
    if ([attr isEqualToString:NSAccessibilityTitleAttribute]) {
        return self.title;
    }
    if ([attr isEqualToString:NSAccessibilityDescriptionAttribute]) {
        return self.label;
    }
    if ([attr isEqualToString:NSAccessibilityValueAttribute]) {
        return [self accessibilityValue];
    }
    if ([attr isEqualToString:NSAccessibilityChildrenAttribute]) {
        return NSAccessibilityUnignoredChildren([self accessibilityChildren]);
    }
    if ([attr isEqualToString:NSAccessibilityParentAttribute]) {
        return NSAccessibilityUnignoredAncestor([self accessibilityParent]);
    }
    if ([attr isEqualToString:NSAccessibilityWindowAttribute]) {
        return targetWindow;
    }
    if ([attr isEqualToString:NSAccessibilityTopLevelUIElementAttribute]) {
        return targetWindow;
    }
    if ([attr isEqualToString:NSAccessibilityPositionAttribute]) {
        NSRect frame = [self accessibilityFrame];
        return [NSValue valueWithPoint:frame.origin];
    }
    if ([attr isEqualToString:NSAccessibilitySizeAttribute]) {
        NSRect frame = [self accessibilityFrame];
        return [NSValue valueWithSize:frame.size];
    }
    if ([attr isEqualToString:NSAccessibilityFocusedAttribute]) {
        return @(self.focused);
    }
    if ([attr isEqualToString:NSAccessibilityEnabledAttribute]) {
        return @(self.enabled);
    }
    return [super accessibilityAttributeValue:attr];
}
#pragma clang diagnostic pop

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

// Swizzled methods for GLFWContentView to fix accessibility hierarchy.
// GLFWContentView reports AXUnknown role and isAccessibilityElement=NO by
// default, which breaks the AX chain from window → content → elements.
static NSAccessibilityRole customAccessibilityRole(id self, SEL _cmd) {
    return NSAccessibilityGroupRole;
}

static BOOL customIsAccessibilityElement(id self, SEL _cmd) {
    return YES;
}

static BOOL customAccessibilityIsNotIgnored(id self, SEL _cmd) {
    return NO;
}

// GLFWWindow needs to return its content view as a child.
// Without this, the window appears empty in the AX hierarchy.
static NSArray* customWindowAccessibilityChildren(id self, SEL _cmd) {
    NSWindow* window = (NSWindow*)self;
    NSView* cv = [window contentView];
    if (cv) {
        return @[cv];
    }
    return @[];
}

static void swizzleContentViewAccessibility(NSView* contentView) {
    if (contentViewSwizzled) {
        return;
    }
    contentViewSwizzled = YES;

    Class viewClass = [contentView class];

    class_addMethod(viewClass, @selector(accessibilityRole),
                    (IMP)customAccessibilityRole, "@@:");
    class_addMethod(viewClass, @selector(isAccessibilityElement),
                    (IMP)customIsAccessibilityElement, "B@:");

#pragma clang diagnostic push
#pragma clang diagnostic ignored "-Wdeprecated-declarations"
    class_addMethod(viewClass, @selector(accessibilityIsIgnored),
                    (IMP)customAccessibilityIsNotIgnored, "B@:");
#pragma clang diagnostic pop
}

static void swizzleAppAccessibility(void) {
    if (appAccessibilitySwizzled) {
        return;
    }
    appAccessibilitySwizzled = YES;

    if (targetWindow) {
        Class winClass = [targetWindow class];
        class_addMethod(winClass, @selector(accessibilityChildren),
                        (IMP)customWindowAccessibilityChildren, "@@:");
    }
}

AccessibilityElementRef AccessibilityElementCreate(
    AccessibilityRole role,
    const char* title,
    const char* label,
    double x, double y, double width, double height,
    AccessibilityElementRef parent,
    AccessibilityActionCallback callback,
    void* callbackContext,
    AccessibilityContextDestroy contextDestroy
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
        elem.contextDestroy = contextDestroy;

        return (void*)elem;
    }
}

void AccessibilityElementSetFrame(AccessibilityElementRef elem, double x, double y, double width, double height) {
    @autoreleasepool {
        AccessibleElement* element = (__bridge AccessibleElement*)elem;
        element.frame = NSMakeRect(x, y, width, height);
        NSAccessibilityPostNotification(element, NSAccessibilityMovedNotification);
        NSAccessibilityPostNotification(element, NSAccessibilityResizedNotification);
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

void AccessibilityElementSetStates(AccessibilityElementRef elem, int stateMask) {
    @autoreleasepool {
        AccessibleElement* element = (__bridge AccessibleElement*)elem;
        BOOL wasSelected = element.selected;
        BOOL wasExpanded = element.expanded;

        element.checked  = (stateMask & AccessibilityStateMaskChecked)  != 0;
        element.enabled  = (stateMask & AccessibilityStateMaskDisabled) == 0;
        element.expanded = (stateMask & AccessibilityStateMaskExpanded) != 0;
        element.focused  = (stateMask & AccessibilityStateMaskFocused)  != 0;
        element.selected = (stateMask & AccessibilityStateMaskSelected) != 0;

        if (wasSelected != element.selected) {
            NSAccessibilityPostNotification(element, NSAccessibilitySelectedChildrenChangedNotification);
        }
        if (wasExpanded != element.expanded) {
            NSAccessibilityPostNotification(element, NSAccessibilityRowExpandedNotification);
        }
    }
}

void AccessibilityElementSetSupportedActions(AccessibilityElementRef elem, int actionMask) {
    @autoreleasepool {
        AccessibleElement* element = (__bridge AccessibleElement*)elem;
        element.supportedActions = actionMask;
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
            [childElem release]; // Transfer create ownership to the parent.
            NSAccessibilityPostNotification(parentElem, NSAccessibilityLayoutChangedNotification);
        }
    }
}

void AccessibilityElementRemoveChild(AccessibilityElementRef parent, AccessibilityElementRef child) {
    @autoreleasepool {
        AccessibleElement* parentElem = (__bridge AccessibleElement*)parent;
        AccessibleElement* childElem = (__bridge AccessibleElement*)child;
        childElem.parentElement = nil;
        [parentElem.children removeObject:childElem];
        if (parentElem) {
            NSAccessibilityPostNotification(parentElem, NSAccessibilityLayoutChangedNotification);
        }
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

            if (targetContentView) {
                swizzleContentViewAccessibility(targetContentView);
            }
            swizzleAppAccessibility();
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
                // Use class_addMethod to add ONLY to GLFWContentView.
                // The original code used method_setImplementation which modified
                // the parent class (NSView/NSResponder), breaking ALL views
                // including NSWindow's accessibility hierarchy.
                class_addMethod(viewClass, selector,
                                (IMP)customAccessibilityChildren, "@@:");
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
        if (globalAccessibilityElements) {
            [globalAccessibilityElements removeObject:element];
        }
        if ([element.parentElement isKindOfClass:[AccessibleElement class]]) {
            AccessibleElement* parent = (AccessibleElement*)element.parentElement;
            element.parentElement = nil;
            [parent.children removeObject:element];
            return;
        }
        element.parentElement = nil;
        [element release];
    }
}
