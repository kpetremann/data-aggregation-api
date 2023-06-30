package dcim

import (
	"fmt"

	"github.com/rs/zerolog/log"

	"github.com/criteo/data-aggregation-api/internal/ingestor/netbox"
	"github.com/criteo/data-aggregation-api/internal/model/dcim"
)

// GetNetworkInventory returns network device inventory from NetBox DCIM.
// TODO: at the moment, to list the device we use the filter "role__n=server".
// > This filter should be configurable by the user.
func GetNetworkInventory() ([]*dcim.NetworkDevice, error) {
	response := netbox.NetboxResponse[dcim.NetworkDevice]{}

	if err := netbox.Get("/api/dcim/devices/?role__n=server", &response); err != nil {
		return nil, fmt.Errorf("network inventory fetching failure: %s", err)
	}

	if response.Count != len(response.Results) {
		log.Warn().Msg("some devices have not been fetched")
	}

	return response.Results, nil
}
