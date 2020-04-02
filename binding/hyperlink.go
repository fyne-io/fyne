package binding

import (
	"net/url"

	"fyne.io/fyne/widget"
)

func BindHyperlinkText(hyperlink *widget.Hyperlink, data *StringBinding) {
	data.AddStringListener(func(s string) {
		if hyperlink.Text != s {
			hyperlink.SetText(s)
		}
	})
}

func BindHyperlinkURL(hyperlink *widget.Hyperlink, data *URLBinding) {
	data.AddURLListener(func(u *url.URL) {
		if hyperlink.URL != u {
			hyperlink.SetURL(u)
		}
	})
}
