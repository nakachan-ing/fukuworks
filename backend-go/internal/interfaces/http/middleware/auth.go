package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
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
		token := c.GetHeader("Authorization")
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			return
		}

		// トークンは "Bearer mock-token-for-<username>" の形式と仮定
		if !strings.HasPrefix(token, "Bearer mock-token-for-") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
			return
		}

		// ユーザー名を抽出
		username := strings.TrimPrefix(token, "Bearer mock-token-for-")
		pathUser := c.Param("user")

		// 認可: パスの :user と一致していなければ拒否
		if username != pathUser {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "You are not authorized to access this resource"})
			return
		}

		// コンテキストにユーザー名を保存（必要に応じて）
		c.Set("username", username)
		c.Next()
	}
}
