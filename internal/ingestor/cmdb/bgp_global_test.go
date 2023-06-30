package cmdb_test

import (
	"encoding/json"
	"testing"

	"github.com/criteo/data-aggregation-api/internal/ingestor/cmdb"
	"github.com/criteo/data-aggregation-api/internal/model/cmdb/bgp"
	"github.com/criteo/data-aggregation-api/internal/model/cmdb/common"
	"github.com/google/go-cmp/cmp"
)

func TestPrecomputeBGPGlobal(t *testing.T) {
	var flagFalse = false
	var flagTrue = true
	var ebgpDistance uint8 = 221
	var ibgpDistance uint8 = 222
	var restartTime uint16 = 60
	var ecmpMax32 uint32 = 32
	var ecmpMax64 uint32 = 64

	tests := []struct {
		name string
		args string
		want map[string]*bgp.BGPGlobal
	}{
		{
			name: "lite configuration",
			args: `
			[
				{
					"id": 1,
					"device": {
						"id": 1,
						"name": "tor01-01"
					},
					"local_asn": {
						"id": 1,
						"number": 65000,
						"organization_name": "Lab-65000"
					},
					"created": "2023-06-23T08:43:32.059591Z",
					"last_updated": "2023-06-23T08:43:32.059617Z",
					"router_id": "",
					"ebgp_administrative_distance": null,
					"ibgp_administrative_distance": null,
					"graceful_restart": false,
					"graceful_restart_time": null,
					"ecmp": false,
					"ecmp_maximum_paths": 32
				}
			]
			`,
			want: map[string]*bgp.BGPGlobal{
				"tor01-01": {
					Device: struct {
						Name string "json:\"name\" validate:\"required\""
					}{
						Name: "tor01-01",
					},
					LocalAsn: common.ASN{
						Number:       &as65000,
						Organization: "Lab-65000",
					},
					RouterID:                   "",
					EBGPAdministrativeDistance: nil,
					IBGPAdministrativeDistance: nil,
					GracefulRestartEnabled:     &flagFalse,
					GracefulRestartTime:        nil,
					EcmpEnabled:                &flagFalse,
					EcmpMaximumPaths:           &ecmpMax32,
				},
			},
		},
		{
			name: "full configuration",
			args: `
			[
				{
					"id": 2,
					"device": {
						"id": 2,
						"name": "spine01-01"
					},
					"local_asn": {
						"id": 2,
						"number": 65001,
						"organization_name": "Lab-65001"
					},
					"created": "2023-06-23T08:43:50.317612Z",
					"last_updated": "2023-06-23T08:43:50.317623Z",
					"router_id": "127.0.0.1:5001",
					"ebgp_administrative_distance": 221,
					"ibgp_administrative_distance": 222,
					"graceful_restart": true,
					"graceful_restart_time": 60,
					"ecmp": true,
					"ecmp_maximum_paths": 64
				}
			]
			`,
			want: map[string]*bgp.BGPGlobal{
				"spine01-01": {
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
					EBGPAdministrativeDistance: &ebgpDistance,
					IBGPAdministrativeDistance: &ibgpDistance,
					GracefulRestartEnabled:     &flagTrue,
					GracefulRestartTime:        &restartTime,
					EcmpEnabled:                &flagTrue,
					EcmpMaximumPaths:           &ecmpMax64,
				},
			},
		},
	}

	for _, test := range tests {
		var cmdbOutput []*bgp.BGPGlobal
		if err := json.Unmarshal([]byte(test.args), &cmdbOutput); err != nil {
			t.Errorf("unable to load test data for '%s': %s", test.name, err)
			continue
		}

		out := cmdb.PrecomputeBGPGlobal(cmdbOutput)
		if diff := cmp.Diff(out, test.want); diff != "" {
			t.Errorf("unexpected diff for '%s': %s\n", test.name, diff)
		}
	}
}
