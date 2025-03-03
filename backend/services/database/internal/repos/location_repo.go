package repos

import (
  "fmt"

  "gorm.io/gorm"
  "github.com/google/uuid"
  "github.com/yungbote/slotter/backend/services/database/internal/models"
)

type LocationFilter struct {
  CompanyID     uuid.UUID
  WarehouseID   uuid.UUID
  ItemID        uuid.UUID
  FileID        uuid.UUID
  RecordID      uuid.UUID
  StartDate     time.Time
  EndDate       time.Time
  SortField     string
  SortDir       string
}

type LRepo interface {
  Create(location models.Location) (*models.Location, error)
  UpdatePath(locationID uuid.UUID, newName string) error
  UpdateNamePath(locationID uuid.UUID, newNumber string) error
  GetByID(locationID uuid.UUID) (*models.Location, error)
  GetByPath(warehouseID uuid.UUID, locationPath string) (*models.Location, error)
  Delete(locationID uuid.UUID) error
  //LINK & UNLINK TO ITEM
  LinkToItem(locationID, itemID uuid.UUID) error
  UnlinkFromItem(locationID, itemID uuid.UUID) error
  //LINK & UNLINK TO TRANSACTIONFILE
  LinkToTransactionFile(locationID, fileID uuid.UUID) error
  UnlinkFromTransactionFile(locationID, fileID uuid.UUID) error
  ListLocations(f LocationFilter) ([]*models.Location, error)
}

type lRepo struct {
  db *gorm.DB
}

func NewLRepo(db *gorm.DB) LRepo {
  return &lRepo{db: db}
}

func (r *lRepo) Create(location models.Location) (*models.Location, error) {
  if err := r.db.Create(&location).Error; err != nil {
    return nil, fmt.Errorf("Failed to create location: %w", err)
  }
  return &location, nil
}

func (r *lRepo) UpdatePath(locationID uuid.UUID, newPath string) error {
  if err := r.db.Model(&models.Location{}).
    Where("id = ?", locationID).
    Update("location_path", newPath).Error; err != nil {
    return fmt.Errorf("Failed to update location name: %w", err)
  }
  return nil
}

func (r *lRepo) UpdateNamePath(locationID uuid.UUID, newNamePath string) error {
  if err := r.db.Model(&models.Location{}).
    Where("id = ?", locationID).
    Update("location_name_path", newNamePath).Error; err != nil {
    return fmt.Errorf("Failed to update location number: %w", err)
  }
  return nil
}

func (r *lRepo) GetByID(locationID uuid.UUID) (*models.Location, error) {
  var loc models.Location
  if err := r.db.First(&loc, "id = ?", locationID).Error; err != nil {
    return nil, fmt.Errorf("Location not found: %w", err)
  }
  return &loc, nil
}

func (r *lRepo) GetByPath(warehouseID uuid.UUID, locationPath string) (*models.Location, error) {
  var loc models.Location
  if err := r.db.Where("warehouse_id = ? AND location_path = ?", warehouseID, locationPath).
    First(&loc).Error; err != nil {
    return nil, fmt.Errorf("Location not found for warehouse with id: '%w' by path: '%w', err: %w", warehouseID, locationPath, err)
  }
  return &loc, nil
}


func (r *lRepo) Delete(locationID uuid.UUID) error {
  loc, err := r.GetByID(locationID)
  if err != nil {
    return err
  }
  if err := r.db.Delete(loc).Error; err != nil {
    return fmt.Errorf("Failed to delete location: %w", err)
  }
  return nil
}

func (r *lRepo) ListByWarehouseID(warehouseID uuid.UUID, sortField, sortDir string) ([]*models.Location, error) {
  dbq := r.db.Model(&models.Location{}).
    Where("warehouse_id = ?", warehouseID)
  dbq = applySorting(dbq, sortField, sortDir, []string{"location_path", "location_name_path", "created_at"})
  var locs []*models.Location
  if err := dbq.Find(&locs).Error; err != nil {
    return nil, Errorf("Failed to list locations by warehouse: %w", err)
  }
  return locs, nil
}

func (r *lRepo) CountByWarehouseID(warehouseID uuid.UUID) (int64, error) {
  var count int64
  if err := r.db.Model(&models.Location{}).
    Where("warehouse_id = ?", warehouseID).
    Count(&count).Error; err != nil {
    return 0, fmt.Errorf("Failed to count locations by warehouseID: %w", err)
  }
  return count, nil
}

func (r *lRepo) ListByTransactionFileID(fileID uuid.UUID) ([]*models.Location, error) {
  dbq := r.db.Model(&models.Location{}).
    Joins("JOIN transaction_file_location tfl ON tfl.location_id = locations.id").
    Where("tfl.transaction_file_id = ?", fileID)
  dbq = applySorting(dbq, sortField, sortDir, []string{"location_name", "location_number", "created_at"})
  var locs []*models.Location
  if err := dbq.Find(&locs).Error; err != nil {
    return nil, fmt.Errorf("Failed to list locations by transaction file: %w", err)
  }
  return locs, nil
}

func (r *lRepo) CountByTransactionFileID(fileID uuid.UUID) (int64, error) {
  var count int64
  if err := r.db.Model(&models.Location{}).
    Joins("JOIN transaction_file_location tfl ON tfl.location_id = locations.id").
    Where("tfl.transaction_file_id = ?", fileID).
    Count(&count).Error; err != nil {
    return 0, fmt.Errorf("Failed to count locations by transaction file: %w", err)
  }
  return count, nil
}

func (r *lRepo) ListByItemID(itemID uuid.UUID, sortField, sortDir string) ([]*models.Location, error) {
  dbq := r.db.Model(&models.Location{}).
    Joins("JOIN item_location il ON il.location_id = locations.id").
    Where("il.item_id = ?", itemID)
  dbq = applySorting(dbq, sortField, sortDir, []string{"location_name", "location_number", "created_at"})
  var locs []*models.Location
  if err := dbq.Find(&locs).Error; err != nil {
    return nil, fmt.Errorf("Failed to list locations by item: %w", err)
  }
  return locs, nil
}

func (r *lRepo) CountByItemID(itemID uuid.UUID) (int64, error) {
  var count int64
  if err := r.db.Model(&models.Location{}).
    Joins("JOIN item_location il ON il.location_id = locations.id").
    Where("il.item_id = ?", itemID).
    Count(&count).Error; err != nil {
    return 0, fmt.Errorf("Failed to count locations by item: %w", err)
  }
  return count, nil
}

func (r *lRepo) LinkToItem(locationID, itemID, uuid.UUID) error {
  var loc models.Location
  if err := r.db.First(&loc, "id = ?", locationID).Error; err != nil {
    return fmt.Errorf("Failed to find location with ID '%s': %w", locationID, err)
  }
  var item models.Item
  if err := r.db.First(&item, "id = ?", itemID).Error; err != nil {
    return fmt.Errorf("Failed to find item with ID '%s': %w", itemID, err)
  }
  var count int64
  if err := r.db.Table("items_locations").Where("item_id = ? AND location_id = ?", itemID, locationID).Count(&count).Error; err != nil {
    return fmt.Errorf("Failed to check existing link: %w", err)
  }
  if count > 0 {
    return nil
  }
  if err := r.db.Model(&loc).Association("Items").Append(&item); err != nil {
    return fmt.Errorf("Failed to append item (ID %s) to location (ID %s): %w", itemID, err)
  }
  return nil
}

func (r *lRepo) UnlinkFromItem(locationID, itemID uuid.UUID) error {
  var loc models.Location
  if err := r.db.First(&loc, "id = ?", locationID).Error; err != nil {
    return fmt.Errorf("Failed to find location with ID '%s': %w", locationID, err)
  }
  var item models.Item
  if err := r.db.First(&item, "id = ?", itemID).Error; err != nil {
    return fmt.Errorf("Failed to find item with ID '%s': %w", itemID, err)
  }
  if err := r.db.Model(&loc).Association("Items").Delete(&item); err != nil {
    return fmt.Errorf("Failed to unlink item (ID %s) from location (ID %s): %w", itemID, locationID, err)
  }
  return nil
}

func (r *lRepo) LinkToTransactionFile(locationID, fileID uuid.UUID) error {
  var loc models.Location
  if err := r.db.First(&loc, "id = ?", locationID).Error; err != nil {
    return fmt.Errorf("Failed to find location with ID '%s': %w", locationID, err)
  }
  var file models.TransactionFile
  if err := r.db.First(&file, "id = ?", fileID).Error; err != nil {
    return fmt.Errorf("Failed to find transaction file with ID '%s': %w", fileID, err)
  }
  var count int64
  if err := r.db.Table("transaction_files_locations").Where("transaction_file_id = ? AND location_id = ?", fileID, locationID).Count(&count).Error; err != nil {
    return fmt.Errorf("Failed to check existing link: %w", err)
  }
  if count > 0 {
    return nil
  }
  if err := r.db.Model(&loc).Association("TransactionFiles").Append(&file); err != nil {
    return fmt.Errorf("Failed to link transaction file (ID %s) to location (ID %s): %w", fileID, locationID, err)
  }
  return nil
}

func (r *lRepo) UnlinkFromTransactionFile(locationID, fileID uuid.UUID) error {
  var loc models.Location
  if err := r.db.First(&loc, "id = ?", locationID).Error; err != nil {
    return fmt.Errorf("Failed to find location with ID: '%s': %w", locationID, err)
  }
  var file models.TransactionFile
  if err := r.db.First(&file, "id = ?", fileID).Error; err != nil {
    return fmt.Errorf("Failed to find transaction file with ID: '%s': %w", fileID, err)
  }
  if err := r.db.Model(&loc).Association("TransactionFiles").Delete(&file); err != nil {
    return fmt.Errorf("Failed to unlink transaction file (ID %s) from location (ID %s): %w", fileID, locationID, err)
  }
  return nil
}

func (r *lRepo) ListLocations(f LocationFilter) ([]*models.Location, error) {
  dbq := r.db.Model(&models.Location{}).Select("DISTINCT locations.*")
  if f.CompanyID != uuid.Nil {
    dbq = dbq.Joins("JOIN warehouses w ON w.id = locations.warehouse_id").Where("w.company_id = ?", f.CompanyID)
  }
  if f.WarehouseID != uuid.Nil {
    dbq = dbq.Where("locations.warehouse_id = ?", f.WarehouseID)
  }
  if f.ItemID != uuid.Nil {
    dbq = dbq.Joins("JOIN item_locations il ON il.location_id = locations.id").Where("il.item_id = ?", f.ItemID)
  }
  if f.FileID != uuid.Nil {
    dbq = dbq.Joins("JOIN transaction_file_location tfl ON tfl.location_id = locations.id").Where("tfl.transaction_file_id = ?", f.FileID)
  }
  if f.RecordID != uuid.Nil {
    dbq = dbq.Joins("JOIN transaction_records tr ON tr.location_id = locations.id").Where("tr.id = ?", f.RecordID)
  }
  if !f.StartDate.IsZero() || !f.EndDate.IsZero() {
    dbq = dbq.Joins("JOIN transaction_records trDate ON trDate.location_id = locations.id")
    if !f.StartDate.IsZero() {
      dbq = dbq.Where("trDate.completed_date >= ?", f.StartDate)
    }
    if !f.EndDate.IsZero() {
      dbq = dbq.Where("trDate.completed_date <= ?", f.EndDate)
    }
  }
  allowed := []string{"location_path", "location_name_path", "created_at", "updated_at"}
  dbq = applySorting(dbq, f.SortField, f.SortDir, allowed)
  var locs []*models.Location
  if err := dbq.Find(&locs).Error; err != nil {
    return nil, err
  }
  return locs, nil
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

