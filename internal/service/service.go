package service

import (
	"errors"
	"urlShortener/internal/lib/linkShortening"
	"urlShortener/internal/storage"
	"urlShortener/utils/e"
)

type Service struct {
	storage.Storager
	linkShortening.Hasher
}

func New(storage storage.Storager, hasher linkShortening.Hasher) *Service {
	return &Service{
		Storager: storage,
		Hasher:   hasher,
	}
}

func (s *Service) GetShortenURL(fullURL string) (string, error) {
	const fn = "service.GetShortenURL"

	shortenURL, err := s.Storager.GetShortenURL(fullURL)
	if err == nil {
		return shortenURL, nil
	} else if !errors.Is(err, storage.ErrURLNotFound) {
		return "", e.WrapError(fn, err)
	}

	shortenURL, err = s.Hash()
	if err != nil {
		return "", e.WrapError(fn, err)
	}
	err = s.SaveURL(fullURL, shortenURL)
	if err != nil {
		return "", e.WrapError(fn, err)
	}

	return shortenURL, nil
}

func (s *Service) GetFullURL(shortenURL string) (string, error) {
	const fn = "service.GetFullURL"

	fullURL, err := s.Storager.GetFullURL(shortenURL)
	if err != nil {
		return "", e.WrapError(fn, err)
	}

	return fullURL, nil
}
