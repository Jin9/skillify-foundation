package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"runtime/debug"
	"syscall"
	"time"

	commonconfig "gitlab.com/b2c-e-commerce-platform/platform/backend/common/config"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/logger"

	"gitlab.com/b2c-e-commerce-platform/platform/backend/services/order/config"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/services/order/router"

	_ "time/tzdata"
)

const (
	gracefulShutdownDuration = 10 * time.Second
	serverReadHeaderTimeout  = 5 * time.Second
	serverReadTimeout        = 5 * time.Second
	serverWriteTimeout       = 10 * time.Second
	handlerTimeout           = serverWriteTimeout - (time.Millisecond * 100)
)

var version = "dev"

func init() {
	if os.Getenv("GOMAXPROCS") != "" {
		runtime.GOMAXPROCS(0)
	} else {
		runtime.GOMAXPROCS(1)
	}
	if os.Getenv("GOMEMLIMIT") != "" {
		debug.SetMemoryLimit(-1)
	}
}

func main() {
	cfg := config.C(commonconfig.Env)
	_ = logger.New(
		logger.AWSKeyReplacer,
		logger.OpenSearchKeyReplacer,
		logger.CensorReplacer,
	)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	d, cleanupDeps := router.NewDeps(ctx, cfg)
	defer cleanupDeps()

	consumerDone, stopConsumer := router.StartSubscriber(ctx, d)

	defer func() {
		cancel()
		stopConsumer()
		router.WaitForSubscriber(consumerDone, gracefulShutdownDuration)
	}()

	// Trigger shutdown if subscriber exits unexpectedly.
	go func() {
		<-consumerDone
		cancel()
	}()

	r := router.New(d, version, "", handlerTimeout)
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
		slog.Info("shutting down")
		shutdownCtx, c := context.WithTimeout(context.Background(), gracefulShutdownDuration)
		defer c()
		if err := srv.Shutdown(shutdownCtx); err != nil {
			slog.Error("server shutdown error", "error", err)
		}
	}()

	slog.Info("order service starting", "port", cfg.Server.Port)

	defer func() {
		if r := recover(); r != nil {
			slog.Error("server panicked", slog.Any("panic", r))
		}
	}()

	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		slog.Error("ListenAndServe error", "error", err)
	}
	slog.Info("order service stopped")
}
