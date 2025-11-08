package ui

import (
	"github.com/frenchsoftware/libhtml/attr"
	"github.com/frenchsoftware/libhtml/html"
)

// Card creates a Basecoat UI card component
func Card(children ...html.Node) html.Node {
	args := []any{attr.Class("card")}
	for _, child := range children {
		args = append(args, child)
	}
	return html.Div(args...)
}

// CardHeaderProps defines props for CardHeader component
type CardHeaderProps struct {
	Title       string
	Description string
}

// CardHeader creates a card header with title and optional description
func CardHeader(props CardHeaderProps) html.Node {
	return html.Header(
		html.H2(html.Text(props.Title)),
		html.If(props.Description != "",
			html.P(html.Text(props.Description)),
		),
	)
}

// CardSection creates a card section
func CardSection(children ...html.Node) html.Node {
	args := []any{attr.Class("flex flex-col gap-4")}
	for _, child := range children {
		args = append(args, child)
	}
	return html.Section(args...)
}

// CardFooter creates a card footer
func CardFooter(children ...html.Node) html.Node {
	args := []any{attr.Class("flex items-center gap-2")}
	for _, child := range children {
		args = append(args, child)
	}
	return html.Footer(args...)
}
