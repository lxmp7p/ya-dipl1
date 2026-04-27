package gophermart

import (
	"github.com/lxmp7p/ya-dipl1/internal/handler"
	"github.com/lxmp7p/ya-dipl1/internal/utils.go"
)

func main() {
	logger := utils.CreateLogger()

	handler := handler.NewHandler(logger)
	handler.InitRoutes()
}
