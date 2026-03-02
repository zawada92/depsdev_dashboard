package main

import (
	"log"
	"net/http"

	_ "modernc.org/sqlite"

	"dependency-dashboard/config"
	"dependency-dashboard/internal/handler"
	"dependency-dashboard/internal/service"
)

func main() {

	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal("Config error: ", err)
	}

	svc := service.New()
	h := handler.New(cfg, svc)

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", h.Routes()))
}
