package services

import (
  "github.com/yungbote/slotter/backend/services/database/internal/models"
  "github.com/yungbote/slotter/backend/services/database/internal/repositories"
)

type TransactionRecordService interface {
  CreateTransactionRecord(record *models.TransactionRecord) (*models.TransactionRecord, error)
  GetTransactionRecordByID(id uint) (*models.TransactionRecord, error)
  UpdateTransactionRecord(existing *models.TransactionRecord, updates *models.TransactionRecord) (*models.TransactionRecord, error)
  DeleteTransactionRecord(existing *models.TransactionRecord) error
}

type transactionRecordService struct {
  repo repositories.TransactionRecordRepository
}

func NewTransactionRecordService(repo repositories.TransactionRecordRepository) TransactionRecordService {
  return &transactionRecordService{repo: repo}
}

func (s *transactionRecordService) CreateTransactionRecord(input *models.TransactionRecord) (*models.TransactionRecord, error) {
  if err := s.repo.Create(input); err != nil {
    return nil, err
  }
  return input, nil
}

func (s *transactionRecordService) GetTransactionRecordByID(id uint) (*models.TransactionRecord, error) {
  return s.repo.GetByID(id)
}

func (s *transactionRecordService) UpdateTransactionRecord(existing *models.TransactionRecord, updates *models.TransactionRecord) (*models.TransactionRecord, error) {
  existing.WarehouseID = updates.WarehouseID
  existing.Warehouse = updates.Warehouse
  existing.TransactionType = updates.TransactionType
  existing.OrderNumber = updates.OrderNumber
  existing.ItemNumber = updates.ItemNumber
  existing.Description = updates.Description
  existing.TransactionQuantity = updates.TransactionQuantity
  existing.Location = updates.Location
  existing.Zone = updates.Zone
  existing.Carousel = updates.Carousel
  existing.Row = updates.Row
  existing.Shelf = updates.Shelf
  existing.Bin = updates.Bin
  existing.CompletedDate = updates.CompletedDate
  existing.CompletedBy = updates.CompletedBy
  existing.CompletedQuantity = updates.CompletedQuantity
  if err := s.repo.Update(existing); err != nil {
    return nil, err
  }
  return existing, nil
}

func (s *transactionRecordService) DeleteTransactionRecord(existing *models.TransactionRecord) error {
  return s.repo.Delete(existing)
}
