package handlers

import (
	"errors"
	"fmt"
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
		respondError(c, http.StatusBadRequest, "The waygate mark is unclear", "That waygate address is not valid.")
		return
	}

	waygateLinkIdInt, err := strconv.Atoi(waygateLinkID)
	if err != nil {
		respondError(c, http.StatusBadRequest, "The path mark is unclear", "That site address is not valid.")
		return
	}

	user := c.MustGet("user").(models.User)

	canAccess, err := h.waygateService.CanUserAccessWaygate(user.ID, waygateIdInt)
	if err != nil {
		respondError(c, http.StatusNotFound, "No waygate here", "That waygate could not be found.")
		return
	}

	if !canAccess {
		respondError(c, http.StatusForbidden, "The gate is warded", "You do not have permission to change this waygate.")
		return
	}

	err = h.waygateLinkService.DeleteWaygateLink(waygateLinkIdInt)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "The path would not fade", err.Error())
		return
	}

	if wantsHTML(c) {
		redirectBack(c, fmt.Sprintf("/waygates/%d", waygateIdInt))
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
		respondError(c, http.StatusBadRequest, "The friend site is missing", "Add a site name and URL, then try again.")
		return
	}

	canAccess, err := h.waygateService.CanUserAccessWaygate(user.ID, createWaygateLinkRequest.WaygateID)
	if err != nil {
		respondError(c, http.StatusNotFound, "No waygate here", "That waygate could not be found.")
		return
	}

	if !canAccess {
		respondError(c, http.StatusForbidden, "The gate is warded", "You do not have permission to change this waygate.")
		return
	}

	waygateLink, err := h.waygateLinkService.CreateWaygateLink(
		createWaygateLinkRequest.Name,
		createWaygateLinkRequest.Link,
		createWaygateLinkRequest.WaygateID,
	)

	if err != nil {
		respondError(c, http.StatusInternalServerError, "The path would not settle", err.Error())
		return
	}

	if wantsHTML(c) {
		c.Header("HX-Redirect", fmt.Sprintf("/waygates/%d", createWaygateLinkRequest.WaygateID))
		c.Redirect(http.StatusSeeOther, fmt.Sprintf("/waygates/%d", createWaygateLinkRequest.WaygateID))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"link": waygateLink,
	})
}

func (h *WaygateLinkHandler) RedirectToNext(c *gin.Context) {
	h.redirectNeighbor(c, true)
}

func (h *WaygateLinkHandler) RedirectToPrevious(c *gin.Context) {
	h.redirectNeighbor(c, false)
}

func (h *WaygateLinkHandler) redirectNeighbor(c *gin.Context, next bool) {
	waygateID, err := strconv.Atoi(c.Param("waygate_id"))
	if err != nil {
		respondError(c, http.StatusBadRequest, "The waygate mark is unclear", "That waygate address is not valid.")
		return
	}

	waygateLinkID, err := strconv.Atoi(c.Param("waygate_link_id"))
	if err != nil {
		respondError(c, http.StatusBadRequest, "The path mark is unclear", "That site address is not valid.")
		return
	}

	nextLink, prevLink, _, err := h.waygateLinkService.ResolveNeighbors(waygateID, waygateLinkID)
	if err != nil {
		if errors.Is(err, services.ErrWaygateLinkNotFound) {
			respondError(c, http.StatusNotFound, "No path in this ring", "That site is not part of this waygate.")
			return
		}

		respondError(c, http.StatusInternalServerError, "The gate would not open", "Could not open the next site.")
		return
	}

	target := prevLink.Link
	if next {
		target = nextLink.Link
	}

	c.Redirect(http.StatusFound, target)
}
