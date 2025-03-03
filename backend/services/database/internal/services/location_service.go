package services

import (
  "fmt"
  "github.com/google/uuid"
  "github.com/yungbote/slotter/backend/services/database/internal/models"
  "github.com/yungbote/slotter/backend/services/database/internal/repos"
)

type LSvc interface {
  //GENERAL CRUD
  CreateLocation(location models.Location) (*models.Location, error)
  UpdateLocationName(locationID uuid.UUID, newName string) error
  UpdateLocationNumber(locationID uuid.UUID, newNumber string) error
  GetLocationByID(locationID uuid.UUID) (*models.Location, error)
  GetLocationByPath(companyID, warehouseID uuid.UUID, locationPath string) (*models.Location, error)
  DeleteLocation(locationID uuid.UUID) error
  
  LinkToItem(locationID, itemID uuid.UUID) error
  UnlinkFromItem(locationID, itemID uuid.UUID) error
  LinkToTransactionFile(locationID, fileID uuid.UUID) error
  UnlinkFromTransactionFile(locationID, fileID uuid.UUID) error

  ListLocations(f repos.LocationFilter) ([]*models.Location, error)

}

type lSvc struct {
  repo            repos.LRepo
}

func NewLSvc(repo repos.LRepo) LSvc {
  return &lSvc{repo: repo}
}

func (s *lSvc) CreateLocation(location models.Location) (*models.Location, error) {
  
  if location.WarehouseID == nil || *location.WarehouseID == uuid.Nil {
    return nil, fmt.Errorf("location must have a warehouseID")
  }
  if location.LocationPath == "" {
    return nil, fmt.Errorf("location path is required")
  }
  created, err := s.repo.Create(location)
  if err != nil {
    return nil, fmt.Errorf("repo create location error: %w", err)
  }
  return created, nil
}

func (s *lSvc) UpdateLocationPath(locationID uuid.UUID, newPath string) error {
  if locationID == uuid.Nil || newPath == "" {
    return fmt.Errorf("invalid input to update location path")
  }
  if err := s.repo.UpdatePath(locationID, newPath); err != nil {
    return fmt.Errorf("Failed to update location path: %w", err)
  }
  return nil
}

func (s *lSvc) UpdateLocationNamePath(locationID uuid.UUID, newNamePath string) error {
  if locationID == uuid.Nil || newNamePath == "" {
    return fmt.Errorf("invalid input to update name path")
  }
  if err := s.repo.UpdateNamePath(locationID, newNamePath); err != nil {
    return fmt.Errorf("Failed to update location name path: %w", err)
  }
  return nil
}

func (s *lSvc) GetLocationByID(locationID uuid.UUID) (*models.Location, error) {
  if locationID == uuid.Nil {
    return nil, fmt.Errorf("invalid locationID")
  }
  location, err := s.repo.GetByID(locationID)
  if err != nil {
    return nil, fmt.Errorf("Failed to get location: %w", err)
  }
  return location, nil
}

func (s *lSvc) DeleteLocation(locationID uuid.UUID) error {
  if locationID == uuid.Nil {
    return fmt.Errorf("Invalid locationID")
  }
  if err := s.repo.Delete(locationID); err != nil {
    return fmt.Errorf("Failed to delete location: %w", err)
  }
  return nil
}

func (s *lSvc) LinkToItem(locationID, itemID uuid.UUID) error {
  if locationID == uuid.Nil {
    return fmt.Errorf("Invalid LocationID")
  }
  if itemID == uuid.Nil {
    return fmt.Errorf("Invalid ItemID")
  }
  return s.repo.LinkToItem(locationID, itemID)
}

func (s *lSvc) UnlinkFromItem(locationID, itemID uuid.UUID) error {
  if locationID == uuid.Nil {
    return fmt.Errorf("Invalid LocationID")
  }
  if itemID == uuid.Nil {
    return fmt.Errorf("Invalid ItemID")
  }
  return s.repo.UnlinkFromItem(locationID, itemID)
}

func (s *lSvc) LinkToTransactionFile(locationID, fileID uuid.UUID) error {
  if locationID == uuid.Nil {
    return fmt.Errorf("Invalid LocationID")
  }
  if fileID == uuid.Nil {
    return fmt.Errorf("Invalid FileID")
  }
  return s.repo.LinkToTransactionFile(locationID, fileID)
}

func (s *lSvc) UnlinkFromTransactionFile(locationID, fileID uuid.UUID) error {
  if locationID == uuid.Nil {
    return fmt.Errorf("Invalid LocationID")
  }
  if fileID == uuid.Nil {
    return fmt.Errorf("Invalid FileID")
  }
  return s.repo.UnlinkFromTransactionFile(locationID, fileID)
}

func (s *lSvc) ListLocations(f repos.LocationFilter) ([]*models.Location, error) {
  locations, err := s.repo.ListLocations(f)
  if err != nil {
    return nil, fmt.Errorf("Failed to list locations: %w", err)
  }
  return locations, nil
}

