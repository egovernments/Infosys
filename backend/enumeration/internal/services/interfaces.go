package services

import (
	"context"
	"enumeration/internal/dto"
	"enumeration/internal/models"

	"github.com/google/uuid"
)

// ApplicationLogService defines business logic for application logs
type ApplicationLogService interface {
	Create(ctx context.Context, req *dto.CreateApplicationLogRequest) (*models.ApplicationLog, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.ApplicationLog, error)
	Update(ctx context.Context, id uuid.UUID, req *dto.UpdateApplicationLogRequest) (*models.ApplicationLog, error)
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, applicationID, action string, page, size int) ([]models.ApplicationLog, int64, error)
	GetByApplicationID(ctx context.Context, applicationID uuid.UUID, action string, page, size int) ([]models.ApplicationLog, int64, error)
}

// FloorDetailsService handles business logic for floor details
type FloorDetailsService interface {
	CreateFloorDetails(ctx context.Context, floorDetails *models.FloorDetails) error
	GetFloorDetailsByID(ctx context.Context, id uuid.UUID) (*models.FloorDetails, error)
	UpdateFloorDetails(ctx context.Context, floorDetails *models.FloorDetails) error
	DeleteFloorDetails(ctx context.Context, id uuid.UUID) error
	GetAllFloorDetails(ctx context.Context, page, size int, constructionDetailsID *uuid.UUID) ([]*models.FloorDetails, int64, error)
	GetFloorDetailsByConstructionDetailsID(ctx context.Context, constructionDetailsID uuid.UUID) ([]*models.FloorDetails, error)
}

// CoordinatesService handles business logic for coordinates
type CoordinatesService interface {
	Create(ctx context.Context, coordinates *models.Coordinates) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Coordinates, error)
	Update(ctx context.Context, coordinates *models.Coordinates) error
	Delete(ctx context.Context, id uuid.UUID) error
	FindAll(ctx context.Context, page, size int, gisDataID *uuid.UUID) ([]models.Coordinates, int64, error)
	CreateBatch(ctx context.Context, coords []*models.Coordinates) error
	GetAll(ctx context.Context, page, size int, gisDataID *uuid.UUID) ([]*models.Coordinates, int64, error)
	ReplaceByGISDataID(ctx context.Context, gisDataID uuid.UUID, coords []*models.Coordinates) error
}

// ApplicationService handles business logic for applications
type ApplicationService interface {
	Create(ctx context.Context, tenantID string, citizenID string, req *dto.CreateApplicationRequest) (*models.Application, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.Application, error)
	GetByApplicationNo(ctx context.Context, applicationNo string) (*models.Application, error)
	Update(ctx context.Context, id uuid.UUID, req *dto.UpdateApplicationRequest) (*models.Application, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, req *dto.UpdateStatusRequest) error
	AssignAgent(ctx context.Context, id uuid.UUID, req *dto.ActionRequest, tenantID string) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, UserID string, tenantID string, Role string, verify string, page, size int) ([]models.Application, int64, error)
	Search(ctx context.Context, criteria *dto.ApplicationSearchCriteria, page, size int) ([]*models.Application, int64, error)
	VerifyApplicationByAgent(ctx context.Context, tenantID, agentID, applicationID string, req *dto.ActionRequest) error
	VerifyApplication(ctx context.Context, tenantID, serviceManagerID, applicationID string, req *dto.ActionRequest) error
	ApproveApplication(ctx context.Context, tenantID, commissionerID, applicationID string, req *dto.ActionRequest) error
	// ReassignAgent(ctx context.Context, id uuid.UUID, req *dto.ActionRequest, tenantID string) error
}

// PropertyOwnerService handles business logic for property owners
type PropertyOwnerService interface {
	Create(ctx context.Context, req *dto.CreatePropertyOwnerRequest) (*models.PropertyOwner, error)
	CreatePropertyOwners(ctx context.Context, reqs []*dto.CreatePropertyOwnerRequest) ([]*models.PropertyOwner, error)
	GetByPropertyID(ctx context.Context, propertyID uuid.UUID, page, size int) ([]*models.PropertyOwner, int64, error)
	Update(ctx context.Context, id uuid.UUID, req *dto.UpdatePropertyOwnerRequest) (*models.PropertyOwner, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

// ConstructionDetailsService handles business logic for construction details
type ConstructionDetailsService interface {
	CreateConstructionDetails(ctx context.Context, constructionDetails *models.ConstructionDetails) error
	GetConstructionDetailsByID(ctx context.Context, id uuid.UUID) (*models.ConstructionDetails, error)
	UpdateConstructionDetails(ctx context.Context, constructionDetails *models.ConstructionDetails) error
	DeleteConstructionDetails(ctx context.Context, id uuid.UUID) error
	GetAllConstructionDetails(ctx context.Context, page, size int, propertyID *uuid.UUID) ([]*models.ConstructionDetails, int64, error)
	GetConstructionDetailsByPropertyID(ctx context.Context, propertyID uuid.UUID) ([]*models.ConstructionDetails, error)
}

// AdditionalPropertyDetailsService handles business logic for additional property details
type AdditionalPropertyDetailsService interface {
	CreateAdditionalPropertyDetails(ctx context.Context, details *models.AdditionalPropertyDetails) error
	GetAdditionalPropertyDetailsByID(ctx context.Context, id uuid.UUID) (*models.AdditionalPropertyDetails, error)
	UpdateAdditionalPropertyDetails(ctx context.Context, details *models.AdditionalPropertyDetails) error
	DeleteAdditionalPropertyDetails(ctx context.Context, id uuid.UUID) error
	GetAllAdditionalPropertyDetails(ctx context.Context, page, size int, propertyID *uuid.UUID, fieldName *string) ([]*models.AdditionalPropertyDetails, int64, error)
	GetAdditionalPropertyDetailsByPropertyID(ctx context.Context, propertyID uuid.UUID) ([]*models.AdditionalPropertyDetails, error)
	GetAdditionalPropertyDetailsByFieldName(ctx context.Context, fieldName string) ([]*models.AdditionalPropertyDetails, error)
}

// AssessmentDetailsService handles business logic for assessment details
type AssessmentDetailsService interface {
	CreateAssessmentDetails(ctx context.Context, assessmentDetails *models.AssessmentDetails) error
	GetAssessmentDetailsByID(ctx context.Context, id uuid.UUID) (*models.AssessmentDetails, error)
	UpdateAssessmentDetails(ctx context.Context, assessmentDetails *models.AssessmentDetails) error
	DeleteAssessmentDetails(ctx context.Context, id uuid.UUID) error
	GetAllAssessmentDetails(ctx context.Context, page, size int, propertyID *uuid.UUID) ([]*models.AssessmentDetails, int64, error)
	GetAssessmentDetailsByPropertyID(ctx context.Context, propertyID uuid.UUID) (*models.AssessmentDetails, error)
}

// PropertyService handles business logic for properties
type PropertyService interface {
	CreateProperty(ctx context.Context, property *models.Property) error
	GetPropertyByID(ctx context.Context, id uuid.UUID) (*models.Property, error)
	UpdateProperty(ctx context.Context, property *models.Property) error
	DeleteProperty(ctx context.Context, id uuid.UUID) error
	GetAllProperties(ctx context.Context, page, size int, propertyType *string) ([]*models.Property, int64, error)
	GetPropertyByPropertyNo(ctx context.Context, propertyNo string) (*models.Property, error)
	SearchProperties(ctx context.Context, params SearchPropertyParams) ([]*models.Property, int64, error)
	GeneratePropertyNo(ctx context.Context) (string, error)
}

// SearchPropertyParams holds search and filter options for properties
type SearchPropertyParams struct {
	Page          int
	Size          int
	PropertyType  *string
	OwnershipType *string
	ComplexName   *string
	Locality      *string
	WardNo        *string
	ZoneNo        *string
	Street        *string
	SortBy        string
	SortOrder     string
}

// PropertyAddressService handles business logic for property addresses
type PropertyAddressService interface {
	CreatePropertyAddress(ctx context.Context, address *models.PropertyAddress) error
	GetPropertyAddressByID(ctx context.Context, id uuid.UUID) (*models.PropertyAddress, error)
	UpdatePropertyAddress(ctx context.Context, address *models.PropertyAddress) error
	DeletePropertyAddress(ctx context.Context, id uuid.UUID) error
	GetAllPropertyAddresses(ctx context.Context, page, size int, propertyID *uuid.UUID) ([]*models.PropertyAddress, int64, error)
	GetPropertyAddressByPropertyID(ctx context.Context, propertyID uuid.UUID) (*models.PropertyAddress, error)
	SearchPropertyAddresses(ctx context.Context, params SearchPropertyAddressParams) ([]*models.PropertyAddress, int64, error)
}

// SearchPropertyAddressParams holds search and filter options for property addresses
type SearchPropertyAddressParams struct {
	Page            int
	Size            int
	PropertyID      *uuid.UUID
	Locality        *string
	ZoneNo          *string
	WardNo          *string
	BlockNo         *string
	Street          *string
	ElectionWard    *string
	SecretariatWard *string
	PinCode         *uint64
	SortBy          string
	SortOrder       string
}

// GISService handles business logic for GIS data
type GISService interface {
	CreateGISData(ctx context.Context, gisData *models.GISData) error
	GetGISDataByID(ctx context.Context, id uuid.UUID) (*models.GISData, error)
	GetGISDataByPropertyID(ctx context.Context, propertyID uuid.UUID) (*models.GISData, error)
	UpdateGISData(ctx context.Context, gisData *models.GISData) error
	DeleteGISData(ctx context.Context, id uuid.UUID) error
	GetAllGISData(ctx context.Context, page, size int) ([]*models.GISData, int64, error)
}

// AmenityService handles business logic for amenities
type AmenityService interface {
	// Paginated GetAll: returns items and total count.
	GetAll(ctx context.Context, page, size int) ([]models.Amenities, int64, error)
	GetAllWithFilters(ctx context.Context, page, size int, amenityType, propertyID string) ([]models.Amenities, int64, error)
	GetByID(ctx context.Context, id string) (*models.Amenities, error)
	Create(ctx context.Context, amenity *models.Amenities) error
	Update(ctx context.Context, id string, amenity *models.Amenities) error
	Delete(ctx context.Context, id string) error
	GetByPropertyID(ctx context.Context, propertyID string) (*models.Amenities, error)
}

// DocumentService handles business logic for documents
type DocumentService interface {
	CreateDocument(ctx context.Context, doc *models.Document) error
	CreateDocuments(ctx context.Context, docs []*models.Document) error
	GetDocumentByID(ctx context.Context, id uuid.UUID) (*models.Document, error)
	GetDocumentsByPropertyID(ctx context.Context, propertyID uuid.UUID) ([]*models.Document, error)
	GetAllDocuments(ctx context.Context, page, size int, propertyID *uuid.UUID) ([]*models.Document, int64, error)
	DeleteDocument(ctx context.Context, id uuid.UUID) error
	UpdateDocument(ctx context.Context, doc *models.Document) error
}

// IGRSService handles business logic for IGRS data
type IGRSService interface {
	Create(ctx context.Context, req *dto.CreateIGRSRequest) (*models.IGRS, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.IGRS, error)
	Update(ctx context.Context, id uuid.UUID, req *dto.UpdateIGRSRequest) (*models.IGRS, error)
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, page, size int) ([]*models.IGRS, int64, error)
}
