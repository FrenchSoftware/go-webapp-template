package router

import (
	"io"
	"net/http"

	"github.com/gorilla/mux"
)

// Router wraps mux.Router with convenient helper methods
type Router struct {
	*mux.Router
	config       Config
	hotReloadMux *mux.Router
}

// WrapRouter wraps a mux.Router to add helper methods
func WrapRouter(r *mux.Router) *Router {
	return &Router{Router: r}
}

// Get registers a GET route with automatic error handling
func (r *Router) Get(path string, handler HandlerFunc) *mux.Route {
	return r.HandleFunc(path, Handle(handler)).Methods(http.MethodGet)
}

// Post registers a POST route with automatic error handling
func (r *Router) Post(path string, handler HandlerFunc) *mux.Route {
	return r.HandleFunc(path, Handle(handler)).Methods(http.MethodPost)
}

// Put registers a PUT route with automatic error handling
func (r *Router) Put(path string, handler HandlerFunc) *mux.Route {
	return r.HandleFunc(path, Handle(handler)).Methods(http.MethodPut)
}

// Patch registers a PATCH route with automatic error handling
func (r *Router) Patch(path string, handler HandlerFunc) *mux.Route {
	return r.HandleFunc(path, Handle(handler)).Methods(http.MethodPatch)
}

// Delete registers a DELETE route with automatic error handling
func (r *Router) Delete(path string, handler HandlerFunc) *mux.Route {
	return r.HandleFunc(path, Handle(handler)).Methods(http.MethodDelete)
}

// Options registers an OPTIONS route with automatic error handling
func (r *Router) Options(path string, handler HandlerFunc) *mux.Route {
	return r.HandleFunc(path, Handle(handler)).Methods(http.MethodOptions)
}

// HTML is a helper for rendering HTML pages with automatic Content-Type header
// Accepts functions that write to io.Writer (like libhtml's Render/Output methods)
func HTML(renderer func(w io.Writer) error) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		return renderer(w)
	}
}

// GetHTML registers a GET route for HTML pages with automatic Content-Type header
// Perfect for libhtml components: r.GetHTML("/", pages.Home().Render)
func (r *Router) GetHTML(path string, renderer func(w io.Writer) error) *mux.Route {
	return r.Get(path, HTML(renderer))
}

// PostHTML registers a POST route for HTML pages with automatic Content-Type header
func (r *Router) PostHTML(path string, renderer func(w io.Writer) error) *mux.Route {
	return r.Post(path, HTML(renderer))
}

// HandleRaw registers a raw http.HandlerFunc directly on the underlying mux router,
// bypassing all middleware. Use this for special endpoints like WebSockets.
func (r *Router) HandleRaw(path string, handler http.HandlerFunc) *mux.Route {
	return r.Router.NewRoute().Path(path).HandlerFunc(handler)
}

// ServeHTTP implements http.Handler, combining hot reload and main routers
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// If hot reload mux exists and path matches hot reload prefix, use it
	if r.hotReloadMux != nil && len(req.URL.Path) >= 12 && req.URL.Path[:12] == "/__hotreload" {
		r.hotReloadMux.ServeHTTP(w, req)
		return
	}
	// Otherwise use main router
	r.Router.ServeHTTP(w, req)
}
