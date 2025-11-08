package pages

import (
	"net/http"

	"github.com/frenchsoftware/libhtml/attr"
	"github.com/frenchsoftware/libhtml/html"
	"github.com/frenchsoftware/libvalidator/validator"
	"github.com/hyperstitieux/template/auth"
	"github.com/hyperstitieux/template/views"
	"github.com/hyperstitieux/template/views/components/ui"
	"github.com/hyperstitieux/template/views/layouts"
)

func Home(w http.ResponseWriter, r *http.Request) error {
	// Get authenticated user from context (nil if not authenticated)
	user := auth.GetCurrentUser(r)

	// Build page
	page := layouts.Base(user, r, "French Software",
		// Page content goes here
	)

	// Render page
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	return page.Render(w)
}

func Settings(w http.ResponseWriter, r *http.Request) error {
	return SettingsWithErrors(w, r, nil)
}

func SettingsWithErrors(w http.ResponseWriter, r *http.Request, errs validator.ValidationErrors) error {
	// Get authenticated user from context (required for settings page)
	user := views.GetUser(r)
	if user == nil {
		// Redirect to OAuth with return URL
		http.Redirect(w, r, "/auth/google?redirect=/settings", http.StatusTemporaryRedirect)
		return nil
	}

	// Build page
	page := layouts.Base(user, r, "Settings - French Software",
		// Main content container
		html.Div(
			attr.Class("max-w-4xl mx-auto px-8 py-8"),

			// Page header
			html.Div(
				attr.Class("mb-8"),
				html.H1(
					attr.Class("text-3xl font-semibold mb-2"),
					html.Text("Settings"),
				),
				html.P(
					attr.Class("text-muted-foreground"),
					html.Text("Manage your account settings and preferences"),
				),
			),

			// Settings cards container
			html.Div(
				attr.Class("flex flex-col gap-6"),

				// General settings form
				html.Form(
					attr.Id("profile-form"),
					attr.Action("/settings/update-profile"),
					attr.Method("POST"),

					ui.Card(
						ui.CardHeader(ui.CardHeaderProps{
							Title:       "General",
							Description: "Update your personal information",
						}),
						ui.CardSection(
							// Email field (readonly)
							html.Div(
								attr.Class("flex flex-col gap-2"),
								html.Label(
									attr.For("email"),
									attr.Class("text-sm font-medium"),
									html.Text("Email"),
								),
								html.Input(
									attr.Type("email"),
									attr.Id("email"),
									attr.Name("email"),
									attr.Value(user.Email),
									attr.Readonly("true"),
									attr.Class("input bg-muted text-muted-foreground cursor-not-allowed"),
								),
								html.P(
									attr.Class("text-xs text-muted-foreground"),
									html.Text("Your email address cannot be changed"),
								),
							),

							// Name field
							html.Div(
								attr.Class("flex flex-col gap-2"),
								html.Label(
									attr.For("name"),
									attr.Class("text-sm font-medium"),
									html.Text("Name"),
								),
								html.Input(
									attr.Type("text"),
									attr.Id("name"),
									attr.Name("name"),
									attr.Value(user.Name),
									attr.Required("true"),
									attr.Maxlength("80"),
									attr.ClassIfElse(errs != nil && errs.Has("name"), "input border-destructive focus:ring-destructive", "input"),
								),
								html.If(errs != nil && errs.Has("name"),
									html.P(
										attr.Class("text-xs text-destructive"),
										html.Text(errs.Get("name")),
									),
								),
							),
						),
						ui.CardFooter(
							html.Button(
								attr.Type("submit"),
								attr.Class("btn-primary"),
								html.Text("Save changes"),
							),
						),
					),
				),

				// Danger zone card
				ui.Card(
					ui.CardHeader(ui.CardHeaderProps{
						Title:       "Danger Zone",
						Description: "Irreversible and destructive actions",
					}),
					ui.CardSection(
						html.Div(
							attr.Class("border border-destructive rounded-lg p-4"),
							html.Div(
								attr.Class("flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4"),
								html.Div(
									html.H3(
										attr.Class("text-sm font-medium"),
										html.Text("Delete Account"),
									),
									html.P(
										attr.Class("text-sm text-muted-foreground"),
										html.Text("Permanently delete your account and all associated data"),
									),
								),
								html.Button(
									attr.Type("button"),
									attr.Class("btn-destructive"),
									html.Attr("onclick", "document.getElementById('delete-account-dialog').showModal()"),
									html.Text("Delete Account"),
								),
							),
						),
					),
				),
			),
		),

		// Delete account confirmation dialog with form in footer
		ui.AlertDialog(ui.AlertDialogProps{
			ID:          "delete-account-dialog",
			Title:       "Are you absolutely sure?",
			Description: "This action cannot be undone. This will permanently delete your account and remove all your data from our servers.",
			Footer: html.Div(
				attr.Class("flex gap-2"),
				html.Button(
					attr.Type("button"),
					attr.Class("btn-outline"),
					html.Attr("onclick", "this.closest('dialog').close()"),
					html.Text("Cancel"),
				),
				html.Form(
					attr.Action("/settings/delete-account"),
					attr.Method("POST"),
					attr.Class("inline"),
					html.Button(
						attr.Type("submit"),
						attr.Class("btn-destructive"),
						html.Text("Delete Account"),
					),
				),
			),
		}),
	)

	// Render page
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	return page.Render(w)
}
