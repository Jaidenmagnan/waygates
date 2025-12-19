// Middleware to check if the user is authenticated.
package middleware

import (
	"net/http"

	"github.com/Jaidenmagnan/waygates/services"
	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	authService *services.AuthService
}

func NewAuthMiddleware(authService *services.AuthService) *AuthMiddleware {
	return &AuthMiddleware{
		authService: authService,
	}
}

// AuthMiddleware checks for a valid JWT token in the request cookies.
func (m *AuthMiddleware) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie("Authorization")
		if err != nil || tokenString == "" {
			c.Redirect(http.StatusFound, "/signin")
			c.Abort()
			return
		}

		user, ok := m.authService.GetUserFromToken(tokenString)
		if !ok {
			c.Redirect(http.StatusFound, "/signin")
			c.Abort()
			return
		}

		c.Set("user", user)
		c.Next()
	}
}

// This middleware redirects authenticated users away from signin and signup pages.
func (m *AuthMiddleware) SigninAndSignupMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie("Authorization")
		if err != nil || tokenString == "" {
			c.Next()
			return
		}

		_, ok := m.authService.GetUserFromToken(tokenString)
		if !ok {
			c.Next()
			return
		}

		c.Redirect(http.StatusFound, "/")
	}
}
