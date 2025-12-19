package services

import (
	"context"
	"enumeration/internal/constants"
	"enumeration/internal/models"
	"enumeration/internal/repositories"
	"fmt"

	"github.com/google/uuid"
)

// Compile-time check for interface implementation
var _ PropertyAddressService = (*propertyAddressService)(nil)

// propertyAddressService handles business logic for property addresses
type propertyAddressService struct {
	propertyAddressRepo repositories.PropertyAddressRepository
}

// NewPropertyAddressService returns a new propertyAddressService
func NewPropertyAddressService(propertyAddressRepo repositories.PropertyAddressRepository) PropertyAddressService {
	return &propertyAddressService{
		propertyAddressRepo: propertyAddressRepo,
	}
}

// CreatePropertyAddress validates and adds a new property address
func (s *propertyAddressService) CreatePropertyAddress(ctx context.Context, address *models.PropertyAddress) error {
	if address == nil {
		return fmt.Errorf("%w: request is nil", ErrValidation)
	}

	if err := s.validatePropertyAddress(address); err != nil {
		return err
	}

	return s.propertyAddressRepo.Create(address)
}

// GetPropertyAddressByID fetches a property address by its ID
func (s *propertyAddressService) GetPropertyAddressByID(ctx context.Context, id uuid.UUID) (*models.PropertyAddress, error) {
	if id == uuid.Nil {
		return nil, fmt.Errorf("invalid property address ID: cannot be nil")
	}

	// Repository layer error - propagate as is
	return s.propertyAddressRepo.GetByID(id)
}

// UpdatePropertyAddress validates and updates an existing property address
func (s *propertyAddressService) UpdatePropertyAddress(ctx context.Context, address *models.PropertyAddress) error {
	if address == nil {
		return fmt.Errorf("%w: request is nil", ErrValidation)
	}

	if address.ID == uuid.Nil {
		return fmt.Errorf("invalid property address ID: cannot be nil")
	}
	return s.propertyAddressRepo.Update(address)
}

// DeletePropertyAddress removes a property address by its ID
func (s *propertyAddressService) DeletePropertyAddress(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return fmt.Errorf("invalid property address ID: cannot be nil")
	}

	// Repository layer error - propagate as is
	return s.propertyAddressRepo.Delete(id)
}

// GetAllPropertyAddresses returns all property addresses with pagination, optionally filtered by propertyID
func (s *propertyAddressService) GetAllPropertyAddresses(ctx context.Context, page, size int, propertyID *uuid.UUID) ([]*models.PropertyAddress, int64, error) {
	if page < 0 {
		page = constants.DefaultPage
	}
	if size <= 0 || size > 100 {
		size = constants.DefaultSize
	}

	return s.propertyAddressRepo.GetAll(page, size, propertyID)
}

// GetPropertyAddressByPropertyID fetches a property address for a property
func (s *propertyAddressService) GetPropertyAddressByPropertyID(ctx context.Context, propertyID uuid.UUID) (*models.PropertyAddress, error) {
	if propertyID == uuid.Nil {
		return nil, fmt.Errorf("invalid property ID: cannot be nil")
	}

	// Repository layer error - propagate as is
	return s.propertyAddressRepo.GetByPropertyID(propertyID)
}

// SearchPropertyAddresses finds property addresses matching the given search parameters
func (s *propertyAddressService) SearchPropertyAddresses(ctx context.Context, params SearchPropertyAddressParams) ([]*models.PropertyAddress, int64, error) {
	if params.Page < 0 {
		params.Page = constants.DefaultPage
	}
	if params.Size <= 0 || params.Size > 100 {
		params.Size = constants.DefaultSize
	}

	repoParams := repositories.SearchPropertyAddressParams{
		Page:            params.Page,
		Size:            params.Size,
		PropertyID:      params.PropertyID,
		Locality:        params.Locality,
		ZoneNo:          params.ZoneNo,
		WardNo:          params.WardNo,
		BlockNo:         params.BlockNo,
		Street:          params.Street,
		ElectionWard:    params.ElectionWard,
		SecretariatWard: params.SecretariatWard,
		PinCode:         params.PinCode,
		SortBy:          params.SortBy,
		SortOrder:       params.SortOrder,
	}

	return s.propertyAddressRepo.Search(repoParams)
}

// validatePropertyAddress checks required fields for a property address and returns an error if validation fails
func (s *propertyAddressService) validatePropertyAddress(address *models.PropertyAddress) error {
	// Validate required fields based on your model
	if address.PropertyID == uuid.Nil {
		return fmt.Errorf("property ID is required")
	}

	if address.Locality == "" {
		return fmt.Errorf("locality is required")
	}

	if address.ZoneNo == "" {
		return fmt.Errorf("zone number is required")
	}

	if address.WardNo == "" {
		return fmt.Errorf("ward number is required")
	}

	if address.BlockNo == "" {
		return fmt.Errorf("block number is required")
	}

	if address.ElectionWard == "" {
		return fmt.Errorf("election ward is required")
	}

	// Validate PIN code (6 digits)
	if address.PinCode != 0 && (address.PinCode < 100000 || address.PinCode > 999999) {
		return fmt.Errorf("PIN code must be 6 digits")
	}

	return nil
}
