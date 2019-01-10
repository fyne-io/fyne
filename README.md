<p align="center">
  <a href="https://godoc.org/fyne.io/fyne" title="GoDoc Reference" rel="nofollow"><img src="https://img.shields.io/badge/go-documentation-blue.svg?style=flat" alt="GoDoc Reference"></a>
  <a href="https://goreportcard.com/report/fyne.io/fyne"><img src="https://goreportcard.com/badge/fyne.io/fyne" alt="Code Status" /></a>
  <a href="https://travis-ci.org/fyne-io/fyne"><img src="https://travis-ci.org/fyne-io/fyne.svg" alt="Build Status" /></a>
  <a href='https://coveralls.io/github/fyne-io/fyne?branch=develop'><img src='https://coveralls.io/repos/github/fyne-io/fyne/badge.svg?branch=develop' alt='Coverage Status' /></a>
  <a href='https://sourcegraph.com/github.com/fyne-io/fyne?badge'><img src='https://sourcegraph.com/github.com/fyne-io/fyne/-/badge.svg' alt='Used By' /></a>
  <a href='https://github.com/avelino/awesome-go'><img src='https://awesome.re/mentioned-badge.svg' alt='Mentioned in Awesome' /></a>
  <a href='http://gophers.slack.com/messages/fyne'><img src='https://img.shields.io/badge/join-us%20on%20slack-gray.svg?longCache=true&logo=slack&colorB=blue' alt='Join us on Slack' /></a>
</p>

# About

[Fyne](http://fyne.io) is an easy to use UI toolkit and app API written in Go. We use the EFL render pipeline to provide cross platform graphics.

This is under heavy development and is not yet capable of supporting a full application

# Getting Started

Fyne is designed to be really easy to code with, here are the steps to your first app.

## Prerequisites

You will require a cgo compiler (typically gcc). On Windows this can be included by installing TDM-GCC from http://tdm-gcc.tdragon.net/download

Fyne uses the [gl golang bindings](https://github.com/go-gl/gl) you may need to also need to install the `libgl1-mesa-dev` package on Debian/Ubuntu based systems.

Using the standard go tools you can install Fyne's core library using:

    go get fyne.io/fyne

## Code

And then you're ready to write your first app!

```go
    package main

    import "fyne.io/fyne/widget"
    import "fyne.io/fyne/app"

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
  <img src="img/hello-normal.png" alt="Fyne Hello Dark Theme" />
</p>

# Scaling

Fyne is built entirely using vector graphics which means that applications
that are written using it will scale to any value beautifully (not just whole number values).
The default scale value is equated from your screen's DPI - and if you move
a window to another screen it will re-calculate and adjust the window size accordingly!

<table style="text-align: center"><tr>
<td><img src="img/hello-normal.png" alt="Hello normal size" />
  <br />Standard size</td>
<td><img src="img/hello-small.png" alt="Hello small size" />
  <br />FYNE_SCALE=0.5</td>
<td><img src="img/hello-large.png" alt="Hello large size" />
  <br />FYNE_SCALE=2.5</td>
</tr></table>

# Themes

Fyne ships with two themes by default, "light" and "dark". You can choose
which to use with the environment variable ```FYNE_THEME```.
The default is dark:

<p align="center" markdown="1">
  <img src="cmd/hello/hello-dark.png" alt="Fyne Hello Dark Theme" />
</p>

If you prefer a light theme then you could run:

    FYNE_THEME=light go run main.go

It should then look like this:

<p align="center" markdown="1">
  <img src="cmd/hello/hello-light.png" alt="Fyne Hello Light Theme" />
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

