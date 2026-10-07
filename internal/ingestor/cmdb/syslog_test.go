package cmdb_test

import (
	"encoding/json"
	"testing"

	"github.com/criteo/data-aggregation-api/internal/ingestor/cmdb"
	"github.com/criteo/data-aggregation-api/internal/model/cmdb/syslog"
	"github.com/google/go-cmp/cmp"
)

// syslogConfig is a syslog configuration as Netbox returns it, Device being an anonymous struct.
func syslogConfig(hostname string, servers ...syslog.Server) *syslog.Syslog {
	config := &syslog.Syslog{ServerList: servers}
	config.Device.Name = hostname
	return config
}

func TestPrecomputeSyslog(t *testing.T) {
	tests := []struct {
		name string
		args string
		want map[string]*syslog.Syslog
	}{
		{
			name: "devices sharing servers",
			args: `[
         {
            "id":1,
            "device":{"id":1, "name":"ra01-01-p01-da1-pnet.crto.io"},
            "server_list":[
               {"id":1, "server_address":"172.30.72.5", "created":"2026-10-06T16:23:11Z", "last_updated":"2026-10-06T16:23:11Z"},
               {"id":2, "server_address":"172.30.72.6", "created":"2026-10-06T16:23:11Z", "last_updated":"2026-10-06T16:23:11Z"}
            ],
            "created":"2026-10-06T16:23:11Z",
            "last_updated":"2026-10-06T16:23:11Z"
         },
         {
            "id":2,
            "device":{"id":2, "name":"sp01-01-p01-da1-pnet.crto.io"},
            "server_list":[
               {"id":1, "server_address":"172.30.72.5", "created":"2026-10-06T16:23:11Z", "last_updated":"2026-10-06T16:23:11Z"}
            ],
            "created":"2026-10-06T16:23:11Z",
            "last_updated":"2026-10-06T16:23:11Z"
         }
      ]`,
			want: map[string]*syslog.Syslog{
				"ra01-01-p01-da1-pnet.crto.io": syslogConfig("ra01-01-p01-da1-pnet.crto.io",
					syslog.Server{ServerAddress: "172.30.72.5"},
					syslog.Server{ServerAddress: "172.30.72.6"},
				),
				"sp01-01-p01-da1-pnet.crto.io": syslogConfig("sp01-01-p01-da1-pnet.crto.io",
					syslog.Server{ServerAddress: "172.30.72.5"},
				),
			},
		},
		{
			name: "no server",
			args: `[
         {
            "id":1,
            "device":{"id":1, "name":"ra01-01-p01-da1-pnet.crto.io"},
            "server_list":[],
            "created":"2026-10-06T16:23:11Z",
            "last_updated":"2026-10-06T16:23:11Z"
         }
      ]`,
			want: map[string]*syslog.Syslog{
				"ra01-01-p01-da1-pnet.crto.io": syslogConfig("ra01-01-p01-da1-pnet.crto.io", []syslog.Server{}...),
			},
		},
	}

	for _, test := range tests {
		var cmdbOutput []*syslog.Syslog
		if err := json.Unmarshal([]byte(test.args), &cmdbOutput); err != nil {
			t.Errorf("unable to load test data for '%s': %s", test.name, err)
			continue
		}

		out := cmdb.PrecomputeSyslog(cmdbOutput)
		if diff := cmp.Diff(out, test.want); diff != "" {
			t.Errorf("unexpected diff for '%s': %s\n", test.name, diff)
		}
	}
}
