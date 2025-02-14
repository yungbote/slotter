package api

import (
  "os"
  "net/http"
  "github.com/gin-gonic/gin"
  "go.uber.org/zap"
  "github.com/yungbote/slotter/backend/services/auth/internal/logger"
  "github.com/yungbote/slotter/backend/services/auth/internal/handlers"
  "github.com/yungbote/slotter/backend/services/auth/internal/services"
  "github.com/yungbote/slotter/backend/services/auth/internal/clients"
)

func NewRouter() *gin.Engine {
  r := gin.Default()
  log := logger.GetLogger()

  r.GET("/health", func(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{"status": "OK"})
  })

  dbServiceURL := os.Getenv("DB_SERVICE_URL")
  if dbServiceURL == "" {
    dbServiceURL = "http://slotter-database:8080"
    log.Warn("DB_SERVICE_URL not set, using default", zap.String("url", dbServiceURL))
  }
  dbClient := clients.NewDatabaseClient(dbServiceURL)
  authService := services.NewAuthService(dbClient)
  authHandler := handlers.NewAuthHandler(authService)

  v1 := r.Group("/v1")
  {
    //Signup
    v1.POST("/signup", authHandler.SignUp)
    //Login
    v1.POST("/login", authHandler.Login)
  }
  return r
}
