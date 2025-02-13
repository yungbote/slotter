package api

import (
  "log"
  "net/http"
  "strconv"
  "github.com/gin-gonic/gin"
  "github.com/yungbote/slotter/backend/services/database/internal/database"
  "github.com/yungbote/slotter/backend/services/database/internal/models"
)

func CreateCompany(c *gin.Context) {
  var input models.Company
  if err := c.ShouldBindJSON(&input); err != nil {
    log.Println("ERROR: CreateCompany bind error:", err)
    c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
    return
  }
  if err := database.DB.Create(&input).Error; err != nil {
    log.Println("ERROR: CreateCompany DB error:", err)
    c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create company"})
    return
  }
  c.JSON(http.StatusCreated, input)
}

func GetAllCompanies(c *gin.Context) {
  var companies []models.Company
  if err := database.DB.Find(&companies).Error; err != nil {
    log.Println("ERROR: GetAllCompanies DB error:", err)
    c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch companies"})
    return
  }
  c.JSON(http.StatusOK, companies)
}

func GetCompanyByID(c *gin.Context) {
  idParam := c.Param("id")
  id, err := strconv.Atoi(idParam)
  if err != nil {
    c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid company ID"})
    return
  }
  var company models.Company
  if err := database.DB.First(&company, id).Error; err != nil {
    log.Println("ERROR: GetCompanyByID DB error:", err)
    c.JSON(http.StatusNotFound, gin.H{"error": "Company not found"})
    return
  }
  c.JSON(http.StatusOK, company)
}

func UpdateCompany(c *gin.Context) {
  idParam := c.Param("id")
  id, err := strconv.Atoi(idParam)
  if err != nil {
    c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid company ID"})
    return
  }
  var existing models.Company
  if err := database.DB.First(&existing, id).Error; err != nil {
    c.JSON(http.StatusNotFound, gin.H{"error": "Company not found"})
    return
  }
  var input models.Company
  if err := c.ShouldBindJSON(&input); err != nil {
    log.Println("ERROR: UpdateCompany bind error:", err)
    c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
    return
  }
  existing.Name = input.Name
  if err := database.DB.Save(&existing).Error; err != nil {
    log.Println("ERROR: UpdateCompany DB error:", err)
    c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update company"})
    return
  }
  c.JSON(http.StatusOK, existing)
}

func DeleteCompany(c *gin.Context) {
  idParam := c.Param("id")
  id, err := strconv.Atoi(idParam)
  if err != nil {
    c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid company ID"})
    return
  }
  var company models.Company
  if err := database.DB.First(&company, id).Error; err != nil {
    c.JSON(http.StatusNotFound, gin.H{"error": "Company not found"})
    return
  }
  if err := database.DB.Delete(&company).Error; err != nil {
    log.Println("ERROR: DeleteCompany DB error:", err)
    c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not delete company"})
    return
  }
  c.JSON(http.StatusOK, gin.H{"message": "Company deleted"})
}
