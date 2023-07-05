package routingpolicy

import (
	"net"

	"github.com/criteo/data-aggregation-api/internal/model/cmdb/common"
)

type RoutingProtocols string

const (
	BGP       RoutingProtocols = "bgp"
	ISIS      RoutingProtocols = "isis"
	Static    RoutingProtocols = "static"
	Connected RoutingProtocols = "connected"
	OSPF      RoutingProtocols = "ospf"
	Unset     RoutingProtocols = ""
)

type BGPRouteType string

const (
	IBGP      BGPRouteType = "ibgp"
	EBGP      BGPRouteType = "ebgp"
	NoBGPType BGPRouteType = ""
)

type RouteProtocolOrigin string

const (
	OriginIGP        RouteProtocolOrigin = "igp"
	OriginEGP        RouteProtocolOrigin = "egp"
	OriginIncomplete RouteProtocolOrigin = "incomplete"
	OriginNone       RouteProtocolOrigin = ""
)

type RoutePolicyTerm struct { //nolint:govet // for readability
	Sequence    int      `json:"sequence"    validate:"required"`
	Decision    Decision `json:"decision"    validate:"required"`
	Description string   `json:"description" validate:"omitempty"`

	FromBGPCommunityList *struct {
		Name string `json:"name" validate:"required"`
	} `json:"from_bgp_community_list" validate:"omitempty"`
	FromPrefixList *struct {
		Name string `json:"name" validate:"required"`
	} `json:"from_prefix_list" validate:"omitempty"`
	FromSourceProtocol RoutingProtocols `json:"from_source_protocol" validate:"omitempty"`
	FromRouteType      BGPRouteType     `json:"from_route_type"      validate:"required"`
	FromLocalPref      uint32           `json:"from_local_pref"      validate:"omitempty"`

	SetOrigin              RouteProtocolOrigin `json:"set_origin"                 validate:"omitempty"`
	SetASPathPrependASN    *common.ASN         `json:"set_as_path_prepend_asn"    validate:"omitempty"`
	SetASPathPrependRepeat uint8               `json:"set_as_path_prepend_repeat" validate:"omitempty"`
	SetCommunity           string              `json:"set_community"              validate:"omitempty"`
	SetLargeCommunity      string              `json:"set_large_community"        validate:"omitempty"`
	SetNextHop             *net.IP             `json:"set_next_hop"               validate:"omitempty"`
	SetLocalPref           uint32              `json:"set_local_pref"             validate:"omitempty"`
	SetMetric              uint32              `json:"set_metric"                 validate:"omitempty"`
}

type RoutePolicy struct {
	Name   string `json:"name" validate:"required"`
	Device struct {
		Name string `json:"name" validate:"required"`
	} `json:"device"`
	Terms []*RoutePolicyTerm `json:"terms" validate:"required"`
}

type RoutePolicyLite struct {
	Name        string `json:"name"        validate:"required"`
	Description string `json:"description" validate:"required"`
}
