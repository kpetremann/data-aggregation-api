package routingpolicy

import (
	"strconv"

	"github.com/criteo/data-aggregation-api/internal/model/cmdb/routingpolicy"
	"github.com/criteo/data-aggregation-api/internal/model/openconfig"
)

var mode = map[routingpolicy.IPVersion]openconfig.E_PrefixSet_Mode{
	routingpolicy.IPv4: openconfig.PrefixSet_Mode_IPV4,
	routingpolicy.IPv6: openconfig.PrefixSet_Mode_IPV6,
}

func getMaskLengthRange(term *routingpolicy.PrefixListTerm) string {
	if term.LessOrEqual == 0 && term.GreaterOrEqual == 0 {
		return "exact"
	}

	var startRange, endRange string

	if term.GreaterOrEqual > 0 {
		startRange = strconv.Itoa(term.GreaterOrEqual)
	} else {
		startRange = strconv.Itoa(term.Prefix.Netmask)
	}

	if term.LessOrEqual > 0 {
		endRange = strconv.Itoa(term.LessOrEqual)
	} else {
		if term.Prefix.IP.To4() != nil {
			endRange = "32"
		} else {
			endRange = "128"
		}
	}

	return startRange + ".." + endRange
}

func extractPrefixListTerms(terms []*routingpolicy.PrefixListTerm) map[openconfig.RoutingPolicy_DefinedSets_PrefixSet_Prefix_Key]*openconfig.RoutingPolicy_DefinedSets_PrefixSet_Prefix {
	var prefixTerms = make(map[openconfig.RoutingPolicy_DefinedSets_PrefixSet_Prefix_Key]*openconfig.RoutingPolicy_DefinedSets_PrefixSet_Prefix)

	for _, term := range terms {
		prefix := term.Prefix.String()
		mask := getMaskLengthRange(term)

		term := openconfig.RoutingPolicy_DefinedSets_PrefixSet_Prefix{
			IpPrefix:        &prefix,
			MasklengthRange: &mask,
		}

		key := openconfig.RoutingPolicy_DefinedSets_PrefixSet_Prefix_Key{
			IpPrefix:        prefix,
			MasklengthRange: mask,
		}

		prefixTerms[key] = &term
	}

	return prefixTerms
}

// PrefixListsToOpenconfig converts precomputed prefix-lists to OpenConfig.
// OpenConfig path: /routing-policy/defined-sets/prefix-sets/.
func PrefixListsToOpenconfig(prefixLists []*routingpolicy.PrefixList) map[string]*openconfig.RoutingPolicy_DefinedSets_PrefixSet {
	prefixSets := make(map[string]*openconfig.RoutingPolicy_DefinedSets_PrefixSet)

	for _, prefixList := range prefixLists {
		prefixSet := openconfig.RoutingPolicy_DefinedSets_PrefixSet{
			Name:   &prefixList.Name,
			Mode:   mode[prefixList.IPVersion],
			Prefix: extractPrefixListTerms(prefixList.Terms),
		}

		prefixSets[prefixList.Name] = &prefixSet
	}

	return prefixSets
}
