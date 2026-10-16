package cmdb

import (
	"net/url"

	"github.com/criteo/data-aggregation-api/internal/config"
	"github.com/rs/zerolog/log"
)

// deviceDatacenterFilter filters an endpoint exposing a device field on the datacenter.
func deviceDatacenterFilter() url.Values {
	return datacenterFilterWithPrefix("device__")
}

// datacenterFilterWithPrefix filters an endpoint on the datacenter, the device
// being reached through prefix (e.g. "device__"). It allows endpoints that do
// not expose a direct device field to be filtered too, without duplicating the
// datacenter filter key handling. Netbox silently ignores an unknown filter and
// returns everything, so the prefix must match the endpoint.
func datacenterFilterWithPrefix(prefix string) url.Values {
	datacenterFilter := ""

	switch string(config.Cfg.NetBox.DatacenterFilterKey) {
	case "site":
		datacenterFilter = prefix + "site__name"
	case "site_group":
		datacenterFilter = prefix + "site__group__name"
	case "region":
		datacenterFilter = prefix + "site__region__name"
	default:
		log.Fatal().Msgf("unknown datacenter filter: %s", config.Cfg.NetBox.DatacenterFilterKey)
	}

	params := url.Values{}
	params.Set(datacenterFilter, config.Cfg.Datacenter)

	return params
}
