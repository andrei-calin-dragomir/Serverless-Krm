package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"scheduler/internal/inbound"
	"scheduler/internal/outbound"
	"scheduler/internal/service"
	"strings"
	"syscall"
	"time"
)

var (
	PORT           string
	DOMAIN         string
	ETCD_ENDPOINTS []string
)

func init() {
	PORT = os.Getenv("PORT")
	DOMAIN = os.Getenv("DOMAIN")
	ETCD_ENDPOINTS = strings.Split(os.Getenv("ETCD_ENDPOINTS"), ",")
	// fmt.Print(PORT)
	// fmt.Print(DOMAIN)
	// fmt.Print(ETCD_ENDPOINTS)
}

// var (
// 	etcdEndpoints = []string{"localhost:2379"}
// )

func monitorEtcdConnection(client *outbound.ETCDClient) {
	for {
		if !client.IsConnected() {
			log.Println("Lost connection to etcd. Retrying...")
			err := client.RetryConnection()
			if err != nil {
				log.Printf("Retry failed: %v. Retrying in 10 seconds...", err)
			} else {
				log.Println("Reconnected to etcd.")
			}
		}
		time.Sleep(10 * time.Second) // Adjust sleep duration based on your needs
	}
}

func main() {

	etcdClient, err := outbound.NewETCDClient(ETCD_ENDPOINTS, 5*time.Second) //, 10*time.Second)
	if err != nil {
		log.Fatalf("Failed to initialize etcd client: %v", err)
	}
	defer etcdClient.Close()

	// Start a Goroutine to monitor connection status
	go monitorEtcdConnection(etcdClient)

	schedulerService := service.NewSchedulerService(etcdClient)
	schedulerHandler := inbound.NewSchedulerHandler(schedulerService)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /scheduler/pod", schedulerHandler.HandlePodBind)

	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", DOMAIN, PORT),
		Handler: mux,
	}
	// Channel to listen for OS signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start the server in a separate goroutine
	go func() {
		slog.Info("Starting server", slog.String("address", server.Addr))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("HTTP server failed", slog.String("error", err.Error()))
		}
	}()

	// Wait for termination signal
	<-sigChan
	slog.Warn("Shutdown signal received, attempting graceful shutdown...")

	// Attempt to gracefully shut down the server
	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownRelease()

	if err := server.Shutdown(shutdownCtx); err != nil {
		slog.Error("Failed to gracefully shut down the server", slog.String("error", err.Error()))
		return
	}
	slog.Info("Server shut down gracefully.")
}
