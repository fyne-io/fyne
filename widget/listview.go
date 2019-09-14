package widget

import "fyne.io/fyne"

type (
	// ListView ListView
	ListView struct {
		ScrollContainer
		box        *Box
		vhs        []*ViewHolder
		onCreate   func(*ViewHolder) fyne.CanvasObject
		onBind     func(*ViewHolder, int)
		onGetCount func() int
		count      int
	}
)

func newListView(onCreate func(*ViewHolder) fyne.CanvasObject, onBind func(*ViewHolder, int), onGetCount func() int, box *Box) *ListView {
	lv := &ListView{
		box:        box,
		onCreate:   onCreate,
		onBind:     onBind,
		onGetCount: onGetCount,
	}
	lv.ScrollContainer = *NewScrollContainer(box)
	lv.count = onGetCount()

	for i := 0; i < lv.count; i++ {
		vh := newViewHolder()
		v := onCreate(vh)
		vh.root = v

		lv.vhs = append(lv.vhs, vh)
		lv.box.Append(vh.root)
	}
	lv.execBindData()
	return lv
}

// NewVListView NewVListView
func NewVListView(onCreate func(*ViewHolder) fyne.CanvasObject, onBind func(*ViewHolder, int), onGetCount func() int) *ListView {
	return newListView(onCreate, onBind, onGetCount, NewVBox())
}

// NewHListView NewHListView
func NewHListView(onCreate func(*ViewHolder) fyne.CanvasObject, onBind func(*ViewHolder, int), onGetCount func() int) *ListView {
	return newListView(onCreate, onBind, onGetCount, NewHBox())
}

func (l *ListView) execBindData() {
	for i := 0; i < l.onGetCount(); i++ {
		l.onBind(l.vhs[i], i)
	}
}

// NotifyDataChange NotifyDataChange
func (l *ListView) NotifyDataChange() {
	originSize := l.count
	newSize := l.onGetCount()
	if newSize > originSize {
		for i := originSize; i < newSize; i++ {

			if i >= len(l.vhs) {
				l.vhs = append(l.vhs, newViewHolder())
			}

			v := l.onCreate(l.vhs[i])
			l.vhs[i].root = v
			l.box.Append(l.vhs[i].root)
		}
	} else {
		for i := newSize; i < originSize; i++ {
			l.vhs[i].root.Hide()
		}
	}

	for i := 0; i < newSize; i++ {
		l.onBind(l.vhs[i], i)
	}
	l.count = newSize
}
