package middleware

import (
    "net/http"
    "strings"

    "github.com/gin-gonic/gin"
    "github.com/google/uuid"

    "github.com/yungbote/slotter/backend/services/database/internal/services"
)

// AuthMiddleware checks for a Bearer <access_token>, validates it,
// and sets user_id in context on success.
func AuthMiddleware(tokenService services.TokenService) gin.HandlerFunc {
    return func(c *gin.Context) {
        header := c.GetHeader("Authorization")
        if header == "" {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing Authorization header"})
            return
        }

        parts := strings.SplitN(header, " ", 2)
        if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization format"})
            return
        }
        accessToken := parts[1]

        // Validate token
        claims, err := tokenService.ValidateAccessToken(accessToken)
        if err != nil {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
            return
        }

        userID, parseErr := uuid.Parse(claims.Subject)
        if parseErr != nil {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid user_id in token"})
            return
        }

        // Store userID for handlers
        c.Set("user_id", userID)
        c.Next()
    }
}

