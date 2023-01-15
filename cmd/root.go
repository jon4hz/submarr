package cmd

import (
	"fmt"
	"log"
	"time"

	"github.com/jon4hz/subrr/internal/config"
	"github.com/jon4hz/subrr/internal/httpclient"
	"github.com/jon4hz/subrr/internal/version"
	"github.com/jon4hz/subrr/pkg/sonarr"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:     "subrr",
	Short:   "subrr is a tui for sonarr, radarr and lidarr",
	Version: version.Version,
	Run:     root,
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

	viper.BindPFlag("sonarr.host", rootCmd.Flags().Lookup("sonarr-host"))
	viper.BindPFlag("sonarr.api_key", rootCmd.Flags().Lookup("sonarr-api-key"))
	viper.BindPFlag("sonarr.ignore_tls", rootCmd.Flags().Lookup("sonarr-ignore-tls"))
	viper.BindPFlag("sonarr.timeout", rootCmd.Flags().Lookup("sonarr-timeout"))
}

func root(cmd *cobra.Command, args []string) {
	// load the config
	var err error
	cfg, err := config.Load(rootCmdFlags.configFile)
	if err != nil {
		log.Fatalln(err)
	}

	sonarrHTTP := httpclient.New(
		httpclient.WithoutTLSVerfiy(cfg.Sonarr.IgnoreTLS),
		httpclient.WithTimeout(time.Duration(cfg.Sonarr.Timeout*int(time.Second))),
	)

	sonarrClient := sonarr.New(sonarrHTTP, cfg.Sonarr)
	_ = sonarrClient

	fmt.Println("Hello World!")
}
