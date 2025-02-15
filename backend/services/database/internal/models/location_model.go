package models

import (
  "time"
)

type Location struct {
  ID                uint        `gorm:"primaryKey"`
  WarehouseID       uint        `gorm:"not null;index"`
  Warehouse         Warehouse   `gorm:"constraint:OnDelete:CASCADE"`
  ParentLocationID  *uint
  ParentLocation    *Location   `gorm:"foreignKey:ParentLocationID"`
  LocationName      string      `gorm:"not null"`
  LocationType      string      `gorm:"not null"`
  CreatedAt         time.Time
  UpdatedAt         time.Time
}
