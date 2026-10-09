package tacacs_test

import (
	"bytes"
	"encoding/json"
	"testing"

	tacacsconvertors "github.com/criteo/data-aggregation-api/internal/convertor/tacacs"
	"github.com/criteo/data-aggregation-api/internal/model/cmdb/tacacs"
	"github.com/criteo/data-aggregation-api/internal/model/openconfig"
	"github.com/google/go-cmp/cmp"
	"github.com/openconfig/ygot/ygot"
)

// emitOpenConfig renders the aaa container the way device.go does, through the
// system container of the openconfig device. It also validates the output against
// the YANG models. The result is compacted to keep the expected values readable.
func emitOpenConfig(t *testing.T, aaa *openconfig.System_Aaa) string {
	t.Helper()

	out, err := ygot.EmitJSON(
		&openconfig.Device{System: &openconfig.System{Aaa: aaa}},
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

func TestTacacsToOpenConfigAAA(t *testing.T) {
	tests := []struct {
		name string
		args *tacacs.Tacacs
		want map[string]uint32 // server address -> priority
	}{
		{
			name: "one server",
			args: &tacacs.Tacacs{
				Passkey: "Mypasword1",
				ServerList: []tacacs.Server{
					{ServerAddress: "1.1.1.1", Priority: 1, TCPPort: 49},
				},
			},
			want: map[string]uint32{"1.1.1.1": 1},
		},
		{
			name: "several servers, each keeps its own priority",
			args: &tacacs.Tacacs{
				Passkey: "Mypasword1",
				ServerList: []tacacs.Server{
					{ServerAddress: "1.1.1.3", Priority: 1, TCPPort: 49},
					{ServerAddress: "1.1.1.1", Priority: 3, TCPPort: 49},
					{ServerAddress: "1.1.1.2", Priority: 2, TCPPort: 4949},
				},
			},
			want: map[string]uint32{"1.1.1.3": 1, "1.1.1.1": 3, "1.1.1.2": 2},
		},
	}

	for _, test := range tests {
		out, err := tacacsconvertors.TacacsToOpenConfigAAA(test.args)
		if err != nil {
			t.Errorf("unexpected error for '%s': %s", test.name, err)
			continue
		}

		serverGroup := out.ServerGroup["TACACS+"]
		got := map[string]uint32{}
		for address, server := range serverGroup.Server {
			if server.Priority != nil {
				got[address] = *server.Priority
			}
		}

		if diff := cmp.Diff(got, test.want); diff != "" {
			t.Errorf("unexpected server priorities for '%s': %s\n", test.name, diff)
		}

		emitOpenConfig(t, out)
	}
}

// TestTacacsToOpenConfigAAAJSON checks the rendered RFC7951 JSON of a full configuration.
func TestTacacsToOpenConfigAAAJSON(t *testing.T) {
	want := `{"system":{"aaa":{"server-groups":{"server-group":[{"config":{"name":"TACACS+","type":"TACACS"},"name":"TACACS+","servers":{"server":[` +
		`{"address":"1.1.1.1","config":{"address":"1.1.1.1","priority":2},"tacacs":{"config":{"port":49,"secret-key":"Mypasword1"}}},` +
		`{"address":"1.1.1.2","config":{"address":"1.1.1.2","priority":1},"tacacs":{"config":{"port":4949,"secret-key":"Mypasword1"}}}` +
		`]}}]}}}}`

	config := &tacacs.Tacacs{
		Passkey: "Mypasword1",
		ServerList: []tacacs.Server{
			{ServerAddress: "1.1.1.2", Priority: 1, TCPPort: 4949},
			{ServerAddress: "1.1.1.1", Priority: 2, TCPPort: 49},
		},
	}

	out, err := tacacsconvertors.TacacsToOpenConfigAAA(config)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if diff := cmp.Diff(emitOpenConfig(t, out), want); diff != "" {
		t.Errorf("unexpected diff: %s\n", diff)
	}
}

// TestTacacsToOpenConfigAAADefaults checks the leaves left out when the CMDB has
// nothing to put in them.
func TestTacacsToOpenConfigAAADefaults(t *testing.T) {
	want := `{"system":{"aaa":{"server-groups":{"server-group":[{"config":{"name":"TACACS+","type":"TACACS"},"name":"TACACS+","servers":{"server":[` +
		`{"address":"1.1.1.1","config":{"address":"1.1.1.1"}}` +
		`]}}]}}}}`

	config := &tacacs.Tacacs{ServerList: []tacacs.Server{{ServerAddress: "1.1.1.1"}}}

	out, err := tacacsconvertors.TacacsToOpenConfigAAA(config)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if diff := cmp.Diff(emitOpenConfig(t, out), want); diff != "" {
		t.Errorf("unexpected diff: %s\n", diff)
	}
}

// TestTacacsToOpenConfigAAANoServer checks nothing is generated without a server:
// the passkey alone has no place in the model.
func TestTacacsToOpenConfigAAANoServer(t *testing.T) {
	tests := map[string]*tacacs.Tacacs{
		"no configuration at all": nil,
		"empty configuration":     {},
		"passkey but no server":   {Passkey: "Mypasword1"},
	}

	for name, config := range tests {
		out, err := tacacsconvertors.TacacsToOpenConfigAAA(config)
		if err != nil {
			t.Errorf("unexpected error for '%s': %s", name, err)
			continue
		}

		if out != nil {
			t.Errorf("unexpected aaa container for '%s': %v", name, out)
		}
	}
}
