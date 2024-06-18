package device

import (
	"encoding/json"
	"errors"

	"github.com/rs/zerolog/log"
)

const emptyJSON string = "{}"

var ErrNotFound = errors.New("not found")

// IsAFKEnabledJSON checks if one device is AFK enabled.
func (s *SafeRepository) IsAFKEnabledJSON(hostname string) ([]byte, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	dev, ok := s.devices[hostname]
	if !ok {
		return nil, ErrNotFound
	}
	if dev == nil {
		// dev is nil when failed or no configuration
		return json.Marshal(afkEnabledResponse{false})
	}

	return json.Marshal(afkEnabledResponse{dev.AFKEnabled})
}

// ListAFKEnabledDevicesJSON returns all AFK enabled devices.
func (s *SafeRepository) ListAFKEnabledDevicesJSON() ([]byte, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	r := make(map[string]afkEnabledResponse, len(s.devices))
	for hostname, dev := range s.devices {
		// dev is nil when failed or no configuration
		if dev == nil {
			r[hostname] = afkEnabledResponse{false}
		} else {
			r[hostname] = afkEnabledResponse{dev.AFKEnabled}
		}
	}

	return json.Marshal(r)
}

// GetAllDevicesOpenConfigJSON copy configuration for all devices to w.
// Compact and wrap configuration in JSON dict with hostname as key and openconfig JSON as value.
// Example: '{"hostname":{"network-instances":{...}}'.
func (s *SafeRepository) GetAllDevicesOpenConfigJSON() ([]byte, error) {
	// We need to wrap an already generated JSON string
	// json.RawMessage instead of string is to avoid escaping the embedded JSON string
	s.mutex.Lock()
	defer s.mutex.Unlock()

	var allConfig = make(map[string]json.RawMessage)
	for hostname, dev := range s.devices {
		if dev == nil {
			// Device which failed to build returns an empty dict.
			allConfig[hostname] = json.RawMessage(emptyJSON)
			continue
		}

		if cfg, err := dev.GetCompactJSON(); err != nil {
			return []byte(emptyJSON), err
		} else {
			allConfig[hostname] = cfg
		}
	}

	r, err := json.Marshal(allConfig)
	if err != nil {
		return []byte(emptyJSON), err
	}
	return r, nil
}

// GetDeviceOpenConfigJSON copy configuration for all devices to w.
// Compact and wrap configuration in JSON dict with hostname as key and openconfig JSON as value
// Example: '{"hostname":{"network-instances":{...}}'.
func (s *SafeRepository) GetDeviceOpenConfigJSON(hostname string) ([]byte, error) {
	// We need to wrap an already generated JSON string
	// json.RawMessage instead of string is to avoid escaping the embedded JSON string
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if dev, ok := s.devices[hostname]; ok {
		if dev == nil {
			return []byte(emptyJSON), errors.New("build failed for this device")
		}

		var config json.RawMessage
		var err error

		if config, err = dev.GetCompactJSON(); err != nil {
			log.Error().Err(err).Send()
			return []byte(emptyJSON), err
		}

		if r, err := json.Marshal(config); err != nil {
			log.Error().Err(err).Send()
			return []byte(emptyJSON), err
		} else {
			return r, nil
		}
	}

	return []byte(emptyJSON), nil
}
