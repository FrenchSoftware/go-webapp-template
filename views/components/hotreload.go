package components

import (
	"os"

	"github.com/frenchsoftware/libhtml/attr"
	"github.com/frenchsoftware/libhtml/html"
)

// HotReloadScript adds the hot reload script in development mode
func HotReloadScript() html.Node {
	// Only add in development mode
	if os.Getenv("GO_ENV") == "production" {
		return nil
	}

	return html.Script(
		attr.Src("/js/hotreload.js"),
		html.Attr("defer", ""),
	)
}
