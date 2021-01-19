/*
Copyright © 2021 Paulson McIntyre <paul@tapestats.com>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/
package cmd

import (
	"github.com/gin-gonic/gin"
	"github.com/gpmidi/TapeStats/ts"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

// migrateCmd represents the migrate command
var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Run Postgres and/or Redis migrations",
	Run: func(cmd *cobra.Command, args []string) {
		zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
		zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

		if strings.ToLower(os.Getenv("DEBUG")) == "true" {
			gin.SetMode(gin.DebugMode)
		}
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
		if gin.IsDebugging() {
			zerolog.SetGlobalLevel(zerolog.DebugLevel)
		}

		log.Logger = log.Output(
			zerolog.ConsoleWriter{
				Out:     os.Stderr,
				NoColor: true,
			},
		).With().
			Str("program", "tape-stats").
			Logger()

		// Local logger
		l := log.With().Logger()

		// Our core
		t, err := ts.NewTapeStatsApp(l)
		if err != nil {
			log.Panic().Err(err).Msg("Problem creating base server setup")
		}

		if err := t.MigrationsRun(args); err != nil {
			log.Panic().Err(err).Msg("Problem running migrations")
		}
	},
}

func init() {
	rootCmd.AddCommand(migrateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// migrateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// migrateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
