package dispatcher

import (
	"context"
	"database/sql"
	"time"

	"github.com/identicalaffiliation/oms-with-events/order-service/internal/infrastructure/logger"
	"github.com/identicalaffiliation/oms-with-events/order-service/internal/models/domain"
)

type dispatcher struct {
	producer         Producer
	eventsRepository EventsRepository
	workersCount     int
	batchSize        int
	retryCount       int
	logger           logger.Logger
	chillDuration    time.Duration
	pool             *sql.DB
}

func NewDispatcher(
	producer Producer,
	eventsRepo EventsRepository,
	workerCount int,
	batchSize int,
	retryCount int,
	logger logger.Logger,
	chillDuration time.Duration,
	pool *sql.DB,
) *dispatcher {
	return &dispatcher{
		producer:         producer,
		eventsRepository: eventsRepo,
		workersCount:     workerCount,
		batchSize:        batchSize,
		retryCount:       retryCount,
		logger:           logger,
		chillDuration:    chillDuration,
		pool:             pool,
	}
}

func (d *dispatcher) Run(ctx context.Context) {
	jobs := make(chan *domain.OrderEvent, d.batchSize*2)

	for i := 0; i < d.workersCount; i++ {
		go d.worker(ctx, jobs)
	}

	go func() {
		ticker := time.NewTicker(d.chillDuration)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				close(jobs)
				return

			case <-ticker.C:
				d.fetchAndDispatch(ctx, jobs)
			}
		}
	}()
}

func (d *dispatcher) worker(ctx context.Context, jobs <-chan *domain.OrderEvent) {
	for e := range jobs {

		var success bool
		for attempt := 0; attempt < d.retryCount; attempt++ {

			err := d.producer.Produce(ctx, e.Payload, e.ID.String())
			if err == nil {
				success = true
				break
			}

			d.logger.Error("failed to send event", "error", err)

			select {
			case <-ctx.Done():
				return
			case <-time.After(d.chillDuration):
			}
		}

		if success {
			if err := d.eventsRepository.MarkEventAsSent(ctx, e.ID); err != nil {
				d.logger.Error("failed to mark event as sent", "error", err, "id", e.ID)
			} else {
				d.logger.Debug("event sent", "id", e.ID)
			}
		}
	}
}

func (d *dispatcher) fetchAndDispatch(ctx context.Context, jobs chan<- *domain.OrderEvent) {
	tx, err := d.pool.BeginTx(ctx, nil)
	if err != nil {
		d.logger.Error("begin tx failed", "error", err)
		return
	}
	defer tx.Rollback()

	events, err := d.eventsRepository.GetUnsentEventsWithTx(ctx, tx, d.batchSize)
	if err != nil {
		d.logger.Error("claim events failed", "error", err)
		return
	}

	if len(events) == 0 {
		return
	}

	if err := tx.Commit(); err != nil {
		d.logger.Error("commit claim tx failed", "error", err)
		return
	}

	for _, e := range events {
		select {
		case jobs <- e:
		case <-ctx.Done():
			return
		}
	}
}
