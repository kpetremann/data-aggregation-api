package routingpolicy

import (
	"github.com/criteo/data-aggregation-api/internal/model/cmdb/routingpolicy"
	"github.com/criteo/data-aggregation-api/internal/model/openconfig"
)

func extractCommunityListTerms(communityList *routingpolicy.CommunityList) []openconfig.RoutingPolicy_DefinedSets_BgpDefinedSets_CommunitySet_CommunityMember_Union {
	terms := make([]openconfig.RoutingPolicy_DefinedSets_BgpDefinedSets_CommunitySet_CommunityMember_Union, 0, len(communityList.Terms))
	for _, term := range communityList.Terms {
		newTerm := openconfig.UnionString(term.Community)
		terms = append(terms, newTerm)
	}
	return terms
}

// CommunityListToOpenconfig converts precomputed prefix-lists to OpenConfig.
// OpenConfig path: /routing-policy/defined-sets/bgp-defined-sets/community-sets/.
func CommunityListToOpenconfig(hostname string, communityLists []*routingpolicy.CommunityList) map[string]*openconfig.RoutingPolicy_DefinedSets_BgpDefinedSets_CommunitySet {
	var communitySets = make(map[string]*openconfig.RoutingPolicy_DefinedSets_BgpDefinedSets_CommunitySet)

	for _, communityList := range communityLists {
		communitySet := openconfig.RoutingPolicy_DefinedSets_BgpDefinedSets_CommunitySet{
			CommunitySetName: &communityList.Name,
			CommunityMember:  extractCommunityListTerms(communityList),
		}

		communitySets[communityList.Name] = &communitySet
	}

	return communitySets
}
