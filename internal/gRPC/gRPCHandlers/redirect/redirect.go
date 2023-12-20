package redirect

import (
	"context"
	"urlShortener/internal/gRPC/proto"
	"urlShortener/utils/e"
)

type HandleRedirect struct {
	fullURLGetter
}

type fullURLGetter interface {
	GetFullURL(shortURL string) (string, error)
}

func New(getter fullURLGetter) *HandleRedirect {
	return &HandleRedirect{getter}
}

func (g *HandleRedirect) Redirect(ctx context.Context, reqShortenURL *proto.ShortURL) (*proto.FullURL, error) {
	const fn = "gRPC.gRPCHandlers.httpRedirect.Redirect"
	shortenURL := reqShortenURL.URL

	fullURL, err := g.GetFullURL(shortenURL)
	if err != nil {
		return nil, e.WrapError(fn, err)
	}

	return &proto.FullURL{URL: fullURL}, nil
}
