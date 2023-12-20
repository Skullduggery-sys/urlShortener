package save

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"urlShortener/internal/gRPC/proto"
)

const getShortenURL = "GetShortenURL"

type mockShortUrlGetter struct {
	mock.Mock
}

func (m *mockShortUrlGetter) GetShortenURL(fullURL string) (string, error) {
	args := m.Called(fullURL)
	return args.String(0), args.Error(1)
}

func TestSaveSuccess(t *testing.T) {
	getter := mockShortUrlGetter{}
	handlerSave := New(&getter)

	fullURL := proto.FullURL{URL: "https://ozon.ru"}
	expectedShortenURL := proto.ShortURL{URL: "iii098iiii"}
	getter.On(getShortenURL, fullURL.URL).Return(expectedShortenURL.URL, nil)

	resultShortenURL, err := handlerSave.Save(context.Background(), &fullURL)
	assert.NoError(t, err)
	assert.Equal(t, resultShortenURL.URL, expectedShortenURL.URL)

	assert.True(t, getter.AssertExpectations(t))
}

func TestSaveInvalidURL(t *testing.T) {
	getter := mockShortUrlGetter{}
	handlerSave := New(&getter)

	fullURL := proto.FullURL{URL: "ozon.ru"}

	_, err := handlerSave.Save(context.Background(), &fullURL)
	assert.Error(t, err)

	assert.True(t, getter.AssertExpectations(t))
}

func TestSaveErrShortening(t *testing.T) {
	getter := mockShortUrlGetter{}
	handlerSave := New(&getter)

	fullURL := proto.FullURL{URL: "https://ozon.ru"}
	getter.On(getShortenURL, fullURL.URL).Return("", errors.New("unknown"))

	_, err := handlerSave.Save(context.Background(), &fullURL)
	assert.Error(t, err)

	assert.True(t, getter.AssertExpectations(t))
}
