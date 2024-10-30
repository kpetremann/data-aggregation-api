package router

import (
	"context"
	"fmt"
	"net/http/pprof"
	"time"

	"github.com/go-fuego/fuego"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"

	"github.com/criteo/data-aggregation-api/internal/api/auth"
	"github.com/criteo/data-aggregation-api/internal/config"
	"github.com/criteo/data-aggregation-api/internal/convertor/device"
	"github.com/criteo/data-aggregation-api/internal/report"
)

const shutdownTimeout = 5 * time.Second

type DevicesRepository interface {
	Set(devices map[string]*device.Device)
	ListAFKEnabledDevicesJSON() ([]byte, error)
	IsAFKEnabledJSON(hostname string) ([]byte, error)
	GetAllDevicesOpenConfigJSON() ([]byte, error)
	GetDeviceOpenConfigJSON(hostname string) ([]byte, error)
	GetAllDevicesIETFConfigJSON() ([]byte, error)
	GetDeviceIETFConfigJSON(hostname string) ([]byte, error)
	GetAllDevicesConfigJSON() ([]byte, error)
	GetDeviceConfigJSON(hostname string) ([]byte, error)
}

type Manager struct {
	devices         DevicesRepository
	reports         *report.Repository
	newBuildRequest chan<- struct{}
}

// NewManager creates and initializes a new API manager.
func NewManager(deviceRepo DevicesRepository, reports *report.Repository, restartRequest chan<- struct{}) *Manager {
	return &Manager{devices: deviceRepo, reports: reports, newBuildRequest: restartRequest}
}

// ListenAndServe starts to serve Web API requests.
func (m *Manager) ListenAndServe(ctx context.Context, address string, port int, enablepprof bool) error {
	defer func() {
		close(m.newBuildRequest)
		log.Warn().Msg("Shutdown.")
	}()

	withAuth, err := auth.NewBasicAuth(ctx, config.Cfg.Authentication)
	if err != nil {
		return err
	}

	listenSocket := fmt.Sprint(address, ":", port)
	srv := fuego.NewServer(
		fuego.WithAddr(listenSocket),
		fuego.WithOpenAPIConfig(fuego.OpenAPIConfig{SwaggerUrl: "/docs"}),
		fuego.WithoutStartupMessages(),
	)

	// / endpoints
	fuego.GetStd(srv, "/metrics", prometheusMetrics(promhttp.Handler())).Summary("Prometheus metrics")

	// /api/ endpoints
	apiGroup := fuego.Group(srv, "/api")
	fuego.GetStd(apiGroup, "/version", getVersion).Summary("version").Description("Details about the running version")
	fuego.GetStd(apiGroup, "/health", healthCheck).Summary("healthcheck").Description("Dummy endpoint for basic healthcheck of the app")

	// /v1/devices endpoints
	devicesGroup := fuego.Group(srv, "/v1/devices")
	fuego.GetStd(devicesGroup, "/{hostname}/afk_enabled", withAuth.Wrap(m.getAFKEnabled)).Summary("afk enabled").Description("Tells if a device should run AFK")
	fuego.GetStd(devicesGroup, "/{hostname}/openconfig", withAuth.Wrap(m.getDeviceOpenConfig)).Summary("device openconfig").Description("Get OpenConfig data for one or all devices ('*' for all devices)")
	fuego.GetStd(devicesGroup, "/{hostname}/ietfconfig", withAuth.Wrap(m.getDeviceIETFConfig)).Summary("device ietf config").Description("Get IETF data for one or all devices ('*' for all devices)")
	fuego.GetStd(devicesGroup, "/{hostname}/config", withAuth.Wrap(m.getDeviceConfig)).Summary("device full config").Description("Returns full config (OpenConfig + IETF) for one or all devices ('*' for all devices")

	// /v1/report endpoints
	reportGroup := fuego.Group(srv, "/v1/report")
	fuego.GetStd(reportGroup, "/last", withAuth.Wrap(m.getLastReport)).Summary("last").Description("Last or ongoing build report")
	fuego.GetStd(reportGroup, "/last/complete", withAuth.Wrap(m.getLastCompleteReport)).Summary("Report of the last complete build (whether it failed or not)")
	fuego.GetStd(reportGroup, "/last/successful", withAuth.Wrap(m.getLastSuccessfulReport)).Summary("Report of the last successful build")

	// /v1/build endpoints
	buildGroup := fuego.Group(srv, "/v1/build")
	fuego.PostStd(buildGroup, "/trigger", withAuth.Wrap(m.triggerBuild)).Summary("trigger").Description("Trigger a new build, only one at a time")

	// /v1/debug endpoints
	if enablepprof {
		debugGroup := fuego.Group(srv, "/debug").Hide()
		fuego.GetStd(debugGroup, "/pprof/", pprof.Index)
		fuego.GetStd(debugGroup, "/pprof/allocs", pprof.Index)
		fuego.GetStd(debugGroup, "/pprof/goroutine", pprof.Index)
		fuego.GetStd(debugGroup, "/pprof/heap", pprof.Index)
		fuego.GetStd(debugGroup, "/pprof/profile", pprof.Profile)
		fuego.GetStd(debugGroup, "/pprof/trace", pprof.Trace)
		fuego.GetStd(debugGroup, "/pprof/symbol", pprof.Symbol)
	}

	// TODO: handle http failure! with a channel
	go func() {
		log.Info().Msgf("Start webserver - listening on %s", listenSocket)
		if err := srv.Run(); err != nil {
			log.Error().Err(err).Msg("stopped to listen and serve")
		}
	}()

	<-ctx.Done()
	ctxCancel, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(ctxCancel); err != nil {
		log.Error().Err(err).Send()
	}

	return nil
}
