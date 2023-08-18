package metrics

import (
	"github.com/criteo/data-aggregation-api/internal/app"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type Registry struct {
	AppInfo            *prometheus.GaugeVec
	BuiltDevicesNumber *prometheus.GaugeVec
	lastBuildStatus    *prometheus.GaugeVec
	buildTotal         *prometheus.CounterVec
}

func NewRegistry() Registry {
	appInfo := promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "data_aggregation_api_info",
			Help: "Version of the application",
		},
		[]string{"version", "commit", "build_time", "built_by"},
	)
	appInfo.WithLabelValues(app.Info.Version, app.Info.Commit, app.Info.BuildTime, app.Info.BuildUser).Set(1)

	return Registry{
		AppInfo: appInfo,
		BuiltDevicesNumber: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "built_devices_number",
				Help: "Number of devices built during last successful build",
			},
			[]string{},
		),

		lastBuildStatus: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "build_status",
				Help: "Last completed build status, 0=Failed, 1=Success",
			},
			[]string{},
		),

		buildTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "completed_build_total",
				Help: "Total number of completed build",
			},
			[]string{"success"},
		),
	}
}

// BuildSuccessful updates all Prometheus metrics related to build.
//
// `build_status` counter is set to 1.
// `completed_build_total` increases with success label set to true.
func (r *Registry) BuildSuccessful() {
	r.lastBuildStatus.WithLabelValues().Set(1)
	r.buildTotal.WithLabelValues("true").Inc()
}

// BuildFailed updates all Prometheus metrics related to build.
//
// `build_status` counter is set to 0.
// `completed_build_total` increases with success label set to false.
func (r *Registry) BuildFailed() {
	r.lastBuildStatus.WithLabelValues().Set(0)
	r.buildTotal.WithLabelValues("false").Inc()
}

// SetBuiltDevices updates the `built_devices` gauge.
func (r *Registry) SetBuiltDevices(count uint32) {
	r.BuiltDevicesNumber.WithLabelValues().Set(float64(count))
}
