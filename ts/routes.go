package ts

import (
	"github.com/gin-gonic/gin"
)

func (ts *TapeStatsApp) AddRoutes(r *gin.Engine) {
	// Main, simple, handlers
	r.GET("/", ts.IndexHandler)
	r.GET("/ping", ts.PingHandler)

	// Auth
	r.POST("/auth/register", ts.RegisterAccountHandler)

	// Submission handlers

	// Stats handlers

}
