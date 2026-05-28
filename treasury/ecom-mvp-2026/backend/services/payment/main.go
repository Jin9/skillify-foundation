package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/logger"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/services/payment/config"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/services/payment/router"

	_ "time/tzdata"
)

const (
	gracefulShutdownDuration = 10 * time.Second
	serverReadHeaderTimeout  = 5 * time.Second
	serverReadTimeout        = 5 * time.Second
	serverWriteTimeout       = 10 * time.Second
	handlerTimeout           = serverWriteTimeout - 100*time.Millisecond
)

func main() {
	cfg := config.C()

	_ = logger.New(
		logger.CensorReplacer,
	)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	d, cleanupDeps := router.NewDeps(ctx, cfg)
	defer cleanupDeps()

	// Start the expiry sweeper in a background goroutine.
	// It stops cleanly when ctx is cancelled (graceful shutdown).
	go d.Sweeper().Start(ctx)

	r := router.New(d, handlerTimeout)

	srv := &http.Server{
		Addr:              ":" + cfg.Server.Port,
		Handler:           r,
		ReadHeaderTimeout: serverReadHeaderTimeout,
		ReadTimeout:       serverReadTimeout,
		WriteTimeout:      serverWriteTimeout,
		MaxHeaderBytes:    1 << 20,
	}

	go func() {
		<-ctx.Done()
		slog.Info("shutting down", "timeout", gracefulShutdownDuration)
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), gracefulShutdownDuration)
		defer shutdownCancel()
		if err := srv.Shutdown(shutdownCtx); err != nil {
			slog.Error("HTTP server shutdown error", "error", err)
		}
	}()

	slog.Info("payment service starting", "port", cfg.Server.Port)

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		slog.Error("HTTP server ListenAndServe", "error", err)
		os.Exit(1)
	}

	slog.Info("bye")
}
