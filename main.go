package main

import (
	"context"
	"errors"
	"github.com/go-playground/webhooks/v6/github"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/ziflex/lecho/v2"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"
)

var DevMode bool
var Version = "DEV"

func init() {
	log.Logger = zerolog.New(os.Stdout).With().Timestamp().Caller().Logger()

}
func main() {

	// Echo instance creation
	e := echo.New()

	// Configure log and debug
	if !DevMode {
		e.Debug = false
		e.HideBanner = true
		e.HidePort = true
		DevMode = false
		e.Use(middleware.Logger())
		e.Logger = lecho.New(log.Logger)
	} else {
		e.Debug = true

		consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC1123Z}
		log.Logger = zerolog.New(consoleWriter).With().Timestamp().Caller().Logger()
		e.Logger = lecho.New(log.Logger)

		e.Use(middleware.LoggerWithConfig(
			middleware.LoggerConfig{
				Format: "[${time_rfc3339}] ${status} ${method} ${path} (${remote_ip}) ${latency_human} [${id}]\n",
				Output: e.Logger.Output(),
			},
		),
		)
	}

	// Middleware
	e.Use(
		middleware.Gzip(),
		middleware.Secure(),
		middleware.Recover(),
		middleware.RequestID(),
	)

	// Control routers
	e.GET("/ready", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	e.GET("/version", func(c echo.Context) error {
		return c.String(http.StatusOK, Version)
	})

	e.POST("/webhook", func(c echo.Context) error {
		hook, _ := github.New(github.Options.Secret("ABC"))
		payload, err := hook.Parse(c.Request(), github.PullRequestEvent)
		if err != nil {
			if err == github.ErrEventNotFound {
				log.Error().Err(err).Msg("Event not found")
			}
		}

		//ctx := context.Background()
		pullRequest := payload.(github.PullRequestPayload)
		err = CloneProject(pullRequest.Repository.Owner.Login, pullRequest.Repository.FullName, pullRequest.Repository.CloneURL, strconv.Itoa(int(pullRequest.Number)))
		if err != nil {
			return err
		}
		//err := gofumpt.Format()
		return err
	})

	go func() {
		if err := e.Start(":8888"); err != nil && !errors.Is(err, http.ErrServerClosed) {
			e.Logger.Fatal("shutting down the server", err)
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server with a timeout of 10 seconds.
	// Use a buffered channel to avoid missing signals as recommended for signal.Notify
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
