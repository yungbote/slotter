package handlers

import (
    "net/http"
    "strconv"
    "github.com/gin-gonic/gin"
    "errors"
    "github.com/yungbote/slotter/backend/services/database/internal/models"
    "github.com/yungbote/slotter/backend/services/database/internal/repositories"
    "github.com/yungbote/slotter/backend/services/database/internal/services"
    "github.com/yungbote/slotter/backend/services/database/internal/logger"
)

type CompanyHandler struct {
    service services.CompanyService
}

func NewCompanyHandler(service service.CompanyService) *CompanyHandler {
    return &CompanyHandler{service: service}
}

func (h *CompanyHandler) CreateCompany(c *gin.Context) {
    var input models.Company
    if err := c.ShouldBindJSON(&input); err != nil {
        logger.GetLogger().Warn("CreateCompany bind error", zap.Error(err))
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON payload"})
        return
    }
    createdCompany, err := h.service.CreateCompany(&input)
    if err != nil {
        logger.GetLogger().Error("CreateCompany service error", zap.Error(err))
        c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create company"})
        return
    }
    c.JSON(http.StatusCreated, createdCompany)
}

func (h *CompanyHandler) GetCompanyByID(c *gin.Context) {
    idParam := c.Param("id")
    id, err := strconv.Atoi(idParam)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid company ID"})
        return
    }
    company, err := h.service.GetCompanyByID(uint(id))
    if errors.Is(err, repositories.ErrNotFound) {
        c.JSON(http.StatusNotFound, gin.H{"error": "company not found"})
        return
    }
    if err != nil {
        logger.GetLogger().Error("GetCompanyByID service error", zap.Error(err))
        c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch company"})
        return
    }
    c.JSON(http.StatusOK, company)
}

func (h *CompanyHandler) UpdateCompany(c *gin.Context) {
    idParam := c.Param("id")
    id, err := strconv.Atoi(idParam)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid company ID"})
        return
    }
    existing, err := h.service.GetCompanyByID(uint(id))
    if errors.Is(err, repositories.ErrNotFound) {
        c.JSON(http.StatusNotFound, gin.H{"error": "company not found"})
        return
    }
    if err != nil {
        logger.GetLogger().Error("UpdateCompany get error", zap.Error(err))
        c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch company"})
        return
    }
    var input models.Company
    if err := c.ShouldBindJSON(&input); err != nil {
        logger.GetLogger().Warn("UpdateUser bind error", zap.Error(err))
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
        return
    }
    updatedCompany, err := h.service.UpdateCompany(existing, &input)
    if err != nil {
        logger.GetLogger().Error("UpdateCompany service error", zap.Error(err))
        c.JSON(http.StatusInternalServerError, gin.H{"error": "could not update company"})
        return
    }
    c.JSON(http.StatusOK, updatedCompany)
}

func (h *CompanyHandler) DeleteCompany(c *gin.Context) {
    idParam := c.Param("id")
    id, err := strconv.Atoi(idParam)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid company ID"})
        return
    }
    existing, err := h.service.GetCompanyByID(uint(id))
    if errors.Is(err, repositories.ErrNotFound) {
        c.JSON(http.StatusNotFound, gin.H{"error": "company not found"})
        return
    }
    if err != nil {
        logger.GetLogger().Error("DeleteCompany get error", zap.Error(err))
        c.JSON(http.StatusInternalServerError, gin.H{"error": "could not delete company"})
        return
    }
    if err := h.service.DeleteCompany(existing); err != nil {
        logger.GetLogger().Error("DeleteCompany service error", zap.Error(err))
        c.JSON(http.StatusInternalServerError, gin.H{"error": "could not delete company"})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "company deleted"})
}
