package repositories

import (
  "gorm.io/gorm"
  "github.com/yungbote/slotter/backend/services/database/internal/models"
)

type LocationRepository interface {
  Create(location *models.Location) error
  GetByID(id uint) (*models.Location, error)
  Update(location *models.Location) error
  GetParentLocations()
  GetType()
  GetName()
}

type locationRepository struct {
  db *gorm.DB
}

func NewLocationRepository(db *gorm.DB) LocationRepository {
  return &locationRepository{db: db}
}

func (r *locationRepository) Create(location *models.Location) error {
  err := r.db.Create(location).Error
  if err != nil {
    return ParseDBError("LocationRepository.Create", err)
  }
  return nil
}

func (r *locationRepository) GetByID(id uint) (*models.Location, error) {
  var l models.Location
  err := r.db.First(&l, id).Error
  if errors.Is(err, gorm.ErrRecordNotFound) {
    return nil, ParseDBError("LocationRepository.GetByID", err)
  }
  if err != nil {
    return nil, ParseDBError("LocationRepository.GetByID", err)
  }
  return &l, nil
}

func (r *locationRepository) Update(location *models.Location) error {
  err := r.db.Save(location).Error
  if err != nil {
    return ParseDBError("LocationRepository.Update", err)
  }
  return nil
}

func (r *locationRepository) Delete(location *models.Location) error {
  err := r.db.Delete(location).Error
  if err != nil {
    return ParseDBError("LocationRepository.Delete", err)
  }
  return nil
}

func (r *locationRepository) ListByWarehouseID(warehouseID uint) ([]*models.Location, error) {
  const op = "LocationRepository.ListByWarehouseID"
  var locs []*models.Location
  err := r.db.Where("warehouse_id = ?", warehouseID).Find(&locs).Error
  if err != nil {
    return nil, ParseDBError(op, err)
  }
  return locs, nil
}

func (r *locationRepository) ListChildrenByParentID(parentLocationID uint) ([]*models.Location, error) {
  const op = "LocationRepository.ListChildrenByParentID"
  var children []*models.Location
  err := r.db.Where("parent_location_id = ?", parentLocationID).Find(&children).Error
  if err != nil {
    return nil, ParseDBError(op, err)
  }
  return children, nil
}

func (r *locationRepository) CountByWarehouseID(warehuseID uint) (int64, error) {
  const op = "LocationRepository.CountByWarehouseID"
  var count int64
  err := r.db.Model(&models.Location{}).Where("warehouse_id = ?", warehouseID).Count(&count).Error
  if err != nil {
    return 0, ParseDBError(op, err)
  }
  return count, nil
}


