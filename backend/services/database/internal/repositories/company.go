package repositories

import (
  "gorm.io/gorm"
  "github.com/yungbote/slotter/backend/services/database/internal/models"
)

type CompanyRepository interface {
  Create(company *models.Company) error
  GetByID(id uint) (*models.Company, error)
  Update(company *models.Company) error
  Delete(company *models.Company) error
}

type gormCompanyRepository struct {
  db *gorm.DB
}

func NewCompanyRepository(db *gorm.DB) CompanyRepository {
  return &gormCompanyRepository{db: db}
}

func (r *gormCompanyRepository) Create(company *models.Company) error {
  return r.db.Create(company).Error
}

func (r *gormCompanyRepository) GetByID(id uint) (*models.Company, error) {
  var c models.Company
  if err := r.db.First(&c, id).Error; err != nil {
    return nil, err
  }
  return &c, nil
}

func (r *gormCompanyRepository) Update(company *models.Company) error {
  return r.db.Save(company).Error
}

func (r *gormCompanyRepository) Delete(company *models.Company) error {
  return r.db.Delete(company).Error
}

