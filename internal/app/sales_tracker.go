package app

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/andreyxaxa/Sales-Tracker/config"
	"github.com/andreyxaxa/Sales-Tracker/internal/controller/restapi"
	"github.com/andreyxaxa/Sales-Tracker/internal/repo/persistent/analytics"
	"github.com/andreyxaxa/Sales-Tracker/internal/repo/persistent/category"
	"github.com/andreyxaxa/Sales-Tracker/internal/repo/persistent/product"
	"github.com/andreyxaxa/Sales-Tracker/internal/repo/persistent/sale"
	"github.com/andreyxaxa/Sales-Tracker/pkg/httpserver"
	"github.com/andreyxaxa/Sales-Tracker/pkg/logger"
	"github.com/andreyxaxa/Sales-Tracker/pkg/postgres"
)

func Run(cfg *config.Config) {
	// Logger
	l := logger.New(cfg.Log.Level)

	// Repository
	pg, err := postgres.New(cfg.PG.URL, postgres.MaxPoolSize(cfg.PG.PoolMax))
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - postgres.New: %w", err))
	}
	defer pg.Close()

	// Use-Case
	categoryUseCase := category.New(pg)
	productUseCase := product.New(pg)
	saleUseCase := sale.New(pg)
	analyticsUseCase := analytics.New(pg)

	// HTTP Server
	httpServer := httpserver.New(
		l,
		httpserver.Port(cfg.HTTP.Port),
		httpserver.Prefork(cfg.HTTP.UsePreforkMode),
	)
	restapi.NewRouter(
		httpServer.App,
		cfg,
		categoryUseCase,
		productUseCase,
		saleUseCase,
		analyticsUseCase,
		l,
	)

	// Start Server
	httpServer.Start()

	// Waiting Signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info("app - Run - signal %s", s.String())
	case err = <-httpServer.Notify():
		l.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	}

	err = httpServer.Shutdown()
	if err != nil {
		l.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}
}
