package dhcprelay

import (
	"fmt"
	"sort"

	"github.com/criteo/data-aggregation-api/internal/model/cmdb/dhcprelay"
	"github.com/criteo/data-aggregation-api/internal/model/openconfig"
	"github.com/openconfig/ygot/ygot"
)

// DHCPRelayToOpenConfig converts the CMDB DHCP relay interfaces of a device
// into the dhcp container of openconfig-relay-agent.
//
// Each relay is keyed by the name of its logical interface. This covers both
// routed ports and SVIs, as openconfig-relay-agent only identifies an interface
// by name: the VRF is the one of the interface, not a property of the relay.
//
// interface-ref is not set on purpose: it is a leafref, which makes the
// validation fail unless the referenced interface is emitted in /interfaces
// too, and AFK does not generate the interfaces yet.
//
// Relays without any server are skipped. It returns nil when there is no relay
// to advertise. Servers are sorted to keep the emitted JSON stable.
func DHCPRelayToOpenConfig(relays []*dhcprelay.Relay) (*openconfig.RelayAgent, error) {
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

		iface, err := dhcp.NewInterface(relay.LogicalInterface.Name)
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
