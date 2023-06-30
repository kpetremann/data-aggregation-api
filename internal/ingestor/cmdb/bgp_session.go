package cmdb

import (
	"fmt"

	"github.com/rs/zerolog/log"

	"github.com/criteo/data-aggregation-api/internal/ingestor/netbox"
	"github.com/criteo/data-aggregation-api/internal/model/cmdb/bgp"
)

// GetBGPSessions returns all BGP sessions from the Network CMDB.
func GetBGPSessions() ([]*bgp.Session, error) {
	response := netbox.NetboxResponse[bgp.Session]{}

	err := netbox.Get("/api/plugins/cmdb/bgp-sessions/", &response)
	if err != nil {
		return nil, fmt.Errorf("BGP Sessions fetching failure: %w", err)
	}

	if response.Count != len(response.Results) {
		log.Warn().Msg("some BGP session have not been fetched")
	}

	return response.Results, nil
}

// PrecomputeBGPSessions links each BGP sessions to the two matching devices.
func PrecomputeBGPSessions(sessions []*bgp.Session) map[string][]*bgp.Session {
	var sessionsPerDevice = make(map[string][]*bgp.Session)
	for _, session := range sessions {
		sessionsPerDevice[session.PeerA.Device.Name] = append(sessionsPerDevice[session.PeerA.Device.Name], session)
		sessionsPerDevice[session.PeerB.Device.Name] = append(sessionsPerDevice[session.PeerB.Device.Name], session)
	}
	return sessionsPerDevice
}
