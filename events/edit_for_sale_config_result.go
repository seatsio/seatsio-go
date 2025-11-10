package events

import "time"

type ForSaleRateLimitInfo struct {
	RateLimitRemainingCalls int        `json:"rateLimitRemainingCalls"`
	RateLimitResetDate      *time.Time `json:"rateLimitResetDate"`
}

type EditForSaleConfigResult struct {
	ForSaleRateLimitInfo *ForSaleRateLimitInfo `json:"rateLimitInfo"`
	ForSaleConfig        *ForSaleConfig        `json:"forSaleConfig"`
}
