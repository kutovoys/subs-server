package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"subs-server/internal/config"
	"subs-server/internal/provider"
	"subs-server/internal/provider/filesystem"
	"subs-server/internal/server"
)

var (
	version = "unknown"
)

func main() {
	config.Parse(version)

	if !config.CLIConfig.Debug {
		log.SetFlags(0)
	}

	var fileProvider provider.Provider
	var err error

	switch config.CLIConfig.Source {
	case "filesystem":
		fileProvider, err = filesystem.NewFilesystemProvider(config.CLIConfig.Location, config.CLIConfig.Debug)
		if err != nil {
			log.Fatalf("failed to create file provider: %v", err)
		}
	default:
		log.Fatalf("unknown source type: %s", config.CLIConfig.Source)
	}

	srv := server.NewServer(fileProvider, config.CLIConfig.Debug, &config.CLIConfig)

	ctx2, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		if err := fileProvider.Watch(ctx2); err != nil {
			log.Printf("Provider stopped: %v", err)
		}
	}()

	if err := fileProvider.LoadExistingFiles(); err != nil {
		log.Fatalf("failed to load existing files: %v", err)
	}

	httpServer := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", config.CLIConfig.Host, config.CLIConfig.Port),
		Handler: srv,
	}

	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		log.Println("Shutting down server...")
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer shutdownCancel()

		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			log.Printf("Server shutdown error: %v", err)
		}
		cancel()
	}()

	log.Printf("Starting server on %s:%d", config.CLIConfig.Host, config.CLIConfig.Port)
	if config.CLIConfig.Debug {
		log.Printf("Debug mode enabled")
		log.Printf("Source type: %s", config.CLIConfig.Source)
		log.Printf("Path: %s", config.CLIConfig.Location)
		log.Printf("Response headers configured:")
		log.Printf("  profile-title: base64:%s", base64.StdEncoding.EncodeToString([]byte(config.CLIConfig.ProfileTitle)))
		log.Printf("  profile-update-interval: %s", config.CLIConfig.ProfileUpdateInterval)
		log.Printf("  profile-web-page-url: %s", config.CLIConfig.ProfileWebPageURL)
		log.Printf("  support-url: %s", config.CLIConfig.SupportURL)
	}

	if err := httpServer.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("failed to start server: %v", err)
	}
}
