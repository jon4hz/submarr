package core

import (
	"github.com/jon4hz/subrr/internal/config"
	coreRadarr "github.com/jon4hz/subrr/internal/core/radarr"
	coreSonarr "github.com/jon4hz/subrr/internal/core/sonarr"
	"github.com/jon4hz/subrr/pkg/lidarr"
	"github.com/jon4hz/subrr/pkg/radarr"
	"github.com/jon4hz/subrr/pkg/sonarr"
)

type Client struct {
	Sonarr *coreSonarr.Client
	Radarr *coreRadarr.Client
	Lidarr *lidarr.Client
}

func New(cfg *config.Config, sonarr *sonarr.Client, radarr *radarr.Client, lidarr *lidarr.Client) *Client {
	return &Client{
		Sonarr: coreSonarr.New(cfg.Sonarr, sonarr),
		Radarr: coreRadarr.New(cfg.Radarr, radarr),
		Lidarr: lidarr,
	}
}
