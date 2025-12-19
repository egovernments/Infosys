package repositories

import (
    "context"
    "enumeration/internal/models"

    "github.com/google/uuid"
    "gorm.io/gorm"
)


var _ ApplicationLogRepository = (*applicationLogRepository)(nil)
type applicationLogRepository struct {
    db *gorm.DB
}

func NewApplicationLogRepository(db *gorm.DB) ApplicationLogRepository {
    return &applicationLogRepository{
        db: db,
    }
}

func (r *applicationLogRepository) Create(ctx context.Context, log *models.ApplicationLog) error {
    return r.db.WithContext(ctx).Create(log).Error
}

func (r *applicationLogRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.ApplicationLog, error) {
    var log models.ApplicationLog
    err := r.db.WithContext(ctx).Where("id = ?", id).First(&log).Error
    if err != nil {
        return nil, err
    }
    return &log, nil
}

func (r *applicationLogRepository) Update(ctx context.Context, log *models.ApplicationLog) error {
    return r.db.WithContext(ctx).Save(log).Error
}

func (r *applicationLogRepository) Delete(ctx context.Context, id uuid.UUID) error {
    return r.db.WithContext(ctx).Delete(&models.ApplicationLog{}, id).Error
}

func (r *applicationLogRepository) List(ctx context.Context, applicationID, action string, page, size int) ([]models.ApplicationLog, int64, error) {
    var logs []models.ApplicationLog
    var total int64

    query := r.db.WithContext(ctx).Model(&models.ApplicationLog{})

    // Apply filters
    if applicationID != "" {
        if appID, err := uuid.Parse(applicationID); err == nil {
            query = query.Where("application_id = ?", appID)
        }
    }
    if action != "" {
        query = query.Where("action = ?", action)
    }

    // Get total count
    err := query.Count(&total).Error
    if err != nil {
        return nil, 0, err
    }

    // Apply pagination and get results
    offset := page * size
    err = query.Offset(offset).Limit(size).Order("performed_date DESC").Find(&logs).Error
    if err != nil {
        return nil, 0, err
    }

    return logs, total, nil
}

func (r *applicationLogRepository) GetByApplicationID(ctx context.Context, applicationID uuid.UUID, action string, page, size int) ([]models.ApplicationLog, int64, error) {
    var logs []models.ApplicationLog
    var total int64

    query := r.db.WithContext(ctx).Model(&models.ApplicationLog{}).Where("application_id = ?", applicationID)

    // Apply action filter if provided
    if action != "" {
        query = query.Where("action = ?", action)
    }

    // Get total count
    err := query.Count(&total).Error
    if err != nil {
        return nil, 0, err
    }

    // Apply pagination and get results
    offset := page * size
    err = query.Offset(offset).Limit(size).Order("performed_date DESC").Find(&logs).Error
    if err != nil {
        return nil, 0, err
    }

    return logs, total, nil
}