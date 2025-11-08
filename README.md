[![Discord](https://img.shields.io/badge/Discord-Join%20Us-7289DA?logo=discord&logoColor=white)](https://discord.gg/MaqTgPF3Ch)

# go-webapp-template

> The French Software Go webapp template

A modern Go web application template with server-side rendering, Google OAuth authentication, and type-safe HTML generation.

## Stack

- **Go 1.25** - Backend language
- **[libhtml](https://github.com/frenchsoftware/libhtml)** - Type-safe HTML generation library
- **[Basecoat UI](https://www.basecoat-ui.com/)** - Tailwind CSS component library
- **TailwindCSS** - Utility-first CSS framework
- **SQLite** - Embedded database (via [modernc.org/sqlite](https://gitlab.com/cznic/sqlite))
- **Google OAuth 2.0** - Authentication
- **[Air](https://github.com/air-verse/air)** - Hot reload for development

## Prerequisites

- **Go 1.25+** - [Download](https://go.dev/dl/)
- **Bun** - Required for TailwindCSS bundling: [Installation](https://bun.sh/docs/installation)

## Installation

1. Clone the repository:
```bash
git clone <your-repo-url>
cd go-webapp-template
```

2. Install Go dependencies:
```bash
go mod download
```

3. Install Air for hot reload:
```bash
go install github.com/air-verse/air@latest
```

4. Set up environment variables:
```bash
cp .env.example .env
# Edit .env with your Google OAuth credentials
```

5. Get Google OAuth credentials:
   - Visit [Google Cloud Console](https://console.cloud.google.com/apis/credentials)
   - Create OAuth 2.0 Client ID
   - Add authorized redirect URI: `http://localhost:8080/auth/google/callback`
   - Copy Client ID and Client Secret to `.env`

## Development

Start the development server with hot reload:

```bash
make dev
```

This will:
- Build and run the Go server
- Automatically rebuild CSS with TailwindCSS on file changes
- Hot reload the browser on Go/HTML/CSS changes
- Server runs at http://localhost:8080

### CSS Development

Build CSS once:
```bash
make css
```

Watch CSS for changes (alternative to `make dev`):
```bash
make css-watch
```

### Other Commands

```bash
make build    # Build production binary to bin/server
make clean    # Remove build artifacts (tmp/, bin/, public/styles.css)
```

## Architecture

### Server-Side Rendering with libhtml

Pages are rendered server-side using Go with the `libhtml` library for type-safe HTML generation:

```go
import "github.com/frenchsoftware/libhtml"

func MyPage(w http.ResponseWriter, r *http.Request) error {
    page := html.Div(
        attr.Class("container mx-auto"),
        html.H1(html.Text("Hello World")),
    )
    return page.Render(w)
}
```

### Basecoat UI Components

The template includes [Basecoat UI](https://www.basecoat-ui.com/), a Tailwind CSS component library that works without React. Use pre-built component classes:

```go
html.Button(
    attr.Class("btn-primary"),
    html.Text("Click Me")
)
```

Available components: buttons, cards, forms, dialogs, badges, alerts, tables, tabs, and more. See `views/basecoat.css` or [Basecoat UI documentation](https://www.basecoat-ui.com/docs).

## Configuration

Environment variables (`.env`):

| Variable | Description | Default |
|----------|-------------|---------|
| `HTTP_ADDR` | Server address and port | `:8080` |
| `DATABASE_URL` | SQLite database file path | `file:app.db` |
| `BASE_URL` | Application base URL (for OAuth) | `http://localhost:8080` |
| `GOOGLE_CLIENT_ID` | Google OAuth Client ID | *Required* |
| `GOOGLE_CLIENT_SECRET` | Google OAuth Client Secret | *Required* |

## Deployment

1. Build the binary:
```bash
make build
```

2. Set production environment variables

3. Run the server:
```bash
./bin/server
```

The SQLite database will be created automatically on first run.

## License

[AGPL-3.0](./LICENSE)
