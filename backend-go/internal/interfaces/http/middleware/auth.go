package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/nakachan-ing/fukuworks/backend-go/internal/utils"
)

// Reserved words (e.g., for static routes like /login, /admin, etc.)
var reservedPaths = map[string]bool{
	"login":  true,
	"admin":  true,
	"health": true,
}

// Middleware to skip reserved paths from user route handling
func ReservedPathGuard() gin.HandlerFunc {
	return func(c *gin.Context) {
		first := strings.Split(strings.TrimLeft(c.Request.URL.Path, "/"), "/")[0]
		if reservedPaths[first] {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Route not found"})
			return
		}
		c.Next()
	}
}

// AuthMiddleware は Authorization ヘッダーをチェックし、ユーザー名を Context に格納します
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization required"})
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			return utils.JwtSecret, nil
		})
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			return
		}

		username, _ := claims["username"].(string)
		role, _ := claims["role"].(string)

		pathUser := c.Param("user")
		if pathUser != "" && username != pathUser {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "You are not authorized to access this resource"})
			return
		}

		c.Set("username", username)
		c.Set("role", role) // 💡ここでコンテキストにroleを渡す
		c.Next()
	}
}
