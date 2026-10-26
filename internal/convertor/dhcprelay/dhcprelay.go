package dhcprelay

import (
	"fmt"
	"sort"
	"strings"

	"github.com/criteo/data-aggregation-api/internal/model/cmdb/dhcprelay"
	cmdbif "github.com/criteo/data-aggregation-api/internal/model/cmdb/interfaces"
	"github.com/criteo/data-aggregation-api/internal/model/openconfig"
	"github.com/openconfig/ygot/ygot"
)

// DHCPRelayToOpenConfig converts the CMDB DHCP relay interfaces of a device
// into the dhcp container of openconfig-relay-agent.
//
// Each relay is keyed by the name on the device of the parent interface of its
// logical interface (see PortLayoutTable.PortName): SONiC enables the relay on
// the interface itself, so only the logical unit 0 can carry it. This covers
// both routed ports and SVIs, as openconfig-relay-agent only identifies an
// interface by name: the VRF is the one of the interface, not a property of the
// relay.
//
// interface-ref is not set on purpose: it is a leafref, which makes the
// validation fail unless the referenced interface is emitted in /interfaces
// too, and AFK does not generate the interfaces yet.
//
// Relays without any server are skipped. It returns nil when there is no relay
// to advertise. Servers are sorted to keep the emitted JSON stable.
//
// portLayout may be nil: the CMDB names are kept.
func DHCPRelayToOpenConfig(relays []*dhcprelay.Relay, portLayout cmdbif.PortLayoutTable) (*openconfig.RelayAgent, error) {
	dhcp := &openconfig.RelayAgent_Dhcp{}

	for _, relay := range relays {
		if relay == nil || len(relay.ServerList) == 0 {
			continue
		}

		helpers := make([]string, 0, len(relay.ServerList))
		for _, server := range relay.ServerList {
			helpers = append(helpers, server.ServerAddress)
		}
		sort.Strings(helpers)

		name, err := relayInterfaceName(relay.LogicalInterface, portLayout)
		if err != nil {
			return nil, err
		}

		iface, err := dhcp.NewInterface(name)
		if err != nil {
			return nil, fmt.Errorf("failed to add the DHCP relay on %s: %w", relay.LogicalInterface.Name, err)
		}

		iface.Enable = ygot.Bool(true)
		iface.HelperAddress = helpers
	}

	if len(dhcp.Interface) == 0 {
		return nil, nil
	}

	return &openconfig.RelayAgent{Dhcp: dhcp}, nil
}

// relayInterfaceName returns the name on the device of the parent interface
// of the logical interface of a relay (e.g. to_spine_01.0 gives Ethernet0).
// The logical interface must be the unit 0 of its parent.
func relayInterfaceName(logicalInterface dhcprelay.LogicalInterface, portLayout cmdbif.PortLayoutTable) (string, error) {
	parent := logicalInterface.ParentInterface.Name

	// the unit is read from the name, as the index is optional
	unit, ok := strings.CutPrefix(logicalInterface.Name, parent)
	if !ok {
		return "", fmt.Errorf("DHCP relay logical interface %s is not named after its parent interface %s", logicalInterface.Name, parent)
	}
	if unit != "" && unit != ".0" {
		return "", fmt.Errorf("DHCP relay on logical interface %s: only the logical unit 0 can carry it", logicalInterface.Name)
	}

	if portLayout.IsAmbiguous(parent) {
		return "", fmt.Errorf("parent interface %s of DHCP relay logical interface %s matches several ports of the port layout", parent, logicalInterface.Name)
	}

	return portLayout.PortName(parent), nil
}
