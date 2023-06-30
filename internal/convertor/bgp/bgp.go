package bgp

import (
	"github.com/criteo/data-aggregation-api/internal/model/cmdb/bgp"
	"github.com/criteo/data-aggregation-api/internal/model/openconfig"
)

// BGPToOpenconfig converts all precomputed assets in BGP/* to OpenConfig.
// OpenConfig path: /network-instances/network-instance/protocols/protocol/bgp/.
func BGPToOpenconfig(hostname string, bgpGlobal *bgp.BGPGlobal, sessions []*bgp.Session, peerGroups []*bgp.PeerGroup) (*openconfig.NetworkInstance_Protocol_Bgp, error) {
	neighbors, err := NeighborsToOpenconfig(hostname, sessions)
	if err != nil {
		return nil, err
	}

	defaultInstance := openconfig.NetworkInstance_Protocol_Bgp{
		Global:    BGPGlobalToOpenconfig(bgpGlobal),
		Neighbor:  neighbors,
		PeerGroup: PeerGroupsToOpenconfig(peerGroups),
	}

	return &defaultInstance, err
}
