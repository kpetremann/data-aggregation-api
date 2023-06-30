package cmdb_test

import (
	"encoding/json"
	"net"
	"testing"

	"github.com/criteo/data-aggregation-api/internal/ingestor/cmdb"
	"github.com/criteo/data-aggregation-api/internal/model/cmdb/common"
	"github.com/criteo/data-aggregation-api/internal/model/cmdb/routingpolicy"
	"github.com/google/go-cmp/cmp"
)

func TestPrecomputeRoutePolicies(t *testing.T) {
	var repeatASN3 uint8 = 3

	tests := []struct {
		name string
		args string
		want map[string][]*routingpolicy.RoutePolicy
	}{
		{
			name: "lite route-map",
			args: `
			[
				{
					"id": 1,
					"name": "SERVER:PROD:OUT",
					"device": {
						"id": 1,
						"name": "tor01-01"
					},
					"description": "Server VLAN prod",
					"terms": [
						{
							"description": "",
							"sequence": 1,
							"decision": "permit",
							"from_bgp_community": "",
							"from_bgp_community_list": null,
							"from_prefix_list": null,
							"from_source_protocol": "",
							"from_route_type": "",
							"from_local_pref": null,
							"set_local_pref": null,
							"set_community": "65000:1000",
							"set_origin": "",
							"set_metric": null,
							"set_large_community": "",
							"set_as_path_prepend_repeat": null,
							"set_as_path_prepend_asn": null,
							"set_next_hop": null
						}
					]
				}
			]
			`,
			want: map[string][]*routingpolicy.RoutePolicy{
				"tor01-01": {
					&routingpolicy.RoutePolicy{
						Name: "SERVER:PROD:OUT",
						Device: struct {
							Name string `json:"name" validate:"required"`
						}{
							Name: "tor01-01",
						},
						Terms: []*routingpolicy.RoutePolicyTerm{
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
								SetCommunity:           "65000:1000",
								SetOrigin:              "",
								SetMetric:              0,
								SetLargeCommunity:      "",
								SetASPathPrependASN:    nil,
								SetASPathPrependRepeat: 0,
								SetNextHop:             net.IP{},
							},
						},
					},
				},
			},
		},
		{
			name: "full route-map",
			args: `
			[
				{
					"id": 2,
					"name": "LAN:IN",
					"device": {
						"id": 1,
						"name": "tor01-01"
					},
					"description": "",
					"terms": [
						{
							"description": "",
							"sequence": 1,
							"decision": "permit",
							"from_bgp_community": "65000:999",
							"from_bgp_community_list": null,
							"from_prefix_list": null,
							"from_source_protocol": "",
							"from_route_type": "",
							"from_local_pref": null,
							"set_local_pref": null,
							"set_community": "",
							"set_origin": "",
							"set_metric": null,
							"set_large_community": "",
							"set_as_path_prepend_repeat": null,
							"set_as_path_prepend_asn": null,
							"set_next_hop": null
						},
						{
							"description": "makes_no_sense_xD",
							"sequence": 2,
							"decision": "permit",
							"from_bgp_community": "65000:1234",
							"from_bgp_community_list": {
								"id": 1,
								"device": 1,
								"name": "SERVERS"
							},
							"from_prefix_list": {
								"id": 1,
								"device": 1,
								"name": "SERVER:VLAN:PROD"
							},
							"from_source_protocol": "bgp",
							"from_route_type": "ibgp",
							"from_local_pref": 100,
							"set_local_pref": 1000,
							"set_community": "65000:12345",
							"set_origin": "igp",
							"set_metric": 123,
							"set_large_community": "",
							"set_as_path_prepend_repeat": 3,
							"set_as_path_prepend_asn": {
								"id": 1,
								"number": 65000,
								"organization_name": "Lab-65000"
							},
							"set_next_hop": "127.0.0.1"
						}
					]
				}
			]
			`,
			want: map[string][]*routingpolicy.RoutePolicy{
				"tor01-01": {
					&routingpolicy.RoutePolicy{
						Name: "LAN:IN",
						Device: struct {
							Name string `json:"name" validate:"required"`
						}{
							Name: "tor01-01",
						},
						Terms: []*routingpolicy.RoutePolicyTerm{
							{
								Sequence:    1,
								Decision:    "permit",
								Description: "",
								FromBGPCommunityList: (*struct {
									Name string `json:"name" validate:"required"`
								})(nil),
								FromPrefixList: (*struct {
									Name string `json:"name" validate:"required"`
								})(nil),
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
								SetNextHop:             net.IP(nil),
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
								FromRouteType:          "ibgp",
								FromLocalPref:          100,
								SetLocalPref:           1000,
								SetCommunity:           "65000:12345",
								SetOrigin:              "igp",
								SetMetric:              123,
								SetLargeCommunity:      "",
								SetASPathPrependASN:    &common.ASN{Number: &as65000, Organization: "Lab-65000"},
								SetASPathPrependRepeat: repeatASN3,
								SetNextHop: net.IP{
									0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 255, 127, 0, 0, 1,
								},
							},
						},
					},
				},
			},
		},
	}

	for _, test := range tests {
		var cmdbOutput []*routingpolicy.RoutePolicy
		if err := json.Unmarshal([]byte(test.args), &cmdbOutput); err != nil {
			t.Errorf("unable to load test data for '%s': %s", test.name, err)
			continue
		}

		out := cmdb.PrecomputeRoutePolicies(cmdbOutput)
		if diff := cmp.Diff(out, test.want); diff != "" {
			t.Errorf("unexpected diff for '%s': %s\n", test.name, diff)
		}
	}
}
