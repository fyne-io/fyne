//go:build !darwin

package widget

func getScrollerPagingBehavior() scrollBarTapBehavior {
	return scrollBarTapBehaviorScrollToPosition
}
