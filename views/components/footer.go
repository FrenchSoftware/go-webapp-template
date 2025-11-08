package components

import (
	"github.com/frenchsoftware/libhtml/attr"
	"github.com/frenchsoftware/libhtml/html"
)

func Footer() html.Node {
	return html.Footer(
		attr.Class("w-full bg-secondary border-t border-border"),
		html.Div(
			attr.Class("flex flex-col md:flex-row justify-between md:items-center py-6 px-8 sm:py-2 gap-4"),
			html.Div(
				attr.Class("flex items-center gap-4 text-sm text-muted-foreground dark:text-white"),
				html.Span(
					html.Text("© Copyright 2025 — French Software"),
				),
			),
		),
	)
}
