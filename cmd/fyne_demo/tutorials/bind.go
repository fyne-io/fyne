package tutorials

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

func bindingScreen(_ fyne.Window) fyne.CanvasObject {
	f := 0.2
	data := binding.BindFloat(&f)
	label := widget.NewLabelWithData(binding.FloatToStringWithFormat(data, "Float value: %0.2f"))
	entry := widget.NewEntryWithData(binding.FloatToString(data))
	floats := container.NewGridWithColumns(2, label, entry)

	slide := widget.NewSliderWithData(0, 1, data)
	slide.Step = 0.01
	bar := widget.NewProgressBarWithData(data)

	buttons := container.NewGridWithColumns(4,
		widget.NewButton("0%", func() {
			data.Set(0)
		}),
		widget.NewButton("30%", func() {
			data.Set(0.3)
		}),
		widget.NewButton("70%", func() {
			data.Set(0.7)
		}),
		widget.NewButton("100%", func() {
			data.Set(1)
		}))

	boolData := binding.NewBool()
	check := widget.NewCheckWithData("Check me!", boolData)
	checkLabel := widget.NewLabelWithData(binding.BoolToString(boolData))
	checkEntry := widget.NewEntryWithData(binding.BoolToString(boolData))
	checks := container.NewGridWithColumns(3, check, checkLabel, checkEntry)
	item := container.NewVBox(floats, slide, bar, buttons, widget.NewSeparator(), checks, widget.NewSeparator())

	dataList := binding.BindFloatList(&[]float64{0.1, 0.2, 0.3})

	button := widget.NewButton("Append", func() {
		dataList.Append(float64(dataList.Length()+1) / 10)
	})

	list := widget.NewListWithData(dataList,
		func() fyne.CanvasObject {
			return container.NewBorder(nil, nil, nil, widget.NewButton("+", nil),
				widget.NewLabel("item x.y"))
		},
		func(item binding.DataItem, obj fyne.CanvasObject) {
			f := item.(binding.Float)
			text := obj.(*fyne.Container).Objects[0].(*widget.Label)
			text.Bind(binding.FloatToStringWithFormat(f, "item %0.1f"))

			btn := obj.(*fyne.Container).Objects[1].(*widget.Button)
			btn.OnTapped = func() {
				val, _ := f.Get()
				_ = f.Set(val + 1)
			}
		})

	formStruct := struct {
		Name, Email string
		Subscribe   bool
	}{}

	formData := binding.BindStruct(&formStruct)
	form := newFormWithData(formData)
	form.OnSubmit = func() {
		fmt.Println("Struct:\n", formStruct)
	}

	listPanel := container.NewBorder(nil, button, nil, nil, list)
	return container.NewBorder(item, nil, nil, nil, container.NewGridWithColumns(2, listPanel, form))
}

func newFormWithData(data binding.DataMap) *widget.Form {
	keys := data.Keys()
	items := make([]*widget.FormItem, len(keys))
	for i, k := range keys {
		data, err := data.GetItem(k)
		if err != nil {
			items[i] = widget.NewFormItem(k, widget.NewLabel(err.Error()))
		}
		items[i] = widget.NewFormItem(k, createBoundItem(data))
	}

	return widget.NewForm(items...)
}

func createBoundItem(v binding.DataItem) fyne.CanvasObject {
	switch val := v.(type) {
	case binding.Bool:
		return widget.NewCheckWithData("", val)
	case binding.Float:
		s := widget.NewSliderWithData(0, 1, val)
		s.Step = 0.01
		return s
	case binding.Int:
		return widget.NewEntryWithData(binding.IntToString(val))
	case binding.String:
		return widget.NewEntryWithData(val)
	default:
		return widget.NewLabel("")
	}
}
