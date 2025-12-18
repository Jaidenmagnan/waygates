package handlers

import (
	"net/http"

	"github.com/Jaidenmagnan/waygates/services"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
  authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
  return &AuthHandler{
    authService: authService,
  }
}

func (h *AuthHandler) Signup(c *gin.Context) {
  type SignUpRequest struct {
    Email   string `json:"email" binding:"required,email"`
    Username string `json:"username" binding:"required"`
    Password string `json:"password" binding:"required"`
  }

  var signUpRequest SignUpRequest
  if err := c.BindJSON(&signUpRequest); err != nil {
    c.JSON(http.StatusBadRequest, gin.H{
      "error": "failed to read body",
    })
    return
  }

  user, err := h.authService.Signup(signUpRequest.Username, signUpRequest.Email, signUpRequest.Password)
  if err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{
      "error": err.Error(),
    })
    return
  }

  c.JSON(http.StatusOK, gin.H{
    "message": user,
  })
}

func (h *AuthHandler) Signin(c *gin.Context) {
  type SignInRequest struct {
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required"`
  }

  var signInRequest SignInRequest

  if err := c.BindJSON(&signInRequest); err != nil {
    c.JSON(http.StatusBadRequest, gin.H{
      "error": "failed to read body",
    })
    return
  }

  user ,err := h.authService.Signin(signInRequest.Email, signInRequest.Password)
  if err != nil {
    c.JSON(http.StatusUnauthorized, gin.H{
      "error": err.Error(),
    })
    return
  }

  tokenString, err := h.authService.GenerateToken(user)
  if err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{
      "error": "failed to generate token",
    })
    return
  }

  c.SetSameSite(http.SameSiteLaxMode)
  c.SetCookie("Authorization", tokenString, 7*24*60*60, "", "", true, true)

  c.JSON(http.StatusOK, gin.H{
    "message": "user signed in successfully",
  })
}

func (h *AuthHandler) Signout(c *gin.Context) {
  c.SetCookie("Authorization", "", -1, "", "", true, true)

  c.JSON(http.StatusOK, gin.H{
    "message": "user signed out successfully",
  })
}