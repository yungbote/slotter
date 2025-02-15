package models

type Warehouse struct {
  ID          uint      `gorm:"primaryKey"`
  Name        string    `gorm:"not null"`
  CompanyID   uint
  Company     Company
}
