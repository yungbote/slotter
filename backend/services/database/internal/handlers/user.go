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

type UserHandler struct {
    service services.UserService
}

func NewUserHandler(service services.UserService) *UserHandler {
    return &UserHandler{service: service}
}

func (h *UserHandler) CreateUser(c *gin.Context) {
    var input models.User
    if err := c.ShouldBindJSON(&input); err != nil {
        log.Println("CreateUser bind error:", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create user"})
        return
    }
    c.JSON(http.StatusCreated, created)
}

func (h *UserHandler) GetUserByID(c *gin.Context) {
    idParam := c.Param("id")
    id, err := strconv.Atoi(idParam)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
        return
    }
    user, err := h.service.GetUserByID(uint(id))
    if errors.Is(err, repositories.ErrNotFound) {
        c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
        return
    }
    if err != nil {
        log.Println("GetUserByID error:", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch user"})
        return
    }
    c.JSON(http.StatusOK, user)
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
    idParam := c.Param("id")
    id, err := strconv.Atoi(idParam)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
        return
    }
    existing, err := h.service.GetUserByID(uint(id))
    if errors.Is(err, repositories.ErrNotFound) {
        c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
        return
    }
    if err != nil {
        log.Println("UpdateUser get error:", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch user"})
        return
    }
    var input models.User
    if err := c.ShouldBindJSON(&input); err != nil {
        log.Println("UpdateUser bind error:", err)
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
        return
    }
    updates, err := h.service.UpdateUser(existing, &input)
    if err != nil {
        log.Prinln("UpdateUser service error:", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "could not update user"})
        return
    }
    c.JSON(http.StatusOK, updated)
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
    idParam := c.Param("id")
    id, err := strconv.Atoi(idParam)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
        return
    }
    existing, err := h.service.GetUserByID(uint(id))
    if errors.Is(err, repositories.ErrNotFound) {
        c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
        return
    }
    if err != nil {
        log.Prinln("DeleteUser get error:", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "could not delete user"})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "user deleted"})
}


