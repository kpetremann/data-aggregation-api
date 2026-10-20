package dcim

import "encoding/json"

// Role is a NetBox device role.
type Role struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type NetworkDevice struct {
	Hostname     string `json:"name" validate:"required"`
	SerialNumber string `json:"serial" validate:"omitempty"`
	Tags         []struct {
		Name string `json:"name" validate:"required"`
	} `json:"tags" validate:"omitempty"`
	// DeviceType and Role select the SONiC HwSKU and device type of
	// DEVICE_METADATA, through the CMDB mappings, and identify the port layout
	// of the device. Not validated: a device without them only lacks its
	// DEVICE_METADATA and interfaces, not its inventory entry.
	DeviceType struct {
		ID    int    `json:"id"`
		Model string `json:"model"`
	} `json:"device_type"`
	// Role is read from "role" (NetBox >= 3.6) or "device_role" (NetBox < 3.6),
	// see UnmarshalJSON.
	Role Role `json:"role"`
}

// UnmarshalJSON decodes a NetBox device, taking the role from "device_role"
// when "role" is absent.
// TODO: drop "device_role" once NetBox < 3.6 is not supported anymore.
func (d *NetworkDevice) UnmarshalJSON(data []byte) error {
	type networkDevice NetworkDevice
	var raw struct {
		networkDevice
		DeviceRole *Role `json:"device_role"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	*d = NetworkDevice(raw.networkDevice)
	if d.Role.ID == 0 && raw.DeviceRole != nil {
		d.Role = *raw.DeviceRole
	}
	return nil
}
