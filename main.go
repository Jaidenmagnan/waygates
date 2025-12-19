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
	auth := r.Group("/auth")
	{
		auth.POST("/signup", authHandler.Signup)
		auth.POST("/signin", authHandler.Signin)
		auth.POST("/signout", authHandler.Signout)
	}

	api := r.Group("/api", authMiddleware.AuthMiddleware())
	{
		waygates := api.Group("/waygates")
		{
			// Waygate routes.
			waygates.GET("/", waygateHandler.ListUserWaygates)
			waygates.POST("/", waygateHandler.CreateWaygate)
			waygates.DELETE("/:id", waygateHandler.DeleteWaygate)

			waygates.GET("/:id", waygateHandler.ViewWaygate)
		}
	}

	r.GET("/signup", authMiddleware.SigninAndSignupMiddleware(), authHandler.SignupPage)
	r.GET("/signin", authMiddleware.SigninAndSignupMiddleware(), authHandler.SigninPage)
	r.GET("/", authMiddleware.AuthMiddleware(), handlers.Dashboard)

	if err := r.Run(); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
