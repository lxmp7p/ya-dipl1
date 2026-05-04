package main

import (
	"net/http"

	"github.com/lxmp7p/ya-dipl1/internal/repository"

	"github.com/lxmp7p/ya-dipl1/internal/config"
	"github.com/lxmp7p/ya-dipl1/internal/handler"
	"github.com/lxmp7p/ya-dipl1/internal/utils.go"
)

func main() {
	logger := utils.CreateLogger()
	cfg := config.NewConfig()
	cfg.InitConfig()
	database, err := repository.InitDB(cfg, logger)
	if err != nil {
		panic("failed init db")
	}

	handler := handler.NewHandler(logger, database, cfg)
	r := handler.InitRoutes()

	logger.Info("Starting server", "addr", cfg.Addr)
	err = http.ListenAndServe(cfg.Addr, r)
	logger.Info("Server stopped", "err", err)
}
