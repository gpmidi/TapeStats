package ts

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gpmidi/TapeStats/ts/tsdb"
)

func (ts *TapeStatsApp) RegisterAccountHandler(c *gin.Context) {
	li, err := Ctxer(c)
	if err != nil {
		ts.Log.Error().Err(c.Error(err)).Msg("Problem with getting ts")
		return
	}
	log := li.Log

	act := &tsdb.Account{}
	pw, err := act.CreateSetPassword()
	if err != nil {
		log.Error().Err(c.Error(err)).Msg("Problem with getting random password and/or hash")
		return
	}

	res, err := ts.DB.Model(act).Insert()
	if err != nil {
		log.Error().Err(c.Error(err)).Msg("Problem saving account to db")
		return
	}
	log = log.With().Int("rows.returned", res.RowsReturned()).Int("rows.affected", res.RowsAffected()).Logger()
	if res.RowsAffected() != 1 {
		err = errors.New("one row not affected")
		log.Error().Err(c.Error(err)).Msg("Problem saving account to db - Rows affected is odd")
		return
	}
	if res.RowsReturned() != 1 {
		err = errors.New("one row not returned")
		log.Error().Err(c.Error(err)).Msg("Problem saving account to db - rows returned is odd")
		return
	}

	c.JSON(200, gin.H{
		"message": "Account Created",
		"account": gin.H{
			"id":       act.Id,
			"password": pw,
		},
		"request": li.Data(),
	})
	li.Log.Info().Msg("Account Created")
}
