package main

import (
	"fyne.io/fyne"
	"fyne.io/fyne/cmd/fyne_demo/screens"
)

type tutorial struct {
	title string // TODO add code intro text and snippets later
	view  func(w fyne.Window) fyne.CanvasObject
}

var (
	tutorials = map[string]tutorial{
		"welcome": {"Welcome",
			welcomeScreen,
		},
		"graphics": {"Graphics",
			screens.GraphicsScreen,
		},
		"widgets": {"Widgets",
			screens.WidgetScreen,
		},
		"containers": {"Containers",
			screens.ContainerScreen,
		},
		"windows": {"Windows",
			screens.DialogScreen,
		},
		"advanced": {"Advanced",
			screens.AdvancedScreen,
		},
	}

	tutorialTree = map[string][]string{
		"": {"welcome", "graphics", "widgets", "containers", "windows", "advanced"},
	}
)
