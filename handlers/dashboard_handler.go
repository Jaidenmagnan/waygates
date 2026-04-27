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
	waygateService     *services.WaygateService
	waygateLinkService *services.WaygateLinkService
}

func NewDashboardHandler(waygateService *services.WaygateService, waygateLinkService *services.WaygateLinkService) *DashboardHandler {
	return &DashboardHandler{
		waygateService:     waygateService,
		waygateLinkService: waygateLinkService,
	}
}

// Returns the template view for the dashboard.
func (h *DashboardHandler) Dashboard(c *gin.Context) {
	value, ok := c.Get("user")
	if !ok {
		respondError(c, http.StatusUnauthorized, "The gate stayed closed", "Sign in before opening your rings.")
		return
	}

	user := value.(models.User)

	s, err := h.waygateService.ListUserWaygates(user.ID)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "The map went dim", "Your waygates could not be loaded.")
		return
	}

	components.Dashboard(user, s).Render(c.Request.Context(), c.Writer)
}

func (h *DashboardHandler) ViewWaygate(c *gin.Context) {
	waygateID := c.Param("id")

	waygateIDInt, err := strconv.Atoi(waygateID)
	if err != nil {
		respondError(c, http.StatusBadRequest, "The waygate mark is unclear", "That waygate address is not valid.")
		return
	}

	waygate, err := h.waygateService.GetWaygateByID(waygateIDInt)
	if err != nil {
		respondError(c, http.StatusNotFound, "No waygate here", "That waygate could not be found.")
		return
	}

	user := c.MustGet("user").(models.User)
	if waygate.UserId != user.ID {
		respondError(c, http.StatusForbidden, "The gate is warded", "You do not have permission to view this waygate.")
		return
	}

	waygateLinks, err := h.waygateLinkService.ListWaygateLinks(waygate.ID)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "The map went dim", "This waygate's links could not be loaded.")
		return
	}

	components.Waygate(waygate, waygateLinks, requestBaseURL(c)).Render(c.Request.Context(), c.Writer)
}
