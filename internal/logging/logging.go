package logging

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/charmbracelet/log"
	"github.com/jon4hz/submarr/internal/config"
)

var (
	Log         *log.Logger
	logFile     *os.File
	initialized bool
)

// Init initializes the logger
func Init(cfg *config.LoggingConfig) error {
	// make sure we only initialize the logger once
	if initialized {
		return nil
	}

	initialized = true
	err := initLogger(cfg.Folder, strings.ToLower(cfg.Level))
	if err != nil {
		initialized = false
	}
	return err
}

// Close closes the log file, sets initialized to false and resets the logger.
func Close() error {
	initialized = false
	Log = log.New(io.Discard)

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
	// and create a subfolder for submarr and its logs
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

	l := log.New(logFile)
	l.SetReportCaller(true)

	// set the log level
	switch level {
	case "debug":
		l.SetLevel(log.DebugLevel)
	case "warn":
		l.SetLevel(log.WarnLevel)
	case "error":
		l.SetLevel(log.ErrorLevel)
	default:
		l.SetLevel(log.InfoLevel)
	}

	// set the global logger
	Log = l
	return nil
}

func getUserConfigDir() (string, error) {
	folder, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(folder, "submarr", "logs"), nil
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
