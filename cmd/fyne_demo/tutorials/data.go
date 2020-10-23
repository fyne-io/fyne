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
		"canvas": {"Canvas",
			"See the canvas capabilities",
			canvasScreen,
		},
		"icons": {"Theme Icons",
			"Browse the embedded icons",
			iconScreen,
		},
		"containers": {"Containers",
			"Containers and collections",
			containerScreen,
		},
		"widgets": {"Widgets",
			"Browse the toolkit widgets",
			widgetScreen,
		},
		"collections": {"Collection Widgets",
			"Learn about the collection widgets",
			collectionScreen,
		},
		"dialogs": {"Dialogs",
			"Work with dialogs",
			dialogScreen,
		},
		"windows": {"Windows",
			"Window function demo",
			windowScreen,
		},
		"advanced": {"Advanced",
			"Debug and advanced information",
			advancedScreen,
		},
	}

	// TutorialIndex  defines how our tutorials should be laid out in the index tree
	TutorialIndex = map[string][]string{
		"": {"welcome", "canvas", "icons", "widgets", "collections", "containers", "dialogs", "windows", "advanced"},
	}
)
