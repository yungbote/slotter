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

type CompanyHandler struct {
    service services.CompanyService
}

func NewCompanyHandler(service service.CompanyService) *CompanyHandler {
    return &CompanyHandler{service: service}
}

func (h *CompanyHandler) CreateCompany(c *gin.Context) {
    var input models.Company
    if err := c.ShouldBindJSON(&input); err != nil {
        log.Println("CreateCompany bind error:", err)
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
        return
    }
    created, err := h.service.CreateCompany(&input)
    if err != nil {
        log.Prinln("CreateCompany service error:", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusCreated, created)
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
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve company"})
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
        log.Println("UpdateCompany get error:", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch company"})
        return
    }
    updated, err := h.service.UpdateCompany(existing, &input)
    if err != nil {
        log.Println("UpdateCompany service error:", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "could not update company"})
        return
    }
    c.JSON(http.StatusOK, updated)
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
        log.Println("DeleteCompany get error:", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "could not delete company"})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "company deleted"})
}
