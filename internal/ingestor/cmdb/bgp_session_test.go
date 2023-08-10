package cmdb_test

import (
	"encoding/json"
	"net"
	"testing"

	"github.com/criteo/data-aggregation-api/internal/ingestor/cmdb"
	"github.com/criteo/data-aggregation-api/internal/model/cmdb/bgp"
	"github.com/criteo/data-aggregation-api/internal/model/cmdb/common"
	"github.com/criteo/data-aggregation-api/internal/model/cmdb/routingpolicy"
	"github.com/criteo/data-aggregation-api/internal/types"
	"github.com/google/go-cmp/cmp"
)

var as65000 uint32 = 65000
var as65001 uint32 = 65001

/*
	BGP session - IPv4
*/

const ipv4BGPSession = `
[
    {
        "id": 1,
        "peer_a": {
            "id": 1,
			"enabled": true,
            "local_address": {
                "id": 1,
                "url": "http://127.0.0.1:8001/api/ipam/ip-addresses/1/",
                "display": "192.0.2.0/31",
                "family": 4,
                "address": "192.0.2.0/31"
            },
            "device": {
                "id": 1,
                "name": "tor01-01"
            },
            "peer_group": null,
            "local_asn": {
                "id": 1,
                "number": 65000,
                "organization_name": "Lab-65000"
            },
            "afi_safis": [
                {
                    "id": 1,
                    "route_policy_in": {
						"id": 5,
						"name": "LAN:IN:MAINTENANCE",
						"description": "MAINTENANCE MODE"
					},
                    "route_policy_out": {
						"id": 6,
						"name": "LAN:OUT:MAINTENANCE",
						"description": "MAINTENANCE MODE"
					},
                    "afi_safi_name": "ipv4-unicast"
                }
            ],
            "route_policy_in": {
                "id": 2,
                "name": "LAN:IN",
                "description": ""
            },
            "route_policy_out": {
                "id": 1,
                "name": "SERVER:PROD:OUT",
                "description": "Server VLAN prod"
            },
            "created": "2023-06-16T14:45:15.361680Z",
            "last_updated": "2023-06-16T14:49:18.095587Z",
            "description": "TO:SPINE01-01",
            "enforce_first_as": true,
            "maximum_prefixes": 128
        },
        "peer_b": {
            "id": 2,
			"enabled": false,
            "local_address": {
                "id": 2,
                "url": "http://127.0.0.1:8001/api/ipam/ip-addresses/2/",
                "display": "192.0.2.1/31",
                "family": 4,
                "address": "192.0.2.1/31"
            },
            "device": {
                "id": 2,
                "name": "spine01-01"
            },
            "peer_group": null,
            "local_asn": {
                "id": 2,
                "number": 65001,
                "organization_name": "Lab-65001"
            },
            "afi_safis": [
                {
                    "id": 2,
                    "route_policy_in": null,
                    "route_policy_out": null,
                    "afi_safi_name": "ipv4-unicast"
                }
            ],
            "route_policy_in": {
                "id": 3,
                "name": "SERVER:PROD:IN",
                "description": ""
            },
            "route_policy_out": {
                "id": 4,
                "name": "LAN:OUT",
                "description": ""
            },
            "created": "2023-06-16T14:45:33.940578Z",
            "last_updated": "2023-06-16T14:49:33.632184Z",
            "description": "TO:TOR01-01",
            "enforce_first_as": true,
            "maximum_prefixes": 10
        },
        "tenant": null,
        "created": "2023-06-16T14:45:42.263661Z",
        "last_updated": "2023-06-16T14:49:34.595313Z",
        "status": "active",
        "password": "awesomepassword",
        "circuit": null
    }
]
`

var ipv4Ingested = map[string][]*bgp.Session{
	"tor01-01": {
		&bgp.Session{
			PeerA: bgp.DeviceSession{
				Device: struct {
					Name string `json:"name" validate:"required"`
				}{
					Name: "tor01-01",
				},
				LocalAddress: bgp.Address{
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
				AfiSafis: []*bgp.AfiSafi{
					{
						Name: "ipv4-unicast",
						RoutePolicyIn: &routingpolicy.RoutePolicyLite{
							Name:        "LAN:IN:MAINTENANCE",
							Description: "MAINTENANCE MODE",
						},
						RoutePolicyOut: &routingpolicy.RoutePolicyLite{
							Name:        "LAN:OUT:MAINTENANCE",
							Description: "MAINTENANCE MODE",
						},
					},
				},
				PeerGroup: nil,
				RoutePolicyIn: &routingpolicy.RoutePolicyLite{
					Name:        "LAN:IN",
					Description: "",
				},
				RoutePolicyOut: &routingpolicy.RoutePolicyLite{
					Name:        "SERVER:PROD:OUT",
					Description: "Server VLAN prod",
				},
				Description:     "TO:SPINE01-01",
				MaximumPrefixes: 128,
				EnforceFirstAs:  true,
				Enabled:         true,
			},
			PeerB: bgp.DeviceSession{
				Device: struct {
					Name string `json:"name" validate:"required"`
				}{
					Name: "spine01-01",
				},
				LocalAddress: bgp.Address{
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
				AfiSafis: []*bgp.AfiSafi{
					{
						Name:           "ipv4-unicast",
						RoutePolicyIn:  nil,
						RoutePolicyOut: nil,
					},
				},
				PeerGroup: nil,
				RoutePolicyIn: &routingpolicy.RoutePolicyLite{
					Name:        "SERVER:PROD:IN",
					Description: "",
				},
				RoutePolicyOut: &routingpolicy.RoutePolicyLite{
					Name:        "LAN:OUT",
					Description: "",
				},
				Description:     "TO:TOR01-01",
				MaximumPrefixes: 10,
				EnforceFirstAs:  true,
				Enabled:         false,
			},
			Password: "awesomepassword",
		},
	},
	"spine01-01": {
		&bgp.Session{
			PeerA: bgp.DeviceSession{
				Device: struct {
					Name string `json:"name" validate:"required"`
				}{
					Name: "tor01-01",
				},
				LocalAddress: bgp.Address{
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
				AfiSafis: []*bgp.AfiSafi{
					{
						Name: "ipv4-unicast",
						RoutePolicyIn: &routingpolicy.RoutePolicyLite{
							Name:        "LAN:IN:MAINTENANCE",
							Description: "MAINTENANCE MODE",
						},
						RoutePolicyOut: &routingpolicy.RoutePolicyLite{
							Name:        "LAN:OUT:MAINTENANCE",
							Description: "MAINTENANCE MODE",
						},
					},
				},
				PeerGroup: nil,
				RoutePolicyIn: &routingpolicy.RoutePolicyLite{
					Name:        "LAN:IN",
					Description: "",
				},
				RoutePolicyOut: &routingpolicy.RoutePolicyLite{
					Name:        "SERVER:PROD:OUT",
					Description: "Server VLAN prod",
				},
				Description:     "TO:SPINE01-01",
				MaximumPrefixes: 128,
				EnforceFirstAs:  true,
				Enabled:         true,
			},
			PeerB: bgp.DeviceSession{
				Device: struct {
					Name string `json:"name" validate:"required"`
				}{
					Name: "spine01-01",
				},
				LocalAddress: bgp.Address{
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
				AfiSafis: []*bgp.AfiSafi{
					{
						Name:           "ipv4-unicast",
						RoutePolicyIn:  nil,
						RoutePolicyOut: nil,
					},
				},
				PeerGroup: nil,
				RoutePolicyIn: &routingpolicy.RoutePolicyLite{
					Name:        "SERVER:PROD:IN",
					Description: "",
				},
				RoutePolicyOut: &routingpolicy.RoutePolicyLite{
					Name:        "LAN:OUT",
					Description: "",
				},
				Description:     "TO:TOR01-01",
				MaximumPrefixes: 10,
				EnforceFirstAs:  true,
				Enabled:         false,
			},
			Password: "awesomepassword",
		},
	},
}

/*
	BGP session - IPv6
*/

const ipv6BGPSession = `
[
	{
		"id": 2,
		"peer_a": {
			"id": 3,
			"enabled": true,
			"local_address": {
				"id": 3,
				"url": "http://127.0.0.1:8001/api/ipam/ip-addresses/3/",
				"display": "2001:db8::100/127",
				"family": 6,
				"address": "2001:db8::100/127"
			},
			"device": {
				"id": 1,
				"name": "tor01-01"
			},
			"peer_group": null,
			"local_asn": {
				"id": 1,
				"number": 65000,
				"organization_name": "Lab-65000"
			},
			"afi_safis": [
				{
					"id": 3,
					"route_policy_in": null,
					"route_policy_out": null,
					"afi_safi_name": "ipv6-unicast"
				}
			],
			"route_policy_in": {
				"id": 2,
				"name": "LAN:IN",
				"description": ""
			},
			"route_policy_out": {
				"id": 1,
				"name": "SERVER:PROD:OUT",
				"description": "Server VLAN prod"
			},
			"created": "2023-06-16T15:39:34.049935Z",
			"last_updated": "2023-06-16T15:40:25.286003Z",
			"description": "TO:SPINE01-01:V6",
			"enforce_first_as": true,
			"maximum_prefixes": null
		},
		"peer_b": {
			"id": 4,
			"enabled": true,
			"local_address": {
				"id": 4,
				"url": "http://127.0.0.1:8001/api/ipam/ip-addresses/4/",
				"display": "2001:db8::101/127",
				"family": 6,
				"address": "2001:db8::101/127"
			},
			"device": {
				"id": 2,
				"name": "spine01-01"
			},
			"peer_group": null,
			"local_asn": {
				"id": 2,
				"number": 65001,
				"organization_name": "Lab-65001"
			},
			"afi_safis": [
				{
					"id": 4,
					"route_policy_in": null,
					"route_policy_out": null,
					"afi_safi_name": "ipv6-unicast"
				}
			],
			"route_policy_in": {
				"id": 3,
				"name": "SERVER:PROD:IN",
				"description": ""
			},
			"route_policy_out": {
				"id": 4,
				"name": "LAN:OUT",
				"description": ""
			},
			"created": "2023-06-16T15:40:15.649865Z",
			"last_updated": "2023-06-16T15:40:15.649875Z",
			"description": "TO:TOR01-01:V6",
			"enforce_first_as": true,
			"maximum_prefixes": null
		},
		"tenant": null,
		"created": "2023-06-16T15:40:26.373165Z",
		"last_updated": "2023-06-16T15:40:26.373175Z",
		"status": "active",
		"password": "awesomepassword",
		"circuit": null
	}
]
`

var ipv6Ingested = map[string][]*bgp.Session{
	"tor01-01": {
		&bgp.Session{
			PeerA: bgp.DeviceSession{
				Device: struct {
					Name string `json:"name" validate:"required"`
				}{
					Name: "tor01-01",
				},
				LocalAddress: bgp.Address{
					Address: types.CIDR{
						IP: net.IP{
							32, 1, 13, 184, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0,
						},
						Netmask: 127,
					},
					Family: 6,
				},
				LocalAsn: common.ASN{
					Number:       &as65000,
					Organization: "Lab-65000",
				},
				AfiSafis: []*bgp.AfiSafi{
					{
						Name:           "ipv6-unicast",
						RoutePolicyIn:  nil,
						RoutePolicyOut: nil,
					},
				},
				PeerGroup: nil,
				RoutePolicyIn: &routingpolicy.RoutePolicyLite{
					Name:        "LAN:IN",
					Description: "",
				},
				RoutePolicyOut: &routingpolicy.RoutePolicyLite{
					Name:        "SERVER:PROD:OUT",
					Description: "Server VLAN prod",
				},
				Description:     "TO:SPINE01-01:V6",
				MaximumPrefixes: 0,
				EnforceFirstAs:  true,
				Enabled:         true,
			},
			PeerB: bgp.DeviceSession{
				Device: struct {
					Name string `json:"name" validate:"required"`
				}{
					Name: "spine01-01",
				},
				LocalAddress: bgp.Address{
					Address: types.CIDR{
						IP: net.IP{
							32, 1, 13, 184, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1,
						},
						Netmask: 127,
					},
					Family: 6,
				},
				LocalAsn: common.ASN{
					Number:       &as65001,
					Organization: "Lab-65001",
				},
				AfiSafis: []*bgp.AfiSafi{
					{
						Name:           "ipv6-unicast",
						RoutePolicyIn:  nil,
						RoutePolicyOut: nil,
					},
				},
				PeerGroup: nil,
				RoutePolicyIn: &routingpolicy.RoutePolicyLite{
					Name:        "SERVER:PROD:IN",
					Description: "",
				},
				RoutePolicyOut: &routingpolicy.RoutePolicyLite{
					Name:        "LAN:OUT",
					Description: "",
				},
				Description:     "TO:TOR01-01:V6",
				MaximumPrefixes: 0,
				EnforceFirstAs:  true,
				Enabled:         true,
			},
			Password: "awesomepassword",
		},
	},
	"spine01-01": {
		&bgp.Session{
			PeerA: bgp.DeviceSession{
				Device: struct {
					Name string `json:"name" validate:"required"`
				}{
					Name: "tor01-01",
				},
				LocalAddress: bgp.Address{
					Address: types.CIDR{
						IP: net.IP{
							32, 1, 13, 184, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0,
						},
						Netmask: 127,
					},
					Family: 6,
				},
				LocalAsn: common.ASN{
					Number:       &as65000,
					Organization: "Lab-65000",
				},
				AfiSafis: []*bgp.AfiSafi{
					{
						Name:           "ipv6-unicast",
						RoutePolicyIn:  nil,
						RoutePolicyOut: nil,
					},
				},
				PeerGroup: nil,
				RoutePolicyIn: &routingpolicy.RoutePolicyLite{
					Name:        "LAN:IN",
					Description: "",
				},
				RoutePolicyOut: &routingpolicy.RoutePolicyLite{
					Name:        "SERVER:PROD:OUT",
					Description: "Server VLAN prod",
				},
				Description:     "TO:SPINE01-01:V6",
				MaximumPrefixes: 0,
				EnforceFirstAs:  true,
				Enabled:         true,
			},
			PeerB: bgp.DeviceSession{
				Device: struct {
					Name string `json:"name" validate:"required"`
				}{
					Name: "spine01-01",
				},
				LocalAddress: bgp.Address{
					Address: types.CIDR{
						IP: net.IP{
							32, 1, 13, 184, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1,
						},
						Netmask: 127,
					},
					Family: 6,
				},
				LocalAsn: common.ASN{
					Number:       &as65001,
					Organization: "Lab-65001",
				},
				AfiSafis: []*bgp.AfiSafi{
					{
						Name:           "ipv6-unicast",
						RoutePolicyIn:  nil,
						RoutePolicyOut: nil,
					},
				},
				PeerGroup: nil,
				RoutePolicyIn: &routingpolicy.RoutePolicyLite{
					Name:        "SERVER:PROD:IN",
					Description: "",
				},
				RoutePolicyOut: &routingpolicy.RoutePolicyLite{
					Name:        "LAN:OUT",
					Description: "",
				},
				Description:     "TO:TOR01-01:V6",
				MaximumPrefixes: 0,
				EnforceFirstAs:  true,
				Enabled:         true,
			},
			Password: "awesomepassword",
		},
	},
}

/*
	Tests
*/

func TestPrecomputeBGPSessions(t *testing.T) {
	tests := []struct {
		name string
		args string
		want map[string][]*bgp.Session
	}{
		{
			name: "IPv4 ToR <=> Spine",
			args: ipv4BGPSession,
			want: ipv4Ingested,
		},
		{
			name: "IPv6 ToR <=> Spine",
			args: ipv6BGPSession,
			want: ipv6Ingested,
		},
	}

	for _, test := range tests {
		var cmdbOutput []*bgp.Session
		if err := json.Unmarshal([]byte(test.args), &cmdbOutput); err != nil {
			t.Errorf("unable to load test data for '%s': %s", test.name, err)
			continue
		}

		out := cmdb.PrecomputeBGPSessions(cmdbOutput)
		if diff := cmp.Diff(out, test.want); diff != "" {
			t.Errorf("unexpected diff for '%s': %s\n", test.name, diff)
		}
	}
}
