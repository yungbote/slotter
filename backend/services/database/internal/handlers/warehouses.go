package handlers

import (
    "net/http"
    "strconv"
    "github.com/gin-gonic/gin"
    "errors"
    "go.uber.org/zap"
    "github.com/yungbote/slotter/backend/services/database/internal/models"
    "github.com/yungbote/slotter/backend/services/database/internal/repositories"
    "github.com/yungbote/slotter/backend/services/database/internal/services"
    "github.com/yungbote/slotter/backend/services/database/internal/logger"
)

type WarehouseHandler struct {
    service services.WarehouseService
}

func NewWarehouseHandler(service services.WarehouseService) *WarehouseHandler {
    return &WarehouseHandler{service: service}
}

func (h *WarehouseHandler) CreateWarehouse(c *gin.Context) {
    var input models.Warehouse
    if err := c.ShouldBindJSON(&input); err != nil {
        logger.GetLogger().Warn("CreateWarehouse bind error", zap.Error(err))
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
        return
    }
    createdWarehouse, err := h.service.CreateWarehouse(&input)
    if err != nil {
        logger.GetLogger().Error("CreateWarehouse service error", zap.Error(err))
        c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create warehouse"})
        return
    }
    c.JSON(http.StatusCreated, createdWarehouse)
}

func (h *WarehouseHandler) GetWarehouseByID(c *gin.Context) {
    idParam := c.Param("id")
    id, err := strconv.Atoi(idParam)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid warehouse ID"})
        return
    }
    warehouse, err := h.service.GetWarehouseByID(uint(id))
    if errors.Is(err, repositories.ErrNotFound) {
        c.JSON(http.StatusNotFound, gin.H{"error": "warehouse not found"})
        return
    }
    if err != nil {
        logger.GetLogger().Error("GetWarehouseByID service error", zap.Error(err))
        c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch warehouse"})
        return
    }
    c.JSON(http.StatusOK, warehouse)
}

func (h *WarehouseHandler) UpdateWarehouse(c *gin.Context) {
    idParam := c.Param("id")
    id, err := strconv.Atoi(idParam)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid warehouse ID"})
        return
    }
    existing, err := h.service.GetWarehouseByID(uint(id))
    if errors.Is(err, repositories.ErrNotFound) {
        c.JSON(http.StatusNotFound, gin.H{"error": "warehouse not found"})
        return
    }
    if err != nil {
        logger.GetLogger().Error("UpdateWarehouse get error", zap.Error(err))
        c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch warehouse"})
        return
    }
    var input models.Warehouse
    if err := c.ShouldBindJSON(&input); err != nil {
        logger.GetLogger().Warn("UpdateWarehouse bind error", zap.Error(err))
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
        return
    }
    updatedWarehouse, err := h.service.UpdateWarehouse(existing, &input)
    if err != nil {
        logger.GetLogger().Error("UpdateUser service error", zap.Error(err))
        c.JSON(http.StatusInternalServerError, gin.H{"error": "could not update warehouse"})
        return
    }
    c.JSON(http.StatusOK, updatedWarehouse)
}

func (h *WarehouseHandler) DeleteWarehouse(c *gin.Context) {
    idParam := c.Param("id")
    id, err := strconv.Atoi(idParam)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid warehouse ID"})
        return
    }
    existing, err := h.service.GetWarehouseByID(uint(id))
    if errors.Is(err, repositories.ErrNotFound) {
        c.JSON(http.StatusNotFound, gin.H{"error": "warehouse not found"})
        return
    }
    if err != nil {
        logger.GetLogger().Error("DeleteWarehouse get error", zap.Error(err))
        c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch warehouse"})
        return
    }
    if err := h.service.DeleteWarehouse(existing); err != nil {
        logger.GetLogger().Error("DeleteWarehouse service error", zap.Error(err))
        c.JSON(http.StatusInternalServerError, gin.H{"error": "could not delete warehouse"})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "warehouse deleted"})
}
