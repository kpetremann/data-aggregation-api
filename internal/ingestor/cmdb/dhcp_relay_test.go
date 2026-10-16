package cmdb_test

import (
	"encoding/json"
	"testing"

	"github.com/criteo/data-aggregation-api/internal/ingestor/cmdb"
	"github.com/criteo/data-aggregation-api/internal/model/cmdb/dhcprelay"
	"github.com/google/go-cmp/cmp"
)

func TestPrecomputeDHCPRelays(t *testing.T) {
	// Shape of /api/plugins/cmdb/dhcp-relay-interface/: two relays on the same
	// switch, and one on a server, which has no device configuration.
	args := `[
		{
			"id": 1,
			"logical_interface": {
				"id": 2270, "index": 0, "name": "vlan506.0",
				"parent_interface": {"id": 2270, "name": "vlan506", "device": {"id": 1, "name": "tor01-01"}}
			},
			"server_list": [{"id": 1, "server_address": "10.1.8.67"}]
		},
		{
			"id": 3,
			"logical_interface": {
				"id": 2247, "index": 0, "name": "to_server_58.0",
				"parent_interface": {"id": 2247, "name": "to_server_58", "device": {"id": 1, "name": "tor01-01"}}
			},
			"server_list": [{"id": 1, "server_address": "10.1.8.67"}]
		},
		{
			"id": 2,
			"logical_interface": {
				"id": 2248, "index": 0, "name": "server058.0",
				"parent_interface": {"id": 2248, "name": "server058", "device": {"id": 2, "name": "generic-server-ty15"}}
			},
			"server_list": [{"id": 1, "server_address": "10.1.8.67"}]
		}
	]`

	var cmdbOutput []*dhcprelay.Relay
	if err := json.Unmarshal([]byte(args), &cmdbOutput); err != nil {
		t.Fatalf("unable to load test data: %s", err)
	}

	out := cmdb.PrecomputeDHCPRelays(cmdbOutput)

	if len(out) != 2 || len(out["tor01-01"]) != 2 || len(out["generic-server-ty15"]) != 1 {
		t.Fatalf("relays are not grouped by device: %v", out)
	}

	want := &dhcprelay.Relay{
		ServerList: []dhcprelay.Server{{ServerAddress: "10.1.8.67"}},
	}
	want.LogicalInterface.Name = "vlan506.0"
	want.LogicalInterface.ParentInterface.Name = "vlan506"
	want.LogicalInterface.ParentInterface.Device.Name = "tor01-01"

	if diff := cmp.Diff(want, out["tor01-01"][0]); diff != "" {
		t.Errorf("unexpected diff: %s\n", diff)
	}
}
