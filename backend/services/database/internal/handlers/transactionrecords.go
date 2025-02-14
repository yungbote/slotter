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

type TransactionRecordHandler struct {
    service services.TransactionRecordService
}

func NewTransactionRecordHandler(service services.TransactionRecordService) *TransactionRecordHandler {
    return &TransactionRecordHandler{service: service}
}

func (h *TransactionRecordHandler) CreateTransactionRecord(c *gin.Context) {
    var input models.TransactionRecord
    if err := c.ShouldBindJSON(&input); err != nil {
        logger.GetLogger().Warn("CreateTransactionRecord bind error", zap.Error(err))
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
        return
    }
    createdTransactionRecord, err := h.service.CreateTransactionRecord(&input)
    if err != nil {
        logger.GetLogger().Warn("CreateTransactionRecord service error", zap.Error(err))
        c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create transaction record"})
        return
    }
    c.JSON(http.StatusCreated, createdTransactionRecord)
}

func (h *TransactionRecordHandler) GetTransactionRecordByID(c *gin.Context) {
    idParam := c.Param("id")
    id, err := strconv.Atoi(idParam)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid transaction record ID"})
        return
    }
    record, err := h.service.GetTransactionRecordByID(uint(id))
    if errors.Is(err, repositories.ErrNotFound) {
        c.JSON(http.StatusNotFound, gin.H{"error": "transaction record not found"})
        return
    }
    if err != nil {
        logger.GetLogger().Error("GetTransactionRecordByID service error", zap.Error(err))
        c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch transaction record"})
        return
    }
    c.JSON(http.StatusOK, record)
}

func (h *TransactionRecordHandler) UpdateTransactionRecord(c *gin.Context) {
    idParam := c.Param("id")
    id, err := strconv.Atoi(idParam)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid transaction record ID"})
        return
    }
    existing, err := h.service.GetTransactionRecordByID(uint(id))
    if errors.Is(err, repositories.ErrNotFound) {
        c.JSON(http.StatusNotFound, gin.H{"error": "transaction record not found"})
        return
    }
    if err != nil {
        logger.GetLogger().Error("UpdateTransactionRecord get error", zap.Error(err))
        c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch transaction record"})
        return
    }
    var input models.TransactionRecord
    if err := c.ShouldBindJSON(&input); err != nil {
        logger.GetLogger().Warn("UpdateTransactionRecord bind error", zap.Error(err))
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
        return
    }
    updatedTransactionRecord, err := h.service.UpdateTransactionRecord(existing, &input)
    if err != nil {
        logger.GetLogger().Error("UpdateTransactionRecord service error", zap.Error(err))
        c.JSON(http.StatusInternalServerError, gin.H{"error": "could not update transaction record"})
        return
    }
    c.JSON(http.StatusOK, updatedTransactionRecord)
}

func (h *TransactionRecordHandler) DeleteTransactionRecord(c *gin.Context) {
    idParam := c.Param("id")
    id, err := strconv.Atoi(idParam)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid transaction record ID"})
        return
    }
    existing, err := h.service.GetTransactionRecordByID(uint(id))
    if errors.Is(err, repositories.ErrNotFound) {
        c.JSON(http.StatusNotFound, gin.H{"error": "transaction record not found"})
        return
    }
    if err != nil {
        logger.GetLogger().Error("DeleteTransactionRecord get error", zap.Error(err))
        c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch transaction record"})
        return
    }
    if err := h.service.DeleteTransactionRecord(existing); err != nil {
        logger.GetLogger().Error("DeleteTransactionRecord service error", zap.Error(err))
        c.JSON(http.StatusInternalServerError, gin.H{"error": "could not delete transaction record"})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "transaction record deleted"})
}
