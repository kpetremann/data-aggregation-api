package bgp

import "github.com/criteo/data-aggregation-api/internal/model/cmdb/common"

type BGPGlobal struct {
	Device struct {
		Name string `json:"name" validate:"required"`
	} `json:"device" validate:"required"`
	LocalAsn common.ASN `json:"local_asn"                    validate:"required"`

	EBGPAdministrativeDistance *uint8  `json:"ebgp_administrative_distance" validate:"omitempty"`
	IBGPAdministrativeDistance *uint8  `json:"ibgp_administrative_distance" validate:"omitempty"`
	GracefulRestartEnabled     *bool   `json:"graceful_restart"             validate:"required"`
	GracefulRestartTime        *uint16 `json:"graceful_restart_time"        validate:"omitempty"`
	EcmpEnabled                *bool   `json:"ecmp"                         validate:"required"`
	EcmpMaximumPaths           *uint32 `json:"ecmp_maximum_paths"           validate:"omitempty"`
	RouterID                   string  `json:"router_id"                    validate:"omitempty"`
}
