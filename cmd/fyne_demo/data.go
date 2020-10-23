package main

import (
	"fyne.io/fyne"
	"fyne.io/fyne/cmd/fyne_demo/screens"
)

type tutorial struct {
	title, intro string
	view         func(w fyne.Window) fyne.CanvasObject
}

var (
	tutorials = map[string]tutorial{
		"welcome": {"Welcome",
			"Welcome to the Fyne toolkit tutorials.\nChoose an area from the options on the left.",
			welcomeScreen,
		},
		"graphics": {"Graphics",
			"See the canvas and graphics capabilities",
			screens.GraphicsScreen,
		},
		"widgets": {"Widgets",
			"Browse the toolkit widgets",
			screens.WidgetScreen,
		},
		"containers": {"Containers",
			"Containers and collections",
			screens.ContainerScreen,
		},
		"windows": {"Windows",
			"Window and dialog demos",
			screens.DialogScreen,
		},
		"advanced": {"Advanced",
			"Debug and advanced information",
			screens.AdvancedScreen,
		},
	}

	tutorialTree = map[string][]string{
		"": {"welcome", "graphics", "widgets", "containers", "windows", "advanced"},
	}
)
