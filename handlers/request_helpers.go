package handlers

import (
	"net/http"

	"github.com/Jaidenmagnan/waygates/components"
	"github.com/gin-gonic/gin"
)

func wantsHTML(c *gin.Context) bool {
	return c.GetHeader("HX-Request") == "true" || c.GetHeader("Accept") == "" || c.NegotiateFormat(gin.MIMEHTML, gin.MIMEJSON) == gin.MIMEHTML
}

func redirectBack(c *gin.Context, fallback string) {
	target := c.GetHeader("Referer")
	if target == "" {
		target = fallback
	}

	if c.GetHeader("HX-Request") == "true" {
		c.Header("HX-Redirect", target)
		c.Status(http.StatusOK)
		return
	}

	c.Redirect(http.StatusSeeOther, target)
}

func respondError(c *gin.Context, status int, title string, message string) {
	if wantsHTML(c) {
		c.Status(status)
		components.ErrorPage(status, title, message).Render(c.Request.Context(), c.Writer)
		return
	}

	c.JSON(status, gin.H{"error": message})
}

func requestBaseURL(c *gin.Context) string {
	scheme := "http"
	if c.Request.TLS != nil || c.GetHeader("X-Forwarded-Proto") == "https" {
		scheme = "https"
	}

	return scheme + "://" + c.Request.Host
}
