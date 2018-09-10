package apps

import "fmt"
import "log"
import "strconv"

import "github.com/fyne-io/fyne"
import "github.com/fyne-io/fyne/layout"
import "github.com/fyne-io/fyne/widget"

import "github.com/Knetic/govaluate"

type calc struct {
	equation  string
	output    *widget.Label
	functions map[string]func()
}

func (c *calc) display(newtext string) {
	c.equation = newtext
	c.output.SetText(newtext)
}

func (c *calc) character(char string) {
	c.display(c.equation + char)
}

func (c *calc) digit(d int) {
	c.character(strconv.Itoa(d))
}

func (c *calc) clear() {
	c.display("")
}

func (c *calc) evaluate() {
	expression, err := govaluate.NewEvaluableExpression(c.output.Text)
	if err == nil {
		result, err := expression.Evaluate(nil)
		if err == nil {
			c.display(strconv.FormatFloat(result.(float64), 'f', -1, 64))
		}
	}

	if err != nil {
		log.Println("Error in calculation", err)
		c.display("error")
	}

	c.equation = ""
}

func (c *calc) digitButton(number int) *widget.Button {
	str := fmt.Sprintf("%d", number)
	action := func() {
		c.digit(number)
	}
	c.functions[str] = action
	return widget.NewButton(str, action)
}

func (c *calc) charButton(char string) *widget.Button {
	action := func() {
		c.character(char)
	}
	c.functions[char] = action
	return widget.NewButton(char, action)
}

func (c *calc) keyDown(ev *fyne.KeyEvent) {
	if ev.String == "=" || ev.Name == "Return" || ev.Name == "KP_Enter" {
		c.evaluate()
		return
	} else if ev.Name == "c" {
		c.clear()
		return
	}

	action := c.functions[ev.String]
	if action != nil {
		action()
	}
}

// Calculator loads a calculator example window for the specified app context
func Calculator(app fyne.App) {
	calc := &calc{}
	calc.functions = make(map[string]func())

	calc.output = widget.NewLabel("")
	calc.output.Alignment = fyne.TextAlignTrailing
	calc.output.TextStyle.Monospace = true
	equals := widget.NewButton("=", func() {
		calc.evaluate()
	})
	equals.Style = widget.PrimaryButton

	window := app.NewWindow("Calc")
	window.SetContent(fyne.NewContainerWithLayout(layout.NewGridLayout(1),
		calc.output,
		fyne.NewContainerWithLayout(layout.NewGridLayout(4),
			calc.charButton("+"),
			calc.charButton("-"),
			calc.charButton("*"),
			calc.charButton("/")),
		fyne.NewContainerWithLayout(layout.NewGridLayout(4),
			calc.digitButton(7),
			calc.digitButton(8),
			calc.digitButton(9),
			widget.NewButton("C", func() {
				calc.clear()
			})),
		fyne.NewContainerWithLayout(layout.NewGridLayout(4),
			calc.digitButton(4),
			calc.digitButton(5),
			calc.digitButton(6),
			calc.charButton("(")),
		fyne.NewContainerWithLayout(layout.NewGridLayout(4),
			calc.digitButton(1),
			calc.digitButton(2),
			calc.digitButton(3),
			calc.charButton(")")),
		fyne.NewContainerWithLayout(layout.NewGridLayout(2),
			fyne.NewContainerWithLayout(layout.NewGridLayout(2),
				calc.digitButton(0),
				calc.charButton(".")),
			equals)),
	)

	window.Canvas().SetOnKeyDown(calc.keyDown)
	window.Show()
}
