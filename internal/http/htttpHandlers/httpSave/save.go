package httpSave

import (
	"encoding/json"
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"net/http"
	"urlShortener/internal/http/httpUtils"
	"urlShortener/internal/storage"
	"urlShortener/utils"
)

type Request struct {
	FullURL string `json:"URL" validate:"required,url"`
}

type Response struct {
	Error      string `json:"error,omitempty"`
	ShortenURL string `json:"shortenURL,omitempty"`
}

type shortURLGetter interface {
	GetShortenURL(fullURL string) (string, error)
}

func New(logger *logrus.Logger, service shortURLGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "httpHandlers.httpSave.New"

		logger := logger.WithField("handler", fn)

		var req Request

		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&req)
		if err != nil {
			logger.Error("can't decode body", "error", err.Error())
			err = httpUtils.RenderJSON(w, Response{Error: "can't decode JSON"}, http.StatusBadRequest)
			if err != nil {
				logger.Error("rendering error", "error", err.Error())
			}
			return
		}
		logger.Info("Incoming URL", req.FullURL)

		if err = validator.New().Struct(&req); err != nil {
			logger.Error("can't validate request struct", "error", err.Error())
			resp := utils.ValidateErrors(err.(validator.ValidationErrors))
			err = httpUtils.RenderJSON(w, Response{Error: resp.Error()}, http.StatusBadRequest)
			if err != nil {
				logger.Error("rendering error", "error", err.Error())
			}
			return
		}

		shortenURL, err := service.GetShortenURL(req.FullURL)
		if errors.Is(err, storage.ErrURLExists) {
			logger.Error("maybe id leaking", "error", err.Error())
			err = httpUtils.RenderJSON(w, Response{
				Error: "error happened while getting URL",
			}, http.StatusInternalServerError)
			if err != nil {
				logger.Error("error while rendering JSON", "error", err.Error())
			}
			return
		} else if err != nil {
			logger.Error("error while getting shortenURL", "error", err.Error())
			err = httpUtils.RenderJSON(w, Response{
				Error: "error happened while trying to get short url sorry",
			}, http.StatusInternalServerError)
			if err != nil {
				logger.Error("error while rendering JSON", "error", err.Error())
			}
			return
		}

		err = httpUtils.RenderJSON(w, Response{
			ShortenURL: shortenURL,
		}, http.StatusOK)
		if err != nil {
			logger.Error("error rendering", "error", err.Error())
		}
	}
}
