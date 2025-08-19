package dcim

import (
	"fmt"
	"net/url"

	"github.com/rs/zerolog/log"

	"github.com/criteo/data-aggregation-api/internal/config"
	"github.com/criteo/data-aggregation-api/internal/ingestor/netbox"
	"github.com/criteo/data-aggregation-api/internal/model/dcim"
)

// GetNetworkInventory returns network device inventory from NetBox DCIM.
func GetNetworkInventory() ([]*dcim.NetworkDevice, error) {
	response := netbox.NetboxResponse[dcim.NetworkDevice]{}

	params := url.Values{}
	params.Set(string(config.Cfg.NetBox.DatacenterFilterKey), config.Cfg.Datacenter)
	for _, filter := range config.Cfg.NetBox.DeviceFilters {
		params.Add(filter.Filter, filter.Value)
	}

	if err := netbox.Get("/api/dcim/devices/", &response, params); err != nil {
		return nil, fmt.Errorf("network inventory fetching failure: %w", err)
	}

	if response.Count != len(response.Results) {
		log.Warn().Msg("some devices have not been fetched")
	}

	return response.Results, nil
}
