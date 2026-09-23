# Accessibility gallery and macOS integration tests

The gallery uses real Fyne widgets and deterministic data. It is useful with
VoiceOver, Accessibility Inspector, and the external AX test runner below.
Actions update status labels, so tests check application behavior as well as
native attributes.

## Explore the gallery

From the repository root:

```sh
go run -tags accessibility ./cmd/fyne_accessibility
go run -tags accessibility ./cmd/fyne_accessibility -scenario controls
```

The `all` screen has tabs for `controls`, `text`, `collections`, `containers`,
and `dynamic`. A scenario starts with fresh state on every launch.

| Screen | Coverage |
| --- | --- |
| Controls | Buttons, disabled buttons, checks, radio/check groups, select menus, slider, progress indicators, activity, icon, separator |
| Text | Labels, rich text, hyperlink, text grid, single-line/multiline/disabled/password entries, validated form |
| Collections | Virtualized list, tree, table with headers, grid wrap, programmatic scrolling |
| Containers | App tabs, document tabs, split, accordion, toolbar |
| Dynamic | Show/hide, enable/disable, dialog, popup |

The hyperlink updates a label instead of opening a browser. The password is
dummy fixture data. The indeterminate progress bar is stopped to avoid constant
animation during tree traversal. No preferences or external data are needed to
operate the scenarios.

## Permission-free tests

These exercise gallery behavior with Fyne's test driver on any supported desktop
platform:

```sh
go test -race -tags ci,migrated_fynedo ./cmd/fyne_accessibility ./widget ./container -run 'Gallery|Accessibility'
```

On macOS, the existing bridge mapping tests can run without a GUI or an
Accessibility grant:

```sh
go test -race -tags accessibility,no_glfw,ci,migrated_fynedo ./internal/driver/glfw \
  -run 'Test(AccessibilityRoleDescription|AccessibilityHitTest|RoleToC|StateMaskFor|ActionMaskFor|ActionFromC)'
```

The macOS platform CI job includes `accessibility` so these bridge tests are
compiled and executed. Isolated native probes check modern and legacy role
descriptions for every bridge role and hit testing through nested containers,
catching recursion and incorrect pointer targets without opening a window.
It does **not** enable `axintegration`.

## External macOS tests

Requirements:

- macOS with a logged-in, unlocked graphical desktop.
- Go and Xcode Command Line Tools (the runner uses cgo and Apple frameworks).
- Accessibility permission for the runner or its launching terminal.

Build both executables into stable, git-ignored paths from the repository root:

```sh
mkdir -p test/accessibility/bin
go build -tags accessibility -o test/accessibility/bin/gallery ./cmd/fyne_accessibility
go test -c -tags accessibility,axintegration \
  -o test/accessibility/bin/ax.test ./test/accessibility
```

Open **System Settings > Privacy & Security > Accessibility** and grant access
to the terminal launching the tests, or add the compiled `ax.test` executable.
The runner checks `AXIsProcessTrusted` and fails with setup instructions if
permission is missing. It does not change permissions or silently skip tests.
Depending on how macOS identifies the executable, rebuilding it can require
renewing the grant. Avoid `go test`'s temporary executable path for regular
external runs; use the stable compiled runner instead.

Run:

```sh
test/accessibility/bin/ax.test -test.v -test.timeout=5m \
  -gallery "$PWD/test/accessibility/bin/gallery"
```

Run one scenario:

```sh
test/accessibility/bin/ax.test -test.v -test.run 'TestAXGallery/dynamic' \
  -gallery "$PWD/test/accessibility/bin/gallery"
```

The tests launch and stop their own gallery processes, query only those PIDs,
and use `AXUIElement` to discover windows, read attributes, perform actions, and
set values. They do not send keyboard/mouse events, inspect unrelated apps, or
take screenshots. Gallery windows remain visible to accessibility clients, but
the runner passes `-background`, which disables GLFW
focus-on-creation and focus-on-show and makes the test windows transparent to
mouse input. You can keep using other applications while the suite runs:
accessibility actions do not require the gallery to be the foreground app.
The runner checks `AXFrontmost` around tree assertions and fails if the gallery
becomes active, including if you deliberately activate it during a run. It never
activates another app to "restore" stolen focus.
Launching the gallery manually without `-background` retains normal interactive
window behavior. Each scenario has a process timeout and cleanup. The
gallery also has a `-quit-after` watchdog so it exits even if the runner is
terminated before its cleanup executes.

Checks include native roles, labels, nonzero bounds, window association,
checked/selected/expanded/disabled states, text editing, password concealment,
action callbacks, tab visibility, collection descendants, and transient
dialogs/menus/popups. The runner polls for expected state and reacquires native
elements after changes rather than relying on fixed startup sleeps or stable
native object identities. Failures include the last complete accessibility
tree and the gallery's output.

Accessibility Inspector currently needs its target reselected after a repaint:
the bridge rebuilds native elements, invalidating Inspector's selection. These
tests reacquire elements and do not establish stable element identity across
updates.
Pointer targeting is covered by the native hit-test probe and interactive
Inspector checks, not this runner: macOS application-wide hit testing skips the
mouse-transparent background windows.

These tests are excluded from normal builds by
`darwin && cgo && accessibility && axintegration`. Do not add `ci` or `no_glfw`
when building the gallery for these tests: it must use the real desktop driver.
Hosted CI is not assumed to have a usable desktop or an Accessibility grant.
A dedicated self-hosted Mac can run the same commands after explicit permission
setup; do not disable macOS privacy protections to automate setup.

## Adding coverage

Add deterministic widgets to `cmd/fyne_accessibility/gallery.go`, give controls
distinct labels/placeholders, and expose callback results through a status
label. Add permission-free assertions in `gallery_test.go`, then an external
assertion/action in `gallery_darwin_test.go`. Keep native expectations semantic:
do not depend on node indices, child ordering, exact pixel coordinates, or
macOS-specific wrapper counts.

The suite is not a replacement for checking VoiceOver announcement quality,
keyboard navigation usability, every widget state, or multi-window behavior.
It intentionally uses one window per scenario to keep failures isolated.
