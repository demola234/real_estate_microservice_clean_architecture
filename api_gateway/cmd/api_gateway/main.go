package main

import (
	"context"
	"job_portal/api_gateway/config"
	"job_portal/api_gateway/interfaces/grpc_clients"
	"job_portal/api_gateway/internal/handler"
	auth "job_portal/api_gateway/routes"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	configs, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("Failed to load env file: %v", err)
	}

	// Initialize gRPC client with dynamic address
	authClient, err := grpc_clients.NewAuthenticationClient("127.0.0.1:9091", 20*time.Second)

	if err != nil {
		log.Fatalf("Failed to connect to Authentication service: %v", err)
	}
	defer authClient.Close() // Ensure gRPC connection is closed on shutdown

	// Create a new Gin router
	router := gin.Default()

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authClient)

	// Group routes under /v1
	v1 := router.Group("/v1")
	{
		// Health check endpoint
		v1.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"status": "OK",
			})
		})
	}

	// Define authentication routes
	auth.RegisterRoutes(v1, authHandler)

	// Create an HTTP server with the configured port
	srv := &http.Server{
		Addr:    configs.Port,
		Handler: router,
	}

	// Start the server in a goroutine
	go func() {
		log.Printf("Starting API Gateway at %s...", configs.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not start server: %v", err)
		}
	}()

	// Wait for an interrupt signal to shut down gracefully
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down API Gateway...")

	// Create a timeout context for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("API Gateway stopped gracefully.")
}
