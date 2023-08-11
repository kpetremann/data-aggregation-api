package cmdb_test

import (
	"encoding/json"
	"testing"

	"github.com/criteo/data-aggregation-api/internal/ingestor/cmdb"
	"github.com/criteo/data-aggregation-api/internal/model/cmdb/routingpolicy"
	"github.com/google/go-cmp/cmp"
)

func TestPrecomputeCommunityLists(t *testing.T) {
	tests := []struct {
		name string
		args string
		want map[string][]*routingpolicy.CommunityList
	}{
		{
			name: "valid community list",
			args: `
			[
				{
					"id": 1,
					"device": {
						"id": 1,
						"name": "tor01-01"
					},
					"terms": [
						{
							"community": "650..:999"
						},
						{
							"community": "650..:1000"
						}
					],
					"created": "2023-06-20T12:23:50.955067Z",
					"last_updated": "2023-06-20T12:23:50.955078Z",
					"name": "SERVERS"
				}
			]
			`,
			want: map[string][]*routingpolicy.CommunityList{
				"tor01-01": {
					&routingpolicy.CommunityList{
						Name: "SERVERS",
						Device: struct {
							Name string `json:"name" validate:"required"`
						}{
							Name: "tor01-01",
						},
						Terms: []*routingpolicy.CommunityListTerm{
							{
								Community: "650..:999",
							},
							{
								Community: "650..:1000",
							},
						},
					},
				},
			},
		},
	}

	for _, test := range tests {
		var cmdbOutput []*routingpolicy.CommunityList
		if err := json.Unmarshal([]byte(test.args), &cmdbOutput); err != nil {
			t.Errorf("unable to load test data for '%s': %s", test.name, err)
			continue
		}

		out := cmdb.PrecomputeCommunityLists(cmdbOutput)
		if diff := cmp.Diff(out, test.want); diff != "" {
			t.Errorf("unexpected diff for '%s': %s\n", test.name, diff)
		}
	}
}
