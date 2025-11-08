package layouts

import (
	"net/http"

	"github.com/frenchsoftware/libhtml/attr"
	"github.com/frenchsoftware/libhtml/html"
	"github.com/hyperstitieux/template/database/models"
	"github.com/hyperstitieux/template/views/components"
)

func Base(user *models.User, r *http.Request, title string, children ...html.Node) html.Node {
	// Build main content with children
	mainArgs := []any{attr.Class("flex-1")}
	for _, child := range children {
		mainArgs = append(mainArgs, child)
	}

	return html.Document(
		html.Html(
			attr.Lang("en"),
			html.Head(
				html.Title(html.Text(title)),

				html.Meta(attr.Charset("utf-8")),
				html.Meta(attr.Name("viewport"), attr.Content("width=device-width, initial-scale=1")),

				// Theme initialization (must be in head before body renders)
				components.ThemeInitScript(),

				// Fonts
				html.Link(attr.Rel("preconnect"), attr.Href("https://fonts.googleapis.com")),
				html.Link(attr.Rel("preconnect"), attr.Href("https://fonts.gstatic.com"), html.Attr("crossorigin", "")),
				html.Link(
					attr.Rel("stylesheet"),
					attr.Href("https://fonts.googleapis.com/css2?family=IBM+Plex+Sans:ital,wght@0,100;0,200;0,300;0,400;0,500;0,600;0,700;1,100;1,200;1,300;1,400;1,500;1,600;1,700&family=Workbench&display=swap"),
				),

				// Stylesheets
				html.Link(attr.Rel("stylesheet"), attr.Href("/styles.css")),

				// Lucide Icons
				html.Script(attr.Src("https://unpkg.com/lucide@latest")),

				// Javascript
				html.Script(attr.Src("/js/app.js"), html.Attr("defer", "")),

				// Hot reload script (development only)
				components.HotReloadScript(),
			),
			html.Body(
				attr.Class("font-sans min-h-screen flex flex-col"),

				components.Banner(),
				components.Header(user, r),

				// Main content area - takes remaining space
				html.Main(mainArgs...),

				components.Footer(),

				html.Script(html.Raw("lucide.createIcons();")),
			),
		),
	)
}
