package main

import (
	"log"
	"net/http"
	adapters "taxApi/internal/adapter"
	"taxApi/internal/configs"
	"taxApi/internal/handlers"
	"taxApi/internal/logs"
	"taxApi/internal/repositories"
	"taxApi/internal/services"
	"taxApi/pkg/database"
)

func main() {
	config, err := configs.InitConfigs()
	if err != nil {
		log.Fatal(err)
	}

	address := config.Server.Host + config.Server.Port
	connToDb := database.InitConnectionToDb(config)
	initLog := logs.NewLogger
	adapter := adapters.NewAdapter(config, initLog)
	repository := repositories.NewRepository(connToDb, initLog)
	service := services.NewService(adapter, repository, initLog)
	handler := handlers.NewHandler(service, initLog)
	router := NewRouter(handler)
	srv := http.Server{
		Addr:    address,
		Handler: router,
	}
	err = srv.ListenAndServe()
	if err != nil {
		log.Println("listen and serve error")
	}
}
