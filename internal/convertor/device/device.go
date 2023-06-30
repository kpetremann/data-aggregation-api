package device

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sync"

	bgpconvertors "github.com/criteo/data-aggregation-api/internal/convertor/bgp"
	rpconvertors "github.com/criteo/data-aggregation-api/internal/convertor/routingpolicy"
	"github.com/criteo/data-aggregation-api/internal/ingestor/repository"
	"github.com/criteo/data-aggregation-api/internal/model/cmdb/bgp"
	"github.com/criteo/data-aggregation-api/internal/model/cmdb/routingpolicy"
	"github.com/criteo/data-aggregation-api/internal/model/dcim"
	"github.com/criteo/data-aggregation-api/internal/model/openconfig"
	"github.com/openconfig/ygot/ygot"
)

const AFKEnabledTag = "afk-enabled"

var defaultInstance = "default"

type GeneratedConfig struct {
	Openconfig *openconfig.Device
	JSON       string
}

type Device struct {
	mutex           *sync.Mutex
	Dcim            *dcim.NetworkDevice
	Config          *GeneratedConfig
	BGPGlobalConfig *bgp.BGPGlobal
	Sessions        []*bgp.Session
	PeerGroups      []*bgp.PeerGroup
	PrefixLists     []*routingpolicy.PrefixList
	CommunityLists  []*routingpolicy.CommunityList
	RoutePolicies   []*routingpolicy.RoutePolicy
	AFKEnabled      bool
}

// isAFKenabled checks if the device contains the AFKEnabledTag.
func isAFKenabled(dcimInfo *dcim.NetworkDevice) bool {
	for _, tag := range dcimInfo.Tags {
		if tag.Name == AFKEnabledTag {
			return true
		}
	}
	return false
}

// NewDevice creates and populates a device with precomputed Ingestor's data.
func NewDevice(dcimInfo *dcim.NetworkDevice, devicesData *repository.AssetsPerDevice) (*Device, error) {
	device := &Device{
		mutex:      &sync.Mutex{},
		Dcim:       dcimInfo,
		AFKEnabled: isAFKenabled(dcimInfo),
	}

	// TODO: be able to set which ingestors are mandatory from the settings.yaml
	var ok bool

	// Check if there is CMDB data for the device
	device.Sessions, ok = devicesData.BGPsessions[dcimInfo.Hostname]
	if !ok {
		return nil, fmt.Errorf("no BGP session found for %s", dcimInfo.Hostname)
	}

	device.BGPGlobalConfig = devicesData.BGPGlobal[dcimInfo.Hostname]
	// TODO: uncomment once mandatory
	// if !ok {
	// 	return nil, fmt.Errorf("no BGP global configuration found for %s", dcimInfo.Hostname)
	// }

	device.PeerGroups, ok = devicesData.PeerGroups[dcimInfo.Hostname]
	if !ok {
		return nil, fmt.Errorf("no peer-groups found for %s", dcimInfo.Hostname)
	}

	device.PrefixLists, ok = devicesData.PrefixLists[dcimInfo.Hostname]
	if !ok {
		return nil, fmt.Errorf("no prefix-lists found for %s", dcimInfo.Hostname)
	}

	device.CommunityLists, ok = devicesData.CommunityLists[dcimInfo.Hostname]
	if !ok {
		return nil, fmt.Errorf("no community-lists found for %s", dcimInfo.Hostname)
	}

	device.RoutePolicies, ok = devicesData.RoutePolicies[dcimInfo.Hostname]
	if !ok {
		return nil, fmt.Errorf("no route-policies found for %s", dcimInfo.Hostname)
	}

	return device, nil
}

// GenerateOpenconfig generate the OpenConfig data for the current device.
// The CMDB data must have been precomputed before running this method.
func (d *Device) GenerateOpenconfig() error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	// Generate sub-configs
	bgpConfig, err := bgpconvertors.BGPToOpenconfig(d.Dcim.Hostname, d.BGPGlobalConfig, d.Sessions, d.PeerGroups)
	if err != nil {
		return fmt.Errorf("convert from BGP to OpenConfig failed: %w", err)
	}

	routingPolicyConfig, err := rpconvertors.RoutingPolicyToOpenconfig(d.Dcim.Hostname, d.PrefixLists, d.CommunityLists, d.RoutePolicies)
	if err != nil {
		return fmt.Errorf("convert from Routing Policy to OpenConfig failed: %w", err)
	}

	// Assemble global configuration
	bgpKey := openconfig.NetworkInstance_Protocol_Key{Identifier: openconfig.PolicyTypes_INSTALL_PROTOCOL_TYPE_BGP, Name: "bgp"}

	config := openconfig.Device{
		RoutingPolicy: routingPolicyConfig,
		NetworkInstance: map[string]*openconfig.NetworkInstance{
			"default": {
				Name: &defaultInstance,
				Protocol: map[openconfig.NetworkInstance_Protocol_Key]*openconfig.NetworkInstance_Protocol{
					bgpKey: {
						Bgp:        bgpConfig,
						Name:       &bgpKey.Name,
						Identifier: openconfig.PolicyTypes_INSTALL_PROTOCOL_TYPE_BGP,
					},
				},
			},
		},
	}

	devJSON, err := ygot.EmitJSON(
		&config,
		&ygot.EmitJSONConfig{
			Format:         ygot.RFC7951,
			SkipValidation: false,
			Indent:         "  ",
		},
	)
	if err != nil {
		return fmt.Errorf("failed to transform an openconfig device specification (%s) into JSON using ygot: %w", d.Dcim.Hostname, err)
	}

	d.Config = &GeneratedConfig{
		Openconfig: &config,
		JSON:       devJSON,
	}

	return nil
}

// GetCompactJSON returns OpenConfig result in not indented JSON format.
// Generated JSON is already indented by Ygot - currently there is no option to not indent the JSON.
func (d *Device) GetCompactJSON() ([]byte, error) {
	out := bytes.NewBuffer(nil)
	err := json.Compact(out, []byte(d.Config.JSON))
	if err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}
