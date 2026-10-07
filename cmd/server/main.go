package main

import (
	"context"
	"job-tracker-api/internal/config"
	"job-tracker-api/internal/handler"
	"job-tracker-api/internal/repository"
	"job-tracker-api/internal/service"
	"log"
	"net/http"

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

	router := gin.Default()

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
	}

	if err := router.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
