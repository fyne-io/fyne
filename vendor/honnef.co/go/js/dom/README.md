# js/dom

Package dom provides Go bindings for the JavaScript DOM APIs.

## Version 2

**API Status:** Alpha, more API changes may be done soon

Version 2 of package dom is implemented on top of the [`syscall/js`](https://godoc.org/syscall/js) API and supports both Go WebAssembly and GopherJS.

It provides an API that is as close as possible to v1, with the following neccessary changes:

- All struct fields with `js:"foo"` tags have been replaced with equivalent methods
- `Underlying()` returns `js.Value` instead of `*js.Object`
- `AddEventListener()` returns `js.Func` instead of `func(*js.Object)`

### Install

    go get honnef.co/go/js/dom/v2

### Documentation

For documentation, see https://godoc.org/honnef.co/go/js/dom/v2.

## Version 1

**API Status:** Stable, changes only due to DOM being a moving target

Version 1 of package dom is implemented on top of the [`github.com/gopherjs/gopherjs/js`](https://godoc.org/github.com/gopherjs/gopherjs/js) API and supports GopherJS only.

### Install

    go get honnef.co/go/js/dom

### Documentation

For documentation, see https://godoc.org/honnef.co/go/js/dom.
