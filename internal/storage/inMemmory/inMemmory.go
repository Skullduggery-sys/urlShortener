package inMemmory

import (
	"sync"
	"urlShortener/internal/storage"
)

type Storage struct {
	mu            sync.RWMutex
	keyShortenURL map[string]string
	keyFullURL    map[string]string
}

func New() *Storage {
	return &Storage{
		mu:            sync.RWMutex{},
		keyShortenURL: make(map[string]string),
		keyFullURL:    make(map[string]string),
	}
}

func (s *Storage) GetFullURL(shortenURL string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	fullURL, ok := s.keyShortenURL[shortenURL]
	if ok {
		return fullURL, nil
	} else {
		return "", storage.ErrURLNotFound
	}
}

func (s *Storage) GetShortenURL(fullURL string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	shortenURL, ok := s.keyFullURL[fullURL]
	if ok {
		return shortenURL, nil
	} else {
		return "", storage.ErrURLNotFound
	}
}

func (s *Storage) SaveURL(fullURL string, shortenURL string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.keyFullURL[fullURL]; ok {
		return storage.ErrURLExists
	}

	s.keyFullURL[fullURL] = shortenURL
	s.keyShortenURL[shortenURL] = fullURL
	return nil
}
