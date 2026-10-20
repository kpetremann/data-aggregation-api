package cmdb

import (
	"fmt"

	"github.com/rs/zerolog/log"

	"github.com/criteo/data-aggregation-api/internal/ingestor/netbox"
	"github.com/criteo/data-aggregation-api/internal/model/cmdb/dhcprelay"
)

// GetDHCPRelays returns all DHCP relay interfaces from the Network CMDB.
func GetDHCPRelays() ([]*dhcprelay.Relay, error) {
	response := netbox.NetboxResponse[dhcprelay.Relay]{}
	params := deviceDatacenterFilter()

	err := netbox.Get("/api/plugins/cmdb/dhcp-relay-interface/", &response, params)
	if err != nil {
		return nil, fmt.Errorf("DHCP relay fetching failure: %w", err)
	}

	if len(response.Results) == 0 {
		log.Warn().Msg("no DHCP relay interface found")
	}

	return response.Results, nil
}

// PrecomputeDHCPRelays groups the relays by the device owning their interface.
func PrecomputeDHCPRelays(relays []*dhcprelay.Relay) map[string][]*dhcprelay.Relay {
	relaysPerDevice := make(map[string][]*dhcprelay.Relay)
	for _, relay := range relays {
		device := relay.LogicalInterface.ParentInterface.Device.Name
		relaysPerDevice[device] = append(relaysPerDevice[device], relay)
	}

	return relaysPerDevice
}
