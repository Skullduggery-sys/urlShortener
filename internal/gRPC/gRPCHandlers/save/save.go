package save

import (
	"context"
	"fmt"
	"net/url"
	"urlShortener/internal/gRPC/proto"
	"urlShortener/utils/e"
)

type HandleSave struct {
	shortURLGetter
}

type shortURLGetter interface {
	GetShortenURL(fullURL string) (string, error)
}

func New(getter shortURLGetter) *HandleSave {
	return &HandleSave{getter}
}

func (g *HandleSave) Save(ctx context.Context, reqFullURL *proto.FullURL) (*proto.ShortURL, error) {
	const fn = "gRPC.gRPCHAndlers.httpSave.Save"
	fullURL := reqFullURL.URL

	_, err := url.ParseRequestURI(fullURL)
	if err != nil {
		return nil, e.WrapError(fn, err)
	}

	shortenURL, err := g.GetShortenURL(fullURL)
	if err != nil {
		return nil, fmt.Errorf("can't get shorten URL")
	}

	return &proto.ShortURL{URL: shortenURL}, nil
}
