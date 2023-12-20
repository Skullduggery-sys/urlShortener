package inMemmory

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
	"urlShortener/internal/storage"
)

func TestGetFullURLSuccess(t *testing.T) {
	st := New()
	fullURL := "ya.ru"
	shortURL := "aaaaaaaaa"
	st.keyFullURL[fullURL] = shortURL
	st.keyShortenURL[shortURL] = fullURL

	resultFullURL, err := st.GetFullURL(shortURL)
	assert.NoError(t, err)
	assert.Equal(t, resultFullURL, fullURL)
}

func TestGetFullURLNotFound(t *testing.T) {
	st := New()
	shortURL := "aaaaaaaaa"

	_, err := st.GetFullURL(shortURL)
	assert.True(t, errors.Is(err, storage.ErrURLNotFound))
}

func TestGetShortenURLSuccess(t *testing.T) {
	st := New()
	fullURL := "ya.ru"
	shortURL := "aaaaaaaaa"
	st.keyFullURL[fullURL] = shortURL
	st.keyShortenURL[shortURL] = fullURL

	resultShortURL, err := st.GetShortenURL(fullURL)
	assert.NoError(t, err)
	assert.Equal(t, resultShortURL, shortURL)
}

func TestGetShortenURLNotFound(t *testing.T) {
	st := New()
	fullURL := "ya.ru"

	_, err := st.GetShortenURL(fullURL)
	assert.True(t, errors.Is(err, storage.ErrURLNotFound))
}

func TestSaveURLSuccess(t *testing.T) {
	st := New()
	fullURL := "ya.ru"
	shortURL := "aaaaaaaaa"

	err := st.SaveURL(fullURL, shortURL)
	assert.NoError(t, err)
}

func TestSaveURLAlreadyExist(t *testing.T) {
	st := New()
	fullURL := "ya.ru"
	shortURL := "aaaaaaaaa"

	st.keyFullURL[fullURL] = shortURL
	st.keyShortenURL[shortURL] = fullURL

	err := st.SaveURL(fullURL, shortURL)
	assert.True(t, errors.Is(err, storage.ErrURLExists))
}
