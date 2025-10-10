package routes

import "github.com/gin-gonic/gin"

func RegisterRoutes(router *gin.Engine) {
	// health route (similar to ping in other applications)
	router.GET("/health", GetHealth)

	api := router.Group("/api/v1")
	{
		// retrieve sorted priority lists for mobile apps
		api.GET("/ad-networks", GetAdNetworks)
		// batch update route for ad networks
		api.POST("/ad-networks/update", UpdateAdNetworks)

		// this route is called from the dashboard
		api.GET("/dashboard", GetDashboard)
	}
}
