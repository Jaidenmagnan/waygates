// This package handles the endpoints associated with the dashboard.
package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Returns the template view for the dashboard.
func Dashboard(c *gin.Context) {
	user, ok := c.Get("user")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": user,
	})
}
