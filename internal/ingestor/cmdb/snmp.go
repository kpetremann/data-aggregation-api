package cmdb

import (
	"fmt"

	"github.com/rs/zerolog/log"

	"github.com/criteo/data-aggregation-api/internal/ingestor/netbox"
	"github.com/criteo/data-aggregation-api/internal/model/cmdb/snmp"
)

// GetSNMP returns all Snmp configuration from the Network CMDB.
func GetSNMP() ([]*snmp.SNMP, error) {
	response := netbox.NetboxResponse[snmp.SNMP]{}
	params := deviceDatacenterFilter()

	err := netbox.Get("/api/plugins/cmdb/snmp/", &response, params)
	if err != nil {
		return nil, fmt.Errorf("SNMP fetching failure: %w", err)
	}

	if len(response.Results) == 0 {
		log.Warn().Msg("no SNMP configuration found")
	}

	return response.Results, nil
}

// PrecomputeSNMP associates each found Snmp configuration to the matching devices.
func PrecomputeSNMP(globalConfigs []*snmp.SNMP) map[string]*snmp.SNMP {
	var SNMPPerDevice = make(map[string]*snmp.SNMP)
	for _, config := range globalConfigs {
		SNMPPerDevice[config.Device.Name] = config
	}

	return SNMPPerDevice
}
