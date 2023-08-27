package config

import (
	"errors"
	"os"

	"github.com/spf13/viper"
)

// Config represents the config
type Config struct {
	Sonarr  *SonarrConfig  `mapstructure:"sonarr"`
	Radarr  *RadarrConfig  `mapstructure:"radarr"`
	Logging *LoggingConfig `mapstructure:"logging"`
	NoMouse bool           `mapstructure:"no_mouse"`
}

type ClientConfig struct {
	Host          string           `mapstructure:"host"`
	APIKey        string           `mapstructure:"api_key"`
	IgnoreTLS     bool             `mapstructure:"ignore_tls"`
	Timeout       int              `mapstructure:"timeout"`
	BasicAuth     *BasicAuthConfig `mapstructure:"basic_auth"`
	HeaderConfigs []HeaderConfig   `mapstructure:"headers"`
}

// SonarrConfig represents the sonarr config
type SonarrConfig struct {
	ClientConfig           `mapstructure:",squash"`
	DefaultQualityProfile  string `mapstructure:"default_quality_profile"`
	DefaultLanguageProfile string `mapstructure:"default_language_profile"`
}

// RadarrConfig represents the radarr config
type RadarrConfig struct {
	ClientConfig `mapstructure:",squash"`
}

// LoggingConfig represents the logging config
type LoggingConfig struct {
	Level  string `mapstructure:"level"`
	Folder string `mapstructure:"folder"`
}

// BasicAuthConfig represents the config for basic authentication
type BasicAuthConfig struct {
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

// HeaderConfig represents an abitrary http header
type HeaderConfig struct {
	Key   string `mapstructure:"key"`
	Value string `mapstructure:"value"`
}

// Load loads the config file.
// It searches in the following locations:
//
// /etc/submarr/config.yml,
// $HOME/.config/submarr/config.yml,
// config.yml
//
// command arguments will overwrite the value from the config
func Load(path string) (cfg *Config, err error) {
	if path != "" {
		return load(path)
	}
	for _, f := range [...]string{
		".config.yml",
		"config.yml",
		".config.yaml",
		"config.yaml",
		"submarr.yml",
		"submarr.yaml",
	} {
		cfg, err = load(f)
		if err != nil && os.IsNotExist(err) {
			err = nil
			continue
		} else if err != nil && errors.As(err, &viper.ConfigFileNotFoundError{}) {
			err = nil
			continue
		}
	}
	if cfg == nil {
		return cfg, viper.Unmarshal(&cfg)
	}
	return
}

func load(file string) (cfg *Config, err error) {
	viper.SetConfigName(file)
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./")
	viper.AddConfigPath("$HOME/.config/submarr/")
	viper.AddConfigPath("/etc/submarr/")
	if err = viper.ReadInConfig(); err != nil {
		return
	}
	if err = viper.Unmarshal(&cfg); err != nil {
		return
	}
	return
}
