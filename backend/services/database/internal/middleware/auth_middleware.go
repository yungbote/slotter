package middleware

import (
  "net/http"
  "strings"
  "github.com/gin-gonic/gin"
  "github.com/golang-jwt/jwt/v4"
  "github.com/google/uuid"
  authsvc "github.com/yungbote/slotter/backend/services/database/services/auth"
)

func AuthMiddleware(tokenService authsvc.TokenService) gin.HandlerFunc {
  return func(c *gin.Context) {
    if header == "" {
      c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing authorization header"})
      return
    }
    parts := strings.SplitN(header, " ", 2)
    if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
      c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization token"})
      return
    }
    tokenStr := parts[1]
    claims, err := tokenService.ValidateToken(tokenStr)
    if err != nil {
      c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
      return
    }
    userID, err := uuid.Parse(claims.Subject)
    if err != nil {
      c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid user uuid"})
      return
    }
    c.Set("user_id", userID)
    c.Next()
  }
}
