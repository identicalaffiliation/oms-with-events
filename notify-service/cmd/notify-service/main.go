package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/identicalaffiliation/oms-with-events/notify-service/internal/infrastructure/broker"
	"github.com/identicalaffiliation/oms-with-events/notify-service/internal/infrastructure/config"
	"github.com/identicalaffiliation/oms-with-events/notify-service/internal/infrastructure/database"
	"github.com/identicalaffiliation/oms-with-events/notify-service/internal/infrastructure/logger"
	"github.com/identicalaffiliation/oms-with-events/notify-service/internal/repository"
	"github.com/identicalaffiliation/oms-with-events/notify-service/pkg/proccesser"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "c", "config.yml", "path to config file")
	flag.Parse()

	cfg, err := config.NewNotifyServiceConfig(configPath)
	if err != nil {
		fmt.Fprintf(os.Stdout, "error: %v\n", err)
		os.Exit(1)
	}

	slogger, err := logger.NewSLogger(cfg)
	if err != nil {
		fmt.Fprintf(os.Stdout, "error: %v\n", err)
		os.Exit(1)
	}

	pool, err := database.NewPool(cfg)
	if err != nil {
		slogger.Error("failed to open new pool", "error", err)
		os.Exit(1)
	}

	defer func() {
		if err := pool.Close(); err != nil {
			slogger.Error("failed to close pool", "error", err)
		}
	}()

	eventsRepository := repository.NewProcEventsRepository(pool)
	consumer := broker.NewKafkaConsumer(cfg)
	defer consumer.Close()
	proccesser := proccesser.NewProccesser(eventsRepository, consumer, slogger)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		<-signalChan
		slogger.Debug("gracefully stopped..")
		cancel()
	}()

	proccesser.Run(ctx)
}
