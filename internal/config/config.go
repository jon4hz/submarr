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
	Lidarr  *LidarrConfig  `mapstructure:"lidarr"`
	Logging *LoggingConfig `mapstructure:"logging"`
}

// SonarrConfig represents the sonarr config
type SonarrConfig struct {
	Host      string `mapstructure:"host"`
	APIKey    string `mapstructure:"api_key"`
	IgnoreTLS bool   `mapstructure:"ignore_tls"`
	Timeout   int    `mapstructure:"timeout"`
}

// RadarrConfig represents the radarr config
type RadarrConfig struct {
	Host      string `mapstructure:"host"`
	APIKey    string `mapstructure:"api_key"`
	IgnoreTLS bool   `mapstructure:"ignore_tls"`
	Timeout   int    `mapstructure:"timeout"`
}

// LidarrConfig represents the lidarr config
type LidarrConfig struct {
	Host      string `mapstructure:"host"`
	APIKey    string `mapstructure:"api_key"`
	IgnoreTLS bool   `mapstructure:"ignore_tls"`
	Timeout   int    `mapstructure:"timeout"`
}

// LoggingConfig represents the logging config
type LoggingConfig struct {
	Level  string `mapstructure:"level"`
	Folder string `mapstructure:"folder"`
}

// Load loads the config file.
// It searches in the following locations:
//
// /etc/subrr/config.yml,
// $HOME/.config/subrr/config.yml,
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
		"subrr.yml",
		"subrr.yaml",
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
	viper.AddConfigPath("$HOME/.config/subrr/")
	viper.AddConfigPath("/etc/subrr/")
	if err = viper.ReadInConfig(); err != nil {
		return
	}
	if err = viper.Unmarshal(&cfg); err != nil {
		return
	}
	return
}
