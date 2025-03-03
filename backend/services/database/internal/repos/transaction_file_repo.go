package repos

import (
  "time"
  "fmt"

  "github.com/google/uuid"
  "gorm.io/gorm"
  "github.com/yungbote/slotter/backend/services/database/internal/models"
)

type TransactionFileFilter struct {
    CompanyID   uuid.UUID
    WarehouseID uuid.UUID
    LocationID  uuid.UUID
    ItemID      uuid.UUID
    RecordID    uuid.UUID
    StartDate   time.Time
    EndDate     time.Time
    SortField   string
    SortDir     string
}

type TFRepo interface {
    //GENERAL CRUD
    Create(file models.TransactionFile) (*models.TransactionFile, error)
    UpdateName(fileID uuid.UUID, newName string) error
    UpdateExtension(fileID uuid.UUID, newExt string) error
    UpdateFilePathURL(fileID uuid.UUID, newPathURL string) error
    GetByID(fileID uuid.UUID) (*models.TransactionFile, error)
    Delete(fileID uuid.UUID) error
    //LINK & UNLINK TO LOCATIONS
    LinkToLocation(fileID, locationID uuid.UUID) error
    UnlinkFromLocation(fileID, locationID uuid.UUID) error
    //LINK & UNLINK TO ITEMS
    LinkToItem(fileID, itemID uuid.UUID) error
    UnlinkFromItem(fileID, itemID uuid.UUID) error
    ListTransactionFiles(f TransactionFileFilter) ([]*models.TransactionFile, error)
}

type tfRepo struct {
    db *gorm.DB
}

func NewTFRepo(db *gorm.DB) TFRepo {
    return &tfRepo{db: db}
}

func (r *tfRepo) Create(file models.TransactionFile) (*models.TransactionFile, error) {
    if err := r.db.Create(&file).Error; err != nil {
        return nil, fmt.Errorf("failed to create transaction file: %w", err)
    }
    return &file, nil
}

func (r *tfRepo) UpdateName(fileID uuid.UUID, newName string) error {
    return r.db.Model(&models.TransactionFile{}).
        Where("id = ?", fileID).
        Update("file_name", newName).Error
}

func (r *tfRepo) UpdateExtension(fileID uuid.UUID, newExt string) error {
    return r.db.Model(&models.TransactionFile{}).
        Where("id = ?", fileID).
        Update("file_extension", newExt).Error
}

func (r *tfRepo) UpdateFilePathURL(fileID uuid.UUID, newPathURL string) error {
    return r.db.Model(&models.TransactionFile{}).
        Where("id = ?", fileID).
        Update("file_path_url", newPathURL).Error
}

func (r *tfRepo) GetByID(fileID uuid.UUID) (*models.TransactionFile, error) {
    var f models.TransactionFile
    if err := r.db.First(&f, "id = ?", fileID).Error; err != nil {
        return nil, fmt.Errorf("transaction file not found: %w", err)
    }
    return &f, nil
}

func (r *tfRepo) Delete(fileID uuid.UUID) error {
    f, err := r.GetByID(fileID)
    if err != nil {
        return err
    }
    if err := r.db.Delete(f).Error; err != nil {
        return fmt.Errorf("failed to delete transaction file: %w", err)
    }
    return nil
}

func (r *tfRepo) LinkToLocation(fileID, locationID uuid.UUID) error {
    var file models.TransactionFile
    if err := r.db.First(&file, "id = ?", fileID).Error; err != nil {
        return fmt.Errorf("Failed to find file with ID: '%s': %w", fileID, err)
    }
    var loc models.Location
    if err := r.db.First(&loc, "id = ?", locationID).Error; err != nil {
        return fmt.Errorf("Failed to find location with ID: '%s': %w", locationID, err)
    }
    var count int64
    if err := r.db.Table("transaction_files_locations").Where("transaction_file_id = ? AND location_id = ?", fileID, locationID).Count(&count).Error; err != nil {
        return fmt.Errorf("Failed to check existing link: %w", err)
    }
    if count > 0 {
        return nil
    }
    if err := r.db.Model(&file).Association("Locations").Append(&loc); err != nil {
        return fmt.Errorf("Failed to link location (ID %s) to transaction file (ID %s): %w", locationID, fileID, err)
    }
    return nil
}

func (r *tfRepo) UnlinkFromLocation(fileID, locationID uuid.UUID) error {
    var file models.TransactionFile
    if err := r.db.First(&file, "id = ?", fileID).Error; err != nil {
        return fmt.Errorf("Failed to find file with ID: '%s': %w", fileID, err)
    }
    var loc models.Location
    if err := r.db.First(&loc, "id = ?", locationID).Error; err != nil {
        return fmt.Errorf("Failed to find location with ID: '%s': %w", locationID, err)
    }
    if err := r.db.Model(&file).Association("Locations").Delete(&loc); err != nil {
        return fmt.Errorf("Failed to unlink location (ID %s) from transaction file (ID %s): %w", locationID, fileID, err)
    }
    return nil
}

func (r *tfRepo) LinkToItem(fileID, itemID uuid.UUID) error {
    var file models.TransactionFile
    if err := r.db.First(&file, "id = ?", fileID).Error; err != nil {
        return fmt.Errorf("Failed to find file with ID: '%s': %w", fileID, err)
    }
    var item models.Item
    if err := r.db.First(&item, "id = ?", itemID).Error; err != nil {
        return fmt.Errorf("Failed to find item with ID: '%s': %w", itemID, err)
    }
    var count int64
    if err := r.db.Table("items_transaction_files").Where("item_id = ? AND transaction_file_id = ?", itemID, fileID).Count(&count).Error; err != nil {
        return fmt.Errorf("Failed to check existing link: %w", err)
    }
    if count > 0 {
        return nil
    }
    if err := r.db.Model(&file).Association("Items").Append(&item); err != nil {
        return fmt.Errorf("Failed to link item (ID %s) to transaction file (ID %s): %w", itemID, fileID, err)
    }
    return nil
}

func (r *tfRepo) UnlinkFromItem(fileID, itemID uuid.UUID) error {
    var file models.TransactionFile
    if err := r.db.First(&file, "id = ?", fileID).Error; err != nil {
        return fmt.Errorf("Failed to find file with ID: '%s': %w", fileID, err)
    }
    var item models.Item
    if err := r.db.First(&item, "id = ?", itemID).Error; err != nil {
        return fmt.Errorf("Failed to find item with ID: '%s': %w", itemID, err)
    }
    if err := r.db.Model(&file).Association("Items").Delete(&item); err != nil {
        return fmt.Errorf("Failed to unlink item (ID %s) from transaction file (ID %s): %w", itemID, fileID, err)
    }
    return nil
}

func (r *tfRepo) ListTransactionFiles(f TransactionFileFilter) ([]*models.TransactionFile, error) {
    dbq := r.db.Model(&models.TransactionFile{}).Select("DISTINCT transaction_files.*")
    if f.CompanyID != uuid.Nil {
        dbq = dbq.Where("transaction_files.company_id = ?", f.CompanyID)
    }
    if f.WarehouseID != uuid.Nil {
        dbq = dbq.Where("transaction_files.warehouse_id = ?", f.WarehouseID)
    }
    if f.LocationID != uuid.Nil {
        dbq = dbq.Joins("JOIN transaction_file_location tfl ON tfl.transaction_file_id = transaction_files.id").Where("tfl.location_id = ?", f.LocationID)
    }
    if f.ItemID != uuid.Nil {
        dbq = dbq.Joins("JOIN item_transaction_file itf ON itf.transaction_file_id = transaction_files.id").Where("itf.item_id = ?", f.ItemID)
    }
    if f.RecordID != uuid.Nil {
        dbq = dbq.Joins("JOIN transaction_records tr ON tr.transaction_file_id = transaction_files.id").Where("tr.id = ?", f.RecordID)
    }
    if !f.StartDate.IsZero() || !f.EndDate.IsZero() {
        dbq = dbq.Joins("JOIN transaction_records trDate ON trDate.transaction_file_id = transaction_files.id")
        if !f.StartDate.IsZero() {
            dbq = dbq.Where("trDate.completed_date >= ?", f.StartDate)
        }
        if !f.EndDate.IsZero() {
            dbq = dbq.Where("trDate.completed_date <= ?", f.EndDate)
        }
    }
    allowed := []string{"file_name", "file_extension", "created_at", "updated_at"}
    dbq = applySorting(dbq, f.SortField, f.SortDir, allowed)
    var files []*models.TransactionFile
    if err := dbq.Find(&files).Error; err != nil {
        return nil, err
    }
    return files, nil
}


