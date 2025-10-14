package models

import "time"

// database models
type PriorityList struct {
	ID          int
	CountryCode string
	AdType      string
	LastUpdated time.Time
}

type PriorityNetwork struct {
	ID             int
	PriorityListId int
	NetworkName    string
	Score          float32
}

// structs for the ad networks update route
type UpdatePayloadNetwork struct {
	NetworkName string  `json:"networkName" binding:"required"`
	Score       float32 `json:"score" binding:"required"`
}

type UpdateNetworksRequest struct {
	Networks []UpdatePayloadNetwork `json:"networks" binding:"required"`
}

// structs for the get ad networks route
type GetNetworksRequest struct {
	CountryCode string `json:"countryCode" binding:"required"`
	AdType      string `json:"adType" binding:"required"`
	Platform    string `json:"platform"`
	OSVersion   string `json:"osVersion"`
	AppName     string `json:"appName"`
	AppVersion  string `json:"appVersion"`
}
