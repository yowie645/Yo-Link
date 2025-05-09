package redirect

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/yowie645/Yo-Link/internal/lib/logger/sl"
	"github.com/yowie645/Yo-Link/internal/storage"
)

type URLGetter interface {
	GetURL(alias string) (string, error)
}

func New(log *slog.Logger, urlGetter URLGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.redirect.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")
		if alias == "" {
			log.Error("empty alias parameter")
			http.Error(w, "alias parameter is required", http.StatusBadRequest)
			return
		}

		resURL, err := urlGetter.GetURL(alias)
		if errors.Is(err, storage.ErrURLNotFound) {
			log.Info("url not found", slog.String("alias", alias))
			http.Error(w, "url not found", http.StatusNotFound)
			return
		}

		if err != nil {
			log.Error("failed to get url", sl.Err(err))
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		log.Info("redirecting",
			slog.String("from", alias),
			slog.String("to", resURL),
		)

		http.Redirect(w, r, resURL, http.StatusFound)
	}
}
