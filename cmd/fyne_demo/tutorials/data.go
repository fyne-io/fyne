//go:generate go run gen.go

package tutorials

import (
	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/styles"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// Tutorial defines the data structure for a tutorial
type Tutorial struct {
	Title, Intro string
	View         func(w fyne.Window) fyne.CanvasObject
}

type Tutorial2 struct {
	title string

	content []string
	code    []func() fyne.CanvasObject
}

var (
	// Tutorials defines the metadata for each tutorial
	Tutorials = map[string]Tutorial{
		"welcome": {"Welcome", "", welcomeScreen},
		"canvas": {"Canvas",
			"See the canvas capabilities.",
			canvasScreen,
		},
		"animations": {"Animations",
			"See how to animate components.",
			makeAnimationScreen,
		},
		"icons": {"Theme Icons",
			"Browse the embedded icons.",
			iconScreen,
		},
		"containers": {"Containers",
			"Containers group other widgets and canvas objects, organising according to their layout.\n" +
				"Standard containers are illustrated in this section, but developers can also provide custom " +
				"layouts using the fyne.NewContainerWithLayout() constructor.",
			containerScreen,
		},
		"apptabs": {"AppTabs",
			"A container to help divide up an application into functional areas.",
			makeAppTabsTab,
		},
		"border": {"Border",
			"A container that positions items around a central content.",
			makeBorderLayout,
		},
		"box": {"Box",
			"A container arranges items in horizontal or vertical list.",
			makeBoxLayout,
		},
		"center": {"Center",
			"A container to that centers child elements.",
			makeCenterLayout,
		},
		"doctabs": {"DocTabs",
			"A container to display a single document from a set of many.",
			makeDocTabsTab,
		},
		"grid": {"Grid",
			"A container that arranges all items in a grid.",
			makeGridLayout,
		},
		"split": {"Split",
			"A split container divides the container in two pieces that the user can resize.",
			makeSplitTab,
		},
		"scroll": {"Scroll",
			"A container that provides scrolling for its content.",
			makeScrollTab,
		},
		"innerwindow": {"InnerWindow",
			"A window that can be used inside a traditional window to contain a document or content.",
			makeInnerWindowTab,
		},
		"widgets": {"Widgets",
			"In this section you can see the features available in the toolkit widget set.\n" +
				"Expand the tree on the left to browse the individual tutorial elements.",
			widgetScreen,
		},
		"accordion": {"Accordion",
			"Expand or collapse content panels.",
			makeAccordionTab,
		},
		"activity": {"Activity",
			"A spinner indicating activity used in buttons etc.",
			makeActivityTab,
		},
		"button": {"Button",
			"",
			func(fyne.Window) fyne.CanvasObject { return makeNewTutorial("widgets/button.md") },
		},
		"card": {"Card",
			"Group content and widgets.",
			makeCardTab,
		},
		"entry": {"Entry",
			"Different ways to use the entry widget.",
			makeEntryTab,
		},
		"form": {"Form",
			"Gathering input widgets for data submission.",
			makeFormTab,
		},
		"input": {"Input",
			"A collection of widgets for user input.",
			makeInputTab,
		},
		"text": {"Text",
			"Text handling widgets.",
			makeTextTab,
		},
		"toolbar": {"Toolbar",
			"A row of shortcut icons for common tasks.",
			makeToolbarTab,
		},
		"progress": {"Progress",
			"Show duration or the need to wait for a task.",
			makeProgressTab,
		},
		"collections": {"Collections",
			"Collection widgets provide an efficient way to present lots of content.\n" +
				"The List, Table, and Tree provide a cache and re-use mechanism that make it possible to scroll through thousands of elements.\n" +
				"Use this for large data sets or for collections that can expand as users scroll.",
			collectionScreen,
		},
		"list": {"List",
			"A vertical arrangement of cached elements with the same styling.",
			makeListTab,
		},
		"table": {"Table",
			"A two dimensional cached collection of cells.",
			makeTableTab,
		},
		"tree": {"Tree",
			"A tree based arrangement of cached elements with the same styling.",
			makeTreeTab,
		},
		"gridwrap": {"GridWrap",
			"A grid based arrangement of cached elements that wraps rows to fit.",
			makeGridWrapTab,
		},
		"dialogs": {"Dialogs",
			"Work with dialogs.",
			dialogScreen,
		},
		"windows": {"Windows",
			"Window function demo.",
			windowScreen,
		},
		"binding": {"Data Binding",
			"Connecting widgets to a data source.",
			bindingScreen,
		},
		"advanced": {"Advanced",
			"Debug and advanced information.",
			advancedScreen,
		},
	}

	// TutorialIndex  defines how our tutorials should be laid out in the index tree
	TutorialIndex = map[string][]string{
		"":            {"welcome", "canvas", "animations", "icons", "widgets", "collections", "containers", "dialogs", "windows", "binding", "advanced"},
		"collections": {"list", "table", "tree", "gridwrap"},
		"containers":  {"apptabs", "border", "box", "center", "doctabs", "grid", "scroll", "split", "innerwindow"},
		"widgets":     {"accordion", "activity", "button", "card", "entry", "form", "input", "progress", "text", "toolbar"},
	}
)

func makeNewTutorial(file string) fyne.CanvasObject {
	tutorial := tutorials[file]
	top := container.NewVBox(
		widget.NewLabel(tutorial.title), widget.NewSeparator())
	details := container.NewVBox()

	for i, p := range tutorial.content {
		if p == "" {
			continue
		}
		if i%2 == 0 {
			text := widget.NewRichTextFromMarkdown(p)
			text.Wrapping = fyne.TextWrapWord

			details.Add(text)
			continue
		}

		usage := widget.NewTextGridFromString(p)
		usage.ShowLineNumbers = true
		highlightTextGrid(usage)

		codeID := (i - 1) / 2
		preview := tutorial.code[codeID]()

		tools := widget.NewToolbar(
			widget.NewToolbarAction(theme.ContentCopyIcon(), func() {
				fyne.CurrentApp().Clipboard().SetContent(usage.Text())
			}),
			widget.NewToolbarAction(theme.MediaPlayIcon(), func() {
				w := fyne.CurrentApp().NewWindow(tutorial.title + " preview")
				w.SetContent(tutorial.code[codeID]())
				w.Show()
			}),
		)

		style := styles.Get("solarized-dark")
		bg := styleBackgroundColor(chroma.Background, style)

		details.Add(container.NewPadded(container.NewPadded(
			container.NewStack(
				canvas.NewRectangle(bg),
				container.NewPadded(usage,
					container.NewHBox(layout.NewSpacer(), container.NewVBox(
						tools)))))))
		details.Add(widget.NewRichTextFromMarkdown("*Preview:*"))
		details.Add(container.NewPadded(container.NewHBox(preview)))
	}

	return container.NewBorder(top, nil, nil, nil, container.NewScroll(details))
}
