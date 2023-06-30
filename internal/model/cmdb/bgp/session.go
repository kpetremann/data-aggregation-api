package bgp

import (
	"github.com/criteo/data-aggregation-api/internal/model/cmdb/common"
	"github.com/criteo/data-aggregation-api/internal/model/cmdb/routingpolicy"
	"github.com/criteo/data-aggregation-api/internal/types"
)

type SessionStatus string

const (
	StatusActive      SessionStatus = "active"
	StatusPlanned     SessionStatus = "planned"
	StatusMaintenance SessionStatus = "maintenance"
	StatusOffline     SessionStatus = "offline"
)

type Address struct {
	Address types.CIDR `json:"address" validate:"required"`
	Family  int        `json:"family"  validate:"required"`
}

type DeviceSession struct {
	Device struct {
		Name string `json:"name" validate:"required"`
	} `json:"device" validate:"required"`

	LocalAsn        common.ASN                     `json:"local_asn"        validate:"required"`
	PeerGroup       *PeerGroupLite                 `json:"peer_group"       validate:"omitempty"`
	RoutePolicyIn   *routingpolicy.RoutePolicyLite `json:"route_policy_in"  validate:"omitempty"`
	RoutePolicyOut  *routingpolicy.RoutePolicyLite `json:"route_policy_out" validate:"omitempty"`
	AfiSafis        []*AfiSafi                     `json:"afi_safis"        validate:"required"`
	Description     string                         `json:"description"      validate:"omitempty"`
	LocalAddress    Address                        `json:"local_address"    validate:"required"`
	MaximumPrefixes uint32                         `json:"maximum_prefixes" validate:"omitempty"`
	EnforceFirstAs  bool                           `json:"enforce_first_as" validate:"omitempty"`
}

type Session struct {
	Status   SessionStatus `json:"status"   validate:"required"`
	Password string        `json:"password" validate:"omitempty"`
	PeerA    DeviceSession `json:"peer_a"   validate:"required"`
	PeerB    DeviceSession `json:"peer_b"   validate:"required"`
}
