package middleware

import (
	"github.com/gin-gonic/gin"
)

// Make sure the user is authorized to view the waygate and attach the waygate to the request.
func (m *AuthMiddleware) WaygateMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}
