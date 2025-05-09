package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/yowie645/Yo-Link/internal/config"
	"github.com/yowie645/Yo-Link/internal/https-server/handlers/redirect"
	"github.com/yowie645/Yo-Link/internal/https-server/handlers/url/save"
	"github.com/yowie645/Yo-Link/internal/lib/logger/handlers/slogpretty"
	"github.com/yowie645/Yo-Link/internal/lib/logger/sl"
	"github.com/yowie645/Yo-Link/internal/storage/sqlite"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()
	log := setupLogger(cfg.Env)
	log.Info("starting server", slog.String("env", cfg.Env), slog.String("address", cfg.Address))

	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		log.Error("failed to init storage", sl.Err(err))
		os.Exit(1)
	}

	// тестовый alias
	if cfg.Env == envLocal {
		if _, err := storage.SaveURL("https://google.com", "test_alias"); err != nil {
			log.Error("failed to save test url", sl.Err(err))
		}
	}

	router := chi.NewRouter()

	// Middleware
	router.Use(
		middleware.RequestID,
		middleware.RealIP,
		middleware.Logger,
		middleware.Recoverer,
		render.SetContentType(render.ContentTypeJSON),
	)

	router.Route("/url", func(r chi.Router) {
		r.Use(middleware.BasicAuth("yo-link", map[string]string{
			cfg.HTTPServer.User: cfg.HTTPServer.Password,
		}))
	})

	router.Post("/url", save.New(log, storage))

	router.Get("/{alias}", redirect.New(log, storage))

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]interface{}{
			"status": "Error",
			"error":  "alias is required",
		})
	})

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Error("server stopped", sl.Err(err))
	}
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case envLocal:
		log = setupPrettySlog()
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}
	return log
}

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{Level: slog.LevelDebug},
	}
	handler := opts.NewPrettyHandler(os.Stdout)
	return slog.New(handler)
}
