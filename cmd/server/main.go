package main

import (
	"log"

	"github.com/Jaidenmagnan/waygates/templates"
	"github.com/gin-gonic/gin"
)

func main() {
  r := gin.Default()

  r.GET("/ping", func(c *gin.Context) {
	  templates.Hello("Jaiden").Render(c.Request.Context(), c.Writer)
  })

  if err := r.Run(); err != nil {
    log.Fatalf("failed to run server: %v", err)
  }
}