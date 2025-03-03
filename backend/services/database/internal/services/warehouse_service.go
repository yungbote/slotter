package services

import (
  "fmt"
  "github.com/google/uuid"
  "github.com/yungbote/slotter/backend/services/database/internal/models"
  "github.com/yungbote/slotter/backend/services/database/internal/repos"
)

type WSvc interface {
  //GENERAL CRUD
  CreateWarehouse(warehouse models.Warehouse) (*models.Warehouse, error)
  UpdateWarehouseName(warehouseID uuid.UUID, newName string) error
  GetWarehouseByID(warehouseID uuid.UUID) (*models.Warehouse, error)
  DeleteWarehouse(warehouseID uuid.UUID) error

  LinkToItem(warehouseID, itemID uuid.UUID) error
  UnlinkFromItem(warehouseID, itemID uuid.UUID) error

  ListWarehouses(f repos.WarehouseFilter) ([]*models.Warehouse, error)

}

type wSvc struct {
  repo            repos.WRepo
}

func NewWSvc(repo repos.WRepo) WSvc {
  return &wSvc{repo: repo}
}

func (s *wSvc) CreateWarehouse(warehouse models.Warehouse) (*models.Warehouse, error) {
  if warehouse.Name == "" {
    return nil, fmt.Errorf("warehouse name is empty")
  }
  if warehouse.CompanyID == nil || *w.CompanyID == uuid.Nil {
    return nil, fmt.Errorf("warehouse must have a valid company ID")
  }
  created, err := s.repo.Create(warehouse)
  if err != nil {
    return nil, fmt.Errorf("failed to create warehouse: %w", err)
  }
  return created, nil
}

func (s *wSvc) UpdateWarehouseName(warehouseID uuid.UUID, newName string) error {
  if warehouseID == uuid.Nil || newName == "" {
    return fmt.Errorf("Invalid input to update warehouse name")
  }
  return s.repo.UpdateName(warehouseID, newName)
}

func (s *wSvc) UpdateWarehouseAvatarURL(warehouseID uuid.UUID, newAvatarURL string) error {
  if warehouseID == uuid.Nil || newAvatarURL == "" {
    return fmt.Errorf("Invalid input to update warehouse avatar url")
  }
  return s.repo.UpdateAvatarURL(warehouseID, newAvatarURL)
}

func (s *wSvc) GetWarehouseByID(warehouseID uuid.UUID) (*models.Warehouse, error) {
  if warehouseID == uuid.Nil {
    return fmt.Errorf("Invalid warehouseID")
  }
  warehouse, err := s.repo.GetByID(warehouseID)
  if err != nil {
    return fmt.Errorf("failed to get warehouse: %w", err)
  }
  return warehouse, nil
}

func (s *wSvc) DeleteWarehouse(warehouseID uuid.UUID) error {
  if warehouseID == uuid.Nil {
    return fmt.Errorf("Invalid warehouseID")
  }
  return s.repo.Delete(warehouseID)
}

func (s *wSvc) LinkToItem(warehouseID, itemID uuid.UUID) error {
  if warehouseID == uuid.Nil {
    return fmt.Errorf("Invalid warehouseID")
  }
  if itemID == uuid.Nil {
    return fmt.Errorf("Invalid itemID")
  }
  return s.repo.LinkToItem(warehouseID, itemID)
}

func (s *wSvc) UnlinkFromItem(warehouseID, itemID uuid.UUID) error {
  if warehouseID == uuid.Nil {
    return fmt.Errorf("Invalid warehouseID")
  }
  if itemID == uuid.Nil {
    return fmt.Errorf("Invalid itemID")
  }
  return s.repo.UnlinkFromItem
}

func (s *wSvc) ListWarehouses(f repos.WarehouseFilter) ([]*models.Warehouse, error) {
  warehouses, err := s.repo.ListWarehouses(f)
  if err != nil {
    return nil, fmt.Errorf("failed to list warehouses: %w", err)
  }
  return warehouses, nil
}
