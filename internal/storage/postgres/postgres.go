package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"urlShortener/internal/config"
	"urlShortener/internal/storage"
	"urlShortener/utils/e"
)

type Storage struct {
	db     *sql.DB
	ctx    context.Context
	logger *logrus.Logger
}

func New(ctx context.Context, cfg *config.PostgresConfig, logger *logrus.Logger) (*Storage, error) {
	const fn = "storage.postgres.New"

	db, err := sql.Open("postgres", connStr(cfg))
	if err != nil {
		return nil, e.WrapError(fn, err)
	}

	err = db.PingContext(ctx)
	if err != nil {
		return nil, e.WrapError(fn, err)
	}

	err = createTable(ctx, db, logger)
	if err != nil {
		return nil, e.WrapError(fn, err)
	}

	res := &Storage{
		db:     db,
		ctx:    ctx,
		logger: logger,
	}

	return res, nil
}

func createTable(ctx context.Context, db *sql.DB, logger *logrus.Logger) error {
	const fn = "storage.postgres.createTable"

	query, err := db.Prepare(
		`CREATE TABLE IF NOT EXISTS url (
    id SERIAL NOT NULL PRIMARY KEY,
    fullURL TEXT NOT NULL UNIQUE ,
    shortenURL TEXT NOT NULL UNIQUE );`)
	if err != nil {
		return e.WrapError(fn, err)
	}
	defer func() {
		err = query.Close()
		if err != nil {
			logger.Errorf("%s: can't close query %v", fn, err)
		}
	}()

	_, err = query.ExecContext(ctx)
	if err != nil {
		return e.WrapError(fn, err)
	}

	return nil
}

func (s *Storage) MaxID() (uint64, error) {
	const fn = "storage.postres.lastID"

	query, err := s.db.Prepare(`SELECT MAX(id) from url`)
	if err != nil {
		return 0, e.WrapError(fn, err)
	}
	defer func() {
		err = query.Close()
		if err != nil {
			s.logger.Errorf("%s: can't close query %v", fn, err)
		}
	}()

	// придется вручную конвертировать в uint64, но тут ничего страшного, так как максимальное число сокращенных ссылок
	// 63 ** 10 - 1 все равно меньше верхней границы int64, переполнения не будет
	var maxID sql.NullInt64
	if err = query.QueryRowContext(s.ctx).Scan(&maxID); err != nil {
		return 0, e.WrapError(fn, err)
	}

	if maxID.Valid {
		return uint64(maxID.Int64), nil
	} else {
		return 0, nil
	}
}

func connStr(cfg *config.PostgresConfig) string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Login, cfg.Password, cfg.DBName, cfg.SSLMode)
}

func (s *Storage) SaveURL(urlToSave string, shortenUrl string) error {
	const fn = "storage.postgres.SaveURL"

	query, err := s.db.Prepare(`INSERT INTO url(fullurl, shortenurl) VALUES ($1,$2)`)
	if err != nil {
		return e.WrapError(fn, err)
	}
	defer func() {
		err = query.Close()
		if err != nil {
			s.logger.Errorf("%s: can't close query %v", fn, err)
		}
	}()

	_, err = query.ExecContext(s.ctx, urlToSave, shortenUrl)
	if err != nil {
		if pqError, ok := err.(*pq.Error); ok {
			switch pqError.Code.Name() {
			case "unique_violation":
				return e.WrapError(fn, storage.ErrURLExists)
			}
		}
		return e.WrapError(fn, err)
	}

	return nil
}

func (s *Storage) GetFullURL(shortenURL string) (string, error) {
	const fn = "storage.postgres.GetURL"

	query, err := s.db.Prepare(`SELECT fullURL FROM url WHERE shortenurl = ($1)`)
	if err != nil {
		return "", e.WrapError(fn, err)
	}
	defer func() {
		err = query.Close()
		if err != nil {
			s.logger.Errorf("%s: can't close query %v", fn, err)
		}
	}()

	var fullURL string
	err = query.QueryRowContext(s.ctx, shortenURL).Scan(&fullURL)
	if errors.Is(err, sql.ErrNoRows) {
		return "", storage.ErrURLNotFound
	} else if err != nil {
		return "", e.WrapError(fn, err)
	}

	return fullURL, nil
}

func (s *Storage) GetShortenURL(fullURL string) (string, error) {
	const fn = "storage.postgres.GetURL"

	query, err := s.db.Prepare(`SELECT shortenurl FROM url WHERE fullurl = ($1)`)
	if err != nil {
		return "", e.WrapError(fn, err)
	}
	defer func() {
		err = query.Close()
		if err != nil {
			s.logger.Errorf("%s: can't close query %v", fn, err)
		}
	}()

	var shortenURL string
	err = query.QueryRowContext(s.ctx, fullURL).Scan(&shortenURL)
	if errors.Is(err, sql.ErrNoRows) {
		return "", storage.ErrURLNotFound
	} else if err != nil {
		return "", e.WrapError(fn, err)
	}

	return shortenURL, nil
}
