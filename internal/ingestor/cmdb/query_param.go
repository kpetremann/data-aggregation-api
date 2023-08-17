package cmdb

import (
	"net/url"

	"github.com/criteo/data-aggregation-api/internal/config"
	"github.com/rs/zerolog/log"
)

func deviceDatacenterFilter() url.Values {
	datacenterFilter := ""

	switch string(config.Cfg.NetBox.DatacenterFilterKey) {
	case "site":
		datacenterFilter = "device__site__name"
	case "site_group":
		datacenterFilter = "device__site__group__name"
	case "region":
		datacenterFilter = "device__site__region__name"
	default:
		log.Fatal().Msgf("unknown datacenter filter: %s", config.Cfg.NetBox.DatacenterFilterKey)
	}

	params := url.Values{}
	params.Set(datacenterFilter, config.Cfg.Datacenter)

	return params
}
