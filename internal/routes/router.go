package routes

import (
	"github.com/TheBizii/outfit7-ad-mediation/internal/controllers"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine) {
	// health route (similar to ping in other applications)
	router.GET("/health", GetHealth)

	api := router.Group("/api/v1")
	{
		// retrieve sorted priority lists for mobile apps
		api.GET("/ad_networks/:countryCode/:adType", controllers.GetAdNetworks)
		// update route for ad networks
		api.PUT("/ad_networks/:countryCode/:adType", controllers.UpdateAdNetworks)

		// this route is called from the dashboard
		api.GET("/dashboard", GetDashboard)
	}
}
