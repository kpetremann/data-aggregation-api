package syslog

import (
	"github.com/criteo/data-aggregation-api/internal/model/cmdb/syslog"
	"github.com/criteo/data-aggregation-api/internal/model/openconfig"
)

// SyslogToOpenconfig converts a syslog configuration to OpenConfig remote log servers.
func SyslogToOpenconfig(config *syslog.Syslog) *openconfig.System_Logging {
	if config == nil || len(config.ServerList) == 0 {
		return nil
	}

	remoteServers := make(map[string]*openconfig.System_Logging_RemoteServer)
	for _, server := range config.ServerList {
		remoteServers[server.ServerAddress] = &openconfig.System_Logging_RemoteServer{Host: &server.ServerAddress}
	}

	return &openconfig.System_Logging{RemoteServer: remoteServers}
}
