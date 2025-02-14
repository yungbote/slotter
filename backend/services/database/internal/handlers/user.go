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

type UserHandler struct {
    service services.UserService
}

func NewUserHandler(service services.UserService) *UserHandler {
    return &UserHandler{service: service}
}

func (h *UserHandler) CreateUser(c *gin.Context) {
    var input models.User
    if err := c.ShouldBindJSON(&input); err != nil {
        logger.GetLogger().Warn("CreateUser bind error", zap.Error(err))
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON payload"})
        return
    }
    createdUser, err := h.service.CreateUser(&input)
    if err != nil {
        logger.GetLogger().Error("CreateUser service error", zap.Error(err))
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create user"})
        return
    }
    c.JSON(http.StatusCreated, createdUser)
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
        logger.GetLogger().Error("GetUserByID service error", zap.Error(err))
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
        logger.GetLogger().Error("UpdateUser get error", zap.Error(err))
        c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch user"})
        return
    }
    var input models.User
    if err := c.ShouldBindJSON(&input); err != nil {
        logger.GetLogger().Warn("UpdateUser bind error", zap.Error(err))
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
        return
    }
    updatedUser, err := h.service.UpdateUser(existing, &input)
    if err != nil {
        logger.GetLogger().Error("UpdateUser service error", zap.Error(err))
        c.JSON(http.StatusInternalServerError, gin.H{"error": "could not update user"})
        return
    }
    c.JSON(http.StatusOK, updatedUser)
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
        logger.GetLogger().Error("DeleteUser get error", zap.Error(err))
        c.JSON(http.StatusInternalServerError, gin.H{"error": "could not delete user"})
        return
    }
    if err := h.service.DeleteUser(existing); err != nil {
        logger.GetLogger().Error("DeleteUser service error", zap.Error(err))
        c.JSON(http.StatusInternalServerError, gin.H{"error": "could not delete user"})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "user deleted"})
}


