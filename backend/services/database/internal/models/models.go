package models

import (
  "time"
  "github.com/google/uuid"
  "gorm.io/datatypes"
)

// ----------------------------------------------------
// User
// ----------------------------------------------------
type User struct {
  ID            uuid.UUID     `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
  Email         string        `gorm:"uniqueIndex;not null"`
  Password      string        `gorm:"not null"`
  FirstName     string        `gorm:"not null"`
  LastName      string        `gorm:"not null"`
  CompanyID     *uuid.UUID    `gorm:"index"`
  Company       *Company
  Role          string        `gorm:"default:'user'"`
  CreatedAt     time.Time     `gorm:"not null;default:now()"`
  UpdatedAt     time.Time     `gorm:"not null;default:now()"`
  AvatarURL     string        `gorm:"column:profile_picture_url"`

  // If you want to load all user actions by user, define a has-many:
  UserActions   []*UserAction `gorm:"foreignKey:UserID"`
}

// ----------------------------------------------------
// Company
// ----------------------------------------------------
type Company struct {
  ID                  uuid.UUID               `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
  Name                string                  `gorm:"not null"`
  Users               []*User                 `gorm:"foreignKey:CompanyID"`
  Warehouses          []*Warehouse            `gorm:"foreignKey:CompanyID"`
  TransactionFiles    []*TransactionFile      `gorm:"foreignKey:CompanyID"`
  TransactionRecords  []*TransactionRecord    `gorm:"foreignKey:CompanyID"`
  Items               []*Item                 `gorm:"foreignKey:CompanyID"`
  AvatarURL           string                  `gorm:"column:avatar_url"`
  CreatedAt           time.Time               `gorm:"not null;default:now()"`
  UpdatedAt           time.Time               `gorm:"not null;default:now()"`
}

// ----------------------------------------------------
// Warehouse
// ----------------------------------------------------
type Warehouse struct {
  ID                  uuid.UUID             `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
  Name                string                `gorm:"not null"`
  CompanyID           *uuid.UUID            `gorm:"not null;index"`
  Company             *Company              `gorm:"constraint:OnDelete:CASCADE"`
  Locations           []*Location           `gorm:"foreignKey:WarehouseID"`
  TransactionFiles    []*TransactionFile    `gorm:"foreignKey:WarehouseID"`
  TransactionRecords  []*TransactionRecord  `gorm:"foreignKey:WarehouseID"`
  Items               []*Item               `gorm:"many2many:items_warehouses;"`
  CreatedAt           time.Time             `gorm:"not null;default:now()"`
  UpdatedAt           time.Time             `gorm:"not null;default:now()"`
}

// ----------------------------------------------------
// Location
// ----------------------------------------------------
type Location struct {
  ID                  uuid.UUID             `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
  WarehouseID         *uuid.UUID            `gorm:"not null;index"`
  Warehouse           *Warehouse            `gorm:"constraint:OnDelete:CASCADE"`
  Items               []*Item               `gorm:"many2many:items_locations;"`
  TransactionFiles    []*TransactionFile    `gorm:"many2many:transaction_files_locations;"`
  TransactionRecords  []*TransactionRecord  `gorm:"foreignKey:LocationID"`
  LocationPath        string                `gorm:"not null"`
  LocationNamePath    string                `gorm:"not null"`
  CreatedAt           time.Time             `gorm:"not null;default:now()"`
  UpdatedAt           time.Time             `gorm:"not null;default:now()"` 
}

// ----------------------------------------------------
// TransactionRecord
// ----------------------------------------------------
type TransactionRecord struct {
  ID                  uuid.UUID         `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
  CompanyID           *uuid.UUID        `gorm:"not null;index"`
  Company             *Company          `gorm:"constraint:OnDelete:CASCADE"`
  WarehouseID         *uuid.UUID        `gorm:"not null;index"`
  Warehouse           *Warehouse        `gorm:"constraint:OnDelete:CASCADE"`
  LocationID          *uuid.UUID        `gorm:"not null;index"`
  Location            *Location         `gorm:"constraint:OnDelete:CASCADE"`
  TransactionFileID   *uuid.UUID        `gorm:"index"`
  TransactionFile     *TransactionFile  `gorm:"constraint:OnDelete:SET NULL"`
  TransactionType     string
  OrderName           string
  ItemID              *uuid.UUID        `gorm:"not null;index"`
  Item                *Item             `gorm:"constraint:OnDelete:CASCADE"`
  Description         string
  TransactionQuantity int
  CompletedDate       time.Time
  CompletedQuantity   int
  CreatedAt           time.Time         `gorm:"not null;default:now()"`
  UpdatedAt           time.Time         `gorm:"not null;default:now()"`

}

// ----------------------------------------------------
// Item
// ----------------------------------------------------
type Item struct {
  ID                 uuid.UUID              `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
  Name               string                 `gorm:"not null"`
  CompanyID          *uuid.UUID             `gorm:"not null;index"`
  Company            *Company               `gorm:"constraint:OnDelete:CASCADE"`
  Warehouses         []*Warehouse           `gorm:"many2many:items_warehouses;"`
  Locations          []*Location            `gorm:"many2many:items_locations;"`
  TransactionRecords []*TransactionRecord   `gorm:"foreignKey:ItemID"`
  TransactionFiles    []*TransactionFile     `gorm:"many2many:items_transaction_files;"`
  CreatedAt          time.Time              `gorm:"not null;default:now()"`
  UpdatedAt          time.Time              `gorm:"not null;default:now()"`
}

// ----------------------------------------------------
// UserAction
// ----------------------------------------------------
// This table captures user actions for auditing or analytics.
type UserAction struct {
  ID          uuid              `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
  UserID      *uuid.UUID        `gorm:"not null;index"`
  User        *User             `gorm:"constraint:OnDelete:CASCADE"` // optional reference to load user details

  ActionType  string            `gorm:"not null"`       // e.g. "CREATE_ITEM", "UPDATE_LOCATION"
  EntityType  string            `gorm:"not null"`       // e.g. "ITEM", "WAREHOUSE", "COMPANY"
  EntityID    *uuid.UUID        `gorm:"index"`          // which record was changed (can be null if none)
  Description string            // short text describing the action
  CreatedAt   time.Time         `gorm:"not null;default:now()"`
  Metadata    datatypes.JSON    `gorm:"type:jsonb"` 
}

type TransactionFile struct {
  ID                  uuid.UUID             `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
  FileName            string                `gorm:"not null"`
  TransactionRecords  []*TransactionRecord  `gorm:"foreignKey:TransactionFileID"`
  Locations           []*Location           `gorm:"many2many:transaction_files_locations;"`
  Items               []*Item               `gorm:"many2many:items_transaction_files;"`
  WarehouseID         *uuid.UUID            `gorm:"not null;index"`
  Warehouse           *Warehouse            `gorm:"constraint:OnDelete:CASCASE"`
  CompanyID           *uuid.UUID            `gorm:"not null;index"`
  Company             *Company              `gorm:"constraint:OnDelete:CASCADE"`
  CreatedAt           time.Time             `gorm:"not null;default:now()"`
  FileExtension       string                `gorm:"column:file_extension"`
  FilePathURL         string                `gorm:"column:file_path_url"`
}

