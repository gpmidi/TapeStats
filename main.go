package main

import (
	"github.com/gin-contrib/logger"
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
	"os"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if gin.IsDebugging() {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	log.Logger = log.Output(
		zerolog.ConsoleWriter{
			Out:     os.Stderr,
			NoColor: true,
		},
	).With().
		Str("program", "tape-stats").
		Logger()

	// Local logger
	l := log.With().Logger()

	r := gin.New()

	// Logging middleware
	r.Use(logger.SetLogger())

	// Request ID middleware
	r.Use(requestid.New(requestid.Config{
		Generator: func() string {
			return uuid.New().String()
		},
	}))

	// Log request w/ Request id and save logger
	r.Use(RequestIDLogMiddleware())

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	l = l.With().Str("listen.port", port).Logger()
	l.Info().Msg("Going to listen on port")

	r.GET("/ping", func(c *gin.Context) {
		li, err := Ctxer(c)
		if err != nil {
			c.JSON(500, gin.H{
				"error": "Unexpected Server Error",
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
	})
	if err := r.Run("0.0.0.0:" + port); err != nil {
		log.Panic().Err(err).Msg("Problem running server")
	}
}
