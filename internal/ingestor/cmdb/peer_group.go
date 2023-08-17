package cmdb

import (
	"fmt"

	"github.com/rs/zerolog/log"

	"github.com/criteo/data-aggregation-api/internal/ingestor/netbox"
	"github.com/criteo/data-aggregation-api/internal/model/cmdb/bgp"
)

// GetPeerGroups returns all peer-groups from the Network CMDB.
//
// Deprecated: peer-groups will be removed from the CMDB in future releases.
// You should migrate to configuration without using peer-groups.
func GetPeerGroups() ([]*bgp.PeerGroup, error) {
	response := netbox.NetboxResponse[bgp.PeerGroup]{}
	params := deviceDatacenterFilter()

	err := netbox.Get("/api/plugins/cmdb/peer-groups/", &response, params)
	if err != nil {
		return nil, fmt.Errorf("peer-groups fetching failure: %w", err)
	}

	if response.Count != len(response.Results) {
		log.Warn().Msg("some peer-groups have not been fetched")
	}

	return response.Results, nil
}

// PrecomputePeerGroups associates each found peer-groups to the matching devices.
//
// Deprecated: peer-groups will be removed from the CMDB in future releases.
// You should migrate to configuration without using peer-groups.
func PrecomputePeerGroups(peerGroups []*bgp.PeerGroup) map[string][]*bgp.PeerGroup {
	var peerGroupsPerDevice = make(map[string][]*bgp.PeerGroup)
	for _, peerGroup := range peerGroups {
		peerGroupsPerDevice[peerGroup.Device.Name] = append(peerGroupsPerDevice[peerGroup.Device.Name], peerGroup)
	}

	return peerGroupsPerDevice
}
