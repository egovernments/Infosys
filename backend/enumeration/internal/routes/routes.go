// Package routes defines all API route registrations for the service
package routes

import (
	"enumeration/internal/handlers"
	"enumeration/internal/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRoutes sets up all application routes
func SetupRoutes(router *gin.Engine, coordinatesHandler *handlers.CoordinatesHandler, floorDetailsHandler *handlers.FloorDetailsHandler, constructionDetailsHandler *handlers.ConstructionDetailsHandler, additionalPropertyDetailsHandler *handlers.AdditionalPropertyDetailsHandler, assessmentDetailsHandler *handlers.AssessmentDetailsHandler, propertyHandler *handlers.PropertyHandler, propertyAddressHandler *handlers.PropertyAddressHandler, gisHandler *handlers.GISHandler, applicationHandler *handlers.ApplicationHandler, propertyOwnerHandler *handlers.PropertyOwnerHandler, amenityHandler *handlers.AmenityHandler, documentHandler *handlers.DocumentHandler, igrsHandler *handlers.IGRSHandler, applicationLogHandler *handlers.ApplicationLogHandler) {
	v1 := router.Group("/v1", middleware.AuthMiddleware())
	{
		SetupCoordinatesRoutes(v1, coordinatesHandler)
		SetupFloorDetailsRoutes(v1, floorDetailsHandler)
		SetupApplicationRoutes(v1, applicationHandler)
		SetupPropertyOwnerRoutes(v1, propertyOwnerHandler)
		SetupConstructionDetailsRoutes(v1, constructionDetailsHandler)
		SetupAdditionalPropertyDetailsRoutes(v1, additionalPropertyDetailsHandler)
		SetupAssessmentDetailsRoutes(v1, assessmentDetailsHandler)
		SetupPropertyRoutes(v1, propertyHandler)
		SetupPropertyAddressRoutes(v1, propertyAddressHandler)
		SetupGISRoutes(v1, gisHandler)
		SetupAmenitiesRoutes(v1, amenityHandler)
		SetupDocumentsRoutes(v1, documentHandler)
		SetupIGRSRoutes(v1, igrsHandler)
		SetupApplicationLogRoutes(v1, applicationLogHandler)
	}
}

// RegisterApplicationLogRoutes registers all application log routes
func SetupApplicationLogRoutes(router *gin.RouterGroup, handler *handlers.ApplicationLogHandler) {
	v1 := router.Group("/application-logs")
	{
		// Application log endpoints
		v1.GET("", handler.GetApplicationLogs)
		v1.POST("", handler.CreateApplicationLog)
		v1.GET("/:id", handler.GetApplicationLogByID)
		v1.PUT("/:id", handler.UpdateApplicationLog)
		v1.DELETE("/:id", handler.DeleteApplicationLog)

		// Get logs by application ID
		v1.GET("/logs/:id", handler.GetApplicationLogsByApplicationID)
	}
}

// SetupApplicationRoutes registers routes for application endpoints
func SetupApplicationRoutes(router *gin.RouterGroup, applicationHandler *handlers.ApplicationHandler) {
	applications := router.Group("/applications", middleware.MDMSRoleMiddleware("POST", "/v1/applications"))
	{
		applications.GET("/:id", middleware.MDMSRoleMiddleware("GET", "/v1/applications/:id"), applicationHandler.GetByID)
		applications.GET("", middleware.MDMSRoleMiddleware("GET", "/v1/applications"), applicationHandler.List)
		applications.GET("/search", middleware.MDMSRoleMiddleware("GET", "/v1/applications/search"), applicationHandler.Search)
		applications.GET("/", middleware.MDMSRoleMiddleware("GET", "/v1/applications/:applicationNo"), applicationHandler.GetByApplicationNo)

		applications.POST("", middleware.MDMSRoleMiddleware("POST", "/v1/applications"), applicationHandler.Create)

		applications.PUT("/:id", middleware.MDMSRoleMiddleware("PUT", "/v1/applications/:id"), applicationHandler.Update)
		applications.PATCH("/:id", middleware.ActionMiddleware(), applicationHandler.Action)

		applications.DELETE("/:id", middleware.MDMSRoleMiddleware("DELETE", "/v1/applications/:id"), applicationHandler.Delete)
	}
}

// SetupPropertyOwnerRoutes registers routes for property owner endpoints
func SetupPropertyOwnerRoutes(router *gin.RouterGroup, propertyOwnerHandler *handlers.PropertyOwnerHandler) {
	propertyOwners := router.Group("/property-owners", middleware.MDMSRoleMiddleware("POST", "/v1/applications"))
	{
		propertyOwners.POST("", propertyOwnerHandler.Create)                              // POST /property-owners - Create property owner
		propertyOwners.POST("/batch", propertyOwnerHandler.CreateBatch)                   // POST /property-owners/batch - Create owners in batch
		propertyOwners.GET("/property/:propertyId", propertyOwnerHandler.GetByPropertyID) // GET /property-owners/property/{propertyId}
		propertyOwners.PUT("/:id", propertyOwnerHandler.Update)                           // PUT /property-owners/{id}
		propertyOwners.DELETE("/:id", propertyOwnerHandler.Delete)                        // DELETE /property-owners/{id}
	}
}

// SetupConstructionDetailsRoutes registers routes for construction details endpoints
func SetupConstructionDetailsRoutes(router *gin.RouterGroup, constructionDetailsHandler *handlers.ConstructionDetailsHandler) {
	constructionDetails := router.Group("/construction-details", middleware.MDMSRoleMiddleware("POST", "/v1/applications"))
	{
		constructionDetails.POST("", constructionDetailsHandler.CreateConstructionDetails)
		constructionDetails.GET("", constructionDetailsHandler.GetAllConstructionDetails)
		constructionDetails.GET("/:id", constructionDetailsHandler.GetConstructionDetailsByID)
		constructionDetails.PUT("/:id", constructionDetailsHandler.UpdateConstructionDetails)
		constructionDetails.DELETE("/:id", constructionDetailsHandler.DeleteConstructionDetails)
		constructionDetails.GET("/property/:propertyId", constructionDetailsHandler.GetConstructionDetailsByPropertyID)
	}
}

func SetupFloorDetailsRoutes(router *gin.RouterGroup, floorDetailsHandler *handlers.FloorDetailsHandler) {
	floorDetails := router.Group("/floor-details", middleware.MDMSRoleMiddleware("POST", "/v1/applications"))
	{
		floorDetails.POST("", floorDetailsHandler.CreateFloorDetails)
		floorDetails.GET("", floorDetailsHandler.GetAllFloorDetails)
		floorDetails.GET("/:id", floorDetailsHandler.GetFloorDetailsByID)
		floorDetails.PUT("/:id", floorDetailsHandler.UpdateFloorDetails)
		floorDetails.DELETE("/:id", floorDetailsHandler.DeleteFloorDetails)
	}
}

func SetupAdditionalPropertyDetailsRoutes(router *gin.RouterGroup, handler *handlers.AdditionalPropertyDetailsHandler) {
	additionalDetails := router.Group("/additional-property-details", middleware.MDMSRoleMiddleware("POST", "/v1/applications"))
	{
		additionalDetails.POST("", handler.CreateAdditionalPropertyDetails)
		additionalDetails.GET("", handler.GetAllAdditionalPropertyDetails)
		additionalDetails.GET("/:id", handler.GetAdditionalPropertyDetailsByID)
		additionalDetails.PUT("/:id", handler.UpdateAdditionalPropertyDetails)
		additionalDetails.DELETE("/:id", handler.DeleteAdditionalPropertyDetails)
		additionalDetails.GET("/property/:propertyId", handler.GetAdditionalPropertyDetailsByPropertyID)
		additionalDetails.GET("/field/:fieldName", handler.GetAdditionalPropertyDetailsByFieldName)
	}
}

// SetupAssessmentDetailsRoutes registers routes for assessment details endpoints
func SetupAssessmentDetailsRoutes(router *gin.RouterGroup, assessmentDetailsHandler *handlers.AssessmentDetailsHandler) {
	assessmentDetails := router.Group("/assessment-details", middleware.MDMSRoleMiddleware("POST", "/v1/applications"))
	{
		assessmentDetails.POST("", assessmentDetailsHandler.CreateAssessmentDetails)
		assessmentDetails.GET("", assessmentDetailsHandler.GetAllAssessmentDetails)
		assessmentDetails.GET("/:id", assessmentDetailsHandler.GetAssessmentDetailsByID)
		assessmentDetails.PUT("/:id", assessmentDetailsHandler.UpdateAssessmentDetails)
		assessmentDetails.DELETE("/:id", assessmentDetailsHandler.DeleteAssessmentDetails)
		assessmentDetails.GET("/property/:propertyId", assessmentDetailsHandler.GetAssessmentDetailsByPropertyID)
	}
}

// SetupPropertyRoutes registers routes for property endpoints
func SetupPropertyRoutes(router *gin.RouterGroup, propertyHandler *handlers.PropertyHandler) {
	properties := router.Group("/properties", middleware.MDMSRoleMiddleware("POST", "/v1/applications"))
	{
		properties.POST("", propertyHandler.CreateProperty)
		properties.GET("", propertyHandler.GetAllProperties)
		properties.GET("/search", propertyHandler.SearchProperties)
		properties.GET("/:id", propertyHandler.GetPropertyByID)
		properties.PUT("/:id", propertyHandler.UpdateProperty)
		properties.DELETE("/:id", propertyHandler.DeleteProperty)
		properties.GET("/property-no/:propertyNo", propertyHandler.GetPropertyByPropertyNo)
	}
}

// SetupPropertyAddressRoutes registers routes for property address endpoints
func SetupPropertyAddressRoutes(router *gin.RouterGroup, propertyAddressHandler *handlers.PropertyAddressHandler) {
	propertyAddresses := router.Group("/property-addresses", middleware.MDMSRoleMiddleware("POST", "/v1/applications"))
	{
		propertyAddresses.POST("", propertyAddressHandler.CreatePropertyAddress)
		propertyAddresses.GET("", propertyAddressHandler.GetAllPropertyAddresses)
		propertyAddresses.GET("/search", propertyAddressHandler.SearchPropertyAddresses)
		propertyAddresses.GET("/:id", propertyAddressHandler.GetPropertyAddressByID)
		propertyAddresses.PUT("/:id", propertyAddressHandler.UpdatePropertyAddress)
		propertyAddresses.DELETE("/:id", propertyAddressHandler.DeletePropertyAddress)
		propertyAddresses.GET("/property/:propertyId", propertyAddressHandler.GetPropertyAddressByPropertyID)
	}
}

// SetupGISRoutes registers routes for GIS data endpoints
func SetupGISRoutes(router *gin.RouterGroup, gisHandler *handlers.GISHandler) {
	gisData := router.Group("/gis-data", middleware.MDMSRoleMiddleware("POST", "/v1/applications"))
	{
		gisData.GET("", gisHandler.GetAll)                               // GET /gis-data - Get all GIS data
		gisData.POST("", gisHandler.Create)                              // POST /gis-data - Create GIS data
		gisData.GET("/:id", gisHandler.GetByID)                          // GET /gis-data/{id} - Get GIS data by ID
		gisData.PUT("/:id", gisHandler.Update)                           // PUT /gis-data/{id} - Update GIS data
		gisData.DELETE("/:id", gisHandler.Delete)                        // DELETE /gis-data/{id} - Delete GIS data
		gisData.GET("/property/:propertyId", gisHandler.GetByPropertyID) // GET /gis-data/property/{propertyId} - Get GIS data by property ID
	}
}

// SetupCoordinatesRoutes sets up the routes for coordinates endpoints
func SetupCoordinatesRoutes(router *gin.RouterGroup, coordinatesHandler *handlers.CoordinatesHandler) {
	coordinates := router.Group("/coordinates", middleware.MDMSRoleMiddleware("POST", "/v1/applications"))
	{
		coordinates.GET("", coordinatesHandler.GetAll)        // GET /coordinates - Get all coordinates
		coordinates.POST("", coordinatesHandler.Create)       // POST /coordinates - Create coordinates
		coordinates.GET("/:id", coordinatesHandler.GetByID)   // GET /coordinates/{id} - Get coordinates by ID
		coordinates.PUT("/:id", coordinatesHandler.Update)    // PUT /coordinates/{id} - Update coordinates
		coordinates.DELETE("/:id", coordinatesHandler.Delete) // DELETE /coordinates/{id} - Delete coordinates
		coordinates.POST("/batch", coordinatesHandler.CreateBatch)
		coordinates.PUT("/gis/:gisDataId", coordinatesHandler.ReplaceByGISDataID)
	}
}

func SetupAmenitiesRoutes(router *gin.RouterGroup, amenityHandler *handlers.AmenityHandler) {
	amenities := router.Group("/amenities", middleware.MDMSRoleMiddleware("POST", "/v1/applications"))
	{
		amenities.GET("", amenityHandler.GetAll)
		amenities.POST("", amenityHandler.Create)
		amenities.GET("/:id", amenityHandler.GetByID)
		amenities.GET("/property/:propertyId", amenityHandler.GetByPropertyID)
		amenities.PUT("/:id", amenityHandler.Update)
		amenities.DELETE("/:id", amenityHandler.Delete)
	}
}

// SetupDocumentsRoutes registers routes for document endpoints
func SetupDocumentsRoutes(router *gin.RouterGroup, handler *handlers.DocumentHandler) {
	docs := router.Group("/documents", middleware.MDMSRoleMiddleware("POST", "/v1/applications"))
	{
		docs.GET("", handler.GetAll)                               // GET /v1/documents
		docs.POST("", handler.CreateBatch)                         // POST /v1/documents (accepts []Document)
		docs.GET("/:id", handler.GetByID)                          // GET /v1/documents/{id}
		docs.GET("/property/:propertyId", handler.GetByPropertyID) // GET /v1/documents/property/{propertyId}
		docs.DELETE("/:id", handler.Delete)
		docs.PUT("/:id", handler.Update) // PUT /v1/documents/{id}
	}
}

func SetupIGRSRoutes(router *gin.RouterGroup, igrsHandler *handlers.IGRSHandler) {
	igrs := router.Group("/igrs", middleware.MDMSRoleMiddleware("POST", "/v1/applications"))
	{
		igrs.GET("", igrsHandler.List)
		igrs.POST("", igrsHandler.Create)

		igrs.GET("/:id", igrsHandler.GetByID)
		igrs.PUT("/:id", igrsHandler.Update)
		igrs.DELETE("/:id", igrsHandler.Delete)
	}
}
