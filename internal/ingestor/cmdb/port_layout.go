package cmdb

import (
	"fmt"
	"net/url"

	"github.com/rs/zerolog/log"

	"github.com/criteo/data-aggregation-api/internal/ingestor/netbox"
	"github.com/criteo/data-aggregation-api/internal/model/cmdb/interfaces"
)

// GetPortLayouts returns all port layouts from the Network CMDB.
//
// Port layouts are attached to a hardware model and a network role, not to a
// device or a datacenter, so the whole table is fetched.
func GetPortLayouts() ([]*interfaces.PortLayout, error) {
	response := netbox.NetboxResponse[interfaces.PortLayout]{}

	err := netbox.Get("/api/plugins/cmdb/port-layouts/", &response, url.Values{})
	if err != nil {
		return nil, fmt.Errorf("port layouts fetching failure: %w", err)
	}

	if response.Count != len(response.Results) {
		log.Warn().Msg("some port layouts have not been fetched")
	}

	return response.Results, nil
}

// PrecomputePortLayouts groups the ports by layout (device type and network
// role) and indexes them by the name the CMDB device interfaces use: the
// logical name (e.g. to_L3-SP01), and the generic name for the ports whose
// device interface is named after the port itself (e.g. eth0).
//
// A logical name shared by several ports (e.g. the UNUSED placeholder) is
// ambiguous: there is no way to know which port is the right one. It is kept
// as a nil entry, which fails the devices whose interfaces use it (see
// ResolveInterfaces).
func PrecomputePortLayouts(portLayouts []*interfaces.PortLayout) map[interfaces.PortLayoutKey]interfaces.PortLayoutTable {
	var tables = make(map[interfaces.PortLayoutKey]interfaces.PortLayoutTable)
	for _, port := range portLayouts {
		key := port.Key()
		table, ok := tables[key]
		if !ok {
			table = make(interfaces.PortLayoutTable)
			tables[key] = table
		}

		if port.LogicalName != "" {
			if previous, duplicate := table[port.LogicalName]; duplicate {
				if previous != nil {
					log.Debug().Msgf("ports %s and %s of the port layout of device type %d, role %d share the logical name %s",
						previous.Name, port.Name, key.DeviceTypeID, key.RoleID, port.LogicalName)
				}
				table[port.LogicalName] = nil
				continue
			}
			table[port.LogicalName] = port
		}
	}

	// generic names never override a logical name
	for _, port := range portLayouts {
		table := tables[port.Key()]
		if _, taken := table[port.Name]; !taken {
			table[port.Name] = port
		}
	}

	return tables
}
