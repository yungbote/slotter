package handlers

import (
    "net/http"
    "strconv"
    "errors"
    "github.com/gin-gonic/gin"
    "go.uber.org/zap"
    "github.com/yungbote/slotter/backend/services/database/internal/logger"
    "github.com/yungbote/slotter/backend/services/database/internal/models"
    "github.com/yungbote/slotter/backend/services/database/internal/repositories"
    "github.com/yungbote/slotter/backend/services/database/internal/services"
)

type LocationHandler struct {
    service services.LocationService
}

func NewLocationHandler(service services.LocationService) *LocationHandler {
    return &LocationHandler{service: service}
}

func (h *LocationHandler) CreateLocation(c *gin.Context) {
    var input models.Location
    if err := c.ShouldBindJSON(&input); err != nil {
        logger.GetLogger().Warn("CreateLocation bind error:", zap.Error(err))
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON payload"})
        return
    }
    created, err := h.service.CreateLocation(&input)
    if err != nil {
        logger.GetLogger().Error("CreateLocation service error", zap.Error(err))
        if dErr, ok := err.(*repositories.DomainError); ok {
            switch dErr.Code {
            case repositories.ErrCodeValidation:
                c.JSON(http.StatusBadRequest, gin.H{"error": dErr.Message})
                return
            case repositories.ErrCodeForeignKey:
                c.JSON(http.StatusBadRequest, gin.H{"error": dErr.Message})
                return
            default:
                c.JSON(http.StatusInternalServerError, gin.H{"error": dErr.Message})
                return
            }
        }
        c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create location"})
        return
    }
    c.JSON(http.StatusCreated, created)
}

func (h *LocationHandler) GetLocationByID(c *gin.Context) {
    idParam := c.Param("id")
    id, err := strconv.Atoi(idParam)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid location ID"})
        return
    }
    loc, err := h.service.GetLocationByID(uint(id))
    if errors.Is(err, repositories.ErrNotFound) {
        c.JSON(http.StatusNotFound, gin.H{"error": "location not found"})
        return
    }
    if err != nil {
        logger.GetLogger().Error("GetLocationByID service error", zap.Error(err))
        c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch location"})
        return
    }
    c.JSON(http.StatusOK, loc)
}

func (h *LocationHandler) UpdateLocation(c *gin.Context) {
    idParam := c.Param("id")
    id, err := strconv.Atoi(idParam)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid location ID"})
        return
    }
    existing, err := h.service.GetLocationByID(uint(id))
    if errors.Is(err, repositories.ErrNotFound) {
        c.JSON(http.StatusNotFound, gin.H{"error": "location not found"})
        return
    }
    if err != nil {
        logger.GetLogger().Error("UpdateLocation get error", zap.Error(err))
        c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch location"})
        return
    }
    var input models.Location
    if err := c.ShouldBindJSON(&input); err != nil {
        logger.GetLogger().Warn("UpdateLocation bind error", zap.Error(err))
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
        return
    }
    updated, err := h.service.UpdateLocation(existing, &input)
    if err != nil {
        logger.GetLogger().Error("UpdateLocation service error", zap.Error(err))
        if dErr, ok := err.(*repositories.DomainError); ok {
            switch dErr.Code {
            case repositories.ErrCodeValidation:
                c.JSON(http.StatusBadRequest, gin.H{"error": dErr.Message})
                return
            case repositories.ErrCodeForeignKey:
                c.JSON(http.StatusBadRequest, gin.H{"error": dErr.Message})
                return
            default:
                c.JSON(http.StatusInternalServerError, gin.H{"error": dErr.Message})
                return
            }
        }
        c.JSON(StatusInternalServerError, gin.H{"error": "could not update location"})
        return
    }
    c.JSON(http.StatusOK, updated)
}

func (h *LocationHandler) DeleteLocation(c *gin.Context) {
    idParam := c.Param("id")
    id, err := strconv.Atoi(idParam)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid location ID"})
        return
    }
    existing, err := h.service.GetLocationByID(uint(id))
    if errors.Is(err, repositories.ErrNotFound) {
        c.JSON(http.StatusNotFound, gin.H{"error": "location not found"})
        return
    }
    if err != nil {
        logger.GetLogger().Error("DeleteLocation get error", zap.Error(err))
        c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch location"})
        return
    }
    if err := h.service.DeleteLocation(existing); err != nil {
        logger.GetLogger().Error("DeleteLocation service error", zap.Error(err))
        c.JSON(http.StatusInternalServerError, gin.H{"error": "could not delete location"})
        return
    }
    c.JSON(http.StatusOK, gin.H{"error": "location deleted"})
}
