package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/lxmp7p/ya-dipl1/internal/repository"
)

type AccrualWorker struct {
	client     *http.Client
	accrualURL string
	repo       *repository.Repository
	logger     *slog.Logger
	interval   time.Duration
	stopCh     chan struct{}
	wg         sync.WaitGroup
}

type accrualResponse struct {
	Order   string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float64 `json:"accrual"`
}

func NewAccrualWorker(
	accrualURL string,
	repo *repository.Repository,
	logger *slog.Logger,
) *AccrualWorker {
	return &AccrualWorker{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		accrualURL: accrualURL,
		repo:       repo,
		logger:     logger,
		interval:   5 * time.Second,
		stopCh:     make(chan struct{}),
	}
}

func (w *AccrualWorker) Start(ctx context.Context) {
	w.wg.Add(1)
	go w.worker(ctx)
	w.logger.Info("Accrual worker started", "interval", w.interval)
}

func (w *AccrualWorker) Stop() {
	close(w.stopCh)
	w.wg.Wait()
	w.logger.Info("Accrual worker stopped")
}

func (w *AccrualWorker) worker(ctx context.Context) {
	defer w.wg.Done()

	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			orders, err := w.repo.ListAllOrdersByStatuses(
				ctx,
				[]string{NewOrderStatus, ProcessingOrderStatus, InvalidOrderStatus},
			)
			if err != nil {
				w.logger.Error("ListAllOrders failed", "error", err)
			}
			for _, order := range orders {
				var response accrualResponse
				url := fmt.Sprintf("%s/api/orders/%s", w.accrualURL, order.OrderNumber)
				resp, err := w.client.Get(url)
				if err != nil {
					w.logger.Error("request failed", "url", url, "error", err)
					continue
				}
				body, err := io.ReadAll(resp.Body)
				resp.Body.Close()
				if err != nil {
					w.logger.Error("failed to read body", "error", err)
					continue
				}

				if err := json.Unmarshal(body, &response); err != nil {
					w.logger.Error("failed to parse JSON", "error", err, "body", string(body))
					continue
				}

				w.logger.Info("got accrual response",
					"order", response.Order,
					"status", response.Status,
					"accrual", response.Accrual)

				if response.Status == "PROCESSED" {
					w.repo.Update(ctx, response.Order, response.Status, response.Accrual)
					w.repo.AddBalance(ctx, order.UserID, response.Accrual)
				}
			}
		case <-w.stopCh:
			return
		}
	}
}
