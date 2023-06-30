package repository

import (
	"fmt"

	"github.com/rs/zerolog/log"

	"github.com/criteo/data-aggregation-api/internal/ingestor/cmdb"
	"github.com/criteo/data-aggregation-api/internal/model/cmdb/bgp"
	"github.com/criteo/data-aggregation-api/internal/model/cmdb/routingpolicy"
	"github.com/criteo/data-aggregation-api/internal/model/dcim"
	"github.com/criteo/data-aggregation-api/internal/report"
)

func statsReport(message string, severity report.Severity) report.Message {
	return report.Message{
		Type:     report.IngestorMessage,
		Severity: severity,
		Text:     message,
	}
}

type AssetsPerDevice struct {
	BGPGlobal      map[string]*bgp.BGPGlobal
	BGPsessions    map[string][]*bgp.Session
	PeerGroups     map[string][]*bgp.PeerGroup
	PrefixLists    map[string][]*routingpolicy.PrefixList
	CommunityLists map[string][]*routingpolicy.CommunityList
	RoutePolicies  map[string][]*routingpolicy.RoutePolicy
}

type Assets struct {
	DeviceInventory    []*dcim.NetworkDevice
	CmdbBGPGlobal      []*bgp.BGPGlobal
	CmdbBGPSessions    []*bgp.Session
	CmdbPeerGroups     []*bgp.PeerGroup
	CmdbRoutePolicies  []*routingpolicy.RoutePolicy
	CmdbPrefixLists    []*routingpolicy.PrefixList
	CmdbCommunityLists []*routingpolicy.CommunityList
}

func (i *Assets) Precompute() *AssetsPerDevice {
	var precomputed AssetsPerDevice
	precomputed.BGPGlobal = cmdb.PrecomputeBGPGlobal(i.CmdbBGPGlobal)
	precomputed.BGPsessions = cmdb.PrecomputeBGPSessions(i.CmdbBGPSessions)
	precomputed.PeerGroups = cmdb.PrecomputePeerGroups(i.CmdbPeerGroups) //nolint:staticcheck // to ignore deprecation notice
	precomputed.PrefixLists = cmdb.PrecomputePrefixLists(i.CmdbPrefixLists)
	precomputed.CommunityLists = cmdb.PrecomputeCommunityLists(i.CmdbCommunityLists)
	precomputed.RoutePolicies = cmdb.PrecomputeRoutePolicies(i.CmdbRoutePolicies)

	return &precomputed
}

func (i *Assets) getStats() map[string]int {
	return map[string]int{
		"devices":        len(i.DeviceInventory),
		"bgpGlobal":      len(i.CmdbBGPGlobal),
		"bgpSessions":    len(i.CmdbBGPSessions),
		"peerGroups":     len(i.CmdbPeerGroups),
		"routePolicies":  len(i.CmdbRoutePolicies),
		"prefixLists":    len(i.CmdbPrefixLists),
		"communityLists": len(i.CmdbCommunityLists),
	}
}

// PrintStats prints number of asset per ingestor.
func (i *Assets) PrintStats() {
	for stat, val := range i.getStats() {
		log.Info().Str("stats", stat).Int("value", val).Msg("new assets found")
	}
}

// ReportStats sends stats to current Report.
func (i *Assets) ReportStats(messageChan chan<- report.Message) {
	for stat, val := range i.getStats() {
		messageChan <- statsReport(fmt.Sprintf("found %d %s", val, stat), report.Info)
	}
}
