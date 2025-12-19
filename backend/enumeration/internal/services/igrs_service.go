// Package services provides business logic and service layer implementations
// for the property tax enumeration system.
package services

import (
	"context"
	"errors"
	"fmt"

	"enumeration/internal/constants"
	"enumeration/internal/dto"
	"enumeration/internal/models"
	"enumeration/internal/repositories"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Compile-time check to ensure igrsService implements IGRSService interface.
var _ IGRSService = (*igrsService)(nil)

// igrsService implements the IGRSService interface and provides methods
// for managing IGRS (Integrated Grievance Redressal System) records.
type igrsService struct {
	repo repositories.IGRSRepository // Repository for IGRS data access
}

// NewIGRSService creates a new instance of IGRSService with the provided repository.
func NewIGRSService(repo repositories.IGRSRepository) IGRSService {
	return &igrsService{repo: repo}
}

// Create validates and creates a new IGRS record.
// Ensures only one IGRS exists per property and required fields are present.
func (s *igrsService) Create(ctx context.Context, req *dto.CreateIGRSRequest) (*models.IGRS, error) {
	if req == nil {
		return nil, fmt.Errorf("%w: request is nil", ErrValidation)
	}

	if req.PropertyID == uuid.Nil {
		return nil, errors.New("propertyId is required")
	}
	if req.Habitation == "" {
		return nil, errors.New("habitation is required")
	}

	// Ensure only one IGRS per property
	if existing, err := s.repo.GetByPropertyID(ctx, req.PropertyID); err == nil && existing != nil {
		return nil, fmt.Errorf("igrs already exists for property %s", req.PropertyID.String())
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		// Propagate unexpected db error
		return nil, fmt.Errorf("failed to check existing igrs: %w", err)
	}

	// Map request fields to IGRS model
	m := &models.IGRS{
		PropertyID:         req.PropertyID,
		Habitation:         req.Habitation,
		IGRSWard:           req.IGRSWard,
		IGRSLocality:       req.IGRSLocality,
		IGRSBlock:          req.IGRSBlock,
		DoorNoFrom:         req.DoorNoFrom,
		DoorNoTo:           req.DoorNoTo,
		IGRSClassification: req.IGRSClassification,
		BuiltUpAreaPct:     req.BuiltUpAreaPct,
		FrontSetback:       req.FrontSetback,
		RearSetback:        req.RearSetback,
		SideSetback:        req.SideSetback,
		TotalPlinthArea:    req.TotalPlinthArea,
	}
	if err := s.repo.Create(ctx, m); err != nil {
		return nil, fmt.Errorf("failed to create igrs: %w", err)
	}
	return m, nil
}

// GetByID retrieves an IGRS record by its unique ID.
func (s *igrsService) GetByID(ctx context.Context, id uuid.UUID) (*models.IGRS, error) {
	return s.repo.GetByID(ctx, id)
}

// Update modifies an existing IGRS record by its ID.
// Only fields provided in the request are updated.
func (s *igrsService) Update(ctx context.Context, id uuid.UUID, req *dto.UpdateIGRSRequest) (*models.IGRS, error) {
	if req == nil {
		return nil, fmt.Errorf("%w: request is nil", ErrValidation)
	}

	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update only provided fields
	if req.Habitation != "" {
		existing.Habitation = req.Habitation
	}
	if req.IGRSWard != "" {
		existing.IGRSWard = req.IGRSWard
	}
	if req.IGRSLocality != "" {
		existing.IGRSLocality = req.IGRSLocality
	}
	if req.IGRSBlock != "" {
		existing.IGRSBlock = req.IGRSBlock
	}
	if req.DoorNoFrom != "" {
		existing.DoorNoFrom = req.DoorNoFrom
	}
	if req.DoorNoTo != "" {
		existing.DoorNoTo = req.DoorNoTo
	}
	if req.IGRSClassification != "" {
		existing.IGRSClassification = req.IGRSClassification
	}
	// Numeric/pointer fields: assign if not nil
	if req.BuiltUpAreaPct != nil {
		existing.BuiltUpAreaPct = req.BuiltUpAreaPct
	}
	if req.FrontSetback != nil {
		existing.FrontSetback = req.FrontSetback
	}
	if req.RearSetback != nil {
		existing.RearSetback = req.RearSetback
	}
	if req.SideSetback != nil {
		existing.SideSetback = req.SideSetback
	}
	if req.TotalPlinthArea != nil {
		existing.TotalPlinthArea = req.TotalPlinthArea
	}

	if err := s.repo.Update(ctx, existing); err != nil {
		return nil, fmt.Errorf("failed to update igrs: %w", err)
	}
	return existing, nil
}

// Delete removes an IGRS record by its ID.
func (s *igrsService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

// List returns a paginated list of IGRS records and the total count.
// If page or size are invalid, defaults are used.
func (s *igrsService) List(ctx context.Context, page, size int) ([]*models.IGRS, int64, error) {
	if page < 0 {
		page = constants.DefaultPage
	}
	if size <= 0 {
		size = constants.DefaultSize
	}
	items, total, err := s.repo.FindAll(ctx, page, size)
	if err != nil {
		return nil, 0, err
	}
	out := make([]*models.IGRS, len(items))
	for i := range items {
		// Take address of slice element is OK here because items is not reused.
		out[i] = &items[i]
	}
	return out, total, nil
}
