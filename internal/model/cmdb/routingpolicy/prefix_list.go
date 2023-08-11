package routingpolicy

import (
	"github.com/criteo/data-aggregation-api/internal/types"
)

type IPVersion string

const (
	IPv4 IPVersion = "ipv4"
	IPv6 IPVersion = "ipv6"
)

type PrefixListTerm struct {
	Prefix         types.CIDR `json:"prefix"   validate:"required"`
	LessOrEqual    int        `json:"le"       validate:"omitempty"`
	GreaterOrEqual int        `json:"ge"       validate:"omitempty"`
}

type PrefixList struct {
	Name   string `json:"name" validate:"required"`
	Device struct {
		Name string `json:"name" validate:"required"`
	} `json:"device" validate:"required"`
	IPVersion IPVersion         `json:"ip_version" validate:"required"`
	Terms     []*PrefixListTerm `json:"terms"      validate:"required,dive"`
}
