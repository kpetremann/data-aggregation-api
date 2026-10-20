package interfaces

// LinkEndpoint is one end of a CMDB link: a device interface.
type LinkEndpoint struct {
	Name   string `json:"name" validate:"required"`
	Device struct {
		Name string `json:"name" validate:"required"`
	} `json:"device" validate:"required"`
}

// Link is a point-to-point link between two device interfaces in the Network CMDB.
//
// CMDB endpoint: /api/plugins/cmdb/links/.
type Link struct {
	InterfaceA LinkEndpoint `json:"interface_a" validate:"required"`
	InterfaceB LinkEndpoint `json:"interface_b" validate:"required"`
}

// Neighbor is the remote end of the link connected to a device interface.
type Neighbor struct {
	// Device is the hostname of the remote device.
	Device string
	// Port is the name of the remote interface on the remote device: see PortName.
	Port string
}

// NeighborTable is the neighbors of one device, indexed by the CMDB name of
// the local device interface.
type NeighborTable map[string]*Neighbor
