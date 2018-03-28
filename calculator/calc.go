package main

import "log"
import "strconv"

import "github.com/fyne-io/fyne/ui"
import "github.com/fyne-io/fyne/ui/event"
import "github.com/fyne-io/fyne/ui/layout"
import "github.com/fyne-io/fyne/ui/widget"
import "github.com/fyne-io/fyne-app"

import "github.com/Knetic/govaluate"

var equation string
var output *widget.Label
var container *ui.Container
var window ui.Window

func display() {
	output.SetText(equation)
	window.Canvas().Refresh(output)
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

	output = widget.NewLabel("")
	row1 := ui.NewContainer(
		widget.NewButton("+", func(*event.MouseEvent) {
			character("+")
		}),
		widget.NewButton("-", func(*event.MouseEvent) {
			character("-")
		}),
		widget.NewButton("*", func(*event.MouseEvent) {
			character("*")
		}),
		widget.NewButton("/", func(*event.MouseEvent) {
			character("/")
		}))
	row1.Layout = layout.NewGridLayout(4)
	row2 := ui.NewContainer(
		widget.NewButton("7", func(*event.MouseEvent) {
			digit(7)
		}),
		widget.NewButton("8", func(*event.MouseEvent) {
			digit(8)
		}),
		widget.NewButton("9", func(*event.MouseEvent) {
			digit(9)
		}),
		widget.NewButton("C", func(*event.MouseEvent) {
			clear()
		}))
	row2.Layout = layout.NewGridLayout(4)
	row3 := ui.NewContainer(
		widget.NewButton("4", func(*event.MouseEvent) {
			digit(4)
		}),
		widget.NewButton("5", func(*event.MouseEvent) {
			digit(5)
		}),
		widget.NewButton("6", func(*event.MouseEvent) {
			digit(6)
		}),
		widget.NewButton("(", func(*event.MouseEvent) {
			character("(")
		}))
	row3.Layout = layout.NewGridLayout(4)
	row4 := ui.NewContainer(
		widget.NewButton("1", func(*event.MouseEvent) {
			digit(1)
		}),
		widget.NewButton("2", func(*event.MouseEvent) {
			digit(2)
		}),
		widget.NewButton("3", func(*event.MouseEvent) {
			digit(3)
		}),
		widget.NewButton(")", func(*event.MouseEvent) {
			character(")")
		}))
	row4.Layout = layout.NewGridLayout(4)
	row5a := ui.NewContainer(
		widget.NewButton("0", func(*event.MouseEvent) {
			digit(0)
		}),
		widget.NewButton(".", func(*event.MouseEvent) {
			character(".")
		}))
	row5a.Layout = layout.NewGridLayout(2)
	equals := widget.NewButton("=", func(*event.MouseEvent) {
		evaluate()
	})
	equals.Style = widget.PrimaryButton
	row5 := ui.NewContainer(
		row5a,
		equals)
	row5.Layout = layout.NewGridLayout(2)

	window = app.NewWindow("Calc")
	container = ui.NewContainer(output, row1, row2, row3, row4, row5)
	container.Layout = layout.NewGridLayout(1)
	window.Canvas().SetContent(container)
	window.Show()
}
