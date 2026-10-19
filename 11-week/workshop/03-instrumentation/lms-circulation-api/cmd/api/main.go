package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/MENESESLUIS-bot/w11/lms-circulation-api/internal/api"
	"github.com/MENESESLUIS-bot/w11/lms-circulation-api/internal/catalogclient"
	"github.com/MENESESLUIS-bot/w11/lms-circulation-api/internal/logging"
	"github.com/MENESESLUIS-bot/w11/lms-circulation-api/internal/metrics"
)

func main() {
	log := logging.New(os.Stdout, "lms-circulation-api", env("APP_VERSION", "2.0.1"), env("APP_ENV", "dev"), slog.LevelDebug)
	m := metrics.New()
	srv := &api.Server{
		Catalog: catalogclient.New(env("CATALOG_SERVICE_URL", "http://localhost:8082"), 5*time.Second, m.CatalogObserver),
		Log:     log,
		Metrics: m,
	}

	ctx := context.Background()
	addr := ":" + env("PORT", "8083")
	logging.Event(ctx, log, slog.LevelInfo, "server.started", "listening on "+addr, 0)
	if err := http.ListenAndServe(addr, srv.Routes()); err != nil {
		logging.Event(ctx, log, slog.LevelError, "server.stopped", err.Error(), 0)
		os.Exit(1)
	}
}

func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
