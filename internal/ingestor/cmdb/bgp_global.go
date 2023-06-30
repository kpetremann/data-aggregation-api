package cmdb

import (
	"fmt"

	"github.com/rs/zerolog/log"

	"github.com/criteo/data-aggregation-api/internal/ingestor/netbox"
	"github.com/criteo/data-aggregation-api/internal/model/cmdb/bgp"
)

// GetBGPGlobal returns all BGP global configuration from the Network CMDB.
func GetBGPGlobal() ([]*bgp.BGPGlobal, error) {
	response := netbox.NetboxResponse[bgp.BGPGlobal]{}

	err := netbox.Get("/api/plugins/cmdb/bgp-global/", &response)
	if err != nil {
		return nil, fmt.Errorf("BGP Global fetching failure: %w", err)
	}

	if response.Count != len(response.Results) {
		log.Warn().Msg("no BGP global configuration found")
	}

	return response.Results, nil
}

// PrecomputeBGPGlobal associates each found BGP global configuration to the matching devices.
func PrecomputeBGPGlobal(globalConfigs []*bgp.BGPGlobal) map[string]*bgp.BGPGlobal {
	var bgpGlobalPerDevice = make(map[string]*bgp.BGPGlobal)
	for _, config := range globalConfigs {
		bgpGlobalPerDevice[config.Device.Name] = config
	}

	return bgpGlobalPerDevice
}
