package dhcprelay

// Server is one DHCP server of the Network CMDB.
type Server struct {
	ServerAddress string `json:"server_address" validate:"required,ip"`
}

// LogicalInterface is the interface a relay is enabled on: a sub-interface
// of a parent interface of a device. Index is a uint32 to match the
// subinterface index of the YANG model it is converted to.
type LogicalInterface struct {
	Name            string `json:"name"  validate:"required"`
	Index           uint32 `json:"index" validate:"omitempty"`
	ParentInterface struct {
		Name   string `json:"name" validate:"required"`
		Device struct {
			Name string `json:"name" validate:"required"`
		} `json:"device" validate:"required"`
	} `json:"parent_interface" validate:"required"`
}

// Relay is one entry of the DhcpRelayInterface table of the Network CMDB:
// the DHCP servers a logical interface relays to.
type Relay struct {
	LogicalInterface LogicalInterface `json:"logical_interface" validate:"required"`
	ServerList       []Server         `json:"server_list"       validate:"omitempty,dive"`
}
