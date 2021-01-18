package main

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"os"
)

func main() {
	log := logrus.New().WithField("program", "tape-stats")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log = log.WithField("listen.port", port)
	log.Info("Going to listen on port")

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	if err := r.Run("0.0.0.0:" + port); err != nil {
		log.WithError(err).Panic("Problem running server")
	}
}
