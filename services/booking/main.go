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
	"github.com/ibnuzaman/porta-pay/pkg/dbx"
	"github.com/ibnuzaman/porta-pay/pkg/logx"
	"github.com/ibnuzaman/porta-pay/pkg/otelx"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	log := logx.New(cfg.AppName, cfg.Env)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	shutdownTracer, err := otelx.SetupTracer(ctx, cfg.AppName, cfg.OTLPEndpoint)
	defer shutdownTracer(context.Background())

	db := dbx.Open(cfg.PostgresDSN)
	defer db.Close()

	//Wiring
	// repo := NewBookingRepository(db)
	// svc := NewBookingService(repo)
	// h := NewHTTPHandler(svc)

	r := setupBasicRouter()

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
			log.Fatal()
		}
	}()

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Server exited")

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
