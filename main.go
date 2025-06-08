package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/adrhrs/gogogo-api/internal/db"
	"github.com/adrhrs/gogogo-api/internal/handler"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type App struct {
	DB *pgxpool.Pool
}

func main() {
	log.Println("starting..")

	pool, err := db.Init()
	if err != nil {
		log.Fatalf("failed to connect to DB: %v", err)
	}

	a := &App{
		DB: pool,
	}

	// Setup router
	r := gin.Default()

	// CORS config
	config := cors.Config{
		AllowOrigins:     []string{"*"}, // Use specific origins in production
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
	r.Use(cors.New(config))

	// Setup router
	// r := gin.Default()
	r.GET("/ping", handler.Ping)
	r.GET("/dbping", a.DBPing)

	// Use port from Render (PORT is set automatically)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)

}

func (a *App) DBPing(c *gin.Context) {
	if a.DB == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB not initialized"})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	if err := a.DB.Ping(ctx); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "DB not reachable", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "DB is reachable"})
}
