package main

import "log"
import "strconv"

import "github.com/fyne-io/fyne/ui"
import "github.com/fyne-io/fyne/ui/layout"
import "github.com/fyne-io/fyne/ui/widget"
import "github.com/fyne-io/fyne-app"

import "github.com/Knetic/govaluate"

var equation string
var output *ui.TextObject
var container *ui.Container
var window ui.Window

func display() {
	output.Text = equation
	window.Canvas().SetContent(container)
}

func character(char string) {
	equation = equation + char
	display()
}

func digit(d int) {
	character(strconv.Itoa(d))
}

func clear() {
	equation = ""
	display()
}

func evaluate() {
	expression, err := govaluate.NewEvaluableExpression(output.Text)
	if err == nil {
		result, err := expression.Evaluate(nil)
		if err == nil {
			equation = strconv.FormatFloat(result.(float64), 'f', -1, 64)
		}
	}

	if err != nil {
		log.Println("Error in calculation", err)
		equation = ""
	}

	display()
	equation = ""
}

func main() {
	app := fyneapp.NewApp()

	output = ui.NewText("")
	row1 := ui.NewContainer(
		widget.NewButton("+", func() {
			character("+")
		}),
		widget.NewButton("-", func() {
			character("-")
		}),
		widget.NewButton("*", func() {
			character("*")
		}),
		widget.NewButton("/", func() {
			character("/")
		}))
	row1.Layout = layout.NewGridLayout(4)
	row2 := ui.NewContainer(
		widget.NewButton("7", func() {
			digit(7)
		}),
		widget.NewButton("8", func() {
			digit(8)
		}),
		widget.NewButton("9", func() {
			digit(9)
		}),
		widget.NewButton("C", func() {
			clear()
		}))
	row2.Layout = layout.NewGridLayout(4)
	row3 := ui.NewContainer(
		widget.NewButton("4", func() {
			digit(4)
		}),
		widget.NewButton("5", func() {
			digit(5)
		}),
		widget.NewButton("6", func() {
			digit(6)
		}),
		widget.NewButton("(", func() {
			character("(")
		}))
	row3.Layout = layout.NewGridLayout(4)
	row4 := ui.NewContainer(
		widget.NewButton("1", func() {
			digit(1)
		}),
		widget.NewButton("2", func() {
			digit(2)
		}),
		widget.NewButton("3", func() {
			digit(3)
		}),
		widget.NewButton(")", func() {
			character(")")
		}))
	row4.Layout = layout.NewGridLayout(4)
	row5a := ui.NewContainer(
		widget.NewButton("0", func() {
			digit(0)
		}),
		widget.NewButton(".", func() {
			character(".")
		}))
	row5a.Layout = layout.NewGridLayout(2)
	row5 := ui.NewContainer(
		row5a,
		widget.NewButton("=", func() {
			evaluate()
		}))
	row5.Layout = layout.NewGridLayout(2)

	window = app.NewWindow("Calc")
	container = ui.NewContainer(output, row1, row2, row3, row4, row5)
	container.Layout = layout.NewGridLayout(1)
	window.Canvas().SetContent(container)
	window.Show()
}
