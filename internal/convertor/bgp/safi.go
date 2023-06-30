package bgp

import (
	"fmt"

	"github.com/criteo/data-aggregation-api/internal/model/cmdb/bgp"
	"github.com/criteo/data-aggregation-api/internal/model/openconfig"
	"github.com/criteo/data-aggregation-api/internal/util"
)

// SafiToOpenconfig converts precomputed AFI/SAFI to OpenConfig.
// OpenConfig path: /network-instances/network-instance/protocols/protocol/bgp/neighbors/neighbor/afi-safis/.
func SafiToOpenconfig(localInfo *bgp.DeviceSession) (map[openconfig.E_BgpTypes_AFI_SAFI_TYPE]*openconfig.NetworkInstance_Protocol_Bgp_Neighbor_AfiSafi, error) {
	var safis = make(map[openconfig.E_BgpTypes_AFI_SAFI_TYPE]*openconfig.NetworkInstance_Protocol_Bgp_Neighbor_AfiSafi)

	for _, safi := range localInfo.AfiSafis {
		enabled := true

		policy := openconfig.NetworkInstance_Protocol_Bgp_Neighbor_AfiSafi_ApplyPolicy{}
		if safi.RoutePolicyIn != nil {
			policy.ImportPolicy = util.AppendIfDefined(policy.ImportPolicy, safi.RoutePolicyIn.Name)
		}

		if safi.RoutePolicyOut != nil {
			policy.ExportPolicy = util.AppendIfDefined(policy.ExportPolicy, safi.RoutePolicyOut.Name)
		}

		newSafi := openconfig.NetworkInstance_Protocol_Bgp_Neighbor_AfiSafi{
			Enabled:     &enabled,
			ApplyPolicy: &policy,
		}

		switch safi.Name {
		case bgp.IPv4Unicast:
			newSafi.AfiSafiName = openconfig.BgpTypes_AFI_SAFI_TYPE_IPV4_UNICAST
			newSafi.Ipv4Unicast = &openconfig.NetworkInstance_Protocol_Bgp_Neighbor_AfiSafi_Ipv4Unicast{
				PrefixLimit: &openconfig.NetworkInstance_Protocol_Bgp_Neighbor_AfiSafi_Ipv4Unicast_PrefixLimit{
					MaxPrefixes: &localInfo.MaximumPrefixes,
				},
			}
		case bgp.IPv6Unicast:
			newSafi.AfiSafiName = openconfig.BgpTypes_AFI_SAFI_TYPE_IPV6_UNICAST
			newSafi.Ipv6Unicast = &openconfig.NetworkInstance_Protocol_Bgp_Neighbor_AfiSafi_Ipv6Unicast{
				PrefixLimit: &openconfig.NetworkInstance_Protocol_Bgp_Neighbor_AfiSafi_Ipv6Unicast_PrefixLimit{
					MaxPrefixes: &localInfo.MaximumPrefixes,
				},
			}
		case bgp.L2vpnEvpn:
			newSafi.AfiSafiName = openconfig.BgpTypes_AFI_SAFI_TYPE_L2VPN_EVPN
			newSafi.L2VpnEvpn = &openconfig.NetworkInstance_Protocol_Bgp_Neighbor_AfiSafi_L2VpnEvpn{
				PrefixLimit: &openconfig.NetworkInstance_Protocol_Bgp_Neighbor_AfiSafi_L2VpnEvpn_PrefixLimit{
					MaxPrefixes: &localInfo.MaximumPrefixes,
				},
			}
		default:
			return nil, fmt.Errorf("unsupported SAFI %s for %s", newSafi.AfiSafiName, localInfo.Device.Name)
		}

		safis[newSafi.AfiSafiName] = &newSafi
	}

	return safis, nil
}
