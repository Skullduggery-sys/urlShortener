package config

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

// FileExists
func TestFileExistsForExistingFile(t *testing.T) {
	tempCfg := createTempFile(t, []byte{})
	defer clearTempFile(t, tempCfg.Name())
	ok := fileExists(tempCfg.Name())
	assert.True(t, ok, "expected file to exist")
}

func TestFileExistsForNonExistingFile(t *testing.T) {
	ok := fileExists("")
	assert.False(t, ok, "expected file to not exist")
}

// readConfig
func createTempFile(t *testing.T, data []byte) *os.File {
	tempFile, err := os.CreateTemp("", "temp-config-*.yaml")
	assert.NoError(t, err, "")
	defer func() {
		err = tempFile.Close()
		assert.NoError(t, err, "error closing file")
	}()
	length, err := tempFile.Write(data)
	assert.Equal(t, length, len(data), "error written len lower than actual")
	assert.NoError(t, err, "error writing file")
	return tempFile
}

func clearTempFile(t *testing.T, path string) {
	assert.NoError(t, os.Remove(path), "error while removing temp file")
}

func TestReadConfigFileNotFound(t *testing.T) {
	cfg, err := readConfig("_")
	assert.Error(t, err, "expected error")
	assert.Nil(t, cfg, "cfg should be nil")
}

func TestReadConfigFileUnmarshalError(t *testing.T) {
	tempCfg := createTempFile(t, []byte("kjsssndknfs"))
	defer clearTempFile(t, tempCfg.Name())

	cfg, err := readConfig(tempCfg.Name())
	assert.Error(t, err)
	assert.Nil(t, cfg)
}

func TestReadConfigFileSuccess(t *testing.T) {
	tempCfg := createTempFile(t, []byte("postgres:\n  login: \"postgres\"\n  "+
		"password: \"dfsdf\"\n  host: \"localhost\"\n  port: \"5432\"\n  dbname:"+
		" \"urlshortener\"\n  sslMode: \"disable\"\nhttpServer:\n  "+
		"address: \"localhost:8081\"\n  timeout: 4s\n  "+
		"idleTimeout: 60s\nenv: \"local\"\ngrpcAddr: \"127.0.0.1:8082\""))
	defer clearTempFile(t, tempCfg.Name())

	cfg, err := readConfig(tempCfg.Name())
	assert.NoError(t, err)
	assert.NotNil(t, cfg)
}

//mustParseConfig

func TestMustParseConfigFileNotExists(t *testing.T) {
	cfg, err := MustParseConfig("_")
	assert.Error(t, err)
	assert.Nil(t, cfg)
}

func TestMustParseConfigInvalidJSON(t *testing.T) {
	tempFile := createTempFile(t, []byte("sdfs"))
	defer clearTempFile(t, tempFile.Name())

	cfg, err := MustParseConfig(tempFile.Name())
	assert.Error(t, err)
	assert.Nil(t, cfg)
}

func TestMustParseConfigValidateErrorRequired(t *testing.T) {
	tempFile := createTempFile(t, []byte("postgres:\n  login: \"postgres\"\n  "+
		"password: \"dfsdf\"\n  host: \"localhost\"\n"))
	defer clearTempFile(t, tempFile.Name())

	cfg, err := MustParseConfig(tempFile.Name())
	assert.Error(t, err)
	assert.Nil(t, cfg)
}

func TestMustParseConfigValidateErrorNumeric(t *testing.T) {
	tempFile := createTempFile(t, []byte("postgres:\n  login: \"postgres\"\n  "+
		"password: \"123123\"\n  host: \"localhost\"\n  port: \"dfdfd\"\n  dbname:"+
		" \"urlshortener\"\n  sslMode: \"disable\"\nhttpServer:\n  "+
		"address: \"localhost:8081\"\n  timeout: 4s\n  "+
		"idleTimeout: 60s\nenv: \"local\"\ngrpcAddr: \"127.0.0.1:8082\""))
	defer clearTempFile(t, tempFile.Name())

	cfg, err := MustParseConfig(tempFile.Name())
	assert.Error(t, err)
	assert.Nil(t, cfg)
}

func TestMustParseConfigSuccess(t *testing.T) {
	tempCfg := createTempFile(t, []byte("postgres:\n  login: \"postgres\"\n  "+
		"password: \"123123\"\n  host: \"localhost\"\n  port: \"5432\"\n  dbname:"+
		" \"urlshortener\"\n  sslMode: \"disable\"\nhttpServer:\n  "+
		"address: \"localhost:8081\"\n  timeout: 4s\n  "+
		"idleTimeout: 60s\nenv: \"local\"\ngrpcAddr: \"127.0.0.1:8082\""))
	defer clearTempFile(t, tempCfg.Name())

	cfg, err := MustParseConfig(tempCfg.Name())
	assert.NoError(t, err)
	assert.NotNil(t, cfg)

	absoluteCfg := Config{
		Postgres: PostgresConfig{
			Login:    "postgres",
			Password: "123123",
			Host:     "localhost",
			Port:     "5432",
			DBName:   "urlshortener",
			SSLMode:  "disable",
		},
		HTTPServer: HTTPServerConfig{
			Address:     "localhost:8081",
			Timeout:     4 * time.Second,
			IdleTimeout: time.Minute,
		},
		GRPCAddr: "127.0.0.1:8082",
	}

	assert.Equal(t, absoluteCfg, *cfg)
}
