package main

import (
	"log"

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
	waygateLinkRepository := repositories.NewWaygateLinkRepository(db.DB)

	authService := services.NewAuthService(userRepository)
	waygateService := services.NewWaygateService(waygateRepository)
	waygateLinkService := services.NewWaygateLinkService(waygateLinkRepository)

	authHandler := handlers.NewAuthHandler(authService)
	waygateHandler := handlers.NewWaygateHandler(waygateService)
	dashboardHandler := handlers.NewDashboardHandler(waygateService, waygateLinkService)
	waygateLinkHandler := handlers.NewWaygateLinkHandler(waygateService, waygateLinkService)

	authMiddleware := middleware.NewAuthMiddleware(authService)

	r := gin.Default()

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
			waygates.GET("/:waygate_id", waygateHandler.ViewWaygate)
			waygates.DELETE("/:waygate_id", waygateHandler.DeleteWaygate)

			links := waygates.Group("/:waygate_id/links")
			{
				links.POST("/", waygateLinkHandler.CreateWaygateLink)
				links.DELETE("/:waygate_link_id", waygateLinkHandler.DeleteWaygateLink)
			}

		}

	}

	r.GET("/signup", authMiddleware.SigninAndSignupMiddleware(), authHandler.SignupPage)
	r.GET("/signin", authMiddleware.SigninAndSignupMiddleware(), authHandler.SigninPage)
	r.GET("/", authMiddleware.AuthMiddleware(), dashboardHandler.Dashboard)

	waygates := r.Group("/waygates", authMiddleware.AuthMiddleware())
	{
		waygates.GET("/:id", dashboardHandler.ViewWaygate)
	}

	if err := r.Run(); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
