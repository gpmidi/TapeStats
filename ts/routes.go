package ts

import (
	"github.com/gin-gonic/gin"
)

func (ts *TapeStatsApp) AddRoutes(r *gin.Engine) {
	// Main, simple, handlers
	r.GET("/", IndexHandler)
	r.GET("/ping", PingHandler)

	// Submission handlers

	// Stats handlers

}
