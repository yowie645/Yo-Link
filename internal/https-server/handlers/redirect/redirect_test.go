package redirect_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"

	"github.com/stretchr/testify/require"
	"github.com/yowie645/Yo-Link/internal/https-server/handlers/redirect"
	"github.com/yowie645/Yo-Link/internal/https-server/handlers/redirect/mocks"
	"github.com/yowie645/Yo-Link/internal/lib/logger/handlers/slogdiscard"
	"github.com/yowie645/Yo-Link/internal/storage"
)

func TestRedirectHandler(t *testing.T) {
	cases := []struct {
		name      string
		alias     string
		url       string
		error     error
		respCode  int
		location  string
		respError string
	}{
		{
			name:     "Success",
			alias:    "test_alias",
			url:      "https://google.com",
			respCode: http.StatusFound,
			location: "https://google.com",
		},
		{
			name:      "Empty alias",
			alias:     "",
			respCode:  http.StatusBadRequest,
			respError: "invalid request",
		},
		{
			name:      "URL not found",
			alias:     "not_found",
			error:     storage.ErrURLNotFound,
			respCode:  http.StatusNotFound,
			respError: "url not found",
		},
		{
			name:      "Internal error",
			alias:     "test_error",
			error:     errors.New("unexpected error"),
			respCode:  http.StatusInternalServerError,
			respError: "internal error",
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			urlGetterMock := mocks.NewURLGetter(t)

			if tc.alias != "" && tc.error != storage.ErrURLNotFound {
				urlGetterMock.On("GetURL", tc.alias).
					Return(tc.url, tc.error).
					Once()
			} else if tc.error == storage.ErrURLNotFound {
				urlGetterMock.On("GetURL", tc.alias).
					Return("", tc.error).
					Once()
			}

			handler := redirect.New(slogdiscard.NewDiscardLogger(), urlGetterMock)

			req, err := http.NewRequest(http.MethodGet, "/"+tc.alias, nil)
			require.NoError(t, err)

			// Устанавливаем параметр маршрута вручную
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("alias", tc.alias)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			require.Equal(t, tc.respCode, rr.Code)

			if tc.error == nil && tc.alias != "" {
				require.Equal(t, tc.location, rr.Header().Get("Location"))
			} else if tc.respError != "" {
				var resp map[string]interface{}
				err := json.Unmarshal(rr.Body.Bytes(), &resp)
				require.NoError(t, err)
				require.Equal(t, tc.respError, resp["error"])
			}
		})
	}
}
