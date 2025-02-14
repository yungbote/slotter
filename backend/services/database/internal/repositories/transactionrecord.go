package repositories

import (
  "errors"
  "gorm.io/gorm"
  "github.com/yungbote/slotter/backend/services/database/internal/models"
)

type TransactionRecordRepository interface {
  Create(record *models.TransactionRecord) error
  GetByID(id uint) (*models.TransactionRecord, error)
  Update(record *models.TransactionRecord) error
  Delete(record *models.TransactionRecord) error
}

type transactionRecordRepository struct {
  db *gorm.DB
}

func NewTransactionRecordRepository(db *gorm.DB) TransactionRecordRepository {
  return &transactionRecordRepository{db: db}
}

func (r *transactionRecordRepository) Create(record *models.TransactionRecord) error {
  return r.db.Create(record).Error
}

func (r *transactionRecordRepository) GetByID(id uint) (*models.TransactionRecord, error) {
  var tr models.TransactionRecord
  err := r.db.First(&tr, id).Error
  if errors.Is(err, gorm.ErrRecordNotFound) {
    return nil, ErrNotFound
  }
  if err != nil {
    return nil, err
  }
  return &tr, nil
}

func (r *transactionRecordRepository) Update(record *models.TransactionRecord) error {
  return r.db.Save(record).Error
}

func (r *transactionRecordRepository) Delete(record *models.TransactionRecord) error {
  return r.db.Delete(record).Error
}
