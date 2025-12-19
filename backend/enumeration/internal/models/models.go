package models

import (
	"enumeration/pkg/utils"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/datatypes"
)

// Application represents a property tax application.
type Application struct {
	ID                 uuid.UUID        `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	ApplicationNo      string           `gorm:"uniqueIndex;size:50"`
	PropertyID         uuid.UUID        `gorm:"type:uuid;not null;uniqueIndex"`
	Priority           string           `gorm:"size:20;not null"`
	TenantID           string           `gorm:"size:100;not null;index"`
	DueDate            time.Time        `gorm:"not null"`
	AssignedAgent      *uuid.UUID       `gorm:"type:uuid;index"`
	Status             string           `gorm:"size:50;index"`
	WorkflowInstanceID string           `gorm:"size:100;index"`
	AppliedBy          string           `gorm:"size:200;not null"`
	AssesseeID         string           `gorm:"type:uuid;index"` // Foreign key to User (external service)
	Property           Property         `gorm:"foreignKey:PropertyID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ApplicationLogs    []ApplicationLog `gorm:"foreignKey:ApplicationID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	IsDraft            bool             `gorm:"default:true;not null;index"`
	CreatedAt          time.Time        `gorm:"autoCreateTime"`
	UpdatedAt          time.Time        `gorm:"autoUpdateTime"`
	ImportantNote		string	  		 `gorm:"size:500" json:"importantNote,omitempty"`
}

func (Application) TableName() string {
	return "DIGIT3.applications"
}

// ApplicationLog represents an action performed on an application.
type ApplicationLog struct {
	ID            uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Action        string     `gorm:"type:text"`
	PerformedBy   string     `gorm:"size:200;not null"`
	PerformedDate time.Time  `gorm:"not null;index"`
	Comments      string     `gorm:"type:text"`
	Actor         string     `gorm:"size:100"`
	Metadata      string     `gorm:"type:json"`
	FileStoreID   *uuid.UUID `gorm:"type:uuid;index"`
	ApplicationID uuid.UUID  `gorm:"type:uuid;not null;index"`
	CreatedAt     time.Time  `gorm:"autoCreateTime"`
}

// Property represents a property entity.
type Property struct {
	ID                  uuid.UUID                  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	PropertyNo          string                     `gorm:"unique;size:50;not null;index"`
	OwnershipType       string                     `gorm:"size:50"`
	PropertyType        string                     `gorm:"size:50"`
	ComplexName         string                     `gorm:"size:200"`
	Address             *PropertyAddress           `gorm:"foreignKey:PropertyID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	AssessmentDetails   *AssessmentDetails         `gorm:"foreignKey:PropertyID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Amenities           *Amenities                 `gorm:"foreignKey:PropertyID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	ConstructionDetails *ConstructionDetails       `gorm:"foreignKey:PropertyID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	AdditionalDetails   *AdditionalPropertyDetails `gorm:"foreignKey:PropertyID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	GISData             *GISData                   `gorm:"foreignKey:PropertyID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	CreatedAt           time.Time                  `gorm:"autoCreateTime"`
	UpdatedAt           time.Time                  `gorm:"autoUpdateTime"`
	Documents           []Document                 `gorm:"foreignKey:PropertyID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	IGRS                *IGRS                      `gorm:"foreignKey:PropertyID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	TypeOfLand          string                     `gorm:"size:50;" json:"typeOfLand,omitempty"`
	NoOfFloors          string                     `gorm:"size:10" json:"noOfFloors,omitempty"`
	NoOfBasements       string                     `gorm:"size:10" json:"noOfBasements,omitempty"`
	NoOfBuildings       string                     `gorm:"size:10" json:"noOfBuildings,omitempty"`
	BuildingNumber      string                     `gorm:"size:50" json:"buildingNumber,omitempty"`
}

func (Property) TableName() string {
	return "DIGIT3.properties"
}

type PropertyOwner struct {
	ID                     uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	PropertyID             uuid.UUID `gorm:"type:uuid;not null;index"`
	AdhaarNo               uint64    `gorm:"not null;index"`
	Name                   string    `gorm:"size:200;not null"`
	ContactNo              string    `gorm:"size:15;not null"`
	Email                  string    `gorm:"size:100"`
	Gender                 string    `gorm:"size:10;not null"`
	Guardian               string    `gorm:"size:200"`
	GuardianType           string    `gorm:"size:10"`
	RelationshipToProperty string    `gorm:"size:20"`
	OwnershipShare         float64   `gorm:"type:decimal(5,2);default:0"`
	IsPrimaryOwner         bool      `gorm:"default:false;index"`
	CreatedAt              time.Time `gorm:"not null;autoCreateTime"`
	UpdatedAt              time.Time `gorm:"not null;autoUpdateTime"`
}

// TableName overrides the table name to use DIGIT3 schema
func (PropertyOwner) TableName() string {
	return "DIGIT3.property_owner"
}

// PropertyAddress represents the address of a property.
type PropertyAddress struct {
	ID                             uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Locality                       string    `gorm:"size:200;not null"`
	ZoneNo                         string    `gorm:"size:50;not null"`
	WardNo                         string    `gorm:"size:50;not null"`
	BlockNo                        string    `gorm:"size:50;not null"`
	Street                         string    `gorm:"size:200"`
	ElectionWard                   string    `gorm:"size:50;not null"`
	SecretariatWard                string    `gorm:"size:50"`
	PinCode                        uint64    `gorm:"index"`
	DifferentCorrespondenceAddress bool      `gorm:"default:false;not null"`
	PropertyID                     uuid.UUID `gorm:"type:uuid;not null;uniqueIndex"`
	CreatedAt                      time.Time `gorm:"autoCreateTime"`
	UpdatedAt                      time.Time `gorm:"autoUpdateTime"`
	CorrespondenceAddress1         string    `gorm:"column:correspondence_address_1;size:500"`
	CorrespondenceAddress2         string    `gorm:"column:correspondence_address_2;size:500"`
	CorrespondenceAddress3         string    `gorm:"column:correspondence_address_3;size:500"`
}

func (PropertyAddress) TableName() string {
	return "DIGIT3.property_addresses"
}

// AssessmentDetails represents assessment details for a property.
type AssessmentDetails struct {
	ID                         uuid.UUID   `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	ReasonOfCreation           string      `gorm:"size:200"`
	OccupancyCertificateNumber string      `gorm:"size:100"`
	OccupancyCertificateDate   *utils.Date `gorm:"type:date"`
	ExtentOfSite               string      `gorm:"size:100;not null"`
	IsLandUnderneathBuilding   string      `gorm:"size:100;not null"`
	IsUnspecifiedShare         bool        `gorm:"default:false;not null"`
	PropertyID                 uuid.UUID   `gorm:"type:uuid;not null;uniqueIndex"`
	CreatedAt                  time.Time   `gorm:"autoCreateTime"`
	UpdatedAt                  time.Time   `gorm:"autoUpdateTime"`
}

// Amenities represents an amenity for a property.
type Amenities struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Type        pq.StringArray `gorm:"type:text[]" json:"type" binding:"required,dive"`
	Description string         `gorm:"type:text"`
	ExpiryDate  *time.Time
	PropertyID  uuid.UUID `gorm:"type:uuid;not null;uniqueIndex"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}

// ConstructionDetails represents construction details for a property.
type ConstructionDetails struct {
	ID           uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	FloorType    string         `gorm:"size:100"`
	WallType     string         `gorm:"size:100"`
	RoofType     string         `gorm:"size:100"`
	WoodType     string         `gorm:"size:100"`
	PropertyID   uuid.UUID      `gorm:"type:uuid;not null;uniqueIndex"`
	FloorDetails []FloorDetails `gorm:"foreignKey:ConstructionDetailsID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	CreatedAt    time.Time      `gorm:"autoCreateTime"`
	UpdatedAt    time.Time      `gorm:"autoUpdateTime"`
}

// FloorDetails represents details for a specific floor.
type FloorDetails struct {
	ID                    uuid.UUID   `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	FloorNo               int         `gorm:"not null"`
	Classification        string      `gorm:"size:100"`
	NatureOfUsage         string      `gorm:"size:100"`
	FirmName              string      `gorm:"size:200"`
	OccupancyType         string      `gorm:"size:100"`
	OccupancyName         string      `gorm:"size:200"`
	ConstructionDate      *utils.Date `gorm:"type:date" json:"constructionDate"`
	EffectiveFromDate     *utils.Date `gorm:"type:date" json:"effectiveFromDate"`
	UnstructuredLand      string      `gorm:"size:100"`
	LengthFt              float64     `gorm:"type:decimal(10,2)"`
	BreadthFt             float64     `gorm:"type:decimal(10,2)"`
	PlinthAreaSqFt        float64     `gorm:"type:decimal(10,2)"`
	BuildingPermissionNo  string      `gorm:"size:100"`
	FloorDetailsEntered   bool        `gorm:"default:false;not null"`
	ConstructionDetailsID uuid.UUID   `gorm:"type:uuid;not null;index"`
	CreatedAt             time.Time   `gorm:"autoCreateTime"`
	UpdatedAt             time.Time   `gorm:"autoUpdateTime"`
}

// AdditionalPropertyDetails represents extra metadata for a property.
type AdditionalPropertyDetails struct {
	ID         uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	FieldName  string         `gorm:"size:100;not null"`
	FieldValue datatypes.JSON `gorm:"type:json" json:"fieldValue"`
	PropertyID uuid.UUID      `gorm:"type:uuid;not null;index"`
	CreatedAt  time.Time      `gorm:"autoCreateTime"`
	UpdatedAt  time.Time      `gorm:"autoUpdateTime"`
}

// GISData represents GIS information for a property.
type GISData struct {
	ID          uuid.UUID     `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Source      string        `gorm:"type:enum('GPS','MANUAL_ENTRY','IMPORT');not null"`
	Type        string        `gorm:"type:enum('POINT','LINE','POLYGON');not null"`
	EntityType  string        `gorm:"size:100"`
	PropertyID  uuid.UUID     `gorm:"type:uuid;not null;uniqueIndex"`
	Latitude    float64       `gorm:"type:decimal(10,8);not null"`
	Longitude   float64       `gorm:"type:decimal(11,8);not null"`
	Coordinates []Coordinates `gorm:"foreignKey:GISDataID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	CreatedAt   time.Time     `gorm:"autoCreateTime"`
	UpdatedAt   time.Time     `gorm:"autoUpdateTime"`
}

// Coordinates represents a latitude/longitude pair.
type Coordinates struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Latitude  float64   `gorm:"type:decimal(10,8);not null"`
	Longitude float64   `gorm:"type:decimal(11,8);not null"`
	GISDataID uuid.UUID `gorm:"type:uuid;not null;index"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

type Document struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	PropertyID   uuid.UUID `gorm:"type:uuid;not null;index"`
	DocumentType string    `gorm:"size:100;not null"`
	DocumentName string    `gorm:"size:255;not null"`
	FileStoreID  string    `gorm:"column:file_store_id"`
	UploadDate   time.Time `gorm:"type:timestamptz;default:now()"`
	Action       string    `gorm:"size:100;not null;default:PENDING" json:"action"`
	UploadedBy   string    `gorm:"size:200" json:"uploadedBy,omitempty"`
	Size         string    `gorm:"not null" json:"size,omitempty"`
}

// ensure GORM uses DIGIT3 schema explicitly (optional; your NamingStrategy may already prefix tables)
func (Document) TableName() string {
	return "DIGIT3.document"
}

type IGRS struct {
	ID                 uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	Habitation         string    `gorm:"size:200;not null" json:"habitation"`
	IGRSWard           string    `gorm:"size:100" json:"igrsWard,omitempty"`
	IGRSLocality       string    `gorm:"size:100" json:"igrsLocality,omitempty"`
	IGRSBlock          string    `gorm:"size:100" json:"igrsBlock,omitempty"`
	DoorNoFrom         string    `gorm:"size:50" json:"doorNoFrom,omitempty"`
	DoorNoTo           string    `gorm:"size:50" json:"doorNoTo,omitempty"`
	IGRSClassification string    `gorm:"size:100" json:"igrsClassification,omitempty"`
	BuiltUpAreaPct     *float64  `gorm:"type:decimal(7,2)" json:"builtUpAreaPct,omitempty"`
	FrontSetback       *float64  `gorm:"type:decimal(8,2)" json:"frontSetback,omitempty"`
	RearSetback        *float64  `gorm:"type:decimal(8,2)" json:"rearSetback,omitempty"`
	SideSetback        *float64  `gorm:"type:decimal(8,2)" json:"sideSetback,omitempty"`
	TotalPlinthArea    *float64  `gorm:"type:decimal(10,2)" json:"totalPlinthArea,omitempty"`
	CreatedAt          time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt          time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
	PropertyID         uuid.UUID `gorm:"type:uuid;index;unique"`
}

func (IGRS) TableName() string {
	return "DIGIT3.igrs"

}
