package main

import (
	"fmt"
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const REQUESTOR_MIDDLEWARE_NAME = "Requestor-Middleware"

type RequestorInstance struct {
	Log zerolog.Logger
}

func RequestIDLogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		l := log.With().Str("request.id", requestid.Get(c)).Logger()
		l.Debug().Msg("New Request")
		ri := RequestorInstance{
			Log: l,
		}
		c.Set(REQUESTOR_MIDDLEWARE_NAME, ri)
	}
}

func Ctxer(c *gin.Context) (*RequestorInstance, error) {
	lraw, ok := c.Get(REQUESTOR_MIDDLEWARE_NAME)
	if !ok {
		return nil, fmt.Errorf("no such key in gin.Context %v", REQUESTOR_MIDDLEWARE_NAME)
	}

	l, ok := lraw.(RequestorInstance)
	if !ok {
		return nil, fmt.Errorf("data in key of gin.Context %v is wrong type (%v)", REQUESTOR_MIDDLEWARE_NAME, lraw)
	}

	return &l, nil
}
