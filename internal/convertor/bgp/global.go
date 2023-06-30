package bgp

import (
	"github.com/criteo/data-aggregation-api/internal/model/cmdb/bgp"
	"github.com/criteo/data-aggregation-api/internal/model/openconfig"
)

// BGPGlobalToOpenconfig converts precomputed prefix-lists to OpenConfig.
// OpenConfig path: /network-instances/network-instance/protocols/protocol/bgp/global/.
func BGPGlobalToOpenconfig(bgpGlobal *bgp.BGPGlobal) *openconfig.NetworkInstance_Protocol_Bgp_Global {
	// TODO: remove once BGP global configuration will be mandatory
	if bgpGlobal == nil {
		return nil
	}

	cfg := openconfig.NetworkInstance_Protocol_Bgp_Global{
		As: bgpGlobal.LocalAsn.Number,
		DefaultRouteDistance: &openconfig.NetworkInstance_Protocol_Bgp_Global_DefaultRouteDistance{
			ExternalRouteDistance: bgpGlobal.EBGPAdministrativeDistance,
			InternalRouteDistance: bgpGlobal.IBGPAdministrativeDistance,
		},
		GracefulRestart: &openconfig.NetworkInstance_Protocol_Bgp_Global_GracefulRestart{
			Enabled:     bgpGlobal.GracefulRestartEnabled,
			RestartTime: bgpGlobal.GracefulRestartTime,
		},
		UseMultiplePaths: &openconfig.NetworkInstance_Protocol_Bgp_Global_UseMultiplePaths{
			Enabled: bgpGlobal.EcmpEnabled,
			Ebgp: &openconfig.NetworkInstance_Protocol_Bgp_Global_UseMultiplePaths_Ebgp{
				MaximumPaths: bgpGlobal.EcmpMaximumPaths,
			},
			Ibgp: &openconfig.NetworkInstance_Protocol_Bgp_Global_UseMultiplePaths_Ibgp{
				MaximumPaths: bgpGlobal.EcmpMaximumPaths,
			},
		},
	}

	if bgpGlobal.RouterID != "" {
		cfg.RouterId = &bgpGlobal.RouterID
	}

	return &cfg
}
