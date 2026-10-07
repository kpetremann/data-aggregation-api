package device

import (
	"encoding/json"
	"sync"
	"testing"

	"github.com/criteo/data-aggregation-api/internal/model/cmdb/bgp"
	"github.com/criteo/data-aggregation-api/internal/model/cmdb/ntp"
	"github.com/criteo/data-aggregation-api/internal/model/cmdb/tacacs"
	"github.com/criteo/data-aggregation-api/internal/model/dcim"
)

// TestGenerateconfigsSystem checks that the NTP and the TACACS+ configurations,
// which both live under openconfig-system, are served side by side.
func TestGenerateconfigsSystem(t *testing.T) {
	tests := []struct {
		name     string
		ntp      *ntp.NTP
		tacacs   *tacacs.Tacacs
		wantKeys []string
	}{
		{
			name:     "NTP and TACACS+",
			ntp:      &ntp.NTP{ServerList: []ntp.Server{{Name: "ntp1", ServerAddress: "10.0.0.1"}}},
			tacacs:   &tacacs.Tacacs{Passkey: "secret", ServerList: []tacacs.Server{{ServerAddress: "10.10.10.10", Priority: 1, TCPPort: 49}}},
			wantKeys: []string{"aaa", "ntp"},
		},
		{
			name:     "NTP only",
			ntp:      &ntp.NTP{ServerList: []ntp.Server{{Name: "ntp1", ServerAddress: "10.0.0.1"}}},
			wantKeys: []string{"ntp"},
		},
		{
			name:     "TACACS+ only",
			tacacs:   &tacacs.Tacacs{ServerList: []tacacs.Server{{ServerAddress: "10.10.10.10"}}},
			wantKeys: []string{"aaa"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Device{
				mutex:           &sync.Mutex{},
				Dcim:            &dcim.NetworkDevice{Hostname: "test-device"},
				BGPGlobalConfig: &bgp.BGPGlobal{},
				NTP:             tt.ntp,
				Tacacs:          tt.tacacs,
			}

			if err := d.Generateconfigs(); err != nil {
				t.Fatalf("Generateconfigs() failed: %s", err)
			}

			var got struct {
				System map[string]json.RawMessage `json:"system"`
			}
			if err := json.Unmarshal([]byte(d.Config.JSONOpenConfig), &got); err != nil {
				t.Fatalf("invalid generated JSON: %s", err)
			}

			system, _ := json.MarshalIndent(got.System, "", "  ")
			t.Logf("system:\n%s", system)

			if len(got.System) != len(tt.wantKeys) {
				t.Errorf("system has %d containers, want %v", len(got.System), tt.wantKeys)
			}
			for _, key := range tt.wantKeys {
				if _, ok := got.System[key]; !ok {
					t.Errorf("system.%s is missing from the output", key)
				}
			}
		})
	}
}
