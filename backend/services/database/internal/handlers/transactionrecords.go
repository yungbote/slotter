package api

import (
  "log"
  "net/http"
  "strconv"
  "github.com/gin-gonic/gin"
  "github.com/yungbote/slotter/backend/services/database/internal/database"
  "github.com/yungbote/slotter/backend/services/database/internal/models"
)

func CreateTransactionRecord(c *gin.Context) {
  var input models.TransactionRecord
  if err := c.ShouldBindJSON(&input); err != nil {
    log.Println("ERROR: CreateTransactionRecord bind error:", err)
    c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
    return
  }
  if err := database.DB.Create(&input).Error; err != nil {
    log.Println("ERROR: CreateTransactionRecord DB error:", err)
    c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create transaction record"})
    return
  }
  c.JSON(http.StatusOK, input)
}

func GetAllTransactionRecords(c *gin.Context) {
  var records []models.TransactionRecord
  if err := database.DB.Find(&records).Error; err != nil {
    log.Println("ERROR: GetAllTransactionRecords DB error:", err)
    c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch transaction records"})
    return
  }
  c.JSON(http.StatusOK, records)
}

func GetTransactionRecordsByID(c *gin.Context) {
  idParam := c.Param("id")
  id, err := strconv.Atoi(idParam)
  if err != nil {
    c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid record ID"})
    return
  }
  var record models.TransactionRecord
  if err := database.DB.First(&record, id).Error; err != nil {
    log.Println("ERROR: GetTransactionRecordByID DB error:", err)
    c.JSON(http.StatusNotFound, gin.H{"error": "Record not found"})
    return
  }
  c.JSON(http.StatusOK, record)
}

func UpdateTransactionRecord(c *gin.Context) {
  idParam := c.Param("id")
  id, err := strconv.Atoi(idParam)
  if err != nil {
    c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid record ID"})
    return
  }
  var existing models.TransactionRecord
  if err := database.DB.First(&existing, id).Error; err != nil {
    c.JSON(http.StatusNotFound, gin.H{"error": "Record not found"})
    return
  }
  var input models.TransactionRecord
  if err := c.ShouldBindJSON(&input); err != nil {
    log.Println("ERROR: UpdateTransactionRecord bind error:", err)
    c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
    return
  }
  existing.CompanyID = input.CompanyID
  existing.TransactionType = input.TransactionType
  existing.OrderNumber = input.OrderNumber
  existing.ItemNumber = input.ItemNumber
  existing.Description = input.Description
  existing.TransactionQuantity = input.TransactionQuantity
  existing.Location = input.Location
  existing.Zone = input.Zone
  existing.Carousel = input.Carousel
  existing.Row = input.Row
  existing.Shelf = input.Shelf
  existing.Bin = input.Bin
  existing.CompletedDate = input.CompletedDate
  existing.CompletedBy = input.CompletedBy
  existing.CompletedQuantity = input.CompletedQuantity

  if err := database.DB.Save(&existing).Error; err != nil {
    log.Println("ERROR: UpdateTransactionRecord DB error:", err)
    c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update transaction record"})
    return
  }
  c.JSON(http.StatusOK, existing)
}

func DeleteTransactionRecord(c *gin.Context) {
  idParam := c.Param("id")
  id, err := strconv.Atoi(idParam)
  if err != nil {
    c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid record ID"})
    return
  }
  var record models.TransactionRecord
  if err := database.DB.First(&record, id).Error; err != nil {
    c.JSON(http.StatusNotFound, gin.H{"error": "Record not found"})
    return
  }
  if err := database.DB.Delete(&record).Error; err != nil {
    log.Println("ERROR: DeleteTransactionRecord DB error:", err)
    c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not delete transaction record"})
    return
  }
  c.JSON(http.StatusOK, gin.H{"message": "Transaction record deleted"})
}
