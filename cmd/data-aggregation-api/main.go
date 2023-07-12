package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/criteo/data-aggregation-api/internal/api/router"
	"github.com/criteo/data-aggregation-api/internal/config"
	"github.com/criteo/data-aggregation-api/internal/convertor/device"
	"github.com/criteo/data-aggregation-api/internal/job"
	"github.com/criteo/data-aggregation-api/internal/report"
)

func configureLogging(logLevel string) error {
	level, err := zerolog.ParseLevel(logLevel)
	if err != nil {
		return fmt.Errorf("failed to parse log level '%s': %w", logLevel, err)
	}
	zerolog.SetGlobalLevel(level)
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix               //nolint:reassign // it is the way
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr}) //nolint:reassign // it is the way

	return nil
}

var (
	version   = ""
	buildTime = "unknown"
	buildUser = "unknown"
)

func run() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := config.LoadConfig(); err != nil {
		return err
	}
	if err := configureLogging(config.Cfg.LogLevel); err != nil {
		return err
	}

	log.Info().Str("version", version).Send()
	log.Info().Str("build-time", buildTime).Send()
	log.Info().Str("build-user", buildUser).Send()

	deviceRepo := device.NewSafeRepository()
	reports := report.NewRepository()

	// TODO: be able to close goroutine when the context is closed (graceful shutdown)
	go job.StartBuildLoop(&deviceRepo, &reports)

	if err := router.NewManager(&deviceRepo, &reports).ListenAndServe(ctx, config.Cfg.API.ListenAddress, config.Cfg.API.ListenPort); err != nil {
		return fmt.Errorf("webserver error: %w", err)
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatal().Err(err).Send()
	}
}
