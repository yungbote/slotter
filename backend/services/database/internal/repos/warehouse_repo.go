package repos

import (
  "time"
  "fmt"

  "gorm.io/gorm"
  "github.com/google/uuid"
  "github.com/yungbote/slotter/backend/services/database/internal/models"
)

type WarehouseFilter struct {
    CompanyID       uuid.UUID
    ItemID          uuid.UUID
    LocationID      uuid.UUID
    FileID          uuid.UUID
    StartDate       time.Time
    EndDate         time.Time
    SortField       string
    SortDir         string
}

type WRepo interface {
    //GENERAL CRUD
    Create(warehouse models.Warehouse) (*models.Warehouse, error)
    UpdateName(warehouseID uuid.UUID, newName string) error
    GetByID(warehouseID uuid.UUID) (*models.Warehouse, error)
    Delete(warehouseID uuid.UUID) error
    //LINK & UNLINK TO ITEMS
    LinkToItem(warehouseID, itemID uuid.UUID) error
    UnlinkFromItem(warehouseID, itemID uuid.UUID) error
    ListWarehouses(f WarehouseFilter) ([]*models.Warehouse, error)
}

type wRepo struct {
    db *gorm.DB
}

func NewWRepo(db *gorm.DB) WRepo {
    return &wRepo{db: db}
}

func (r *wRepo) Create(warehouse models.Warehouse) (*models.Warehouse, error) {
    if err := r.db.Create(&warehouse).Error; err != nil {
        return nil, fmt.Errorf("failed to create warehouse: %w", err)
    }
    return &warehouse, nil
}

func (r *wRepo) UpdateName(warehouseID uuid.UUID, newName string) error {
    return r.db.Model(&models.Warehouse{}).
        Where("id = ?", warehouseID).
        Update("name", newName).Error
}

func (r *wRepo) GetByID(warehouseID uuid.UUID) (*models.Warehouse, error) {
    var wh models.Warehouse
    if err := r.db.First(&wh, "id = ?", warehouseID).Error; err != nil {
        return nil, fmt.Errorf("warehouse not found: %w", err)
    }
    return &wh, nil
}

func (r *wRepo) Delete(warehouseID uuid.UUID) error {
    wh, err := r.GetByID(warehouseID)
    if err != nil {
        return err
    }
    if err := r.db.Delete(wh).Error; err != nil {
        return fmt.Errorf("failed to delete warehouse: %w", err)
    }
    return nil
}

func (r *wRepo) LinkToItem(warehouseID, itemID uuid.UUID) error {
    var warehouse models.Warehouse
    if err := r.db.First(&warehouse, "id = ?", warehouseID).Error; err != nil {
        return fmt.Errorf("Failed to find warehouse with ID: '%s': %w", warehouseID, err)
    }
    var item models.Item
    if err := r.db.First(&item, "id = ?", itemID).Error; err != nil {
        return fmt.Errorf("Failed to find item with ID: '%s': %w", itemID, err)
    }
    var count int64
    if err := r.db.Table("items_warehouses").Where("item_id = ? AND warehouse_id = ?", itemID, warehouseID).Count(&count).Error; err != nil {
        return fmt.Errorf("Failed to check existing link: %w", err)
    }
    if count > 0 {
        return nil
    }
    if err := r.db.Model(&warehouse).Association("Items").Append(&item); err != nil {
        return fmt.Errorf("Failed to append item (ID %s) to warehouse (ID %s): %w", itemID, warehouseID, err)
    }
    return nil
}

func (r *wRepo) UnlinkFromItem(warehouseID, itemID uuid.UUID) error {
    var warehouse models.Warehouse
    if err := r.db.First(&warehouse, "id = ?", warehouseID).Error; err != nil {
        return fmt.Errorf("Failed to find warehouse with ID: '%s': %w", warehouseID, err)
    }
    var item models.Item
    if err := r.db.First(&item, "id = ?", itemID).Error; err != nil {
        return fmt.Errorf("Failed to find item with ID: '%s': %w", itemID, err)
    }
    if err := r.db.Model(&warehouse).Association("Items").Delete(&item); err != nil {
        return fmt.Errorf("Failed to unlink item (ID %s) from warehouse (ID %s): %w", itemID, warehouseID, err)
    }
    return nil
}

func (r *wRepo) ListWarehouses(f WarehouseFilter) ([]*models.Warehouse, error) {
    dbq := r.db.Model(&models.Warehouse{}).Select("DISTINCT warehouses.*")
    if f.CompanyID != uuid.Nil {
        dbq = dbq.Where("warehouses.company_id = ?", f.CompanyID)
    }
    if f.ItemID != uuid.Nil {
        dbq = dbq.Joins("JOIN item_warehouses iw ON iw.warehouse_id = warehouses.id").Where("iw.item_id = ?", f.ItemID)
    }
    if f.LocationID != uuid.Nil {
        dbq = dbq.Joins("JOIN locations l on l.warehouses_id = warehouses.id").Where("l.id = ?", f.LocationID)
    }
    if f.FileID != uuid.Nil {
        dbq = dbq.Joins("JOIN transaction_files tf ON tf.warehouse_id = warehouses.id").Where("tf.id = ?", f.FileID)
    }
    if !f.StartDate.IsZero() || !f.EndDate.IsZero() {
        dbq = dbq.Joins("JOIN transaction_records tr ON tr.warehouse_id = warehouses.id")
        if !f.StartDate.IsZero() {
            dbq = dbq.Where("tr.completed_date >= ?", f.StartDate)
        }
        if !f.EndDate.IsZero() {
            dbq = dbq.Where("tr.completed_date <= ?", f.EndDate)
        }
    }
    allowed := []string{"name", "created_at", "updated_at"}
    dbq = applySorting(dbq, f.SortField, f.SortDir, allowed)
    var warehouses []*models.Warehouse
    if err := dbq.Find(&warehouses).Error; err != nil {
        return nil, err
    }
    return warehouses, nil
}


