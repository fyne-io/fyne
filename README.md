<p align="center">
  <a href="https://travis-ci.org/fyne-io/fyne"><img src="https://travis-ci.org/fyne-io/fyne.svg" alt="Build Status" /></a>
  <a href="https://goreportcard.com/report/github.com/fyne-io/fyne"><img src="https://goreportcard.com/badge/github.com/fyne-io/fyne" alt="Code Status" /></a>
  <a href='https://coveralls.io/github/fyne-io/fyne?branch=develop'><img src='https://coveralls.io/repos/github/fyne-io/fyne/badge.svg?branch=develop' alt='Coverage Status' /></a>
</p>

# About

[Fyne](http://fyne.io) is an easy to use UI toolkit and app API written in Go. We use the EFL render pipeline to provide cross platform graphics.

This is under heavy development and is not yet capable of supporting a full application

# Prerequisites

Before you can use the Fyne tools you need to have a stable copy of EFL installed. This will be automated by our [bootstrap](https://github.com/fyne-io/bootstrap/) scripts soon, but for now you need to install at least EFL 1.20 from your package manager or the [Enlightenment downloads](https://download.enlightenment.org/rel/libs/efl/) page.

# Getting Started

Using standard go tools you can install Fyne's core library using:

    go get github.com/fyne-io/fyne-app

And then you're ready to write your first app - this example shows how:

    package main

    import "github.com/fyne-io/fyne/ui"
    import "github.com/fyne-io/fyne/ui/event"
    import "github.com/fyne-io/fyne/ui/layout"
    import "github.com/fyne-io/fyne/ui/widget"
    import "github.com/fyne-io/fyne-app"

    func main() {
    	app := fyneapp.NewApp()

    	w := app.NewWindow("Hello")
    	quit := widget.NewButton("Quit", func(*event.MouseEvent) {
    		app.Quit()
    	})
    	w.Canvas().SetContent(ui.NewContainer(
    		[]ui.CanvasObject{
    			widget.NewLabel("Hello Fyne!"),
    			quit,
    		},
    		layout.NewGridLayout(1)))

    	w.Show()
    }

And you can run that simply as:

    go run main.go

It should look like this:

<p align="center" markdown="1">
  <img src="hello.png" alt="Fyne Screenshot" />
</p>
