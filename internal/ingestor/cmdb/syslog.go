package cmdb

import (
	"fmt"

	"github.com/rs/zerolog/log"

	"github.com/criteo/data-aggregation-api/internal/ingestor/netbox"
	"github.com/criteo/data-aggregation-api/internal/model/cmdb/syslog"
)

// GetSyslog returns all syslog configurations from the Network CMDB.
func GetSyslog() ([]*syslog.Syslog, error) {
	response := netbox.NetboxResponse[syslog.Syslog]{}
	params := deviceDatacenterFilter()

	err := netbox.Get("/api/plugins/cmdb/syslog/", &response, params)
	if err != nil {
		return nil, fmt.Errorf("syslog fetching failure: %w", err)
	}

	if len(response.Results) == 0 {
		log.Warn().Msg("no syslog configuration found")
	}

	return response.Results, nil
}

// PrecomputeSyslog associates each syslog configuration with its device.
func PrecomputeSyslog(configs []*syslog.Syslog) map[string]*syslog.Syslog {
	var syslogPerDevice = make(map[string]*syslog.Syslog)
	for _, config := range configs {
		syslogPerDevice[config.Device.Name] = config
	}

	return syslogPerDevice
}
