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
	equals := widget.NewButton("=", func(*event.MouseEvent) {
		evaluate()
	})
	equals.Style = widget.PrimaryButton

	window := app.NewWindow("Calc")
	window.Canvas().SetContent(ui.NewContainerWithLayout(layout.NewGridLayout(1),
		output,
		ui.NewContainerWithLayout(layout.NewGridLayout(4),
			charButton("+"),
			charButton("-"),
			charButton("*"),
			charButton("/")),
		ui.NewContainerWithLayout(layout.NewGridLayout(4),
			digitButton(7),
			digitButton(8),
			digitButton(9),
			widget.NewButton("C", func(*event.MouseEvent) {
				clear()
			})),
		ui.NewContainerWithLayout(layout.NewGridLayout(4),
			digitButton(4),
			digitButton(5),
			digitButton(6),
			charButton("(")),
		ui.NewContainerWithLayout(layout.NewGridLayout(4),
			digitButton(1),
			digitButton(2),
			digitButton(3),
			charButton(")")),
		ui.NewContainerWithLayout(layout.NewGridLayout(2),
			ui.NewContainerWithLayout(layout.NewGridLayout(2),
				digitButton(0),
				charButton(".")),
			equals)),
	)
	window.Show()
}
