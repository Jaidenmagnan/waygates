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

	authService := services.NewAuthService(userRepository)

	authHandler := handlers.NewAuthHandler(authService)

	authMiddleware := middleware.NewAuthMiddleware(authService)

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		components.Hello("Jaiden").Render(c.Request.Context(), c.Writer)
	})

	r.POST("/signup", authHandler.Signup)
	r.POST("/signin", authHandler.Signin)
	r.POST("/signout", authHandler.Signout)
	r.GET("/dashboard", authMiddleware.AuthMiddleware(), handlers.Dashboard)

	if err := r.Run(); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
