package main

import (
	"log"

	"github.com/Jaidenmagnan/waygates/components"
	"github.com/Jaidenmagnan/waygates/db"
	auth "github.com/Jaidenmagnan/waygates/handlers"
	"github.com/gin-gonic/gin"
)



func main() {
  if err := db.Connect(); err != nil {
    log.Fatal("Failed to connect to database:", err)
  }

  defer db.DB.Close()

  r := gin.Default()


  r.GET("/ping", func(c *gin.Context) {
	  components.Hello("Jaiden").Render(c.Request.Context(), c.Writer)
  })

  r.POST("/signup", auth.Signup)
  r.POST("/signin", auth.Signin)

  if err := r.Run(); err != nil {
    log.Fatalf("failed to run server: %v", err)
  }
}