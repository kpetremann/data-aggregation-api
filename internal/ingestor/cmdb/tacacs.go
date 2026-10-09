package cmdb

import (
	"fmt"

	"github.com/rs/zerolog/log"

	"github.com/criteo/data-aggregation-api/internal/ingestor/netbox"
	"github.com/criteo/data-aggregation-api/internal/model/cmdb/tacacs"
)

// GetTacacs returns all TACACS configurations from the Network CMDB.
func GetTacacs() ([]*tacacs.Tacacs, error) {
	response := netbox.NetboxResponse[tacacs.Tacacs]{}
	params := deviceDatacenterFilter()

	err := netbox.Get("/api/plugins/cmdb/tacacs/", &response, params)
	if err != nil {
		return nil, fmt.Errorf("TACACS fetching failure: %w", err)
	}

	if len(response.Results) == 0 {
		log.Warn().Msg("no TACACS configuration found")
	}

	return response.Results, nil
}

// PrecomputeTacacs associates each found TACACS configuration to the matching device.
func PrecomputeTacacs(globalConfigs []*tacacs.Tacacs) map[string]*tacacs.Tacacs {
	var tacacsPerDevice = make(map[string]*tacacs.Tacacs)
	for _, config := range globalConfigs {
		tacacsPerDevice[config.Device.Name] = config
	}

	return tacacsPerDevice
}
