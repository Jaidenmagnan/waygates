// This package handles the endpoints associated with the dashboard.
package handlers

import (
	"net/http"

	"github.com/Jaidenmagnan/waygates/components"
	"github.com/Jaidenmagnan/waygates/models"
	"github.com/gin-gonic/gin"
)

// Returns the template view for the dashboard.
func Dashboard(c *gin.Context) {
	value, ok := c.Get("user")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	user := value.(models.User)

	components.Dashboard(user).Render(c.Request.Context(), c.Writer)
}
