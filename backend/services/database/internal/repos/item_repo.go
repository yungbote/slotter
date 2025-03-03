package repos

import (
  "time"
  "fmt"

  "gorm.io/gorm"
  "github.com/google/uuid"
  "github.com/yungbote/slotter/backend/services/database/internal/models"
)

type ItemFilter struct {
    CompanyID   uuid.UUID
    WarehouseID uuid.UUID
    LocationID  uuid.UUID
    FileID      uuid.UUID
    RecordID    uuid.UUID
    StartDate   time.Time
    EndDate     time.Time
    SortField   string
    SortDir     string
}

type IRepo interface {
    //GENERAL CRUD
    Create(item models.Item) (*models.Item, error)
    UpdateName(itemID uuid.UUID, newName string) error
    GetByID(itemID uuid.UUID) (*models.Item, error)
    GetByNameAndCompanyID(companyID uuid.UUID, name string) (*models.Item, error)
    Delete(itemID uuid.UUID) error
    //GET'S BY TRANSACTIONRECORDID
    GetByTransactionRecordID(recordID uuid.UUID) (*models.Item, error)
    //LinkItemToLocation
    LinkToLocation(itemID, locationID uuid.UUID) error
    UnlinkFromLocation(itemID, locationID uuid.UUID) error
    //LINK & UNLINK ITEM TO WAREHOUSE
    LinkToWarehouse(itemID, warehouseID uuid.UUID) error
    UnlinkFromWarehouse(itemID, warehouseID uuid.UUID) error
    //LINK & UNLINK ITEM TO TRANSACTIONRECORD
    LinkToTransactionRecord(itemID, recordID uuid.UUID) error
    UnlinkFromTransactionRecord(itemID, recordID uuid.UUID) error
    //LINK & UNLINK ITEM TO TRANSACTIONFILE
    LinkToTransactionFile(itemID, fileID uuid.UUID) error
    UnlinkFromTransactionFile(itemID, fileID uuid.UUID) error
    ListItems(f ItemFilter) ([]*models.Item, error)
}

type iRepo struct {
    db *gorm.DB
}

func NewIRepo(db *gorm.DB) IRepo {
    return &iRepo{db: db}
}

func (r *iRepo) Create(item models.Item) (*models.Item, error) {
    if err := r.db.Create(&item).Error; err != nil {
        return nil, fmt.Errorf("failed to create item: %w", err)
    }
    return &item, nil
}

func (r *iRepo) UpdateName(itemID uuid.UUID, newName string) error {
    return r.db.Model(&models.Item{}).
        Where("id = ?", itemID).
        Update("name", newName).Error
}

func (r *iRepo) GetByID(itemID uuid.UUID) (*models.Item, error) {
    var i models.Item
    if err := r.db.First(&i, "id = ?", itemID).Error; err != nil {
        return nil, fmt.Errorf("item not found: %w", err)
    }
    return &i, nil
}

func (r *iRepo) Delete(itemID uuid.UUID) error {
    i, err := r.GetByID(itemID)
    if err != nil {
        return err
    }
    if err := r.db.Delete(i).Error; err != nil {
        return fmt.Errorf("failed to delete item: %w", err)
    }
    return nil
}

func (r *iRepo) GetByTransactionRecordID(recordID uuid.UUID) (*models.Item, error) {
    // transaction record references itemID, so let's join or do a two-step:
    var tr models.TransactionRecord
    if err := r.db.First(&tr, "id = ?", recordID).Error; err != nil {
        return nil, fmt.Errorf("transaction record not found: %w", err)
    }
    var i models.Item
    if err := r.db.First(&i, "id = ?", tr.ItemID).Error; err != nil {
        return nil, fmt.Errorf("item not found for transaction record: %w", err)
    }
    return &i, nil
}

func (r *iRepo) LinkToLocation(itemID, locationID uuid.UUID) error {
    var item models.Item
    if err := r.db.First(&item, "id = ?", itemID).Error; err != nil {
        return fmt.Errorf("Failed to find item with ID: '%s': %w", itemID, err)
    }
    var loc models.Location
    if err := r.db.First(&loc, "id = ?", locationID).Error; err != nil {
        return fmt.Errorf("Failed to find location with ID: '%s': %w", locationID, err)
    }
    var count int64
    if err := r.db.Table("items_locations").Where("item_id = ? AND location_id = ?", itemID, locationID).Count(&count).Error; err != nil {
        return fmt.Errorf("Failed to check existing link: %w", err)
    }
    if count > 0 {
        return nil
    }
    if err := r.db.Model(&item).Association("Locations").Append(&loc); err != nil {
        return fmt.Errorf("Failed to append location (ID %s) to item (ID %s): %w", locationID, itemID, err)
    }
    return nil
}

func (r *iRepo) UnlinkFromLocation(itemID, locationID uuid.UUID) error {
    var item models.Item
    if err := r.db.First(&item, "id = ?", itemID).Error; err != nil {
        return fmt.Errorf("Failed to find item with ID: '%s': %w", itemID, err)
    }
    var loc models.Location
    if err := r.db.First(&loc, "id = ?", locationID).Error; err != nil {
        return fmt.Errorf("Failed to find location with ID: '%s': %w", locationID, err)
    }
    if err := r.db.Model(&item).Association("Locations").Delete(&loc); err != nil {
        return fmt.Errorf("Failed to unlink location (ID %s) from item (ID %s): %w", locationID, itemID, err)
    }
    return nil
}

func (r *iRepo) LinkToWarehouse(itemID, warehouseID uuid.UUID) error {
    var item models.Item
    if err := r.db.First(&item, "id = ?", itemID).Error; err != nil {
        return fmt.Errorf("Failed to find item with ID: '%s': %w", itemID, err)
    }
    var warehouse models.Warehouse
    if err := r.db.First(&warehouse, "id = ?", warehouseID).Error; err != nil {
        return fmt.Errorf("Failed to find warehouse with ID: '%s': %w", warehouseID, err)
    }
    var count int64
    if err := r.db.Table("items_warehouses").Where("item_id = ? AND warehouse_id = ?", itemID, warehouseID).Count(&count).Error; err != nil {
        return fmt.Errorf("Failed to check existing link: %w", err)
    }
    if count > 0 {
        return nil
    }
    if err := r.db.Model(&item).Association("Warehouses").Append(&warehouse); err != nil {
        return fmt.Errorf("Failed to append warehouse (ID %s) to item (ID %s): %w", warehouseID, itemID, err)
    }
    return nil
}

func (r *iRepo) UnlinkFromWarehouse(itemID, warehouseID uuid.UUID) error {
    var item models.Item
    if err := r.db.First(&item, "id = ?", itemID).Error; err != nil {
        return fmt.Errorf("Failed to find item with ID: '%s': %w", itemID, err)
    }
    var warehouse models.Warehouse
    if err := r.db.First(&warehouse, "id = ?", warehouseID).Error; err != nil {
        return fmt.Errorf("Failed to find warehouse with ID: '%s': %w", warehouseID, err)
    }
    if err := r.db.Model(&item).Association("Warehouses").Delete(&warehouse); err != nil {
        return fmt.Errorf("Failed to unlink warehouse (ID %s) from item (ID %s): %w", warehouseID, itemID, err)
    }
    return nil
}

func (r *iRepo) LinkToTransactionRecord(itemID, recordID uuid.UUID) error {
    var item models.Item
    if err := r.db.First(&item, "id = ?", itemID).Error; err != nil {
        return fmt.Errorf("Failed to find item with ID: '%s': %w", itemID, err)
    }
    var record models.TransactionRecord
    if err := r.db.First(&record, "id = ?", recordID).Error; err != nil {
        return fmt.Errorf("Failed to find transaction record with ID: '%s': %w", recordID, err)
    }
    if err := r.db.Model(&item).Association("TransactionRecords").Append(&record); err != nil {
        return fmt.Errorf("Failed to append transaction record (ID %s) to item (ID %s): %w", recordID, itemID, err)
    }
    return nil
}

func (r *iRepo) UnlinkFromTransactionRecord(itemID, recordID uuid.UUID) error {
    var item models.Item
    if err := r.db.First(&item, "id = ?", itemID).Error; err != nil {
        return fmt.Errorf("Failed to find item with ID: '%s': %w", itemID, err)
    }
    var record models.TransactionRecord
    if err := r.db.First(&record, "id = ?", recordID).Error; err != nil {
        return fmt.Errorf("Failed to find transaction record with ID: '%s': %w", recordID, err)
    }
    if err := r.db.Model(&item).Association("TransactionRecords").Delete(&record); err != nil {
        return fmt.Errorf("Failed to unlink transaction record (ID %s) from item (ID %s): %w", recordID, itemID, err)
    }
    return nil
}

func (r *iRepo) LinkToTransactionFile(itemID, fileID uuid.UUID) error {
    var item models.Item
    if err := r.db.First(&item, "id = ?", itemID).Error; err != nil {
        return fmt.Errorf("Failed to find item with ID: '%s': %w", itemID, err)
    }
    var file models.TransactionFile
    if err := r.db.First(&file, "id = ?", fileID).Error; err != nil {
        return fmt.Errorf("Failed to find transaction file with ID: '%s': %w", fileID, err)
    }
    var count int64
    if err := r.db.Table("items_transaction_files").Where("item_id = ? AND transaction_file_id = ?", itemID, fileID).Count(&count).Error; err != nil {
        return fmt.Errorf("Failed to check existing link: %w", err)
    }
    if count > 0 {
        return nil
    }
    if err := r.db.Model(&item).Association("TransactionFiles").Append(&file); err != nil {
        return fmt.Errorf("Failed to link transaction file (ID %s) to item (ID %s): %w", fileID, itemID, err)
    }
    return nil
}

func (r *iRepo) UnlinkFromTransactionFile(itemID, fileID uuid.UUID) error {
    var item models.Item
    if err := r.db.First(&item, "id = ?", itemID).Error; err != nil {
        return fmt.Errorf("Failed to find item with ID: '%s': %w", itemID, err)
    }
    var file models.TransactionFile
    if err := r.db.First(&file, "id = ?", fileID).Error; err != nil {
        return fmt.Errorf("Failed to find transaction file with ID: '%s': %w", fileID, err)
    }
    if err := r.db.Model(&item).Association("TransactionFiles").Delete(&file); err != nil {
        return fmt.Errorf("Failed to unlink transaction file (ID %s) from item (ID %s): %w", fileID, itemID, err)
    }
    return nil
}


func (r *iRepo) GetByNameAndCompanyID(companyID uuid.UUID, name string) (*models.Item, error) {
    var item models.Item
    if err := r.db.Where("company_id = ? AND name = ?", companyID, name).First(&item).Error; err != nil {
        return nil, fmt.Errorf("could not find item with name=%s in company_id=%s: %w", name, companyID, err)
    }
    return &item, nil
}

func (r *iRepo) ListItems(f ItemFilter) ([]*models.Item, error) {
    dbq := r.db.Model(&models.Item{}).Select("DISTINCT items.*")
    if f.CompanyID != uuid.Nil {
        dbq = dbq.Where("items.company_id = ?", f.CompanyID)
    }
    if f.WarehouseID != uuid.Nil {
        dbq = dbq.Joins("JOIN item_warehouses iw ON iw.item_id = items.id").Where("iw.warehouse_id = ?", f.WarehouseID)
    }
    if f.LocationID != uuid.Nil {
        dbq = dbq.Joins("JOIN item_locations il ON il.item_id = items.id").Where("il.location_id = ?", f.LocationID)
    }
    if f.FileID != uuid.Nil {
        dbq = dbq.Joins("JOIN item_transaction_files itf ON itf.item_id = items.id").Where("itf.transaction_file_id = ?", f.FileID)
    }
    if f.RecordID != uuid.Nil {
        dbq = dbq.Joins("JOIN transaction_records tr ON tr.item_id = items.id").Where("tr.id = ?", f.RecordID)
    }
    if !f.StartDate.IsZero() || !f.EndDate.IsZero() {
        dbq = dbq.Joins("JOIN transaction_records trDate ON trDate.item_id = items.id")
        if !f.StartDate.IsZero() {
            dbq = dbq.Where("trDate.completed_date >= ?", f.StartDate)
        }
        if !f.EndDate.IsZero() {
            dbq = dbq.Where("trDate.completed_date <= ?", f.EndDate)
        }
    }
    allowed := []string{"name", "created_at", "updated_at"}
    dbq = applySorting(dbq, f.SortField, f.SortDir, allowed)
    var items []*models.Item
    if err := dbq.Find(&items).Error; err != nil {
        return nil, err
    }
    return items, nil
}

func applySorting(dbq *gorm.DB, sortField, sortDir string, allowed []string) *gorm.DB {
    found := false
    for _, f := range allowed {
        if f == sortField {
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

