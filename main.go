package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/ritikchawla/net-canary/internal/config"
	"github.com/ritikchawla/net-canary/internal/logging"
	"github.com/ritikchawla/net-canary/internal/services"
)

var (
	configFile = flag.String("config", "config.yaml", "path to configuration file")
)

func main() {
	flag.Parse()

	// Load configuration
	cfg, err := config.Load(*configFile)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize logger
	logger, err := logging.NewLogger(cfg.Logging.File)
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Close()

	// Context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize services
	var activeServices []services.Service

	// Initialize SSH service if enabled
	if cfg.Services.SSH.Enabled {
		sshService := services.NewSSHService(
			cfg.Services.SSH.Port,
			cfg.Services.SSH.Host,
			cfg.Services.SSH.Banner,
			logger,
		)
		activeServices = append(activeServices, sshService)
	}

	// Start all services
	var wg sync.WaitGroup
	for _, svc := range activeServices {
		wg.Add(1)
		go func(s services.Service) {
			defer wg.Done()
			if err := s.Start(ctx); err != nil {
				log.Printf("Error starting %s service: %v", s.Name(), err)
			}
		}(svc)
	}

	fmt.Println("NetCanary - A Stealthy Network Threat Tripwire")
	fmt.Printf("Started %d services. Logging to: %s\n", len(activeServices), cfg.Logging.File)

	// Setup signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Wait for shutdown signal
	<-sigChan
	fmt.Println("\nShutting down gracefully...")

	// Cancel context to stop all services
	cancel()

	// Stop all services
	for _, svc := range activeServices {
		if err := svc.Stop(); err != nil {
			log.Printf("Error stopping %s service: %v", svc.Name(), err)
		}
	}

	// Wait for all services to stop
	wg.Wait()
	fmt.Println("All services stopped. Goodbye!")
}
