package models

import (
  "time"
)

type TransactionRecord struct {
  ID                    uint          `gorm:"primaryKey"`
  WarehouseID           uint          `gorm:"not null;index"`
  Warehouse             Warehouse     `gorm:"constraint:OnDelete:CASCADE"`
  LocationID            uint          `gorm:"not null;index"`
  Location              Location      `gorm:"constraint:OnDelete:CASCADE"`
  TransactionType       string
  OrderNumber           string
  ItemNumber            string
  Description           string
  TransactionQuantity   int
  CompletedDate         time.Time
  CompletedQuantity     int
  CreatedAt             time.Time
  UpdatedAt             time.Time
}
