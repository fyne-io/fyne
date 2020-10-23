package tutorials

import (
	"fyne.io/fyne"
)

// Tutorial defines the data structure for a tutorial
type Tutorial struct {
	Title, Intro string
	View         func(w fyne.Window) fyne.CanvasObject
}

var (
	// Tutorials defines the metadata for each tutorial
	Tutorials = map[string]Tutorial{
		"welcome": {"Welcome",
			"Welcome to the Fyne toolkit tutorials.\nChoose an area from the options on the left.",
			welcomeScreen,
		},
		"graphics": {"Graphics",
			"See the canvas and graphics capabilities",
			GraphicsScreen,
		},
		"widgets": {"Widgets",
			"Browse the toolkit widgets",
			WidgetScreen,
		},
		"containers": {"Containers",
			"Containers and collections",
			ContainerScreen,
		},
		"windows": {"Windows",
			"Window and dialog demos",
			DialogScreen,
		},
		"advanced": {"Advanced",
			"Debug and advanced information",
			AdvancedScreen,
		},
	}

	// TutorialIndex  defines how our tutorials should be laid out in the index tree
	TutorialIndex = map[string][]string{
		"": {"welcome", "graphics", "widgets", "containers", "windows", "advanced"},
	}
)
