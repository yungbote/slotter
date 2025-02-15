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
  err := r.db.Create(company).Error
  if err != nil {
    return ParseDBError("CompanyRepository.Create", err)
  }
  return nil
}

func (r *gormCompanyRepository) GetByID(id uint) (*models.Company, error) {
  var c models.Company
  err := r.db.First(&c, id).Error
  if errors.Is(err, gorm.ErrRecordNotFound) {
    return nil, ParseDBError("CompanyRepository.GetByID", err)
  }
  if err != nil {
    return nil, ParseDBError("CompanyRepository.GetByID", err)
  }
  return &c, nil
}

func (r *gormCompanyRepository) Update(company *models.Company) error {
  err := r.db.Save(company).Error
  if err != nil {
    return nil, ParseDBError("CompanyRepository.Update", err)
  }
  return nil
}

func (r *gormCompanyRepository) Delete(company *models.Company) error {
  err := r.db.Delete(company).Error
  if err != nil {
    return nil, ParseDBError("CompanyRepository.Delete", err)
  }
  return nil
}



