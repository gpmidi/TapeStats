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
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"strconv"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "TapeStats",
	Short: "Runs tapestats.com!",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.TapeStats.yaml)")

	rootCmd.PersistentFlags().Bool("debug", true, "Enable debugging mode")
	if err := viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug")); err != nil {
		panic(fmt.Sprintf("Error while creating pflag: %v", err))
	}

	rootCmd.PersistentFlags().String("database.url", os.Getenv("DATABASE_URL"), "Postgres Database URL")
	if err := viper.BindPFlag("database.url", rootCmd.PersistentFlags().Lookup("database.url")); err != nil {
		panic(fmt.Sprintf("Error while creating pflag: %v", err))
	}

	var port int64 = 8080
	var err error
	if os.Getenv("PORT") != "" {
		port, err = strconv.ParseInt(os.Getenv("PORT"), 10, 64)
		if err != nil {
			port = 8080 // Don't assume it won't get overwritten
		}
	}
	rootCmd.PersistentFlags().Int("listen.port", (int)(port), "Port to listen on")
	if err := viper.BindPFlag("listen.port", rootCmd.PersistentFlags().Lookup("listen.port")); err != nil {
		panic(fmt.Sprintf("Error while creating pflag: %v", err))
	}
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".TapeStats" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".TapeStats")
	}

	viper.SetEnvPrefix("ts") // Use the prefix "TS_"
	viper.AutomaticEnv()     // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
