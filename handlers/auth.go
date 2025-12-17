package auth

import (
	"net/http"

	"github.com/Jaidenmagnan/waygates/models"
	"github.com/Jaidenmagnan/waygates/repositories"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func Signup(c *gin.Context) {
  type SignUpRequest struct {
    Email   string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required"`
  }

  var signUpRequest SignUpRequest

  if err := c.BindJSON(&signUpRequest); err != nil {
    c.JSON(http.StatusBadRequest, gin.H{
      "error": "failed to read body",
    })
    return
  }

  hash, err := bcrypt.GenerateFromPassword([]byte(signUpRequest.Password), bcrypt.DefaultCost)
  if err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{
      "error": "failed to hash password",
    })
    return
  }

  userRepository := repositories.NewUserRepository()

  _, err = userRepository.Create(models.User{
    Email:    signUpRequest.Email,
    Password: string(hash),
  })

  if err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{
      "error": "failed to create user",
    })
    return
  }

  c.JSON(http.StatusOK, gin.H{
    "message": "user signed up successfully",
  })
}

func Signin(c *gin.Context) {
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

  userRepository := repositories.NewUserRepository()
  user, err := userRepository.GetByEmail(signInRequest.Email)

  if err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{
      "error": "failed to get user",
    })
    return
  }

  if user == nil {
    c.JSON(http.StatusUnauthorized, gin.H{
      "error": "invalid email or password",
    })
    return
  }
  
  err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(signInRequest.Password))
  if err != nil {
    c.JSON(http.StatusUnauthorized, gin.H{
      "error": "invalid email or password",
    })
    return
  }

  c.JSON(http.StatusOK, gin.H{
    "message": "user signed in successfully",
  })
}