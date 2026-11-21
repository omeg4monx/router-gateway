package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"
	"time"

	"github.com/omeg4monx/router-gateway/internal/config"
	"github.com/omeg4monx/router-gateway/internal/database"
	"github.com/omeg4monx/router-gateway/internal/gateway"
)

const banner = "OMEG4MON ROUTER · Go Gateway\n"

func main() {
	// Optimize runtime memory scavenger: reclaim heap aggressively back to the OS
	debug.SetGCPercent(50)
	if os.Getenv("GOMEMLIMIT") == "" {
		debug.SetMemoryLimit(64 * 1024 * 1024) // 64MiB soft limit
	}

	fmt.Print(banner)
	cfg := config.Load()

	log.Printf("[OMEG4MON ROUTER] Initializing OMEG4MON ROUTER (Go Core v%s)...", cfg.Version)
	log.Printf("[OMEG4MON ROUTER] Data Directory: %s", cfg.DataDir)
	log.Printf("[OMEG4MON ROUTER] Port: %s", cfg.Port)

	db, err := database.Open(cfg.DBPath)
	if err != nil {
		log.Fatalf("[OMEG4MON ROUTER FATAL] Database failed to open: %v", err)
	}
	defer db.Close()

	gw := gateway.New(cfg, db)

	server := &http.Server{
		Addr:         net.JoinHostPort(cfg.Host, cfg.Port),
		Handler:      gw,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 0, // 0 for streaming SSE
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("[OMEG4MON ROUTER] Server listening on http://%s", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("[OMEG4MON ROUTER FATAL] ListenAndServe failed: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("[OMEG4MON ROUTER] Shutting down server gracefully...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("[OMEG4MON ROUTER ERROR] Server forced to shutdown: %v", err)
	}

	log.Println("[OMEG4MON ROUTER] Server exited cleanly.")
}
