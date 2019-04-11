<p align="center">
  <a href="https://godoc.org/fyne.io/fyne" title="GoDoc Reference" rel="nofollow"><img src="https://img.shields.io/badge/go-documentation-blue.svg?style=flat" alt="GoDoc Reference"></a>
  <a href="https://github.com/fyne-io/fyne/releases/tag/v1.0.0" title="1.0.0 Release" rel="nofollow"><img src="https://img.shields.io/badge/version-1.0.0-blue.svg?style=flat" alt="1.0.0 release"></a>
  <a href='http://gophers.slack.com/messages/fyne'><img src='https://img.shields.io/badge/join-us%20on%20slack-gray.svg?longCache=true&logo=slack&colorB=blue' alt='Join us on Slack' /></a>
  <a href='https://fossfi.sh/support-fyneio'><img src='https://img.shields.io/badge/$-support_us-orange.svg?labelWidth=20&logo=data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAADAAAAAkCAYAAADPRbkKAAAABmJLR0QA7wAyAD/CTveyAAAACXBIWXMAAA9hAAAPYQGoP6dpAAAAB3RJTUUH4wMVCQ4LeuPReAAABDFJREFUWMPVmX9oVWUYxz/3tr5jbX9oYiZWSmNgoYYlUVqraYvIamlCg8pqUBTZr38qMhZZ0A8oIqKIiLJoJtEIsR+zki1mGqQwZTYkImtlgk5Wy7jPnbv90XPgcD3n3Lvdy731wOGc+77nec73Oe/zfp/nOTeVy+WolphZC3CVH9N9OAsMAzuBLkkHkmykKu2AmQGsAjYAC4pQ+QJ4WNJg1R0ws3pgI3DTJFVPAM8AGyRNVMUBB98DLCvBzLtAh6QTwUC6gmGzuUTwAGuBR8MD6QpFzwPAypi5v4HXgRXAEg+vLQm2njSzCyoWQmbWAPwAzIqY/gm4UdJAhN5qoAuojdCbAD4BVlViBdbGgJ8A2gPwZtZsZvea2XwASd1AZ4zNNHA90FgJB9ri6FHStw7+WaAPeA3YbWaBzovAiF/nfMUC+QD4sRIOLIkZ73Hws4BHQuOneY7A2aY/CHdgL3AtsBhYD9TVVMCB0+O2h59nRJBJfej6WOj6Bj8C2VcpFoqSC/38PbArb24zQDabBWhOsFFbNAuZ2UygxZcvcHwPsEPScILecaAuYmoUaJR01MymAU8AFzm7vCxp3MxagW0xpv8EzinogLPC40A7cGpMKHQB6yX9FqE/ACyKMb8VWC0pG6E3D+gF5sbo9km6Ml0A/N3Ad8BtMeABBNwB9JrZ7Ij5PQmPuA7YZWZXmFnan1ljZrf65p2boDuQmInNrAN4I29DJUmT02C+bCtiL/QCo2Y2BIwB7wFzCuj1xGZiM7vcb9gOTAMunUTZMVPSkZCtJuBAmQngIHCupImaCPApoBVokvSrj7UDm4o0Pg84Evo97EkoVUYHXgjK6pNWwB1I5dfdZtbtjUghWSppZ0ivDvirjA4MAosCfCeFhaRcPvhQLV6QbYGhvLHZZQRvwO1hfJNJZFu9qkySNyUdyxu7pkzgc8BDknZPuSMzs8XA58AZEdPbgTZJY3ld2GABOixWXpH0YMk9sZlN96zZ6hn2F+AdYFM4IZlZLfBRQiMzGXlaUmdFmnpPSC3A814alCJjwDpJG+NuqCkC0HzgLi+L651RvgTeAu4ELnHazAJnedkwpwzvoh+4RdLPid+FMpnMDElHY8CvAd4GGiKmR4A1DvwpYHkZabLTO7KCkgYeiwG/3Iu0hoQ6fwvwh6QVwFKn2tEpgB4BPgSaJS0oFnywAnuB5yR15bHNV6HPfUnyqaSVIV0BlwEXe51zNtAYehGHgMNOyfu9F+iXND6V5UplMplXgfuAj4HPgIXAPcXsj5CcKelwNbqiVCaTOd/jrhRZJumbajiQlrQf6OZ/KkEpsQ74vQQ7w1V1QNIh4Gbg+BRs9BXi6kqsAJK+Bq52hihWxoH7/wshFDixw786vF9kmm+TtK+aDsTWQma2EOjg379/zgNO8akh73NfknSw2pv4H3Ayg0FmbTMRAAAAAElFTkSuQmCC' alt='Support Fyne.io' /></a>
  <br />
  <a href="https://goreportcard.com/report/fyne.io/fyne"><img src="https://goreportcard.com/badge/fyne.io/fyne" alt="Code Status" /></a>
  <a href="https://travis-ci.org/fyne-io/fyne"><img src="https://travis-ci.org/fyne-io/fyne.svg" alt="Build Status" /></a>
  <a href='https://coveralls.io/github/fyne-io/fyne?branch=develop'><img src='https://coveralls.io/repos/github/fyne-io/fyne/badge.svg?branch=develop' alt='Coverage Status' /></a>
  <!--a href='https://sourcegraph.com/github.com/fyne-io/fyne?badge'><img src='https://sourcegraph.com/github.com/fyne-io/fyne/-/badge.svg' alt='Used By' /></a-->
</p>

# About

[Fyne](http://fyne.io) is an easy to use UI toolkit and app API written in Go. We use OpenGL (through the go-gl and go-glfw projects) to provide cross platform graphics.

The 1.0 release is now out and we encourage feedback and requests for the next major release :).

# Widget demo

To run a showcase of the features of Fyne execute the following:

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

# Getting Started

Fyne is designed to be really easy to code with. Here are the steps to your first app.

## Prerequisites

As Fyne uses CGo you will require a C compiler (typically gcc).
If you don't have one set up the instructions at [Compiling](https://github.com/fyne-io/fyne/wiki/Compiling) may help.

By default Fyne uses the [gl golang bindings](https://github.com/go-gl/gl) which means you need a working OpenGL configuration.
Debian/Ubuntu based systems may also need to install the `libgl1-mesa-dev` and `xorg-dev` packages.

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

> Note that Windows applications load from a command prompt by default, which means if you click an icon you may see a command window.
> To fix this add the parameters `-ldflags -H=windowsgui` to your run or build commands.

# Documentation

More documentation is available at the [Fyne developer website](https://fyne.io/develop/) or on [godoc.org](https://godoc.org/fyne.io/fyne).

# Examples

You can find many example applications in the [examples repository](https://github.com/fyne-io/examples/).
Alternatively a list of applications using fyne can be found at [our website](https://fyne.io/develop/applist.html)
