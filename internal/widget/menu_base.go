package widget

type menuBase struct {
	DismissAction func()

	activeChild *Menu
}

func (b *menuBase) dismiss() {
	if b.activeChild != nil {
		defer b.activeChild.dismiss()
		b.activeChild.Hide()
		b.activeChild = nil
	}
	if b.DismissAction != nil {
		b.DismissAction()
	}
}
