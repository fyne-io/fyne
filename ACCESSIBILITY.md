# Fyne Accessibility Support

Fyne provides optional accessibility support that exposes widgets to platform
accessibility APIs — enabling VoiceOver, System Events (AppleScript), Windows
UI Automation, and AI computer-use tools to discover and interact with Fyne
applications.

Accessibility is enabled at build time with the `accessibility` build tag:

```bash
go run -tags=accessibility .
go build -tags=accessibility -o myapp .
```

Without this tag, all accessibility code is compiled out and the application
behaves exactly as before.

> **Status**: Accessibility support is new in Fyne 2.8 and is under active
> development. macOS (via AppKit/NSAccessibility) and Windows (via
> UI Automation) are the primary desktop targets.

## Table of Contents

- [Quick Start](#quick-start)
- [Architecture](#architecture)
- [Public API](#public-api)
- [Making a Widget Accessible](#making-a-widget-accessible)
- [Platform Bridges](#platform-bridges)
- [Debugging Accessibility on macOS](#debugging-accessibility-on-macos)
- [Known Issues and Pitfalls](#known-issues-and-pitfalls)
- [File Reference](#file-reference)

## Quick Start

1. **Build with the accessibility tag:**

   ```bash
   go build -tags=accessibility -o myapp .
   ```

2. **Run the application.** Standard Fyne widgets (Button, Entry, Check, Slider,
   etc.) will automatically expose themselves to the platform accessibility
   system.

3. **Verify** with your platform's accessibility inspector:
   - **macOS**: Xcode → Accessibility Inspector (requires a `.app` bundle — see
     [Debugging on macOS](#debugging-accessibility-on-macos))
   - **Windows**: Accessibility Insights for Windows or Inspect.exe

## Architecture

Accessibility in Fyne has three layers:

```
┌─────────────────────────────────┐
│  Public API (accessibility.go)  │  Roles, Actions, States, Interfaces
├─────────────────────────────────┤
│  Tree Walker (common/)          │  Traverses widget tree → flat AX elements
├─────────────────────────────────┤
│  Platform Bridge (glfw/)        │  Obj-C / C / Win32 code that registers
│                                 │  elements with the OS accessibility API
└─────────────────────────────────┘
```

### Accessibility Tree Update Flow

When the window content changes, Fyne walks the widget tree and rebuilds the
accessibility element list:

1. `window.Show()` / content change triggers `window.updateAccessibility()`
2. `collectAccessibilityElements()` walks the tree recursively
3. For each `fyne.Accessible` widget, a native accessibility element is created
   with the widget's label, role, frame, actions, states, and value
4. Elements are attached to the window's content view via the platform bridge
5. Platform accessibility notifications are posted to inform assistive
   technologies of the new tree

### Build Tag Gating

All accessibility code is behind `//go:build accessibility && <platform>`.
When the tag is absent, stub no-op functions are compiled instead
(`accessibility_notdarwin.go`), ensuring zero overhead.

## Public API

All accessibility types and interfaces are defined in `accessibility.go`.

### Interfaces

| Interface | Methods | Purpose |
|-----------|---------|---------|
| `Accessible` | `AccessibilityLabel() string`, `AccessibilityRole() AccessibleRole` | **Required.** Provides the element's name and semantic role. |
| `AccessibleChildren` | `AccessibilityChildren() []CanvasObject` | Override which children appear in the AX tree (e.g., only the active tab's content). |
| `AccessibleValue` | `AccessibilityValue() string` | Expose a current value distinct from the label (slider position, entry text). |
| `AccessibleValueSetter` | `AccessibilitySetValue(value string) bool` | Allow assistive tech to set the value (voice dictation into a text field). |
| `AccessibleActions` | `AccessibilityActions() []AccessibleAction`, `AccessibilityPerformAction(AccessibleAction) bool` | Declare and handle user-triggered actions (press, increment, show menu). |
| `AccessibleStates` | `AccessibilityStates() []AccessibleState` | Report dynamic state flags (checked, disabled, expanded, focused). |

### Roles

| Role | Typical Widget |
|------|---------------|
| `AccessibleRoleButton` | Button |
| `AccessibleRoleCheckbox` | Check |
| `AccessibleRoleRadio` | RadioItem |
| `AccessibleRoleTextField` | Entry |
| `AccessibleRoleSlider` | Slider |
| `AccessibleRoleProgressBar` | ProgressBar |
| `AccessibleRoleText` | Label, RichText |
| `AccessibleRoleLink` | Hyperlink |
| `AccessibleRoleList` | List |
| `AccessibleRoleListItem` | List items |
| `AccessibleRoleTable` | Table |
| `AccessibleRoleTree` | Tree |
| `AccessibleRoleTreeItem` | Tree items |
| `AccessibleRoleTab` | Tab button |
| `AccessibleRoleTabList` | AppTabs, DocTabs |
| `AccessibleRoleContainer` | Generic container |
| `AccessibleRoleImage` | Icon, FileIcon |
| `AccessibleRoleSeparator` | Separator |
| `AccessibleRoleHeading` | Card title |

### Actions

| Action | Meaning |
|--------|---------|
| `AccessibleActionPress` | Click / tap |
| `AccessibleActionIncrement` | Increase value (slider) |
| `AccessibleActionDecrement` | Decrease value (slider) |
| `AccessibleActionSetValue` | Replace the widget's value |
| `AccessibleActionSelect` | Select this item |
| `AccessibleActionShowMenu` | Open a context menu |

### States

| State | Meaning |
|-------|---------|
| `AccessibleStateChecked` | Checkbox/radio is checked |
| `AccessibleStateDisabled` | Widget is disabled |
| `AccessibleStateExpanded` | Disclosure/tree node is expanded |
| `AccessibleStateFocused` | Widget has keyboard focus |
| `AccessibleStateSelected` | Item is selected |
| `AccessibleStateInvalid` | Validation failed |
| `AccessibleStateRequired` | Input is required |

## Making a Widget Accessible

### Minimum: Implement Accessible

Every accessible widget must implement `fyne.Accessible`:

```go
func (w *MyWidget) AccessibilityLabel() string {
    return "Descriptive label for assistive technology"
}

func (w *MyWidget) AccessibilityRole() fyne.AccessibleRole {
    return fyne.AccessibleRoleButton
}
```

### Adding Actions

For interactive widgets, implement `fyne.AccessibleActions`:

```go
func (w *MyWidget) AccessibilityActions() []fyne.AccessibleAction {
    return []fyne.AccessibleAction{fyne.AccessibleActionPress}
}

func (w *MyWidget) AccessibilityPerformAction(action fyne.AccessibleAction) bool {
    if action == fyne.AccessibleActionPress {
        w.handleTap()
        return true
    }
    return false
}
```

### Adding a Value

For widgets with a current value (entries, sliders, progress bars):

```go
func (w *MySlider) AccessibilityValue() string {
    return fmt.Sprintf("%.0f", w.Value)
}
```

To allow assistive tech to set the value:

```go
func (w *MyEntry) AccessibilitySetValue(value string) bool {
    w.SetText(value)
    return true
}
```

### Controlling the Child Tree

By default, the accessibility walker follows `Container.Objects`. To override
which children appear (e.g., only the visible tab content):

```go
func (t *MyTabs) AccessibilityChildren() []fyne.CanvasObject {
    if t.selected >= 0 && t.selected < len(t.items) {
        return []fyne.CanvasObject{t.items[t.selected].Content}
    }
    return nil
}
```

### Reporting State

```go
func (c *MyCheck) AccessibilityStates() []fyne.AccessibleState {
    var states []fyne.AccessibleState
    if c.Checked {
        states = append(states, fyne.AccessibleStateChecked)
    }
    if c.disabled {
        states = append(states, fyne.AccessibleStateDisabled)
    }
    return states
}
```

### Real Example: Button

From `widget/button.go`:

```go
func (b *Button) AccessibilityLabel() string {
    if b.Text != "" {
        return b.Text
    }
    if b.Icon != nil {
        return b.Icon.Name()
    }
    return ""
}

func (b *Button) AccessibilityRole() fyne.AccessibleRole {
    return fyne.AccessibleRoleButton
}

func (b *Button) AccessibilityActions() []fyne.AccessibleAction {
    return []fyne.AccessibleAction{fyne.AccessibleActionPress}
}

func (b *Button) AccessibilityPerformAction(action fyne.AccessibleAction) bool {
    if action == fyne.AccessibleActionPress && b.OnTapped != nil {
        b.OnTapped()
        return true
    }
    return false
}

func (b *Button) AccessibilityStates() []fyne.AccessibleState {
    if b.Disabled() {
        return []fyne.AccessibleState{fyne.AccessibleStateDisabled}
    }
    return nil
}
```

## Platform Bridges

### macOS (AppKit / NSAccessibility)

**Files:** `internal/driver/glfw/accessibility_darwin.{go,h,m}`

The macOS bridge creates `AccessibleElement` objects (an `NSAccessibilityElement`
subclass) for each Fyne widget that implements `fyne.Accessible`. Elements are
organized into a tree that mirrors the widget hierarchy:

```
NSApplication
 └── GLFWWindow (AXWindow)
      └── GLFWContentView (AXGroup)  ← swizzled to fix default AXUnknown role
           ├── AccessibleElement (AXButton, "Save")
           ├── AccessibleElement (AXTextField, "Name")
           └── AccessibleElement (AXGroup, "Panel")
                └── AccessibleElement (AXStaticText, "Hello")
```

Key implementation details:

- **Coordinate conversion**: Fyne uses top-left origin; macOS accessibility
  uses bottom-left (screen coordinates). The `accessibilityFrame` method on
  `AccessibleElement` performs the conversion dynamically.
- **Action routing**: Each element holds a `cgo.Handle` to its Go widget.
  When the OS invokes an accessibility action (press, increment, etc.), the
  Obj-C callback crosses the cgo boundary to call the widget's
  `AccessibilityPerformAction`.
- **GLFWContentView fixes**: GLFW's content view reports `AXUnknown` role and
  `isAccessibilityElement=NO` by default. The bridge uses `class_addMethod` to
  override these on the `GLFWContentView` class specifically (not the parent
  `NSView`), setting role to `AXGroup` and making it a visible accessibility
  element.
- **Legacy API support**: Custom `AccessibleElement` instances implement both
  the modern NSAccessibility protocol (`accessibilityRole`,
  `accessibilityChildren`, etc.) and the deprecated legacy API
  (`accessibilityAttributeValue:`, `accessibilityAttributeNames`,
  `accessibilityIsIgnored`) because some external AX clients — including
  System Events and VoiceOver — still query via the legacy path.

### Windows (UI Automation)

**Files:** `internal/driver/glfw/accessibility_windows.{go,h,c}`

The Windows bridge uses the UI Automation COM API to expose accessibility
elements to screen readers (Narrator, NVDA, JAWS) and automation tools.

### Mobile (iOS / Android)

**Files:** `internal/driver/mobile/accessibility_ios.{go,m}`,
`internal/driver/mobile/accessibility_android.go`

Mobile accessibility maps Fyne roles to platform-native traits
(`UIAccessibilityTraits` on iOS, `AccessibilityNodeInfo` on Android).

## Debugging Accessibility on macOS

> For a comprehensive AI-agent-oriented reference on macOS accessibility
> debugging, see the skill documentation in
> `.agents/skills/macos-accessibility/SKILL.md`.

### 1. Check Permissions First

The **querying process** (not the target app) must have Accessibility
permission. If `osascript` / Terminal returns 0 windows for ALL apps (including
Chrome), the problem is permissions:

```bash
# This should return >0 for Chrome if permissions are correct
osascript -e 'tell application "System Events" to return count of windows of process "Google Chrome"'
```

Fix: **System Settings → Privacy & Security → Accessibility** → add
Terminal.app (or whatever hosts your `osascript`).

### 2. Bundle for Accessibility Inspector

Xcode's Accessibility Inspector only lists `.app` bundles in its target
dropdown. A bare Go binary won't appear. Create a minimal bundle:

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

> **Note:** If the app requires a specific working directory (e.g., to find
> config files), you'll need a launcher wrapper. See the `GUIBundle` mage
> target in `copilot-api` for an example.

### 3. AppleScript Quick Checks

```bash
# Check if System Events sees the window
osascript -e 'tell application "System Events"
    tell process "MyApp"
        return "windows=" & (count of windows) & " elements=" & (count of UI elements)
    end tell
end tell'

# Get window title and element count
osascript -e 'tell application "System Events"
    tell process "MyApp"
        tell window 1
            return "title=" & title & " contents=" & (count of entire contents)
        end tell
    end tell
end tell'

# Bring window to front
osascript -e 'tell application "System Events"
    set frontmost of process "MyApp" to true
    tell process "MyApp"
        perform action "AXRaise" of window 1
    end tell
end tell'
```

### 4. Python AXUIElement API

For lower-level inspection that bypasses System Events caching:

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
cf.CFStringGetCString.restype = ctypes.c_bool
cf.CFStringGetCString.argtypes = [c_void_p, ctypes.c_char_p, c_long, c_int32]

def cfstr(s):
    return cf.CFStringCreateWithCString(None, s.encode(), 0x08000100)

def get_str(ref):
    if not ref: return None
    buf = ctypes.create_string_buffer(512)
    return buf.value.decode() if cf.CFStringGetCString(ref, buf, 512, 0x08000100) else '?'

pid = int(subprocess.run(['pgrep', 'MyApp'],
          capture_output=True, text=True).stdout.strip())
app = ax.AXUIElementCreateApplication(pid)

val = c_void_p()
ax.AXUIElementCopyAttributeValue(app, cfstr('AXRole'), byref(val))
print(f'Role: {get_str(val.value)}')

ax.AXUIElementCopyAttributeValue(app, cfstr('AXWindows'), byref(val))
if val.value:
    print(f'Windows: {cf.CFArrayGetCount(val.value)}')
```

### 5. Internal Diagnostic Logging

Add temporary logging to `AccessibilitySetTargetWindow` in
`accessibility_darwin.m` to verify the bridge state:

```objc
NSLog(@"[AX] contentView class=%@ role=%@ isAXElem=%d isIgnored=%d",
      NSStringFromClass([targetContentView class]),
      [targetContentView accessibilityRole],
      [targetContentView isAccessibilityElement],
      [targetContentView accessibilityIsIgnored]);
NSLog(@"[AX] window role=%@ children=%lu",
      [targetWindow accessibilityRole],
      (unsigned long)[[targetWindow accessibilityChildren] count]);
NSLog(@"[AX] app activationPolicy=%ld windows=%lu",
      (long)[NSApp activationPolicy],
      (unsigned long)[[NSApp accessibilityWindows] count]);
```

**What to look for:**
- Content view role must be `AXGroup` (not `AXUnknown`)
- Content view `isAccessibilityElement` must be `1` (YES)
- Content view `accessibilityIsIgnored` must be `0` (NO)
- Window `accessibilityChildren` count must be ≥ 1
- App `activationPolicy` must be `0` (regular foreground app)

### 6. Common Failure Modes

| Symptom | Cause | Fix |
|---------|-------|-----|
| System Events: "0 windows" for all apps | `osascript` lacks Accessibility permission | Grant permission in System Settings |
| System Events: "0 windows" for this app only | `method_setImplementation` on inherited method broke NSWindow's AX hierarchy | Use `class_addMethod` instead |
| Accessibility Inspector doesn't list the app | Running as bare binary, not `.app` bundle | Create a `.app` bundle (see above) |
| Elements visible but window not recognized | Content view returns `AXUnknown` role | Swizzle `accessibilityRole` on GLFWContentView |
| Elements visible but can't navigate between them | Missing `accessibilityWindow` / `accessibilityTopLevelUIElement` | Implement on AccessibleElement |
| Actions don't work from assistive tech | Legacy `accessibilityPerformPress` not wired | Ensure AccessibleElement implements legacy action methods |

## Known Issues and Pitfalls

### class_addMethod vs method_setImplementation

When adding accessibility overrides to GLFW classes, **always use
`class_addMethod`**, which adds the method only to the specific subclass.
Never use `method_setImplementation` on methods found via
`class_getInstanceMethod`, as this modifies the method on whatever parent
class defines it (typically `NSView` or `NSResponder`), breaking
accessibility for ALL views — including `NSWindow`.

```objc
// ✅ SAFE — adds only to GLFWContentView
class_addMethod([contentView class], @selector(accessibilityRole),
                (IMP)customRole, "@@:");

// ❌ DANGEROUS — modifies NSView/NSResponder globally
Method m = class_getInstanceMethod([contentView class], @selector(accessibilityRole));
method_setImplementation(m, (IMP)customRole);
```

### macOS Dual API Paths

macOS has two accessibility dispatch paths:

- **Modern** (`accessibilityRole`, `accessibilityChildren`, etc.) — used by
  some internal framework paths
- **Legacy** (`accessibilityAttributeValue:`, `accessibilityIsIgnored`,
  `accessibilityAttributeNames`) — used by external AX clients like System
  Events and VoiceOver

Custom `AccessibleElement` classes must implement **both** paths. The modern
methods are typically sufficient for the framework, but external clients may
only query the legacy path.

### NSView Defaults

`NSView` defaults to `accessibilityIsIgnored = YES`, which means subclasses
(including `GLFWContentView`) are invisible to the accessibility hierarchy
unless explicitly overridden. This is independent of the modern
`isAccessibilityElement` property.

### Coordinate System

macOS accessibility uses screen coordinates with **bottom-left origin**. Fyne
uses top-left origin internally. The `accessibilityFrame` method in
`AccessibleElement` must convert between these systems, taking into account
the window's position on screen and the content view's bounds.

## File Reference

| File | Purpose |
|------|---------|
| `accessibility.go` | Public API: roles, actions, states, interfaces |
| `accessibility_test.go` | Tests for the accessibility interfaces |
| `internal/driver/common/accessibility.go` | Shared tree-walking helper (`AccessibilityChildren`) |
| `internal/driver/glfw/accessibility_darwin.go` | macOS Go bridge: tree walker, cgo callbacks |
| `internal/driver/glfw/accessibility_darwin.h` | macOS C header: element types, action enums |
| `internal/driver/glfw/accessibility_darwin.m` | macOS Obj-C: element class, swizzling, AX notifications |
| `internal/driver/glfw/accessibility_darwin_test.go` | macOS accessibility tests |
| `internal/driver/glfw/accessibility_notdarwin.go` | No-op stubs when accessibility tag is absent |
| `internal/driver/glfw/accessibility_windows.go` | Windows Go bridge |
| `internal/driver/glfw/accessibility_windows.{h,c}` | Windows UI Automation bridge |
| `internal/driver/mobile/accessibility_ios.{go,m}` | iOS accessibility bridge |
| `internal/driver/mobile/accessibility_android.go` | Android accessibility bridge |
| `widget/accessibility_test.go` | Widget-level accessibility tests |
| `.agents/skills/macos-accessibility/SKILL.md` | AI-agent skill for diagnosing macOS AX issues |
