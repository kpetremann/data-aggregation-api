package cmdb_test

import (
	"encoding/json"
	"net"
	"testing"

	"github.com/criteo/data-aggregation-api/internal/ingestor/cmdb"
	"github.com/criteo/data-aggregation-api/internal/model/cmdb/routingpolicy"
	"github.com/criteo/data-aggregation-api/internal/types"
	"github.com/google/go-cmp/cmp"
)

func TestPrecomputePrefixLists(t *testing.T) {
	tests := []struct {
		name string
		args string
		want map[string][]*routingpolicy.PrefixList
	}{
		{
			name: "valid community list",
			args: `
			[
				{
					"id": 1,
					"name": "SERVER:VLAN:PROD",
					"device": {
						"id": 1,
						"name": "tor01-01"
					},
					"ip_version": "ipv4",
					"terms": [
						{
							"prefix": "192.0.2.0/28",
							"le": null,
							"ge": null
						},
						{
							"prefix": "192.0.2.16/28",
							"le": 32,
							"ge": 30
						}
					]
				}
			]
			`,
			want: map[string][]*routingpolicy.PrefixList{
				"tor01-01": {
					&routingpolicy.PrefixList{
						Name: "SERVER:VLAN:PROD",
						Device: struct {
							Name string `json:"name" validate:"required"`
						}{
							Name: "tor01-01",
						},
						Terms: []*routingpolicy.PrefixListTerm{
							{
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
				},
			},
		},
	}

	for _, test := range tests {
		var cmdbOutput []*routingpolicy.PrefixList
		if err := json.Unmarshal([]byte(test.args), &cmdbOutput); err != nil {
			t.Errorf("unable to load test data for '%s': %s", test.name, err)
			continue
		}

		out := cmdb.PrecomputePrefixLists(cmdbOutput)
		if diff := cmp.Diff(out, test.want); diff != "" {
			t.Errorf("unexpected diff for '%s': %s\n", test.name, diff)
		}
	}
}
