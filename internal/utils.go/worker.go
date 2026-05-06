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
	"golang.org/x/sync/errgroup"
)

type AccrualWorker struct {
	client     *http.Client
	accrualURL string
	repo       *repository.Repository
	logger     *slog.Logger
	interval   time.Duration
	wg         sync.WaitGroup
	workers    int
	cancel     context.CancelFunc
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
		workers:    WorkersCount,
	}
}

func (w *AccrualWorker) Start(ctx context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	w.cancel = cancel

	w.wg.Add(1)
	go w.worker(ctx)
	w.logger.Info("Accrual worker started", "interval", w.interval, "workers", w.workers)
}

func (w *AccrualWorker) Stop() {
	w.cancel()
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
				[]string{NewOrderStatus, ProcessingOrderStatus},
			)
			if err != nil {
				w.logger.Error("ListAllOrders failed", "error", err)
				continue
			}
			if len(orders) == 0 {
				w.logger.Debug("Not found orders by filter")
				continue
			}

			g, gCtx := errgroup.WithContext(ctx)
			g.SetLimit(w.workers)

			for _, order := range orders {
				order := order
				g.Go(func() error {
					select {
					case <-gCtx.Done():
						return gCtx.Err()
					default:
						return w.processOrder(gCtx, order)
					}
				})
			}
			if err := g.Wait(); err != nil {
				w.logger.Error("processing batch failed", "error", err)
			}
		case <-ctx.Done():
			w.logger.Info("worker stopped by ctx")
			return
		}
	}
}

func (w *AccrualWorker) processOrder(ctx context.Context, order repository.OrderData) error {
	var response accrualResponse
	url := fmt.Sprintf("%s/api/orders/%s", w.accrualURL, order.OrderNumber)
	resp, err := w.client.Get(url)
	if err != nil {
		w.logger.Error("request failed", "url", url, "error", err)
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		w.logger.Error("failed to read body", "error", err)
		return err
	}

	if err := json.Unmarshal(body, &response); err != nil {
		w.logger.Error("failed to parse JSON", "error", err, "body", string(body))
		return err
	}

	w.logger.Info("got accrual response",
		"order", response.Order,
		"status", response.Status,
		"accrual", response.Accrual)

	if response.Status == "PROCESSED" {
		if err = w.repo.Update(ctx, response.Order, response.Status, response.Accrual); err != nil {
			w.logger.Error("failed to update order", "error", err)
			return err
		}

		_, err = w.repo.AddBalance(ctx, order.UserID, response.Accrual)
		if err != nil {
			w.logger.Error("failed to add balance", "error", err)
			return err
		}
	}
	return nil
}
