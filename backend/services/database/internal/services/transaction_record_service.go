package services

import (
  "fmt"
  "time"
  "github.com/google/uuid"
  "github.com/yungbote/slotter/backend/services/database/internal/models"
  "github.com/yungbote/slotter/backend/services/database/internal/repos"
)

type TRSvc interface {
  //GENERAL CRUD
  CreateTransactionRecord(record models.TransactionRecord) (*models.TransactionRecord, error)
  UpdateTransactionRecordOrderName(recordID uuid.UUID, newOrderName string) error
  UpdateTransactionRecordDescription(recordID uuid.UUID, newDescription string) error
  UpdateTransactionRecordTransactionQuantity(recordID uuid.UUID, newTQuantity int64) error
  UpdateTransactionRecordCompletedQuantity(recordID uuid.UUID, newCQuantity int64) error
  UpdateTransactionRecordCompletedDate(recordID uuid.UUID, newDate time.Time) error
  UpdateTransactionRecordTransactionType(recordID uuid.UUID, newType string) error
  GetTransactionRecordByID(recordID uuid.UUID) (*models.TransactionRecord, error)

  ListTransactionRecords(f repos.TransactionRecordFilter) ([]*models.TransactionRecord, error)
 }

type trSvc struct {
  repo              repos.TRRepo
}

func NewTRSvc(repo repos.TRRepo) TRSvc {
  return &trSvc{repo: repo}
}

func (s *trSvc) CreateTransactionRecord(record models.TransactionRecord) (*models.TransactionRecord, error) {
  if record.CompanyID == nil || *rec.CompanyID == uuid.Nil {
    return nil, fmt.Errorf("Transaction record must have a valid companyID")
  }
  created, err := s.repo.Create(record)
  if err != nil {
    return nil, fmt.Errorf("Failed to create transaction record: %w", err)
  }
  return created, nil
} 

func (s *trSvc) UpdateTransactionRecordOrderName(recordID uuid.UUID, newOrderName string) error {
  if recordID == uuid.Nil {
    return fmt.Errorf("Invalid recordID")
  }
  if newOrderName == "" {
    return fmt.Errorf("Invalid New Order Name")
  }
  return s.repo.UpdateOrderName(recordID, newOrderName)
}

func (s *trSvc) UpdateTransactionRecordDescription(recordID uuid.UUID, newDescription string) error {
  if recordID == uuid.Nil {
    return fmt.Errorf("Invalid recordID")
  }
  if newDescription == "" {
    return fmt.Errorf("Invalid New Description")
  }
  return s.repo.UpdateDescription(recordID, newDescription)
}

func (s *trSvc) UpdateTransactionRecordTransactionQuantity(recordID uuid.UUID, newTQuantity int64) error {
  if recordID == uuid.Nil {
    return fmt.Errorf("Invalid recordID")
  }
  return s.repo.UpdateTransactionQuantity(recordID, newTQuantity)
}

func (s *trSvc) UpdateTransactionRecordCompletedQuantity(recordID uuid.UUID, newCQuantity int64) error {
  if recordID == uuid.Nil {
    return fmt.Errorf("Invalid recordID")
  }
  return s.repo.UpdateCompletedQuantity(recordID, newCQuantity)
}

func (s *trSvc) UpdateTransactionRecordCompletedDate(recordID uuid.UUID, newDate time.Time) error {
  if recordID == uuid.Nil {
    return fmt.Errorf("Invalid recordID")
  }
  return s.repo.UpdateCompletedDate(recordID, newDate)
}

func (s *trSvc) UpdateTransactionRecordTransactionType(recordID uuid.UUID, newType string) error {

}

func (s *trSvc) GetTransactionRecordByID(recordID uuid.UUID) (*models.TransactionRecord, error) {
  if recordID == uuid.Nil {
    return nil, fmt.Errorf("InvalidRecordID")
  }
  rec, err := s.repo.GetByID(recordID)
  if err != nil {
    return nil, fmt.Errorf("repo get transaction record error: %w", err)
  }
  return rec, nil
}

func (s *trSvc) ListTransactionRecords(f repos.TransactionRecordFilter) ([]*models.TransactionRecord, error) {
  records, err := s.repo.ListTransactionRecords(f)
  if err != nil {
    return nil, fmt.Errorf("repo list transaction records error: %w", err)
  }
  return records, nil
}


