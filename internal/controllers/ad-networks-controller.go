package controllers

import (
	"net/http"

	"github.com/TheBizii/outfit7-ad-mediation/internal/models"
	"github.com/TheBizii/outfit7-ad-mediation/internal/services"
	"github.com/gin-gonic/gin"
)

func GetAdNetworks(ctx *gin.Context) {
	ctx.JSON(http.StatusNotImplemented, gin.H{
		"message": "This route has no implementation.",
	})
}

func UpdateAdNetworks(ctx *gin.Context) {
	countryCode := ctx.Param("country-code")
	adType := ctx.Param("ad-type")

	var body models.UpdateNetworksRequest
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	err := services.UpsertPriorityList(countryCode, adType, body)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":             "ok",
		"countryCode":        countryCode,
		"adType":             adType,
		"numUpdatedNetworks": len(body.Networks),
	})
}
