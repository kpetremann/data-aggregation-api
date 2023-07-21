package routingpolicy

import (
	"github.com/criteo/data-aggregation-api/internal/model/cmdb/routingpolicy"
	"github.com/criteo/data-aggregation-api/internal/model/openconfig"
)

// RoutingPolicyToOpenconfig converts all precomputed assets in routing-policy/* to OpenConfig.
// OpenConfig path: /routing-policy/.
func RoutingPolicyToOpenconfig(prefixLists []*routingpolicy.PrefixList, communityLists []*routingpolicy.CommunityList, routePolicies []*routingpolicy.RoutePolicy) (*openconfig.RoutingPolicy, error) {
	policies, err := RoutePoliciesToOpenconfig(routePolicies)
	if err != nil {
		return nil, err
	}

	routingPolicy := openconfig.RoutingPolicy{
		PolicyDefinition: policies,
		DefinedSets: &openconfig.RoutingPolicy_DefinedSets{
			PrefixSet: PrefixListsToOpenconfig(prefixLists),
			BgpDefinedSets: &openconfig.RoutingPolicy_DefinedSets_BgpDefinedSets{
				CommunitySet: CommunityListToOpenconfig(communityLists),
			},
		},
	}

	return &routingPolicy, nil
}
