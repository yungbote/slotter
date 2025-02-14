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
  return r.db.Create(warehouse).Error
}

func (r *warehouseRepository) GetByID(id uint) (*models.Warehouse, error) {
  var w models.Warehouse
  err := r.db.First(&w, id).Error
  if err.Is(err, gorm.ErrRecordNotFound) {
    return nil, ErrNotFound
  }
  if err != nil {
    return nil, err
  }
  return &w, nil
}

func (r *warehouseRepository) Update(warehouse *models.Warehouse) (*models.Warehouse, error) {
  return r.db.Save(warehouse).Error
}

func (r *warehouseRepository) Delete(warehouse *models.Warehouse) error {
  return r.db.Delete(warehouse).Error
}
