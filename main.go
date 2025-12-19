package main

import (
	"log"

	"github.com/Jaidenmagnan/waygates/components"
	"github.com/Jaidenmagnan/waygates/db"
	"github.com/Jaidenmagnan/waygates/handlers"
	"github.com/Jaidenmagnan/waygates/middleware"
	"github.com/Jaidenmagnan/waygates/repositories"
	"github.com/Jaidenmagnan/waygates/services"
	"github.com/gin-gonic/gin"
)

func main() {
	if err := db.Connect(); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	defer db.DB.Close()

	userRepository := repositories.NewUserRepository(db.DB)
	waygateRepository := repositories.NewWaygateRepository(db.DB)

	authService := services.NewAuthService(userRepository)
	waygateService := services.NewWaygateService(waygateRepository)

	authHandler := handlers.NewAuthHandler(authService)
	waygateHandler := handlers.NewWaygateHandler(waygateService)

	authMiddleware := middleware.NewAuthMiddleware(authService)

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		components.Hello("Jaiden").Render(c.Request.Context(), c.Writer)
	})

	// Authentication routes.
	r.POST("/signup", authHandler.Signup)
	r.POST("/signin", authHandler.Signin)
	r.POST("/signout", authHandler.Signout)

	// User routes.
	r.GET("/user/:id/waygates", waygateHandler.ListUserWaygates)

	// Waygate routes.
	r.POST("/waygate/create", authMiddleware.AuthMiddleware(), waygateHandler.CreateWaygate)
	r.POST("/waygate/:id/delete", authMiddleware.AuthMiddleware(), waygateHandler.DeleteWaygate)

	r.GET("/waygate/:id/view", authMiddleware.AuthMiddleware(), waygateHandler.WaygatePage)

	r.GET("/signup", authMiddleware.SigninAndSignupMiddleware(), authHandler.SignupPage)
	r.GET("/signin", authMiddleware.SigninAndSignupMiddleware(), authHandler.SigninPage)

	r.GET("/", authMiddleware.AuthMiddleware(), handlers.Dashboard)

	if err := r.Run(); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
