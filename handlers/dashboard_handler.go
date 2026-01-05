// This package handles the endpoints associated with the dashboard.
package handlers

import (
	"net/http"
	"strconv"

	"github.com/Jaidenmagnan/waygates/components"
	"github.com/Jaidenmagnan/waygates/models"
	"github.com/Jaidenmagnan/waygates/services"
	"github.com/gin-gonic/gin"
)

type DashboardHandler struct {
	waygateService *services.WaygateService
}

func NewDashboardHandler(waygateService *services.WaygateService) *DashboardHandler {
	return &DashboardHandler{
		waygateService: waygateService,
	}
}

// Returns the template view for the dashboard.
func (h *DashboardHandler) Dashboard(c *gin.Context) {
	value, ok := c.Get("user")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	user := value.(models.User)

	s, err := h.waygateService.ListUserWaygates(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load waygates"})
		return
	}

	components.Dashboard(user, s).Render(c.Request.Context(), c.Writer)
}

func (h *DashboardHandler) ViewWaygate(c *gin.Context) {
	waygateID := c.Param("id")

	waygateIDInt, err := strconv.Atoi(waygateID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid waygate id",
		})
		return
	}

	waygate, err := h.waygateService.GetWaygateByID(waygateIDInt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load waygate"})
		return
	}

	components.Waygate(waygate).Render(c.Request.Context(), c.Writer)
}
