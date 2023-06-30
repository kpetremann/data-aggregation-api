package bgp

import (
	"github.com/criteo/data-aggregation-api/internal/model/cmdb/common"
	"github.com/criteo/data-aggregation-api/internal/model/cmdb/routingpolicy"
)

type PeerGroupLite struct {
	Name string `json:"name" validate:"required"`
}

type PeerGroup struct {
	Device struct {
		Name string `json:"name" validate:"required"`
	} `json:"device" validate:"required"`
	Name           string                         `json:"name"             validate:"required"`
	RemoteAsn      *common.ASN                    `json:"remote_asn"       validate:"omitempty"`
	LocalAsn       *common.ASN                    `json:"local_asn"        validate:"omitempty"`
	RoutePolicyIn  *routingpolicy.RoutePolicyLite `json:"route_policy_in"  validate:"omitempty"`
	RoutePolicyOut *routingpolicy.RoutePolicyLite `json:"route_policy_out" validate:"omitempty"`
	Description    string                         `json:"description"      validate:"omitempty"`
	EnforceFirstAs bool                           `json:"enforce_first_as" validate:"omitempty"`
}
