package api

import (
  "log"
  "net/http"
  "strconv"
  "github.com/gin-gonic/gin"
  "github.com/yungbote/slotter/backend/services/database/internal/database"
  "github.com/yungbote/slotter/backend/services/database/internal/models"
)

func CreateUser(c *gin.Context) {
  var input models.User
  if err := c.ShouldBindJSON(&input); err != nil {
    log.Println("ERROR: CreateUser bind error:", err)
    c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
    return
  }
  if err := database.DB.Create(&input).Error; err != nil {
    log.Println("ERROR: CreateUser DB error:", err)
    c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create user"})
    return
  }
  c.JSON(http.StatusCreated, input)
}

func GetAllUsers(c *gin.Context) {
  var users []models.User
  if err := database.DB.Find(&users).Error; err != nil {
    log.Println("ERROR: GetAllUsers DB error:", err)
    c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch users"})
    return
  }
  c.JSON(http.StatusOK, users)
}

func GetUserByID(c *gin.Context) {
  idParam := c.Param("id")
  id, err := strconv.Atoi(idParam)
  if err != nil {
    c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
    return
  }
  var user models.User
  if err := database.DB.First(&user, id).Error; err != nil {
    log.Println("ERROR: GetUserByID DB error:", err)
    c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
    return
  }
  c.JSON(http.StatusOK, user)
}

func UpdateUser(c *gin.Context) {
  idParam := c.Param("id")
  id, err := strconv.Atoi(idParam)
  if err != nil {
    c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
    return
  }
  var existing models.User
  if err := database.DB.First(&existing, id).Error; err != nil {
    c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
    return
  }
  var input models.User
  if err := c.ShouldBindJSON(&input); err != nil {
    log.Println("ERROR: UpdateUser bind error:", err)
    c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
    return
  }
  existing.Email = input.Email
  existing.PasswordHash = input.PasswordHash
  existing.FullName = input.FullName
  existing.Role = input.Role
  existing.Status = input.Status
  existing.CompanyID = input.CompanyID

  if err := database.DB.Save(&existing).Error; err != nil {
    log.Println("ERROR: UpdateUser DB error:", err)
    c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update user"})
    return
  }
  c.JSON(http.StatusOK, existing)
}

func DeleteUser(c *gin.Context) {
  idParam := c.Param("id")
  id, err := strconv.Atoi(idParam)
  if err != nil {
    c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
    return
  }
  var user models.User
  if err := database.DB.First(&user, id).Error; err != nil {
    c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
    return
  }
  if err := database.DB.Delete(&user).Error; err != nil {
    log.Println("ERROR: DeleteUser DB error:", err)
    c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not delete user"})
    return
  }
  c.JSON(http.StatusOK, gin.H{"message": "User deleted"})
}
