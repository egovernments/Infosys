package repositories

import (
	"context"
	"enumeration/internal/dto"
	"enumeration/internal/models"

	"github.com/google/uuid"
)

// ApplicationLogRepository handles CRUD for ApplicationLog entities
type ApplicationLogRepository interface {
	Create(ctx context.Context, log *models.ApplicationLog) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.ApplicationLog, error)
	Update(ctx context.Context, log *models.ApplicationLog) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, applicationID, action string, page, size int) ([]models.ApplicationLog, int64, error)
	GetByApplicationID(ctx context.Context, applicationID uuid.UUID, action string, page, size int) ([]models.ApplicationLog, int64, error)
}

// FloorDetailsRepository handles CRUD for FloorDetails entities
type FloorDetailsRepository interface {
	Create(ctx context.Context, floorDetails *models.FloorDetails) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.FloorDetails, error)
	Update(ctx context.Context, floorDetails *models.FloorDetails) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetAll(ctx context.Context, page, size int, constructionDetailsID *uuid.UUID) ([]*models.FloorDetails, int64, error)
	GetByConstructionDetailsID(ctx context.Context, constructionDetailsID uuid.UUID) ([]*models.FloorDetails, error)
}

// CoordinatesRepository handles CRUD and batch operations for Coordinates entities
type CoordinatesRepository interface {
	Create(ctx context.Context, coordinates *models.Coordinates) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Coordinates, error)
	Update(ctx context.Context, coordinates *models.Coordinates) error
	Delete(ctx context.Context, id uuid.UUID) error
	FindAll(ctx context.Context, page, size int, gisDataID *uuid.UUID) ([]models.Coordinates, int64, error)
	FindByGISDataID(ctx context.Context, gisDataID uuid.UUID) ([]models.Coordinates, error)
	DeleteByGISDataID(ctx context.Context, gisDataID uuid.UUID) error
	CreateBatch(ctx context.Context, coords []*models.Coordinates) error
	ReplaceByGISDataID(ctx context.Context, gisDataID uuid.UUID, coords []*models.Coordinates) error
}

// ApplicationRepository handles CRUD and search for Application entities
type ApplicationRepository interface {
	Create(ctx context.Context, application *models.Application) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Application, error)
	GetByApplicationNo(ctx context.Context, applicationNo string) (*models.Application, error)
	Update(ctx context.Context, application *models.Application) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, page, size int) ([]models.Application, int64, error)
	Search(ctx context.Context, criteria *dto.ApplicationSearchCriteria, page, size int) ([]*models.Application, int64, error)
	ExistsByApplicationNo(ctx context.Context, applicationNo string) (bool, error)
	GetApplicationsByIDs(ctx context.Context, ids []string, state string, page, size int) ([]models.Application, int64, error)
	GetByTenantIDAndStatus(ctx context.Context, tenantID string, status string, page, size int) ([]models.Application, int64, error)
	GetByAssignedAgent(ctx context.Context, agentID string, status string, page, size int) ([]models.Application, int64, error)
	CheckUserExists(ctx context.Context, userId string) error
}

// PropertyOwnerRepository handles CRUD and batch operations for PropertyOwner entities
type PropertyOwnerRepository interface {
	Create(ctx context.Context, owner *models.PropertyOwner) error
	CreateBatch(ctx context.Context, owners []*models.PropertyOwner) error
	GetAll(ctx context.Context) ([]*models.PropertyOwner, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.PropertyOwner, error)
	GetByPropertyID(ctx context.Context, propertyID uuid.UUID, page, size int) ([]*models.PropertyOwner, int64, error)
	Update(ctx context.Context, owner *models.PropertyOwner) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// ConstructionDetailsRepository handles CRUD for ConstructionDetails entities
type ConstructionDetailsRepository interface {
	Create(ctx context.Context, constructionDetails *models.ConstructionDetails) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.ConstructionDetails, error)
	Update(ctx context.Context, constructionDetails *models.ConstructionDetails) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetAll(ctx context.Context, page, size int, propertyID *uuid.UUID) ([]*models.ConstructionDetails, int64, error)
	GetByPropertyID(ctx context.Context, propertyID uuid.UUID) ([]*models.ConstructionDetails, error)
}

// AdditionalPropertyDetailsRepository handles CRUD and filtering for AdditionalPropertyDetails entities
type AdditionalPropertyDetailsRepository interface {
	Create(ctx context.Context, details *models.AdditionalPropertyDetails) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.AdditionalPropertyDetails, error)
	Update(ctx context.Context, details *models.AdditionalPropertyDetails) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetAll(ctx context.Context, page, size int, propertyID *uuid.UUID, fieldName *string) ([]*models.AdditionalPropertyDetails, int64, error)
	GetByPropertyID(ctx context.Context, propertyID uuid.UUID) ([]*models.AdditionalPropertyDetails, error)
	GetByFieldName(ctx context.Context, fieldName string) ([]*models.AdditionalPropertyDetails, error)
}

// AssessmentDetailsRepository handles CRUD for AssessmentDetails entities
type AssessmentDetailsRepository interface {
	Create(ctx context.Context, assessmentDetails *models.AssessmentDetails) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.AssessmentDetails, error)
	Update(ctx context.Context, assessmentDetails *models.AssessmentDetails) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetAll(ctx context.Context, page, size int, propertyID *uuid.UUID) ([]*models.AssessmentDetails, int64, error)
	GetByPropertyID(ctx context.Context, propertyID uuid.UUID) (*models.AssessmentDetails, error)
}

// PropertyRepository handles CRUD and search for Property entities
type PropertyRepository interface {
	Create(ctx context.Context, property *models.Property) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Property, error)
	Update(ctx context.Context, property *models.Property) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetAll(ctx context.Context, page, size int, propertyType *string) ([]*models.Property, int64, error)
	GetByPropertyNo(ctx context.Context, propertyNo string) (*models.Property, error)
	Search(ctx context.Context, params SearchPropertyParams) ([]*models.Property, int64, error)
}

// SearchPropertyParams holds filters and options for property search
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
	// SearchPropertyParams holds parameters for searching properties with various filters and sorting options.
	SortOrder string
}

// PropertyAddressRepository handles CRUD and search for PropertyAddress entities
type PropertyAddressRepository interface {
	Create(address *models.PropertyAddress) error
	GetByID(id uuid.UUID) (*models.PropertyAddress, error)
	Update(address *models.PropertyAddress) error
	Delete(id uuid.UUID) error
	GetAll(page, size int, propertyID *uuid.UUID) ([]*models.PropertyAddress, int64, error)
	GetByPropertyID(propertyID uuid.UUID) (*models.PropertyAddress, error)
	Search(params SearchPropertyAddressParams) ([]*models.PropertyAddress, int64, error)
}

// SearchPropertyAddressParams holds filters and options for property address search
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

// GISRepository handles CRUD and pagination for GISData entities
type GISRepository interface {
	Create(ctx context.Context, gisData *models.GISData) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.GISData, error)
	GetByPropertyID(ctx context.Context, propertyID uuid.UUID) (*models.GISData, error)
	Update(ctx context.Context, gisData *models.GISData) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetAll(ctx context.Context, page, size int) ([]*models.GISData, int64, error)
}

// AmenityRepository handles CRUD and filtering for Amenities entities
type AmenityRepository interface {
	GetAll(ctx context.Context) ([]models.Amenities, error)
	GetAllWithFilters(ctx context.Context, page, size int, amenityType, propertyID string) ([]models.Amenities, int64, error)
	GetByID(ctx context.Context, id string) (*models.Amenities, error)
	Create(ctx context.Context, amenity *models.Amenities) error
	Update(ctx context.Context, id string, amenity *models.Amenities) error
	Delete(ctx context.Context, id string) error
	GetByPropertyID(ctx context.Context, propertyID string) (*models.Amenities, error)
}

// DocumentRepository handles CRUD and batch operations for Document entities
type DocumentRepository interface {
	Create(ctx context.Context, doc *models.Document) error
	CreateBatch(ctx context.Context, docs []*models.Document) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Document, error)
	GetByPropertyID(ctx context.Context, propertyID uuid.UUID) ([]*models.Document, error)
	GetAll(ctx context.Context, page, size int, propertyID *uuid.UUID) ([]*models.Document, int64, error)
	Delete(ctx context.Context, id uuid.UUID) error
	Update(ctx context.Context, doc *models.Document) error
}

// IGRSRepository handles CRUD and pagination for IGRS entities
type IGRSRepository interface {
	Create(ctx context.Context, igrs *models.IGRS) error
	GetByPropertyID(ctx context.Context, propertyID uuid.UUID) (*models.IGRS, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.IGRS, error)
	Update(ctx context.Context, igrs *models.IGRS) error
	Delete(ctx context.Context, id uuid.UUID) error
	FindAll(ctx context.Context, page, size int) ([]models.IGRS, int64, error)
}
