package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetAdNetworks(ctx *gin.Context) {
	ctx.JSON(http.StatusNotImplemented, gin.H{
		"message": "This route has no implementation.",
	})
}

func UpdateAdNetworks(ctx *gin.Context) {
	ctx.JSON(http.StatusNotImplemented, gin.H{
		"message": "This route has no implementation.",
	})
}
