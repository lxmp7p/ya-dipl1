package gophermart

import (
	"net/http"

	"github.com/lxmp7p/ya-dipl1/internal/config"
	"github.com/lxmp7p/ya-dipl1/internal/handler"
	"github.com/lxmp7p/ya-dipl1/internal/utils.go"
)

func main() {
	logger := utils.CreateLogger()
	cfg := config.NewConfig()
	cfg.InitConfig()

	handler := handler.NewHandler(logger)
	r := handler.InitRoutes()

	logger.Printf("Starting server on %s", cfg.Addr)
	err := http.ListenAndServe(cfg.Addr, r)
	logger.Fatalf("Server stopped: %v", err)
}
