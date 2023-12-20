package httpRedirect

import (
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
	"urlShortener/internal/http/htttpHandlers"
	"urlShortener/internal/storage"
)

type mockURLGetter struct {
	mock.Mock
}

func (m *mockURLGetter) GetFullURL(shortenURL string) (string, error) {
	args := m.Called(shortenURL)
	return args.String(0), args.Error(1)
}

func TestNewSuccess(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.PanicLevel)
	getter := &mockURLGetter{}
	handler := New(logger, getter)

	getter.On("GetFullURL", "known").Return("https://911.com", nil)

	req := httptest.NewRequest(http.MethodGet, "/known", nil)
	req = mux.SetURLVars(req, map[string]string{htttpHandlers.ShortenURLQuery: "known"})

	w := httptest.NewRecorder()
	handler(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusFound, resp.StatusCode)
	assert.Equal(t, "https://911.com", resp.Header.Get("Location"))

	getter.AssertExpectations(t)
}

func TestNewURLNotFound(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.PanicLevel)
	getter := &mockURLGetter{}
	handler := New(logger, getter)

	getter.On("GetFullURL", "bbbbb").Return("", storage.ErrURLNotFound)

	req := httptest.NewRequest(http.MethodGet, "/unknown", nil)
	req = mux.SetURLVars(req, map[string]string{htttpHandlers.ShortenURLQuery: "bbbbb"})

	w := httptest.NewRecorder()
	handler(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)

	getter.AssertExpectations(t)
}

func TestNewQueryError(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.PanicLevel)
	getter := &mockURLGetter{}
	handler := New(logger, getter)

	req := httptest.NewRequest(http.MethodGet, "/notok", nil)

	w := httptest.NewRecorder()
	handler(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

	getter.AssertExpectations(t)
}
