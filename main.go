package main

import (
	"log"
	"os"

	"github.com/adrhrs/gogogo-api/internal/db"
	"github.com/adrhrs/gogogo-api/internal/handler"
	"github.com/gin-gonic/gin"
)

func main() {
	log.Println("starting..")

	// Initialize DB
	if err := db.Init(); err != nil {
		log.Fatalf("failed to connect to DB: %v", err)
	}

	// Setup router
	r := gin.Default()
	r.GET("/ping", handler.Ping)
	r.GET("/dbping", handler.DBPing)

	// Use port from Render (PORT is set automatically)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)

}
