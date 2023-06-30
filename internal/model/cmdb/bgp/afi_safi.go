package bgp

import (
	"github.com/criteo/data-aggregation-api/internal/model/cmdb/routingpolicy"
)

type AfiSafiChoice string

const (
	IPv4Unicast AfiSafiChoice = "ipv4-unicast"
	IPv6Unicast AfiSafiChoice = "ipv6-unicast"
	L2vpnEvpn   AfiSafiChoice = "l2vpn-evpn"
)

type AfiSafi struct {
	RoutePolicyIn  *routingpolicy.RoutePolicyLite `json:"route_policy_in"  validate:"omitempty"`
	RoutePolicyOut *routingpolicy.RoutePolicyLite `json:"route_policy_out" validate:"omitempty"`
	Name           AfiSafiChoice                  `json:"afi_safi_name"    validate:"required"`
}
