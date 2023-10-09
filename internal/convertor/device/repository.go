package device

import (
	"sync"
)

type SafeRepository struct {
	devices map[string]*Device
	mutex   *sync.Mutex
}

type afkEnabledResponse struct {
	AFKEnabled bool `json:"afk_enabled"`
}

func NewSafeRepository() SafeRepository {
	return SafeRepository{
		mutex:   &sync.Mutex{},
		devices: map[string]*Device{},
	}
}

// Set new device configuration in the repository.
// This method is concurrent-safe.
func (s *SafeRepository) Set(devices map[string]*Device) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.devices = devices
}
