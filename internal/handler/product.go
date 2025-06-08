package handler

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CreateProductRequest struct {
	Name      string  `json:"name" binding:"required"`
	SellPrice float64 `json:"sell_price" binding:"required"`
}

func CreateProductHandler(db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req CreateProductRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		code := generateProductCode()

		query := `
	INSERT INTO products (product_code, name, sell_price, ip_address)
	VALUES ($1, $2, $3, $4)
	RETURNING product_id, product_code, name, sell_price, ip_address, created_at, updated_at
`

		ip := getClientIP(c)

		row := db.QueryRow(c, query, code, req.Name, req.SellPrice, ip)

		var (
			productID   int
			productCode string
			name        string
			sellPrice   float64
			ipAddress   string
			createdAt   time.Time
			updatedAt   time.Time
		)

		if err := row.Scan(&productID, &productCode, &name, &sellPrice, &ipAddress, &createdAt, &updatedAt); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "DB insert failed", "details": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"product_id":   productID,
			"product_code": productCode,
			"name":         name,
			"sell_price":   sellPrice,
			"ip_address":   ipAddress,
			"created_at":   createdAt,
			"updated_at":   updatedAt,
		})
	}
}

// generateProductCode creates a 10-char random alphanumeric string
func generateProductCode() string {
	b := make([]byte, 5)
	_, err := rand.Read(b)
	if err != nil {
		return "XXXXXXXXXX" // fallback
	}
	return strings.ToUpper(hex.EncodeToString(b)) // hex gives 10 chars
}

func getClientIP(c *gin.Context) string {
	ip := c.ClientIP()
	if ip == "" {
		return "UNKNOWN"
	}
	return ip
}
