.PHONY: dev build css css-watch css-build clean

# Development mode with hot reload
dev:
	@air

# Build Go binary
build:
	@go build -o bin/server ./cmd/server

# Build TailwindCSS once
css-build:
	@bunx tailwindcss -i views/styles.css -o public/styles.css

# Watch TailwindCSS for changes
css-watch:
	@bunx tailwindcss -i views/styles.css -o public/styles.css --watch

# Alias for css-build
css: css-build

# Clean build artifacts
clean:
	@rm -rf tmp/ bin/ public/styles.css
