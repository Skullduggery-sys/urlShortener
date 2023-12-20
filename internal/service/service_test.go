package service

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"urlShortener/internal/lib/linkShortening/hashByID"
	"urlShortener/internal/storage"
)

const (
	getFullURL    = "GetFullURL"
	saveURL       = "SaveURL"
	getShortenURL = "GetShortenURL"
	hash          = "Hash"
)

type mockStorager struct {
	mock.Mock
}

func (m *mockStorager) SaveURL(urlToSave string, shortenURL string) error {
	args := m.Called(urlToSave, shortenURL)
	return args.Error(0)
}

func (m *mockStorager) GetFullURL(shortenURL string) (string, error) {
	args := m.Called(shortenURL)
	return args.String(0), args.Error(1)
}

func (m *mockStorager) GetShortenURL(fullURL string) (string, error) {
	args := m.Called(fullURL)
	return args.String(0), args.Error(1)
}

type mockHasher struct {
	mock.Mock
}

func (m *mockHasher) Hash() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

func TestGetShortenURLNotFound(t *testing.T) {
	mockStorage := &mockStorager{}
	mockHash := &mockHasher{}
	service := New(mockStorage, mockHash)

	fullurl := "ozon.ru"
	expextedShortenURL := "aaaaaaaaaa"
	mockStorage.On(getShortenURL, fullurl).Return("", storage.ErrURLNotFound)
	mockHash.On(hash).Return(expextedShortenURL, nil)
	mockStorage.On(saveURL, fullurl, expextedShortenURL).Return(nil)

	resultShortenURL, err := service.GetShortenURL(fullurl)
	assert.NoError(t, err)
	assert.Equal(t, resultShortenURL, expextedShortenURL)

	assert.True(t, mockStorage.AssertExpectations(t))
	assert.True(t, mockHash.AssertExpectations(t))
}

func TestGetShortenURLFound(t *testing.T) {
	mockStorage := &mockStorager{}
	mockHash := &mockHasher{}
	service := New(mockStorage, mockHash)

	fullurl := "ozon.ru"
	expextedShortenURL := "aaaaaaaaaa"
	mockStorage.On(getShortenURL, fullurl).Return(expextedShortenURL, nil)

	resultShortenURL, err := service.GetShortenURL(fullurl)
	assert.NoError(t, err)
	assert.Equal(t, resultShortenURL, expextedShortenURL)

	assert.True(t, mockStorage.AssertExpectations(t))
}

func TestGetShortenUnexpectedError(t *testing.T) {
	mockStorage := &mockStorager{}
	mockHash := &mockHasher{}
	service := New(mockStorage, mockHash)

	fullurl := "ozon.ru"
	mockStorage.On(getShortenURL, fullurl).Return("", errors.New("unknown"))

	_, err := service.GetShortenURL(fullurl)
	assert.Error(t, err)

	assert.True(t, mockStorage.AssertExpectations(t))
	assert.True(t, mockHash.AssertExpectations(t))
}

func TestGetShortenOverflow(t *testing.T) {
	mockStorage := &mockStorager{}
	mockHash := &mockHasher{}
	service := New(mockStorage, mockHash)

	fullurl := "ozon.ru"
	mockStorage.On(getShortenURL, fullurl).Return("", storage.ErrURLNotFound)
	mockHash.On(hash).Return("", hashByID.ErrOverFlow)

	_, err := service.GetShortenURL(fullurl)
	assert.True(t, errors.Is(err, hashByID.ErrOverFlow))

	assert.True(t, mockStorage.AssertExpectations(t))
	assert.True(t, mockHash.AssertExpectations(t))
}

func TestGetShortenUnknownError(t *testing.T) {
	mockStorage := &mockStorager{}
	mockHash := &mockHasher{}
	service := New(mockStorage, mockHash)

	fullurl := "ozon.ru"
	mockStorage.On(getShortenURL, fullurl).Return("", storage.ErrURLNotFound)
	mockHash.On(hash).Return("", errors.New("wtf just happend i fell asleep"))

	_, err := service.GetShortenURL(fullurl)
	assert.Error(t, err)

	assert.True(t, mockStorage.AssertExpectations(t))
	assert.True(t, mockHash.AssertExpectations(t))
}

func TestGetShortenSavingError(t *testing.T) {
	mockStorage := &mockStorager{}
	mockHash := &mockHasher{}
	service := New(mockStorage, mockHash)

	fullurl := "ozon.ru"
	expextedShortenURL := "aaaaaaaaaa"
	mockStorage.On(getShortenURL, fullurl).Return("", storage.ErrURLNotFound)
	mockHash.On(hash).Return(expextedShortenURL, nil)
	mockStorage.On(saveURL, fullurl, expextedShortenURL).Return(errors.New("unknown"))

	_, err := service.GetShortenURL(fullurl)
	assert.Error(t, err)

	assert.True(t, mockStorage.AssertExpectations(t))
	assert.True(t, mockHash.AssertExpectations(t))
}

func TestGetFullURLSuccess(t *testing.T) {
	mockStorage := &mockStorager{}
	mockHash := &mockHasher{}
	service := New(mockStorage, mockHash)

	expectedFullURL := "ozon.ru"
	shortenURL := "aaaaaaaaaa"
	mockStorage.On(getFullURL, shortenURL).Return(expectedFullURL, nil)

	resultFullURL, err := service.GetFullURL(shortenURL)
	assert.NoError(t, err)
	assert.Equal(t, resultFullURL, expectedFullURL)

	assert.True(t, mockStorage.AssertExpectations(t))
	assert.True(t, mockHash.AssertExpectations(t))
}

func TestGetFullURLNotFound(t *testing.T) {
	mockStorage := &mockStorager{}
	mockHash := &mockHasher{}
	service := New(mockStorage, mockHash)

	shortenURL := "aaaaaaaaaa"
	mockStorage.On(getFullURL, shortenURL).Return("", storage.ErrURLNotFound)

	_, err := service.GetFullURL(shortenURL)
	assert.True(t, errors.Is(err, storage.ErrURLNotFound))

	assert.True(t, mockStorage.AssertExpectations(t))
	assert.True(t, mockHash.AssertExpectations(t))
}

func TestGetFullURLUnknownError(t *testing.T) {
	mockStorage := &mockStorager{}
	mockHash := &mockHasher{}
	service := New(mockStorage, mockHash)

	shortenURL := "aaaaaaaaaa"
	mockStorage.On(getFullURL, shortenURL).Return("", errors.New("unknown"))

	_, err := service.GetFullURL(shortenURL)
	assert.Error(t, err)

	assert.True(t, mockStorage.AssertExpectations(t))
	assert.True(t, mockHash.AssertExpectations(t))
}
