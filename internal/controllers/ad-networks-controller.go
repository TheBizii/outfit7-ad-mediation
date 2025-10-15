package controllers

import (
	"net/http"

	"github.com/TheBizii/outfit7-ad-mediation/internal/models"
	"github.com/TheBizii/outfit7-ad-mediation/internal/services"
	"github.com/gin-gonic/gin"
)

func GetNetworksDashboard(ctx *gin.Context) {
	lists, err := services.GetDashboardPriorityLists()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"lists":  lists,
	})
}

func GetAdNetworks(ctx *gin.Context) {
	body := models.GetNetworksRequest{
		CountryCode: ctx.Param("countryCode"),
		AdType:      ctx.Param("adType"),
		Platform:    ctx.Query("platform"),
		OSVersion:   ctx.Query("osVersion"),
		AppName:     ctx.Query("appName"),
		AppVersion:  ctx.Query("appVersion"),
	}

	if body.CountryCode == "" || body.AdType == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "You must supply the countryCode and adType parameters.",
		})
		return
	}

	networks, err := services.GetAdNetworks(body)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":      "ok",
		"countryCode": body.CountryCode,
		"adType":      body.AdType,
		"networks":    networks,
	})
}

func UpdateAdNetworks(ctx *gin.Context) {
	countryCode := ctx.Param("countryCode")
	adType := ctx.Param("adType")

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
