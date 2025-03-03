package services

import (
  "fmt"
  "github.com/google/uuid"
  "github.com/yungbote/slotter/backend/services/database/internal/models"
  "github.com/yungbote/slotter/backend/services/database/internal/repos"
)

type ISvc interface {
  //GENERAL CRUD
  CreateItem(item models.Item) (*models.Item, error)
  UpdateItemName(itemID uuid.UUID, newName string) error
  GetItemByID(itemID uuid.UUID) (*models.Item, error)
  GetByItemNameAndCompanyID(companyID uuid.UUID, name string) (*models.Item, error)
  DeleteItem(itemID uuid.UUID) error
  
  LinkToLocation(itemID, locationID uuid.UUID) error
  UnlinkFromLocation(itemID, locationID uuid.UUID) error
  LinkToWarehouse(itemID, warehouseID uuid.UUID) error
  UnlinkFromWarehouse(itemID, warehouseID uuid.UUID) error
  LinkToTransactionRecord(itemID, recordID uuid.UUID) error
  UnlinkFromTransactionRecord(itemID, recordID uuid.UUID) error
  LinkToTransactionFile(itemID, fileID uuid.UUID) error
  UnlinkFromTransactionFile(itemID, fileID uuid.UUID) error

  ListItems(f repos.ItemFilter) ([]*models.Item, error)
}

type iSvc struct {
  repo      repos.IRepo
}

func NewISvc(repo repos.IRepo) ISvc {
  return &iSvc{repo: repo}
}

func (s *iSvc) CreateItem(item models.Item) (*models.Item, error) {
  if item.Name == "" {
    return nil, fmt.Errorf("item name is required")
  }
  created, err := s.repo.Create(item)
  if err != nil {
    return nil, fmt.Errorf("repo create item error: %w", err)
  }
  return created, nil
}

func (s *iSvc) UpdateItemName(itemID uuid.UUID, newName string) error {
  if itemID == uuid.Nil {
    return nil, fmt.Errorf("invalid itemID")
  }
  if newName == "" {
    return fmt.Errorf("new name cannot be empty")
  }
  if err := s.repo.UpdateName(itemID, newName); err != nil {
    return fmt.Errorf("Failed to update item name: %w", err)
  }
  return nil
}

func (s *iSvc) GetItemByID(itemID uuid.UUID) (*models.Item, error) {
  if itemID == uuid.Nil {
    return nil, fmt.Errorf("Invalid itemID")
  }
  item, err := s.repo.GetByID(itemID)
  if err != nil {
    return nil, fmt.Errorf("get item by id error: %w", err)
  }
  return item, nil
}

func (s *iSvc) DeleteItem(itemID uuid.UUID) error {
  if itemID == uuid.Nil {
    return fmt.Errorf("invalid itemID")
  }
  if err := s.repo.Delete(itemID); err != nil {
    return fmt.Errorf("failed to delete item: %w", err)
  }
  return nil
}

func (s *iSvc) LinkToLocation(itemID, locationID uuid.UUID) error {
  if itemID == uuid.Nil {
    return fmt.Errorf("invalid itemID")
  }
  if locationID == uuid.Nil {
    return fmt.Errorf("invalid locationID")
  }
  return s.repo.LinkToLocation(itemID, locationID)
}

func (s *iSvc) UnlinkFromLocation(itemID, locationID uuid.UUID) error {
  if itemID == uuid.Nil {
    return fmt.Errorf("invalid ItemID")
  }
  if locationID == uuid.Nil {
    return fmt.Errorf("invalid LocationID")
  }
  return s.repo.UnlinkFromLocation(itemID, locationID)
}

func (s *iSvc) LinkToWarehouse(itemID, warehouseID uuid.UUID) error {
  if itemID == uuid.Nil {
    return fmt.Errorf("Invalid ItemID")
  }
  if warehouseID == uuid.Nil {
    return fmt.Errorf("Invalid WarehouseID")
  }
  return s.repo.LinkToWarehouse(itemID, warehouseID)
}

func (s *iSvc) UnlinkFromWarehouse(itemID, warehouseID uuid.UUID) error {
  if itemID == uuid.Nil {
    return fmt.Errorf("Invalid ItemID")
  }
  if warehouseID == uuid.Nil {
    return fmt.Errorf("Invalid WarehouseID")
  }
  return s.repo.UnlinkFromWarehouse(itemID, warehouseID)
}

func (s *iSvc) LinkToTransactionRecord(itemID, recordID uuid.UUID) error {
  if itemID == uuid.Nil {
    return fmt.Errorf("Invalid ItemID")
  }
  if recordID == uuid.Nil {
    return fmt.Errorf("Invalid RecordID")
  }
  return s.repo.LinkToTransactionRecord(itemID, recordID)
}

func (s *iSvc) UnlinkFromTransactionRecord(itemID, recordID uuid.UUID) error {
  if itemID == uuid.Nil {
    return fmt.Errorf("Invalid ItemID")
  }
  if recordID == uuid.Nil {
    return fmt.Errorf("Invalid RecordID")
  }
  return s.repo.UnlinkFromTransactionRecord(itemID, recordID)
}

func (s *iSvc) LinkToTransactionFile(itemID, fileID uuid.UUID) error {
  if itemID == uuid.Nil {
    return fmt.Errorf("Invalid ItemID")
  }
  if fileID == uuid.Nil {
    return fmt.Errorf("Invalid FileID")
  }
  return s.repo.LinkToTransactionFile(itemID, fileID)
}

func (s *iSvc) UnlinkFromTransactionFile(itemID, fileID uuid.UUID) error {
  if itemID == uuid.Nil {
    return fmt.Errorf("Invalid ItemID")
  }
  if fileID == uuid.Nil {
    return fmt.Errorf("Invalid FileID")
  }
  return s.repo.UnlinkFromTransactionFile(itemID, fileID)
}

func (s *iSvc) ListItems(f repos.ItemFilter) ([]*models.Item, error) {
  items, err := s.repo.ListItems(f)
  if err != nil {
    return nil, fmt.Errorf("Failed to list items: %w", err)
  }
  return items, nil
}
