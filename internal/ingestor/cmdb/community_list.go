package cmdb

import (
	"fmt"

	"github.com/rs/zerolog/log"

	"github.com/criteo/data-aggregation-api/internal/ingestor/netbox"
	"github.com/criteo/data-aggregation-api/internal/model/cmdb/routingpolicy"
)

// GetCommunityLists returns all community-lists from the Network CMDB.
func GetCommunityLists() ([]*routingpolicy.CommunityList, error) {
	response := netbox.NetboxResponse[routingpolicy.CommunityList]{}

	err := netbox.Get("/api/plugins/cmdb/bgp-community-lists/", &response)
	if err != nil {
		return nil, fmt.Errorf("BGP Community Lists fetching failure: %w", err)
	}

	if response.Count != len(response.Results) {
		log.Warn().Msg("some community-lists have not been fetched")
	}

	return response.Results, nil
}

// PrecomputeCommunityLists associates each found community-lists to the matching devices.
func PrecomputeCommunityLists(communityLists []*routingpolicy.CommunityList) map[string][]*routingpolicy.CommunityList {
	var communityListsPerDevice = make(map[string][]*routingpolicy.CommunityList)
	for _, communityList := range communityLists {
		communityListsPerDevice[communityList.Device.Name] = append(communityListsPerDevice[communityList.Device.Name], communityList)
	}

	return communityListsPerDevice
}
