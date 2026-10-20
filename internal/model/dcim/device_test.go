package dcim_test

import (
	"encoding/json"
	"testing"

	"github.com/go-playground/validator/v10"

	"github.com/criteo/data-aggregation-api/internal/ingestor/netbox"
	"github.com/criteo/data-aggregation-api/internal/model/dcim"
)

// validate checks a /api/dcim/devices/ response the way netbox.Get does.
func validate(t *testing.T, results string) error {
	t.Helper()
	var response netbox.NetboxResponse[dcim.NetworkDevice]
	if err := json.Unmarshal([]byte(`{"results": `+results+`}`), &response); err != nil {
		t.Fatal(err)
	}
	return validator.New().Struct(&response)
}

// TestNetworkDevice_TypeAndRoleAreOptional: a device without them stays in
// the inventory; it only lacks its DEVICE_METADATA.
func TestNetworkDevice_TypeAndRoleAreOptional(t *testing.T) {
	for name, results := range map[string]string{
		"complete": `[{"name": "tor01-01",
		  "device_type": {"id": 10, "model": "MSN2700"},
		  "device_role": {"id": 20, "name": "ToR"}}]`,
		"no device type": `[{"name": "tor01-01", "device_role": {"id": 20, "name": "ToR"}}]`,
		"no device role": `[{"name": "tor01-01", "device_type": {"id": 10, "model": "MSN2700"}}]`,
	} {
		if err := validate(t, results); err != nil {
			t.Errorf("%s: the device was rejected: %v", name, err)
		}
	}
}

// TestNetworkDevice_Role: the role is read from "role" (NetBox >= 3.6) and
// from "device_role" (NetBox < 3.6), "role" first.
func TestNetworkDevice_Role(t *testing.T) {
	for name, test := range map[string]struct {
		device string
		want   dcim.Role
	}{
		"role":        {`{"name": "tor01-01", "role": {"id": 20, "name": "ToR"}}`, dcim.Role{ID: 20, Name: "ToR"}},
		"device_role": {`{"name": "tor01-01", "device_role": {"id": 21, "name": "Spine"}}`, dcim.Role{ID: 21, Name: "Spine"}},
		"both": {
			`{"name": "tor01-01", "role": {"id": 20, "name": "ToR"}, "device_role": {"id": 21, "name": "Spine"}}`,
			dcim.Role{ID: 20, Name: "ToR"},
		},
		"none": {`{"name": "tor01-01"}`, dcim.Role{}},
	} {
		var device dcim.NetworkDevice
		if err := json.Unmarshal([]byte(test.device), &device); err != nil {
			t.Fatalf("%s: %v", name, err)
		}
		if device.Hostname != "tor01-01" {
			t.Errorf("%s: hostname not decoded: %q", name, device.Hostname)
		}
		if device.Role != test.want {
			t.Errorf("%s: got role %+v, want %+v", name, device.Role, test.want)
		}
	}
}
