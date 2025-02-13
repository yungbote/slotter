package api

import (
  "log"
  "net/http"
  "strconv"
  "github.com/gin-gonic/gin"
  "github.com/yungbote/slotter/backend/services/database/internal/database"
  "github.com/yungbote/slotter/backend/services/database/internal/models"
)

func CreateWarehouse(c *gin.Context) {
  var input models.Warehouse
  if err := c.ShouldBindJSON(&input); err != nil {
    log.Println("ERROR: CreateWarehouse bind error:", err)
    c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
    return
  }
  if err := database.DB.Create(&input).Error; err != nil {
    log.Println("ERROR: CreateWarehouse DB error:", err)
    c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create warehouse"})
    return
  }
  c.JSON(http.StatusOK, input)
}

func GetAllWarehouses(c *gin.Context) {
  var warehouses []models.Warehouse
  if err := database.DB.Find(&warehouses).Error; err != nil {
    log.Println("ERROR: GetAllWarehouses DB error:", err)
    c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch warehouses"})
    return
  }
  c.JSON(http.StatusOK, warehouses)
}

func GetWarehouseByID(c *gin.Context) {
  idParam := c.Param("id")
  id, err := strconv.Atoi(idParam)
  if err != nil {
    c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid warehouse ID"})
    return
  }
  var warehouse models.Warehouse
  if err := database.DB.First(&warehouse, id).Error; err != nil {
    log.Println("ERROR: GetWarehouseByID DB error:", err)
    c.JSON(http.StatusNotFound, gin.H{"error": "Warehouse not found"})
    return
  }
  c.JSON(http.StatusOK, warehouse)
}

func UpdateWarehouse(c *gin.Context) {
  idParam := c.Param("id")
  id, err := strconv.Atoi(idParam)
  if err != nil {
    c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid warehouse ID"})
    return
  }
  var existing models.Warehouse
  if err := database.DB.First(&existing, id).Error; err != nil {
    c.JSON(http.StatusNotFound, gin.H{"error": "Warehouse not found"})
    return
  }
  var input models.Warehouse
  if err := c.ShouldBindJSON(&input); err != nil {
    log.Println("ERROR: UpdateWarehouse bind error:", err)
    c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
    return
  }
  existing.Name = input.Name
  existing.CompanyID = input.CompanyID

  if err := database.DB.Save(&existing).Error; err != nil {
    log.Println("ERROR: UpdateWarehouse DB error:", err)
    c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update warehouse"})
    return
  }
  c.JSON(http.StatusOK, existing)
}

func DeleteWarehouse(c *gin.Context) {
  idParam := c.Param("id")
  id, err := strconv.Atoi(idParam)
  if err != nil {
    c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid warehouse ID"})
    return
  }
  var warehouse models.Warehouse
  if err := database.DB.First(&warehouse, id).Error; err != nil {
    c.JSON(http.StatusNotFound, gin.H{"error": "Warehouse not found"})
    return
  }
  if err := database.DB.Delete(&warehouse).Error; err != nil {
    log.Println("ERROR: DeleteWarehouse DB error:", err)
    c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not delete warehouse"})
    return
  }
  c.JSON(http.StatusOK, gin.H{"message": "Warehouse deleted"})
}
