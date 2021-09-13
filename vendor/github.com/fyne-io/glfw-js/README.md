glfw
====

[![Build Status](https://travis-ci.org/goxjs/glfw.svg?branch=master)](https://travis-ci.org/goxjs/glfw) [![GoDoc](https://godoc.org/github.com/goxjs/glfw?status.svg)](https://godoc.org/github.com/goxjs/glfw)

Package glfw experimentally provides a glfw-like API
with desktop (via glfw) and browser (via HTML5 canvas) backends.

It is used for creating a GL context and receiving events.

**Note:** This package is currently in development. The API is incomplete and may change.

Installation
------------

```bash
go get -u github.com/goxjs/glfw
GOARCH=js go get -u -d github.com/goxjs/glfw
```

Directories
-----------

| Path                                                               | Synopsis                                                           |
|--------------------------------------------------------------------|--------------------------------------------------------------------|
| [test/events](https://godoc.org/github.com/goxjs/glfw/test/events) | events hooks every available callback and outputs their arguments. |

License
-------

-	[MIT License](https://opensource.org/licenses/mit-license.php)
