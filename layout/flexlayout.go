package layout

import (
	"log"

	"fyne.io/fyne/v2"
)

// Declare conformity with Layout interface
var _ fyne.Layout = (*flexLayout)(nil)

type flexLayout struct {
	Axis               fyne.Axis
	MainAxisAlignment  fyne.MainAxisAlignment
	CrossAxisAlignment fyne.CrossAxisAlignment
}

// NewFlexLayout creates a new FlexLayout instance.
func NewFlexLayout(axis fyne.Axis, axisAlignment *fyne.AxisAlignment) fyne.Layout {
	return &flexLayout{
		Axis:               axis,
		MainAxisAlignment:  axisAlignment.MainAxisAlignment,
		CrossAxisAlignment: axisAlignment.CrossAxisAlignment,
	}
}

// ===============================================================
// Privates
// ===============================================================

// Kind returns what kind of flex layout is this object based on Axis field.
func (l *flexLayout) kind() string {
	if l.Axis == fyne.AxisVertical {
		return "ColumnFlexLayout"
	}
	return "RowFlexLayout"
}

// getFlex returns the flex factor if any, otherwise it returns 0.
func (l *flexLayout) getFlex(obj fyne.CanvasObject) int {
	type flexibleobj interface{ Flex() int }
	flex := 0
	if f, ok := obj.(flexibleobj); ok {
		flex = f.Flex()
	}
	return flex
}

// getCrossSize returns the cross size based on the layout Axis.
func (l *flexLayout) getCrossSize(obj fyne.CanvasObject, min bool) float32 {
	minSize := obj.MinSize()
	curSize := obj.Size()
	switch l.Axis {
	case fyne.AxisHorizontal:
		if min {
			return minSize.Height
		}
		return fyne.Max(curSize.Height, minSize.Height)
	case fyne.AxisVertical:
		if min {
			return minSize.Width
		}
		return fyne.Max(curSize.Width, minSize.Width)
	}
	if min {
		return minSize.Height
	}
	return fyne.Max(curSize.Height, minSize.Height)
}

// getMainSize returns the main size based on the layout Axis.
func (l *flexLayout) getMainSize(obj fyne.CanvasObject, min bool) float32 {
	minSize := obj.MinSize()
	curSize := obj.Size()
	switch l.Axis {
	case fyne.AxisHorizontal:
		if min {
			return minSize.Width
		}
		return fyne.Max(curSize.Width, minSize.Width)
	case fyne.AxisVertical:
		if min {
			return minSize.Height
		}
		return fyne.Max(curSize.Height, minSize.Height)
	}
	if min {
		return minSize.Width
	}
	return fyne.Max(curSize.Width, minSize.Width)
}

// getDistanceToBaseline returns the distance from the top of the object until its
// text baseline.
func (l *flexLayout) getDistanceToBaseline(obj fyne.CanvasObject) float32 {
	// TODO check
	type baseliner interface{ DistanceToTextBaseline() float32 }
	distance := float32(0)
	if bl, ok := obj.(baseliner); ok {
		distance = bl.DistanceToTextBaseline()
	}
	return distance
}

// ===============================================================
// Fyne.Layout implementation
// ===============================================================

// Layout is called to pack all child objects into a specified size.
//
func (l *flexLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	// Determine used flex factor, size inflexible items, calculate free space.
	totalFlex := 0
	maxMainSize := size.Width
	if l.Axis == fyne.AxisVertical {
		maxMainSize = size.Height
	}
	crossSize := float32(0)
	allocatedSize := float32(0) // Sum of the sizes of the non-flexible children.
	var lastFlexChild fyne.CanvasObject

	// --- Resizing non-flex objects
	for _, obj := range objects {
		if !obj.Visible() {
			continue
		}
		flex := l.getFlex(obj)
		if flex > 0 {
			totalFlex += flex
			lastFlexChild = obj
		} else {
			// TODO mainAxisSize min and max??
			if _, ok := obj.(*fyne.Container); ok {
				switch l.Axis {
				case fyne.AxisHorizontal:
					obj.Resize(fyne.NewSize(l.getMainSize(obj, false), size.Height))
				case fyne.AxisVertical:
					obj.Resize(fyne.NewSize(size.Width, l.getMainSize(obj, false)))
				}
			} else {
				obj.Resize(obj.MinSize())
			}
			allocatedSize += l.getMainSize(obj, false)
			crossSize = fyne.Max(crossSize, l.getCrossSize(obj, false))
		}
	}

	// --- Resizing flex objects
	freeSpace := maxMainSize - allocatedSize
	allocatedFlexSpace := float32(0)
	maxBaselineDistance := float32(0)
	if totalFlex > 0 || l.CrossAxisAlignment == fyne.CrossAxisAlignmentBaseline {
		spacePerFlex := freeSpace / float32(totalFlex)
		maxSizeAboveBaseline := float32(0)
		maxSizeBelowBaseline := float32(0)

		for _, obj := range objects {
			if !obj.Visible() {
				continue
			}
			flex := l.getFlex(obj)
			if flex > 0 {
				maxChildExtent := spacePerFlex * float32(flex)
				if obj == lastFlexChild {
					maxChildExtent = freeSpace - allocatedFlexSpace
				}
				switch l.Axis {
				case fyne.AxisHorizontal:
					obj.Resize(fyne.NewSize(maxChildExtent, l.getCrossSize(obj, false)))
				case fyne.AxisVertical:
					obj.Resize(fyne.NewSize(l.getCrossSize(obj, false), maxChildExtent))
				}
				objChildMainSize := l.getMainSize(obj, false)
				if objChildMainSize > maxChildExtent {
					log.Printf("WARNING: Widget %s can not be shrinked to %.1f (current size: %.1f)\n",
						"WidgetName", maxChildExtent, objChildMainSize)
				}
				allocatedSize += objChildMainSize
				allocatedFlexSpace += maxChildExtent
				crossSize = fyne.Max(crossSize, l.getCrossSize(obj, false))
			}
			if l.CrossAxisAlignment == fyne.CrossAxisAlignmentBaseline {
				distance := l.getDistanceToBaseline(obj)
				maxBaselineDistance = fyne.Max(maxBaselineDistance, distance)
				maxSizeAboveBaseline = fyne.Max(
					distance,
					maxSizeAboveBaseline,
				)
				maxSizeBelowBaseline = fyne.Max(
					fyne.Max(obj.Size().Height, obj.MinSize().Height)-distance,
					maxSizeBelowBaseline,
				)
				crossSize = fyne.Max(maxSizeAboveBaseline+maxSizeBelowBaseline, crossSize)
			}
		}
	}

	// --- MaxAxisAlignment
	actualSizeDelta := maxMainSize - allocatedSize
	if overflow := fyne.Max(0, -actualSizeDelta); overflow > 0 {
		log.Printf(`[WARNING OVERFLOW]:
 %s overflows by %.1f.
 This warning is commonly triggered when you wrap widgets with a Flexible widget and they cannot
 be shrinked beyond to its minimum size as required for the layout calculations.
 `, l.kind(), overflow)
	}
	remainingSpace := fyne.Max(0, actualSizeDelta)
	leadingSpace := float32(0)
	betweenSpace := float32(0)
	switch l.MainAxisAlignment {
	case fyne.MainAxisAlignmentStart:
		leadingSpace = 0
		betweenSpace = 0
	case fyne.MainAxisAlignmentEnd:
		leadingSpace = remainingSpace
		betweenSpace = 0
	case fyne.MainAxisAlignmentCenter:
		leadingSpace = remainingSpace / 2
		betweenSpace = 0
	case fyne.MainAxisAlignmentSpaceBetween:
		leadingSpace = 0
		betweenSpace = 0
		if len(objects) > 1 {
			betweenSpace = remainingSpace / float32(len(objects)-1)
		}
	case fyne.MainAxisAlignmentSpaceAround:
		betweenSpace = 0
		if len(objects) > 0 {
			betweenSpace = remainingSpace / float32(len(objects))
		}
		leadingSpace = betweenSpace / 2
	case fyne.MainAxisAlignmentSpaceEvenly:
		betweenSpace = 0
		if len(objects) > 0 {
			betweenSpace = remainingSpace / float32(len(objects)+1)
		}
		leadingSpace = betweenSpace
	}

	// --- Position objects
	objMainPos := leadingSpace
	for _, obj := range objects {
		if !obj.Visible() {
			continue
		}
		objCrossPos := float32(0)
		switch l.CrossAxisAlignment {
		case fyne.CrossAxisAlignmentStart:
			objCrossPos = 0
		case fyne.CrossAxisAlignmentEnd:
			objCrossPos = crossSize - l.getCrossSize(obj, false)
		case fyne.CrossAxisAlignmentCenter:
			objCrossPos = crossSize/2 - l.getCrossSize(obj, false)/2
		case fyne.CrossAxisAlignmentBaseline:
			objCrossPos = 0
			if l.Axis == fyne.AxisHorizontal {
				distance := l.getDistanceToBaseline(obj)
				objCrossPos = maxBaselineDistance - distance
			}
		}
		switch l.Axis {
		case fyne.AxisHorizontal:
			obj.Move(fyne.NewPos(objMainPos, objCrossPos))
		case fyne.AxisVertical:
			obj.Move(fyne.NewPos(objCrossPos, objMainPos))
		}
		objMainPos += l.getMainSize(obj, false) + betweenSpace
	}
}

// MinSize calculates the smallest size that will fit the listed
// CanvasObjects using this Layout algorithm.
func (l *flexLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	// this is calculated by dividing minSize with flex factor and get the greater one,
	// so we can guarantee minSize for all widgets
	minSpacePerFlex := float32(0)
	totalFlex := 0
	allocatedSize := float32(0)
	crossSize := float32(0)
	maxSizeAboveBaseline := float32(0)
	maxSizeBelowBaseline := float32(0)
	for _, obj := range objects {
		if !obj.Visible() {
			continue
		}
		flex := l.getFlex(obj)
		if flex > 0 {
			totalFlex += flex
			minOneFlexSize := l.getMainSize(obj, true) / float32(flex)
			minSpacePerFlex = fyne.Max(minSpacePerFlex, minOneFlexSize)
		} else {
			allocatedSize += l.getMainSize(obj, true)
		}
		crossSize = fyne.Max(crossSize, l.getCrossSize(obj, true))
		if l.CrossAxisAlignment == fyne.CrossAxisAlignmentBaseline {
			distance := l.getDistanceToBaseline(obj)
			maxSizeAboveBaseline = fyne.Max(
				distance,
				maxSizeAboveBaseline,
			)
			maxSizeBelowBaseline = fyne.Max(
				fyne.Max(obj.Size().Height, obj.MinSize().Height)-distance,
				maxSizeBelowBaseline,
			)
			crossSize = fyne.Max(maxSizeAboveBaseline+maxSizeBelowBaseline, crossSize)
		}
	}

	allocatedFlexSpace := minSpacePerFlex * float32(totalFlex)

	if l.Axis == fyne.AxisHorizontal {
		return fyne.NewSize(allocatedSize+allocatedFlexSpace, crossSize)
	}
	return fyne.NewSize(crossSize, allocatedSize+allocatedFlexSpace)
}
