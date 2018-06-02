<p align="center">
  <a href="https://godoc.org/github.com/fyne-io/fyne" title="GoDoc Reference" rel="nofollow"><img src="https://img.shields.io/badge/go-documentation-blue.svg?style=flat" alt="GoDoc Reference"></a>
  <a href="https://goreportcard.com/report/github.com/fyne-io/fyne"><img src="https://goreportcard.com/badge/github.com/fyne-io/fyne" alt="Code Status" /></a>
  <a href="https://travis-ci.org/fyne-io/fyne"><img src="https://travis-ci.org/fyne-io/fyne.svg" alt="Build Status" /></a>
  <a href='https://coveralls.io/github/fyne-io/fyne?branch=develop'><img src='https://coveralls.io/repos/github/fyne-io/fyne/badge.svg?branch=develop' alt='Coverage Status' /></a>
</p>

# About

[Fyne](http://fyne.io) is an easy to use UI toolkit and app API written in Go. We use the EFL render pipeline to provide cross platform graphics.

This is under heavy development and is not yet capable of supporting a full application

# Prerequisites

Before you can use the Fyne tools you need to have a stable copy of EFL installed. This will be automated by our [bootstrap](https://github.com/fyne-io/bootstrap/) scripts soon, but for now you can follow our [setup instructions](https://github.com/fyne-io/bootstrap/blob/master/README.md).

# Getting Started

Using standard go tools you can install Fyne's core library using:

    go get github.com/fyne-io/fyne

And then you're ready to write your first app - this example shows how:

    package main

    import "github.com/fyne-io/fyne/api/ui/widget"
    import "github.com/fyne-io/fyne/desktop"

    func main() {
    	app := desktop.NewApp()

    	w := app.NewWindow("Hello")
    	quit := widget.NewButton("Quit", func() {
    		app.Quit()
    	})
    	w.Canvas().SetContent(widget.NewList(
    		widget.NewLabel("Hello Fyne!"),
    		quit))

    	w.Show()
    }

And you can run that simply as:

    go run main.go

It should look like this:

<p align="center" markdown="1">
  <img src="hello.png" alt="Fyne Screenshot" />
</p>

# Examples

To see the examples you can run examples/main.go and optionally specify an example, like this:

    cd ~/go/src/github.com/fyne-io/fyne/examples/
    go run main.go -example calculator

It should look like one of these:

|       | Linux | Max OS X | Windows |
| -----:|:-----:|:--------:|:-------:|
|  dark | ![Calculator on Linux](img/calc-linux-dark.png) | ![Calculator on OS X](img/calc-osx-dark.png) | ![Calculator on Windows](img/calc-windows-dark.png) |
| light | ![Calculator (light) on Linux](img/calc-linux-light.png) | ![Calculator (light) on OS X](img/calc-osx-light.png) | ![Calculator (light) on Windows](img/calc-windows-light.png) |

## Clock

    go run main.go -example clock

![Clock dark](img/clock-dark.png)
![Clock light](img/clock-light.png)

## Fractal (Mandelbrot)

    go run main.go -example fractal

![Fractal](img/fractal-dark.png)
