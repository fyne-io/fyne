Thanks very much for your interest in contributing to Fyne!
The community is what makes this project successful and we are glad to welcome you on board.

There are various ways to contribute, perhaps the following helps you know how to get started.

## Reporting a bug

If you've found something wrong we want to know about it, please help us understand the problem so we can resolve it.

1. Check to see if this already is recorded, if so add some more information [issue list](https://github.com/fyne-io/fyne/issues)
2. If not then create a new issue using the [bug report template](https://github.com/fyne-io/fyne/issues/new?assignees=&labels=&template=bug_report.md&title=)
3. Stay involved in the conversation on the issue as it is triaged and progressed.


## Fixing an issue

Great! You found an issue and figured you can fix it for us.
If you can follow these steps then your code should get accepted fast.

1. Read through the "Contributing Code" section further down this page.
2. Write a unit test to show it is broken.
3. Create the fix and you should see the test passes.
4. Run the tests and make sure everything still works as expected using `go test ./...`.
5. [Open a PR](https://github.com/fyne-io/fyne/compare) and work through the review checklist.


## Adding a feature

It's always good news to hear that people want to contribute functionality.
But first of all check that it fits within our [Vision](https://github.com/fyne-io/fyne/wiki/Vision) and if we are already considering it on our [Roadmap](https://github.com/fyne-io/fyne/wiki/Roadmap).
If you're not sure then you should join our #fyne-contributors channel on the [Gophers Slack server](https://gophers.slack.com/app_redirect?channel=fyne-contributors).

Once you are ready to code then the following steps should give you a smooth process:

1. Read through the [Contributing Code](#contributing-code) section further down this page.
2. Think about how you would structure your code and how it can be tested.
3. Write some code and enjoy the ease of writing Go code for even a complex project :).
4. Run the tests and make sure everything still works as expected using `go test ./...`.
5. [Open a PR](https://github.com/fyne-io/fyne/compare) and work through the review checklist.


# Contributing Code

We aim to maintain a very high standard of code, through design, test and implementation.
To manage this we have various checks and processes in place that everyone should follow, including:

* We use the Go standard format (with tabs not spaces) - you can run `gofmt` before committing
* Imports should be ordered according to the GoImports spec - you can use the `goimports` tool instead of `gofmt`.
* Everything should have a unit test attached (as much as possible, to keep our coverage up)

For desktop event-loop changes, run the main-queue and scheduler tests with
`go test -race -tags no_glfw,ci,migrated_fynedo ./internal/async ./internal/driver/glfw`.
On macOS, add `accessibility` to check native element identity and in-place
attribute updates without requiring Accessibility permission. The tagged suite
also launches a real desktop fixture to check queued `fyne.Do` work, window
teardown, menus, dialogs, shutdown, and rendering cadence.
The native accessibility probes also check ownership, window teardown, and
isolation between simultaneous windows and different content-view classes.
The unchanged-tree benchmark measures synchronization of collected nodes, not
the complete canvas traversal or repaint.

To measure actual AX request latency, build the same fixture on each revision:

```sh
go build -tags accessibility,migrated_fynedo -ldflags=-linkmode=external -o /tmp/fyne-main-loop ./internal/driver/glfw/testdata/main_loop
clang -framework Cocoa -framework ApplicationServices internal/driver/glfw/testdata/accessibility_query.m -o /tmp/fyne-ax-query
/tmp/fyne-main-loop -duration=30s -elements=16 -repaint -repaint-delay=1s
```

After the fixture prints `READY <pid>`, run `/tmp/fyne-ax-query <pid> 200`
from another terminal. The query process (or its launching application) needs
macOS Accessibility permission. The probe reports median, p95, and maximum
window-attribute and hit-test latency, plus changes or invalidation of a retained
text element and the number of registered roots. The repaint delay gives the
client time to retain an element before repainting starts. Repeat without
`-repaint` and with `-elements=256`. The fixture reports
animation frames, accessibility updates, and CPU time; its `-exercise` flag also
checks queued work, menus, dialogs, and a second window.

For detailed Code style, check [Contributing](https://github.com/fyne-io/fyne/wiki/Contributing#code-style) in our wiki please.

# Decision Process

The following points apply to our decision making process:

* Any decisions or votes will be opened on the #fyne-votes Slack channel and follows lazy consensus.
* Any contributors not responding in 4 days will be deemed in agreement.
* Any PR that has not been responded to within 7 days can be automatically approved.
* No functionality will be added unless at least 2 developers agree it belongs.

Bear in mind that this is a cross platform project so any new features would normally
be required to work on multiple desktop and mobile platforms.
