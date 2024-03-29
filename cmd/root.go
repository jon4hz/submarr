package cmd

import (
	"log"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jon4hz/submarr/internal/config"
	"github.com/jon4hz/submarr/internal/core"
	"github.com/jon4hz/submarr/internal/httpclient"
	"github.com/jon4hz/submarr/internal/logging"
	"github.com/jon4hz/submarr/internal/tui"
	"github.com/jon4hz/submarr/internal/version"
	"github.com/jon4hz/submarr/pkg/radarr"
	"github.com/jon4hz/submarr/pkg/sonarr"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:               "submarr",
	Short:             "submarr is a tui for sonarr and radarr",
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

	for _, v := range []string{"sonarr", "radarr"} {
		bindClientFlags(rootCmd, v)
	}

	rootCmd.Flags().String("logging-level", "info", "log level")
	rootCmd.Flags().String("logging-folder", "", "log folder")
	mustBindPFlag("logging.level", rootCmd.Flags().Lookup("logging-level"))
	mustBindPFlag("logging.folder", rootCmd.Flags().Lookup("logging-folder"))

	rootCmd.Flags().Bool("no-mouse", false, "disable mouse support")
	mustBindPFlag("no_mouse", rootCmd.Flags().Lookup("no-mouse"))
}

func bindClientFlags(cmd *cobra.Command, client string) {
	cmd.Flags().String(client+"-host", "", client+" host")
	cmd.Flags().String(client+"-api-key", "", client+" api key")
	cmd.Flags().Bool(client+"-ignore-tls", false, "ignore tls verification")
	cmd.Flags().Int(client+"-timeout", 30, "timeout in seconds")
	mustBindPFlag(client+".host", cmd.Flags().Lookup(client+"-host"))
	mustBindPFlag(client+".api_key", cmd.Flags().Lookup(client+"-api-key"))
	mustBindPFlag(client+".ignore_tls", cmd.Flags().Lookup(client+"-ignore-tls"))
	mustBindPFlag(client+".timeout", cmd.Flags().Lookup(client+"-timeout"))
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
	logging.Log.Debug("starting submarr", "version", version.Version)

	var (
		sonarrClient *sonarr.Client
		radarrClient *radarr.Client
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

	client := core.New(
		cfg,
		sonarrClient,
		radarrClient,
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
