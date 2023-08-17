package cmdb

import (
	"fmt"

	"github.com/rs/zerolog/log"

	"github.com/criteo/data-aggregation-api/internal/ingestor/netbox"
	"github.com/criteo/data-aggregation-api/internal/model/cmdb/routingpolicy"
)

// GetRoutePolicies returns all route-policies defined in the CDMB.
func GetRoutePolicies() ([]*routingpolicy.RoutePolicy, error) {
	response := netbox.NetboxResponse[routingpolicy.RoutePolicy]{}
	params := deviceDatacenterFilter()

	err := netbox.Get("/api/plugins/cmdb/route-policies/", &response, params)
	if err != nil {
		return nil, fmt.Errorf("route-policies fetching failure: %w", err)
	}

	if response.Count != len(response.Results) {
		log.Warn().Msg("some route-policies have not been fetched")
	}

	return response.Results, nil
}

// PrecomputeRoutePolicies associates each found route-policies to the matching devices.
func PrecomputeRoutePolicies(routePolicies []*routingpolicy.RoutePolicy) map[string][]*routingpolicy.RoutePolicy {
	var routePoliciesPerDevice = make(map[string][]*routingpolicy.RoutePolicy)
	for _, routePolicy := range routePolicies {
		routePoliciesPerDevice[routePolicy.Device.Name] = append(routePoliciesPerDevice[routePolicy.Device.Name], routePolicy)
	}

	return routePoliciesPerDevice
}
