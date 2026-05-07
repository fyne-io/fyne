# macOS Accessibility Fix Patterns for GLFW/Fyne Apps

## Complete Swizzle Implementation Pattern

This is the proven pattern for fixing GLFW-based app accessibility on macOS.
The fixes target three GLFW classes: GLFWContentView, GLFWWindow, and optionally NSApplication.

### 1. GLFWContentView Fixes

GLFWContentView (an NSView subclass) has three problems by default:
- `accessibilityRole` returns `AXUnknown`
- `isAccessibilityElement` returns `NO`
- `accessibilityIsIgnored` returns `YES` (NSView default)

All three must be fixed for the content view to appear in the AX hierarchy.

```objc
// Static functions to be used as method implementations
static NSAccessibilityRole customAccessibilityRole(id self, SEL _cmd) {
    return NSAccessibilityGroupRole;
}

static BOOL customIsAccessibilityElement(id self, SEL _cmd) {
    return YES;
}

static BOOL customAccessibilityIsNotIgnored(id self, SEL _cmd) {
    return NO;
}

// Apply to GLFWContentView class
static void swizzleContentViewAccessibility(NSView* contentView) {
    Class viewClass = [contentView class];

    // class_addMethod adds ONLY to the specific class, not parent classes
    class_addMethod(viewClass, @selector(accessibilityRole),
                    (IMP)customAccessibilityRole, "@@:");
    class_addMethod(viewClass, @selector(isAccessibilityElement),
                    (IMP)customIsAccessibilityElement, "B@:");

    // Legacy API (still used by external AX clients)
    class_addMethod(viewClass, @selector(accessibilityIsIgnored),
                    (IMP)customAccessibilityIsNotIgnored, "B@:");
}
```

### 2. GLFWWindow Fix

GLFWWindow's `accessibilityChildren` returns empty because the content view is ignored.
Even after fixing the content view, the window may still return empty children.

```objc
static NSArray* customWindowAccessibilityChildren(id self, SEL _cmd) {
    NSWindow* window = (NSWindow*)self;
    NSView* cv = [window contentView];
    if (cv) {
        return @[cv];
    }
    return @[];
}

// Apply to GLFWWindow class
Class winClass = [targetWindow class];
class_addMethod(winClass, @selector(accessibilityChildren),
                (IMP)customWindowAccessibilityChildren, "@@:");
```

### 3. Custom AccessibleElement Requirements

Each custom `NSAccessibilityElement` subclass must implement:

```objc
// Modern protocol methods
- (id)accessibilityWindow { return targetWindow; }
- (id)accessibilityTopLevelUIElement { return targetWindow; }
- (id)accessibilityParent { return self.parentElement ?: targetContentView; }

// Legacy protocol methods (for external AX clients)
- (BOOL)accessibilityIsIgnored {
    return ![self isAccessibilityElement];
}

- (NSArray*)accessibilityAttributeNames {
    return @[
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
    ];
}

- (id)accessibilityAttributeValue:(NSString*)attr {
    if ([attr isEqualToString:NSAccessibilityRoleAttribute])
        return [self accessibilityRole];
    if ([attr isEqualToString:NSAccessibilityChildrenAttribute])
        return NSAccessibilityUnignoredChildren([self accessibilityChildren]);
    if ([attr isEqualToString:NSAccessibilityParentAttribute])
        return NSAccessibilityUnignoredAncestor([self accessibilityParent]);
    if ([attr isEqualToString:NSAccessibilityWindowAttribute])
        return targetWindow;
    if ([attr isEqualToString:NSAccessibilityPositionAttribute]) {
        NSRect frame = [self accessibilityFrame];
        return [NSValue valueWithPoint:frame.origin];
    }
    if ([attr isEqualToString:NSAccessibilitySizeAttribute]) {
        NSRect frame = [self accessibilityFrame];
        return [NSValue valueWithSize:frame.size];
    }
    // ... handle other attributes
    return [super accessibilityAttributeValue:attr];
}
```

### 4. Coordinate System

macOS accessibility uses **screen coordinates** with **bottom-left origin** (like Cocoa),
NOT top-left origin (like most UI frameworks including Fyne).

```objc
- (NSRect)accessibilityFrame {
    // Convert from local (top-left origin) to screen (bottom-left origin)
    NSRect contentBounds = [targetContentView bounds];
    NSPoint contentBottomLeft = [targetContentView convertPoint:NSMakePoint(0,0) toView:nil];
    contentBottomLeft = [targetWindow convertPointToScreen:contentBottomLeft];

    double screenX = contentBottomLeft.x + self.localFrame.origin.x;
    double screenY = contentBottomLeft.y +
        (contentBounds.size.height - self.localFrame.origin.y - self.localFrame.size.height);

    return NSMakeRect(screenX, screenY,
                      self.localFrame.size.width, self.localFrame.size.height);
}
```

## Method Type Encodings Quick Reference

When using `class_addMethod`, the type encoding string must match the method signature:

| Return Type | Encoding |
|------------|----------|
| `id` (object) | `@` |
| `BOOL` | `B` |
| `void` | `v` |
| `NSRect` | `{CGRect={CGPoint=dd}{CGSize=dd}}` |

Common full encodings:
- `"@@:"` → `id method(id self, SEL _cmd)` — returns object, no args
- `"B@:"` → `BOOL method(id self, SEL _cmd)` — returns BOOL, no args
- `"@@:@"` → `id method(id self, SEL _cmd, id arg)` — returns object, one object arg
- `"v@:"` → `void method(id self, SEL _cmd)` — returns void, no args

## NSAccessibility Role Mappings

| Widget Type | NSAccessibility Role |
|------------|---------------------|
| Button | `NSAccessibilityButtonRole` |
| Checkbox | `NSAccessibilityCheckBoxRole` |
| Text (static) | `NSAccessibilityStaticTextRole` |
| Text field | `NSAccessibilityTextFieldRole` |
| Image | `NSAccessibilityImageRole` |
| Link | `NSAccessibilityLinkRole` |
| List | `NSAccessibilityListRole` |
| Progress bar | `NSAccessibilityProgressIndicatorRole` |
| Radio button | `NSAccessibilityRadioButtonRole` |
| Slider | `NSAccessibilitySliderRole` |
| Tab group | `NSAccessibilityTabGroupRole` |
| Table | `NSAccessibilityTableRole` |
| Outline/Tree | `NSAccessibilityOutlineRole` |
| Group/Container | `NSAccessibilityGroupRole` |
| Window | `NSAccessibilityWindowRole` |
| Application | `NSAccessibilityApplicationRole` |
