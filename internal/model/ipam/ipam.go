package ipam

import (
	"github.com/criteo/data-aggregation-api/internal/types"
)

type IPAddress struct {
	ID          string     `json:"id"`
	Description string     `json:"description"`
	DNSName     string     `json:"dns_name"`
	Address     types.CIDR `json:"address"`
}
