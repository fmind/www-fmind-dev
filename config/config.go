// Package config parses and validates the application's environment-driven
// settings (Twelve-Factor) into a typed struct at the process boundary.
package config

import (
	"fmt"
	"io"
	"log/slog"

	"github.com/caarlos0/env/v11"
)

// Environment enumerates the supported runtime environments. Parsing external
// input into a small enum keeps the "development"/"production" strings out of
// call sites and lets the type system reject anything else at the boundary.
type Environment string

const (
	// Development is the local, human-readable, debug-level environment.
	Development Environment = "development"
	// Production is the deployed, structured-logging, info-level environment.
	Production Environment = "production"
)

// Config holds environment-driven settings. Every operational value the app
// needs is parsed and validated here once, so misconfiguration fails fast at
// startup rather than surfacing mid-request.
type Config struct {
	Environment Environment `env:"ENVIRONMENT" envDefault:"development"`
	// AnalyticsID is the Google Analytics measurement ID. Empty disables the
	// analytics snippet entirely; it is only ever injected in Production.
	AnalyticsID string `env:"GOOGLE_ANALYTICS_ID" envDefault:""`
	Port        int    `env:"PORT"                envDefault:"8080"`
}

// Load parses configuration from the environment and validates it, failing fast
// so misconfiguration surfaces at startup rather than mid-request.
func Load() (Config, error) {
	cfg, err := env.ParseAs[Config]()
	if err != nil {
		return Config{}, fmt.Errorf("parsing config: %w", err)
	}
	switch cfg.Environment {
	case Development, Production:
	default:
		return Config{}, fmt.Errorf("invalid ENVIRONMENT %q (want %q or %q)", cfg.Environment, Development, Production)
	}
	if cfg.Port < 1 || cfg.Port > 65535 {
		return Config{}, fmt.Errorf("invalid PORT %d (want 1-65535)", cfg.Port)
	}
	return cfg, nil
}

// IsProduction reports whether the app runs in the production environment.
func (c Config) IsProduction() bool { return c.Environment == Production }

// NewHandler builds the slog handler for the environment: human-readable text at
// debug level in development, structured JSON at info level in production.
func (c Config) NewHandler(w io.Writer) slog.Handler {
	if c.Environment == Production {
		return slog.NewJSONHandler(w, &slog.HandlerOptions{Level: slog.LevelInfo})
	}
	return slog.NewTextHandler(w, &slog.HandlerOptions{Level: slog.LevelDebug})
}
