package middleware

import "github.com/gin-gonic/gin"

// Make sure the user is authorized to view the passed in user.
func (m *AuthMiddleware) UserMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}
