package services

import (
  "github.com/yungbote/slotter/backend/services/database/internal/models"
  "github.com/yungbote/slotter/backend/services/database/internal/repositories"
)

type WarehouseService interface {
  CreateWarehouse(input *models.Warehouse) (*models.Warehouse, error)
  GetWarehouseByID(id uint) (*models.Warehouse, error)
  UpdateWarehouse(existing *models.Warehouse, updates *models.Warehouse) (*models.Warehouse, error)
  DeleteWarehouse(existing *models.Warehouse) error
}

type warehouseService struct {
  repo repositories.WarehouseRepository
}

func NewWarehouseService(repo repositories.WarehouseRepository) WarehouseService {
  return &warehouseService{repo: repo}
}

func (s *warehouseService) CreateWarehouse(input *models.Warehouse) (*models.Warehouse, error) {
  if err := s.repo.Create(input); err != nil {
    return nil, err
  }
  return input, nil
}

func (s *warehouseService) GetWarehouseByID(id uint) (*models.Warehouse, error) {
  return s.repo.GetByID(id)
}

func (s *warehouseService) UpdateWarehouse(existing *models.Warehouse, updates *models.Warehouse) (*models.Warehouse, error) {
  existing.Name = updates.Name
  existing.CompanyID = updates.CompanyID
  if err := s.repo.Update(existing); err != nil {
    return nil, err
  }
  return existing, nil
}

func (s *warehouseService) DeleteWarehouse(existing *models.Warehouse) error {
  return s.repo.Delete(existing)
}
