package router

import (
	"context"
	"fmt"
	"net/http"
	"net/http/pprof"
	"time"

	"github.com/a-h/rest"
	"github.com/a-h/rest/swaggerui"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"

	"github.com/criteo/data-aggregation-api/internal/api/auth"
	"github.com/criteo/data-aggregation-api/internal/app"
	"github.com/criteo/data-aggregation-api/internal/config"
	"github.com/criteo/data-aggregation-api/internal/convertor/device"
	"github.com/criteo/data-aggregation-api/internal/report"
)

const shutdownTimeout = 5 * time.Second
const httpReadHeaderTimeout = 60 * time.Second

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

	mux := http.NewServeMux()
	api := rest.NewAPI("AFK Data Aggregation API")

	// internal endpoints
	mux.HandleFunc("GET /metrics", prometheusMetrics(promhttp.Handler()))
	mux.HandleFunc("GET /api/version", getVersion)
	mux.HandleFunc("GET /api/health", healthCheck)

	api.Get("/metrics").
		HasResponseModel(http.StatusOK, rest.ModelOf[string]()).
		HasTags([]string{"internal"}).HasDescription("Prometheus metrics")
	api.Get("/api/version").
		HasResponseModel(http.StatusOK, rest.ModelOf[string]()).
		HasTags([]string{"internal"}).HasDescription("Details about the running version")
	api.Get("/api/health").
		HasResponseModel(http.StatusOK, rest.ModelOf[string]()).
		HasTags([]string{"internal"}).HasDescription("Dummy endpoint for basic healthcheck of the app")

	// devices endpoints
	mux.HandleFunc("GET /v1/devices/{hostname}/afk_enabled", withAuth.Wrap(m.getAFKEnabled))
	mux.HandleFunc("GET /v1/devices/{hostname}/openconfig", withAuth.Wrap(m.getDeviceOpenConfig))
	mux.HandleFunc("GET /v1/devices/{hostname}/ietfconfig", withAuth.Wrap(m.getDeviceIETFConfig))
	mux.HandleFunc("GET /v1/devices/{hostname}/config", withAuth.Wrap(m.getDeviceConfig))

	api.Get("/v1/devices/*/afk_enabled").
		HasResponseModel(http.StatusOK, rest.ModelOf[device.AFKEnabledResponse]()).
		HasTags([]string{"devices"}).HasDescription("Give all devices that should run AFK")
	api.Get("/v1/devices/{hostname}/afk_enabled").
		HasResponseModel(http.StatusOK, rest.ModelOf[map[string]device.AFKEnabledResponse]()).
		HasPathParameter("hostname", rest.PathParam{Description: "Device hostname", Type: rest.PrimitiveTypeString}).
		HasTags([]string{"devices"}).HasDescription("Tells if a device should run AFK")

	api.Get("/v1/devices/*/openconfig").
		HasResponseModel(http.StatusOK, rest.ModelOf[map[string]struct{}]()).
		HasTags([]string{"devices"}).HasDescription("Get OpenConfig data for all devices")
	api.Get("/v1/devices/{hostname}/openconfig").
		HasResponseModel(http.StatusOK, rest.ModelOf[struct{}]()).
		HasPathParameter("hostname", rest.PathParam{Description: "Device hostname", Type: rest.PrimitiveTypeString}).
		HasTags([]string{"devices"}).HasDescription("Get OpenConfig data for one specific device")

	api.Get("/v1/devices/*/ietf").
		HasResponseModel(http.StatusOK, rest.ModelOf[map[string]struct{}]()).
		HasTags([]string{"devices"}).HasDescription("Get IETF data for all devices")
	api.Get("/v1/devices/{hostname}/ietfconfig").
		HasResponseModel(http.StatusOK, rest.ModelOf[struct{}]()).
		HasPathParameter("hostname", rest.PathParam{Description: "Device hostname", Type: rest.PrimitiveTypeString}).
		HasTags([]string{"devices"}).HasDescription("Get IETF data for one or all devices")

	api.Get("/v1/devices/*/config").
		HasResponseModel(http.StatusOK, rest.ModelOf[map[string]struct{}]()).
		HasTags([]string{"devices"}).HasDescription("Get full config (OpenConfig + IETF) for one specific device")
	api.Get("/v1/devices/{hostname}/config").
		HasResponseModel(http.StatusOK, rest.ModelOf[struct{}]()).
		HasPathParameter("hostname", rest.PathParam{Description: "Device hostname", Type: rest.PrimitiveTypeString}).
		HasTags([]string{"devices"}).HasDescription("Get full config (OpenConfig + IETF) for one specific device")

	// report endpoints
	mux.HandleFunc("GET /v1/report/last", withAuth.Wrap(m.getLastReport))
	mux.HandleFunc("GET /v1/report/last/complete", withAuth.Wrap(m.getLastCompleteReport))
	mux.HandleFunc("GET /v1/report/last/successful", withAuth.Wrap(m.getLastSuccessfulReport))

	api.Get("/v1/report/last").
		HasResponseModel(http.StatusOK, rest.ModelOf[report.Report]()).
		HasTags([]string{"report"}).HasDescription("Last or ongoing build report")
	api.Get("/v1/report/last/complete").
		HasResponseModel(http.StatusOK, rest.ModelOf[report.Report]()).
		HasTags([]string{"report"}).HasDescription("Report of the last complete build (whether it failed or not)")
	api.Get("/v1/report/last/successful").
		HasResponseModel(http.StatusOK, rest.ModelOf[report.Report]()).
		HasTags([]string{"report"}).HasDescription("Report of the last successful build")

	// build endpoints
	mux.HandleFunc("POST /v1/build/trigger", withAuth.Wrap(m.triggerBuild))

	api.Post("/v1/build/trigger").
		HasResponseModel(http.StatusOK, rest.ModelOf[string]()).
		HasTags([]string{"build"}).HasDescription("Trigger a new build, only one at a time")

	if enablepprof {
		mux.HandleFunc("GET /debug/pprof/", pprof.Index)
		mux.HandleFunc("GET /debug/pprof/allocs", pprof.Index)
		mux.HandleFunc("GET /debug/pprof/goroutine", pprof.Index)
		mux.HandleFunc("GET /debug/pprof/heap", pprof.Index)
		mux.HandleFunc("GET /debug/pprof/profile", pprof.Profile)
		mux.HandleFunc("GET /debug/pprof/trace", pprof.Trace)
		mux.HandleFunc("GET /debug/pprof/symbol", pprof.Symbol)
	}

	if handler, err := swaggerHandler(api); err != nil {
		panic(fmt.Sprintf("failed to bootstrap Swagger UI: %s", err)) // this should never happen
	} else {
		mux.Handle("GET /swagger-ui", handler)
		mux.Handle("GET /swagger-ui/", handler)
	}

	listenSocket := fmt.Sprint(address, ":", port)
	log.Info().Msgf("start webserver - listening on %s", listenSocket)

	httpServer := http.Server{Addr: listenSocket, Handler: mux, ReadHeaderTimeout: httpReadHeaderTimeout}

	// TODO: handle http failure! with a channel
	go func() {
		if err := httpServer.ListenAndServe(); err != nil {
			log.Error().Err(err).Msg("stopped to listen and serve")
		}
	}()

	<-ctx.Done()
	ctxCancel, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := httpServer.Shutdown(ctxCancel); err != nil {
		log.Error().Err(err).Send()
	}

	return nil
}

func swaggerHandler(api *rest.API) (http.Handler, error) {
	spec, err := api.Spec()
	if err != nil {
		return nil, fmt.Errorf("failed to generate OpenAPI spec: %w", err)
	}

	spec.Info.Version = app.Info.Version
	spec.ExternalDocs = &openapi3.ExternalDocs{URL: "https://criteo.github.io/AFK/Data-Aggregation-API/configuration", Description: "Official documentation"}

	ui, err := swaggerui.New(spec)
	if err != nil {
		return nil, fmt.Errorf("failed to generate handler: %w", err)
	}

	return ui, nil
}
