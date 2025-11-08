package ui

import (
	"github.com/frenchsoftware/libhtml/attr"
	"github.com/frenchsoftware/libhtml/html"
)

// AlertDialogProps defines props for AlertDialog component
type AlertDialogProps struct {
	ID            string    // Unique identifier for the dialog
	Title         string    // Dialog title
	Description   string    // Dialog description text
	CancelText    string    // Text for cancel button
	ConfirmText   string    // Text for confirm button
	ConfirmAction string    // JavaScript code to execute on confirm (e.g., "handleDelete()")
	Footer        html.Node // Optional: Custom footer content (overrides default buttons)
}

// AlertDialog creates a Basecoat UI alert dialog (modal that requires user response)
// Note: Alert dialogs don't show close button and don't close on backdrop click
func AlertDialog(props AlertDialogProps) html.Node {
	// Build footer - will contain buttons directly
	var footerArgs []any
	if props.Footer != nil {
		// Use custom footer content
		footerArgs = []any{props.Footer}
	} else {
		// Use default buttons
		footerArgs = []any{
			html.Button(
				attr.Type("button"),
				attr.Class("btn-outline"),
				html.Attr("onclick", "this.closest('dialog').close()"),
				html.Text(props.CancelText),
			),
			html.Button(
				attr.Type("button"),
				attr.Class("btn-destructive"),
				html.Attr("onclick", props.ConfirmAction+"; this.closest('dialog').close()"),
				html.Text(props.ConfirmText),
			),
		}
	}

	return html.Dialog(
		attr.Id(props.ID),
		attr.Class("dialog"),
		html.Attr("aria-labelledby", props.ID+"-title"),
		html.Attr("aria-describedby", props.ID+"-description"),
		html.Div(
			html.Header(
				html.H2(
					attr.Id(props.ID+"-title"),
					html.Text(props.Title),
				),
				html.P(
					attr.Id(props.ID+"-description"),
					html.Text(props.Description),
				),
			),
			html.Footer(footerArgs...),
		),
	)
}

// Dialog creates a standard Basecoat UI dialog (can close on backdrop click)
type DialogProps struct {
	ID          string      // Unique identifier for the dialog
	Title       string      // Dialog title
	Description string      // Dialog description text
	Content     html.Node   // Dialog content (section)
	Footer      html.Node   // Footer content with actions
}

func Dialog(props DialogProps) html.Node {
	return html.Dialog(
		attr.Id(props.ID),
		attr.Class("dialog"),
		html.Attr("aria-labelledby", props.ID+"-title"),
		html.Attr("aria-describedby", props.ID+"-description"),
		html.Attr("onclick", "if (event.target === this) this.close()"),
		html.Div(
			html.Header(
				html.H2(
					attr.Id(props.ID+"-title"),
					html.Text(props.Title),
				),
				html.If(props.Description != "",
					html.P(
						attr.Id(props.ID+"-description"),
						html.Text(props.Description),
					),
				),
			),
			html.If(props.Content != nil,
				html.Section(props.Content),
			),
			html.If(props.Footer != nil,
				html.Footer(props.Footer),
			),
			html.Button(
				attr.Type("button"),
				html.Attr("aria-label", "Close dialog"),
				html.Attr("onclick", "this.closest('dialog').close()"),
				attr.Class("absolute top-4 right-4 text-muted-foreground hover:text-foreground"),
				html.Text("✕"),
			),
		),
	)
}
