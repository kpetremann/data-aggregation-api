package device

import (
	"encoding/json"
	"errors"

	"github.com/rs/zerolog/log"
)

const emptyJSON string = "{}"

var ErrBuidFailed = errors.New("build failed for this device")
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

		if cfg, err := dev.GetCompactOpenconfigJSON(); err != nil {
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

// GetAllDevicesIETFConfigJSON copy configuration for all devices to w.
// Compact and wrap configuration in JSON dict with hostname as key and openconfig JSON as value.
// Example: '{"hostname":{"network-instances":{...}}'.
func (s *SafeRepository) GetAllDevicesIETFConfigJSON() ([]byte, error) {
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

		if cfg, err := dev.GetCompactIETFJSON(); err != nil {
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

// GetAllDevicesConfigJSON copy configuration for all devices to w.
// Compact and wrap configuration in JSON dict with hostname as key and openconfig JSON as value.
// Example: '{"hostname":{"network-instances":{...}}'.
func (s *SafeRepository) GetAllDevicesConfigJSON() ([]byte, error) {
	// We need to wrap an already generated JSON string
	// json.RawMessage instead of string is to avoid escaping the embedded JSON string
	s.mutex.Lock()
	defer s.mutex.Unlock()

	var allConfig = make(map[string]map[string]json.RawMessage)
	for hostname, dev := range s.devices {
		if dev == nil {
			// Device which failed to build returns an empty dict.
			allConfig[hostname] = map[string]json.RawMessage{}
			continue
		}
		allConfig[hostname] = make(map[string]json.RawMessage)
		if cfg, err := dev.GetCompactIETFJSON(); err != nil {
			return []byte(emptyJSON), err
		} else {
			allConfig[hostname]["ietfconfig"] = cfg
		}
		if cfg, err := dev.GetCompactOpenconfigJSON(); err != nil {
			return []byte(emptyJSON), err
		} else {
			allConfig[hostname]["openconfig"] = cfg
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
			return []byte(emptyJSON), ErrBuidFailed
		}

		var config json.RawMessage
		var err error

		if config, err = dev.GetCompactOpenconfigJSON(); err != nil {
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

// GetDeviceIETFConfigJSON copy configuration for all devices to w.
// Compact and wrap configuration in JSON dict with hostname as key and ietf JSON as value.
func (s *SafeRepository) GetDeviceIETFConfigJSON(hostname string) ([]byte, error) {
	// We need to wrap an already generated JSON string
	// json.RawMessage instead of string is to avoid escaping the embedded JSON string
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if dev, ok := s.devices[hostname]; ok {
		if dev == nil {
			return []byte(emptyJSON), ErrBuidFailed
		}

		var config json.RawMessage
		var err error

		if config, err = dev.GetCompactIETFJSON(); err != nil {
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

// GetDeviceConfigJSON copy configuration for all devices to w.
// Compact and wrap configuration in JSON dict with hostname as key and ietf JSON as value.
func (s *SafeRepository) GetDeviceConfigJSON(hostname string) ([]byte, error) {
	// We need to wrap an already generated JSON string
	// json.RawMessage instead of string is to avoid escaping the embedded JSON string
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if dev, ok := s.devices[hostname]; ok {
		if dev == nil {
			return []byte(emptyJSON), ErrBuidFailed
		}

		var configIETF json.RawMessage
		var configOpenconfig json.RawMessage
		var err error

		if configIETF, err = dev.GetCompactIETFJSON(); err != nil {
			log.Error().Err(err).Send()
			return []byte(emptyJSON), err
		}

		if configOpenconfig, err = dev.GetCompactOpenconfigJSON(); err != nil {
			log.Error().Err(err).Send()
			return []byte(emptyJSON), err
		}

		config := map[string]json.RawMessage{
			"openconfig": configOpenconfig,
			"ietf":       configIETF,
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
