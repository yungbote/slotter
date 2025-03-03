package repos

import (
  "fmt"

  "gorm.io/gorm"
  "github.com/google/uuid"
  "github.com/yungbote/slotter/backend/services/database/internal/models"
)

type TransactionRecordFilter struct {
  CompanyID       uuid.UUID
  WarehouseID     uuid.UUID
  LocationID      uuid.UUID
  FileID          uuid.UUID
  ItemID          uuid.UUID
  TransactionType string
  OrderNameLike   string
  StartDate       time.Time
  EndDate         time.Time
  SortField       string
  SortDir         string
}

type TRRepo interface {
  //GENERAL CRUD
  Create(record models.TransactionRecord) (*models.TransactionRecord, error)
  UpdateOrderName(recordID uuid.UUID, newOrderName string) error
  UpdateDescription(recordID uuid.UUID, newDescription string) error
  UpdateTransactionQuantity(recordID uuid.UUID, newTQuantity int64) error
  UpdateCompletedQuantity(recordID uuid.UUID, newCQuantity int64) error
  UpdateCompletedDate(recordID uuid.UUID, newDate time.Time) error
  UpdateTransactionType(recordID uuid.UUID, newType string) error
  GetByID(recordID uuid.UUID) (*models.TransactionRecord, error)
  ListTransactionRecords(f TransactionRecordFilter) ([]*models.TransactionRecord, error)
}

type trRepo struct {
  db *gorm.DB
}

func NewTRRepo(db *gorm.DB) TRRepo {
  return &trRepo{db: db}
}

func (r *trRepo) Create(record models.TransactionRecord) (*models.TransactionRecord, error) {
  if err := r.db.Create(&record).Error; err != nil {
    return nil, fmt.Errorf("Failed to create transaction record: %w", err)
  }
  return &record, nil
}

func (r *trRepo) UpdateOrderName(recordID uuid.UUID, newOrderName string) error {
  return r.db.Model(&models.TransactionRecord{}).
    Where("id = ?", recordID).
    Update("order_name", newOrderName).Error
}

func (r *trRepo) UpdateDescription(recordID uuid.UUID, newDescription string) error {
  return r.db.Model(&models.TransactionRecord{}).
    Where("id = ?", recordID).
    Update("description", newDescription).Error
}

func (r *trRepo) UpdateTransactionQuantity(recordID uuid.UUID, newTQuantity int64) error {
  return r.db.Model(&models.TransactionRecord{}).
    Where("id = ?", recordID).
    Update("transaction_quantity", newTQuantity).Error
}

func (r *trRepo) UpdateCompletedQuantity(recordID uuid.UUID, newCQuantity int64) error {
  return r.db.Model(&models.TransactionRecord{}).
    Where("id = ?", recordID).
    Update("completed_quantity", newCQuantity).Error
}

func (r *trRepo) UpdateCompletedDate(recordID uuid.UUID, newDate time.Time) error {
  return r.db.Model(&models.TransactionRecord{}).
    Where("id = ?", recordID).
    Update("completed_date", newDate).Error
}

func (r *trRepo) UpdateTransactionType(recordID uuid.UUID, newType string) error {
  return r.db.Model(&models.TransactionRecord{}).
    Where("id = ?", recordID).
    Update("transaction_type", newType).Error
}

func (r *trRepo) GetByID(recordID uuid.UUID) (*models.TransactionRecord, error) {
  var rec models.TransactionRecord
  if err := r.db.First(&rec, "id = ?", recordID).Error; err != nil {
    return nil, fmt.Errorf("transaction record not found: %w", err)
  }
  return &rec, nil
}

func (r *trRepo) ListTransactionRecords(f TransactionRecordFilter) ([]*models.TransactionRecord, error) {
  dbq := r.db.Model(&models.TransactionRecord{}).Select("DISTINCT transaction_records.*")
  if f.CompanyID != uuid.Nil {
    dbq = dbq.Where("transaction_records.company_id = ?", f.CompanyID)
  }
  if f.WarehouseID != uuid.Nil {
    dbq = dbq.Where("transaction_records.warehouse_id = ?", f.WarehouseID)
  }
  if f.LocationID != uuid.Nil {
    dbq = dbq.Where("transaction_records.location_id = ?", f.LocationID)
  }
  if f.FileID != uuid.Nil {
    dbq = dbq.Where("transaction_records.transaction_file_id = ?", f.FileID)
  }
  if f.ItemID != uuid.Nil {
    dbq = dbq.Where("transaction_records.item_id = ?", f.ItemID)
  }
  if f.TransactionType != "" {
    dbq = dbq.Where("transaction_records.transaction_type = ?", f.TransactionType)
  }
  if f.OrderNameLike != "" {
    dbq = dbq.Where("transaction_records.order_name ILIKE ?", "%"+f.OrderNameLike+"%")
  }
  if !f.StartDate.IsZero() {
    dbq = dbq.Where("transaction_records.completed_date >= ?", f.StartDate)
  }
  if !f.EndDate.IsZero() {
    dbq = dbq.Where("transaction_records.completed_date <= ?", f.EndDate)
  }
  allowed := []string{"order_name", "transaction_type", "created_at", "updated_at", "completed_date"}
  dbq = applySorting(dbq, f.SortField, f.SortDir, allowed)
  var recs []*models.TransactionRecord
  if err := dbq.Find(&recs).Error; err != nil {
    return nil, err
  }
  return recs, nil
}

func applySorting(dbq *gorm.DB, sortField, sortDir string, allowedFields []string) *gorm.DB {
  found := false
  for _, af := range allowedFields {
    if af == sortField {
      found = true
      break
    }
  }
  if !found {
    sortField = "created_at"
  }
  if sortDir != "asc" && sortDir != "ASC" {
    sortDir = "DESC"
  }
  return dbq.Order(fmt.Sprintf("%s %s", sortField, sortDir))
}
