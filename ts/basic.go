package ts

import (
	"github.com/gin-gonic/gin"
)

func (ts *TapeStatsApp) PingHandler(c *gin.Context) {
	li, err := ts.Ctxer(c)
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

func (ts *TapeStatsApp) IndexHandler(c *gin.Context) {
	li, err := ts.Ctxer(c)
	if err != nil {
		c.HTML(500, "error.html", gin.H{
			"error": "Unknown Server Error",
			"title": "Error: Server Error",
		})
		return
	}

	c.HTML(200, "index.html", gin.H{
		"title": "Welcome!",
	})
	li.Log.Debug().Msg("Ping-Pong")
}
