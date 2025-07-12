package getAuth

import (
	"auth_service/internal/jwt"
	"auth_service/internal/jwt/codec"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func GetAuthMiddleware(jwt jwt.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "missing or invalid Authorization header"})
			c.Abort()
			return
		}

		tokenBase64 := strings.TrimPrefix(authHeader, "Bearer ")

		tokenString, err := codec.Decode(tokenBase64)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid token",
			})
			c.Abort()
			return
		}

		claims, err := jwt.VerifyAccessToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid access token"})
			c.Abort()
			return
		}

		c.Set("email", claims.Subject)

		c.Next()
	}

}
