package gRPCHandlers

import (
	"context"
	"urlShortener/internal/gRPC/gRPCHandlers/redirect"
	"urlShortener/internal/gRPC/gRPCHandlers/save"
	"urlShortener/internal/gRPC/proto"
)

type Handlers struct {
	*redirect.HandleRedirect
	*save.HandleSave

	proto.UnimplementedURLShortenerServer
}

type Service interface {
	GetShortenURL(fullURL string) (string, error)
	GetFullURL(shortenURL string) (string, error)
}

func New(service Service) *Handlers {
	return &Handlers{
		HandleRedirect: redirect.New(service),
		HandleSave:     save.New(service),
	}
}

func (h Handlers) Save(ctx context.Context, req *proto.FullURL) (*proto.ShortURL, error) {
	return h.HandleSave.Save(ctx, req)
}

func (h Handlers) Redirect(ctx context.Context, req *proto.ShortURL) (*proto.FullURL, error) {
	return h.HandleRedirect.Redirect(ctx, req)
}
