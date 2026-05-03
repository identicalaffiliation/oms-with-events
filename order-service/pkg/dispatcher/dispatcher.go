package dispatcher

import (
	"context"
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
}

func NewDispatcher(
	producer Producer,
	eventsRepo EventsRepository,
	workerCount int,
	batchSize int,
	retryCount int,
	logger logger.Logger,
	chillDuration time.Duration,
) *dispatcher {
	return &dispatcher{
		producer:         producer,
		eventsRepository: eventsRepo,
		workersCount:     workerCount,
		batchSize:        batchSize,
		retryCount:       retryCount,
		logger:           logger,
		chillDuration:    chillDuration,
	}
}

func (d *dispatcher) Run(ctx context.Context) {
	jobs := make(chan *domain.OrderEvent)

	go func() {
		<-ctx.Done()
		close(jobs)
	}()

	for i := 0; i < d.workersCount; i++ {
		go d.worker(ctx, jobs)
	}

	for {
		events, err := d.eventsRepository.GetUnsentEvents(ctx, d.batchSize)
		if err != nil {
			d.logger.Error("failed to get events", "error", err)
			time.Sleep(d.chillDuration)
			continue
		}

		if len(events) == 0 {
			time.Sleep(d.chillDuration)
			continue
		}

		for _, e := range events {
			select {
			case jobs <- e:
			case <-ctx.Done():
				return
			}
		}
	}
}

func (d *dispatcher) worker(ctx context.Context, jobs <-chan *domain.OrderEvent) {
	for e := range jobs {
		for attempt := 0; attempt < d.retryCount; attempt++ {
			select {
			case <-ctx.Done():
				return
			default:
			}

			err := d.producer.Produce(ctx, e.Payload, e.ID.String(), e.EventType)
			if err == nil {
				if err := d.eventsRepository.MarkEventAsSent(ctx, e.ID); err != nil {
					d.logger.Error("failed to mark event as sent", "event id", e.ID, "error", err)
				}
				break
			}

			d.logger.Error("failed to send event to broker", "error", err)

			select {
			case <-ctx.Done():
				return
			case <-time.After(d.chillDuration):
			}
		}

	}
}
