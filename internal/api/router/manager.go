package router

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"

	"github.com/rs/zerolog/log"

	"github.com/criteo/data-aggregation-api/internal/api/auth"
	"github.com/criteo/data-aggregation-api/internal/config"
	"github.com/criteo/data-aggregation-api/internal/convertor/device"
	"github.com/criteo/data-aggregation-api/internal/report"
	"github.com/julienschmidt/httprouter"
)

type DevicesRepository interface {
	Set(devices map[string]*device.Device)
	ListAFKEnabledDevicesJSON() ([]byte, error)
	IsAFKEnabledJSON(hostname string) ([]byte, error)
	GetAllDevicesOpenConfigJSON() ([]byte, error)
	GetDeviceOpenConfigJSON(hostname string) ([]byte, error)
}

type Manager struct {
	devices DevicesRepository
	reports *report.Repository
}

// NewManager creates and initializes a new API manager.
func NewManager(deviceRepo DevicesRepository, reports *report.Repository) *Manager {
	return &Manager{devices: deviceRepo, reports: reports}
}

// ListenAndServe starts to serve Web API requests.
func (m *Manager) ListenAndServe(ctx context.Context, address string, port int) error {
	defer func() {
		log.Warn().Msg("Shutdown.")
	}()

	tlsConfig := &tls.Config{
		InsecureSkipVerify: config.Cfg.LDAP.InsecureSkipVerify, //nolint:gosec // configurable on purpose
	}

	withAuth := auth.NewBasicAuth()
	if config.Cfg.LDAP.Enabled {
		ldap := auth.NewLDAPAuth(config.Cfg.LDAP.URL, config.Cfg.LDAP.BindDN, config.Cfg.LDAP.Password, config.Cfg.LDAP.BaseDN, tlsConfig)
		if err := withAuth.ConfigureLdap(ldap); err != nil {
			return fmt.Errorf("failed to configure the request authenticator: %w", err)
		}
		withAuth.SetMode(auth.LDAPMode)
	}

	router := httprouter.New()

	router.GET("/api/health", healthCheck)
	router.GET("/v1/devices/:hostname/afk_enabled", withAuth.Wrap(m.getAFKEnabled))
	router.GET("/v1/devices/:hostname/openconfig", withAuth.Wrap(m.getDeviceOpenConfig))
	router.GET("/v1/report/last", withAuth.Wrap(m.getLastReport))
	router.GET("/v1/report/last/complete", withAuth.Wrap(m.getLastCompleteReport))
	router.GET("/v1/report/last/successful", withAuth.Wrap(m.getLastSuccessfulReport))

	listenSocket := fmt.Sprint(address, ":", port)
	log.Info().Msgf("Start webserver - listening on %s", listenSocket)

	httpServer := http.Server{Addr: listenSocket, Handler: router}

	// TODO: handle http failure! with a channel
	go func() {
		if err := httpServer.ListenAndServe(); err != nil {
			log.Error().Err(err).Send()
		}
	}()

	<-ctx.Done()
	if err := httpServer.Shutdown(context.Background()); err != nil {
		log.Error().Err(err).Send()
	}

	return nil
}
