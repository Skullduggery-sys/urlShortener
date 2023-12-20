package router

import (
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
	"urlShortener/internal/http/htttpHandlers"
	"urlShortener/internal/http/htttpHandlers/httpRedirect"
	"urlShortener/internal/http/htttpHandlers/httpSave"
	"urlShortener/internal/http/htttpHandlers/middleware"
)

const (
	saveRoute     = "/"
	redirectRoute = "/{" + htttpHandlers.ShortenURLQuery + "}"
)

type Service interface {
	GetShortenURL(fullURL string) (string, error)
	GetFullURL(shortenURL string) (string, error)
}

func New(log *logrus.Logger, service Service) *mux.Router {
	r := mux.NewRouter()

	r.Handle(saveRoute, httpSave.New(log, service)).Methods(http.MethodPost)
	r.Handle(redirectRoute, httpRedirect.New(log, service)).Methods(http.MethodGet)
	r.Use(middleware.LoggingMiddleware(log))

	return r
}
