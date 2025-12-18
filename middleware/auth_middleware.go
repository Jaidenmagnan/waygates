package middleware

import (
	"github.com/Jaidenmagnan/waygates/services"
	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct{
	authService *services.AuthService 
}

func NewAuthMiddleware(authService *services.AuthService) *AuthMiddleware {
	return &AuthMiddleware{
		authService: authService,
	}
}

func (m *AuthMiddleware) AuthMiddleware() gin.HandlerFunc {
  return func(c *gin.Context) {
	tokenString, err := c.Cookie("Authorization")
	if err != nil || tokenString == "" {
	  c.AbortWithStatusJSON(401, gin.H{
		"error": "unauthorized",
	  })
	  return
	}

	user, ok := m.authService.GetUserFromToken(tokenString)
	if !ok {
	  c.AbortWithStatusJSON(401, gin.H{
		"error": "unauthorized",
	  })
	  return
	}

	c.Set("user", user)
	c.Next()
  }
}