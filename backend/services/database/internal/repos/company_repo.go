package repos

import (
  "fmt"

  "gorm.io/gorm"
  "github.com/google/uuid"
  "github.com/yungbote/slotter/backend/services/database/internal/models"
  customErr "github.com/yungbote/slotter/backend/services/database/internal/errors"
)

type CRepo interface {
  Create(company models.Company) (*models.Company, error)
  UpdateName(companyID uuid.UUID, newName string) error
  UpdateAvatarURL(companyID uuid.UUID, newAvatarURL string) error
  GetByID(companyID uuid.UUID) (*models.Company, error)
  Delete(companyID uuid.UUID) error
}

type cRepo struct {
  db *gorm.DB
}

func NewCRepo(db *gorm.DB) CRepo {
  return &cRepo{db: db}
}

func (r *cRepo) Create(company models.Company) (*models.Company, error) {
  if err := r.db.Create(&company).Error; err != nil {
    return nil, fmt.Errorf("Failed to create company: %w", err)
  }
  return &company, nil
}

func (r *cRepo) UpdateName(companyID uuid.UUID, newName string) error {
  if err := r.db.Model(&models.Company{}).
    Where("id = ?", companyID).
    Update("name", newName).Error; err != nil {
    return fmt.Errorf("Failed to update company name: %w", err)
  }
  return nil
}

func (r *cRepo) UpdateAvatarURL(companyID uuid.UUID, newAvatarURL string) error {
  if err := r.db.Model(&models.Company{}).
    Where("id = ?", companyID).
    Update("avatar_url", newAvatarURL).Error; err != nil {
    return fmt.Errorf("Failed to update company avatar url: %w", err)
  }
  return nil
}

func (r *cRepo) GetByID(companyID uuid.UUID) (*models.Company, error) {
  var c models.Company
  if err := r.db.First(&c, "id = ?", companyID).Error; err != nil {
    return nil, fmt.Errorf("Company not found: %w", err)
  }
  return &c, nil
}

func (r *cRepo) Delete(companyID uuid.UUID) error {
  company, err := r.GetByID(companyID)
  if err != nil {
    return err
  }
  if err := r.db.Delete(company).Error; err != nil {
    return fmt.Errorf("Failed to delete company: %w", err)
  }
  return nil
}





