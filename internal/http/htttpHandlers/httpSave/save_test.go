package httpSave

import (
	"bytes"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
	"urlShortener/internal/storage"
)

type mockShortURLGetter struct {
	mock.Mock
}

func (m *mockShortURLGetter) GetShortenURL(fullURL string) (string, error) {
	args := m.Called(fullURL)
	return args.String(0), args.Error(1)
}

func TestNewSuccess(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.PanicLevel)
	service := mockShortURLGetter{}
	handler := New(logger, &service)

	reqBody := `{"URL": "https://bmstu.com"}`
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")

	service.On("GetShortenURL", "https://bmstu.com").Return("abcabcabc", nil)

	w := httptest.NewRecorder()
	handler(w, req)

	resp := w.Result()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var response Response
	err := json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(t, err)

	assert.Equal(t, "abcabcabc", response.ShortenURL)

	service.AssertExpectations(t)
}

func TestNewValidationError(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.PanicLevel)
	service := mockShortURLGetter{}
	handler := New(logger, &service)

	reqBody := `{"URL": "123456789"}`
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	handler(w, req)

	resp := w.Result()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	var response Response
	err := json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(t, err)

	assert.Equal(t, "field FullURL url is wrong", response.Error)

	service.AssertExpectations(t)
}

func TestNewStorageError(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.PanicLevel)
	service := mockShortURLGetter{}
	handler := New(logger, &service)

	reqBody := `{"URL": "https://opposite.com"}`
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")

	service.On("GetShortenURL", "https://opposite.com").Return("", storage.ErrURLExists)

	w := httptest.NewRecorder()
	handler(w, req)

	resp := w.Result()

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

	var response Response
	err := json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(t, err)

	assert.Equal(t, "error happened while getting URL", response.Error)

	service.AssertExpectations(t)
}

func TestNewDecodeError(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.PanicLevel)
	service := mockShortURLGetter{}
	handler := New(logger, &service)

	reqBody := `invalid json`
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	handler(w, req)

	resp := w.Result()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	var response Response
	err := json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(t, err)

	assert.Equal(t, "can't decode JSON", response.Error)

	service.AssertExpectations(t)
}
