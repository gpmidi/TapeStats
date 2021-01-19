package ts

import (
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
)

func PingHandler(c *gin.Context) {
	li, err := Ctxer(c)
	if err != nil {
		li.Log.Error().Err(c.Error(err)).Msg("Problem with getting ts")
		return
	}

	c.JSON(200, gin.H{
		"message": "pong",
		"request": gin.H{
			"id": requestid.Get(c),
		},
	})
	li.Log.Debug().Msg("Ping-Pong")
}

func IndexHandler(c *gin.Context) {
	li, err := Ctxer(c)
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
