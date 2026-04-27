// Implements utilities for authenticating the user and interfacing with the auth service.
package handlers

import (
	"net/http"

	"github.com/Jaidenmagnan/waygates/components"
	"github.com/Jaidenmagnan/waygates/services"
	"github.com/gin-gonic/gin"
)

// AuthHandler handles user authentication requests.
type AuthHandler struct {
	authService *services.AuthService
}

// NewAuthHandler creates a new AuthHandler.
func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Signup handles user signup requests.
func (h *AuthHandler) Signup(c *gin.Context) {
	type SignUpRequest struct {
		Email           string `form:"email" binding:"required,email"`
		Username        string `form:"username" binding:"required"`
		Password        string `form:"password" binding:"required"`
		ConfirmPassword string `form:"confirm_password" binding:"required,eqfield=Password"`
	}

	var signUpRequest SignUpRequest
	if err := c.ShouldBind(&signUpRequest); err != nil {
		respondError(c, http.StatusBadRequest, "The gate would not open", "Check the signup form and try again.")
		return
	}

	user, err := h.authService.Signup(signUpRequest.Username, signUpRequest.Email, signUpRequest.Password)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "The gate would not open", err.Error())
		return
	}

	tokenString, err := h.authService.GenerateToken(&user)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "The gate would not open", "Your account was created, but sign-in could not be completed.")
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 7*24*60*60, "/", "", c.Request.TLS != nil, true)
	c.Header("HX-Redirect", "/")
	c.Redirect(http.StatusSeeOther, "/")
}

// Signin handles user signin requests by creating a JWT token.
func (h *AuthHandler) Signin(c *gin.Context) {
	type SigninRequest struct {
		Email    string `form:"email" binding:"required,email"`
		Password string `form:"password" binding:"required"`
	}

	var signinRequest SigninRequest

	if err := c.ShouldBind(&signinRequest); err != nil {
		respondError(c, http.StatusBadRequest, "The gate would not open", "Check your email and password, then try again.")
		return
	}

	user, err := h.authService.Signin(signinRequest.Email, signinRequest.Password)
	if err != nil {
		respondError(c, http.StatusUnauthorized, "The gate stayed closed", "That email and password did not match.")
		return
	}

	tokenString, err := h.authService.GenerateToken(user)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "The gate would not open", "Sign-in could not be completed.")
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 7*24*60*60, "/", "", c.Request.TLS != nil, true)
	c.Header("HX-Redirect", "/")
	c.Redirect(http.StatusSeeOther, "/")
}

// Signout handles user signout requests by clearing the Authorization cookie.
func (h *AuthHandler) Signout(c *gin.Context) {
	c.SetCookie("Authorization", "", -1, "/", "", c.Request.TLS != nil, true)

	c.Header("HX-Redirect", "/signin")
	c.Redirect(http.StatusSeeOther, "/signin")
}

// Renders the signup page.
func (h *AuthHandler) SignupPage(c *gin.Context) {
	components.SignupPage().Render(c.Request.Context(), c.Writer)
}

// Renders the signin page.
func (h *AuthHandler) SigninPage(c *gin.Context) {
	components.SigninPage().Render(c.Request.Context(), c.Writer)
}
