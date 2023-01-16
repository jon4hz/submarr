package core

import (
	"github.com/jon4hz/subrr/pkg/lidarr"
	"github.com/jon4hz/subrr/pkg/radarr"
	"github.com/jon4hz/subrr/pkg/sonarr"
)

type Client struct {
	Sonarr *sonarr.Client
	Radarr *radarr.Client
	Lidarr *lidarr.Client
}

func New(sonarr *sonarr.Client, radarr *radarr.Client, lidarr *lidarr.Client) *Client {
	return &Client{
		Sonarr: sonarr,
		Radarr: radarr,
		Lidarr: lidarr,
	}
}
