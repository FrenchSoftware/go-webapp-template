package components

import (
	"net/http"

	"github.com/frenchsoftware/libhtml/attr"
	"github.com/frenchsoftware/libhtml/html"
	"github.com/hyperstitieux/template/database/models"
	"github.com/hyperstitieux/template/views/components/icons"
)

func Header(user *models.User, r *http.Request) html.Node {
	currentPath := r.URL.Path
	// Determine right section content based on authentication
	var rightSection html.Node

	if user != nil {
		// User is authenticated - show dropdown menu with user info
		rightSection = html.Div(
			attr.Class("dropdown-menu"),
			// Trigger button
			html.Button(
				html.Attr("aria-haspopup", "menu"),
				html.Attr("aria-controls", "user-menu"),
				html.Attr("aria-expanded", "false"),
				attr.Class("flex items-center justify-center size-9 rounded-lg transition-colors cursor-pointer overflow-hidden"),
				html.If(user.Picture != nil && *user.Picture != "",
					html.Img(
						attr.Src(*user.Picture),
						attr.Alt(user.Name),
						attr.Class("w-full h-full object-cover"),
						attr.Referrerpolicy("no-referrer"),
					),
				),
			),
			// Dropdown popover
			html.Div(
				html.Attr("data-popover", ""),
				html.Attr("data-side", "bottom"),
				html.Attr("data-align", "end"),
				html.Attr("aria-hidden", "true"),
				attr.Class("w-56"),
				html.Div(
					html.Attr("role", "menu"),
					attr.Id("user-menu"),
					// Account name header
					html.Div(
						html.Attr("role", "heading"),
						attr.Class("px-2 py-1.5 max-w-56 flex flex-col"),
						html.Div(
							attr.Class("text-sm font-medium break-words"),
							html.Text(user.Name),
						),
						html.Div(
							attr.Class("text-xs text-muted-foreground break-words"),
							html.Text(user.Email),
						),
					),
					// Separator
					html.Hr(html.Attr("role", "separator")),
					// Settings menu item
					html.A(
						html.Attr("role", "menuitem"),
						attr.Href("/settings"),
						attr.Class("flex cursor-pointer items-center gap-2"),
						html.I(html.Attr("data-lucide", "settings")),
						html.Text("Settings"),
					),
					// Separator
					html.Hr(html.Attr("role", "separator")),
					// Logout menu item
					html.A(
						html.Attr("role", "menuitem"),
						attr.Href("/auth/sign-out"),
						attr.Class("flex cursor-pointer items-center gap-2 text-destructive"),
						html.I(html.Attr("data-lucide", "log-out"), attr.Class("text-destructive")),
						html.Text("Log out"),
					),
				),
			),
		)
	} else {
		// User is not authenticated - show sign in button
		rightSection = html.Div(
			attr.Class("flex items-center"),
			html.A(
				attr.Class("btn-primary h-9 flex items-center"),
				attr.Href("/auth/google"),
				icons.Google(),
				html.Span(
					attr.Class("font-medium"),
					html.Text("Sign in with Google"),
				),
			),
		)
	}

	return html.Header(
		attr.Class("flex flex-wrap justify-between items-center px-8 py-4 gap-4 border-b border-border"),

		// Left section: Logo + Navigation
		html.Div(
			attr.Class("flex flex-wrap items-center gap-4"),
			// Logo
			html.A(
				attr.Href("/"),
				html.Attr("data-tooltip", "Go back home"),
				html.Attr("data-side", "bottom"),
				html.H1(
					attr.Class("text-2xl text-black dark:text-white cursor-pointer hover:opacity-80 transition-opacity antialiased"),
					attr.Style("font-family: var(--font-workbench); font-weight: 400;"),
					html.Text("French Software"),
				),
			),
			// Navigation links
			html.Div(
				attr.Class("flex flex-wrap items-center gap-2"),
				html.A(
					attr.Href("/"),
					attr.ClassIfElse(currentPath == "/", "btn-ghost bg-accent", "btn-ghost"),
					html.Text("Home"),
				),
			),
		),

		// Right section: Theme switcher + Dynamic based on auth status
		html.Div(
			attr.Class("flex items-center gap-2"),
			ThemeSwitcher(),
			rightSection,
		),
	)
}
