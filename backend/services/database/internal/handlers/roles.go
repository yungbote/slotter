package handlers

import (
  "log"
  "net/http"
  "strconv"
  "github.com/gin-gonic/gin"
  "github.com/yungbote/slotter/backend/services/database/internal/database"
  "github.com/yungbote/slotter/backend/services/database/internal/models"
)

func CreateRole(c *gin.Context) {
  var input models.Role
  if err := c.ShouldBindJSON(&input); err != nil {
    log.Println("ERROR: CreateRole bind error:", err)
    c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
    return
  }
  if err := database.DB.Create(&input).Error; err != nil {
    log.Println("ERROR: CreateRoleDB error:", err)
    c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create role"})
    return
  }
  c.JSON(http.StatusOK, input)
}

func GetAllRoles(c *gin.Context) {
  var roles []models.Role
  if err := database.DB.Find(&roles).Error; err != nil {
    log.Println("ERROR: GetAllRoles DB error:", err)
    c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch roles"})
    return
  }
  c.JSON(http.StatusOK, roles)
}

func GetRoleByID(c *gin.Context) {
  idParam := c.Param("id")
  id, err := strconv.Atoi(idParam)
  if err != nil {
    c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role ID"})
    return
  }
  var role models.Role
  if err := database.DB.First(&role, id).Error; err != nil {
    log.Println("ERROR: GetRoleByID DB error:", err)
    c.JSON(http.StatusNotFound, gin.H{"error": "Role not found"})
    return
  }
  c.JSON(http.StatusOK, role)
}

func UpdateRole(c *gin.Context) {
  idParam := c.Param("id")
  id, err := strconv.Atoi(idParam)
  if err != nil {
    c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role ID"})
    return
  }
  var existing models.Role
  if err := database.DB.First(&existing, id).Error; err != nil {
    c.JSON(http.StatusNotFound, gin.H{"error": "Role not found"})
    return
  }
  var input models.Role
  if err := c.ShouldBindJSON(&input); err != nil {
    log.Println("ERROR: UpdateRole bind error:", err)
    c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
    return
  }
  existing.ID = input.ID
  existing.Name = input.Name
  existing.Permissions = input.Permissions
  existing.CompanyID = input.CompanyID
  existing.Company = input.Company

  if err := database.DB.Save(&existing).Error; err != nil {
    log.Println("ERROR: UpdateRole DB error:", err)
    c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update role"})
    return
  }
  c.JSON(http.StatusOK, existing)
}

func DeleteRole(c *gin.Context) {
  idParam := c.Param("id")
  id, err := strconv.Atoi(idParam)
  if err != nil {
    c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role ID"})
    return
  }
  var role models.Role
  if err := database.DB.Delete(&role).Error; err != nil {
    log.Println("ERROR: DeleteRole DB error:", err)
    c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not delete role"})
    return
  }
  c.JSON(http.StatusOK, gin.H{"message": "Role deleted"})
}
