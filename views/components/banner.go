package components

import (
	"time"

	"github.com/frenchsoftware/libhtml/attr"
	"github.com/frenchsoftware/libhtml/html"
)

func Banner() html.Node {
	return html.Div(
		attr.Class("flex justify-end w-full bg-black text-white dark:bg-secondary border-b border-transparent dark:border-border px-8 py-2"),
		html.Div(
			attr.Class("text-sm"),
			attr.Id("banner-time"),
			html.Text(time.Now().Format("Mon 2 Jan 2006 15:04")),
		),
	)
}
