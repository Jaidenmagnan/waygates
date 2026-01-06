package handlers

import (
	"net/http"
	"strconv"

	"github.com/Jaidenmagnan/waygates/models"
	"github.com/Jaidenmagnan/waygates/services"
	"github.com/gin-gonic/gin"
)

type WaygateHandler struct {
	waygateService *services.WaygateService
}

func NewWaygateHandler(waygateService *services.WaygateService) *WaygateHandler {
	return &WaygateHandler{
		waygateService: waygateService,
	}
}

// Create a waygate.
func (h *WaygateHandler) CreateWaygate(c *gin.Context) {
	user := c.MustGet("user").(models.User)

	type CreateWaygateRequest struct {
		Name string `form:"name" binding:"required"`
	}

	var createWaygateRequest CreateWaygateRequest
	if err := c.ShouldBind(&createWaygateRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to read body",
		})
		return
	}

	waygate, err := h.waygateService.CreateWaygate(createWaygateRequest.Name, user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"waygate": waygate,
	})
}

// View a waygate.
func (h *WaygateHandler) ViewWaygate(c *gin.Context) {
	waygateID := c.Param("waygate_id")

	waygateIdInt, err := strconv.Atoi(waygateID)
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
			"error": "you do not have permission to view this waygate",
		})
		return
	}

	waygate, err := h.waygateService.GetWaygateByID(waygateIdInt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"waygate": waygate,
	})
}

// List all waygates for a user.
func (h *WaygateHandler) ListUserWaygates(c *gin.Context) {
	user := c.MustGet("user").(models.User)

	waygates, err := h.waygateService.ListUserWaygates(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"waygates": waygates,
	})
}

func (h *WaygateHandler) DeleteWaygate(c *gin.Context) {
	waygateID := c.Param("waygate_id")

	waygateIdInt, err := strconv.Atoi(waygateID)
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

	err = h.waygateService.DeleteWaygate(waygateIdInt)
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
