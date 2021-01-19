package tsdb

import (
	"context"
	"github.com/go-pg/pg/v10"
)

func Connect(url string) (*pg.DB, error) {
	opt, err := pg.ParseURL(url)
	if err != nil {
		return nil, err
	}

	db := pg.Connect(opt)

	// Validate we can connect
	ctx := context.Background()
	if err := db.Ping(ctx); err != nil {
		panic(err)
	}

	return db, nil
}
