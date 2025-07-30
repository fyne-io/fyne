//go:build darwin

package widget

/*
int getScrollerPagingBehavior();
*/
import "C"

func isScrollerPageOnTap() bool {
	if C.getScrollerPagingBehavior() == 0 {
		return true
	}
	return false
}
