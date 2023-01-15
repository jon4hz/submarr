package config_test

import (
	"testing"

	"github.com/jon4hz/subrr/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	cfg, err := config.Load("testdata/config.yml")
	assert.NoError(t, err)
	assert.NotNil(t, cfg)

	assert.Equal(t, "https://sonarr.local/", cfg.Sonarr.Host)
	assert.Equal(t, "123456a", cfg.Sonarr.APIKey)
}

func TestLoadNoFileConfig(t *testing.T) {
	cfg, err := config.Load("")
	assert.NoError(t, err)
	assert.NotNil(t, cfg)
}

func TestLoadInvalidFile(t *testing.T) {
	_, err := config.Load("testdata/invalid.txt")
	assert.Error(t, err)
}
