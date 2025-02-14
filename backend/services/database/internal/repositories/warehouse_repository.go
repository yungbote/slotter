package repositories

import (
  "errors"
  "gorm.io/gorm"
  "github.com/yungbote/slotter/backend/services/database/internal/models"
)

type WarehouseRepository interface {
  Create(warehouse *models.Warehouse) error
  GetByID(id uint) (*models.Warehouse, error)
  Update(warehouse *models.Warehouse) error
  Delete(warehouse *models.Warehouse) error
}

type warehouseRepository struct {
  db *gorm.DB
}

func NewWarehouseRepository(db *gorm.DB) WarehouseRepository {
  return &warehouseRepository{db: db}
}

func (r *warehouseRepository) Create(warehouse *models.Warehouse) error {
  err := r.db.Create(warehouse).Error
  if err != nil {
    return ParseDBError("WarehouseRepository.Create", err)
  }
  return nil
}

func (r *warehouseRepository) GetByID(id uint) (*models.Warehouse, error) {
  var w models.Warehouse
  err := r.db.First(&w, id).Error
  if errors.Is(err, gorm.ErrRecordNotFound) {
    return nil, ParseDBError("WarehouseRepository.GetByID", err)
  }
  if err != nil {
    return nil, ParseDBError("WarehouseRepository.GetByID", err)
  }
  return &w, nil
}

func (r *warehouseRepository) Update(warehouse *models.Warehouse) (*models.Warehouse, error) {
  err := r.db.Save(warehouse).Error
  if err != nil {
    return nil, ParseDBError("WarehouseRepository.Update", err)
  }
  return nil
}

func (r *warehouseRepository) Delete(warehouse *models.Warehouse) error {
  err := r.db.Delete(warehouse).Error
  if err != nil {
    return nil, ParseDBError("WarehouseRepository.Delete", err)
  }
  return nil
}

func (r *warehouseRepository) ListByCompanyID(companyID uint) ([]*models.Warehouse, error) {
  const op = "WarehouseRepository.ListByCompanyID"
  var warehouses []*models.Warehouse
  err := r.db.Where("company_id = ?", companyID).Find(&warehouses).Error
  if err != nil {
    return nil, ParseDBError(op, err)
  }
  return warehouses, nil
}

func (r *warehouseRepository) CountByCompanyID(companyID uint) (int64, error) {
  const op = "WarehouseRepository.CountByCompanyID"
  var count int64
  err := r.db.Model(&models.Warehouse{}).Where("company_id = ?", companyID).Count(&count).Error
  if err != nil {
    return 0, ParseDBError(op, err)
  }
  return count, nil
}


