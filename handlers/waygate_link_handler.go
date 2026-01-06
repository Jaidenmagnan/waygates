package handlers

import (
	"net/http"
	"strconv"

	"github.com/Jaidenmagnan/waygates/models"
	"github.com/Jaidenmagnan/waygates/services"
	"github.com/gin-gonic/gin"
)

type WaygateLinkHandler struct {
	waygateService     *services.WaygateService
	waygateLinkService *services.WaygateLinkService
}

func NewWaygateLinkHandler(waygateService *services.WaygateService, waygateLinkService *services.WaygateLinkService) *WaygateLinkHandler {
	return &WaygateLinkHandler{
		waygateService:     waygateService,
		waygateLinkService: waygateLinkService,
	}
}

func (h *WaygateLinkHandler) DeleteWaygateLink(c *gin.Context) {
	waygateID := c.Param("waygate_id")
	waygateLinkID := c.Param("waygate_link_id")

	waygateIdInt, err := strconv.Atoi(waygateID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid waygate id",
		})
		return
	}

	waygateLinkIdInt, err := strconv.Atoi(waygateLinkID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid waygate id",
		})
		return
	}

	user := c.MustGet("user").(models.User)

	canAccess, err := h.waygateService.CanUserAccessWaygate(user.ID, waygateIdInt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if !canAccess {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "you do not have permission to delete this waygate",
		})
		return
	}

	err = h.waygateLinkService.DeleteWaygateLink(waygateLinkIdInt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "waygate deleted successfully",
	})

}

// Create a link for a waygate.
func (h *WaygateLinkHandler) CreateWaygateLink(c *gin.Context) {
	user := c.MustGet("user").(models.User)

	type CreateWaygateLinkRequest struct {
		Name      string `form:"name" binding:"required"`
		Link      string `form:"link" binding:"required"`
		WaygateID int    `form:"waygate_id" binding:"required"`
	}

	var createWaygateLinkRequest CreateWaygateLinkRequest
	if err := c.ShouldBind(&createWaygateLinkRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to read body",
		})
		return
	}

	canAccess, err := h.waygateService.CanUserAccessWaygate(user.ID, createWaygateLinkRequest.WaygateID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if !canAccess {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "you do not have permission to view this waygate",
		})
		return
	}

	waygateLink, err := h.waygateLinkService.CreateWaygateLink(
		createWaygateLinkRequest.Name,
		createWaygateLinkRequest.Link,
		createWaygateLinkRequest.WaygateID,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"link": waygateLink,
	})
}
