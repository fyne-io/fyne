package widget

import "fyne.io/fyne"

type menuItemBase struct {
	Child *Menu
}

func (b *menuItemBase) activateChild(parent *menuBase, updateChildPos func()) {
	if b.Child != nil && b.Child.activeChild != nil {
		b.Child.activeChild.Hide()
		b.Child.activeChild = nil
	}

	if parent.activeChild == b.Child {
		return
	}

	if parent.activeChild != nil {
		parent.activeChild.Hide()
	}
	parent.activeChild = b.Child
	if b.Child != nil {
		if b.Child.Size().IsZero() {
			b.Child.Resize(b.Child.MinSize())
			updateChildPos()
		}
		b.Child.Show()
	}
}

func (b *menuItemBase) initChildWidget(menu *fyne.Menu, dismissAction func()) {
	if b.Child != nil {
		return
	}
	b.Child = NewMenu(menu)
	b.Child.Hide()
	b.Child.DismissAction = dismissAction
}
