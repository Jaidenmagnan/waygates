package services

import (
	"errors"
	"log/slog"
	"os"
	"time"

	"github.com/Jaidenmagnan/waygates/models"
	"github.com/Jaidenmagnan/waygates/repositories"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepository *repositories.UserRepository 
}

func NewAuthService(userRepository *repositories.UserRepository) *AuthService {
	return &AuthService{
		userRepository: userRepository,
	}
}

func (s* AuthService) Signup(username string, email string, password string) (models.User, error) {
  	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
  	if err != nil {
		slog.Error("failed to hash password", "error", err)
		return models.User{}, err
  	}

	user, err := s.userRepository.Create(models.CreateUser{
		Username: username,
		Email:    email,
		Password: string(hash),
	})
	if err != nil {
		slog.Error("failed to create user", "error", err)
		return models.User{}, err
	}

	return user, nil
}

func (s* AuthService) Signin(email string, password string) (*models.User, error) {
  	user, ok := s.userRepository.GetByEmail(email)
	if !ok {
		slog.Error("failed to get user by email")
		return nil, errors.New("invalid email or password")
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	return user, nil
}

func (s* AuthService) GenerateToken(user *models.User) (string, error) {
	expTime := jwt.NewNumericDate(jwt.TimeFunc().Add(7 * 24 * time.Hour))
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"exp":     expTime.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s* AuthService) GetUserFromToken(tokenString string) (*models.User, bool) {
	token, err := s.validateToken(tokenString)
	if err != nil {
		return nil, false
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		exp := claims["exp"].(float64)
		if float64(time.Now().Unix()) > exp {
			slog.Error("token expired")
			return nil, false
		}

		userId := int(claims["user_id"].(float64))
		if userId == 0 {
			slog.Error("invalid user id in token")
			return nil, false
		}

		user, ok := s.userRepository.GetById(userId)
		if !ok {
			slog.Error("user not found")
			return nil, false
		}

		return user, true;
	} else {
		slog.Error("invalid token claims")
		return nil, false
	}

}

func (s* AuthService) validateToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		return nil, err
	}

	return token, nil
}