package core

import (
	"github.com/jon4hz/submarr/internal/config"
	coreRadarr "github.com/jon4hz/submarr/internal/core/radarr"
	coreSonarr "github.com/jon4hz/submarr/internal/core/sonarr"
	"github.com/jon4hz/submarr/pkg/radarr"
	"github.com/jon4hz/submarr/pkg/sonarr"
)

type Client struct {
	Sonarr *coreSonarr.Client
	Radarr *coreRadarr.Client
}

func New(cfg *config.Config, sonarr *sonarr.Client, radarr *radarr.Client) *Client {
	return &Client{
		Sonarr: coreSonarr.New(cfg.Sonarr, sonarr),
		Radarr: coreRadarr.New(cfg.Radarr, radarr),
	}
}
