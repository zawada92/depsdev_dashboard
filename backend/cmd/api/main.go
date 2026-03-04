package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	_ "modernc.org/sqlite"

	"dependency-dashboard/config"
	"dependency-dashboard/internal/depsdev"
	"dependency-dashboard/internal/handler"
	"dependency-dashboard/internal/logger"
	"dependency-dashboard/internal/repository"
	"dependency-dashboard/internal/service"

	"github.com/rs/zerolog/log"
)

func main() {

	cfg, err := config.NewConfig()
	if err != nil {
		panic(fmt.Sprintf("Config error: %v", err))
	}

	// TODO move LOG_LEVEL to cfg
	logger.Setup(os.Getenv("LOG_LEVEL"))
	db, err := sql.Open("sqlite", cfg.DbPath)
	if err != nil {
		log.Error().Err(err).Send()
	}

	repo := repository.New(db)
	if err := repo.InitSchema(); err != nil {
		log.Error().Err(err).Send()
	}

	client := depsdev.New(cfg.HttpClientTimeoutSec)
	svc := service.New(repo, client)

	h := handler.New(cfg, svc)

	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second*10, "gracuflly waiting time for server shutdown")
	flag.Parse()

	srv := &http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: h.Routes(),
	}

	go func() {
		log.Info().Msg("Server running on :8080")
		defer func() {
			log.Info().Msg("Server closed")
		}()
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error().Err(err).Send()
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error().Err(err).Send()
	}
	os.Exit(0)
}
