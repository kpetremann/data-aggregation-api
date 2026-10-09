package tacacs

import (
	"fmt"
	"sort"

	"github.com/criteo/data-aggregation-api/internal/model/cmdb/tacacs"
	"github.com/criteo/data-aggregation-api/internal/model/openconfig"
	"github.com/openconfig/ygot/ygot"
)

// tacacsServerGroupName is the fixed name of the TACACS+ server group: the
// CMDB has no notion of multiple named server groups, only one set of
// servers per device.
const tacacsServerGroupName = "TACACS+"

// TacacsToOpenConfigAAA converts a CMDB TACACS configuration into the aaa
// container of openconfig-system, under a single TACACS+ server group.
//
// It returns nil when there is no server to advertise: the CMDB holds a
// single passkey for the whole device, and the model has no place for it
// outside of a server entry.
//
// OpenConfig's server list has no ordering or priority leaf of its own, so
// the CMDB priority is carried by the criteo-aaa-ext priority leaf instead of
// the position in the list. Servers are still appended by descending
// priority, ties broken on the address, to keep the emitted JSON stable.
//
// The CMDB holds one passkey per device where the model has one shared
// secret per server, so it is set on each of them.
func TacacsToOpenConfigAAA(config *tacacs.Tacacs) (*openconfig.System_Aaa, error) {
	if config == nil || len(config.ServerList) == 0 {
		return nil, nil
	}

	servers := make([]tacacs.Server, len(config.ServerList))
	copy(servers, config.ServerList)
	sort.Slice(servers, func(i, j int) bool {
		if servers[i].Priority != servers[j].Priority {
			return servers[i].Priority > servers[j].Priority
		}
		return servers[i].ServerAddress < servers[j].ServerAddress
	})

	aaa := &openconfig.System_Aaa{}
	serverGroup, err := aaa.NewServerGroup(tacacsServerGroupName)
	if err != nil {
		return nil, fmt.Errorf("failed to create the TACACS+ server group: %w", err)
	}
	serverGroup.Type = openconfig.AaaTypes_AAA_SERVER_TYPE_TACACS

	for _, server := range servers {
		entry, err := serverGroup.NewServer(server.ServerAddress)
		if err != nil {
			return nil, fmt.Errorf("failed to add the TACACS+ server %s: %w", server.ServerAddress, err)
		}

		if server.Priority != 0 {
			entry.Priority = ygot.Uint32(server.Priority)
		}

		if server.TCPPort != 0 || config.Passkey != "" {
			entry.Tacacs = &openconfig.System_Aaa_ServerGroup_Server_Tacacs{}
			if server.TCPPort != 0 {
				entry.Tacacs.Port = ygot.Uint16(server.TCPPort)
			}
			if config.Passkey != "" {
				entry.Tacacs.SecretKey = ygot.String(config.Passkey)
			}
		}
	}

	return aaa, nil
}
