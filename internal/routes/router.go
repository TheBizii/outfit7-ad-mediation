package routes

import (
	"github.com/TheBizii/outfit7-ad-mediation/internal/controllers"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine) {
	// health route (similar to ping in other applications)
	router.GET("/health", GetHealth)

	api := router.Group("/api/v1/ad_networks")
	{
		// this route is called from the dashboard
		api.GET("/dashboard", controllers.GetNetworksDashboard)

		// retrieve sorted priority lists for mobile apps
		api.GET("/:countryCode/:adType", controllers.GetAdNetworks)
		// update route for ad networks
		api.PUT("/:countryCode/:adType", controllers.UpdateAdNetworks)
	}
}
