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
	"github.com/gin-contrib/logger"
	"github.com/gin-contrib/multitemplate"
	"github.com/gin-contrib/requestid"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gpmidi/TapeStats/ts"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
	"github.com/spf13/viper"
	"github.com/unrolled/secure"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

func createRender(templatesDir string, t *ts.TapeStatsApp) multitemplate.Renderer {
	r := multitemplate.NewRenderer()

	layouts, err := filepath.Glob(templatesDir + "/layouts/*.html")
	if err != nil {
		panic(err.Error())
	}

	includes, err := filepath.Glob(templatesDir + "/includes/*.html")
	if err != nil {
		panic(err.Error())
	}

	// Generate our templates map from our layouts/ and includes/ directories
	for _, include := range includes {
		layoutCopy := make([]string, len(layouts))
		copy(layoutCopy, layouts)
		files := append(layoutCopy, include)
		r.AddFromFilesFuncs(filepath.Base(include), t.GetTemplateContext(), files...)
	}
	return r
}

func ParseAllowedHosts() []string {
	ret := make([]string, 0)
	for _, host := range viper.GetStringSlice("hosts.allowed") {
		if host != "" {
			// TODO: Do more FQDN validation
			ret = append(ret, host)
		}
	}
	return ret
}

// webCmd represents the web command
var webCmd = &cobra.Command{
	Use:   "web",
	Short: "Runs main web server",
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

		t := ts.NewTapeStatsApp()

		// Local logger
		l := log.With().Logger()

		r := gin.New()

		// Multi-template render
		// https://github.com/gin-contrib/multitemplate
		r.HTMLRender = createRender("templates", t)

		// Logging middleware
		r.Use(logger.SetLogger())

		// Request ID middleware
		r.Use(requestid.New(requestid.Config{
			Generator: func() string {
				return uuid.New().String()
			},
		}))

		// Security Middleware (mainly http->https)
		secureMiddleware := secure.New(secure.Options{
			FrameDeny:    true,
			SSLRedirect:  true,
			AllowedHosts: ParseAllowedHosts(),
			SSLProxyHeaders: map[string]string{
				"X-Forwarded-Proto": "https",
			},
		})

		secureFunc := func() gin.HandlerFunc {
			return func(c *gin.Context) {
				err := secureMiddleware.Process(c.Writer, c.Request)

				// If there was an error, do not continue.
				if err != nil {
					c.Abort()
					return
				}

				// Avoid header rewrite if response is a redirection.
				if status := c.Writer.Status(); status > 300 && status < 399 {
					c.Abort()
				}
			}
		}()
		r.Use(secureFunc)

		// Host static files
		r.Use(static.Serve("/static", static.LocalFile("static", false)))

		// Log request w/ Request id and save logger
		r.Use(ts.RequestIDLogMiddleware(t))

		// FIXME: Use viper for PORT
		port := os.Getenv("PORT")
		if port == "" {
			port = "8080"
		}
		l = l.With().Str("listen.port", port).Logger()
		l.Info().Msg("Going to listen on port")

		t.AddRoutes(r)

		if err := r.Run("0.0.0.0:" + port); err != nil {
			log.Panic().Err(err).Msg("Problem running server")
		}

	},
}

func init() {
	rootCmd.AddCommand(webCmd)

	webCmd.PersistentFlags().Int32("port", 8080, "Listen port")
	if err := viper.BindPFlag("port", webCmd.PersistentFlags().Lookup("port")); err != nil {
		panic(fmt.Sprintf("Error while creating pflag: %v", err))
	}
}
