package ts

import (
	"github.com/gin-gonic/gin"
	"github.com/gpmidi/TapeStats/ts/handlers"
)

func (ts *TapeStatsApp) AddRoutes(r *gin.Engine) {
	// Main, simple, handlers
	r.GET("/", handlers.IndexHandler)
	r.GET("/ping", handlers.PingHandler)

	// Submission handlers

	// Stats handlers

}
