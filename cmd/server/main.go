package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/TheBizii/outfit7-ad-mediation/internal/config"
	"github.com/TheBizii/outfit7-ad-mediation/internal/db"
)

func main() {
	cfg := config.Load()
	database := db.Connect(cfg.PSQLUrl)
	defer database.Close() // will run when main() stops running

	router := gin.Default()
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Still breathing!",
		})
	})

	log.Println("Starting server on http://localhost:" + cfg.AppPort + "...")
	if err := router.Run(":" + cfg.AppPort); err != nil {
		log.Fatal("Server failed:", err)
	}
}
