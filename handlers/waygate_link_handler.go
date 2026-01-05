
type WaygateLinkHandler struct {
	WaygateService *services.WaygateService
}

func NewWaygateLinkHandler(waygateService *services.WaygateService) *WaygateLinkHandler {
	return &WaygateLinkHandler{
		WaygateService: waygateService,
	}
}

// List links for a waygate.
func (h *WaygateLinkHandler) ListLinks(c *gin.Context) {
	waygateID := c.Param("waygate_id")

	waygateIdInt, err := strconv.Atoi(waygateID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid waygate id",
		})
		return
	}

	user := c.MustGet("user").(models.User)

	canAccess, err := h.WaygateService.CanUserAccessWaygate(user.ID, waygateIdInt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if !canAccess {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "forbidden",
		})
		return
	}

	links, err := h.WaygateService.ListLinks(waygateIdInt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"links": links,
	})
}