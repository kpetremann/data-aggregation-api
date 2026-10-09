package tacacs

// Server is one entry of the TacacsServer table of the Network CMDB.
//
// The server with the highest Priority is the preferred one. TCPPort is a uint16
// to match the inet:port-number of the YANG model it is converted to.
type Server struct {
	ServerAddress string `json:"server_address" validate:"required,ip"`
	Priority      uint32 `json:"priority"       validate:"omitempty"`
	TCPPort       uint16 `json:"tcp_port"       validate:"omitempty"`
}

// Tacacs is the TACACS configuration of a single device: its global settings and its servers.
type Tacacs struct {
	Device struct {
		Name string `json:"name" validate:"required"`
	} `json:"device" validate:"required"`
	Passkey    string   `json:"passkey"     validate:"omitempty"`
	ServerList []Server `json:"server_list" validate:"omitempty,dive"`
}
