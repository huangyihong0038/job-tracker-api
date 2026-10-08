package main

import (
	"context"
	"errors"
	"job-tracker-api/internal/config"
	"job-tracker-api/internal/handler"
	"job-tracker-api/internal/repository"
	"job-tracker-api/internal/service"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"job-tracker-api/internal/middleware"

	"github.com/gin-gonic/gin"
)

func main() {

	ctx := context.Background()
	db, err := config.OpenMySQL(ctx)
	if err != nil {
		log.Fatalf("database connection failed: %v", err)
	}
	defer db.Close()

	applicationRepository := repository.NewMySQLApplicationRepository(db)

	applicationService := service.NewApplicationService(applicationRepository)

	applicationHandler := handler.NewApplicationHandler(applicationService)

	router := gin.New()

	router.Use(
		middleware.RequestLogger(
			log.Default(),
		),
		gin.Recovery(),
	)

	router.HandleMethodNotAllowed = true

	router.NoRoute(
		handler.NotFound,
	)

	router.NoMethod(
		handler.MethodNotAllowed,
	)

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	api := router.Group("/api/v1")
	{
		api.POST(
			"/applications",
			applicationHandler.Create,
		)

		api.GET(
			"/applications/:id",
			applicationHandler.GetByID,
		)

		api.DELETE(
			"/applications/:id",
			applicationHandler.Delete,
		)

		api.GET(
			"/applications",
			applicationHandler.List,
		)

		api.PATCH(
			"/applications/:id",
			applicationHandler.Update,
		)
	}

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	server := &http.Server{
		Addr:              ":" + port,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	serverError := make(chan error, 1)

	go func() {
		log.Printf(
			"server listening on %s",
			server.Addr,
		)

		serverError <- server.ListenAndServe()
	}()

	signalContext, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	select {
	case err := <-serverError:
		if err != nil &&
			!errors.Is(
				err,
				http.ErrServerClosed,
			) {
			log.Printf(
				"server stopped unexpectedly: %v",
				err,
			)
		}
		return

	case <-signalContext.Done():
		log.Print(
			"shutdown signal received",
		)
	}

	shutdownContext, cancel :=
		context.WithTimeout(
			context.Background(),
			10*time.Second,
		)
	defer cancel()

	if err := server.Shutdown(
		shutdownContext,
	); err != nil {
		log.Printf(
			"graceful shutdown failed: %v",
			err,
		)

		if closeErr := server.Close(); closeErr != nil {
			log.Printf(
				"force close failed: %v",
				closeErr,
			)
		}

		return
	}

	log.Print("server stopped")

}
