package interfaces

// ManagementInterfaceName is the name of the management interface of a
// device in the NOS, the one SONiC configures through MGMT_PORT and
// MGMT_INTERFACE. It is matched against Interface.Name, i.e. after the port
// layout translation when the interface has a layout entry.
const ManagementInterfaceName = "eth0"

// Interface is a CMDB device interface of one device, precomputed for the
// convertor: its names are resolved and its logical interfaces and management
// routes are attached to it.
type Interface struct {
	DeviceInterface *DeviceInterface
	// Name is the name of the interface on the device: see PortLayoutTable.PortName.
	Name string
	// Port is the port layout entry of the interface, nil when the interface
	// is not a front-panel port or when the device has no port layout.
	Port *PortLayout
	// Neighbor is the remote end of the link connected to the interface, nil
	// when there is none.
	Neighbor *Neighbor
	// LogicalInterfaces are in CMDB order, each with a unique index.
	LogicalInterfaces []*LogicalInterface
	// ManagementRoutes are indexed by logical interface index, in CMDB order.
	ManagementRoutes map[uint32][]*ManagementRoute
}
