package router

import (
	"context"
	"log/slog"
	"strings"
	"time"

	"github.com/IBM/sarama"
	"github.com/gin-gonic/gin"

	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/app"
	commonconfig "gitlab.com/b2c-e-commerce-platform/platform/backend/common/config"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/health"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/kafka"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/middleware"

	"gitlab.com/b2c-e-commerce-platform/platform/backend/services/order/app/order"
)

// New constructs a gin.Engine with all order service routes registered.
func New(d Deps, version, commit string, timeoutDuration time.Duration) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())

	if commonconfig.IsLocalEnv() {
		r.Use(gin.Logger())
	}

	r.GET("/liveness", health.Liveness(version, commit))
	r.GET("/metrics", health.Metrics())
	r.GET("/readiness", health.Readiness())

	jwtParser, jwtVerifier := d.jwtMW()

	r.Use(
		middleware.SecurityHeaders(),
		middleware.AccessControl(d.cfg.AccessControl.AllowOrigin, allowedHeaders(d.cfg.Header.RefIDHeaderKey)),
		middleware.TraceContextTraceIDMiddleware(""),
		middleware.RefIDMiddleware(d.cfg.Header.RefIDHeaderKey),
		middleware.AutoLoggingMiddleware(app.CodeSuccess),
		middleware.Timeout(timeoutDuration),
		middleware.AccessLog(),
	)

	api := r.Group("/api/v1/order/order")
	{
		// Internal endpoints: Checkout → Order (shared-secret guard, no JWT).
		internalGroup := api.Group("")
		internalGroup.Use(order.InternalAuthMiddleware(d.cfg.Internal.SharedSecret))
		{
			internalGroup.POST("/create-from-checkout", d.service.HandleCreateFromCheckout)
			internalGroup.POST("/cancel-on-checkout-failure", d.service.HandleCancelOnCheckoutFailure)
		}

		// Customer endpoints (JWT required).
		customerGroup := api.Group("")
		customerGroup.Use(middleware.JWT(jwtParser, jwtVerifier))
		{
			customerGroup.POST("/detail", d.service.HandleDetail)
			customerGroup.POST("/list-mine", d.service.HandleListMine)
			customerGroup.POST("/cancel-mine", d.service.HandleCancelMine)
		}

		// Admin endpoints (JWT required; role=ADMIN checked inside handler).
		adminGroup := api.Group("")
		adminGroup.Use(middleware.JWT(jwtParser, jwtVerifier))
		{
			adminGroup.POST("/list-admin", d.service.HandleListAdmin)
			adminGroup.POST("/update-status-admin", d.service.HandleUpdateStatusAdmin)
		}
	}

	return r
}

// StartSubscriber starts the Kafka consumer group for the order service.
// Returns a done channel and a stop func for graceful shutdown coordination.
// If Kafka is not enabled (KAFKA_ENABLED=false or brokers empty), the done
// channel is closed immediately and a no-op stop func is returned.
func StartSubscriber(ctx context.Context, d Deps) (<-chan struct{}, func()) {
	done := make(chan struct{})

	if !subscriberEnabled(d) {
		slog.Info("kafka consumer disabled (KAFKA_ENABLED=false or brokers/group/topics empty)")
		close(done)
		return done, func() {}
	}

	eventHandlers := d.service.RegisterConsumerRoutes()

	consumerCtx, cancel := context.WithCancel(ctx)

	go func() {
		defer close(done)
		defer func() {
			if r := recover(); r != nil {
				slog.Error("subscriber panicked", slog.Any("panic", r))
			}
		}()

		group := newConsumerGroup(d)
		defer func() {
			if err := group.Close(); err != nil {
				slog.Error("subscriber group close error", "error", err)
			}
		}()

		topics := splitCSV(d.cfg.Consumer.Topics)
		processor := kafka.NewEventRouter(eventHandlers)
		handler := kafka.NewConsumerGroupHandler(consumerCtx, processor)

		slog.Info("kafka subscriber started", "topics", topics)
		for {
			if err := group.Consume(consumerCtx, topics, handler); err != nil {
				slog.Error("subscriber consume error", "error", err)
				return
			}
			if consumerCtx.Err() != nil {
				return
			}
		}
	}()

	return done, cancel
}

// WaitForSubscriber waits for the subscriber goroutine to exit or times out.
func WaitForSubscriber(done <-chan struct{}, timeout time.Duration) {
	select {
	case <-done:
	case <-time.After(timeout):
		slog.Warn("consumer shutdown timeout", "timeout", timeout)
	}
}

func newConsumerGroup(d Deps) sarama.ConsumerGroup {
	return kafka.MustNewConsumerGroup(kafka.ConsumerConfig{
		Brokers:                  splitCSV(d.cfg.Consumer.Brokers),
		GroupID:                  d.cfg.Consumer.GroupID,
		OffsetsInitial:           d.cfg.Consumer.OffsetsInitial,
		RebalanceGroupStrategies: d.cfg.Consumer.RebalanceStrategy,
		KafkaConf:                kafka.NewConsumerConfigAtLeastOnce(),
	})
}

func subscriberEnabled(d Deps) bool {
	return d.cfg.Consumer.Enabled &&
		strings.TrimSpace(d.cfg.Consumer.Brokers) != "" &&
		strings.TrimSpace(d.cfg.Consumer.GroupID) != "" &&
		strings.TrimSpace(d.cfg.Consumer.Topics) != ""
}

func allowedHeaders(refIDKey string) []string {
	return []string{
		"Content-Type",
		"Content-Length",
		"Accept-Encoding",
		"X-CSRF-Token",
		"Authorization",
		"accept",
		"origin",
		"Cache-Control",
		"X-Requested-With",
		order.HeaderInternalSecret,
		refIDKey,
	}
}
