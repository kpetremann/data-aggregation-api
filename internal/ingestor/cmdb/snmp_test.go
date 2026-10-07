package cmdb_test

import (
	"encoding/json"
	"testing"

	"github.com/criteo/data-aggregation-api/internal/ingestor/cmdb"
	"github.com/criteo/data-aggregation-api/internal/model/cmdb/snmp"
	"github.com/google/go-cmp/cmp"
)

// snmpConfig is an SNMP configuration as Netbox returns it, Device being an anonymous struct.
func snmpConfig(hostname, location, contact string, communities ...snmp.Community) *snmp.SNMP {
	config := &snmp.SNMP{Location: location, Contact: contact, CommunityList: communities}
	config.Device.Name = hostname
	return config
}

func TestPrecomputeSNMP(t *testing.T) {
	tests := []struct {
		name string
		args string
		want map[string]*snmp.SNMP
	}{
		{
			name: "lite configuration",
			args: `[
         {
            "id":1,
            "device":{
               "id":1,
               "name":"tor01-01"
            },
            "community_list":[
               {
                  "community":"my_community_read",
                  "type":"readonly"
               },
               {
                  "community":"my_community_write",
                  "type":"readwrite"
               }
            ],
            "created":"2024-06-10T07:37:29.238195Z",
            "last_updated":"2024-06-11T13:08:11.378498Z",
            "location":"location1",
            "contact":"best team"
         }
      ]`,
			want: map[string]*snmp.SNMP{
				"tor01-01": snmpConfig("tor01-01", "location1", "best team",
					snmp.Community{Community: "my_community_read", Type: "readonly"},
					snmp.Community{Community: "my_community_write", Type: "readwrite"},
				),
			},
		},
	}

	for _, test := range tests {
		var cmdbOutput []*snmp.SNMP
		if err := json.Unmarshal([]byte(test.args), &cmdbOutput); err != nil {
			t.Errorf("unable to load test data for '%s': %s", test.name, err)
			continue
		}

		out := cmdb.PrecomputeSNMP(cmdbOutput)
		if diff := cmp.Diff(out, test.want); diff != "" {
			t.Errorf("unexpected diff for '%s': %s\n", test.name, diff)
		}
	}
}
