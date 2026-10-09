package cmdb_test

import (
	"encoding/json"
	"testing"

	"github.com/criteo/data-aggregation-api/internal/ingestor/cmdb"
	"github.com/criteo/data-aggregation-api/internal/model/cmdb/tacacs"
	"github.com/google/go-cmp/cmp"
)

func TestPrecomputeTacacs(t *testing.T) {
	tests := []struct {
		name string
		args string
		want map[string]*tacacs.Tacacs
	}{
		{
			name: "one configuration with two servers",
			args: `[
         {
            "id":1,
            "device":{
               "id":1,
               "name":"tor01-01"
            },
            "passkey":"Mypasword1",
            "server_list":[
               {
                  "id":1,
                  "server_address":"1.1.1.1",
                  "priority":1,
                  "tcp_port":49
               },
               {
                  "id":2,
                  "server_address":"1.1.1.2",
                  "priority":2,
                  "tcp_port":4949
               }
            ],
            "created":"2026-08-10T07:37:29.238195Z",
            "last_updated":"2026-08-11T13:08:11.378498Z"
         }
      ]`,
			want: map[string]*tacacs.Tacacs{
				"tor01-01": {
					Device: struct {
						Name string "json:\"name\" validate:\"required\""
					}{
						Name: "tor01-01",
					},
					Passkey: "Mypasword1",
					ServerList: []tacacs.Server{
						{
							ServerAddress: "1.1.1.1",
							Priority:      1,
							TCPPort:       49,
						},
						{
							ServerAddress: "1.1.1.2",
							Priority:      2,
							TCPPort:       4949,
						},
					},
				},
			},
		},
	}

	for _, test := range tests {
		var cmdbOutput []*tacacs.Tacacs
		if err := json.Unmarshal([]byte(test.args), &cmdbOutput); err != nil {
			t.Errorf("unable to load test data for '%s': %s", test.name, err)
			continue
		}

		out := cmdb.PrecomputeTacacs(cmdbOutput)
		if diff := cmp.Diff(out, test.want); diff != "" {
			t.Errorf("unexpected diff for '%s': %s\n", test.name, diff)
		}
	}
}
