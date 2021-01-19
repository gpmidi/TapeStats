package migrations

import (
	"github.com/go-pg/migrations/v8"
)

func init() {
	migrations.MustRegisterTx(func(db migrations.DB) error {
		_, err := db.Exec(`CREATE TABLE test_table()`)
		return err
	}, func(db migrations.DB) error {
		_, err := db.Exec(`DROP TABLE test_table`)
		return err
	})
}
