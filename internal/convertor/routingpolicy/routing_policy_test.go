package routingpolicy_test

import (
	"net"
	"testing"

	"github.com/criteo/data-aggregation-api/internal/convertor/routingpolicy"
	"github.com/criteo/data-aggregation-api/internal/model/cmdb/common"
	cmdbRP "github.com/criteo/data-aggregation-api/internal/model/cmdb/routingpolicy"
	"github.com/criteo/data-aggregation-api/internal/model/openconfig"
	"github.com/criteo/data-aggregation-api/internal/types"
	"github.com/google/go-cmp/cmp"
)

func TestRoutingPolicyToOpenconfig(t *testing.T) {
	var as65000 uint32 = 65000
	var serverCommunitySetName = "SERVERS"
	var serverVlanProdPFName = "SERVER:VLAN:PROD"
	var prefixSet0 = "192.0.2.0/28"
	var prefixSet16 = "192.0.2.16/28"
	var exactRange = "exact"
	var range30To32 = "30..32"
	var pfLanInName = "LAN:IN"
	var seq1 = "1"
	var seq2 = "2"
	var localPref100 uint32 = 100
	var localPref1000 uint32 = 1000
	var nexthopLocalhost = openconfig.UnionString("127.0.0.1")
	var repeatASN3 uint8 = 3

	prefixLists := []*cmdbRP.PrefixList{
		{
			Name: "SERVER:VLAN:PROD",
			Device: struct {
				Name string `json:"name" validate:"required"`
			}{
				Name: "tor01-01",
			},
			Terms: []*cmdbRP.PrefixListTerm{
				{
					Sequence: 10,
					Prefix: types.CIDR{
						IP: net.IP{
							0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 255, 192, 0, 2, 0,
						},
						Netmask: 28,
					},
					LessOrEqual:    0,
					GreaterOrEqual: 0,
				},
				{
					Sequence: 20,
					Prefix: types.CIDR{
						IP: net.IP{
							0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 255, 192, 0, 2, 16,
						},
						Netmask: 28,
					},
					LessOrEqual:    32,
					GreaterOrEqual: 30,
				},
			},
			IPVersion: "ipv4",
		},
	}

	communityLists := []*cmdbRP.CommunityList{
		{
			Name: "SERVERS",
			Device: struct {
				Name string `json:"name" validate:"required"`
			}{
				Name: "tor01-01",
			},
			Terms: []*cmdbRP.CommunityListTerm{
				{
					Sequence:  10,
					Community: "650..:999",
				},
				{
					Sequence:  20,
					Community: "650..:1000",
				},
			},
		},
	}

	routingPolicies := []*cmdbRP.RoutePolicy{
		{
			Name: "LAN:IN",
			Device: struct {
				Name string `json:"name" validate:"required"`
			}{
				Name: "tor01-01",
			},
			Terms: []*cmdbRP.RoutePolicyTerm{
				{
					Sequence:               1,
					Decision:               "permit",
					Description:            "",
					FromBGPCommunityList:   nil,
					FromPrefixList:         nil,
					FromSourceProtocol:     "",
					FromRouteType:          "",
					FromLocalPref:          0,
					SetLocalPref:           0,
					SetCommunity:           "",
					SetOrigin:              "",
					SetMetric:              0,
					SetLargeCommunity:      "",
					SetASPathPrependASN:    nil,
					SetASPathPrependRepeat: 0,
					SetNextHop:             nil,
				},
				{
					Sequence:    2,
					Decision:    "permit",
					Description: "makes_no_sense_xD",
					FromBGPCommunityList: &struct {
						Name string `json:"name" validate:"required"`
					}{
						Name: "SERVERS",
					},
					FromPrefixList: &struct {
						Name string `json:"name" validate:"required"`
					}{
						Name: "SERVER:VLAN:PROD",
					},
					FromSourceProtocol:     "bgp",
					FromRouteType:          "igp",
					FromLocalPref:          100,
					SetLocalPref:           1000,
					SetCommunity:           "65000:12345",
					SetOrigin:              "igp",
					SetMetric:              123,
					SetLargeCommunity:      "",
					SetASPathPrependASN:    &common.ASN{Number: &as65000, Organization: "Lab-65000"},
					SetASPathPrependRepeat: repeatASN3,
					SetNextHop: &net.IP{
						0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 255, 127, 0, 0, 1,
					},
				},
			},
		},
	}

	want := &openconfig.RoutingPolicy{
		DefinedSets: &openconfig.RoutingPolicy_DefinedSets{
			BgpDefinedSets: &openconfig.RoutingPolicy_DefinedSets_BgpDefinedSets{
				CommunitySet: map[string]*openconfig.RoutingPolicy_DefinedSets_BgpDefinedSets_CommunitySet{
					"SERVERS": {
						CommunityMember: []openconfig.RoutingPolicy_DefinedSets_BgpDefinedSets_CommunitySet_CommunityMember_Union{
							openconfig.UnionString("650..:999"),
							openconfig.UnionString("650..:1000"),
						},
						CommunitySetName: &serverCommunitySetName,
					},
				},
			},
			PrefixSet: map[string]*openconfig.RoutingPolicy_DefinedSets_PrefixSet{
				"SERVER:VLAN:PROD": {
					Mode: openconfig.PrefixSet_Mode_IPV4,
					Name: &serverVlanProdPFName,
					Prefix: map[openconfig.RoutingPolicy_DefinedSets_PrefixSet_Prefix_Key]*openconfig.RoutingPolicy_DefinedSets_PrefixSet_Prefix{
						{
							IpPrefix:        "192.0.2.0/28",
							MasklengthRange: "exact",
						}: {
							IpPrefix:        &prefixSet0,
							MasklengthRange: &exactRange,
						},
						{
							IpPrefix:        "192.0.2.16/28",
							MasklengthRange: "30..32",
						}: {
							IpPrefix:        &prefixSet16,
							MasklengthRange: &range30To32,
						},
					},
				},
			},
		},
		PolicyDefinition: map[string]*openconfig.RoutingPolicy_PolicyDefinition{
			"LAN:IN": {
				Name: &pfLanInName,
				Statement: map[string]*openconfig.RoutingPolicy_PolicyDefinition_Statement{
					"1": {
						Actions: &openconfig.RoutingPolicy_PolicyDefinition_Statement_Actions{
							BgpActions: &openconfig.RoutingPolicy_PolicyDefinition_Statement_Actions_BgpActions{
								SetAsPathPrepend: &openconfig.RoutingPolicy_PolicyDefinition_Statement_Actions_BgpActions_SetAsPathPrepend{
									Asn:     nil,
									RepeatN: nil,
								},
								SetCommunity:   nil,
								SetLocalPref:   nil,
								SetMed:         nil,
								SetNextHop:     nil,
								SetRouteOrigin: openconfig.BgpTypes_BgpOriginAttrType_UNSET,
							},
							PolicyResult: openconfig.RoutingPolicy_PolicyResultType_ACCEPT_ROUTE,
						},
						Conditions: &openconfig.RoutingPolicy_PolicyDefinition_Statement_Conditions{
							BgpConditions: &openconfig.RoutingPolicy_PolicyDefinition_Statement_Conditions_BgpConditions{
								CommunitySet: nil,
								LocalPrefEq:  nil,
							},
							InstallProtocolEq: 0,
							MatchPrefixSet:    nil,
						},
						Name: &seq1,
					},
					"2": {
						Actions: &openconfig.RoutingPolicy_PolicyDefinition_Statement_Actions{
							BgpActions: &openconfig.RoutingPolicy_PolicyDefinition_Statement_Actions_BgpActions{
								SetAsPathPrepend: &openconfig.RoutingPolicy_PolicyDefinition_Statement_Actions_BgpActions_SetAsPathPrepend{
									Asn:     &as65000,
									RepeatN: &repeatASN3,
								},
								SetCommunity: &openconfig.RoutingPolicy_PolicyDefinition_Statement_Actions_BgpActions_SetCommunity{
									Inline: &openconfig.RoutingPolicy_PolicyDefinition_Statement_Actions_BgpActions_SetCommunity_Inline{
										Communities: []openconfig.RoutingPolicy_PolicyDefinition_Statement_Actions_BgpActions_SetCommunity_Inline_Communities_Union{
											openconfig.UnionString("65000:12345"),
										},
									},
									Method:  openconfig.SetCommunity_Method_INLINE,
									Options: openconfig.BgpPolicy_BgpSetCommunityOptionType_REPLACE,
								},
								SetLocalPref:   &localPref1000,
								SetMed:         openconfig.UnionUint32(123),
								SetNextHop:     &nexthopLocalhost,
								SetRouteOrigin: openconfig.BgpTypes_BgpOriginAttrType_IGP,
							},
							PolicyResult: 1,
						},
						Conditions: &openconfig.RoutingPolicy_PolicyDefinition_Statement_Conditions{
							BgpConditions: &openconfig.RoutingPolicy_PolicyDefinition_Statement_Conditions_BgpConditions{
								AfiSafiIn:    []openconfig.E_BgpTypes_AFI_SAFI_TYPE(nil),
								CommunitySet: &serverCommunitySetName,
								LocalPrefEq:  &localPref100,
							},
							InstallProtocolEq: 1,
							MatchPrefixSet: &openconfig.RoutingPolicy_PolicyDefinition_Statement_Conditions_MatchPrefixSet{
								MatchSetOptions: 0,
								PrefixSet:       &serverVlanProdPFName,
							},
						},
						Name: &seq2,
					},
				},
			},
		},
	}

	ret, err := routingpolicy.RoutingPolicyToOpenconfig("tor01-01", prefixLists, communityLists, routingPolicies)
	if err != nil {
		t.Errorf("failed to convert routing policies to OpenConfig")
	}

	if diff := cmp.Diff(ret, want); diff != "" {
		t.Errorf("unexpected diff for '%s': %s\n", "BGP integration test", diff)
	}
}
