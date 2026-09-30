package ntp

import (
	"github.com/criteo/data-aggregation-api/internal/model/cmdb/ntp"
	"github.com/criteo/data-aggregation-api/internal/model/openconfig"
)

func NTPToOpenconfig(config *ntp.NTP) *openconfig.System_Ntp {
	if config == nil {
		return nil
	}
	servers := make(map[string]*openconfig.System_Ntp_Server)

	for _, s := range config.ServerList {
		servers[s.ServerAddress] = &openconfig.System_Ntp_Server{Address: &s.ServerAddress}
	}

	return &openconfig.System_Ntp{
		Server: servers,
	}
}
