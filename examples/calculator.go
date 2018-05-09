package examples

import "fmt"
import "log"
import "strconv"

import "github.com/fyne-io/fyne/app"
import "github.com/fyne-io/fyne/ui"
import "github.com/fyne-io/fyne/ui/event"
import "github.com/fyne-io/fyne/ui/layout"
import "github.com/fyne-io/fyne/ui/widget"

import "github.com/Knetic/govaluate"

var equation string
var output *widget.Label
var container *ui.Container

func display(newtext string) {
	equation = newtext
	output.SetText(newtext)
}

func character(char string) {
	display(equation + char)
}

func digit(d int) {
	character(strconv.Itoa(d))
}

func clear() {
	display("")
}

func evaluate() {
	expression, err := govaluate.NewEvaluableExpression(output.Text)
	if err == nil {
		result, err := expression.Evaluate(nil)
		if err == nil {
			display(strconv.FormatFloat(result.(float64), 'f', -1, 64))
		}
	}

	if err != nil {
		log.Println("Error in calculation", err)
		display("error")
	}

	equation = ""
}

func digitButton(number int) *widget.Button {
	return widget.NewButton(fmt.Sprintf("%d", number), func(e *event.MouseEvent) {
		digit(number)
	})
}

func charButton(char string) *widget.Button {
	return widget.NewButton(char, func(e *event.MouseEvent) {
		character(char)
	})
}

func Calculator(app app.App) {
	output = widget.NewLabel("")
	row1 := ui.NewContainer(
		charButton("+"),
		charButton("-"),
		charButton("*"),
		charButton("/"))
	row1.Layout = layout.NewGridLayout(4)
	row2 := ui.NewContainer(
		digitButton(7),
		digitButton(8),
		digitButton(9),
		widget.NewButton("C", func(*event.MouseEvent) {
			clear()
		}))
	row2.Layout = layout.NewGridLayout(4)
	row3 := ui.NewContainer(
		digitButton(4),
		digitButton(5),
		digitButton(6),
		charButton("("))
	row3.Layout = layout.NewGridLayout(4)
	row4 := ui.NewContainer(
		digitButton(1),
		digitButton(2),
		digitButton(3),
		charButton(")"))
	row4.Layout = layout.NewGridLayout(4)
	row5a := ui.NewContainer(
		digitButton(0),
		charButton("."))
	row5a.Layout = layout.NewGridLayout(2)
	equals := widget.NewButton("=", func(*event.MouseEvent) {
		evaluate()
	})
	equals.Style = widget.PrimaryButton
	row5 := ui.NewContainer(
		row5a,
		equals)
	row5.Layout = layout.NewGridLayout(2)

	window := app.NewWindow("Calc")
	container = ui.NewContainer(output, row1, row2, row3, row4, row5)
	container.Layout = layout.NewGridLayout(1)
	window.Canvas().SetContent(container)
	window.Show()
}
