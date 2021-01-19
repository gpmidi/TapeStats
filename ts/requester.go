package ts

import (
	"fmt"
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"time"
)

const REQUESTOR_MIDDLEWARE_NAME = "Requestor-Middleware"

type RequestorInstance struct {
	Log       zerolog.Logger
	TS        *TapeStatsApp
	RequestId string
}

func (ri *RequestorInstance) OurData() gin.H {
	ret := gin.H{
		"id":  ri.RequestId,
		"now": time.Now().UTC().String(),
	}

	return ret
}

func (ri *RequestorInstance) Data(overrides ...gin.H) gin.H {
	ret := ri.OurData()

	// Apply overrides - Do last
	for _, override := range overrides {
		for k, v := range override {
			ret[k] = v
		}
	}

	return ret
}

func RequestIDLogMiddleware(app *TapeStatsApp) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Before
		ri := app.GetSetRI(c)

		t := time.Now() // Runs just before request

		c.Next()

		// After
		ri.Log.Debug().TimeDiff("latency", time.Now(), t).Msg("Request inner runtime") // Right after request
		// More after....

	}
}

func (ts *TapeStatsApp) Ctxer(c *gin.Context) (*RequestorInstance, error) {
	lraw, ok := c.Get(REQUESTOR_MIDDLEWARE_NAME)
	if !ok {
		return nil, fmt.Errorf("no such key in gin.Context %v", REQUESTOR_MIDDLEWARE_NAME)
	}

	l, ok := lraw.(*RequestorInstance)
	if !ok {
		return nil, fmt.Errorf("data in key of gin.Context %v is wrong type (%v)", REQUESTOR_MIDDLEWARE_NAME, lraw)
	}

	// Update the log in case it went bad
	l.Log = ts.GetLogGinCtx(c)

	return l, nil
}

func (ts *TapeStatsApp) GetLogGinCtx(c *gin.Context) zerolog.Logger {
	return log.With().Str("request.id", requestid.Get(c)).Logger()
}

func (ts *TapeStatsApp) GetRI(c *gin.Context) *RequestorInstance {
	l := ts.GetLogGinCtx(c)
	l.Debug().Msg("New Request")
	ri := RequestorInstance{
		Log:       l,
		TS:        ts,
		RequestId: requestid.Get(c),
	}
	return &ri
}

func (ts *TapeStatsApp) GetSetRI(c *gin.Context) *RequestorInstance {
	ri := ts.GetRI(c)
	c.Set(REQUESTOR_MIDDLEWARE_NAME, ri)
	return ri
}
