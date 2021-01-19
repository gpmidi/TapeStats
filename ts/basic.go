package ts

import (
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
)

func PingHandler(c *gin.Context) {
	li, err := Ctxer(c)
	if err != nil {
		c.JSON(500, gin.H{
			"error": "Unexpected Server Error",
			"code":  err.Error(),
			"request": gin.H{
				"id": requestid.Get(c),
			},
		})
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
