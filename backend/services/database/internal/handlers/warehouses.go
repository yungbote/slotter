package handlers

import (
    "errors"
    "log"
    "net/http"
    "strconv"
    "github.com/gin-gonic/gin"
    "github.com/yungbote/slotter/backend/services/database/internal/models"
    "github.com/yungbote/slotter/backend/services/database/internal/repositories"
    "github.com/yungbote/slotter/backend/services/database/internal/services"
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
        log.Println("CreateWarehouse bind error:", err)
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
        return
    }
    created, err := h.service.CreateWarehouse(&input)
    if err != nil {
        log.Println("CreateWarehouse service error:", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create warehouse"})
        return
    }
    c.JSON(http.StatusCreated, created)
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
        log.Println("GetWarehouseByID service error:", err)
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
        log.Println("UpdateWarehouse get error:", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch warehouse"})
        return
    }
    var input models.Warehouse
    if err := c.ShouldBindJSON(&input); err != nil {
        log.Println("UpdateWarehouse bind error:", err)
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
        return
    }
    updates, err := h.service.UpdateWarehouse(existing, &input)
    if err != nil {
        log.Println("UpdateWarehouse service error:", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "could not update warehouse"})
        return
    }
    c.JSON(http.StatusOK, updated)
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
        log.Println("DeleteWarehouse get error:", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch warehouse"})
        return
    }
    if err := h.service.DeleteWarehouse(existing); err != nil {
        log.Println("DeleteWarehouse service error:", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "could not delete warehouse"})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "warehouse deleted"})
}
