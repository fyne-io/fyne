package apps

import "fmt"
import "log"
import "strconv"

import "github.com/fyne-io/fyne"
import "github.com/fyne-io/fyne/layout"
import "github.com/fyne-io/fyne/widget"

import "github.com/Knetic/govaluate"

var equation string
var output *widget.Label
var functions = make(map[string]func())

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
	str := fmt.Sprintf("%d", number)
	action := func() {
		digit(number)
	}
	functions[str] = action
	return widget.NewButton(str, action)
}

func charButton(char string) *widget.Button {
	action := func() {
		character(char)
	}
	functions[char] = action
	return widget.NewButton(char, action)
}

func keyDown(ev *fyne.KeyEvent) {
	if ev.String == "=" || ev.Name == "Return" || ev.Name == "KP_Enter" {
		evaluate()
		return
	} else if ev.Name == "c" {
		clear()
		return
	}

	action := functions[ev.String]
	if action != nil {
		action()
	}
}

// Calculator loads a calculator example window for the specified app context
func Calculator(app fyne.App) {
	output = widget.NewLabel("")
	output.Alignment = fyne.TextAlignTrailing
	equals := widget.NewButton("=", func() {
		evaluate()
	})
	equals.Style = widget.PrimaryButton

	window := app.NewWindow("Calc")
	window.Canvas().SetContent(fyne.NewContainerWithLayout(layout.NewGridLayout(1),
		output,
		fyne.NewContainerWithLayout(layout.NewGridLayout(4),
			charButton("+"),
			charButton("-"),
			charButton("*"),
			charButton("/")),
		fyne.NewContainerWithLayout(layout.NewGridLayout(4),
			digitButton(7),
			digitButton(8),
			digitButton(9),
			widget.NewButton("C", func() {
				clear()
			})),
		fyne.NewContainerWithLayout(layout.NewGridLayout(4),
			digitButton(4),
			digitButton(5),
			digitButton(6),
			charButton("(")),
		fyne.NewContainerWithLayout(layout.NewGridLayout(4),
			digitButton(1),
			digitButton(2),
			digitButton(3),
			charButton(")")),
		fyne.NewContainerWithLayout(layout.NewGridLayout(2),
			fyne.NewContainerWithLayout(layout.NewGridLayout(2),
				digitButton(0),
				charButton(".")),
			equals)),
	)

	window.Canvas().SetOnKeyDown(keyDown)
	window.Show()
}
