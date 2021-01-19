package migrations

import (
	"github.com/go-pg/migrations/v8"
)

func init() {
	migrations.MustRegisterTx(
		func(db migrations.DB) error {
			_, err := db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`)
			return err
		},
		func(db migrations.DB) error {
			_, err := db.Exec(`
CREATE OR REPLACE FUNCTION trigger_set_timestamp()
RETURNS TRIGGER AS $$
BEGIN
  NEW.modified = NOW();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;
`)
			return err
		},
		func(db migrations.DB) error {
			_, err := db.Exec(`
CREATE TABLE accounts (
	-- Our info
	id NOT NULL PRIMARY KEY DEFAULT uuid_generate_v4(),

	-- Whens (auto set/updated)
	created TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	modified TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		
	-- Auth info
	salt VARCHAR(1024),
	hashed VARCHAR(1024),
);
`)
			return err
		},
		func(db migrations.DB) error {
			_, err := db.Exec(`
CREATE TRIGGER trigger_accounts_set_modified
BEFORE UPDATE ON accounts
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();
`)
			return err
		},
	)
}
