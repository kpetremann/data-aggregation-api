package cmdb

import (
	"fmt"

	"github.com/criteo/data-aggregation-api/internal/model/cmdb/interfaces"
)

// ResolveInterfaces precomputes the interfaces of one device for the
// convertor: each device interface gets its name on the device and its port
// layout entry, its neighbor, its logical interfaces and its management routes.
//
// It fails on CMDB data no configuration can be generated from: duplicate,
// dangling or ambiguous references, and interfaces sharing a port.
//
// portLayout and neighbors may be nil. The result is in CMDB order.
func ResolveInterfaces(hostname string, deviceInterfaces []*interfaces.DeviceInterface, logicalInterfaces []*interfaces.LogicalInterface,
	managementRoutes []*interfaces.ManagementRoute, portLayout interfaces.PortLayoutTable, neighbors interfaces.NeighborTable,
) ([]*interfaces.Interface, error) {
	names, err := resolveNames(hostname, deviceInterfaces, portLayout)
	if err != nil {
		return nil, err
	}

	logicalPerParent, err := groupLogicalInterfaces(hostname, logicalInterfaces, names)
	if err != nil {
		return nil, err
	}

	routesPerParent, err := groupManagementRoutes(hostname, managementRoutes, names, logicalPerParent)
	if err != nil {
		return nil, err
	}

	resolved := make([]*interfaces.Interface, 0, len(deviceInterfaces))
	for _, deviceInterface := range deviceInterfaces {
		resolved = append(resolved, &interfaces.Interface{
			DeviceInterface:   deviceInterface,
			Name:              names[deviceInterface.Name],
			Port:              portLayout[deviceInterface.Name],
			Neighbor:          neighbors[deviceInterface.Name],
			LogicalInterfaces: logicalPerParent[deviceInterface.Name],
			ManagementRoutes:  routesPerParent[deviceInterface.Name],
		})
	}

	return resolved, nil
}

// resolveNames returns the name on the device of each device interface, keyed
// by CMDB name.
func resolveNames(hostname string, deviceInterfaces []*interfaces.DeviceInterface, portLayout interfaces.PortLayoutTable) (map[string]string, error) {
	names := make(map[string]string, len(deviceInterfaces))
	// CMDB name per name on the device
	cmdbNames := make(map[string]string, len(deviceInterfaces))

	for _, deviceInterface := range deviceInterfaces {
		if _, duplicate := names[deviceInterface.Name]; duplicate {
			return nil, fmt.Errorf("duplicate device interface %s on %s", deviceInterface.Name, hostname)
		}

		if portLayout.IsAmbiguous(deviceInterface.Name) {
			return nil, fmt.Errorf("device interface %s of %s matches several ports of the port layout", deviceInterface.Name, hostname)
		}

		name := portLayout.PortName(deviceInterface.Name)
		if other, duplicate := cmdbNames[name]; duplicate {
			return nil, fmt.Errorf("device interfaces %s and %s of %s both map to port %s", deviceInterface.Name, other, hostname, name)
		}

		names[deviceInterface.Name] = name
		cmdbNames[name] = deviceInterface.Name
	}

	return names, nil
}

// groupLogicalInterfaces groups the logical interfaces by CMDB name of their
// parent, which must be in names, keeping the CMDB order.
func groupLogicalInterfaces(hostname string, logicalInterfaces []*interfaces.LogicalInterface, names map[string]string) (map[string][]*interfaces.LogicalInterface, error) {
	perParent := make(map[string][]*interfaces.LogicalInterface)
	type unit struct {
		parent string
		index  uint32
	}
	seen := make(map[unit]struct{}, len(logicalInterfaces))

	for _, logicalInterface := range logicalInterfaces {
		name := logicalInterface.Name()
		parent := logicalInterface.ParentInterface.Name

		if logicalInterface.Index == nil {
			return nil, fmt.Errorf("missing index on logical interface of %s on %s", parent, hostname)
		}
		index := *logicalInterface.Index

		if _, ok := names[parent]; !ok {
			return nil, fmt.Errorf("logical interface %s of %s references unknown device interface %s", name, hostname, parent)
		}

		if _, duplicate := seen[unit{parent, index}]; duplicate {
			return nil, fmt.Errorf("duplicate logical interface %s on %s", name, hostname)
		}
		seen[unit{parent, index}] = struct{}{}

		// SONiC configures it on the interface itself: only the unit 0 can carry it
		if index != 0 && logicalInterface.UseIPv6LinkLocalOnly != nil && *logicalInterface.UseIPv6LinkLocalOnly {
			return nil, fmt.Errorf("use_ipv6_link_local_only is set on logical interface %s of %s: only the logical unit 0 can carry it", name, hostname)
		}

		perParent[parent] = append(perParent[parent], logicalInterface)
	}

	return perParent, nil
}

// groupManagementRoutes groups the management routes by CMDB name of their
// device interface, then by logical interface index, keeping the CMDB order.
// Their logical interface must belong to the management interface and carry
// an IPv4 address.
func groupManagementRoutes(hostname string, managementRoutes []*interfaces.ManagementRoute, names map[string]string,
	logicalPerParent map[string][]*interfaces.LogicalInterface,
) (map[string]map[uint32][]*interfaces.ManagementRoute, error) {
	perParent := make(map[string]map[uint32][]*interfaces.ManagementRoute)

	for _, route := range managementRoutes {
		parent := route.InterfaceName()
		index := route.SubinterfaceIndex()
		name := fmt.Sprintf("%s.%d", parent, index)

		if route.Kind != interfaces.ManagementRouteStatic && route.Kind != interfaces.ManagementRouteForced {
			return nil, fmt.Errorf("unsupported management route kind %q for %s on %s", route.Kind, route.Prefix.String(), hostname)
		}

		parentName, ok := names[parent]
		if !ok {
			return nil, fmt.Errorf("management route %s on %s references unknown device interface %s", route.Prefix.String(), hostname, parent)
		}
		if parentName != interfaces.ManagementInterfaceName {
			return nil, fmt.Errorf("management route %s on %s is attached to %s, which is not the management interface", route.Prefix.String(), hostname, name)
		}

		unit := logicalInterfaceByIndex(logicalPerParent[parent], index)
		if unit == nil {
			return nil, fmt.Errorf("management route %s on %s references unknown logical interface %s", route.Prefix.String(), hostname, name)
		}
		if unit.IPv4Address == nil {
			return nil, fmt.Errorf("management route %s on %s: no IPv4 address on %s", route.Prefix.String(), hostname, name)
		}

		if perParent[parent] == nil {
			perParent[parent] = make(map[uint32][]*interfaces.ManagementRoute)
		}
		perParent[parent][index] = append(perParent[parent][index], route)
	}

	return perParent, nil
}

// logicalInterfaceByIndex returns the logical interface of the given index, nil if none.
func logicalInterfaceByIndex(logicalInterfaces []*interfaces.LogicalInterface, index uint32) *interfaces.LogicalInterface {
	for _, logicalInterface := range logicalInterfaces {
		if *logicalInterface.Index == index {
			return logicalInterface
		}
	}
	return nil
}
