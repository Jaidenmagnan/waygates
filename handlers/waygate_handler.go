package handlers

import (
	"fmt"
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
		respondError(c, http.StatusBadRequest, "The gate name is missing", "Give your web ring a name and try again.")
		return
	}

	waygate, err := h.waygateService.CreateWaygate(createWaygateRequest.Name, user.ID)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "The gate stones would not settle", err.Error())
		return
	}

	if wantsHTML(c) {
		c.Header("HX-Redirect", fmt.Sprintf("/waygates/%d", waygate.ID))
		c.Redirect(http.StatusSeeOther, fmt.Sprintf("/waygates/%d", waygate.ID))
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
		respondError(c, http.StatusBadRequest, "The waygate mark is unclear", "That waygate address is not valid.")
		return
	}

	user := c.MustGet("user").(models.User)

	canAccess, err := h.waygateService.CanUserAccessWaygate(user.ID, waygateIdInt)
	if err != nil {
		respondError(c, http.StatusNotFound, "No waygate here", "That waygate could not be found.")
		return
	}

	if !canAccess {
		respondError(c, http.StatusForbidden, "The gate is warded", "You do not have permission to view this waygate.")
		return
	}

	waygate, err := h.waygateService.GetWaygateByID(waygateIdInt)
	if err != nil {
		respondError(c, http.StatusNotFound, "No waygate here", "That waygate could not be found.")
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
		respondError(c, http.StatusInternalServerError, "The map went dim", err.Error())
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
		respondError(c, http.StatusBadRequest, "The waygate mark is unclear", "That waygate address is not valid.")
		return
	}

	user := c.MustGet("user").(models.User)

	canAccess, err := h.waygateService.CanUserAccessWaygate(user.ID, waygateIdInt)
	if err != nil {
		respondError(c, http.StatusNotFound, "No waygate here", "That waygate could not be found.")
		return
	}

	if !canAccess {
		respondError(c, http.StatusForbidden, "The gate is warded", "You do not have permission to delete this waygate.")
		return
	}

	err = h.waygateService.DeleteWaygate(waygateIdInt)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "The stones would not move", err.Error())
		return
	}

	if wantsHTML(c) {
		c.Header("HX-Redirect", "/")
		c.Redirect(http.StatusSeeOther, "/")
		return
	}

	c.JSON(200, gin.H{
		"message": "waygate deleted successfully",
	})
}

func (h *WaygateHandler) DeleteWaygateFromForm(c *gin.Context) {
	h.DeleteWaygate(c)
}
