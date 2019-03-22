<p align="center">
  <a href="https://godoc.org/fyne.io/fyne" title="GoDoc Reference" rel="nofollow"><img src="https://img.shields.io/badge/go-documentation-blue.svg?style=flat" alt="GoDoc Reference"></a>
  <a href="https://github.com/fyne-io/fyne/releases/tag/v1.0.0" title="1.0.0 Release" rel="nofollow"><img src="https://img.shields.io/badge/version-1.0.0-blue.svg?style=flat" alt="1.0.0 release"></a>
  <a href='https://github.com/avelino/awesome-go'><img src='https://awesome.re/mentioned-badge.svg' alt='Mentioned in Awesome' /></a>
  <a href='http://gophers.slack.com/messages/fyne'><img src='https://img.shields.io/badge/join-us%20on%20slack-gray.svg?longCache=true&logo=slack&colorB=blue' alt='Join us on Slack' /></a>
  <br />
  <a href="https://goreportcard.com/report/fyne.io/fyne"><img src="https://goreportcard.com/badge/fyne.io/fyne" alt="Code Status" /></a>
  <a href="https://travis-ci.org/fyne-io/fyne"><img src="https://travis-ci.org/fyne-io/fyne.svg" alt="Build Status" /></a>
  <a href='https://coveralls.io/github/fyne-io/fyne?branch=develop'><img src='https://coveralls.io/repos/github/fyne-io/fyne/badge.svg?branch=develop' alt='Coverage Status' /></a>
  <!--a href='https://sourcegraph.com/github.com/fyne-io/fyne?badge'><img src='https://sourcegraph.com/github.com/fyne-io/fyne/-/badge.svg' alt='Used By' /></a-->
</p>

# About

[Fyne](http://fyne.io) is an easy to use UI toolkit and app API written in Go. We use OpenGL (through the go-gl and go-glfw projects) to provide cross platform graphics.

The 1.0 release is now out and we encourage feedback and requests for the next major releae :).

# Getting Started

Fyne is designed to be really easy to code with, here are the steps to your first app.

## Prerequisites

As Fyne uses CGo you will require a C compiler (typically gcc).
If you don't have one set up the instructions at [Compiling](https://github.com/fyne-io/fyne/wiki/Compiling) may help.

By default Fyne uses the [gl golang bindings](https://github.com/go-gl/gl) which means you need a working OpenGL configuration.
Debian/Ubuntu based systems may need to also need to install the `libgl1-mesa-dev` and `xorg-dev` packages.

Using the standard go tools you can install Fyne's core library using:

    go get fyne.io/fyne

## Code

And then you're ready to write your first app!

```go
package main

import (
	"fyne.io/fyne/widget"
	"fyne.io/fyne/app"
)

func main() {
	app := app.New()

	w := app.NewWindow("Hello")
	w.SetContent(widget.NewVBox(
		widget.NewLabel("Hello Fyne!"),
		widget.NewButton("Quit", func() {
			app.Quit()
		}),
	))

	w.ShowAndRun()
}
```

And you can run that simply as:

    go run main.go

It should look like this:

<p align="center" markdown="1">
  <img src="img/hello-normal.png" width="207" height="204" alt="Fyne Hello Dark Theme" />
</p>

> Note that windows applications load from a command prompt by default, which means if you click an icon you may see a command window.
> To fix this add the parameters `-ldflags -H=windowsgui` to your run or build commands.

# Scaling

Fyne is built entirely using vector graphics, which means applications
written with Fyne will scale to any value beautifully (not just whole number values).
The default scale value is calculated from your screen's DPI - and if you move
a window to another screen it will re-scale and adjust the window size accordingly!
We call this "auto scaling", and it is designed to keep an app GUI the same size as you change monitor.
You can override this behaviour by setting a specific scale using the FYNE_SCALE environment variable.

<table style="text-align: center"><tr>
<td><img src="img/hello-normal.png" width="207" height="204" alt="Hello normal size" />
  <br />Standard size</td>
<td><img src="img/hello-small.png" width="160" height="169" alt="Hello small size" />
  <br />FYNE_SCALE=0.5</td>
<td><img src="img/hello-large.png" width="350" height="309" alt="Hello large size" />
  <br />FYNE_SCALE=2.5</td>
</tr></table>

# Themes

Fyne ships with two themes by default, "light" and "dark". You can choose
which to use with the environment variable ```FYNE_THEME```.
The default is dark:

<p align="center" markdown="1">
  <img src="img/hello-normal.png" width="207" height="204" alt="Fyne Hello Dark Theme" />
</p>

If you prefer a light theme then you could run:

    FYNE_THEME=light go run main.go

It should then look like this:

<p align="center" markdown="1">
  <img src="img/hello-light.png" width="207" height="204" alt="Fyne Hello Light Theme" />
</p>

# Widget demo

To run a showcase of the features of fyne execute the following:

    cd $GOPATH/src/fyne.io/fyne/cmd/fyne_demo/
    go build
    ./fyne_demo

And you should see something like this (after you click a few buttons):

<p align="center" markdown="1" style="max-width: 100%">
  <img src="img/widgets-dark.png" alt="Fyne Hello Light Theme" style="max-width: 100%" />
</p>

Or if you are using the light theme:

<p align="center" markdown="1" style="max-width: 100%">
  <img src="img/widgets-light.png" alt="Fyne Hello Light Theme" />
</p>

# Declarative API

If you prefer a more declarative API then that is provided too.
The following is exactly the same as the code above but in this different style.

```go
package main

import (
	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/widget"
)

func main() {
	app := app.New()

	w := app.NewWindow("Hello")
	w.SetContent(&widget.Box{Children: []fyne.CanvasObject{
		&widget.Label{Text: "Hello Fyne!"},
		&widget.Button{Text: "Quit", OnTapped: func() {
			app.Quit()
		}},
	}})

	w.ShowAndRun()
}
```

# Examples

The main examples have been moved - you can find them in their [own repository](https://github.com/fyne-io/examples/).

