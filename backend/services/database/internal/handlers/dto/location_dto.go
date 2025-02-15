package dto

import (
    "github.com/yungbote/slotter/backend/services/database/internal/models"
)

//Input DTO's
type CreateLocationRequest struct {
    WarehouseID         uint        `json:"warehouse_id"    binding:"required"`
    ParentLocationID    *uint       `json:"parent_location_id,omitempty"`
    LocationName        string      `json:"location_name"   binding:"required"`
    LocationType        string      `json:"location_type"   binding:"required"`
}

type UpdateLocationRequest struct {
    ParentLocationID     *uint       `json:"parent_location_id,omitempty"`
    LocationName         *string     `json:"location_name,omitempty"`
    LocationType         *string     `json:"location_type,omitempty"`
}

type  GetLocationByIDRequest struct {
    ID      uint        `json:"id" binding:"required"`
}

type DeleteLocationRequest struct {
    ID      uint        `json:"id" binding:"required"`
}

type ListLocationsByWarehouseIDRequest struct {
    WarehouseID     uint    `json:"warehouse_id" binding:"required"`
}

type ListLocationChilrenByParentIDRequest struct {
    ParentLocationID    uint    `json:"parent_location_id" binding:"required"`
}

type CountLocationsByWarehouseIDRequest struct {
    WarehouseID     uint    `json:"warehouse_id" binding:"required"`
}

//Output DTO's
type LocationResponse struct {
    ID               uint        `json:"id"`
    WarehouseID      uint        `json:"warehouse_id"`
    ParentLocationID *uint       `json:"parent_location_id,omitempty"`
    LocationName     string      `json:"location_name"`
    LocationType     string      `json:"location_type"`
}

//Mapping Functions
func (req *CreateLocationRequest) ToModel() *models.Location {
    return &models.Location{
        WarehouseID:        req.WarehouseID,
        ParentLocationID:   req.ParentLocationID,
        LocationName:       req.LocationName,
        LocationType:       req.LocationType,
    }
}

func (req *UpdateLocationRequest) ApplyToModel(loc *models.Location) {
    if req.ParentLocationID != nil {
        loc.ParentLocationID = req.ParentLocationID
    }
    if req.LocationName != nil {
        loc.LocationName = *req.LocationName
    }
    if req.LocationType != nil {
        loc.LocationType = *req.LocationType
    }
}

func ToLocationResponse(loc *models.Location) *LocationResponse {
    if loc == nil {
        return nil
    }
    return &LocationResponse{
        ID:                 loc.ID,
        WarehouseID:        loc.WarehouseID,
        ParentLocationID:   loc.ParentLocationID,
        LocationName:       loc.LocationName,
        LocationType:       loc.LocationType,
    }
}
