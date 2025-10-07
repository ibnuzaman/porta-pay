package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/ibnuzaman/porta-pay/pkg/config"
	"github.com/ibnuzaman/porta-pay/pkg/database"
	"github.com/ibnuzaman/porta-pay/pkg/logger"
	"github.com/ibnuzaman/porta-pay/pkg/tracer"

	"github.com/ibnuzaman/porta-pay/services/booking/internal/delivery/http/handler"
	"github.com/ibnuzaman/porta-pay/services/booking/internal/delivery/http/router"
	"github.com/ibnuzaman/porta-pay/services/booking/internal/repository"
	"github.com/ibnuzaman/porta-pay/services/booking/internal/usecase"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	log := logger.New(cfg.AppName, cfg.Env)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Setup tracing
	shutdownTracer, err := tracer.SetupTracer(ctx, cfg.AppName, cfg.OTLPEndpoint)
	if err != nil {
		log.Warn().Err(err).Msg("Failed to setup tracer")
	} else {
		defer shutdownTracer(context.Background())
	}

	// Setup database (skip if POSTGRES_DSN is not set)
	var r chi.Router
	if cfg.GetDSN() != "" && cfg.GetDSN() != "postgres://::@:0/?sslmode=disable" {
		db := database.Open(cfg.GetDSN())
		defer db.Close()

		// Dependency injection - Clean Architecture wiring
		bookingRepo := repository.NewPostgresBookingRepository(db)
		bookingUsecase := usecase.NewBookingUsecase(bookingRepo)
		bookingHandler := handler.NewBookingHandler(bookingUsecase)

		// Setup router with all middleware applied
		r = router.NewBookingRouter(bookingHandler)
	} else {
		log.Warn().Msg("Database not configured, running in health-check mode only")
		r = setupHealthOnlyRouter()
	}

	server := &http.Server{
		Addr:         cfg.HTTPAddr,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("Booking service starting on %s", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("Server failed to start")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Graceful shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatal().Err(err).Msg("Server forced to shutdown")
	}

	log.Println("Server exited")
}

// setupHealthOnlyRouter creates a minimal router for health checks only
func setupHealthOnlyRouter() chi.Router {
	r := chi.NewRouter()

	// Basic middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	// Health check endpoints only
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok","service":"booking","mode":"health-only"}`))
	})

	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"pong"}`))
	})

	return r
}

func setupBasicRouter() chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/health", healthHandler)
	r.Get("/", homeHandler)

	return r
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Welcome to the Booking Service"}`))
}
