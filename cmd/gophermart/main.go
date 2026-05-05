package main

import (
	"context"
	"net/http"
	"os"

	"github.com/lxmp7p/ya-dipl1/internal/repository"

	"github.com/lxmp7p/ya-dipl1/internal/config"
	"github.com/lxmp7p/ya-dipl1/internal/handler"
	"github.com/lxmp7p/ya-dipl1/internal/utils.go"
)

func main() {
	logger := utils.CreateLogger()
	cfg := config.NewConfig()
	cfg.InitConfig()
	repository, err := repository.InitDB(cfg, logger)
	if err != nil {
		logger.Error("failed init db")
		os.Exit(1)
	}

	worker := utils.NewAccrualWorker(cfg.BalanceSystemAddr, &repository, logger)
	worker.Start(context.Background())
	defer worker.Stop()

	handler := handler.NewHandler(logger, repository, cfg)
	r := handler.InitRoutes()

	logger.Info("Starting server", "addr", cfg.Addr)
	err = http.ListenAndServe(cfg.Addr, r)
	logger.Info("Server stopped", "err", err)
}
