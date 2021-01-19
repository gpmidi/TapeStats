package ts

import (
	"context"
	"github.com/go-pg/pg/v10"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

type TapeStatsApp struct {
	DB  *pg.DB
	Log zerolog.Logger
}

func NewTapeStatsApp(log zerolog.Logger) (*TapeStatsApp, error) {
	// Get a good connection
	opt, err := pg.ParseURL(viper.GetString("database.url"))
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse url")
		return nil, err
	}

	db := pg.Connect(opt)

	// Validate we can connect
	ctx := context.Background()
	if err := db.Ping(ctx); err != nil {
		log.Error().Err(err).Msg("Failed to ping database")
		return nil, err
	}

	// TODO: Add context info to log
	return &TapeStatsApp{
		DB:  db,
		Log: log,
	}, nil
}
