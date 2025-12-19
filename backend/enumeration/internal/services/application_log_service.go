package services

import (
	"context"
	"encoding/json"
	"enumeration/internal/constants"
	"enumeration/internal/dto"
	"enumeration/internal/models"
	"enumeration/internal/repositories"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var _ ApplicationLogService = (*applicationLogService)(nil)

type applicationLogService struct {
    repo repositories.ApplicationLogRepository
}

func NewApplicationLogService(repo repositories.ApplicationLogRepository) ApplicationLogService {
    return &applicationLogService{
        repo: repo,
    }
}

func (s *applicationLogService) Create(ctx context.Context, req *dto.CreateApplicationLogRequest) (*models.ApplicationLog, error) {
    // Convert metadata to JSON string if provided
     if req == nil {
        return nil, fmt.Errorf("%w: request is nil", ErrValidation)
    }
    var metadataJSON string
    if req.Metadata != nil {
        metadataBytes, err := json.Marshal(req.Metadata)
        if err != nil {
            return nil, fmt.Errorf("failed to marshal metadata: %w", err)
        }
        metadataJSON = string(metadataBytes)
    }

    log := &models.ApplicationLog{
        Action:        req.Action,
        PerformedBy:   req.PerformedBy,
        PerformedDate: time.Now(),
        Comments:      req.Comments,
        Metadata:      metadataJSON,
        FileStoreID:   req.FileStoreID,
        ApplicationID: req.ApplicationID,
        Actor: req.Actor,
    }

    err := s.repo.Create(ctx, log)
    if err != nil {
        return nil, fmt.Errorf("failed to create application log: %w", err)
    }

    return log, nil
}

func (s *applicationLogService) GetByID(ctx context.Context, id uuid.UUID) (*models.ApplicationLog, error) {
    log, err := s.repo.GetByID(ctx, id)
    if err != nil {
        if err == gorm.ErrRecordNotFound {
            return nil, fmt.Errorf("record not found")
        }
        return nil, fmt.Errorf("failed to get application log: %w", err)
    }
    return log, nil
}

func (s *applicationLogService) Update(ctx context.Context, id uuid.UUID, req *dto.UpdateApplicationLogRequest) (*models.ApplicationLog, error) {
    // Check if log exists
    if req == nil {
        return nil, fmt.Errorf("%w: request is nil", ErrValidation)
    }
    existingLog, err := s.repo.GetByID(ctx, id)
    if err != nil {
        if err == gorm.ErrRecordNotFound {
            return nil, fmt.Errorf("record not found")
        }
        return nil, fmt.Errorf("failed to get application log: %w", err)
    }

    // Update fields if provided
    if req.Action != "" {
        existingLog.Action = req.Action
    }
    if req.PerformedBy != "" {
        existingLog.PerformedBy = req.PerformedBy
    }
    if req.Comments != "" {
        existingLog.Comments = req.Comments
    }
    if req.FileStoreID != nil {
        existingLog.FileStoreID = req.FileStoreID
    }
    if req.Metadata != nil {
        metadataBytes, err := json.Marshal(req.Metadata)
        if err != nil {
            return nil, fmt.Errorf("failed to marshal metadata: %w", err)
        }
        existingLog.Metadata = string(metadataBytes)
    }

    err = s.repo.Update(ctx, existingLog)
    if err != nil {
        return nil, fmt.Errorf("failed to update application log: %w", err)
    }

    return existingLog, nil
}

func (s *applicationLogService) Delete(ctx context.Context, id uuid.UUID) error {
    // Check if log exists
    _, err := s.repo.GetByID(ctx, id)
    if err != nil {
        if err == gorm.ErrRecordNotFound {
            return fmt.Errorf("record not found")
        }
        return fmt.Errorf("failed to get application log: %w", err)
    }

    err = s.repo.Delete(ctx, id)
    if err != nil {
        return fmt.Errorf("failed to delete application log: %w", err)
    }

    return nil
}

func (s *applicationLogService) List(ctx context.Context, applicationID, action string, page, size int) ([]models.ApplicationLog, int64, error) {
    if page < 0 {
        page = constants.DefaultPage
    }
    if size <= 0 {
        size = constants.DefaultSize
    }

    logs, total, err := s.repo.List(ctx, applicationID, action, page, size)
    if err != nil {
        return nil, 0, fmt.Errorf("failed to list application logs: %w", err)
    }

    return logs, total, nil
}

func (s *applicationLogService) GetByApplicationID(ctx context.Context, applicationID uuid.UUID, action string, page, size int) ([]models.ApplicationLog, int64, error) {
    if page < 0 {
        page = constants.DefaultPage
    }
    if size <= 0 {
        size = constants.DefaultSize
    }

    logs, total, err := s.repo.GetByApplicationID(ctx, applicationID, action, page, size)
    if err != nil {
        return nil, 0, fmt.Errorf("failed to get application logs by application ID: %w", err)
    }

    return logs, total, nil
}