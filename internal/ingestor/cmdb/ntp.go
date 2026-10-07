package cmdb

import (
	"fmt"

	"github.com/rs/zerolog/log"

	"github.com/criteo/data-aggregation-api/internal/ingestor/netbox"
	"github.com/criteo/data-aggregation-api/internal/model/cmdb/ntp"
)

// GetNTP returns all NTP configuration from the Network CMDB.
func GetNTP() ([]*ntp.NTP, error) {
	response := netbox.NetboxResponse[ntp.NTP]{}
	params := deviceDatacenterFilter()

	err := netbox.Get("/api/plugins/cmdb/ntp/", &response, params)
	if err != nil {
		return nil, fmt.Errorf("NTP fetching failure: %w", err)
	}

	if len(response.Results) == 0 {
		log.Warn().Msg("no NTP configuration found")
	}

	return response.Results, nil
}

// PrecomputeNTP associates each found NTP configuration to the matching devices.
func PrecomputeNTP(configs []*ntp.NTP) map[string]*ntp.NTP {
	var ntpPerDevice = make(map[string]*ntp.NTP)
	for _, config := range configs {
		ntpPerDevice[config.Device.Name] = config
	}

	return ntpPerDevice
}
