---
name: macos-accessibility
description: Diagnose and fix macOS accessibility issues for GUI applications, especially GLFW/Fyne-based apps. Use this skill when an app's windows are not visible to System Events, AppleScript, VoiceOver, or computer-use MCP servers, or when accessibility hierarchy is broken. Triggers include "accessibility", "System Events can't see windows", "AXWindow", "VoiceOver", "computer-use can't interact", or "accessibilityIsIgnored".
---

# macOS Accessibility Diagnostics & Fixes

## Overview

This skill helps diagnose and fix macOS accessibility issues, particularly for apps built with custom rendering frameworks (GLFW, SDL, Electron, etc.) that don't use standard AppKit controls. These apps often have broken accessibility hierarchies that prevent System Events, VoiceOver, AppleScript, and AI computer-use tools from interacting with them.

## Quick Diagnostic Checklist

Before diving deep, check these common issues in order:

### 1. Permissions (check FIRST — most common false alarm)

The **querying process** must have Accessibility permission, not just the target app.

```bash
# Check if osascript/Terminal has accessibility permission
osascript -e 'tell application "System Events"
    tell process "Google Chrome"
        return count of windows
    end tell
end tell'
```

If this returns 0 for a known-good app like Chrome or Finder, the problem is **permissions**, not code. Fix:
- **System Settings → Privacy & Security → Accessibility**
- Add Terminal.app, iTerm, VS Code, or whatever app runs the osascript command
- Note: `osascript` itself is not an app — add the PARENT app that invokes it

### 2. App Bundling (required for Accessibility Inspector)

Xcode's Accessibility Inspector only lists `.app` bundles in its target dropdown. A bare Go binary (`go run .` or `./myapp`) will **not** appear. System Events (AppleScript) can interact with bare binaries, but Accessibility Inspector cannot.

To make the app appear, wrap the binary in a minimal `.app` bundle:

```bash
APP=/tmp/MyApp.app
mkdir -p "$APP/Contents/MacOS"
cp myapp "$APP/Contents/MacOS/MyApp"
cat > "$APP/Contents/Info.plist" << 'EOF'
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN"
  "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>CFBundleExecutable</key><string>MyApp</string>
    <key>CFBundleIdentifier</key><string>com.example.myapp</string>
    <key>CFBundleName</key><string>MyApp</string>
    <key>CFBundlePackageType</key><string>APPL</string>
    <key>NSHighResolutionCapable</key><true/>
</dict>
</plist>
EOF
open "$APP"
```

> **Note:** If the app requires a specific working directory (e.g., to find config files relative to the repo root), you'll need a small C or shell launcher wrapper inside the bundle that `chdir`s before exec-ing the real binary.

### 3. Window Visibility to CGWindow API

```bash
# Check if the window exists at the CGWindow (Core Graphics) level
# This is lower-level than accessibility and should always work
python3 -c "
import subprocess
result = subprocess.run(['osascript', '-e',
    'tell application \"System Events\" to return name of every process whose visible is true'],
    capture_output=True, text=True)
print(result.stdout.strip())
"
```

If the app appears in the visible process list but not in `list_applications` from computer-use MCP, check if:
- The window is on a different display/Space
- The window is behind other windows
- The window hasn't finished rendering yet (timing issue)

### 4. AX Hierarchy Integrity

Use the diagnostic commands from the **Diagnostic AppleScript Commands** section below.

## macOS Accessibility Architecture

### Two API Paths (Critical Knowledge)

macOS has **two** accessibility dispatch paths that coexist:

| Path | Methods | Used By |
|------|---------|---------|
| **Modern** (10.10+) | `accessibilityRole`, `accessibilityChildren`, `isAccessibilityElement`, etc. | Internal Obj-C calls, some framework paths |
| **Legacy** | `accessibilityAttributeValue:`, `accessibilityAttributeNames`, `accessibilityIsIgnored` | External AX clients (System Events, VoiceOver, AXUIElement API) |

**Key insight**: Overriding ONLY the modern protocol methods may NOT affect what external AX clients see. For full compatibility, implement BOTH paths.

### AX Hierarchy Structure

A well-formed macOS AX hierarchy:

```
AXApplication (NSApplication)
 ├── AXWindow (NSWindow) — role=AXWindow, subrole=AXStandardWindow
 │   └── AXGroup (content view) — role=AXGroup
 │       ├── AXButton "Save" — role=AXButton
 │       ├── AXTextField "Name" — role=AXTextField
 │       └── AXStaticText "Hello" — role=AXStaticText
 └── AXMenuBar (NSMenu)
```

System Events discovers "windows" by looking for `AXWindow`-role elements in the app's `AXChildren`.

### Critical NSAccessibility Methods

For each level of the hierarchy, these methods must return correct values:

**NSApplication level:**
- `accessibilityWindows` → array of visible NSWindow objects
- `accessibilityChildren` → windows + menu bar
- `accessibilityRole` → `NSAccessibilityApplicationRole`

**NSWindow level:**
- `accessibilityChildren` → content view (and any other child views)
- `accessibilityRole` → `NSAccessibilityWindowRole`
- `accessibilityTitle` → window title string
- `isAccessibilityElement` → YES
- `accessibilityIsIgnored` → NO

**Content View level:**
- `accessibilityChildren` → child UI elements
- `accessibilityRole` → `NSAccessibilityGroupRole` (NOT `AXUnknown`)
- `isAccessibilityElement` → YES
- `accessibilityIsIgnored` → NO (NSView defaults to YES!)

**Each AccessibleElement must implement:**
- `accessibilityParent` → parent element
- `accessibilityWindow` → containing window
- `accessibilityTopLevelUIElement` → containing window
- `accessibilityFrame` → screen coordinates (bottom-left origin)
- `accessibilityRole` → appropriate role
- `isAccessibilityElement` → YES for leaf elements, NO for containers with children
- `accessibilityIsIgnored` → inverse of isAccessibilityElement

## Common Issues with GLFW/Custom Framework Apps

### Issue 0: method_setImplementation on Inherited Methods Breaks ALL Views (CRITICAL)

**Symptom**: Window disappears from System Events entirely; `count of windows` returns 0 for the app.

**Cause**: Using `method_setImplementation` on a method found via `class_getInstanceMethod` modifies the method **on the class where it's defined** (typically NSView or NSResponder), not the target subclass. This breaks `accessibilityChildren` for ALL NSView subclasses, including NSWindow.

**Example of the bug**:
```objc
// DANGEROUS: This modifies NSView/NSResponder, not GLFWContentView!
Method m = class_getInstanceMethod([contentView class], @selector(accessibilityChildren));
method_setImplementation(m, (IMP)customAccessibilityChildren);
```

**Fix**: Always use `class_addMethod` which adds ONLY to the specific subclass:
```objc
// SAFE: This adds only to GLFWContentView
class_addMethod([contentView class], @selector(accessibilityChildren),
                (IMP)customAccessibilityChildren, "@@:");
```

This is the **#1 cause** of broken accessibility in GLFW apps with custom AX bridges. If `class_addMethod` returns NO (the class already has its own method), use `class_copyMethodList` to verify the method exists on the class directly before using `method_setImplementation`.

### Issue 1: Content View Returns AXUnknown Role

**Symptom**: Elements appear but no window is recognized.

**Cause**: GLFW's `GLFWContentView` (NSView subclass) doesn't override `accessibilityRole`, so it returns `AXUnknown`.

**Fix**: Override `accessibilityRole` on GLFWContentView to return `NSAccessibilityGroupRole`:
```objc
// Use class_addMethod to add ONLY to the specific class
class_addMethod(viewClass, @selector(accessibilityRole),
                (IMP)customAccessibilityRole, "@@:");
```

### Issue 2: Content View Is Ignored (accessibilityIsIgnored = YES)

**Symptom**: Content view and its children don't appear in AX hierarchy.

**Cause**: NSView defaults `accessibilityIsIgnored` to YES. The modern `isAccessibilityElement` may return YES but the legacy `accessibilityIsIgnored` takes precedence for external clients.

**Fix**: Override BOTH:
```objc
class_addMethod(viewClass, @selector(isAccessibilityElement),
                (IMP)customIsAccessibilityElement, "B@:");
class_addMethod(viewClass, @selector(accessibilityIsIgnored),
                (IMP)customAccessibilityIsNotIgnored, "B@:");
```

### Issue 3: Window Has No Accessible Children

**Symptom**: Window exists in AX tree but appears empty; System Events may filter it.

**Cause**: GLFWWindow doesn't override `accessibilityChildren`, and the default returns empty because the content view was ignored.

**Fix**: Override `accessibilityChildren` on GLFWWindow to return the content view:
```objc
class_addMethod(winClass, @selector(accessibilityChildren),
                (IMP)customWindowAccessibilityChildren, "@@:");
```

### Issue 4: Elements Missing accessibilityWindow/accessibilityTopLevelUIElement

**Symptom**: Elements appear but can't be correlated to a window; navigation broken.

**Fix**: Each AccessibleElement must implement:
```objc
- (id)accessibilityWindow { return targetWindow; }
- (id)accessibilityTopLevelUIElement { return targetWindow; }
```

### Issue 5: Legacy API Not Implemented on Custom Elements

**Symptom**: Elements work with modern protocol but not with external AX clients.

**Fix**: Implement `accessibilityAttributeValue:`, `accessibilityAttributeNames`, and `accessibilityIsIgnored` on custom AccessibleElement classes.

## Safe Method Swizzling for AX

### SAFE: class_addMethod (class-specific)

```objc
// Adds the method ONLY to GLFWContentView — does NOT modify NSView/NSResponder
class_addMethod([contentView class], @selector(accessibilityRole),
                (IMP)myCustomRole, "@@:");
```

`class_addMethod` succeeds only if the class doesn't already have its OWN implementation (inherited methods don't count). This is safe because it shadows the parent's method without modifying it.

### DANGEROUS: method_setImplementation on inherited methods

```objc
// DANGER: If GLFWContentView inherits accessibilityRole from NSResponder,
// this modifies NSResponder's method — breaking ALL NSResponder subclasses!
Method m = class_getInstanceMethod([contentView class], @selector(accessibilityRole));
method_setImplementation(m, (IMP)myCustomRole); // Modifies NSResponder!
```

**Rule**: Always use `class_addMethod` first. Only use `method_setImplementation` when:
1. `class_addMethod` fails (class has its own method), AND
2. You've verified with `class_copyMethodList` that the method is on the target class specifically

### Pattern for safe swizzling:

```objc
if (!class_addMethod(targetClass, selector, newIMP, typeEncoding)) {
    // class has its own method — safe to modify directly
    unsigned int count;
    Method* methods = class_copyMethodList(targetClass, &count);
    for (unsigned int i = 0; i < count; i++) {
        if (method_getName(methods[i]) == selector) {
            method_setImplementation(methods[i], newIMP);
            break;
        }
    }
    free(methods);
}
```

## Diagnostic AppleScript Commands

### Check if a process has windows
```applescript
tell application "System Events"
    tell process "MyApp"
        return "windows=" & (count of windows) & " UI=" & (count of UI elements)
    end tell
end tell
```

### Get window details
```applescript
tell application "System Events"
    tell process "MyApp"
        tell window 1
            return {title:title, position:position, size:size}
        end tell
    end tell
end tell
```

### Inspect AX hierarchy
```applescript
tell application "System Events"
    tell process "MyApp"
        set elems to every UI element
        repeat with e in elems
            log (role of e) & " " & (name of e)
        end repeat
    end tell
end tell
```

### Full accessibility tree dump
```applescript
tell application "System Events"
    tell process "MyApp"
        return entire contents
    end tell
end tell
```

## Diagnostic Python Script (AXUIElement API)

For lower-level AX tree inspection that bypasses System Events:

```python
import ctypes, ctypes.util, subprocess
from ctypes import c_void_p, c_int32, byref, POINTER, c_long

cf = ctypes.cdll.LoadLibrary(ctypes.util.find_library('CoreFoundation'))
ax = ctypes.cdll.LoadLibrary(ctypes.util.find_library('ApplicationServices'))

ax.AXUIElementCreateApplication.restype = c_void_p
ax.AXUIElementCreateApplication.argtypes = [c_int32]
ax.AXUIElementCopyAttributeValue.restype = c_int32
ax.AXUIElementCopyAttributeValue.argtypes = [c_void_p, c_void_p, POINTER(c_void_p)]

cf.CFStringCreateWithCString.restype = c_void_p
cf.CFStringCreateWithCString.argtypes = [c_void_p, ctypes.c_char_p, c_int32]
cf.CFArrayGetCount.restype = c_long
cf.CFArrayGetCount.argtypes = [c_void_p]
cf.CFArrayGetValueAtIndex.restype = c_void_p
cf.CFArrayGetValueAtIndex.argtypes = [c_void_p, c_long]
cf.CFStringGetCString.restype = ctypes.c_bool
cf.CFStringGetCString.argtypes = [c_void_p, ctypes.c_char_p, c_long, c_int32]

def cfstr(s):
    return cf.CFStringCreateWithCString(None, s.encode(), 0x08000100)

def get_str(ref):
    if not ref: return None
    buf = ctypes.create_string_buffer(512)
    return buf.value.decode() if cf.CFStringGetCString(ref, buf, 512, 0x08000100) else '?'

def ax_get(elem, attr):
    val = c_void_p()
    err = ax.AXUIElementCopyAttributeValue(elem, cfstr(attr), byref(val))
    return err, val.value

# Usage: replace PID with target app's PID
pid = int(subprocess.run(['pgrep', 'MyApp'], capture_output=True, text=True).stdout.strip())
app = ax.AXUIElementCreateApplication(pid)

# Check role
err, val = ax_get(app, 'AXRole')
print(f'App role: {get_str(val)}')

# Check windows
err, val = ax_get(app, 'AXWindows')
if err == 0 and val:
    count = cf.CFArrayGetCount(val)
    print(f'AXWindows: {count}')

# Check children
err, val = ax_get(app, 'AXChildren')
if err == 0 and val:
    count = cf.CFArrayGetCount(val)
    print(f'AXChildren: {count}')
```

## Adding Internal Diagnostics (Obj-C Logging)

When debugging from inside the app, add this to `AccessibilitySetTargetWindow` or similar:

```objc
NSLog(@"[AX] NSWindow class = %@", NSStringFromClass([targetWindow class]));
NSLog(@"[AX] NSWindow title = %@", [targetWindow title]);
NSLog(@"[AX] NSWindow isVisible = %d", [targetWindow isVisible]);
NSLog(@"[AX] NSWindow accessibilityRole = %@", [targetWindow accessibilityRole]);
NSLog(@"[AX] NSWindow isAccessibilityElement = %d", [targetWindow isAccessibilityElement]);
NSLog(@"[AX] contentView class = %@", NSStringFromClass([cv class]));
NSLog(@"[AX] contentView accessibilityRole = %@", [cv accessibilityRole]);
NSLog(@"[AX] contentView isAccessibilityElement = %d", [cv isAccessibilityElement]);
NSLog(@"[AX] contentView accessibilityIsIgnored = %d", [cv accessibilityIsIgnored]);
NSLog(@"[AX] [NSApp windows] count = %lu", (unsigned long)[[NSApp windows] count]);
NSLog(@"[AX] NSApp activationPolicy = %ld", (long)[NSApp activationPolicy]);
```

Key things to look for:
- contentView role should be `AXGroup`, NOT `AXUnknown`
- contentView isAccessibilityElement should be `1` (YES)
- contentView accessibilityIsIgnored should be `0` (NO)
- NSApp activationPolicy should be `0` (NSApplicationActivationPolicyRegular)

## Relevant Files in Fyne

- `internal/driver/glfw/accessibility_darwin.go` — Go-side accessibility bridge
- `internal/driver/glfw/accessibility_darwin.h` — C header for AX bridge
- `internal/driver/glfw/accessibility_darwin.m` — Obj-C implementation with swizzling
- `accessibility.go` — Public accessibility API (roles, actions, states)
- `internal/driver/common/accessibility.go` — Common accessibility tree walking
