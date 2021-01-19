package main

import (
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
	"github.com/unrolled/secure"
	"os"
	"path/filepath"
	"strings"
)

const ENV_ALLOWED_HOSTS = "ALLOWED_HOSTS"

func createRender(templatesDir string) multitemplate.Renderer {
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
		r.AddFromFiles(filepath.Base(include), files...)
	}
	return r
}

func ParseAllowedHosts() []string {
	ret := make([]string, 0)
	for _, host := range strings.Split(os.Getenv(ENV_ALLOWED_HOSTS), ",") {
		if host != "" {
			ret = append(ret, host)
		}
	}
	return ret
}

func main() {
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

	// Set template context
	t.SetTemplateContext(r)

	// Multi-template render
	// https://github.com/gin-contrib/multitemplate
	r.HTMLRender = createRender("templates")

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
}
