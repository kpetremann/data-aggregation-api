package interfaces

import (
	"fmt"
	"net"

	cmdbif "github.com/criteo/data-aggregation-api/internal/model/cmdb/interfaces"
	"github.com/criteo/data-aggregation-api/internal/model/openconfig"
	"github.com/rs/zerolog/log"
)

// StaticProtocolName is the name of the static routing protocol instance.
// OpenConfig path: /network-instances/network-instance/protocols/protocol[STATIC][static].
const StaticProtocolName = "static"

// StaticProtocolKey is the key of the static routing protocol instance of a
// network instance.
var StaticProtocolKey = openconfig.NetworkInstance_Protocol_Key{
	Identifier: openconfig.PolicyTypes_INSTALL_PROTOCOL_TYPE_STATIC,
	Name:       StaticProtocolName,
}

// NetworkInstanceProtocols maps a network instance name to its routing
// protocol instances.
// OpenConfig path: /network-instances/network-instance/protocols/protocol/.
type NetworkInstanceProtocols map[string]map[openconfig.NetworkInstance_Protocol_Key]*openconfig.NetworkInstance_Protocol

// ManagementRoutesToOpenconfig converts the static management routes of the
// precomputed interfaces of one device to static routes of the network
// instance their logical interface is bound to, with the logical interface as
// interface-ref. Forced routes are emitted with their subinterface, see
// SubinterfaceToOpenconfig.
//
// The consumer contract for the management default route is "the 0.0.0.0/0
// static route whose interface-ref is the management interface".
//
// It returns the routing protocol instances per network instance.
func ManagementRoutesToOpenconfig(hostname string, ifaces []*cmdbif.Interface) (NetworkInstanceProtocols, error) {
	protocols := make(NetworkInstanceProtocols)

	for _, iface := range ifaces {
		for _, logicalInterface := range iface.LogicalInterfaces {
			routes := iface.ManagementRoutes[*logicalInterface.Index]
			if len(routes) == 0 {
				continue
			}
			name := logicalInterface.Name()

			var hasStatic, hasForced bool
			for _, route := range routes {
				if route.Kind == cmdbif.ManagementRouteForced {
					hasForced = true
					continue
				}
				hasStatic = true

				networkInstance, ok := networkInstanceOf(logicalInterface)
				if !ok {
					return nil, fmt.Errorf("static route on %s of %s: the interface is bound to no network instance", name, hostname)
				}
				prefix, static, err := staticRouteToOpenconfig(hostname, name, iface.Name, logicalInterface, route)
				if err != nil {
					return nil, err
				}

				if protocols[networkInstance] == nil {
					protocols[networkInstance] = make(map[openconfig.NetworkInstance_Protocol_Key]*openconfig.NetworkInstance_Protocol)
				}
				protocol, ok := protocols[networkInstance][StaticProtocolKey]
				if !ok {
					protocolName := StaticProtocolName
					protocol = &openconfig.NetworkInstance_Protocol{
						Identifier: openconfig.PolicyTypes_INSTALL_PROTOCOL_TYPE_STATIC,
						Name:       &protocolName,
						Static:     make(map[string]*openconfig.NetworkInstance_Protocol_Static),
					}
					protocols[networkInstance][StaticProtocolKey] = protocol
				}
				if _, duplicate := protocol.Static[prefix]; duplicate {
					return nil, fmt.Errorf("duplicate static route %s in network instance %s of %s", prefix, networkInstance, hostname)
				}
				protocol.Static[prefix] = static
			}

			if hasForced && !hasStatic {
				log.Warn().Msgf("forced routes on %s of %s without a static route: they resolve through nothing", name, hostname)
			}
		}
	}

	return protocols, nil
}

// ipv4Prefix returns the route prefix as a string, or an error when it is not IPv4.
func ipv4Prefix(hostname string, route *cmdbif.ManagementRoute) (string, error) {
	if route.Prefix.IP.To4() == nil || route.Prefix.Netmask < 0 || route.Prefix.Netmask > 32 {
		return "", fmt.Errorf("management route %s on %s: IPv6 management routes are not supported", route.Prefix.String(), hostname)
	}
	return route.Prefix.String(), nil
}

// forcedRoutesToOpenconfig returns the prefixes of the forced routes among the
// management routes of the logical interface name, in CMDB order and without
// duplicates.
// OpenConfig path: /interfaces/interface/subinterfaces/subinterface/ipv4/config/forced-routes.
func forcedRoutesToOpenconfig(hostname string, name string, managementRoutes []*cmdbif.ManagementRoute) ([]string, error) {
	var forcedRoutes []string
	seen := make(map[string]struct{})

	for _, route := range managementRoutes {
		if route.Kind != cmdbif.ManagementRouteForced {
			continue
		}

		prefix, err := ipv4Prefix(hostname, route)
		if err != nil {
			return nil, err
		}
		if route.NextHop != nil {
			return nil, fmt.Errorf("forced route %s on %s of %s has a next-hop: forced routes resolve through the static route", prefix, name, hostname)
		}

		if _, duplicate := seen[prefix]; duplicate {
			log.Warn().Msgf("duplicate forced route %s on %s of %s", prefix, name, hostname)
			continue
		}
		seen[prefix] = struct{}{}
		forcedRoutes = append(forcedRoutes, prefix)
	}

	return forcedRoutes, nil
}

// staticRouteToOpenconfig converts a static management route of the logical
// interface name, with the logical interface as interface-ref. The next-hop
// must be an IPv4 address of the subnet of the logical interface, other than
// its own address. parentName is the name of the parent interface on the device.
// OpenConfig path: /network-instances/network-instance/protocols/protocol/static-routes/static/.
func staticRouteToOpenconfig(hostname string, name string, parentName string, logicalInterface *cmdbif.LogicalInterface, route *cmdbif.ManagementRoute) (string, *openconfig.NetworkInstance_Protocol_Static, error) {
	prefix, err := ipv4Prefix(hostname, route)
	if err != nil {
		return "", nil, err
	}

	if route.NextHop == nil {
		return "", nil, fmt.Errorf("static route %s on %s of %s has no next-hop", prefix, name, hostname)
	}
	nextHop := route.NextHop.To4()
	if nextHop == nil {
		return "", nil, fmt.Errorf("static route %s on %s of %s: next-hop %s is not an IPv4 address", prefix, name, hostname, route.NextHop.String())
	}
	if err := checkNextHopInSubnet(hostname, prefix, name, nextHop, logicalInterface); err != nil {
		return "", nil, err
	}

	// A single next-hop per static route: the management gateway.
	nextHopIndex := "0"
	return prefix, &openconfig.NetworkInstance_Protocol_Static{
		Prefix: &prefix,
		NextHop: map[string]*openconfig.NetworkInstance_Protocol_Static_NextHop{
			nextHopIndex: {
				Index:   &nextHopIndex,
				NextHop: openconfig.UnionString(nextHop.String()),
				InterfaceRef: &openconfig.NetworkInstance_Protocol_Static_NextHop_InterfaceRef{
					Interface:    &parentName,
					Subinterface: logicalInterface.Index,
				},
			},
		},
	}, nil
}

// checkNextHopInSubnet verifies the next-hop is reachable on the logical
// interface name: in the subnet of its IPv4 address, and not that address.
func checkNextHopInSubnet(hostname string, prefix string, name string, nextHop net.IP, logicalInterface *cmdbif.LogicalInterface) error {
	if address := logicalInterface.IPv4Address; address != nil {
		if ip := address.Address.IP.To4(); ip != nil {
			if ip.Equal(nextHop) {
				return fmt.Errorf("static route %s on %s of %s: next-hop %s is the address of the interface", prefix, name, hostname, nextHop.String())
			}
			subnet := net.IPNet{IP: ip, Mask: net.CIDRMask(address.Address.Netmask, 32)}
			if subnet.Contains(nextHop) {
				return nil
			}
		}
	}

	return fmt.Errorf("static route %s on %s of %s: next-hop %s is not in the subnet of the interface", prefix, name, hostname, nextHop.String())
}
