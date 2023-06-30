package cmdb

import (
	"fmt"

	"github.com/rs/zerolog/log"

	"github.com/criteo/data-aggregation-api/internal/ingestor/netbox"
	"github.com/criteo/data-aggregation-api/internal/model/cmdb/routingpolicy"
)

// GetPrefixLists returns all prefix-lists from the Network CMDB.
func GetPrefixLists() ([]*routingpolicy.PrefixList, error) {
	response := netbox.NetboxResponse[routingpolicy.PrefixList]{}

	err := netbox.Get("/api/plugins/cmdb/prefix-lists/", &response)
	if err != nil {
		return nil, fmt.Errorf("prefix-lists fetching failure: %w", err)
	}

	if response.Count != len(response.Results) {
		log.Warn().Msg("some prefix-lists have not been fetched")
	}

	return response.Results, nil
}

// PrecomputePrefixLists associates each found prefix-lists to the matching devices.
func PrecomputePrefixLists(prefixLists []*routingpolicy.PrefixList) map[string][]*routingpolicy.PrefixList {
	var prefixListsPerDevice = make(map[string][]*routingpolicy.PrefixList)
	for _, prefixList := range prefixLists {
		prefixListsPerDevice[prefixList.Device.Name] = append(prefixListsPerDevice[prefixList.Device.Name], prefixList)
	}

	return prefixListsPerDevice
}
