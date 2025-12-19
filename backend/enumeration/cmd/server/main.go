// Entry point for the Property Tax Enumeration Service
package main

import (
	workflow "enumeration/internal/clients"
	"enumeration/internal/config"
	"enumeration/internal/database"
	"enumeration/internal/handlers"
	"enumeration/internal/middleware"
	"enumeration/internal/repositories"
	"enumeration/internal/routes"
	"enumeration/internal/services"
	"enumeration/pkg/logger"

	// Standard and third-party packages
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// init loads environment variables and initializes middleware
func init() {
	err := godotenv.Load()
	if err != nil {
		logger.Error("\nError loading .env file using default environment variables\n")
	}
	middleware.Init()
}

// main sets up dependencies and starts the HTTP server
func main() {
	cfg := config.GetConfig() // Load configuration

	logger.Info("Starting Property Tax Enumeration Service...")

	// Initialize workflow client for external workflow integration
	workflowClient := workflow.NewClient(cfg.WorkflowURL)

	// Connect to the database
	db := database.GetDB()

	// Initialize repositories for data access
	coordinatesRepo := repositories.NewCoordinatesRepository(db)
	floorDetailsRepo := repositories.NewFloorDetailsRepository(db)
	applicationRepo := repositories.NewApplicationRepository(db)
	applicationLogRepo:=repositories.NewApplicationLogRepository(db)
	propertyOwnerRepo := repositories.NewPropertyOwnerRepository(db)
	constructionDetailsRepo := repositories.NewConstructionDetailsRepository(db)
	additionalPropertyDetailsRepo := repositories.NewAdditionalPropertyDetailsRepository(db)
	assessmentDetailsRepo := repositories.NewAssessmentDetailsRepository(db)
	propertyRepo := repositories.NewPropertyRepository(db)
	propertyAddressRepo := repositories.NewPropertyAddressRepository(db)
	gisRepo := repositories.NewGISRepository(db)
	amenityRepo := repositories.NewAmenityRepository(db)
	documentRepo := repositories.NewDocumentRepository(db)
	igrsRepo := repositories.NewIGRSRepository(db)

	// Initialize services for business logic
	coordinatesService := services.NewCoordinatesService(coordinatesRepo)
	floorDetailsService := services.NewFloorDetailsService(floorDetailsRepo)
	applicationLogService := services.NewApplicationLogService(applicationLogRepo)
	applicationService := services.NewApplicationService(applicationRepo, propertyOwnerRepo, workflowClient, cfg,applicationLogService)
	

	propertyOwnerService := services.NewPropertyOwnerService(propertyOwnerRepo)
	constructionDetailsService := services.NewConstructionDetailsService(constructionDetailsRepo)
	additionalPropertyDetailsService := services.NewAdditionalPropertyDetailsService(additionalPropertyDetailsRepo)
	assessmentDetailsService := services.NewAssessmentDetailsService(assessmentDetailsRepo)
	propertyService := services.NewPropertyService(propertyRepo)
	propertyAddressService := services.NewPropertyAddressService(propertyAddressRepo)
	gisService := services.NewGISService(gisRepo)
	amenityService := services.NewAmenityService(amenityRepo)
	documentService := services.NewDocumentService(documentRepo)
	igrsService := services.NewIGRSService(igrsRepo)

	// Initialize handlers
	coordinatesHandler := handlers.NewCoordinatesHandler(coordinatesService)
	floorDetailsHandler := handlers.NewFloorDetailsHandler(floorDetailsService)
	applicationHandler := handlers.NewApplicationHandler(applicationService)
	applicationLogHandler:= handlers.NewApplicationLogHandler(applicationLogService)
	propertyOwnerHandler := handlers.NewPropertyOwnerHandler(propertyOwnerService)
	constructionDetailsHandler := handlers.NewConstructionDetailsHandler(constructionDetailsService)
	additionalPropertyDetailsHandler := handlers.NewAdditionalPropertyDetailsHandler(additionalPropertyDetailsService)
	assessmentDetailsHandler := handlers.NewAssessmentDetailsHandler(assessmentDetailsService)
	propertyHandler := handlers.NewPropertyHandler(propertyService)
	propertyAddressHandler := handlers.NewPropertyAddressHandler(propertyAddressService)
	gisHandler := handlers.NewGISHandler(gisService)
	amenityHandler := handlers.NewAmenityHandler(amenityService)
	documentHandler := handlers.NewDocumentHandler(documentService)
	igrsHandler := handlers.NewIGRSHandler(igrsService)

	// Setup Gin router for HTTP requests
	router := gin.Default()

	// Add CORS and logging middleware
	router.Use(middleware.CorsMiddleware(), middleware.LoggerMiddleware())

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "UP",
			"service": "property-tax-enumeration",
		})
	})

	// Setup routes
	routes.SetupRoutes(router, coordinatesHandler, floorDetailsHandler, constructionDetailsHandler, additionalPropertyDetailsHandler, assessmentDetailsHandler, propertyHandler, propertyAddressHandler, gisHandler, applicationHandler, propertyOwnerHandler, amenityHandler, documentHandler, igrsHandler,applicationLogHandler)

	// Start the HTTP server
	address := fmt.Sprintf(":%s", cfg.Port)
	logger.Info("Server starting on", address)
	if err := router.Run(address); err != nil {
		logger.Fatal("Failed to start server:", err)
	}
}
