package services

import (
	"context"
	"enumeration/internal/constants"
	"enumeration/internal/dto"
	"enumeration/internal/models"
	"enumeration/internal/repositories"
	"fmt"

	"github.com/google/uuid"
)

// Compile-time check for interface implementation
var _ PropertyOwnerService = (*propertyOwnerService)(nil)

// propertyOwnerService handles business logic for property owners
type propertyOwnerService struct {
	repo repositories.PropertyOwnerRepository
}

// NewPropertyOwnerService returns a new propertyOwnerService
func NewPropertyOwnerService(repo repositories.PropertyOwnerRepository) PropertyOwnerService {
	return &propertyOwnerService{
		repo: repo,
	}
}

// Create adds a new property owner after validation
func (s *propertyOwnerService) Create(ctx context.Context, req *dto.CreatePropertyOwnerRequest) (*models.PropertyOwner, error) {
	if req == nil {
		return nil, fmt.Errorf("%w: request is nil", ErrValidation)
	}
	owner := &models.PropertyOwner{
		ID:                     uuid.New(),
		PropertyID:             req.PropertyID,
		AdhaarNo:               req.AdhaarNo,
		Name:                   req.Name,
		ContactNo:              req.ContactNo,
		Email:                  req.Email,
		Gender:                 req.Gender,
		Guardian:               req.Guardian,
		GuardianType:           req.GuardianType,
		RelationshipToProperty: req.RelationshipToProperty,
		OwnershipShare:         req.OwnershipShare,
		IsPrimaryOwner:         req.IsPrimaryOwner,
	}

	if err := s.repo.Create(ctx, owner); err != nil {
		return nil, err
	}
	return owner, nil
}

// CreatePropertyOwners adds multiple property owners in a batch after validation
func (s *propertyOwnerService) CreatePropertyOwners(ctx context.Context, reqs []*dto.CreatePropertyOwnerRequest) ([]*models.PropertyOwner, error) {
	if len(reqs) == 0 {
		return []*models.PropertyOwner{}, nil
	}

	owners := make([]*models.PropertyOwner, 0, len(reqs))
	for i, r := range reqs {

		if r == nil {
			return nil, fmt.Errorf("%w: owner at index %d is nil", ErrValidation, i)
		}

		// Basic validation
		if r.PropertyID == uuid.Nil {
			return nil, fmt.Errorf("%w: propertyId is required for owner at index %d", ErrValidation, i)
		}
		if r.Name == "" {
			return nil, fmt.Errorf("%w: name is required for owner at index %d", ErrValidation, i)
		}
		if r.ContactNo == "" {
			return nil, fmt.Errorf("%w: contactNo is required for owner at index %d", ErrValidation, i)
		}
		if r.Gender == "" {
			return nil, fmt.Errorf("%w: gender is required for owner at index %d", ErrValidation, i)
		}

		owner := &models.PropertyOwner{
			ID:                     uuid.New(),
			PropertyID:             r.PropertyID,
			AdhaarNo:               r.AdhaarNo,
			Name:                   r.Name,
			ContactNo:              r.ContactNo,
			Email:                  r.Email,
			Gender:                 r.Gender,
			Guardian:               r.Guardian,
			GuardianType:           r.GuardianType,
			RelationshipToProperty: r.RelationshipToProperty,
			OwnershipShare:         r.OwnershipShare,
			IsPrimaryOwner:         r.IsPrimaryOwner,
		}
		owners = append(owners, owner)
	}

	if err := s.repo.CreateBatch(ctx, owners); err != nil {
		return nil, err
	}
	return owners, nil
}

// Update modifies an existing property owner by ID
func (s *propertyOwnerService) Update(ctx context.Context, id uuid.UUID, req *dto.UpdatePropertyOwnerRequest) (*models.PropertyOwner, error) {
	if req == nil {
		return nil, fmt.Errorf("%w: request is nil", ErrValidation)
	}

	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Name != "" {
		existing.Name = req.Name
	}
	if req.ContactNo != "" {
		existing.ContactNo = req.ContactNo
	}
	if req.Email != "" {
		existing.Email = req.Email
	}
	if req.Gender != "" {
		existing.Gender = req.Gender
	}
	if req.Guardian != "" {
		existing.Guardian = req.Guardian
	}
	if req.GuardianType != "" {
		existing.GuardianType = req.GuardianType
	}
	if req.RelationshipToProperty != "" {
		existing.RelationshipToProperty = req.RelationshipToProperty
	}
	if req.OwnershipShare > 0 {
		existing.OwnershipShare = req.OwnershipShare
	}
	existing.IsPrimaryOwner = req.IsPrimaryOwner

	if err := s.repo.Update(ctx, existing); err != nil {
		return nil, err
	}
	return existing, nil
}

// Delete removes a property owner by ID
func (s *propertyOwnerService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

// GetByPropertyID fetches property owners for a property with pagination
func (s *propertyOwnerService) GetByPropertyID(ctx context.Context, propertyID uuid.UUID, page, size int) ([]*models.PropertyOwner, int64, error) {
	// clamp
	if page < 0 {
		page = constants.DefaultPage
	}
	if size <= 0 || size > 100 {
		size = constants.DefaultSize
	}

	owners, total, err := s.repo.GetByPropertyID(ctx, propertyID, page, size)
	if err != nil {
		return nil, 0, err
	}
	if owners == nil {
		owners = []*models.PropertyOwner{}
	}
	return owners, total, nil
}
