package httpRedirect

import (
	"errors"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
	"urlShortener/internal/http/httpUtils"
	"urlShortener/internal/http/htttpHandlers"
	"urlShortener/internal/storage"
)

type FullURLGetter interface {
	GetFullURL(shortenURL string) (string, error)
}

type Response struct {
	Error   string `json:"error,omitempty"`
	FullURL string `json:"fullURL,omitempty"`
}

func New(logger *logrus.Logger, getter FullURLGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "httpHandlers.httpRedirect.New"

		logger := logger.WithField(
			"function", fn,
		)

		shortenURL, ok := mux.Vars(r)[htttpHandlers.ShortenURLQuery]
		if !ok {
			http.Error(w, "something went wrong", http.StatusInternalServerError)
			return
		}

		fullURL, err := getter.GetFullURL(shortenURL)
		if errors.Is(err, storage.ErrURLNotFound) {
			logger.Info("url not found")
			err = httpUtils.RenderJSON(w, Response{Error: "can't find url"}, http.StatusNotFound)
			if err != nil {
				logger.Error("can't render JSON", "error", err.Error())
			}
			return
		} else if err != nil {
			logger.Error("error while getting full URL")
			err = httpUtils.RenderJSON(w, Response{Error: "can't get url sorry"}, http.StatusInternalServerError)
			if err != nil {
				logger.Error("error while rendering JSON")
			}
			return
		}

		http.Redirect(w, r, fullURL, http.StatusFound)
	}
}
