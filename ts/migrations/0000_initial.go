package migrations

import (
	"github.com/go-pg/migrations/v8"
)

func init() {
	migrations.MustRegisterTx(func(db migrations.DB) error {
		_, err := db.Exec(`CREATE EXTENSION hstore;`)
		return err
	}, func(db migrations.DB) error {
		_, err := db.Exec(`DROP EXTENSION hstore;`)
		return err
	})
}
