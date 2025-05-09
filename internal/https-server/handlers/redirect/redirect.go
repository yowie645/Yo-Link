package redirect

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"

	resp "github.com/yowie645/Yo-Link/internal/lib/api/response"
	"github.com/yowie645/Yo-Link/internal/lib/logger/sl"
	"github.com/yowie645/Yo-Link/internal/storage"
)

type URLGetter interface {
	GetURL(alias string) (string, error)
}

//go⁡⁣⁢⁡⁣⁢⁣:generate go run github.com/vektra/mockery/v2@v2.53.3 --name=URLGetter⁡

func New(log *slog.Logger, urlGetter URLGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.redirect.New"

		log = log.With(slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")

		if alias == "" {
			log.Info("alias is empty")

			render.JSON(w, r, resp.Error("invalid request"))

			return
		}

		resURL, err := urlGetter.GetURL(alias)
		if errors.Is(err, storage.ErrURLNotFound) {
			log.Info("url not found", slog.String("alias", alias))

			render.JSON(w, r, resp.Error("url not found"))

			return

		}

		if err != nil {
			log.Error("failed to get url", sl.Err(err))

			render.JSON(w, r, resp.Error("internal error"))

			return
		}

		log.Info("got url", slog.String("url", resURL))

		http.Redirect(w, r, resURL, http.StatusFound)

	}
}
