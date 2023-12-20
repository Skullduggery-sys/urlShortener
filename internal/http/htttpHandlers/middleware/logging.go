package middleware

import (
	"github.com/sirupsen/logrus"
	"net/http"
)

func LoggingMiddleware(logger *logrus.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger.Infof("Incomig route: %s, with Method: %s", r.URL.Path, r.Method)
			logger.Infof("IP: %s", r.RemoteAddr)

			next.ServeHTTP(w, r)
			logger.Info("HTTP method completed successfully")
		})
	}
}
