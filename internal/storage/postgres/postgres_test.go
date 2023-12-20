package postgres

import (
	"context"
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"testing"
	st "urlShortener/internal/storage"
)

func TestMaxIDdbNotEmpty(t *testing.T) {
	db, mock, err := sqlmock.New()

	storage := &Storage{
		db:  db,
		ctx: context.Background(),
	}

	var maxValue uint64 = 90
	mock.ExpectPrepare(`SELECT MAX\(id\) from url`).
		ExpectQuery().WillReturnRows(sqlmock.NewRows([]string{"max"}).AddRow(maxValue))

	result, err := storage.MaxID()
	assert.NoError(t, err)
	assert.Equal(t, result, maxValue)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMaxIDErr(t *testing.T) {
	db, mock, err := sqlmock.New()

	storage := &Storage{
		db:  db,
		ctx: context.Background(),
	}

	rc := "hello"
	mock.ExpectPrepare(`SELECT MAX\(id\) from url`).
		ExpectQuery().WillReturnRows(sqlmock.NewRows([]string{"max"}).AddRow(rc))

	_, err = storage.MaxID()
	assert.Error(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSaveURLSuccess(t *testing.T) {
	db, mock, err := sqlmock.New()

	storage := &Storage{
		db:  db,
		ctx: context.Background(),
	}

	fullURL := "https://ya.ru"
	shortURL := "qewqeqwe"
	mock.ExpectPrepare(`INSERT INTO url\(fullurl, shortenurl\) VALUES \(\$1,\$2\)`).
		ExpectExec().WithArgs(fullURL, shortURL).WillReturnResult(sqlmock.NewResult(1, 1))

	err = storage.SaveURL(fullURL, shortURL)
	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSaveURLRepeatedURL(t *testing.T) {
	db, mock, err := sqlmock.New()

	storage := &Storage{
		db:  db,
		ctx: context.Background(),
	}

	fullURL := "https://ya.ru"
	shortURL := "qewqeqwe"
	mock.ExpectPrepare(`INSERT INTO url\(fullurl, shortenurl\) VALUES \(\$1,\$2\)`).
		ExpectExec().WithArgs(fullURL, shortURL).WillReturnError(&pq.Error{Code: "23505"})

	err = storage.SaveURL(fullURL, shortURL)
	assert.True(t, errors.Is(err, st.ErrURLExists))

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSaveURLUnknownError(t *testing.T) {
	db, mock, err := sqlmock.New()

	storage := &Storage{
		db:  db,
		ctx: context.Background(),
	}

	fullURL := "https://ya.ru"
	shortURL := "qewqeqwe"
	mock.ExpectPrepare(`INSERT INTO url\(fullurl, shortenurl\) VALUES \(\$1,\$2\)`).
		ExpectExec().WithArgs(fullURL, shortURL).WillReturnError(errors.New("unknown error"))

	err = storage.SaveURL(fullURL, shortURL)
	assert.Error(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetFullURLSuccess(t *testing.T) {
	db, mock, err := sqlmock.New()

	storage := &Storage{
		db:  db,
		ctx: context.Background(),
	}

	fullURL := "https://ya.ru"
	shortURL := "qewqeqwe"
	mock.ExpectPrepare(`SELECT fullURL FROM url WHERE shortenurl = \(\$1\)`).
		ExpectQuery().WithArgs(shortURL).WillReturnRows(sqlmock.NewRows([]string{"fullurl"}).AddRow(fullURL))

	resultFullURL, err := storage.GetFullURL(shortURL)
	assert.NoError(t, err)
	assert.Equal(t, resultFullURL, fullURL)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetFullURLNotFound(t *testing.T) {
	db, mock, err := sqlmock.New()

	storage := &Storage{
		db:  db,
		ctx: context.Background(),
	}

	shortURL := "qewqeqwe"
	mock.ExpectPrepare(`SELECT fullURL FROM url WHERE shortenurl = \(\$1\)`).
		ExpectQuery().WithArgs(shortURL).WillReturnError(sql.ErrNoRows)

	_, err = storage.GetFullURL(shortURL)
	assert.True(t, errors.Is(err, st.ErrURLNotFound))

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetFullURLUnexpectedError(t *testing.T) {
	db, mock, err := sqlmock.New()

	storage := &Storage{
		db:  db,
		ctx: context.Background(),
	}

	shortURL := "qewqeqwe"
	mock.ExpectPrepare(`SELECT fullURL FROM url WHERE shortenurl = \(\$1\)`).
		ExpectQuery().WithArgs(shortURL).WillReturnError(errors.New("error"))

	_, err = storage.GetFullURL(shortURL)
	assert.Error(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetShortenURLSuccess(t *testing.T) {
	db, mock, err := sqlmock.New()

	storage := &Storage{
		db:  db,
		ctx: context.Background(),
	}

	fullURL := "qewqeqwe"
	shortURL := "yaaaaaz"
	mock.ExpectPrepare(`SELECT shortenurl FROM url WHERE fullurl = \(\$1\)`).
		ExpectQuery().WithArgs(fullURL).WillReturnRows(sqlmock.NewRows([]string{"shortenurl"}).AddRow(shortURL))

	resultShortURL, err := storage.GetShortenURL(fullURL)
	assert.NoError(t, err)
	assert.Equal(t, resultShortURL, shortURL)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetShortenURLEmpty(t *testing.T) {
	db, mock, err := sqlmock.New()

	storage := &Storage{
		db:  db,
		ctx: context.Background(),
	}

	fullURL := "qewqeqwe"
	mock.ExpectPrepare(`SELECT shortenurl FROM url WHERE fullurl = \(\$1\)`).
		ExpectQuery().WithArgs(fullURL).WillReturnError(sql.ErrNoRows)

	_, err = storage.GetShortenURL(fullURL)
	assert.True(t, errors.Is(err, st.ErrURLNotFound))

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetShortenURLUnknownError(t *testing.T) {
	db, mock, err := sqlmock.New()

	storage := &Storage{
		db:  db,
		ctx: context.Background(),
	}

	fullURL := "qewqeqwe"
	mock.ExpectPrepare(`SELECT shortenurl FROM url WHERE fullurl = \(\$1\)`).
		ExpectQuery().WithArgs(fullURL).WillReturnError(errors.New("unknown"))

	_, err = storage.GetShortenURL(fullURL)
	assert.Error(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}
