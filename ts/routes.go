package ts

import (
	"github.com/gin-gonic/gin"
)

func (ts *TapeStatsApp) AddRoutes(r *gin.Engine) {
	// Main, simple, handlers
	r.GET("/", ts.IndexHandler)
	r.GET("/ping", ts.PingHandler)

	// Auth

	// Account Mgmt
	r.POST("/auth/register/org", ts.RegisterOrgHandler)         //Org+1st user
	r.POST("/auth/register/user", ts.RegisterUserHandler)       // Add user to existing org
	r.POST("/auth/register/account", ts.RegisterAccountHandler) // Add account to existing org

	// Submission handlers
	//r.POST("/submit/record", ts.LoadRecordHandler)
	r.POST("/submit/unparsed", ts.LoadUnparsedHandler)

	// Stats handlers

}
