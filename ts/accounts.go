package ts

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg/v10"
	"github.com/gpmidi/TapeStats/ts/tsdb"
)

func (ts *TapeStatsApp) RegisterOrgHandler(c *gin.Context) {
	li, err := ts.Ctxer(c)
	if err != nil {
		ts.Log.Error().Err(c.Error(err)).Msg("Problem with getting ts")
		return
	}
	log := li.Log

	log.Error().Msg("Not implemented yet")

	c.JSON(500, gin.H{"error": "Not implemented yet"})
}

func (ts *TapeStatsApp) RegisterUserHandler(c *gin.Context) {
	// FIXME: Add auth
	li, err := ts.Ctxer(c)
	if err != nil {
		ts.Log.Error().Err(c.Error(err)).Msg("Problem with getting ts")
		return
	}
	log := li.Log

	log.Error().Msg("Not implemented yet")

	c.JSON(500, gin.H{"error": "Not implemented yet"})
}

func (ts *TapeStatsApp) RegisterAccountHandler(c *gin.Context) {
	// FIXME: Add auth
	li, err := ts.Ctxer(c)
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
			"guid":     act.Guid,
			"password": pw,
		},
		"request": li.Data(),
	})
	li.Log.Info().Msg("Account Created")
}

// getAccount fetches the request account with the given PK using the given TX
func (ts *TapeStatsApp) getAccount(tx *pg.Tx, accountGuid string) (*tsdb.Account, error) {
	account := new(tsdb.Account)
	err := tx.Model(account).Where("guid = ?", accountGuid).Select()
	if err != nil {
		return nil, err
	}
	return account, nil
}
