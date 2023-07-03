package cmd

import (
	"log"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jon4hz/subrr/internal/config"
	"github.com/jon4hz/subrr/internal/core"
	"github.com/jon4hz/subrr/internal/httpclient"
	"github.com/jon4hz/subrr/internal/logging"
	"github.com/jon4hz/subrr/internal/tui"
	"github.com/jon4hz/subrr/internal/version"
	"github.com/jon4hz/subrr/pkg/lidarr"
	"github.com/jon4hz/subrr/pkg/radarr"
	"github.com/jon4hz/subrr/pkg/sonarr"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:               "subrr",
	Short:             "subrr is a tui for sonarr, radarr and lidarr",
	Version:           version.Version,
	Run:               root,
	CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
}

var rootCmdFlags struct {
	configFile string
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(versionCmd)

	rootCmd.Flags().StringVarP(&rootCmdFlags.configFile, "config", "c", "", "path to the config file")

	rootCmd.Flags().String("sonarr-host", "", "sonarr host")
	rootCmd.Flags().String("sonarr-api-key", "", "sonarr api key")
	rootCmd.Flags().Bool("sonarr-ignore-tls", false, "ignore tls verification")
	rootCmd.Flags().Int("sonarr-timeout", 30, "timeout in seconds")
	mustBindPFlag("sonarr.host", rootCmd.Flags().Lookup("sonarr-host"))
	mustBindPFlag("sonarr.api_key", rootCmd.Flags().Lookup("sonarr-api-key"))
	mustBindPFlag("sonarr.ignore_tls", rootCmd.Flags().Lookup("sonarr-ignore-tls"))
	mustBindPFlag("sonarr.timeout", rootCmd.Flags().Lookup("sonarr-timeout"))

	rootCmd.Flags().String("radarr-host", "", "radarr host")
	rootCmd.Flags().String("radarr-api-key", "", "radarr api key")
	rootCmd.Flags().Bool("radarr-ignore-tls", false, "ignore tls verification")
	rootCmd.Flags().Int("radarr-timeout", 30, "timeout in seconds")
	mustBindPFlag("radarr.host", rootCmd.Flags().Lookup("radarr-host"))
	mustBindPFlag("radarr.api_key", rootCmd.Flags().Lookup("radarr-api-key"))
	mustBindPFlag("radarr.ignore_tls", rootCmd.Flags().Lookup("radarr-ignore-tls"))
	mustBindPFlag("radarr.timeout", rootCmd.Flags().Lookup("radarr-timeout"))

	rootCmd.Flags().String("lidarr-host", "", "lidarr host")
	rootCmd.Flags().String("lidarr-api-key", "", "lidarr api key")
	rootCmd.Flags().Bool("lidarr-ignore-tls", false, "ignore tls verification")
	rootCmd.Flags().Int("lidarr-timeout", 30, "timeout in seconds")
	mustBindPFlag("lidarr.host", rootCmd.Flags().Lookup("lidarr-host"))
	mustBindPFlag("lidarr.api_key", rootCmd.Flags().Lookup("lidarr-api-key"))
	mustBindPFlag("lidarr.ignore_tls", rootCmd.Flags().Lookup("lidarr-ignore-tls"))
	mustBindPFlag("lidarr.timeout", rootCmd.Flags().Lookup("lidarr-timeout"))

	rootCmd.Flags().String("logging-level", "info", "log level")
	rootCmd.Flags().String("logging-folder", "", "log folder")
	mustBindPFlag("logging.level", rootCmd.Flags().Lookup("logging-level"))
	mustBindPFlag("logging.folder", rootCmd.Flags().Lookup("logging-folder"))

	rootCmd.Flags().Bool("no-mouse", false, "disable mouse support")
	mustBindPFlag("no_mouse", rootCmd.Flags().Lookup("no-mouse"))
}

func mustBindPFlag(key string, flag *pflag.Flag) {
	if err := viper.BindPFlag(key, flag); err != nil {
		log.Fatalf("unable to bind flag %q: %v", key, err)
	}
}

func root(cmd *cobra.Command, args []string) {
	// load the config
	var err error
	cfg, err := config.Load(rootCmdFlags.configFile)
	if err != nil {
		log.Fatalln(err)
	}

	// init the logger
	if err := logging.Init(cfg.Logging); err != nil {
		log.Fatalln(err)
	}
	defer func() {
		if err := logging.Close(); err != nil {
			log.Fatalln(err)
		}
	}()
	logging.Log.Debug().Str("version", version.Version).Msg("starting subrr")

	var (
		sonarrClient *sonarr.Client
		radarrClient *radarr.Client
		lidarrClient *lidarr.Client
	)

	if cfg.Sonarr.Host != "" {
		opts := []httpclient.ClientOpts{
			httpclient.WithAPIKey(cfg.Sonarr.APIKey),
			httpclient.WithoutTLSVerfiy(cfg.Sonarr.IgnoreTLS),
			httpclient.WithTimeout(time.Duration(cfg.Sonarr.Timeout * int(time.Second))),
		}
		if cfg.Sonarr.BasicAuth != nil {
			opts = append(opts, httpclient.WithBasicAuth(cfg.Sonarr.BasicAuth.Username, cfg.Sonarr.BasicAuth.Password))
		}
		for _, v := range cfg.Sonarr.HeaderConfigs {
			opts = append(opts, httpclient.WithHeader(v.Key, v.Value))
		}
		sonarrHTTP := httpclient.New(opts...)
		sonarrClient = sonarr.New(sonarrHTTP, cfg.Sonarr)
	}

	if cfg.Radarr.Host != "" {
		opts := []httpclient.ClientOpts{
			httpclient.WithAPIKey(cfg.Radarr.APIKey),
			httpclient.WithoutTLSVerfiy(cfg.Radarr.IgnoreTLS),
			httpclient.WithTimeout(time.Duration(cfg.Radarr.Timeout * int(time.Second))),
		}
		if cfg.Radarr.BasicAuth != nil {
			opts = append(opts, httpclient.WithBasicAuth(cfg.Radarr.BasicAuth.Username, cfg.Radarr.BasicAuth.Password))
		}
		for _, v := range cfg.Radarr.HeaderConfigs {
			opts = append(opts, httpclient.WithHeader(v.Key, v.Value))
		}
		radarrHTTP := httpclient.New(opts...)
		radarrClient = radarr.New(radarrHTTP, cfg.Radarr)
	}

	if cfg.Lidarr.Host != "" {
		opts := []httpclient.ClientOpts{
			httpclient.WithAPIKey(cfg.Lidarr.APIKey),
			httpclient.WithoutTLSVerfiy(cfg.Lidarr.IgnoreTLS),
			httpclient.WithTimeout(time.Duration(cfg.Lidarr.Timeout * int(time.Second))),
		}
		if cfg.Lidarr.BasicAuth != nil {
			opts = append(opts, httpclient.WithBasicAuth(cfg.Lidarr.BasicAuth.Username, cfg.Lidarr.BasicAuth.Password))
		}
		for _, v := range cfg.Lidarr.HeaderConfigs {
			opts = append(opts, httpclient.WithHeader(v.Key, v.Value))
		}
		lidarrHTTP := httpclient.New(opts...)
		lidarrClient = lidarr.New(lidarrHTTP, cfg.Lidarr)
	}

	client := core.New(
		sonarrClient,
		radarrClient,
		lidarrClient,
	)

	tui := tui.New(client)
	opts := []tea.ProgramOption{
		tea.WithAltScreen(),
	}
	if !cfg.NoMouse {
		opts = append(opts, tea.WithMouseCellMotion())
	}
	if err := tui.Run(opts...); err != nil {
		log.Fatalln(err)
	}
}
