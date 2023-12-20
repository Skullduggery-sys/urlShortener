package gRPCServer

import (
	"context"
	"errors"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net"
	"net/http"
	"sync"
	"testing"
	"time"
)

type mockShortService struct {
	mock.Mock
}

func (m *mockShortService) GetShortenURL(fullURL string) (string, error) {
	args := m.Called(fullURL)
	return args.String(0), args.Error(1)
}

func (m *mockShortService) GetFullURL(shortURL string) (string, error) {
	args := m.Called(shortURL)
	return args.String(0), args.Error(1)
}

func TestRunSuccess(t *testing.T) {
	service := &mockShortService{}
	logger := logrus.New()
	logger.SetLevel(logrus.PanicLevel)
	srv := New(logger)
	testAddr := "localhost:8090"
	ctx, final := context.WithCancel(context.Background())
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		err := srv.Run(ctx, testAddr, service)
		assert.NoError(t, err)
		wg.Done()
	}()

	<-time.After(time.Millisecond * 10)
	final()
	wg.Wait()
}

func TestRunPortInUse(t *testing.T) {
	service := &mockShortService{}
	logger := logrus.New()
	logger.Out = nil
	srv := New(logger)
	// Не знаю что делать, если порт уже занят
	testAddr := "localhost:8090"
	lis, err := net.Listen("tcp", testAddr)
	assert.NoError(t, err)

	wg := sync.WaitGroup{}

	srvHTTP := http.Server{Addr: testAddr}
	wg.Add(1)
	go func() {
		err := srvHTTP.Serve(lis)
		assert.True(t, errors.Is(err, http.ErrServerClosed))
		wg.Done()
	}()

	<-time.After(time.Millisecond * 10)

	err = srv.Run(context.Background(), testAddr, service)
	assert.Error(t, err)

	err = srvHTTP.Close()
	assert.NoError(t, err)
	wg.Wait()
}
