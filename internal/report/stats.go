package report

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
)

// PerformanceStats contains durations of each step of the build pipeline.
type PerformanceStats struct {
	DataFetchingDuration time.Duration `json:"data_fetching_duration"`
	PrecomputeDuration   time.Duration `json:"precomputeduration"`
	ComputeDuration      time.Duration `json:"compute_duration"`
	BuildDuration        time.Duration `json:"build_duration"`
}

// MarshalJSON overrides PerformanceStats JSON to pretty print the duration.
func (p *PerformanceStats) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		DataFetchingDuration string `json:"data_fetching_duration"`
		PrecomputeDuration   string `json:"precomputeduration"`
		ComputeDuration      string `json:"compute_duration"`
		BuildDuration        string `json:"build_duration"`
	}{
		DataFetchingDuration: fmt.Sprint(p.DataFetchingDuration),
		PrecomputeDuration:   fmt.Sprint(p.PrecomputeDuration),
		ComputeDuration:      fmt.Sprint(p.ComputeDuration),
		BuildDuration:        fmt.Sprint(p.BuildDuration),
	})
}

// Log print stats to terminal.
func (p PerformanceStats) Log() {
	const valueKey = "value"
	const statsKey = "stats"
	log.Info().Str(statsKey, "data fetching duration").Str(valueKey, p.DataFetchingDuration.String()).Send()
	log.Info().Str(statsKey, "precompute duration").Str(valueKey, p.PrecomputeDuration.String()).Send()
	log.Info().Str(statsKey, "compute duration").Str(valueKey, p.ComputeDuration.String()).Send()
	log.Info().Str(statsKey, "total build duration").Str(valueKey, p.BuildDuration.String()).Send()
}
