package dto

import (
    "github.com/yungbote/slotter/backend/services/database/internal/models"
)

//Input DTO's
type CreateTransactionRecordRequest struct {
    WarehouseID         uint            `json:"warehouse_id" binding:"required"`
    LocationID          uint            `json:"location_id" binding:"required"`
    TransactionType     string          `json:"transaction_type" binding:"required"`
    OrderNumber         string          `json:"order_number" binding:"required"`
    ItemNumber          string          `json:"item_number" binding:"required"`
    Description         string          `json:"description,omitempty"`
    TransactionQuantity int             `json:"transaction_quantity" binding:"required"`
    CompletedDate       string          `json:"completed_date" binding:"required"`
    CompletedQuantity   int             `json:"completed_quantity" binding:"required"`
}

type UpdateTransactionRecordRequest struct {
    WarehouseID         *uint           `json:"warehouse_id,omitempty"`
    LocationID          *uint           `json:"location_id,omitempty"`
    TransactionType     *string         `json:"transaction_type,omitempty"`
    OrderNumber         *string         `json:"order_number,omitempty"`
    ItemNumber          *string         `json:"item_number,omitempty"`
    Description         *string         `json:"description,omitempty"`
    TransactionQuantity *int            `json:"transaction_quantity,omitempty"`
    CompletedDate       *string         `json:"completed_date,omitempty"`
    CompletedQuantity   *int            `json:"completed_quantity,omitempty"`
}

type GetTransactionRecordByIDRequest struct {
    ID          uint            `json:"id" binding:"required"`
}

type DeleteTransactionRecordRequest struct {
    ID          uint            `json:"id" binding:"required"`
}

type ListTransactionRecordsByWarehouseIDRequest struct {
    WarehouseID         uint            `json:"warehouse_id" binding:"required"`
}

type ListTransactionRecordsByLocationIDRequest struct {
    LocationID          uint            `json:"location_id" binding:"required"`
}

//Output DTO's
type TransactionRecordResponse struct {
    ID                      uint        `json:"id"`
    WarehouseID             uint        `json:"warehouse_id"`
    LocationID              uint        `json:"location_id"`
    TransactionType         string      `json:"transaction_type,omitempty"`
    OrderNumber             string      `json:"order_number,omitempty"`
    ItemNumber              string      `json:"item_number,omitempty"`
    Description             string      `json:"description,omitempty"`
    TransactionQuantity     int         `json:"transaction_quantity,omitempty"`
    CompletedDate           string      `json:"completed_date,omitempty"`
    CompletedQuantity       int         `json:"completed_quantity,omitempty"`
}

//Mapping Functions
func (req *CreateTransactionRecordRequest) ToModel() *models.TransactionRecord {
    return &models.TransactionRecord{
        WarehouseID:            req.WarehouseID,
        LocationID:             req.LocationID,
        TransactionType:        req.TransactionType,
        OrderNumber:            req.OrderNumber,
        ItemNumber:             req.ItemNumber,
        Description:            req.Description,
        TransactionQuantity:    req.TransactionQuantity,
        CompletedDate:          req.CompletedDate,
        CompletedQuantity:      req.CompletedQuantity,
    }
}

func (req *UpdateTransactionRecordRequest) ApplyToModel(tr *models.TransactionRecord) {
    if req.WarehouseID != nil {
        tr.WarehouseID = *req.WarehouseID
    }
    if req.LocationID != nil {
        tr.LocationID = *req.LocationID
    }
    if req.TransactionType != nil {
        tr.TransactionType = *req.TransactionType
    }
    if req.OrderNumber != nil {
        tr.OrderNumber = *req.OrderNumber
    }
    if req.ItemNumber != nil {
        tr.ItemNumber = *req.ItemNumber
    }
    if req.TransactionQuantity != nil {
        tr.TransactionQuantity = *req.TransactionQuantity
    }
    if req.CompletedDate != nil {
        tr.CompletedDate = *req.CompletedDate
    }
    if req.CompletedQuantity != nil {
        tr.CompletedQuantity = *req.CompletedQuantity
    }
}

func ToTransactionRecordResponse(tr *models.TransactionRecord) *TransactionRecordResponse {
    if tr == nil {
        return nil
    }
    return &TransactionRecordResponse{
        ID:                     tr.ID,
        WarehouseID:            tr.WarehouseID,
        LocationID:             tr.LocationID,
        TransactionType:        tr.TransactionType,
        OrderNumber:            tr.OrderNumber,
        ItemNumber:             tr.ItemNumber,
        Description:            tr.Description,
        TransactionQuantity:    tr.TransactionQuantity,
        CompletedDate:          tr.CompletedDate,
        CompletedQuantity:      tr.CompletedQuantity,
    }
} 
