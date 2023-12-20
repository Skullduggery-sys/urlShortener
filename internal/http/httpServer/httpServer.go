package httpServer

import (
	"context"
	"errors"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
	"sync"
	"time"
	"urlShortener/internal/config"
)

type Server struct {
	srv    *http.Server
	ctx    context.Context
	logger *logrus.Logger
}

type ServerOption func(*Server)

func New(ctx context.Context, cfg config.HTTPServerConfig, router *mux.Router, logger *logrus.Logger) *Server {
	server := &Server{
		srv: &http.Server{
			Addr:         cfg.Address,
			Handler:      router,
			ReadTimeout:  cfg.Timeout,
			WriteTimeout: cfg.Timeout,
			IdleTimeout:  cfg.IdleTimeout,
		},
		ctx:    ctx,
		logger: logger,
	}

	return server
}

func (s *Server) Run() {

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		if err := s.srv.ListenAndServe(); err != nil && !errors.Is(http.ErrServerClosed, err) {
			if s.logger != nil {
				s.logger.Fatal(err)
			} else {
				panic(err)
			}
		}
	}()

	<-s.ctx.Done()

	ctx, final := context.WithTimeout(context.Background(), time.Second)
	defer final()

	err := s.srv.Shutdown(ctx)
	if err != nil {
		s.logger.Error("can't stop http Server: ", err.Error())
	}

	wg.Wait()
}
