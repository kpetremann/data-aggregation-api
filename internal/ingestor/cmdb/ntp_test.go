package cmdb_test

import (
	"encoding/json"
	"testing"

	"github.com/criteo/data-aggregation-api/internal/ingestor/cmdb"
	"github.com/criteo/data-aggregation-api/internal/model/cmdb/ntp"
	"github.com/google/go-cmp/cmp"
)

// ntpConfig is an NTP configuration as Netbox returns it, Device being an anonymous struct.
func ntpConfig(hostname string, servers ...ntp.Server) *ntp.NTP {
	config := &ntp.NTP{ServerList: servers}
	config.Device.Name = hostname
	return config
}

func TestPrecomputeNTP(t *testing.T) {
	tests := []struct {
		name string
		args string
		want map[string]*ntp.NTP
	}{
		{
			name: "devices sharing servers",
			args: `[
         {
            "id":1,
            "device":{"id":1, "name":"tor01-01"},
            "server_list":[
               {"id":1, "name":"ntp1", "server_address":"192.0.2.1", "created":"2026-10-07T07:37:29Z", "last_updated":"2026-10-07T07:37:29Z"},
               {"id":2, "name":"ntp2", "server_address":"192.0.2.2", "created":"2026-10-07T07:37:29Z", "last_updated":"2026-10-07T07:37:29Z"}
            ],
            "created":"2026-10-07T07:37:29Z",
            "last_updated":"2026-10-07T07:37:29Z"
         },
         {
            "id":2,
            "device":{"id":2, "name":"tor01-02"},
            "server_list":[
               {"id":1, "name":"ntp1", "server_address":"192.0.2.1", "created":"2026-10-07T07:37:29Z", "last_updated":"2026-10-07T07:37:29Z"}
            ],
            "created":"2026-10-07T07:37:29Z",
            "last_updated":"2026-10-07T07:37:29Z"
         }
      ]`,
			want: map[string]*ntp.NTP{
				"tor01-01": ntpConfig("tor01-01",
					ntp.Server{Name: "ntp1", ServerAddress: "192.0.2.1"},
					ntp.Server{Name: "ntp2", ServerAddress: "192.0.2.2"},
				),
				"tor01-02": ntpConfig("tor01-02",
					ntp.Server{Name: "ntp1", ServerAddress: "192.0.2.1"},
				),
			},
		},
	}

	for _, test := range tests {
		var cmdbOutput []*ntp.NTP
		if err := json.Unmarshal([]byte(test.args), &cmdbOutput); err != nil {
			t.Errorf("unable to load test data for '%s': %s", test.name, err)
			continue
		}

		out := cmdb.PrecomputeNTP(cmdbOutput)
		if diff := cmp.Diff(out, test.want); diff != "" {
			t.Errorf("unexpected diff for '%s': %s\n", test.name, diff)
		}
	}
}
