package ts

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/gpmidi/TapeStats/ts/mam"
	"github.com/rs/zerolog"
)

func (ts *TapeStatsApp) LoadRecordHandler(c *gin.Context) {
	li, err := Ctxer(c)
	if err != nil {
		ts.Log.Error().Err(c.Error(err)).Msg("Problem with getting ts")
		return
	}
	log := li.Log

	log.Error().Msg("Not implemented yet")

	c.JSON(500, gin.H{"error": "Not implemented yet"})
}

func (ts *TapeStatsApp) LoadUnparsedHandler(c *gin.Context) {
	li, err := Ctxer(c)
	if err != nil {
		ts.Log.Error().Err(c.Error(err)).Msg("Problem with getting ts")
		return
	}
	l := li.Log

	// Ready the body
	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(c.Request.Body); err != nil {
		l.Error().Err(c.Error(err)).Msg("Problem reading body")
		return
	}
	// Parse 'em
	fields := mam.NewParser(l).ParseString(buf.String())

	if err := ts.loadFields(l, fields); err != nil {
		l.Error().Err(c.Error(err)).Msg("Problem loading")
		return
	}

	c.JSON(200, gin.H{"message": "ok"})
}

func (ts *TapeStatsApp) loadFields(l zerolog.Logger, fields map[string]*mam.Field) error {
	tx, err := ts.DB.Begin()
	if err != nil {
		l.Warn().Err(err).Msg("Problem starting db transaction")
		return err
	}
	defer func() {
		if err := tx.Close(); err != nil {
			l.Warn().Err(err).Msg("Problem closing db transaction")
		}
	}()

	for name, field := range fields {
		l := l.With().Str("field.name", name).Interface("field", field).Logger()
		l.Info().Msg("Found k:v field")
	}

	return nil
}
