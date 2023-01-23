package logging

import (
	"os"
	"path"
	"testing"

	"github.com/jon4hz/subrr/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestEnsureFolderExists(t *testing.T) {
	tmpFolder, err := os.MkdirTemp(os.TempDir(), "subrr-test")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpFolder)

	assert.NotPanics(t, func() {
		ensureFolderExists(tmpFolder)
	})

	assert.NoError(t, os.Remove(tmpFolder))

	assert.NotPanics(t, func() {
		ensureFolderExists(tmpFolder)
	})

	assert.DirExists(t, tmpFolder)

	// create a directory without write permissions and ensure it panics
	assert.NoError(t, os.Chmod(tmpFolder, 0o444))

	assert.Panics(t, func() {
		ensureFolderExists(path.Join(tmpFolder, "panics"))
	})
}

func TestGetUserConfigDir(t *testing.T) {
	dir, err := getUserConfigDir()
	assert.NoError(t, err)
	assert.NotEmpty(t, dir)
}

func TestInit(t *testing.T) {
	tmpDir, err := os.MkdirTemp(os.TempDir(), "subrr-test")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpDir)
	cfg := &config.LoggingConfig{
		Folder: tmpDir,
		Level:  "debug",
	}
	assert.NoError(t, Init(cfg))
	assert.NoError(t, Close())
}

func TestInitUserCfgDir(t *testing.T) {
	cfg := &config.LoggingConfig{
		Folder: "",
		Level:  "info",
	}
	assert.NoError(t, Init(cfg))
	assert.NoError(t, Close())
}

func TestRmIfEmptyRm(t *testing.T) {
	var err error
	logFile, err = os.CreateTemp(os.TempDir(), "subrr-test")
	assert.NoError(t, err)
	defer os.Remove(logFile.Name())

	assert.NoError(t, rmIfEmpty())
	assert.NoDirExists(t, logFile.Name())
}

func TestRmIfEmpty(t *testing.T) {
	var err error
	logFile, err = os.CreateTemp(os.TempDir(), "subrr-test")
	assert.NoError(t, err)
	defer os.Remove(logFile.Name())

	logFile.Write([]byte("test"))

	assert.NoError(t, rmIfEmpty())
}
