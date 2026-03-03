package main

import (
	"database/sql"
	"log"
	"net/http"

	_ "modernc.org/sqlite"

	"dependency-dashboard/config"
	"dependency-dashboard/internal/depsdev"
	"dependency-dashboard/internal/handler"
	"dependency-dashboard/internal/repository"
	"dependency-dashboard/internal/service"
)

func main() {

	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %v", err)
	}

	db, err := sql.Open("sqlite", cfg.DbPath)
	if err != nil {
		log.Fatal(err)
	}

	repo := repository.New(db)
	if err := repo.InitSchema(); err != nil {
		log.Fatal(err)
	}

	client := depsdev.New(cfg.HttpClientTimeoutSec)
	svc := service.New(repo, client)

	h := handler.New(cfg, svc)

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", h.Routes()))
	// TODO_TOM gracefull shutdown
}
