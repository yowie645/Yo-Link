package redirect

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"

	resp "github.com/yowie645/Yo-Link/internal/lib/api/response"
	"github.com/yowie645/Yo-Link/internal/lib/logger/sl"
	"github.com/yowie645/Yo-Link/internal/storage"
)

type URLGetter interface {
	GetURL(alias string) (string, error)
}

// go⁡⁣⁢⁣:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=URLGetter⁡
func New(log *slog.Logger, urlGetter URLGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.redirect.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")
		log.Debug("Handling redirect request",
			slog.String("alias", alias),
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
		)

		if alias == "" {
			log.Error("empty alias parameter")
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, resp.Error("alias parameter is required"))
			return
		}

		resURL, err := urlGetter.GetURL(alias)
		if errors.Is(err, storage.ErrURLNotFound) {
			log.Info("url not found", slog.String("alias", alias))
			render.Status(r, http.StatusNotFound)
			render.JSON(w, r, resp.Error("url not found"))
			return
		}

		if err != nil {
			log.Error("failed to get url", sl.Err(err))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("internal server error"))
			return
		}

		log.Info("redirecting",
			slog.String("from", alias),
			slog.String("to", resURL),
		)
		http.Redirect(w, r, resURL, http.StatusFound)
	}
}
