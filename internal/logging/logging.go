package logging

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/jon4hz/subrr/internal/config"
	"github.com/rs/zerolog"
)

var (
	Log            *zerolog.Logger
	logFile        *os.File
	logDir         string
	initialized    bool
	initializeOnce sync.Once = sync.Once{}
)

// Init initializes the logger
func Init(cfg *config.LoggingConfig) error {
	var err error
	// make sure we only initialize the logger once
	initializeOnce.Do(func() {
		err = initLogger(cfg.Folder, strings.ToLower(cfg.Level))
		initialized = true
	})

	return err
}

// Close closes the log file, sets initialized to false and resets the logger.
func Close() error {
	initialized = false
	Log = &zerolog.Logger{}
	initializeOnce = sync.Once{}

	if err := rmIfEmpty(); err != nil {
		return err
	}

	return logFile.Close()
}

// rmIfEmpty removes the logFile if it is empty
func rmIfEmpty() error {
	if err := logFile.Sync(); err != nil {
		return err
	}
	stats, err := logFile.Stat()
	if err != nil {
		return err
	}
	if stats.Size() == 0 {
		return os.Remove(logFile.Name())
	}
	return nil
}

// openFile opens the log file
// if the folder does not exist, it will be created
func openFile(folder string) error {
	ensureFolderExists(folder)
	f := filepath.Join(folder, fmt.Sprintf("%s.log", time.Now().UTC().Format("20060102150405")))
	var err error
	logFile, err = os.OpenFile(f, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0o644)
	if err != nil {
		return err
	}
	return nil
}

func initLogger(folder, level string) error {
	// if folder is empty, try to get the user config dir
	// and create a subfolder for subrr and its logs
	if folder == "" {
		var err error
		folder, err = getUserConfigDir()
		if err != nil {
			return err
		}
	}

	// open the log file
	if err := openFile(folder); err != nil {
		return err
	}

	// create a new zerolog logger
	zerolog.TimeFieldFormat = time.RFC3339
	zl := zerolog.New(logFile)
	zl = zl.With().Caller().Timestamp().Logger()

	// set the log level
	switch level {
	case "debug":
		zl = zl.Level(zerolog.DebugLevel)
	case "info":
		zl = zl.Level(zerolog.InfoLevel)
	case "warn":
		zl = zl.Level(zerolog.WarnLevel)
	case "error":
		zl = zl.Level(zerolog.ErrorLevel)
	default:
		zl = zl.Level(zerolog.InfoLevel)
	}

	// set the global logger
	Log = &zl
	return nil
}

func getUserConfigDir() (string, error) {
	folder, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(folder, "subrr", "logs"), nil
}

func ensureFolderExists(folder string) {
	if _, err := os.Stat(folder); os.IsNotExist(err) {
		if err := os.MkdirAll(folder, 0o755); err != nil {
			panic(err)
		}
	} else if err != nil {
		panic(err)
	}
}
