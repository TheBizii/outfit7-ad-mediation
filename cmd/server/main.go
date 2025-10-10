package main

import (
	"log"

	"github.com/TheBizii/outfit7-ad-mediation/internal/config"
	"github.com/TheBizii/outfit7-ad-mediation/internal/db"
	"github.com/TheBizii/outfit7-ad-mediation/internal/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()
	database := db.Connect(cfg.PSQLUrl)
	defer database.Close() // will run when main() stops running

	r := gin.Default()
	routes.RegisterRoutes(r)

	log.Println("Starting server on http://localhost:" + cfg.AppPort + "...")
	if err := r.Run(":" + cfg.AppPort); err != nil {
		log.Fatal("Server failed:", err)
	}
}
