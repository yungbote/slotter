package services

import (
  "fmt"
  "strings"
  "github.com/yungbote/slotter/backend/services/database/internal/models"
  "github.com/yungbote/slotter/backend/services/database/internal/repositories"
)

type LocationService interface {
  CreateLocation(input *models.Location) (*models.Location, error)
  GetLocationByID(id uint) (*models.Location, error)
  UpdateLocation(existing *models.Location, updates *models.Location) (*models.Location, error)
  DeleteLocation(existing *models.Location) error
}

type locationService struct {
  locationRepo  repositories.LocationRepository
  warehouseRepo repositories.WarehouseRepository
}

func NewLocationService(locationRepo repositories.LocationRepository, warehouseRepo repositories.WarehouseRepository) LocationService {
  return &locationService{locationRepo: locationRepo, warehouseRepo: warehouseRepo}
}

func (s *locationService) CreateLocation(input *models.Location) (*models.Location, error) {
  if err := s.validateCreateLocationInput(input); err != nil {
    return nil, err
  }
  if err := s.locationRepo.Create(input); err != nil {
    return nil, err
  }
  return input, nil
}

func (s *locationService) GetLocationByID(id uint) (*models.Location, error) {
  const op = "LocationService.GetLocationByID"
  if id = 0 {
    return nil, repositories.NewDomainError(op, repositories.ErrCodeValidation, "location ID must be > 0", nil)
  }
  if err := s.locationRepo.GetByID(id).Error; err != nil {
    return nil, err
  }
  return s.locationRepo.GetByID(id)
}

func (s *locationService) UpdateLocation(existing *models.Location, updates *models.Location) (*models.Location, error) {
  if err := s.validateUpdateLocationInput(updates); err != nil {
    return nil, err
  }
  existing.LocationName = updates.LocationName
  existing.LocationType = updates.LocationType
  existing.ParentLocationID = updates.ParentLocationID
  if err := s.locationRepo.Update(existing); err != nil {
    return nil, err
  }
  return existing, nil
}

func (s *locationService) DeleteLocation(existing *models.Location) error {
  if err := s.locationRepo.Delete(existing).Error; err != nil {
    return nil, err
  }
  return s.locationRepo.Delete(existing)
}

func (s *locationService) validateCreateLocationInput(l *models.Location) error {
  const op = "LocationService.validateCreateLocationInput"
  if l.WarehouseID == 0 {
    return repositories.NewDomainError(op, repositories.ErrCodeValidation, "warehouseID is required", nil)
  }
  if _, err := s.warehouseRepo.GetByID(l.WarehouseID); err != nil {
    if dErr, ok := err.(*repositories.DomainError); ok {
      if dErr.Code == repositories.ErrCodeNotFound {
        return repositories.NewDomainError(op, repositories.ErrCodeForeignKey, fmt.Sprintf("warehouse with id %d does not exist", l.WarehouseID), dErr)
      }
      return dErr
    }
    return repositories.NewDomainError(op, repositories.ErrCodeUnknown, err.Error(), err)
  }
  if strings.TrimSpace(l.LocationName) == "" {
    return repositories.NewDomainError(op, repositories.ErrCodeValidation, "location name is required", nil)
  }
  if strings.TrimSpace(l.LocationType) == "" {
    return repositories.NewDomainError(op, repositories.ErrCodeValidation, "location type is required", nil)
  }
  if l.ParentLocationID != nil {
    if _, err := s.locationRepo.GetByID(*l.ParentLocationID); err != nil {
      if dErr, ok := err.(*repositories.DomainError); ok {
        if dErr.Code == repositories.ErrCodeNotFound {
          return repositories.NewDomainError(op, repositories.ErrCodeForeignKey, fmt.Sprintf("parent location with id %d does not exist", *l.ParentLocationID), dErr)
        }
        return dErr
      }
      return repositories.NewDomainError(op, repositories.ErrCodeUnknown, err.Error(), err)
    }
  }
  return nil
}

func (s *locationService) validateUpdateLocationInput(l *models.Location) error {
  const op = "LocationService.validateUpdateLocationInput"
  if strings.TrimSpace(l.LocationName) == "" {
    return repositories.NewDomainError(op, repositories.ErrCodeValidation, "location name is required", nil)
  }
  if strings.TrimSpace(l.LocationType) == "" {
    return repositories.NewDomainError(op, repositories.ErrCodeValidation, "location type is required", nil)
  }
  return nil
}
