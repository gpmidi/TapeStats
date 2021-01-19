package ts

import (
	"github.com/go-pg/pg/v10"
	"github.com/gpmidi/TapeStats/ts/tsdb"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

type TapeStatsApp struct {
	DB  *pg.DB
	Log zerolog.Logger
}

func NewTapeStatsApp(log zerolog.Logger) (*TapeStatsApp, error) {
	// Get a good connection
	db, err := tsdb.Connect(viper.GetString("database.url"))
	if err != nil {
		return nil, err
	}
	// TODO: Add context info to log
	return &TapeStatsApp{
		DB:  db,
		Log: log,
	}, nil
}
