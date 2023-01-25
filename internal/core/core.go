package core

import (
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

func New(sonarr *sonarr.Client, radarr *radarr.Client, lidarr *lidarr.Client) *Client {
	return &Client{
		Sonarr: coreSonarr.New(sonarr),
		Radarr: coreRadarr.New(radarr),
		Lidarr: lidarr,
	}
}
