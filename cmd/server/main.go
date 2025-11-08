package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/hyperstitieux/template/auth"
	"github.com/hyperstitieux/template/config"
	"github.com/hyperstitieux/template/controllers"
	"github.com/hyperstitieux/template/database"
	"github.com/hyperstitieux/template/database/repositories"
	"github.com/hyperstitieux/template/pages"
	"github.com/hyperstitieux/template/router"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		slog.Warn("no .env file found or error loading it", "error", err)
	}

	cfg := config.New()

	// Initialize database
	db, err := database.New(cfg.DatabaseURL)
	if err != nil {
		slog.Error("failed to initialize database", "error", err)
		panic(err)
	}
	defer db.Close()

	// Initialize repositories
	users := repositories.NewUsersRepository(db.DB)

	// Initialize controllers
	googleOAuthController := controllers.NewGoogleOAuthController(users, cfg.GoogleOAuthConfig)
	signOutController := controllers.NewSignOutController(users)
	settingsController := controllers.NewSettingsController(users)

	// Initialize router with default configuration
	// Note: Hot reload endpoints are registered separately to bypass middleware
	r, hotReloadMux := router.NewWithHotReload()

	// Register hot reload endpoints (development only) on separate mux
	// These bypass timeout/compression middleware to allow WebSocket connections
	hotReloadMux.HandleFunc("/__hotreload", router.HotReloadHandler)
	hotReloadMux.HandleFunc("/__hotreload_trigger", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			router.NotifyReload()
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	// Apply authentication middleware globally
	r.Use(auth.AuthMiddleware(users))

	// Serve static files from public directory (without /public/ prefix)
	fileServer := http.FileServer(http.Dir("./public"))
	r.PathPrefix("/js/").Handler(fileServer)
	r.PathPrefix("/css/").Handler(fileServer)
	r.PathPrefix("/images/").Handler(fileServer)
	r.PathPrefix("/fonts/").Handler(fileServer)
	r.HandleFunc("/styles.css", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./public/styles.css")
	})
	r.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./public/favicon.ico")
	})

	// Register routes
	r.Get("/", pages.Home)
	r.Get("/settings", pages.Settings)

	// OAuth routes
	r.Get("/auth/google", googleOAuthController.Redirect)
	r.Get("/auth/google/callback", googleOAuthController.Callback)
	r.Get("/auth/sign-out", signOutController.Handle)

	// Settings routes
	r.Post("/settings/update-profile", settingsController.UpdateProfile)
	r.Post("/settings/delete-account", settingsController.DeleteAccount)

	// Start HTTP server
	slog.Info("http server listening", "addr", cfg.HTTPAddr)
	if err := http.ListenAndServe(cfg.HTTPAddr, r); err != nil {
		slog.Error("failed to start http server", "error", err)
		os.Exit(1)
	}
}
