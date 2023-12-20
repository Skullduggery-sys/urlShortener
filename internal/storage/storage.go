package storage

import "errors"

var (
	ErrURLExists   = errors.New("URL already exists")
	ErrURLNotFound = errors.New("URL not found")
)

type Storager interface {
	SaveURL(urlToSave string, shortenUrl string) error
	GetFullURL(shortenURL string) (string, error)
	GetShortenURL(fullURL string) (string, error)
}
