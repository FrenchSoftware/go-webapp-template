# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Go web application template built by "The French Software" for creating server-rendered web applications with Google OAuth authentication. The stack uses Go with server-side rendering via the `libhtml` library (French Software's HTML generation library), SQLite for persistence, and Basecoat UI for styling.

## Development Commands

### Running the Application

- `make dev` - Start development server with hot reload (uses Air)
- `make build` - Build Go binary to `bin/server`
- `go run ./cmd/server` - Run server directly without hot reload

### CSS/Styling

- `make css` or `make css-build` - Build TailwindCSS once (views/styles.css → public/styles.css)
- `make css-watch` - Watch and rebuild TailwindCSS on changes
- `make clean` - Remove build artifacts (tmp/, bin/, public/styles.css)

### Hot Reload Configuration

The project uses Air for hot reload. Configuration is in `.air.toml`:
- Runs `bunx tailwindcss` as pre-command to rebuild CSS before each reload
- Triggers hot reload via WebSocket at `/__hotreload` endpoint
- Watches `.go`, `.tpl`, `.tmpl`, `.html`, `.css` files
- Build output: `./tmp/main`

## Architecture

### Project Structure

```
cmd/server/main.go          # Application entry point
config/                     # Configuration management
database/
  ├── database.go           # SQLite connection and schema migration
  ├── schema.sql            # Embedded SQL schema
  ├── models/               # Database models (User, etc.)
  └── repositories/         # Data access layer (UsersRepository)
router/
  ├── router.go             # Router initialization with middleware chain
  ├── helpers.go            # Router wrapper with Get/Post/HTML helpers
  ├── middleware.go         # Custom middleware (recovery, logging, etc.)
  ├── errors.go             # Error handling
  └── hotreload.go          # WebSocket-based hot reload
auth/
  ├── auth.go               # Context-based auth utilities
  ├── middleware.go         # Auth middleware
  └── cookie.go             # Session cookie management
controllers/                # HTTP controllers for features
pages/                      # Page handlers (Home, Settings, etc.)
views/
  ├── layouts/base.go       # Base HTML layout
  ├── components/           # Reusable UI components
  │   └── ui/               # Basecoat UI wrappers (Card, Dialog, etc.)
  ├── basecoat.css          # Basecoat UI CSS framework
  └── styles.css            # Custom Tailwind CSS (input file)
public/                     # Static assets (CSS, JS, images)
```

### Key Architectural Patterns

**Server-Side Rendering with libhtml**
- Uses `github.com/frenchsoftware/libhtml` for type-safe HTML generation in Go
- Components return `html.Node` types
- Pages use `layouts.Base()` wrapper for consistent structure
- Example: `html.Div(attr.Class("..."), html.Text("Hello"))`

**Router Pattern**
- Custom router wrapper around `gorilla/mux` with helper methods
- Routes defined as: `r.Get("/path", handlerFunc)` or `r.Post("/path", handlerFunc)`
- Handler signature: `func(w http.ResponseWriter, r *http.Request) error`
- Error handling wrapped automatically via `Handle()` middleware
- HTML helpers: `router.HTML(renderer)` sets Content-Type automatically

**Middleware Chain** (via `justinas/alice`)
1. Recovery (panic handling)
2. Request ID
3. Structured logging (slog)
4. Security headers (CSP, HSTS, etc.)
5. CORS (configurable origins)
6. Compression (gzip)
7. Rate limiting (token bucket)
8. Timeout (30s default)

Note: Hot reload endpoints (`/__hotreload`, `/__hotreload_trigger`) bypass middleware to allow WebSocket connections

**Authentication Flow**
- Google OAuth 2.0 via `golang.org/x/oauth2`
- Session tokens stored in SQLite `sessions` table
- Auth middleware checks session cookie and loads user into request context
- Access user: `auth.GetCurrentUser(r)` returns `*models.User` or `nil`
- Protected pages: redirect to `/auth/google?redirect=/settings`

**Database Pattern**
- SQLite via `modernc.org/sqlite` (pure Go, no CGo)
- Schema embedded via `//go:embed schema.sql`
- Repository pattern for data access
- Foreign keys enabled explicitly (disabled by default in SQLite)
- Migrations run automatically on startup

**Environment Configuration**
- Uses `github.com/joho/godotenv` to load `.env` file
- Helper: `env.GetVar("KEY", "default")` in `env/env.go`
- Required vars: `GOOGLE_CLIENT_ID`, `GOOGLE_CLIENT_SECRET`, `BASE_URL`
- Optional vars: `HTTP_ADDR` (default `:8080`), `DATABASE_URL` (default `file:app.db`)

## Basecoat UI Usage

This project uses Basecoat UI, a Tailwind CSS-based component library that works without React. The CSS is in `views/basecoat.css`.

### Available Components and Classes

**Buttons**
- `.btn`, `.btn-primary` - Primary button (default size)
- `.btn-secondary`, `.btn-outline`, `.btn-ghost`, `.btn-link`, `.btn-destructive`
- Size variants: `.btn-sm`, `.btn-lg`
- Icon buttons: `.btn-icon`, `.btn-sm-icon`, `.btn-lg-icon`
- Combine: `.btn-sm-outline`, `.btn-lg-destructive`, etc.

**Cards**
- `.card` - Card container
- Use semantic `<header>`, `<section>`, `<footer>` inside cards
- Go helpers: `ui.Card()`, `ui.CardHeader()`, `ui.CardSection()`, `ui.CardFooter()`

**Forms & Inputs**
- `.field` - Form field container
- `.fieldset` - Group of fields with `<legend>`
- `.input` - Apply to `<input>` elements for styling (or wrap in `.field` for automatic styling)
- `.label` - Form label styling
- Input types supported: text, email, password, number, file, tel, url, search, date, etc.
- `.textarea` - Textarea styling
- Checkboxes and radio buttons styled automatically inside `.field`

**Switch**
- `<input type="checkbox" role="switch">` - Toggle switch

**Select**
- `<select class="select">` - Styled select dropdown
- Custom select: `.select` wrapper with `[data-popover]` and `[role="option"]`

**Dialog**
- `.dialog` - Modal dialog
- Use `<dialog>` element with `showModal()` method
- Go helper: `ui.AlertDialog()` with props

**Badges**
- `.badge`, `.badge-primary`, `.badge-secondary`, `.badge-destructive`, `.badge-outline`

**Alerts**
- `.alert` - Alert box
- `.alert-destructive` - Error/destructive alert

**Tables**
- `.table` - Styled table

**Tabs**
- `.tabs` - Tab container
- `[role="tablist"]` and `[role="tab"]` inside

**Tooltips**
- `[data-tooltip="Text"]` - Hover tooltip
- Positioning: `data-side="top|bottom|left|right"`, `data-align="start|center|end"`

**Other Components**
- `.kbd` - Keyboard shortcut display
- `.dropdown-menu` - Dropdown menu
- `.command` - Command palette
- `.popover` - Popover container
- `.sidebar` - Sidebar navigation
- `.toaster` and `.toast` - Toast notifications

### Color Tokens
- `bg-background`, `text-foreground`
- `bg-card`, `text-card-foreground`
- `bg-primary`, `text-primary-foreground`
- `bg-secondary`, `text-secondary-foreground`
- `bg-muted`, `text-muted-foreground`
- `bg-accent`, `text-accent-foreground`
- `bg-destructive`
- `border-border`, `border-input`

## Common Patterns

### Creating a New Page

1. Add handler in `pages/pages.go`:
```go
func MyPage(w http.ResponseWriter, r *http.Request) error {
    user := auth.GetCurrentUser(r)
    page := layouts.Base(user, r, "Page Title",
        html.Div(
            attr.Class("container"),
            html.Text("Content"),
        ),
    )
    w.Header().Set("Content-Type", "text/html; charset=utf-8")
    return page.Render(w)
}
```

2. Register route in `cmd/server/main.go`:
```go
r.Get("/mypage", pages.MyPage)
```

### Adding a New Controller

1. Create file in `controllers/` with struct and constructor:
```go
type MyController struct {
    users *repositories.UsersRepository
}

func NewMyController(users *repositories.UsersRepository) *MyController {
    return &MyController{users: users}
}

func (c *MyController) Handle(w http.ResponseWriter, r *http.Request) error {
    // Handle request
    return nil
}
```

2. Initialize and register in `cmd/server/main.go`

### Working with Forms

Use `github.com/frenchsoftware/libvalidator` for validation:
```go
func UpdateProfile(w http.ResponseWriter, r *http.Request) error {
    v := validator.New()
    name := v.Required("name", r.FormValue("name"), "Name is required")

    if v.HasErrors() {
        return pages.SettingsWithErrors(w, r, v.Errors)
    }

    // Process valid data
}
```

### Creating UI Components

In `views/components/ui/`:
```go
func MyComponent(children ...html.Node) html.Node {
    args := []any{attr.Class("my-classes")}
    for _, child := range children {
        args = append(args, child)
    }
    return html.Div(args...)
}
```

## Testing

No test framework is currently set up. Add tests using Go's built-in `testing` package or a framework like `testify`.

## Deployment

1. Build: `make build` (outputs to `bin/server`)
2. Set environment variables (especially `GOOGLE_CLIENT_ID`, `GOOGLE_CLIENT_SECRET`, `BASE_URL`)
3. Run: `./bin/server`
4. Database file created automatically at `DATABASE_URL` location

## Dependencies

Key dependencies:
- `github.com/frenchsoftware/libhtml` - HTML generation library
- `github.com/frenchsoftware/libvalidator` - Form validation
- `github.com/gorilla/mux` - HTTP router
- `github.com/justinas/alice` - Middleware chaining
- `golang.org/x/oauth2` - OAuth 2.0 client
- `modernc.org/sqlite` - Pure Go SQLite driver
- `github.com/joho/godotenv` - Environment variable loading

All dependencies managed via Go modules (`go.mod`).
