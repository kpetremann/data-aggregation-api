package dhcprelay_test

import (
	"bytes"
	"encoding/json"
	"testing"

	dhcprelayconvertors "github.com/criteo/data-aggregation-api/internal/convertor/dhcprelay"
	"github.com/criteo/data-aggregation-api/internal/model/cmdb/dhcprelay"
	"github.com/criteo/data-aggregation-api/internal/model/openconfig"
	"github.com/google/go-cmp/cmp"
	"github.com/openconfig/ygot/ygot"
)

// emitOpenConfig renders the relay-agent container the way device.go does,
// through the openconfig device. It also validates the output against the YANG
// models. The result is compacted to keep the expected values readable.
func emitOpenConfig(t *testing.T, relayAgent *openconfig.RelayAgent) string {
	t.Helper()

	out, err := ygot.EmitJSON(
		&openconfig.Device{RelayAgent: relayAgent},
		&ygot.EmitJSONConfig{Format: ygot.RFC7951, SkipValidation: false, Indent: "  "},
	)
	if err != nil {
		t.Fatalf("failed to emit the openconfig configuration: %s", err)
	}

	compacted := bytes.NewBuffer(nil)
	if err := json.Compact(compacted, []byte(out)); err != nil {
		t.Fatalf("failed to compact the emitted configuration: %s", err)
	}

	return compacted.String()
}

// relay builds a CMDB relay on the sub-interface of a parent interface.
func relay(parent string, servers ...string) *dhcprelay.Relay {
	r := &dhcprelay.Relay{}
	r.LogicalInterface.Name = parent + ".0"
	r.LogicalInterface.ParentInterface.Name = parent
	r.LogicalInterface.ParentInterface.Device.Name = "tor01-01"

	for _, address := range servers {
		r.ServerList = append(r.ServerList, dhcprelay.Server{ServerAddress: address})
	}

	return r
}

func TestDHCPRelayToOpenConfig(t *testing.T) {
	tests := []struct {
		name string
		args []*dhcprelay.Relay
		want string
	}{
		{
			name: "nothing to advertise",
			args: nil,
			want: "{}",
		},
		{
			name: "a relay without server is skipped",
			args: []*dhcprelay.Relay{relay("vlan506")},
			want: "{}",
		},
		{
			name: "one SVI",
			args: []*dhcprelay.Relay{relay("vlan506", "10.1.8.67")},
			want: `{"relay-agent":{"dhcp":{"interfaces":{"interface":[{"config":{"enable":true,"helper-address":["10.1.8.67"],"id":"vlan506.0"},"id":"vlan506.0"}]}}}}`,
		},
		{
			name: "a routed port and an SVI, servers sorted",
			args: []*dhcprelay.Relay{
				relay("vlan506", "10.2.8.67", "10.1.8.67"),
				relay("Ethernet0", "10.1.8.67"),
			},
			want: `{"relay-agent":{"dhcp":{"interfaces":{"interface":[` +
				`{"config":{"enable":true,"helper-address":["10.1.8.67"],"id":"Ethernet0.0"},"id":"Ethernet0.0"},` +
				`{"config":{"enable":true,"helper-address":["10.1.8.67","10.2.8.67"],"id":"vlan506.0"},"id":"vlan506.0"}` +
				`]}}}}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			relayAgent, err := dhcprelayconvertors.DHCPRelayToOpenConfig(test.args)
			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			if diff := cmp.Diff(test.want, emitOpenConfig(t, relayAgent)); diff != "" {
				t.Errorf("unexpected diff: %s\n", diff)
			}
		})
	}
}
