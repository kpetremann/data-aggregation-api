package bgp

import (
	"github.com/criteo/data-aggregation-api/internal/model/cmdb/bgp"
	"github.com/criteo/data-aggregation-api/internal/model/openconfig"
	"github.com/criteo/data-aggregation-api/internal/util"
)

// PeerGroupsToOpenconfig converts precomputed AFI/SAFI to OpenConfig.
// OpenConfig path: /network-instances/network-instance/protocols/protocol/bgp/peer-groups.
//
// Deprecated: peer-groups will be removed from the CMDB in future releases.
// You should migrate to configuration without using peer-groups.
func PeerGroupsToOpenconfig(peerGroups []*bgp.PeerGroup) map[string]*openconfig.NetworkInstance_Protocol_Bgp_PeerGroup {
	var groups = make(map[string]*openconfig.NetworkInstance_Protocol_Bgp_PeerGroup)

	for _, peerGroup := range peerGroups {
		policy := openconfig.NetworkInstance_Protocol_Bgp_PeerGroup_ApplyPolicy{}
		if peerGroup.RoutePolicyIn != nil {
			policy.ImportPolicy = util.AppendIfDefined(policy.ImportPolicy, peerGroup.RoutePolicyIn.Name)
		}
		if peerGroup.RoutePolicyOut != nil {
			policy.ExportPolicy = util.AppendIfDefined(policy.ExportPolicy, peerGroup.RoutePolicyOut.Name)
		}

		group := openconfig.NetworkInstance_Protocol_Bgp_PeerGroup{
			PeerGroupName: &peerGroup.Name,
			Description:   &peerGroup.Description,
			ApplyPolicy:   &policy,
		}

		if peerGroup.LocalAsn != nil {
			group.LocalAs = peerGroup.LocalAsn.Number
		}
		if peerGroup.RemoteAsn != nil {
			group.PeerAs = peerGroup.RemoteAsn.Number
		}

		groups[peerGroup.Name] = &group
	}

	return groups
}
