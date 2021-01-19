package ts

import (
	"github.com/go-pg/migrations/v8"
	"github.com/go-pg/pg/v10"
	_ "github.com/gpmidi/TapeStats/ts/migrations" // Load any migrations files
	"github.com/gpmidi/TapeStats/ts/tsdb"
)

func (ts *TapeStatsApp) MigrationsRun(args ...string) error {
	l := ts.Log.With().Strs("args", args).Logger()

	if len(args) == 0 {
		l.Info().Msg("No args given - Running init+up")
		if err := ts.MigrationsRun("init"); err != nil {
			return err
		}
		if err := ts.MigrationsRun("up"); err != nil {
			return err
		}
		l.Info().Msg("No args given - Done running init+up")
		return nil
	}

	l.Info().Msg("Starting Migration")
	if err := migrations.DefaultCollection.DiscoverSQLMigrations("migrations"); err != nil {
		l.Error().Err(err).Msg("Failed to read/discover SQL migrations from FS")
		return err
	}

	oldVersion, newVersion, err := migrations.Run(ts.DB, args...)
	l = l.With().Int64("version.old", oldVersion).Int64("version.new", newVersion).Logger()
	l.Info().Msg("Ending Migration")
	if err != nil {
		l.Error().Err(err).Msg("Failed Migration")
		return err
	}

	l.Info().Msg("Migration successful")
	return nil
}

//tapeExists returns if the tape is already in the tapes table or not
func (ts *TapeStatsApp) tapeExists(tx *pg.Tx, accountId string, manufacturer string, manufactureDT string,
	serialNumber string, densityCode string, mediumType string, ltoVersion int) (bool, error) {
	tape := new(tsdb.Tape)
	err := tx.Model(tape).Where("id = ?", accountId).Where("manufacturer = ?", manufacturer).
		Where("manufacture_dt = ?", manufactureDT).Where("serial_number = ?", serialNumber).
		Where("density_code = ?", densityCode).Where("medium_type = ?", mediumType).
		Where("lto_version = ?", ltoVersion).Select()
	if err != nil {
		return false, err
	}
	return tape.Id != "", nil
}
