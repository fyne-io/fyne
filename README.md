# About

[Fyne](http://fyne.io) is an easy to use UI toolkit and app API written in Go. We use the EFL render pipeline to provide cross platform graphics.

This is under heavy development and is not yet capable of supporting a full application

# Prerequisites

Before you can use the Fyne tools you need to have a stable copy of EFL installed. This will be automated by our [bootstrap](https://github.com/fyne-io/bootstrap/) scripts soon, but for now you need to install at least EFL 1.20 from your package manager or the [Enlightenment downloads](https://download.enlightenment.org/rel/libs/efl/) page.

# Getting Started

Using standard go tools you can install Fyne's core library using:

    go get github.com/fyne-io/fyne

And then you're ready to write your first app - this example shows how:

    package main

    import "github.com/fyne-io/fyne/app"

    func main() {
        w := app.NewWindow("Hello")
        w.Canvas().NewText("Hello Fyne!")

        app.Run()
    }

And you can run that simply as:

    go run main.go

