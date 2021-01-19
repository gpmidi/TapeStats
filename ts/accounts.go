package ts

import (
	"github.com/gin-gonic/gin"
)

func (ts *TapeStatsApp) RegisterAccountHandler(c *gin.Context) {
	li, err := Ctxer(c)
	if err != nil {
		ts.Log.Error().Err(c.Error(err)).Msg("Problem with getting ts")
		return
	}

	c.JSON(200, gin.H{
		"message": "pong",
		"request": li.Data(),
	})
	li.Log.Debug().Msg("Ping-Pong")
}
