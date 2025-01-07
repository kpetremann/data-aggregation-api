package bgp_test

import (
	"net"
	"testing"

	"github.com/criteo/data-aggregation-api/internal/convertor/bgp"
	cmdbBGP "github.com/criteo/data-aggregation-api/internal/model/cmdb/bgp"
	"github.com/criteo/data-aggregation-api/internal/model/cmdb/common"
	cmdbRP "github.com/criteo/data-aggregation-api/internal/model/cmdb/routingpolicy"
	"github.com/criteo/data-aggregation-api/internal/model/openconfig"
	"github.com/criteo/data-aggregation-api/internal/types"
	"github.com/google/go-cmp/cmp"
)

func TestBGPToOpenConfig(t *testing.T) {
	// some vars because we need pointers!
	var as65000 uint32 = 65000
	var as65001 uint32 = 65001
	var flagTrue = true
	var ebgpDistance221 uint8 = 221
	var ibgpDistance222 uint8 = 222
	var restartTime uint16 = 60
	var routerID5001 = "127.0.0.1:5001"
	var ecmpMaximumPaths64 uint32 = 64
	var maxPrefixes10 uint32 = 10
	var awesomePassword = "awesomepassword"
	var descriptionToTOR01 = "TO:TOR01-01"
	var neighborAddress0 = "192.0.2.0"

	globalConfig := cmdbBGP.BGPGlobal{
		Device: struct {
			Name string "json:\"name\" validate:\"required\""
		}{
			Name: "spine01-01",
		},
		LocalAsn: common.ASN{
			Number:       &as65001,
			Organization: "Lab-65001",
		},
		RouterID:                   "127.0.0.1:5001",
		EBGPAdministrativeDistance: &ebgpDistance221,
		IBGPAdministrativeDistance: &ibgpDistance222,
		GracefulRestartEnabled:     &flagTrue,
		GracefulRestartTime:        &restartTime,
		EcmpEnabled:                &flagTrue,
		EcmpMaximumPaths:           &ecmpMaximumPaths64,
	}

	sessions := []*cmdbBGP.Session{
		{
			PeerA: cmdbBGP.DeviceSession{
				Device: struct {
					Name string `json:"name" validate:"required"`
				}{
					Name: "tor01-01",
				},
				LocalAddress: cmdbBGP.Address{
					Address: types.CIDR{
						IP:      net.ParseIP("192.0.2.0"),
						Netmask: 31,
					},
					Family: 4,
				},
				LocalAsn: common.ASN{
					Number:       &as65000,
					Organization: "Lab-65000",
				},
				AfiSafis: []*cmdbBGP.AfiSafi{
					{
						Name: "ipv4-unicast",
						RoutePolicyIn: &cmdbRP.RoutePolicyLite{
							Name:        "LAN:IN:MAINTENANCE",
							Description: "MAINTENANCE MODE",
						},
						RoutePolicyOut: &cmdbRP.RoutePolicyLite{
							Name:        "LAN:OUT:MAINTENANCE",
							Description: "MAINTENANCE MODE",
						},
					},
				},
				PeerGroup: nil,
				RoutePolicyIn: &cmdbRP.RoutePolicyLite{
					Name:        "LAN:IN",
					Description: "",
				},
				RoutePolicyOut: &cmdbRP.RoutePolicyLite{
					Name:        "SERVER:PROD:OUT",
					Description: "Server VLAN prod",
				},
				Description:     "TO:SPINE01-01",
				MaximumPrefixes: 128,
				EnforceFirstAs:  true,
				Enabled:         &flagTrue,
			},
			PeerB: cmdbBGP.DeviceSession{
				Device: struct {
					Name string `json:"name" validate:"required"`
				}{
					Name: "spine01-01",
				},
				LocalAddress: cmdbBGP.Address{
					Address: types.CIDR{
						IP:      net.ParseIP("192.0.2.1"),
						Netmask: 31,
					},
					Family: 4,
				},
				LocalAsn: common.ASN{
					Number:       &as65001,
					Organization: "Lab-65001",
				},
				AfiSafis: []*cmdbBGP.AfiSafi{
					{
						Name: "ipv4-unicast",
						RoutePolicyIn: &cmdbRP.RoutePolicyLite{
							Name:        "SERVER:PROD:MAINTENANCE:IN",
							Description: "",
						},
						RoutePolicyOut: &cmdbRP.RoutePolicyLite{
							Name:        "LAN:MAINTENANCE:OUT",
							Description: "",
						},
					},
				},
				PeerGroup: nil,
				RoutePolicyIn: &cmdbRP.RoutePolicyLite{
					Name:        "SERVER:PROD:IN",
					Description: "",
				},
				RoutePolicyOut: &cmdbRP.RoutePolicyLite{
					Name:        "LAN:OUT",
					Description: "",
				},
				Description:     "TO:TOR01-01",
				MaximumPrefixes: 10,
				EnforceFirstAs:  true,
				Enabled:         &flagTrue,
			},
			Password: awesomePassword,
		},
	}

	want := &openconfig.NetworkInstance_Protocol_Bgp{
		Global: &openconfig.NetworkInstance_Protocol_Bgp_Global{
			As: &as65001,
			DefaultRouteDistance: &openconfig.NetworkInstance_Protocol_Bgp_Global_DefaultRouteDistance{
				ExternalRouteDistance: &ebgpDistance221,
				InternalRouteDistance: &ibgpDistance222,
			},
			GracefulRestart: &openconfig.NetworkInstance_Protocol_Bgp_Global_GracefulRestart{
				Enabled:     &flagTrue,
				RestartTime: &restartTime,
			},
			RouterId: &routerID5001,
			UseMultiplePaths: &openconfig.NetworkInstance_Protocol_Bgp_Global_UseMultiplePaths{
				Enabled: &flagTrue,
				Ebgp: &openconfig.NetworkInstance_Protocol_Bgp_Global_UseMultiplePaths_Ebgp{
					MaximumPaths: &ecmpMaximumPaths64,
				},
				Ibgp: &openconfig.NetworkInstance_Protocol_Bgp_Global_UseMultiplePaths_Ibgp{
					MaximumPaths: &ecmpMaximumPaths64,
				},
			},
		},
		Neighbor: map[string]*openconfig.NetworkInstance_Protocol_Bgp_Neighbor{
			"192.0.2.0": {
				AfiSafi: map[openconfig.E_BgpTypes_AFI_SAFI_TYPE]*openconfig.NetworkInstance_Protocol_Bgp_Neighbor_AfiSafi{
					openconfig.BgpTypes_AFI_SAFI_TYPE_IPV4_UNICAST: {
						AfiSafiName: openconfig.BgpTypes_AFI_SAFI_TYPE_IPV4_UNICAST,
						ApplyPolicy: &openconfig.NetworkInstance_Protocol_Bgp_Neighbor_AfiSafi_ApplyPolicy{
							ExportPolicy: []string{"LAN:MAINTENANCE:OUT"},
							ImportPolicy: []string{"SERVER:PROD:MAINTENANCE:IN"},
						},
						Enabled: &flagTrue,
						Ipv4Unicast: &openconfig.NetworkInstance_Protocol_Bgp_Neighbor_AfiSafi_Ipv4Unicast{
							PrefixLimit: &openconfig.NetworkInstance_Protocol_Bgp_Neighbor_AfiSafi_Ipv4Unicast_PrefixLimit{
								MaxPrefixes: &maxPrefixes10,
							},
						},
					},
				},
				ApplyPolicy: &openconfig.NetworkInstance_Protocol_Bgp_Neighbor_ApplyPolicy{
					ExportPolicy: []string{"LAN:OUT"},
					ImportPolicy: []string{"SERVER:PROD:IN"},
				},
				AuthPassword:    &awesomePassword,
				Description:     &descriptionToTOR01,
				Enabled:         &flagTrue,
				LocalAs:         &as65001,
				NeighborAddress: &neighborAddress0,
				PeerAs:          &as65000,
				PeerGroup:       nil,
				PeerType:        0,
				RemovePrivateAs: 0,
				SendCommunity:   0,
			},
		},
		PeerGroup: map[string]*openconfig.NetworkInstance_Protocol_Bgp_PeerGroup{},
	}

	ret, err := bgp.BGPToOpenconfig("spine01-01", &globalConfig, sessions, nil)
	if err != nil {
		t.Errorf("failed to convert BGP to OpenConfig")
	}

	if diff := cmp.Diff(ret, want); diff != "" {
		t.Errorf("unexpected diff for '%s': %s\n", "BGP integration test", diff)
	}
}
