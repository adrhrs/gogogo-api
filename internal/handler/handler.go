package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/adrhrs/gogogo-api/internal/db"

	"github.com/gin-gonic/gin"
)

func Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func DBPing(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	err := db.Pool.Ping(ctx)
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "DB not reachable"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "DB is reachable"})
}
