package gRPCServer

import (
	"context"
	"errors"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"log"
	"net"
	"urlShortener/internal/gRPC/gRPCHandlers"
	"urlShortener/internal/gRPC/gRPCHandlers/interceptors"
	"urlShortener/internal/gRPC/proto"
	"urlShortener/utils/e"
)

type Service interface {
	GetShortenURL(fullURL string) (string, error)
	GetFullURL(shortenURL string) (string, error)
}

type GRPCServer struct {
	*grpc.Server
	logger *logrus.Logger
}

func New(logger *logrus.Logger) *GRPCServer {
	return &GRPCServer{
		grpc.NewServer(grpc.UnaryInterceptor(interceptors.LoggerInterceptor(logger))),
		logger,
	}
}

func (g *GRPCServer) Run(ctx context.Context, addr string, service Service) error {
	const fn = "grpc.gRPCServer.Run"
	handlers := gRPCHandlers.New(service)

	proto.RegisterURLShortenerServer(g.Server, handlers)

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return e.WrapError(fn, err)
	}

	go func() {
		err := g.Serve(lis)
		if errors.Is(err, grpc.ErrServerStopped) && err != nil {
			log.Fatalf("%s: %v", fn, err)
		}
	}()
	<-ctx.Done()
	g.GracefulStop()
	return nil
}
