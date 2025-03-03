package repos

import (
  "fmt"

  "gorm.io/gorm"
  "github.com/google/uuid"
  "github.com/yungbote/slotter/backend/services/database/internal/models"
)

type UARepo interface {
    Create(action models.UserAction) (*models.UserAction, error)
    UpdateActionType(actionID uuid.UUID, newAType string) error
    UpdateEntityType(actionID uuid.UUID, newEType string) error
    UpdateDescription(actionID uuid.UUID, newDescription string) error
    GetByID(actionID uuid.UUID) (*models.UserAction, error)
}

type uaRepo struct {
    db *gorm.DB
}

func NewUARepo(db *gorm.DB) UARepo {
    return &uaRepo{db: db}
}

func (r *uaRepo) Create(action models.UserAction) (*models.UserAction, error) {
    if err := r.db.Create(&action).Error; err != nil {
        return nil, fmt.Errorf("failed to create user action: %w", err)
    }
    return &action, nil
}

func (r *uaRepo) UpdateActionType(actionID uuid.UUID, newAType string) error {
    return r.db.Model(&models.UserAction{}).
        Where("id = ?", actionID).
        Update("action_type", newAType).Error
}

func (r *uaRepo) UpdateEntityType(actionID uuid.UUID, newEType string) error {
    return r.db.Model(&models.UserAction{}).
        Where("id = ?", actionID).
        Update("entity_type", newEType).Error
}

func (r *uaRepo) UpdateDescription(actionID uuid.UUID, newDescription string) error {
    return r.db.Model(&models.UserAction{}).
        Where("id = ?", actionID).
        Update("description", newDescription).Error
}

func (r *uaRepo) GetByID(actionID uuid.UUID) (*models.UserAction, error) {
    var ua models.UserAction
    if err := r.db.First(&ua, "id = ?", actionID).Error; err != nil {
        return nil, fmt.Errorf("user action not found: %w", err)
    }
    return &ua, nil
}


