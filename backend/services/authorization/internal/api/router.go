package api

import (
  "github.com/gin-gonic/gin"
  "github.com/yungbote/slotter/backend/services/authorization/internal/handlers"
)

func NewRouter() *gin.Engine {
  r := gin.Default()

  r.GET("/health", func(c *gin.Context) {
    c.JSON(200, gin.H{"status": "OK"})
  })

  v1 := r.Group("/v1")
  {
    v1.POST("/signup", handlers.SignUp)
    v1.POST("/login", handlers.Login)
    v1.POST("/refresh", handlers.RefreshToken)
  

    protected := v1.Group("/")
    protected.Use(handlers.JWTAuthMiddleware)
    {
      protected.GET("protected", handlers.ProtectedEndpoint)
    }
  }
  return r
}
