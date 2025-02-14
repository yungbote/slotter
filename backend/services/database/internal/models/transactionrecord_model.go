package models

import (
  "time"
)

type TransactionRecord struct {
  ID                    uint          `gorm:"primaryKey"`
  WarehouseID           uint          `gorm:"not null;index"`
  Warehouse             Warehouse
  TransactionType       string
  OrderNumber           string
  ItemNumber            string
  Description           string
  TransactionQuantity   int
  Location              string
  Zone                  string
  Carousel              string
  Row                   string
  Shelf                 string
  Bin                   string
  CompletedDate         time.Time
  CompletedBy           string
  CompletedQuantity     int
  CreatedAt             time.Time
  UpdatedAt             time.Time
}
