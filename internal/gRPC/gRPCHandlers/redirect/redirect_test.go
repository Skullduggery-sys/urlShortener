package redirect

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"urlShortener/internal/gRPC/proto"
)

const getFullURL = "GetFullURL"

type mockFullUrlGetter struct {
	mock.Mock
}

func (m *mockFullUrlGetter) GetFullURL(shortURL string) (string, error) {
	args := m.Called(shortURL)
	return args.String(0), args.Error(1)
}

func TestRedirectSuccess(t *testing.T) {
	getter := &mockFullUrlGetter{}
	handler := New(getter)

	expectedFullURL := proto.FullURL{URL: "ozon.ru"}
	shortenURL := proto.ShortURL{URL: "aaaadaaaa"}
	getter.On(getFullURL, shortenURL.URL).Return(expectedFullURL.URL, nil)

	resultFullURL, err := handler.Redirect(context.Background(), &shortenURL)
	assert.NoError(t, err)
	assert.Equal(t, resultFullURL.URL, expectedFullURL.URL)

	assert.True(t, getter.AssertExpectations(t))
}

func TestRedirectErr(t *testing.T) {
	getter := &mockFullUrlGetter{}
	handler := New(getter)

	shortenURL := proto.ShortURL{URL: "aaaadaaaa"}
	getter.On(getFullURL, shortenURL.URL).Return("", errors.New("unknown"))

	_, err := handler.Redirect(context.Background(), &shortenURL)
	assert.Error(t, err)

	assert.True(t, getter.AssertExpectations(t))
}
